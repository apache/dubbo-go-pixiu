/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package listener

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/filterchain"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const snapshotTrackingFilterName = "dgp.test.listener-snapshot-tracking"

var snapshotFilters sync.Map

type snapshotTrackingConfig struct {
	ID string `yaml:"id" json:"id" mapstructure:"id"`
}

type snapshotTrackingPlugin struct{}

func (*snapshotTrackingPlugin) Kind() string { return snapshotTrackingFilterName }
func (*snapshotTrackingPlugin) Config() any  { return &snapshotTrackingConfig{} }
func (*snapshotTrackingPlugin) CreateFilter(config any) (filter.NetworkFilter, error) {
	created := &snapshotTrackingFilter{}
	snapshotFilters.Store(config.(*snapshotTrackingConfig).ID, created)
	return created, nil
}

type snapshotTrackingFilter struct {
	filter.EmptyNetworkFilter
	closed atomic.Int32
}

func (f *snapshotTrackingFilter) Close() error {
	f.closed.Add(1)
	return nil
}

func init() {
	filter.RegisterNetworkFilterPlugin(&snapshotTrackingPlugin{})
}

func buildSnapshotTrackingChain(t *testing.T, id string) (*filterchain.NetworkFilterChain, *snapshotTrackingFilter) {
	t.Helper()
	chain, err := filterchain.BuildNetworkFilterChain(model.FilterChain{Filters: []model.NetworkFilter{{
		Name:   snapshotTrackingFilterName,
		Config: map[string]any{"id": id},
	}}})
	require.NoError(t, err)
	value, ok := snapshotFilters.Load(id)
	require.True(t, ok)
	return chain, value.(*snapshotTrackingFilter)
}

func TestLongLivedRequestDoesNotBlockUpdateAndDefersRetiredChainClose(t *testing.T) {
	oldChain, oldFilter := buildSnapshotTrackingChain(t, "long-lived-old")
	newChain, newFilter := buildSnapshotTrackingChain(t, "long-lived-new")
	base := NewBaseListenerService(&model.Listener{}, oldChain)

	entered := make(chan struct{})
	release := make(chan struct{})
	requestDone := make(chan error, 1)
	go func() {
		requestDone <- base.WithFilterChain(func(chain *filterchain.NetworkFilterChain) error {
			close(entered)
			if chain != oldChain {
				return fmt.Errorf("long-lived request acquired the wrong filter chain")
			}
			<-release
			return nil
		})
	}()
	<-entered

	committed := make(chan *RetiredUpdate, 1)
	go func() {
		gate := base.gate()
		gate.Lock()
		retired := base.CommitRefresh(&PreparedUpdate{snapshot: newFilterChainSnapshot(newChain)})
		gate.Unlock()
		committed <- retired
	}()

	var retired *RetiredUpdate
	select {
	case retired = <-committed:
	case <-time.After(time.Second):
		t.Fatal("listener update was blocked by a long-lived request")
	}
	require.NoError(t, retired.Close())
	require.Zero(t, oldFilter.closed.Load(), "retired chain closed while a request still referenced it")
	require.NoError(t, base.WithFilterChain(func(chain *filterchain.NetworkFilterChain) error {
		require.Same(t, newChain, chain)
		return nil
	}))

	close(release)
	require.NoError(t, <-requestDone)
	require.Eventually(t, func() bool { return oldFilter.closed.Load() == 1 }, time.Second, time.Millisecond)
	require.NoError(t, base.CloseFilterChain())
	require.Equal(t, int32(1), newFilter.closed.Load())
}

func TestSharedUpdateGatePublishesListenerChainsTogether(t *testing.T) {
	gate := &sync.RWMutex{}
	oldOne := &filterchain.NetworkFilterChain{}
	oldTwo := &filterchain.NetworkFilterChain{}
	newOne := &filterchain.NetworkFilterChain{}
	newTwo := &filterchain.NetworkFilterChain{}
	one := NewBaseListenerService(&model.Listener{}, oldOne)
	two := NewBaseListenerService(&model.Listener{}, oldTwo)
	one.setUpdateGate(gate)
	two.setUpdateGate(gate)

	gate.Lock()
	requestStarted := make(chan struct{})
	requestDone := make(chan struct{})
	type observation struct {
		first  *filterchain.NetworkFilterChain
		second *filterchain.NetworkFilterChain
		err    error
	}
	observed := make(chan observation, 1)
	go func() {
		close(requestStarted)
		var result observation
		result.err = one.WithFilterChain(func(first *filterchain.NetworkFilterChain) error {
			result.first = first
			return two.WithFilterChain(func(second *filterchain.NetworkFilterChain) error {
				result.second = second
				return nil
			})
		})
		observed <- result
		close(requestDone)
	}()
	<-requestStarted

	retiredOne := one.CommitRefresh(&PreparedUpdate{snapshot: newFilterChainSnapshot(newOne)})
	retiredTwo := two.CommitRefresh(&PreparedUpdate{snapshot: newFilterChainSnapshot(newTwo)})
	require.Same(t, oldOne, retiredOne.snapshot.filterChain)
	require.Same(t, oldTwo, retiredTwo.snapshot.filterChain)
	select {
	case <-requestDone:
		t.Fatal("request entered while the multi-listener commit gate was held")
	case <-time.After(20 * time.Millisecond):
	}

	gate.Unlock()
	select {
	case <-requestDone:
		result := <-observed
		require.NoError(t, result.err)
		require.Same(t, newOne, result.first)
		require.Same(t, newTwo, result.second)
	case <-time.After(time.Second):
		t.Fatal("request did not resume after the complete transaction was published")
	}
	require.NoError(t, retiredOne.Close())
	require.NoError(t, retiredTwo.Close())
}

func TestInactiveListenerRejectsRequests(t *testing.T) {
	base := NewBaseListenerService(&model.Listener{}, &filterchain.NetworkFilterChain{})
	base.SetActive(false)

	require.ErrorContains(t, base.WithFilterChain(func(*filterchain.NetworkFilterChain) error {
		t.Fatal("inactive listener entered its filter chain")
		return nil
	}), "not active")
}

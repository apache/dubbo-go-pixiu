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

package logic

import (
	"context"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/require"

	clientv3 "go.etcd.io/etcd/client/v3"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
)

type blockingEtcdRangeGetter struct {
	started chan struct{}
}

func (g *blockingEtcdRangeGetter) Get(ctx context.Context, _ string, _ ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	close(g.started)
	<-ctx.Done()
	return nil, ctx.Err()
}

func TestGetChildrenKVListStopsWhenRequestContextIsCanceled(t *testing.T) {
	getter := &blockingEtcdRangeGetter{started: make(chan struct{})}
	ctx, cancel := context.WithCancel(context.Background())
	done := make(chan error, 1)
	go func() {
		_, _, err := getChildrenKVList(ctx, getter, "/pixiu/config/listeners")
		done <- err
	}()

	select {
	case <-getter.started:
	case <-time.After(time.Second):
		t.Fatal("etcd range request did not start")
	}
	cancel()

	select {
	case err := <-done:
		require.ErrorIs(t, err, context.Canceled)
	case <-time.After(time.Second):
		t.Fatal("etcd range request did not stop after cancellation")
	}
}

func TestDecodeXDSResourcesPropagatesMalformedYAML(t *testing.T) {
	previous := adminconfig.Bootstrap
	adminconfig.Bootstrap = &adminconfig.AdminBootstrap{EtcdConfig: adminconfig.EtcdConfig{Path: "/pixiu/config/api"}}
	t.Cleanup(func() { adminconfig.Bootstrap = previous })

	_, err := decodeClusters(
		[]string{getClusterKey("1")},
		[]string{"name: ["},
	)
	require.ErrorContains(t, err, "decode cluster configuration")

	_, err = decodeListeners(
		[]string{getListenerKey("http")},
		[]string{"name: ["},
	)
	require.ErrorContains(t, err, "decode listener configuration")
}

func TestDecodeXDSResourcesKeepsLegitimateEmptySets(t *testing.T) {
	clusters, err := decodeClusters(nil, nil)
	require.NoError(t, err)
	require.Empty(t, clusters)

	listeners, err := decodeListeners(nil, nil)
	require.NoError(t, err)
	require.Empty(t, listeners)
}

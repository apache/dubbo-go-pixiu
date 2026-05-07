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

package loadbalancer

import (
	"sync/atomic"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type legacyLoadBalancer struct {
	seenEndpoints []*model.Endpoint
}

type legacyCursorLoadBalancer struct{}

var _ LoadBalancer = (*legacyLoadBalancer)(nil)

type blockingLegacyLoadBalancer struct {
	entered chan int
	release chan struct{}
	calls   int32
}

func (l *legacyLoadBalancer) Handler(c *model.ClusterConfig, _ model.LbPolicy) *model.Endpoint {
	l.seenEndpoints = c.Endpoints
	if len(c.Endpoints) == 0 {
		return nil
	}
	return c.Endpoints[0]
}

func (legacyCursorLoadBalancer) Handler(c *model.ClusterConfig, _ model.LbPolicy) *model.Endpoint {
	index := atomic.AddUint32(&c.PrePickEndpointIndex, 1) - 1
	return c.Endpoints[int(index%uint32(len(c.Endpoints)))]
}

func (b *blockingLegacyLoadBalancer) Handler(c *model.ClusterConfig, _ model.LbPolicy) *model.Endpoint {
	call := atomic.AddInt32(&b.calls, 1)
	b.entered <- int(call)
	<-b.release
	if len(c.Endpoints) == 0 {
		return nil
	}
	return c.Endpoints[0]
}

func TestPickEndpointAdaptsLegacyLoadBalancer(t *testing.T) {
	healthy := &model.Endpoint{ID: "healthy"}
	unhealthy := &model.Endpoint{ID: "unhealthy"}
	cluster := &model.ClusterConfig{
		Name:      "legacy-load-balancer",
		Endpoints: []*model.Endpoint{healthy, unhealthy},
	}
	balancer := &legacyLoadBalancer{}

	got := PickEndpoint(balancer, PickContext{
		Config:           cluster,
		HealthyEndpoints: []*model.Endpoint{healthy},
	}, nil)

	assert.Same(t, healthy, got)
	assert.Equal(t, []*model.Endpoint{healthy}, balancer.seenEndpoints)
	assert.Equal(t, []*model.Endpoint{healthy, unhealthy}, cluster.Endpoints)
}

func TestPickEndpointReconcilesLegacyCursorState(t *testing.T) {
	first := &model.Endpoint{ID: "first"}
	second := &model.Endpoint{ID: "second"}
	cluster := &model.ClusterConfig{
		Name:                 "legacy-cursor-load-balancer",
		Endpoints:            []*model.Endpoint{first, second},
		PrePickEndpointIndex: 3,
	}

	got := PickEndpoint(legacyCursorLoadBalancer{}, PickContext{
		Config:           cluster,
		HealthyEndpoints: []*model.Endpoint{first, second},
	}, nil)

	assert.Same(t, second, got)
	assert.Equal(t, uint32(4), atomic.LoadUint32(&cluster.PrePickEndpointIndex))
}

func TestPickEndpointSerializesLegacyLoadBalancerHandlers(t *testing.T) {
	first := &model.Endpoint{ID: "first"}
	cluster := &model.ClusterConfig{
		Name:      "blocking-legacy-load-balancer",
		Endpoints: []*model.Endpoint{first},
	}
	balancer := &blockingLegacyLoadBalancer{
		entered: make(chan int, 2),
		release: make(chan struct{}),
	}
	pickContext := PickContext{
		Config:           cluster,
		HealthyEndpoints: []*model.Endpoint{first},
	}

	firstDone := make(chan struct{})
	go func() {
		defer close(firstDone)
		_ = PickEndpoint(balancer, pickContext, nil)
	}()
	assert.Equal(t, 1, waitLegacyHandlerEntry(t, balancer.entered))

	secondDone := make(chan struct{})
	go func() {
		defer close(secondDone)
		_ = PickEndpoint(balancer, pickContext, nil)
	}()

	select {
	case call := <-balancer.entered:
		t.Fatalf("legacy handler call %d entered before the first call returned", call)
	case <-time.After(50 * time.Millisecond):
	}

	close(balancer.release)
	assert.Equal(t, 2, waitLegacyHandlerEntry(t, balancer.entered))
	waitClosed(t, firstDone)
	waitClosed(t, secondDone)
}

func waitLegacyHandlerEntry(t *testing.T, entered <-chan int) int {
	t.Helper()
	select {
	case call := <-entered:
		return call
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for legacy handler entry")
		return 0
	}
}

func waitClosed(t *testing.T, done <-chan struct{}) {
	t.Helper()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for legacy pick to finish")
	}
}

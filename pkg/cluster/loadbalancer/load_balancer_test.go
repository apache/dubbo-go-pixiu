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
	"sync"
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

type clusterScopedBlockingLegacyLoadBalancer struct {
	*blockingLegacyLoadBalancer
}

type observableLocker struct {
	mu      sync.Mutex
	waiter  chan struct{}
	locked  bool
	blocked chan struct{}
}

type legacyPickHarness struct {
	t           *testing.T
	balancer    LoadBalancer
	blocking    *blockingLegacyLoadBalancer
	releaseOnce sync.Once
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

func (b *clusterScopedBlockingLegacyLoadBalancer) UseClusterScopedLegacyLock() bool {
	return true
}

func newObservableLocker() *observableLocker {
	return &observableLocker{
		blocked: make(chan struct{}, 1),
	}
}

func (l *observableLocker) Lock() {
	l.mu.Lock()
	if !l.locked {
		l.locked = true
		l.mu.Unlock()
		return
	}
	waiter := make(chan struct{})
	l.waiter = waiter
	l.blocked <- struct{}{}
	l.mu.Unlock()
	<-waiter

	l.mu.Lock()
	l.locked = true
	l.mu.Unlock()
}

func (l *observableLocker) Unlock() {
	l.mu.Lock()
	l.locked = false
	if l.waiter != nil {
		close(l.waiter)
		l.waiter = nil
	}
	l.mu.Unlock()
}

func newLegacyPickHarness(t *testing.T, scoped bool) *legacyPickHarness {
	t.Helper()
	blocking := &blockingLegacyLoadBalancer{
		entered: make(chan int, 2),
		release: make(chan struct{}),
	}
	harness := &legacyPickHarness{
		t:        t,
		balancer: blocking,
		blocking: blocking,
	}
	if scoped {
		harness.balancer = &clusterScopedBlockingLegacyLoadBalancer{
			blockingLegacyLoadBalancer: blocking,
		}
	}
	t.Cleanup(harness.release)
	return harness
}

func newLegacyPickContext(clusterName, endpointID string) PickContext {
	endpoint := &model.Endpoint{ID: endpointID}
	return PickContext{
		Config: &model.ClusterConfig{
			Name:      clusterName,
			Endpoints: []*model.Endpoint{endpoint},
		},
		HealthyEndpoints: []*model.Endpoint{endpoint},
	}
}

func (h *legacyPickHarness) startDirectPick(context PickContext) <-chan struct{} {
	h.t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = PickEndpoint(h.balancer, context, nil)
	}()
	return done
}

func (h *legacyPickHarness) startLockedPick(lock sync.Locker, context PickContext) <-chan struct{} {
	h.t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = PickEndpointWithLegacyLock(h.balancer, lock, context, nil)
	}()
	return done
}

func (h *legacyPickHarness) release() {
	h.releaseOnce.Do(func() {
		close(h.blocking.release)
	})
}

func (h *legacyPickHarness) waitEntry() int {
	h.t.Helper()
	return waitLegacyHandlerEntry(h.t, h.blocking.entered)
}

func (h *legacyPickHarness) assertNoEntry() {
	h.t.Helper()
	assertNoLegacyHandlerEntry(h.t, h.blocking.entered)
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
	harness := newLegacyPickHarness(t, false)
	pickContext := newLegacyPickContext("blocking-legacy-load-balancer", "first")

	firstDone := harness.startDirectPick(pickContext)
	assert.Equal(t, 1, harness.waitEntry())

	secondDone := harness.startDirectPick(pickContext)
	harness.assertNoEntry()

	harness.release()
	assert.Equal(t, 2, harness.waitEntry())
	waitClosed(t, firstDone)
	waitClosed(t, secondDone)
}

func TestPickEndpointKeepsCompatibilityLockForOptInLegacyLoadBalancer(t *testing.T) {
	harness := newLegacyPickHarness(t, true)
	firstContext := newLegacyPickContext("blocking-legacy-load-balancer-first-direct-pick", "first")
	secondContext := newLegacyPickContext("blocking-legacy-load-balancer-second-direct-pick", "second")

	firstDone := harness.startDirectPick(firstContext)
	assert.Equal(t, 1, harness.waitEntry())

	secondDone := harness.startDirectPick(secondContext)
	harness.assertNoEntry()

	harness.release()
	assert.Equal(t, 2, harness.waitEntry())
	waitClosed(t, firstDone)
	waitClosed(t, secondDone)
}

func TestPickEndpointWithLegacyLockSerializesOptInLegacyLoadBalancerHandlersWithSameLock(t *testing.T) {
	harness := newLegacyPickHarness(t, true)
	legacyPickLock := newObservableLocker()
	pickContext := newLegacyPickContext("blocking-legacy-load-balancer-same-lock", "first")

	firstDone := harness.startLockedPick(legacyPickLock, pickContext)
	assert.Equal(t, 1, harness.waitEntry())

	secondDone := harness.startLockedPick(legacyPickLock, pickContext)
	waitLegacyLockBlocked(t, legacyPickLock.blocked)

	harness.release()
	assert.Equal(t, 2, harness.waitEntry())
	waitClosed(t, firstDone)
	waitClosed(t, secondDone)
}

func TestPickEndpointWithLegacyLockSerializesNonOptInLegacyLoadBalancerHandlersAcrossDifferentLocks(t *testing.T) {
	harness := newLegacyPickHarness(t, false)
	firstContext := newLegacyPickContext("blocking-legacy-load-balancer-first-lock", "first")
	secondContext := newLegacyPickContext("blocking-legacy-load-balancer-second-lock", "second")
	var firstLock sync.Mutex
	var secondLock sync.Mutex

	firstDone := harness.startLockedPick(&firstLock, firstContext)
	assert.Equal(t, 1, harness.waitEntry())

	secondDone := harness.startLockedPick(&secondLock, secondContext)
	harness.assertNoEntry()

	harness.release()
	assert.Equal(t, 2, harness.waitEntry())
	waitClosed(t, firstDone)
	waitClosed(t, secondDone)
}

func TestPickEndpointWithLegacyLockAllowsOptInLegacyLoadBalancerHandlersWithDifferentLocksToRunConcurrently(t *testing.T) {
	harness := newLegacyPickHarness(t, true)
	firstContext := newLegacyPickContext("blocking-legacy-load-balancer-first-lock", "first")
	secondContext := newLegacyPickContext("blocking-legacy-load-balancer-second-lock", "second")
	var firstLock sync.Mutex
	var secondLock sync.Mutex

	firstDone := harness.startLockedPick(&firstLock, firstContext)
	assert.Equal(t, 1, harness.waitEntry())

	secondDone := harness.startLockedPick(&secondLock, secondContext)
	assert.Equal(t, 2, harness.waitEntry())

	harness.release()
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

func assertNoLegacyHandlerEntry(t *testing.T, entered <-chan int) {
	t.Helper()
	select {
	case call := <-entered:
		t.Fatalf("legacy handler call %d entered before the first call returned", call)
	case <-time.After(50 * time.Millisecond):
	}
}

func waitLegacyLockBlocked(t *testing.T, blocked <-chan struct{}) {
	t.Helper()
	select {
	case <-blocked:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for legacy pick to block on the runtime lock")
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

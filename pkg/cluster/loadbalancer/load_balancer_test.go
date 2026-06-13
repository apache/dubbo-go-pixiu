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
	"fmt"
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
type mutatingLegacyLoadBalancer struct{}
type unhealthyLegacyLoadBalancer struct{}
type mutatingSnapshotLoadBalancer struct{}
type unhealthySnapshotLoadBalancer struct{}
type healthyOnlySnapshotLoadBalancer struct{}

var _ LoadBalancer = (*legacyLoadBalancer)(nil)

type blockingLegacyLoadBalancer struct {
	entered chan int
	release chan struct{}
	calls   int32
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

func (mutatingLegacyLoadBalancer) Handler(c *model.ClusterConfig, _ model.LbPolicy) *model.Endpoint {
	if len(c.Endpoints) == 0 {
		return nil
	}
	c.Endpoints[0].Metadata["weight"] = "99"
	c.Endpoints[0].LLMMeta.APIKey = "mutated-key"
	return c.Endpoints[0]
}

func (unhealthyLegacyLoadBalancer) Handler(c *model.ClusterConfig, _ model.LbPolicy) *model.Endpoint {
	if len(c.Endpoints) < 2 {
		return nil
	}
	return c.Endpoints[1]
}

func (mutatingSnapshotLoadBalancer) Handler(_ *model.ClusterConfig, _ model.LbPolicy) *model.Endpoint {
	return nil
}

func (mutatingSnapshotLoadBalancer) HandlerWithSnapshot(c PickContext, _ model.LbPolicy) *model.Endpoint {
	if len(c.HealthyEndpoints) == 0 {
		return nil
	}
	c.HealthyEndpoints[0].Metadata["weight"] = "99"
	c.HealthyEndpoints[0].LLMMeta.APIKey = "mutated-key"
	if len(c.AllEndpoints) > 1 {
		c.AllEndpoints[1].Metadata["weight"] = "42"
	}
	return c.HealthyEndpoints[0]
}

func (unhealthySnapshotLoadBalancer) Handler(_ *model.ClusterConfig, _ model.LbPolicy) *model.Endpoint {
	return nil
}

func (unhealthySnapshotLoadBalancer) HandlerWithSnapshot(c PickContext, _ model.LbPolicy) *model.Endpoint {
	if len(c.AllEndpoints) < 2 {
		return nil
	}
	return c.AllEndpoints[1]
}

func (healthyOnlySnapshotLoadBalancer) Handler(_ *model.ClusterConfig, _ model.LbPolicy) *model.Endpoint {
	return nil
}

func (healthyOnlySnapshotLoadBalancer) HandlerWithSnapshot(_ PickContext, _ model.LbPolicy) *model.Endpoint {
	return nil
}

func (healthyOnlySnapshotLoadBalancer) UseHealthyEndpointsOnly() bool {
	return true
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

func newLegacyPickHarness(t *testing.T) *legacyPickHarness {
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

func (h *legacyPickHarness) startPick(context PickContext) <-chan struct{} {
	h.t.Helper()
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = PickEndpoint(h.balancer, context, nil)
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
	unhealthy := &model.Endpoint{ID: "unhealthy", UnHealthy: true}
	cluster := &model.ClusterConfig{
		Name:      "legacy-load-balancer",
		Endpoints: []*model.Endpoint{healthy, unhealthy},
	}
	balancer := &legacyLoadBalancer{}

	got := PickEndpoint(balancer, PickContext{
		AllEndpoints:     []*model.Endpoint{healthy, unhealthy},
		Config:           cluster,
		HealthyEndpoints: []*model.Endpoint{healthy},
	}, nil)

	assert.NotSame(t, healthy, got)
	assert.Equal(t, healthy, got)
	assert.Equal(t, []*model.Endpoint{healthy, unhealthy}, balancer.seenEndpoints)
	assert.NotSame(t, healthy, balancer.seenEndpoints[0])
	assert.NotSame(t, unhealthy, balancer.seenEndpoints[1])
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

	assert.NotSame(t, second, got)
	assert.Equal(t, second, got)
	assert.Equal(t, uint32(4), atomic.LoadUint32(&cluster.PrePickEndpointIndex))
}

func TestPickEndpointLegacyMutationsDoNotEscape(t *testing.T) {
	healthy := &model.Endpoint{
		ID:       "healthy",
		Metadata: map[string]string{"weight": "1"},
		LLMMeta:  &model.LLMMeta{APIKey: "original-key"},
	}
	cluster := &model.ClusterConfig{
		Name:      "legacy-mutation",
		Endpoints: []*model.Endpoint{healthy},
	}

	got := PickEndpoint(mutatingLegacyLoadBalancer{}, PickContext{
		AllEndpoints:     []*model.Endpoint{healthy},
		Config:           cluster,
		HealthyEndpoints: []*model.Endpoint{healthy},
	}, nil)

	if assert.NotNil(t, got) {
		assert.Equal(t, map[string]string{"weight": "1"}, got.Metadata)
		assert.Equal(t, "original-key", got.LLMMeta.APIKey)
	}
	assert.Equal(t, map[string]string{"weight": "1"}, healthy.Metadata)
	assert.Equal(t, "original-key", healthy.LLMMeta.APIKey)
}

func TestPickEndpointRejectsUnhealthyLegacyReturn(t *testing.T) {
	healthy := &model.Endpoint{ID: "healthy"}
	unhealthy := &model.Endpoint{ID: "unhealthy", UnHealthy: true}
	cluster := &model.ClusterConfig{
		Name:      "legacy-health-guard",
		Endpoints: []*model.Endpoint{healthy, unhealthy},
	}

	got := PickEndpoint(unhealthyLegacyLoadBalancer{}, PickContext{
		AllEndpoints:     []*model.Endpoint{healthy, unhealthy},
		Config:           cluster,
		HealthyEndpoints: []*model.Endpoint{healthy},
	}, nil)

	assert.Nil(t, got)
}

func TestPickEndpointSnapshotMutationsDoNotEscape(t *testing.T) {
	healthy := &model.Endpoint{
		ID:       "healthy",
		Metadata: map[string]string{"weight": "1"},
		LLMMeta:  &model.LLMMeta{APIKey: "original-key"},
	}
	unhealthy := &model.Endpoint{
		ID:        "unhealthy",
		Metadata:  map[string]string{"weight": "2"},
		UnHealthy: true,
	}
	cluster := &model.ClusterConfig{
		Name:      "snapshot-mutation",
		Endpoints: []*model.Endpoint{healthy, unhealthy},
	}

	got := PickEndpoint(mutatingSnapshotLoadBalancer{}, PickContext{
		AllEndpoints:     []*model.Endpoint{healthy, unhealthy},
		Config:           cluster,
		HealthyEndpoints: []*model.Endpoint{healthy},
	}, nil)

	if assert.NotNil(t, got) {
		assert.NotSame(t, healthy, got)
		assert.Equal(t, map[string]string{"weight": "1"}, got.Metadata)
		assert.Equal(t, "original-key", got.LLMMeta.APIKey)
	}
	assert.Equal(t, map[string]string{"weight": "1"}, healthy.Metadata)
	assert.Equal(t, "original-key", healthy.LLMMeta.APIKey)
	assert.Equal(t, map[string]string{"weight": "2"}, unhealthy.Metadata)
}

func TestPickEndpointRejectsUnhealthySnapshotReturn(t *testing.T) {
	healthy := &model.Endpoint{ID: "healthy"}
	unhealthy := &model.Endpoint{ID: "unhealthy", UnHealthy: true}
	cluster := &model.ClusterConfig{
		Name:      "snapshot-health-guard",
		Endpoints: []*model.Endpoint{healthy, unhealthy},
	}

	got := PickEndpoint(unhealthySnapshotLoadBalancer{}, PickContext{
		AllEndpoints:     []*model.Endpoint{healthy, unhealthy},
		Config:           cluster,
		HealthyEndpoints: []*model.Endpoint{healthy},
	}, nil)

	assert.Nil(t, got)
}

func TestNeedsAllEndpointsKeepsCompatibilityForUnmarkedSnapshotLoadBalancer(t *testing.T) {
	assert.True(t, NeedsAllEndpoints(mutatingSnapshotLoadBalancer{}))
	assert.False(t, NeedsAllEndpoints(healthyOnlySnapshotLoadBalancer{}))
}

func TestPickEndpointSerializesLegacyLoadBalancerHandlers(t *testing.T) {
	harness := newLegacyPickHarness(t)
	pickContext := newLegacyPickContext("blocking-legacy-load-balancer", "first")

	firstDone := harness.startPick(pickContext)
	assert.Equal(t, 1, harness.waitEntry())

	secondDone := harness.startPick(pickContext)
	harness.assertNoEntry()

	harness.release()
	assert.Equal(t, 2, harness.waitEntry())
	waitClosed(t, firstDone)
	waitClosed(t, secondDone)
}

func TestPickEndpointSerializesLegacyHandlersThroughPackageLock(t *testing.T) {
	harness := newLegacyPickHarness(t)
	firstContext := newLegacyPickContext("blocking-legacy-load-balancer-first-context", "first")
	secondContext := newLegacyPickContext("blocking-legacy-load-balancer-second-context", "second")

	firstDone := harness.startPick(firstContext)
	assert.Equal(t, 1, harness.waitEntry())

	secondDone := harness.startPick(secondContext)
	harness.assertNoEntry()

	harness.release()
	assert.Equal(t, 2, harness.waitEntry())
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

func waitClosed(t *testing.T, done <-chan struct{}) {
	t.Helper()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for legacy pick to finish")
	}
}

// TestSameEndpointIdentityBlankDomainWildcard locks the narrow placeholder
// wildcard in sameEndpointIdentity. The wildcard fires when the SNAPSHOT
// endpoint (first argument, per call-site convention in
// healthyEndpointFromSnapshot) has a fully-empty SocketAddress — a single
// blank-domain marker with no IP/port. Any non-zero address field on the
// snapshot side revokes the wildcard so a stray blank domain in a config
// does not silently widen the trust boundary. See the contract comment on
// sameEndpointIdentity for the trust model.
//
// Field names below match the function signature exactly (snapshot is
// `candidate`, balancer return is `endpoint`) so a reader cannot lose
// track of which side carries the placeholder.
func TestSameEndpointIdentityBlankDomainWildcard(t *testing.T) {
	cases := []struct {
		name           string
		snapshot       *model.Endpoint // function param: candidate
		balancerReturn *model.Endpoint // function param: endpoint
		want           bool
		why            string
	}{
		{
			name: "snapshot_placeholder_accepts_resolved_balancer_address",
			snapshot: &model.Endpoint{
				ID:      "shared-id",
				Address: model.SocketAddress{Domains: []string{""}},
			},
			balancerReturn: &model.Endpoint{
				ID:      "shared-id",
				Address: model.SocketAddress{Address: "127.0.0.1", Port: 8080},
			},
			want: true,
			why:  "snapshot side is the fully-empty placeholder; legacy resolver wildcard fires on ID alone",
		},
		{
			name: "snapshot_real_domain_requires_address_equality_even_with_matching_id",
			snapshot: &model.Endpoint{
				ID:      "shared-id",
				Address: model.SocketAddress{Domains: []string{"openai.com"}},
			},
			balancerReturn: &model.Endpoint{
				ID:      "shared-id",
				Address: model.SocketAddress{Address: "127.0.0.1", Port: 8080},
			},
			want: false,
			why:  "snapshot side is a real domain, not a placeholder; address must match",
		},
		{
			name: "different_id_with_snapshot_placeholder_is_not_a_match",
			snapshot: &model.Endpoint{
				ID:      "snapshot-id",
				Address: model.SocketAddress{Domains: []string{""}},
			},
			balancerReturn: &model.Endpoint{
				ID:      "different-id",
				Address: model.SocketAddress{Domains: []string{""}},
			},
			want: false,
			why:  "ID is the trust boundary; placeholder address does not paper over an ID mismatch",
		},
		{
			name: "snapshot_blank_domain_with_real_port_revokes_wildcard",
			snapshot: &model.Endpoint{
				ID: "shared-id",
				Address: model.SocketAddress{
					Domains: []string{""},
					Address: "10.0.0.1",
					Port:    8080,
				},
			},
			balancerReturn: &model.Endpoint{
				ID:      "shared-id",
				Address: model.SocketAddress{Address: "10.0.0.2", Port: 8080},
			},
			want: false,
			why:  "snapshot address fields are populated, so this is not a placeholder; address must match",
		},
		{
			name: "snapshot_with_two_domains_including_blank_revokes_wildcard",
			snapshot: &model.Endpoint{
				ID: "shared-id",
				Address: model.SocketAddress{
					Domains: []string{"", "fallback.example.com"},
				},
			},
			balancerReturn: &model.Endpoint{
				ID:      "shared-id",
				Address: model.SocketAddress{Address: "127.0.0.1", Port: 8080},
			},
			want: false,
			why:  "snapshot Domains has more than one entry, so this is not the placeholder form",
		},
		{
			name: "balancer_return_placeholder_does_NOT_trigger_wildcard_when_snapshot_is_real",
			snapshot: &model.Endpoint{
				ID:      "shared-id",
				Address: model.SocketAddress{Address: "127.0.0.1", Port: 8080},
			},
			balancerReturn: &model.Endpoint{
				ID:      "shared-id",
				Address: model.SocketAddress{Domains: []string{""}},
			},
			want: false,
			why:  "wildcard is keyed on the SNAPSHOT side, not the balancer return; reversing roles is a no-op",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			// Call-site convention from healthyEndpointFromSnapshot:
			//   sameEndpointIdentity(candidate /* snapshot */, endpoint /* return */)
			got := sameEndpointIdentity(tc.snapshot, tc.balancerReturn)
			assert.Equal(t, tc.want, got, tc.why)
		})
	}
}

// TestHealthyEndpointFromSnapshotAcceptsResolvedAddressForBlankPlaceholder
// is the integration test for the wildcard. It exercises the real call
// path (healthyEndpointFromSnapshot iterating the healthy snapshot slice)
// so the parameter-order contract between the caller and
// sameEndpointIdentity is verified end-to-end. A previous regression
// passed the unit test by reversing the args; this test would have
// caught it because there is no way to mis-name the arguments when they
// come from a real snapshot.
func TestHealthyEndpointFromSnapshotAcceptsResolvedAddressForBlankPlaceholder(t *testing.T) {
	// Snapshot holds a legacy resolver placeholder: ID set, address absent.
	snapshotPlaceholder := &model.Endpoint{
		ID:      "legacy-resolver-id",
		Address: model.SocketAddress{Domains: []string{""}},
	}
	healthyEndpoints := []*model.Endpoint{snapshotPlaceholder}

	// Balancer (or downstream resolver) returns the resolved real address
	// for the same ID. The healthy lookup MUST accept this and return a
	// clone of the snapshot entry.
	balancerReturn := &model.Endpoint{
		ID:      "legacy-resolver-id",
		Address: model.SocketAddress{Address: "127.0.0.1", Port: 8080},
	}

	got := healthyEndpointFromSnapshot(balancerReturn, healthyEndpoints)
	if !assert.NotNil(t, got, "balancer return with same ID as snapshot placeholder must match via wildcard") {
		return
	}
	assert.Equal(t, snapshotPlaceholder.ID, got.ID)
	assert.NotSame(t, snapshotPlaceholder, got, "returned endpoint must be a clone, not the snapshot pointer")
}

func TestHealthyEndpointFromSnapshotPointerFastPathReturnsClone(t *testing.T) {
	endpoint := &model.Endpoint{
		ID: "zero-copy-id",
		Address: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    8080,
		},
	}

	got := healthyEndpointFromSnapshot(endpoint, []*model.Endpoint{endpoint})

	if !assert.NotNil(t, got) {
		return
	}
	assert.Equal(t, endpoint, got)
	assert.NotSame(t, endpoint, got, "request path must not return the snapshot-owned endpoint pointer")
}

// TestHealthyEndpointFromSnapshotRejectsMismatchedRealAddress ensures the
// wildcard is not a free pass: when the snapshot has a real address and
// the balancer returns a different real address for the same ID, the
// lookup must reject so the pick path returns nil instead of a stale
// hostname/port.
func TestHealthyEndpointFromSnapshotRejectsMismatchedRealAddress(t *testing.T) {
	snapshot := &model.Endpoint{
		ID:      "shared-id",
		Address: model.SocketAddress{Address: "127.0.0.1", Port: 8080},
	}
	healthyEndpoints := []*model.Endpoint{snapshot}

	balancerReturn := &model.Endpoint{
		ID:      "shared-id",
		Address: model.SocketAddress{Address: "10.0.0.1", Port: 9090},
	}

	got := healthyEndpointFromSnapshot(balancerReturn, healthyEndpoints)
	assert.Nil(t, got, "real-address mismatch must not match even when IDs agree")
}

func BenchmarkHealthyEndpointFromSnapshot(b *testing.B) {
	const endpointCount = 1024
	healthyEndpoints := make([]*model.Endpoint, endpointCount)
	for i := range healthyEndpoints {
		healthyEndpoints[i] = &model.Endpoint{
			ID: fmt.Sprintf("ep-%d", i),
			Address: model.SocketAddress{
				Address: "127.0.0.1",
				Port:    10000 + i,
			},
		}
	}
	zeroCopyEndpoint := healthyEndpoints[endpointCount-1]
	defensiveCopyEndpoint := model.CloneEndpoint(zeroCopyEndpoint)

	b.Run("pointer-eq-fast-path", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if healthyEndpointFromSnapshot(zeroCopyEndpoint, healthyEndpoints) == nil {
				b.Fatal("expected match")
			}
		}
	})

	b.Run("identity-scan", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			if healthyEndpointFromSnapshot(defensiveCopyEndpoint, healthyEndpoints) == nil {
				b.Fatal("expected match")
			}
		}
	})
}

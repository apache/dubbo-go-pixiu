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

package cluster

import (
	"fmt"
	"sync"
	"sync/atomic"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/healthcheck"
	_ "github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/maglev"   // Register Maglev for snapshot hash tests.
	_ "github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/ringhash" // Register RingHash for snapshot hash tests.
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestClusterEndpointSnapshotSeedsFromEndpointHealth(t *testing.T) {
	healthy := testEndpoint("ep-1", "127.0.0.1", 18080)
	unhealthy := testEndpoint("ep-2", "127.0.0.1", 18081)
	unhealthy.UnHealthy = true

	runtimeCluster := NewCluster(testCluster("snapshot-seed", healthy, unhealthy))
	snapshot := runtimeCluster.EndpointSnapshot()

	assert.Equal(t, []*model.Endpoint{healthy, unhealthy}, snapshot.AllEndpoints())
	assert.Equal(t, []*model.Endpoint{healthy}, snapshot.HealthyEndpoints())
	assert.NotSame(t, healthy, snapshot.EndpointByID(healthy.ID))
	assert.Equal(t, healthy, snapshot.EndpointByID(healthy.ID))
	assert.NotSame(t, healthy, snapshot.HealthyEndpointByID(healthy.ID))
	assert.Equal(t, healthy, snapshot.HealthyEndpointByID(healthy.ID))
	assert.Nil(t, snapshot.HealthyEndpointByID(unhealthy.ID))
}

func TestClusterEndpointSnapshotReturnsDefensiveEndpointSlices(t *testing.T) {
	first := testEndpoint("ep-1", "127.0.0.1", 18080)
	second := testEndpoint("ep-2", "127.0.0.1", 18081)

	runtimeCluster := NewCluster(testCluster("snapshot-defensive-copy", first, second))
	snapshot := runtimeCluster.EndpointSnapshot()

	all := snapshot.AllEndpoints()
	healthy := snapshot.HealthyEndpoints()
	all[0] = nil
	healthy[0] = nil

	assert.Equal(t, []*model.Endpoint{first, second}, snapshot.AllEndpoints())
	assert.Equal(t, []*model.Endpoint{first, second}, snapshot.HealthyEndpoints())
	assert.Equal(t, 2, snapshot.EndpointCount())
	assert.NotSame(t, first, snapshot.EndpointByID(first.ID))
	assert.Equal(t, first, snapshot.EndpointByID(first.ID))
	endpointByID := snapshot.EndpointByID(first.ID)
	healthyByID := snapshot.HealthyEndpointByID(first.ID)
	assert.NotSame(t, endpointByID, healthyByID)
	assert.Equal(t, endpointByID, healthyByID)
}

func TestClusterEndpointSnapshotEndpointCountIsNilSafe(t *testing.T) {
	var snapshot *EndpointSnapshot

	assert.Zero(t, snapshot.EndpointCount())
}

// TestClusterEndpointSnapshotForPickAccessorsExposeStableView locks the
// per-call semantic contract of HealthyEndpointsForPick / AllEndpointsForPick,
// without locking the implementation detail of returning the same backing
// slice. Concretely: repeated calls return slices of the same length, same
// IDs in the same order, and the same endpoint *pointers* per slot (so
// snapshot consumers can compare by identity). This deliberately stops
// short of asserting that callers see the snapshot's internal slice header
// itself, leaving room for future lazy-build / metric / hook indirections.
func TestClusterEndpointSnapshotForPickAccessorsExposeStableView(t *testing.T) {
	first := testEndpoint("ep-1", "127.0.0.1", 18080)
	second := testEndpoint("ep-2", "127.0.0.1", 18081)

	runtimeCluster := NewCluster(testCluster("snapshot-for-pick-alias", first, second))
	snapshot := runtimeCluster.EndpointSnapshot()

	healthy1 := snapshot.HealthyEndpointsForPick()
	healthy2 := snapshot.HealthyEndpointsForPick()
	if assert.Len(t, healthy1, 2) && assert.Len(t, healthy2, 2) {
		assert.Equal(t, healthy1[0].ID, healthy2[0].ID)
		assert.Equal(t, healthy1[1].ID, healthy2[1].ID)
		// Same endpoint identity per slot: snapshot guarantees the for-pick
		// accessors are stable for the snapshot's lifetime. We compare the
		// endpoint pointer (zero-copy view of the snapshot's owned endpoint)
		// rather than the slice header address.
		assert.Same(t, healthy1[0], healthy2[0])
		assert.Same(t, healthy1[1], healthy2[1])
	}

	all1 := snapshot.AllEndpointsForPick()
	all2 := snapshot.AllEndpointsForPick()
	if assert.Len(t, all1, 2) && assert.Len(t, all2, 2) {
		assert.Equal(t, all1[0].ID, all2[0].ID)
		assert.Equal(t, all1[1].ID, all2[1].ID)
		assert.Same(t, all1[0], all2[0])
		assert.Same(t, all1[1], all2[1])
	}

	var nilSnapshot *EndpointSnapshot
	assert.Nil(t, nilSnapshot.HealthyEndpointsForPick())
	assert.Nil(t, nilSnapshot.AllEndpointsForPick())
}

func TestClusterEndpointSnapshotBuildsConsistentHashFromRuntimeHealthyEndpoints(t *testing.T) {
	tests := []model.LbPolicyType{
		model.LoadBalancerRingHashing,
		model.LoadBalancerMaglevHashing,
	}

	for _, lb := range tests {
		t.Run(string(lb), func(t *testing.T) {
			first := testEndpoint("ep-1", "127.0.0.1", 18080)
			second := testEndpoint("ep-2", "127.0.0.1", 18081)
			temporarilyUnhealthy := testEndpoint("ep-3", "127.0.0.1", 18082)
			config := testCluster("snapshot-healthy-hash", first, second, temporarilyUnhealthy)
			config.LbStr = lb
			config.ConsistentHash = model.ConsistentHash{
				ReplicaNum:      10,
				MaxVnodeNum:     1023,
				MaglevTableSize: 521,
			}
			runtimeCluster := NewCluster(config)

			assert.True(t, runtimeCluster.UpdateEndpointHealth(
				temporarilyUnhealthy.ID,
				temporarilyUnhealthy.Address.GetAddress(),
				false,
			))

			snapshot := runtimeCluster.EndpointSnapshot()
			hash := snapshot.HealthyConsistentHash()
			if !assert.NotNil(t, hash) {
				return
			}
			for i := 0; i < 20; i++ {
				host, err := hash.Get(fmt.Sprintf("request-%d", i))
				if assert.NoError(t, err) {
					assert.NotEqual(t, temporarilyUnhealthy.GetHost(), host)
				}
			}
		})
	}
}

func TestClusterEndpointSnapshotExposesReadOnlyConsistentHash(t *testing.T) {
	first := testEndpoint("ep-1", "127.0.0.1", 18080)
	second := testEndpoint("ep-2", "127.0.0.1", 18081)
	config := testCluster("snapshot-readonly-hash", first, second)
	config.LbStr = model.LoadBalancerRingHashing
	config.ConsistentHash = model.ConsistentHash{
		ReplicaNum:  10,
		MaxVnodeNum: 1023,
	}

	hashView := NewCluster(config).EndpointSnapshot().HealthyConsistentHash()
	if !assert.NotNil(t, hashView) {
		return
	}
	_, mutable := hashView.(model.LbConsistentHash)
	assert.False(t, mutable)

	host, err := hashView.Get("request-key")
	assert.NoError(t, err)
	assert.Contains(t, []string{first.GetHost(), second.GetHost()}, host)
}

func TestClusterEndpointSnapshotBuildsConsistentHashLazily(t *testing.T) {
	var builds int32
	lbPolicy := model.LbPolicyType("test-lazy-hash")
	// Mutates the global ConsistentHashInitMap; keep usage serial
	// and do not add t.Parallel.
	previousInit, hadPreviousInit := model.ConsistentHashInitMap[lbPolicy]
	model.ConsistentHashInitMap[lbPolicy] = func(_ model.ConsistentHash, endpoints []*model.Endpoint) model.LbConsistentHash {
		atomic.AddInt32(&builds, 1)
		if len(endpoints) == 0 {
			return countingConsistentHash{}
		}
		return countingConsistentHash{host: endpoints[0].GetHost()}
	}
	t.Cleanup(func() {
		if hadPreviousInit {
			model.ConsistentHashInitMap[lbPolicy] = previousInit
			return
		}
		delete(model.ConsistentHashInitMap, lbPolicy)
	})

	first := testEndpoint("ep-1", "127.0.0.1", 18080)
	second := testEndpoint("ep-2", "127.0.0.1", 18081)
	config := testCluster("snapshot-lazy-hash", first, second)
	config.LbStr = lbPolicy
	runtimeCluster := NewCluster(config)

	// Locks: snapshot must NOT eagerly build the consistent hash. A
	// runtime that toggles health repeatedly without ever serving a
	// hash-based pick must not build a single hash.
	assert.Zero(t, atomic.LoadInt32(&builds),
		"snapshot must not eagerly build consistent hash before first access")
	assert.True(t, runtimeCluster.UpdateEndpointHealth(second.ID, second.Address.GetAddress(), false))
	assert.Zero(t, atomic.LoadInt32(&builds),
		"snapshot must not build consistent hash on health toggle alone")

	// First access populates lazily.
	snapshot := runtimeCluster.EndpointSnapshot()
	assert.NotNil(t, snapshot.HealthyConsistentHash())
	firstBuilds := atomic.LoadInt32(&builds)
	assert.Greater(t, firstBuilds, int32(0), "first HealthyConsistentHash() call must build")

	// Repeated reads of the SAME snapshot must not re-build.
	assert.NotNil(t, snapshot.HealthyConsistentHash())
	assert.Equal(t, firstBuilds, atomic.LoadInt32(&builds),
		"repeat HealthyConsistentHash() on the same snapshot must not rebuild")

	// Toggling health produces a new snapshot. The new snapshot is allowed
	// to (a) build on first access, or (b) reuse the hash from a sibling
	// snapshot when the healthy set is equivalent. Both are valid
	// implementations; this test no longer locks the count. See follow-up
	// PR-4 for the "reuse on identical healthy set" optimization.
	assert.True(t, runtimeCluster.UpdateEndpointHealth(second.ID, second.Address.GetAddress(), true))
	newSnapshot := runtimeCluster.EndpointSnapshot()
	assert.NotNil(t, newSnapshot.HealthyConsistentHash())
	// New snapshot must serve a hash view; we do not assert how many builds
	// it took to get there.
	assert.NotNil(t, newSnapshot.HealthyConsistentHash())
	finalBuilds := atomic.LoadInt32(&builds)
	assert.GreaterOrEqual(t, finalBuilds, firstBuilds,
		"build count is monotonic, but per-flap rebuild count is not locked here")
}

func TestClusterEndpointSnapshotReusesConsistentHashForUnchangedHealthySet(t *testing.T) {
	var builds int32
	lbPolicy := registerCountingConsistentHash(t, "test-reuse-hash-unchanged", &builds)

	first := testEndpoint("ep-1", "127.0.0.1", 18080)
	second := testEndpoint("ep-2", "127.0.0.1", 18081)
	config := testCluster("snapshot-reuse-hash-unchanged", first, second)
	config.LbStr = lbPolicy
	config.ConsistentHash = model.ConsistentHash{ReplicaNum: 10, MaxVnodeNum: 1023, MaglevTableSize: 521}

	previous := NewCluster(config).EndpointSnapshot()
	previousHash := previous.HealthyConsistentHash()
	assert.NotNil(t, previousHash)
	assert.Equal(t, int32(1), atomic.LoadInt32(&builds))

	next := newEndpointSnapshot(config, previous, false)
	assert.Equal(t, previousHash, next.HealthyConsistentHash(),
		"equivalent snapshot should reuse the previous hash view")
	assert.Equal(t, int32(1), atomic.LoadInt32(&builds),
		"reused snapshot must not rebuild on first access")
}

func TestClusterEndpointSnapshotDoesNotReuseConsistentHashWhenHealthChangesHealthySet(t *testing.T) {
	var builds int32
	lbPolicy := registerCountingConsistentHash(t, "test-reuse-hash-health-change", &builds)

	first := testEndpoint("ep-1", "127.0.0.1", 18080)
	second := testEndpoint("ep-2", "127.0.0.1", 18081)
	config := testCluster("snapshot-reuse-hash-health-change", first, second)
	config.LbStr = lbPolicy

	previous := NewCluster(config).EndpointSnapshot()
	assert.NotNil(t, previous.HealthyConsistentHash())
	assert.Equal(t, int32(1), atomic.LoadInt32(&builds))

	next, ok := previous.withEndpointHealthForIDs(map[string]struct{}{second.ID: {}}, false)
	assert.True(t, ok)
	assert.NotNil(t, next.HealthyConsistentHash())
	assert.Equal(t, int32(2), atomic.LoadInt32(&builds),
		"health change that changes the healthy set must force a fresh consistent hash")
}

func TestClusterEndpointSnapshotDoesNotReuseConsistentHashWhenEndpointAddressChanges(t *testing.T) {
	var builds int32
	lbPolicy := registerCountingConsistentHash(t, "test-reuse-hash-address-change", &builds)

	endpoint := testEndpoint("ep-1", "127.0.0.1", 18080)
	config := testCluster("snapshot-reuse-hash-address-change", endpoint)
	config.LbStr = lbPolicy

	previous := NewCluster(config).EndpointSnapshot()
	assert.NotNil(t, previous.HealthyConsistentHash())
	assert.Equal(t, int32(1), atomic.LoadInt32(&builds))

	config.Endpoints[0] = testEndpoint("ep-1", "127.0.0.2", 18080)
	next := newEndpointSnapshot(config, previous, false)
	assert.NotNil(t, next.HealthyConsistentHash())
	assert.Equal(t, int32(2), atomic.LoadInt32(&builds),
		"address change must force a fresh consistent hash")
}

func TestClusterEndpointSnapshotDoesNotReuseConsistentHashWhenEndpointCountChanges(t *testing.T) {
	tests := []struct {
		name      string
		endpoints func(first, second, third *model.Endpoint) []*model.Endpoint
	}{
		{
			name: "add",
			endpoints: func(first, second, third *model.Endpoint) []*model.Endpoint {
				return []*model.Endpoint{first, second, third}
			},
		},
		{
			name: "delete",
			endpoints: func(first, second, third *model.Endpoint) []*model.Endpoint {
				return []*model.Endpoint{first}
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var builds int32
			lbPolicy := registerCountingConsistentHash(t, "test-reuse-hash-count-change"+tt.name, &builds)

			first := testEndpoint("ep-1", "127.0.0.1", 18080)
			second := testEndpoint("ep-2", "127.0.0.1", 18081)
			third := testEndpoint("ep-3", "127.0.0.1", 18082)
			config := testCluster("snapshot-reuse-hash-count-change-"+tt.name, first, second)
			config.LbStr = lbPolicy

			previous := NewCluster(config).EndpointSnapshot()
			assert.NotNil(t, previous.HealthyConsistentHash())
			assert.Equal(t, int32(1), atomic.LoadInt32(&builds))

			config.Endpoints = tt.endpoints(first, second, third)
			next := newEndpointSnapshot(config, previous, false)
			assert.NotNil(t, next.HealthyConsistentHash())
			assert.Equal(t, int32(2), atomic.LoadInt32(&builds),
				"count change must force a fresh consistent hash")
		})
	}
}

func TestClusterEndpointSnapshotDoesNotReuseConsistentHashWhenHashConfigChanges(t *testing.T) {
	var builds int32
	lbPolicy := registerCountingConsistentHash(t, "test-reuse-hash-config-change", &builds)

	endpoint := testEndpoint("ep-1", "127.0.0.1", 18080)
	config := testCluster("snapshot-reuse-hash-config-change", endpoint)
	config.LbStr = lbPolicy
	config.ConsistentHash = model.ConsistentHash{ReplicaNum: 10, MaxVnodeNum: 1023, MaglevTableSize: 521}

	previous := NewCluster(config).EndpointSnapshot()
	assert.NotNil(t, previous.HealthyConsistentHash())
	assert.Equal(t, int32(1), atomic.LoadInt32(&builds))

	config.ConsistentHash.ReplicaNum = 20
	next := newEndpointSnapshot(config, previous, false)
	assert.NotNil(t, next.HealthyConsistentHash())
	assert.Equal(t, int32(2), atomic.LoadInt32(&builds),
		"hash config change must force a fresh consistent hash")
}

func TestClusterEndpointSnapshotReusesConsistentHashWhenMetadataChanges(t *testing.T) {
	var builds int32
	lbPolicy := registerCountingConsistentHash(t, "test-reuse-hash-metadata-change", &builds)

	endpoint := testEndpoint("ep-1", "127.0.0.1", 18080)
	endpoint.Metadata = map[string]string{"weight": "10"}
	config := testCluster("snapshot-reuse-hash-metadata-change", endpoint)
	config.LbStr = lbPolicy

	previous := NewCluster(config).EndpointSnapshot()
	assert.NotNil(t, previous.HealthyConsistentHash())
	assert.Equal(t, int32(1), atomic.LoadInt32(&builds))

	movedWeight := testEndpoint("ep-1", "127.0.0.1", 18080)
	movedWeight.Metadata = map[string]string{"weight": "20"}
	config.Endpoints[0] = movedWeight
	next := newEndpointSnapshot(config, previous, false)
	assert.NotNil(t, next.HealthyConsistentHash())
	assert.Equal(t, int32(1), atomic.LoadInt32(&builds),
		"metadata-only change with unchanged hosts must reuse the previous consistent hash instead of rebuilding")
}

func TestClusterEndpointSnapshotClonesConfigEndpointObjects(t *testing.T) {
	endpoint := testSnapshotEndpointWithLLMMeta()
	runtimeCluster := NewCluster(testCluster("snapshot-clone", endpoint))
	snapshotEndpoint := runtimeCluster.EndpointSnapshot().EndpointByID(endpoint.ID)
	if !assert.NotNil(t, snapshotEndpoint) {
		return
	}

	assert.NotSame(t, endpoint, snapshotEndpoint)
	assert.NotSame(t, endpoint.LLMMeta, snapshotEndpoint.LLMMeta)

	endpoint.Address.Address = "127.0.0.2"
	endpoint.Address.Domains[0] = "changed.example.com"
	endpoint.Metadata["weight"] = "9"
	endpoint.LLMMeta.APIKey = "new-key"
	endpoint.LLMMeta.RetryPolicy.Config["attempts"] = 2
	endpoint.LLMMeta.RetryPolicy.Config["nested"].(map[string]any)["delays"].([]any)[0] = "200ms"
	endpoint.UnHealthy = true

	assertEndpointMatchesOriginalSnapshot(t, runtimeCluster.EndpointSnapshot().EndpointByID(endpoint.ID))
}

func TestClusterEndpointSnapshotReturnsDefensiveEndpointObjects(t *testing.T) {
	endpoint := testSnapshotEndpointWithLLMMeta()
	runtimeCluster := NewCluster(testCluster("snapshot-defensive-endpoint", endpoint))
	snapshot := runtimeCluster.EndpointSnapshot()

	all := snapshot.AllEndpoints()
	healthy := snapshot.HealthyEndpoints()
	byID := snapshot.EndpointByID(endpoint.ID)
	healthyByID := snapshot.HealthyEndpointByID(endpoint.ID)

	assert.NotSame(t, all[0], healthy[0])
	assert.NotSame(t, all[0], byID)
	assert.NotSame(t, byID, healthyByID)

	mutateReturnedEndpoint(all[0])
	mutateReturnedEndpoint(healthy[0])
	mutateReturnedEndpoint(byID)
	mutateReturnedEndpoint(healthyByID)

	assertEndpointMatchesOriginalSnapshot(t, snapshot.EndpointByID(endpoint.ID))
	assertEndpointMatchesOriginalSnapshot(t, snapshot.HealthyEndpointByID(endpoint.ID))

	assert.True(t, runtimeCluster.UpdateEndpointHealth(endpoint.ID, endpoint.Address.GetAddress(), false))
	assert.Nil(t, runtimeCluster.EndpointSnapshot().HealthyEndpointByID(endpoint.ID))
	assert.True(t, runtimeCluster.UpdateEndpointHealth(endpoint.ID, endpoint.Address.GetAddress(), true))
	assertEndpointMatchesOriginalSnapshot(t, runtimeCluster.EndpointSnapshot().HealthyEndpointByID(endpoint.ID))
}

func TestNewClusterAssignsUniqueRuntimeIDsForAnonymousEndpoints(t *testing.T) {
	first := testEndpoint("", "127.0.0.1", 18080)
	second := testEndpoint("", "127.0.0.2", 18081)
	runtimeCluster := NewCluster(testCluster("snapshot-anonymous-ids", first, second))

	snapshot := runtimeCluster.EndpointSnapshot()
	all := snapshot.AllEndpoints()
	if !assert.Len(t, all, 2) {
		return
	}
	assert.NotEmpty(t, all[0].ID)
	assert.NotEmpty(t, all[1].ID)
	assert.NotEqual(t, all[0].ID, all[1].ID)

	assert.True(t, runtimeCluster.UpdateEndpointAddressHealth(first.Address.GetAddress(), false))
	updated := runtimeCluster.EndpointSnapshot().AllEndpoints()
	if assert.Len(t, updated, 2) {
		assert.True(t, updated[0].UnHealthy)
		assert.False(t, updated[1].UnHealthy)
	}
	healthy := runtimeCluster.EndpointSnapshot().HealthyEndpoints()
	if assert.Len(t, healthy, 1) {
		assert.Equal(t, second.Address.GetAddress(), healthy[0].Address.GetAddress())
	}
}

func TestClusterEndpointHealthEventUpdatesSnapshotWithoutMutatingEndpoint(t *testing.T) {
	endpoint := testEndpoint("ep-1", "127.0.0.1", 18082)
	runtimeCluster := NewCluster(testCluster("snapshot-health", endpoint))

	assert.True(t, runtimeCluster.UpdateEndpointHealth(endpoint.ID, endpoint.Address.GetAddress(), false))

	assert.False(t, endpoint.UnHealthy)
	assert.True(t, runtimeCluster.EndpointSnapshot().EndpointByID(endpoint.ID).UnHealthy)
	assert.Empty(t, runtimeCluster.EndpointSnapshot().HealthyEndpoints())
	assert.Nil(t, runtimeCluster.EndpointSnapshot().HealthyEndpointByID(endpoint.ID))

	assert.True(t, runtimeCluster.UpdateEndpointHealth(endpoint.ID, endpoint.Address.GetAddress(), true))
	assert.False(t, endpoint.UnHealthy)
	healthyEndpoints := runtimeCluster.EndpointSnapshot().HealthyEndpoints()
	assert.Equal(t, []*model.Endpoint{endpoint}, healthyEndpoints)
	if assert.Len(t, healthyEndpoints, 1) {
		assert.False(t, healthyEndpoints[0].UnHealthy)
	}
	assert.False(t, runtimeCluster.EndpointSnapshot().EndpointByID(endpoint.ID).UnHealthy)
}

func TestClusterEndpointHealthEventUpdatesSameAddressEndpoints(t *testing.T) {
	first := testEndpoint("ep-1", "127.0.0.1", 18082)
	second := testEndpoint("ep-2", "127.0.0.1", 18082)
	otherAddress := testEndpoint("ep-3", "127.0.0.1", 18083)
	runtimeCluster := NewCluster(testCluster("snapshot-shared-address-health", first, second, otherAddress))

	runtimeCluster.handleEndpointHealth(healthcheck.EndpointHealthEvent{
		EndpointID:      first.ID,
		EndpointAddress: first.Address.GetAddress(),
		Healthy:         false,
	})

	assert.False(t, first.UnHealthy)
	assert.False(t, second.UnHealthy)
	assert.True(t, runtimeCluster.EndpointSnapshot().EndpointByID(first.ID).UnHealthy)
	assert.True(t, runtimeCluster.EndpointSnapshot().EndpointByID(second.ID).UnHealthy)
	assert.False(t, runtimeCluster.EndpointSnapshot().EndpointByID(otherAddress.ID).UnHealthy)
	assert.Nil(t, runtimeCluster.EndpointSnapshot().HealthyEndpointByID(first.ID))
	assert.Nil(t, runtimeCluster.EndpointSnapshot().HealthyEndpointByID(second.ID))
	assert.NotNil(t, runtimeCluster.EndpointSnapshot().HealthyEndpointByID(otherAddress.ID))

	runtimeCluster.handleEndpointHealth(healthcheck.EndpointHealthEvent{
		EndpointID:      first.ID,
		EndpointAddress: first.Address.GetAddress(),
		Healthy:         true,
	})

	assert.Equal(t, []*model.Endpoint{first, second, otherAddress}, runtimeCluster.EndpointSnapshot().HealthyEndpoints())
}

func TestClusterEndpointHealthEventRestoresRuntimeEndpointHealthFlag(t *testing.T) {
	endpoint := testEndpoint("ep-1", "127.0.0.1", 18082)
	endpoint.UnHealthy = true
	runtimeCluster := NewCluster(testCluster("snapshot-health-flag", endpoint))

	assert.Empty(t, runtimeCluster.EndpointSnapshot().HealthyEndpoints())
	assert.True(t, runtimeCluster.EndpointSnapshot().EndpointByID(endpoint.ID).UnHealthy)

	assert.True(t, runtimeCluster.UpdateEndpointHealth(endpoint.ID, endpoint.Address.GetAddress(), true))

	assert.True(t, endpoint.UnHealthy)
	healthyEndpoint := runtimeCluster.EndpointSnapshot().HealthyEndpointByID(endpoint.ID)
	if assert.NotNil(t, healthyEndpoint) {
		assert.False(t, healthyEndpoint.UnHealthy)
	}
	allEndpoint := runtimeCluster.EndpointSnapshot().EndpointByID(endpoint.ID)
	if assert.NotNil(t, allEndpoint) {
		assert.False(t, allEndpoint.UnHealthy)
	}
}

func TestClusterEndpointHealthEventIgnoresStaleAddress(t *testing.T) {
	oldEndpoint := testEndpoint("ep-1", "127.0.0.1", 18083)
	oldAddress := oldEndpoint.Address.GetAddress()
	config := testCluster("snapshot-stale-address", oldEndpoint)
	runtimeCluster := NewCluster(config)

	newEndpoint := testEndpoint("ep-1", "127.0.0.2", 18084)
	config.Endpoints[0] = newEndpoint
	runtimeCluster.RefreshEndpoints()

	assert.False(t, runtimeCluster.UpdateEndpointHealth(newEndpoint.ID, oldAddress, false))
	assert.Equal(t, []*model.Endpoint{newEndpoint}, runtimeCluster.EndpointSnapshot().HealthyEndpoints())

	assert.True(t, runtimeCluster.UpdateEndpointHealth(newEndpoint.ID, newEndpoint.Address.GetAddress(), false))
	assert.Empty(t, runtimeCluster.EndpointSnapshot().HealthyEndpoints())
}

func TestClusterRefreshEndpointsFromDoesNotOverwriteNewerHealthUpdate(t *testing.T) {
	endpoint := testEndpoint("ep-1", "127.0.0.1", 18084)
	runtimeCluster := NewCluster(testClusterWithHealthCheck("snapshot-refresh-cas", endpoint))
	t.Cleanup(runtimeCluster.Stop)
	staleSnapshot := runtimeCluster.EndpointSnapshot()

	assert.True(t, runtimeCluster.UpdateEndpointHealth(endpoint.ID, endpoint.Address.GetAddress(), false))
	runtimeCluster.RefreshEndpointsFrom(staleSnapshot)

	assert.Empty(t, runtimeCluster.EndpointSnapshot().HealthyEndpoints())
	assert.Nil(t, runtimeCluster.EndpointSnapshot().HealthyEndpointByID(endpoint.ID))
}

func TestClusterSnapshotForRuntimeReplacementFreezesHealthEvents(t *testing.T) {
	endpoint := testEndpoint("ep-1", "127.0.0.1", 18085)
	runtimeCluster := NewCluster(testClusterWithHealthCheck("snapshot-freeze", endpoint))
	t.Cleanup(runtimeCluster.Stop)

	assert.True(t, runtimeCluster.UpdateEndpointHealth(endpoint.ID, endpoint.Address.GetAddress(), false))
	frozen := runtimeCluster.SnapshotForRuntimeReplacement()

	assert.Nil(t, frozen.HealthyEndpointByID(endpoint.ID))
	assert.False(t, runtimeCluster.UpdateEndpointHealth(endpoint.ID, endpoint.Address.GetAddress(), true))
	assert.Nil(t, runtimeCluster.EndpointSnapshot().HealthyEndpointByID(endpoint.ID))

	replacement := testEndpoint(endpoint.ID, "127.0.0.1", 18085)
	replacementRuntime := NewClusterWithEndpointSnapshot(
		testClusterWithHealthCheck("snapshot-freeze-replacement", replacement),
		frozen,
	)
	t.Cleanup(replacementRuntime.Stop)

	assert.Nil(t, replacementRuntime.EndpointSnapshot().HealthyEndpointByID(replacement.ID))
	assert.True(t, replacementRuntime.UpdateEndpointHealth(replacement.ID, replacement.Address.GetAddress(), true))
	assert.NotSame(t, replacement, replacementRuntime.EndpointSnapshot().HealthyEndpointByID(replacement.ID))
	assert.Equal(t, replacement, replacementRuntime.EndpointSnapshot().HealthyEndpointByID(replacement.ID))
}

func TestNewClusterWithEndpointSnapshotInheritsAddressHealthForChangedEndpointID(t *testing.T) {
	endpoint := testEndpoint("ep-old", "127.0.0.1", 18085)
	oldRuntime := NewCluster(testClusterWithHealthCheck("snapshot-address-old", endpoint))
	t.Cleanup(oldRuntime.Stop)
	oldRuntime.handleEndpointHealth(healthcheck.EndpointHealthEvent{
		EndpointID:      endpoint.ID,
		EndpointAddress: endpoint.Address.GetAddress(),
		Healthy:         false,
	})
	previous := oldRuntime.SnapshotForRuntimeReplacement()

	replacement := testEndpoint("ep-new", "127.0.0.1", 18085)
	replacementRuntime := NewClusterWithEndpointSnapshot(
		testClusterWithHealthCheck("snapshot-address-new", replacement),
		previous,
	)
	t.Cleanup(replacementRuntime.Stop)

	assert.Nil(t, replacementRuntime.EndpointSnapshot().HealthyEndpointByID(replacement.ID))
	if got := replacementRuntime.EndpointSnapshot().EndpointByID(replacement.ID); assert.NotNil(t, got) {
		assert.True(t, got.UnHealthy)
	}

	noHealthCheckReplacement := testEndpoint("ep-new-no-healthcheck", "127.0.0.1", 18085)
	noHealthCheckRuntime := NewClusterWithEndpointSnapshot(
		testCluster("snapshot-address-no-healthcheck", noHealthCheckReplacement),
		previous,
	)
	assert.NotNil(t, noHealthCheckRuntime.EndpointSnapshot().HealthyEndpointByID(noHealthCheckReplacement.ID))
}

func TestNewClusterWithEndpointSnapshotInheritsHealthOnlyWhenHealthCheckEnabled(t *testing.T) {
	endpoint := testEndpoint("ep-1", "127.0.0.1", 18085)
	oldRuntime := NewCluster(testCluster("snapshot-constructor-old", endpoint))
	assert.True(t, oldRuntime.UpdateEndpointHealth(endpoint.ID, endpoint.Address.GetAddress(), false))
	previous := oldRuntime.EndpointSnapshot()

	replacement := testEndpoint(endpoint.ID, "127.0.0.1", 18085)
	replacementRuntime := NewClusterWithEndpointSnapshot(
		testClusterWithHealthCheck("snapshot-constructor-same", replacement),
		previous,
	)
	t.Cleanup(replacementRuntime.Stop)
	assert.Nil(t, replacementRuntime.EndpointSnapshot().HealthyEndpointByID(replacement.ID))

	noHealthCheckReplacement := testEndpoint(endpoint.ID, "127.0.0.1", 18085)
	noHealthCheckRuntime := NewClusterWithEndpointSnapshot(
		testCluster("snapshot-constructor-no-healthcheck", noHealthCheckReplacement),
		previous,
	)
	assert.NotSame(t, noHealthCheckReplacement, noHealthCheckRuntime.EndpointSnapshot().HealthyEndpointByID(noHealthCheckReplacement.ID))
	assert.Equal(t, noHealthCheckReplacement, noHealthCheckRuntime.EndpointSnapshot().HealthyEndpointByID(noHealthCheckReplacement.ID))

	moved := testEndpoint(endpoint.ID, "127.0.0.2", 18086)
	movedRuntime := NewClusterWithEndpointSnapshot(
		testCluster("snapshot-constructor-moved", moved),
		previous,
	)
	assert.NotSame(t, moved, movedRuntime.EndpointSnapshot().HealthyEndpointByID(moved.ID))
	assert.Equal(t, moved, movedRuntime.EndpointSnapshot().HealthyEndpointByID(moved.ID))
}

func testCluster(name string, endpoints ...*model.Endpoint) *model.ClusterConfig {
	return &model.ClusterConfig{
		Name:      name,
		Endpoints: endpoints,
	}
}

func testClusterWithHealthCheck(name string, endpoints ...*model.Endpoint) *model.ClusterConfig {
	config := testCluster(name, endpoints...)
	config.HealthChecks = []model.HealthCheckConfig{{
		Protocol:       "tcp",
		TimeoutConfig:  "1h",
		IntervalConfig: "1h",
	}}
	return config
}

func testEndpoint(id string, host string, port int) *model.Endpoint {
	return &model.Endpoint{
		ID:   id,
		Name: "endpoint-" + id,
		Address: model.SocketAddress{
			Address: host,
			Port:    port,
		},
	}
}

func testSnapshotEndpointWithLLMMeta() *model.Endpoint {
	endpoint := testEndpoint("ep-1", "127.0.0.1", 18080)
	endpoint.Address.Domains = []string{"api.example.com"}
	endpoint.Metadata = map[string]string{"weight": "3"}
	endpoint.LLMMeta = &model.LLMMeta{
		Provider: "openai",
		APIKey:   "old-key",
		RetryPolicy: model.RetryPolicy{
			Config: map[string]any{
				"attempts": 1,
				"nested": map[string]any{
					"delays": []any{"100ms"},
				},
			},
		},
	}
	return endpoint
}

// Mutates the global ConsistentHashInitMap; keep usage serial
// and do not add t.Parallel.
func registerCountingConsistentHash(t *testing.T, name string, builds *int32) model.LbPolicyType {
	t.Helper()

	lbPolicy := model.LbPolicyType(name)
	previousInit, hadPreviousInit := model.ConsistentHashInitMap[lbPolicy]
	model.ConsistentHashInitMap[lbPolicy] = func(_ model.ConsistentHash, endpoints []*model.Endpoint) model.LbConsistentHash {
		atomic.AddInt32(builds, 1)
		if len(endpoints) == 0 {
			return countingConsistentHash{}
		}
		return countingConsistentHash{host: endpoints[0].GetHost()}
	}
	t.Cleanup(func() {
		if hadPreviousInit {
			model.ConsistentHashInitMap[lbPolicy] = previousInit
			return
		}
		delete(model.ConsistentHashInitMap, lbPolicy)
	})
	return lbPolicy
}

type countingConsistentHash struct {
	host string
}

func (h countingConsistentHash) Hash(string) uint32 {
	return 0
}

func (h countingConsistentHash) Get(string) (string, error) {
	return h.host, nil
}

func (h countingConsistentHash) GetHash(uint32) (string, error) {
	return h.host, nil
}

func (h countingConsistentHash) Add(string) {}

func (h countingConsistentHash) Remove(string) bool {
	return false
}

func mutateReturnedEndpoint(endpoint *model.Endpoint) {
	endpoint.Name = "mutated"
	endpoint.Address.Address = "127.0.0.2"
	endpoint.Address.Domains[0] = "changed.example.com"
	endpoint.Metadata["weight"] = "9"
	endpoint.LLMMeta.APIKey = "new-key"
	endpoint.LLMMeta.RetryPolicy.Config["attempts"] = 2
	endpoint.LLMMeta.RetryPolicy.Config["nested"].(map[string]any)["delays"].([]any)[0] = "200ms"
	endpoint.UnHealthy = true
}

func assertEndpointMatchesOriginalSnapshot(t *testing.T, endpoint *model.Endpoint) {
	t.Helper()

	if !assert.NotNil(t, endpoint) {
		return
	}
	assert.Equal(t, "endpoint-ep-1", endpoint.Name)
	assert.Equal(t, model.SocketAddress{Address: "127.0.0.1", Port: 18080, Domains: []string{"api.example.com"}}, endpoint.Address)
	assert.Equal(t, map[string]string{"weight": "3"}, endpoint.Metadata)
	assert.Equal(t, "old-key", endpoint.LLMMeta.APIKey)
	assert.Equal(t, 1, endpoint.LLMMeta.RetryPolicy.Config["attempts"])
	assert.Equal(t, "100ms", endpoint.LLMMeta.RetryPolicy.Config["nested"].(map[string]any)["delays"].([]any)[0])
	assert.False(t, endpoint.UnHealthy)
}

// TestClusterRefreshAndHealthUpdateConcurrent locks the CAS contract that
// RefreshEndpoints and UpdateEndpointHealth must commute under concurrent
// writers. The final snapshot's health for ep-1 must reflect the LAST
// UpdateEndpointHealth call (`(N-1)%2 == 0` → healthy=true at completion),
// and ep-2 must remain healthy because no update ever toggled it.
//
// Run with -race -count=20 in CI; flakiness here means the snapshot
// inheritance / CAS invariants are broken.
func TestClusterRefreshAndHealthUpdateConcurrent(t *testing.T) {
	ep1 := testEndpoint("ep-1", "127.0.0.1", 18080)
	ep2 := testEndpoint("ep-2", "127.0.0.1", 18081)
	config := testClusterWithHealthCheck("race", ep1, ep2)
	runtimeCluster := NewCluster(config)
	t.Cleanup(runtimeCluster.Stop)

	const N = 200
	var wg sync.WaitGroup
	wg.Add(2)

	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			runtimeCluster.RefreshEndpoints()
		}
	}()
	go func() {
		defer wg.Done()
		for i := 0; i < N; i++ {
			runtimeCluster.UpdateEndpointHealth(ep1.ID, ep1.Address.GetAddress(), i%2 == 0)
		}
	}()
	wg.Wait()

	snapshot := runtimeCluster.EndpointSnapshot()
	healthy := snapshot.HealthyEndpointByID(ep1.ID)
	if (N-1)%2 == 0 {
		assert.NotNil(t, healthy, "last update was healthy=true")
	} else {
		assert.Nil(t, healthy, "last update was healthy=false")
	}
	assert.NotNil(t, snapshot.HealthyEndpointByID(ep2.ID), "ep2 must remain healthy")
}

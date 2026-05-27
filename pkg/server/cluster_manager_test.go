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

package server

import (
	"fmt"
	"reflect"
	"sync"
	"sync/atomic"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cluster"
	_ "github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/maglev"     // Register Maglev for cluster-manager tests.
	_ "github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/rand"       // Register Rand for cluster-manager tests.
	_ "github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/ringhash"   // Register RingHash for cluster-manager tests.
	_ "github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/roundrobin" // Register RoundRobin for cluster-manager tests.
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestClusterManager(t *testing.T) {
	cm := testClusterManager(
		testCluster("test", model.LoadBalancerRoundRobin, []*model.Endpoint{
			testEndpoint("1", "127.0.0.1", 18080),
		}),
	)

	assert.Len(t, cm.store.Config, 1)

	cm.AddCluster(testCluster("test2", model.LoadBalancerRoundRobin, []*model.Endpoint{
		testEndpoint("1", "127.0.0.1", 18081),
	}))

	assert.Len(t, cm.store.Config, 2)

	cm.SetEndpoint("test2", testEndpoint("2", "127.0.0.1", 18082))
	assert.Equal(t, "1", cm.PickEndpoint("test", nil).ID)
	cm.DeleteEndpoint("test2", "1")
}

func TestClusterManager_PickEndpointReturnsNilForMissingCluster(t *testing.T) {
	cm := testClusterManager()
	assert.Nil(t, cm.PickEndpoint("missing-cluster", nil))
}

func TestClusterManager_PickEndpointUsesRuntimeClusterMap(t *testing.T) {
	runtimeConfig := testCluster("runtime-lookup", model.LoadBalancerRoundRobin, []*model.Endpoint{
		testEndpoint("runtime-ep", "127.0.0.1", 18083),
	})
	cm := testClusterManager(runtimeConfig)
	cm.store.Config = []*model.ClusterConfig{
		testCluster(runtimeConfig.Name, model.LoadBalancerRoundRobin, []*model.Endpoint{
			testEndpoint("stale-ep", "127.0.0.1", 18084),
		}),
	}

	endpoint := cm.PickEndpoint(runtimeConfig.Name, nil)

	if assert.NotNil(t, endpoint) {
		assert.Equal(t, "runtime-ep", endpoint.ID)
	}
}

func TestClusterManager_PickNextEndpointUsesRuntimeClusterMap(t *testing.T) {
	runtimeConfig := testCluster("runtime-next-lookup", model.LoadBalancerRoundRobin, []*model.Endpoint{
		testEndpoint("runtime-1", "127.0.0.1", 18085),
		testEndpoint("runtime-2", "127.0.0.1", 18086),
	})
	cm := testClusterManager(runtimeConfig)
	cm.store.Config = []*model.ClusterConfig{
		testCluster(runtimeConfig.Name, model.LoadBalancerRoundRobin, []*model.Endpoint{
			testEndpoint("stale-1", "127.0.0.1", 18087),
			testEndpoint("stale-2", "127.0.0.1", 18088),
		}),
	}

	endpoint := cm.PickNextEndpoint(runtimeConfig.Name, "runtime-1")

	if assert.NotNil(t, endpoint) {
		assert.Equal(t, "runtime-2", endpoint.ID)
	}
}

func TestClusterManager_GetEndpointByIDUsesRuntimeClusterMap(t *testing.T) {
	runtimeConfig := testCluster("runtime-id-lookup", model.LoadBalancerRoundRobin, []*model.Endpoint{
		testEndpoint("runtime-ep", "127.0.0.1", 18089),
	})
	cm := testClusterManager(runtimeConfig)
	cm.store.Config = []*model.ClusterConfig{
		testCluster(runtimeConfig.Name, model.LoadBalancerRoundRobin, []*model.Endpoint{
			testEndpoint("stale-ep", "127.0.0.1", 18090),
		}),
	}

	endpoint := cm.GetEndpointByID(runtimeConfig.Name, "runtime-ep")

	if assert.NotNil(t, endpoint) {
		assert.Equal(t, "runtime-ep", endpoint.ID)
	}
}

func TestClusterManager_PickEndpointSingleUnhealthyReturnsNil(t *testing.T) {
	cm := testClusterManager(
		testCluster("single-unhealthy", model.LoadBalancerRoundRobin, []*model.Endpoint{
			{
				ID:        "ep-1",
				Name:      "endpoint-ep-1",
				UnHealthy: true,
				Address: model.SocketAddress{
					Address: "127.0.0.1",
					Port:    18090,
				},
			},
		}),
	)

	assert.Nil(t, cm.PickEndpoint("single-unhealthy", nil))
}

func TestClusterManager_PickEndpointAllUnhealthyReturnsNil(t *testing.T) {
	cm := testClusterManager(
		testCluster("all-unhealthy", model.LoadBalancerRoundRobin, []*model.Endpoint{
			{
				ID:        "ep-1",
				Name:      "endpoint-ep-1",
				UnHealthy: true,
				Address: model.SocketAddress{
					Address: "127.0.0.1",
					Port:    18091,
				},
			},
			{
				ID:        "ep-2",
				Name:      "endpoint-ep-2",
				UnHealthy: true,
				Address: model.SocketAddress{
					Address: "127.0.0.1",
					Port:    18092,
				},
			},
		}),
	)

	assert.Nil(t, cm.PickEndpoint("all-unhealthy", nil))
}

func TestClusterManager_CompareAndSetStorePreservesRoundRobinCursorAcrossRefresh(t *testing.T) {
	cluster := testCluster("refresh-round-robin", model.LoadBalancerRoundRobin, []*model.Endpoint{
		testEndpoint("ep-1", "127.0.0.1", 19200),
		testEndpoint("ep-2", "127.0.0.1", 19201),
		testEndpoint("ep-3", "127.0.0.1", 19202),
	})
	cm := testClusterManager(cluster)

	const expectedCursor uint32 = 5
	atomic.StoreUint32(&cm.store.Config[0].PrePickEndpointIndex, expectedCursor)

	oldStore, err := cm.CloneStore()
	if !assert.NoError(t, err) {
		return
	}
	newStore := cm.NewStore(oldStore.Version)
	for _, endpoint := range oldStore.Config[0].Endpoints {
		copied := *endpoint
		newStore.SetEndpoint(cluster.Name, &copied)
	}

	assert.True(t, cm.CompareAndSetStore(newStore))
	if assert.Len(t, cm.store.Config, 1) {
		assert.Equal(t, expectedCursor, atomic.LoadUint32(&cm.store.Config[0].PrePickEndpointIndex))
	}

	endpoint := cm.PickEndpoint(cluster.Name, nil)
	if assert.NotNil(t, endpoint) {
		assert.Equal(t, "ep-3", endpoint.ID)
	}
}

func TestClusterManager_UpdateClusterRebuildsRuntimeCluster(t *testing.T) {
	oldConfig := testCluster("runtime-update", model.LoadBalancerRoundRobin, []*model.Endpoint{
		testEndpoint("ep-1", "127.0.0.1", 19300),
	}, testHealthCheck())
	cm := testClusterManager(oldConfig)
	defer stopStoreRuntimes(cm.store)

	oldRuntime := cm.store.clustersMap[oldConfig.Name]
	if !assert.NotNil(t, oldRuntime) {
		return
	}
	assert.Same(t, oldConfig, oldRuntime.Config)
	assert.Greater(t, healthCheckersLen(oldRuntime), 0)

	const expectedCursor uint32 = 11
	atomic.StoreUint32(&oldConfig.PrePickEndpointIndex, expectedCursor)

	newConfig := testCluster(oldConfig.Name, model.LoadBalancerRoundRobin, []*model.Endpoint{
		testEndpoint("ep-2", "127.0.0.1", 19301),
	}, testHealthCheck())

	cm.UpdateCluster(newConfig)

	newRuntime := cm.store.clustersMap[oldConfig.Name]
	if !assert.NotNil(t, newRuntime) {
		return
	}
	assert.NotSame(t, oldRuntime, newRuntime)
	assert.Same(t, newConfig, newRuntime.Config)
	assert.Same(t, newConfig, cm.store.Config[0])
	assert.Equal(t, expectedCursor, atomic.LoadUint32(&newConfig.PrePickEndpointIndex))
	assert.Equal(t, 0, healthCheckersLen(oldRuntime))
	assert.Greater(t, healthCheckersLen(newRuntime), 0)
}

func TestClusterManager_CompareAndSetStoreVersionMismatchHasNoSideEffects(t *testing.T) {
	currentConfig := testCluster("cas-version-mismatch", model.LoadBalancerRoundRobin, []*model.Endpoint{
		testEndpoint("ep-1", "127.0.0.1", 19310),
	}, testHealthCheck())
	cm := testClusterManager(currentConfig)
	defer stopStoreRuntimes(cm.store)

	oldStore := cm.store
	oldRuntime := oldStore.clustersMap[currentConfig.Name]
	if !assert.NotNil(t, oldRuntime) {
		return
	}
	assert.Greater(t, healthCheckersLen(oldRuntime), 0)

	candidateConfig := testCluster(currentConfig.Name, model.LoadBalancerRoundRobin, []*model.Endpoint{
		testEndpoint("ep-2", "127.0.0.1", 19311),
	})
	candidate := &ClusterStore{
		Config:  []*model.ClusterConfig{candidateConfig},
		Version: oldStore.Version + 1,
	}

	assert.False(t, cm.CompareAndSetStore(candidate))
	assert.Same(t, oldStore, cm.store)
	assert.Nil(t, candidate.clustersMap)
	assert.Greater(t, healthCheckersLen(oldRuntime), 0)
}

func TestClusterManager_CompareAndSetStoreEnsuresRuntimeAndStopsOld(t *testing.T) {
	oldConfig := testCluster("cas-runtime-refresh", model.LoadBalancerRoundRobin, []*model.Endpoint{
		testEndpoint("ep-1", "127.0.0.1", 19320),
	}, testHealthCheck())
	cm := testClusterManager(oldConfig)

	oldRuntime := cm.store.clustersMap[oldConfig.Name]
	if !assert.NotNil(t, oldRuntime) {
		return
	}
	assert.Greater(t, healthCheckersLen(oldRuntime), 0)

	const expectedCursor uint32 = 17
	atomic.StoreUint32(&oldConfig.PrePickEndpointIndex, expectedCursor)

	newConfig := testCluster(oldConfig.Name, model.LoadBalancerRoundRobin, []*model.Endpoint{
		testEndpoint("ep-2", "127.0.0.1", 19321),
	}, testHealthCheck())
	candidate := &ClusterStore{
		Config:  []*model.ClusterConfig{newConfig},
		Version: cm.store.Version,
	}

	assert.True(t, cm.CompareAndSetStore(candidate))
	defer stopStoreRuntimes(cm.store)

	newRuntime := candidate.clustersMap[newConfig.Name]
	if !assert.NotNil(t, newRuntime) {
		return
	}
	assert.Same(t, candidate, cm.store)
	assert.NotSame(t, oldRuntime, newRuntime)
	assert.Same(t, newConfig, newRuntime.Config)
	assert.Equal(t, expectedCursor, atomic.LoadUint32(&newConfig.PrePickEndpointIndex))
	assert.Equal(t, 0, healthCheckersLen(oldRuntime))
	assert.Greater(t, healthCheckersLen(newRuntime), 0)
}

func TestClusterManager_CompareAndSetStoreStopsRemovedRuntime(t *testing.T) {
	oldConfig := testCluster("cas-runtime-removed", model.LoadBalancerRoundRobin, []*model.Endpoint{
		testEndpoint("ep-1", "127.0.0.1", 19325),
	}, testHealthCheck())
	cm := testClusterManager(oldConfig)

	oldRuntime := cm.store.clustersMap[oldConfig.Name]
	if !assert.NotNil(t, oldRuntime) {
		return
	}
	assert.Greater(t, healthCheckersLen(oldRuntime), 0)

	candidate := &ClusterStore{
		Version: cm.store.Version,
	}

	assert.True(t, cm.CompareAndSetStore(candidate))
	assert.Same(t, candidate, cm.store)
	assert.Empty(t, cm.store.Config)
	assert.NotContains(t, cm.store.clustersMap, oldConfig.Name)
	assert.Equal(t, 0, healthCheckersLen(oldRuntime))
}

func TestClusterStore_EnsureRuntimeClustersRepairsRuntimeMap(t *testing.T) {
	t.Run("nil map", func(t *testing.T) {
		config := testCluster("ensure-nil-map", model.LoadBalancerRoundRobin, []*model.Endpoint{
			testEndpoint("ep-1", "127.0.0.1", 19330),
		})
		store := &ClusterStore{Config: []*model.ClusterConfig{config}}

		replaced := store.ensureRuntimeClusters()
		defer stopStoreRuntimes(store)

		assert.Empty(t, replaced)
		if assert.NotNil(t, store.clustersMap[config.Name]) {
			assert.Same(t, config, store.clustersMap[config.Name].Config)
		}
	})

	t.Run("missing mismatched stale and idempotent", func(t *testing.T) {
		correctConfig := testCluster("ensure-correct", model.LoadBalancerRoundRobin, []*model.Endpoint{
			testEndpoint("ep-1", "127.0.0.1", 19331),
		})
		missingConfig := testCluster("ensure-missing", model.LoadBalancerRoundRobin, []*model.Endpoint{
			testEndpoint("ep-2", "127.0.0.1", 19332),
		})
		oldMismatchedConfig := testCluster("ensure-mismatch", model.LoadBalancerRoundRobin, []*model.Endpoint{
			testEndpoint("old", "127.0.0.1", 19333),
		})
		newMismatchedConfig := testCluster("ensure-mismatch", model.LoadBalancerRoundRobin, []*model.Endpoint{
			testEndpoint("new", "127.0.0.1", 19334),
		})
		staleConfig := testCluster("ensure-stale", model.LoadBalancerRoundRobin, []*model.Endpoint{
			testEndpoint("stale", "127.0.0.1", 19335),
		})

		correctRuntime := cluster.NewCluster(correctConfig)
		mismatchedRuntime := cluster.NewCluster(oldMismatchedConfig)
		staleRuntime := cluster.NewCluster(staleConfig)
		store := &ClusterStore{
			Config: []*model.ClusterConfig{
				correctConfig,
				missingConfig,
				newMismatchedConfig,
			},
			clustersMap: map[string]*cluster.Cluster{
				correctConfig.Name:       correctRuntime,
				newMismatchedConfig.Name: mismatchedRuntime,
				staleConfig.Name:         staleRuntime,
			},
		}
		defer stopStoreRuntimes(store)

		replaced := store.ensureRuntimeClusters()
		stopClusters(replaced)

		assert.Same(t, correctRuntime, store.clustersMap[correctConfig.Name])
		if assert.NotNil(t, store.clustersMap[missingConfig.Name]) {
			assert.Same(t, missingConfig, store.clustersMap[missingConfig.Name].Config)
		}
		if assert.NotNil(t, store.clustersMap[newMismatchedConfig.Name]) {
			assert.NotSame(t, mismatchedRuntime, store.clustersMap[newMismatchedConfig.Name])
			assert.Same(t, newMismatchedConfig, store.clustersMap[newMismatchedConfig.Name].Config)
		}
		assert.NotContains(t, store.clustersMap, staleConfig.Name)
		assert.Contains(t, replaced, mismatchedRuntime)
		assert.Contains(t, replaced, staleRuntime)

		runtimesAfterRepair := map[string]*cluster.Cluster{}
		for name, runtime := range store.clustersMap {
			runtimesAfterRepair[name] = runtime
		}

		assert.Empty(t, store.ensureRuntimeClusters())
		assert.Equal(t, runtimesAfterRepair, store.clustersMap)
	})
}

// TestClusterManager_SetEndpointSuffixAppendRebuildsConsistentHash locks the
// invariant that when a SetEndpoint call collides on explicit ID with an
// existing endpoint but differs in routing identity, BOTH endpoints survive
// (the second gets a -2 suffix) and the consistent hash reflects both hosts.
// Previously the colliding call silently overwrote the original endpoint,
// which made the PR-3 dedup rule effective only for static assembly and not
// for the dynamic LLM/Nacos registration path.
func TestClusterManager_SetEndpointSuffixAppendRebuildsConsistentHash(t *testing.T) {
	tests := []model.LbPolicyType{
		model.LoadBalancerRingHashing,
		model.LoadBalancerMaglevHashing,
	}

	for _, lb := range tests {
		t.Run(string(lb), func(t *testing.T) {
			oldEndpoint := testEndpoint("ep-1", "127.0.0.1", 19340)
			oldHost := oldEndpoint.GetHost()
			config := testCluster(fmt.Sprintf("hash-update-%s", lb), lb, []*model.Endpoint{oldEndpoint})
			cm := testClusterManager(config)
			defer stopStoreRuntimes(cm.store)

			newEndpoint := testEndpoint("ep-1", "127.0.0.2", 19341)
			newHost := newEndpoint.GetHost()
			cm.SetEndpoint(config.Name, newEndpoint)

			endpoints := cm.store.Config[0].Endpoints
			if !assert.Len(t, endpoints, 2, "duplicate explicit ID must produce a suffixed sibling, not overwrite") {
				return
			}
			assert.Equal(t, "ep-1", endpoints[0].ID)
			assert.Equal(t, "ep-1-2", endpoints[1].ID)

			hash := cm.store.Config[0].ConsistentHash.Hash
			if !assert.NotNil(t, hash) {
				return
			}
			if hostList, ok := hash.(interface{ Hosts() []string }); ok {
				hosts := hostList.Hosts()
				assert.Contains(t, hosts, oldHost, "old host stays in the hash because the slot is preserved")
				assert.Contains(t, hosts, newHost, "new host is rebuilt into the hash because the slot is appended")
				return
			}
			assert.True(t, hash.Remove(oldHost), "old host stays in the hash because the slot is preserved")
			assert.True(t, hash.Remove(newHost), "new host is rebuilt into the hash because the slot is appended")
		})
	}
}

func TestClusterManager_DeleteEndpointRepairsRuntimeAndConsistentHash(t *testing.T) {
	deletedEndpoint := testEndpoint("ep-1", "127.0.0.1", 19350)
	remainingEndpoint := testEndpoint("ep-2", "127.0.0.1", 19351)
	config := testCluster("delete-runtime-repair", model.LoadBalancerRingHashing, []*model.Endpoint{
		deletedEndpoint,
		remainingEndpoint,
	})
	cm := testClusterManager(config)
	defer stopStoreRuntimes(cm.store)

	staleConfig := testCluster(config.Name, model.LoadBalancerRingHashing, []*model.Endpoint{
		testEndpoint("stale", "127.0.0.1", 19352),
	})
	staleRuntime := cluster.NewCluster(staleConfig)
	cm.store.clustersMap[config.Name] = staleRuntime

	deletedHost := deletedEndpoint.GetHost()
	remainingHost := remainingEndpoint.GetHost()

	cm.DeleteEndpoint(config.Name, deletedEndpoint.ID)

	runtime := cm.store.clustersMap[config.Name]
	if !assert.NotNil(t, runtime) {
		return
	}
	assert.NotSame(t, staleRuntime, runtime)
	assert.Same(t, config, runtime.Config)
	if assert.Len(t, config.Endpoints, 1) {
		assert.Equal(t, remainingEndpoint, config.Endpoints[0])
		assert.NotSame(t, remainingEndpoint, config.Endpoints[0])
	}

	hash := config.ConsistentHash.Hash
	if !assert.NotNil(t, hash) {
		return
	}
	hostList, ok := hash.(interface{ Hosts() []string })
	if !assert.True(t, ok) {
		return
	}
	hosts := hostList.Hosts()
	assert.NotContains(t, hosts, deletedHost)
	assert.Contains(t, hosts, remainingHost)
}

func TestClusterManager_Race_RoundRobinPickEndpoint(t *testing.T) {
	cluster := testCluster("race-round-robin", model.LoadBalancerRoundRobin, []*model.Endpoint{
		testEndpoint("ep-1", "127.0.0.1", 19100),
		testEndpoint("ep-2", "127.0.0.1", 19101),
		testEndpoint("ep-3", "127.0.0.1", 19102),
		testEndpoint("ep-4", "127.0.0.1", 19103),
	})
	cm := testClusterManager(cluster)

	start := make(chan struct{})
	var wg sync.WaitGroup

	for i := 0; i < 8; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			<-start
			for j := 0; j < 4000; j++ {
				_ = cm.PickEndpoint(cluster.Name, nil)
			}
		}()
	}

	close(start)
	wg.Wait()
}

func TestClusterManager_AssembleEndpointsAssignsDeterministicID(t *testing.T) {
	build := func() *model.ClusterConfig {
		return &model.ClusterConfig{
			Name:  "assemble-id",
			LbStr: model.LoadBalancerRoundRobin,
			Endpoints: []*model.Endpoint{
				{Address: model.SocketAddress{Address: "127.0.0.1", Port: 21080}},
				{Address: model.SocketAddress{Address: "127.0.0.1", Port: 21081}},
			},
		}
	}

	first := testClusterManager(build())
	defer stopStoreRuntimes(first.store)
	second := testClusterManager(build())
	defer stopStoreRuntimes(second.store)

	firstEPs := first.store.Config[0].Endpoints
	secondEPs := second.store.Config[0].Endpoints

	if !assert.Len(t, firstEPs, 2) || !assert.Len(t, secondEPs, 2) {
		return
	}

	// Deterministic: identical static config produces identical IDs across
	// independent constructions, so dashboards keyed on endpoint.ID survive
	// process restart.
	assert.Equal(t, firstEPs[0].ID, secondEPs[0].ID)
	assert.Equal(t, firstEPs[1].ID, secondEPs[1].ID)

	// Generated prefix indicates the deterministic helper, not the legacy
	// random UUID fallback.
	assert.Contains(t, firstEPs[0].ID, "pixiu-generated-endpoint-")

	// Endpoints differing only by port must not collide within the same cluster.
	assert.NotEqual(t, firstEPs[0].ID, firstEPs[1].ID)
}

func TestClusterManager_AddClusterDoesNotMutateUserSuppliedEndpoint(t *testing.T) {
	endpoint := &model.Endpoint{
		Address: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    21079,
		},
	}
	cluster := &model.ClusterConfig{
		Name:      "add-clone-boundary",
		LbStr:     model.LoadBalancerRoundRobin,
		Endpoints: []*model.Endpoint{endpoint},
	}

	cm := testClusterManager(cluster)
	defer stopStoreRuntimes(cm.store)

	assert.Empty(t, endpoint.ID)
	assert.Empty(t, endpoint.Name)
	assert.NotSame(t, endpoint, cm.store.Config[0].Endpoints[0])
	assert.NotEmpty(t, cm.store.Config[0].Endpoints[0].ID)
	assert.NotEmpty(t, cm.store.Config[0].Endpoints[0].Name)
}

func TestClusterManager_SetEndpointDoesNotMutateInputEndpoint(t *testing.T) {
	cm := testClusterManager(testCluster("set-clone-boundary", model.LoadBalancerRoundRobin, nil))
	defer stopStoreRuntimes(cm.store)

	endpoint := &model.Endpoint{
		Address: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    21078,
		},
	}

	cm.SetEndpoint("set-clone-boundary", endpoint)

	assert.Empty(t, endpoint.ID)
	assert.Empty(t, endpoint.Name)
	if assert.Len(t, cm.store.Config[0].Endpoints, 1) {
		assert.NotSame(t, endpoint, cm.store.Config[0].Endpoints[0])
		assert.NotEmpty(t, cm.store.Config[0].Endpoints[0].ID)
		assert.NotEmpty(t, cm.store.Config[0].Endpoints[0].Name)
	}
}

// TestClusterManager_SetEndpointSuffixAppendKeepsOldLLMMeta verifies that
// when an incoming SetEndpoint call collides on explicit ID but differs in
// routing identity, the original endpoint's LLMMeta is left untouched
// because the new endpoint is appended with a -2 suffix instead of merged
// into the old slot. This is the dynamic-path counterpart to the static
// dedup test TestAssembleEndpointsDeduplicatesExplicitID.
func TestClusterManager_SetEndpointSuffixAppendKeepsOldLLMMeta(t *testing.T) {
	oldEndpoint := testEndpoint("ep-1", "127.0.0.1", 21077)
	oldEndpoint.LLMMeta = &model.LLMMeta{APIKey: "old-key"}
	config := testCluster("set-llm-meta-clone", model.LoadBalancerRoundRobin, []*model.Endpoint{oldEndpoint})
	cm := testClusterManager(config)
	defer stopStoreRuntimes(cm.store)

	storedOld := cm.store.Config[0].Endpoints[0]
	incoming := testEndpoint("ep-1", "127.0.0.2", 21078)

	cm.SetEndpoint(config.Name, incoming)

	endpoints := cm.store.Config[0].Endpoints
	if !assert.Len(t, endpoints, 2,
		"same explicit ID + different address must append a suffixed endpoint, not overwrite") {
		return
	}
	assert.Equal(t, "ep-1", endpoints[0].ID)
	assert.Equal(t, "ep-1-2", endpoints[1].ID)

	if assert.NotNil(t, endpoints[0].LLMMeta) {
		assert.Equal(t, "old-key", endpoints[0].LLMMeta.APIKey,
			"the old slot's LLMMeta survives an additive SetEndpoint untouched")
	}
	assert.Equal(t, storedOld.Address, endpoints[0].Address,
		"the old slot's address must not be overwritten by the colliding call")
	assert.Nil(t, endpoints[1].LLMMeta, "the appended endpoint inherits nothing from the old slot")
	assert.Nil(t, incoming.LLMMeta, "input endpoint is not mutated")
}

// TestClusterManager_SetEndpointExplicitSameIDDifferentAddressAppendsSuffix
// is the reviewer's direct repro for PR-3 dynamic-path dedup. Two
// SetEndpoint calls that share an explicit ID but route to different
// addresses must produce two endpoints (the second suffixed -2), not a
// silent merge.
func TestClusterManager_SetEndpointExplicitSameIDDifferentAddressAppendsSuffix(t *testing.T) {
	cm := testClusterManager(testCluster("c", model.LoadBalancerRoundRobin, nil))
	defer stopStoreRuntimes(cm.store)

	cm.SetEndpoint("c", testEndpoint("foo", "127.0.0.1", 21100))
	cm.SetEndpoint("c", testEndpoint("foo", "127.0.0.2", 21101))

	endpoints := cm.store.Config[0].Endpoints
	if !assert.Len(t, endpoints, 2,
		"two SetEndpoint calls with the same explicit ID but different addresses "+
			"must keep both endpoints, not silently overwrite the first") {
		return
	}
	assert.Equal(t, "foo", endpoints[0].ID)
	assert.Equal(t, "foo-2", endpoints[1].ID)
	assert.Equal(t, "127.0.0.1", endpoints[0].Address.Address)
	assert.Equal(t, "127.0.0.2", endpoints[1].Address.Address)
}

// TestClusterManager_SetEndpointExplicitSameIDSameContentIsIdempotent locks
// the other side of the dedup contract: when two calls share the same
// explicit ID AND the same routing-relevant content, the second call is a
// re-registration (idempotent), not a duplicate that should accumulate.
// This is the safety net so Nacos heartbeats don't grow the endpoint slice.
func TestClusterManager_SetEndpointExplicitSameIDSameContentIsIdempotent(t *testing.T) {
	cm := testClusterManager(testCluster("c", model.LoadBalancerRoundRobin, nil))
	defer stopStoreRuntimes(cm.store)

	cm.SetEndpoint("c", testEndpoint("foo", "127.0.0.1", 21110))
	cm.SetEndpoint("c", testEndpoint("foo", "127.0.0.1", 21110))

	endpoints := cm.store.Config[0].Endpoints
	if !assert.Len(t, endpoints, 1,
		"same-content re-registration must be idempotent, never accumulate -N entries") {
		return
	}
	assert.Equal(t, "foo", endpoints[0].ID)
}

// TestClusterManager_SetEndpointEmptyIDHashCollisionDifferentContentSuffixes
// covers the empty-ID variant of the dedup contract. Two SetEndpoint calls
// with empty IDs and the same hash material (so they resolve to the same
// generated-* ID) but different non-hash content (Metadata) must produce
// two endpoints — the second suffixed -2 — to stay consistent with how
// assembleClusterEndpoints handles two static entries with identical hash
// material.
func TestClusterManager_SetEndpointEmptyIDHashCollisionDifferentContentSuffixes(t *testing.T) {
	cm := testClusterManager(testCluster("c", model.LoadBalancerRoundRobin, nil))
	defer stopStoreRuntimes(cm.store)

	addr := model.SocketAddress{Address: "127.0.0.1", Port: 21120}
	cm.SetEndpoint("c", &model.Endpoint{
		Address:  addr,
		Metadata: map[string]string{"region": "us-east"},
	})
	cm.SetEndpoint("c", &model.Endpoint{
		Address:  addr,
		Metadata: map[string]string{"region": "us-west"},
	})

	endpoints := cm.store.Config[0].Endpoints
	if !assert.Len(t, endpoints, 2,
		"empty-ID calls with matching hash but differing Metadata must survive as siblings") {
		return
	}
	baseID := model.GenerateEndpointID("c", &model.Endpoint{Address: addr})
	assert.Equal(t, baseID, endpoints[0].ID)
	assert.Equal(t, baseID+"-2", endpoints[1].ID)
}

// TestClusterManager_SetEndpointEmptyIDSameContentIsIdempotent guarantees
// the natural idempotency of dynamic re-registration when callers leave the
// ID empty: two byte-equal empty-ID SetEndpoint calls collapse to one
// runtime endpoint, just like an explicit-ID idempotent call would.
func TestClusterManager_SetEndpointEmptyIDSameContentIsIdempotent(t *testing.T) {
	cm := testClusterManager(testCluster("c", model.LoadBalancerRoundRobin, nil))
	defer stopStoreRuntimes(cm.store)

	addr := model.SocketAddress{Address: "127.0.0.1", Port: 21130}
	cm.SetEndpoint("c", &model.Endpoint{Address: addr})
	cm.SetEndpoint("c", &model.Endpoint{Address: addr})

	endpoints := cm.store.Config[0].Endpoints
	assert.Len(t, endpoints, 1,
		"empty-ID re-registration with identical content must collapse to one entry")
}

// TestClusterManager_SetEndpointEmptyIDReusesOperatorPinnedID locks truth-
// table row 3: when an empty-ID SetEndpoint call hash-matches an existing
// endpoint that carries an operator-pinned (non-generated) ID and is
// content-equal, the pinned ID is reused — the dynamic path does NOT
// override the operator's choice with a generated-* hash. Without this
// test the empty-ID branch of resolveSetEndpointSlot could be silently
// regressed to `return incomingHash, true` in a future refactor.
//
// Non-obvious dependency: assembleClusterEndpoints defaults the existing
// pinned slot's Name to "endpoint-1" on first load while the incoming
// arrives with Name="". The test passes because endpointContentEqualForSet
// excludes Name from the equality check — see that helper's godoc. A
// refactor that re-includes Name in equality would make this test fail
// for a confusing reason (apparent content mismatch even though Address
// and LLMMeta are identical); update endpointContentEqualForSet's
// contract first if you change that.
func TestClusterManager_SetEndpointEmptyIDReusesOperatorPinnedID(t *testing.T) {
	addr := model.SocketAddress{Address: "127.0.0.1", Port: 21140}
	pinned := &model.Endpoint{ID: "operator-pinned", Address: addr}
	cm := testClusterManager(testCluster("c", model.LoadBalancerRoundRobin, []*model.Endpoint{pinned}))
	defer stopStoreRuntimes(cm.store)

	cm.SetEndpoint("c", &model.Endpoint{Address: addr})

	endpoints := cm.store.Config[0].Endpoints
	if !assert.Len(t, endpoints, 1,
		"empty-ID call with matching hash and content must reuse the pinned slot, not append") {
		return
	}
	assert.Equal(t, "operator-pinned", endpoints[0].ID,
		"the operator's pinned ID must survive empty-ID re-registration")
}

// TestClusterManager_SetEndpointNameOnlyChangeIsIdempotent locks the
// deliberate trade-off documented on ClusterManager.SetEndpoint: Name is
// excluded from the dedup content check, so a SetEndpoint call whose only
// difference from an existing slot is Name is treated as an idempotent
// re-registration. The original Name stays in place. Callers that need to
// rename must DeleteEndpoint + SetEndpoint, or update the cluster config.
//
// Without this test we cannot tell whether the Name-not-applied behavior
// is a deliberate contract or a latent bug, and the next refactor could
// flip it silently.
func TestClusterManager_SetEndpointNameOnlyChangeIsIdempotent(t *testing.T) {
	addr := model.SocketAddress{Address: "127.0.0.1", Port: 21150}
	original := &model.Endpoint{ID: "ep-1", Name: "original-name", Address: addr}
	cm := testClusterManager(testCluster("c", model.LoadBalancerRoundRobin, []*model.Endpoint{original}))
	defer stopStoreRuntimes(cm.store)

	cm.SetEndpoint("c", &model.Endpoint{ID: "ep-1", Name: "renamed", Address: addr})

	endpoints := cm.store.Config[0].Endpoints
	if !assert.Len(t, endpoints, 1,
		"Name-only change must NOT append a -2 sibling — that would grow the slice on harmless renames") {
		return
	}
	assert.Equal(t, "ep-1", endpoints[0].ID)
	assert.Equal(t, "original-name", endpoints[0].Name,
		"the SetEndpoint godoc commits to not applying Name-only updates; "+
			"if this assertion ever flips, update the SetEndpoint contract documentation too")
}

func TestClusterManager_AssembleEndpointsPreservesExplicitID(t *testing.T) {
	cluster := &model.ClusterConfig{
		Name:  "explicit-id",
		LbStr: model.LoadBalancerRoundRobin,
		Endpoints: []*model.Endpoint{
			{ID: "operator-pinned", Address: model.SocketAddress{Address: "127.0.0.1", Port: 21082}},
		},
	}

	cm := testClusterManager(cluster)
	defer stopStoreRuntimes(cm.store)

	if assert.Len(t, cm.store.Config[0].Endpoints, 1) {
		assert.Equal(t, "operator-pinned", cm.store.Config[0].Endpoints[0].ID)
	}
}

// TestClusterManager_AssembleEndpointsDeduplicatesDuplicateGenerated locks
// the PR-3 dedup behavior: when two static endpoints share every input that
// feeds GenerateEndpointID (cluster name, address, LLM credential), the
// second one's ID gets a -2 suffix so runtime health-key lookups stay
// unambiguous. A warning is logged so the operator notices.
func TestClusterManager_AssembleEndpointsDeduplicatesDuplicateGenerated(t *testing.T) {
	cluster := &model.ClusterConfig{
		Name:  "collapse-id",
		LbStr: model.LoadBalancerRoundRobin,
		Endpoints: []*model.Endpoint{
			{Address: model.SocketAddress{Address: "127.0.0.1", Port: 21083}},
			{Address: model.SocketAddress{Address: "127.0.0.1", Port: 21083}},
		},
	}

	cm := testClusterManager(cluster)
	defer stopStoreRuntimes(cm.store)

	endpoints := cm.store.Config[0].Endpoints
	if !assert.Len(t, endpoints, 2) {
		return
	}

	assert.NotEqual(t, endpoints[0].ID, endpoints[1].ID,
		"PR-3 deduplicates colliding generated IDs instead of collapsing them")
	assert.Contains(t, endpoints[0].ID, "pixiu-generated-endpoint-")
	assert.Equal(t, endpoints[0].ID+"-2", endpoints[1].ID,
		"second endpoint must be suffixed -2 relative to the first")
}

// TestAssembleEndpointsDeduplicatesExplicitID locks Blocker 6: when an
// operator writes the same `id:` on two static endpoints, the second one
// becomes `<id>-2` rather than `generated-<hash>-2`. Most-readable
// (least-surprising) choice for dashboards and log correlation.
func TestAssembleEndpointsDeduplicatesExplicitID(t *testing.T) {
	cluster := &model.ClusterConfig{
		Name:  "explicit-id-collision",
		LbStr: model.LoadBalancerRoundRobin,
		Endpoints: []*model.Endpoint{
			{ID: "foo", Address: model.SocketAddress{Address: "127.0.0.1", Port: 21090}},
			{ID: "foo", Address: model.SocketAddress{Address: "127.0.0.1", Port: 21091}},
		},
	}

	cm := testClusterManager(cluster)
	defer stopStoreRuntimes(cm.store)

	endpoints := cm.store.Config[0].Endpoints
	if !assert.Len(t, endpoints, 2) {
		return
	}

	assert.Equal(t, "foo", endpoints[0].ID, "first explicit ID is preserved")
	assert.Equal(t, "foo-2", endpoints[1].ID,
		"colliding second endpoint must be suffixed off the operator's ID, "+
			"not the generated- hash, so the operator's choice stays readable")
}

func testClusterManager(clusters ...*model.ClusterConfig) *ClusterManager {
	return CreateDefaultClusterManager(&model.Bootstrap{
		StaticResources: model.StaticResources{
			Clusters: clusters,
		},
	})
}

func testCluster(name string, lb model.LbPolicyType, endpoints []*model.Endpoint, healthChecks ...model.HealthCheckConfig) *model.ClusterConfig {
	cluster := &model.ClusterConfig{
		Name:      name,
		LbStr:     lb,
		Endpoints: endpoints,
	}
	if len(healthChecks) > 0 {
		cluster.HealthChecks = append([]model.HealthCheckConfig(nil), healthChecks...)
	}
	return cluster
}

func testEndpoint(id string, host string, port int) *model.Endpoint {
	return &model.Endpoint{
		ID:   id,
		Name: fmt.Sprintf("endpoint-%s", id),
		Address: model.SocketAddress{
			Address: host,
			Port:    port,
		},
	}
}

func testHealthCheck() model.HealthCheckConfig {
	return model.HealthCheckConfig{
		Protocol:       "tcp",
		TimeoutConfig:  "1h",
		IntervalConfig: "1h",
	}
}

func stopStoreRuntimes(store *ClusterStore) {
	if store == nil {
		return
	}
	for _, runtime := range store.clustersMap {
		if runtime != nil {
			runtime.Stop()
		}
	}
}

func healthCheckersLen(runtime *cluster.Cluster) int {
	if runtime == nil || runtime.HealthCheck == nil {
		return 0
	}
	return reflect.ValueOf(runtime.HealthCheck).Elem().FieldByName("checkers").Len()
}

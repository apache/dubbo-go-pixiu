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

func TestClusterManager_SetEndpointUpdateRebuildsConsistentHash(t *testing.T) {
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

			hash := cm.store.Config[0].ConsistentHash.Hash
			if !assert.NotNil(t, hash) {
				return
			}
			if hostList, ok := hash.(interface{ Hosts() []string }); ok {
				hosts := hostList.Hosts()
				assert.NotContains(t, hosts, oldHost)
				assert.Contains(t, hosts, newHost)
				return
			}
			assert.False(t, hash.Remove(oldHost))
			assert.True(t, hash.Remove(newHost))
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
		assert.Same(t, remainingEndpoint, config.Endpoints[0])
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

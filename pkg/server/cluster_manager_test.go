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
	"sync"
	"sync/atomic"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
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

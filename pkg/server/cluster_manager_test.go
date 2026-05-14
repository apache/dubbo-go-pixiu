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
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cluster"
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer"
	_ "github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/maglev"     // Register Maglev for cluster-manager tests.
	_ "github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/rand"       // Register Rand for cluster-manager tests.
	_ "github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/ringhash"   // Register RingHash for cluster-manager tests.
	_ "github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/roundrobin" // Register RoundRobin for cluster-manager tests.
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	testLegacyCompatibilityLockLB model.LbPolicyType = "test-legacy-compatibility-lock"
	testLegacyScopedLockLB        model.LbPolicyType = "test-legacy-scoped-lock"
	testLegacyHealthFilteringLB   model.LbPolicyType = "test-legacy-health-filtering"
	testSnapshotAllEndpointsLB    model.LbPolicyType = "test-snapshot-all-endpoints"
	testSnapshotMutatingLB        model.LbPolicyType = "test-snapshot-mutating"
)

type serverBlockingLegacyLoadBalancer struct {
	entered chan string
	release chan struct{}
}

type serverClusterScopedBlockingLegacyLoadBalancer struct {
	*serverBlockingLegacyLoadBalancer
}

type serverHealthFilteringLegacyLoadBalancer struct{}

type serverSnapshotAllEndpointsLoadBalancer struct {
	seenAll     []*model.Endpoint
	seenHealthy []*model.Endpoint
}

type serverSnapshotMutatingLoadBalancer struct{}

type serverLegacyPickHarness struct {
	t           *testing.T
	balancer    *serverBlockingLegacyLoadBalancer
	releaseOnce sync.Once
}

func (b *serverBlockingLegacyLoadBalancer) Handler(c *model.ClusterConfig, _ model.LbPolicy) *model.Endpoint {
	b.entered <- c.Name
	<-b.release
	if len(c.Endpoints) == 0 {
		return nil
	}
	return c.Endpoints[0]
}

func (b *serverClusterScopedBlockingLegacyLoadBalancer) UseClusterScopedLegacyLock() bool {
	return true
}

func (serverHealthFilteringLegacyLoadBalancer) Handler(c *model.ClusterConfig, _ model.LbPolicy) *model.Endpoint {
	endpoints := c.GetEndpoint(true)
	if len(endpoints) == 0 {
		return nil
	}
	return endpoints[0]
}

func (b *serverSnapshotAllEndpointsLoadBalancer) Handler(_ *model.ClusterConfig, _ model.LbPolicy) *model.Endpoint {
	return nil
}

func (b *serverSnapshotAllEndpointsLoadBalancer) HandlerWithSnapshot(c loadbalancer.PickContext, _ model.LbPolicy) *model.Endpoint {
	b.seenAll = c.AllEndpoints
	b.seenHealthy = c.HealthyEndpoints
	if len(c.HealthyEndpoints) == 0 {
		return nil
	}
	return c.HealthyEndpoints[0]
}

func (serverSnapshotMutatingLoadBalancer) Handler(_ *model.ClusterConfig, _ model.LbPolicy) *model.Endpoint {
	return nil
}

func (serverSnapshotMutatingLoadBalancer) HandlerWithSnapshot(c loadbalancer.PickContext, _ model.LbPolicy) *model.Endpoint {
	if len(c.HealthyEndpoints) > 0 {
		c.HealthyEndpoints[0].Metadata["weight"] = "99"
	}
	if len(c.AllEndpoints) > 1 {
		return c.AllEndpoints[1]
	}
	return nil
}

func registerServerLegacyBalancer(t *testing.T, policy model.LbPolicyType, scoped bool) *serverLegacyPickHarness {
	t.Helper()
	blocking := &serverBlockingLegacyLoadBalancer{
		entered: make(chan string, 2),
		release: make(chan struct{}),
	}
	var registered loadbalancer.LoadBalancer = blocking
	if scoped {
		registered = &serverClusterScopedBlockingLegacyLoadBalancer{
			serverBlockingLegacyLoadBalancer: blocking,
		}
	}

	previous, hadPrevious := loadbalancer.LoadBalancerStrategy[policy]
	loadbalancer.LoadBalancerStrategy[policy] = registered
	harness := &serverLegacyPickHarness{
		t:        t,
		balancer: blocking,
	}
	t.Cleanup(func() {
		harness.release()
		if hadPrevious {
			loadbalancer.LoadBalancerStrategy[policy] = previous
			return
		}
		delete(loadbalancer.LoadBalancerStrategy, policy)
	})
	return harness
}

func testLegacyLockCluster(name string, policy model.LbPolicyType, portBase int) *model.ClusterConfig {
	return testCluster(name, policy, []*model.Endpoint{
		testEndpoint(name+"-1", "127.0.0.1", portBase),
		testEndpoint(name+"-2", "127.0.0.1", portBase+1),
	})
}

func startServerPick(cm *ClusterManager, clusterName string) <-chan struct{} {
	done := make(chan struct{})
	go func() {
		defer close(done)
		_ = cm.PickEndpoint(clusterName, nil)
	}()
	return done
}

func (h *serverLegacyPickHarness) release() {
	h.releaseOnce.Do(func() {
		close(h.balancer.release)
	})
}

func (h *serverLegacyPickHarness) waitEntry() string {
	h.t.Helper()
	return waitServerLegacyHandlerEntry(h.t, h.balancer.entered)
}

func (h *serverLegacyPickHarness) assertNoEntry() {
	h.t.Helper()
	assertNoServerLegacyHandlerEntry(h.t, h.balancer.entered)
}

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

func TestClusterManagerRepairsDuplicateEndpointIDsBeforeRuntimeSnapshot(t *testing.T) {
	first := testEndpoint("duplicate", "127.0.0.1", 18082)
	second := testEndpoint("duplicate", "127.0.0.2", 18083)
	config := testCluster("duplicate-endpoint-id", model.LoadBalancerRoundRobin, []*model.Endpoint{first, second})
	cm := testClusterManager(config)

	if !assert.Len(t, config.Endpoints, 2) {
		return
	}
	assert.Equal(t, "duplicate", first.ID)
	assert.NotEmpty(t, second.ID)
	assert.NotEqual(t, first.ID, second.ID)

	runtimeCluster := cm.store.clustersMap[config.Name]
	if !assert.NotNil(t, runtimeCluster) {
		return
	}
	assert.True(t, runtimeCluster.UpdateEndpointHealth(first.ID, first.Address.GetAddress(), false))
	assert.True(t, runtimeCluster.UpdateEndpointHealth(second.ID, second.Address.GetAddress(), false))
	assert.Empty(t, runtimeCluster.EndpointSnapshot().HealthyEndpoints())
	assert.NotNil(t, runtimeCluster.EndpointSnapshot().EndpointByID(first.ID))
	assert.NotNil(t, runtimeCluster.EndpointSnapshot().EndpointByID(second.ID))
}

func TestClusterManager_SetEndpointFreshAnonymousEndpointUpdatesStableEndpoint(t *testing.T) {
	config := testCluster("set-anonymous-endpoint", model.LoadBalancerRoundRobin, nil)
	cm := testClusterManager(config)

	first := testEndpoint("", "127.0.0.1", 18084)
	first.Metadata = map[string]string{"version": "1"}
	cm.SetEndpoint(config.Name, first)
	firstID := first.ID
	assert.NotEmpty(t, firstID)

	second := testEndpoint("", "127.0.0.1", 18084)
	second.Metadata = map[string]string{"version": "2"}
	cm.SetEndpoint(config.Name, second)

	if assert.Len(t, config.Endpoints, 1) {
		assert.Equal(t, firstID, config.Endpoints[0].ID)
		assert.Equal(t, map[string]string{"version": "2"}, config.Endpoints[0].Metadata)
	}
	snapshot := cm.store.clustersMap[config.Name].EndpointSnapshot()
	assert.Len(t, snapshot.AllEndpoints(), 1)
	if picked := cm.PickEndpoint(config.Name, nil); assert.NotNil(t, picked) {
		assert.Equal(t, firstID, picked.ID)
		assert.Equal(t, map[string]string{"version": "2"}, picked.Metadata)
	}
}

func TestClusterManager_SetEndpointAnonymousIDAvoidsExistingExplicitGeneratedID(t *testing.T) {
	clusterName := "set-anonymous-id-collision"
	incoming := testEndpoint("", "127.0.0.2", 18085)
	collidingID := model.GeneratedEndpointID(clusterName, incoming)
	existing := testEndpoint(collidingID, "127.0.0.1", 18084)
	config := testCluster(clusterName, model.LoadBalancerRoundRobin, []*model.Endpoint{existing})
	cm := testClusterManager(config)

	cm.SetEndpoint(config.Name, incoming)

	if !assert.Len(t, config.Endpoints, 2) {
		return
	}
	assert.Equal(t, collidingID, config.Endpoints[0].ID)
	assert.Equal(t, model.SocketAddress{Address: "127.0.0.1", Port: 18084}, config.Endpoints[0].Address)
	assert.NotEmpty(t, incoming.ID)
	assert.NotEqual(t, collidingID, incoming.ID)
	assert.Equal(t, incoming.ID, config.Endpoints[1].ID)
	assert.Equal(t, model.SocketAddress{Address: "127.0.0.2", Port: 18085}, config.Endpoints[1].Address)
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

func TestClusterManager_PickNextEndpointSkipsSnapshotUnhealthyEndpoints(t *testing.T) {
	ep1 := testEndpoint("fallback-1", "127.0.0.1", 18087)
	ep2 := testEndpoint("fallback-2", "127.0.0.1", 18088)
	ep3 := testEndpoint("fallback-3", "127.0.0.1", 18089)
	cm := testClusterManager(
		testCluster("fallback-health", model.LoadBalancerRoundRobin, []*model.Endpoint{ep1, ep2, ep3}),
	)
	runtimeCluster := cm.store.clustersMap["fallback-health"]

	assert.True(t, runtimeCluster.UpdateEndpointHealth(ep2.ID, ep2.Address.GetAddress(), false))
	assert.False(t, ep2.UnHealthy)

	endpoint := cm.PickNextEndpoint("fallback-health", ep1.ID)
	if assert.NotNil(t, endpoint) {
		assert.Equal(t, ep3.ID, endpoint.ID)
	}

	assert.True(t, runtimeCluster.UpdateEndpointHealth(ep3.ID, ep3.Address.GetAddress(), false))
	assert.Nil(t, cm.PickNextEndpoint("fallback-health", ep1.ID))
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

func TestClusterManager_PickEndpointUsesHealthySnapshot(t *testing.T) {
	endpoint := testEndpoint("snapshot-ep", "127.0.0.1", 18088)
	cm := testClusterManager(
		testCluster("snapshot-pick", model.LoadBalancerRoundRobin, []*model.Endpoint{endpoint}),
	)
	runtimeCluster := cm.store.clustersMap["snapshot-pick"]

	assert.True(t, runtimeCluster.UpdateEndpointHealth(endpoint.ID, endpoint.Address.GetAddress(), false))
	assert.False(t, endpoint.UnHealthy)
	assert.Nil(t, cm.PickEndpoint("snapshot-pick", nil))

	assert.True(t, runtimeCluster.UpdateEndpointHealth(endpoint.ID, endpoint.Address.GetAddress(), true))
	if picked := cm.PickEndpoint("snapshot-pick", nil); assert.NotNil(t, picked) {
		assert.Equal(t, endpoint.ID, picked.ID)
	}
}

func TestClusterManager_PickEndpointSkipsSameAddressEndpointsAfterAddressHealthEvent(t *testing.T) {
	first := testEndpoint("snapshot-shared-address-1", "127.0.0.1", 18088)
	second := testEndpoint("snapshot-shared-address-2", "127.0.0.1", 18088)
	fallback := testEndpoint("snapshot-shared-address-fallback", "127.0.0.1", 18089)
	config := testCluster("snapshot-shared-address", model.LoadBalancerRoundRobin, []*model.Endpoint{
		first,
		second,
		fallback,
	})
	cm := testClusterManager(config)
	runtimeCluster := cm.store.clustersMap[config.Name]

	assert.True(t, runtimeCluster.UpdateEndpointAddressHealth(first.Address.GetAddress(), false))
	assert.Nil(t, cm.GetHealthyEndpointByID(config.Name, first.ID))
	assert.Nil(t, cm.GetHealthyEndpointByID(config.Name, second.ID))
	assert.NotNil(t, cm.GetHealthyEndpointByID(config.Name, fallback.ID))

	picked := cm.PickEndpoint(config.Name, nil)
	if assert.NotNil(t, picked) {
		assert.Equal(t, fallback.ID, picked.ID)
	}

	assert.True(t, runtimeCluster.UpdateEndpointAddressHealth(first.Address.GetAddress(), true))
	assert.NotNil(t, cm.GetHealthyEndpointByID(config.Name, first.ID))
	assert.NotNil(t, cm.GetHealthyEndpointByID(config.Name, second.ID))
}

func TestClusterManager_PickEndpointLegacyLBSeesRestoredRuntimeHealth(t *testing.T) {
	previous, hadPrevious := loadbalancer.LoadBalancerStrategy[testLegacyHealthFilteringLB]
	loadbalancer.LoadBalancerStrategy[testLegacyHealthFilteringLB] = serverHealthFilteringLegacyLoadBalancer{}
	t.Cleanup(func() {
		if hadPrevious {
			loadbalancer.LoadBalancerStrategy[testLegacyHealthFilteringLB] = previous
			return
		}
		delete(loadbalancer.LoadBalancerStrategy, testLegacyHealthFilteringLB)
	})

	restored := testEndpoint("snapshot-restored", "127.0.0.1", 18088)
	restored.UnHealthy = true
	fallback := testEndpoint("snapshot-fallback", "127.0.0.1", 18089)
	config := testCluster("snapshot-legacy-health", testLegacyHealthFilteringLB, []*model.Endpoint{restored, fallback})
	cm := testClusterManager(config)
	runtimeCluster := cm.store.clustersMap[config.Name]

	assert.True(t, runtimeCluster.UpdateEndpointHealth(restored.ID, restored.Address.GetAddress(), true))

	picked := cm.PickEndpoint(config.Name, nil)
	if assert.NotNil(t, picked) {
		assert.Equal(t, restored.ID, picked.ID)
		assert.False(t, picked.UnHealthy)
	}
	assert.True(t, restored.UnHealthy)
}

func TestClusterManager_PickEndpointSingleHealthyInMultiEndpointUsesLoadBalancer(t *testing.T) {
	ep1 := testEndpoint("snapshot-lb-1", "127.0.0.1", 18089)
	ep2 := testEndpoint("snapshot-lb-2", "127.0.0.1", 18090)
	config := testCluster("snapshot-lb-degraded", model.LoadBalancerRoundRobin, []*model.Endpoint{ep1, ep2})
	cm := testClusterManager(config)
	runtimeCluster := cm.store.clustersMap[config.Name]

	assert.True(t, runtimeCluster.UpdateEndpointHealth(ep1.ID, ep1.Address.GetAddress(), false))
	picked := cm.PickEndpoint(config.Name, nil)

	if assert.NotNil(t, picked) {
		assert.Equal(t, ep2.ID, picked.ID)
	}
	assert.Equal(t, uint32(1), atomic.LoadUint32(&config.PrePickEndpointIndex))
}

func TestClusterManager_SnapshotLoadBalancerReceivesAllEndpoints(t *testing.T) {
	balancer := &serverSnapshotAllEndpointsLoadBalancer{}
	previous, hadPrevious := loadbalancer.LoadBalancerStrategy[testSnapshotAllEndpointsLB]
	loadbalancer.LoadBalancerStrategy[testSnapshotAllEndpointsLB] = balancer
	t.Cleanup(func() {
		if hadPrevious {
			loadbalancer.LoadBalancerStrategy[testSnapshotAllEndpointsLB] = previous
			return
		}
		delete(loadbalancer.LoadBalancerStrategy, testSnapshotAllEndpointsLB)
	})

	healthy := testEndpoint("snapshot-all-healthy", "127.0.0.1", 18091)
	unhealthy := testEndpoint("snapshot-all-unhealthy", "127.0.0.1", 18092)
	config := testCluster("snapshot-all-endpoints", testSnapshotAllEndpointsLB, []*model.Endpoint{healthy, unhealthy})
	cm := testClusterManager(config)
	runtimeCluster := cm.store.clustersMap[config.Name]

	assert.True(t, runtimeCluster.UpdateEndpointHealth(unhealthy.ID, unhealthy.Address.GetAddress(), false))
	picked := cm.PickEndpoint(config.Name, nil)

	if assert.NotNil(t, picked) {
		assert.Equal(t, healthy.ID, picked.ID)
	}
	assert.Equal(t, []*model.Endpoint{healthy}, balancer.seenHealthy)
	if assert.Len(t, balancer.seenAll, 2) {
		assert.Equal(t, healthy.ID, balancer.seenAll[0].ID)
		assert.False(t, balancer.seenAll[0].UnHealthy)
		assert.Equal(t, unhealthy.ID, balancer.seenAll[1].ID)
		assert.True(t, balancer.seenAll[1].UnHealthy)
	}
}

func TestClusterManager_PickEndpointToleratesNilEntryInAllEndpoints(t *testing.T) {
	balancer := &serverSnapshotAllEndpointsLoadBalancer{}
	previous, hadPrevious := loadbalancer.LoadBalancerStrategy[testSnapshotAllEndpointsLB]
	loadbalancer.LoadBalancerStrategy[testSnapshotAllEndpointsLB] = balancer
	t.Cleanup(func() {
		if hadPrevious {
			loadbalancer.LoadBalancerStrategy[testSnapshotAllEndpointsLB] = previous
			return
		}
		delete(loadbalancer.LoadBalancerStrategy, testSnapshotAllEndpointsLB)
	})

	healthy := testEndpoint("snapshot-nil-tolerant-healthy", "127.0.0.1", 18099)
	config := testCluster("snapshot-nil-tolerant", testSnapshotAllEndpointsLB, []*model.Endpoint{healthy, nil})
	cm := testClusterManager(config)

	picked := cm.PickEndpoint(config.Name, nil)
	if assert.NotNil(t, picked) {
		assert.Equal(t, healthy.ID, picked.ID)
	}
	if assert.Len(t, balancer.seenAll, 2) {
		assert.NotNil(t, balancer.seenAll[0])
		assert.Nil(t, balancer.seenAll[1])
	}
}

func TestClusterManager_SnapshotLoadBalancerCannotMutateOrBypassHealthySnapshot(t *testing.T) {
	previous, hadPrevious := loadbalancer.LoadBalancerStrategy[testSnapshotMutatingLB]
	loadbalancer.LoadBalancerStrategy[testSnapshotMutatingLB] = serverSnapshotMutatingLoadBalancer{}
	t.Cleanup(func() {
		if hadPrevious {
			loadbalancer.LoadBalancerStrategy[testSnapshotMutatingLB] = previous
			return
		}
		delete(loadbalancer.LoadBalancerStrategy, testSnapshotMutatingLB)
	})

	healthy := testEndpoint("snapshot-safe-healthy", "127.0.0.1", 18093)
	healthy.Metadata = map[string]string{"weight": "1"}
	unhealthy := testEndpoint("snapshot-safe-unhealthy", "127.0.0.1", 18094)
	config := testCluster("snapshot-safe", testSnapshotMutatingLB, []*model.Endpoint{healthy, unhealthy})
	cm := testClusterManager(config)
	runtimeCluster := cm.store.clustersMap[config.Name]

	assert.True(t, runtimeCluster.UpdateEndpointHealth(unhealthy.ID, unhealthy.Address.GetAddress(), false))
	assert.Nil(t, cm.PickEndpoint(config.Name, nil))

	if got := runtimeCluster.EndpointSnapshot().HealthyEndpointByID(healthy.ID); assert.NotNil(t, got) {
		assert.Equal(t, map[string]string{"weight": "1"}, got.Metadata)
	}
	assert.Nil(t, runtimeCluster.EndpointSnapshot().HealthyEndpointByID(unhealthy.ID))
}

func TestClusterManager_GetEndpointByIDUsesHealthySnapshot(t *testing.T) {
	endpoint := testEndpoint("snapshot-id-ep", "127.0.0.1", 18089)
	cm := testClusterManager(
		testCluster("snapshot-id", model.LoadBalancerRoundRobin, []*model.Endpoint{endpoint}),
	)
	runtimeCluster := cm.store.clustersMap["snapshot-id"]

	assert.True(t, runtimeCluster.UpdateEndpointHealth(endpoint.ID, endpoint.Address.GetAddress(), false))
	assert.False(t, endpoint.UnHealthy)
	assert.Nil(t, cm.GetEndpointByID("snapshot-id", endpoint.ID))
	assert.Nil(t, cm.GetHealthyEndpointByID("snapshot-id", endpoint.ID))
	if got := cm.GetAnyEndpointByID("snapshot-id", endpoint.ID); assert.NotNil(t, got) {
		assert.Equal(t, endpoint.ID, got.ID)
		assert.True(t, got.UnHealthy)
	}
	if got := runtimeCluster.EndpointSnapshot().EndpointByID(endpoint.ID); assert.NotNil(t, got) {
		assert.Equal(t, endpoint.ID, got.ID)
	}

	assert.True(t, runtimeCluster.UpdateEndpointHealth(endpoint.ID, endpoint.Address.GetAddress(), true))
	if got := cm.GetEndpointByID("snapshot-id", endpoint.ID); assert.NotNil(t, got) {
		assert.Equal(t, endpoint.ID, got.ID)
	}
	if got := cm.GetHealthyEndpointByID("snapshot-id", endpoint.ID); assert.NotNil(t, got) {
		assert.Equal(t, endpoint.ID, got.ID)
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
	}, testHealthCheck())
	cm := testClusterManager(cluster)
	defer stopStoreRuntimes(cm.store)

	const expectedCursor uint32 = 5
	atomic.StoreUint32(&cm.store.Config[0].PrePickEndpointIndex, expectedCursor)
	oldRuntime := cm.store.clustersMap[cluster.Name]
	temporarilyUnhealthy := cluster.Endpoints[1]
	assert.True(t, oldRuntime.UpdateEndpointHealth(
		temporarilyUnhealthy.ID,
		temporarilyUnhealthy.Address.GetAddress(),
		false,
	))

	oldStore, err := cm.CloneStore()
	if !assert.NoError(t, err) {
		return
	}
	newStore := cm.NewStore(oldStore.Version)
	defer stopStoreRuntimes(newStore)
	newStore.AddCluster(testCluster(cluster.Name, model.LoadBalancerRoundRobin, nil, testHealthCheck()))
	for _, endpoint := range oldStore.Config[0].Endpoints {
		copied := *endpoint
		newStore.SetEndpoint(cluster.Name, &copied)
	}

	assert.True(t, cm.CompareAndSetStore(newStore))
	if assert.Len(t, cm.store.Config, 1) {
		assert.Equal(t, expectedCursor, atomic.LoadUint32(&cm.store.Config[0].PrePickEndpointIndex))
	}
	assert.Nil(t, cm.GetEndpointByID(cluster.Name, temporarilyUnhealthy.ID))
	assert.Nil(t, cm.GetHealthyEndpointByID(cluster.Name, temporarilyUnhealthy.ID))
	if got := cm.store.clustersMap[cluster.Name].EndpointSnapshot().EndpointByID(temporarilyUnhealthy.ID); assert.NotNil(t, got) {
		assert.Equal(t, temporarilyUnhealthy.ID, got.ID)
	}
	assert.False(t, oldRuntime.UpdateEndpointHealth(
		temporarilyUnhealthy.ID,
		temporarilyUnhealthy.Address.GetAddress(),
		true,
	))
	assert.Nil(t, cm.GetEndpointByID(cluster.Name, temporarilyUnhealthy.ID))
	assert.Nil(t, cm.GetHealthyEndpointByID(cluster.Name, temporarilyUnhealthy.ID))

	endpoint := cm.PickEndpoint(cluster.Name, nil)
	if assert.NotNil(t, endpoint) {
		assert.Equal(t, "ep-3", endpoint.ID)
	}
}

func TestClusterManager_CompareAndSetStorePreservesAnonymousEndpointHealthAcrossRefresh(t *testing.T) {
	initialEndpoint := testEndpoint("", "127.0.0.1", 19203)
	cluster := testCluster("refresh-anonymous-health", model.LoadBalancerRoundRobin, []*model.Endpoint{
		initialEndpoint,
	}, testHealthCheck())
	cm := testClusterManager(cluster)
	defer stopStoreRuntimes(cm.store)

	initialID := initialEndpoint.ID
	assert.NotEmpty(t, initialID)

	oldRuntime := cm.store.clustersMap[cluster.Name]
	assert.True(t, oldRuntime.UpdateEndpointHealth(
		initialID,
		initialEndpoint.Address.GetAddress(),
		false,
	))

	newStore := cm.NewStore(cm.store.Version)
	defer stopStoreRuntimes(newStore)
	refreshedEndpoint := testEndpoint("", "127.0.0.1", 19203)
	refreshedEndpoint.Name = "renamed-endpoint"
	newStore.AddCluster(testCluster(cluster.Name, model.LoadBalancerRoundRobin, []*model.Endpoint{
		refreshedEndpoint,
	}, testHealthCheck()))
	assert.Equal(t, initialID, refreshedEndpoint.ID)

	assert.True(t, cm.CompareAndSetStore(newStore))
	assert.Nil(t, cm.GetHealthyEndpointByID(cluster.Name, initialID))
	if got := cm.store.clustersMap[cluster.Name].EndpointSnapshot().EndpointByID(initialID); assert.NotNil(t, got) {
		assert.True(t, got.UnHealthy)
	}
}

func TestClusterManager_GeneratedEndpointIDsIncludeLLMIdentity(t *testing.T) {
	first := testLLMIdentityEndpoint("127.0.0.1", 19204, "openai", "key-a")
	second := testLLMIdentityEndpoint("127.0.0.1", 19204, "openai", "key-b")
	cluster := testCluster("refresh-llm-identity", model.LoadBalancerRoundRobin, []*model.Endpoint{
		first,
		second,
	}, testHealthCheck())
	cm := testClusterManager(cluster)
	defer stopStoreRuntimes(cm.store)

	firstID := first.ID
	secondID := second.ID
	assert.NotEmpty(t, firstID)
	assert.NotEmpty(t, secondID)
	assert.NotEqual(t, firstID, secondID)

	oldRuntime := cm.store.clustersMap[cluster.Name]
	assert.True(t, oldRuntime.UpdateEndpointHealth(secondID, second.Address.GetAddress(), false))

	newStore := cm.NewStore(cm.store.Version)
	defer stopStoreRuntimes(newStore)
	refreshedSecond := testLLMIdentityEndpoint("127.0.0.1", 19204, "openai", "key-b")
	refreshedFirst := testLLMIdentityEndpoint("127.0.0.1", 19204, "openai", "key-a")
	newStore.AddCluster(testCluster(cluster.Name, model.LoadBalancerRoundRobin, []*model.Endpoint{
		refreshedSecond,
		refreshedFirst,
	}, testHealthCheck()))

	assert.Equal(t, secondID, refreshedSecond.ID)
	assert.Equal(t, firstID, refreshedFirst.ID)
	assert.True(t, cm.CompareAndSetStore(newStore))
	assert.Nil(t, cm.GetHealthyEndpointByID(cluster.Name, secondID))
	assert.NotNil(t, cm.GetHealthyEndpointByID(cluster.Name, firstID))
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

			assertRuntimeEndpointAddress(t, cm.store.clustersMap[config.Name], newEndpoint)
			assertConsistentHashRebuilt(t, cm.store.Config[0].ConsistentHash.Hash, oldHost, newHost)
		})
	}
}

func assertRuntimeEndpointAddress(t *testing.T, runtime *cluster.Cluster, endpoint *model.Endpoint) {
	t.Helper()
	if !assert.NotNil(t, runtime) {
		return
	}
	snapshot := runtime.EndpointSnapshot()
	runtimeEndpoint := snapshot.EndpointByID(endpoint.ID)
	if assert.NotNil(t, runtimeEndpoint) {
		assert.Equal(t, endpoint.Address, runtimeEndpoint.Address)
	}
	healthyEndpoint := snapshot.HealthyEndpointByID(endpoint.ID)
	if assert.NotNil(t, healthyEndpoint) {
		assert.Equal(t, endpoint.Address, healthyEndpoint.Address)
	}
}

func assertConsistentHashRebuilt(t *testing.T, hash model.LbConsistentHash, oldHost, newHost string) {
	t.Helper()
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
}

func TestClusterManager_SetEndpointUpdateMergesWithoutMutatingOldEndpoint(t *testing.T) {
	oldEndpoint := testEndpoint("ep-1", "127.0.0.1", 19345)
	oldEndpoint.Metadata = map[string]string{"stable": "old"}
	oldEndpoint.LLMMeta = &model.LLMMeta{
		Provider: "openai",
		APIKey:   "old-key",
	}
	oldEndpoint.UnHealthy = true
	config := testCluster("endpoint-merge", model.LoadBalancerRoundRobin, []*model.Endpoint{oldEndpoint})
	cm := testClusterManager(config)

	incoming := testEndpoint("ep-1", "127.0.0.2", 19346)
	incoming.Metadata = map[string]string{"discovered": "new"}
	cm.SetEndpoint(config.Name, incoming)

	if !assert.Len(t, config.Endpoints, 1) {
		return
	}
	merged := config.Endpoints[0]
	assert.NotSame(t, oldEndpoint, merged)
	assert.NotSame(t, incoming, merged)
	assert.Equal(t, incoming.ID, merged.ID)
	assert.Equal(t, incoming.Name, merged.Name)
	assert.Equal(t, incoming.Address, merged.Address)
	assert.Equal(t, incoming.Metadata, merged.Metadata)
	assert.Same(t, oldEndpoint.LLMMeta, merged.LLMMeta)
	assert.True(t, merged.UnHealthy)

	assert.Equal(t, "endpoint-ep-1", oldEndpoint.Name)
	assert.Equal(t, model.SocketAddress{Address: "127.0.0.1", Port: 19345}, oldEndpoint.Address)
	assert.Equal(t, map[string]string{"stable": "old"}, oldEndpoint.Metadata)

	runtimeEndpoint := cm.store.clustersMap[config.Name].EndpointSnapshot().EndpointByID(incoming.ID)
	assert.NotSame(t, merged, runtimeEndpoint)
	assert.Equal(t, merged, runtimeEndpoint)
	assert.Nil(t, cm.store.clustersMap[config.Name].EndpointSnapshot().HealthyEndpointByID(incoming.ID))
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
	assert.Equal(t, []*model.Endpoint{remainingEndpoint}, runtime.EndpointSnapshot().AllEndpoints())
	assert.Equal(t, []*model.Endpoint{remainingEndpoint}, runtime.EndpointSnapshot().HealthyEndpoints())

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

func TestClusterManager_Race_PickEndpointWithHealthUpdates(t *testing.T) {
	cluster := testCluster("race-health-snapshot", model.LoadBalancerRoundRobin, []*model.Endpoint{
		testEndpoint("ep-1", "127.0.0.1", 19400),
		testEndpoint("ep-2", "127.0.0.1", 19401),
		testEndpoint("ep-3", "127.0.0.1", 19402),
		testEndpoint("ep-4", "127.0.0.1", 19403),
	})
	cm := testClusterManager(cluster)
	runtimeCluster := cm.store.clustersMap[cluster.Name]

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

	wg.Add(1)
	go func() {
		defer wg.Done()
		<-start
		target := cluster.Endpoints[0]
		address := target.Address.GetAddress()
		for j := 0; j < 4000; j++ {
			_ = runtimeCluster.UpdateEndpointHealth(target.ID, address, j%2 == 0)
		}
	}()

	close(start)
	wg.Wait()
}

func TestClusterManager_LegacyLoadBalancerPicksUseCompatibilityLockAcrossClusters(t *testing.T) {
	harness := registerServerLegacyBalancer(t, testLegacyCompatibilityLockLB, false)
	cm := testClusterManager(
		testLegacyLockCluster("legacy-cluster-a", testLegacyCompatibilityLockLB, 19500),
		testLegacyLockCluster("legacy-cluster-b", testLegacyCompatibilityLockLB, 19502),
	)

	firstDone := startServerPick(cm, "legacy-cluster-a")
	assert.Equal(t, "legacy-cluster-a", harness.waitEntry())

	secondDone := startServerPick(cm, "legacy-cluster-b")
	harness.assertNoEntry()

	harness.release()
	assert.Equal(t, "legacy-cluster-b", harness.waitEntry())
	waitServerClosed(t, firstDone)
	waitServerClosed(t, secondDone)
}

func TestClusterManager_OptInLegacyLoadBalancerPicksUseRuntimeScopedLocksAcrossClusters(t *testing.T) {
	harness := registerServerLegacyBalancer(t, testLegacyScopedLockLB, true)
	cm := testClusterManager(
		testLegacyLockCluster("legacy-cluster-a", testLegacyScopedLockLB, 19500),
		testLegacyLockCluster("legacy-cluster-b", testLegacyScopedLockLB, 19502),
	)

	firstDone := startServerPick(cm, "legacy-cluster-a")
	assert.Equal(t, "legacy-cluster-a", harness.waitEntry())

	secondDone := startServerPick(cm, "legacy-cluster-b")
	assert.Equal(t, "legacy-cluster-b", harness.waitEntry())

	harness.release()
	waitServerClosed(t, firstDone)
	waitServerClosed(t, secondDone)
}

func TestClusterManager_OptInLegacyLoadBalancerPicksSerializeWithinSameCluster(t *testing.T) {
	harness := registerServerLegacyBalancer(t, testLegacyScopedLockLB, true)
	cm := testClusterManager(
		testLegacyLockCluster("legacy-cluster-a", testLegacyScopedLockLB, 19500),
	)

	firstDone := startServerPick(cm, "legacy-cluster-a")
	assert.Equal(t, "legacy-cluster-a", harness.waitEntry())

	secondDone := startServerPick(cm, "legacy-cluster-a")
	harness.assertNoEntry()

	harness.release()
	assert.Equal(t, "legacy-cluster-a", harness.waitEntry())
	waitServerClosed(t, firstDone)
	waitServerClosed(t, secondDone)
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

func testLLMIdentityEndpoint(host string, port int, provider, apiKey string) *model.Endpoint {
	endpoint := testEndpoint("", host, port)
	endpoint.Name = "shared-llm"
	endpoint.LLMMeta = &model.LLMMeta{
		Provider: provider,
		APIKey:   apiKey,
	}
	return endpoint
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

func waitServerLegacyHandlerEntry(t *testing.T, entered <-chan string) string {
	t.Helper()
	select {
	case clusterName := <-entered:
		return clusterName
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for legacy handler entry")
		return ""
	}
}

func assertNoServerLegacyHandlerEntry(t *testing.T, entered <-chan string) {
	t.Helper()
	select {
	case clusterName := <-entered:
		t.Fatalf("legacy handler for %s entered before the first cluster pick returned", clusterName)
	case <-time.After(50 * time.Millisecond):
	}
}

func waitServerClosed(t *testing.T, done <-chan struct{}) {
	t.Helper()
	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("timed out waiting for legacy pick to finish")
	}
}

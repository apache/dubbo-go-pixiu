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
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
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

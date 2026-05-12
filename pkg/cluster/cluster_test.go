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
	assert.Same(t, snapshot.EndpointByID(first.ID), snapshot.HealthyEndpointByID(first.ID))
}

func TestClusterEndpointSnapshotEndpointCountIsNilSafe(t *testing.T) {
	var snapshot *EndpointSnapshot

	assert.Zero(t, snapshot.EndpointCount())
}

func TestClusterEndpointSnapshotClonesConfigEndpointObjects(t *testing.T) {
	endpoint := testEndpoint("ep-1", "127.0.0.1", 18080)
	endpoint.Address.Domains = []string{"api.example.com"}
	endpoint.Metadata = map[string]string{"weight": "3"}
	endpoint.LLMMeta = &model.LLMMeta{
		Provider: "openai",
		APIKey:   "old-key",
		RetryPolicy: model.RetryPolicy{
			Config: map[string]any{"attempts": 1},
		},
	}

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
	endpoint.UnHealthy = true

	assert.Equal(t, model.SocketAddress{Address: "127.0.0.1", Port: 18080, Domains: []string{"api.example.com"}}, snapshotEndpoint.Address)
	assert.Equal(t, "api.example.com", snapshotEndpoint.Address.GetAddress())
	assert.Equal(t, map[string]string{"weight": "3"}, snapshotEndpoint.Metadata)
	assert.Equal(t, "old-key", snapshotEndpoint.LLMMeta.APIKey)
	assert.Equal(t, 1, snapshotEndpoint.LLMMeta.RetryPolicy.Config["attempts"])
	assert.False(t, snapshotEndpoint.UnHealthy)
}

func TestClusterEndpointHealthEventUpdatesSnapshotWithoutMutatingEndpoint(t *testing.T) {
	endpoint := testEndpoint("ep-1", "127.0.0.1", 18082)
	runtimeCluster := NewCluster(testCluster("snapshot-health", endpoint))

	assert.True(t, runtimeCluster.UpdateEndpointHealth(endpoint.ID, endpoint.Address.GetAddress(), false))

	assert.False(t, endpoint.UnHealthy)
	assert.False(t, runtimeCluster.EndpointSnapshot().EndpointByID(endpoint.ID).UnHealthy)
	assert.Empty(t, runtimeCluster.EndpointSnapshot().HealthyEndpoints())
	assert.Nil(t, runtimeCluster.EndpointSnapshot().HealthyEndpointByID(endpoint.ID))

	assert.True(t, runtimeCluster.UpdateEndpointHealth(endpoint.ID, endpoint.Address.GetAddress(), true))
	assert.False(t, endpoint.UnHealthy)
	assert.Equal(t, []*model.Endpoint{endpoint}, runtimeCluster.EndpointSnapshot().HealthyEndpoints())
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

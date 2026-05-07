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
	assert.Same(t, healthy, snapshot.EndpointByID(healthy.ID))
	assert.Same(t, healthy, snapshot.HealthyEndpointByID(healthy.ID))
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
	assert.Same(t, first, snapshot.EndpointByID(first.ID))
	assert.Same(t, first, snapshot.HealthyEndpointByID(first.ID))
}

func TestClusterEndpointSnapshotEndpointCountIsNilSafe(t *testing.T) {
	var snapshot *EndpointSnapshot

	assert.Zero(t, snapshot.EndpointCount())
}

func TestClusterEndpointHealthEventUpdatesSnapshotWithoutMutatingEndpoint(t *testing.T) {
	endpoint := testEndpoint("ep-1", "127.0.0.1", 18082)
	runtimeCluster := NewCluster(testCluster("snapshot-health", endpoint))

	assert.True(t, runtimeCluster.UpdateEndpointHealth(endpoint.ID, endpoint.Address.GetAddress(), false))

	assert.False(t, endpoint.UnHealthy)
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
	assert.Same(t, replacement, replacementRuntime.EndpointSnapshot().HealthyEndpointByID(replacement.ID))
}

func TestNewClusterWithEndpointSnapshotInheritsHealthOnlyWhenHealthCheckEnabled(t *testing.T) {
	endpoint := testEndpoint("ep-1", "127.0.0.1", 18085)
	oldRuntime := NewCluster(testCluster("snapshot-constructor-old", endpoint))
	state := oldRuntime.EndpointRuntimeState(endpoint.ID, endpoint.Address.GetAddress())
	if !assert.NotNil(t, state) {
		return
	}
	state.Store("cooldown", "true")
	assert.True(t, oldRuntime.UpdateEndpointHealth(endpoint.ID, endpoint.Address.GetAddress(), false))
	previous := oldRuntime.EndpointSnapshot()

	replacement := testEndpoint(endpoint.ID, "127.0.0.1", 18085)
	replacementRuntime := NewClusterWithEndpointSnapshot(
		testClusterWithHealthCheck("snapshot-constructor-same", replacement),
		previous,
	)
	t.Cleanup(replacementRuntime.Stop)
	assert.Nil(t, replacementRuntime.EndpointSnapshot().HealthyEndpointByID(replacement.ID))
	if carriedState := replacementRuntime.EndpointRuntimeState(replacement.ID, replacement.Address.GetAddress()); assert.Same(t, state, carriedState) {
		value, ok := carriedState.Load("cooldown")
		assert.True(t, ok)
		assert.Equal(t, "true", value)
	}

	noHealthCheckReplacement := testEndpoint(endpoint.ID, "127.0.0.1", 18085)
	noHealthCheckRuntime := NewClusterWithEndpointSnapshot(
		testCluster("snapshot-constructor-no-healthcheck", noHealthCheckReplacement),
		previous,
	)
	assert.Same(t, noHealthCheckReplacement, noHealthCheckRuntime.EndpointSnapshot().HealthyEndpointByID(noHealthCheckReplacement.ID))
	if carriedState := noHealthCheckRuntime.EndpointRuntimeState(noHealthCheckReplacement.ID, noHealthCheckReplacement.Address.GetAddress()); assert.Same(t, state, carriedState) {
		value, ok := carriedState.Load("cooldown")
		assert.True(t, ok)
		assert.Equal(t, "true", value)
	}

	moved := testEndpoint(endpoint.ID, "127.0.0.2", 18086)
	movedRuntime := NewClusterWithEndpointSnapshot(
		testCluster("snapshot-constructor-moved", moved),
		previous,
	)
	assert.Same(t, moved, movedRuntime.EndpointSnapshot().HealthyEndpointByID(moved.ID))
	movedState := movedRuntime.EndpointRuntimeState(moved.ID, moved.Address.GetAddress())
	if assert.NotNil(t, movedState) {
		_, ok := movedState.Load("cooldown")
		assert.False(t, ok)
	}
}

func TestClusterEndpointRuntimeStateFollowsSameAddressOnly(t *testing.T) {
	endpoint := testEndpoint("ep-1", "127.0.0.1", 18085)
	config := testCluster("runtime-state", endpoint)
	runtimeCluster := NewCluster(config)
	address := endpoint.Address.GetAddress()

	state := runtimeCluster.EndpointRuntimeState(endpoint.ID, address)
	if !assert.NotNil(t, state) {
		return
	}
	state.Store("cooldown", "true")

	replacement := testEndpoint(endpoint.ID, "127.0.0.1", 18085)
	config.Endpoints[0] = replacement
	runtimeCluster.RefreshEndpoints()

	carriedState := runtimeCluster.EndpointRuntimeState(replacement.ID, replacement.Address.GetAddress())
	if assert.Same(t, state, carriedState) {
		value, ok := carriedState.Load("cooldown")
		assert.True(t, ok)
		assert.Equal(t, "true", value)
	}

	moved := testEndpoint(endpoint.ID, "127.0.0.2", 18086)
	config.Endpoints[0] = moved
	runtimeCluster.RefreshEndpoints()

	assert.Nil(t, runtimeCluster.EndpointRuntimeState(endpoint.ID, address))
	movedState := runtimeCluster.EndpointRuntimeState(moved.ID, moved.Address.GetAddress())
	if assert.NotNil(t, movedState) {
		_, ok := movedState.Load("cooldown")
		assert.False(t, ok)
	}

	config.Endpoints = nil
	runtimeCluster.RefreshEndpoints()
	assert.Nil(t, runtimeCluster.EndpointRuntimeState(moved.ID, moved.Address.GetAddress()))
}

func TestEndpointRuntimeStateLoadManyReturnsCopy(t *testing.T) {
	state := newEndpointRuntimeState()
	state.StoreMany(map[string]string{
		"cooldown": "true",
		"checked":  "now",
	})

	values := state.LoadMany("cooldown", "checked", "missing")
	assert.Equal(t, map[string]string{
		"cooldown": "true",
		"checked":  "now",
	}, values)

	values["cooldown"] = "false"
	value, ok := state.Load("cooldown")
	assert.True(t, ok)
	assert.Equal(t, "true", value)
}

func TestEndpointRuntimeStateStoreManyStoresPair(t *testing.T) {
	state := newEndpointRuntimeState()

	state.StoreMany(map[string]string{
		"unhealthy": "true",
		"checked":   "now",
	})

	values := state.LoadMany("unhealthy", "checked")
	assert.Equal(t, "true", values["unhealthy"])
	assert.Equal(t, "now", values["checked"])
}

func TestEndpointRuntimeStateDeleteIfMatchesDeletesMatchingValues(t *testing.T) {
	state := newEndpointRuntimeState()
	state.StoreMany(map[string]string{
		"unhealthy": "true",
		"checked":   "old",
		"stable":    "keep",
	})

	assert.True(t, state.DeleteIfMatches(
		map[string]string{
			"unhealthy": "true",
			"checked":   "old",
		},
		"unhealthy",
		"checked",
	))

	_, ok := state.Load("unhealthy")
	assert.False(t, ok)
	_, ok = state.Load("checked")
	assert.False(t, ok)
	stable, ok := state.Load("stable")
	assert.True(t, ok)
	assert.Equal(t, "keep", stable)
}

func TestEndpointRuntimeStateDeleteIfMatchesPreservesRefreshedValues(t *testing.T) {
	state := newEndpointRuntimeState()
	state.StoreMany(map[string]string{
		"unhealthy": "true",
		"checked":   "old",
	})
	expected := state.LoadMany("unhealthy", "checked")

	state.Store("checked", "new")

	assert.False(t, state.DeleteIfMatches(expected, "unhealthy", "checked"))
	unhealthy, ok := state.Load("unhealthy")
	assert.True(t, ok)
	assert.Equal(t, "true", unhealthy)
	checked, ok := state.Load("checked")
	assert.True(t, ok)
	assert.Equal(t, "new", checked)
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

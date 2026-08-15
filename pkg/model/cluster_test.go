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

package model_test

import (
	"encoding/json"
	"strings"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	_ "github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/ringhash" // Register RingHash for consistent-hash tests.
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestClusterConfig_GetEndpointFiltersByHealth(t *testing.T) {
	cluster := &model.ClusterConfig{
		Endpoints: []*model.Endpoint{
			{ID: "ep-1"},
			{ID: "ep-2", UnHealthy: true},
			{ID: "ep-3"},
		},
	}

	all := cluster.GetEndpoint(false)
	healthy := cluster.GetEndpoint(true)

	require.Len(t, all, 3)
	require.Len(t, healthy, 2)
	assert.Equal(t, []string{"ep-1", "ep-3"}, []string{healthy[0].ID, healthy[1].ID})
}

func TestClusterConfig_CreateConsistentHashRegistersHash(t *testing.T) {
	cluster := &model.ClusterConfig{
		LbStr: model.LoadBalancerRingHashing,
		ConsistentHash: model.ConsistentHash{
			ReplicaNum:  32,
			MaxVnodeNum: 1023,
		},
		Endpoints: []*model.Endpoint{
			{
				ID: "ep-1",
				Address: model.SocketAddress{
					Address: "127.0.0.1",
					Port:    20880,
				},
			},
		},
	}

	cluster.CreateConsistentHash()

	require.NotNil(t, cluster.ConsistentHash.Hash)
	hash, err := cluster.ConsistentHash.Hash.GetHash(cluster.ConsistentHash.Hash.Hash("coverage-key"))
	require.NoError(t, err)
	assert.Equal(t, "127.0.0.1:20880", hash)
}

func TestClusterConfig_HasConsistentHashFactory(t *testing.T) {
	assert.True(t, (&model.ClusterConfig{LbStr: model.LoadBalancerRingHashing}).HasConsistentHashFactory())
	assert.False(t, (&model.ClusterConfig{LbStr: model.LbPolicyType("CustomHash")}).HasConsistentHashFactory())
}

func TestClusterConfig_EnsureConsistentHashBuildsOnce(t *testing.T) {
	cluster := &model.ClusterConfig{
		LbStr: model.LoadBalancerRingHashing,
		ConsistentHash: model.ConsistentHash{
			ReplicaNum:  32,
			MaxVnodeNum: 1023,
		},
		Endpoints: []*model.Endpoint{
			{
				ID: "ep-1",
				Address: model.SocketAddress{
					Address: "127.0.0.1",
					Port:    20880,
				},
			},
		},
	}

	cluster.EnsureConsistentHash()
	first := cluster.ConsistentHash.Hash
	cluster.EnsureConsistentHash()

	require.NotNil(t, first)
	assert.Same(t, first, cluster.ConsistentHash.Hash)
}

func TestCloneClusterConfigHandlesConsistentHashByFactoryAvailability(t *testing.T) {
	customHash := &testConsistentHash{}

	registered := &model.ClusterConfig{
		LbStr:          model.LoadBalancerRingHashing,
		ConsistentHash: model.ConsistentHash{Hash: customHash},
	}
	registeredClone := model.CloneClusterConfig(registered)
	assert.Nil(t, registeredClone.ConsistentHash.Hash,
		"a registered factory can rebuild the hash for the runtime snapshot")
	assert.Same(t, customHash, registered.ConsistentHash.Hash,
		"cloning must not mutate the source config")

	unregistered := &model.ClusterConfig{
		LbStr:          model.LbPolicyType("ProgrammaticCustomHash"),
		ConsistentHash: model.ConsistentHash{Hash: customHash},
	}
	unregisteredClone := model.CloneClusterConfig(unregistered)
	assert.Same(t, customHash, unregisteredClone.ConsistentHash.Hash,
		"a programmatic hash must survive when no factory can rebuild it")
}

type testConsistentHash struct{}

func (*testConsistentHash) Hash(string) uint32             { return 0 }
func (*testConsistentHash) Get(string) (string, error)     { return "", nil }
func (*testConsistentHash) GetHash(uint32) (string, error) { return "", nil }
func (*testConsistentHash) Add(string)                     {}
func (*testConsistentHash) Remove(string) bool             { return false }

func TestClusterConfig_PrePickEndpointIndexIsRuntimeOnly(t *testing.T) {
	cluster := &model.ClusterConfig{
		Name:                 "runtime-cursor",
		PrePickEndpointIndex: 7,
	}

	jsonBytes, err := json.Marshal(cluster)
	require.NoError(t, err)
	assert.NotContains(t, strings.ToLower(string(jsonBytes)), "prepickendpointindex")

	yamlBytes, err := yaml.MarshalYML(cluster)
	require.NoError(t, err)
	assert.NotContains(t, strings.ToLower(string(yamlBytes)), "prepickendpointindex")
}

func TestEndpoint_GetHost(t *testing.T) {
	endpoint := model.Endpoint{
		Address: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    20880,
		},
	}

	assert.Equal(t, "127.0.0.1:20880", endpoint.GetHost())
}

func TestGenerateEndpointIDIsDeterministic(t *testing.T) {
	build := func(apiKey string) *model.Endpoint {
		return &model.Endpoint{
			Name:    "shared-llm",
			Address: model.SocketAddress{Address: "127.0.0.1", Port: 20880},
			LLMMeta: &model.LLMMeta{Provider: "openai", APIKey: apiKey},
		}
	}

	first := model.GenerateEndpointID("cluster-a", build("key-a"))
	second := model.GenerateEndpointID("cluster-a", build("key-a"))
	other := model.GenerateEndpointID("cluster-a", build("key-b"))

	assert.Equal(t, first, second, "same material must produce same ID")
	assert.NotEqual(t, first, other, "different credential must produce different ID")
	assert.True(t, strings.HasPrefix(first, "pixiu-generated-endpoint-"),
		"ID must use pixiu-generated-endpoint- prefix")
	assert.NotContains(t, first, "key-a", "raw API key must not leak into the ID")
}

func TestGenerateEndpointIDIgnoresEndpointName(t *testing.T) {
	base := &model.Endpoint{
		Name:    "old-name",
		Address: model.SocketAddress{Address: "127.0.0.1", Port: 20880},
		LLMMeta: &model.LLMMeta{Provider: "openai", APIKey: "k"},
	}
	renamed := *base
	renamed.Name = "new-name"

	assert.Equal(t,
		model.GenerateEndpointID("cluster-a", base),
		model.GenerateEndpointID("cluster-a", &renamed),
	)
}

func TestGenerateEndpointIDSeparatesClusters(t *testing.T) {
	endpoint := &model.Endpoint{
		Address: model.SocketAddress{Address: "127.0.0.1", Port: 20880},
		LLMMeta: &model.LLMMeta{Provider: "openai", APIKey: "k"},
	}

	a := model.GenerateEndpointID("cluster-a", endpoint)
	b := model.GenerateEndpointID("cluster-b", endpoint)

	assert.NotEqual(t, a, b, "different cluster names must produce different IDs for the same endpoint")
}

func TestGenerateEndpointIDSeparatesAddresses(t *testing.T) {
	build := func(host string, port int) *model.Endpoint {
		return &model.Endpoint{
			Address: model.SocketAddress{Address: host, Port: port},
			LLMMeta: &model.LLMMeta{Provider: "openai", APIKey: "k"},
		}
	}

	base := model.GenerateEndpointID("cluster-a", build("127.0.0.1", 20880))

	assert.NotEqual(t, base,
		model.GenerateEndpointID("cluster-a", build("127.0.0.2", 20880)),
		"different host must produce different ID")
	assert.NotEqual(t, base,
		model.GenerateEndpointID("cluster-a", build("127.0.0.1", 20881)),
		"different port must produce different ID")
}

func TestGenerateEndpointIDHandlesNilEndpoint(t *testing.T) {
	a := model.GenerateEndpointID("cluster-a", nil)
	b := model.GenerateEndpointID("cluster-b", nil)
	assert.True(t, strings.HasPrefix(a, "pixiu-generated-endpoint-"))
	assert.NotEqual(t, a, b)
}

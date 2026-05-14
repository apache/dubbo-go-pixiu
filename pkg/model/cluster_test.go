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

func TestGeneratedEndpointIDIsDeterministic(t *testing.T) {
	build := func(apiKey string) *model.Endpoint {
		return &model.Endpoint{
			Name:    "shared-llm",
			Address: model.SocketAddress{Address: "127.0.0.1", Port: 20880},
			LLMMeta: &model.LLMMeta{Provider: "openai", APIKey: apiKey},
		}
	}

	first := model.GeneratedEndpointID("cluster-a", build("key-a"))
	second := model.GeneratedEndpointID("cluster-a", build("key-a"))
	other := model.GeneratedEndpointID("cluster-a", build("key-b"))

	assert.Equal(t, first, second, "same material must produce same ID")
	assert.NotEqual(t, first, other, "different credential must produce different ID")
	assert.True(t, strings.HasPrefix(first, "generated-"), "ID must use generated- prefix")
	assert.NotContains(t, first, "key-a", "raw API key must not leak into the ID")
}

func TestGeneratedEndpointIDIgnoresEndpointName(t *testing.T) {
	base := &model.Endpoint{
		Name:    "old-name",
		Address: model.SocketAddress{Address: "127.0.0.1", Port: 20880},
		LLMMeta: &model.LLMMeta{Provider: "openai", APIKey: "k"},
	}
	renamed := *base
	renamed.Name = "new-name"

	assert.Equal(t,
		model.GeneratedEndpointID("cluster-a", base),
		model.GeneratedEndpointID("cluster-a", &renamed),
	)
}

func TestGeneratedEndpointIDHandlesNilEndpoint(t *testing.T) {
	a := model.GeneratedEndpointID("cluster-a", nil)
	b := model.GeneratedEndpointID("cluster-b", nil)
	assert.True(t, strings.HasPrefix(a, "generated-"))
	assert.NotEqual(t, a, b)
}

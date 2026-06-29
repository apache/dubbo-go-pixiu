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

package mcpserver

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	filtermcp "github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type endpointOp struct {
	action    string
	cluster   string
	id        string
	address   string
	toolCount int
}

type recordingEndpointSink struct {
	ops []endpointOp
}

func (s *recordingEndpointSink) SetEndpoint(clusterName string, endpoint *model.Endpoint) {
	s.ops = append(s.ops, endpointOp{
		action:  "set",
		cluster: clusterName,
		id:      endpoint.ID,
		address: endpoint.Address.GetAddress(),
	})
}

func (s *recordingEndpointSink) DeleteEndpoint(clusterName, endpointID string) {
	s.ops = append(s.ops, endpointOp{
		action:  "delete",
		cluster: clusterName,
		id:      endpointID,
	})
}

func mcpConfigWithEndpointTools(tools ...model.ToolConfig) *model.McpServerConfig {
	return &model.McpServerConfig{Tools: tools}
}

func endpointTool(name, cluster, backendURL string) model.ToolConfig {
	return model.ToolConfig{Name: name, Cluster: cluster, BackendURL: backendURL}
}

func TestEndpointReconcilerBackendURLChangeReplacesStableOwnerEndpoint(t *testing.T) {
	sink := &recordingEndpointSink{}
	r := newEndpointReconciler(sink)
	source := filtermcp.NewServerSource("nacos", "server-a")

	err := r.ApplyServerConfig(source, mcpConfigWithEndpointTools(
		endpointTool("tool", "cluster-a", "http://127.0.0.1:8080"),
	))
	require.NoError(t, err)
	require.Len(t, sink.ops, 1)
	firstID := sink.ops[0].id

	err = r.ApplyServerConfig(source, mcpConfigWithEndpointTools(
		endpointTool("tool", "cluster-a", "http://127.0.0.1:9090"),
	))
	require.NoError(t, err)

	require.Len(t, sink.ops, 2)
	assert.Equal(t, "set", sink.ops[1].action)
	assert.Equal(t, firstID, sink.ops[1].id)
	assert.Equal(t, "127.0.0.1:9090", sink.ops[1].address)
}

func TestEndpointReconcilerToolAndServerDeletion(t *testing.T) {
	sink := &recordingEndpointSink{}
	r := newEndpointReconciler(sink)
	source := filtermcp.NewServerSource("nacos", "server-a")

	err := r.ApplyServerConfig(source, mcpConfigWithEndpointTools(
		endpointTool("a", "cluster-a", "http://127.0.0.1:8080"),
		endpointTool("b", "cluster-b", "http://127.0.0.1:8081"),
	))
	require.NoError(t, err)
	require.Len(t, sink.ops, 2)
	idA := stableEndpointID(source, "a")
	idB := stableEndpointID(source, "b")
	assert.ElementsMatch(t, []endpointOp{
		{action: "set", cluster: "cluster-a", id: idA, address: "127.0.0.1:8080"},
		{action: "set", cluster: "cluster-b", id: idB, address: "127.0.0.1:8081"},
	}, sink.ops)

	err = r.ApplyServerConfig(source, mcpConfigWithEndpointTools(
		endpointTool("a", "cluster-a", "http://127.0.0.1:8080"),
	))
	require.NoError(t, err)
	require.Len(t, sink.ops, 3)
	assert.Equal(t, endpointOp{action: "delete", cluster: "cluster-b", id: idB}, sink.ops[2])

	err = r.ApplyServerConfig(source, nil)
	require.NoError(t, err)
	require.Len(t, sink.ops, 4)
	assert.Equal(t, endpointOp{action: "delete", cluster: "cluster-a", id: idA}, sink.ops[3])

	err = r.ApplyServerConfig(source, nil)
	require.NoError(t, err)
	assert.Len(t, sink.ops, 4, "repeated tombstone must be idempotent")
}

func TestEndpointReconcilerSharedAddressDifferentOwners(t *testing.T) {
	sink := &recordingEndpointSink{}
	r := newEndpointReconciler(sink)

	err := r.ApplyServerConfig(filtermcp.NewServerSource("nacos", "server-a"), mcpConfigWithEndpointTools(
		endpointTool("tool", "cluster", "http://127.0.0.1:8080"),
	))
	require.NoError(t, err)
	err = r.ApplyServerConfig(filtermcp.NewServerSource("nacos", "server-b"), mcpConfigWithEndpointTools(
		endpointTool("tool", "cluster", "http://127.0.0.1:8080"),
	))
	require.NoError(t, err)

	require.Len(t, sink.ops, 2)
	assert.NotEqual(t, sink.ops[0].id, sink.ops[1].id)

	err = r.ApplyServerConfig(filtermcp.NewServerSource("nacos", "server-a"), nil)
	require.NoError(t, err)
	require.Len(t, sink.ops, 3)
	assert.Equal(t, sink.ops[0].id, sink.ops[2].id)
	assert.NotEqual(t, sink.ops[1].id, sink.ops[2].id)
}

func TestEndpointReconcilerSameServerIDDifferentRegistriesDoNotCollide(t *testing.T) {
	sink := &recordingEndpointSink{}
	r := newEndpointReconciler(sink)

	err := r.ApplyServerConfig(filtermcp.NewServerSource("registry1", "serverA"), mcpConfigWithEndpointTools(
		endpointTool("tool", "cluster", "http://127.0.0.1:8080"),
	))
	require.NoError(t, err)
	err = r.ApplyServerConfig(filtermcp.NewServerSource("registry2", "serverA"), mcpConfigWithEndpointTools(
		endpointTool("tool", "cluster", "http://127.0.0.1:8081"),
	))
	require.NoError(t, err)

	require.Len(t, sink.ops, 2)
	assert.Equal(t, "mcp/registry1/serverA/tool", sink.ops[0].id)
	assert.Equal(t, "mcp/registry2/serverA/tool", sink.ops[1].id)
}

func TestEndpointReconcilerRemoveAllDeletesPublishedEndpoints(t *testing.T) {
	sink := &recordingEndpointSink{}
	r := newEndpointReconciler(sink)

	require.NoError(t, r.ApplyServerConfig(filtermcp.NewServerSource("nacos", "server-a"), mcpConfigWithEndpointTools(
		endpointTool("a", "cluster-a", "http://127.0.0.1:8080"),
	)))
	require.NoError(t, r.ApplyServerConfig(filtermcp.NewServerSource("nacos", "server-b"), mcpConfigWithEndpointTools(
		endpointTool("b", "cluster-b", "http://127.0.0.1:8081"),
	)))
	require.Len(t, sink.ops, 2)

	r.RemoveAll()

	require.Len(t, sink.ops, 4)
	assert.ElementsMatch(t, []endpointOp{
		{action: "delete", cluster: "cluster-a", id: "mcp/nacos/server-a/a"},
		{action: "delete", cluster: "cluster-b", id: "mcp/nacos/server-b/b"},
	}, sink.ops[2:])
	assert.Empty(t, r.published)
}

func TestEndpointReconcilerInvalidDesiredDoesNotMutatePublishedState(t *testing.T) {
	sink := &recordingEndpointSink{}
	r := newEndpointReconciler(sink)
	source := filtermcp.NewServerSource("nacos", "server-a")

	err := r.ApplyServerConfig(source, mcpConfigWithEndpointTools(
		endpointTool("tool", "cluster", "http://127.0.0.1:8080"),
	))
	require.NoError(t, err)
	require.Len(t, sink.ops, 1)

	err = r.ApplyServerConfig(source, mcpConfigWithEndpointTools(
		endpointTool("tool", "cluster", "://bad-url"),
	))
	require.Error(t, err)
	assert.Len(t, sink.ops, 1)

	err = r.ApplyServerConfig(source, nil)
	require.NoError(t, err)
	require.Len(t, sink.ops, 2)
	assert.Equal(t, "delete", sink.ops[1].action)
	assert.Equal(t, sink.ops[0].id, sink.ops[1].id)
}

func TestEndpointReconcilerValidateDoesNotPublish(t *testing.T) {
	sink := &recordingEndpointSink{}
	r := newEndpointReconciler(sink)
	source := filtermcp.NewServerSource("nacos", "server-a")

	err := r.ValidateServerConfig(source, mcpConfigWithEndpointTools(
		endpointTool("tool", "cluster", "http://127.0.0.1:8080"),
	))
	require.NoError(t, err)
	assert.Empty(t, sink.ops)

	err = r.ValidateServerConfig(source, mcpConfigWithEndpointTools(
		endpointTool("tool", "cluster", "://bad-url"),
	))
	require.Error(t, err)
	assert.Empty(t, sink.ops)
}

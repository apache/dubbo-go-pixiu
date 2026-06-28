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
)

import (
	filtermcp "github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type testPublicationSink struct {
	runtimeID string
	applied   int
}

func (s *testPublicationSink) RuntimeID() string {
	return s.runtimeID
}

func (s *testPublicationSink) ApplyMcpServerConfigByServer(_ string, _ *model.McpServerConfig) error {
	s.applied++
	return nil
}

func TestAdapterBindPublicationSinkRequiresRuntime(t *testing.T) {
	filtermcp.ResetGlobalState()
	defer filtermcp.ResetGlobalState()

	a := &Adapter{}
	sink, err := a.bindPublicationSink()

	assert.ErrorIs(t, err, filtermcp.ErrDynamicConsumerUnavailable)
	assert.Nil(t, sink)
	assert.Nil(t, a.sink)
}

func TestAdapterApplyServerConfigEventNoBoundSinkDoesNothing(t *testing.T) {
	a := &Adapter{}
	reconciler := newEndpointReconciler(&recordingEndpointSink{})

	a.applyServerConfigEvent(reconciler, nil, "nacos", "server-a", &model.McpServerConfig{
		Tools: []model.ToolConfig{{Name: "tool", Cluster: "cluster", BackendURL: "http://127.0.0.1:8080"}},
	})

	assert.Empty(t, reconciler.published)
}

func TestAdapterApplyServerConfigEventUsesBoundRuntime(t *testing.T) {
	reconciler := newEndpointReconciler(&recordingEndpointSink{})
	sink := &testPublicationSink{runtimeID: "runtime-1"}
	a := &Adapter{sink: sink}

	a.applyServerConfigEvent(reconciler, sink, "nacos", "server-a", &model.McpServerConfig{
		Tools: []model.ToolConfig{{Name: "tool", Cluster: "cluster", BackendURL: "http://127.0.0.1:8080"}},
	})

	assert.Equal(t, 1, sink.applied)
}

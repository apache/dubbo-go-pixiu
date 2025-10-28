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

package a2a

import (
	"bytes"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/mock"
)

// createTestConfig creates a test configuration
func createTestConfig() *Config {
	return &Config{
		Endpoint: "/a2a",
		AgentInfo: &AgentConfig{
			AgentID:     "test-agent",
			Name:        "Test Agent",
			Version:     "1.0.0",
			Description: "Test agent for unit tests",
			Endpoint:    "http://localhost:8888/a2a",
			Status:      StatusOnline,
			Capabilities: []CapabilityConfig{
				{
					Name:        "test_capability",
					Description: "Test capability",
					Tags:        []string{"test"},
				},
			},
		},
		KnownAgents: []AgentConfig{
			{
				AgentID:  "agent-1",
				Name:     "Agent 1",
				Version:  "1.0.0",
				Endpoint: "http://localhost:9001/a2a",
				Status:   StatusOnline,
			},
		},
		TaskConfig: &TaskConfig{
			DefaultTimeout:     5000,
			MaxConcurrentTasks: 10,
			CleanupInterval:    1000,
			TaskRetention:      5000,
		},
	}
}

// createTestFilterFactory creates a filter factory for testing
func createTestFilterFactory() *FilterFactory {
	factory := &FilterFactory{
		cfg: createTestConfig(),
	}
	factory.Apply()
	return factory
}

// buildJSONRPCRequest builds a JSON-RPC request body
func buildJSONRPCRequest(method string, params interface{}) []byte {
	req := JSONRPCRequest{
		JSONRPC: JSONRPCVersion,
		Method:  method,
		Params:  params,
		ID:      1,
	}
	body, _ := json.Marshal(req)
	return body
}

// TestFilterFactory_Apply tests filter factory initialization
func TestFilterFactory_Apply(t *testing.T) {
	factory := createTestFilterFactory()

	assert.NotNil(t, factory.agentMap)
	assert.NotNil(t, factory.taskManager)
	assert.Equal(t, 2, len(factory.agentMap)) // test-agent + agent-1
}

// TestA2AFilter_HandlePing tests the ping method
func TestA2AFilter_HandlePing(t *testing.T) {
	factory := createTestFilterFactory()

	reqBody := buildJSONRPCRequest(MethodPing, nil)
	request, err := http.NewRequest("POST", "http://localhost:8888/a2a", bytes.NewReader(reqBody))
	assert.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	ctx := mock.GetMockHTTPContext(request)
	f := &A2AFilter{
		cfg:         factory.cfg,
		agentMap:    factory.agentMap,
		taskManager: factory.taskManager,
		mutex:       &factory.mutex,
	}

	status := f.Decode(ctx)
	assert.Equal(t, filter.Stop, status)

	// Verify response was written (we can't easily check the actual response in this test setup)
	// In a real scenario, you'd inspect ctx.Writer or use httptest.ResponseRecorder
}

// TestA2AFilter_HandleGetAgentCard tests the get_agent_card method
func TestA2AFilter_HandleGetAgentCard(t *testing.T) {
	factory := createTestFilterFactory()

	reqBody := buildJSONRPCRequest(MethodGetAgentCard, nil)
	request, err := http.NewRequest("POST", "http://localhost:8888/a2a", bytes.NewReader(reqBody))
	assert.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	ctx := mock.GetMockHTTPContext(request)
	f := &A2AFilter{
		cfg:         factory.cfg,
		agentMap:    factory.agentMap,
		taskManager: factory.taskManager,
		mutex:       &factory.mutex,
	}

	status := f.Decode(ctx)
	assert.Equal(t, filter.Stop, status)
}

// TestA2AFilter_HandleDiscoverAgents tests the discover_agents method
func TestA2AFilter_HandleDiscoverAgents(t *testing.T) {
	factory := createTestFilterFactory()

	reqBody := buildJSONRPCRequest(MethodDiscoverAgents, nil)
	request, err := http.NewRequest("POST", "http://localhost:8888/a2a", bytes.NewReader(reqBody))
	assert.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	ctx := mock.GetMockHTTPContext(request)
	f := &A2AFilter{
		cfg:         factory.cfg,
		agentMap:    factory.agentMap,
		taskManager: factory.taskManager,
		mutex:       &factory.mutex,
	}

	status := f.Decode(ctx)
	assert.Equal(t, filter.Stop, status)
}

// TestA2AFilter_HandleCreateTask tests task creation
func TestA2AFilter_HandleCreateTask(t *testing.T) {
	factory := createTestFilterFactory()

	params := CreateTaskRequest{
		To:   "agent-1",
		Type: "test_task",
		Content: map[string]interface{}{
			"message": "test",
		},
		Timeout: 3000,
	}

	reqBody := buildJSONRPCRequest(MethodCreateTask, params)
	request, err := http.NewRequest("POST", "http://localhost:8888/a2a", bytes.NewReader(reqBody))
	assert.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	ctx := mock.GetMockHTTPContext(request)
	f := &A2AFilter{
		cfg:         factory.cfg,
		agentMap:    factory.agentMap,
		taskManager: factory.taskManager,
		mutex:       &factory.mutex,
	}

	status := f.Decode(ctx)
	assert.Equal(t, filter.Stop, status)

	// Verify task was created
	assert.Equal(t, 1, factory.taskManager.GetTaskCount())
}

// TestA2AFilter_HandleCreateTask_MissingParams tests task creation with missing parameters
func TestA2AFilter_HandleCreateTask_MissingParams(t *testing.T) {
	factory := createTestFilterFactory()

	// Missing required fields
	params := CreateTaskRequest{
		To: "agent-1",
		// Missing Type and Content
	}

	reqBody := buildJSONRPCRequest(MethodCreateTask, params)
	request, err := http.NewRequest("POST", "http://localhost:8888/a2a", bytes.NewReader(reqBody))
	assert.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	ctx := mock.GetMockHTTPContext(request)
	f := &A2AFilter{
		cfg:         factory.cfg,
		agentMap:    factory.agentMap,
		taskManager: factory.taskManager,
		mutex:       &factory.mutex,
	}

	status := f.Decode(ctx)
	assert.Equal(t, filter.Stop, status)
}

// TestA2AFilter_HandleSendMessage tests sending a message with mock server
func TestA2AFilter_HandleSendMessage(t *testing.T) {
	// Create a mock agent server
	mockServer := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Simulate agent receiving message
		response := JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      1,
			Result:  map[string]string{"status": "received"},
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer mockServer.Close()

	// Create factory with mock server endpoint
	factory := createTestFilterFactory()
	factory.agentMap["agent-1"].Endpoint = mockServer.URL

	params := SendMessageRequest{
		To:   "agent-1",
		Type: "test_message",
		Content: map[string]interface{}{
			"text": "Hello, agent!",
		},
	}

	reqBody := buildJSONRPCRequest(MethodSendMessage, params)
	request, err := http.NewRequest("POST", "http://localhost:8888/a2a", bytes.NewReader(reqBody))
	assert.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	ctx := mock.GetMockHTTPContext(request)
	f := &A2AFilter{
		cfg:         factory.cfg,
		agentMap:    factory.agentMap,
		taskManager: factory.taskManager,
		mutex:       &factory.mutex,
	}

	status := f.Decode(ctx)
	assert.Equal(t, filter.Stop, status)
}

// TestA2AFilter_HandleSendMessage_AgentNotFound tests sending message to non-existent agent
func TestA2AFilter_HandleSendMessage_AgentNotFound(t *testing.T) {
	factory := createTestFilterFactory()

	params := SendMessageRequest{
		To:   "non-existent-agent",
		Type: "test_message",
		Content: map[string]interface{}{
			"text": "Hello",
		},
	}

	reqBody := buildJSONRPCRequest(MethodSendMessage, params)
	request, err := http.NewRequest("POST", "http://localhost:8888/a2a", bytes.NewReader(reqBody))
	assert.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	ctx := mock.GetMockHTTPContext(request)
	f := &A2AFilter{
		cfg:         factory.cfg,
		agentMap:    factory.agentMap,
		taskManager: factory.taskManager,
		mutex:       &factory.mutex,
	}

	status := f.Decode(ctx)
	assert.Equal(t, filter.Stop, status)
}

// TestA2AFilter_NonA2ARequest tests that non-A2A requests are ignored
func TestA2AFilter_NonA2ARequest(t *testing.T) {
	factory := createTestFilterFactory()

	request, err := http.NewRequest("GET", "http://localhost:8888/other-path", nil)
	assert.NoError(t, err)

	ctx := mock.GetMockHTTPContext(request)
	f := &A2AFilter{
		cfg:         factory.cfg,
		agentMap:    factory.agentMap,
		taskManager: factory.taskManager,
		mutex:       &factory.mutex,
	}

	status := f.Decode(ctx)
	assert.Equal(t, filter.Continue, status)
}

// TestA2AFilter_InvalidJSON tests handling of invalid JSON
func TestA2AFilter_InvalidJSON(t *testing.T) {
	factory := createTestFilterFactory()

	request, err := http.NewRequest("POST", "http://localhost:8888/a2a", bytes.NewReader([]byte("invalid json")))
	assert.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	ctx := mock.GetMockHTTPContext(request)
	f := &A2AFilter{
		cfg:         factory.cfg,
		agentMap:    factory.agentMap,
		taskManager: factory.taskManager,
		mutex:       &factory.mutex,
	}

	status := f.Decode(ctx)
	assert.Equal(t, filter.Stop, status)
}

// TestA2AFilter_UnknownMethod tests handling of unknown methods
func TestA2AFilter_UnknownMethod(t *testing.T) {
	factory := createTestFilterFactory()

	reqBody := buildJSONRPCRequest("a2a.unknown_method", nil)
	request, err := http.NewRequest("POST", "http://localhost:8888/a2a", bytes.NewReader(reqBody))
	assert.NoError(t, err)
	request.Header.Set("Content-Type", "application/json")

	ctx := mock.GetMockHTTPContext(request)
	f := &A2AFilter{
		cfg:         factory.cfg,
		agentMap:    factory.agentMap,
		taskManager: factory.taskManager,
		mutex:       &factory.mutex,
	}

	status := f.Decode(ctx)
	assert.Equal(t, filter.Stop, status)
}

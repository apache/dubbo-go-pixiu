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
	"encoding/json"
	"io"
	"strings"
	"sync"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// FilterFactory is the factory for creating A2A filter instances
type FilterFactory struct {
	cfg         *Config
	agentMap    map[string]*AgentInfo
	taskManager *TaskManager
	mutex       sync.RWMutex
}

// Config returns the configuration struct for the factory
func (f *FilterFactory) Config() any {
	return f.cfg
}

// Apply initializes the filter factory from its configuration
func (f *FilterFactory) Apply() error {
	logger.Infof("[dubbo-go-pixiu] A2A Server Filter factory applying configuration")

	// Initialize agent map
	f.agentMap = make(map[string]*AgentInfo)

	// Load current agent info
	agentInfo := f.cfg.GetDefaultAgentInfo()
	f.agentMap[agentInfo.AgentID] = agentInfo

	// Load known agents from configuration
	for _, agentConfig := range f.cfg.KnownAgents {
		agent := agentConfig.ToAgentInfo(f.cfg.Endpoint)
		f.agentMap[agent.AgentID] = agent
		logger.Infof("[dubbo-go-pixiu] A2A Server loaded known agent: %s (%s)", agent.Name, agent.AgentID)
	}

	// Initialize task manager
	taskConfig := f.cfg.GetTaskConfig()
	f.taskManager = NewTaskManager(taskConfig)

	logger.Infof("[dubbo-go-pixiu] A2A Server Filter factory initialized with %d agents", len(f.agentMap))
	return nil
}

// PrepareFilterChain prepares the filter chain for a new request
func (f *FilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	a2aFilter := &A2AFilter{
		cfg:         f.cfg,
		agentMap:    f.agentMap,
		taskManager: f.taskManager,
		mutex:       &f.mutex,
	}

	chain.AppendDecodeFilters(a2aFilter)
	return nil
}

// A2AFilter is the actual filter that processes A2A requests
type A2AFilter struct {
	cfg         *Config
	agentMap    map[string]*AgentInfo
	taskManager *TaskManager
	mutex       *sync.RWMutex
}

// Decode processes incoming HTTP requests for A2A protocol
func (f *A2AFilter) Decode(ctx *contexthttp.HttpContext) filter.FilterStatus {
	// Check if this is an A2A request
	if !f.isA2ARequest(ctx) {
		return filter.Continue
	}

	logger.Infof("[dubbo-go-pixiu] A2A Server processing request: %s %s", ctx.Request.Method, ctx.Request.URL.Path)

	// Read request body
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] A2A Server failed to read request body: %v", err)
		f.sendErrorResponse(ctx, nil, ErrorCodeInternalError, "Failed to read request body")
		return filter.Stop
	}

	// Parse JSON-RPC request
	var jsonrpcReq JSONRPCRequest
	if err := json.Unmarshal(body, &jsonrpcReq); err != nil {
		logger.Errorf("[dubbo-go-pixiu] A2A Server failed to parse JSON-RPC request: %v", err)
		f.sendErrorResponse(ctx, nil, ErrorCodeParseError, "Invalid JSON-RPC request")
		return filter.Stop
	}

	// Validate JSON-RPC version
	if jsonrpcReq.JSONRPC != JSONRPCVersion {
		logger.Warnf("[dubbo-go-pixiu] A2A Server invalid JSON-RPC version: %s", jsonrpcReq.JSONRPC)
		f.sendErrorResponse(ctx, jsonrpcReq.ID, ErrorCodeInvalidRequest, "Invalid JSON-RPC version")
		return filter.Stop
	}

	logger.Infof("[dubbo-go-pixiu] A2A Server received method: %s (id: %v)", jsonrpcReq.Method, jsonrpcReq.ID)

	// Process the A2A method
	response := f.processA2AMethod(&jsonrpcReq)

	// Send response
	f.sendJSONResponse(ctx, response)
	return filter.Stop
}

// isA2ARequest checks if the incoming request is an A2A request
func (f *A2AFilter) isA2ARequest(ctx *contexthttp.HttpContext) bool {
	// Check if the request path matches the A2A endpoint
	if !strings.HasPrefix(ctx.Request.URL.Path, f.cfg.Endpoint) {
		return false
	}

	// Check if it's a POST request with JSON content
	if ctx.Request.Method != "POST" {
		return false
	}

	contentType := ctx.Request.Header.Get("Content-Type")
	return strings.Contains(contentType, "application/json")
}

// processA2AMethod processes different A2A protocol methods
func (f *A2AFilter) processA2AMethod(req *JSONRPCRequest) *JSONRPCResponse {
	switch req.Method {
	case MethodPing:
		return f.handlePing(req)
	case MethodGetAgentCard:
		return f.handleGetAgentCard(req)
	case MethodDiscoverAgents:
		return f.handleDiscoverAgents(req)
	case MethodCreateTask:
		return f.handleCreateTask(req)
	case MethodGetTaskStatus:
		return f.handleGetTaskStatus(req)
	case MethodUpdateTask:
		return f.handleUpdateTask(req)
	case MethodCancelTask:
		return f.handleCancelTask(req)
	case MethodSendMessage:
		return f.handleSendMessage(req)
	case MethodBroadcastMessage:
		return f.handleBroadcastMessage(req)
	default:
		logger.Warnf("[dubbo-go-pixiu] A2A Server unsupported method: %s", req.Method)
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeMethodNotFound,
				Message: "Method not found",
			},
		}
	}
}

// Basic method handlers (placeholder implementations for phase 1)

// handlePing handles ping requests
func (f *A2AFilter) handlePing(req *JSONRPCRequest) *JSONRPCResponse {
	f.mutex.RLock()
	agentInfo := f.cfg.GetDefaultAgentInfo()
	f.mutex.RUnlock()

	response := PingResponse{
		Status:    "pong",
		AgentID:   agentInfo.AgentID,
		Timestamp: time.Now().UnixMilli(),
	}

	return &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		ID:      req.ID,
		Result:  response,
	}
}

// handleGetAgentCard handles get agent card requests
func (f *A2AFilter) handleGetAgentCard(req *JSONRPCRequest) *JSONRPCResponse {
	f.mutex.RLock()
	agentInfo := f.cfg.GetDefaultAgentInfo()
	f.mutex.RUnlock()

	return &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		ID:      req.ID,
		Result:  agentInfo,
	}
}

// handleDiscoverAgents handles agent discovery requests
func (f *A2AFilter) handleDiscoverAgents(req *JSONRPCRequest) *JSONRPCResponse {
	f.mutex.RLock()
	defer f.mutex.RUnlock()

	// For now, return all known agents
	var agents []AgentInfo
	for _, agent := range f.agentMap {
		agents = append(agents, *agent)
	}

	response := DiscoverAgentsResponse{
		Agents: agents,
	}

	return &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		ID:      req.ID,
		Result:  response,
	}
}

// Placeholder implementations for other methods (to be implemented in later phases)

func (f *A2AFilter) handleCreateTask(req *JSONRPCRequest) *JSONRPCResponse {
	// TODO: Implement in phase 3
	return &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		ID:      req.ID,
		Error: &RPCError{
			Code:    ErrorCodeInternalError,
			Message: "Task management not implemented yet",
		},
	}
}

func (f *A2AFilter) handleGetTaskStatus(req *JSONRPCRequest) *JSONRPCResponse {
	// TODO: Implement in phase 3
	return f.handleCreateTask(req)
}

func (f *A2AFilter) handleUpdateTask(req *JSONRPCRequest) *JSONRPCResponse {
	// TODO: Implement in phase 3
	return f.handleCreateTask(req)
}

func (f *A2AFilter) handleCancelTask(req *JSONRPCRequest) *JSONRPCResponse {
	// TODO: Implement in phase 3
	return f.handleCreateTask(req)
}

func (f *A2AFilter) handleSendMessage(req *JSONRPCRequest) *JSONRPCResponse {
	// TODO: Implement in phase 4
	return &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		ID:      req.ID,
		Error: &RPCError{
			Code:    ErrorCodeInternalError,
			Message: "Message routing not implemented yet",
		},
	}
}

func (f *A2AFilter) handleBroadcastMessage(req *JSONRPCRequest) *JSONRPCResponse {
	// TODO: Implement in phase 4
	return f.handleSendMessage(req)
}

// Utility methods

// sendErrorResponse sends an error response to the client
func (f *A2AFilter) sendErrorResponse(ctx *contexthttp.HttpContext, id interface{}, code int, message string) {
	response := &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		ID:      id,
		Error: &RPCError{
			Code:    code,
			Message: message,
		},
	}
	f.sendJSONResponse(ctx, response)
}

// sendJSONResponse sends a JSON response to the client
func (f *A2AFilter) sendJSONResponse(ctx *contexthttp.HttpContext, response *JSONRPCResponse) {
	responseBody, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] A2A Server failed to marshal response: %v", err)
		ctx.SendLocalReply(500, []byte("Internal Server Error"))
		return
	}

	// SendLocalReply automatically sets Content-Type to application/json for valid JSON
	ctx.SendLocalReply(200, responseBody)
}

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
	"fmt"
	"io"
	"net/http"
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
func (f *FilterFactory) PrepareFilterChain(_ *contexthttp.HttpContext, chain filter.FilterChain) error {
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
	// Parse request parameters
	paramsBytes, err := json.Marshal(req.Params)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] A2A Server failed to marshal create task params: %v", err)
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInvalidParams,
				Message: "Invalid parameters",
			},
		}
	}

	var params CreateTaskRequest
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		logger.Errorf("[dubbo-go-pixiu] A2A Server failed to parse create task params: %v", err)
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInvalidParams,
				Message: "Invalid task parameters",
			},
		}
	}

	// Validate required fields
	if params.To == "" || params.Type == "" || params.Content == nil {
		logger.Warnf("[dubbo-go-pixiu] A2A Server create task missing required fields")
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInvalidParams,
				Message: "Missing required fields: to, type, content",
			},
		}
	}

	// Get current agent info as the sender
	f.mutex.RLock()
	agentInfo := f.cfg.GetDefaultAgentInfo()
	f.mutex.RUnlock()

	// Create task
	task, err := f.taskManager.CreateTask(agentInfo.AgentID, params.To, params.Type, params.Content, params.Timeout)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] A2A Server failed to create task: %v", err)
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInternalError,
				Message: fmt.Sprintf("Failed to create task: %v", err),
			},
		}
	}

	logger.Infof("[dubbo-go-pixiu] A2A Server created task %s: %s -> %s", task.TaskID, task.From, task.To)

	// Build response
	response := CreateTaskResponse{
		TaskID: task.TaskID,
		Status: task.Status,
	}

	return &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		ID:      req.ID,
		Result:  response,
	}
}

func (f *A2AFilter) handleGetTaskStatus(req *JSONRPCRequest) *JSONRPCResponse {
	// Parse request parameters
	paramsBytes, err := json.Marshal(req.Params)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] A2A Server failed to marshal get task status params: %v", err)
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInvalidParams,
				Message: "Invalid parameters",
			},
		}
	}

	var params GetTaskStatusRequest
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		logger.Errorf("[dubbo-go-pixiu] A2A Server failed to parse get task status params: %v", err)
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInvalidParams,
				Message: "Invalid task status parameters",
			},
		}
	}

	// Validate task ID
	if params.TaskID == "" {
		logger.Warnf("[dubbo-go-pixiu] A2A Server get task status missing task_id")
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInvalidParams,
				Message: "Missing required field: task_id",
			},
		}
	}

	// Get task from task manager
	task, err := f.taskManager.GetTask(params.TaskID)
	if err != nil {
		logger.Warnf("[dubbo-go-pixiu] A2A Server task not found: %s", params.TaskID)
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeTaskNotFound,
				Message: fmt.Sprintf("Task not found: %s", params.TaskID),
			},
		}
	}

	logger.Infof("[dubbo-go-pixiu] A2A Server retrieved task status: %s (status: %s)", params.TaskID, task.Status)

	// Build response
	response := GetTaskStatusResponse{
		Task: *task,
	}

	return &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		ID:      req.ID,
		Result:  response,
	}
}

func (f *A2AFilter) handleUpdateTask(req *JSONRPCRequest) *JSONRPCResponse {
	// Parse request parameters
	paramsBytes, err := json.Marshal(req.Params)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] A2A Server failed to marshal update task params: %v", err)
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInvalidParams,
				Message: "Invalid parameters",
			},
		}
	}

	var params UpdateTaskRequest
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		logger.Errorf("[dubbo-go-pixiu] A2A Server failed to parse update task params: %v", err)
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInvalidParams,
				Message: "Invalid update task parameters",
			},
		}
	}

	// Validate task ID
	if params.TaskID == "" {
		logger.Warnf("[dubbo-go-pixiu] A2A Server update task missing task_id")
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInvalidParams,
				Message: "Missing required field: task_id",
			},
		}
	}

	// Update task
	err = f.taskManager.UpdateTask(params.TaskID, params.Status, params.Result, params.Error)
	if err != nil {
		logger.Warnf("[dubbo-go-pixiu] A2A Server failed to update task %s: %v", params.TaskID, err)
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeTaskNotFound,
				Message: fmt.Sprintf("Failed to update task: %v", err),
			},
		}
	}

	logger.Infof("[dubbo-go-pixiu] A2A Server updated task %s to status %s", params.TaskID, params.Status)

	// Build response
	response := UpdateTaskResponse{
		Success: true,
		Message: "Task updated successfully",
	}

	return &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		ID:      req.ID,
		Result:  response,
	}
}

func (f *A2AFilter) handleCancelTask(req *JSONRPCRequest) *JSONRPCResponse {
	// Parse request parameters
	paramsBytes, err := json.Marshal(req.Params)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] A2A Server failed to marshal cancel task params: %v", err)
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInvalidParams,
				Message: "Invalid parameters",
			},
		}
	}

	var params struct {
		TaskID string `json:"task_id"`
	}
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		logger.Errorf("[dubbo-go-pixiu] A2A Server failed to parse cancel task params: %v", err)
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInvalidParams,
				Message: "Invalid cancel task parameters",
			},
		}
	}

	// Validate task ID
	if params.TaskID == "" {
		logger.Warnf("[dubbo-go-pixiu] A2A Server cancel task missing task_id")
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInvalidParams,
				Message: "Missing required field: task_id",
			},
		}
	}

	// Cancel task
	err = f.taskManager.CancelTask(params.TaskID)
	if err != nil {
		logger.Warnf("[dubbo-go-pixiu] A2A Server failed to cancel task %s: %v", params.TaskID, err)
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeTaskNotFound,
				Message: fmt.Sprintf("Failed to cancel task: %v", err),
			},
		}
	}

	logger.Infof("[dubbo-go-pixiu] A2A Server cancelled task %s", params.TaskID)

	// Build response
	response := UpdateTaskResponse{
		Success: true,
		Message: "Task cancelled successfully",
	}

	return &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		ID:      req.ID,
		Result:  response,
	}
}

func (f *A2AFilter) handleSendMessage(req *JSONRPCRequest) *JSONRPCResponse {
	// Parse request parameters
	paramsBytes, err := json.Marshal(req.Params)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] A2A Server failed to marshal send message params: %v", err)
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInvalidParams,
				Message: "Invalid parameters",
			},
		}
	}

	var params SendMessageRequest
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		logger.Errorf("[dubbo-go-pixiu] A2A Server failed to parse send message params: %v", err)
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInvalidParams,
				Message: "Invalid send message parameters",
			},
		}
	}

	// Validate required fields
	if params.To == "" || params.Type == "" || params.Content == nil {
		logger.Warnf("[dubbo-go-pixiu] A2A Server send message missing required fields")
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInvalidParams,
				Message: "Missing required fields: to, type, content",
			},
		}
	}

	// Find target agent
	f.mutex.RLock()
	targetAgent, exists := f.agentMap[params.To]
	agentInfo := f.cfg.GetDefaultAgentInfo()
	f.mutex.RUnlock()

	if !exists {
		logger.Warnf("[dubbo-go-pixiu] A2A Server target agent not found: %s", params.To)
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeAgentNotFound,
				Message: fmt.Sprintf("Agent not found: %s", params.To),
			},
		}
	}

	// Generate message ID
	messageID := fmt.Sprintf("msg_%d_%d", time.Now().UnixMilli(), time.Now().UnixNano()%1000)

	// Build message
	message := &Message{
		MessageID: messageID,
		From:      agentInfo.AgentID,
		To:        params.To,
		Type:      params.Type,
		Content:   params.Content,
		Timestamp: time.Now().UnixMilli(),
	}

	// Forward message to agent
	err = f.forwardMessageToAgent(targetAgent, message)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] A2A Server failed to forward message to %s: %v", params.To, err)
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInternalError,
				Message: fmt.Sprintf("Failed to send message: %v", err),
			},
		}
	}

	logger.Infof("[dubbo-go-pixiu] A2A Server sent message %s to agent %s", messageID, params.To)

	// Build response
	response := SendMessageResponse{
		MessageID: messageID,
		Success:   true,
		Message:   "Message sent successfully",
	}

	return &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		ID:      req.ID,
		Result:  response,
	}
}

func (f *A2AFilter) handleBroadcastMessage(req *JSONRPCRequest) *JSONRPCResponse {
	// Parse request parameters
	paramsBytes, err := json.Marshal(req.Params)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] A2A Server failed to marshal broadcast message params: %v", err)
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInvalidParams,
				Message: "Invalid parameters",
			},
		}
	}

	var params SendMessageRequest
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		logger.Errorf("[dubbo-go-pixiu] A2A Server failed to parse broadcast message params: %v", err)
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInvalidParams,
				Message: "Invalid broadcast message parameters",
			},
		}
	}

	// Validate required fields
	if params.Type == "" || params.Content == nil {
		logger.Warnf("[dubbo-go-pixiu] A2A Server broadcast message missing required fields")
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Error: &RPCError{
				Code:    ErrorCodeInvalidParams,
				Message: "Missing required fields: type, content",
			},
		}
	}

	// Get current agent info
	f.mutex.RLock()
	agentInfo := f.cfg.GetDefaultAgentInfo()

	// Collect target agents (online agents excluding self)
	var targetAgents []*AgentInfo
	for _, agent := range f.agentMap {
		if agent.AgentID != agentInfo.AgentID && agent.Status == StatusOnline {
			targetAgents = append(targetAgents, agent)
		}
	}
	f.mutex.RUnlock()

	totalAgents := len(targetAgents)
	if totalAgents == 0 {
		logger.Infof("[dubbo-go-pixiu] A2A Server no online agents to broadcast to")
		return &JSONRPCResponse{
			JSONRPC: JSONRPCVersion,
			ID:      req.ID,
			Result: BroadcastMessageResponse{
				TotalAgents:  0,
				SuccessCount: 0,
				MessageIDs:   []string{},
			},
		}
	}

	// Use channel to collect results
	type result struct {
		agentID   string
		messageID string
		err       error
	}
	resultCh := make(chan result, totalAgents)

	// Broadcast to all target agents concurrently
	for _, agent := range targetAgents {
		go func(targetAgent *AgentInfo) {
			// Generate unique message ID
			messageID := fmt.Sprintf("msg_%d_%s_%d", time.Now().UnixMilli(), targetAgent.AgentID, time.Now().UnixNano()%1000)

			// Build message
			message := &Message{
				MessageID: messageID,
				From:      agentInfo.AgentID,
				To:        targetAgent.AgentID,
				Type:      params.Type,
				Content:   params.Content,
				Timestamp: time.Now().UnixMilli(),
			}

			// Forward message
			err := f.forwardMessageToAgent(targetAgent, message)
			resultCh <- result{
				agentID:   targetAgent.AgentID,
				messageID: messageID,
				err:       err,
			}
		}(agent)
	}

	// Collect results
	var successCount int
	var failedAgents []string
	var messageIDs []string

	for i := 0; i < totalAgents; i++ {
		res := <-resultCh
		if res.err != nil {
			logger.Errorf("[dubbo-go-pixiu] A2A Server failed to broadcast to agent %s: %v", res.agentID, res.err)
			failedAgents = append(failedAgents, res.agentID)
		} else {
			successCount++
			messageIDs = append(messageIDs, res.messageID)
		}
	}
	close(resultCh)

	logger.Infof("[dubbo-go-pixiu] A2A Server broadcast message to %d agents: %d succeeded, %d failed",
		totalAgents, successCount, len(failedAgents))

	// Build response
	response := BroadcastMessageResponse{
		TotalAgents:  totalAgents,
		SuccessCount: successCount,
		FailedAgents: failedAgents,
		MessageIDs:   messageIDs,
	}

	return &JSONRPCResponse{
		JSONRPC: JSONRPCVersion,
		ID:      req.ID,
		Result:  response,
	}
}

// Utility methods

// sendErrorResponse sends an error response to the client
func (f *A2AFilter) sendErrorResponse(ctx *contexthttp.HttpContext, id any, code int, message string) {
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

// forwardMessageToAgent forwards a message to a specific agent via HTTP
func (f *A2AFilter) forwardMessageToAgent(agent *AgentInfo, message *Message) error {
	// Build JSON-RPC request for receiving message
	jsonrpcReq := JSONRPCRequest{
		JSONRPC: JSONRPCVersion,
		Method:  "a2a.receive_message",
		Params:  message,
		ID:      fmt.Sprintf("msg_%d", time.Now().UnixMilli()),
	}

	// Marshal request body
	reqBody, err := json.Marshal(jsonrpcReq)
	if err != nil {
		return fmt.Errorf("failed to marshal message request: %v", err)
	}

	// Create HTTP request
	httpReq, err := http.NewRequest("POST", agent.Endpoint, bytes.NewReader(reqBody))
	if err != nil {
		return fmt.Errorf("failed to create HTTP request: %v", err)
	}

	// Set headers
	httpReq.Header.Set("Content-Type", "application/json")

	// Create HTTP client with timeout
	httpClient := &http.Client{
		Timeout: 30 * time.Second,
	}

	// Send request
	resp, err := httpClient.Do(httpReq)
	if err != nil {
		return fmt.Errorf("failed to send message to agent %s: %v", agent.AgentID, err)
	}
	defer resp.Body.Close()

	// Check response status
	if resp.StatusCode != 200 {
		bodyBytes, _ := io.ReadAll(resp.Body)
		return fmt.Errorf("agent %s returned status %d: %s", agent.AgentID, resp.StatusCode, string(bodyBytes))
	}

	// Parse response
	var jsonrpcResp JSONRPCResponse
	if err := json.NewDecoder(resp.Body).Decode(&jsonrpcResp); err != nil {
		return fmt.Errorf("failed to parse response from agent %s: %v", agent.AgentID, err)
	}

	// Check for JSON-RPC error
	if jsonrpcResp.Error != nil {
		return fmt.Errorf("agent %s returned error: %s", agent.AgentID, jsonrpcResp.Error.Message)
	}

	logger.Infof("[dubbo-go-pixiu] A2A Server successfully forwarded message %s to agent %s", message.MessageID, agent.AgentID)
	return nil
}

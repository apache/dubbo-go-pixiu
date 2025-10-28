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

// JSONRPCRequest represents a JSON-RPC 2.0 request
type JSONRPCRequest struct {
	JSONRPC string `json:"jsonrpc"`
	Method  string `json:"method"`
	Params  any    `json:"params,omitempty"`
	ID      any    `json:"id"`
}

// JSONRPCResponse represents a JSON-RPC 2.0 response
type JSONRPCResponse struct {
	JSONRPC string    `json:"jsonrpc"`
	Result  any       `json:"result,omitempty"`
	Error   *RPCError `json:"error,omitempty"`
	ID      any       `json:"id"`
}

// RPCError represents a JSON-RPC error object
type RPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// AgentStatus represents the status of an agent
type AgentStatus string

const (
	StatusOnline  AgentStatus = AgentStatusOnline
	StatusOffline AgentStatus = AgentStatusOffline
	StatusBusy    AgentStatus = AgentStatusBusy
)

// TaskStatus represents the status of a task
type TaskStatus string

const (
	TaskPending   TaskStatus = TaskStatusPending
	TaskRunning   TaskStatus = TaskStatusRunning
	TaskCompleted TaskStatus = TaskStatusCompleted
	TaskFailed    TaskStatus = TaskStatusFailed
	TaskCancelled TaskStatus = TaskStatusCancelled
)

// AgentInfo represents information about an agent
type AgentInfo struct {
	AgentID      string         `json:"agent_id"`
	Name         string         `json:"name"`
	Version      string         `json:"version"`
	Description  string         `json:"description,omitempty"`
	Endpoint     string         `json:"endpoint"`
	Status       AgentStatus    `json:"status"`
	Capabilities []Capability   `json:"capabilities,omitempty"`
	Metadata     map[string]any `json:"metadata,omitempty"`
}

// Capability represents a capability that an agent provides
type Capability struct {
	Name        string      `json:"name"`
	Description string      `json:"description,omitempty"`
	InputTypes  []string    `json:"input_types,omitempty"`
	OutputTypes []string    `json:"output_types,omitempty"`
	Tags        []string    `json:"tags,omitempty"`
	Parameters  []Parameter `json:"parameters,omitempty"`
}

// Parameter represents a capability parameter
type Parameter struct {
	Name        string `json:"name"`
	Type        string `json:"type"`
	Description string `json:"description,omitempty"`
	Required    bool   `json:"required"`
	Default     any    `json:"default,omitempty"`
}

// Task represents a task that can be executed by an agent
type Task struct {
	TaskID    string         `json:"task_id"`
	From      string         `json:"from"`
	To        string         `json:"to"`
	Type      string         `json:"type"`
	Content   map[string]any `json:"content"`
	Status    TaskStatus     `json:"status"`
	Result    map[string]any `json:"result,omitempty"`
	Error     string         `json:"error,omitempty"`
	CreatedAt int64          `json:"created_at"`
	UpdatedAt int64          `json:"updated_at"`
	Timeout   int64          `json:"timeout,omitempty"`
}

// Message represents a message sent between agents
type Message struct {
	MessageID string         `json:"message_id"`
	From      string         `json:"from"`
	To        string         `json:"to"`
	Type      string         `json:"type"`
	Content   map[string]any `json:"content"`
	Timestamp int64          `json:"timestamp"`
}

// PingRequest represents a ping request
type PingRequest struct {
	Timestamp int64 `json:"timestamp,omitempty"`
}

// PingResponse represents a ping response
type PingResponse struct {
	Status    string `json:"status"`
	AgentID   string `json:"agent_id"`
	Timestamp int64  `json:"timestamp"`
}

// DiscoverAgentsRequest represents a request to discover agents
type DiscoverAgentsRequest struct {
	Capabilities []string          `json:"capabilities,omitempty"`
	Tags         []string          `json:"tags,omitempty"`
	Status       AgentStatus       `json:"status,omitempty"`
	Filters      map[string]string `json:"filters,omitempty"`
}

// DiscoverAgentsResponse represents the response to an agent discovery request
type DiscoverAgentsResponse struct {
	Agents []AgentInfo `json:"agents"`
}

// CreateTaskRequest represents a request to create a task
type CreateTaskRequest struct {
	To      string         `json:"to"`
	Type    string         `json:"type"`
	Content map[string]any `json:"content"`
	Timeout int64          `json:"timeout,omitempty"`
}

// CreateTaskResponse represents the response to a create task request
type CreateTaskResponse struct {
	TaskID string     `json:"task_id"`
	Status TaskStatus `json:"status"`
}

// GetTaskStatusRequest represents a request to get task status
type GetTaskStatusRequest struct {
	TaskID string `json:"task_id"`
}

// GetTaskStatusResponse represents the response to a get task status request
type GetTaskStatusResponse struct {
	Task Task `json:"task"`
}

// UpdateTaskRequest represents a request to update a task
type UpdateTaskRequest struct {
	TaskID string         `json:"task_id"`
	Status TaskStatus     `json:"status,omitempty"`
	Result map[string]any `json:"result,omitempty"`
	Error  string         `json:"error,omitempty"`
}

// UpdateTaskResponse represents the response to an update task request
type UpdateTaskResponse struct {
	Success bool   `json:"success"`
	Message string `json:"message,omitempty"`
}

// SendMessageRequest represents a request to send a message
type SendMessageRequest struct {
	To      string         `json:"to"`
	Type    string         `json:"type"`
	Content map[string]any `json:"content"`
}

// SendMessageResponse represents the response to a send message request
type SendMessageResponse struct {
	MessageID string `json:"message_id"`
	Success   bool   `json:"success"`
	Message   string `json:"message,omitempty"`
}

// BroadcastMessageResponse represents the response to a broadcast message request
type BroadcastMessageResponse struct {
	TotalAgents  int      `json:"total_agents"`
	SuccessCount int      `json:"success_count"`
	FailedAgents []string `json:"failed_agents,omitempty"`
	MessageIDs   []string `json:"message_ids"`
}

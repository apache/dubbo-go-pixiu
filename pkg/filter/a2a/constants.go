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
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
)

const (
	// Kind is the type identifier for A2A Server Filter
	Kind = constant.A2AServerFilter

	// DefaultEndpoint is the default A2A service endpoint
	DefaultEndpoint = "/a2a"

	// DefaultAgentVersion is the default agent version
	DefaultAgentVersion = "1.0.0"

	// JSON-RPC version constant
	JSONRPCVersion = "2.0"

	// A2A Protocol Methods

	// Core methods
	MethodPing           = "a2a.ping"
	MethodGetAgentCard   = "a2a.get_agent_card"
	MethodDiscoverAgents = "a2a.discover_agents"

	// Task management methods
	MethodCreateTask    = "a2a.create_task"
	MethodGetTaskStatus = "a2a.get_task_status"
	MethodUpdateTask    = "a2a.update_task"
	MethodCancelTask    = "a2a.cancel_task"

	// Communication methods
	MethodSendMessage      = "a2a.send_message"
	MethodBroadcastMessage = "a2a.broadcast_message"

	// Agent status constants
	AgentStatusOnline  = "online"
	AgentStatusOffline = "offline"
	AgentStatusBusy    = "busy"

	// Task status constants
	TaskStatusPending   = "pending"
	TaskStatusRunning   = "running"
	TaskStatusCompleted = "completed"
	TaskStatusFailed    = "failed"
	TaskStatusCancelled = "canceled"

	// Error codes for JSON-RPC errors
	ErrorCodeParseError     = -32700
	ErrorCodeInvalidRequest = -32600
	ErrorCodeMethodNotFound = -32601
	ErrorCodeInvalidParams  = -32602
	ErrorCodeInternalError  = -32603

	// A2A specific error codes
	ErrorCodeAgentNotFound = -32001
	ErrorCodeTaskNotFound  = -32002
	ErrorCodeTaskTimeout   = -32003

	// Default configuration values
	DefaultTaskTimeout         = 30000 // 30 seconds in milliseconds
	DefaultMaxConcurrentTasks  = 100
	DefaultTaskCleanupInterval = 300000 // 5 minutes in milliseconds
)

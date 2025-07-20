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
	"fmt"

	"github.com/mark3labs/mcp-go/mcp"
)

// ResponseBuilder provides methods to create standardized MCP responses
type ResponseBuilder struct{}

// NewResponseBuilder creates a new response builder
func NewResponseBuilder() *ResponseBuilder {
	return &ResponseBuilder{}
}

// Success creates a successful JSON-RPC response
func (rb *ResponseBuilder) Success(id any, result any) mcp.JSONRPCResponse {
	return mcp.JSONRPCResponse{
		JSONRPC: mcp.JSONRPC_VERSION,
		ID:      mcp.NewRequestId(id),
		Result:  result,
	}
}

// Error creates an error JSON-RPC response
func (rb *ResponseBuilder) Error(id any, code int, message string) mcp.JSONRPCError {
	return mcp.NewJSONRPCError(mcp.NewRequestId(id), code, message, nil)
}

// ToolCallSuccess creates a successful tool call response
func (rb *ResponseBuilder) ToolCallSuccess(id any, content string) mcp.JSONRPCResponse {
	// Use mcp-go API to create text content
	textContent := mcp.NewTextContent(content)

	// Build MCP tool call result using mcp-go structures
	result := mcp.CallToolResult{
		Content: []mcp.Content{textContent},
		IsError: false,
	}

	return rb.Success(id, result)
}

// ToolCallError creates an error tool call response
func (rb *ResponseBuilder) ToolCallError(id any, message string) mcp.JSONRPCResponse {
	errorText := fmt.Sprintf("Error: %s", message)
	textContent := mcp.NewTextContent(errorText)

	// Build MCP tool call error result using mcp-go structures
	result := mcp.CallToolResult{
		Content: []mcp.Content{textContent},
		IsError: true,
	}

	return rb.Success(id, result)
}

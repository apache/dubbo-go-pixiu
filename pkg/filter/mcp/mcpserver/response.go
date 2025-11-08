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
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
)

import (
	"github.com/mark3labs/mcp-go/mcp"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// ResponseBuilder provides methods to create standardized MCP responses
type ResponseBuilder struct{}

var (
	responseBuilderInstance *ResponseBuilder
	responseBuilderOnce     sync.Once
)

// NewResponseBuilder returns the singleton ResponseBuilder instance
func NewResponseBuilder() *ResponseBuilder {
	responseBuilderOnce.Do(func() {
		responseBuilderInstance = &ResponseBuilder{}
	})
	return responseBuilderInstance
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

// ErrorHandler provides centralized error handling for MCP responses
type ErrorHandler struct {
	responseBuilder *ResponseBuilder
}

var (
	errorHandlerInstance *ErrorHandler
	errorHandlerOnce     sync.Once
)

// NewErrorHandler returns the singleton ErrorHandler instance
func NewErrorHandler() *ErrorHandler {
	errorHandlerOnce.Do(func() {
		errorHandlerInstance = &ErrorHandler{
			responseBuilder: NewResponseBuilder(),
		}
	})
	return errorHandlerInstance
}

// SendInternalError sends an internal server error response
func (eh *ErrorHandler) SendInternalError(ctx *MCPContext, id any, message string) filter.FilterStatus {
	response := eh.responseBuilder.Error(id, mcp.INTERNAL_ERROR, message)
	return eh.sendResponse(ctx, response)
}

// SendMethodNotFound sends a method not found error response
func (eh *ErrorHandler) SendMethodNotFound(ctx *MCPContext, id any) filter.FilterStatus {
	response := eh.responseBuilder.Error(id, mcp.METHOD_NOT_FOUND, "Method not found")
	return eh.sendResponse(ctx, response)
}

// SendInvalidParams sends an invalid parameters error response
func (eh *ErrorHandler) SendInvalidParams(ctx *MCPContext, id any, message string) filter.FilterStatus {
	response := eh.responseBuilder.Error(id, mcp.INVALID_PARAMS, fmt.Sprintf("Invalid params: %s", message))
	return eh.sendResponse(ctx, response)
}

// SendToolCallError sends a tool call error response
func (eh *ErrorHandler) SendToolCallError(ctx *MCPContext, id any, message string) filter.FilterStatus {
	response := eh.responseBuilder.ToolCallError(id, message)
	return eh.sendResponse(ctx, response)
}

// sendResponse sends any response and handles Content-Length cleanup
func (eh *ErrorHandler) sendResponse(ctx *MCPContext, response any) filter.FilterStatus {
	responseBody, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to marshal response: %v", err)
		errResp := contexthttp.InternalError.WithError(fmt.Errorf("marshal response failed: %w", err))
		ctx.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	// Critical: Clear Content-Length header to prevent mismatch errors
	ctx.ClearContentLengthHeader()
	ctx.SendLocalReply(http.StatusOK, responseBody)
	return filter.Stop
}

// ServerNotification creates a JSON-RPC notification from server
func (rb *ResponseBuilder) ServerNotification(method string, params map[string]any) mcp.JSONRPCNotification {
	notificationParams := mcp.NotificationParams{
		AdditionalFields: params,
	}

	return mcp.JSONRPCNotification{
		JSONRPC: mcp.JSONRPC_VERSION,
		Notification: mcp.Notification{
			Method: method,
			Params: notificationParams,
		},
	}
}

// ServerRequest creates a JSON-RPC request from server
func (rb *ResponseBuilder) ServerRequest(id any, method string, params any) mcp.JSONRPCRequest {
	return mcp.JSONRPCRequest{
		JSONRPC: mcp.JSONRPC_VERSION,
		ID:      mcp.NewRequestId(id),
		Params:  params,
		Request: mcp.Request{
			Method: method,
		},
	}
}

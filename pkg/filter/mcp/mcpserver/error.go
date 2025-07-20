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

	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// ErrorHandler provides centralized error handling for MCP responses
type ErrorHandler struct {
	responseBuilder *ResponseBuilder
}

// NewErrorHandler creates a new error handler
func NewErrorHandler() *ErrorHandler {
	return &ErrorHandler{
		responseBuilder: NewResponseBuilder(),
	}
}

// SendInternalError sends an internal server error response
func (eh *ErrorHandler) SendInternalError(ctx *MCPContext, id any, message string) filter.FilterStatus {
	response := eh.responseBuilder.Error(id, -32603, message)
	return eh.sendResponse(ctx, response)
}

// SendMethodNotFound sends a method not found error response
func (eh *ErrorHandler) SendMethodNotFound(ctx *MCPContext, id any) filter.FilterStatus {
	response := eh.responseBuilder.Error(id, -32601, "Method not found")
	return eh.sendResponse(ctx, response)
}

// SendInvalidParams sends an invalid parameters error response
func (eh *ErrorHandler) SendInvalidParams(ctx *MCPContext, id any, message string) filter.FilterStatus {
	response := eh.responseBuilder.Error(id, -32602, fmt.Sprintf("Invalid params: %s", message))
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
		ctx.SendLocalReply(http.StatusInternalServerError, []byte("internal server error"))
		return filter.Stop
	}

	// Critical: Clear Content-Length header to prevent mismatch errors
	ctx.Writer.Header().Del("Content-Length")
	ctx.SendLocalReply(http.StatusOK, responseBody)
	return filter.Stop
}

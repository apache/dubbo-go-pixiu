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
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"net/http"
	"strings"

	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/mark3labs/mcp-go/mcp"
)

// Note: ToolCallProcessor removed - using ResponseBuilder.ToolCallSuccess() directly for cleaner code

// handleToolCallResponse handles tool call responses, wrapping backend responses in MCP format
func (f *MCPServerFilter) handleToolCallResponse(ctx *MCPContext) filter.FilterStatus {
	logger.Debugf("[dubbo-go-pixiu] mcp server handling tool call response")

	// Extract request information
	requestID := ctx.GetMCPRequestID()
	if requestID == nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server missing request ID for tool call response")
		return filter.Continue
	}

	// Extract backend response
	responseBody, statusCode, err := f.extractBackendResponse(ctx)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to extract backend response: %v", err)
		return f.errorHandler.SendToolCallError(ctx, requestID, "failed to process backend response")
	}

	// Check for duplicate wrapping
	if f.isAlreadyMCPFormat(responseBody) {
		logger.Debugf("[dubbo-go-pixiu] mcp server response already in MCP format, skipping processing")
		return filter.Continue
	}

	// Process the response
	return f.processToolCallResponse(ctx, requestID, responseBody, statusCode)
}

// extractBackendResponse extracts response data from the context
func (f *MCPServerFilter) extractBackendResponse(ctx *MCPContext) ([]byte, int, error) {
	if ctx.TargetResp == nil {
		return nil, 0, fmt.Errorf("no target response available")
	}

	unaryResp, ok := ctx.TargetResp.(*client.UnaryResponse)
	if !ok {
		return nil, 0, fmt.Errorf("unexpected response type")
	}

	responseBody := unaryResp.Data
	statusCode := ctx.GetStatusCode()

	if len(responseBody) == 0 {
		return nil, statusCode, fmt.Errorf("empty response body")
	}

	logger.Debugf("[dubbo-go-pixiu] mcp server backend response: status=%d, size=%d bytes", statusCode, len(responseBody))
	return responseBody, statusCode, nil
}

// isAlreadyMCPFormat checks if the response is already in MCP JSON-RPC format
func (f *MCPServerFilter) isAlreadyMCPFormat(responseBody []byte) bool {
	return bytes.Contains(responseBody, []byte(`"jsonrpc":"2.0"`))
}

// processToolCallResponse processes the tool call response and sends the result
func (f *MCPServerFilter) processToolCallResponse(ctx *MCPContext, requestID any, responseBody []byte, statusCode int) filter.FilterStatus {
	// Check for backend errors
	if statusCode >= httpStatusClientErrorStart {
		logger.Errorf("[dubbo-go-pixiu] mcp server backend returned error status: %d", statusCode)
		return f.errorHandler.SendToolCallError(ctx, requestID, fmt.Sprintf("backend error: %d", statusCode))
	}

	// Build successful response using ToolCallSuccess method
	content := strings.TrimSpace(string(responseBody))
	mcpResponse := f.responseBuilder.ToolCallSuccess(requestID, content)
	return f.sendMCPResponse(ctx, mcpResponse)
}

// sendMCPResponse sends an MCP response and updates the target response
func (f *MCPServerFilter) sendMCPResponse(ctx *MCPContext, response mcp.JSONRPCResponse) filter.FilterStatus {
	mcpResponseBody, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to marshal MCP response: %v", err)
		return filter.Continue
	}

	// Override TargetResp to ensure MCP format response is sent
	ctx.TargetResp = &client.UnaryResponse{Data: mcpResponseBody}
	ctx.StatusCode(http.StatusOK)
	ctx.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueApplicationJson)

	// Critical: Clear Content-Length header to prevent mismatch errors
	ctx.Writer.Header().Del(constant.HeaderKeyContentLength)

	logger.Debugf("[dubbo-go-pixiu] mcp server successfully wrapped backend response in MCP format")
	return filter.Continue
}

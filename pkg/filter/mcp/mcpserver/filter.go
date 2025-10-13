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
	"io"
	"net/http"
)

import (
	"github.com/mark3labs/mcp-go/mcp"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// FilterFactory and MCPServerFilter types
type (
	// FilterFactory is a factory to create MCP server filters.
	FilterFactory struct {
		cfg      *Config
		registry *ToolRegistry
	}

	// MCPServerFilter is a filter that handles MCP protocol.
	MCPServerFilter struct {
		cfg             *Config
		registry        *ToolRegistry
		errorHandler    *ErrorHandler
		responseBuilder *ResponseBuilder
	}
)

// Apply prepares the MCP server and tool registry.
func (f *FilterFactory) Apply() error {
	// Initialize tool registry
	f.registry = NewToolRegistry()

	// Register statically configured tools
	for _, tool := range f.cfg.Tools {
		if err := f.registry.RegisterTool(tool); err != nil {
			return fmt.Errorf("failed to register tool %s: %v", tool.Name, err)
		}
		logger.Debugf("[dubbo-go-pixiu] mcp server registered tool '%s' -> cluster:%s", tool.Name, tool.Cluster)
	}

	// Register statically configured resources
	for _, resource := range f.cfg.Resources {
		if err := f.registry.RegisterResource(resource); err != nil {
			return fmt.Errorf("failed to register resource %s: %v", resource.Name, err)
		}
		logger.Debugf("[dubbo-go-pixiu] mcp server registered resource '%s' -> uri:%s", resource.Name, resource.URI)
	}

	// Register statically configured resource templates
	for _, template := range f.cfg.ResourceTemplates {
		if err := f.registry.RegisterResourceTemplate(template); err != nil {
			return fmt.Errorf("failed to register resource template %s: %v", template.Name, err)
		}
		logger.Debugf("[dubbo-go-pixiu] mcp server registered template '%s' -> pattern:%s", template.Name, template.URITemplate)
	}

	// Register statically configured prompts
	for _, prompt := range f.cfg.Prompts {
		if err := f.registry.RegisterPrompt(prompt); err != nil {
			return fmt.Errorf("failed to register prompt %s: %v", prompt.Name, err)
		}
		logger.Debugf("[dubbo-go-pixiu] mcp server registered prompt '%s'", prompt.Name)
	}

	return nil
}

// Config returns the configuration struct
func (f *FilterFactory) Config() any {
	return f.cfg
}

// PrepareFilterChain prepares the filter chain
func (f *FilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	mcpFilter := &MCPServerFilter{
		cfg:             f.cfg,
		registry:        f.registry,
		errorHandler:    NewErrorHandler(),
		responseBuilder: NewResponseBuilder(),
	}
	chain.AppendDecodeFilters(mcpFilter)
	chain.AppendEncodeFilters(mcpFilter) // Add to Encode chain
	return nil
}

// Decode processes incoming HTTP requests for MCP protocol.
func (f *MCPServerFilter) Decode(ctx *contexthttp.HttpContext) filter.FilterStatus {
	// Check if it's an MCP request
	if !f.isMCPRequest(ctx) {
		return filter.Continue
	}

	// Create MCP context wrapper
	mcpCtx := NewMCPContext(ctx)

	// Read request body
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to read request body: %v", err)
		return f.errorHandler.SendInternalError(mcpCtx, nil, "failed to read request body")
	}

	// Parse JSON-RPC request
	var jsonrpcReq mcp.JSONRPCRequest
	if err := json.Unmarshal(body, &jsonrpcReq); err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to parse JSON-RPC request: %v", err)
		return f.errorHandler.SendInternalError(mcpCtx, nil, "invalid JSON-RPC request")
	}

	logger.Infof("[dubbo-go-pixiu] mcp server received request: %s (id: %v)", jsonrpcReq.Method, jsonrpcReq.ID)

	// Store information in MCP context
	mcpCtx.SetMCPMethod(jsonrpcReq.Method)
	mcpCtx.SetMCPRequestID(jsonrpcReq.ID)

	// Handle terminal methods (methods that don't need forwarding to backend)
	if f.isTerminalMethod(jsonrpcReq.Method) {
		return f.handleTerminalMethod(mcpCtx, jsonrpcReq)
	} else if jsonrpcReq.Method == string(mcp.MethodToolsCall) {
		// Tool call will be processed in Encode stage (IsMCPToolCall() checks method)
		return f.handleToolCall(mcpCtx, jsonrpcReq)
	} else {
		// Unknown method
		logger.Warnf("[dubbo-go-pixiu] mcp server unsupported method: %s", jsonrpcReq.Method)
		return f.errorHandler.SendMethodNotFound(mcpCtx, jsonrpcReq.ID)
	}
}

// isTerminalMethod checks if it's a terminal method (methods that don't need forwarding to backend)
func (f *MCPServerFilter) isTerminalMethod(method string) bool {
	switch method {
	case string(mcp.MethodInitialize), string(mcp.MethodToolsList), string(mcp.MethodResourcesList), string(mcp.MethodResourcesRead),
		"resources/templates/list", string(mcp.MethodPromptsList), string(mcp.MethodPromptsGet),
		"notifications/initialized", string(mcp.MethodPing):
		return true
	default:
		return false
	}
}

// handleTerminalMethod handles terminal methods
func (f *MCPServerFilter) handleTerminalMethod(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	switch req.Method {
	case string(mcp.MethodInitialize):
		return f.handleInitialize(ctx, req)
	case string(mcp.MethodToolsList):
		return f.handleToolsList(ctx, req)
	case string(mcp.MethodResourcesList):
		return f.handleResourcesList(ctx, req)
	case string(mcp.MethodResourcesRead):
		return f.handleResourceRead(ctx, req)
	case "resources/templates/list":
		return f.handleResourceTemplatesList(ctx, req)
	case string(mcp.MethodPromptsList):
		return f.handlePromptsList(ctx, req)
	case string(mcp.MethodPromptsGet):
		return f.handlePromptsGet(ctx, req)
	case "notifications/initialized":
		return f.handleNotificationsInitialized(ctx, req)
	case string(mcp.MethodPing):
		return f.handlePing(ctx, req)
	default:
		logger.Warnf("[dubbo-go-pixiu] mcp server unsupported method: %s", req.Method)
		return f.errorHandler.SendMethodNotFound(ctx, req.ID)
	}
}

// Encode processes outgoing HTTP responses.
func (f *MCPServerFilter) Encode(ctx *contexthttp.HttpContext) filter.FilterStatus {
	// Create MCP context wrapper and load stored MCP data
	mcpCtx := NewMCPContextFromHttpContext(ctx)

	// Check if it's a tool call response
	if mcpCtx.IsMCPToolCall() {
		logger.Debugf("[dubbo-go-pixiu] mcp server processing tool call response: %s", ctx.Request.URL.Path)
		return f.handleToolCallResponse(mcpCtx)
	}

	// For regular MCP requests, no special processing needed
	if mcpCtx.IsMCPRequest() {
		logger.Debugf("[dubbo-go-pixiu] mcp server regular MCP request, no special processing needed")
	}

	return filter.Continue
}

// isMCPRequest checks if it's an MCP request
func (f *MCPServerFilter) isMCPRequest(ctx *contexthttp.HttpContext) bool {
	return ctx.Request.URL.Path == f.cfg.Endpoint
}

// sendJSONResponse sends a JSON response
func (f *MCPServerFilter) sendJSONResponse(ctx *MCPContext, response any) filter.FilterStatus {
	responseBody, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to marshal response: %v", err)
		errResp := contexthttp.InternalError.WithError(fmt.Errorf("marshal response failed: %w", err))
		ctx.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	// Get method and request ID for logging
	method := ctx.McpMethod()
	requestID := ctx.McpRequestID()

	logger.Infof("[dubbo-go-pixiu] mcp server response sent: %s (id: %v)", method, requestID)

	// Critical: Clear Content-Length header to prevent mismatch errors
	ctx.ClearContentLengthHeader()
	ctx.SendLocalReply(http.StatusOK, responseBody)
	return filter.Stop
}

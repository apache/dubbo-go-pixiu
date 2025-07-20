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
	"strings"

	"github.com/apache/dubbo-go-pixiu/pkg/model"

	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	h "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/creasty/defaults"
	"github.com/mark3labs/mcp-go/mcp"
)

// MCP method constants using mcp-go library constants
const (
	methodInitialize            = string(mcp.MethodInitialize)
	methodToolsList             = string(mcp.MethodToolsList)
	methodToolsCall             = string(mcp.MethodToolsCall)
	methodResourcesList         = string(mcp.MethodResourcesList)
	methodResourcesRead         = string(mcp.MethodResourcesRead)
	methodResourceTemplatesList = "resources/templates/list"
	methodPromptsList           = string(mcp.MethodPromptsList)
	methodPromptsGet            = string(mcp.MethodPromptsGet)
	methodNotificationsInit     = "notifications/initialized"
	methodPing                  = string(mcp.MethodPing)
)

// HTTP status code ranges
const (
	httpStatusClientErrorStart = 400
)

// MCP protocol constants
const (
	mcpProtocolVersion = "2024-11-05"
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
	// Set configuration default values
	if err := defaults.Set(f.cfg); err != nil {
		return fmt.Errorf("failed to set config defaults: %v", err)
	}

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

	toolCount, resourceCount, templateCount, promptCount := f.registry.Count()
	logger.Infof("[dubbo-go-pixiu] mcp server initialized successfully - tools:%d, resources:%d, templates:%d, prompts:%d",
		toolCount, resourceCount, templateCount, promptCount)

	return nil
}

// Config returns the configuration struct
func (f *FilterFactory) Config() any {
	return f.cfg
}

// PrepareFilterChain prepares the filter chain
func (f *FilterFactory) PrepareFilterChain(ctx *h.HttpContext, chain filter.FilterChain) error {
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
func (f *MCPServerFilter) Decode(ctx *h.HttpContext) filter.FilterStatus {
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
	} else if jsonrpcReq.Method == methodToolsCall {
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
	case methodInitialize, methodToolsList, methodResourcesList, methodResourcesRead,
		methodResourceTemplatesList, methodPromptsList, methodPromptsGet,
		methodNotificationsInit, methodPing:
		return true
	default:
		return false
	}
}

// handleTerminalMethod handles terminal methods
func (f *MCPServerFilter) handleTerminalMethod(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	switch req.Method {
	case methodInitialize:
		return f.handleInitialize(ctx, req)
	case methodToolsList:
		return f.handleToolsList(ctx, req)
	case methodResourcesList:
		return f.handleResourcesList(ctx, req)
	case methodResourcesRead:
		return f.handleResourceRead(ctx, req)
	case methodResourceTemplatesList:
		return f.handleResourceTemplatesList(ctx, req)
	case methodPromptsList:
		return f.handlePromptsList(ctx, req)
	case methodPromptsGet:
		return f.handlePromptsGet(ctx, req)
	case methodNotificationsInit:
		return f.handleNotificationsInitialized(ctx, req)
	case methodPing:
		return f.handlePing(ctx, req)
	default:
		logger.Warnf("[dubbo-go-pixiu] mcp server unsupported method: %s", req.Method)
		return f.errorHandler.SendMethodNotFound(ctx, req.ID)
	}
}

// Encode processes outgoing HTTP responses.
func (f *MCPServerFilter) Encode(ctx *h.HttpContext) filter.FilterStatus {
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
func (f *MCPServerFilter) isMCPRequest(ctx *h.HttpContext) bool {
	return ctx.Request.URL.Path == f.cfg.Endpoint
}

// sendJSONResponse sends a JSON response
func (f *MCPServerFilter) sendJSONResponse(ctx *MCPContext, response any) filter.FilterStatus {
	responseBody, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to marshal response: %v", err)
		ctx.SendLocalReply(http.StatusInternalServerError, []byte("internal server error"))
		return filter.Stop
	}

	// Get method and request ID for logging
	method := ctx.GetMCPMethod()
	requestID := ctx.GetMCPRequestID()

	logger.Infof("[dubbo-go-pixiu] mcp server response sent: %s (id: %v)", method, requestID)

	// Critical: Clear Content-Length header to prevent mismatch errors
	ctx.Writer.Header().Del("Content-Length")
	ctx.SendLocalReply(http.StatusOK, responseBody)
	return filter.Stop
}

// handleToolCall handles tool call requests by forwarding to backend
func (f *MCPServerFilter) handleToolCall(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	// Parse tool call parameters
	paramsBytes, err := json.Marshal(req.Params)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to marshal tool call params: %v", err)
		return f.errorHandler.SendInternalError(ctx, req.ID, "invalid tool call parameters")
	}

	var params struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments,omitempty"`
	}
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to parse tool call params: %v", err)
		return f.errorHandler.SendInvalidParams(ctx, req.ID, "invalid tool call parameters")
	}

	// Find tool configuration
	toolConfig, exists := f.registry.GetTool(params.Name)
	if !exists {
		logger.Warnf("[dubbo-go-pixiu] mcp server tool not found: %s", params.Name)
		return f.errorHandler.SendToolCallError(ctx, req.ID, fmt.Sprintf("tool not found: %s", params.Name))
	}

	// Build backend request
	err = f.buildBackendRequest(ctx, toolConfig, params.Arguments)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to build backend request: %v", err)
		return f.errorHandler.SendToolCallError(ctx, req.ID, "failed to build backend request")
	}

	// Set cluster information for routing
	if ctx.Params == nil {
		ctx.Params = make(map[string]any)
	}

	logger.Infof("[dubbo-go-pixiu] mcp server forwarding tool call: %s -> %s %s (cluster: %s)",
		params.Name, toolConfig.Request.Method, ctx.Request.URL.Path, toolConfig.Cluster)

	// Store MCP data for Encode stage processing
	ctx.StoreMCPDataInParams()

	ctx.Route = &model.RouteAction{
		Cluster: toolConfig.Cluster,
	}

	// Continue to next filter for backend forwarding
	return filter.Continue
}

// buildBackendRequest builds the complete backend request including path, body, and headers
func (f *MCPServerFilter) buildBackendRequest(ctx *MCPContext, toolConfig ToolConfig, arguments map[string]any) error {
	// Set HTTP method
	ctx.Request.Method = toolConfig.Request.Method

	// Build request path and body based on argument locations
	path := toolConfig.Request.Path
	bodyParams := make(map[string]any)
	queryParams := make(map[string]string)

	// Process arguments based on their location (path, query, body)
	if arguments != nil {
		for argName, argValue := range arguments {
			// Find argument configuration
			var argConfig *ArgConfig
			for _, arg := range toolConfig.Args {
				if arg.Name == argName {
					argConfig = &arg
					break
				}
			}

			if argConfig == nil {
				continue // Skip unknown arguments
			}

			switch argConfig.In {
			case "path":
				// Replace path parameters
				placeholder := fmt.Sprintf("{%s}", argName)
				replacement := fmt.Sprintf("%v", argValue)
				path = strings.ReplaceAll(path, placeholder, replacement)

			case "query":
				// Add to query parameters
				queryParams[argName] = fmt.Sprintf("%v", argValue)

			case "body":
				// Add to request body
				bodyParams[argName] = argValue
			}
		}
	}

	// Set the request path
	ctx.Request.URL.Path = path

	// Add query parameters
	if len(queryParams) > 0 {
		query := ctx.Request.URL.Query()
		for key, value := range queryParams {
			query.Set(key, value)
		}
		ctx.Request.URL.RawQuery = query.Encode()
	}

	// Build request body for POST/PUT requests
	if len(bodyParams) > 0 && (toolConfig.Request.Method == "POST" || toolConfig.Request.Method == "PUT") {
		bodyJSON, err := json.Marshal(bodyParams)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %v", err)
		}

		// Set request body
		ctx.Request.Body = io.NopCloser(strings.NewReader(string(bodyJSON)))
		ctx.Request.ContentLength = int64(len(bodyJSON))

		// Set Content-Type header
		ctx.Request.Header.Set("Content-Type", "application/json")

		logger.Debugf("[dubbo-go-pixiu] mcp server built request body: %s", string(bodyJSON))
	}

	return nil
}

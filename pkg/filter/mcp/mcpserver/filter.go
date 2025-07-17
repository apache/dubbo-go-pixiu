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
	"io"
	"net/http"
	"strings"

	"github.com/apache/dubbo-go-pixiu/pkg/model"

	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	h "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/creasty/defaults"
	"github.com/mark3labs/mcp-go/mcp"
)

// MCP method constants - using constants from mcp-go library
const (
	methodInitialize            = string(mcp.MethodInitialize)    // "initialize"
	methodToolsList             = string(mcp.MethodToolsList)     // "tools/list"
	methodToolsCall             = string(mcp.MethodToolsCall)     // "tools/call"
	methodResourcesList         = string(mcp.MethodResourcesList) // "resources/list"
	methodResourcesRead         = string(mcp.MethodResourcesRead) // "resources/read"
	methodResourceTemplatesList = "resources/templates/list"      // "resources/templates/list"
	methodPromptsList           = string(mcp.MethodPromptsList)   // "prompts/list"
	methodPromptsGet            = string(mcp.MethodPromptsGet)    // "prompts/get"
	methodNotificationsInit     = "notifications/initialized"     // initialization completion notification
	methodPing                  = string(mcp.MethodPing)          // connection keep-alive
)

type (
	// FilterFactory is the factory of mcp_server filter.
	FilterFactory struct {
		cfg      *Config
		registry *ToolRegistry
	}

	// MCPServerFilter is a filter that handles MCP protocol.
	MCPServerFilter struct {
		cfg      *Config
		registry *ToolRegistry
	}
)

// Apply prepares the MCP server and tool registry.
func (f *FilterFactory) Apply() error {
	// Set configuration default values
	if err := defaults.Set(f.cfg); err != nil {
		return fmt.Errorf("failed to set config defaults: %v", err)
	}

	f.registry = NewToolRegistry()

	// Note: method handling is now unified through handleTerminalMethod

	// Validate and register statically configured tools
	for i, tool := range f.cfg.Tools {
		// Set default values for each tool
		if err := defaults.Set(&tool); err != nil {
			return fmt.Errorf("failed to set defaults for tool %s: %v", tool.Name, err)
		}

		// Set default values for tool parameters
		for j := range tool.Args {
			if err := defaults.Set(&tool.Args[j]); err != nil {
				return fmt.Errorf("failed to set defaults for arg %s in tool %s: %v",
					tool.Args[j].Name, tool.Name, err)
			}
		}

		// Set default values for response configuration
		if tool.Response != nil {
			if err := defaults.Set(tool.Response); err != nil {
				return fmt.Errorf("failed to set defaults for response in tool %s: %v", tool.Name, err)
			}
		}

		// Update the tool in configuration (since we modified a copy)
		f.cfg.Tools[i] = tool

		// Validate tool configuration
		if err := tool.Validate(); err != nil {
			return fmt.Errorf("invalid tool configuration for %s: %v", tool.Name, err)
		}

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
		cfg:      f.cfg,
		registry: f.registry,
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
		return f.sendError(mcpCtx, nil, "failed to read request body")
	}

	// Parse JSON-RPC request
	var jsonrpcReq mcp.JSONRPCRequest
	if err := json.Unmarshal(body, &jsonrpcReq); err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to parse JSON-RPC request: %v", err)
		return f.sendError(mcpCtx, nil, "invalid JSON-RPC request")
	}

	logger.Infof("[dubbo-go-pixiu] mcp server: %s [id:%v]", jsonrpcReq.Method, jsonrpcReq.ID)

	// Store information in MCP context
	mcpCtx.SetMCPMethod(jsonrpcReq.Method)
	mcpCtx.SetMCPRequestID(jsonrpcReq.ID)

	// Store MCP data in HttpContext.Params for passing through the filter chain
	mcpCtx.StoreMCPDataInParams()

	// Distinguish between terminal and forwarding MCP methods
	if f.isTerminalMethod(jsonrpcReq.Method) {
		// Terminal methods: handle directly in Decode stage and return response
		return f.handleTerminalMethod(mcpCtx, jsonrpcReq)
	} else if jsonrpcReq.Method == methodToolsCall {
		// Tool calls: need to forward to backend service, wrap response in Encode stage
		return f.handleToolCall(mcpCtx, jsonrpcReq)
	} else {
		// Unknown method
		logger.Warnf("[dubbo-go-pixiu] mcp server unsupported method: %s", jsonrpcReq.Method)
		return f.sendMethodNotFound(mcpCtx, jsonrpcReq.ID)
	}
}

// isTerminalMethod checks if it's a terminal method (methods that don't need forwarding to backend)
func (f *MCPServerFilter) isTerminalMethod(method string) bool {
	switch method {
	case methodInitialize, methodToolsList, methodResourcesList, methodResourcesRead, methodResourceTemplatesList,
		methodPromptsList, methodPromptsGet, methodNotificationsInit, methodPing:
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
		return f.sendMethodNotFound(ctx, req.ID)
	}
}

// Encode processes outgoing HTTP responses.
func (f *MCPServerFilter) Encode(ctx *h.HttpContext) filter.FilterStatus {
	// Create MCP context wrapper and load data
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

// JSONRPCError JSON-RPC error structure
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// JSONRPCResponse generic JSON-RPC response structure
type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      any           `json:"id"`
	Result  any           `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

// NewJSONRPCResponse creates a success response
func NewJSONRPCResponse(id any, result any) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

// NewJSONRPCError creates an error response
func NewJSONRPCError(id any, code int, message string) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Error: &JSONRPCError{
			Code:    code,
			Message: message,
		},
	}
}

// sendError sends an error response
func (f *MCPServerFilter) sendError(ctx *MCPContext, id any, message string) filter.FilterStatus {
	response := NewJSONRPCError(id, -32603, message)
	return f.sendJSONResponse(ctx, response)
}

// sendMethodNotFound sends a method not found error
func (f *MCPServerFilter) sendMethodNotFound(ctx *MCPContext, id any) filter.FilterStatus {
	response := NewJSONRPCError(id, -32601, "Method not found")
	return f.sendJSONResponse(ctx, response)
}

// sendJSONResponse sends a JSON response
func (f *MCPServerFilter) sendJSONResponse(ctx *MCPContext, response any) filter.FilterStatus {
	responseBody, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to marshal response: %v", err)
		ctx.SendLocalReply(http.StatusInternalServerError, []byte("internal server error"))
		return filter.Stop
	}

	// Get method name and request ID for logging
	method := ctx.GetMCPMethod()
	requestID := ctx.GetMCPRequestID()

	// Analyze response content type
	resultType := f.getResponseType(response)

	logger.Infof("[dubbo-go-pixiu] mcp server: %s [id:%v] -> %s", method, requestID, resultType)
	ctx.SendLocalReply(http.StatusOK, responseBody)
	return filter.Stop
}

// getResponseType analyzes response type for log display
func (f *MCPServerFilter) getResponseType(response any) string {
	// Check if it's an initialization response
	if _, ok := response.(InitializeResponse); ok {
		// Show total capability type count (tools, resources, prompts = 3)
		return "initialize(caps:3)"
	}

	// Check if it's an error response
	if resp, ok := response.(*JSONRPCResponse); ok {
		if resp.Error != nil {
			return "error"
		}
		if resp.Result != nil {
			if resultMap, ok := resp.Result.(map[string]any); ok {
				// Check if it's a tools list
				if tools, exists := resultMap["tools"]; exists {
					if toolsArray, ok := tools.([]any); ok {
						return fmt.Sprintf("tools(%d)", len(toolsArray))
					}
				}
				// Check if it's a resources list
				if resources, exists := resultMap["resources"]; exists {
					if resourcesArray, ok := resources.([]any); ok {
						return fmt.Sprintf("resources(%d)", len(resourcesArray))
					}
				}
				// Check if it's a prompts list
				if prompts, exists := resultMap["prompts"]; exists {
					if promptsArray, ok := prompts.([]any); ok {
						return fmt.Sprintf("prompts(%d)", len(promptsArray))
					}
				}
				// Check if it's a tool call result
				if content, exists := resultMap["content"]; exists {
					if isError, errorExists := resultMap["isError"]; errorExists {
						if isErr, ok := isError.(bool); ok && isErr {
							return "tool_error"
						}
					}
					if contentArray, ok := content.([]any); ok {
						return fmt.Sprintf("tool_result(%d)", len(contentArray))
					}
				}
				// Check if it's an initialization response
				if capabilities, exists := resultMap["capabilities"]; exists {
					if capMap, ok := capabilities.(map[string]any); ok {
						return fmt.Sprintf("initialize(caps:%d)", len(capMap))
					}
				}
				// Empty result object (like ping)
				if len(resultMap) == 0 {
					return "empty"
				}
				return fmt.Sprintf("object(%d)", len(resultMap))
			}
		}
	}
	return "unknown"
}

// buildRequestPath builds request path, replacing path parameters
func (f *MCPServerFilter) buildRequestPath(pathTemplate string, arguments map[string]any) (string, error) {
	path := pathTemplate

	// Extract all path parameters
	pathParams := GetPathParameterNames(pathTemplate)

	// Replace each placeholder
	for _, paramName := range pathParams {
		placeholder := "{" + paramName + "}"

		// Get parameter value
		value, exists := arguments[paramName]
		if !exists {
			return "", fmt.Errorf("required path parameter '%s' is missing", paramName)
		}

		// Replace placeholder
		var strValue string
		if str, ok := value.(string); ok {
			strValue = str
		} else {
			strValue = fmt.Sprintf("%v", value)
		}

		path = strings.Replace(path, placeholder, strValue, -1)
	}

	// Check if there are still unresolved placeholders
	if strings.Contains(path, "{") && strings.Contains(path, "}") {
		return "", fmt.Errorf("unresolved placeholders in path: %s", path)
	}

	return path, nil
}

// InitializeResponse initialization response structure
type InitializeResponse struct {
	JSONRPC string `json:"jsonrpc"`
	ID      any    `json:"id"`
	Result  struct {
		ProtocolVersion string `json:"protocolVersion"`
		Capabilities    struct {
			Tools struct {
				ListChanged bool `json:"listChanged"`
			} `json:"tools"`
			Resources struct {
				ListChanged bool `json:"listChanged"`
			} `json:"resources"`
			Prompts struct {
				ListChanged bool `json:"listChanged"`
			} `json:"prompts"`
		} `json:"capabilities"`
		ServerInfo struct {
			Name    string `json:"name"`
			Version string `json:"version"`
		} `json:"serverInfo"`
		Instructions string `json:"instructions,omitempty"`
	} `json:"result"`
}

// handleInitialize handles the initialize method
func (f *MCPServerFilter) handleInitialize(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	// Create initialization response, must include capabilities field
	response := InitializeResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
	}

	response.Result.ProtocolVersion = "2024-11-05"

	// TODO: Dynamic update capabilities - enable after Nacos integration
	// Currently set to false, future Nacos integration will support:
	// 1. Dynamic discovery and registration of new backend services
	// 2. Automatic generation of corresponding MCP tools
	// 3. Send notifications/tools/list_changed notifications
	// 4. Send notifications/resources/list_changed notifications
	// 5. Send notifications/prompts/list_changed notifications
	// Reference: https://nacos.io/docs/latest/manual/user/ai/api-to-mcp/
	response.Result.Capabilities.Tools.ListChanged = false
	response.Result.Capabilities.Resources.ListChanged = false
	response.Result.Capabilities.Prompts.ListChanged = false

	response.Result.ServerInfo.Name = f.cfg.ServerInfo.Name
	response.Result.ServerInfo.Version = f.cfg.ServerInfo.Version

	// If there are instructions in the configuration, add them to the response
	if f.cfg.ServerInfo.Instructions != "" {
		response.Result.Instructions = f.cfg.ServerInfo.Instructions
	}

	return f.sendJSONResponse(ctx, response)
}

// handleToolsList handles the tools/list method
func (f *MCPServerFilter) handleToolsList(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	// Get all tools
	mcpTools, err := f.registry.ToMCPTools()
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to get MCP tools: %v", err)
		return f.sendError(ctx, req.ID, "failed to get tools")
	}

	// Build tools list response
	result := map[string]any{
		"tools": mcpTools,
	}

	response := NewJSONRPCResponse(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// handleResourcesList handles the resources/list method
func (f *MCPServerFilter) handleResourcesList(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	// Get all resources
	mcpResources, err := f.registry.ToMCPResources()
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to get MCP resources: %v", err)
		return f.sendError(ctx, req.ID, "failed to get resources")
	}

	// Build resources list response
	result := map[string]any{
		"resources": mcpResources,
	}

	response := NewJSONRPCResponse(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// handleResourceRead handles the resources/read method
func (f *MCPServerFilter) handleResourceRead(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Debugf("[dubbo-go-pixiu] mcp server handling resources/read")

	// Parse request parameters
	paramsBytes, err := json.Marshal(req.Params)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to marshal params: %v", err)
		return f.sendError(ctx, req.ID, "invalid parameters")
	}

	var params struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to parse resource read params: %v", err)
		return f.sendError(ctx, req.ID, "invalid parameters")
	}

	// Find resource (by URI)
	resource, exists := f.registry.GetResourceByURI(params.URI)
	if !exists {
		logger.Warnf("[dubbo-go-pixiu] mcp server resource not found: %s", params.URI)
		return f.sendError(ctx, req.ID, fmt.Sprintf("resource not found: %s", params.URI))
	}

	// Read resource content (simplified handling here, should read according to source configuration)
	var content string
	if resource.Source.Type == "static" {
		if staticContent, ok := resource.Source.Config["content"].(string); ok {
			content = staticContent
		}
	}

	result := map[string]any{
		"contents": []map[string]any{
			{
				"uri":      resource.URI,
				"mimeType": resource.MIMEType,
				"text":     content,
			},
		},
	}

	response := NewJSONRPCResponse(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// handleResourceTemplatesList handles the resources/templates/list method
func (f *MCPServerFilter) handleResourceTemplatesList(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	// Get all resource templates (parameterized resource patterns)
	mcpTemplates, err := f.registry.ToMCPResourceTemplates()
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to get MCP resource templates: %v", err)
		return f.sendError(ctx, req.ID, "failed to get resource templates")
	}

	result := map[string]any{
		"resourceTemplates": mcpTemplates,
	}

	response := NewJSONRPCResponse(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// handleToolCall handles the tools/call method
func (f *MCPServerFilter) handleToolCall(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {

	// Parse request parameters
	paramsBytes, err := json.Marshal(req.Params)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to marshal params: %v", err)
		return f.sendError(ctx, req.ID, "invalid parameters")
	}

	var params struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to parse tool call params: %v", err)
		return f.sendError(ctx, req.ID, "invalid parameters")
	}

	// Find tool configuration
	toolConfig, exists := f.registry.GetTool(params.Name)
	if !exists {
		logger.Warnf("[dubbo-go-pixiu] mcp server tool not found: %s", params.Name)
		return f.sendError(ctx, req.ID, fmt.Sprintf("tool not found: %s", params.Name))
	}

	// Process parameters and build request
	err = f.processParameters(ctx, toolConfig, params.Arguments)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to process parameters: %v", err)
		return f.sendError(ctx, req.ID, "failed to process parameters")
	}

	// Set HttpContext routing information (critical!)
	ctx.Request.Method = toolConfig.Request.Method

	// Set request headers
	for key, value := range toolConfig.Request.Headers {
		ctx.Request.Header.Set(key, value)
	}

	// Important: set correct routing information so HttpProxy Filter knows the target cluster
	routeAction := &model.RouteAction{
		Cluster: toolConfig.Cluster,
	}
	ctx.RouteEntry(routeAction)

	// Store MCP data in HttpContext.Params for passing through the filter chain
	ctx.StoreMCPDataInParams()

	// Continue to next filter (httpproxy)
	return filter.Continue
}

// processParameters processes request parameters according to parameter configuration
func (f *MCPServerFilter) processParameters(ctx *MCPContext, toolConfig ToolConfig, arguments map[string]any) error {
	// Get all parameters (including auto-inferred path parameters)
	allParams, err := toolConfig.GetAllParameters()
	if err != nil {
		return fmt.Errorf("failed to get parameters: %v", err)
	}

	// Validate required parameters and apply default values
	processedArgs := make(map[string]any)
	for key, value := range arguments {
		processedArgs[key] = value
	}

	// Collect request body parameters
	bodyParams := make(map[string]any)

	// Process each parameter
	for _, param := range allParams {
		value := processedArgs[param.Name]

		// Check required parameters
		if value == nil {
			if param.Default != nil {
				processedArgs[param.Name] = param.Default
				value = param.Default
			} else if param.Required {
				return fmt.Errorf("required parameter %s is missing", param.Name)
			} else {
				continue // Skip optional parameters without default values
			}
		}

		// Process according to parameter location
		switch param.In {
		case "path":
			// Path parameters are handled through path building
			continue
		case "query":
			// Add to query parameters
			query := ctx.Request.URL.Query()
			query.Set(param.Name, fmt.Sprintf("%v", value))
			ctx.Request.URL.RawQuery = query.Encode()
		case "header":
			// Set request header
			ctx.Request.Header.Set(param.Name, fmt.Sprintf("%v", value))
		case "body":
			// Collect request body parameters
			bodyParams[param.Name] = value
		default:
			logger.Warnf("[dubbo-go-pixiu] mcp server unsupported parameter location: %s for parameter %s", param.In, param.Name)
		}
	}

	// Process request body parameters
	if len(bodyParams) > 0 {
		if err := f.setRequestBody(ctx, bodyParams); err != nil {
			return fmt.Errorf("failed to set request body: %v", err)
		}
	}

	// Build path (handle path parameters)
	path, err := f.buildRequestPath(toolConfig.Request.Path, processedArgs)
	if err != nil {
		return fmt.Errorf("failed to build path: %v", err)
	}
	ctx.Request.URL.Path = path

	return nil
}

// setRequestBody sets request body parameters
func (f *MCPServerFilter) setRequestBody(ctx *MCPContext, bodyParams map[string]any) error {
	// Serialize request body parameters to JSON
	bodyData, err := json.Marshal(bodyParams)
	if err != nil {
		return fmt.Errorf("failed to marshal body parameters: %v", err)
	}

	// Set request body
	ctx.Request.Body = io.NopCloser(bytes.NewReader(bodyData))
	ctx.Request.ContentLength = int64(len(bodyData))

	// Set Content-Type header
	if ctx.Request.Header.Get("Content-Type") == "" {
		ctx.Request.Header.Set("Content-Type", "application/json")
	}

	return nil
}

// handlePromptsList handles the prompts/list method
func (f *MCPServerFilter) handlePromptsList(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Debugf("[dubbo-go-pixiu] mcp server handling prompts/list request")

	// Get all prompts
	mcpPrompts, err := f.registry.ToMCPPrompts()
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to get prompts: %v", err)
		return f.sendError(ctx, req.ID, "failed to get prompts")
	}

	result := map[string]any{
		"prompts": mcpPrompts,
	}

	response := NewJSONRPCResponse(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// handlePromptsGet handles the prompts/get method
func (f *MCPServerFilter) handlePromptsGet(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Debugf("[dubbo-go-pixiu] mcp server handling prompts/get request")

	// Parse request parameters
	paramsBytes, err := json.Marshal(req.Params)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to marshal params: %v", err)
		return f.sendError(ctx, req.ID, "invalid parameters")
	}

	var params struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments,omitempty"`
	}
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to parse prompts/get params: %v", err)
		return f.sendError(ctx, req.ID, "invalid parameters")
	}

	// Find prompt configuration
	promptConfig, exists := f.registry.GetPrompt(params.Name)
	if !exists {
		logger.Warnf("[dubbo-go-pixiu] mcp server prompt not found: %s", params.Name)
		return f.sendError(ctx, req.ID, fmt.Sprintf("prompt not found: %s", params.Name))
	}

	// Build prompt messages
	messages, err := f.buildPromptMessages(promptConfig, params.Arguments)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to build prompt messages: %v", err)
		return f.sendError(ctx, req.ID, "failed to build prompt messages")
	}

	result := map[string]any{
		"description": promptConfig.Description,
		"messages":    messages,
	}

	response := NewJSONRPCResponse(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// buildPromptMessages builds prompt messages with parameter replacement support
func (f *MCPServerFilter) buildPromptMessages(promptConfig PromptConfig, arguments map[string]any) ([]map[string]any, error) {
	messages := make([]map[string]any, 0, len(promptConfig.Messages))

	for _, msgConfig := range promptConfig.Messages {
		// Replace parameters in message content
		content := f.replacePromptArguments(msgConfig.Content, arguments)

		message := map[string]any{
			"role": msgConfig.Role,
			"content": map[string]any{
				"type": "text",
				"text": content,
			},
		}

		messages = append(messages, message)
	}

	return messages, nil
}

// replacePromptArguments replaces parameter placeholders in prompt content
func (f *MCPServerFilter) replacePromptArguments(content string, arguments map[string]any) string {
	if arguments == nil {
		return content
	}

	result := content
	for key, value := range arguments {
		placeholder := fmt.Sprintf("{{%s}}", key)
		replacement := fmt.Sprintf("%v", value)
		result = strings.ReplaceAll(result, placeholder, replacement)
	}

	return result
}

// handleNotificationsInitialized handles notifications/initialized notification
func (f *MCPServerFilter) handleNotificationsInitialized(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Debugf("[dubbo-go-pixiu] mcp server received initialized notification from client")

	// notifications/initialized is a notification, no response needed
	// According to MCP protocol, notifications should not have responses
	// We just need to log this event, indicating the client is ready for normal operations

	logger.Infof("[dubbo-go-pixiu] mcp server client initialization completed, ready for normal operations")

	// For notifications, we don't send a response, just return Stop
	return filter.Stop
}

// handlePing handles the ping method
func (f *MCPServerFilter) handlePing(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Debugf("[dubbo-go-pixiu] mcp server handling ping request")

	// ping method is usually used for connection keep-alive, return a simple pong response
	result := map[string]any{} // Empty result object

	response := NewJSONRPCResponse(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// handleToolCallResponse handles tool call responses, wrapping backend responses in MCP format
func (f *MCPServerFilter) handleToolCallResponse(ctx *MCPContext) filter.FilterStatus {
	logger.Debugf("[dubbo-go-pixiu] mcp server handling tool call response")

	// Get original JSON-RPC request ID
	requestID := ctx.GetMCPRequestID()
	if requestID == nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server missing request ID for tool call response")
		return filter.Continue
	}

	logger.Debugf("[dubbo-go-pixiu] mcp server processing response for request ID: %v", requestID)

	// Get backend response (from TargetResp get processed data)
	var responseBody []byte
	var statusCode int

	// Get processed data from TargetResp by buildTargetResponse
	if ctx.TargetResp != nil {
		if unaryResp, ok := ctx.TargetResp.(*client.UnaryResponse); ok {
			responseBody = unaryResp.Data
			statusCode = ctx.GetStatusCode()

			// Check if response is already in MCP format (prevent duplicate wrapping)
			if bytes.Contains(responseBody, []byte(`"jsonrpc":"2.0"`)) {
				logger.Debugf("[dubbo-go-pixiu] mcp server response already in MCP format, skipping processing")
				return filter.Continue
			}

			logger.Debugf("[dubbo-go-pixiu] mcp server backend response: status=%d, size=%d bytes", statusCode, len(responseBody))
		} else {
			logger.Errorf("[dubbo-go-pixiu] mcp server unexpected TargetResp type: %T", ctx.TargetResp)
			return f.sendToolCallError(ctx, requestID, "unexpected response type")
		}
	} else {
		logger.Errorf("[dubbo-go-pixiu] mcp server no TargetResp available")
		return f.sendToolCallError(ctx, requestID, "no response data available")
	}

	// Check if response is empty
	if len(responseBody) == 0 {
		logger.Errorf("[dubbo-go-pixiu] mcp server empty response body from backend")
		return f.sendToolCallError(ctx, requestID, "empty response from backend")
	}

	// Check HTTP status code
	if statusCode >= 400 {
		logger.Errorf("[dubbo-go-pixiu] mcp server backend returned error status: %d", statusCode)
		return f.sendToolCallError(ctx, requestID, fmt.Sprintf("backend error: %d", statusCode))
	}

	// Build MCP tool call response
	result := f.buildToolCallResult(responseBody, statusCode)

	// Create JSON-RPC response
	mcpResponse := NewJSONRPCResponse(requestID, result)

	// Serialize MCP response to JSON
	mcpResponseBody, err := json.Marshal(mcpResponse)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to marshal MCP response: %v", err)
		return f.sendToolCallError(ctx, requestID, "failed to marshal response")
	}

	// Override TargetResp to ensure MCP format response is sent
	ctx.TargetResp = &client.UnaryResponse{Data: mcpResponseBody}
	ctx.StatusCode(http.StatusOK)
	ctx.AddHeader("Content-Type", "application/json")

	// Clear Content-Length header, let HTTP library auto-calculate new length
	ctx.Writer.Header().Del("Content-Length")

	logger.Debugf("[dubbo-go-pixiu] mcp server response prepared: size=%d bytes", len(mcpResponseBody))
	return filter.Continue
}

// buildToolCallResult builds tool call result
func (f *MCPServerFilter) buildToolCallResult(responseBody []byte, statusCode int) map[string]any {
	// Clean response body, remove trailing newlines and whitespace
	cleanedBody := strings.TrimSpace(string(responseBody))

	// Try to parse response as JSON
	var jsonData any
	if err := json.Unmarshal(responseBody, &jsonData); err == nil {
		// If it's valid JSON, return structured content
		return map[string]any{
			"content": []map[string]any{
				{
					"type": "text",
					"text": cleanedBody,
				},
			},
			"isError": false,
		}
	} else {
		// If it's not valid JSON, treat as plain text
		return map[string]any{
			"content": []map[string]any{
				{
					"type": "text",
					"text": cleanedBody,
				},
			},
			"isError": false,
		}
	}
}

// sendToolCallError sends tool call error response
func (f *MCPServerFilter) sendToolCallError(ctx *MCPContext, requestID any, message string) filter.FilterStatus {
	// Use MCP tool call result format
	result := map[string]any{
		"isError": true,
		"content": []map[string]any{
			{
				"type": "text",
				"text": fmt.Sprintf("Error: %s", message),
			},
		},
	}

	response := NewJSONRPCResponse(requestID, result)
	return f.sendJSONResponse(ctx, response)
}

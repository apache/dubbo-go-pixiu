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
)

import (
	"github.com/mark3labs/mcp-go/mcp"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	inPath  = "path"
	inQuery = "query"
	inBody  = "body"

	typeString  = "string"
	typeInteger = "integer"
	typeNumber  = "number"
	typeBoolean = "boolean"
)

// handleInitialize handles the initialize method
func (f *MCPServerFilter) handleInitialize(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	// Build server capabilities using mcp-go structures
	capabilities := mcp.ServerCapabilities{
		Tools: &struct {
			ListChanged bool `json:"listChanged,omitempty"`
		}{
			// TODO: Dynamic update capabilities - enable after Nacos integration
			// Currently set to false, future Nacos integration will support:
			// 1. Dynamic discovery and registration of new backend services
			// 2. Automatic generation of corresponding MCP tools
			// 3. Send notifications/tools/list_changed notifications
			ListChanged: false,
		},
		Resources: &struct {
			Subscribe   bool `json:"subscribe,omitempty"`
			ListChanged bool `json:"listChanged,omitempty"`
		}{
			Subscribe:   false,
			ListChanged: false,
		},
		Prompts: &struct {
			ListChanged bool `json:"listChanged,omitempty"`
		}{
			ListChanged: false,
		},
	}

	// Build server info using mcp-go structures
	serverInfo := mcp.Implementation{
		Name:    f.cfg.ServerInfo.Name,
		Version: f.cfg.ServerInfo.Version,
	}

	// Create initialization result using mcp-go API
	instructions := f.cfg.ServerInfo.Instructions
	if instructions == "" {
		instructions = "This MCP server provides API access through tools, documentation through resources, and AI assistance through prompts."
	}
	result := mcp.NewInitializeResult(mcp.LATEST_PROTOCOL_VERSION, capabilities, serverInfo, instructions)

	// Create JSON-RPC response
	response := f.responseBuilder.Success(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// handleToolsList handles the tools/list method using mcp-go APIs
func (f *MCPServerFilter) handleToolsList(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	tools := make([]mcp.Tool, 0, len(f.cfg.Tools))

	// Build tools using mcp-go API for standard compliance
	for _, toolCfg := range f.cfg.Tools {
		// Start with basic tool options
		toolOptions := []mcp.ToolOption{
			mcp.WithDescription(toolCfg.Description),
		}

		// Add parameter definitions using mcp-go APIs
		for _, arg := range toolCfg.Args {
			opts := f.buildToolParameterOptions(&arg)
			switch arg.Type {
			case typeString:
				toolOptions = append(toolOptions, mcp.WithString(arg.Name, opts...))
			case typeInteger, typeNumber:
				toolOptions = append(toolOptions, mcp.WithNumber(arg.Name, opts...))
			case typeBoolean:
				toolOptions = append(toolOptions, mcp.WithBoolean(arg.Name, opts...))
			}
		}

		// Create tool using mcp-go API
		tool := mcp.NewTool(toolCfg.Name, toolOptions...)
		tools = append(tools, tool)
	}

	// Build standard MCP tools list response using mcp-go structures
	result := mcp.NewListToolsResult(tools, "") // empty cursor for no pagination

	response := f.responseBuilder.Success(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// buildToolParameterOptions builds the mcp.PropertyOption slice for a given tool argument
func (f *MCPServerFilter) buildToolParameterOptions(arg *ArgConfig) []mcp.PropertyOption {
	opts := []mcp.PropertyOption{mcp.Description(arg.Description)}

	if arg.Required {
		opts = append(opts, mcp.Required())
	}

	if arg.Default != nil {
		switch arg.Type {
		case typeString:
			if defaultStr, ok := arg.Default.(string); ok {
				opts = append(opts, mcp.DefaultString(defaultStr))
			}
		case typeInteger, typeNumber:
			switch defaultVal := arg.Default.(type) {
			case float64:
				opts = append(opts, mcp.DefaultNumber(defaultVal))
			case int:
				opts = append(opts, mcp.DefaultNumber(float64(defaultVal)))
			case int64:
				opts = append(opts, mcp.DefaultNumber(float64(defaultVal)))
			}
		case typeBoolean:
			if defaultBool, ok := arg.Default.(bool); ok {
				opts = append(opts, mcp.DefaultBool(defaultBool))
			}
		}
	}

	if len(arg.Enum) > 0 && arg.Type == typeString {
		opts = append(opts, mcp.Enum(arg.Enum...))
	}

	return opts
}

// handleResourcesList handles the resources/list method
func (f *MCPServerFilter) handleResourcesList(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	// Get all resources
	mcpResources, err := f.registry.ToMCPResources()
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to get MCP resources: %v", err)
		return f.errorHandler.SendInternalError(ctx, req.ID, "failed to get resources")
	}

	// Build resources list response using mcp-go structures
	result := mcp.NewListResourcesResult(mcpResources, "") // empty cursor for no pagination

	response := f.responseBuilder.Success(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// handlePing handles the ping method
func (f *MCPServerFilter) handlePing(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Debugf("[dubbo-go-pixiu] mcp server handling ping request")

	// Simple ping response
	result := map[string]any{}

	response := f.responseBuilder.Success(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// handleNotificationsInitialized handles notifications/initialized notification
func (f *MCPServerFilter) handleNotificationsInitialized(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Debugf("[dubbo-go-pixiu] mcp server received initialized notification from client")

	// Store client initialization state
	// This notification indicates that the client has completed initialization
	// and is ready to receive requests

	// For notifications, we don't send a response, just return Stop
	return filter.Stop
}

// handleResourceRead handles the resources/read method
func (f *MCPServerFilter) handleResourceRead(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Debugf("[dubbo-go-pixiu] mcp server handling resources/read")

	// Parse request parameters
	paramsBytes, err := json.Marshal(req.Params)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to marshal params: %v", err)
		return f.errorHandler.SendInvalidParams(ctx, req.ID, "invalid parameters")
	}

	var params struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to parse resource read params: %v", err)
		return f.errorHandler.SendInvalidParams(ctx, req.ID, "invalid parameters")
	}

	// Find resource (by URI)
	resource, exists := f.registry.GetResourceByURI(params.URI)
	if !exists {
		logger.Warnf("[dubbo-go-pixiu] mcp server resource not found: %s", params.URI)
		return f.errorHandler.SendInternalError(ctx, req.ID, fmt.Sprintf("resource not found: %s", params.URI))
	}

	// Build resource content response
	// TODO: Implement actual resource content loading from source
	content := fmt.Sprintf("Resource content for %s (source: %s)", resource.URI, resource.Source.Type)

	result := map[string]any{
		"contents": []map[string]any{
			{
				"uri":      resource.URI,
				"mimeType": resource.MIMEType,
				"text":     content,
			},
		},
	}

	response := f.responseBuilder.Success(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// handleResourceTemplatesList handles the resources/templates/list method
func (f *MCPServerFilter) handleResourceTemplatesList(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	// Get all resource templates (parameterized resource patterns)
	mcpResourceTemplates, err := f.registry.ToMCPResourceTemplates()
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to get MCP resource templates: %v", err)
		return f.errorHandler.SendInternalError(ctx, req.ID, "failed to get resource templates")
	}

	// Build resource templates list response
	result := map[string]any{
		"resourceTemplates": mcpResourceTemplates,
	}

	response := f.responseBuilder.Success(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// handlePromptsList handles the prompts/list method
func (f *MCPServerFilter) handlePromptsList(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Debugf("[dubbo-go-pixiu] mcp server handling prompts/list request")

	// Get all prompts
	mcpPrompts, err := f.registry.ToMCPPrompts()
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to get prompts: %v", err)
		return f.errorHandler.SendInternalError(ctx, req.ID, "failed to get prompts")
	}

	// Build prompts list response
	result := map[string]any{
		"prompts": mcpPrompts,
	}

	response := f.responseBuilder.Success(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// handlePromptsGet handles the prompts/get method
func (f *MCPServerFilter) handlePromptsGet(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Debugf("[dubbo-go-pixiu] mcp server handling prompts/get request")

	// Parse request parameters
	paramsBytes, err := json.Marshal(req.Params)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to marshal params: %v", err)
		return f.errorHandler.SendInvalidParams(ctx, req.ID, "invalid parameters")
	}

	var params struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments,omitempty"`
	}
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to parse prompts/get params: %v", err)
		return f.errorHandler.SendInvalidParams(ctx, req.ID, "invalid parameters")
	}

	// Find prompt configuration
	promptConfig, exists := f.registry.GetPrompt(params.Name)
	if !exists {
		logger.Warnf("[dubbo-go-pixiu] mcp server prompt not found: %s", params.Name)
		return f.errorHandler.SendInternalError(ctx, req.ID, fmt.Sprintf("prompt not found: %s", params.Name))
	}

	// Build prompt messages with parameter replacement
	messages, err := f.buildPromptMessages(promptConfig, params.Arguments)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to build prompt messages: %v", err)
		return f.errorHandler.SendInternalError(ctx, req.ID, "failed to build prompt messages")
	}

	// Build prompts/get response
	result := map[string]any{
		"description": promptConfig.Description,
		"messages":    messages,
	}

	response := f.responseBuilder.Success(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// buildPromptMessages builds prompt messages with parameter replacement support
func (f *MCPServerFilter) buildPromptMessages(promptConfig PromptConfig, arguments map[string]any) ([]map[string]any, error) {
	messages := make([]map[string]any, 0, len(promptConfig.Messages))

	for _, msg := range promptConfig.Messages {
		// Replace parameter placeholders in content
		content := f.replacePromptArguments(msg.Content, arguments)

		message := map[string]any{
			"role":    msg.Role,
			"content": content,
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
			case inPath:
				// Replace path parameters
				placeholder := fmt.Sprintf("{%s}", argName)
				replacement := fmt.Sprintf("%v", argValue)
				path = strings.ReplaceAll(path, placeholder, replacement)

			case inQuery:
				// Add to query parameters
				queryParams[argName] = fmt.Sprintf("%v", argValue)

			case inBody:
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
	if len(bodyParams) > 0 && (toolConfig.Request.Method == constant.Post || toolConfig.Request.Method == constant.Put) {
		bodyJSON, err := json.Marshal(bodyParams)
		if err != nil {
			return fmt.Errorf("failed to marshal request body: %v", err)
		}

		// Set request body
		ctx.Request.Body = io.NopCloser(strings.NewReader(string(bodyJSON)))

		// Set Content-Type header
		ctx.Request.Header.Set(constant.HeaderKeyContextType, constant.HeaderValueApplicationJson)

		logger.Debugf("[dubbo-go-pixiu] mcp server built request body: %s", string(bodyJSON))
	}

	return nil
}

// handleToolCallResponse handles tool call responses, wrapping backend responses in MCP format
func (f *MCPServerFilter) handleToolCallResponse(ctx *MCPContext) filter.FilterStatus {
	logger.Debugf("[dubbo-go-pixiu] mcp server handling tool call response")

	// Extract request information
	requestID := ctx.McpRequestID()
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

// processToolCallResponse processes the tool call response and sends the result
func (f *MCPServerFilter) processToolCallResponse(ctx *MCPContext, requestID any, responseBody []byte, statusCode int) filter.FilterStatus {
	// Check for backend errors
	if statusCode >= 400 {
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
	ctx.ClearContentLengthHeader()

	logger.Debugf("[dubbo-go-pixiu] mcp server successfully wrapped backend response in MCP format")
	return filter.Continue
}

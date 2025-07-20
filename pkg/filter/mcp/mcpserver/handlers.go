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
	"strings"

	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/mark3labs/mcp-go/mcp"
)

// Note: Using mcp-go library's standard InitializeResult structure instead of custom InitializeResponse

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
	result := mcp.NewInitializeResult(mcpProtocolVersion, capabilities, serverInfo, instructions)

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
			switch arg.Type {
			case "string":
				opts := []mcp.PropertyOption{mcp.Description(arg.Description)}
				if arg.Required {
					opts = append(opts, mcp.Required())
				}
				if arg.Default != nil {
					if defaultStr, ok := arg.Default.(string); ok {
						opts = append(opts, mcp.DefaultString(defaultStr))
					}
				}
				if len(arg.Enum) > 0 {
					opts = append(opts, mcp.Enum(arg.Enum...))
				}
				toolOptions = append(toolOptions, mcp.WithString(arg.Name, opts...))

			case "integer", "number":
				opts := []mcp.PropertyOption{mcp.Description(arg.Description)}
				if arg.Required {
					opts = append(opts, mcp.Required())
				}
				if arg.Default != nil {
					if defaultNum, ok := arg.Default.(float64); ok {
						opts = append(opts, mcp.DefaultNumber(defaultNum))
					}
				}
				toolOptions = append(toolOptions, mcp.WithNumber(arg.Name, opts...))

			case "boolean":
				opts := []mcp.PropertyOption{mcp.Description(arg.Description)}
				if arg.Required {
					opts = append(opts, mcp.Required())
				}
				if arg.Default != nil {
					if defaultBool, ok := arg.Default.(bool); ok {
						opts = append(opts, mcp.DefaultBool(defaultBool))
					}
				}
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

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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"sort"
	"strings"
)

import (
	"github.com/mark3labs/mcp-go/mcp"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/router"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/transport"
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

	toolsListChangedMethod = "notifications/tools/list_changed"
)

type initializeParams struct {
	ProtocolVersion string `json:"protocolVersion"`
}

// handleInitialize handles the initialize method
func (f *MCPServerFilter) handleInitialize(ctx *MCPContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	if f.governanceEnabled && ctx.SessionID() != "" {
		return f.sendBadRequest(ctx, "initialize must not include Mcp-Session-Id")
	}

	// Parse client's protocol version from request params
	var initParams initializeParams

	if req.Params != nil {
		if paramsBytes, err := json.Marshal(req.Params); err == nil {
			json.Unmarshal(paramsBytes, &initParams)
		}
	}

	clientVersion := initParams.ProtocolVersion
	if clientVersion == "" {
		clientVersion = ctx.ProtocolVersion()
	}

	logger.Infof("[dubbo-go-pixiu] mcp server initialize: client_protocol=%s server_protocol=%s",
		clientVersion, constant.MCPProtocolVersion20250618)

	capabilities := mcp.ServerCapabilities{
		Tools: &struct {
			ListChanged bool `json:"listChanged,omitempty"`
		}{
			ListChanged: true,
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

	serverInfo := mcp.Implementation{
		Name:    f.cfg.ServerInfo.Name,
		Version: f.cfg.ServerInfo.Version,
	}

	instructions := f.cfg.ServerInfo.Instructions
	if instructions == "" {
		instructions = "This MCP server provides API access through tools, documentation through resources, and AI assistance through prompts."
	}
	result := mcp.NewInitializeResult(constant.MCPProtocolVersion20250618, capabilities, serverInfo, instructions)

	response := f.responseBuilder.Success(req.ID, result)

	// Per MCP spec: assign a session ID at initialization time for Streamable HTTP transport.
	// Governance mode rejects client-supplied IDs to avoid session fixation; legacy
	// non-governance mode preserves the existing behavior of reusing a valid ID.
	var (
		session *transport.MCPSession
		err     error
	)
	if f.governanceEnabled {
		session, err = f.sessionManager.CreateSession()
	} else {
		session, _, err = f.sessionManager.EnsureSession(ctx.SessionID())
	}
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to create session: %v", err)
		return f.errorHandler.SendInternalError(ctx, req.ID, "failed to create session")
	}
	if f.governanceEnabled && f.plans != nil {
		if err := f.plans.ActivateSession(session.ID, session.Generation); err != nil {
			f.sessionManager.RemoveSession(session.ID)
			logger.Errorf("[dubbo-go-pixiu] mcp server failed to activate router session plan owner: %v", err)
			return f.errorHandler.SendInternalError(ctx, req.ID, "failed to create session")
		}
	}

	// Add Mcp-Session-Id header to the response
	ctx.AddHeader(constant.HeaderKeyMCPSessionId, session.ID)

	logger.Infof("[dubbo-go-pixiu] mcp server created session for client")

	return f.sendJSONResponse(ctx, response)
}

// handleToolsList handles the tools/list method using mcp-go APIs
func (f *MCPServerFilter) handleToolsList(ctx *MCPContext, req mcp.JSONRPCRequest, responseFormat transport.ResponseFormat) filter.FilterStatus {
	response, err := f.buildToolsListResponseObject(ctx, req)
	if err != nil {
		logger.Warnf("[dubbo-go-pixiu] mcp tool router failed to build tools/list response: %v", err)
		return f.errorHandler.SendInternalError(ctx, req.ID, "failed to persist session plan")
	}
	return f.sendResponseWithFormat(ctx, response, responseFormat)
}

// buildToolsListResponseObject builds the tools/list response object (for SSE)
func (f *MCPServerFilter) buildToolsListResponseObject(ctx *MCPContext, req mcp.JSONRPCRequest) (mcp.JSONRPCResponse, error) {
	// Read immutable catalog snapshot to reflect dynamic updates.
	snapshot := f.registry.toolCatalogSnapshotUnsafe()
	toolCfgs := snapshot.orderedToolsUnsafe()

	// Router hookpoint: trim the candidate set to a session-scoped plan.
	// On internal error or invalid session, fail closed so governance failures
	// never expose the raw candidate catalog.
	if f.governanceEnabled {
		if !f.routerSessionValid(ctx) {
			logger.Warnf("[dubbo-go-pixiu] mcp tool router rejected tools/list for invalid session")
			return mcp.JSONRPCResponse{}, fmt.Errorf("mcp session invalidated during request")
		}
		sc := f.buildSelectionContextWithCatalog(ctx, string(mcp.MethodToolsList), "", snapshot.Version)
		plan, err := f.selector.Select(ctx.Ctx, sc, toolCfgs)
		if err != nil {
			logger.Warnf("[dubbo-go-pixiu] mcp tool router Select failed: %v", err)
			toolCfgs, err = f.handleSelectionFailure(ctx, sc, toolCfgs, err)
			if err != nil {
				return mcp.JSONRPCResponse{}, err
			}
		} else {
			toolCfgs = filterByPlan(toolCfgs, plan)
		}
		if !f.routerSessionValid(ctx) {
			return mcp.JSONRPCResponse{}, fmt.Errorf("mcp session invalidated during request")
		}
	}

	tools := f.buildMCPTools(toolCfgs)

	// Build standard MCP tools list response using mcp-go structures
	result := mcp.NewListToolsResult(tools, "")

	return f.responseBuilder.Success(req.ID, result), nil
}

func (f *MCPServerFilter) buildMCPTools(toolCfgs []model.ToolConfig) []mcp.Tool {
	tools := make([]mcp.Tool, 0, len(toolCfgs))
	for _, toolCfg := range toolCfgs {
		tools = append(tools, f.buildMCPTool(toolCfg))
	}
	return tools
}

func (f *MCPServerFilter) buildMCPTool(toolCfg model.ToolConfig) mcp.Tool {
	toolOptions := []mcp.ToolOption{
		mcp.WithDescription(toolCfg.Description),
	}

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

	return mcp.NewTool(toolCfg.Name, toolOptions...)
}

func visibleToolsFingerprint(tools []model.ToolConfig) string {
	mcpTools := (&MCPServerFilter{}).buildMCPTools(discoverableToolConfigs(tools))
	if len(mcpTools) == 0 {
		return router.VisibleToolNamesFingerprint(nil)
	}
	items := make([]string, 0, len(mcpTools))
	for _, tool := range mcpTools {
		data, err := json.Marshal(tool)
		if err != nil {
			data = []byte(tool.Name)
		}
		items = append(items, string(data))
	}
	sort.Strings(items)
	hash := sha256.New()
	for _, item := range items {
		_, _ = hash.Write([]byte(item))
		_, _ = hash.Write([]byte{0})
	}
	return hex.EncodeToString(hash.Sum(nil))
}

func discoverableToolConfigs(tools []model.ToolConfig) []model.ToolConfig {
	if len(tools) == 0 {
		return nil
	}
	out := make([]model.ToolConfig, 0, len(tools))
	for _, tool := range tools {
		if tool.Meta != nil && tool.Meta.DiscoveryVisibility != nil && !*tool.Meta.DiscoveryVisibility {
			continue
		}
		out = append(out, tool)
	}
	return out
}

func (f *MCPServerFilter) routerSessionValid(ctx *MCPContext) bool {
	if !f.governanceEnabled {
		return true
	}
	session := ctx.ValidatedSession()
	if session == nil || ctx.SessionID() == "" || ctx.SessionGeneration() == 0 {
		return false
	}
	return session.ID == ctx.SessionID() && session.Generation == ctx.SessionGeneration() && !session.IsClosed()
}

// buildToolParameterOptions builds the mcp.PropertyOption slice for a given tool argument
func (f *MCPServerFilter) buildToolParameterOptions(arg *model.ArgConfig) []mcp.PropertyOption {
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
func (f *MCPServerFilter) handleResourcesList(ctx *MCPContext, req mcp.JSONRPCRequest, responseFormat transport.ResponseFormat) filter.FilterStatus {
	mcpResources, err := f.registry.ToMCPResources()
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to get MCP resources: %v", err)
		return f.errorHandler.SendInternalError(ctx, req.ID, "failed to get resources")
	}

	response := f.buildResourcesListResponseObject(req, mcpResources)
	return f.sendResponseWithFormat(ctx, response, responseFormat)
}

// buildResourcesListResponseObject builds the resources/list response object (for SSE)
func (f *MCPServerFilter) buildResourcesListResponseObject(req mcp.JSONRPCRequest, mcpResources []mcp.Resource) mcp.JSONRPCResponse {
	// Build resources list response using mcp-go structures
	result := mcp.NewListResourcesResult(mcpResources, "")
	return f.responseBuilder.Success(req.ID, result)
}

// handlePing handles the ping method
func (f *MCPServerFilter) handlePing(ctx *MCPContext, req mcp.JSONRPCRequest, responseFormat transport.ResponseFormat) filter.FilterStatus {
	response := f.buildPingResponseObject(req)
	return f.sendResponseWithFormat(ctx, response, responseFormat)
}

// buildPingResponseObject builds the ping response object (for SSE)
func (f *MCPServerFilter) buildPingResponseObject(req mcp.JSONRPCRequest) mcp.JSONRPCResponse {
	logger.Debugf("[dubbo-go-pixiu] mcp server handling ping request")
	return f.responseBuilder.Success(req.ID, map[string]any{})
}

// handleNotificationsInitialized handles notifications/initialized notification
func (f *MCPServerFilter) handleNotificationsInitialized(ctx *MCPContext, _ mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Debugf("[dubbo-go-pixiu] mcp server received initialized notification from client")

	// Per MCP spec, notifications MUST return 202 Accepted with no body
	ctx.SendLocalReply(http.StatusAccepted, nil)
	return filter.Stop
}

// handleResourceRead handles the resources/read method
func (f *MCPServerFilter) handleResourceRead(ctx *MCPContext, req mcp.JSONRPCRequest, responseFormat transport.ResponseFormat) filter.FilterStatus {
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
	return f.sendResponseWithFormat(ctx, response, responseFormat)
}

// handleResourceTemplatesList handles the resources/templates/list method
func (f *MCPServerFilter) handleResourceTemplatesList(ctx *MCPContext, req mcp.JSONRPCRequest, responseFormat transport.ResponseFormat) filter.FilterStatus {
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
	return f.sendResponseWithFormat(ctx, response, responseFormat)
}

// handlePromptsList handles the prompts/list method
func (f *MCPServerFilter) handlePromptsList(ctx *MCPContext, req mcp.JSONRPCRequest, responseFormat transport.ResponseFormat) filter.FilterStatus {
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
	return f.sendResponseWithFormat(ctx, response, responseFormat)
}

// handlePromptsGet handles the prompts/get method
func (f *MCPServerFilter) handlePromptsGet(ctx *MCPContext, req mcp.JSONRPCRequest, responseFormat transport.ResponseFormat) filter.FilterStatus {
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
	return f.sendResponseWithFormat(ctx, response, responseFormat)
}

// buildPromptMessages builds prompt messages with parameter replacement support
func (f *MCPServerFilter) buildPromptMessages(promptConfig model.PromptConfig, arguments map[string]any) ([]map[string]any, error) {
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

	// Read a single live tool snapshot and use it for both lookup and router
	// authorization so tools/call cannot authorize against stale metadata.
	snapshot := f.registry.toolCatalogSnapshotUnsafe()
	toolCfgs := snapshot.orderedToolsUnsafe()

	// Router hookpoint: enforce that the tool is authorized for this session.
	// This implements discovery/execution separation: even a tool name learned
	// out-of-band cannot be invoked unless it is part of the session plan.
	if f.governanceEnabled {
		if !f.routerSessionValid(ctx) {
			logger.Warnf("[dubbo-go-pixiu] mcp tool router denied tool call for invalid session")
			return f.errorHandler.SendToolCallError(ctx, req.ID, "tool not authorized for this session")
		}
		sc := f.buildSelectionContextWithCatalog(ctx, string(mcp.MethodToolsCall), params.Name, snapshot.Version)
		receipt, err := f.selector.AuthorizeCall(ctx.Ctx, sc, toolCfgs)
		if err != nil {
			logger.Warnf("[dubbo-go-pixiu] mcp tool router denied tool call: %v", err)
			// The client-facing message is intentionally generic and decoupled from
			// the internal error: it does not reveal whether the tool exists, only
			// that it is not callable in this session.
			return f.errorHandler.SendToolCallError(ctx, req.ID, "tool not authorized for this session")
		}
		ctx.SetAuthorizationReceipt(receipt)
	}
	receiptForwarded := false
	defer func() {
		if !receiptForwarded {
			f.finalizeAuthorizationReceipt(ctx, router.ReceiptAborted)
		}
	}()

	toolConfig, exists := snapshot.lookup(params.Name)
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
	receiptForwarded = true

	// Continue to next filter for backend forwarding
	return filter.Continue
}

// buildBackendRequest builds the complete backend request including path, body, and headers
func (f *MCPServerFilter) buildBackendRequest(ctx *MCPContext, toolConfig model.ToolConfig, arguments map[string]any) error {
	// Set HTTP method
	ctx.Request.Method = toolConfig.Request.Method

	// Build request path and body based on argument locations
	path := toolConfig.Request.Path
	bodyParams := make(map[string]any)
	queryParams := make(map[string]string)

	// Process arguments based on their location (path, query, body)
	for argName, argValue := range arguments {
		// Find argument configuration
		var argConfig *model.ArgConfig
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

		logger.Debugf("[dubbo-go-pixiu] mcp server built backend request body")
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
		f.finalizeAuthorizationReceipt(ctx, router.ReceiptAborted)
		return filter.Continue
	}

	// Extract backend response
	responseBody, statusCode, err := f.extractBackendResponse(ctx)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to extract backend response: %v", err)
		f.finalizeAuthorizationReceipt(ctx, router.ReceiptAborted)
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
	if f.governanceEnabled && !f.routerSessionValid(ctx) {
		logger.Warnf("[dubbo-go-pixiu] mcp tool router rejected tool response for invalid session")
		f.finalizeAuthorizationReceipt(ctx, router.ReceiptAborted)
		return f.errorHandler.SendToolCallError(ctx, requestID, "MCP session invalidated during request")
	}
	// Check for backend errors
	if statusCode >= 400 {
		logger.Errorf("[dubbo-go-pixiu] mcp server backend returned error status: %d", statusCode)
		f.finalizeAuthorizationReceipt(ctx, router.ReceiptAborted)
		return f.errorHandler.SendToolCallError(ctx, requestID, fmt.Sprintf("backend error: %d", statusCode))
	}

	// Build successful response using ToolCallSuccess method
	content := strings.TrimSpace(string(responseBody))
	mcpResponse := f.responseBuilder.ToolCallSuccess(requestID, content)
	transitioned := f.finalizeAuthorizationReceipt(ctx, router.ReceiptSucceeded)
	status := f.sendMCPResponse(ctx, mcpResponse)
	if transitioned {
		f.notifyToolsListChanged(ctx.SessionID())
	}
	return status
}

func (f *MCPServerFilter) finalizeAuthorizationReceipt(ctx *MCPContext, outcome router.ReceiptOutcome) bool {
	if !f.governanceEnabled {
		return false
	}
	receipt := ctx.AuthorizationReceipt()
	if receipt == nil {
		return false
	}
	finalizer, ok := f.selector.(router.ReceiptFinalizer)
	if !ok {
		logger.Errorf("[dubbo-go-pixiu] mcp tool router cannot finalize authorization receipt: selector does not implement ReceiptFinalizer")
		ctx.SetAuthorizationReceipt(nil)
		return false
	}
	result, err := finalizer.FinalizeReceipt(ctx.Ctx, *receipt, outcome)
	ctx.SetAuthorizationReceipt(nil)
	if err != nil {
		logger.Warnf("[dubbo-go-pixiu] mcp tool router failed to finalize authorization receipt: %v", err)
		return false
	}
	if !result.Transitioned {
		return false
	}
	snapshot := f.registry.toolCatalogSnapshotUnsafe()
	sc := f.buildSelectionContextWithCatalog(ctx, string(mcp.MethodToolsCall), receipt.ToolName, snapshot.Version)
	if _, err := f.selector.Select(ctx.Ctx, sc, snapshot.orderedToolsUnsafe()); err != nil {
		logger.Warnf("[dubbo-go-pixiu] mcp tool router failed to refresh expanded plan after transition: %v", err)
		return false
	}
	return true
}

// sendMCPResponse sends an MCP response and updates the target response
func (f *MCPServerFilter) sendMCPResponse(ctx *MCPContext, response mcp.JSONRPCResponse) filter.FilterStatus {
	// Check if we should send via SSE (when session exists with active SSE stream)
	sessionID := ctx.SessionID()
	if sessionID != "" {
		session, exists := f.sessionManager.Session(sessionID)
		if exists && session.HasPipeWriter() {
			// Send via SSE stream
			if sseErr := f.sseHandler.SendSSEMessage(session, response); sseErr == nil {
				logger.Debugf("[dubbo-go-pixiu] mcp server sent tool call response via SSE")
				// Return 202 Accepted without body per MCP spec
				ctx.SendLocalReply(http.StatusAccepted, nil)
				return filter.Stop
			} else {
				// If SSE send failed, fall through to JSON response
				logger.Warnf("[dubbo-go-pixiu] mcp server SSE send failed, falling back to JSON: %v", sseErr)
			}
		}
	}

	// Default: send as JSON response (no SSE stream available)
	mcpResponseBody, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to marshal MCP response: %v", err)
		return filter.Continue
	}

	// Override TargetResp to ensure MCP format response is sent
	ctx.TargetResp = &client.UnaryResponse{Data: mcpResponseBody}
	ctx.StatusCode(http.StatusOK)
	ctx.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueApplicationJson)
	ctx.ClearContentLengthHeader()

	logger.Debugf("[dubbo-go-pixiu] mcp server successfully wrapped backend response in MCP format")
	return filter.Continue
}

func (f *MCPServerFilter) notifyToolsListChanged(sessionID string) {
	if sessionID == "" {
		return
	}
	session, exists := f.sessionManager.GetSession(sessionID)
	if !exists {
		return
	}
	if _, _, err := markToolsListChangedAndFlush(f.sseHandler, session); err != nil {
		logger.Warnf("[dubbo-go-pixiu] mcp server failed to flush pending tools/list_changed: %v", err)
	}
}

func (f *MCPServerFilter) flushPendingToolsListChanged(session *transport.MCPSession) {
	if err := flushToolsListChangedNotification(f.sseHandler, session); err != nil {
		logger.Warnf("[dubbo-go-pixiu] mcp server failed to flush pending tools/list_changed: %v", err)
	}
}

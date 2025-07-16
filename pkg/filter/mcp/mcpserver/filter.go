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
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	h "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/creasty/defaults"
	"github.com/mark3labs/mcp-go/mcp"
	"io"
	"net/http"
	"strings"
)

// MCP 方法常量 - 使用 mcp-go 库的常量
const (
	methodInitialize        = string(mcp.MethodInitialize)    // "initialize"
	methodToolsList         = string(mcp.MethodToolsList)     // "tools/list"
	methodToolsCall         = string(mcp.MethodToolsCall)     // "tools/call"
	methodResourcesList     = string(mcp.MethodResourcesList) // "resources/list"
	methodResourcesRead     = string(mcp.MethodResourcesRead) // "resources/read"
	methodPromptsList       = string(mcp.MethodPromptsList)   // "prompts/list"
	methodPromptsGet        = string(mcp.MethodPromptsGet)    // "prompts/get"
	methodNotificationsInit = "notifications/initialized"     // 初始化完成通知
	methodPing              = string(mcp.MethodPing)          // 连接保活
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
	// 设置配置默认值
	if err := defaults.Set(f.cfg); err != nil {
		return fmt.Errorf("failed to set config defaults: %v", err)
	}

	f.registry = NewToolRegistry()

	// 注意：方法处理现在通过 handleTerminalMethod 统一处理

	// 验证并注册静态配置的工具
	for i, tool := range f.cfg.Tools {
		// 为每个工具设置默认值
		if err := defaults.Set(&tool); err != nil {
			return fmt.Errorf("failed to set defaults for tool %s: %v", tool.Name, err)
		}

		// 为工具的参数设置默认值
		for j := range tool.Args {
			if err := defaults.Set(&tool.Args[j]); err != nil {
				return fmt.Errorf("failed to set defaults for arg %s in tool %s: %v",
					tool.Args[j].Name, tool.Name, err)
			}
		}

		// 为响应配置设置默认值
		if tool.Response != nil {
			if err := defaults.Set(tool.Response); err != nil {
				return fmt.Errorf("failed to set defaults for response in tool %s: %v", tool.Name, err)
			}
		}

		// 更新配置中的工具（因为我们修改了副本）
		f.cfg.Tools[i] = tool

		// 验证工具配置
		if err := tool.Validate(); err != nil {
			return fmt.Errorf("invalid tool configuration for %s: %v", tool.Name, err)
		}

		if err := f.registry.RegisterTool(tool); err != nil {
			return fmt.Errorf("failed to register tool %s: %v", tool.Name, err)
		}
		logger.Infof("mcp_server: registered tool %s", tool.Name)
	}

	// 注册静态配置的资源
	for _, resource := range f.cfg.Resources {
		if err := f.registry.RegisterResource(resource); err != nil {
			return fmt.Errorf("failed to register resource %s: %v", resource.Name, err)
		}
		logger.Infof("mcp_server: registered resource %s", resource.Name)
	}

	// 注册静态配置的提示词
	for _, prompt := range f.cfg.Prompts {
		if err := f.registry.RegisterPrompt(prompt); err != nil {
			return fmt.Errorf("failed to register prompt %s: %v", prompt.Name, err)
		}
		logger.Infof("mcp_server: registered prompt %s", prompt.Name)
	}

	toolCount, resourceCount, promptCount := f.registry.Count()
	logger.Infof("mcp_server: initialized with %d tools, %d resources, and %d prompts",
		toolCount, resourceCount, promptCount)

	return nil
}

// Config 返回配置结构体
func (f *FilterFactory) Config() any {
	return f.cfg
}

// PrepareFilterChain 准备过滤器链
func (f *FilterFactory) PrepareFilterChain(ctx *h.HttpContext, chain filter.FilterChain) error {
	mcpFilter := &MCPServerFilter{
		cfg:      f.cfg,
		registry: f.registry,
	}
	chain.AppendDecodeFilters(mcpFilter)
	chain.AppendEncodeFilters(mcpFilter) // 添加到 Encode 链
	return nil
}

// Decode processes incoming HTTP requests for MCP protocol.
func (f *MCPServerFilter) Decode(ctx *h.HttpContext) filter.FilterStatus {
	// 检查是否是 MCP 请求
	if !f.isMCPRequest(ctx) {
		return filter.Continue
	}

	logger.Infof("mcp_server: [DECODE] processing MCP request: %s %s", ctx.Request.Method, ctx.Request.URL.Path)

	// 标记为 MCP 请求
	ctx.SetMCPRequest(true)

	// 读取请求体
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Errorf("mcp_server: failed to read request body: %v", err)
		return f.sendError(ctx, nil, "failed to read request body")
	}

	// 解析 JSON-RPC 请求
	var jsonrpcReq mcp.JSONRPCRequest
	if err := json.Unmarshal(body, &jsonrpcReq); err != nil {
		logger.Errorf("mcp_server: failed to parse JSON-RPC request: %v", err)
		return f.sendError(ctx, nil, "invalid JSON-RPC request")
	}

	logger.Debugf("mcp_server: received method: %s, id: %v", jsonrpcReq.Method, jsonrpcReq.ID)

	// 在 HttpContext 中存储 MCP 信息
	ctx.SetMCPMethod(jsonrpcReq.Method)
	ctx.SetMCPRequestID(jsonrpcReq.ID)

	// 区分终端处理和转发处理的 MCP 方法
	if f.isTerminalMethod(jsonrpcReq.Method) {
		// 终端方法：直接在 Decode 阶段处理并返回响应
		return f.handleTerminalMethod(ctx, jsonrpcReq)
	} else if jsonrpcReq.Method == methodToolsCall {
		// 工具调用：需要转发到后端服务，在 Encode 阶段包装响应
		return f.handleToolCall(ctx, jsonrpcReq)
	} else {
		// 未知方法
		logger.Warnf("mcp_server: unsupported method: %s", jsonrpcReq.Method)
		return f.sendMethodNotFound(ctx, jsonrpcReq.ID)
	}
}

// isTerminalMethod 检查是否是终端方法（不需要转发到后端的方法）
func (f *MCPServerFilter) isTerminalMethod(method string) bool {
	switch method {
	case methodInitialize, methodToolsList, methodResourcesList, methodResourcesRead,
		methodPromptsList, methodPromptsGet, methodNotificationsInit, methodPing:
		return true
	default:
		return false
	}
}

// handleTerminalMethod 处理终端方法
func (f *MCPServerFilter) handleTerminalMethod(ctx *h.HttpContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	switch req.Method {
	case methodInitialize:
		return f.handleInitialize(ctx, req)
	case methodToolsList:
		return f.handleToolsList(ctx, req)
	case methodResourcesList:
		return f.handleResourcesList(ctx, req)
	case methodResourcesRead:
		return f.handleResourceRead(ctx, req)
	case methodPromptsList:
		return f.handlePromptsList(ctx, req)
	case methodPromptsGet:
		return f.handlePromptsGet(ctx, req)
	case methodNotificationsInit:
		return f.handleNotificationsInitialized(ctx, req)
	case methodPing:
		return f.handlePing(ctx, req)
	default:
		logger.Warnf("mcp_server: unsupported terminal method: %s", req.Method)
		return f.sendMethodNotFound(ctx, req.ID)
	}
}

// Encode processes outgoing HTTP responses.
func (f *MCPServerFilter) Encode(ctx *h.HttpContext) filter.FilterStatus {
	// 检查是否是工具调用响应（通过 HttpContext.Params 检查）
	if ctx.IsMCPToolCall() {
		// 检查是否已经处理过，防止重复处理
		if ctx.IsMCPProcessed() {
			logger.Debugf("mcp_server: [ENCODE] Tool call already processed, skipping")
			return filter.Continue
		}

		logger.Infof("mcp_server: [ENCODE] Processing tool call response for URL: %s", ctx.Request.URL.Path)

		// 设置已处理标记
		ctx.SetMCPProcessed(true)

		return f.handleToolCallResponse(ctx)
	}

	// 对于普通的 MCP 请求，不需要特殊处理
	if ctx.IsMCPRequest() {
		logger.Debugf("mcp_server: [ENCODE] Regular MCP request, no special processing needed")
	}

	return filter.Continue
}

// isMCPRequest 检查是否是 MCP 请求
func (f *MCPServerFilter) isMCPRequest(ctx *h.HttpContext) bool {
	return ctx.Request.URL.Path == f.cfg.Endpoint
}

// JSONRPCError JSON-RPC 错误结构
type JSONRPCError struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
	Data    any    `json:"data,omitempty"`
}

// JSONRPCResponse 通用 JSON-RPC 响应结构
type JSONRPCResponse struct {
	JSONRPC string        `json:"jsonrpc"`
	ID      any           `json:"id"`
	Result  any           `json:"result,omitempty"`
	Error   *JSONRPCError `json:"error,omitempty"`
}

// NewJSONRPCResponse 创建成功响应
func NewJSONRPCResponse(id any, result any) *JSONRPCResponse {
	return &JSONRPCResponse{
		JSONRPC: "2.0",
		ID:      id,
		Result:  result,
	}
}

// NewJSONRPCError 创建错误响应
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

// sendError 发送错误响应
func (f *MCPServerFilter) sendError(ctx *h.HttpContext, id any, message string) filter.FilterStatus {
	response := NewJSONRPCError(id, -32603, message)
	return f.sendJSONResponse(ctx, response)
}

// sendMethodNotFound 发送方法未找到错误
func (f *MCPServerFilter) sendMethodNotFound(ctx *h.HttpContext, id any) filter.FilterStatus {
	response := NewJSONRPCError(id, -32601, "Method not found")
	return f.sendJSONResponse(ctx, response)
}

// sendJSONResponse 发送 JSON 响应
func (f *MCPServerFilter) sendJSONResponse(ctx *h.HttpContext, response any) filter.FilterStatus {
	responseBody, err := json.Marshal(response)
	if err != nil {
		logger.Errorf("mcp_server: failed to marshal response: %v", err)
		ctx.SendLocalReply(http.StatusInternalServerError, []byte("internal server error"))
		return filter.Stop
	}

	logger.Infof("mcp_server: marshaled response body: %s", string(responseBody))
	ctx.SendLocalReply(http.StatusOK, responseBody)
	return filter.Stop
}

// buildRequestPath 构建请求路径，替换路径参数
func (f *MCPServerFilter) buildRequestPath(pathTemplate string, arguments map[string]any) (string, error) {
	path := pathTemplate

	// 提取所有路径参数
	pathParams := GetPathParameterNames(pathTemplate)

	// 替换每个占位符
	for _, paramName := range pathParams {
		placeholder := "{" + paramName + "}"

		// 获取参数值
		value, exists := arguments[paramName]
		if !exists {
			return "", fmt.Errorf("required path parameter '%s' is missing", paramName)
		}

		// 替换占位符
		var strValue string
		if str, ok := value.(string); ok {
			strValue = str
		} else {
			strValue = fmt.Sprintf("%v", value)
		}

		path = strings.Replace(path, placeholder, strValue, -1)
	}

	// 检查是否还有未替换的占位符
	if strings.Contains(path, "{") && strings.Contains(path, "}") {
		return "", fmt.Errorf("unresolved placeholders in path: %s", path)
	}

	return path, nil
}

// InitializeResponse 初始化响应结构
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

// handleInitialize 处理 initialize 方法
func (f *MCPServerFilter) handleInitialize(ctx *h.HttpContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Infof("mcp_server: handling initialize request - USING STRUCT VERSION")

	// 创建初始化响应，必须包含 capabilities 字段
	response := InitializeResponse{
		JSONRPC: "2.0",
		ID:      req.ID,
	}

	response.Result.ProtocolVersion = "2024-11-05"
	response.Result.Capabilities.Tools.ListChanged = true
	response.Result.Capabilities.Resources.ListChanged = true
	response.Result.Capabilities.Prompts.ListChanged = true
	response.Result.ServerInfo.Name = f.cfg.ServerInfo.Name
	response.Result.ServerInfo.Version = f.cfg.ServerInfo.Version

	// 如果配置中有说明，添加到响应中
	if f.cfg.ServerInfo.Instructions != "" {
		response.Result.Instructions = f.cfg.ServerInfo.Instructions
	}

	logger.Infof("mcp_server: sending initialize response: %+v", response)
	return f.sendJSONResponse(ctx, response)
}

// handleToolsList 处理 tools/list 方法
func (f *MCPServerFilter) handleToolsList(ctx *h.HttpContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Debugf("mcp_server: handling tools/list request")

	// 获取所有工具
	mcpTools, err := f.registry.ToMCPTools()
	if err != nil {
		logger.Errorf("mcp_server: failed to get MCP tools: %v", err)
		return f.sendError(ctx, req.ID, "failed to get tools")
	}

	// 构建工具列表响应
	result := map[string]any{
		"tools": mcpTools,
	}

	response := NewJSONRPCResponse(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// handleResourcesList 处理 resources/list 方法
func (f *MCPServerFilter) handleResourcesList(ctx *h.HttpContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Debugf("mcp_server: handling resources/list request")

	// 获取所有资源
	mcpResources, err := f.registry.ToMCPResources()
	if err != nil {
		logger.Errorf("mcp_server: failed to get MCP resources: %v", err)
		return f.sendError(ctx, req.ID, "failed to get resources")
	}

	// 构建资源列表响应
	result := map[string]any{
		"resources": mcpResources,
	}

	response := NewJSONRPCResponse(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// handleResourceRead 处理 resources/read 方法
func (f *MCPServerFilter) handleResourceRead(ctx *h.HttpContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Debugf("mcp_server: handling resources/read request")

	// 解析请求参数
	paramsBytes, err := json.Marshal(req.Params)
	if err != nil {
		logger.Errorf("mcp_server: failed to marshal params: %v", err)
		return f.sendError(ctx, req.ID, "invalid parameters")
	}

	var params struct {
		URI string `json:"uri"`
	}
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		logger.Errorf("mcp_server: failed to parse resource read params: %v", err)
		return f.sendError(ctx, req.ID, "invalid parameters")
	}

	// 查找资源
	resource, exists := f.registry.GetResource(params.URI)
	if !exists {
		logger.Warnf("mcp_server: resource not found: %s", params.URI)
		return f.sendError(ctx, req.ID, fmt.Sprintf("resource not found: %s", params.URI))
	}

	// 读取资源内容（这里简化处理，实际应该根据 source 配置读取）
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

// handleToolCall 处理 tools/call 方法
func (f *MCPServerFilter) handleToolCall(ctx *h.HttpContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Infof("mcp_server: [TOOL_CALL] handling tools/call request with ID: %v", req.ID)

	// 解析请求参数
	paramsBytes, err := json.Marshal(req.Params)
	if err != nil {
		logger.Errorf("mcp_server: failed to marshal params: %v", err)
		return f.sendError(ctx, req.ID, "invalid parameters")
	}

	var params struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments"`
	}
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		logger.Errorf("mcp_server: failed to parse tool call params: %v", err)
		return f.sendError(ctx, req.ID, "invalid parameters")
	}

	// 查找工具配置
	toolConfig, exists := f.registry.GetTool(params.Name)
	if !exists {
		logger.Warnf("mcp_server: tool not found: %s", params.Name)
		return f.sendError(ctx, req.ID, fmt.Sprintf("tool not found: %s", params.Name))
	}

	// 处理参数并构建请求
	err = f.processParameters(ctx, toolConfig, params.Arguments)
	if err != nil {
		logger.Errorf("mcp_server: failed to process parameters: %v", err)
		return f.sendError(ctx, req.ID, "failed to process parameters")
	}

	// 设置 HttpContext 路由信息（关键！）
	ctx.Request.Method = toolConfig.Request.Method

	// 设置请求头
	for key, value := range toolConfig.Request.Headers {
		ctx.Request.Header.Set(key, value)
	}

	// 设置 MCP 相关信息到 HttpContext
	ctx.SetMCPCluster(toolConfig.Cluster)
	ctx.SetMCPToolCall(true)
	ctx.SetMCPToolName(params.Name)

	logger.Infof("mcp_server: [TOOL_CALL] set MCP context - tool=%s, cluster=%s, requestID=%v",
		params.Name, toolConfig.Cluster, req.ID)
	logger.Infof("mcp_server: [TOOL_CALL] forwarding to cluster %s, method %s, path %s",
		toolConfig.Cluster, ctx.Request.Method, ctx.Request.URL.Path)

	// 继续到下一个 filter（httpproxy）
	logger.Infof("mcp_server: [TOOL_CALL] returning filter.Continue to proceed to httpproxy")
	return filter.Continue
}

// processParameters 根据参数配置处理请求参数
func (f *MCPServerFilter) processParameters(ctx *h.HttpContext, toolConfig ToolConfig, arguments map[string]any) error {
	// 获取所有参数（包括自动推断的路径参数）
	allParams, err := toolConfig.GetAllParameters()
	if err != nil {
		return fmt.Errorf("failed to get parameters: %v", err)
	}

	// 验证必需参数并应用默认值
	processedArgs := make(map[string]any)
	for key, value := range arguments {
		processedArgs[key] = value
	}

	// 收集请求体参数
	bodyParams := make(map[string]any)

	// 处理每个参数
	for _, param := range allParams {
		value := processedArgs[param.Name]

		// 检查必需参数
		if value == nil {
			if param.Default != nil {
				processedArgs[param.Name] = param.Default
				value = param.Default
			} else if param.Required {
				return fmt.Errorf("required parameter %s is missing", param.Name)
			} else {
				continue // 跳过可选且无默认值的参数
			}
		}

		// 根据参数位置处理
		switch param.In {
		case "path":
			// 路径参数通过路径构建处理
			continue
		case "query":
			// 添加到查询参数
			query := ctx.Request.URL.Query()
			query.Set(param.Name, fmt.Sprintf("%v", value))
			ctx.Request.URL.RawQuery = query.Encode()
		case "header":
			// 设置请求头
			ctx.Request.Header.Set(param.Name, fmt.Sprintf("%v", value))
		case "body":
			// 收集请求体参数
			bodyParams[param.Name] = value
		default:
			logger.Warnf("mcp_server: unsupported parameter location: %s for parameter %s", param.In, param.Name)
		}
	}

	// 处理请求体参数
	if len(bodyParams) > 0 {
		if err := f.setRequestBody(ctx, bodyParams); err != nil {
			return fmt.Errorf("failed to set request body: %v", err)
		}
	}

	// 构建路径（处理路径参数）
	path, err := f.buildRequestPath(toolConfig.Request.Path, processedArgs)
	if err != nil {
		return fmt.Errorf("failed to build path: %v", err)
	}
	ctx.Request.URL.Path = path

	return nil
}

// setRequestBody 设置请求体参数
func (f *MCPServerFilter) setRequestBody(ctx *h.HttpContext, bodyParams map[string]any) error {
	// 将请求体参数序列化为 JSON
	bodyData, err := json.Marshal(bodyParams)
	if err != nil {
		return fmt.Errorf("failed to marshal body parameters: %v", err)
	}

	// 设置请求体
	ctx.Request.Body = io.NopCloser(bytes.NewReader(bodyData))
	ctx.Request.ContentLength = int64(len(bodyData))

	// 设置 Content-Type 头部
	if ctx.Request.Header.Get("Content-Type") == "" {
		ctx.Request.Header.Set("Content-Type", "application/json")
	}

	return nil
}

// handlePromptsList 处理 prompts/list 方法
func (f *MCPServerFilter) handlePromptsList(ctx *h.HttpContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Debugf("mcp_server: handling prompts/list request")

	// 获取所有提示词
	mcpPrompts, err := f.registry.ToMCPPrompts()
	if err != nil {
		logger.Errorf("mcp_server: failed to get prompts: %v", err)
		return f.sendError(ctx, req.ID, "failed to get prompts")
	}

	result := map[string]any{
		"prompts": mcpPrompts,
	}

	response := NewJSONRPCResponse(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// handlePromptsGet 处理 prompts/get 方法
func (f *MCPServerFilter) handlePromptsGet(ctx *h.HttpContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Debugf("mcp_server: handling prompts/get request")

	// 解析请求参数
	paramsBytes, err := json.Marshal(req.Params)
	if err != nil {
		logger.Errorf("mcp_server: failed to marshal params: %v", err)
		return f.sendError(ctx, req.ID, "invalid parameters")
	}

	var params struct {
		Name      string         `json:"name"`
		Arguments map[string]any `json:"arguments,omitempty"`
	}
	if err := json.Unmarshal(paramsBytes, &params); err != nil {
		logger.Errorf("mcp_server: failed to parse prompts/get params: %v", err)
		return f.sendError(ctx, req.ID, "invalid parameters")
	}

	// 查找提示词配置
	promptConfig, exists := f.registry.GetPrompt(params.Name)
	if !exists {
		logger.Warnf("mcp_server: prompt not found: %s", params.Name)
		return f.sendError(ctx, req.ID, fmt.Sprintf("prompt not found: %s", params.Name))
	}

	// 构建提示词消息
	messages, err := f.buildPromptMessages(promptConfig, params.Arguments)
	if err != nil {
		logger.Errorf("mcp_server: failed to build prompt messages: %v", err)
		return f.sendError(ctx, req.ID, "failed to build prompt messages")
	}

	result := map[string]any{
		"description": promptConfig.Description,
		"messages":    messages,
	}

	response := NewJSONRPCResponse(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// buildPromptMessages 构建提示词消息，支持参数替换
func (f *MCPServerFilter) buildPromptMessages(promptConfig PromptConfig, arguments map[string]any) ([]map[string]any, error) {
	messages := make([]map[string]any, 0, len(promptConfig.Messages))

	for _, msgConfig := range promptConfig.Messages {
		// 替换消息内容中的参数
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

// replacePromptArguments 替换提示词内容中的参数占位符
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

// handleNotificationsInitialized 处理 notifications/initialized 通知
func (f *MCPServerFilter) handleNotificationsInitialized(ctx *h.HttpContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Debugf("mcp_server: received initialized notification from client")

	// notifications/initialized 是一个通知，不需要响应
	// 根据 MCP 协议，通知不应该有响应
	// 我们只需要记录这个事件，表示客户端已经准备好进行正常操作

	logger.Infof("mcp_server: client initialization completed, ready for normal operations")

	// 对于通知，我们不发送响应，直接返回 Stop
	return filter.Stop
}

// handlePing 处理 ping 方法
func (f *MCPServerFilter) handlePing(ctx *h.HttpContext, req mcp.JSONRPCRequest) filter.FilterStatus {
	logger.Debugf("mcp_server: handling ping request")

	// ping 方法通常用于连接保活，返回一个简单的 pong 响应
	result := map[string]any{} // 空的结果对象

	response := NewJSONRPCResponse(req.ID, result)
	return f.sendJSONResponse(ctx, response)
}

// handleToolCallResponse 处理工具调用的响应，将后端响应包装成 MCP 格式
func (f *MCPServerFilter) handleToolCallResponse(ctx *h.HttpContext) filter.FilterStatus {
	logger.Infof("mcp_server: [RESPONSE] handleToolCallResponse called")

	// 获取原始的 JSON-RPC 请求 ID（从 HttpContext.Params）
	requestID := ctx.GetMCPRequestID()
	if requestID == nil {
		logger.Errorf("mcp_server: [RESPONSE] missing request ID for tool call response")
		return filter.Continue
	}

	logger.Infof("mcp_server: [RESPONSE] processing response for request ID: %v", requestID)

	// 获取后端响应（从 TargetResp 获取已处理的数据）
	var responseBody []byte
	var statusCode int

	// 从 TargetResp 获取 buildTargetResponse 已处理的数据
	if ctx.TargetResp != nil {
		if unaryResp, ok := ctx.TargetResp.(*client.UnaryResponse); ok {
			responseBody = unaryResp.Data
			statusCode = ctx.GetStatusCode()

			// 检查是否已经是 MCP 格式的响应（防止重复包装）
			if bytes.Contains(responseBody, []byte(`"jsonrpc":"2.0"`)) {
				logger.Infof("mcp_server: [RESPONSE] response already in MCP format, skipping processing")
				return filter.Continue
			}

			logger.Infof("mcp_server: [RESPONSE] got processed response: status=%d, body=%s", statusCode, string(responseBody))
		} else {
			logger.Errorf("mcp_server: [RESPONSE] unexpected TargetResp type: %T", ctx.TargetResp)
			return f.sendToolCallError(ctx, requestID, "unexpected response type")
		}
	} else {
		logger.Errorf("mcp_server: [RESPONSE] no TargetResp available")
		return f.sendToolCallError(ctx, requestID, "no response data available")
	}

	// 检查响应是否为空
	if len(responseBody) == 0 {
		logger.Errorf("mcp_server: empty response body from backend")
		return f.sendToolCallError(ctx, requestID, "empty response from backend")
	}

	// 检查 HTTP 状态码
	if statusCode >= 400 {
		logger.Errorf("mcp_server: backend returned error status: %d", statusCode)
		return f.sendToolCallError(ctx, requestID, fmt.Sprintf("backend error: %d", statusCode))
	}

	// 构建 MCP 工具调用响应
	result := f.buildToolCallResult(responseBody, statusCode)

	// 创建 JSON-RPC 响应
	mcpResponse := NewJSONRPCResponse(requestID, result)

	// 将 MCP 响应序列化为 JSON
	mcpResponseBody, err := json.Marshal(mcpResponse)
	if err != nil {
		logger.Errorf("mcp_server: [RESPONSE] failed to marshal MCP response: %v", err)
		return f.sendToolCallError(ctx, requestID, "failed to marshal response")
	}

	// 覆盖 TargetResp 以确保发送 MCP 格式的响应
	ctx.TargetResp = &client.UnaryResponse{Data: mcpResponseBody}
	ctx.StatusCode(http.StatusOK)
	ctx.AddHeader("Content-Type", "application/json")

	// 清除响应头部的 Content-Length，让 HTTP 库自动计算新的长度
	ctx.Writer.Header().Del("Content-Length")

	logger.Infof("mcp_server: [RESPONSE] MCP response set in TargetResp: %s", string(mcpResponseBody))
	return filter.Continue
}

// buildToolCallResult 构建工具调用结果
func (f *MCPServerFilter) buildToolCallResult(responseBody []byte, statusCode int) map[string]any {
	// 尝试解析响应为 JSON
	var jsonData any
	if err := json.Unmarshal(responseBody, &jsonData); err == nil {
		// 如果是有效的 JSON，返回结构化内容
		return map[string]any{
			"content": []map[string]any{
				{
					"type": "text",
					"text": string(responseBody),
				},
			},
			"isError": false,
		}
	} else {
		// 如果不是有效的 JSON，作为纯文本处理
		return map[string]any{
			"content": []map[string]any{
				{
					"type": "text",
					"text": string(responseBody),
				},
			},
			"isError": false,
		}
	}
}

// sendToolCallError 发送工具调用错误响应
func (f *MCPServerFilter) sendToolCallError(ctx *h.HttpContext, requestID any, message string) filter.FilterStatus {
	// 使用 MCP 工具调用结果格式
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

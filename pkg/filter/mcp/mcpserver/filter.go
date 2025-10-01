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

	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/mark3labs/mcp-go/mcp"

	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/transport"
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
		cfg               *Config
		registry          *ToolRegistry
		errorHandler      *ErrorHandler
		responseBuilder   *ResponseBuilder
		sessionManager    *transport.SessionManager
		sseHandler        *transport.SSEHandler
		contentNegotiator *transport.ContentNegotiator
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
func (f *FilterFactory) PrepareFilterChain(_ *contexthttp.HttpContext, chain filter.FilterChain) error {
	// Initialize transport components
	sessionManager := transport.NewSessionManager()
	sseHandler := transport.NewSSEHandler(sessionManager)
	contentNegotiator := transport.NewContentNegotiator()

	mcpFilter := &MCPServerFilter{
		cfg:               f.cfg,
		registry:          f.registry,
		errorHandler:      NewErrorHandler(),
		responseBuilder:   NewResponseBuilder(),
		sessionManager:    sessionManager,
		sseHandler:        sseHandler,
		contentNegotiator: contentNegotiator,
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

	// Parse and validate HTTP headers
	mcpCtx.ParseAndSetProtocolVersionHeader()
	mcpCtx.ParseAndSetSessionHeader()
	mcpCtx.ParseAndSetAcceptHeader()

	// Validate MCP protocol version
	if !f.validateMCPProtocolVersion(mcpCtx) {
		return f.sendBadRequest(mcpCtx, "Invalid or missing MCP-Protocol-Version header")
	}

	// Dispatch based on HTTP method
	switch ctx.Request.Method {
	case constant.Get:
		return f.handleGetRequest(mcpCtx)
	case constant.Post:
		return f.handlePostRequest(mcpCtx)
	default:
		return f.sendMethodNotAllowed(mcpCtx)
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
		ctx.SendLocalReply(http.StatusInternalServerError, []byte("internal server error"))
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

// validateMCPProtocolVersion validates the MCP protocol version header
func (f *MCPServerFilter) validateMCPProtocolVersion(ctx *MCPContext) bool {
	version := ctx.ProtocolVersion()
	// Allow empty version for backward compatibility
	if version == "" {
		return true
	}
	// Currently only support the 2025-06-18 version
	return version == constant.MCPProtocolVersion2025
}

// handleGetRequest handles HTTP GET requests for SSE stream establishment
func (f *MCPServerFilter) handleGetRequest(ctx *MCPContext) filter.FilterStatus {
	logger.Infof("[dubbo-go-pixiu] mcp server handling GET request for SSE stream")

	// Validate Accept header includes text/event-stream
	if !ctx.AcceptSSE() {
		logger.Warnf("[dubbo-go-pixiu] mcp server GET request must accept text/event-stream")
		return f.sendNotAcceptable(ctx, "GET request must accept text/event-stream")
	}

	// Get or create session
	sessionIDHeader := ctx.SessionID()
	session, isNewSession := f.sessionManager.EnsureSession(sessionIDHeader)

	// Store session ID in context
	ctx.SetSessionID(session.ID)

	// Establish SSE connection
	if err := f.sseHandler.EstablishSSEConnection(ctx.Writer, ctx.Request, session); err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to establish SSE connection: %v", err)
		return f.sendInternalError(ctx, "failed to establish SSE connection")
	}

	if isNewSession {
		logger.Infof("[dubbo-go-pixiu] mcp server established new SSE stream for session: %s", session.ID)
	} else {
		logger.Infof("[dubbo-go-pixiu] mcp server resumed SSE stream for existing session: %s", session.ID)
	}

	return filter.Stop
}

// handlePostRequest handles HTTP POST requests with JSON-RPC messages
func (f *MCPServerFilter) handlePostRequest(ctx *MCPContext) filter.FilterStatus {
	logger.Debugf("[dubbo-go-pixiu] mcp server handling POST request")

	// Read request body
	body, err := io.ReadAll(ctx.Request.Body)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to read request body: %v", err)
		return f.errorHandler.SendInternalError(ctx, nil, "failed to read request body")
	}

	// Parse JSON-RPC request
	var jsonrpcReq mcp.JSONRPCRequest
	if err := json.Unmarshal(body, &jsonrpcReq); err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to parse JSON-RPC request: %v", err)
		return f.errorHandler.SendInternalError(ctx, nil, "invalid JSON-RPC request")
	}

	logger.Infof("[dubbo-go-pixiu] mcp server received POST request: %s (id: %v)", jsonrpcReq.Method, jsonrpcReq.ID)

	// Store information in MCP context
	ctx.SetMCPMethod(jsonrpcReq.Method)
	ctx.SetMCPRequestID(jsonrpcReq.ID)

	// Determine response format based on content negotiation
	sessionID := ctx.SessionID()
	hasSession := sessionID != "" && f.sessionExists(sessionID)
	responseFormat := f.contentNegotiator.NegotiateResponse(
		ctx.Request.Header.Get(constant.HeaderKeyAccept), jsonrpcReq, hasSession)

	// Store response format decision in context for later use
	ctx.StoreMCPDataInParams()

	// Handle different request types
	if f.isTerminalMethod(jsonrpcReq.Method) {
		return f.handleTerminalMethodWithNegotiation(ctx, jsonrpcReq, responseFormat)
	} else if jsonrpcReq.Method == string(mcp.MethodToolsCall) {
		return f.handleToolCallWithNegotiation(ctx, jsonrpcReq, responseFormat)
	} else {
		// Unknown method
		logger.Warnf("[dubbo-go-pixiu] mcp server unsupported method: %s", jsonrpcReq.Method)
		return f.errorHandler.SendMethodNotFound(ctx, jsonrpcReq.ID)
	}
}

// handleTerminalMethodWithNegotiation handles terminal methods with response format negotiation
func (f *MCPServerFilter) handleTerminalMethodWithNegotiation(ctx *MCPContext, req mcp.JSONRPCRequest, responseFormat transport.ResponseFormat) filter.FilterStatus {
	// For terminal methods, we handle them directly and apply response format
	response := f.processTerminalMethod(ctx, req)
	if response == nil {
		return filter.Stop // Error already handled
	}

	return f.sendResponseWithFormat(ctx, response, responseFormat)
}

// handleToolCallWithNegotiation handles tool calls with response format negotiation
func (f *MCPServerFilter) handleToolCallWithNegotiation(ctx *MCPContext, req mcp.JSONRPCRequest, responseFormat transport.ResponseFormat) filter.FilterStatus {
	// For notifications (no response needed), send 202 Accepted immediately
	if responseFormat == transport.ResponseFormatAccepted {
		ctx.SendLocalReply(http.StatusAccepted, nil)
		return filter.Stop
	}

	// For tool calls, we need to forward to backend, so continue with existing logic
	// but store the response format for use in Encode stage
	return f.handleToolCall(ctx, req)
}

// processTerminalMethod processes a terminal method and returns the response
func (f *MCPServerFilter) processTerminalMethod(ctx *MCPContext, req mcp.JSONRPCRequest) any {
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
		logger.Warnf("[dubbo-go-pixiu] mcp server unsupported terminal method: %s", req.Method)
		f.errorHandler.SendMethodNotFound(ctx, req.ID)
		return nil
	}
}

// sendResponseWithFormat sends response in the negotiated format
func (f *MCPServerFilter) sendResponseWithFormat(ctx *MCPContext, response any, format transport.ResponseFormat) filter.FilterStatus {
	switch format {
	case transport.ResponseFormatJSON:
		return f.sendJSONResponse(ctx, response)
	case transport.ResponseFormatSSE:
		return f.sendSSEResponse(ctx, response)
	case transport.ResponseFormatAccepted:
		ctx.SendLocalReply(http.StatusAccepted, nil)
		return filter.Stop
	default:
		return f.sendJSONResponse(ctx, response)
	}
}

// sendSSEResponse sends response via SSE stream
func (f *MCPServerFilter) sendSSEResponse(ctx *MCPContext, response any) filter.FilterStatus {
	sessionID := ctx.SessionID()
	if sessionID == "" {
		// No session, fall back to JSON
		return f.sendJSONResponse(ctx, response)
	}

	session, exists := f.sessionManager.Session(sessionID)
	if !exists {
		// Session not found, fall back to JSON
		logger.Warnf("[dubbo-go-pixiu] mcp server session not found: %s, falling back to JSON", sessionID)
		return f.sendJSONResponse(ctx, response)
	}

	// Send via SSE
	if err := f.sseHandler.SendSSEMessage(session, response); err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to send SSE message: %v", err)
		// Fall back to JSON
		return f.sendJSONResponse(ctx, response)
	}

	// SSE message sent successfully
	logger.Debugf("[dubbo-go-pixiu] mcp server sent response via SSE for session: %s", sessionID)
	return filter.Stop
}

// sessionExists checks if a session exists
func (f *MCPServerFilter) sessionExists(sessionID string) bool {
	_, exists := f.sessionManager.Session(sessionID)
	return exists
}

// sendBadRequest sends a 400 Bad Request response
func (f *MCPServerFilter) sendBadRequest(ctx *MCPContext, message string) filter.FilterStatus {
	logger.Warnf("[dubbo-go-pixiu] mcp server bad request: %s", message)
	ctx.SendLocalReply(http.StatusBadRequest, []byte(message))
	return filter.Stop
}

// sendMethodNotAllowed sends a 405 Method Not Allowed response
func (f *MCPServerFilter) sendMethodNotAllowed(ctx *MCPContext) filter.FilterStatus {
	logger.Warnf("[dubbo-go-pixiu] mcp server method not allowed: %s", ctx.Request.Method)
	ctx.Writer.Header().Set("Allow", "GET, POST")
	ctx.SendLocalReply(http.StatusMethodNotAllowed, []byte("Method Not Allowed"))
	return filter.Stop
}

// sendNotAcceptable sends a 406 Not Acceptable response
func (f *MCPServerFilter) sendNotAcceptable(ctx *MCPContext, message string) filter.FilterStatus {
	logger.Warnf("[dubbo-go-pixiu] mcp server not acceptable: %s", message)
	ctx.SendLocalReply(http.StatusNotAcceptable, []byte(message))
	return filter.Stop
}

// sendInternalError sends a 500 Internal Server Error response
func (f *MCPServerFilter) sendInternalError(ctx *MCPContext, message string) filter.FilterStatus {
	logger.Errorf("[dubbo-go-pixiu] mcp server internal error: %s", message)
	ctx.SendLocalReply(http.StatusInternalServerError, []byte("Internal Server Error"))
	return filter.Stop
}

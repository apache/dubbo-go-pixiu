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
	"time"
)

import (
	"github.com/mark3labs/mcp-go/mcp"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/router"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/transport"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// FilterFactory and MCPServerFilter types
type (
	// FilterFactory is a factory to create MCP server filters.
	FilterFactory struct {
		cfg      *model.McpServerConfig
		registry *ToolRegistry
		selector router.ToolSelector
	}

	// MCPServerFilter is a filter that handles MCP protocol.
	MCPServerFilter struct {
		cfg               *model.McpServerConfig
		registry          *ToolRegistry
		errorHandler      *ErrorHandler
		responseBuilder   *ResponseBuilder
		sessionManager    *transport.SessionManager
		sseHandler        *transport.SSEHandler
		contentNegotiator *transport.ContentNegotiator
		selector          router.ToolSelector
	}
)

// Apply prepares the MCP server and tool registry.
func (f *FilterFactory) Apply() error {
	// Initialize tool registry (singleton)
	f.registry = GetOrInitRegistry()

	if err := f.registerConfiguredTools(); err != nil {
		return err
	}
	if err := f.registerConfiguredResources(); err != nil {
		return err
	}
	if err := f.registerConfiguredResourceTemplates(); err != nil {
		return err
	}
	if err := f.registerConfiguredPrompts(); err != nil {
		return err
	}

	return f.configureSelector()
}

func (f *FilterFactory) registerConfiguredTools() error {
	if err := router.ValidateTools(f.cfg.Tools); err != nil {
		return fmt.Errorf("invalid mcp tool router metadata: %v", err)
	}

	// Sync statically configured tools into registry (full replace)
	if err := f.registry.ReplaceAllTools(f.cfg.Tools); err != nil {
		return fmt.Errorf("failed to register mcp tools: %v", err)
	}
	for _, tool := range f.cfg.Tools {
		logger.Debugf("[dubbo-go-pixiu] mcp server registered tool '%s' -> cluster:%s", tool.Name, tool.Cluster)
	}

	return nil
}

func (f *FilterFactory) registerConfiguredResources() error {
	// Register statically configured resources
	for _, resource := range f.cfg.Resources {
		if err := f.registry.RegisterResource(resource); err != nil {
			return fmt.Errorf("failed to register resource %s: %v", resource.Name, err)
		}
		logger.Debugf("[dubbo-go-pixiu] mcp server registered resource '%s' -> uri:%s", resource.Name, resource.URI)
	}

	return nil
}

func (f *FilterFactory) registerConfiguredResourceTemplates() error {
	// Register statically configured resource templates
	for _, template := range f.cfg.ResourceTemplates {
		if err := f.registry.RegisterResourceTemplate(template); err != nil {
			return fmt.Errorf("failed to register resource template %s: %v", template.Name, err)
		}
		logger.Debugf("[dubbo-go-pixiu] mcp server registered template '%s' -> pattern:%s", template.Name, template.URITemplate)
	}

	return nil
}

func (f *FilterFactory) registerConfiguredPrompts() error {
	// Register statically configured prompts
	for _, prompt := range f.cfg.Prompts {
		if err := f.registry.RegisterPrompt(prompt); err != nil {
			return fmt.Errorf("failed to register prompt %s: %v", prompt.Name, err)
		}
		logger.Debugf("[dubbo-go-pixiu] mcp server registered prompt '%s'", prompt.Name)
	}

	return nil
}

func (f *FilterFactory) configureSelector() error {
	f.selector = nil
	if f.cfg.Router != nil && f.cfg.Router.Enabled {
		// Build the tool selector lazily so disabled routing preserves the
		// pre-router passthrough behavior with zero plan-store overhead.
		selector, err := router.Build(f.cfg.Router, GetOrInitPlanStoreWithMaxEntries(f.cfg.Router.Session.MaxEntries))
		if err != nil {
			return fmt.Errorf("failed to build mcp tool router: %v", err)
		}
		f.selector = selector
	}

	return nil
}

// Config returns the configuration struct
func (f *FilterFactory) Config() any {
	return f.cfg
}

// PrepareFilterChain prepares the filter chain
func (f *FilterFactory) PrepareFilterChain(_ *contexthttp.HttpContext, chain filter.FilterChain) error {
	// Get global session manager singleton
	sessionManager := GetOrInitSessionManager()
	sseHandler := transport.NewSSEHandler(sessionManager)
	contentNegotiator := transport.NewContentNegotiator()

	// Deep copy config to avoid pointer sharing (factory.cfg may change at runtime)
	mcpFilter := &MCPServerFilter{
		cfg:               f.cfg.DeepCopy(),
		registry:          f.registry,
		errorHandler:      NewErrorHandler(),
		responseBuilder:   NewResponseBuilder(),
		sessionManager:    sessionManager,
		sseHandler:        sseHandler,
		contentNegotiator: contentNegotiator,
		selector:          f.selector,
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

	// Parse HTTP headers
	mcpCtx.ParseAndSetProtocolVersionHeader()
	mcpCtx.ParseAndSetSessionHeader()
	mcpCtx.ParseAndSetAcceptHeader()

	// Protocol version is logged for debugging
	version := mcpCtx.ProtocolVersion()
	if version != "" {
		logger.Debugf("[dubbo-go-pixiu] mcp server client protocol version: %s", version)
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

// handleGetRequest handles HTTP GET requests for SSE stream establishment
func (f *MCPServerFilter) handleGetRequest(ctx *MCPContext) filter.FilterStatus {
	logger.Infof("[dubbo-go-pixiu] mcp server handling GET request for SSE stream")

	// Validate Accept header includes text/event-stream
	if !ctx.AcceptSSE() {
		logger.Warnf("[dubbo-go-pixiu] mcp server GET request must accept text/event-stream")
		return f.sendNotAcceptable(ctx, "GET request must accept text/event-stream")
	}

	sessionIDHeader := ctx.SessionID()
	if sessionIDHeader == "" {
		return f.sendBadRequest(ctx, "Mcp-Session-Id header is required")
	}
	session, exists := f.sessionManager.GetSession(sessionIDHeader)
	if !exists {
		return f.sendNotFound(ctx, "MCP session not found")
	}
	ctx.SetSessionID(session.ID)

	// Create io.Pipe for SSE message transport
	pipeReader, pipeWriter := io.Pipe()
	streamToken := session.AttachStream(pipeWriter)

	// Create virtual HTTP response with pipe as body
	virtualResp := &http.Response{
		StatusCode: http.StatusOK,
		Header: http.Header{
			constant.HeaderKeyContextType:              []string{constant.HeaderValueTextEventStream},
			constant.HeaderKeyCacheControl:             []string{constant.HeaderValueNoCache},
			constant.HeaderKeyConnection:               []string{constant.HeaderValueKeepAlive},
			constant.HeaderKeyMCPSessionId:             []string{session.ID},
			constant.HeaderKeyAccessControlAllowOrigin: []string{constant.HeaderValueAll},
		},
		Body: pipeReader,
	}

	// Set SourceResp to let buildTargetResponse convert to StreamResponse
	ctx.SourceResp = virtualResp
	ctx.StatusCode(http.StatusOK)

	// Start background goroutine to maintain the SSE connection
	go f.maintainSSEPipe(ctx, session, streamToken)
	go f.flushPendingToolsListChanged(session)

	logger.Infof("[dubbo-go-pixiu] mcp server attached SSE stream")

	// Return Stop to skip remaining filters and backend call
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

	if status := f.validateSessionForMethod(ctx, jsonrpcReq.Method); status != filter.Continue {
		return status
	}

	// Determine response format based on content negotiation
	sessionID := ctx.SessionID()
	hasSession := sessionID != "" && f.sessionExists(sessionID)
	responseFormat := f.contentNegotiator.NegotiateResponse(
		ctx.Request.Header.Get(constant.HeaderKeyAccept), hasSession)

	// Store response format decision in context for later use
	ctx.StoreMCPDataInParams()

	// Handle different request types
	if f.isTerminalMethod(jsonrpcReq.Method) {
		return f.handleTerminalMethodWithNegotiation(ctx, jsonrpcReq, responseFormat)
	} else if jsonrpcReq.Method == string(mcp.MethodToolsCall) {
		return f.handleToolCall(ctx, jsonrpcReq)
	} else {
		// Unknown method
		logger.Warnf("[dubbo-go-pixiu] mcp server unsupported method: %s", jsonrpcReq.Method)
		return f.errorHandler.SendMethodNotFound(ctx, jsonrpcReq.ID)
	}
}

// handleTerminalMethodWithNegotiation handles terminal methods with response format negotiation
func (f *MCPServerFilter) handleTerminalMethodWithNegotiation(ctx *MCPContext, req mcp.JSONRPCRequest, responseFormat transport.ResponseFormat) filter.FilterStatus {
	// Special case 1: initialize always returns JSON immediately (no SSE option per MCP spec)
	if req.Method == string(mcp.MethodInitialize) {
		return f.handleInitialize(ctx, req)
	}

	// Special case 2: notifications/initialized always returns 202 Accepted with no body
	if req.Method == "notifications/initialized" {
		logger.Infof("[dubbo-go-pixiu] mcp server received initialized notification, returning 202 Accepted")
		ctx.SendLocalReply(http.StatusAccepted, nil)
		return filter.Stop
	}

	// For other terminal methods, dispatch to specific handlers with response format
	switch req.Method {
	case string(mcp.MethodToolsList):
		return f.handleToolsList(ctx, req, responseFormat)
	case string(mcp.MethodResourcesList):
		return f.handleResourcesList(ctx, req, responseFormat)
	case string(mcp.MethodResourcesRead):
		return f.handleResourceRead(ctx, req, responseFormat)
	case "resources/templates/list":
		return f.handleResourceTemplatesList(ctx, req, responseFormat)
	case string(mcp.MethodPromptsList):
		return f.handlePromptsList(ctx, req, responseFormat)
	case string(mcp.MethodPromptsGet):
		return f.handlePromptsGet(ctx, req, responseFormat)
	case string(mcp.MethodPing):
		return f.handlePing(ctx, req, responseFormat)
	default:
		logger.Warnf("[dubbo-go-pixiu] mcp server unsupported terminal method: %s", req.Method)
		return f.errorHandler.SendMethodNotFound(ctx, req.ID)
	}
}

// sendResponseWithFormat sends response in the negotiated format
func (f *MCPServerFilter) sendResponseWithFormat(ctx *MCPContext, response any, format transport.ResponseFormat) filter.FilterStatus {
	switch format {
	case transport.ResponseFormatJSON:
		return f.sendJSONResponse(ctx, response)
	case transport.ResponseFormatSSE:
		return f.sendSSEResponse(ctx, response)
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
		logger.Warnf("[dubbo-go-pixiu] mcp server session not found, falling back to JSON")
		return f.sendJSONResponse(ctx, response)
	}

	// Send via SSE
	if err := f.sseHandler.SendSSEMessage(session, response); err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to send SSE message: %v", err)
		// Fall back to JSON
		return f.sendJSONResponse(ctx, response)
	}

	// SSE message sent successfully, return 202 Accepted per MCP spec
	logger.Debugf("[dubbo-go-pixiu] mcp server sent response via SSE")
	ctx.SendLocalReply(http.StatusAccepted, nil)
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

// sendNotFound sends a 404 Not Found response
func (f *MCPServerFilter) sendNotFound(ctx *MCPContext, message string) filter.FilterStatus {
	logger.Warnf("[dubbo-go-pixiu] mcp server not found: %s", message)
	ctx.SendLocalReply(http.StatusNotFound, []byte(message))
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

// maintainSSEPipe maintains the SSE pipe connection with keepalive
func (f *MCPServerFilter) maintainSSEPipe(ctx *MCPContext, session *transport.MCPSession, token transport.StreamToken) {
	// Ensure cleanup on exit
	defer func() {
		session.DetachStream(token)
		logger.Debugf("[dubbo-go-pixiu] mcp server detached SSE stream")
	}()

	ticker := time.NewTicker(transport.KeepaliveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Send keepalive comment (ignored by SSE clients)
			keepalive := f.sseHandler.FormatSSEKeepalive(time.Now().Unix())
			if err := session.WriteSSEDataForStream(token, []byte(keepalive), time.Now()); err != nil {
				logger.Warnf("[dubbo-go-pixiu] mcp server keepalive write failed: %v", err)
				return
			}
			logger.Debugf("[dubbo-go-pixiu] mcp server sent keepalive")

		case <-session.Done:
			// Server-initiated close
			logger.Infof("[dubbo-go-pixiu] mcp server closing SSE stream (server initiated)")
			return

		case <-ctx.Ctx.Done():
			// Client disconnected
			logger.Infof("[dubbo-go-pixiu] mcp server closing SSE stream (client disconnected)")
			return
		}
	}
}

// SendServerNotification sends a server notification to the client via SSE stream
func (f *MCPServerFilter) SendServerNotification(sessionID string, method string, params map[string]any) error {
	session, exists := f.sessionManager.Session(sessionID)
	if !exists {
		return fmt.Errorf("session not found")
	}

	if !session.HasPipeWriter() {
		return fmt.Errorf("SSE pipe not established")
	}

	// Use ResponseBuilder to create notification
	notification := f.responseBuilder.ServerNotification(method, params)

	// Send via SSE
	if err := f.sseHandler.SendSSEMessage(session, notification); err != nil {
		return err
	}

	logger.Debugf("[dubbo-go-pixiu] mcp server sent notification: %s", method)
	return nil
}

// SendServerRequest sends a server request to the client via SSE stream
func (f *MCPServerFilter) SendServerRequest(sessionID string, id any, method string, params any) error {
	session, exists := f.sessionManager.Session(sessionID)
	if !exists {
		return fmt.Errorf("session not found")
	}

	if !session.HasPipeWriter() {
		return fmt.Errorf("SSE pipe not established")
	}

	// Use ResponseBuilder to create request
	request := f.responseBuilder.ServerRequest(id, method, params)

	// Send via SSE
	if err := f.sseHandler.SendSSEMessage(session, request); err != nil {
		return err
	}

	logger.Debugf("[dubbo-go-pixiu] mcp server sent request: %s (id: %v)", method, id)
	return nil
}

func (f *MCPServerFilter) validateSessionForMethod(ctx *MCPContext, method string) filter.FilterStatus {
	if method == string(mcp.MethodInitialize) {
		return filter.Continue
	}

	sessionID := ctx.SessionID()
	if sessionID != "" {
		if _, exists := f.sessionManager.GetSession(sessionID); !exists {
			return f.sendNotFound(ctx, "MCP session not found")
		}
		return filter.Continue
	}

	if f.selector != nil && (method == string(mcp.MethodToolsList) || method == string(mcp.MethodToolsCall)) {
		return f.sendBadRequest(ctx, "Mcp-Session-Id header is required")
	}
	return filter.Continue
}

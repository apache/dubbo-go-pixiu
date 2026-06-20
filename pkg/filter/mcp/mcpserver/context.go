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
	"strings"
)

import (
	"github.com/mark3labs/mcp-go/mcp"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/router"
)

const MCPDataKey = "mcp_data"

// MCPData stores MCP-related data
type MCPData struct {
	// Method stores MCP method name
	Method string
	// RequestID stores JSON-RPC request ID
	RequestID any
	// ToolName stores the requested tool name for tools/call requests
	ToolName string
	// SessionID stores MCP session ID for SSE connections
	SessionID string
	// AcceptSSE indicates if client accepts text/event-stream
	AcceptSSE bool
	// AcceptJSON indicates if client accepts application/json
	AcceptJSON bool
	// ProtocolVersion stores MCP protocol version from header
	ProtocolVersion string
	// AuthorizationReceipt binds a tools/call completion to its authorization.
	AuthorizationReceipt *router.AuthorizationReceipt
}

// MCPContext MCP context wrapper that composes HttpContext and provides MCP-specific operations
type MCPContext struct {
	*contexthttp.HttpContext
	mcpData *MCPData
}

// NewMCPContext creates a new MCP context
func NewMCPContext(httpCtx *contexthttp.HttpContext) *MCPContext {
	return &MCPContext{
		HttpContext: httpCtx,
		mcpData:     &MCPData{},
	}
}

// IsMCPRequest checks if it's an MCP request (by method name)
func (ctx *MCPContext) IsMCPRequest() bool {
	return ctx.mcpData.Method != ""
}

// SetMCPMethod sets MCP method name
func (ctx *MCPContext) SetMCPMethod(method string) {
	ctx.mcpData.Method = method
}

// McpMethod gets MCP method name
func (ctx *MCPContext) McpMethod() string {
	return ctx.mcpData.Method
}

// SetMCPRequestID sets JSON-RPC request ID
func (ctx *MCPContext) SetMCPRequestID(id any) {
	ctx.mcpData.RequestID = id
}

// McpRequestID gets JSON-RPC request ID
func (ctx *MCPContext) McpRequestID() any {
	return ctx.mcpData.RequestID
}

// SetMCPToolName stores the current tools/call tool name.
func (ctx *MCPContext) SetMCPToolName(name string) {
	ctx.mcpData.ToolName = name
}

// McpToolName gets the current tools/call tool name.
func (ctx *MCPContext) McpToolName() string {
	return ctx.mcpData.ToolName
}

func (ctx *MCPContext) SetAuthorizationReceipt(receipt *router.AuthorizationReceipt) {
	ctx.mcpData.AuthorizationReceipt = receipt
}

func (ctx *MCPContext) AuthorizationReceipt() *router.AuthorizationReceipt {
	return ctx.mcpData.AuthorizationReceipt
}

// IsMCPToolCall checks if it's a tool call request (by method name)
func (ctx *MCPContext) IsMCPToolCall() bool {
	return ctx.mcpData.Method == string(mcp.MethodToolsCall)
}

// StoreMCPDataInParams stores MCP data in HttpContext.Params for passing through the filter chain
func (ctx *MCPContext) StoreMCPDataInParams() {
	if ctx.Params == nil {
		ctx.Params = make(map[string]any)
	}
	ctx.Params[MCPDataKey] = ctx.mcpData
}

// LoadMCPDataFromParams loads MCP data from HttpContext.Params
func (ctx *MCPContext) LoadMCPDataFromParams() {
	if ctx.Params == nil {
		return
	}
	if data, ok := ctx.Params[MCPDataKey].(*MCPData); ok {
		ctx.mcpData = data
	}
}

// NewMCPContextFromHttpContext creates MCPContext from existing HttpContext and tries to load stored MCP data
func NewMCPContextFromHttpContext(httpCtx *contexthttp.HttpContext) *MCPContext {
	mcpCtx := NewMCPContext(httpCtx)
	mcpCtx.LoadMCPDataFromParams()
	return mcpCtx
}

// ClearContentLengthHeader removes the Content-Length header to prevent conflicts with chunked transfer encoding.
func (ctx *MCPContext) ClearContentLengthHeader() {
	ctx.Writer.Header().Del(constant.HeaderKeyContentLength)
}

// SessionID management methods

// SetSessionID sets MCP session ID
func (ctx *MCPContext) SetSessionID(sessionID string) {
	ctx.mcpData.SessionID = sessionID
}

// SessionID gets MCP session ID
func (ctx *MCPContext) SessionID() string {
	return ctx.mcpData.SessionID
}

// HasSession checks if context has a session ID
func (ctx *MCPContext) HasSession() bool {
	return ctx.mcpData.SessionID != ""
}

// Accept header management methods

// SetAcceptSSE sets if client accepts text/event-stream
func (ctx *MCPContext) SetAcceptSSE(acceptSSE bool) {
	ctx.mcpData.AcceptSSE = acceptSSE
}

// AcceptSSE gets if client accepts text/event-stream
func (ctx *MCPContext) AcceptSSE() bool {
	return ctx.mcpData.AcceptSSE
}

// SetAcceptJSON sets if client accepts application/json
func (ctx *MCPContext) SetAcceptJSON(acceptJSON bool) {
	ctx.mcpData.AcceptJSON = acceptJSON
}

// AcceptJSON gets if client accepts application/json
func (ctx *MCPContext) AcceptJSON() bool {
	return ctx.mcpData.AcceptJSON
}

// Protocol version management methods

// SetProtocolVersion sets MCP protocol version
func (ctx *MCPContext) SetProtocolVersion(version string) {
	ctx.mcpData.ProtocolVersion = version
}

// ProtocolVersion gets MCP protocol version
func (ctx *MCPContext) ProtocolVersion() string {
	return ctx.mcpData.ProtocolVersion
}

// ParseAndSetAcceptHeader parses Accept header and sets AcceptSSE/AcceptJSON flags
func (ctx *MCPContext) ParseAndSetAcceptHeader() {
	acceptHeader := ctx.Request.Header.Get(constant.HeaderKeyAccept)
	ctx.mcpData.AcceptJSON = acceptHeader == "" || // Default to JSON for backward compatibility
		strings.Contains(acceptHeader, constant.HeaderValueApplicationJson) ||
		strings.Contains(acceptHeader, constant.MediaTypeApplicationWild) ||
		strings.Contains(acceptHeader, constant.MediaTypeWildcard)

	ctx.mcpData.AcceptSSE = strings.Contains(acceptHeader, constant.HeaderValueTextEventStream) ||
		strings.Contains(acceptHeader, constant.MediaTypeTextWild) ||
		strings.Contains(acceptHeader, constant.MediaTypeWildcard)
}

// ParseAndSetSessionHeader parses Mcp-Session-Id header and sets session ID
func (ctx *MCPContext) ParseAndSetSessionHeader() {
	sessionID := ctx.Request.Header.Get(constant.HeaderKeyMCPSessionId)
	ctx.mcpData.SessionID = sessionID
}

// ParseAndSetProtocolVersionHeader parses MCP-Protocol-Version header
func (ctx *MCPContext) ParseAndSetProtocolVersionHeader() {
	version := ctx.Request.Header.Get(constant.HeaderKeyMCPProtocolVersion)
	ctx.mcpData.ProtocolVersion = version
}

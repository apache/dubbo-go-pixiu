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
	"github.com/mark3labs/mcp-go/mcp"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

const MCPDataKey = "mcp_data"

// MCPData stores MCP-related data
type MCPData struct {
	// Method stores MCP method name
	Method string
	// RequestID stores JSON-RPC request ID
	RequestID any
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

// IsMCPToolCall checks if it's a tool call request (by method name)
func (ctx *MCPContext) IsMCPToolCall() bool {
	return ctx.mcpData.Method == string(mcp.MethodToolsCall)
}

// StoreMCPDataInParams stores MCP data in HttpContext.Params for passing through the filter chain
func (ctx *MCPContext) StoreMCPDataInParams() {
	if ctx.HttpContext.Params == nil {
		ctx.HttpContext.Params = make(map[string]any)
	}
	ctx.HttpContext.Params[MCPDataKey] = ctx.mcpData
}

// LoadMCPDataFromParams loads MCP data from HttpContext.Params
func (ctx *MCPContext) LoadMCPDataFromParams() {
	if ctx.HttpContext.Params == nil {
		return
	}
	if data, ok := ctx.HttpContext.Params[MCPDataKey].(*MCPData); ok {
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

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

package http

import (
	"context"
	"encoding/json"
	"math"
	"net"
	"net/http"
	"net/url"
	"strings"
	"time"

	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/router"
)

const abortIndex int8 = math.MaxInt8 / 2

// MCP 相关的上下文键常量
const (
	// MCPRequestKey 标识这是一个 MCP 请求
	MCPRequestKey = "mcp.request"
	// MCPMethodKey 存储 MCP 方法名
	MCPMethodKey = "mcp.method"
	// MCPRequestIDKey 存储 JSON-RPC 请求 ID
	MCPRequestIDKey = "mcp.request_id"
	// MCPToolCallKey 标识这是一个工具调用请求
	MCPToolCallKey = "mcp.tool_call"
	// MCPToolNameKey 存储工具名称
	MCPToolNameKey = "mcp.tool_name"
	// MCPClusterKey 存储目标集群信息
	MCPClusterKey = "mcp.cluster"
	// MCPProcessedKey 标识请求已被处理过
	MCPProcessedKey = "mcp.processed"
)

// HttpContext http context
type HttpContext struct {
	//Deprecated: waiting to delete
	Index int8
	//Deprecated: waiting to delete
	Filters FilterChain
	Timeout time.Duration
	Ctx     context.Context
	Params  map[string]any

	// localReply Means that the request was interrupted,
	// which may occur in the Decode or Encode stage
	localReply bool
	// statusCode code will be return
	statusCode int
	// localReplyBody: happen error
	localReplyBody []byte
	// the response context will return.
	TargetResp any
	// client call response.
	SourceResp any

	HttpConnectionManager model.HttpConnectionManagerConfig
	Route                 *model.RouteAction
	Api                   *router.API

	Request *http.Request
	Writer  http.ResponseWriter
}

type (
	// ErrResponse err response.
	ErrResponse struct {
		Message string `json:"message"`
	}

	// FilterFunc filter func, filter
	FilterFunc func(c *HttpContext)

	// FilterChain filter chain
	FilterChain []FilterFunc
)

// Deprecated: Next logic for lookup filter
func (hc *HttpContext) Next() {
}

// Reset reset http context
func (hc *HttpContext) Reset() {
	hc.Ctx = nil
	hc.Index = -1
	hc.Filters = []FilterFunc{}
	hc.Route = nil
	hc.Api = nil

	hc.TargetResp = nil
	hc.SourceResp = nil
	hc.statusCode = 0
	hc.localReply = false
	hc.localReplyBody = nil
}

// RouteEntry set route
func (hc *HttpContext) RouteEntry(r *model.RouteAction) {
	hc.Route = r
}

// GetRouteEntry get route
func (hc *HttpContext) GetRouteEntry() *model.RouteAction {
	return hc.Route
}

// AddHeader add header
func (hc *HttpContext) AddHeader(k, v string) {
	hc.Writer.Header().Add(k, v)
}

// GetHeader get header
func (hc *HttpContext) GetHeader(k string) string {
	return hc.Request.Header.Get(k)
}

// AllHeaders  get all headers
func (hc *HttpContext) AllHeaders() http.Header {
	return hc.Request.Header
}

// GetUrl get http request url
func (hc *HttpContext) GetUrl() string {
	return hc.Request.URL.Path
}

func (hc *HttpContext) SetUrl(url string) {
	hc.Request.URL.Path = url
}

// GetMethod get method, POST/GET ...
func (hc *HttpContext) GetMethod() string {
	return hc.Request.Method
}

// GetClientIP get client IP
func (hc *HttpContext) GetClientIP() string {
	xForwardedFor := hc.Request.Header.Get("X-Forwarded-For")
	ip := strings.TrimSpace(strings.Split(xForwardedFor, ",")[0])
	if len(ip) != 0 {
		return ip
	}

	ip = strings.TrimSpace(hc.Request.Header.Get("X-Real-Ip"))
	if len(ip) != 0 {
		return ip
	}

	if ip, _, err := net.SplitHostPort(strings.TrimSpace(hc.Request.RemoteAddr)); err == nil && len(ip) != 0 {
		return ip
	}

	return ""
}

// GetApplicationName get application name
func (hc *HttpContext) GetApplicationName() string {
	if u, err := url.Parse(hc.Request.RequestURI); err == nil {
		return strings.Split(u.Path, "/")[0]
	}

	return ""
}

// SendLocalReply Means that the request was interrupted and UnaryResponse will be sent directly
// Even if it’s currently in to Decode stage
func (hc *HttpContext) SendLocalReply(status int, body []byte) {
	hc.localReply = true
	hc.statusCode = status
	hc.localReplyBody = body
	hc.TargetResp = &client.UnaryResponse{Data: body}
	if json.Valid(body) {
		hc.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueApplicationJson)
	} else {
		hc.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
	}
	writer := hc.Writer
	writer.WriteHeader(status)
	_, err := writer.Write(body)
	if err != nil {
		logger.Errorf("write fail: %v", err)
	}
}

func (hc *HttpContext) GetLocalReplyBody() []byte {
	return hc.localReplyBody
}

func (hc *HttpContext) GetStatusCode() int {
	return hc.statusCode
}

func (hc *HttpContext) StatusCode(code int) {
	hc.statusCode = code
}

func (hc *HttpContext) LocalReply() bool {
	return hc.localReply
}

// API sets the API to http context
func (hc *HttpContext) API(api router.API) {
	if hc.Timeout > api.Timeout {
		hc.Timeout = api.Timeout
	}
	hc.Api = &api
}

// GetAPI get api
func (hc *HttpContext) GetAPI() *router.API {
	return hc.Api
}

// Deprecated: Abort  filter chain break , filter after the current filter will not executed.
func (hc *HttpContext) Abort() {
	hc.Index = abortIndex
}

// Deprecated: AppendFilterFunc append filter func.
func (hc *HttpContext) AppendFilterFunc(ff ...FilterFunc) {
	for _, v := range ff {
		hc.Filters = append(hc.Filters, v)
	}
}

func (hc *HttpContext) GenerateHash() string {
	req := hc.Request
	return req.Method + "." + req.RequestURI
}

// MCP 相关的辅助方法

// SetMCPRequest 标记这是一个 MCP 请求
func (hc *HttpContext) SetMCPRequest(isMCP bool) {
	if hc.Params == nil {
		hc.Params = make(map[string]any)
	}
	hc.Params[MCPRequestKey] = isMCP
}

// IsMCPRequest 检查是否是 MCP 请求
func (hc *HttpContext) IsMCPRequest() bool {
	if hc.Params == nil {
		return false
	}
	if val, exists := hc.Params[MCPRequestKey]; exists {
		if isMCP, ok := val.(bool); ok {
			return isMCP
		}
	}
	return false
}

// SetMCPMethod 设置 MCP 方法名
func (hc *HttpContext) SetMCPMethod(method string) {
	if hc.Params == nil {
		hc.Params = make(map[string]any)
	}
	hc.Params[MCPMethodKey] = method
}

// GetMCPMethod 获取 MCP 方法名
func (hc *HttpContext) GetMCPMethod() string {
	if hc.Params == nil {
		return ""
	}
	if val, exists := hc.Params[MCPMethodKey]; exists {
		if method, ok := val.(string); ok {
			return method
		}
	}
	return ""
}

// SetMCPRequestID 设置 JSON-RPC 请求 ID
func (hc *HttpContext) SetMCPRequestID(id any) {
	if hc.Params == nil {
		hc.Params = make(map[string]any)
	}
	hc.Params[MCPRequestIDKey] = id
}

// GetMCPRequestID 获取 JSON-RPC 请求 ID
func (hc *HttpContext) GetMCPRequestID() any {
	if hc.Params == nil {
		return nil
	}
	return hc.Params[MCPRequestIDKey]
}

// SetMCPToolCall 标记这是一个工具调用请求
func (hc *HttpContext) SetMCPToolCall(isToolCall bool) {
	if hc.Params == nil {
		hc.Params = make(map[string]any)
	}
	hc.Params[MCPToolCallKey] = isToolCall
}

// IsMCPToolCall 检查是否是工具调用请求
func (hc *HttpContext) IsMCPToolCall() bool {
	if hc.Params == nil {
		return false
	}
	if val, exists := hc.Params[MCPToolCallKey]; exists {
		if isToolCall, ok := val.(bool); ok {
			return isToolCall
		}
	}
	return false
}

// SetMCPToolName 设置工具名称
func (hc *HttpContext) SetMCPToolName(toolName string) {
	if hc.Params == nil {
		hc.Params = make(map[string]any)
	}
	hc.Params[MCPToolNameKey] = toolName
}

// GetMCPToolName 获取工具名称
func (hc *HttpContext) GetMCPToolName() string {
	if hc.Params == nil {
		return ""
	}
	if val, exists := hc.Params[MCPToolNameKey]; exists {
		if toolName, ok := val.(string); ok {
			return toolName
		}
	}
	return ""
}

// SetMCPCluster 设置目标集群信息
func (hc *HttpContext) SetMCPCluster(cluster string) {
	if hc.Params == nil {
		hc.Params = make(map[string]any)
	}
	hc.Params[MCPClusterKey] = cluster
}

// GetMCPCluster 获取目标集群信息
func (hc *HttpContext) GetMCPCluster() string {
	if hc.Params == nil {
		return ""
	}
	if val, exists := hc.Params[MCPClusterKey]; exists {
		if cluster, ok := val.(string); ok {
			return cluster
		}
	}
	return ""
}

// SetMCPProcessed 标记请求已被处理过
func (hc *HttpContext) SetMCPProcessed(processed bool) {
	if hc.Params == nil {
		hc.Params = make(map[string]any)
	}
	hc.Params[MCPProcessedKey] = processed
}

// IsMCPProcessed 检查请求是否已被处理过
func (hc *HttpContext) IsMCPProcessed() bool {
	if hc.Params == nil {
		return false
	}
	if val, exists := hc.Params[MCPProcessedKey]; exists {
		if processed, ok := val.(bool); ok {
			return processed
		}
	}
	return false
}

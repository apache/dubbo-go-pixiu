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
)

import (
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/router"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

const abortIndex int8 = math.MaxInt8 / 2

// HttpContext http context
type HttpContext struct {
	//Deprecated: waiting to delete
	Index int8
	//Deprecated: waiting to delete
	Filters FilterChain
	Timeout time.Duration
	Ctx     context.Context
	Params  map[string]interface{}

	// localReply Means that the request was interrupted,
	// which may occur in the Decode or Encode stage
	localReply bool
	// statusCode code will be return
	statusCode int
	// localReplyBody: happen error
	localReplyBody []byte
	// the response context will return.
	TargetResp *client.Response
	// client call response.
	SourceResp interface{}

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

// SendLocalReply Means that the request was interrupted and Response will be sent directly
// Even if it’s currently in to Decode stage
func (hc *HttpContext) SendLocalReply(status int, body []byte) {
	hc.localReply = true
	hc.statusCode = status
	hc.localReplyBody = body
	hc.TargetResp = &client.Response{Data: body}
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

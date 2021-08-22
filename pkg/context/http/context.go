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
	"encoding/json"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"
	"net"
	"net/http"
	"net/url"
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/context"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// HttpContext http context
type HttpContext struct {
	*context.BaseContext
	HttpConnectionManager model.HttpConnectionManager
	FilterChains          []model.FilterChain
	Listener              *model.Listener
	Route                 *model.RouteAction
	Api                   router.API

	Request   *http.Request
	writermem responseWriter
	Writer    ResponseWriter
}

// Next logic for lookup filter
func (hc *HttpContext) Next() {
	hc.Index++
	for hc.Index < int8(len(hc.Filters)) {
		hc.Filters[hc.Index](hc)
		hc.Index++
	}
}

// Reset reset http context
func (hc *HttpContext) Reset() {
	hc.Writer = &hc.writermem
	hc.BaseContext = context.NewBaseContext()
}

// Status set header status code
func (hc *HttpContext) Status(code int) {
	hc.Writer.WriteHeader(code)
}

// StatusCode get header status code
func (hc *HttpContext) StatusCode() int {
	return hc.Writer.Status()
}

// Write write body data
func (hc *HttpContext) Write(b []byte) (int, error) {
	return hc.Writer.Write(b)
}

// WriteHeaderNow write header now
func (hc *HttpContext) WriteHeaderNow() {
	hc.writermem.WriteHeaderNow()
}

// WriteWithStatus status must set first
func (hc *HttpContext) WriteWithStatus(code int, b []byte) (int, error) {
	hc.Writer.WriteHeader(code)
	return hc.Writer.Write(b)
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

// WriteJSONWithStatus write fail, auto add context-type json.
func (hc *HttpContext) WriteJSONWithStatus(code int, res interface{}) {
	hc.doWriteJSON(nil, code, res)
}

// WriteErr
func (hc *HttpContext) WriteErr(p interface{}) {
	hc.doWriteJSON(nil, http.StatusInternalServerError, p)
}

// WriteSuccess
func (hc *HttpContext) WriteSuccess() {
	hc.doWriteJSON(nil, http.StatusOK, nil)
}

// WriteResponse
func (hc *HttpContext) WriteResponse(resp client.Response) {
	hc.doWriteJSON(nil, http.StatusOK, resp.Data)
}

func (hc *HttpContext) doWriteJSON(h map[string]string, code int, d interface{}) {
	if h == nil {
		h = make(map[string]string, 1)
	}
	h[constant.HeaderKeyContextType] = constant.HeaderValueJsonUtf8
	hc.doWrite(h, code, d)
}

func (hc *HttpContext) doWrite(h map[string]string, code int, d interface{}) {
	for k, v := range h {
		hc.Writer.Header().Set(k, v)
	}

	hc.Writer.WriteHeader(code)

	if d != nil {
		byt, ok := d.([]byte)
		if ok {
			hc.Writer.Write(byt)
			return
		}

		if b, err := json.Marshal(d); err != nil {
			hc.Writer.Write([]byte(err.Error()))
		} else {
			hc.Writer.Write(b)
		}
	}
}

// ResetWritermen reset writermen
func (hc *HttpContext) ResetWritermen(w http.ResponseWriter) {
	hc.writermem.reset(w)
}

// API sets the API to http context
func (hc *HttpContext) API(api router.API) {
	hc.Timeout = api.Timeout
	hc.Api = api
}

// GetAPI get api
func (hc *HttpContext) GetAPI() *router.API {
	return &hc.Api
}
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
	"net/http"
	"sync"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
	"github.com/dubbogo/dubbo-go-proxy/pkg/router"
)

// HttpContext http context
type HttpContext struct {
	*context.BaseContext
	HttpConnectionManager model.HttpConnectionManager
	FilterChains          []model.FilterChain
	Listener              *model.Listener
	api                   router.API

	Request   *http.Request
	writermem responseWriter
	Writer    ResponseWriter
	Lock      sync.Mutex
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
	hc.Filters = nil
	hc.Index = -1
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

// WriteHeaderNow
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

// GetUrl get http request url
func (hc *HttpContext) GetUrl() string {
	return hc.Request.URL.Path
}

// GetMethod get method, POST/GET ...
func (hc *HttpContext) GetMethod() string {
	return hc.Request.Method
}

// Api
func (hc *HttpContext) Api(api *model.Api) {
	// hc.api = api
}

// API sets the API to http context
func (hc *HttpContext) API(api router.API) {
	hc.api = api
}

// GetAPI get api
func (hc *HttpContext) GetAPI() *router.API {
	return &hc.api
}

// WriteFail
func (hc *HttpContext) WriteFail() {
	hc.doWriteJSON(nil, http.StatusInternalServerError, nil)
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
		if b, err := json.Marshal(d); err != nil {
			hc.Writer.Write([]byte(err.Error()))
		} else {
			hc.Writer.Write(b)
		}
	}
}

// BuildFilters build filter, from config http_filters
func (hc *HttpContext) BuildFilters() {
	var filterFuncs []context.FilterFunc
	api := hc.GetAPI()

	if api == nil {
		return
	}
	for _, v := range api.Method.Filters {
		filterFuncs = append(filterFuncs, extension.GetMustFilterFunc(v))
	}

	switch api.Method.IntegrationRequest.RequestType {
	case config.DubboRequest:
		hc.AppendFilterFunc(extension.GetMustFilterFunc(constant.HttpTransferDubboFilter))
	case config.HTTPRequest:
		break
	}

	hc.AppendFilterFunc(filterFuncs...)
}

// ResetWritermen reset writermen
func (hc *HttpContext) ResetWritermen(w http.ResponseWriter) {
	hc.writermem.reset(w)
}

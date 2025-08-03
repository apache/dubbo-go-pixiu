/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 * http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package resolver

import (
	"net/http"
	"net/http/httptest"
	"testing"
)

import (
	apiConf "github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"

	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
)

// TestBaseResolver_PreCheck tests the PreCheck method of the BaseResolver.
func TestBaseResolver_PreCheck(t *testing.T) {
	tests := []struct {
		name        string
		setupReq    func() *http.Request
		expectError bool
		errorMsg    string
	}{
		{
			name: "Valid Request",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("POST", "/app/service/method", nil)
				req.Header.Set(constant.DubboHttpDubboVersion, "1.0.0")
				return req
			},
			expectError: false,
		},
		{
			name: "Invalid HTTP Method (GET)",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("GET", "/app/service/method", nil)
				req.Header.Set(constant.DubboHttpDubboVersion, "1.0.0")
				return req
			},
			expectError: true,
			errorMsg:    "http request must be POST and have x-dubbo-http1.1-dubbo-version header",
		},
		{
			name: "Missing Dubbo Version Header",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("POST", "/app/service/method", nil)
				return req
			},
			expectError: true,
			errorMsg:    "http request must be POST and have x-dubbo-http1.1-dubbo-version header",
		},
		{
			name: "Invalid Path (too short)",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("POST", "/app/service", nil)
				req.Header.Set(constant.DubboHttpDubboVersion, "1.0.0")
				return req
			},
			expectError: true,
			errorMsg:    "http request path must be in {application}/{service}/{method} format",
		},
		{
			name: "Invalid Path (too long)",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("POST", "/app/service/method/extra", nil)
				req.Header.Set(constant.DubboHttpDubboVersion, "1.0.0")
				return req
			},
			expectError: true,
			errorMsg:    "http request path must be in {application}/{service}/{method} format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &BaseResolver{}
			req := tt.setupReq()
			err := resolver.PreCheck(req)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}
		})
	}
}

// TestBaseResolver_BuildAPI tests the BuildAPI method of the BaseResolver.
func TestBaseResolver_BuildAPI(t *testing.T) {
	sampleMappingParams := []apiConf.MappingParam{
		{Name: "requestBody.name", MapTo: "opt.name"},
	}

	tests := []struct {
		name                string
		setupReq            func() *http.Request
		mappingParams       []apiConf.MappingParam
		expectError         bool
		expectedRequestType apiConf.RequestType
		expectedVersion     string
		expectedGroup       string
		errorMsg            string
	}{
		{
			name: "Dubbo Protocol Request",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("POST", "/app/service/method", nil)
				req.Header.Set(constant.DubboServiceProtocol, string(apiConf.DubboRequest))
				req.Header.Set(constant.DubboServiceVersion, "1.0.0")
				req.Header.Set(constant.DubboGroup, "test-group")
				return req
			},
			mappingParams:       sampleMappingParams,
			expectError:         false,
			expectedRequestType: apiConf.DubboRequest,
			expectedVersion:     "1.0.0",
			expectedGroup:       "test-group",
		},
		{
			name: "Triple Protocol Request",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("POST", "/app/service/method", nil)
				req.Header.Set(constant.DubboServiceProtocol, "triple")
				return req
			},
			mappingParams:       sampleMappingParams,
			expectError:         false,
			expectedRequestType: "triple",
		},
		{
			name: "HTTP Protocol Request",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("POST", "/app/service/method", nil)
				req.Header.Set(constant.DubboServiceProtocol, string(apiConf.HTTPRequest))
				return req
			},
			mappingParams:       sampleMappingParams,
			expectError:         false,
			expectedRequestType: apiConf.HTTPRequest,
		},
		{
			name: "No Protocol Specified (Defaults to Dubbo)",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("POST", "/app/service/method", nil)
				return req
			},
			mappingParams:       sampleMappingParams,
			expectError:         false,
			expectedRequestType: apiConf.DubboRequest,
		},
		{
			name: "Unknown Protocol",
			setupReq: func() *http.Request {
				req := httptest.NewRequest("POST", "/app/service/method", nil)
				req.Header.Set(constant.DubboServiceProtocol, "unknown-protocol")
				return req
			},
			mappingParams: sampleMappingParams,
			expectError:   true,
			errorMsg:      "http request has unknown protocol in x-dubbo-service-protocol",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &BaseResolver{}
			req := tt.setupReq()
			api, err := resolver.BuildAPI(req, tt.mappingParams)

			if tt.expectError {
				assert.Error(t, err)
				assert.Nil(t, api)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, api)
				assert.Equal(t, "/:application/:interface/:method", api.URLPattern)
				assert.Equal(t, string(http.MethodPost), string(api.Method.HTTPVerb))
				assert.True(t, api.Method.Enable)
				assert.Equal(t, apiConf.HTTPRequest, api.Method.InboundRequest.RequestType)
				assert.Equal(t, tt.expectedRequestType, api.Method.IntegrationRequest.RequestType)
				assert.Equal(t, tt.expectedVersion, api.Method.IntegrationRequest.DubboBackendConfig.Version)
				assert.Equal(t, tt.expectedGroup, api.Method.IntegrationRequest.DubboBackendConfig.Group)
				assert.Equal(t, tt.mappingParams, api.Method.IntegrationRequest.MappingParams)
			}
		})
	}
}

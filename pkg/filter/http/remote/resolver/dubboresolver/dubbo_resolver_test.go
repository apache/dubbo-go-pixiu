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

package dubboresolver

import (
	"context"
	"net/http/httptest"
	"testing"
)

import (
	apiConf "github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/router"

	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

// TestStandardDubboResolver_Resolve tests the Resolve method of the StandardDubboResolver.
func TestStandardDubboResolver_Resolve(t *testing.T) {
	tests := []struct {
		name        string
		setupCtx    func() *contexthttp.HttpContext
		expectError bool
		expectAPI   bool
		checkAPI    func(t *testing.T, api *router.API)
		errorMsg    string
	}{
		{
			name: "Successful Resolution",
			setupCtx: func() *contexthttp.HttpContext {
				req := httptest.NewRequest("POST", "/my-app/my.service.name/myMethod", nil)
				req.Header.Set(constant.DubboHttpDubboVersion, "1.0.0")
				ctx := contexthttp.HttpContext{Ctx: context.Background()}
				ctx.Request = req
				return &ctx
			},
			expectError: false,
			expectAPI:   true,
			checkAPI: func(t *testing.T, api *router.API) {
				assert.NotNil(t, api)
				// The key check: ensure the specific mapping params for this resolver are present.
				expectedParams := []apiConf.MappingParam{
					{Name: "requestBody.values", MapTo: "opt.values"},
					{Name: "requestBody.types", MapTo: "opt.types"},
					{Name: "uri.application", MapTo: "opt.application"},
					{Name: "uri.interface", MapTo: "opt.interface"},
					{Name: "uri.method", MapTo: "opt.method"},
				}
				assert.Equal(t, expectedParams, api.Method.IntegrationRequest.MappingParams)
				assert.Equal(t, apiConf.DubboRequest, api.Method.IntegrationRequest.RequestType)
			},
		},
		{
			name: "Failed Resolution - Invalid HTTP Method",
			setupCtx: func() *contexthttp.HttpContext {
				req := httptest.NewRequest("GET", "/my-app/my.service/myMethod", nil)
				req.Header.Set(constant.DubboHttpDubboVersion, "1.0.0")
				ctx := contexthttp.HttpContext{Ctx: context.Background()}
				ctx.Request = req
				return &ctx
			},
			expectError: true,
			expectAPI:   false,
			errorMsg:    "http request must be POST and have x-dubbo-http1.1-dubbo-version header",
		},
		{
			name: "Failed Resolution - Missing Header",
			setupCtx: func() *contexthttp.HttpContext {
				req := httptest.NewRequest("POST", "/my-app/my.service/myMethod", nil)
				ctx := contexthttp.HttpContext{Ctx: context.Background()}
				ctx.Request = req
				return &ctx
			},
			expectError: true,
			expectAPI:   false,
			errorMsg:    "http request must be POST and have x-dubbo-http1.1-dubbo-version header",
		},
		{
			name: "Failed Resolution - Invalid Path",
			setupCtx: func() *contexthttp.HttpContext {
				req := httptest.NewRequest("POST", "/my-app/my.service", nil)
				req.Header.Set(constant.DubboHttpDubboVersion, "1.0.0")
				ctx := contexthttp.HttpContext{Ctx: context.Background()}
				ctx.Request = req
				return &ctx
			},
			expectError: true,
			expectAPI:   false,
			errorMsg:    "http request path must be in {application}/{service}/{method} format",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			resolver := &StandardDubboResolver{}
			ctx := tt.setupCtx()
			api, err := resolver.Resolve(ctx)

			if tt.expectError {
				assert.Error(t, err)
				assert.Equal(t, tt.errorMsg, err.Error())
			} else {
				assert.NoError(t, err)
			}

			if tt.expectAPI {
				assert.NotNil(t, api)
				if tt.checkAPI != nil {
					tt.checkAPI(t, api)
				}
			} else {
				assert.Nil(t, api)
			}
		})
	}
}

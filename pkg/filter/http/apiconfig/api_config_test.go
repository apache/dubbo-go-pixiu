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

package apiconfig

import (
	"net/http"
	"net/http/httptest"
	"os"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	extfilter "github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/http/apiconfig/api"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/http/apiconfig/openapi"
	"github.com/apache/dubbo-go-pixiu/pkg/router"
)

func TestDecode_StopsOnOpenAPIValidationFailure(t *testing.T) {
	apiService := api.NewLocalMemoryAPIDiscoveryService()
	err := apiService.AddAPI(router.API{
		URLPattern: "/users",
		Method: config.Method{
			Enable:   true,
			HTTPVerb: constant.Get,
		},
		Metadata: map[string]any{
			openapi.ValidationPlanMetadataKey: &openapi.ValidationPlan{
				QueryParameters: []openapi.ParameterValidation{
					{Name: "id", Required: true},
				},
			},
		},
	})
	require.NoError(t, err)

	filterInstance := &Filter{apiService: apiService}
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	recorder := httptest.NewRecorder()
	ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

	status := filterInstance.Decode(ctx)

	assert.Equal(t, extfilter.Stop, status)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.True(t, ctx.LocalReply())
	assert.Nil(t, ctx.GetAPI())
}

func TestDecode_ContinuesOnOpenAPIValidationSuccess(t *testing.T) {
	apiService := api.NewLocalMemoryAPIDiscoveryService()
	err := apiService.AddAPI(router.API{
		URLPattern: "/users",
		Method: config.Method{
			Enable:   true,
			HTTPVerb: constant.Get,
		},
		Metadata: map[string]any{
			openapi.ValidationPlanMetadataKey: &openapi.ValidationPlan{
				QueryParameters: []openapi.ParameterValidation{
					{Name: "id", Required: true},
				},
			},
		},
	})
	require.NoError(t, err)

	filterInstance := &Filter{apiService: apiService}
	req := httptest.NewRequest(http.MethodGet, "/users?id=123", nil)
	recorder := httptest.NewRecorder()
	ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

	status := filterInstance.Decode(ctx)

	assert.Equal(t, extfilter.Continue, status)
	require.NotNil(t, ctx.GetAPI())
	assert.Equal(t, "/users", ctx.GetAPI().URLPattern)
}

func TestDecode_StopsOnOpenAPIParameterTypeFailure(t *testing.T) {
	apiService := api.NewLocalMemoryAPIDiscoveryService()
	err := apiService.AddAPI(router.API{
		URLPattern: "/users",
		Method: config.Method{
			Enable:   true,
			HTTPVerb: constant.Get,
		},
		Metadata: map[string]any{
			openapi.ValidationPlanMetadataKey: &openapi.ValidationPlan{
				QueryParameters: []openapi.ParameterValidation{
					{Name: "page", Type: "integer", Required: true},
				},
			},
		},
	})
	require.NoError(t, err)

	filterInstance := &Filter{apiService: apiService}
	req := httptest.NewRequest(http.MethodGet, "/users?page=abc", nil)
	recorder := httptest.NewRecorder()
	ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

	status := filterInstance.Decode(ctx)

	assert.Equal(t, extfilter.Stop, status)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.True(t, ctx.LocalReply())
	assert.Nil(t, ctx.GetAPI())
}

func TestApply_MergesOpenAPIRoutesFromFile(t *testing.T) {
	specFile, err := os.CreateTemp(t.TempDir(), "openapi-*.yaml")
	require.NoError(t, err)
	defer specFile.Close()

	_, err = specFile.WriteString(`
openapi: 3.0.3
info:
  title: users
  version: "1.0.0"
paths:
  /users/{id}:
    get:
      parameters:
        - name: id
          in: path
          required: true
          schema:
            type: string
      responses:
        "200":
          description: ok
  /users:
    post:
      parameters:
        - name: source
          in: query
          required: true
          schema:
            type: string
            enum: [web, app]
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required: [name]
              properties:
                name:
                  type: string
                  maxLength: 32
      responses:
        "200":
          description: ok
`)
	require.NoError(t, err)

	apiConfigFile, err := os.CreateTemp(t.TempDir(), "api-config-*.yaml")
	require.NoError(t, err)
	defer apiConfigFile.Close()

	_, err = apiConfigFile.WriteString(`
name: api name
resources:
  - path: /users
    type: restful
    methods:
      - httpVerb: POST
        enable: true
    resources:
      - path: /:id
        type: restful
        methods:
          - httpVerb: GET
            enable: true
`)
	require.NoError(t, err)

	factory := &FilterFactory{
		cfg: &ApiConfigConfig{
			Path:                    apiConfigFile.Name(),
			OpenAPIPath:             specFile.Name(),
			EnableOpenAPIValidation: true,
		},
	}

	err = factory.Apply()
	require.NoError(t, err)

	matched, err := factory.apiService.MatchAPI("/users", http.MethodPost)
	require.NoError(t, err)
	require.NotNil(t, openapi.ExtractValidationPlan(matched))
	assert.Equal(t, "/users", matched.URLPattern)
	assert.Equal(t, http.MethodPost, matched.HTTPVerb)

	pathMatched, err := factory.apiService.MatchAPI("/users/42", http.MethodGet)
	require.NoError(t, err)
	require.NotNil(t, openapi.ExtractValidationPlan(pathMatched))
	assert.Equal(t, "/users/:id", pathMatched.URLPattern)
}

func TestApply_RejectsDynamicOpenAPIValidationCombination(t *testing.T) {
	factory := &FilterFactory{
		cfg: &ApiConfigConfig{
			Dynamic:                 true,
			DynamicAdapter:          "mock",
			OpenAPIPath:             "configs/openapi_users.yaml",
			EnableOpenAPIValidation: true,
		},
	}

	err := factory.Apply()
	require.Error(t, err)
	assert.ErrorContains(t, err, "dynamic")
	assert.ErrorContains(t, err, "openapi")
}

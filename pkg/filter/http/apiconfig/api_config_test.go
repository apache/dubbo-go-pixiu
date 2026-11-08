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
	"path/filepath"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/yaml.v3"
)

import (
	extfilter "github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

func TestDecode_ContinuesWhenRouteMatches(t *testing.T) {
	factory := newAPIConfigFactory(t, `
name: api name
resources:
  - path: /users
    type: restful
    methods:
      - httpVerb: GET
        enable: true
`)

	filterInstance := &Filter{apiService: factory.apiService}
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	recorder := httptest.NewRecorder()
	ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

	status := filterInstance.Decode(ctx)

	assert.Equal(t, extfilter.Continue, status)
	require.NotNil(t, ctx.GetAPI())
	assert.Equal(t, "/users", ctx.GetAPI().URLPattern)
	assert.Equal(t, http.MethodGet, ctx.GetAPI().HTTPVerb)
	assert.False(t, ctx.LocalReply())
}

func TestDecode_StopsWhenRouteIsMissing(t *testing.T) {
	factory := newAPIConfigFactory(t, `
name: api name
resources:
  - path: /users
    type: restful
    methods:
      - httpVerb: GET
        enable: true
`)

	filterInstance := &Filter{apiService: factory.apiService}
	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	recorder := httptest.NewRecorder()
	ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

	status := filterInstance.Decode(ctx)

	assert.Equal(t, extfilter.Stop, status)
	assert.Equal(t, http.StatusNotFound, recorder.Code)
	assert.True(t, ctx.LocalReply())
	assert.Nil(t, ctx.GetAPI())
}

func TestDecode_StopsWhenRouteIsDisabled(t *testing.T) {
	factory := newAPIConfigFactory(t, `
name: api name
resources:
  - path: /users
    type: restful
    methods:
      - httpVerb: GET
        enable: false
`)

	filterInstance := &Filter{apiService: factory.apiService}
	req := httptest.NewRequest(http.MethodGet, "/users", nil)
	recorder := httptest.NewRecorder()
	ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

	status := filterInstance.Decode(ctx)

	assert.Equal(t, extfilter.Stop, status)
	assert.Equal(t, http.StatusNotAcceptable, recorder.Code)
	assert.True(t, ctx.LocalReply())
	assert.Nil(t, ctx.GetAPI())
}

func TestApply_KeepsNestedRouteMatching(t *testing.T) {
	factory := newAPIConfigFactory(t, `
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

	matched, err := factory.apiService.MatchAPI("/users", http.MethodPost)
	require.NoError(t, err)
	require.NotNil(t, matched)
	assert.Equal(t, "/users", matched.URLPattern)
	assert.Equal(t, http.MethodPost, matched.HTTPVerb)

	pathMatched, err := factory.apiService.MatchAPI("/users/42", http.MethodGet)
	require.NoError(t, err)
	require.NotNil(t, pathMatched)
	assert.Equal(t, "/users/:id", pathMatched.URLPattern)
}

func TestApply_RejectsDeprecatedOpenAPIConfig(t *testing.T) {
	factory := &FilterFactory{
		cfg: &ApiConfigConfig{
			OpenAPIPath:             "configs/openapi.yaml",
			EnableOpenAPIValidation: true,
		},
	}

	err := factory.Apply()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "dgp.filter.http.openapi")
}

func TestApply_RejectsDeprecatedOpenAPIConfigPresence(t *testing.T) {
	tests := []struct {
		name string
		yaml string
	}{
		{
			name: "empty openapi path",
			yaml: `
path: configs/api_config.yaml
openapi_path: ""
`,
		},
		{
			name: "disabled openapi validation",
			yaml: `
path: configs/api_config.yaml
enable_openapi_validation: false
`,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			cfg := &ApiConfigConfig{}
			require.NoError(t, yaml.Unmarshal([]byte(tt.yaml), cfg))
			factory := &FilterFactory{cfg: cfg}

			err := factory.Apply()

			require.Error(t, err)
			assert.Contains(t, err.Error(), "dgp.filter.http.openapi")
		})
	}
}

func newAPIConfigFactory(t *testing.T, apiConfig string) *FilterFactory {
	t.Helper()

	dir := t.TempDir()
	apiConfigPath := filepath.Join(dir, "api-config.yaml")
	require.NoError(t, os.WriteFile(apiConfigPath, []byte(apiConfig), 0o600))

	factory := &FilterFactory{
		cfg: &ApiConfigConfig{
			Path: apiConfigPath,
		},
	}
	require.NoError(t, factory.Apply())
	return factory
}

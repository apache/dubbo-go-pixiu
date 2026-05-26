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

package openapi

import (
	"io"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"strings"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	extfilter "github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

func TestDecode_StopsOnOpenAPIValidationFailure(t *testing.T) {
	filterInstance := newOpenAPIFilter(t, usersSpecWithQueryAndBody())
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"tom"}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

	status := filterInstance.Decode(ctx)

	assert.Equal(t, extfilter.Stop, status)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.True(t, ctx.LocalReply())
}

func TestDecode_ContinuesOnOpenAPIValidationSuccess(t *testing.T) {
	filterInstance := newOpenAPIFilter(t, usersSpecWithQueryAndBody())
	req := httptest.NewRequest(http.MethodPost, "/users?source=web", strings.NewReader(`{"name":"tom"}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

	status := filterInstance.Decode(ctx)

	assert.Equal(t, extfilter.Continue, status)
	body, err := io.ReadAll(req.Body)
	require.NoError(t, err)
	assert.JSONEq(t, `{"name":"tom"}`, string(body))
}

func TestDecode_IgnoresOpenAPISecurityValidation(t *testing.T) {
	filterInstance := newOpenAPIFilter(t, `
openapi: 3.0.3
info:
  title: users
  version: "1.0.0"
components:
  securitySchemes:
    ApiKeyAuth:
      type: apiKey
      in: header
      name: X-API-Key
security:
  - ApiKeyAuth: []
paths:
  /users:
    post:
      parameters:
        - name: source
          in: query
          required: true
          schema:
            type: string
      responses:
        "200":
          description: ok
`)
	req := httptest.NewRequest(http.MethodPost, "/users?source=web", nil)
	recorder := httptest.NewRecorder()
	ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

	status := filterInstance.Decode(ctx)

	assert.Equal(t, extfilter.Continue, status)
	assert.False(t, ctx.LocalReply())
}

func TestDecode_SkipsRoutesNotDeclaredInOpenAPI(t *testing.T) {
	filterInstance := newOpenAPIFilter(t, usersSpecWithQueryAndBody())
	req := httptest.NewRequest(http.MethodGet, "/orders", nil)
	recorder := httptest.NewRecorder()
	ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

	status := filterInstance.Decode(ctx)

	assert.Equal(t, extfilter.Continue, status)
	assert.False(t, ctx.LocalReply())
}

func TestDecode_ValidatesTemplatedPaths(t *testing.T) {
	filterInstance := newOpenAPIFilter(t, `
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
            type: integer
      responses:
        "200":
          description: ok
`)
	req := httptest.NewRequest(http.MethodGet, "/users/abc", nil)
	recorder := httptest.NewRecorder()
	ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

	status := filterInstance.Decode(ctx)

	assert.Equal(t, extfilter.Stop, status)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.True(t, ctx.LocalReply())
}

func TestDecode_ValidatesHeaderRequiredAndType(t *testing.T) {
	filterInstance := newOpenAPIFilter(t, `
openapi: 3.0.3
info:
  title: reports
  version: "1.0.0"
paths:
  /reports:
    get:
      parameters:
        - name: X-Tenant-ID
          in: header
          required: true
          schema:
            type: integer
      responses:
        "200":
          description: ok
`)

	tests := []struct {
		name        string
		headerValue string
		wantStatus  extfilter.FilterStatus
	}{
		{
			name:       "missing required header",
			wantStatus: extfilter.Stop,
		},
		{
			name:        "invalid header type",
			headerValue: "tenant-a",
			wantStatus:  extfilter.Stop,
		},
		{
			name:        "valid header",
			headerValue: "42",
			wantStatus:  extfilter.Continue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodGet, "/reports", nil)
			if tt.headerValue != "" {
				req.Header.Set("X-Tenant-ID", tt.headerValue)
			}
			recorder := httptest.NewRecorder()
			ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

			status := filterInstance.Decode(ctx)

			assert.Equal(t, tt.wantStatus, status)
			assert.Equal(t, tt.wantStatus == extfilter.Stop, ctx.LocalReply())
		})
	}
}

func TestDecode_ValidatesQueryEnum(t *testing.T) {
	filterInstance := newOpenAPIFilter(t, usersSpecWithQueryAndBody())
	req := httptest.NewRequest(http.MethodPost, "/users?source=cli", strings.NewReader(`{"name":"tom"}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

	status := filterInstance.Decode(ctx)

	assert.Equal(t, extfilter.Stop, status)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.True(t, ctx.LocalReply())
}

func TestDecode_ValidatesRequestBodyMinAndMaxLength(t *testing.T) {
	filterInstance := newOpenAPIFilter(t, usersSpecWithQueryAndBody())

	tests := []struct {
		name       string
		body       string
		wantStatus extfilter.FilterStatus
	}{
		{
			name:       "below minLength",
			body:       `{"name":"Al"}`,
			wantStatus: extfilter.Stop,
		},
		{
			name:       "above maxLength",
			body:       `{"name":"Alexandria"}`,
			wantStatus: extfilter.Stop,
		},
		{
			name:       "within length bounds",
			body:       `{"name":"Alice"}`,
			wantStatus: extfilter.Continue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/users?source=web", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

			status := filterInstance.Decode(ctx)

			assert.Equal(t, tt.wantStatus, status)
			assert.Equal(t, tt.wantStatus == extfilter.Stop, ctx.LocalReply())
		})
	}
}

func TestDecode_ValidatesRequestBodyMinimumAndMaximum(t *testing.T) {
	filterInstance := newOpenAPIFilter(t, usersSpecWithQueryAndBody())

	tests := []struct {
		name       string
		body       string
		wantStatus extfilter.FilterStatus
	}{
		{
			name:       "below minimum",
			body:       `{"name":"Alice","age":12}`,
			wantStatus: extfilter.Stop,
		},
		{
			name:       "above maximum",
			body:       `{"name":"Alice","age":121}`,
			wantStatus: extfilter.Stop,
		},
		{
			name:       "within numeric bounds",
			body:       `{"name":"Alice","age":18}`,
			wantStatus: extfilter.Continue,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			req := httptest.NewRequest(http.MethodPost, "/users?source=web", strings.NewReader(tt.body))
			req.Header.Set("Content-Type", "application/json")
			recorder := httptest.NewRecorder()
			ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

			status := filterInstance.Decode(ctx)

			assert.Equal(t, tt.wantStatus, status)
			assert.Equal(t, tt.wantStatus == extfilter.Stop, ctx.LocalReply())
		})
	}
}

func newOpenAPIFilter(t *testing.T, spec string) *Filter {
	t.Helper()

	dir := t.TempDir()
	specPath := filepath.Join(dir, "openapi.yaml")
	require.NoError(t, os.WriteFile(specPath, []byte(spec), 0o600))

	factory := &FilterFactory{
		cfg: &Config{
			Path: specPath,
		},
	}
	require.NoError(t, factory.Apply())

	return &Filter{
		validator: factory.validator,
		model:     factory.model,
	}
}

func usersSpecWithQueryAndBody() string {
	return `
openapi: 3.0.3
info:
  title: users
  version: "1.0.0"
paths:
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
                  minLength: 3
                  maxLength: 8
                age:
                  type: integer
                  minimum: 13
                  maximum: 120
      responses:
        "200":
          description: ok
`
}

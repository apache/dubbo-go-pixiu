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
	v3 "github.com/pb33f/libopenapi/datamodel/high/v3"

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

func TestDecode_HidesOpenAPIValidationDetailsFromClient(t *testing.T) {
	filterInstance := newOpenAPIFilter(t, usersSpecWithQueryAndBody())
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"tom"}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

	status := filterInstance.Decode(ctx)

	assert.Equal(t, extfilter.Stop, status)
	assert.Equal(t, http.StatusBadRequest, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "openapi request validation failed")
	assert.NotContains(t, recorder.Body.String(), "source")
}

func TestDecode_RejectsOpenAPIRequestBodyLargerThanLimit(t *testing.T) {
	filterInstance := newOpenAPIFilterWithConfig(t, usersSpecWithQueryAndBody(), func(cfg *Config) {
		cfg.MaxRequestBodyBytes = 8
	})
	req := httptest.NewRequest(http.MethodPost, "/users?source=web", strings.NewReader(`{"name":"tom"}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

	status := filterInstance.Decode(ctx)

	assert.Equal(t, extfilter.Stop, status)
	assert.Equal(t, http.StatusRequestEntityTooLarge, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "openapi request body too large")
}

func TestDecode_RejectsChunkedOpenAPIRequestBodyLargerThanLimit(t *testing.T) {
	filterInstance := newOpenAPIFilterWithConfig(t, usersSpecWithQueryAndBody(), func(cfg *Config) {
		cfg.MaxRequestBodyBytes = 8
	})
	req := httptest.NewRequest(http.MethodPost, "/users?source=web", io.NopCloser(strings.NewReader(`{"name":"tom"}`)))
	req.ContentLength = -1
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

	status := filterInstance.Decode(ctx)

	assert.Equal(t, extfilter.Stop, status)
	assert.Equal(t, http.StatusRequestEntityTooLarge, recorder.Code)
	assert.Contains(t, recorder.Body.String(), "openapi request body too large")
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

func TestDecode_SkipsMethodsNotDeclaredInOpenAPI(t *testing.T) {
	filterInstance := newOpenAPIFilter(t, `
openapi: 3.0.3
info:
  title: users
  version: "1.0.0"
paths:
  /users:
    get:
      responses:
        "200":
          description: ok
`)
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"tom"}`))
	req.Header.Set("Content-Type", "application/json")
	recorder := httptest.NewRecorder()
	ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

	status := filterInstance.Decode(ctx)

	assert.Equal(t, extfilter.Continue, status)
	assert.False(t, ctx.LocalReply())
}

func TestDecode_SkipsHeadWhenOnlyGetDeclared(t *testing.T) {
	// When only GET is declared (no HEAD), the SDK FindPath treats HEAD as
	// an undeclared method and returns nil. The filter skips validation and
	// lets the request continue to downstream filters / handlers.
	filterInstance := newOpenAPIFilter(t, `
openapi: 3.0.3
info:
  title: users
  version: "1.0.0"
paths:
  /users:
    get:
      parameters:
        - name: source
          in: query
          required: true
          schema:
            type: string
            enum: [web, app]
      responses:
        "200":
          description: ok
`)

	t.Run("HEAD without required query param is skipped", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodHead, "/users", nil)
		recorder := httptest.NewRecorder()
		ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

		status := filterInstance.Decode(ctx)

		assert.Equal(t, extfilter.Continue, status)
		assert.False(t, ctx.LocalReply())
	})

	t.Run("HEAD with valid query param is also skipped", func(t *testing.T) {
		req := httptest.NewRequest(http.MethodHead, "/users?source=web", nil)
		recorder := httptest.NewRecorder()
		ctx := &contexthttp.HttpContext{Request: req, Writer: recorder}

		status := filterInstance.Decode(ctx)

		assert.Equal(t, extfilter.Continue, status)
		assert.False(t, ctx.LocalReply())
	})
}

func TestHasRequestOperation_HeadRequiresExplicitHeadOperation(t *testing.T) {
	req := httptest.NewRequest(http.MethodHead, "/users", nil)

	assert.False(t, hasRequestOperation(req, &v3.PathItem{
		Get: &v3.Operation{},
	}))
	assert.True(t, hasRequestOperation(req, &v3.PathItem{
		Head: &v3.Operation{},
	}))
}

func TestApply_RejectsInvalidOpenAPIModel(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	require.NoError(t, os.WriteFile("openapi.yaml", []byte(`
openapi: 3.0.3
info:
  title: users
  version: "1.0.0"
paths:
  /users:
    post:
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: "#/components/schemas/User"
      responses:
        "200":
          description: ok
components:
  schemas:
    User:
      $ref: "#/components/schemas/User2"
`), 0o600))

	factory := &FilterFactory{
		cfg: &Config{
			Path: "openapi.yaml",
		},
	}

	err := factory.Apply()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "build openapi model")
}

func TestApply_RejectsAbsoluteOpenAPIPath(t *testing.T) {
	factory := &FilterFactory{
		cfg: &Config{
			Path: "/etc/passwd",
		},
	}

	err := factory.Apply()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "openapi path must be relative")
}

func TestApply_RejectsParentDirectoryOpenAPIPath(t *testing.T) {
	factory := &FilterFactory{
		cfg: &Config{
			Path: "../openapi.yaml",
		},
	}

	err := factory.Apply()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "openapi path must not contain parent directory")
}

func TestApply_RejectsSymlinkPointingToSensitiveDirectory(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)

	// Create a symlink: specs/current -> /etc
	require.NoError(t, os.MkdirAll("specs", 0o700))
	require.NoError(t, os.Symlink("/etc", filepath.Join(dir, "specs", "current")))

	factory := &FilterFactory{
		cfg: &Config{
			Path: "specs/current/passwd",
		},
	}

	err := factory.Apply()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "openapi path base directory is not allowed")
}

func TestApply_AllowsNestedLocalExternalRefs(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	require.NoError(t, os.MkdirAll("specs/schemas", 0o700))
	require.NoError(t, os.WriteFile("specs/openapi.yaml", []byte(`
openapi: 3.0.3
info:
  title: users
  version: "1.0.0"
paths:
  /users:
    post:
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: ./schemas/user.yaml#/User
      responses:
        "200":
          description: ok
`), 0o600))
	require.NoError(t, os.WriteFile("specs/schemas/user.yaml", []byte(`
User:
  type: object
  required: [name]
  properties:
    name:
      $ref: ./name.yaml#/Name
`), 0o600))
	require.NoError(t, os.WriteFile("specs/schemas/name.yaml", []byte(`
Name:
  type: string
  minLength: 3
`), 0o600))

	factory := &FilterFactory{
		cfg: &Config{
			Path: "specs/openapi.yaml",
		},
	}

	require.NoError(t, factory.Apply())
}

func TestApply_RejectsExternalRefWithParentDirectory(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	require.NoError(t, os.MkdirAll("specs", 0o700))
	require.NoError(t, os.WriteFile("specs/openapi.yaml", []byte(`
openapi: 3.0.3
info:
  title: users
  version: "1.0.0"
paths:
  /users:
    post:
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: ../schemas/user.yaml#/User
      responses:
        "200":
          description: ok
`), 0o600))

	factory := &FilterFactory{
		cfg: &Config{
			Path: "specs/openapi.yaml",
		},
	}

	err := factory.Apply()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "openapi external ref must not contain parent directory")
}

func TestApply_RejectsAbsoluteExternalRef(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	require.NoError(t, os.WriteFile("openapi.yaml", []byte(`
openapi: 3.0.3
info:
  title: users
  version: "1.0.0"
paths:
  /users:
    post:
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: /etc/passwd#/User
      responses:
        "200":
          description: ok
`), 0o600))

	factory := &FilterFactory{
		cfg: &Config{
			Path: "openapi.yaml",
		},
	}

	err := factory.Apply()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "openapi external ref must be relative")
}

func TestApply_RejectsRemoteExternalRef(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	require.NoError(t, os.WriteFile("openapi.yaml", []byte(`
openapi: 3.0.3
info:
  title: users
  version: "1.0.0"
paths:
  /users:
    post:
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: https://example.com/schemas/user.yaml#/User
      responses:
        "200":
          description: ok
`), 0o600))

	factory := &FilterFactory{
		cfg: &Config{
			Path: "openapi.yaml",
		},
	}

	err := factory.Apply()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "openapi external ref must be relative")
}

func TestApply_RejectsNestedExternalRefEscapingAllowedDirectoryViaSymlink(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	require.NoError(t, os.MkdirAll("specs/schemas", 0o700))
	require.NoError(t, os.Symlink("/etc", filepath.Join(dir, "specs", "schemas", "current")))
	require.NoError(t, os.WriteFile("specs/openapi.yaml", []byte(`
openapi: 3.0.3
info:
  title: users
  version: "1.0.0"
paths:
  /users:
    post:
      requestBody:
        required: true
        content:
          application/json:
            schema:
              $ref: ./schemas/user.yaml#/User
      responses:
        "200":
          description: ok
`), 0o600))
	require.NoError(t, os.WriteFile("specs/schemas/user.yaml", []byte(`
User:
  type: object
  properties:
    name:
      $ref: ./current/passwd#/Name
`), 0o600))

	factory := &FilterFactory{
		cfg: &Config{
			Path: "specs/openapi.yaml",
		},
	}

	err := factory.Apply()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "openapi external ref escapes allowed directory")
}

func TestApply_RejectsNegativeMaxRequestBodyBytes(t *testing.T) {
	dir := t.TempDir()
	t.Chdir(dir)
	require.NoError(t, os.WriteFile("openapi.yaml", []byte(usersSpecWithQueryAndBody()), 0o600))

	factory := &FilterFactory{
		cfg: &Config{
			Path:                "openapi.yaml",
			MaxRequestBodyBytes: -1,
		},
	}

	err := factory.Apply()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "max_request_body_bytes must not be negative")
}

func TestApply_UsesDefaultMaxRequestBodyBytesForZeroValue(t *testing.T) {
	factory := newOpenAPIFactory(t, usersSpecWithQueryAndBody(), func(cfg *Config) {
		cfg.MaxRequestBodyBytes = 0
	})

	assert.Equal(t, int64(defaultMaxRequestBodyBytes), factory.maxRequestBody)
}

func TestApply_BuildsReusablePathLookupOptions(t *testing.T) {
	factory := newOpenAPIFactory(t, usersSpecWithQueryAndBody(), nil)

	require.NotNil(t, factory.validationOptions)
	require.NotNil(t, factory.validationOptions.PathTree)

	filterInstance := &Filter{
		validator:         factory.validator,
		model:             factory.model,
		validationOptions: factory.validationOptions,
	}
	req := httptest.NewRequest(http.MethodPost, "/users?source=web", strings.NewReader(`{"name":"tom"}`))

	pathItem, foundPath, ok := filterInstance.findRequestOperation(req)

	require.True(t, ok)
	require.NotNil(t, pathItem)
	assert.Equal(t, "/users", foundPath)
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
	factory := newOpenAPIFactory(t, spec, nil)
	return &Filter{
		validator:         factory.validator,
		model:             factory.model,
		validationOptions: factory.validationOptions,
		maxRequestBody:    factory.maxRequestBody,
	}
}

func newOpenAPIFilterWithConfig(t *testing.T, spec string, configure func(*Config)) *Filter {
	t.Helper()
	factory := newOpenAPIFactory(t, spec, configure)
	return &Filter{
		validator:         factory.validator,
		model:             factory.model,
		validationOptions: factory.validationOptions,
		maxRequestBody:    factory.maxRequestBody,
	}
}

func newOpenAPIFactory(t *testing.T, spec string, configure func(*Config)) *FilterFactory {
	t.Helper()

	dir := t.TempDir()
	t.Chdir(dir)
	specPath := "openapi.yaml"
	require.NoError(t, os.WriteFile(filepath.Join(dir, specPath), []byte(spec), 0o600))

	cfg := &Config{
		Path: specPath,
	}
	if configure != nil {
		configure(cfg)
	}

	factory := &FilterFactory{
		cfg: cfg,
	}
	require.NoError(t, factory.Apply())

	return factory
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

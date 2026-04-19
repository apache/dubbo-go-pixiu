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
	"net/http"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

const minimalOpenAPISpec = `
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
        - name: trace-id
          in: header
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
            enum:
              - web
              - app
        - name: page
          in: query
          required: true
          schema:
            type: integer
            minimum: 1
            maximum: 100
        - name: x-trace
          in: header
          required: true
          schema:
            type: string
            minLength: 4
            maxLength: 12
      requestBody:
        required: true
        content:
          application/json:
            schema:
              type: object
              required:
                - name
                - role
              properties:
                name:
                  type: string
                  maxLength: 32
                role:
                  type: string
                  enum:
                    - admin
                    - member
      responses:
        "200":
          description: ok
`

const pathLevelParametersSpec = `
openapi: 3.0.3
info:
  title: users
  version: "1.0.0"
paths:
  /users/{id}:
    parameters:
      - name: id
        in: path
        required: true
        schema:
          type: string
      - name: trace-id
        in: header
        required: true
        schema:
          type: string
    get:
      responses:
        "200":
          description: ok
`

func TestCompileOpenAPIRoute(t *testing.T) {
	doc, err := LoadDocumentFromBytes([]byte(minimalOpenAPISpec))
	require.NoError(t, err)

	compiled, err := CompileRoutes(doc)
	require.NoError(t, err)
	require.Len(t, compiled, 2)

	compiledByMethod := map[string]CompiledRoute{}
	for _, item := range compiled {
		compiledByMethod[item.Route.HTTPVerb+" "+item.Route.URLPattern] = item
	}

	getRoute, ok := compiledByMethod[http.MethodGet+" /users/:id"]
	require.True(t, ok)
	assert.Equal(t, "/users/:id", getRoute.Route.URLPattern)
	require.NotNil(t, getRoute.Validation)
	assert.Equal(t, "/users/:id", getRoute.Validation.RoutePattern)
	require.Len(t, getRoute.Validation.PathParameters, 1)
	assert.Equal(t, "id", getRoute.Validation.PathParameters[0].Name)
	require.Len(t, getRoute.Validation.HeaderParameters, 1)
	assert.Equal(t, "trace-id", getRoute.Validation.HeaderParameters[0].Name)
	assert.True(t, getRoute.Validation.HeaderParameters[0].Required)

	postRoute, ok := compiledByMethod[http.MethodPost+" /users"]
	require.True(t, ok)
	assert.Equal(t, "/users", postRoute.Route.URLPattern)
	assert.NotNil(t, postRoute.Validation)
	require.Len(t, postRoute.Validation.QueryParameters, 2)
	assert.Equal(t, "source", postRoute.Validation.QueryParameters[0].Name)
	assert.True(t, postRoute.Validation.QueryParameters[0].Required)
	assert.Equal(t, []string{"web", "app"}, postRoute.Validation.QueryParameters[0].Enum)
	pageParam := postRoute.Validation.QueryParameters[1]
	assert.Equal(t, "page", pageParam.Name)
	assert.Equal(t, "integer", pageParam.Type)
	require.NotNil(t, pageParam.Minimum)
	assert.Equal(t, 1.0, *pageParam.Minimum)
	require.NotNil(t, pageParam.Maximum)
	assert.Equal(t, 100.0, *pageParam.Maximum)

	require.Len(t, postRoute.Validation.HeaderParameters, 1)
	headerParam := postRoute.Validation.HeaderParameters[0]
	assert.Equal(t, "x-trace", headerParam.Name)
	assert.Equal(t, "string", headerParam.Type)
	require.NotNil(t, headerParam.MinLength)
	assert.Equal(t, 4, *headerParam.MinLength)
	require.NotNil(t, headerParam.MaxLength)
	assert.Equal(t, 12, *headerParam.MaxLength)

	require.NotNil(t, postRoute.Validation.RequestBody)
	assert.True(t, postRoute.Validation.RequestBody.Required)
	require.NotNil(t, postRoute.Validation.RequestBody.Schema)
	assert.Equal(t, "object", postRoute.Validation.RequestBody.Schema.Type)
	assert.Equal(t, []string{"name", "role"}, postRoute.Validation.RequestBody.Schema.Required)
	require.Contains(t, postRoute.Validation.RequestBody.Schema.Properties, "name")
	require.Contains(t, postRoute.Validation.RequestBody.Schema.Properties, "role")
	require.NotNil(t, postRoute.Validation.RequestBody.Schema.Properties["name"].MaxLength)
	assert.Equal(t, 32, *postRoute.Validation.RequestBody.Schema.Properties["name"].MaxLength)
	assert.Equal(t, []string{"admin", "member"}, postRoute.Validation.RequestBody.Schema.Properties["role"].Enum)
}

func TestCompileOpenAPIRoute_UsesPathLevelParameters(t *testing.T) {
	doc, err := LoadDocumentFromBytes([]byte(pathLevelParametersSpec))
	require.NoError(t, err)

	compiled, err := CompileRoutes(doc)
	require.NoError(t, err)
	require.Len(t, compiled, 1)

	getRoute := compiled[0]
	assert.Equal(t, http.MethodGet, getRoute.Route.HTTPVerb)
	assert.Equal(t, "/users/:id", getRoute.Route.URLPattern)
	require.NotNil(t, getRoute.Validation)
	require.Len(t, getRoute.Validation.PathParameters, 1)
	assert.Equal(t, "id", getRoute.Validation.PathParameters[0].Name)
	assert.True(t, getRoute.Validation.PathParameters[0].Required)
	require.Len(t, getRoute.Validation.HeaderParameters, 1)
	assert.Equal(t, "trace-id", getRoute.Validation.HeaderParameters[0].Name)
	assert.True(t, getRoute.Validation.HeaderParameters[0].Required)
}

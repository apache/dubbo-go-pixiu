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
	"net/http/httptest"
	"strings"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestValidateRequest_MissingRequiredQuery(t *testing.T) {
	req := httptest.NewRequest("GET", "/users", nil)
	plan := &ValidationPlan{
		RoutePattern: "/users",
		QueryParameters: []ParameterValidation{
			{Name: "id", Required: true},
		},
	}

	err := ValidateRequest(req, plan)
	require.Error(t, err)
	assert.ErrorContains(t, err, "id")
}

func TestValidateRequest_QueryEnumMismatch(t *testing.T) {
	req := httptest.NewRequest("GET", "/users?role=superadmin", nil)
	plan := &ValidationPlan{
		RoutePattern: "/users",
		QueryParameters: []ParameterValidation{
			{Name: "role", Enum: []string{"admin", "member"}},
		},
	}

	err := ValidateRequest(req, plan)
	require.Error(t, err)
	assert.ErrorContains(t, err, "role")
}

func TestValidateRequest_QueryTypeMismatch(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users?page=abc", nil)
	plan := &ValidationPlan{
		RoutePattern: "/users",
		QueryParameters: []ParameterValidation{
			{Name: "page", Type: "integer", Required: true},
		},
	}

	err := ValidateRequest(req, plan)
	require.Error(t, err)
	assert.ErrorContains(t, err, "integer")
}

func TestValidateRequest_QueryRangeViolation(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users?page=0", nil)
	plan := &ValidationPlan{
		RoutePattern: "/users",
		QueryParameters: []ParameterValidation{
			{Name: "page", Type: "integer", Minimum: float64Ptr(1), Maximum: float64Ptr(100)},
		},
	}

	err := ValidateRequest(req, plan)
	require.Error(t, err)
	assert.ErrorContains(t, err, ">=")
}

func TestValidateRequest_InvalidJSONBody(t *testing.T) {
	req := httptest.NewRequest("POST", "/users", strings.NewReader(`{"role":"member"}`))
	req.Header.Set("Content-Type", "application/json")

	plan := &ValidationPlan{
		RoutePattern: "/users",
		RequestBody: &BodyValidation{
			Required: true,
			Schema: &SchemaValidation{
				Type:     "object",
				Required: []string{"name", "role"},
				Properties: map[string]*SchemaValidation{
					"name": {
						Type:      "string",
						MaxLength: intPtr(32),
					},
					"role": {
						Type: "string",
						Enum: []string{"admin", "member"},
					},
				},
			},
		},
	}

	err := ValidateRequest(req, plan)
	require.Error(t, err)
	assert.ErrorContains(t, err, "name")
}

func TestValidateRequest_MissingRequiredHeader(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	plan := &ValidationPlan{
		RoutePattern: "/users/:id",
		PathParameters: []ParameterValidation{
			{Name: "id", Required: true},
		},
		HeaderParameters: []ParameterValidation{
			{Name: "trace-id", Required: true},
		},
	}

	err := ValidateRequest(req, plan)
	require.Error(t, err)
	assert.ErrorContains(t, err, "trace-id")
}

func TestValidateRequest_HeaderLengthViolation(t *testing.T) {
	req := httptest.NewRequest(http.MethodGet, "/users/123", nil)
	req.Header.Set("trace-id", "abc")
	plan := &ValidationPlan{
		RoutePattern: "/users/:id",
		HeaderParameters: []ParameterValidation{
			{Name: "trace-id", Type: "string", MinLength: intPtr(4)},
		},
	}

	err := ValidateRequest(req, plan)
	require.Error(t, err)
	assert.ErrorContains(t, err, "length")
}

func TestValidateRequest_InvalidNumericBody(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/products", strings.NewReader(`{"price":-10}`))
	req.Header.Set("Content-Type", "application/json")

	plan := &ValidationPlan{
		RoutePattern: "/products",
		RequestBody: &BodyValidation{
			Required: true,
			Schema: &SchemaValidation{
				Type:     "object",
				Required: []string{"price"},
				Properties: map[string]*SchemaValidation{
					"price": {
						Type:    "number",
						Minimum: float64Ptr(0),
					},
				},
			},
		},
	}

	err := ValidateRequest(req, plan)
	require.Error(t, err)
	assert.ErrorContains(t, err, "price")
}

func TestValidateRequest_RejectsUnsupportedBodyContentType(t *testing.T) {
	req := httptest.NewRequest(http.MethodPost, "/users", strings.NewReader(`{"name":"tom"}`))
	req.Header.Set("Content-Type", "text/plain")

	plan := &ValidationPlan{
		RoutePattern: "/users",
		RequestBody: &BodyValidation{
			Required: true,
			Schema: &SchemaValidation{
				Type: "object",
			},
		},
	}

	err := ValidateRequest(req, plan)
	require.Error(t, err)
	assert.ErrorContains(t, err, "application/json")
}

func float64Ptr(v float64) *float64 {
	return &v
}

func intPtr(v int) *int {
	return &v
}

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

package opa

import (
	"context"
	"encoding/json"
	"net/http"
	"net/http/httptest"
	"testing"
	"time"
)

import (
	"github.com/open-policy-agent/opa/rego"

	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contextHttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

const testPolicy = `
package test
import future.keywords.if

default allow := false

allow if {
    input.headers[Test_Header][0] == "1"
}
`

// setupFilterWithoutFile is a helper function for testing. It simulates the core logic of
// PrepareFilterChain by creating an OPA Rego instance and preparing a query directly.
func setupFilterWithoutFile(t *testing.T, policy string) *Filter {
	r := rego.New(
		rego.Query("data.test.allow"),
		rego.Module("policy.rego", policy),
	)

	preparedQuery, err := r.PrepareForEval(context.Background())
	assert.Nil(t, err)

	return &Filter{
		cfg: &Config{
			Policy:     policy,
			Entrypoint: "data.test.allow",
		},
		preparedQuery: &preparedQuery,
	}
}

// TestEmbeddedAllowedRule tests embedded mode with allowed request
func TestEmbeddedAllowedRule(t *testing.T) {
	f := setupFilterWithoutFile(t, testPolicy)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Test_Header", "1")

	rec := httptest.NewRecorder()
	ctx := &contextHttp.HttpContext{
		Writer:  rec,
		Request: req,
		Ctx:     context.Background(),
	}

	result := f.Decode(ctx)
	assert.Equal(t, filter.Continue, result)
}

// TestEmbeddedDeniedRule tests embedded mode with denied request
func TestEmbeddedDeniedRule(t *testing.T) {
	f := setupFilterWithoutFile(t, testPolicy)
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Test_Header", "0")

	rec := httptest.NewRecorder()
	ctx := &contextHttp.HttpContext{
		Writer:  rec,
		Request: req,
		Ctx:     context.Background(),
	}

	result := f.Decode(ctx)
	assert.Equal(t, filter.Stop, result)
}

// TestServerModeAllowed tests server mode with allowed request
func TestServerModeAllowed(t *testing.T) {
	// Create mock OPA server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		assert.Equal(t, "POST", r.Method)
		assert.Equal(t, "application/json", r.Header.Get("Content-Type"))
		assert.Equal(t, "/v1/data/test/allow", r.URL.Path)

		// Read and verify request body
		var reqBody map[string]any
		err := json.NewDecoder(r.Body).Decode(&reqBody)
		assert.Nil(t, err)
		assert.NotNil(t, reqBody["input"])

		input := reqBody["input"].(map[string]any)

		// Simulate policy: Test_Header == "1" allows
		// After JSON unmarshaling, headers is map[string][]string
		// Note: HTTP header keys are canonicalized (e.g., "Test_Header" -> "Test_header")
		allow := false
		if headersMap, ok := input["headers"].(map[string]any); ok {
			// Check for "Test_header" (canonicalized form)
			if testHeader, ok := headersMap["Test_header"]; ok {
				if headerArray, ok := testHeader.([]any); ok && len(headerArray) > 0 {
					if strVal, ok := headerArray[0].(string); ok {
						allow = strVal == "1"
					}
				}
			}
		}

		// Return OPA response format
		response := map[string]any{
			"result": allow,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	// Create FilterFactory with server mode configuration
	factory := &FilterFactory{
		cfg: &Config{
			ServerURL:    server.URL,
			DecisionPath: "/v1/data/test/allow",
			TimeoutMs:    100,
		},
	}

	err := factory.Apply()
	assert.Nil(t, err)
	assert.NotNil(t, factory.httpClient)

	// Prepare filter chain
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Test_Header", "1")
	rec := httptest.NewRecorder()
	ctx := &contextHttp.HttpContext{
		Writer:  rec,
		Request: req,
		Ctx:     context.Background(),
	}

	chain := &mockFilterChain{}
	err = factory.PrepareFilterChain(ctx, chain)
	assert.Nil(t, err)
	assert.Len(t, chain.filters, 1)

	// Execute filter
	result := chain.filters[0].Decode(ctx)
	assert.Equal(t, filter.Continue, result)
}

// TestServerModeDenied tests server mode with denied request
func TestServerModeDenied(t *testing.T) {
	// Create mock OPA server
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody map[string]any
		json.NewDecoder(r.Body).Decode(&reqBody)
		input := reqBody["input"].(map[string]any)

		// After JSON unmarshaling, headers is map[string][]string
		// Note: HTTP header keys are canonicalized
		allow := false
		if headersMap, ok := input["headers"].(map[string]any); ok {
			if testHeader, ok := headersMap["Test_header"]; ok {
				if headerArray, ok := testHeader.([]any); ok && len(headerArray) > 0 {
					if strVal, ok := headerArray[0].(string); ok {
						allow = strVal == "1"
					}
				}
			}
		}

		response := map[string]any{
			"result": allow,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	factory := &FilterFactory{
		cfg: &Config{
			ServerURL:    server.URL,
			DecisionPath: "/v1/data/test/allow",
			TimeoutMs:    100,
		},
	}

	err := factory.Apply()
	assert.Nil(t, err)

	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Test_Header", "0") // This will be denied
	rec := httptest.NewRecorder()
	ctx := &contextHttp.HttpContext{
		Writer:  rec,
		Request: req,
		Ctx:     context.Background(),
	}

	chain := &mockFilterChain{}
	err = factory.PrepareFilterChain(ctx, chain)
	assert.Nil(t, err)

	result := chain.filters[0].Decode(ctx)
	assert.Equal(t, filter.Stop, result)
}

// TestServerModeWithBearerToken tests server mode with authentication token
func TestServerModeWithBearerToken(t *testing.T) {
	expectedToken := "test-token-123"

	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Verify authentication header
		authHeader := r.Header.Get("Authorization")
		assert.Equal(t, "Bearer "+expectedToken, authHeader)

		response := map[string]any{
			"result": true,
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	factory := &FilterFactory{
		cfg: &Config{
			ServerURL:    server.URL,
			DecisionPath: "/v1/data/test/allow",
			TimeoutMs:    100,
			BearerToken:  expectedToken,
		},
	}

	err := factory.Apply()
	assert.Nil(t, err)

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	ctx := &contextHttp.HttpContext{
		Writer:  rec,
		Request: req,
		Ctx:     context.Background(),
	}

	chain := &mockFilterChain{}
	factory.PrepareFilterChain(ctx, chain)
	result := chain.filters[0].Decode(ctx)
	assert.Equal(t, filter.Continue, result)
}

// TestServerModeError tests error handling in server mode
func TestServerModeError(t *testing.T) {
	// Create a server that returns error
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusInternalServerError)
		w.Write([]byte("Internal server error"))
	}))
	defer server.Close()

	factory := &FilterFactory{
		cfg: &Config{
			ServerURL:    server.URL,
			DecisionPath: "/v1/data/test/allow",
			TimeoutMs:    100,
		},
	}

	err := factory.Apply()
	assert.Nil(t, err)

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	ctx := &contextHttp.HttpContext{
		Writer:  rec,
		Request: req,
		Ctx:     context.Background(),
	}

	chain := &mockFilterChain{}
	factory.PrepareFilterChain(ctx, chain)
	result := chain.filters[0].Decode(ctx)
	assert.Equal(t, filter.Stop, result) // Should deny on error
}

// TestServerModeTimeout tests timeout handling in server mode
func TestServerModeTimeout(t *testing.T) {
	// Create a server that delays response
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// Sleep longer than the timeout
		time.Sleep(200 * time.Millisecond)
		response := map[string]any{
			"result": true,
		}
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	factory := &FilterFactory{
		cfg: &Config{
			ServerURL:    server.URL,
			DecisionPath: "/v1/data/test/allow",
			TimeoutMs:    50, // Very short timeout to trigger timeout
		},
	}

	err := factory.Apply()
	assert.Nil(t, err)

	req := httptest.NewRequest("GET", "/test", nil)
	rec := httptest.NewRecorder()
	ctx := &contextHttp.HttpContext{
		Writer:  rec,
		Request: req,
		Ctx:     context.Background(),
	}

	chain := &mockFilterChain{}
	factory.PrepareFilterChain(ctx, chain)
	result := chain.filters[0].Decode(ctx)

	// Should return Stop on timeout
	assert.Equal(t, filter.Stop, result)

	// Check that the response contains timeout error (504)
	assert.Equal(t, 504, ctx.GetStatusCode())
}

// TestServerModeObjectResponse tests server mode with object response format
func TestServerModeObjectResponse(t *testing.T) {
	// Test object response with "allow" field
	server := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		var reqBody map[string]any
		json.NewDecoder(r.Body).Decode(&reqBody)
		input := reqBody["input"].(map[string]any)

		headersMap := input["headers"].(map[string]any)
		testHeaderValue := ""
		if testHeader, ok := headersMap["Test_header"]; ok {
			if headerArray, ok := testHeader.([]any); ok && len(headerArray) > 0 {
				if strVal, ok := headerArray[0].(string); ok {
					testHeaderValue = strVal
				}
			}
		}

		// Return object format: {allow: true}
		response := map[string]any{
			"result": map[string]any{
				"allow":  testHeaderValue == "1",
				"reason": "test policy",
			},
		}
		w.Header().Set("Content-Type", "application/json")
		json.NewEncoder(w).Encode(response)
	}))
	defer server.Close()

	factory := &FilterFactory{
		cfg: &Config{
			ServerURL:    server.URL,
			DecisionPath: "/v1/data/test/allow",
			TimeoutMs:    100,
		},
	}

	err := factory.Apply()
	assert.Nil(t, err)

	// Test allow case
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Test_Header", "1")
	rec := httptest.NewRecorder()
	ctx := &contextHttp.HttpContext{
		Writer:  rec,
		Request: req,
		Ctx:     context.Background(),
	}

	chain := &mockFilterChain{}
	factory.PrepareFilterChain(ctx, chain)
	result := chain.filters[0].Decode(ctx)
	assert.Equal(t, filter.Continue, result)

	// Test deny case
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.Header.Set("Test_Header", "0")
	rec2 := httptest.NewRecorder()
	ctx2 := &contextHttp.HttpContext{
		Writer:  rec2,
		Request: req2,
		Ctx:     context.Background(),
	}

	chain2 := &mockFilterChain{}
	factory.PrepareFilterChain(ctx2, chain2)
	result2 := chain2.filters[0].Decode(ctx2)
	assert.Equal(t, filter.Stop, result2)
}

// TestEmbeddedObjectResponse tests embedded mode with object response format
func TestEmbeddedObjectResponse(t *testing.T) {
	// Policy that returns object with allow field
	objectPolicy := `
package test
import future.keywords.if

allow if {
    input.headers[Test_Header][0] == "1"
}

decision := {
    "allow": allow,
    "reason": "test policy"
}
`

	r := rego.New(
		rego.Query("data.test.decision"),
		rego.Module("policy.rego", objectPolicy),
	)

	preparedQuery, err := r.PrepareForEval(context.Background())
	assert.Nil(t, err)

	f := &Filter{
		cfg: &Config{
			Policy:     objectPolicy,
			Entrypoint: "data.test.decision",
		},
		preparedQuery: &preparedQuery,
	}

	// Test allow case
	req := httptest.NewRequest("GET", "/test", nil)
	req.Header.Set("Test_Header", "1")
	rec := httptest.NewRecorder()
	ctx := &contextHttp.HttpContext{
		Writer:  rec,
		Request: req,
		Ctx:     context.Background(),
	}

	result := f.Decode(ctx)
	assert.Equal(t, filter.Continue, result)

	// Test deny case
	req2 := httptest.NewRequest("GET", "/test", nil)
	req2.Header.Set("Test_Header", "0")
	rec2 := httptest.NewRecorder()
	ctx2 := &contextHttp.HttpContext{
		Writer:  rec2,
		Request: req2,
		Ctx:     context.Background(),
	}

	result2 := f.Decode(ctx2)
	assert.Equal(t, filter.Stop, result2)
}

// TestConfigValidation tests configuration validation
func TestConfigValidation(t *testing.T) {
	// Test missing decision_path in server mode
	factory := &FilterFactory{
		cfg: &Config{
			ServerURL: "http://localhost:8181",
		},
	}
	err := factory.Apply()
	assert.NotNil(t, err)
	assert.Contains(t, err.Error(), "decision_path is required")

	// Test missing policy in embedded mode
	factory2 := &FilterFactory{
		cfg: &Config{
			Entrypoint: "data.test.allow",
		},
	}
	err2 := factory2.Apply()
	assert.NotNil(t, err2)
	assert.Contains(t, err2.Error(), "server_url")

	// Test missing entrypoint in embedded mode
	factory3 := &FilterFactory{
		cfg: &Config{
			Policy: "package test",
		},
	}
	err3 := factory3.Apply()
	assert.NotNil(t, err3)
	assert.Contains(t, err3.Error(), "entrypoint is required")
}

// mockFilterChain is a mock implementation of filter.FilterChain for testing
type mockFilterChain struct {
	filters []filter.HttpDecodeFilter
}

func (m *mockFilterChain) AppendDecodeFilters(f ...filter.HttpDecodeFilter) {
	m.filters = append(m.filters, f...)
}

func (m *mockFilterChain) AppendEncodeFilters(f ...filter.HttpEncodeFilter) {
	// Not needed for testing
}

func (m *mockFilterChain) OnDecode(ctx *contextHttp.HttpContext) {
	// Not needed for testing
}

func (m *mockFilterChain) OnEncode(ctx *contextHttp.HttpContext) {
	// Not needed for testing
}

// Backward compatibility test names
func TestAllowedRule(t *testing.T) {
	TestEmbeddedAllowedRule(t)
}

func TestDeniedRule(t *testing.T) {
	TestEmbeddedDeniedRule(t)
}

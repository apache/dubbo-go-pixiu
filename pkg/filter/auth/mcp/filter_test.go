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

package mcp

import (
	"encoding/base64"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"os"
	"path/filepath"
	"testing"
	"time"
)

import (
	"github.com/lestrrat-go/jwx/v3/jwa"
	"github.com/lestrrat-go/jwx/v3/jws"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	dgpfilter "github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/auth/mcp/internal/validator"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// =============================================================================
// Test Helper Functions
// =============================================================================
// writeTempJWKS creates a temporary JWKS file with empty keys for testing
func writeTempJWKS(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "jwks.json")
	if err := os.WriteFile(p, []byte(`{"keys":[]}`), 0o644); err != nil {
		t.Fatalf("write temp jwks: %v", err)
	}
	return "file://" + p
}

// writeHS256JWKS creates a temporary JWKS file with HS256 key for testing
func writeHS256JWKS(t *testing.T, secret []byte) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "jwks.json")
	k := base64.RawURLEncoding.EncodeToString(secret)
	jwks := []byte(`{"keys":[{"kty":"oct","k":"` + k + `","alg":"HS256","use":"sig","kid":"test"}]}`)
	if err := os.WriteFile(p, jwks, 0o644); err != nil {
		t.Fatalf("write jwks: %v", err)
	}
	return "file://" + p
}

// buildFactory creates a FilterFactory for testing with default configuration
func buildFactory(t *testing.T) *FilterFactory {
	t.Helper()
	cfg := &Config{
		ResourceMetadata: ResourceMetadata{
			Path:                 "/.well-known/oauth-protected-resource",
			Resource:             "https://mcp.example.com",
			AuthorizationServers: []string{"https://auth.example.com/.well-known/oauth-authorization-server"},
		},
		Providers: []validator.Provider{
			{
				Name:     "p1",
				Issuer:   "https://issuer.example.com",
				Audience: "mcp-aud",
				JWKS:     writeTempJWKS(t),
			},
		},
		Rules: []Rule{{Cluster: "protected-cluster"}},
	}

	ff := &FilterFactory{cfg: cfg}
	if err := ff.Apply(); err != nil {
		t.Fatalf("apply factory: %v", err)
	}
	return ff
}

// newHttpContext creates a test HTTP context with request and response recorder
func newHttpContext(method, url, host string) (*contexthttp.HttpContext, *httptest.ResponseRecorder) {
	req := httptest.NewRequest(method, url, nil)
	if host != "" {
		req.Host = host
	}
	rr := httptest.NewRecorder()
	hc := &contexthttp.HttpContext{Request: req, Writer: rr}
	return hc, rr
}

// makeHS256JWT creates a signed JWT token for testing
func makeHS256JWT(t *testing.T, secret []byte, iss, aud, scope string) string {
	t.Helper()
	now := time.Now()
	payload := map[string]any{
		"iss": iss,
		"aud": []string{aud},
		"iat": now.Add(-time.Minute).Unix(),
		"exp": now.Add(10 * time.Minute).Unix(),
	}
	if scope != "" {
		payload["scope"] = scope
	}
	pb, _ := json.Marshal(payload)
	hdr := jws.NewHeaders()
	_ = hdr.Set(jws.AlgorithmKey, jwa.HS256())
	_ = hdr.Set(jws.TypeKey, "JWT")
	_ = hdr.Set(jws.KeyIDKey, "test")
	signed, err := jws.Sign(pb, jws.WithKey(jwa.HS256(), secret, jws.WithProtectedHeaders(hdr)))
	if err != nil {
		t.Fatalf("sign jws: %v", err)
	}
	return string(signed)
}

// buildFactoryWithHS256 creates a FilterFactory with HS256 JWKS for testing
func buildFactoryWithHS256(t *testing.T, secret []byte) *FilterFactory {
	t.Helper()
	jwksURI := writeHS256JWKS(t, secret)

	cfg := &Config{
		ResourceMetadata: ResourceMetadata{
			Path:                 "/.well-known/oauth-protected-resource",
			Resource:             "https://mcp.example.com",
			AuthorizationServers: []string{"https://auth.example.com/.well-known/oauth-authorization-server"},
		},
		Providers: []validator.Provider{{
			Name:     "p1",
			Issuer:   "https://issuer.example.com",
			Audience: "mcp-aud",
			JWKS:     jwksURI,
		}},
		Rules: []Rule{{Cluster: "protected-cluster"}},
	}

	ff := &FilterFactory{cfg: cfg}
	if err := ff.Apply(); err != nil {
		t.Fatalf("apply factory: %v", err)
	}
	return ff
}

// =============================================================================
// Core MCP Filter Tests
// =============================================================================

func TestMCPAuth_MetadataEndpoint(t *testing.T) {
	ff := buildFactory(t)
	hc, rr := newHttpContext(http.MethodGet, ff.state.metaPath, "mcp.example.com")

	chain := dgpfilter.NewDefaultFilterChain()
	_ = ff.PrepareFilterChain(hc, chain)
	chain.OnDecode(hc)

	if rr.Code != http.StatusOK {
		t.Fatalf("status = %d, want 200", rr.Code)
	}
	var got struct {
		AuthorizationServers []string `json:"authorization_servers"`
		Resource             string   `json:"resource"`
	}
	if err := json.Unmarshal(rr.Body.Bytes(), &got); err != nil {
		t.Fatalf("unmarshal body: %v", err)
	}
	if got.Resource != ff.cfg.ResourceMetadata.Resource {
		t.Fatalf("resource = %q, want %q", got.Resource, ff.cfg.ResourceMetadata.Resource)
	}
	if len(got.AuthorizationServers) != len(ff.cfg.ResourceMetadata.AuthorizationServers) {
		t.Fatalf("authorization_servers len = %d, want %d", len(got.AuthorizationServers), len(ff.cfg.ResourceMetadata.AuthorizationServers))
	}
}

func TestMCPAuth_MissingToken_Unauthorized(t *testing.T) {
	ff := buildFactory(t)
	// path matches rule, but no Authorization header
	hc, rr := newHttpContext(http.MethodGet, "/api/hello", "mcp.example.com")
	// simulate router matched cluster for protected route
	hc.RouteEntry(&model.RouteAction{Cluster: "protected-cluster"})

	chain := dgpfilter.NewDefaultFilterChain()
	_ = ff.PrepareFilterChain(hc, chain)
	chain.OnDecode(hc)

	if rr.Code != http.StatusUnauthorized {
		t.Fatalf("status = %d, want 401", rr.Code)
	}
	// WWW-Authenticate should include resource_metadata
	wa := rr.Header().Get(constant.WWWAuthenticate)
	if wa == "" {
		t.Fatalf("missing WWW-Authenticate header")
	}
	// response body should be oauth error json
	var body map[string]string
	if err := json.Unmarshal(rr.Body.Bytes(), &body); err != nil {
		t.Fatalf("unmarshal oauth error: %v", err)
	}
	if body["error"] == "" {
		t.Fatalf("missing error in body")
	}
}

// =============================================================================
// JWT Token Tests
// =============================================================================

func TestMCPAuth_NoScopeEnforcement_AllowsRequest(t *testing.T) {
	secret := []byte("secret123")
	ff := buildFactoryWithHS256(t, secret)

	// Token has only read scope but scope is not enforced
	token := makeHS256JWT(t, secret, "https://issuer.example.com", "mcp-aud", "read")
	hc, rr := newHttpContext(http.MethodGet, "/api/hello", "mcp.example.com")
	hc.RouteEntry(&model.RouteAction{Cluster: "protected-cluster"})
	hc.Request.Header.Set("Authorization", "Bearer "+token)

	chain := dgpfilter.NewDefaultFilterChain()
	_ = ff.PrepareFilterChain(hc, chain)
	chain.OnDecode(hc)

	if rr.Code != 0 && rr.Code != http.StatusOK { // no local reply expected
		t.Fatalf("unexpected status = %d", rr.Code)
	}
	if v := hc.Request.Header.Get("Authorization"); v != "" {
		t.Fatalf("Authorization header not removed on success")
	}
}

func TestMCPAuth_Success_RemoveAuthorizationHeader(t *testing.T) {
	secret := []byte("secret123")
	ff := buildFactoryWithHS256(t, secret)

	token := makeHS256JWT(t, secret, "https://issuer.example.com", "mcp-aud", "read write")
	hc, rr := newHttpContext(http.MethodGet, "/api/hello", "mcp.example.com")
	hc.RouteEntry(&model.RouteAction{Cluster: "protected-cluster"})
	hc.Request.Header.Set("Authorization", "Bearer "+token)

	chain := dgpfilter.NewDefaultFilterChain()
	_ = ff.PrepareFilterChain(hc, chain)
	chain.OnDecode(hc)

	if rr.Code != 0 && rr.Code != http.StatusOK { // no local reply expected
		t.Fatalf("unexpected status code: %d", rr.Code)
	}
	if v := hc.Request.Header.Get(constant.Authorization); v != "" {
		t.Fatalf("Authorization header not removed on success")
	}
}

// =============================================================================
// Unit Tests for asValidationError Function
// =============================================================================

func TestAsValidationError_WithValidationError(t *testing.T) {
	// Save original errorsAs function
	originalErrorsAs := errorsAs
	defer func() { errorsAs = originalErrorsAs }()

	// Mock errorsAs to return true (simulating successful type assertion)
	errorsAs = func(err error, target any) bool {
		if verr, ok := target.(*validator.ValidationError); ok {
			verr.Code = "test_code"
			verr.Message = "test message"
		}
		return true
	}

	// Test the function
	testErr := errors.New("some error")
	var verr validator.ValidationError
	result := asValidationError(testErr, &verr)

	if !result {
		t.Fatalf("expected asValidationError to return true")
	}
	if verr.Code != "test_code" {
		t.Fatalf("expected Code to be 'test_code', got %q", verr.Code)
	}
	if verr.Message != "test message" {
		t.Fatalf("expected Message to be 'test message', got %q", verr.Message)
	}
}

func TestAsValidationError_WithNonValidationError(t *testing.T) {
	// Save original errorsAs function
	originalErrorsAs := errorsAs
	defer func() { errorsAs = originalErrorsAs }()

	// Mock errorsAs to return false (simulating failed type assertion)
	errorsAs = func(err error, target any) bool {
		return false
	}

	// Test the function
	testErr := errors.New("some error")
	var verr validator.ValidationError
	result := asValidationError(testErr, &verr)

	if result {
		t.Fatalf("expected asValidationError to return false")
	}
	// verr should remain unchanged
	if verr.Code != "" {
		t.Fatalf("expected Code to remain empty, got %q", verr.Code)
	}
	if verr.Message != "" {
		t.Fatalf("expected Message to remain empty, got %q", verr.Message)
	}
}

func TestAsValidationError_MockBehavior(t *testing.T) {
	// Save original errorsAs function
	originalErrorsAs := errorsAs
	defer func() { errorsAs = originalErrorsAs }()

	// Mock errorsAs to always return true with specific values
	errorsAs = func(err error, target any) bool {
		if verr, ok := target.(*validator.ValidationError); ok {
			verr.Code = "mocked_code"
			verr.Message = "mocked message"
			return true
		}
		return false
	}

	regularErr := errors.New("any error")
	var target validator.ValidationError
	result := asValidationError(regularErr, &target)

	if !result {
		t.Fatalf("expected mocked asValidationError to return true")
	}
	if target.Code != "mocked_code" {
		t.Fatalf("expected mocked Code 'mocked_code', got %q", target.Code)
	}
	if target.Message != "mocked message" {
		t.Fatalf("expected mocked Message 'mocked message', got %q", target.Message)
	}
}

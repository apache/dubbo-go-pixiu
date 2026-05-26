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

package initialize

import (
	"bytes"
	"encoding/json"
	"io"
	"mime/multipart"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

import (
	"github.com/gin-gonic/gin"

	"github.com/golang-jwt/jwt/v4"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/admin/controller/auth"
)

// This file is the only place in the test suite that exercises the full
// admin startup pipeline that handles OPA requests:
//
//   YAML → adminconfig.Bootstrap.OPA
//          → initialize.Routers() registers /config/api/opa/policy
//          → auth.JWTAuth() middleware validates "token" header
//          → opa.PutOPAPolicy/GetOPAPolicy/DeleteOPAPolicy handlers
//          → logic.BizPut/Get/Delete... → real HTTP call to OPA server
//
// We stand up an httptest server as the OPA backend, point
// adminconfig.Bootstrap.OPA.ServerURL at it, and sign JWTs with the
// SAME hardcoded SignKey ("dubbo-go-pixiu") the middleware reads from
// admin/controller/auth/auth.go:83.

// ----- helpers ---------------------------------------------------------------

type recordedRequest struct {
	method string
	path   string
	ctype  string
	auth   string
	body   string
}

type mockOPA struct {
	srv  *httptest.Server
	mu   sync.Mutex
	recs []recordedRequest
	// policies stores PUT'd bodies keyed by policy ID so GETs return them.
	policies map[string]string
	// nextStatus, if non-zero, overrides the success status on the next request.
	nextStatus int
	nextBody   string
}

func newMockOPA(t *testing.T) *mockOPA {
	t.Helper()
	m := &mockOPA{policies: map[string]string{}}
	m.srv = httptest.NewServer(http.HandlerFunc(m.handle))
	t.Cleanup(m.srv.Close)
	return m
}

func (m *mockOPA) handle(w http.ResponseWriter, r *http.Request) {
	b, _ := io.ReadAll(r.Body)
	m.mu.Lock()
	m.recs = append(m.recs, recordedRequest{
		method: r.Method,
		path:   r.URL.Path,
		ctype:  r.Header.Get("Content-Type"),
		auth:   r.Header.Get("Authorization"),
		body:   string(b),
	})
	override, overrideBody := m.nextStatus, m.nextBody
	m.nextStatus, m.nextBody = 0, ""
	m.mu.Unlock()

	if override != 0 {
		w.WriteHeader(override)
		if overrideBody != "" {
			_, _ = w.Write([]byte(overrideBody))
		}
		return
	}

	if strings.HasPrefix(r.URL.Path, "/v1/policies/") {
		id := strings.TrimPrefix(r.URL.Path, "/v1/policies/")
		switch r.Method {
		case http.MethodPut:
			m.mu.Lock()
			m.policies[id] = string(b)
			m.mu.Unlock()
			w.WriteHeader(http.StatusOK)
			_, _ = w.Write([]byte("{}"))
		case http.MethodGet:
			m.mu.Lock()
			raw, ok := m.policies[id]
			m.mu.Unlock()
			if !ok {
				w.WriteHeader(http.StatusNotFound)
				return
			}
			w.Header().Set("Content-Type", "application/json")
			w.WriteHeader(http.StatusOK)
			_ = json.NewEncoder(w).Encode(map[string]any{
				"result": map[string]any{"id": id, "raw": raw},
			})
		case http.MethodDelete:
			m.mu.Lock()
			delete(m.policies, id)
			m.mu.Unlock()
			w.WriteHeader(http.StatusOK)
		}
		return
	}
	w.WriteHeader(http.StatusNotFound)
}

func (m *mockOPA) records() []recordedRequest {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]recordedRequest, len(m.recs))
	copy(out, m.recs)
	return out
}

// signToken signs a JWT with the same hardcoded key the middleware uses,
// so the resulting token passes JWTAuth without any DB lookup.
func signToken(t *testing.T) string {
	t.Helper()
	claims := auth.CustomClaims{
		Username: "e2e",
		StandardClaims: jwt.StandardClaims{
			ExpiresAt: time.Now().Add(time.Hour).Unix(),
			Issuer:    "router-test",
		},
	}
	tok := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	s, err := tok.SignedString([]byte(auth.GetSignKey()))
	if err != nil {
		t.Fatalf("sign: %v", err)
	}
	return s
}

// installRouter builds the real Routers() and points OPAConfig at the mock.
// Restores both gin mode and the global Bootstrap on cleanup.
func installRouter(t *testing.T, m *mockOPA, opaCfg adminconfig.OPAConfig) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	t.Cleanup(func() { gin.SetMode(gin.DebugMode) })

	prev := adminconfig.Bootstrap
	cfg := opaCfg
	if cfg.ServerURL == "" {
		cfg.ServerURL = m.srv.URL
	}
	adminconfig.Bootstrap = &adminconfig.AdminBootstrap{OPA: cfg}
	t.Cleanup(func() { adminconfig.Bootstrap = prev })

	return Routers()
}

func putMultipart(t *testing.T, fields map[string]string) (string, *bytes.Buffer) {
	t.Helper()
	body := &bytes.Buffer{}
	mw := multipart.NewWriter(body)
	for k, v := range fields {
		_ = mw.WriteField(k, v)
	}
	_ = mw.Close()
	return mw.FormDataContentType(), body
}

func doReq(t *testing.T, r *gin.Engine, req *http.Request) *httptest.ResponseRecorder {
	t.Helper()
	w := httptest.NewRecorder()
	r.ServeHTTP(w, req)
	return w
}

// ----- tests ----------------------------------------------------------------

// 1. JWT middleware actually gates /config/api/opa/policy.
func TestOPARoutes_NoTokenRejected(t *testing.T) {
	m := newMockOPA(t)
	r := installRouter(t, m, adminconfig.OPAConfig{PolicyID: "p", RequestTimeout: time.Second})

	w := doReq(t, r, httptest.NewRequest(http.MethodGet, "/config/api/opa/policy", nil))

	if w.Code != http.StatusOK {
		t.Fatalf("status: %d", w.Code)
	}
	if !strings.Contains(w.Body.String(), "does not carry token") {
		t.Fatalf("body should mention token: %s", w.Body.String())
	}
	if len(m.records()) != 0 {
		t.Fatalf("mock OPA should not be called when JWT fails, got %d calls", len(m.records()))
	}
}

// 2. With valid JWT, PUT with no form policy_id uses Bootstrap default → OPA
// gets /v1/policies/<bootstrap-policy-id> with text/plain body.
func TestOPARoutes_PutUsesBootstrapDefaults(t *testing.T) {
	m := newMockOPA(t)
	r := installRouter(t, m, adminconfig.OPAConfig{PolicyID: "from-config", RequestTimeout: 2 * time.Second})

	ctype, body := putMultipart(t, map[string]string{
		"content": "package pixiu\ndefault allow = false",
	})
	req := httptest.NewRequest(http.MethodPut, "/config/api/opa/policy", body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("token", signToken(t))

	w := doReq(t, r, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "Update Success") {
		t.Fatalf("expected Update Success, got %s", w.Body.String())
	}

	recs := m.records()
	if len(recs) != 1 {
		t.Fatalf("expected 1 call to OPA, got %d", len(recs))
	}
	if recs[0].method != http.MethodPut || recs[0].path != "/v1/policies/from-config" {
		t.Errorf("wrong upstream request: %+v", recs[0])
	}
	if recs[0].ctype != "text/plain" {
		t.Errorf("ctype: want text/plain, got %s", recs[0].ctype)
	}
	if recs[0].body != "package pixiu\ndefault allow = false" {
		t.Errorf("body: %q", recs[0].body)
	}
}

// 3. Form policy_id overrides Bootstrap; bearer_token is forwarded as
// Authorization: Bearer ...
func TestOPARoutes_PutFormOverridesAndBearer(t *testing.T) {
	m := newMockOPA(t)
	r := installRouter(t, m, adminconfig.OPAConfig{PolicyID: "from-config", RequestTimeout: time.Second})

	ctype, body := putMultipart(t, map[string]string{
		"policy_id":    "override-id",
		"bearer_token": "secret-123",
		"content":      "package over\nallow = true",
	})
	req := httptest.NewRequest(http.MethodPut, "/config/api/opa/policy", body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("token", signToken(t))

	w := doReq(t, r, req)
	if w.Code != http.StatusOK || !strings.Contains(w.Body.String(), "Update Success") {
		t.Fatalf("status=%d body=%s", w.Code, w.Body.String())
	}

	recs := m.records()
	if recs[0].path != "/v1/policies/override-id" {
		t.Errorf("override failed: path=%s", recs[0].path)
	}
	if recs[0].auth != "Bearer secret-123" {
		t.Errorf("bearer not forwarded: auth=%q", recs[0].auth)
	}
}

// 4. GET via the full chain decodes raw from OPA.
func TestOPARoutes_GetReadsBack(t *testing.T) {
	m := newMockOPA(t)
	m.policies["from-config"] = "package pixiu\nallow = true"
	r := installRouter(t, m, adminconfig.OPAConfig{PolicyID: "from-config", RequestTimeout: time.Second})

	req := httptest.NewRequest(http.MethodGet, "/config/api/opa/policy", nil)
	req.Header.Set("token", signToken(t))

	w := doReq(t, r, req)
	if w.Code != http.StatusOK {
		t.Fatalf("status %d", w.Code)
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode: %v", err)
	}
	if resp["data"] != "package pixiu\nallow = true" {
		t.Errorf("data: %v", resp["data"])
	}
}

// 5. DELETE via the full chain hits the right URL.
func TestOPARoutes_DeleteRoute(t *testing.T) {
	m := newMockOPA(t)
	r := installRouter(t, m, adminconfig.OPAConfig{PolicyID: "from-config", RequestTimeout: time.Second})

	req := httptest.NewRequest(http.MethodDelete, "/config/api/opa/policy", nil)
	req.Header.Set("token", signToken(t))

	if w := doReq(t, r, req); w.Code != http.StatusOK {
		t.Fatalf("status %d body=%s", w.Code, w.Body.String())
	}
	recs := m.records()
	if recs[0].method != http.MethodDelete || recs[0].path != "/v1/policies/from-config" {
		t.Errorf("wrong call: %+v", recs[0])
	}
}

// 6. OPAConfig.RequestTimeout actually fires through the full chain.
// Slow mock + 200ms config → ~200ms-ish elapsed and a context-deadline error
// surfaces in the response body.
func TestOPARoutes_RequestTimeoutThroughFullChain(t *testing.T) {
	slow := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second)
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(slow.Close)

	// Use newMockOPA only to satisfy installRouter's signature; override URL.
	m := newMockOPA(t)
	r := installRouter(t, m, adminconfig.OPAConfig{
		ServerURL:      slow.URL,
		PolicyID:       "p",
		RequestTimeout: 200 * time.Millisecond,
	})

	req := httptest.NewRequest(http.MethodGet, "/config/api/opa/policy", nil)
	req.Header.Set("token", signToken(t))

	start := time.Now()
	w := doReq(t, r, req)
	elapsed := time.Since(start)

	if !strings.Contains(w.Body.String(), "context deadline exceeded") {
		t.Errorf("expected timeout error in body, got: %s", w.Body.String())
	}
	if elapsed > 1500*time.Millisecond {
		t.Errorf("timeout fired too late (%v); config likely ignored", elapsed)
	}
	if elapsed < 150*time.Millisecond {
		t.Errorf("timeout fired too early (%v)", elapsed)
	}
}

// 7. End-to-end loop: PUT a policy through the route, then GET it back —
// proves the admin REST → OPA → admin REST round-trip works with real
// JWT + Gin + httpClient + httptest OPA.
func TestOPARoutes_RoundTrip(t *testing.T) {
	m := newMockOPA(t)
	r := installRouter(t, m, adminconfig.OPAConfig{PolicyID: "rt-policy", RequestTimeout: time.Second})
	token := signToken(t)

	rego := "package rt\nimport rego.v1\ndefault allow := false\nallow if input.method == \"GET\""

	// PUT
	ctype, body := putMultipart(t, map[string]string{"content": rego})
	putReq := httptest.NewRequest(http.MethodPut, "/config/api/opa/policy", body)
	putReq.Header.Set("Content-Type", ctype)
	putReq.Header.Set("token", token)
	if w := doReq(t, r, putReq); w.Code != http.StatusOK ||
		!strings.Contains(w.Body.String(), "Update Success") {
		t.Fatalf("PUT failed: status=%d body=%s", w.Code, w.Body.String())
	}

	// GET
	getReq := httptest.NewRequest(http.MethodGet, "/config/api/opa/policy", nil)
	getReq.Header.Set("token", token)
	w := doReq(t, r, getReq)
	if w.Code != http.StatusOK {
		t.Fatalf("GET status %d", w.Code)
	}
	var resp map[string]any
	if err := json.Unmarshal(w.Body.Bytes(), &resp); err != nil {
		t.Fatalf("decode GET response: %v", err)
	}
	if resp["data"] != rego {
		t.Errorf("round-tripped policy mismatch:\nput:  %q\nback: %v", rego, resp["data"])
	}
}

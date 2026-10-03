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
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"net/url"
	"strings"
	"sync"
	"testing"
	"time"
)

import (
	"github.com/gin-gonic/gin"

	"github.com/open-policy-agent/opa/rego"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contextHttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	opaFilter "github.com/apache/dubbo-go-pixiu/pkg/filter/opa"
)

// This file is the PR2 deliverable: a CI-runnable end-to-end test that exercises
// the full admin → OPA Server → gateway filter chain in a single process.
//
//	admin REST PUT /config/api/opa/policy
//	     → JWT auth → controller → logic → HTTP PUT /v1/policies/<id>
//	         ↓
//	     regoMockOPA (httptest) — stores rego module text AND compiles it
//	         ↑
//	gateway OPA filter POST {server_url}/v1/data/<path>
//	     ← input(method, path, headers, ...) → rego evaluation
//	     → filter.Continue (allow) | filter.Stop+403 (deny)
//
// The mock OPA is "smart" — it uses the real github.com/open-policy-agent/opa
// rego library that pkg/filter/opa already depends on, so the policy
// evaluation in the test is identical to what a real OPA server would do.
// No docker, no etcd, no external process needed.
// ---------------------------------------------------------------------------
// regoMockOPA: an in-process OPA server that speaks the subset of the OPA REST
// API exercised by the admin + gateway: PUT/GET/DELETE /v1/policies/{id} and
// POST /v1/data/<any/path>.
// ---------------------------------------------------------------------------
type regoMockOPA struct {
	srv *httptest.Server

	mu       sync.Mutex
	policies map[string]string // policy_id -> rego module text

	// putRequests records every PUT for assertions (auth header, body, etc.).
	putRequests []recordedRequest

	// decisionDelay, if set, makes the POST /v1/data/... handler sleep before
	// evaluating. Lets tests exercise gateway-side timeouts.
	decisionDelay time.Duration
}

func newRegoMockOPA(t *testing.T) *regoMockOPA {
	t.Helper()
	useDirectHTTPTransport(t)
	m := &regoMockOPA{policies: map[string]string{}}
	m.srv = startLoopbackHTTPServer(t, http.HandlerFunc(m.handle))
	return m
}

func (m *regoMockOPA) URL() string { return m.srv.URL }

func (m *regoMockOPA) handle(w http.ResponseWriter, r *http.Request) {
	switch {
	case strings.HasPrefix(r.URL.Path, "/v1/policies/"):
		m.handlePolicy(w, r)
	case strings.HasPrefix(r.URL.Path, "/v1/data/"):
		m.handleDecision(w, r)
	default:
		w.WriteHeader(http.StatusNotFound)
	}
}

func (m *regoMockOPA) handlePolicy(w http.ResponseWriter, r *http.Request) {
	id := strings.TrimPrefix(r.URL.Path, "/v1/policies/")
	body, _ := io.ReadAll(r.Body)

	switch r.Method {
	case http.MethodPut:
		// Sanity-check the rego compiles before accepting; mirrors real OPA's
		// behavior of returning 400 with a compile error transcript.
		if _, err := rego.New(
			rego.Query("data"),
			rego.Module(id, string(body)),
		).PrepareForEval(context.Background()); err != nil {
			w.WriteHeader(http.StatusBadRequest)
			_, _ = w.Write([]byte(err.Error()))
			return
		}
		m.mu.Lock()
		m.policies[id] = string(body)
		m.putRequests = append(m.putRequests, recordedRequest{
			method: r.Method,
			path:   r.URL.Path,
			ctype:  r.Header.Get("Content-Type"),
			auth:   r.Header.Get("Authorization"),
			body:   string(body),
		})
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
		_ = json.NewEncoder(w).Encode(map[string]any{
			"result": map[string]any{"id": id, "raw": raw},
		})

	case http.MethodDelete:
		m.mu.Lock()
		delete(m.policies, id)
		m.mu.Unlock()
		w.WriteHeader(http.StatusOK)

	default:
		w.WriteHeader(http.StatusMethodNotAllowed)
	}
}

// handleDecision implements POST /v1/data/<rule path>. It bundles every stored
// policy into a single rego module set, runs the query, and replies with
// {"result": <value>} — matching the OPA REST contract that
// pkg/filter/opa/opa.go expects in evaluateServer().
func (m *regoMockOPA) handleDecision(w http.ResponseWriter, r *http.Request) {
	if r.Method != http.MethodPost {
		w.WriteHeader(http.StatusMethodNotAllowed)
		return
	}

	if m.decisionDelay > 0 {
		time.Sleep(m.decisionDelay)
	}

	var reqBody map[string]any
	if err := json.NewDecoder(r.Body).Decode(&reqBody); err != nil {
		w.WriteHeader(http.StatusBadRequest)
		return
	}
	input := reqBody["input"]

	// /v1/data/http/authz/allow → data.http.authz.allow
	rulePath := strings.TrimPrefix(r.URL.Path, "/v1/data/")
	query := "data." + strings.ReplaceAll(rulePath, "/", ".")

	m.mu.Lock()
	modules := make(map[string]string, len(m.policies))
	for id, raw := range m.policies {
		modules[id] = raw
	}
	m.mu.Unlock()

	opts := []func(r *rego.Rego){rego.Query(query)}
	for id, raw := range modules {
		opts = append(opts, rego.Module(id, raw))
	}
	pq, err := rego.New(opts...).PrepareForEval(context.Background())
	if err != nil {
		// Real OPA returns 500 on bundle compile errors at decision time; the
		// filter treats non-200 as BadGateway, which is what we want.
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	results, err := pq.Eval(r.Context(), rego.EvalInput(input))
	if err != nil {
		w.WriteHeader(http.StatusInternalServerError)
		_, _ = w.Write([]byte(err.Error()))
		return
	}

	resp := map[string]any{}
	// No rules matched → omit "result" entirely. This is exactly what real
	// OPA does, and it triggers the gateway's "missing 'result' field" branch
	// — the same fail-closed behavior documented in test_opa.md §6.6.
	if len(results) > 0 && len(results[0].Expressions) > 0 {
		resp["result"] = results[0].Expressions[0].Value
	}
	w.Header().Set("Content-Type", "application/json")
	_ = json.NewEncoder(w).Encode(resp)
}

func (m *regoMockOPA) recordedPUTs() []recordedRequest {
	m.mu.Lock()
	defer m.mu.Unlock()
	out := make([]recordedRequest, len(m.putRequests))
	copy(out, m.putRequests)
	return out
}

// ---------------------------------------------------------------------------
// Helpers for wiring admin router and gateway OPA filter against the mock.
// ---------------------------------------------------------------------------

// installAdminRouterWithRegoMock mounts the real admin router with
// adminconfig.Bootstrap pointed at the given mock, mirroring installRouter()
// in router_opa_test.go but accepting a regoMockOPA instead.
func installAdminRouterWithRegoMock(t *testing.T, m *regoMockOPA, opaCfg adminconfig.OPAConfig) *gin.Engine {
	t.Helper()
	gin.SetMode(gin.TestMode)
	t.Cleanup(func() { gin.SetMode(gin.DebugMode) })

	prev := adminconfig.Bootstrap
	cfg := opaCfg
	if cfg.ServerURL == "" {
		cfg.ServerURL = m.URL()
	}
	adminconfig.Bootstrap = &adminconfig.AdminBootstrap{OPA: cfg}
	t.Cleanup(func() { adminconfig.Bootstrap = prev })

	return Routers()
}

// adminPutPolicy uses the real /config/api/opa/policy PUT route, signed with
// the same JWT key the middleware reads. The full request travels through
// gin → JWT middleware → controller → logic → mock OPA, just like in prod.
func adminPutPolicy(t *testing.T, r *gin.Engine, m *regoMockOPA, policyID, content string) {
	t.Helper()
	before := len(m.recordedPUTs())
	fields := map[string]string{
		"content":    content,
		"server_url": m.URL(),
	}
	if policyID != "" {
		fields["policy_id"] = policyID
	}
	ctype, body := putMultipart(t, fields)
	req := httptest.NewRequest(http.MethodPut, "/config/api/opa/policy", body)
	req.Header.Set("Content-Type", ctype)
	req.Header.Set("token", signToken(t))

	w := doReq(t, r, req)
	if w.Code != http.StatusOK {
		t.Fatalf("admin PUT failed: status=%d body=%s", w.Code, w.Body.String())
	}
	if !strings.Contains(w.Body.String(), "Update Success") {
		t.Fatalf("admin PUT expected Update Success, got %s", w.Body.String())
	}

	puts := m.recordedPUTs()
	if len(puts) != before+1 {
		t.Fatalf("admin PUT did not reach mock OPA: before=%d after=%d calls=%+v", before, len(puts), puts)
	}
	if policyID != "" && puts[len(puts)-1].path != "/v1/policies/"+policyID {
		t.Fatalf("admin PUT reached wrong OPA policy path: want %s got %s", "/v1/policies/"+policyID, puts[len(puts)-1].path)
	}
}

// adminDeletePolicy hits DELETE /config/api/opa/policy.
func adminDeletePolicy(t *testing.T, r *gin.Engine, serverURL, policyID string) {
	t.Helper()
	target := "/config/api/opa/policy"
	query := url.Values{}
	if serverURL != "" {
		query.Set("server_url", serverURL)
	}
	if policyID != "" {
		query.Set("policy_id", policyID)
	}
	if len(query) > 0 {
		target = target + "?" + query.Encode()
	}
	req := httptest.NewRequest(http.MethodDelete, target, nil)
	req.Header.Set("token", signToken(t))
	w := doReq(t, r, req)
	if w.Code != http.StatusOK {
		t.Fatalf("admin DELETE failed: status=%d body=%s", w.Code, w.Body.String())
	}
}

// buildGatewayFilter constructs and applies the real gateway OPA filter
// (pkg/filter/opa) pointed at the same mock OPA the admin writes to.
func buildGatewayFilter(t *testing.T, mockURL, decisionPath string, timeoutMs int) filter.HttpDecodeFilter {
	t.Helper()
	plugin := &opaFilter.Plugin{}
	factory, err := plugin.CreateFilterFactory()
	if err != nil {
		t.Fatalf("create filter factory: %v", err)
	}
	cfg := factory.Config().(*opaFilter.Config)
	cfg.ServerURL = mockURL
	cfg.DecisionPath = decisionPath
	cfg.TimeoutMs = timeoutMs
	if err := factory.Apply(); err != nil {
		t.Fatalf("apply gateway filter: %v", err)
	}

	chain := &e2eFilterChain{}
	ctxStub := &contextHttp.HttpContext{
		Request: httptest.NewRequest(http.MethodGet, "/", nil),
		Writer:  httptest.NewRecorder(),
		Ctx:     context.Background(),
	}
	if err := factory.PrepareFilterChain(ctxStub, chain); err != nil {
		t.Fatalf("prepare gateway filter chain: %v", err)
	}
	if len(chain.filters) != 1 {
		t.Fatalf("expected 1 decode filter, got %d", len(chain.filters))
	}
	return chain.filters[0]
}

// driveGatewayRequest runs one HTTP request through the gateway OPA filter
// and returns the FilterStatus plus the captured HttpContext for status code
// and response body inspection.
func driveGatewayRequest(t *testing.T, f filter.HttpDecodeFilter, method, path string, headers map[string]string) (filter.FilterStatus, *contextHttp.HttpContext) {
	t.Helper()
	req := httptest.NewRequest(method, path, nil)
	for k, v := range headers {
		req.Header.Set(k, v)
	}
	ctx := &contextHttp.HttpContext{
		Writer:  httptest.NewRecorder(),
		Request: req,
		Ctx:     context.Background(),
	}
	return f.Decode(ctx), ctx
}

type e2eFilterChain struct {
	filters []filter.HttpDecodeFilter
}

func (c *e2eFilterChain) AppendDecodeFilters(f ...filter.HttpDecodeFilter) {
	c.filters = append(c.filters, f...)
}
func (c *e2eFilterChain) AppendEncodeFilters(f ...filter.HttpEncodeFilter) {}
func (c *e2eFilterChain) OnDecode(ctx *contextHttp.HttpContext)            {}
func (c *e2eFilterChain) OnEncode(ctx *contextHttp.HttpContext)            {}

// ---------------------------------------------------------------------------
// Scenarios
// ---------------------------------------------------------------------------

const (
	e2ePolicyID     = "pixiu-authz"
	e2eDecisionPath = "/v1/data/pixiu/authz/allow"

	// "allow GET only" — covers the headline allow case.
	allowGETPolicy = `package pixiu.authz
import future.keywords.if
default allow := false
allow if input.method == "GET"
`

	// "default allow := false" only — every request denied.
	denyAllPolicy = `package pixiu.authz
import future.keywords.if
default allow := false
`

	// "allow if header X-Role == admin" — exercises header propagation.
	headerRolePolicy = `package pixiu.authz
import future.keywords.if
default allow := false
allow if input.headers["X-Role"][0] == "admin"
`
)

//  1. Admin PUTs an "allow GET" policy. Gateway GET request is allowed
//     end-to-end: filter returns Continue, no local reply written.
func TestE2E_AllowedThroughFullChain(t *testing.T) {
	mock := newRegoMockOPA(t)
	r := installAdminRouterWithRegoMock(t, mock, adminconfig.OPAConfig{
		PolicyID:       e2ePolicyID,
		RequestTimeout: 2 * time.Second,
	})

	adminPutPolicy(t, r, mock, e2ePolicyID, allowGETPolicy)

	// Verify the PUT actually reached OPA (full admin chain works).
	puts := mock.recordedPUTs()
	if len(puts) != 1 || puts[0].path != "/v1/policies/"+e2ePolicyID {
		t.Fatalf("admin PUT didn't reach mock OPA correctly: %+v", puts)
	}

	gw := buildGatewayFilter(t, mock.URL(), e2eDecisionPath, 2000)
	status, ctx := driveGatewayRequest(t, gw, http.MethodGet, "/anything", nil)

	if status != filter.Continue {
		t.Fatalf("GET should be allowed, got status=%v code=%d body=%s",
			status, ctx.GetStatusCode(), string(ctx.GetLocalReplyBody()))
	}
	if ctx.LocalReply() {
		t.Errorf("Continue must not write a local reply")
	}
}

//  2. Same policy, gateway POST is denied — proves the deny path returns
//     filter.Stop with 403 and that the gateway short-circuits before reaching
//     any upstream.
func TestE2E_DeniedThroughFullChain(t *testing.T) {
	mock := newRegoMockOPA(t)
	r := installAdminRouterWithRegoMock(t, mock, adminconfig.OPAConfig{
		PolicyID:       e2ePolicyID,
		RequestTimeout: 2 * time.Second,
	})
	adminPutPolicy(t, r, mock, e2ePolicyID, allowGETPolicy)

	gw := buildGatewayFilter(t, mock.URL(), e2eDecisionPath, 2000)
	status, ctx := driveGatewayRequest(t, gw, http.MethodPost, "/anything", nil)

	if status != filter.Stop {
		t.Fatalf("POST should be denied, got %v", status)
	}
	if ctx.GetStatusCode() != http.StatusForbidden {
		t.Errorf("deny status: want 403, got %d body=%s",
			ctx.GetStatusCode(), string(ctx.GetLocalReplyBody()))
	}
}

//  3. "default allow := false" with no allow rule — every method/path denied.
//     Subtests share one mock+filter to confirm the deny is policy-driven, not
//     request-shape-dependent.
func TestE2E_DefaultDenyForAllRequests(t *testing.T) {
	mock := newRegoMockOPA(t)
	r := installAdminRouterWithRegoMock(t, mock, adminconfig.OPAConfig{
		PolicyID:       e2ePolicyID,
		RequestTimeout: 2 * time.Second,
	})
	adminPutPolicy(t, r, mock, e2ePolicyID, denyAllPolicy)
	gw := buildGatewayFilter(t, mock.URL(), e2eDecisionPath, 2000)

	cases := []struct {
		method, path string
	}{
		{http.MethodGet, "/"},
		{http.MethodGet, "/users/1"},
		{http.MethodPost, "/api/x"},
		{http.MethodPut, "/anything"},
		{http.MethodDelete, "/secret"},
	}
	for _, c := range cases {
		c := c
		t.Run(c.method+" "+c.path, func(t *testing.T) {
			status, ctx := driveGatewayRequest(t, gw, c.method, c.path, nil)
			if status != filter.Stop || ctx.GetStatusCode() != http.StatusForbidden {
				t.Errorf("expected deny+403, got status=%v code=%d",
					status, ctx.GetStatusCode())
			}
		})
	}
}

//  4. Policy hot-reload through the admin REST API — gateway sees the new
//     decision on the *next* request, with no restart. This is the key value
//     proposition of OPA server mode vs. embedded mode.
func TestE2E_PolicyHotReload(t *testing.T) {
	mock := newRegoMockOPA(t)
	r := installAdminRouterWithRegoMock(t, mock, adminconfig.OPAConfig{
		PolicyID:       e2ePolicyID,
		RequestTimeout: 2 * time.Second,
	})
	gw := buildGatewayFilter(t, mock.URL(), e2eDecisionPath, 2000)

	// v1: only GET allowed.
	adminPutPolicy(t, r, mock, e2ePolicyID, allowGETPolicy)
	if status, _ := driveGatewayRequest(t, gw, http.MethodGet, "/", nil); status != filter.Continue {
		t.Fatalf("v1: GET should be allowed, got %v", status)
	}
	if status, _ := driveGatewayRequest(t, gw, http.MethodPost, "/", nil); status != filter.Stop {
		t.Fatalf("v1: POST should be denied, got %v", status)
	}

	// v2: flip — only POST allowed.
	adminPutPolicy(t, r, mock, e2ePolicyID, `package pixiu.authz
import future.keywords.if
default allow := false
allow if input.method == "POST"
`)

	if status, _ := driveGatewayRequest(t, gw, http.MethodGet, "/", nil); status != filter.Stop {
		t.Errorf("v2: GET should now be denied, got %v", status)
	}
	if status, _ := driveGatewayRequest(t, gw, http.MethodPost, "/", nil); status != filter.Continue {
		t.Errorf("v2: POST should now be allowed, got %v", status)
	}
}

//  5. After DELETE, no rules are loaded → mock OPA returns a body with no
//     "result" field, the gateway filter must fail closed with BadGateway.
//     This locks in the §6.6 invariant from test_opa.md.
func TestE2E_DeleteCausesMissingResultFailClosed(t *testing.T) {
	mock := newRegoMockOPA(t)
	r := installAdminRouterWithRegoMock(t, mock, adminconfig.OPAConfig{
		PolicyID:       e2ePolicyID,
		RequestTimeout: 2 * time.Second,
	})
	adminPutPolicy(t, r, mock, e2ePolicyID, allowGETPolicy)
	gw := buildGatewayFilter(t, mock.URL(), e2eDecisionPath, 2000)

	// Sanity: allowed before delete.
	if status, _ := driveGatewayRequest(t, gw, http.MethodGet, "/", nil); status != filter.Continue {
		t.Fatalf("pre-delete: GET should be allowed, got %v", status)
	}

	adminDeletePolicy(t, r, mock.URL(), e2ePolicyID)

	status, ctx := driveGatewayRequest(t, gw, http.MethodGet, "/", nil)
	if status != filter.Stop {
		t.Fatalf("post-delete: expected Stop, got %v", status)
	}
	if ctx.GetStatusCode() != http.StatusBadGateway {
		t.Errorf("post-delete: expected 502 (missing 'result'), got %d body=%s",
			ctx.GetStatusCode(), string(ctx.GetLocalReplyBody()))
	}
}

//  6. Header-based policy: gateway forwards input.headers to OPA, and headers
//     are canonicalised by net/http to the Title-Case form
//     (X-Role, not x-role). This proves the gateway request shape contract.
func TestE2E_HeaderBasedAllowDeny(t *testing.T) {
	mock := newRegoMockOPA(t)
	r := installAdminRouterWithRegoMock(t, mock, adminconfig.OPAConfig{
		PolicyID:       e2ePolicyID,
		RequestTimeout: 2 * time.Second,
	})
	adminPutPolicy(t, r, mock, e2ePolicyID, headerRolePolicy)
	gw := buildGatewayFilter(t, mock.URL(), e2eDecisionPath, 2000)

	t.Run("admin header allowed", func(t *testing.T) {
		status, ctx := driveGatewayRequest(t, gw, http.MethodGet, "/admin",
			map[string]string{"X-Role": "admin"})
		if status != filter.Continue {
			t.Errorf("expected Continue, got %v code=%d body=%s",
				status, ctx.GetStatusCode(), string(ctx.GetLocalReplyBody()))
		}
	})

	t.Run("user header denied", func(t *testing.T) {
		status, ctx := driveGatewayRequest(t, gw, http.MethodGet, "/admin",
			map[string]string{"X-Role": "user"})
		if status != filter.Stop || ctx.GetStatusCode() != http.StatusForbidden {
			t.Errorf("expected Stop+403, got status=%v code=%d",
				status, ctx.GetStatusCode())
		}
	})

	t.Run("missing header denied", func(t *testing.T) {
		status, ctx := driveGatewayRequest(t, gw, http.MethodGet, "/admin", nil)
		if status != filter.Stop || ctx.GetStatusCode() != http.StatusForbidden {
			t.Errorf("expected Stop+403, got status=%v code=%d",
				status, ctx.GetStatusCode())
		}
	})
}

//  7. Slow mock OPA + short gateway timeout → gateway returns 504 GatewayTimeout
//     on the next decision. Verifies the gateway-side timeout config really
//     fires under network slowness, matching the §6.4 manual finding.
func TestE2E_GatewayTimeoutFailClosed(t *testing.T) {
	mock := newRegoMockOPA(t)
	r := installAdminRouterWithRegoMock(t, mock, adminconfig.OPAConfig{
		PolicyID:       e2ePolicyID,
		RequestTimeout: 2 * time.Second,
	})
	adminPutPolicy(t, r, mock, e2ePolicyID, allowGETPolicy)

	// 200ms delay on decision; 50ms filter timeout → must time out.
	mock.decisionDelay = 200 * time.Millisecond
	gw := buildGatewayFilter(t, mock.URL(), e2eDecisionPath, 50)

	start := time.Now()
	status, ctx := driveGatewayRequest(t, gw, http.MethodGet, "/", nil)
	elapsed := time.Since(start)

	if status != filter.Stop {
		t.Fatalf("expected Stop on timeout, got %v", status)
	}
	if ctx.GetStatusCode() != http.StatusGatewayTimeout {
		t.Errorf("expected 504, got %d body=%s",
			ctx.GetStatusCode(), string(ctx.GetLocalReplyBody()))
	}
	if elapsed > 180*time.Millisecond {
		t.Errorf("timeout fired too late (%v); 50ms config likely ignored", elapsed)
	}
}

//  8. PUT a policy with a different policy_id via the admin REST form override,
//     then point the gateway's decision path at the rule of *that* policy. This
//     proves the override flag in PR1's controller propagates all the way
//     through to a working gateway decision.
func TestE2E_PolicyIDOverrideRoutesThroughGateway(t *testing.T) {
	mock := newRegoMockOPA(t)
	r := installAdminRouterWithRegoMock(t, mock, adminconfig.OPAConfig{
		PolicyID:       e2ePolicyID, // Bootstrap default; we'll override below.
		RequestTimeout: 2 * time.Second,
	})

	const overrideID = "tenant-a-policy"
	const overridePackage = `package tenant.a
import future.keywords.if
default allow := false
allow if input.method == "GET"
`
	adminPutPolicy(t, r, mock, overrideID, overridePackage)

	puts := mock.recordedPUTs()
	if len(puts) != 1 || puts[0].path != "/v1/policies/"+overrideID {
		t.Fatalf("admin override didn't land at overridden policy_id: %+v", puts)
	}

	// Gateway points at the override package's rule.
	gw := buildGatewayFilter(t, mock.URL(), "/v1/data/tenant/a/allow", 2000)
	if status, _ := driveGatewayRequest(t, gw, http.MethodGet, "/", nil); status != filter.Continue {
		t.Errorf("override GET should be allowed")
	}
	if status, _ := driveGatewayRequest(t, gw, http.MethodPost, "/", nil); status != filter.Stop {
		t.Errorf("override POST should be denied")
	}
}

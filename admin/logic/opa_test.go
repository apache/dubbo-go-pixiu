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

package logic

import (
	"io"
	"net/http"
	"net/http/httptest"
	"strings"
	"sync"
	"testing"
	"time"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
)

// setBootstrap installs a temporary Bootstrap for the test and restores it after.
func setBootstrap(t *testing.T, b *adminconfig.AdminBootstrap) {
	t.Helper()
	prev := adminconfig.Bootstrap
	adminconfig.Bootstrap = b
	t.Cleanup(func() { adminconfig.Bootstrap = prev })
}

func TestGetOPATimeout_FallsBackToDefaultWhenUnset(t *testing.T) {
	setBootstrap(t, nil)
	if got := getOPATimeout(); got != adminconfig.DefaultOPAPolicyTimeout {
		t.Fatalf("nil Bootstrap: want default %v, got %v", adminconfig.DefaultOPAPolicyTimeout, got)
	}

	setBootstrap(t, &adminconfig.AdminBootstrap{}) // RequestTimeout == 0
	if got := getOPATimeout(); got != adminconfig.DefaultOPAPolicyTimeout {
		t.Fatalf("zero RequestTimeout: want default %v, got %v", adminconfig.DefaultOPAPolicyTimeout, got)
	}
}

func TestGetOPATimeout_UsesConfigValue(t *testing.T) {
	setBootstrap(t, &adminconfig.AdminBootstrap{
		OPA: adminconfig.OPAConfig{RequestTimeout: 3 * time.Second},
	})
	if got := getOPATimeout(); got != 3*time.Second {
		t.Fatalf("want 3s from config, got %v", got)
	}
}

func TestBuildOPAPolicyURL(t *testing.T) {
	cases := []struct {
		name, server, policy, want string
		wantErr                    bool
	}{
		{"basic", "http://opa:8181", "pid", "http://opa:8181/v1/policies/pid", false},
		{"trims trailing slash", "http://opa:8181/", "pid", "http://opa:8181/v1/policies/pid", false},
		{"trims whitespace", "  http://opa:8181  ", "  pid  ", "http://opa:8181/v1/policies/pid", false},
		{"empty server", "", "pid", "", true},
		{"empty policy", "http://opa:8181", "", "", true},
	}
	for _, c := range cases {
		t.Run(c.name, func(t *testing.T) {
			got, err := buildOPAPolicyURL(c.server, c.policy)
			if c.wantErr {
				if err == nil {
					t.Fatalf("want error, got url=%q", got)
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if got != c.want {
				t.Fatalf("want %q, got %q", c.want, got)
			}
		})
	}
}

type recordedRequest struct {
	method        string
	path          string
	contentType   string
	authorization string
	body          string
}

// startMockOPA spins up an httptest server that records each request and
// responds with `status` and `body` (body may be empty).
func startMockOPA(t *testing.T, status int, body string) (*httptest.Server, *[]recordedRequest, *sync.Mutex) {
	t.Helper()
	var (
		mu   sync.Mutex
		recs []recordedRequest
	)
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		b, _ := io.ReadAll(r.Body)
		mu.Lock()
		recs = append(recs, recordedRequest{
			method:        r.Method,
			path:          r.URL.Path,
			contentType:   r.Header.Get("Content-Type"),
			authorization: r.Header.Get("Authorization"),
			body:          string(b),
		})
		mu.Unlock()
		w.WriteHeader(status)
		if body != "" {
			_, _ = w.Write([]byte(body))
		}
	}))
	t.Cleanup(srv.Close)
	return srv, &recs, &mu
}

func TestBizPutOPAPolicy_SendsCorrectRequest(t *testing.T) {
	srv, recs, mu := startMockOPA(t, http.StatusNoContent, "")

	err := BizPutOPAPolicy(srv.URL, "my-policy", "tok123", "package pixiu\r\ndefault allow = false")
	if err != nil {
		t.Fatalf("BizPutOPAPolicy: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if len(*recs) != 1 {
		t.Fatalf("want 1 request, got %d", len(*recs))
	}
	r := (*recs)[0]
	if r.method != http.MethodPut {
		t.Errorf("method: want PUT, got %s", r.method)
	}
	if r.path != "/v1/policies/my-policy" {
		t.Errorf("path: want /v1/policies/my-policy, got %s", r.path)
	}
	if r.contentType != "text/plain" {
		t.Errorf("content-type: want text/plain, got %s", r.contentType)
	}
	if r.authorization != "Bearer tok123" {
		t.Errorf("authorization: want 'Bearer tok123', got %q", r.authorization)
	}
	// \r\n must be normalized to \n
	if r.body != "package pixiu\ndefault allow = false" {
		t.Errorf("body not normalized: %q", r.body)
	}
}

func TestBizPutOPAPolicy_NoBearerTokenOmitsHeader(t *testing.T) {
	srv, recs, mu := startMockOPA(t, http.StatusNoContent, "")

	if err := BizPutOPAPolicy(srv.URL, "p", "", "package x"); err != nil {
		t.Fatalf("BizPutOPAPolicy: %v", err)
	}

	mu.Lock()
	defer mu.Unlock()
	if (*recs)[0].authorization != "" {
		t.Errorf("expected no Authorization header, got %q", (*recs)[0].authorization)
	}
}

func TestBizGetOPAPolicy_DecodesRaw(t *testing.T) {
	body := `{"result":{"id":"p","raw":"package pixiu\ndefault allow = true"}}`
	srv, recs, mu := startMockOPA(t, http.StatusOK, body)

	got, err := BizGetOPAPolicy(srv.URL, "p", "")
	if err != nil {
		t.Fatalf("BizGetOPAPolicy: %v", err)
	}
	if got != "package pixiu\ndefault allow = true" {
		t.Errorf("decoded raw mismatch: %q", got)
	}

	mu.Lock()
	defer mu.Unlock()
	if (*recs)[0].method != http.MethodGet || (*recs)[0].path != "/v1/policies/p" {
		t.Errorf("unexpected request: %+v", (*recs)[0])
	}
}

func TestBizGetOPAPolicy_NotFoundReturnsEmpty(t *testing.T) {
	srv, _, _ := startMockOPA(t, http.StatusNotFound, "")
	got, err := BizGetOPAPolicy(srv.URL, "missing", "")
	if err != nil {
		t.Fatalf("404 should not be an error: %v", err)
	}
	if got != "" {
		t.Errorf("404 should yield empty string, got %q", got)
	}
}

func TestBizDeleteOPAPolicy_SendsDelete(t *testing.T) {
	srv, recs, mu := startMockOPA(t, http.StatusNoContent, "")
	if err := BizDeleteOPAPolicy(srv.URL, "p", ""); err != nil {
		t.Fatalf("BizDeleteOPAPolicy: %v", err)
	}
	mu.Lock()
	defer mu.Unlock()
	if (*recs)[0].method != http.MethodDelete {
		t.Errorf("want DELETE, got %s", (*recs)[0].method)
	}
}

func TestBizDeleteOPAPolicy_NotFoundIsNil(t *testing.T) {
	srv, _, _ := startMockOPA(t, http.StatusNotFound, "")
	if err := BizDeleteOPAPolicy(srv.URL, "gone", ""); err != nil {
		t.Errorf("404 on DELETE should be nil error, got %v", err)
	}
}

// TestOPARequestTimeout_HonorsConfig is the critical test: proves that
// adminconfig.Bootstrap.OPA.RequestTimeout actually wraps the OPA request
// context, rather than being ignored in favor of DefaultOPAPolicyTimeout (8s)
// or opaHTTPClient.Timeout (30s).
func TestOPARequestTimeout_HonorsConfig(t *testing.T) {
	srv := httptest.NewServer(http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		time.Sleep(2 * time.Second) // far longer than configured timeout
		w.WriteHeader(http.StatusOK)
	}))
	t.Cleanup(srv.Close)

	setBootstrap(t, &adminconfig.AdminBootstrap{
		OPA: adminconfig.OPAConfig{RequestTimeout: 200 * time.Millisecond},
	})

	start := time.Now()
	_, err := BizGetOPAPolicy(srv.URL, "p", "")
	elapsed := time.Since(start)

	if err == nil {
		t.Fatal("want timeout error, got nil")
	}
	if !strings.Contains(err.Error(), "context deadline exceeded") {
		t.Errorf("want 'context deadline exceeded', got %v", err)
	}
	// Should fire near 200ms; allow generous upper bound to avoid CI flake but
	// well below the 8s default and 30s client fallback.
	if elapsed > 1500*time.Millisecond {
		t.Errorf("timeout fired too late (%v); config likely ignored", elapsed)
	}
	if elapsed < 150*time.Millisecond {
		t.Errorf("timeout fired too early (%v); something other than config drove it", elapsed)
	}
}

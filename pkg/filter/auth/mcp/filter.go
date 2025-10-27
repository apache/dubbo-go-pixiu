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
	"encoding/json"
	"errors"
	"fmt"
	"net/http"
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/auth/mcp/internal/validator"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	// Kind is the filter kind key for MCP auth
	Kind = constant.HTTPMCPAuthFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is http filter plugin.
	Plugin struct{}

	// runtimeState holds read-only runtime data for filters
	runtimeState struct {
		validator *validator.Validator
		metaPath  string
		metaBody  []byte
		rules     []Rule
	}

	// FilterFactory holds immutable state for creating filters
	FilterFactory struct {
		cfg   *Config
		state *runtimeState
	}

	// Filter is the runtime decode filter
	Filter struct {
		state *runtimeState
	}
)

func (p *Plugin) Kind() string { return Kind }

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

func (factory *FilterFactory) Config() any { return factory.cfg }

// Apply initializes the validator and prebuilds resource metadata body
func (factory *FilterFactory) Apply() error {
	if err := factory.cfg.Validate(); err != nil {
		return err
	}

	v, err := validator.NewValidator(validator.Config{Providers: factory.cfg.Providers})
	if err != nil {
		return fmt.Errorf("init validator: %w", err)
	}
	metaPath := factory.cfg.ResourceMetadata.Path
	// Minimal RFC9728 document: resource + authorization_servers
	meta := struct {
		AuthorizationServers []string `json:"authorization_servers"`
		Resource             string   `json:"resource"`
	}{
		AuthorizationServers: factory.cfg.ResourceMetadata.AuthorizationServers,
		Resource:             factory.cfg.ResourceMetadata.Resource,
	}
	body, err := json.Marshal(meta)
	if err != nil {
		return fmt.Errorf("marshal resource metadata: %w", err)
	}
	rules := factory.cfg.Rules

	factory.state = &runtimeState{
		validator: v,
		metaPath:  metaPath,
		metaBody:  body,
		rules:     rules,
	}
	return nil
}

// PrepareFilterChain appends the decode filter to chain
func (factory *FilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	f := &Filter{state: factory.state}
	chain.AppendDecodeFilters(f)
	return nil
}

// Decode implements MCP auth flow
func (f *Filter) Decode(hc *contexthttp.HttpContext) filter.FilterStatus {
	path := hc.GetUrl()

	// Well-known metadata endpoint
	if path == f.state.metaPath {
		logger.Debugf("[dubbo-go-pixiu] mcp auth filter meta path: %s", path)
		hc.SendLocalReply(http.StatusOK, f.state.metaBody)
		return filter.Stop
	}

	// Resolve rule by framework route entry's cluster
	var rule *Rule
	if rEntry := hc.GetRouteEntry(); rEntry != nil {
		for i := range f.state.rules {
			if rEntry.Cluster == f.state.rules[i].Cluster {
				rule = &f.state.rules[i]
				break
			}
		}
	}
	if rule == nil {
		return filter.Continue
	}

	// Extract bearer token
	token := extractBearer(hc.GetHeader(constant.Authorization))
	if token == "" {
		f.unauthorized(hc, "invalid_token", "missing bearer token")
		return filter.Stop
	}

	// Determine provider by token issuer (do not trust token issuer blindly)
	providerName, err := f.state.validator.ProviderByTokenIssuer(token)
	if err != nil {
		logger.Warnf("[dubbo-go-pixiu] provider lookup by token issuer failed: %v", err)
		f.unauthorized(hc, "invalid_token", "untrusted token issuer")
		return filter.Stop
	}

	// Validate token using provider derived from issuer
	_, err = f.state.validator.Validate(providerName, token)
	if err != nil {
		// Map validator.ValidationError if possible
		verr := validator.ValidationError{}
		code := "invalid_token"
		msg := "invalid token"
		if ok := asValidationError(err, &verr); ok {
			if verr.Code != "" {
				code = verr.Code
			}
			if verr.Message != "" {
				msg = verr.Message
			}
		} else {
			msg = err.Error()
		}
		f.unauthorized(hc, code, msg)
		return filter.Stop
	}

	// remove Authorization header to avoid leaking token to downstream services
	hc.Request.Header.Del(constant.Authorization)

	return filter.Continue
}

// unauthorized writes 401 with WWW-Authenticate including resource metadata URL
func (f *Filter) unauthorized(hc *contexthttp.HttpContext, code, desc string) {
	// Build absolute metadata URL from request
	scheme := "http"
	if hc.Request.TLS != nil {
		scheme = "https"
	}
	metaURL := scheme + constant.ProtocolSlash + hc.Request.Host + f.state.metaPath
	// Per RFC9728, include resource_metadata parameter; include OAuth error fields
	header := fmt.Sprintf("Bearer resource_metadata=\"%s\", error=\"%s\", error_description=\"%s\"", metaURL, escapeParam(code), escapeParam(desc))
	hc.AddHeader(constant.WWWAuthenticate, header)
	writeOAuthError(hc, http.StatusUnauthorized, code, desc)
}

// writeOAuthError responds with JSON {error, error_description}
func writeOAuthError(hc *contexthttp.HttpContext, status int, code, desc string) {
	resp := map[string]string{
		"error":             code,
		"error_description": desc,
	}
	b, _ := json.Marshal(resp)
	hc.SendLocalReply(status, b)
}

// extractBearer pulls the token from Authorization header value
func extractBearer(v string) string {
	if v == "" {
		return ""
	}
	if len(v) < 7 {
		return ""
	}
	if strings.HasPrefix(strings.ToLower(v), "bearer ") {
		return strings.TrimSpace(v[7:])
	}
	return ""
}

// asValidationError performs a typed unwrap without importing errors in callers
func asValidationError(err error, target *validator.ValidationError) bool {
	// local re-implementation to avoid importing errors here again
	// keep it straightforward using standard errors.As
	return errorsAs(err, target)
}

// errorsAs isolates usage to enable unit testing in this file easily
var errorsAs = func(err error, target any) bool { return errors.As(err, target) }

// escapeParam performs minimal escaping suitable for WWW-Authenticate param values
func escapeParam(s string) string {
	// Replace embedded quotes and backslashes per RFC6750 guidance
	s = strings.ReplaceAll(s, "\\", "\\\\")
	s = strings.ReplaceAll(s, "\"", "\\\"")
	return s
}

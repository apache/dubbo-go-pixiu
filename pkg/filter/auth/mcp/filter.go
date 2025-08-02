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
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"strings"

	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"

	pixiu_http "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/auth/mcp/internal/validator"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/lestrrat-go/jwx/v2/jwt"
)

const (
	// Kind is the filter's kind.
	Kind = constant.HTTPMCPAuthFilter
	// WellKnownOAuthPath is the path for OAuth protected resource metadata.
	WellKnownOAuthPath = "/.well-known/oauth-protected-resource"
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is the HTTP filter plugin.
	Plugin struct{}

	// FilterFactory is the HTTP filter factory.
	FilterFactory struct {
		cfg       *Config
		validator *validator.Validator
	}

	// Filter is the HTTP filter instance.
	Filter struct {
		cfg       *Config
		validator *validator.Validator
	}

	// ErrorResponse defines the standard error response body.
	ErrorResponse struct {
		Error            string `json:"error"`
		ErrorDescription string `json:"error_description,omitempty"`
	}
)

// Kind returns the filter's kind.
func (p *Plugin) Kind() string {
	return Kind
}

// CreateFilterFactory creates a new FilterFactory.
func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

// PrepareFilterChain adds the filter to the filter chain.
func (ff *FilterFactory) PrepareFilterChain(ctx *pixiu_http.HttpContext, chain filter.FilterChain) error {
	f := &Filter{cfg: ff.cfg, validator: ff.validator}
	chain.AppendDecodeFilters(f)
	return nil
}

// Apply initializes the validator from the configuration.
func (ff *FilterFactory) Apply() error {
	v, err := validator.NewValidator(context.Background(), ff.cfg.Providers)
	if err != nil {
		return fmt.Errorf("failed to create MCP auth validator: %w", err)
	}
	ff.validator = v
	return nil
}

// Config returns the filter's configuration.
func (ff *FilterFactory) Config() interface{} {
	return ff.cfg
}

// Decode is the main logic of the filter.
func (f *Filter) Decode(ctx *pixiu_http.HttpContext) filter.FilterStatus {
	path := ctx.Request.URL.Path

	if f.cfg.ResourceMetadata.Enabled && path == WellKnownOAuthPath {
		f.handleWellKnownRequest(ctx)
		return filter.Stop
	}

	rule := f.findMatchingRule(path)
	if rule == nil {
		return filter.Continue // No rule for this path, pass through.
	}

	tokenStr := f.extractBearerToken(ctx)
	if tokenStr == "" {
		f.handleMissingToken(ctx)
		return filter.Stop
	}

	validatedToken, err := f.validator.Validate(ctx.Request.Context(), rule.ProviderName, tokenStr)
	if err != nil {
		logger.Warnf("Token validation failed: %v", err)
		f.sendError(ctx, http.StatusUnauthorized, "invalid_token", err.Error())
		return filter.Stop
	}

	if !f.checkScopes(validatedToken, rule.RequiredScopes) {
		f.sendError(ctx, http.StatusForbidden, "insufficient_scope", "Token has insufficient scope")
		return filter.Stop
	}

	return filter.Continue
}

func (f *Filter) findMatchingRule(path string) *Rule {
	for i := range f.cfg.Rules {
		rule := &f.cfg.Rules[i]
		if strings.HasPrefix(path, rule.Match.Prefix) {
			return rule
		}
	}
	return nil
}

func (f *Filter) extractBearerToken(ctx *pixiu_http.HttpContext) string {
	authHeader := ctx.Request.Header.Get("Authorization")
	if !strings.HasPrefix(authHeader, "Bearer ") {
		return ""
	}
	return strings.TrimPrefix(authHeader, "Bearer ")
}

func (f *Filter) handleWellKnownRequest(ctx *pixiu_http.HttpContext) {
	metadata := map[string]interface{}{
		"authorization_servers": f.cfg.ResourceMetadata.AuthorizationServers,
	}
	body, _ := json.Marshal(metadata)
	ctx.SendLocalReply(http.StatusOK, body)
}

func (f *Filter) handleMissingToken(ctx *pixiu_http.HttpContext) {
	realm := f.cfg.ResourceMetadata.ServerHost
	uri := realm + WellKnownOAuthPath
	authHeader := fmt.Sprintf(`Bearer realm="%s", resource_metadata_uri="%s"`, realm, uri)
	ctx.AddHeader("WWW-Authenticate", authHeader)
	f.sendError(ctx, http.StatusUnauthorized, "unauthorized", "Authorization header is missing or invalid")
}

func (f *Filter) checkScopes(token jwt.Token, requiredScopes []string) bool {
	if len(requiredScopes) == 0 {
		return true // No scopes required.
	}

	scopeClaim, ok := token.Get("scope")
	if !ok {
		return false // Scope claim is missing.
	}

	scopeStr, ok := scopeClaim.(string)
	if !ok {
		return false // Scope claim is not a string.
	}

	providedScopes := make(map[string]struct{})
	for _, s := range strings.Split(scopeStr, " ") {
		providedScopes[s] = struct{}{}
	}

	for _, required := range requiredScopes {
		if _, ok := providedScopes[required]; !ok {
			return false // A required scope is missing.
		}
	}

	return true
}

func (f *Filter) sendError(ctx *pixiu_http.HttpContext, statusCode int, err, desc string) {
	body, _ := json.Marshal(ErrorResponse{
		Error:            err,
		ErrorDescription: desc,
	})
	ctx.SendLocalReply(statusCode, body)
}

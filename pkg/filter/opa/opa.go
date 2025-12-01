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
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"time"
)

import (
	"github.com/open-policy-agent/opa/rego"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contextHttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	Kind = constant.HTTPAuthOPAFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	Plugin struct{}

	FilterFactory struct {
		cfg           *Config
		rego          *rego.Rego              // For embedded mode (deprecated, kept for backward compat)
		preparedQuery *rego.PreparedEvalQuery // Pre-compiled query for embedded mode
		httpClient    *http.Client            // For server mode
	}

	Filter struct {
		cfg           *Config
		preparedQuery *rego.PreparedEvalQuery // For embedded mode
		httpClient    *http.Client            // For server mode
	}

	Config struct {
		// Server mode configuration (recommended for production)
		ServerURL    string `yaml:"server_url" json:"server_url" mapstructure:"server_url"`          // OPA Server address, e.g., http://opa-server:8181
		DecisionPath string `yaml:"decision_path" json:"decision_path" mapstructure:"decision_path"` // Decision path, e.g., /v1/data/http/authz/allow
		TimeoutMs    int    `yaml:"timeout_ms" json:"timeout_ms" mapstructure:"timeout_ms"`          // Request timeout in milliseconds, default 100
		BearerToken  string `yaml:"bearer_token" json:"bearer_token" mapstructure:"bearer_token"`    // Optional authentication token

		// Embedded mode configuration (for backward compatibility)
		Policy     string `yaml:"policy" json:"policy" mapstructure:"policy"`             // Policy content
		Entrypoint string `yaml:"entrypoint" json:"entrypoint" mapstructure:"entrypoint"` // Policy entrypoint
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

func (factory *FilterFactory) Config() any {
	return factory.cfg
}

// Apply is called after the configuration is loaded and is used to prepare the OPA query.
func (factory *FilterFactory) Apply() error {
	cfg := factory.cfg

	// Server mode (recommended for production)
	if cfg.ServerURL != "" {
		if cfg.DecisionPath == "" {
			return fmt.Errorf("decision_path is required when using OPA server mode")
		}

		timeout := 100 * time.Millisecond
		if cfg.TimeoutMs > 0 {
			timeout = time.Duration(cfg.TimeoutMs) * time.Millisecond
		}

		factory.httpClient = &http.Client{
			Timeout: timeout,
			Transport: &http.Transport{
				MaxIdleConns:        100,
				MaxIdleConnsPerHost: 10,
				IdleConnTimeout:     90 * time.Second,
			},
		}

		logger.Infof("OPA filter initialized in server mode: %s%s (timeout: %v)", cfg.ServerURL, cfg.DecisionPath, timeout)
		return nil
	}

	// Embedded mode (for backward compatibility)
	if cfg.Policy != "" {
		logger.Warnf("OPA filter using embedded mode. Consider using server mode for better maintainability and performance.")

		if cfg.Entrypoint == "" {
			return fmt.Errorf("entrypoint is required when using embedded mode")
		}

		r := rego.New(
			rego.Query(cfg.Entrypoint),
			rego.Module("policy.rego", cfg.Policy),
		)

		// Pre-compile the query once at initialization for better performance
		preparedQuery, err := r.PrepareForEval(context.Background())
		if err != nil {
			return fmt.Errorf("failed to prepare OPA query: %w", err)
		}

		factory.rego = r
		factory.preparedQuery = &preparedQuery

		logger.Infof("OPA filter initialized in embedded mode")
		return nil
	}

	return fmt.Errorf("OPA filter requires either server_url (recommended) or policy configuration")
}

// PrepareFilterChain prepares the filter chain for a new request by dynamically creating a Filter
func (factory *FilterFactory) PrepareFilterChain(ctx *contextHttp.HttpContext, chain filter.FilterChain) error {
	// Shallow copy cfg (copy the struct value; inner reference fields remain shared)
	cfgCopy := *factory.cfg
	var f *Filter

	// Server mode (priority)
	if factory.httpClient != nil {
		f = &Filter{
			cfg:        &cfgCopy,
			httpClient: factory.httpClient,
		}
	} else if factory.preparedQuery != nil {
		// Embedded mode (backward compatibility) - reuse pre-compiled query
		f = &Filter{
			cfg:           &cfgCopy,
			preparedQuery: factory.preparedQuery,
		}
	} else {
		return fmt.Errorf("OPA filter not properly initialized")
	}

	chain.AppendDecodeFilters(f)
	return nil
}

// Decode is the core logic of the filter. It converts HTTP request data into a standard OPA input format and evaluates the policy.
func (f *Filter) Decode(c *contextHttp.HttpContext) filter.FilterStatus {
	input := map[string]any{
		"method":      c.Request.Method,
		"path":        c.Request.URL.Path,
		"headers":     c.Request.Header,
		"client_ip":   c.GetClientIP(),
		"query":       c.Request.URL.Query(),
		"host":        c.Request.Host,
		"remote_addr": c.Request.RemoteAddr,
		"user_agent":  c.Request.UserAgent(),
		"route":       c.GetRouteEntry(),
		"api":         c.GetAPI(),
		"params":      c.Params,
	}

	// Server mode (priority)
	if f.httpClient != nil {
		return f.evaluateServer(c, input)
	}

	// Embedded mode (backward compatibility)
	if f.preparedQuery != nil {
		return f.evaluateEmbedded(c, input)
	}

	logger.Error("OPA filter not initialized properly")
	errResp := contextHttp.InternalError.WithError(fmt.Errorf("OPA filter not initialized"))
	c.SendLocalReply(errResp.Status, errResp.ToJSON())
	return filter.Stop
}

// evaluateEmbedded evaluates the policy using embedded OPA engine
func (f *Filter) evaluateEmbedded(c *contextHttp.HttpContext, input map[string]any) filter.FilterStatus {
	results, err := f.preparedQuery.Eval(c.Ctx, rego.EvalInput(input))
	if err != nil {
		logger.Errorf("OPA embedded evaluation error: %v", err)
		errResp := contextHttp.InternalError.WithError(err)
		c.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	if len(results) == 0 {
		logger.Debugf("OPA embedded policy returned empty result for request: %s %s", input["method"], input["path"])
		errResp := contextHttp.Forbidden.New()
		c.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	// Extract decision from result (supports both boolean and object formats)
	allow := extractDecision(results[0].Expressions[0].Value)
	if !allow {
		logger.Debugf("OPA embedded policy denied request: %s %s", input["method"], input["path"])
		errResp := contextHttp.Forbidden.New()
		c.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	return filter.Continue
}

// evaluateServer evaluates the policy by calling OPA server REST API
func (f *Filter) evaluateServer(c *contextHttp.HttpContext, input map[string]any) filter.FilterStatus {
	// Construct request body according to OPA REST API specification
	requestBody := map[string]any{
		"input": input,
	}

	jsonData, err := json.Marshal(requestBody)
	if err != nil {
		logger.Errorf("Failed to marshal OPA request: %v", err)
		errResp := contextHttp.InternalError.WithError(err)
		c.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	// Construct HTTP request
	url := f.cfg.ServerURL + f.cfg.DecisionPath
	req, err := http.NewRequestWithContext(c.Ctx, "POST", url, bytes.NewBuffer(jsonData))
	if err != nil {
		logger.Errorf("Failed to create OPA request: %v", err)
		errResp := contextHttp.InternalError.WithError(err)
		c.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	req.Header.Set("Content-Type", "application/json")
	if f.cfg.BearerToken != "" {
		req.Header.Set("Authorization", "Bearer "+f.cfg.BearerToken)
	}

	// Send request to OPA server
	resp, err := f.httpClient.Do(req)
	if err != nil {
		logger.Errorf("OPA server request failed: %v", err)
		// Check if it's a timeout error using net.Error interface
		if netErr, ok := err.(interface{ Timeout() bool }); ok && netErr.Timeout() {
			errResp := contextHttp.GatewayTimeout.WithError(err)
			c.SendLocalReply(errResp.Status, errResp.ToJSON())
		} else {
			errResp := contextHttp.ServiceUnavailable.WithError(err)
			c.SendLocalReply(errResp.Status, errResp.ToJSON())
		}
		return filter.Stop
	}
	defer resp.Body.Close()

	// Check HTTP status code
	if resp.StatusCode != 200 {
		body, err := io.ReadAll(resp.Body)
		if err != nil {
			logger.Errorf("Failed to read OPA server response body: %v", err)
			body = []byte("")
		}
		logger.Errorf("OPA server returned status %d: %s", resp.StatusCode, string(body))
		errResp := contextHttp.BadGateway.WithError(fmt.Errorf("OPA server returned status %d", resp.StatusCode))
		c.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	// Parse response according to OPA REST API specification
	var result map[string]any
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		logger.Errorf("Failed to decode OPA response: %v", err)
		errResp := contextHttp.BadGateway.WithError(err)
		c.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	// Extract decision from result (supports both boolean and object formats)
	if resultValue, exists := result["result"]; exists {
		allow := extractDecision(resultValue)
		if !allow {
			logger.Debugf("OPA server policy denied request: %s %s", input["method"], input["path"])
			errResp := contextHttp.Forbidden.New()
			c.SendLocalReply(errResp.Status, errResp.ToJSON())
			return filter.Stop
		}
		return filter.Continue
	}

	logger.Errorf("Invalid OPA response format: missing 'result' field")
	errResp := contextHttp.BadGateway.WithError(fmt.Errorf("missing 'result' field in OPA response"))
	c.SendLocalReply(errResp.Status, errResp.ToJSON())
	return filter.Stop
}

// extractDecision extracts the allow decision from OPA result
// Supports multiple formats:
// 1. Boolean: true/false
// 2. Object with "allow" field: {allow: true}
// 3. Object with "result" field: {result: true}
func extractDecision(value any) bool {
	// Format 1: Direct boolean value
	if boolVal, ok := value.(bool); ok {
		return boolVal
	}

	// Format 2 & 3: Object with allow/result field
	if objVal, ok := value.(map[string]any); ok {
		// Try "allow" field first (common OPA pattern)
		if allowVal, exists := objVal["allow"]; exists {
			if boolVal, ok := allowVal.(bool); ok {
				return boolVal
			}
		}
		// Try "result" field as fallback
		if resultVal, exists := objVal["result"]; exists {
			if boolVal, ok := resultVal.(bool); ok {
				return boolVal
			}
		}
	}

	// Default to deny if format is unrecognized
	logger.Warnf("Unrecognized OPA decision format, defaulting to deny. Value: %v", value)
	return false
}

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
	"fmt"
)

import (
	"github.com/open-policy-agent/opa/rego"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
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
		cfg  *Config
		rego *rego.Rego
	}

	Filter struct {
		cfg           *Config
		preparedQuery *rego.PreparedEvalQuery
	}

	Config struct {
		Policy     string `yaml:"policy" json:"policy" `
		Entrypoint string `yaml:"entrypoint" json:"entrypoint" `
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
	policy := factory.cfg.Policy
	if policy == "" {
		return fmt.Errorf("OPA policy is empty in the configuration")
	}

	r := rego.New(
		rego.Query(factory.cfg.Entrypoint),
		rego.Module("policy.rego", policy),
	)

	factory.rego = r

	return nil
}

// PrepareFilterChain prepares the filter chain for a new request by dynamically creating a Filter
func (factory *FilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
	if factory.rego == nil {
		return fmt.Errorf("rego instance not initialized in factory")
	}

	preparedQuery, err := factory.rego.PrepareForEval(ctx.Ctx)
	if err != nil {
		return fmt.Errorf("failed to prepare OPA query: %w", err)
	}

	// Make a shallow copy of the factory config to avoid sharing the factory's pointer.
	cfgCopy := *factory.cfg
	f := &Filter{cfg: &cfgCopy, preparedQuery: &preparedQuery}
	chain.AppendDecodeFilters(f)
	return nil
}

// Decode is the core logic of the filter. It converts HTTP request data into a standard OPA input format and evaluates the policy.
func (f *Filter) Decode(c *http.HttpContext) filter.FilterStatus {
	if f.preparedQuery == nil {
		logger.Error("OPA filter not initialized properly.")
		return filter.Stop
	}

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

	results, err := f.preparedQuery.Eval(c.Ctx, rego.EvalInput(input))
	if err != nil {
		logger.Error("OPA evaluation error: %v\n", err)
		return filter.Stop
	}

	if len(results) == 0 || results[0].Expressions[0].Value != true {
		return filter.Stop
	}

	return filter.Continue
}

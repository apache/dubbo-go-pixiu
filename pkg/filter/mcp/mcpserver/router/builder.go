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

package router

import (
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"math"
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// Build constructs a ToolSelector from the router configuration and a shared
// session plan store.
//
// It returns (nil, nil) when routing is disabled or unconfigured, signalling
// the MCP server filter to keep its passthrough behavior with zero overhead.
// Configuration errors (invalid regex, unknown default bundle) fail fast so a
// broken policy chain never silently degrades to the wrong default.
func Build(cfg *model.RouterConfig, store *SessionPlanStore) (ToolSelector, error) {
	if cfg == nil || !cfg.Enabled {
		return nil, nil
	}
	if store == nil {
		return nil, fmt.Errorf("router enabled but session plan store is nil")
	}
	if err := validateFallback(cfg.Fallback); err != nil {
		return nil, err
	}
	if err := validateSampleRate(cfg.Audit.SampleRate); err != nil {
		return nil, err
	}
	if err := validatePolicyRisks(cfg.Policy); err != nil {
		return nil, err
	}

	// Register Prometheus collectors on first enabled build.
	initMetrics()

	opts := CompositeOptions{
		Store:         store,
		Log:           NewDecisionLogger(cfg.Audit.SampleRate, cfg.Audit.PayloadLogging),
		Fallback:      cfg.Fallback,
		DefaultBundle: cfg.DefaultBundle,
		EnforceOnCall: enforceOnCall(cfg),
		ConfigHash:    configHash(cfg),
	}

	// Workflow selector is needed both as a pipeline stage and to resolve named
	// bundles for fallback/progressive, so build it up front.
	var wf *WorkflowSelector
	if len(cfg.Workflows) > 0 {
		var err error
		wf, err = NewWorkflowSelector(cfg.Workflows)
		if err != nil {
			return nil, err
		}
		opts.Bundles = wf
	}

	if stageEnabled(cfg.Stages.Policy, true) && len(cfg.Policy.Rules) > 0 {
		pf, err := NewPolicyFilter(cfg.Policy)
		if err != nil {
			return nil, err
		}
		opts.Policy = pf
	}

	if stageEnabled(cfg.Stages.Workflow, true) {
		opts.Workflow = wf
	}

	if cfg.Stages.Progressive {
		initialBundle := strings.TrimSpace(cfg.Progressive.InitialBundle)
		if initialBundle == "" {
			return nil, fmt.Errorf("router progressive.initial_bundle is required when progressive stage is enabled")
		}
		progressiveCfg := cfg.Progressive
		progressiveCfg.InitialBundle = initialBundle
		opts.Progressive = NewProgressiveGate(progressiveCfg, wf)
	}

	// Validate that a bundle_default fallback references a real non-empty workflow
	// bundle. This avoids silently producing an empty "fallback" plan.
	if opts.Fallback == "" || opts.Fallback == FallbackBundleDefault {
		defaultBundle := strings.TrimSpace(cfg.DefaultBundle)
		if defaultBundle == "" {
			return nil, fmt.Errorf("router default_bundle is required when fallback is %q", FallbackBundleDefault)
		}
		if wf == nil {
			return nil, fmt.Errorf("router default_bundle %q set but no workflows defined", defaultBundle)
		}
		bundle, ok := wf.bundleTools(defaultBundle)
		if !ok {
			return nil, fmt.Errorf("router default_bundle %q does not match any workflow", defaultBundle)
		}
		if len(bundle) == 0 {
			return nil, fmt.Errorf("router default_bundle %q references an empty workflow", defaultBundle)
		}
	}

	// Validate that progressive.initial_bundle references a real workflow bundle.
	// A misconfigured bundle would otherwise fail open (reveal all tools) at
	// runtime, defeating the purpose of limiting the initial exposure surface.
	if cfg.Stages.Progressive {
		if wf == nil {
			return nil, fmt.Errorf("router progressive.initial_bundle %q set but no workflows defined", cfg.Progressive.InitialBundle)
		}
		if _, ok := wf.bundleTools(cfg.Progressive.InitialBundle); !ok {
			return nil, fmt.Errorf("router progressive.initial_bundle %q does not match any workflow", cfg.Progressive.InitialBundle)
		}
	}

	logger.Infof("[dubbo-go-pixiu] mcp tool router enabled (policy=%v workflow=%v progressive=%v enforce_on_call=%v fallback=%s)",
		opts.Policy != nil, opts.Workflow != nil, opts.Progressive != nil, opts.EnforceOnCall, firstNonEmpty(opts.Fallback, FallbackBundleDefault))

	return NewCompositeSelector(opts), nil
}

// enforceOnCall resolves the enforce_on_call setting, defaulting to true.
func enforceOnCall(cfg *model.RouterConfig) bool {
	if cfg.EnforceOnCall == nil {
		return true
	}
	return *cfg.EnforceOnCall
}

// stageEnabled resolves a *bool stage toggle with a default.
func stageEnabled(v *bool, def bool) bool {
	if v == nil {
		return def
	}
	return *v
}

func firstNonEmpty(a, b string) string {
	if a != "" {
		return a
	}
	return b
}

func validateFallback(fallback string) error {
	switch fallback {
	case "", FallbackBundleDefault, FallbackFailClosed:
		return nil
	default:
		return fmt.Errorf("router fallback %q is not supported", fallback)
	}
}

func validateSampleRate(rate float64) error {
	if math.IsNaN(rate) || rate < 0 || rate > 1 {
		return fmt.Errorf("router audit.sample_rate must be between 0 and 1")
	}
	return nil
}

// configHash produces a stable hash of the router config so plans recompute
// when the config changes (e.g. a dynamic update). It is best-effort: a
// marshalling failure falls back to an empty hash, which simply disables reuse.
func configHash(cfg *model.RouterConfig) string {
	data, err := json.Marshal(cfg)
	if err != nil {
		return ""
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:8])
}

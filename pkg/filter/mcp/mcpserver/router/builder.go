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
	"sync/atomic"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

var routerInstanceCounter uint64

// Build constructs a ToolSelector from the normalized router configuration and
// an instance-owned session plan store.
//
// Configuration errors (invalid regex, unknown default bundle) fail fast so a
// broken policy chain never silently degrades to the wrong default.
func Build(cfg *model.RouterConfig, store *SessionPlanStore) (ToolSelector, error) {
	return BuildWithOptions(cfg, store, BuildOptions{})
}

type BuildOptions struct {
	VisibleFP VisibleFingerprintFunc
}

func BuildWithOptions(cfg *model.RouterConfig, store *SessionPlanStore, buildOpts BuildOptions) (ToolSelector, error) {
	cfg = normalizeConfig(cfg)
	if store == nil {
		return nil, fmt.Errorf("router session plan store is nil")
	}
	if err := validateFallback(cfg.Fallback); err != nil {
		return nil, err
	}
	if err := validateSampleRate(cfg.Audit.SampleRate); err != nil {
		return nil, err
	}
	if err := validateSessionConfig(cfg.Session); err != nil {
		return nil, err
	}
	if err := validatePolicyRisks(cfg.Policy); err != nil {
		return nil, err
	}

	// Register Prometheus collectors on first governance build.
	initMetrics()

	cfgHash, err := configHash(cfg)
	if err != nil {
		return nil, err
	}
	opts := CompositeOptions{
		Store:         store,
		Log:           NewDecisionLogger(cfg.Audit.SampleRate, cfg.Audit.DecisionDetailLogging),
		Fallback:      cfg.Fallback,
		DefaultBundle: cfg.DefaultBundle,
		ConfigHash:    cfgHash,
		RouterID:      newRouterInstanceID(cfgHash),
		VisibleFP:     buildOpts.VisibleFP,
	}

	wf, err := buildWorkflowSelector(cfg, &opts)
	if err != nil {
		return nil, err
	}
	if err := buildPolicyFilter(cfg, &opts); err != nil {
		return nil, err
	}
	if err := buildProgressiveGate(cfg, wf, &opts); err != nil {
		return nil, err
	}
	if err := validateDefaultFallback(cfg, wf, opts.Fallback); err != nil {
		return nil, err
	}

	logger.Infof("[dubbo-go-pixiu] mcp tool governance active (policy=%v workflow=%v progressive=%v fallback=%s)",
		opts.Policy != nil, opts.Workflow != nil, opts.Progressive != nil, opts.Fallback)

	return NewCompositeSelector(opts)
}

func normalizeConfig(cfg *model.RouterConfig) *model.RouterConfig {
	if cfg == nil {
		return &model.RouterConfig{Fallback: DefaultFallback}
	}
	cp := cfg.DeepCopy()
	if cp.Fallback == "" {
		cp.Fallback = DefaultFallback
	}
	return cp
}

// buildWorkflowSelector creates the shared workflow selector used by workflow,
// fallback, and progressive stages.
func buildWorkflowSelector(cfg *model.RouterConfig, opts *CompositeOptions) (*WorkflowSelector, error) {
	if len(cfg.Workflows) == 0 {
		return nil, nil
	}
	wf, err := NewWorkflowSelector(cfg.Workflows)
	if err != nil {
		return nil, err
	}
	opts.Bundles = wf
	if stageEnabled(cfg.Stages.Workflow, true) {
		opts.Workflow = wf
	}
	return wf, nil
}

func buildPolicyFilter(cfg *model.RouterConfig, opts *CompositeOptions) error {
	if !stageEnabled(cfg.Stages.Policy, true) || len(cfg.Policy.Rules) == 0 {
		return nil
	}
	pf, err := NewPolicyFilter(cfg.Policy)
	if err != nil {
		return err
	}
	opts.Policy = pf
	return nil
}

func buildProgressiveGate(cfg *model.RouterConfig, wf *WorkflowSelector, opts *CompositeOptions) error {
	if !cfg.Stages.Progressive {
		return nil
	}
	if cfg.Progressive.ExpandAfterCalls <= 0 {
		return fmt.Errorf("router progressive.expand_after_calls must be greater than 0")
	}
	initialBundle := strings.TrimSpace(cfg.Progressive.InitialBundle)
	if err := validateProgressiveBundle(initialBundle, wf); err != nil {
		return err
	}
	progressiveCfg := cfg.Progressive
	progressiveCfg.InitialBundle = initialBundle
	opts.Progressive = NewProgressiveGate(progressiveCfg, wf)
	return nil
}

// validateDefaultFallback verifies bundle_default references a real non-empty
// workflow bundle, avoiding an empty fallback plan at runtime.
func validateDefaultFallback(cfg *model.RouterConfig, wf *WorkflowSelector, fallback string) error {
	if fallback != FallbackBundleDefault {
		return nil
	}
	defaultBundle := strings.TrimSpace(cfg.DefaultBundle)
	if defaultBundle == "" {
		return fmt.Errorf("router default_bundle is required when fallback is %q", FallbackBundleDefault)
	}
	if wf == nil {
		return fmt.Errorf("router default_bundle %q set but no workflows defined", defaultBundle)
	}
	bundle, ok := wf.bundleTools(defaultBundle)
	if !ok {
		return fmt.Errorf("router default_bundle %q does not match any workflow", defaultBundle)
	}
	if len(bundle) == 0 {
		return fmt.Errorf("router default_bundle %q references an empty workflow", defaultBundle)
	}
	return nil
}

// validateProgressiveBundle verifies the initial bundle before progressive
// disclosure starts, so a typo cannot reveal all tools by accident.
func validateProgressiveBundle(initialBundle string, wf *WorkflowSelector) error {
	if initialBundle == "" {
		return fmt.Errorf("router progressive.initial_bundle is required when progressive disclosure is active")
	}
	if wf == nil {
		return fmt.Errorf("router progressive.initial_bundle %q set but no workflows defined", initialBundle)
	}
	if _, ok := wf.bundleTools(initialBundle); !ok {
		return fmt.Errorf("router progressive.initial_bundle %q does not match any workflow", initialBundle)
	}
	return nil
}

// stageEnabled resolves a *bool stage toggle with a default.
func stageEnabled(v *bool, def bool) bool {
	if v == nil {
		return def
	}
	return *v
}

func newRouterInstanceID(cfgHash string) string {
	n := atomic.AddUint64(&routerInstanceCounter, 1)
	if cfgHash == "" {
		return fmt.Sprintf("router-%d", n)
	}
	return fmt.Sprintf("router-%d-%s", n, cfgHash)
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
	if math.IsNaN(rate) || math.IsInf(rate, 0) || rate < 0 || rate > 1 {
		return fmt.Errorf("router audit.sample_rate must be between 0 and 1")
	}
	return nil
}

func validateSessionConfig(cfg model.RouterSessionConfig) error {
	if cfg.MaxEntries < 0 {
		return fmt.Errorf("router session.max_entries must be greater than or equal to 0")
	}
	return nil
}

// configHash produces a stable hash of the normalized authorization-affecting
// router config so plans recompute when governance behavior changes.
func configHash(cfg *model.RouterConfig) (string, error) {
	authCfg := cfg.DeepCopy()
	authCfg.Audit = model.AuditConfig{}
	authCfg.Session = model.RouterSessionConfig{}
	data, err := json.Marshal(authCfg)
	if err != nil {
		logger.Warnf("[dubbo-go-pixiu] mcp router config hash failed")
		return "", fmt.Errorf("router config hash failed: %w", err)
	}
	sum := sha256.Sum256(data)
	return hex.EncodeToString(sum[:8]), nil
}

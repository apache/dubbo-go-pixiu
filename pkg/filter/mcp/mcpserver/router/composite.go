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
	"context"
	"encoding/json"
	"fmt"
	"hash/fnv"
	"io"
	"sort"
	"strconv"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// Fallback strategies.
const (
	DefaultFallback       = FallbackFailClosed
	FallbackBundleDefault = "bundle_default"
	FallbackFailClosed    = "fail_closed"
)

// CompositeSelector runs the deterministic selection pipeline:
//
//	policy -> workflow -> progressive
//
// Each stage is optional; a nil stage is skipped. Results are cached per session
// in the SessionPlanStore and reused while the metadata/config version is stable.
type CompositeSelector struct {
	policy      *PolicyFilter
	workflow    *WorkflowSelector
	progressive *ProgressiveGate
	bundles     *WorkflowSelector

	store *SessionPlanStore
	log   *DecisionLogger

	routerID        string
	fallback        string
	defaultBundle   string
	configHash      string
	progressiveHash string
	visibleFP       VisibleFingerprintFunc
}

// CompositeOptions configures a CompositeSelector. The builder populates it.
type CompositeOptions struct {
	Policy        *PolicyFilter
	Workflow      *WorkflowSelector
	Progressive   *ProgressiveGate
	Bundles       *WorkflowSelector
	Store         *SessionPlanStore
	Log           *DecisionLogger
	Fallback      string
	DefaultBundle string
	ConfigHash    string
	RouterID      string
	VisibleFP     VisibleFingerprintFunc
}

// NewCompositeSelector assembles the pipeline from validated options.
func NewCompositeSelector(opts CompositeOptions) (*CompositeSelector, error) {
	if opts.Store == nil {
		return nil, fmt.Errorf("router session plan store is nil")
	}
	fallback := opts.Fallback
	if fallback == "" {
		fallback = DefaultFallback
	}
	if err := validateFallback(fallback); err != nil {
		return nil, err
	}
	log := opts.Log
	if log == nil {
		log = NewDecisionLogger(0, false)
	}
	visibleFP := opts.VisibleFP
	if visibleFP == nil {
		visibleFP = DefaultVisibleFingerprint
	}
	return &CompositeSelector{
		policy:          opts.Policy,
		workflow:        opts.Workflow,
		progressive:     opts.Progressive,
		bundles:         opts.Bundles,
		store:           opts.Store,
		log:             log,
		routerID:        opts.RouterID,
		fallback:        fallback,
		defaultBundle:   opts.DefaultBundle,
		configHash:      opts.ConfigHash,
		progressiveHash: progressiveConfigHash(opts.Progressive),
		visibleFP:       visibleFP,
	}, nil
}

type selectionPipeline struct {
	tools         []model.ToolConfig
	policyAllowed []model.ToolConfig
	traces        []DecisionTrace
	stageCounts   map[string]StageCount
	outcome       string
	expanded      bool
}

func newSelectionPipeline(candidates []model.ToolConfig) *selectionPipeline {
	return &selectionPipeline{
		tools:         candidates,
		policyAllowed: candidates,
		stageCounts:   make(map[string]StageCount),
		outcome:       SelectionOutcomeSelected,
	}
}

// Select runs the pipeline, reusing a cached plan when the version is unchanged.
func (c *CompositeSelector) Select(_ context.Context, sc SelectionContext, candidates []model.ToolConfig) (*SelectionPlan, error) {
	start := time.Now()
	identityHash, err := identityFingerprint(sc)
	if err != nil {
		return nil, err
	}
	version := c.version(candidates, identityHash, sc.SessionID, sc.CatalogVersion)
	key := c.planKey(sc.SessionID)

	if cached, ok := c.cachedSelectionPlan(key, version, candidates, start); ok {
		return cached, nil
	}

	plan := c.computeSelectionPlan(key, identityHash, sc, candidates, version)
	plan, err = c.storeSelectionPlan(key, plan, sc)
	if err != nil {
		recordFallback("plan_persistence_failed")
		return nil, err
	}

	result := "ok"
	if outcomeAllowsFallback(plan.Outcome) {
		result = "fallback"
		recordFallback(plan.Outcome)
	}

	c.recordSelectionResult(result, plan, candidates, start)
	c.log.Log(sc, plan, len(candidates))

	return plan, nil
}

func (c *CompositeSelector) cachedSelectionPlan(key PlanKey, version string, candidates []model.ToolConfig, start time.Time) (*SelectionPlan, bool) {
	cached, ok := c.store.Get(key)
	if !ok || cached.Version != version {
		return nil, false
	}
	elapsedMS := float64(time.Since(start).Microseconds()) / 1000.0
	recordSelection("cached", cached.Mode, len(candidates), len(cached.ToolNames), elapsedMS)
	return cached, true
}

func (c *CompositeSelector) computeSelectionPlan(key PlanKey, identityHash string, sc SelectionContext, candidates []model.ToolConfig, version string) *SelectionPlan {
	pipeline := c.runSelectionPipeline(key, identityHash, sc, candidates)
	if outcomeAllowsFallback(pipeline.outcome) {
		return c.applyFallback(sc, pipeline.policyAllowed, version, pipeline.traces, pipeline.stageCounts, pipeline.outcome)
	}
	return c.newSelectionPlan(sc, candidates, version, identityHash, pipeline)
}

func (c *CompositeSelector) runSelectionPipeline(key PlanKey, identityHash string, sc SelectionContext, candidates []model.ToolConfig) *selectionPipeline {
	pipeline := newSelectionPipeline(candidates)
	pipeline.applyPolicy(c.policy, sc)
	pipeline.applyWorkflow(c.workflow, sc)
	pipeline.applyProgressive(c.progressive, c.store, key, identityHash, c.progressiveHash)
	pipeline.finalizeOutcome()
	return pipeline
}

func (p *selectionPipeline) applyPolicy(policy *PolicyFilter, sc SelectionContext) {
	if policy == nil {
		return
	}
	input := len(p.tools)
	var traces []DecisionTrace
	p.tools, traces = policy.Filter(p.tools, sc)
	recordStageCount(p.stageCounts, StagePolicy, input, len(p.tools))
	p.policyAllowed = p.tools
	p.traces = append(p.traces, traces...)
	if input > 0 && len(p.tools) == 0 {
		p.outcome = SelectionOutcomeExplicitDeny
	}
}

func (p *selectionPipeline) applyWorkflow(workflow *WorkflowSelector, sc SelectionContext) {
	if !p.selected() || workflow == nil {
		return
	}
	input := len(p.tools)
	result := workflow.filter(p.tools, sc)
	p.tools = result.tools
	recordStageCount(p.stageCounts, StageWorkflow, input, len(p.tools))
	p.traces = append(p.traces, result.traces...)
	if !result.matched && workflow.hasMatchableWorkflows() {
		p.outcome = SelectionOutcomeNoMatch
		return
	}
	if input > 0 && len(p.tools) == 0 {
		p.outcome = SelectionOutcomeEmptyByConfiguration
	}
}

func (p *selectionPipeline) applyProgressive(progressive *ProgressiveGate, store *SessionPlanStore, key PlanKey, identityHash, progressiveHash string) {
	if !p.selected() || progressive == nil {
		return
	}
	input := len(p.tools)
	var traces []DecisionTrace
	_, p.expanded = store.CallState(key, identityHash, progressiveHash)
	p.tools, traces = progressive.Apply(p.tools, p.expanded)
	recordStageCount(p.stageCounts, StageProgressive, input, len(p.tools))
	p.traces = append(p.traces, traces...)
	if input > 0 && len(p.tools) == 0 {
		p.outcome = SelectionOutcomeEmptyByConfiguration
	}
}

func (p *selectionPipeline) finalizeOutcome() {
	if p.selected() && len(p.tools) == 0 {
		p.outcome = SelectionOutcomeNoMatch
	}
}

func (p *selectionPipeline) selected() bool {
	return p.outcome == SelectionOutcomeSelected
}

func (c *CompositeSelector) newSelectionPlan(sc SelectionContext, candidates []model.ToolConfig, version, identityHash string, pipeline *selectionPipeline) *SelectionPlan {
	plan := &SelectionPlan{
		SessionID:          sc.SessionID,
		ToolNames:          toolNames(pipeline.tools),
		VisibleToolNames:   visibleToolNames(pipeline.tools),
		VisibleFingerprint: c.visibleFP(pipeline.tools),
		Mode:               ModeSelected,
		Outcome:            pipeline.outcome,
		StageCounts:        pipeline.stageCounts,
		Reasons:            pipeline.traces,
		Version:            version,
		CreatedAt:          time.Now().UnixNano(),
		IdentityHash:       identityHash,
		ProgressiveHash:    c.progressiveHash,
		ConfigHash:         c.configHash,
		CatalogVersion:     c.catalogVersion(candidates, sc.CatalogVersion),
		Expanded:           pipeline.expanded,
	}
	plan.toolSet = toolNameSet(plan.ToolNames)
	return plan
}

func (c *CompositeSelector) storeSelectionPlan(key PlanKey, plan *SelectionPlan, sc SelectionContext) (*SelectionPlan, error) {
	return c.store.Commit(key, plan, sc)
}

// RefreshPlan recomputes a plan from a cloned store context and commits it only
// if the original plan generation is still current. Dynamic catalog updates use
// this to avoid reviving deleted sessions or overwriting newer request plans.
func (c *CompositeSelector) RefreshPlan(_ context.Context, base SessionPlanContext, candidates []model.ToolConfig) (*SelectionPlan, bool, error) {
	start := time.Now()
	sc := base.Context
	identityHash, err := identityFingerprint(sc)
	if err != nil {
		return nil, false, err
	}
	version := c.version(candidates, identityHash, sc.SessionID, sc.CatalogVersion)
	plan := c.computeSelectionPlan(base.Key, identityHash, sc, candidates, version)
	committedPlan, committed, err := c.store.SetIfCurrent(base, plan)
	if err != nil || !committed {
		if err != nil {
			recordFallback("plan_persistence_failed")
		}
		return nil, committed, err
	}
	result := "ok"
	if outcomeAllowsFallback(committedPlan.Outcome) {
		result = "fallback"
		recordFallback(committedPlan.Outcome)
	}
	c.recordSelectionResult(result, committedPlan, candidates, start)
	c.log.Log(sc, committedPlan, len(candidates))
	return committedPlan, true, nil
}

func (c *CompositeSelector) recordSelectionResult(result string, plan *SelectionPlan, candidates []model.ToolConfig, start time.Time) {
	elapsedMS := float64(time.Since(start).Microseconds()) / 1000.0
	recordSelection(result, plan.Mode, len(candidates), len(plan.ToolNames), elapsedMS)
}

// HandleSelectionFailure returns a safe plan when Select fails. fail_closed
// exposes nothing; bundle_default is intersected with the hard-policy result.
func (c *CompositeSelector) HandleSelectionFailure(_ context.Context, sc SelectionContext, candidates []model.ToolConfig, _ error) (*SelectionPlan, error) {
	start := time.Now()
	identityHash, err := identityFingerprint(sc)
	if err != nil {
		identityHash = ""
	}
	version := c.version(candidates, identityHash, sc.SessionID, sc.CatalogVersion)
	key := c.planKey(sc.SessionID)

	allowed := candidates
	stageCounts := make(map[string]StageCount)
	traces := []DecisionTrace{{Stage: StageInternal, Detail: "selection_failure"}}
	if c.policy != nil {
		input := len(candidates)
		var t []DecisionTrace
		allowed, t = c.policy.Filter(candidates, sc)
		recordStageCount(stageCounts, StagePolicy, input, len(allowed))
		traces = append(traces, t...)
	}

	plan := c.applyFallback(sc, allowed, version, traces, stageCounts, SelectionOutcomeInternalError)
	plan.IdentityHash = identityHash
	plan.ConfigHash = c.configHash
	plan.CatalogVersion = c.catalogVersion(candidates, sc.CatalogVersion)
	plan, err = c.storeSelectionPlan(key, plan, sc)
	if err != nil {
		recordFallback("plan_persistence_failed")
		return nil, err
	}

	elapsedMS := float64(time.Since(start).Microseconds()) / 1000.0
	recordSelection("fallback", plan.Mode, len(candidates), len(plan.ToolNames), elapsedMS)
	recordFallback(SelectionOutcomeInternalError)
	c.log.Log(sc, plan, len(candidates))
	return plan, nil
}

// applyFallback builds a plan according to the configured fallback strategy.
// allowed is the policy-filtered candidate set, so the bundle_default branch
// can never re-expose a tool that policy explicitly denied.
func (c *CompositeSelector) applyFallback(sc SelectionContext, allowed []model.ToolConfig, version string, traces []DecisionTrace, stageCounts map[string]StageCount, outcome string) *SelectionPlan {
	identityHash, _ := identityFingerprint(sc)
	plan := &SelectionPlan{
		SessionID:       sc.SessionID,
		Outcome:         outcome,
		StageCounts:     stageCounts,
		Reasons:         traces,
		Version:         version,
		CreatedAt:       time.Now().UnixNano(),
		IdentityHash:    identityHash,
		ProgressiveHash: c.progressiveHash,
		ConfigHash:      c.configHash,
		CatalogVersion:  sc.CatalogVersion,
	}

	if c.fallback == FallbackFailClosed {
		plan.Mode = ModeFailClosed
		plan.ToolNames = nil
		plan.VisibleToolNames = []string{}
		plan.VisibleFingerprint = c.visibleFP(nil)
		plan.toolSet = nil
		return plan
	}

	// bundle_default: expose the named safe bundle (intersected with the
	// policy-allowed set).
	plan.Mode = ModeFallbackBundle
	bundle, ok := c.defaultBundleTools()
	if !ok {
		return plan
	}
	plan.ToolNames, plan.VisibleToolNames = fallbackBundleToolNames(allowed, bundle)
	plan.VisibleFingerprint = c.visibleFP(filterToolsByName(allowed, plan.VisibleToolNames))
	plan.toolSet = toolNameSet(plan.ToolNames)
	return plan
}

func (c *CompositeSelector) defaultBundleTools() (map[string]struct{}, bool) {
	if c.bundles == nil || c.defaultBundle == "" {
		return nil, false
	}
	return c.bundles.bundleTools(c.defaultBundle)
}

func fallbackBundleToolNames(allowed []model.ToolConfig, bundle map[string]struct{}) ([]string, []string) {
	toolNames := make([]string, 0, len(allowed))
	visibleToolNames := make([]string, 0, len(allowed))
	for _, t := range allowed {
		if _, in := bundle[t.Name]; !in {
			continue
		}
		toolNames = append(toolNames, t.Name)
		if toolVisible(t) {
			visibleToolNames = append(visibleToolNames, t.Name)
		}
	}
	return toolNames, visibleToolNames
}

func filterToolsByName(tools []model.ToolConfig, names []string) []model.ToolConfig {
	if len(tools) == 0 || len(names) == 0 {
		return nil
	}
	keep := make(map[string]struct{}, len(names))
	for _, name := range names {
		keep[name] = struct{}{}
	}
	out := make([]model.ToolConfig, 0, len(names))
	for _, tool := range tools {
		if _, ok := keep[tool.Name]; ok {
			out = append(out, tool)
		}
	}
	return out
}

func toolVisible(t model.ToolConfig) bool {
	return t.Meta == nil || t.Meta.DiscoveryVisibility == nil || *t.Meta.DiscoveryVisibility
}

// AuthorizeCall enforces that the requested tool is part of the current session
// plan. A cached plan is reused only while its version matches the live
// candidates and request claims; stale plans are recomputed before checking.
func (c *CompositeSelector) AuthorizeCall(ctx context.Context, sc SelectionContext, candidates []model.ToolConfig) (*AuthorizationReceipt, error) {
	key := c.planKey(sc.SessionID)
	identityHash, err := identityFingerprint(sc)
	if err != nil {
		recordCallDenied("identity_hash_error")
		return nil, ErrToolNotAuthorized
	}
	catalogVersion := c.catalogVersion(candidates, sc.CatalogVersion)
	expectedVersion := c.version(candidates, identityHash, sc.SessionID, catalogVersion)
	receipt, stale, err := c.store.IssueReceiptForVersion(ReceiptVersionRequest{
		Key:             key,
		Requested:       sc.Requested,
		ExpectedVersion: expectedVersion,
		IdentityHash:    identityHash,
		ConfigHash:      c.configHash,
		CatalogVersion:  catalogVersion,
		ProgressiveHash: c.progressiveHash,
		RouterID:        c.routerID,
	})
	if stale {
		receipt, err := c.recomputeAndIssueReceipt(ctx, key, identityHash, sc, candidates, expectedVersion)
		if err != nil {
			recordCallDenied("stale_plan_recompute_failed")
			return nil, ErrToolNotAuthorized
		}
		return receipt, nil
	}
	if err == nil {
		return &receipt, nil
	}
	recordCallDenied("not_in_plan")
	return nil, ErrToolNotAuthorized
}

func (c *CompositeSelector) recomputeAndIssueReceipt(_ context.Context, key PlanKey, identityHash string, sc SelectionContext, candidates []model.ToolConfig, version string) (*AuthorizationReceipt, error) {
	start := time.Now()
	plan := c.computeSelectionPlan(key, identityHash, sc, candidates, version)
	committedPlan, receipt, err := c.store.CommitAndIssueReceipt(key, plan, sc, sc.Requested, c.routerID)
	if err != nil {
		return nil, err
	}
	result := "ok"
	if outcomeAllowsFallback(committedPlan.Outcome) {
		result = "fallback"
		recordFallback(committedPlan.Outcome)
	}
	c.recordSelectionResult(result, committedPlan, candidates, start)
	c.log.Log(sc, committedPlan, len(candidates))
	return &receipt, nil
}

// RecordCallSuccess counts completed tool calls for progressive disclosure.
func (c *CompositeSelector) RecordCallSuccess(_ context.Context, receipt AuthorizationReceipt) (CallSuccessResult, error) {
	if c.progressive == nil {
		return CallSuccessResult{}, nil
	}
	return c.store.RecordCallSuccess(receipt, c.progressive.expandAfter), nil
}

func (c *CompositeSelector) FinalizeReceipt(_ context.Context, receipt AuthorizationReceipt, outcome ReceiptOutcome) (CallSuccessResult, error) {
	if c.progressive == nil {
		return c.store.FinalizeReceipt(receipt, outcome, 0), nil
	}
	return c.store.FinalizeReceipt(receipt, outcome, c.progressive.expandAfter), nil
}

// version combines the static config hash with a fingerprint of the candidate
// tools, validated claims, and progressive-disclosure state so cached plans
// invalidate when the registry, tool definitions, identity/policy inputs, or
// the expansion threshold changes.
func (c *CompositeSelector) version(candidates []model.ToolConfig, identityHash, sessionID, catalogVersion string) string {
	h := fnv.New64a()
	_, _ = h.Write([]byte("catalog:"))
	_, _ = h.Write([]byte(c.catalogVersion(candidates, catalogVersion)))
	_, _ = h.Write([]byte{0})

	_, _ = h.Write([]byte("identity:"))
	_, _ = h.Write([]byte(identityHash))
	_, _ = h.Write([]byte{0})

	// Include progressive expansion state so crossing the threshold invalidates cache.
	if c.progressive != nil {
		_, expanded := c.store.CallState(c.planKey(sessionID), identityHash, c.progressiveHash)
		if expanded {
			_, _ = h.Write([]byte("expanded"))
		} else {
			_, _ = h.Write([]byte("initial"))
		}
		_, _ = h.Write([]byte{0})
	}

	return c.configHash + ":" + strconv.FormatUint(h.Sum64(), 16)
}

func (c *CompositeSelector) catalogVersion(candidates []model.ToolConfig, catalogVersion string) string {
	if catalogVersion != "" {
		return catalogVersion
	}
	if len(candidates) == 0 {
		return "empty"
	}
	return legacyCatalogVersion(candidates)
}

func (c *CompositeSelector) planKey(sessionID string) PlanKey {
	return NewPlanKey(c.routerID, sessionID)
}

func identityFingerprint(sc SelectionContext) (string, error) {
	h := fnv.New64a()
	if err := writeClaimsFingerprint(h, sc); err != nil {
		return "", err
	}
	return strconv.FormatUint(h.Sum64(), 16), nil
}

func progressiveConfigHash(g *ProgressiveGate) string {
	if g == nil {
		return ""
	}
	h := fnv.New64a()
	_, _ = h.Write([]byte(g.initialBundle))
	_, _ = h.Write([]byte{0})
	_, _ = h.Write([]byte(strconv.Itoa(g.expandAfter)))
	return strconv.FormatUint(h.Sum64(), 16)
}

func recordStageCount(counts map[string]StageCount, stage string, input, output int) {
	if counts == nil {
		return
	}
	counts[stage] = StageCount{Input: input, Output: output}
}

func outcomeAllowsFallback(outcome string) bool {
	return outcome == SelectionOutcomeNoMatch || outcome == SelectionOutcomeInternalError
}

// writeClaimsFingerprint folds every validated claim into the plan version so
// policies and workflows that match arbitrary claim keys cannot reuse another
// claim set's plan. Only the final hash is exposed on SelectionPlan.Version.
func writeClaimsFingerprint(h io.Writer, sc SelectionContext) error {
	if sc.UserID != "" {
		_, _ = h.Write([]byte("sub:"))
		_, _ = h.Write([]byte(sc.UserID))
		_, _ = h.Write([]byte{0})
	}
	if sc.Tenant != "" {
		_, _ = h.Write([]byte("tenant:"))
		_, _ = h.Write([]byte(sc.Tenant))
		_, _ = h.Write([]byte{0})
	}
	if len(sc.Claims) == 0 {
		return nil
	}
	keys := make([]string, 0, len(sc.Claims))
	for k := range sc.Claims {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		if (k == "sub" && sc.UserID != "") || (k == "tenant" && sc.Tenant != "") {
			continue
		}
		_, _ = h.Write([]byte("claim:"))
		_, _ = h.Write([]byte(k))
		_, _ = h.Write([]byte("="))
		data, err := json.Marshal(sc.Claims[k])
		if err != nil {
			return fmt.Errorf("claim %q is not JSON-canonicalizable: %w", k, err)
		}
		_, _ = h.Write(data)
		_, _ = h.Write([]byte{0})
	}
	return nil
}

func legacyCatalogVersion(candidates []model.ToolConfig) string {
	prints := make([]string, len(candidates))
	for i, t := range candidates {
		prints[i] = legacyToolFingerprint(t)
	}
	sort.Strings(prints)

	h := fnv.New64a()
	for _, p := range prints {
		_, _ = h.Write([]byte(p))
		_, _ = h.Write([]byte{0})
	}
	return strconv.FormatUint(h.Sum64(), 16)
}

// legacyToolFingerprint is used only when a direct test or external caller
// invokes the selector without a registry snapshot version.
func legacyToolFingerprint(t model.ToolConfig) string {
	data, err := json.Marshal(t)
	if err != nil {
		return fmt.Sprintf("%#v", t)
	}
	return string(data)
}

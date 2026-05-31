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
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// Fallback strategies.
const (
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

	fallback      string
	defaultBundle string
	enforceOnCall bool
	configHash    string
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
	EnforceOnCall bool
	ConfigHash    string
}

// NewCompositeSelector assembles the pipeline from its options.
func NewCompositeSelector(opts CompositeOptions) *CompositeSelector {
	fallback := opts.Fallback
	if fallback == "" {
		fallback = FallbackBundleDefault
	}
	return &CompositeSelector{
		policy:        opts.Policy,
		workflow:      opts.Workflow,
		progressive:   opts.Progressive,
		bundles:       opts.Bundles,
		store:         opts.Store,
		log:           opts.Log,
		fallback:      fallback,
		defaultBundle: opts.DefaultBundle,
		enforceOnCall: opts.EnforceOnCall,
		configHash:    opts.ConfigHash,
	}
}

// Name identifies this selector.
func (c *CompositeSelector) Name() string { return "composite" }

// Select runs the pipeline, reusing a cached plan when the version is unchanged.
func (c *CompositeSelector) Select(_ context.Context, sc SelectionContext, candidates []model.ToolConfig) (*SelectionPlan, error) {
	version := c.version(candidates, sc)

	// Reuse an existing plan when metadata/config has not changed.
	if cached, ok := c.store.Get(sc.SessionID); ok && cached.Version == version {
		recordSelection("cached", cached.Mode, len(candidates), len(cached.ToolNames), 0)
		return cached, nil
	}

	start := time.Now()
	cur := candidates
	var traces []DecisionTrace

	// policyAllowed is the candidate set after the hard policy filter. The
	// fallback bundle is intersected with this (not raw candidates) so a tool
	// explicitly denied by policy is never re-exposed through fallback.
	policyAllowed := candidates
	if c.policy != nil {
		var t []DecisionTrace
		cur, t = c.policy.Filter(cur, sc)
		policyAllowed = cur
		traces = append(traces, t...)
	}

	if c.workflow != nil {
		var t []DecisionTrace
		cur, t = c.workflow.Filter(cur, sc)
		traces = append(traces, t...)
	}

	if c.progressive != nil {
		var t []DecisionTrace
		cur, t = c.progressive.Apply(cur, sc, c.store)
		traces = append(traces, t...)
	}

	plan := &SelectionPlan{
		SessionID:        sc.SessionID,
		ToolNames:        toolNames(cur),
		VisibleToolNames: visibleToolNames(cur),
		Mode:             ModeHybrid,
		Reasons:          traces,
		Version:          version,
		CreatedAt:        time.Now().UnixNano(),
	}

	result := "ok"
	// Empty selection triggers the configured fallback so discovery never
	// returns an unexpectedly empty tool set due to overly strict rules.
	if len(plan.ToolNames) == 0 {
		plan = c.applyFallback(sc, policyAllowed, version, traces)
		result = "fallback"
		recordFallback("empty")
	}

	c.store.Set(plan)
	setPlansActive(c.store.Len())

	elapsedMS := float64(time.Since(start).Microseconds()) / 1000.0
	recordSelection(result, plan.Mode, len(candidates), len(plan.ToolNames), elapsedMS)

	// Backfill AgentID from store for decision logs (it was set during initialize).
	if sc.AgentID == "" {
		sc.AgentID = c.store.AgentID(sc.SessionID)
	}
	c.log.Log(sc, plan, len(candidates))

	logger.Debugf("[dubbo-go-pixiu] mcp router selected %d/%d tools for session %s (mode=%s)",
		len(plan.ToolNames), len(candidates), sc.SessionID, plan.Mode)
	return plan, nil
}

// applyFallback builds a plan according to the configured fallback strategy.
// allowed is the policy-filtered candidate set, so the bundle_default branch
// can never re-expose a tool that policy explicitly denied.
func (c *CompositeSelector) applyFallback(sc SelectionContext, allowed []model.ToolConfig, version string, traces []DecisionTrace) *SelectionPlan {
	plan := &SelectionPlan{
		SessionID: sc.SessionID,
		Reasons:   traces,
		Version:   version,
		CreatedAt: time.Now().UnixNano(),
	}

	if c.fallback == FallbackFailClosed {
		plan.Mode = ModeFailClosed
		plan.ToolNames = nil
		plan.VisibleToolNames = []string{}
		return plan
	}

	// bundle_default: expose the named safe bundle (intersected with the
	// policy-allowed set).
	plan.Mode = ModeFallbackBundle
	if c.bundles != nil && c.defaultBundle != "" {
		if bundle, ok := c.bundles.bundleTools(c.defaultBundle); ok {
			for _, t := range allowed {
				if _, in := bundle[t.Name]; in {
					plan.ToolNames = append(plan.ToolNames, t.Name)
					if t.Meta == nil || t.Meta.DiscoveryVisibility == nil || *t.Meta.DiscoveryVisibility {
						plan.VisibleToolNames = append(plan.VisibleToolNames, t.Name)
					}
				}
			}
		}
	}
	return plan
}

// AuthorizeCall enforces that the requested tool is part of the current session
// plan. A cached plan is reused only while its version matches the live
// candidates and request claims; stale plans are recomputed before checking.
func (c *CompositeSelector) AuthorizeCall(ctx context.Context, sc SelectionContext, candidates []model.ToolConfig) error {
	if !c.enforceOnCall {
		return nil
	}
	plan, ok := c.store.Get(sc.SessionID)
	if !ok {
		// No plan yet: the client must call tools/list first. Deny to preserve
		// the discovery/execution separation guarantee.
		recordCallDenied("no_session_plan")
		return ErrToolNotAuthorized
	}
	if plan.Version != c.version(candidates, sc) {
		var err error
		plan, err = c.Select(ctx, sc, candidates)
		if err != nil {
			recordCallDenied("stale_plan_recompute_failed")
			return ErrToolNotAuthorized
		}
	}
	if plan.Contains(sc.Requested) {
		return nil
	}
	recordCallDenied("not_in_plan")
	return ErrToolNotAuthorized
}

// RecordCallSuccess counts completed tool calls for progressive disclosure.
func (c *CompositeSelector) RecordCallSuccess(_ context.Context, sc SelectionContext) error {
	plan, ok := c.store.Get(sc.SessionID)
	if !ok || !plan.Contains(sc.Requested) {
		return nil
	}
	c.store.IncrementCallCount(sc.SessionID)
	return nil
}

// OnInitialize stores the agent identifier for later decision logs. Plans are
// computed lazily at tools/list.
func (c *CompositeSelector) OnInitialize(_ context.Context, sc SelectionContext, _ []model.ToolConfig) error {
	if sc.SessionID != "" && sc.AgentID != "" {
		c.store.SetAgentID(sc.SessionID, sc.AgentID)
	}
	return nil
}

// InspectPlan returns the current cached plan for a session, for the admin
// debug endpoint. It satisfies the PlanInspector interface.
func (c *CompositeSelector) InspectPlan(sessionID string) (*SelectionPlan, bool) {
	return c.store.Get(sessionID)
}

// version combines the static config hash with a fingerprint of the candidate
// tools, validated claims, and progressive-disclosure state so cached plans
// invalidate when the registry, tool definitions, identity/policy inputs, or
// the expansion threshold changes.
func (c *CompositeSelector) version(candidates []model.ToolConfig, sc SelectionContext) string {
	prints := make([]string, len(candidates))
	for i, t := range candidates {
		prints[i] = toolFingerprint(t)
	}
	sort.Strings(prints)

	h := fnv.New64a()
	for _, p := range prints {
		_, _ = h.Write([]byte(p))
		_, _ = h.Write([]byte{0})
	}

	writeClaimsFingerprint(h, sc)

	// Include progressive expansion state so crossing the threshold invalidates cache.
	if c.progressive != nil {
		callCount := c.store.CallCount(sc.SessionID)
		expandAfter := c.progressive.expandAfter
		if callCount >= expandAfter {
			_, _ = h.Write([]byte("expanded"))
		} else {
			_, _ = h.Write([]byte("initial"))
		}
		_, _ = h.Write([]byte{0})
	}

	return c.configHash + ":" + strconv.FormatUint(h.Sum64(), 16)
}

// writeClaimsFingerprint folds every validated claim into the plan version so
// policies and workflows that match arbitrary claim keys cannot reuse another
// claim set's plan. Only the final hash is exposed on SelectionPlan.Version.
func writeClaimsFingerprint(h io.Writer, sc SelectionContext) {
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
		return
	}
	keys := make([]string, 0, len(sc.Claims))
	for k := range sc.Claims {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, k := range keys {
		_, _ = h.Write([]byte("claim:"))
		_, _ = h.Write([]byte(k))
		_, _ = h.Write([]byte("="))
		data, err := json.Marshal(sc.Claims[k])
		if err != nil {
			data = []byte(fmt.Sprintf("%v", sc.Claims[k]))
		}
		_, _ = h.Write(data)
		_, _ = h.Write([]byte{0})
	}
}

// toolFingerprint builds a stable string for a complete tool definition. The
// router authorizes against the live catalog, so backend-shape changes on a
// same-named tool must invalidate cached plans too.
func toolFingerprint(t model.ToolConfig) string {
	data, err := json.Marshal(t)
	if err != nil {
		return fmt.Sprintf("%#v", t)
	}
	return string(data)
}

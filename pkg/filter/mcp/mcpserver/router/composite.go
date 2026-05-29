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
	"hash/fnv"
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
//	policy -> workflow -> schema -> progressive
//
// Each stage is optional; a nil stage is skipped. Results are cached per session
// in the SessionPlanStore and reused while the metadata/config version is stable.
type CompositeSelector struct {
	policy      *PolicyFilter
	workflow    *WorkflowSelector
	schema      *SchemaMatcher   // Stage 3; nil in MVP-core
	progressive *ProgressiveGate // Stage 3; nil in MVP-core

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
	Schema        *SchemaMatcher
	Progressive   *ProgressiveGate
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
		schema:        opts.Schema,
		progressive:   opts.Progressive,
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
	version := c.version(candidates)

	// Reuse an existing plan when metadata/config has not changed.
	if cached, ok := c.store.Get(sc.SessionID); ok && cached.Version == version {
		recordSelection("cached", cached.Mode, len(candidates), len(cached.ToolNames), 0)
		return cached, nil
	}

	start := time.Now()
	cur := candidates
	var traces []DecisionTrace

	if c.policy != nil {
		var t []DecisionTrace
		cur, t = c.policy.Filter(cur, sc)
		traces = append(traces, t...)
	}

	if c.workflow != nil {
		var t []DecisionTrace
		cur, t = c.workflow.Filter(cur, sc)
		traces = append(traces, t...)
	}

	if c.schema != nil {
		var t []DecisionTrace
		cur, t = c.schema.Rank(cur, sc)
		traces = append(traces, t...)
	}

	if c.progressive != nil {
		var t []DecisionTrace
		cur, t = c.progressive.Apply(cur, sc, c.store)
		traces = append(traces, t...)
	}

	plan := &SelectionPlan{
		SessionID: sc.SessionID,
		ToolNames: toolNames(cur),
		Mode:      ModeHybrid,
		Reasons:   traces,
		Version:   version,
		CreatedAt: time.Now().UnixNano(),
	}

	result := "ok"
	// Empty selection triggers the configured fallback so discovery never
	// returns an unexpectedly empty tool set due to overly strict rules.
	if len(plan.ToolNames) == 0 {
		plan = c.applyFallback(sc, candidates, version, traces)
		result = "fallback"
		recordFallback("empty")
	}

	c.store.Set(plan)
	setPlansActive(c.store.Len())

	elapsedMS := float64(time.Since(start).Microseconds()) / 1000.0
	recordSelection(result, plan.Mode, len(candidates), len(plan.ToolNames), elapsedMS)
	c.log.Log(sc, plan, len(candidates))

	logger.Debugf("[dubbo-go-pixiu] mcp router selected %d/%d tools for session %s (mode=%s)",
		len(plan.ToolNames), len(candidates), sc.SessionID, plan.Mode)
	return plan, nil
}

// applyFallback builds a plan according to the configured fallback strategy.
func (c *CompositeSelector) applyFallback(sc SelectionContext, candidates []model.ToolConfig, version string, traces []DecisionTrace) *SelectionPlan {
	plan := &SelectionPlan{
		SessionID: sc.SessionID,
		Mode:      ModePassthrough,
		Reasons:   traces,
		Version:   version,
		CreatedAt: time.Now().UnixNano(),
	}

	if c.fallback == FallbackFailClosed {
		plan.Mode = "fail_closed"
		plan.ToolNames = nil
		return plan
	}

	// bundle_default: expose the named safe bundle (intersected with candidates).
	plan.Mode = "fallback_bundle"
	if c.workflow != nil && c.defaultBundle != "" {
		if bundle, ok := c.workflow.bundleTools(c.defaultBundle); ok {
			for _, t := range candidates {
				if _, in := bundle[t.Name]; in {
					plan.ToolNames = append(plan.ToolNames, t.Name)
				}
			}
		}
	}
	return plan
}

// AuthorizeCall enforces that the requested tool is part of the session plan.
func (c *CompositeSelector) AuthorizeCall(_ context.Context, sc SelectionContext) error {
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
	if plan.Contains(sc.Requested) {
		// Count successful authorizations to drive progressive disclosure.
		c.store.IncrementCallCount(sc.SessionID)
		return nil
	}
	recordCallDenied("not_in_plan")
	return ErrToolNotAuthorized
}

// OnInitialize is a no-op; plans are computed lazily at tools/list.
func (c *CompositeSelector) OnInitialize(_ context.Context, _ SelectionContext, _ []model.ToolConfig) error {
	return nil
}

// InspectPlan returns the current cached plan for a session, for the admin
// debug endpoint. It satisfies the PlanInspector interface.
func (c *CompositeSelector) InspectPlan(sessionID string) (*SelectionPlan, bool) {
	return c.store.Get(sessionID)
}

// version combines the static config hash with a fingerprint of the candidate
// tool names so a dynamic registry update invalidates cached plans.
func (c *CompositeSelector) version(candidates []model.ToolConfig) string {
	names := make([]string, len(candidates))
	for i, t := range candidates {
		names[i] = t.Name
	}
	sort.Strings(names)

	h := fnv.New64a()
	for _, n := range names {
		_, _ = h.Write([]byte(n))
		_, _ = h.Write([]byte{0})
	}
	return c.configHash + ":" + strconv.FormatUint(h.Sum64(), 16)
}

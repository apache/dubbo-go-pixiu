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
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// ProgressiveGate implements session-level progressive disclosure: a fresh
// session sees only the initial bundle; after enough successful tool calls the
// full (already policy/workflow-filtered) set is revealed.
//
// The gate is intentionally simple in the MVP: a single expansion threshold
// flips a session from "initial bundle" to "everything passed in". Tier-based
// staging beyond two levels is deferred until there is demand.
type ProgressiveGate struct {
	workflow      *WorkflowSelector
	initialBundle string
	expandAfter   int
}

// NewProgressiveGate builds the gate from config. The workflow selector is used
// to resolve the named initial bundle to a tool set.
func NewProgressiveGate(cfg model.ProgressiveConfig, wf *WorkflowSelector) *ProgressiveGate {
	expand := cfg.ExpandAfterCalls
	if expand <= 0 {
		expand = 1
	}
	return &ProgressiveGate{
		workflow:      wf,
		initialBundle: cfg.InitialBundle,
		expandAfter:   expand,
	}
}

// Apply restricts the candidate set to the initial bundle until the session has
// accrued enough successful calls, after which the full set passes through.
func (g *ProgressiveGate) Apply(tools []model.ToolConfig, sc SelectionContext, store *SessionPlanStore) ([]model.ToolConfig, []DecisionTrace) {
	// Once expanded, expose the full set.
	if store.CallCount(sc.SessionID) >= g.expandAfter {
		return tools, nil
	}

	// Before expansion, only the initial bundle is visible.
	bundle, ok := g.bundleSet()
	if !ok {
		// No usable initial bundle configured: behave as passthrough so a
		// misconfiguration never hides every tool.
		return tools, nil
	}

	kept := make([]model.ToolConfig, 0, len(bundle))
	var traces []DecisionTrace
	for _, t := range tools {
		if _, in := bundle[t.Name]; in {
			kept = append(kept, t)
		} else {
			traces = append(traces, DecisionTrace{
				Tool:   t.Name,
				Kept:   false,
				Stage:  StageProgressive,
				Detail: "tier_locked",
			})
		}
	}
	return kept, traces
}

// bundleSet resolves the initial bundle name to its tool-name set.
func (g *ProgressiveGate) bundleSet() (map[string]struct{}, bool) {
	if g.workflow == nil || g.initialBundle == "" {
		return nil, false
	}
	return g.workflow.bundleTools(g.initialBundle)
}

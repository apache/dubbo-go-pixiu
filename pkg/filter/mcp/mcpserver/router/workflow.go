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
	"fmt"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// compiledWorkflow is a WorkflowConfig with its When clause precompiled and its
// tool membership hoisted into a set.
type compiledWorkflow struct {
	name    string
	when    *matcher
	tools   map[string]struct{}
	hasWhen bool
}

// WorkflowSelector narrows the candidate set to a matched workflow bundle.
//
// The first workflow whose When clause matches the context wins; only tools in
// that bundle (and present in the candidate set) are kept. When no workflow
// matches, the candidates pass through unchanged. Workflows with an empty When
// clause are treated as bundles addressable by name (e.g. fallback bundles) and
// are never auto-matched.
type WorkflowSelector struct {
	workflows []compiledWorkflow
	byName    map[string]compiledWorkflow
}

// NewWorkflowSelector compiles the workflow definitions.
func NewWorkflowSelector(cfgs []model.WorkflowConfig) (*WorkflowSelector, error) {
	ws := &WorkflowSelector{
		workflows: make([]compiledWorkflow, 0, len(cfgs)),
		byName:    make(map[string]compiledWorkflow, len(cfgs)),
	}
	for _, c := range cfgs {
		m, err := newMatcher(c.When)
		if err != nil {
			return nil, fmt.Errorf("workflow %q: %w", c.Name, err)
		}
		cw := compiledWorkflow{
			name:    c.Name,
			when:    m,
			tools:   toSet(c.Tools),
			hasWhen: !m.alwaysMatches(),
		}
		ws.workflows = append(ws.workflows, cw)
		ws.byName[c.Name] = cw
	}
	return ws, nil
}

// Filter keeps only tools belonging to the first matching workflow bundle.
func (w *WorkflowSelector) Filter(tools []model.ToolConfig, sc SelectionContext) ([]model.ToolConfig, []DecisionTrace) {
	matched, ok := w.matchWorkflow(sc)
	if !ok {
		return tools, nil
	}

	kept := make([]model.ToolConfig, 0, len(matched.tools))
	var traces []DecisionTrace
	for _, tool := range tools {
		if _, in := matched.tools[tool.Name]; in {
			kept = append(kept, tool)
		} else {
			traces = append(traces, DecisionTrace{
				Tool:   tool.Name,
				Kept:   false,
				Stage:  StageWorkflow,
				Rule:   matched.name,
				Detail: "not_in_workflow",
			})
		}
	}
	return kept, traces
}

// matchWorkflow returns the first workflow whose When clause matches. Workflows
// without a When clause are skipped (they are name-addressable bundles only).
func (w *WorkflowSelector) matchWorkflow(sc SelectionContext) (compiledWorkflow, bool) {
	for _, wf := range w.workflows {
		if wf.hasWhen && wf.when.matches(sc) {
			return wf, true
		}
	}
	return compiledWorkflow{}, false
}

// bundleTools returns the tool-name set for a named workflow bundle, used by
// fallback and progressive disclosure. The second return is false if unknown.
func (w *WorkflowSelector) bundleTools(name string) (map[string]struct{}, bool) {
	wf, ok := w.byName[name]
	if !ok {
		return nil, false
	}
	return wf.tools, true
}

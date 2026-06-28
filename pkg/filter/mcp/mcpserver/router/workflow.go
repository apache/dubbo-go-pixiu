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
	"strings"
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

type workflowResult struct {
	tools   []model.ToolConfig
	traces  []DecisionTrace
	matched bool
	rule    string
}

// NewWorkflowSelector compiles the workflow definitions.
func NewWorkflowSelector(cfgs []model.WorkflowConfig) (*WorkflowSelector, error) {
	ws := &WorkflowSelector{
		workflows: make([]compiledWorkflow, 0, len(cfgs)),
		byName:    make(map[string]compiledWorkflow, len(cfgs)),
	}
	for i, c := range cfgs {
		name := strings.TrimSpace(c.Name)
		if name == "" {
			return nil, fmt.Errorf("workflow at index %d: name is required", i)
		}
		if _, exists := ws.byName[name]; exists {
			return nil, fmt.Errorf("workflow %q is defined more than once", name)
		}
		m, err := newMatcher(c.When)
		if err != nil {
			return nil, fmt.Errorf("workflow %q: %w", name, err)
		}
		tools, err := compileWorkflowTools(name, c.Tools)
		if err != nil {
			return nil, err
		}
		cw := compiledWorkflow{
			name:    name,
			when:    m,
			tools:   tools,
			hasWhen: !m.alwaysMatches(),
		}
		ws.workflows = append(ws.workflows, cw)
		ws.byName[name] = cw
	}
	return ws, nil
}

func compileWorkflowTools(workflowName string, raw []string) (map[string]struct{}, error) {
	tools := make(map[string]struct{}, len(raw))
	for i, value := range raw {
		name := strings.TrimSpace(value)
		if name == "" {
			return nil, fmt.Errorf("workflow %q: tool at index %d is empty", workflowName, i)
		}
		if _, exists := tools[name]; exists {
			return nil, fmt.Errorf("workflow %q: duplicate tool %q at index %d", workflowName, name, i)
		}
		tools[name] = struct{}{}
	}
	return tools, nil
}

// Filter keeps only tools belonging to the first matching workflow bundle.
func (w *WorkflowSelector) Filter(tools []model.ToolConfig, sc SelectionContext) ([]model.ToolConfig, []DecisionTrace) {
	result := w.filter(tools, sc)
	return result.tools, result.traces
}

func (w *WorkflowSelector) filter(tools []model.ToolConfig, sc SelectionContext) workflowResult {
	matched, ok := w.matchWorkflow(sc)
	if !ok {
		return workflowResult{tools: tools}
	}

	kept := make([]model.ToolConfig, 0, len(matched.tools))
	var traces []DecisionTrace
	for _, tool := range tools {
		if _, in := matched.tools[tool.Name]; in {
			kept = append(kept, tool)
		} else {
			if len(traces) < maxDecisionTraceSamples {
				traces = append(traces, DecisionTrace{
					Tool:   tool.Name,
					Stage:  StageWorkflow,
					Rule:   matched.name,
					Detail: "not_in_workflow",
				})
			}
		}
	}
	return workflowResult{tools: kept, traces: traces, matched: true, rule: matched.name}
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

func (w *WorkflowSelector) hasMatchableWorkflows() bool {
	for _, wf := range w.workflows {
		if wf.hasWhen {
			return true
		}
	}
	return false
}

// bundleTools returns the tool-name set for a named workflow bundle, used by
// fallback and progressive disclosure. The second return is false if unknown.
func (w *WorkflowSelector) bundleTools(name string) (map[string]struct{}, bool) {
	wf, ok := w.byName[strings.TrimSpace(name)]
	if !ok {
		return nil, false
	}
	return wf.tools, true
}

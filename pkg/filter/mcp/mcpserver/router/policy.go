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

// riskLevel maps a risk string to an ordinal for max_risk comparison. An empty
// or unknown value is treated as low so untagged tools are never over-blocked.
func riskLevel(risk string) int {
	switch risk {
	case "high":
		return 3
	case "medium":
		return 2
	case "low", "":
		return 1
	default:
		return 1
	}
}

// compiledRule is a PolicyRule with its When clause precompiled and allow/deny
// tag sets hoisted into maps for O(1) lookup.
type compiledRule struct {
	name      string
	when      *matcher
	allowTags map[string]struct{}
	denyTags  map[string]struct{}
	maxRisk   int // 0 = no limit
}

// PolicyFilter applies hard-filter rules. A tool is dropped when any applicable
// rule denies it (deny tag, missing required allow tag, or risk over limit).
// Rules whose When clause does not match the context are skipped entirely.
type PolicyFilter struct {
	rules []compiledRule
}

// NewPolicyFilter compiles the policy rules. It fails fast on invalid regex so
// configuration errors are caught at startup, not per request.
func NewPolicyFilter(cfg model.PolicyConfig) (*PolicyFilter, error) {
	rules := make([]compiledRule, 0, len(cfg.Rules))
	for _, r := range cfg.Rules {
		m, err := newMatcher(r.When)
		if err != nil {
			return nil, fmt.Errorf("policy rule %q: %w", r.Name, err)
		}
		rules = append(rules, compiledRule{
			name:      r.Name,
			when:      m,
			allowTags: toSet(r.AllowTags),
			denyTags:  toSet(r.DenyTags),
			maxRisk:   maxRiskOrdinal(r.MaxRisk),
		})
	}
	return &PolicyFilter{rules: rules}, nil
}

// maxRiskOrdinal converts a configured max_risk to an ordinal, 0 meaning unset.
func maxRiskOrdinal(risk string) int {
	if risk == "" {
		return 0
	}
	return riskLevel(risk)
}

// Filter returns the subset of tools allowed by all applicable rules, along
// with a DecisionTrace for every dropped tool.
func (p *PolicyFilter) Filter(tools []model.ToolConfig, sc SelectionContext) ([]model.ToolConfig, []DecisionTrace) {
	if len(p.rules) == 0 {
		return tools, nil
	}

	// Pre-select the rules that apply to this context.
	applicable := make([]compiledRule, 0, len(p.rules))
	for _, r := range p.rules {
		if r.when.matches(sc) {
			applicable = append(applicable, r)
		}
	}
	if len(applicable) == 0 {
		return tools, nil
	}

	kept := make([]model.ToolConfig, 0, len(tools))
	var traces []DecisionTrace

	for _, tool := range tools {
		if rule, reason := denyReason(applicable, tool); reason != "" {
			traces = append(traces, DecisionTrace{
				Tool:   tool.Name,
				Kept:   false,
				Stage:  StagePolicy,
				Rule:   rule,
				Detail: reason,
			})
			continue
		}
		kept = append(kept, tool)
	}

	return kept, traces
}

// denyReason returns the first rule name and reason that denies the tool, or an
// empty reason if every applicable rule allows it.
func denyReason(rules []compiledRule, tool model.ToolConfig) (string, string) {
	tags := toolTags(tool)
	risk := toolRisk(tool)

	for _, r := range rules {
		// Deny tags take precedence: a single match blocks the tool.
		for tag := range r.denyTags {
			if _, has := tags[tag]; has {
				return r.name, "deny_tag:" + tag
			}
		}
		// Allow tags, when present, require at least one intersection.
		if len(r.allowTags) > 0 && !intersects(tags, r.allowTags) {
			return r.name, "no_allow_tag"
		}
		// Risk ceiling.
		if r.maxRisk > 0 && riskLevel(risk) > r.maxRisk {
			return r.name, "risk_exceeds_" + risk
		}
	}
	return "", ""
}

// toolTags returns the tool's tag set, empty if the tool has no Meta.
func toolTags(tool model.ToolConfig) map[string]struct{} {
	if tool.Meta == nil {
		return nil
	}
	return toSet(tool.Meta.Tags)
}

// toolRisk returns the tool's declared risk, defaulting to low.
func toolRisk(tool model.ToolConfig) string {
	if tool.Meta == nil || tool.Meta.Risk == "" {
		return "low"
	}
	return tool.Meta.Risk
}

func toSet(items []string) map[string]struct{} {
	if len(items) == 0 {
		return nil
	}
	s := make(map[string]struct{}, len(items))
	for _, item := range items {
		s[item] = struct{}{}
	}
	return s
}

func intersects(a, b map[string]struct{}) bool {
	// Iterate the smaller set for efficiency.
	if len(a) > len(b) {
		a, b = b, a
	}
	for k := range a {
		if _, ok := b[k]; ok {
			return true
		}
	}
	return false
}

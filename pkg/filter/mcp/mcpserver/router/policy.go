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

const policyRuleErrorFormat = "policy rule %q: %w"

// ValidateTools validates routing metadata on configured tools.
func ValidateTools(tools []model.ToolConfig) error {
	if err := ValidateUniqueNames(tools); err != nil {
		return err
	}
	for _, tool := range tools {
		if tool.Meta == nil {
			continue
		}
		if err := ValidateRisk(tool.Meta.Risk); err != nil {
			return fmt.Errorf("tool %q meta.risk: %w", tool.Name, err)
		}
	}
	return nil
}

// ValidateRisk checks the public risk enum. Empty is allowed and treated as low.
func ValidateRisk(risk string) error {
	switch risk {
	case "", "low", "medium", "high":
		return nil
	default:
		return fmt.Errorf("unsupported risk %q (allowed: low, medium, high)", risk)
	}
}

// validatePolicyRisks validates max_risk values even when the policy stage is
// off, because an unknown risk enum is a configuration error.
func validatePolicyRisks(cfg model.PolicyConfig) error {
	for _, r := range cfg.Rules {
		if _, err := maxRiskOrdinal(r.MaxRisk); err != nil {
			return fmt.Errorf(policyRuleErrorFormat, r.Name, err)
		}
	}
	return nil
}

// riskLevel maps a validated risk string to an ordinal for max_risk comparison.
// Empty means low so untagged tools stay backward compatible.
func riskLevel(risk string) int {
	switch risk {
	case "high":
		return 3
	case "medium":
		return 2
	case "low", "":
		return 1
	default:
		// Callers validate configured risk values before reaching runtime. Keep
		// this defensive fallback conservative if future inputs bypass validation.
		return 3
	}
}

// compiledRule is a PolicyRule with its When clause precompiled and allow/deny
// tag sets hoisted into maps for O(1) lookup.
type compiledRule struct {
	name      string
	when      *matcher
	allowTags StringSet
	denyTags  StringSet
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
			return nil, fmt.Errorf(policyRuleErrorFormat, r.Name, err)
		}
		maxRisk, err := maxRiskOrdinal(r.MaxRisk)
		if err != nil {
			return nil, fmt.Errorf(policyRuleErrorFormat, r.Name, err)
		}
		rules = append(rules, compiledRule{
			name:      r.Name,
			when:      m,
			allowTags: NewStringSet(r.AllowTags),
			denyTags:  NewStringSet(r.DenyTags),
			maxRisk:   maxRisk,
		})
	}
	return &PolicyFilter{rules: rules}, nil
}

// maxRiskOrdinal converts a configured max_risk to an ordinal, 0 meaning unset.
func maxRiskOrdinal(risk string) (int, error) {
	if risk == "" {
		return 0, nil
	}
	if err := ValidateRisk(risk); err != nil {
		return 0, fmt.Errorf("max_risk: %w", err)
	}
	return riskLevel(risk), nil
}

// Filter returns the subset of tools allowed by all applicable rules, along
// with bounded DecisionTrace samples for dropped tools.
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
			if len(traces) < maxDecisionTraceSamples {
				traces = append(traces, DecisionTrace{
					Tool:   tool.Name,
					Stage:  StagePolicy,
					Rule:   rule,
					Detail: reason,
				})
			}
			continue
		}
		kept = append(kept, tool)
	}

	return kept, traces
}

// denyReason returns the first rule name and reason that denies the tool, or an
// empty reason if every applicable rule allows it.
//
// Note the combination semantics: every applicable rule must allow the tool for
// it to be kept (logical AND). In particular, if two rules both specify
// allow_tags, the tool must carry a tag matching each rule's allow_tags set;
// satisfying only one rule is not enough.
func denyReason(rules []compiledRule, tool model.ToolConfig) (string, string) {
	tags := toolTags(tool)
	risk := toolRisk(tool)

	for _, r := range rules {
		// Deny tags take precedence: a single match blocks the tool.
		for tag := range r.denyTags {
			if tags.Contains(tag) {
				return r.name, "deny_tag:" + tag
			}
		}
		// Allow tags, when present, require at least one intersection.
		if len(r.allowTags) > 0 && !tags.Intersects(r.allowTags) {
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
func toolTags(tool model.ToolConfig) StringSet {
	if tool.Meta == nil {
		return nil
	}
	return NewStringSet(tool.Meta.Tags)
}

// toolRisk returns the tool's declared risk, defaulting to low.
func toolRisk(tool model.ToolConfig) string {
	if tool.Meta == nil || tool.Meta.Risk == "" {
		return "low"
	}
	return tool.Meta.Risk
}

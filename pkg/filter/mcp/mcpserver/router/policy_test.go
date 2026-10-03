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
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func toolWithMeta(name string, meta *model.ToolMeta) model.ToolConfig {
	return model.ToolConfig{Name: name, Cluster: "c", Meta: meta}
}

func keptNames(tools []model.ToolConfig) []string {
	return toolNames(tools)
}

func TestPolicyFilter_NoRulesAllAllowed(t *testing.T) {
	pf, err := NewPolicyFilter(model.PolicyConfig{})
	require.NoError(t, err)

	tools := testTools("a", "b")
	out, traces := pf.Filter(tools, SelectionContext{})

	assert.Equal(t, tools, out)
	assert.Nil(t, traces)
}

func TestPolicyFilter_DenyTagBlocks(t *testing.T) {
	pf, err := NewPolicyFilter(model.PolicyConfig{Rules: []model.PolicyRule{
		{Name: "no-admin", DenyTags: []string{"admin"}},
	}})
	require.NoError(t, err)

	tools := []model.ToolConfig{
		toolWithMeta("safe", &model.ToolMeta{Tags: []string{"user"}}),
		toolWithMeta("danger", &model.ToolMeta{Tags: []string{"admin"}}),
	}

	out, traces := pf.Filter(tools, SelectionContext{})

	assert.Equal(t, []string{"safe"}, keptNames(out))
	require.Len(t, traces, 1)
	assert.Equal(t, "danger", traces[0].Tool)
	assert.Equal(t, "deny_tag:admin", traces[0].Detail)
}

func TestPolicyFilter_AllowTagsRequireIntersection(t *testing.T) {
	pf, err := NewPolicyFilter(model.PolicyConfig{Rules: []model.PolicyRule{
		{Name: "tenant", AllowTags: []string{"acme", "shared"}},
	}})
	require.NoError(t, err)

	tools := []model.ToolConfig{
		toolWithMeta("a", &model.ToolMeta{Tags: []string{"acme"}}),
		toolWithMeta("b", &model.ToolMeta{Tags: []string{"shared", "x"}}),
		toolWithMeta("c", &model.ToolMeta{Tags: []string{"other"}}),
		toolWithMeta("d", nil), // no tags -> no allow tag -> dropped
	}

	out, _ := pf.Filter(tools, SelectionContext{})

	assert.ElementsMatch(t, []string{"a", "b"}, keptNames(out))
}

func TestPolicyFilter_MaxRisk(t *testing.T) {
	pf, err := NewPolicyFilter(model.PolicyConfig{Rules: []model.PolicyRule{
		{Name: "anon", When: model.PolicyMatch{MissingClaim: "sub"}, MaxRisk: "low"},
	}})
	require.NoError(t, err)

	tools := []model.ToolConfig{
		toolWithMeta("low1", &model.ToolMeta{Risk: "low"}),
		toolWithMeta("med1", &model.ToolMeta{Risk: "medium"}),
		toolWithMeta("untagged", nil), // defaults to low -> kept
	}

	// Anonymous (no sub) -> rule applies -> only low risk allowed.
	out, _ := pf.Filter(tools, SelectionContext{})
	assert.ElementsMatch(t, []string{"low1", "untagged"}, keptNames(out))

	// Authenticated (sub present) -> rule does not apply -> all kept.
	out2, _ := pf.Filter(tools, SelectionContext{UserID: "user-1"})
	assert.ElementsMatch(t, []string{"low1", "med1", "untagged"}, keptNames(out2))
}

func TestPolicyFilter_InvalidMaxRiskFailsFast(t *testing.T) {
	_, err := NewPolicyFilter(model.PolicyConfig{Rules: []model.PolicyRule{
		{Name: "bad-risk", MaxRisk: "hihg"},
	}})
	assert.ErrorContains(t, err, "max_risk")
	assert.ErrorContains(t, err, "unsupported risk")
}

func TestValidateTools_RiskEnum(t *testing.T) {
	valid := []model.ToolConfig{
		toolWithMeta("empty", &model.ToolMeta{}),
		toolWithMeta("low", &model.ToolMeta{Risk: "low"}),
		toolWithMeta("medium", &model.ToolMeta{Risk: "medium"}),
		toolWithMeta("high", &model.ToolMeta{Risk: "high"}),
		{Name: "no-meta"},
	}
	require.NoError(t, ValidateTools(valid))

	err := ValidateTools([]model.ToolConfig{
		toolWithMeta("typo", &model.ToolMeta{Risk: "hihg"}),
	})
	assert.ErrorContains(t, err, `tool "typo" meta.risk`)
	assert.ErrorContains(t, err, "unsupported risk")
}

func TestPolicyFilter_WhenClaimEquals(t *testing.T) {
	pf, err := NewPolicyFilter(model.PolicyConfig{Rules: []model.PolicyRule{
		{
			Name:     "acme-only",
			When:     model.PolicyMatch{Claim: "tenant", Equals: "acme"},
			DenyTags: []string{"internal"},
		},
	}})
	require.NoError(t, err)

	tools := []model.ToolConfig{
		toolWithMeta("pub", &model.ToolMeta{Tags: []string{"public"}}),
		toolWithMeta("int", &model.ToolMeta{Tags: []string{"internal"}}),
	}

	// tenant=acme -> rule applies -> internal dropped.
	out, _ := pf.Filter(tools, SelectionContext{Tenant: "acme"})
	assert.Equal(t, []string{"pub"}, keptNames(out))

	// tenant=globex -> rule does not apply -> all kept.
	out2, _ := pf.Filter(tools, SelectionContext{Tenant: "globex"})
	assert.ElementsMatch(t, []string{"pub", "int"}, keptNames(out2))
}

func TestPolicyFilter_InvalidRegexFailsFast(t *testing.T) {
	_, err := NewPolicyFilter(model.PolicyConfig{Rules: []model.PolicyRule{
		{Name: "bad", When: model.PolicyMatch{Claim: "x", Regex: "([a-z"}},
	}})
	assert.Error(t, err)
}

func TestPolicyFilter_InvalidMatchFailsFast(t *testing.T) {
	tests := []struct {
		name    string
		match   model.PolicyMatch
		wantErr string
	}{
		{
			name:    "operator without claim",
			match:   model.PolicyMatch{Equals: "admin"},
			wantErr: "claim is required",
		},
		{
			name:    "multiple operators",
			match:   model.PolicyMatch{Claim: "role", Equals: "admin", Regex: "^admin"},
			wantErr: "only one of equals, in, or regex",
		},
		{
			name:    "missing claim combined with claim",
			match:   model.PolicyMatch{MissingClaim: "sub", Claim: "role"},
			wantErr: "missing_claim cannot be combined",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewPolicyFilter(model.PolicyConfig{Rules: []model.PolicyRule{
				{Name: "bad", When: tt.match},
			}})
			assert.ErrorContains(t, err, tt.wantErr)
		})
	}
}

func TestPolicyFilter_RegexMatch(t *testing.T) {
	pf, err := NewPolicyFilter(model.PolicyConfig{Rules: []model.PolicyRule{
		{Name: "role", When: model.PolicyMatch{Claim: "role", Regex: "^admin.*"}, DenyTags: []string{"secret"}},
	}})
	require.NoError(t, err)

	tools := []model.ToolConfig{
		toolWithMeta("s", &model.ToolMeta{Tags: []string{"secret"}}),
	}

	out, _ := pf.Filter(tools, SelectionContext{Claims: map[string]any{"role": "administrator"}})
	assert.Empty(t, keptNames(out))

	out2, _ := pf.Filter(tools, SelectionContext{Claims: map[string]any{"role": "viewer"}})
	assert.Equal(t, []string{"s"}, keptNames(out2))
}

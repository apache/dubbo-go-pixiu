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

func TestPolicyFilter_NoRulesPassthrough(t *testing.T) {
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

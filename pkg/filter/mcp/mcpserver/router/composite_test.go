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
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// buildComposite is a test helper that wires a CompositeSelector from a config.
func buildComposite(t *testing.T, cfg *model.RouterConfig) (*CompositeSelector, *SessionPlanStore) {
	store := NewSessionPlanStoreWithTTL(time.Minute)
	sel, err := Build(cfg, store)
	require.NoError(t, err)
	require.NotNil(t, sel)
	cs, ok := sel.(*CompositeSelector)
	require.True(t, ok)
	return cs, store
}

func tenantTools() []model.ToolConfig {
	return []model.ToolConfig{
		toolWithMeta("acme_read", &model.ToolMeta{Tags: []string{"acme"}}),
		toolWithMeta("acme_write", &model.ToolMeta{Tags: []string{"acme"}}),
		toolWithMeta("globex_read", &model.ToolMeta{Tags: []string{"globex"}}),
		toolWithMeta("shared_ping", &model.ToolMeta{Tags: []string{"shared"}}),
	}
}

func TestComposite_PolicyTenantIsolation(t *testing.T) {
	cfg := &model.RouterConfig{
		Enabled: true,
		Policy: model.PolicyConfig{Rules: []model.PolicyRule{
			{Name: "acme", When: model.PolicyMatch{Claim: "tenant", Equals: "acme"}, AllowTags: []string{"acme", "shared"}},
		}},
	}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	plan, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1", Tenant: "acme"}, tenantTools())
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{"acme_read", "acme_write", "shared_ping"}, plan.ToolNames)
	assert.NotContains(t, plan.ToolNames, "globex_read")
}

func TestComposite_PlanReuseByVersion(t *testing.T) {
	cfg := &model.RouterConfig{Enabled: true}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	tools := testTools("a", "b")
	p1, _ := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, tools)
	p2, _ := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, tools)

	// Same version => same cached plan instance returned.
	assert.Same(t, p1, p2)
}

func TestComposite_PlanRecomputesOnRegistryChange(t *testing.T) {
	cfg := &model.RouterConfig{Enabled: true}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	p1, _ := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, testTools("a", "b"))
	p2, _ := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, testTools("a", "b", "c"))

	assert.NotEqual(t, p1.Version, p2.Version)
	assert.Len(t, p2.ToolNames, 3)
}

func TestComposite_EnforceOnCallDeniesWithoutPlan(t *testing.T) {
	cfg := &model.RouterConfig{Enabled: true}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	err := cs.AuthorizeCall(context.Background(), SelectionContext{SessionID: "no-plan", Requested: "x"})
	assert.ErrorIs(t, err, ErrToolNotAuthorized)
}

func TestComposite_EnforceOnCallAllowsInPlan(t *testing.T) {
	cfg := &model.RouterConfig{Enabled: true}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, testTools("a", "b"))

	assert.NoError(t, cs.AuthorizeCall(context.Background(), SelectionContext{SessionID: "s1", Requested: "a"}))
	assert.ErrorIs(t, cs.AuthorizeCall(context.Background(), SelectionContext{SessionID: "s1", Requested: "ghost"}), ErrToolNotAuthorized)
}

func TestComposite_EnforceOnCallDisabledAllowsAll(t *testing.T) {
	disabled := false
	cfg := &model.RouterConfig{Enabled: true, EnforceOnCall: &disabled}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	assert.NoError(t, cs.AuthorizeCall(context.Background(), SelectionContext{SessionID: "anything", Requested: "x"}))
}

func TestComposite_FallbackBundleDefault(t *testing.T) {
	cfg := &model.RouterConfig{
		Enabled:       true,
		DefaultBundle: "safe-minimal",
		Workflows: []model.WorkflowConfig{
			{Name: "safe-minimal", Tools: []string{"ping"}},
			{Name: "acme-flow", Tools: []string{"acme_only"}, When: model.PolicyMatch{Claim: "tenant", Equals: "acme"}},
		},
		Policy: model.PolicyConfig{Rules: []model.PolicyRule{
			// Deny everything for this tenant to force an empty selection.
			{Name: "block-all", When: model.PolicyMatch{Claim: "tenant", Equals: "acme"}, DenyTags: []string{"x"}},
		}},
	}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	// acme-flow matches and keeps only acme_only, but acme_only is not in the
	// candidate set, producing an empty selection -> fallback to safe-minimal.
	plan, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1", Tenant: "acme"},
		[]model.ToolConfig{toolWithMeta("ping", nil)})
	require.NoError(t, err)

	assert.Equal(t, []string{"ping"}, plan.ToolNames)
	assert.Equal(t, "fallback_bundle", plan.Mode)
}

func TestComposite_FallbackFailClosed(t *testing.T) {
	cfg := &model.RouterConfig{
		Enabled:  true,
		Fallback: FallbackFailClosed,
		Workflows: []model.WorkflowConfig{
			{Name: "acme-flow", Tools: []string{"acme_only"}, When: model.PolicyMatch{Claim: "tenant", Equals: "acme"}},
		},
	}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	plan, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1", Tenant: "acme"},
		[]model.ToolConfig{toolWithMeta("ping", nil)})
	require.NoError(t, err)

	assert.Empty(t, plan.ToolNames)
	assert.Equal(t, "fail_closed", plan.Mode)
}

func TestBuild_DefaultBundleMustExist(t *testing.T) {
	_, err := Build(&model.RouterConfig{
		Enabled:       true,
		DefaultBundle: "ghost",
		Workflows:     []model.WorkflowConfig{{Name: "real", Tools: []string{"a"}}},
	}, NewSessionPlanStoreWithTTL(time.Minute))
	assert.Error(t, err)
}

func TestBuild_DefaultBundleWithoutWorkflowsFails(t *testing.T) {
	_, err := Build(&model.RouterConfig{
		Enabled:       true,
		DefaultBundle: "ghost",
	}, NewSessionPlanStoreWithTTL(time.Minute))
	assert.Error(t, err)
}

func TestComposite_PolicyWorkflowPipeline(t *testing.T) {
	cfg := &model.RouterConfig{
		Enabled: true,
		Policy: model.PolicyConfig{Rules: []model.PolicyRule{
			{Name: "no-admin", DenyTags: []string{"admin"}},
		}},
		Workflows: []model.WorkflowConfig{
			{Name: "support", Tools: []string{"search", "ticket", "admin_tool"}, When: model.PolicyMatch{Claim: "role", Equals: "support"}},
		},
	}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	tools := []model.ToolConfig{
		toolWithMeta("search", &model.ToolMeta{Tags: []string{"read"}}),
		toolWithMeta("ticket", &model.ToolMeta{Tags: []string{"write"}}),
		toolWithMeta("admin_tool", &model.ToolMeta{Tags: []string{"admin"}}),
		toolWithMeta("unrelated", &model.ToolMeta{Tags: []string{"read"}}),
	}

	plan, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1", Claims: map[string]any{"role": "support"}}, tools)
	require.NoError(t, err)

	// Policy drops admin_tool; workflow keeps only support bundle members.
	assert.ElementsMatch(t, []string{"search", "ticket"}, plan.ToolNames)
}

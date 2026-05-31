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

func testRouterConfig() *model.RouterConfig {
	return &model.RouterConfig{Enabled: true, Fallback: FallbackFailClosed}
}

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
		Enabled:  true,
		Fallback: FallbackFailClosed,
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
	cfg := testRouterConfig()
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	tools := testTools("a", "b")
	p1, _ := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, tools)
	p2, _ := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, tools)

	// Same version => same cached plan instance returned.
	assert.Same(t, p1, p2)
}

func TestComposite_PlanRecomputesOnRegistryChange(t *testing.T) {
	cfg := testRouterConfig()
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	p1, _ := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, testTools("a", "b"))
	p2, _ := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, testTools("a", "b", "c"))

	assert.NotEqual(t, p1.Version, p2.Version)
	assert.Len(t, p2.ToolNames, 3)
}

func TestComposite_EnforceOnCallDeniesWithoutPlan(t *testing.T) {
	cfg := testRouterConfig()
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	err := cs.AuthorizeCall(context.Background(), SelectionContext{SessionID: "no-plan", Requested: "x"}, testTools("x"))
	assert.ErrorIs(t, err, ErrToolNotAuthorized)
}

func TestComposite_EnforceOnCallAllowsInPlan(t *testing.T) {
	cfg := testRouterConfig()
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	tools := testTools("a", "b")
	cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, tools)

	assert.NoError(t, cs.AuthorizeCall(context.Background(), SelectionContext{SessionID: "s1", Requested: "a"}, tools))
	assert.Equal(t, 0, store.CallCount("s1"), "authorization alone must not advance progressive disclosure")
	assert.ErrorIs(t, cs.AuthorizeCall(context.Background(), SelectionContext{SessionID: "s1", Requested: "ghost"}, tools), ErrToolNotAuthorized)
}

func TestComposite_HiddenDiscoveryStillAuthorized(t *testing.T) {
	cfg := testRouterConfig()
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	hidden := false
	tools := []model.ToolConfig{
		{Name: "visible", Cluster: "test"},
		{Name: "hidden", Cluster: "test", Meta: &model.ToolMeta{DiscoveryVisibility: &hidden}},
	}

	plan, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, tools)
	require.NoError(t, err)

	assert.ElementsMatch(t, []string{"visible", "hidden"}, plan.ToolNames)
	assert.Equal(t, []string{"visible"}, plan.VisibleToolNames)
	assert.NoError(t, cs.AuthorizeCall(context.Background(), SelectionContext{SessionID: "s1", Requested: "hidden"}, tools))
}

func TestComposite_RecordCallSuccessIncrementsCallCount(t *testing.T) {
	cfg := testRouterConfig()
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	tools := testTools("a")
	_, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, tools)
	require.NoError(t, err)

	require.NoError(t, cs.RecordCallSuccess(context.Background(), SelectionContext{SessionID: "s1", Requested: "a"}))

	assert.Equal(t, 1, store.CallCount("s1"))
}

func TestComposite_RecordCallSuccessSkipsWithoutPlanOrOutsidePlan(t *testing.T) {
	cfg := testRouterConfig()
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	require.NoError(t, cs.RecordCallSuccess(context.Background(), SelectionContext{SessionID: "missing", Requested: "a"}))
	assert.Equal(t, 0, store.CallCount("missing"))

	require.NoError(t, cs.OnInitialize(context.Background(), SelectionContext{SessionID: "init-only", AgentID: "agent"}, nil))
	require.NoError(t, cs.RecordCallSuccess(context.Background(), SelectionContext{SessionID: "init-only", Requested: "a"}))
	assert.Equal(t, 0, store.CallCount("init-only"))

	_, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, testTools("a"))
	require.NoError(t, err)
	require.NoError(t, cs.RecordCallSuccess(context.Background(), SelectionContext{SessionID: "s1", Requested: "ghost"}))
	assert.Equal(t, 0, store.CallCount("s1"))
}

func TestComposite_EnforceOnCallDisabledAllowsAll(t *testing.T) {
	disabled := false
	cfg := testRouterConfig()
	cfg.EnforceOnCall = &disabled
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	assert.NoError(t, cs.AuthorizeCall(context.Background(), SelectionContext{SessionID: "anything", Requested: "x"}, testTools("x")))
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
	assert.Equal(t, ModeFallbackBundle, plan.Mode)
}

func TestComposite_FallbackBundleRespectsPolicyDeniedTools(t *testing.T) {
	cfg := &model.RouterConfig{
		Enabled:       true,
		DefaultBundle: "safe-minimal",
		Policy: model.PolicyConfig{Rules: []model.PolicyRule{
			{Name: "block-admin", DenyTags: []string{"admin"}},
		}},
		Workflows: []model.WorkflowConfig{
			{Name: "safe-minimal", Tools: []string{"ping", "admin_tool"}},
			{Name: "empty-flow", Tools: []string{"missing"}, When: model.PolicyMatch{Claim: "role", Equals: "empty"}},
		},
	}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	plan, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1", Claims: map[string]any{"role": "empty"}},
		[]model.ToolConfig{
			toolWithMeta("ping", &model.ToolMeta{Tags: []string{"safe"}}),
			toolWithMeta("admin_tool", &model.ToolMeta{Tags: []string{"admin"}}),
		})
	require.NoError(t, err)

	assert.Equal(t, []string{"ping"}, plan.ToolNames)
	assert.NotContains(t, plan.ToolNames, "admin_tool")
	assert.Equal(t, ModeFallbackBundle, plan.Mode)
}

func TestComposite_FallbackBundleWorksWhenWorkflowStageDisabled(t *testing.T) {
	workflowDisabled := false
	cfg := &model.RouterConfig{
		Enabled:       true,
		DefaultBundle: "safe-minimal",
		Stages: model.RouterStages{
			Workflow:    &workflowDisabled,
			Progressive: true,
		},
		Progressive: model.ProgressiveConfig{InitialBundle: "locked-empty"},
		Workflows: []model.WorkflowConfig{
			{Name: "safe-minimal", Tools: []string{"ping"}},
			{Name: "locked-empty", Tools: []string{"missing"}},
		},
	}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	plan, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1"},
		[]model.ToolConfig{toolWithMeta("ping", nil)})
	require.NoError(t, err)

	assert.Equal(t, []string{"ping"}, plan.ToolNames)
	assert.Equal(t, ModeFallbackBundle, plan.Mode)
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
	assert.Empty(t, plan.VisibleToolNames)
	assert.Equal(t, ModeFailClosed, plan.Mode)
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

func TestBuild_BundleDefaultRequiresDefaultBundle(t *testing.T) {
	_, err := Build(&model.RouterConfig{
		Enabled: true,
	}, NewSessionPlanStoreWithTTL(time.Minute))
	assert.ErrorContains(t, err, "default_bundle is required")
}

func TestBuild_DefaultBundleEmptyWorkflowFails(t *testing.T) {
	_, err := Build(&model.RouterConfig{
		Enabled:       true,
		DefaultBundle: "empty",
		Workflows:     []model.WorkflowConfig{{Name: "empty"}},
	}, NewSessionPlanStoreWithTTL(time.Minute))
	assert.ErrorContains(t, err, "empty workflow")
}

func TestBuild_FailClosedDoesNotRequireDefaultBundle(t *testing.T) {
	sel, err := Build(&model.RouterConfig{
		Enabled:  true,
		Fallback: FallbackFailClosed,
	}, NewSessionPlanStoreWithTTL(time.Minute))
	assert.NoError(t, err)
	assert.NotNil(t, sel)
}

func TestBuild_ProgressiveRequiresInitialBundle(t *testing.T) {
	_, err := Build(&model.RouterConfig{
		Enabled: true,
		Stages:  model.RouterStages{Progressive: true},
	}, NewSessionPlanStoreWithTTL(time.Minute))

	assert.ErrorContains(t, err, "progressive.initial_bundle is required")
}

func TestBuild_DuplicateWorkflowNameFails(t *testing.T) {
	_, err := Build(&model.RouterConfig{
		Enabled: true,
		Workflows: []model.WorkflowConfig{
			{Name: "support", Tools: []string{"search"}},
			{Name: "support", Tools: []string{"ticket"}},
		},
	}, NewSessionPlanStoreWithTTL(time.Minute))

	assert.ErrorContains(t, err, "defined more than once")
}

func TestBuild_UnknownFallbackFails(t *testing.T) {
	_, err := Build(&model.RouterConfig{
		Enabled:  true,
		Fallback: "fail_open",
	}, NewSessionPlanStoreWithTTL(time.Minute))
	assert.Error(t, err)
}

func TestBuild_InvalidSampleRateFails(t *testing.T) {
	_, err := Build(&model.RouterConfig{
		Enabled:  true,
		Fallback: FallbackFailClosed,
		Audit:    model.AuditConfig{SampleRate: 1.5},
	}, NewSessionPlanStoreWithTTL(time.Minute))
	assert.ErrorContains(t, err, "sample_rate")

	_, err = Build(&model.RouterConfig{
		Enabled:  true,
		Fallback: FallbackFailClosed,
		Audit:    model.AuditConfig{SampleRate: -0.1},
	}, NewSessionPlanStoreWithTTL(time.Minute))
	assert.ErrorContains(t, err, "sample_rate")
}

func TestBuild_InvalidMaxRiskFailsEvenWhenPolicyStageDisabled(t *testing.T) {
	disabled := false
	_, err := Build(&model.RouterConfig{
		Enabled: true,
		Stages:  model.RouterStages{Policy: &disabled},
		Policy: model.PolicyConfig{Rules: []model.PolicyRule{
			{Name: "bad-risk", MaxRisk: "hihg"},
		}},
	}, NewSessionPlanStoreWithTTL(time.Minute))
	assert.ErrorContains(t, err, "max_risk")
	assert.ErrorContains(t, err, "unsupported risk")
}

func TestBuild_EnabledRequiresStore(t *testing.T) {
	_, err := Build(&model.RouterConfig{Enabled: true}, nil)
	assert.Error(t, err)
}

func TestBuild_DisabledDoesNotRequireStore(t *testing.T) {
	sel, err := Build(&model.RouterConfig{Enabled: false}, nil)
	assert.NoError(t, err)
	assert.Nil(t, sel)
}

func TestComposite_PolicyWorkflowPipeline(t *testing.T) {
	cfg := &model.RouterConfig{
		Enabled:  true,
		Fallback: FallbackFailClosed,
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

func TestComposite_AuthorizeCallRecomputesOnClaimChange(t *testing.T) {
	cfg := &model.RouterConfig{
		Enabled:  true,
		Fallback: FallbackFailClosed,
		Workflows: []model.WorkflowConfig{
			{Name: "support", Tools: []string{"support_search"}, When: model.PolicyMatch{Claim: "role", Equals: "support"}},
			{Name: "billing", Tools: []string{"billing_search"}, When: model.PolicyMatch{Claim: "role", Equals: "billing"}},
		},
	}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	tools := testTools("support_search", "billing_search")
	supportCtx := SelectionContext{SessionID: "s1", Claims: map[string]any{"role": "support"}}
	plan, err := cs.Select(context.Background(), supportCtx, tools)
	require.NoError(t, err)
	assert.Equal(t, []string{"support_search"}, plan.ToolNames)

	billingCtx := SelectionContext{SessionID: "s1", Claims: map[string]any{"role": "billing"}, Requested: "support_search"}
	err = cs.AuthorizeCall(context.Background(), billingCtx, tools)
	assert.ErrorIs(t, err, ErrToolNotAuthorized)

	got, ok := store.Get("s1")
	require.True(t, ok)
	assert.Equal(t, []string{"billing_search"}, got.ToolNames)
}

func TestComposite_AuthorizeCallRecomputesOnToolMetadataChange(t *testing.T) {
	cfg := &model.RouterConfig{
		Enabled:  true,
		Fallback: FallbackFailClosed,
		Policy: model.PolicyConfig{Rules: []model.PolicyRule{
			{Name: "safe-only", AllowTags: []string{"safe"}},
		}},
	}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	safeTools := []model.ToolConfig{toolWithMeta("export_data", &model.ToolMeta{Tags: []string{"safe"}})}
	_, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, safeTools)
	require.NoError(t, err)

	riskyTools := []model.ToolConfig{toolWithMeta("export_data", &model.ToolMeta{Tags: []string{"risky"}})}
	err = cs.AuthorizeCall(context.Background(), SelectionContext{SessionID: "s1", Requested: "export_data"}, riskyTools)
	assert.ErrorIs(t, err, ErrToolNotAuthorized)
}

func TestComposite_AuthorizeCallRecomputesOnToolDefinitionChange(t *testing.T) {
	cfg := testRouterConfig()
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	initialTools := testTools("export_data")
	_, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, initialTools)
	require.NoError(t, err)
	before, ok := store.Get("s1")
	require.True(t, ok)

	updatedTools := testTools("export_data")
	updatedTools[0].Cluster = "new-cluster"
	updatedTools[0].Request.Path = "/api/v2/export"
	updatedTools[0].Args = append(updatedTools[0].Args, model.ArgConfig{Name: "format", Type: "string", In: "query"})
	err = cs.AuthorizeCall(context.Background(), SelectionContext{SessionID: "s1", Requested: "export_data"}, updatedTools)
	require.NoError(t, err)

	after, ok := store.Get("s1")
	require.True(t, ok)
	assert.NotEqual(t, before.Version, after.Version)
}

func TestComposite_AuthorizeCallAllowsAfterStaleRecompute(t *testing.T) {
	cfg := testRouterConfig()
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	initialTools := testTools("a")
	_, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, initialTools)
	require.NoError(t, err)

	updatedTools := testTools("a", "b")
	err = cs.AuthorizeCall(context.Background(), SelectionContext{SessionID: "s1", Requested: "a"}, updatedTools)
	assert.NoError(t, err)

	got, ok := store.Get("s1")
	require.True(t, ok)
	assert.ElementsMatch(t, []string{"a", "b"}, got.ToolNames)
}

func TestComposite_IdentityIsolatedCache(t *testing.T) {
	cfg := &model.RouterConfig{
		Enabled:  true,
		Fallback: FallbackFailClosed,
		Policy: model.PolicyConfig{Rules: []model.PolicyRule{
			{Name: "acme", When: model.PolicyMatch{Claim: "tenant", Equals: "acme"}, AllowTags: []string{"acme"}},
			{Name: "globex", When: model.PolicyMatch{Claim: "tenant", Equals: "globex"}, AllowTags: []string{"globex"}},
		}},
	}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	tools := []model.ToolConfig{
		toolWithMeta("acme_tool", &model.ToolMeta{Tags: []string{"acme"}}),
		toolWithMeta("globex_tool", &model.ToolMeta{Tags: []string{"globex"}}),
	}

	// Same session ID, different tenants should get different plans.
	p1, _ := cs.Select(context.Background(), SelectionContext{SessionID: "shared", Tenant: "acme"}, tools)
	assert.Equal(t, []string{"acme_tool"}, p1.ToolNames)

	p2, _ := cs.Select(context.Background(), SelectionContext{SessionID: "shared", Tenant: "globex"}, tools)
	assert.Equal(t, []string{"globex_tool"}, p2.ToolNames, "different tenant should not reuse cached plan")

	// Verify versions differ due to identity fingerprint.
	assert.NotEqual(t, p1.Version, p2.Version)
}

func TestComposite_AgentIDPropagatedToLogs(t *testing.T) {
	cfg := testRouterConfig()
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	// Simulate initialize storing agentID (happens before tools/list in real flow).
	err := cs.OnInitialize(context.Background(), SelectionContext{SessionID: "s1", AgentID: "test-agent"}, nil)
	require.NoError(t, err)

	// tools/list should backfill agentID from store for decision logs.
	tools := testTools("a", "b")
	plan, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, tools)
	require.NoError(t, err)
	assert.NotNil(t, plan)

	// Verify agentID was stored and survives the Select call.
	assert.Equal(t, "test-agent", store.AgentID("s1"))
}

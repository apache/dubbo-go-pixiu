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
	"math"
	"sync"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func testRouterConfig() *model.RouterConfig {
	return &model.RouterConfig{Fallback: FallbackFailClosed}
}

// buildComposite is a test helper that wires a CompositeSelector from a config.
func buildComposite(t *testing.T, cfg *model.RouterConfig) (*CompositeSelector, *SessionPlanStore) {
	store := NewSessionPlanStore()
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

	// Same version => cached content is reused, but callers receive isolated copies.
	assert.Equal(t, p1.Version, p2.Version)
	assert.Equal(t, p1.ToolNames, p2.ToolNames)
	assert.NotSame(t, p1, p2)
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

	err := authorizeCall(cs, context.Background(), SelectionContext{SessionID: "no-plan", Requested: "x"}, testTools("x"))
	assert.ErrorIs(t, err, ErrToolNotAuthorized)
}

func TestComposite_EnforceOnCallAllowsInPlan(t *testing.T) {
	cfg := testRouterConfig()
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	tools := testTools("a", "b")
	cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, tools)

	assert.NoError(t, authorizeCall(cs, context.Background(), SelectionContext{SessionID: "s1", Requested: "a"}, tools))
	assert.Equal(t, int64(0), store.CallCount(cs.planKey("s1")), "authorization alone must not advance progressive disclosure")
	assert.ErrorIs(t, authorizeCall(cs, context.Background(), SelectionContext{SessionID: "s1", Requested: "ghost"}, tools), ErrToolNotAuthorized)
}

func TestComposite_AuthorizeCallRecomputeReceiptBindsCommittedPlan(t *testing.T) {
	cfg := &model.RouterConfig{
		Fallback: FallbackFailClosed,
		Policy: model.PolicyConfig{Rules: []model.PolicyRule{
			{Name: "acme", When: model.PolicyMatch{Claim: "tenant", Equals: "acme"}, AllowTags: []string{"acme"}},
			{Name: "globex", When: model.PolicyMatch{Claim: "tenant", Equals: "globex"}, AllowTags: []string{"globex"}},
		}},
	}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	tools := []model.ToolConfig{
		toolWithMeta("acme_read", &model.ToolMeta{Tags: []string{"acme"}}),
		toolWithMeta("globex_read", &model.ToolMeta{Tags: []string{"globex"}}),
	}

	_, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1", Tenant: "globex", Claims: map[string]any{"tenant": "globex"}}, tools)
	require.NoError(t, err)

	acmeReceipt, err := cs.AuthorizeCall(context.Background(), SelectionContext{
		SessionID: "s1",
		Tenant:    "acme",
		Claims:    map[string]any{"tenant": "acme"},
		Requested: "acme_read",
	}, tools)
	require.NoError(t, err)
	require.NotNil(t, acmeReceipt)
	assert.Equal(t, "acme_read", acmeReceipt.ToolName)

	_, err = cs.Select(context.Background(), SelectionContext{SessionID: "s1", Tenant: "globex", Claims: map[string]any{"tenant": "globex"}}, tools)
	require.NoError(t, err)

	got, ok := store.Get(cs.planKey("s1"))
	require.True(t, ok)
	assert.Equal(t, []string{"globex_read"}, got.ToolNames)
	assert.NotEqual(t, got.IdentityHash, acmeReceipt.IdentityHash)
	assert.Equal(t, CallSuccessResult{}, store.FinalizeReceipt(*acmeReceipt, ReceiptSucceeded, 1), "receipt from acme generation must not advance the later globex plan")
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
	assert.NoError(t, authorizeCall(cs, context.Background(), SelectionContext{SessionID: "s1", Requested: "hidden"}, tools))
}

func TestComposite_RecordCallSuccessIncrementsCallCount(t *testing.T) {
	cfg := &model.RouterConfig{
		Fallback: FallbackFailClosed,
		Stages:   model.RouterStages{Progressive: true},
		Workflows: []model.WorkflowConfig{
			{Name: "starter", Tools: []string{"a"}},
		},
		Progressive: model.ProgressiveConfig{InitialBundle: "starter", ExpandAfterCalls: 2},
	}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	tools := testTools("a")
	_, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, tools)
	require.NoError(t, err)

	result, err := recordCallSuccessForTest(t, cs, SelectionContext{SessionID: "s1", Requested: "a"})
	require.NoError(t, err)

	assert.Equal(t, int64(1), result.Count)
	assert.False(t, result.Transitioned)
	assert.Equal(t, int64(1), store.CallCount(cs.planKey("s1")))
}

func TestComposite_RecordCallSuccessSkipsWithoutPlanOrOutsidePlan(t *testing.T) {
	cfg := &model.RouterConfig{
		Fallback: FallbackFailClosed,
		Stages:   model.RouterStages{Progressive: true},
		Workflows: []model.WorkflowConfig{
			{Name: "starter", Tools: []string{"a"}},
		},
		Progressive: model.ProgressiveConfig{InitialBundle: "starter", ExpandAfterCalls: 1},
	}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	result, err := recordCallSuccessForTest(t, cs, SelectionContext{SessionID: "missing", Requested: "a"})
	require.NoError(t, err)
	assert.False(t, result.Transitioned)
	assert.Equal(t, int64(0), store.CallCount(cs.planKey("missing")))

	result, err = recordCallSuccessForTest(t, cs, SelectionContext{SessionID: "init-only", Requested: "a"})
	require.NoError(t, err)
	assert.False(t, result.Transitioned)
	assert.Equal(t, int64(0), store.CallCount(cs.planKey("init-only")))

	_, err = cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, testTools("a"))
	require.NoError(t, err)
	result, err = recordCallSuccessForTest(t, cs, SelectionContext{SessionID: "s1", Requested: "ghost"})
	require.NoError(t, err)
	assert.False(t, result.Transitioned)
	assert.Equal(t, int64(0), store.CallCount(cs.planKey("s1")))
}

func TestComposite_FallbackBundleDefault(t *testing.T) {
	cfg := &model.RouterConfig{
		Fallback:      FallbackBundleDefault,
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
	// candidate set. An explicit empty workflow selection must not fall back.
	plan, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1", Tenant: "acme"},
		[]model.ToolConfig{toolWithMeta("ping", nil)})
	require.NoError(t, err)

	assert.Empty(t, plan.ToolNames)
	assert.Equal(t, ModeSelected, plan.Mode)
	assert.Equal(t, SelectionOutcomeEmptyByConfiguration, plan.Outcome)
}

func TestComposite_NoWorkflowMatchUsesDefaultBundle(t *testing.T) {
	cfg := &model.RouterConfig{
		Fallback:      FallbackBundleDefault,
		DefaultBundle: "safe-minimal",
		Workflows: []model.WorkflowConfig{
			{Name: "safe-minimal", Tools: []string{"ping"}},
			{Name: "support", Tools: []string{"ticket"}, When: model.PolicyMatch{Claim: "role", Equals: "support"}},
		},
	}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	plan, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1", Claims: map[string]any{"role": "unknown"}},
		[]model.ToolConfig{toolWithMeta("ping", nil), toolWithMeta("ticket", nil)})
	require.NoError(t, err)

	assert.Equal(t, []string{"ping"}, plan.ToolNames)
	assert.Equal(t, ModeFallbackBundle, plan.Mode)
	assert.Equal(t, SelectionOutcomeNoMatch, plan.Outcome)
}

func TestComposite_ExplicitPolicyDenyDoesNotUseDefaultBundle(t *testing.T) {
	cfg := &model.RouterConfig{
		Fallback:      FallbackBundleDefault,
		DefaultBundle: "safe-minimal",
		Policy: model.PolicyConfig{Rules: []model.PolicyRule{
			{Name: "block-admin", DenyTags: []string{"admin"}},
		}},
		Workflows: []model.WorkflowConfig{
			{Name: "safe-minimal", Tools: []string{"admin_tool"}},
		},
	}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	plan, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1"},
		[]model.ToolConfig{toolWithMeta("admin_tool", &model.ToolMeta{Tags: []string{"admin"}})})
	require.NoError(t, err)

	assert.Empty(t, plan.ToolNames)
	assert.Equal(t, ModeSelected, plan.Mode)
	assert.Equal(t, SelectionOutcomeExplicitDeny, plan.Outcome)
}

func TestComposite_SelectionFailureBundleDefaultRespectsPolicy(t *testing.T) {
	cfg := &model.RouterConfig{
		Fallback:      FallbackBundleDefault,
		DefaultBundle: "safe-minimal",
		Policy: model.PolicyConfig{Rules: []model.PolicyRule{
			{Name: "block-admin", DenyTags: []string{"admin"}},
		}},
		Workflows: []model.WorkflowConfig{
			{Name: "safe-minimal", Tools: []string{"ping", "admin_tool"}},
		},
	}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	plan, err := cs.HandleSelectionFailure(context.Background(), SelectionContext{SessionID: "s1"},
		[]model.ToolConfig{
			toolWithMeta("ping", &model.ToolMeta{Tags: []string{"safe"}}),
			toolWithMeta("admin_tool", &model.ToolMeta{Tags: []string{"admin"}}),
		}, assert.AnError)
	require.NoError(t, err)

	assert.Equal(t, []string{"ping"}, plan.ToolNames)
	assert.Equal(t, ModeFallbackBundle, plan.Mode)
	assert.Equal(t, SelectionOutcomeInternalError, plan.Outcome)
}

func TestComposite_FallbackBundleRespectsPolicyDeniedTools(t *testing.T) {
	cfg := &model.RouterConfig{
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

	assert.Empty(t, plan.ToolNames)
	assert.NotContains(t, plan.ToolNames, "admin_tool")
	assert.Equal(t, ModeSelected, plan.Mode)
	assert.Equal(t, SelectionOutcomeEmptyByConfiguration, plan.Outcome)
}

func TestComposite_FallbackBundleWorksWhenWorkflowStageDisabled(t *testing.T) {
	workflowDisabled := false
	cfg := &model.RouterConfig{
		DefaultBundle: "safe-minimal",
		Stages: model.RouterStages{
			Workflow:    &workflowDisabled,
			Progressive: true,
		},
		Progressive: model.ProgressiveConfig{InitialBundle: "locked-empty", ExpandAfterCalls: 1},
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

	assert.Empty(t, plan.ToolNames)
	assert.Equal(t, ModeSelected, plan.Mode)
	assert.Equal(t, SelectionOutcomeEmptyByConfiguration, plan.Outcome)
}

func TestComposite_FallbackFailClosed(t *testing.T) {
	cfg := &model.RouterConfig{
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
	assert.Equal(t, ModeSelected, plan.Mode)
	assert.Equal(t, SelectionOutcomeEmptyByConfiguration, plan.Outcome)
}

func TestBuild_DefaultBundleMustExist(t *testing.T) {
	_, err := Build(&model.RouterConfig{
		Fallback:      FallbackBundleDefault,
		DefaultBundle: "ghost",
		Workflows:     []model.WorkflowConfig{{Name: "real", Tools: []string{"a"}}},
	}, NewSessionPlanStore())
	assert.Error(t, err)
}

func TestBuild_DefaultBundleWithoutWorkflowsFails(t *testing.T) {
	_, err := Build(&model.RouterConfig{
		Fallback:      FallbackBundleDefault,
		DefaultBundle: "ghost",
	}, NewSessionPlanStore())
	assert.Error(t, err)
}

func TestBuild_BundleDefaultRequiresDefaultBundle(t *testing.T) {
	_, err := Build(&model.RouterConfig{Fallback: FallbackBundleDefault}, NewSessionPlanStore())
	assert.ErrorContains(t, err, "default_bundle is required")
}

func TestBuild_DefaultBundleEmptyWorkflowFails(t *testing.T) {
	_, err := Build(&model.RouterConfig{
		Fallback:      FallbackBundleDefault,
		DefaultBundle: "empty",
		Workflows:     []model.WorkflowConfig{{Name: "empty"}},
	}, NewSessionPlanStore())
	assert.ErrorContains(t, err, "empty workflow")
}

func TestBuild_FailClosedDoesNotRequireDefaultBundle(t *testing.T) {
	sel, err := Build(&model.RouterConfig{
		Fallback: FallbackFailClosed,
	}, NewSessionPlanStore())
	assert.NoError(t, err)
	assert.NotNil(t, sel)
}

func TestBuild_EmptyRouterDefaultsToFailClosed(t *testing.T) {
	sel, err := Build(&model.RouterConfig{}, NewSessionPlanStore())
	require.NoError(t, err)
	cs := sel.(*CompositeSelector)
	assert.Equal(t, FallbackFailClosed, cs.fallback)

	plan, err := cs.HandleSelectionFailure(context.Background(), SelectionContext{SessionID: "s1"}, []model.ToolConfig{toolWithMeta("ping", nil)}, assert.AnError)
	require.NoError(t, err)
	assert.Equal(t, ModeFailClosed, plan.Mode)
	assert.Empty(t, plan.ToolNames)
}

func TestNewCompositeSelectorDefaultsAndValidatesFallback(t *testing.T) {
	store := NewSessionPlanStore()
	defer store.Stop()

	cs, err := NewCompositeSelector(CompositeOptions{Store: store})
	require.NoError(t, err)
	assert.Equal(t, FallbackFailClosed, cs.fallback)

	_, err = NewCompositeSelector(CompositeOptions{Store: store, Fallback: "fail_open"})
	assert.ErrorContains(t, err, "not supported")
}

func TestBuild_ProgressiveRequiresInitialBundle(t *testing.T) {
	_, err := Build(&model.RouterConfig{
		Stages: model.RouterStages{Progressive: true},
		Progressive: model.ProgressiveConfig{
			ExpandAfterCalls: 1,
		},
	}, NewSessionPlanStore())

	assert.ErrorContains(t, err, "progressive.initial_bundle is required")
}

func TestBuild_DuplicateWorkflowNameFails(t *testing.T) {
	_, err := Build(&model.RouterConfig{
		Workflows: []model.WorkflowConfig{
			{Name: "support", Tools: []string{"search"}},
			{Name: "support", Tools: []string{"ticket"}},
		},
	}, NewSessionPlanStore())

	assert.ErrorContains(t, err, "defined more than once")
}

func TestBuild_UnknownFallbackFails(t *testing.T) {
	_, err := Build(&model.RouterConfig{
		Fallback: "fail_open",
	}, NewSessionPlanStore())
	assert.Error(t, err)
}

func TestBuild_InvalidSampleRateFails(t *testing.T) {
	_, err := Build(&model.RouterConfig{
		Fallback: FallbackFailClosed,
		Audit:    model.AuditConfig{SampleRate: 1.5},
	}, NewSessionPlanStore())
	assert.ErrorContains(t, err, "sample_rate")

	_, err = Build(&model.RouterConfig{
		Fallback: FallbackFailClosed,
		Audit:    model.AuditConfig{SampleRate: -0.1},
	}, NewSessionPlanStore())
	assert.ErrorContains(t, err, "sample_rate")

	_, err = Build(&model.RouterConfig{
		Fallback: FallbackFailClosed,
		Audit:    model.AuditConfig{SampleRate: math.NaN()},
	}, NewSessionPlanStore())
	assert.ErrorContains(t, err, "sample_rate")

	_, err = Build(&model.RouterConfig{
		Fallback: FallbackFailClosed,
		Audit:    model.AuditConfig{SampleRate: math.Inf(1)},
	}, NewSessionPlanStore())
	assert.ErrorContains(t, err, "sample_rate")
}

func TestBuild_InvalidSessionMaxEntriesFails(t *testing.T) {
	_, err := Build(&model.RouterConfig{
		Fallback: FallbackFailClosed,
		Session:  model.RouterSessionConfig{MaxEntries: -1},
	}, NewSessionPlanStore())
	assert.ErrorContains(t, err, "max_entries")
}

func TestBuild_ConfigHashIgnoresNonAuthorizationFields(t *testing.T) {
	base := &model.RouterConfig{
		Fallback: FallbackFailClosed,
		Policy: model.PolicyConfig{Rules: []model.PolicyRule{
			{Name: "readers", AllowTags: []string{"read"}, MaxRisk: "low"},
		}},
	}
	withOpsOnly := base.DeepCopy()
	withOpsOnly.Audit = model.AuditConfig{SampleRate: 1, DecisionDetailLogging: true}
	withOpsOnly.Session = model.RouterSessionConfig{MaxEntries: 7}

	baseHash, err := configHash(normalizeConfig(base))
	require.NoError(t, err)
	opsHash, err := configHash(normalizeConfig(withOpsOnly))
	require.NoError(t, err)
	assert.Equal(t, baseHash, opsHash)

	authChanged := base.DeepCopy()
	authChanged.Policy.Rules[0].MaxRisk = "medium"
	authHash, err := configHash(normalizeConfig(authChanged))
	require.NoError(t, err)
	assert.NotEqual(t, baseHash, authHash)
}

func TestBuild_ProgressiveRejectsNonPositiveExpandAfter(t *testing.T) {
	_, err := Build(&model.RouterConfig{
		Fallback: FallbackFailClosed,
		Stages:   model.RouterStages{Progressive: true},
		Workflows: []model.WorkflowConfig{
			{Name: "starter", Tools: []string{"t0"}},
		},
		Progressive: model.ProgressiveConfig{InitialBundle: "starter"},
	}, NewSessionPlanStore())
	assert.ErrorContains(t, err, "expand_after_calls")
}

func TestBuild_InvalidMaxRiskFailsEvenWhenPolicyStageOff(t *testing.T) {
	off := false
	_, err := Build(&model.RouterConfig{
		Stages: model.RouterStages{Policy: &off},
		Policy: model.PolicyConfig{Rules: []model.PolicyRule{
			{Name: "bad-risk", MaxRisk: "hihg"},
		}},
	}, NewSessionPlanStore())
	assert.ErrorContains(t, err, "max_risk")
	assert.ErrorContains(t, err, "unsupported risk")
}

func TestBuild_DefaultGovernanceRequiresStore(t *testing.T) {
	_, err := Build(&model.RouterConfig{}, nil)
	assert.Error(t, err)
}

func TestBuild_DefaultConfigRequiresStore(t *testing.T) {
	sel, err := Build(&model.RouterConfig{}, NewSessionPlanStore())
	assert.NoError(t, err)
	assert.NotNil(t, sel)
}

func TestComposite_PolicyWorkflowPipeline(t *testing.T) {
	cfg := &model.RouterConfig{
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
	err = authorizeCall(cs, context.Background(), billingCtx, tools)
	assert.ErrorIs(t, err, ErrToolNotAuthorized)

	got, ok := store.Get(cs.planKey("s1"))
	require.True(t, ok)
	assert.Equal(t, []string{"billing_search"}, got.ToolNames)
}

func TestComposite_AuthorizeCallRecomputesOnToolMetadataChange(t *testing.T) {
	cfg := &model.RouterConfig{
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
	err = authorizeCall(cs, context.Background(), SelectionContext{SessionID: "s1", Requested: "export_data"}, riskyTools)
	assert.ErrorIs(t, err, ErrToolNotAuthorized)
}

func TestComposite_AuthorizeCallRecomputesOnToolDefinitionChange(t *testing.T) {
	cfg := testRouterConfig()
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	initialTools := testTools("export_data")
	_, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, initialTools)
	require.NoError(t, err)
	before, ok := store.Get(cs.planKey("s1"))
	require.True(t, ok)

	updatedTools := testTools("export_data")
	updatedTools[0].Cluster = "new-cluster"
	updatedTools[0].Request.Path = "/api/v2/export"
	updatedTools[0].Args = append(updatedTools[0].Args, model.ArgConfig{Name: "format", Type: "string", In: "query"})
	err = authorizeCall(cs, context.Background(), SelectionContext{SessionID: "s1", Requested: "export_data"}, updatedTools)
	require.NoError(t, err)

	after, ok := store.Get(cs.planKey("s1"))
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
	err = authorizeCall(cs, context.Background(), SelectionContext{SessionID: "s1", Requested: "a"}, updatedTools)
	assert.NoError(t, err)

	got, ok := store.Get(cs.planKey("s1"))
	require.True(t, ok)
	assert.ElementsMatch(t, []string{"a", "b"}, got.ToolNames)
}

func TestComposite_IdentityIsolatedCache(t *testing.T) {
	cfg := &model.RouterConfig{
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

func TestComposite_InitializeDoesNotPersistRawIdentity(t *testing.T) {
	cfg := &model.RouterConfig{
		Fallback: FallbackFailClosed,
		Stages:   model.RouterStages{Progressive: true},
		Workflows: []model.WorkflowConfig{
			{Name: "starter", Tools: []string{"a"}},
		},
		Progressive: model.ProgressiveConfig{InitialBundle: "starter", ExpandAfterCalls: 1},
	}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	assert.Equal(t, 0, store.Len())

	tools := testTools("a", "b")
	plan, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, tools)
	require.NoError(t, err)
	assert.NotNil(t, plan)
}

func TestComposite_ProgressiveBoundaryAndPolicyAfterExpansion(t *testing.T) {
	cfg := &model.RouterConfig{
		Fallback: FallbackFailClosed,
		Stages:   model.RouterStages{Progressive: true},
		Policy: model.PolicyConfig{Rules: []model.PolicyRule{
			{Name: "no-admin", DenyTags: []string{"admin"}},
		}},
		Workflows: []model.WorkflowConfig{
			{Name: "starter", Tools: []string{"ping"}},
		},
		Progressive: model.ProgressiveConfig{InitialBundle: "starter", ExpandAfterCalls: 2},
	}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	tools := []model.ToolConfig{
		toolWithMeta("ping", &model.ToolMeta{Tags: []string{"safe"}}),
		toolWithMeta("advanced", &model.ToolMeta{Tags: []string{"safe"}}),
		toolWithMeta("admin", &model.ToolMeta{Tags: []string{"admin"}}),
	}

	sc := SelectionContext{SessionID: "s1"}
	p1, err := cs.Select(context.Background(), sc, tools)
	require.NoError(t, err)
	assert.Equal(t, []string{"ping"}, p1.ToolNames)

	r1, err := recordCallSuccessForTest(t, cs, SelectionContext{SessionID: "s1", Requested: "ping"})
	require.NoError(t, err)
	assert.Equal(t, int64(1), r1.Count)
	assert.False(t, r1.Transitioned)

	r2, err := recordCallSuccessForTest(t, cs, SelectionContext{SessionID: "s1", Requested: "ping"})
	require.NoError(t, err)
	assert.Equal(t, int64(2), r2.Count)
	assert.True(t, r2.Transitioned)

	r3, err := recordCallSuccessForTest(t, cs, SelectionContext{SessionID: "s1", Requested: "ping"})
	require.NoError(t, err)
	assert.Equal(t, int64(3), r3.Count)
	assert.False(t, r3.Transitioned)

	p2, err := cs.Select(context.Background(), sc, tools)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"ping", "advanced"}, p2.ToolNames)
	assert.NotContains(t, p2.ToolNames, "admin")
}

func TestComposite_ProgressiveConcurrentTransitionOnce(t *testing.T) {
	cfg := &model.RouterConfig{
		Fallback: FallbackFailClosed,
		Stages:   model.RouterStages{Progressive: true},
		Workflows: []model.WorkflowConfig{
			{Name: "starter", Tools: []string{"ping"}},
		},
		Progressive: model.ProgressiveConfig{InitialBundle: "starter", ExpandAfterCalls: 10},
	}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	_, err := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, testTools("ping", "advanced"))
	require.NoError(t, err)

	var wg sync.WaitGroup
	transitions := make(chan bool, 100)
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			result, err := recordCallSuccessForTest(t, cs, SelectionContext{SessionID: "s1", Requested: "ping"})
			require.NoError(t, err)
			if result.Transitioned {
				transitions <- true
			}
		}()
	}
	wg.Wait()
	close(transitions)

	assert.Len(t, transitions, 1)
}

func TestComposite_ProgressiveCallCountIsRouterScoped(t *testing.T) {
	cfg := &model.RouterConfig{
		Fallback: FallbackFailClosed,
		Stages:   model.RouterStages{Progressive: true},
		Workflows: []model.WorkflowConfig{
			{Name: "starter", Tools: []string{"ping"}},
		},
		Progressive: model.ProgressiveConfig{InitialBundle: "starter", ExpandAfterCalls: 1},
	}
	store := NewSessionPlanStore()
	defer store.Stop()

	selA, err := Build(cfg, store)
	require.NoError(t, err)
	selB, err := Build(cfg, store)
	require.NoError(t, err)
	csA := selA.(*CompositeSelector)
	csB := selB.(*CompositeSelector)

	tools := testTools("ping", "advanced")
	sc := SelectionContext{SessionID: "shared"}
	_, err = csA.Select(context.Background(), sc, tools)
	require.NoError(t, err)
	_, err = csB.Select(context.Background(), sc, tools)
	require.NoError(t, err)

	result, err := recordCallSuccessForTest(t, csA, SelectionContext{SessionID: "shared", Requested: "ping"})
	require.NoError(t, err)
	assert.True(t, result.Transitioned)
	assert.Equal(t, int64(1), store.CallCount(csA.planKey("shared")))
	assert.Equal(t, int64(0), store.CallCount(csB.planKey("shared")))

	expandedA, err := csA.Select(context.Background(), sc, tools)
	require.NoError(t, err)
	initialB, err := csB.Select(context.Background(), sc, tools)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"ping", "advanced"}, expandedA.ToolNames)
	assert.Equal(t, []string{"ping"}, initialB.ToolNames)
}

func TestComposite_ClaimsChangeResetsProgressiveState(t *testing.T) {
	cfg := &model.RouterConfig{
		Fallback: FallbackFailClosed,
		Stages:   model.RouterStages{Progressive: true},
		Policy: model.PolicyConfig{Rules: []model.PolicyRule{
			{Name: "acme", When: model.PolicyMatch{Claim: "tenant", Equals: "acme"}, AllowTags: []string{"safe"}},
			{Name: "globex", When: model.PolicyMatch{Claim: "tenant", Equals: "globex"}, AllowTags: []string{"safe"}},
		}},
		Workflows: []model.WorkflowConfig{
			{Name: "starter", Tools: []string{"ping"}},
		},
		Progressive: model.ProgressiveConfig{InitialBundle: "starter", ExpandAfterCalls: 1},
	}
	cs, store := buildComposite(t, cfg)
	defer store.Stop()

	tools := []model.ToolConfig{
		toolWithMeta("ping", &model.ToolMeta{Tags: []string{"safe"}}),
		toolWithMeta("advanced", &model.ToolMeta{Tags: []string{"safe"}}),
	}
	acme := SelectionContext{SessionID: "s1", Tenant: "acme"}
	plan, err := cs.Select(context.Background(), acme, tools)
	require.NoError(t, err)
	assert.Equal(t, []string{"ping"}, plan.ToolNames)

	result, err := recordCallSuccessForTest(t, cs, SelectionContext{SessionID: "s1", Tenant: "acme", Requested: "ping"})
	require.NoError(t, err)
	assert.True(t, result.Transitioned)
	expanded, err := cs.Select(context.Background(), acme, tools)
	require.NoError(t, err)
	assert.ElementsMatch(t, []string{"ping", "advanced"}, expanded.ToolNames)

	globex := SelectionContext{SessionID: "s1", Tenant: "globex"}
	reset, err := cs.Select(context.Background(), globex, tools)
	require.NoError(t, err)
	assert.Equal(t, []string{"ping"}, reset.ToolNames)
	assert.Equal(t, int64(0), store.CallCount(cs.planKey("s1")))
}

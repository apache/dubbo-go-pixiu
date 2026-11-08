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

package mcpserver

import (
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"net/http/httptest"
	"testing"
)

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/router"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// stubSelector is a controllable ToolSelector for hookpoint tests.
type stubSelector struct {
	keep                []string // tool names to keep in Select; nil = keep all
	selectErr           error
	authorizeErr        error // returned by AuthorizeCall
	selectCalled        bool
	authorizeCandidates []model.ToolConfig
	recordSuccessCalls  []router.AuthorizationReceipt
	finalizedReceipts   []router.AuthorizationReceipt
	finalizedOutcomes   []router.ReceiptOutcome
}

func (s *stubSelector) Select(_ context.Context, sc router.SelectionContext, candidates []model.ToolConfig) (*router.SelectionPlan, error) {
	s.selectCalled = true
	if s.selectErr != nil {
		return nil, s.selectErr
	}
	names := s.keep
	if names == nil {
		names = make([]string, len(candidates))
		for i, c := range candidates {
			names[i] = c.Name
		}
	}
	return &router.SelectionPlan{SessionID: sc.SessionID, ToolNames: names, Mode: router.ModeSelected}, nil
}

func (s *stubSelector) AuthorizeCall(_ context.Context, sc router.SelectionContext, candidates []model.ToolConfig) (*router.AuthorizationReceipt, error) {
	s.authorizeCandidates = candidates
	if s.authorizeErr != nil {
		return nil, s.authorizeErr
	}
	return &router.AuthorizationReceipt{SessionID: sc.SessionID, ToolName: sc.Requested, PlanGeneration: 1, ReceiptID: 1}, nil
}

func (s *stubSelector) RecordCallSuccess(_ context.Context, receipt router.AuthorizationReceipt) (router.CallSuccessResult, error) {
	s.recordSuccessCalls = append(s.recordSuccessCalls, receipt)
	return router.CallSuccessResult{Count: int64(len(s.recordSuccessCalls))}, nil
}

func (s *stubSelector) FinalizeReceipt(_ context.Context, receipt router.AuthorizationReceipt, outcome router.ReceiptOutcome) (router.CallSuccessResult, error) {
	s.finalizedReceipts = append(s.finalizedReceipts, receipt)
	s.finalizedOutcomes = append(s.finalizedOutcomes, outcome)
	if outcome != router.ReceiptSucceeded {
		return router.CallSuccessResult{}, nil
	}
	s.recordSuccessCalls = append(s.recordSuccessCalls, receipt)
	return router.CallSuccessResult{Count: int64(len(s.recordSuccessCalls))}, nil
}

type failingSelectionSelector struct {
	*router.CompositeSelector
	err error
}

func (s *failingSelectionSelector) Select(_ context.Context, _ router.SelectionContext, _ []model.ToolConfig) (*router.SelectionPlan, error) {
	if s.err != nil {
		return nil, s.err
	}
	return nil, errors.New("selector failed")
}

func buildToolsListResult(t *testing.T, f *MCPServerFilter, tools []model.ToolConfig) *mcp.ListToolsResult {
	t.Helper()
	require.NoError(t, f.registry.ReplaceAllTools(tools))

	req := mcp.JSONRPCRequest{}
	req.ID = mcp.NewRequestId(int64(1))

	httpReq := httptest.NewRequest("POST", "/mcp", nil)
	ctx := NewMCPContext(createTestContext(httpReq, httptest.NewRecorder()))
	if f.governanceEnabled {
		session, _ := f.sessionManager.CreateSession()
		if f.plans != nil {
			require.NoError(t, f.plans.ActivateSession(session.ID, session.Generation))
		}
		ctx.SetValidatedSession(session)
	}

	resp, err := f.buildToolsListResponseObject(ctx, req)
	require.NoError(t, err)
	result, ok := resp.Result.(*mcp.ListToolsResult)
	require.True(t, ok)
	return result
}

func alphaBetaTools() []model.ToolConfig {
	return []model.ToolConfig{
		createTestToolConfig("alpha", "A"),
		createTestToolConfig("beta", "B"),
	}
}

func TestFilterByPlan_KeepsPlanOrder(t *testing.T) {
	tools := []model.ToolConfig{
		{Name: "a"}, {Name: "b"}, {Name: "c"},
	}
	plan := &router.SelectionPlan{ToolNames: []string{"c", "a"}}

	out := filterByPlan(tools, plan)

	require.Len(t, out, 2)
	assert.Equal(t, "c", out[0].Name)
	assert.Equal(t, "a", out[1].Name)
}

func TestFilterByPlan_NilPlanReturnsAll(t *testing.T) {
	tools := []model.ToolConfig{{Name: "a"}, {Name: "b"}}
	assert.Equal(t, tools, filterByPlan(tools, nil))
}

func TestFilterByPlan_UnknownNamesIgnored(t *testing.T) {
	tools := []model.ToolConfig{{Name: "a"}}
	plan := &router.SelectionPlan{ToolNames: []string{"a", "ghost"}}

	out := filterByPlan(tools, plan)

	require.Len(t, out, 1)
	assert.Equal(t, "a", out[0].Name)
}

func TestFilterByPlan_UsesVisibleToolNames(t *testing.T) {
	tools := []model.ToolConfig{{Name: "visible"}, {Name: "hidden"}}
	plan := &router.SelectionPlan{
		ToolNames:        []string{"visible", "hidden"},
		VisibleToolNames: []string{"visible"},
	}

	out := filterByPlan(tools, plan)

	require.Len(t, out, 1)
	assert.Equal(t, "visible", out[0].Name)
}

func TestVisibleToolsFingerprintTracksActualToolsListProjection(t *testing.T) {
	base := model.ToolConfig{
		Name:        "tool",
		Description: "desc",
		Cluster:     "cluster-a",
		BackendURL:  "http://127.0.0.1:8080",
		Request:     model.RequestConfig{Method: "GET", Path: "/v1/{id}"},
		Args: []model.ArgConfig{{
			Name:        "id",
			Type:        "string",
			In:          "path",
			Description: "identifier",
			Required:    true,
			Enum:        []string{"a", "b"},
			Default:     "a",
		}},
		Meta: &model.ToolMeta{Risk: "low"},
	}
	baseFP := visibleToolsFingerprint([]model.ToolConfig{base})

	descriptionChanged := base
	descriptionChanged.Description = "changed"
	assert.NotEqual(t, baseFP, visibleToolsFingerprint([]model.ToolConfig{descriptionChanged}))

	schemaChanged := base
	schemaChanged.Args = append([]model.ArgConfig(nil), base.Args...)
	schemaChanged.Args[0].Type = "number"
	assert.NotEqual(t, baseFP, visibleToolsFingerprint([]model.ToolConfig{schemaChanged}))

	requiredChanged := base
	requiredChanged.Args = append([]model.ArgConfig(nil), base.Args...)
	requiredChanged.Args[0].Required = false
	assert.NotEqual(t, baseFP, visibleToolsFingerprint([]model.ToolConfig{requiredChanged}))

	defaultChanged := base
	defaultChanged.Args = append([]model.ArgConfig(nil), base.Args...)
	defaultChanged.Args[0].Default = "b"
	assert.NotEqual(t, baseFP, visibleToolsFingerprint([]model.ToolConfig{defaultChanged}))

	enumChanged := base
	enumChanged.Args = append([]model.ArgConfig(nil), base.Args...)
	enumChanged.Args[0].Enum = []string{"a", "c"}
	assert.NotEqual(t, baseFP, visibleToolsFingerprint([]model.ToolConfig{enumChanged}))

	pathOnlyChanged := base
	pathOnlyChanged.Request.Path = "/v2/{id}"
	assert.Equal(t, baseFP, visibleToolsFingerprint([]model.ToolConfig{pathOnlyChanged}), "fingerprint follows tools/list projection, not backend path templates")

	backendChanged := base
	backendChanged.BackendURL = "http://127.0.0.1:9090"
	backendChanged.Cluster = "cluster-b"
	backendChanged.Meta = &model.ToolMeta{Risk: "high", Tags: []string{"internal"}}
	assert.Equal(t, baseFP, visibleToolsFingerprint([]model.ToolConfig{backendChanged}))

	hidden := false
	hiddenTool := base
	hiddenTool.Name = "hidden"
	hiddenTool.Meta = &model.ToolMeta{DiscoveryVisibility: &hidden}
	assert.Equal(t, baseFP, visibleToolsFingerprint([]model.ToolConfig{base, hiddenTool}))
}

func TestVisibleToolsFingerprintMatchesToolsListToolJSON(t *testing.T) {
	tool := createTestToolConfig("tool", "desc")
	tool.Args = []model.ArgConfig{{
		Name:        "id",
		Type:        "string",
		Description: "identifier",
		Required:    true,
	}}
	f := createTestFilter(t)
	result := buildToolsListResult(t, f, []model.ToolConfig{tool})
	require.Len(t, result.Tools, 1)

	got, err := json.Marshal(result.Tools[0])
	require.NoError(t, err)
	rebuilt := (&MCPServerFilter{}).buildMCPTools([]model.ToolConfig{tool})
	require.Len(t, rebuilt, 1)
	want, err := json.Marshal(rebuilt[0])
	require.NoError(t, err)

	assert.JSONEq(t, string(want), string(got))
	assert.Equal(t, visibleToolsFingerprint([]model.ToolConfig{tool}), visibleToolsFingerprint([]model.ToolConfig{tool}))
}

func TestVisibleToolsFingerprintIgnoresToolOrder(t *testing.T) {
	a := createTestToolConfig("a", "A")
	b := createTestToolConfig("b", "B")

	assert.Equal(t,
		visibleToolsFingerprint([]model.ToolConfig{a, b}),
		visibleToolsFingerprint([]model.ToolConfig{b, a}))
}

func TestBuildSelectionContext_PopulatesFields(t *testing.T) {
	req := httptest.NewRequest("POST", "/mcp", nil)
	ctx := NewMCPContext(createTestContext(req, httptest.NewRecorder()))
	ctx.SetSessionID("sess-xyz")

	f := createTestFilter(t)
	sc := f.buildSelectionContext(ctx, "tools/call", "get_user")

	assert.Equal(t, "sess-xyz", sc.SessionID)
	assert.Equal(t, "tools/call", sc.Method)
	assert.Equal(t, "get_user", sc.Requested)
}

// TestToolsList_ResponseBuilderWithoutGovernanceReturnsAll covers the no-router
// response helper path.
func TestToolsList_ResponseBuilderWithoutSelectorReturnsAll(t *testing.T) {
	f := createTestFilter(t)
	f.governanceEnabled = false
	f.selector = nil
	result := buildToolsListResult(t, f, alphaBetaTools())
	assert.Len(t, result.Tools, 2)
}

func TestFilterFactory_OmittedRouterLeavesGovernanceStateNil(t *testing.T) {
	ResetGlobalState()
	defer ResetGlobalState()

	cfg := &model.McpServerConfig{
		ServerInfo: model.ServerInfo{Name: "Test", Version: "1.0.0"},
		Endpoint:   "/mcp",
		Tools: []model.ToolConfig{
			createTestToolConfig("alpha", "A"),
		},
	}
	factory := &FilterFactory{cfg: cfg}

	require.NoError(t, factory.Apply())
	assert.NotNil(t, factory.runtime)
	assert.False(t, factory.runtime.governanceEnabled)
	assert.Nil(t, factory.runtime.selector)
	assert.Nil(t, factory.runtime.plans)
}

func TestFilterFactory_ConfiguredRouterInitializesGovernanceState(t *testing.T) {
	ResetGlobalState()
	defer ResetGlobalState()

	cfg := &model.McpServerConfig{
		ServerInfo: model.ServerInfo{Name: "Test", Version: "1.0.0"},
		Endpoint:   "/mcp",
		Tools: []model.ToolConfig{
			createTestToolConfig("alpha", "A"),
		},
		Router: &model.RouterConfig{Fallback: router.FallbackFailClosed},
	}
	factory := &FilterFactory{cfg: cfg}

	require.NoError(t, factory.Apply())
	assert.True(t, factory.runtime.governanceEnabled)
	assert.NotNil(t, factory.runtime.selector)
	assert.NotNil(t, factory.runtime.plans)
	assert.NotEmpty(t, factory.runtime.dynamic.RuntimeID())
}

func TestRuntimePublicationSinkBinding(t *testing.T) {
	ResetGlobalState()
	defer ResetGlobalState()

	cfg := &model.McpServerConfig{
		ServerInfo: model.ServerInfo{Name: "Test", Version: "1.0.0"},
		Endpoint:   "/mcp",
		Tools:      []model.ToolConfig{createTestToolConfig("alpha", "A")},
	}
	factory := &FilterFactory{cfg: cfg}
	require.NoError(t, factory.Apply())

	sink, err := ServerPublicationSinkForSingleRuntime()
	require.NoError(t, err)
	assert.Same(t, factory.runtime.dynamic, sink)
	assert.Equal(t, factory.runtime.id, sink.RuntimeID())

	factory.runtime.Stop()
	assert.Empty(t, sink.RuntimeID())
}

func TestRuntimePublicationSinkAmbiguousWhenMultipleRuntimes(t *testing.T) {
	ResetGlobalState()
	defer ResetGlobalState()

	cfg := &model.McpServerConfig{
		ServerInfo: model.ServerInfo{Name: "Test", Version: "1.0.0"},
		Endpoint:   "/mcp",
	}
	factoryA := &FilterFactory{cfg: cfg}
	require.NoError(t, factoryA.Apply())
	factoryB := &FilterFactory{cfg: cfg}
	require.NoError(t, factoryB.Apply())

	sink, err := ServerPublicationSinkForSingleRuntime()
	assert.ErrorIs(t, err, ErrDynamicConsumerAmbiguous)
	assert.Nil(t, sink)
}

func TestFilterFactory_InvalidToolRiskFailsFast(t *testing.T) {
	ResetGlobalState()
	defer ResetGlobalState()

	tool := createTestToolConfig("alpha", "A")
	tool.Meta = &model.ToolMeta{Risk: "hihg"}
	cfg := &model.McpServerConfig{
		ServerInfo: model.ServerInfo{Name: "Test", Version: "1.0.0"},
		Endpoint:   "/mcp",
		Tools:      []model.ToolConfig{tool},
		Router:     &model.RouterConfig{},
	}
	factory := &FilterFactory{cfg: cfg}

	err := factory.Apply()
	assert.ErrorContains(t, err, "invalid mcp tool router metadata")
	assert.ErrorContains(t, err, "unsupported risk")
}

func TestFilterFactory_OmittedRouterSkipsGovernanceMetadataValidation(t *testing.T) {
	ResetGlobalState()
	defer ResetGlobalState()

	tool := createTestToolConfig("alpha", "A")
	tool.Meta = &model.ToolMeta{Risk: "legacy-risk-value"}
	cfg := &model.McpServerConfig{
		ServerInfo: model.ServerInfo{Name: "Test", Version: "1.0.0"},
		Endpoint:   "/mcp",
		Tools:      []model.ToolConfig{tool},
	}
	factory := &FilterFactory{cfg: cfg}

	require.NoError(t, factory.Apply())
	require.NotNil(t, factory.runtime)
	assert.False(t, factory.runtime.governanceEnabled)
	assert.Nil(t, factory.runtime.selector)
}

func TestFilterFactory_DuplicateStaticToolFailsFast(t *testing.T) {
	ResetGlobalState()
	defer ResetGlobalState()

	cfg := &model.McpServerConfig{
		ServerInfo: model.ServerInfo{Name: "Test", Version: "1.0.0"},
		Endpoint:   "/mcp",
		Tools: []model.ToolConfig{
			createTestToolConfig("dup", "A"),
			createTestToolConfig("dup", "B"),
		},
	}
	factory := &FilterFactory{cfg: cfg}

	err := factory.Apply()
	assert.ErrorContains(t, err, "duplicate tool name")
}

// TestToolsList_SelectorTrimsTools confirms the Select hookpoint trims the set.
func TestToolsList_SelectorTrimsTools(t *testing.T) {
	f := createTestFilter(t)
	f.selector = &stubSelector{keep: []string{"alpha"}}
	result := buildToolsListResult(t, f, alphaBetaTools())
	require.Len(t, result.Tools, 1)
	assert.Equal(t, "alpha", result.Tools[0].Name)
}

func TestToolsList_SelectorErrorFailClosed(t *testing.T) {
	f := createTestFilter(t)
	f.selector = &failingSelectionSelector{
		CompositeSelector: f.selector.(*router.CompositeSelector),
		err:               errors.New("selector exploded"),
	}

	result := buildToolsListResult(t, f, alphaBetaTools())

	require.Empty(t, result.Tools)
}

func TestToolsList_SelectorErrorBundleDefaultDoesNotExposeFullCatalog(t *testing.T) {
	f := createTestFilter(t)
	store := router.NewSessionPlanStore()
	defer store.Stop()

	sel, err := router.Build(&model.RouterConfig{
		Fallback:      router.FallbackBundleDefault,
		DefaultBundle: "safe-minimal",
		Policy: model.PolicyConfig{Rules: []model.PolicyRule{
			{Name: "block-admin", DenyTags: []string{"admin"}},
		}},
		Workflows: []model.WorkflowConfig{
			{Name: "safe-minimal", Tools: []string{"alpha", "beta"}},
		},
	}, store)
	require.NoError(t, err)
	f.plans = store
	f.selector = &failingSelectionSelector{
		CompositeSelector: sel.(*router.CompositeSelector),
		err:               errors.New("selector exploded"),
	}

	alpha := createTestToolConfig("alpha", "A")
	alpha.Meta = &model.ToolMeta{Tags: []string{"safe"}}
	beta := createTestToolConfig("beta", "B")
	beta.Meta = &model.ToolMeta{Tags: []string{"admin"}}
	gamma := createTestToolConfig("gamma", "G")
	gamma.Meta = &model.ToolMeta{Tags: []string{"safe"}}

	result := buildToolsListResult(t, f, []model.ToolConfig{alpha, beta, gamma})

	require.Len(t, result.Tools, 1)
	assert.Equal(t, "alpha", result.Tools[0].Name)
}

func TestToolsList_PlanStoreFullReturnsInternalError(t *testing.T) {
	f := createTestFilter(t)
	store := router.NewSessionPlanStoreWithOptions(router.SessionPlanStoreOptions{MaxEntries: 1})
	defer store.Stop()
	sel, err := router.Build(&model.RouterConfig{}, store)
	require.NoError(t, err)
	f.plans = store
	f.selector = sel
	require.NoError(t, f.registry.ReplaceAllTools(alphaBetaTools()))

	occupied := "occupied"
	require.NoError(t, store.Set(router.NewPlanKey("occupied-router", occupied),
		&router.SelectionPlan{SessionID: occupied, ToolNames: []string{"alpha"}},
		router.SelectionContext{SessionID: occupied}))

	session, err := f.sessionManager.CreateSession()
	require.NoError(t, err)
	require.NoError(t, store.ActivateSession(session.ID, session.Generation))
	req := mcp.JSONRPCRequest{}
	req.ID = mcp.NewRequestId(int64(30))
	ctx := NewMCPContext(createTestContext(httptest.NewRequest("POST", "/mcp", nil), httptest.NewRecorder()))
	ctx.SetValidatedSession(session)

	_, err = f.buildToolsListResponseObject(ctx, req)
	assert.ErrorIs(t, err, router.ErrPlanStoreFull)
}

func TestToolsList_FallbackPlanStoreFailureIsReturned(t *testing.T) {
	f := createTestFilter(t)
	store := router.NewSessionPlanStoreWithOptions(router.SessionPlanStoreOptions{MaxEntries: 1})
	defer store.Stop()
	sel, err := router.Build(&model.RouterConfig{}, store)
	require.NoError(t, err)
	f.plans = store
	f.selector = &failingSelectionSelector{
		CompositeSelector: sel.(*router.CompositeSelector),
		err:               errors.New("selector exploded"),
	}
	require.NoError(t, f.registry.ReplaceAllTools(alphaBetaTools()))

	occupied := "occupied"
	require.NoError(t, store.Set(router.NewPlanKey("occupied-router", occupied),
		&router.SelectionPlan{SessionID: occupied, ToolNames: []string{"alpha"}},
		router.SelectionContext{SessionID: occupied}))

	session, err := f.sessionManager.CreateSession()
	require.NoError(t, err)
	require.NoError(t, store.ActivateSession(session.ID, session.Generation))
	req := mcp.JSONRPCRequest{}
	req.ID = mcp.NewRequestId(int64(31))
	ctx := NewMCPContext(createTestContext(httptest.NewRequest("POST", "/mcp", nil), httptest.NewRecorder()))
	ctx.SetValidatedSession(session)

	_, err = f.buildToolsListResponseObject(ctx, req)
	assert.ErrorIs(t, err, router.ErrPlanStoreFull)
}

// TestToolCall_SelectorDeniesUnauthorized confirms AuthorizeCall rejection
// produces a tool call error and the backend is never reached.
func TestToolCall_SelectorDeniesUnauthorized(t *testing.T) {
	f := createTestFilter(t)
	f.selector = &stubSelector{authorizeErr: errors.New("denied")}
	require.NoError(t, f.registry.ReplaceAllTools([]model.ToolConfig{
		createTestToolConfig("get_user", "get user"),
	}))

	req := mcp.JSONRPCRequest{Request: mcp.Request{Method: string(mcp.MethodToolsCall)}}
	req.ID = mcp.NewRequestId(int64(7))
	req.Params = map[string]any{"name": "get_user", "arguments": map[string]any{}}

	httpReq := httptest.NewRequest("POST", "/mcp", nil)
	ctx := NewMCPContext(createTestContext(httpReq, httptest.NewRecorder()))
	ctx.SetMCPRequestID(req.ID)
	session, _ := f.sessionManager.CreateSession()
	require.NoError(t, f.plans.ActivateSession(session.ID, session.Generation))
	ctx.SetValidatedSession(session)

	status := f.handleToolCall(ctx, req)

	// Denied calls are written as a local reply and stop the chain.
	assert.Equal(t, filter.Stop, status)
	assert.Nil(t, ctx.Route)
}

func TestToolCall_InvalidRouterSessionDeniedBeforeLookup(t *testing.T) {
	f := createTestFilter(t)
	f.selector = &stubSelector{}
	require.NoError(t, f.registry.ReplaceAllTools([]model.ToolConfig{
		createTestToolConfig("get_user", "get user"),
	}))

	req := mcp.JSONRPCRequest{Request: mcp.Request{Method: string(mcp.MethodToolsCall)}}
	req.ID = mcp.NewRequestId(int64(17))
	req.Params = map[string]any{"name": "get_user", "arguments": map[string]any{}}

	httpReq := httptest.NewRequest("POST", "/mcp", nil)
	ctx := NewMCPContext(createTestContext(httpReq, httptest.NewRecorder()))
	ctx.SetMCPRequestID(req.ID)
	ctx.SetSessionID("unknown-session")

	status := f.handleToolCall(ctx, req)

	assert.Equal(t, filter.Stop, status)
	assert.Nil(t, ctx.Route)
}

func TestPostToolCall_WithSSEAcceptStillAuthorizes(t *testing.T) {
	f := createTestFilter(t)
	sel := &stubSelector{authorizeErr: errors.New("denied")}
	f.selector = sel
	require.NoError(t, f.registry.ReplaceAllTools([]model.ToolConfig{
		createTestToolConfig("get_user", "get user"),
	}))
	session, _ := f.sessionManager.CreateSession()

	reqBody := []byte(`{"jsonrpc":"2.0","id":11,"method":"tools/call","params":{"name":"get_user","arguments":{}}}`)
	httpReq := httptest.NewRequest("POST", "/mcp", bytes.NewReader(reqBody))
	httpReq.Header.Set(constant.HeaderKeyMCPSessionId, session.ID)
	httpReq.Header.Set(constant.HeaderKeyAccept, constant.HeaderValueTextEventStream)
	ctx := NewMCPContext(createTestContext(httpReq, httptest.NewRecorder()))
	ctx.ParseAndSetSessionHeader()

	status := f.handlePostRequest(ctx)

	assert.Equal(t, filter.Stop, status)
	assert.Nil(t, ctx.Route)
	require.Len(t, sel.authorizeCandidates, 1)
	assert.Equal(t, "get_user", sel.authorizeCandidates[0].Name)
}

// TestToolCall_SelectorAllowsAuthorized confirms an allowed call proceeds to
// backend forwarding (filter.Continue with a route set).
func TestToolCall_SelectorAllowsAuthorized(t *testing.T) {
	f := createTestFilter(t)
	sel := &stubSelector{} // authorizeErr nil = allow
	f.selector = sel
	require.NoError(t, f.registry.ReplaceAllTools([]model.ToolConfig{
		createTestToolConfig("get_user", "get user"),
	}))

	req := mcp.JSONRPCRequest{Request: mcp.Request{Method: string(mcp.MethodToolsCall)}}
	req.ID = mcp.NewRequestId(int64(8))
	req.Params = map[string]any{"name": "get_user", "arguments": map[string]any{"param": "v"}}

	httpReq := httptest.NewRequest("POST", "/mcp", nil)
	ctx := NewMCPContext(createTestContext(httpReq, httptest.NewRecorder()))
	ctx.SetMCPRequestID(req.ID)
	session, _ := f.sessionManager.CreateSession()
	require.NoError(t, f.plans.ActivateSession(session.ID, session.Generation))
	ctx.SetValidatedSession(session)

	f.handleToolCall(ctx, req)

	require.NotNil(t, ctx.Route)
	assert.Equal(t, "test-cluster", ctx.Route.Cluster)
	require.Len(t, sel.authorizeCandidates, 1)
	assert.Equal(t, "get_user", sel.authorizeCandidates[0].Name)
}

func TestToolCall_LookupFailureAfterAuthorizationAbortsReceipt(t *testing.T) {
	f := createTestFilter(t)
	sel := &stubSelector{}
	f.selector = sel
	require.NoError(t, f.registry.ReplaceAllTools(nil))

	req := mcp.JSONRPCRequest{Request: mcp.Request{Method: string(mcp.MethodToolsCall)}}
	req.ID = mcp.NewRequestId(int64(18))
	req.Params = map[string]any{"name": "missing_tool", "arguments": map[string]any{}}

	httpReq := httptest.NewRequest("POST", "/mcp", nil)
	ctx := NewMCPContext(createTestContext(httpReq, httptest.NewRecorder()))
	ctx.SetMCPRequestID(req.ID)
	session, _ := f.sessionManager.CreateSession()
	require.NoError(t, f.plans.ActivateSession(session.ID, session.Generation))
	ctx.SetValidatedSession(session)

	status := f.handleToolCall(ctx, req)

	assert.Equal(t, filter.Stop, status)
	require.Equal(t, []router.ReceiptOutcome{router.ReceiptAborted}, sel.finalizedOutcomes)
	require.Len(t, sel.finalizedReceipts, 1)
	assert.Equal(t, "missing_tool", sel.finalizedReceipts[0].ToolName)
}

func TestProcessToolCallResponse_BackendErrorDoesNotRecordSuccess(t *testing.T) {
	f := createTestFilter(t)
	sel := &stubSelector{}
	f.selector = sel

	reqID := mcp.NewRequestId(int64(9))
	httpReq := httptest.NewRequest("POST", "/mcp", nil)
	ctx := NewMCPContext(createTestContext(httpReq, httptest.NewRecorder()))
	ctx.SetMCPMethod(string(mcp.MethodToolsCall))
	ctx.SetMCPRequestID(reqID)
	session, _ := f.sessionManager.CreateSession()
	require.NoError(t, f.plans.ActivateSession(session.ID, session.Generation))
	ctx.SetValidatedSession(session)
	ctx.SetAuthorizationReceipt(&router.AuthorizationReceipt{SessionID: session.ID, ToolName: "get_user", PlanGeneration: 1, ReceiptID: 1})

	status := f.processToolCallResponse(ctx, reqID, []byte("backend failed"), 500)

	assert.Equal(t, filter.Stop, status)
	assert.Empty(t, sel.recordSuccessCalls)
	require.Equal(t, []router.ReceiptOutcome{router.ReceiptAborted}, sel.finalizedOutcomes)
}

func TestProcessToolCallResponse_SuccessRecordsToolCall(t *testing.T) {
	f := createTestFilter(t)
	sel := &stubSelector{}
	f.selector = sel

	reqID := mcp.NewRequestId(int64(10))
	httpReq := httptest.NewRequest("POST", "/mcp", nil)
	ctx := NewMCPContext(createTestContext(httpReq, httptest.NewRecorder()))
	ctx.SetMCPMethod(string(mcp.MethodToolsCall))
	ctx.SetMCPRequestID(reqID)
	session, _ := f.sessionManager.CreateSession()
	require.NoError(t, f.plans.ActivateSession(session.ID, session.Generation))
	ctx.SetValidatedSession(session)
	ctx.SetAuthorizationReceipt(&router.AuthorizationReceipt{SessionID: session.ID, ToolName: "get_user", PlanGeneration: 1, ReceiptID: 1})

	status := f.processToolCallResponse(ctx, reqID, []byte("ok"), 200)

	assert.Equal(t, filter.Continue, status)
	require.Len(t, sel.recordSuccessCalls, 1)
	assert.Equal(t, session.ID, sel.recordSuccessCalls[0].SessionID)
	assert.Equal(t, "get_user", sel.recordSuccessCalls[0].ToolName)
	require.Equal(t, []router.ReceiptOutcome{router.ReceiptSucceeded}, sel.finalizedOutcomes)
	assert.Nil(t, ctx.AuthorizationReceipt())
}

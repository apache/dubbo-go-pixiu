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
		ctx.SetSessionID(session.ID)
	}

	resp := f.buildToolsListResponseObject(ctx, req)
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
	f.selector = &stubSelector{selectErr: errors.New("selector exploded")}

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
	ctx.SetSessionID(session.ID)

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
	ctx.SetSessionID(session.ID)

	f.handleToolCall(ctx, req)

	require.NotNil(t, ctx.Route)
	assert.Equal(t, "test-cluster", ctx.Route.Cluster)
	require.Len(t, sel.authorizeCandidates, 1)
	assert.Equal(t, "get_user", sel.authorizeCandidates[0].Name)
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
	ctx.SetMCPToolName("get_user")
	ctx.SetSessionID("s1")

	status := f.processToolCallResponse(ctx, reqID, []byte("backend failed"), 500)

	assert.Equal(t, filter.Stop, status)
	assert.Empty(t, sel.recordSuccessCalls)
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
	ctx.SetMCPToolName("get_user")
	ctx.SetSessionID("s1")
	ctx.SetAuthorizationReceipt(&router.AuthorizationReceipt{SessionID: "s1", ToolName: "get_user", PlanGeneration: 1, ReceiptID: 1})

	status := f.processToolCallResponse(ctx, reqID, []byte("ok"), 200)

	assert.Equal(t, filter.Continue, status)
	require.Len(t, sel.recordSuccessCalls, 1)
	assert.Equal(t, "s1", sel.recordSuccessCalls[0].SessionID)
	assert.Equal(t, "get_user", sel.recordSuccessCalls[0].ToolName)
}

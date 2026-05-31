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
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/mark3labs/mcp-go/mcp"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/router"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// stubSelector is a controllable ToolSelector for hookpoint tests.
type stubSelector struct {
	keep                []string // tool names to keep in Select; nil = keep all
	authorizeErr        error    // returned by AuthorizeCall
	onInitCalled        bool
	selectCalled        bool
	authorizeCandidates []model.ToolConfig
	recordSuccessCalls  []router.SelectionContext
}

func (s *stubSelector) Select(_ context.Context, sc router.SelectionContext, candidates []model.ToolConfig) (*router.SelectionPlan, error) {
	s.selectCalled = true
	names := s.keep
	if names == nil {
		names = make([]string, len(candidates))
		for i, c := range candidates {
			names[i] = c.Name
		}
	}
	return &router.SelectionPlan{SessionID: sc.SessionID, ToolNames: names, Mode: router.ModeHybrid}, nil
}

func (s *stubSelector) AuthorizeCall(_ context.Context, _ router.SelectionContext, candidates []model.ToolConfig) error {
	s.authorizeCandidates = candidates
	return s.authorizeErr
}

func (s *stubSelector) OnInitialize(_ context.Context, _ router.SelectionContext, _ []model.ToolConfig) error {
	s.onInitCalled = true
	return nil
}

func (s *stubSelector) RecordCallSuccess(_ context.Context, sc router.SelectionContext) error {
	s.recordSuccessCalls = append(s.recordSuccessCalls, sc)
	return nil
}

func (s *stubSelector) Name() string { return "stub" }

func buildToolsListResult(t *testing.T, f *MCPServerFilter, tools []model.ToolConfig) *mcp.ListToolsResult {
	t.Helper()
	f.registry.ReplaceAllTools(tools)

	req := mcp.JSONRPCRequest{}
	req.ID = mcp.NewRequestId(int64(1))

	httpReq := httptest.NewRequest("POST", "/mcp", nil)
	ctx := NewMCPContext(createTestContext(httpReq, httptest.NewRecorder()))

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

// TestToolsList_NilSelectorReturnsAll confirms that with no selector the
// tools/list response contains every registered tool (passthrough behavior).
func TestToolsList_NilSelectorReturnsAll(t *testing.T) {
	f := createTestFilter(t)
	result := buildToolsListResult(t, f, alphaBetaTools())
	assert.Len(t, result.Tools, 2)
}

func TestFilterFactory_RouterDisabledDoesNotInitPlanStore(t *testing.T) {
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
	assert.Nil(t, factory.selector)
	assert.Nil(t, globalPlanStore)

	cfg.Router = &model.RouterConfig{Enabled: false}
	require.NoError(t, factory.Apply())
	assert.Nil(t, factory.selector)
	assert.Nil(t, globalPlanStore)
}

func TestFilterFactory_RouterEnabledInitializesPlanStore(t *testing.T) {
	ResetGlobalState()
	defer ResetGlobalState()

	cfg := &model.McpServerConfig{
		ServerInfo: model.ServerInfo{Name: "Test", Version: "1.0.0"},
		Endpoint:   "/mcp",
		Tools: []model.ToolConfig{
			createTestToolConfig("alpha", "A"),
		},
		Router: &model.RouterConfig{Enabled: true, Fallback: router.FallbackFailClosed},
	}
	factory := &FilterFactory{cfg: cfg}

	require.NoError(t, factory.Apply())
	assert.NotNil(t, factory.selector)
	assert.NotNil(t, globalPlanStore)
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
	}
	factory := &FilterFactory{cfg: cfg}

	err := factory.Apply()
	assert.ErrorContains(t, err, "invalid mcp tool router metadata")
	assert.ErrorContains(t, err, "unsupported risk")
}

// TestToolsList_SelectorTrimsTools confirms the Select hookpoint trims the set.
func TestToolsList_SelectorTrimsTools(t *testing.T) {
	f := createTestFilter(t)
	f.selector = &stubSelector{keep: []string{"alpha"}}
	result := buildToolsListResult(t, f, alphaBetaTools())
	require.Len(t, result.Tools, 1)
	assert.Equal(t, "alpha", result.Tools[0].Name)
}

// TestToolCall_SelectorDeniesUnauthorized confirms AuthorizeCall rejection
// produces a tool call error and the backend is never reached.
func TestToolCall_SelectorDeniesUnauthorized(t *testing.T) {
	f := createTestFilter(t)
	f.selector = &stubSelector{authorizeErr: errors.New("denied")}
	f.registry.ReplaceAllTools([]model.ToolConfig{
		createTestToolConfig("get_user", "get user"),
	})

	req := mcp.JSONRPCRequest{Request: mcp.Request{Method: string(mcp.MethodToolsCall)}}
	req.ID = mcp.NewRequestId(int64(7))
	req.Params = map[string]any{"name": "get_user", "arguments": map[string]any{}}

	httpReq := httptest.NewRequest("POST", "/mcp", nil)
	ctx := NewMCPContext(createTestContext(httpReq, httptest.NewRecorder()))
	ctx.SetMCPRequestID(req.ID)

	status := f.handleToolCall(ctx, req)

	// Denied calls are written as a local reply and stop the chain.
	assert.Equal(t, filter.Stop, status)
	assert.Nil(t, ctx.Route)
}

func TestPostToolCall_WithSSEAcceptStillAuthorizes(t *testing.T) {
	f := createTestFilter(t)
	sel := &stubSelector{authorizeErr: errors.New("denied")}
	f.selector = sel
	f.registry.ReplaceAllTools([]model.ToolConfig{
		createTestToolConfig("get_user", "get user"),
	})
	session, _ := f.sessionManager.EnsureSession("")

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
	f.registry.ReplaceAllTools([]model.ToolConfig{
		createTestToolConfig("get_user", "get user"),
	})

	req := mcp.JSONRPCRequest{Request: mcp.Request{Method: string(mcp.MethodToolsCall)}}
	req.ID = mcp.NewRequestId(int64(8))
	req.Params = map[string]any{"name": "get_user", "arguments": map[string]any{"param": "v"}}

	httpReq := httptest.NewRequest("POST", "/mcp", nil)
	ctx := NewMCPContext(createTestContext(httpReq, httptest.NewRecorder()))
	ctx.SetMCPRequestID(req.ID)

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

	status := f.processToolCallResponse(ctx, reqID, []byte("ok"), 200)

	assert.Equal(t, filter.Continue, status)
	require.Len(t, sel.recordSuccessCalls, 1)
	assert.Equal(t, "s1", sel.recordSuccessCalls[0].SessionID)
	assert.Equal(t, string(mcp.MethodToolsCall), sel.recordSuccessCalls[0].Method)
	assert.Equal(t, "get_user", sel.recordSuccessCalls[0].Requested)
}

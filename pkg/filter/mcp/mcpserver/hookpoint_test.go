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
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/router"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// stubSelector is a controllable ToolSelector for hookpoint tests.
type stubSelector struct {
	keep         []string // tool names to keep in Select; nil = keep all
	authorizeErr error    // returned by AuthorizeCall
	onInitCalled bool
	selectCalled bool
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
	return &router.SelectionPlan{SessionID: sc.SessionID, ToolNames: names, Mode: router.ModePolicy}, nil
}

func (s *stubSelector) AuthorizeCall(_ context.Context, _ router.SelectionContext) error {
	return s.authorizeErr
}

func (s *stubSelector) OnInitialize(_ context.Context, _ router.SelectionContext, _ []model.ToolConfig) error {
	s.onInitCalled = true
	return nil
}

func (s *stubSelector) Name() string { return "stub" }

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

func TestSanitizeHeaders_StripsSensitive(t *testing.T) {
	req := httptest.NewRequest("POST", "/mcp", nil)
	req.Header.Set("Authorization", "Bearer secret")
	req.Header.Set("Cookie", "sid=abc")
	req.Header.Set("X-Tenant", "acme")

	ctx := NewMCPContext(createTestContext(req, httptest.NewRecorder()))
	headers := sanitizeHeaders(ctx)

	assert.Equal(t, "acme", headers["X-Tenant"])
	_, hasAuth := headers["Authorization"]
	assert.False(t, hasAuth)
	_, hasCookie := headers["Cookie"]
	assert.False(t, hasCookie)
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
	f.registry.ReplaceAllTools([]model.ToolConfig{
		createTestToolConfig("alpha", "A"),
		createTestToolConfig("beta", "B"),
	})

	req := mcp.JSONRPCRequest{}
	req.ID = mcp.NewRequestId(int64(1))

	httpReq := httptest.NewRequest("POST", "/mcp", nil)
	ctx := NewMCPContext(createTestContext(httpReq, httptest.NewRecorder()))

	resp := f.buildToolsListResponseObject(ctx, req)

	result, ok := resp.Result.(*mcp.ListToolsResult)
	require.True(t, ok)
	assert.Len(t, result.Tools, 2)
}

// TestToolsList_SelectorTrimsTools confirms the Select hookpoint trims the set.
func TestToolsList_SelectorTrimsTools(t *testing.T) {
	f := createTestFilter(t)
	f.selector = &stubSelector{keep: []string{"alpha"}}
	f.registry.ReplaceAllTools([]model.ToolConfig{
		createTestToolConfig("alpha", "A"),
		createTestToolConfig("beta", "B"),
	})

	req := mcp.JSONRPCRequest{}
	req.ID = mcp.NewRequestId(int64(1))

	httpReq := httptest.NewRequest("POST", "/mcp", nil)
	ctx := NewMCPContext(createTestContext(httpReq, httptest.NewRecorder()))

	resp := f.buildToolsListResponseObject(ctx, req)

	result, ok := resp.Result.(*mcp.ListToolsResult)
	require.True(t, ok)
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

// TestToolCall_SelectorAllowsAuthorized confirms an allowed call proceeds to
// backend forwarding (filter.Continue with a route set).
func TestToolCall_SelectorAllowsAuthorized(t *testing.T) {
	f := createTestFilter(t)
	f.selector = &stubSelector{} // authorizeErr nil = allow
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
}

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
	"fmt"
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
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/transport"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// buildTenantFixture creates 50 tools across 3 tenants, each tool tagged with
// its tenant name so tenant-isolation policy can act on it.
func buildTenantFixture() []model.ToolConfig {
	tenants := []string{"acme", "globex", "initech"}
	tools := make([]model.ToolConfig, 0, 50)
	for i := 0; i < 50; i++ {
		tenant := tenants[i%len(tenants)]
		tools = append(tools, model.ToolConfig{
			Name:    fmt.Sprintf("%s_tool_%d", tenant, i),
			Cluster: "test-cluster",
			Request: model.RequestConfig{Method: "GET", Path: "/api/test"},
			Meta:    &model.ToolMeta{Tags: []string{tenant}},
		})
	}
	return tools
}

// newRoutedFilter creates an MCPServerFilter with a tenant-isolation router
// enabled, returning the filter and its plan store.
func newRoutedFilter(t *testing.T, tools []model.ToolConfig) *MCPServerFilter {
	enforce := true
	cfg := &model.McpServerConfig{
		ServerInfo: model.ServerInfo{Name: "Test", Version: "1.0.0"},
		Endpoint:   "/mcp",
		Tools:      tools,
		Router: &model.RouterConfig{
			Enabled:       true,
			Fallback:      router.FallbackFailClosed,
			EnforceOnCall: &enforce,
			Policy: model.PolicyConfig{Rules: []model.PolicyRule{
				{Name: "acme", When: model.PolicyMatch{Claim: "tenant", Equals: "acme"}, AllowTags: []string{"acme"}},
				{Name: "globex", When: model.PolicyMatch{Claim: "tenant", Equals: "globex"}, AllowTags: []string{"globex"}},
				{Name: "initech", When: model.PolicyMatch{Claim: "tenant", Equals: "initech"}, AllowTags: []string{"initech"}},
			}},
		},
	}

	factory := &FilterFactory{cfg: cfg}
	require.NoError(t, factory.Apply())

	sm := transport.NewSessionManager()
	return &MCPServerFilter{
		cfg:               cfg,
		registry:          factory.registry,
		errorHandler:      NewErrorHandler(),
		responseBuilder:   NewResponseBuilder(),
		sessionManager:    sm,
		sseHandler:        transport.NewSSEHandler(sm),
		contentNegotiator: transport.NewContentNegotiator(),
		selector:          factory.selector,
	}
}

// listToolsForTenant runs tools/list for a session bound to a tenant and
// returns the visible tool names.
func listToolsForTenant(f *MCPServerFilter, sessionID, tenant string) []string {
	httpReq := httptest.NewRequest("POST", "/mcp", nil)
	ctx := NewMCPContext(createTestContext(httpReq, httptest.NewRecorder()))
	ctx.SetSessionID(sessionID)
	ctx.Params[constant.MCPAuthClaimsParamKey] = map[string]any{"tenant": tenant, "sub": "u-" + tenant}

	req := mcp.JSONRPCRequest{}
	req.ID = mcp.NewRequestId(int64(1))

	resp := f.buildToolsListResponseObject(ctx, req)
	result := resp.Result.(*mcp.ListToolsResult)
	names := make([]string, len(result.Tools))
	for i, tool := range result.Tools {
		names[i] = tool.Name
	}
	return names
}

func TestIntegration_TenantIsolation_ListAndCall(t *testing.T) {
	ResetGlobalState()
	defer ResetGlobalState()

	tools := buildTenantFixture()
	f := newRoutedFilter(t, tools)
	require.NotNil(t, f.selector)

	// 1. acme session sees only acme tools.
	acmeSession, _ := f.sessionManager.CreateSession()
	acmeNames := listToolsForTenant(f, acmeSession.ID, "acme")
	require.NotEmpty(t, acmeNames)
	for _, name := range acmeNames {
		assert.Contains(t, name, "acme_", "acme session must only see acme tools, got %s", name)
	}

	// 2. globex session sees only globex tools (no cross-tenant leakage).
	globexSession, _ := f.sessionManager.CreateSession()
	globexNames := listToolsForTenant(f, globexSession.ID, "globex")
	require.NotEmpty(t, globexNames)
	for _, name := range globexNames {
		assert.Contains(t, name, "globex_")
	}

	// 3. acme session cannot call a globex tool, even knowing its name.
	globexTool := globexNames[0]
	status := callTool(f, acmeSession.ID, globexTool)
	assert.Equal(t, filter.Stop, status, "cross-tenant call must be denied (Stop with error reply)")

	// 4. acme session CAN call one of its own tools.
	acmeTool := acmeNames[0]
	statusOK := callTool(f, acmeSession.ID, acmeTool)
	assert.Equal(t, filter.Continue, statusOK, "in-tenant call must be forwarded")
}

// callTool runs tools/call for a session and returns the filter status.
func callTool(f *MCPServerFilter, sessionID, toolName string) filter.FilterStatus {
	httpReq := httptest.NewRequest("POST", "/mcp", nil)
	ctx := NewMCPContext(createTestContext(httpReq, httptest.NewRecorder()))
	ctx.SetSessionID(sessionID)
	ctx.Params[constant.MCPAuthClaimsParamKey] = map[string]any{"tenant": "acme", "sub": "u-acme"}

	req := mcp.JSONRPCRequest{Request: mcp.Request{Method: string(mcp.MethodToolsCall)}}
	req.ID = mcp.NewRequestId(int64(2))
	req.Params = map[string]any{"name": toolName, "arguments": map[string]any{}}
	ctx.SetMCPRequestID(req.ID)

	return f.handleToolCall(ctx, req)
}

func TestIntegration_BypassListDirectCallDenied(t *testing.T) {
	ResetGlobalState()
	defer ResetGlobalState()

	tools := buildTenantFixture()
	f := newRoutedFilter(t, tools)

	// Client skips tools/list entirely and calls directly -> no plan -> denied.
	status := callTool(f, "fresh-session", "acme_tool_0")
	assert.Equal(t, filter.Stop, status)
}

func TestIntegration_RouterDisabledIsPassthrough(t *testing.T) {
	ResetGlobalState()
	defer ResetGlobalState()

	cfg := &model.McpServerConfig{
		ServerInfo: model.ServerInfo{Name: "Test", Version: "1.0.0"},
		Endpoint:   "/mcp",
		Tools:      buildTenantFixture(),
		// No Router block.
	}
	factory := &FilterFactory{cfg: cfg}
	require.NoError(t, factory.Apply())
	assert.Nil(t, factory.selector, "no router config => nil selector => passthrough")
}

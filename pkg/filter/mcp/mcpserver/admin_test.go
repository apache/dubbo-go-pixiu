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
	"encoding/json"
	"net/http/httptest"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/router"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/transport"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// newAdminFilter builds a routed filter with audit.payload_logging toggled.
func newAdminFilter(t *testing.T, payloadLogging bool) *MCPServerFilter {
	cfg := &model.McpServerConfig{
		ServerInfo: model.ServerInfo{Name: "Test", Version: "1.0.0"},
		Endpoint:   "/mcp",
		Tools: []model.ToolConfig{
			{Name: "acme_a", Cluster: "c", Meta: &model.ToolMeta{Tags: []string{"acme"}}},
		},
		Router: &model.RouterConfig{
			Enabled:  true,
			Fallback: router.FallbackFailClosed,
			Policy: model.PolicyConfig{Rules: []model.PolicyRule{
				{Name: "acme", When: model.PolicyMatch{Claim: "tenant", Equals: "acme"}, AllowTags: []string{"acme"}},
			}},
			Audit: model.AuditConfig{PayloadLogging: payloadLogging},
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

func TestAdmin_DisabledReturns404(t *testing.T) {
	ResetGlobalState()
	defer ResetGlobalState()

	f := newAdminFilter(t, false) // payload logging off

	req := httptest.NewRequest("GET", "/__mcp/router/plan/some-session", nil)
	rec := httptest.NewRecorder()
	ctx := createTestContext(req, rec)

	f.Decode(ctx)
	assert.Equal(t, 404, rec.Code)
}

func TestAdmin_EnabledReturnsPlan(t *testing.T) {
	ResetGlobalState()
	defer ResetGlobalState()

	f := newAdminFilter(t, true) // payload logging on

	// First produce a plan via tools/list for a session.
	listReq := httptest.NewRequest("POST", "/mcp", nil)
	listCtx := NewMCPContext(createTestContext(listReq, httptest.NewRecorder()))
	listCtx.SetSessionID("sess-admin")
	listCtx.Params[constant.MCPAuthClaimsParamKey] = map[string]any{"tenant": "acme"}
	plan, err := f.selector.Select(listCtx.Ctx, f.buildSelectionContext(listCtx, "tools/list", ""), f.registry.ListTools())
	require.NoError(t, err)
	require.NotEmpty(t, plan.ToolNames)

	// Now query the admin endpoint.
	req := httptest.NewRequest("GET", "/__mcp/router/plan/sess-admin", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()
	ctx := createTestContext(req, rec)

	f.Decode(ctx)
	require.Equal(t, 200, rec.Code)

	var got router.SelectionPlan
	require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &got))
	assert.Equal(t, "sess-admin", got.SessionID)
	assert.Contains(t, got.ToolNames, "acme_a")
}

func TestAdmin_UnknownSessionReturns404(t *testing.T) {
	ResetGlobalState()
	defer ResetGlobalState()

	f := newAdminFilter(t, true)

	req := httptest.NewRequest("GET", "/__mcp/router/plan/ghost", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()
	ctx := createTestContext(req, rec)

	f.Decode(ctx)
	assert.Equal(t, 404, rec.Code)
}

func TestAdmin_PostReturns405(t *testing.T) {
	ResetGlobalState()
	defer ResetGlobalState()

	f := newAdminFilter(t, true)

	req := httptest.NewRequest("POST", "/__mcp/router/plan/x", nil)
	req.RemoteAddr = "127.0.0.1:12345"
	rec := httptest.NewRecorder()
	ctx := createTestContext(req, rec)

	f.Decode(ctx)
	assert.Equal(t, 405, rec.Code)
}

func TestAdmin_NonLoopbackReturns404(t *testing.T) {
	ResetGlobalState()
	defer ResetGlobalState()

	f := newAdminFilter(t, true)

	req := httptest.NewRequest("GET", "/__mcp/router/plan/x", nil)
	req.RemoteAddr = "203.0.113.10:12345"
	rec := httptest.NewRecorder()
	ctx := createTestContext(req, rec)

	f.Decode(ctx)
	assert.Equal(t, 404, rec.Code)
}

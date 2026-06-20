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
	"encoding/json"
	"io"
	"net/http"
	"net/http/httptest"
	"sync"
	"testing"
	"time"
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

type lifecycleClock struct {
	mu  sync.Mutex
	now time.Time
}

func newLifecycleClock() *lifecycleClock {
	return &lifecycleClock{now: time.Unix(1000, 0)}
}

func (c *lifecycleClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *lifecycleClock) Advance(d time.Duration) {
	c.mu.Lock()
	c.now = c.now.Add(d)
	c.mu.Unlock()
}

func newLifecycleFilter(t *testing.T) (*MCPServerFilter, *router.SessionPlanStore) {
	return newLifecycleFilterWithSessionManager(t, transport.NewSessionManager())
}

func newLifecycleFilterWithSessionManager(t *testing.T, sm *transport.SessionManager) (*MCPServerFilter, *router.SessionPlanStore) {
	t.Helper()

	tools := []model.ToolConfig{createTestToolConfig("ping", "ping")}
	reg := NewToolRegistry()
	require.NoError(t, reg.ReplaceAllTools(tools))

	cfg := &model.McpServerConfig{
		ServerInfo: model.ServerInfo{Name: "Test", Version: "1.0.0"},
		Endpoint:   "/mcp",
		Tools:      tools,
		Router: &model.RouterConfig{
			Fallback: router.FallbackFailClosed,
		},
	}
	store := router.NewSessionPlanStoreWithOptions(router.SessionPlanStoreOptions{
		MaxEntries: 10,
	})
	sel, err := router.Build(cfg.Router, store)
	require.NoError(t, err)

	sm.AddSessionRemovedHandler(store.DeleteSession)
	return &MCPServerFilter{
		cfg:               cfg,
		registry:          reg,
		errorHandler:      NewErrorHandler(),
		responseBuilder:   NewResponseBuilder(),
		sessionManager:    sm,
		plans:             store,
		sseHandler:        transport.NewSSEHandler(sm),
		contentNegotiator: transport.NewContentNegotiator(),
		selector:          sel,
		governanceEnabled: true,
	}, store
}

func newLifecycleAllAllowedFilter(t *testing.T, routerCfg *model.RouterConfig) *MCPServerFilter {
	t.Helper()

	tools := []model.ToolConfig{createTestToolConfig("ping", "ping"), createTestToolConfig("pong", "pong")}
	cfg := &model.McpServerConfig{
		ServerInfo: model.ServerInfo{Name: "Test", Version: "1.0.0"},
		Endpoint:   "/mcp",
		Tools:      tools,
		Router:     routerCfg,
	}
	factory := &FilterFactory{cfg: cfg}
	require.NoError(t, factory.Apply())
	if routerCfg != nil {
		require.NotNil(t, factory.runtime.selector)
		require.True(t, factory.runtime.governanceEnabled)
	} else {
		require.Nil(t, factory.runtime.selector)
		require.False(t, factory.runtime.governanceEnabled)
	}
	return &MCPServerFilter{
		cfg:               cfg,
		registry:          factory.runtime.registry,
		errorHandler:      NewErrorHandler(),
		responseBuilder:   NewResponseBuilder(),
		sessionManager:    factory.runtime.sessionManager,
		plans:             factory.runtime.plans,
		sseHandler:        factory.runtime.sseHandler,
		contentNegotiator: transport.NewContentNegotiator(),
		selector:          factory.runtime.selector,
		governanceEnabled: factory.runtime.governanceEnabled,
	}
}

func initializeSession(t *testing.T, f *MCPServerFilter) string {
	t.Helper()

	req := mcp.JSONRPCRequest{Request: mcp.Request{Method: string(mcp.MethodInitialize)}}
	req.ID = mcp.NewRequestId(int64(1))
	req.Params = map[string]any{
		"protocolVersion": constant.MCPProtocolVersion20250618,
		"clientInfo":      map[string]any{"name": "client", "version": "1.0"},
		"capabilities":    map[string]any{},
	}

	httpReq := httptest.NewRequest("POST", "/mcp", nil)
	rec := httptest.NewRecorder()
	ctx := NewMCPContext(createTestContext(httpReq, rec))

	status := f.handleInitialize(ctx, req)
	require.Equal(t, filter.Stop, status)
	sessionID := rec.Header().Get(constant.HeaderKeyMCPSessionId)
	require.NotEmpty(t, sessionID)
	return sessionID
}

func postMCP(t *testing.T, f *MCPServerFilter, sessionID string, body []byte) (*httptest.ResponseRecorder, filter.FilterStatus) {
	t.Helper()
	httpReq := httptest.NewRequest(http.MethodPost, "/mcp", bytes.NewReader(body))
	httpReq.Header.Set(constant.HeaderKeyAccept, constant.HeaderValueApplicationJson)
	if sessionID != "" {
		httpReq.Header.Set(constant.HeaderKeyMCPSessionId, sessionID)
	}
	rec := httptest.NewRecorder()
	ctx := NewMCPContext(createTestContext(httpReq, rec))
	ctx.ParseAndSetSessionHeader()
	ctx.ParseAndSetAcceptHeader()
	return rec, f.handlePostRequest(ctx)
}

func TestInitializeCreatesNewSessionForStandardClientCapabilities(t *testing.T) {
	cases := []struct {
		name         string
		capabilities map[string]any
	}{
		{name: "empty", capabilities: map[string]any{}},
		{name: "roots", capabilities: map[string]any{"roots": map[string]any{"listChanged": true}}},
		{name: "sampling", capabilities: map[string]any{"sampling": map[string]any{}}},
		{name: "roots and sampling", capabilities: map[string]any{
			"roots":    map[string]any{"listChanged": true},
			"sampling": map[string]any{},
		}},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			f, store := newLifecycleFilter(t)
			defer f.sessionManager.Stop()
			defer store.Stop()

			body, err := json.Marshal(map[string]any{
				"jsonrpc": "2.0",
				"id":      1,
				"method":  string(mcp.MethodInitialize),
				"params": map[string]any{
					"protocolVersion": constant.MCPProtocolVersion20250618,
					"clientInfo":      map[string]any{"name": "standard-client", "version": "1.0"},
					"capabilities":    tc.capabilities,
				},
			})
			require.NoError(t, err)

			rec, status := postMCP(t, f, "", body)

			require.Equal(t, filter.Stop, status)
			require.Equal(t, http.StatusOK, rec.Code)
			require.NotEmpty(t, rec.Header().Get(constant.HeaderKeyMCPSessionId))

			var response map[string]any
			require.NoError(t, json.Unmarshal(rec.Body.Bytes(), &response))
			result := response["result"].(map[string]any)
			capabilities := result["capabilities"].(map[string]any)
			tools := capabilities["tools"].(map[string]any)
			assert.Equal(t, true, tools["listChanged"])
		})
	}
}

func TestInitializeWithSuppliedSessionIDRejected(t *testing.T) {
	f, store := newLifecycleFilter(t)
	defer f.sessionManager.Stop()
	defer store.Stop()

	existing, _ := f.sessionManager.CreateSession()
	manualKey := router.NewPlanKey("manual", existing.ID)
	store.Set(manualKey, &router.SelectionPlan{SessionID: existing.ID, ToolNames: []string{"ping"}}, router.SelectionContext{SessionID: manualKey.SessionID})

	body := []byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","clientInfo":{"name":"client","version":"1.0"},"capabilities":{}}}`)
	rec, status := postMCP(t, f, existing.ID, body)

	require.Equal(t, filter.Stop, status)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, 1, f.sessionManager.ActiveSessionCount())
	_, ok := store.Get(manualKey)
	assert.True(t, ok, "rejected initialize must not overwrite the existing plan")
}

func TestInitializeWithUnknownSuppliedSessionIDRejected(t *testing.T) {
	f, store := newLifecycleFilter(t)
	defer f.sessionManager.Stop()
	defer store.Stop()

	body := []byte(`{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","clientInfo":{"name":"client","version":"1.0"},"capabilities":{}}}`)
	rec, status := postMCP(t, f, "caller-supplied", body)

	require.Equal(t, filter.Stop, status)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
	assert.Equal(t, 0, f.sessionManager.ActiveSessionCount())
}

func TestRouterAbsentPreservesPostWithoutSessionBehavior(t *testing.T) {
	f := newLifecycleAllAllowedFilter(t, nil)
	defer f.sessionManager.Stop()

	listBody := []byte(`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`)
	rec, status := postMCP(t, f, "", listBody)
	require.Equal(t, filter.Stop, status)
	require.Equal(t, http.StatusOK, rec.Code)

	callBody := []byte(`{"jsonrpc":"2.0","id":3,"method":"tools/call","params":{"name":"ping","arguments":{"param":"v"}}}`)
	rec, status = postMCP(t, f, "", callBody)
	require.Equal(t, filter.Continue, status)
	assert.Equal(t, http.StatusOK, rec.Code)
	assert.Equal(t, 0, f.sessionManager.ActiveSessionCount())
}

func TestPostUnknownSessionReturns404AndDoesNotCreate(t *testing.T) {
	f, store := newLifecycleFilter(t)
	defer f.sessionManager.Stop()
	defer store.Stop()

	body := []byte(`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`)
	rec, status := postMCP(t, f, "missing-session", body)

	require.Equal(t, filter.Stop, status)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, 0, f.sessionManager.ActiveSessionCount())
}

func TestExpiredSessionReturns404AndDoesNotCreate(t *testing.T) {
	clock := newLifecycleClock()
	sm := transport.NewSessionManagerWithNow(clock.Now)
	f, store := newLifecycleFilterWithSessionManager(t, sm)
	defer f.sessionManager.Stop()
	defer store.Stop()

	session, _ := f.sessionManager.CreateSession()
	manualKey := router.NewPlanKey("manual", session.ID)
	store.Set(manualKey, &router.SelectionPlan{SessionID: session.ID, ToolNames: []string{"ping"}}, router.SelectionContext{SessionID: manualKey.SessionID})
	clock.Advance(transport.SessionTimeout + time.Nanosecond)

	getReq := httptest.NewRequest(http.MethodGet, "/mcp", nil)
	getReq.Header.Set(constant.HeaderKeyAccept, constant.HeaderValueTextEventStream)
	getReq.Header.Set(constant.HeaderKeyMCPSessionId, session.ID)
	getRec := httptest.NewRecorder()
	getCtx := NewMCPContext(createTestContext(getReq, getRec))
	getCtx.ParseAndSetAcceptHeader()
	getCtx.ParseAndSetSessionHeader()
	require.Equal(t, filter.Stop, f.handleGetRequest(getCtx))
	assert.Equal(t, http.StatusNotFound, getRec.Code)

	body := []byte(`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`)
	rec, status := postMCP(t, f, session.ID, body)
	require.Equal(t, filter.Stop, status)
	assert.Equal(t, http.StatusNotFound, rec.Code)
	assert.Equal(t, 0, f.sessionManager.ActiveSessionCount())
	_, ok := store.Get(manualKey)
	assert.False(t, ok)
}

func TestPostRouterMissingSessionReturns400(t *testing.T) {
	f, store := newLifecycleFilter(t)
	defer f.sessionManager.Stop()
	defer store.Stop()

	body := []byte(`{"jsonrpc":"2.0","id":2,"method":"tools/list","params":{}}`)
	rec, status := postMCP(t, f, "", body)

	require.Equal(t, filter.Stop, status)
	assert.Equal(t, http.StatusBadRequest, rec.Code)
}

func buildToolsListForSession(t *testing.T, f *MCPServerFilter, sessionID string) []string {
	t.Helper()

	httpReq := httptest.NewRequest("POST", "/mcp", nil)
	ctx := NewMCPContext(createTestContext(httpReq, httptest.NewRecorder()))
	ctx.SetSessionID(sessionID)
	setValidatedSessionForTest(f, ctx, sessionID)

	req := mcp.JSONRPCRequest{Request: mcp.Request{Method: string(mcp.MethodToolsList)}}
	req.ID = mcp.NewRequestId(int64(2))
	resp, err := f.buildToolsListResponseObject(ctx, req)
	require.NoError(t, err)
	result := resp.Result.(*mcp.ListToolsResult)
	names := make([]string, len(result.Tools))
	for i, tool := range result.Tools {
		names[i] = tool.Name
	}
	return names
}

func TestSessionPlanRemovedWithTransportSession(t *testing.T) {
	f, store := newLifecycleFilter(t)
	defer f.sessionManager.Stop()
	defer store.Stop()

	sessionID := initializeSession(t, f)
	assert.Equal(t, []string{"ping"}, buildToolsListForSession(t, f, sessionID))
	require.Equal(t, 1, store.Len())

	f.sessionManager.RemoveSession(sessionID)
	assert.Equal(t, 0, store.Len())

	status := callTool(f, sessionID, "ping")
	assert.Equal(t, filter.Stop, status, "removed transport session must not authorize old plan")

	newSession, _ := f.sessionManager.CreateSession()
	status = callTool(f, newSession.ID, "ping")
	assert.Equal(t, filter.Stop, status, "new session must not inherit old plan")
}

func TestTransportTTLExpiryDeletesPlan(t *testing.T) {
	clock := newLifecycleClock()
	sm := transport.NewSessionManagerWithNow(clock.Now)
	defer sm.Stop()

	store := router.NewSessionPlanStoreWithOptions(router.SessionPlanStoreOptions{MaxEntries: 10})
	defer store.Stop()
	sm.AddSessionRemovedHandler(store.DeleteSession)

	session, _ := sm.CreateSession()
	manualKey := router.NewPlanKey("manual", session.ID)
	store.Set(manualKey, &router.SelectionPlan{SessionID: session.ID, ToolNames: []string{"ping"}}, router.SelectionContext{SessionID: manualKey.SessionID})

	clock.Advance(transport.SessionTimeout + time.Nanosecond)
	_, exists := sm.Session(session.ID)
	assert.False(t, exists)
	_, ok := store.Get(manualKey)
	assert.False(t, ok)
}

func TestToolsListChangedPendingFlushesOnReconnect(t *testing.T) {
	f, store := newLifecycleFilter(t)
	defer f.sessionManager.Stop()
	defer store.Stop()

	session, _ := f.sessionManager.CreateSession()

	f.notifyToolsListChanged(session.ID)
	_, pending := session.PendingToolsListChangedVersion()
	assert.True(t, pending)

	reader, writer := io.Pipe()
	defer reader.Close()
	defer writer.Close()
	_, _ = session.AttachStream(writer)

	notificationCh := make(chan string, 1)
	go func() {
		buf := make([]byte, 512)
		n, err := reader.Read(buf)
		if err == nil {
			notificationCh <- string(buf[:n])
		}
	}()

	f.flushPendingToolsListChanged(session)

	select {
	case notification := <-notificationCh:
		assert.Contains(t, notification, toolsListChangedMethod)
	case <-time.After(time.Second):
		t.Fatal("expected pending tools/list_changed notification")
	}
	_, pending = session.PendingToolsListChangedVersion()
	assert.False(t, pending)
}

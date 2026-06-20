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
	"io"
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
	t.Helper()

	tools := []model.ToolConfig{createTestToolConfig("ping", "ping")}
	reg := NewToolRegistry()
	require.NoError(t, reg.ReplaceAllTools(tools))

	cfg := &model.McpServerConfig{
		ServerInfo: model.ServerInfo{Name: "Test", Version: "1.0.0"},
		Endpoint:   "/mcp",
		Tools:      tools,
		Router: &model.RouterConfig{
			Enabled:  true,
			Fallback: router.FallbackFailClosed,
		},
	}
	store := router.NewSessionPlanStoreWithOptions(router.SessionPlanStoreOptions{
		TTL:             time.Minute,
		MaxEntries:      10,
		CleanupInterval: time.Hour,
	})
	sel, err := router.Build(cfg.Router, store)
	require.NoError(t, err)

	sm := transport.NewSessionManager()
	sm.AddSessionRemovedHandler(store.Delete)
	return &MCPServerFilter{
		cfg:               cfg,
		registry:          reg,
		errorHandler:      NewErrorHandler(),
		responseBuilder:   NewResponseBuilder(),
		sessionManager:    sm,
		sseHandler:        transport.NewSSEHandler(sm),
		contentNegotiator: transport.NewContentNegotiator(),
		selector:          sel,
	}, store
}

func initializeSession(t *testing.T, f *MCPServerFilter) string {
	t.Helper()

	req := mcp.JSONRPCRequest{Request: mcp.Request{Method: string(mcp.MethodInitialize)}}
	req.ID = mcp.NewRequestId(int64(1))
	req.Params = map[string]any{
		"protocolVersion": constant.MCPProtocolVersion20250618,
		"clientInfo":      map[string]any{"name": "client", "version": "1.0"},
		"capabilities":    map[string]any{"tools": map[string]any{"listChanged": true}},
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

func buildToolsListForSession(t *testing.T, f *MCPServerFilter, sessionID string) []string {
	t.Helper()

	httpReq := httptest.NewRequest("POST", "/mcp", nil)
	ctx := NewMCPContext(createTestContext(httpReq, httptest.NewRecorder()))
	ctx.SetSessionID(sessionID)

	req := mcp.JSONRPCRequest{Request: mcp.Request{Method: string(mcp.MethodToolsList)}}
	req.ID = mcp.NewRequestId(int64(2))
	resp := f.buildToolsListResponseObject(ctx, req)
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
	_, ok := store.Get(sessionID)
	require.True(t, ok)

	f.sessionManager.RemoveSession(sessionID)
	_, ok = store.Get(sessionID)
	assert.False(t, ok)

	status := callTool(f, sessionID, "ping")
	assert.Equal(t, filter.Stop, status, "removed transport session must not authorize old plan")

	newSession, _ := f.sessionManager.EnsureSession("")
	status = callTool(f, newSession.ID, "ping")
	assert.Equal(t, filter.Stop, status, "new session must not inherit old plan")
}

func TestTransportTTLExpiryDeletesPlan(t *testing.T) {
	clock := newLifecycleClock()
	sm := transport.NewSessionManagerWithNow(clock.Now)
	defer sm.Stop()

	store := router.NewSessionPlanStoreWithOptions(router.SessionPlanStoreOptions{
		TTL:             time.Hour,
		MaxEntries:      10,
		Now:             clock.Now,
		CleanupInterval: time.Hour,
	})
	defer store.Stop()
	sm.AddSessionRemovedHandler(store.Delete)

	session, _ := sm.EnsureSession("")
	store.Set(&router.SelectionPlan{SessionID: session.ID, ToolNames: []string{"ping"}})

	clock.Advance(transport.SessionTimeout + time.Nanosecond)
	_, exists := sm.Session(session.ID)
	assert.False(t, exists)
	_, ok := store.Get(session.ID)
	assert.False(t, ok)
}

func TestToolsListChangedPendingFlushesOnReconnect(t *testing.T) {
	f, store := newLifecycleFilter(t)
	defer f.sessionManager.Stop()
	defer store.Stop()

	session, _ := f.sessionManager.EnsureSession("")
	session.SetToolsListChangedSupported(true)

	f.notifyToolsListChanged(session.ID)
	assert.True(t, session.ConsumeToolsListChangedPending())
	session.MarkToolsListChangedPending()

	reader, writer := io.Pipe()
	defer reader.Close()
	defer writer.Close()
	session.SetPipeWriter(writer)

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
	assert.False(t, session.ConsumeToolsListChangedPending())
}

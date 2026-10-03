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
	"io"
	"net/http"
	"net/http/httptest"
	"runtime"
	"strings"
	"testing"
	"time"
)

import (
	"github.com/mark3labs/mcp-go/mcp"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/transport"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestHandleGetRequest_SSEStream(t *testing.T) {
	// Create filter
	mcpFilter := createTestFilter(t)

	// Create GET request
	session, err := mcpFilter.sessionManager.CreateSession()
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}
	req := httptest.NewRequest(constant.Get, "/mcp", nil)
	req.Header.Set(constant.HeaderKeyAccept, constant.HeaderValueTextEventStream)
	req.Header.Set(constant.HeaderKeyMCPProtocolVersion, constant.MCPProtocolVersion20250618)
	req.Header.Set(constant.HeaderKeyMCPSessionId, session.ID)

	recorder := httptest.NewRecorder()
	ctx := createTestContext(req, recorder)
	mcpCtx := NewMCPContext(ctx)

	// Parse headers
	mcpCtx.ParseAndSetProtocolVersionHeader()
	mcpCtx.ParseAndSetSessionHeader()
	mcpCtx.ParseAndSetAcceptHeader()

	// Execute handleGetRequest
	status := mcpFilter.handleGetRequest(mcpCtx)

	// Verify filter status
	if status != filter.Stop {
		t.Errorf("Expected filter.Stop, got %v", status)
	}

	// Verify SourceResp is set
	if ctx.SourceResp == nil {
		t.Fatal("SourceResp should be set")
	}

	httpResp, ok := ctx.SourceResp.(*http.Response)
	if !ok {
		t.Fatal("SourceResp should be *http.Response")
	}

	// Verify response headers
	if httpResp.StatusCode != http.StatusOK {
		t.Errorf("Expected status 200, got %d", httpResp.StatusCode)
	}

	contentType := httpResp.Header.Get(constant.HeaderKeyContextType)
	if contentType != constant.HeaderValueTextEventStream {
		t.Errorf("Expected Content-Type %s, got %s", constant.HeaderValueTextEventStream, contentType)
	}

	sessionID := httpResp.Header.Get(constant.HeaderKeyMCPSessionId)
	if sessionID != session.ID {
		t.Errorf("Mcp-Session-Id should be the existing session, got %q", sessionID)
	}

	// Verify session stream was attached
	session, exists := mcpFilter.sessionManager.Session(sessionID)
	if !exists {
		t.Error("Session should exist")
	}
	if !session.HasPipeWriter() {
		t.Error("Session SSE stream should be attached")
	}

	// Verify response body is pipe reader
	if httpResp.Body == nil {
		t.Error("Response body should be set")
	}

	// Cleanup
	mcpFilter.sessionManager.Stop()
}

func TestHandleGetRequest_MissingAcceptHeader(t *testing.T) {
	mcpFilter := createTestFilter(t)
	defer mcpFilter.sessionManager.Stop()

	// Create GET request without Accept header
	req := httptest.NewRequest(constant.Get, "/mcp", nil)
	recorder := httptest.NewRecorder()
	ctx := createTestContext(req, recorder)
	mcpCtx := NewMCPContext(ctx)

	mcpCtx.ParseAndSetAcceptHeader()

	status := mcpFilter.handleGetRequest(mcpCtx)

	// Should return error
	if status != filter.Stop {
		t.Error("Expected filter.Stop for missing Accept header")
	}

	// Check for error response
	if recorder.Code != http.StatusNotAcceptable {
		t.Errorf("Expected status 406, got %d", recorder.Code)
	}
}

func TestHandleGetRequest_ResumeExistingSession(t *testing.T) {
	mcpFilter := createTestFilter(t)
	defer mcpFilter.sessionManager.Stop()

	// Create first session
	session1, _ := mcpFilter.sessionManager.CreateSession()
	sessionID := session1.ID

	// Create GET request with existing session ID
	req := httptest.NewRequest(constant.Get, "/mcp", nil)
	req.Header.Set(constant.HeaderKeyAccept, constant.HeaderValueTextEventStream)
	req.Header.Set(constant.HeaderKeyMCPSessionId, sessionID)

	recorder := httptest.NewRecorder()
	ctx := createTestContext(req, recorder)
	mcpCtx := NewMCPContext(ctx)

	mcpCtx.ParseAndSetSessionHeader()
	mcpCtx.ParseAndSetAcceptHeader()

	status := mcpFilter.handleGetRequest(mcpCtx)

	if status != filter.Stop {
		t.Errorf("Expected filter.Stop, got %v", status)
	}

	// Verify same session is reused
	session2, exists := mcpFilter.sessionManager.Session(sessionID)
	if !exists {
		t.Error("Session should exist")
	}
	if session1.ID != session2.ID {
		t.Error("Should reuse existing session")
	}
}

func TestHandleGetRequest_MissingSession(t *testing.T) {
	mcpFilter := createTestFilter(t)
	defer mcpFilter.sessionManager.Stop()

	req := httptest.NewRequest(constant.Get, "/mcp", nil)
	req.Header.Set(constant.HeaderKeyAccept, constant.HeaderValueTextEventStream)
	recorder := httptest.NewRecorder()
	mcpCtx := NewMCPContext(createTestContext(req, recorder))
	mcpCtx.ParseAndSetSessionHeader()
	mcpCtx.ParseAndSetAcceptHeader()

	status := mcpFilter.handleGetRequest(mcpCtx)
	if status != filter.Stop {
		t.Fatalf("expected filter.Stop, got %v", status)
	}
	if recorder.Code != http.StatusBadRequest {
		t.Fatalf("expected 400, got %d", recorder.Code)
	}
}

func TestHandleGetRequest_UnknownSession(t *testing.T) {
	mcpFilter := createTestFilter(t)
	defer mcpFilter.sessionManager.Stop()

	req := httptest.NewRequest(constant.Get, "/mcp", nil)
	req.Header.Set(constant.HeaderKeyAccept, constant.HeaderValueTextEventStream)
	req.Header.Set(constant.HeaderKeyMCPSessionId, "missing")
	recorder := httptest.NewRecorder()
	mcpCtx := NewMCPContext(createTestContext(req, recorder))
	mcpCtx.ParseAndSetSessionHeader()
	mcpCtx.ParseAndSetAcceptHeader()

	status := mcpFilter.handleGetRequest(mcpCtx)
	if status != filter.Stop {
		t.Fatalf("expected filter.Stop, got %v", status)
	}
	if recorder.Code != http.StatusNotFound {
		t.Fatalf("expected 404, got %d", recorder.Code)
	}
	if mcpFilter.sessionManager.ActiveSessionCount() != 0 {
		t.Fatalf("unknown GET created a session, count=%d", mcpFilter.sessionManager.ActiveSessionCount())
	}
}

func TestHandleGetRequest_NoRouterMissingSessionCreatesSession(t *testing.T) {
	mcpFilter := createNoRouterTestFilter(t)
	defer mcpFilter.sessionManager.Stop()

	req := httptest.NewRequest(constant.Get, "/mcp", nil)
	req.Header.Set(constant.HeaderKeyAccept, constant.HeaderValueTextEventStream)
	recorder := httptest.NewRecorder()
	ctx := createTestContext(req, recorder)
	mcpCtx := NewMCPContext(ctx)
	mcpCtx.ParseAndSetAcceptHeader()
	mcpCtx.ParseAndSetSessionHeader()

	status := mcpFilter.handleGetRequest(mcpCtx)

	if status != filter.Stop {
		t.Fatalf("expected filter.Stop, got %v", status)
	}
	if ctx.SourceResp == nil {
		t.Fatal("SourceResp should be set")
	}
	httpResp := ctx.SourceResp.(*http.Response)
	sessionID := httpResp.Header.Get(constant.HeaderKeyMCPSessionId)
	if sessionID == "" {
		t.Fatal("legacy GET should create a session when no router is configured")
	}
	if mcpFilter.sessionManager.ActiveSessionCount() != 1 {
		t.Fatalf("expected one active session, got %d", mcpFilter.sessionManager.ActiveSessionCount())
	}
}

func TestInitialize_NoRouterAllowsSuppliedSessionHeader(t *testing.T) {
	mcpFilter := createNoRouterTestFilter(t)
	defer mcpFilter.sessionManager.Stop()

	body := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","clientInfo":{"name":"client","version":"1.0"},"capabilities":{}}}`
	req := httptest.NewRequest(constant.Post, "/mcp", strings.NewReader(body))
	req.Header.Set(constant.HeaderKeyMCPSessionId, "caller-supplied")
	recorder := httptest.NewRecorder()
	ctx := NewMCPContext(createTestContext(req, recorder))
	ctx.ParseAndSetSessionHeader()
	ctx.SetMCPMethod(string(mcp.MethodInitialize))
	ctx.SetMCPRequestID(mcp.NewRequestId(int64(1)))

	status := mcpFilter.handleInitialize(ctx, mcp.JSONRPCRequest{
		Request: mcp.Request{Method: string(mcp.MethodInitialize)},
		ID:      mcp.NewRequestId(int64(1)),
	})

	if status != filter.Stop {
		t.Fatalf("expected filter.Stop, got %v", status)
	}
	if recorder.Code != http.StatusOK {
		t.Fatalf("expected 200, got %d", recorder.Code)
	}
	if recorder.Header().Get(constant.HeaderKeyMCPSessionId) == "" {
		t.Fatal("initialize should return a server-issued session")
	}
}

func TestInitialize_NoRouterReusesExistingSuppliedSessionHeader(t *testing.T) {
	mcpFilter := createNoRouterTestFilter(t)
	defer mcpFilter.sessionManager.Stop()

	existing, err := mcpFilter.sessionManager.CreateSession()
	require.NoError(t, err)

	body := `{"jsonrpc":"2.0","id":1,"method":"initialize","params":{"protocolVersion":"2025-06-18","clientInfo":{"name":"client","version":"1.0"},"capabilities":{}}}`
	req := httptest.NewRequest(constant.Post, "/mcp", strings.NewReader(body))
	req.Header.Set(constant.HeaderKeyMCPSessionId, existing.ID)
	recorder := httptest.NewRecorder()
	ctx := NewMCPContext(createTestContext(req, recorder))
	ctx.ParseAndSetSessionHeader()
	ctx.SetMCPMethod(string(mcp.MethodInitialize))
	ctx.SetMCPRequestID(mcp.NewRequestId(int64(1)))

	status := mcpFilter.handleInitialize(ctx, mcp.JSONRPCRequest{
		Request: mcp.Request{Method: string(mcp.MethodInitialize)},
		ID:      mcp.NewRequestId(int64(1)),
	})

	require.Equal(t, filter.Stop, status)
	require.Equal(t, http.StatusOK, recorder.Code)
	assert.Equal(t, existing.ID, recorder.Header().Get(constant.HeaderKeyMCPSessionId))
	assert.Equal(t, 1, mcpFilter.sessionManager.ActiveSessionCount())
}

func TestMaintainSSEPipe_ContextCancellationDetachesStream(t *testing.T) {
	mcpFilter := createTestFilter(t)
	defer mcpFilter.sessionManager.Stop()

	session, _ := mcpFilter.sessionManager.CreateSession()
	pipeReader, pipeWriter := io.Pipe()
	defer pipeReader.Close()
	token, _ := session.AttachStream(pipeWriter)

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	httpCtx := &contexthttp.HttpContext{
		Ctx: ctx,
	}
	mcpCtx := NewMCPContext(httpCtx)

	go mcpFilter.maintainSSEPipe(mcpCtx, session, token)

	cancel()
	waitForNoPipeWriter(t, session)

	_, exists := mcpFilter.sessionManager.Session(session.ID)
	if !exists {
		t.Error("Session should remain after context cancellation")
	}
	if session.HasPipeWriter() {
		t.Error("SSE stream should be detached after context cancellation")
	}
}

func TestHandleGetRequest_ReconnectFlushesPendingListChanged(t *testing.T) {
	mcpFilter := createTestFilter(t)
	defer mcpFilter.sessionManager.Stop()

	session, _ := mcpFilter.sessionManager.CreateSession()
	session.MarkToolsListChangedPending()

	req := httptest.NewRequest(constant.Get, "/mcp", nil)
	req.Header.Set(constant.HeaderKeyAccept, constant.HeaderValueTextEventStream)
	req.Header.Set(constant.HeaderKeyMCPSessionId, session.ID)
	recorder := httptest.NewRecorder()
	ctx := createTestContext(req, recorder)
	mcpCtx := NewMCPContext(ctx)
	mcpCtx.ParseAndSetSessionHeader()
	mcpCtx.ParseAndSetAcceptHeader()

	status := mcpFilter.handleGetRequest(mcpCtx)
	if status != filter.Stop {
		t.Fatalf("expected filter.Stop, got %v", status)
	}

	httpResp, ok := ctx.SourceResp.(*http.Response)
	if !ok {
		t.Fatal("SourceResp should be *http.Response")
	}
	buf := make([]byte, 512)
	n, err := httpResp.Body.Read(buf)
	if err != nil {
		t.Fatalf("expected pending notification on reconnect: %v", err)
	}
	notification := string(buf[:n])
	if !strings.Contains(notification, toolsListChangedMethod) {
		t.Fatalf("expected tools/list_changed notification, got %q", notification)
	}
	if strings.Contains(notification, `"id"`) {
		t.Fatalf("notification must not contain JSON-RPC id: %q", notification)
	}
}

func TestHandleGetRequest_OverlappingStreamsOldDetachDoesNotClearNew(t *testing.T) {
	mcpFilter := createTestFilter(t)
	defer mcpFilter.sessionManager.Stop()

	session, _ := mcpFilter.sessionManager.CreateSession()

	ctxA, cancelA := context.WithCancel(context.Background())
	defer cancelA()
	reqA := httptest.NewRequest(constant.Get, "/mcp", nil)
	reqA.Header.Set(constant.HeaderKeyAccept, constant.HeaderValueTextEventStream)
	reqA.Header.Set(constant.HeaderKeyMCPSessionId, session.ID)
	httpCtxA := createTestContext(reqA, httptest.NewRecorder())
	httpCtxA.Ctx = ctxA
	mcpCtxA := NewMCPContext(httpCtxA)
	mcpCtxA.ParseAndSetSessionHeader()
	mcpCtxA.ParseAndSetAcceptHeader()
	if status := mcpFilter.handleGetRequest(mcpCtxA); status != filter.Stop {
		t.Fatalf("expected first GET to stop, got %v", status)
	}

	ctxB, cancelB := context.WithCancel(context.Background())
	defer cancelB()
	reqB := httptest.NewRequest(constant.Get, "/mcp", nil)
	reqB.Header.Set(constant.HeaderKeyAccept, constant.HeaderValueTextEventStream)
	reqB.Header.Set(constant.HeaderKeyMCPSessionId, session.ID)
	httpCtxB := createTestContext(reqB, httptest.NewRecorder())
	httpCtxB.Ctx = ctxB
	mcpCtxB := NewMCPContext(httpCtxB)
	mcpCtxB.ParseAndSetSessionHeader()
	mcpCtxB.ParseAndSetAcceptHeader()
	if status := mcpFilter.handleGetRequest(mcpCtxB); status != filter.Stop {
		t.Fatalf("expected second GET to stop, got %v", status)
	}

	cancelA()
	runtime.Gosched()
	if !session.HasPipeWriter() {
		t.Fatal("old stream goroutine detached the replacement stream")
	}

	respB, ok := httpCtxB.SourceResp.(*http.Response)
	if !ok {
		t.Fatal("second SourceResp should be *http.Response")
	}
	notificationCh := make(chan string, 1)
	go func() {
		buf := make([]byte, 512)
		n, err := respB.Body.Read(buf)
		if err == nil {
			notificationCh <- string(buf[:n])
		}
	}()

	mcpFilter.notifyToolsListChanged(session.ID)
	select {
	case notification := <-notificationCh:
		if !strings.Contains(notification, toolsListChangedMethod) {
			t.Fatalf("expected notification on replacement stream, got %q", notification)
		}
	case <-time.After(time.Second):
		t.Fatal("replacement stream did not receive notification")
	}
}

func TestSendServerNotification(t *testing.T) {
	mcpFilter := createTestFilter(t)
	defer mcpFilter.sessionManager.Stop()

	// Create session with pipe
	session, _ := mcpFilter.sessionManager.CreateSession()
	pipeReader, pipeWriter := io.Pipe()
	_, _ = session.AttachStream(pipeWriter)
	sessionID := session.ID

	// Send notification in goroutine
	params := map[string]any{
		"message": "test notification",
		"value":   123,
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- mcpFilter.SendServerNotification(sessionID, "notifications/test", params)
	}()

	// Read from pipe
	buf := make([]byte, 1024)
	n, err := pipeReader.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}

	received := string(buf[:n])

	// Verify SSE format
	if !strings.HasPrefix(received, "data:") {
		t.Error("Should start with 'data:'")
	}
	if !strings.HasSuffix(received, "\n\n") {
		t.Error("Should end with \\n\\n")
	}
	if !strings.Contains(received, "notifications/test") {
		t.Error("Should contain method name")
	}
	if !strings.Contains(received, "jsonrpc") {
		t.Error("Should contain jsonrpc field")
	}

	// Verify no error
	if err := <-errCh; err != nil {
		t.Errorf("SendServerNotification failed: %v", err)
	}

	pipeWriter.Close()
	pipeReader.Close()
}

func TestSendServerNotification_NoSession(t *testing.T) {
	mcpFilter := createTestFilter(t)
	defer mcpFilter.sessionManager.Stop()

	err := mcpFilter.SendServerNotification("non-existent-id", "test", map[string]any{})
	if err == nil {
		t.Error("Expected error for non-existent session")
	}
	if !strings.Contains(err.Error(), "session not found") {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestSendServerRequest(t *testing.T) {
	mcpFilter := createTestFilter(t)
	defer mcpFilter.sessionManager.Stop()

	// Create session with pipe
	session, _ := mcpFilter.sessionManager.CreateSession()
	pipeReader, pipeWriter := io.Pipe()
	_, _ = session.AttachStream(pipeWriter)
	sessionID := session.ID

	// Send request in goroutine
	params := map[string]any{"key": "value"}

	errCh := make(chan error, 1)
	go func() {
		errCh <- mcpFilter.SendServerRequest(sessionID, 123, "prompts/get", params)
	}()

	// Read from pipe
	buf := make([]byte, 1024)
	n, err := pipeReader.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}

	received := string(buf[:n])

	// Verify SSE format
	if !strings.HasPrefix(received, "data:") {
		t.Error("Should start with 'data:'")
	}
	if !strings.Contains(received, "prompts/get") {
		t.Error("Should contain method name")
	}
	if !strings.Contains(received, `"id"`) {
		t.Error("Should contain id field")
	}

	// Verify no error
	if err := <-errCh; err != nil {
		t.Errorf("SendServerRequest failed: %v", err)
	}

	pipeWriter.Close()
	pipeReader.Close()
}

// Note: TestValidateMCPProtocolVersion removed - version negotiation now happens
// during initialize request/response per MCP spec, not at HTTP header level.

// Helper functions

func createTestFilter(t *testing.T) *MCPServerFilter {
	return createTestFilterWithRouter(t, &model.RouterConfig{})
}

func createNoRouterTestFilter(t *testing.T) *MCPServerFilter {
	return createTestFilterWithRouter(t, nil)
}

func createTestFilterWithRouter(t *testing.T, routerCfg *model.RouterConfig) *MCPServerFilter {
	cfg := &model.McpServerConfig{
		ServerInfo: model.ServerInfo{
			Name:    "Test Server",
			Version: "1.0.0",
		},
		Endpoint: "/mcp",
		Tools:    []model.ToolConfig{},
		Router:   routerCfg,
	}

	factory := &FilterFactory{cfg: cfg}
	if err := factory.Apply(); err != nil {
		t.Fatalf("Failed to apply filter factory: %v", err)
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

func waitForNoPipeWriter(t *testing.T, session *transport.MCPSession) {
	t.Helper()
	deadline := time.After(time.Second)
	for session.HasPipeWriter() {
		select {
		case <-deadline:
			t.Fatal("SSE stream was not detached")
		default:
			runtime.Gosched()
		}
	}
}

func createTestContext(req *http.Request, recorder *httptest.ResponseRecorder) *contexthttp.HttpContext {
	return &contexthttp.HttpContext{
		Request: req,
		Writer:  recorder,
		Ctx:     context.Background(),
		Params:  make(map[string]any),
	}
}

// TestSendMCPResponse_ClearsContentLength verifies that sendMCPResponse
// clears the stale Content-Length header (which was copied from the backend
// response by buildTargetResponse) before returning filter.Continue on the
// JSON path. Without this, the larger MCP-wrapped body would exceed the
// declared Content-Length and Go's net/http would abort the connection.
func TestSendMCPResponse_ClearsContentLength(t *testing.T) {
	mcpFilter := createTestFilter(t)

	req := httptest.NewRequest(http.MethodPost, "/mcp", nil)
	rec := httptest.NewRecorder()
	ctx := createTestContext(req, rec)
	mcpCtx := NewMCPContext(ctx)

	// Simulate what buildTargetResponse does: set a stale Content-Length
	// from a hypothetical backend response (e.g. backend returned 47 bytes).
	rec.Header().Set("Content-Length", "47")

	// Build a response whose wrapped JSON body will be larger than 47 bytes.
	resp := mcpFilter.responseBuilder.ToolCallSuccess(
		"test-id",
		`{"name":"test","age":30}`,
	)

	status := mcpFilter.sendMCPResponse(mcpCtx, resp)

	// Must return Continue for the JSON path (no active SSE session).
	if status != filter.Continue {
		t.Fatalf("expected filter.Continue, got %v", status)
	}

	// The stale Content-Length header must be cleared.
	if cl := rec.Header().Get("Content-Length"); cl != "" {
		t.Errorf("Content-Length should have been cleared, but got %q", cl)
	}

	// Sanity: the wrapped body must be larger than the old Content-Length,
	// otherwise the test scenario is not meaningful.
	if unary, ok := ctx.TargetResp.(*client.UnaryResponse); ok {
		if len(unary.Data) <= 47 {
			t.Errorf("wrapped body too short for test (%d bytes)", len(unary.Data))
		}
	}
}

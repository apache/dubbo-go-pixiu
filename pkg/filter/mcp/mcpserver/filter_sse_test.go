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
	"strings"
	"testing"
	"time"
)

import (
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
	req := httptest.NewRequest(constant.Get, "/mcp", nil)
	req.Header.Set(constant.HeaderKeyAccept, constant.HeaderValueTextEventStream)
	req.Header.Set(constant.HeaderKeyMCPProtocolVersion, constant.MCPProtocolVersion20250618)

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
	if sessionID == "" {
		t.Error("Mcp-Session-Id should be set")
	}

	// Verify session was created
	session, exists := mcpFilter.sessionManager.Session(sessionID)
	if !exists {
		t.Error("Session should be created")
	}
	if session.PipeWriter == nil {
		t.Error("Session PipeWriter should be set")
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
	session1, _ := mcpFilter.sessionManager.EnsureSession("")
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

func TestMaintainSSEPipe_Keepalive(t *testing.T) {
	mcpFilter := createTestFilter(t)
	defer mcpFilter.sessionManager.Stop()

	// Create session with pipe
	session, _ := mcpFilter.sessionManager.EnsureSession("")
	pipeReader, pipeWriter := io.Pipe()
	session.PipeWriter = pipeWriter

	// Create context with timeout
	ctx, cancel := context.WithTimeout(context.Background(), 2*time.Second)
	defer cancel()

	httpCtx := &contexthttp.HttpContext{
		Ctx: ctx,
	}
	mcpCtx := NewMCPContext(httpCtx)

	// Start maintainSSEPipe with very short interval for testing
	// Note: This test uses production KeepaliveInterval (30s), so we won't actually receive keepalive in 2s
	go mcpFilter.maintainSSEPipe(mcpCtx, session)

	// Read from pipe in background
	dataCh := make(chan string, 10)
	go func() {
		buf := make([]byte, 1024)
		for {
			n, err := pipeReader.Read(buf)
			if err != nil {
				return
			}
			if n > 0 {
				dataCh <- string(buf[:n])
			}
		}
	}()

	// Wait a bit to see if we receive anything (we shouldn't with 30s interval)
	select {
	case data := <-dataCh:
		// If we somehow received data (unlikely with 30s interval), verify it's keepalive
		if !strings.HasPrefix(data, ":") {
			t.Errorf("Expected keepalive comment, got: %s", data)
		}
	case <-time.After(500 * time.Millisecond):
		// Expected: no data yet due to 30s interval
	}

	// Cancel context to stop maintenance
	cancel()
	time.Sleep(100 * time.Millisecond)

	// Verify session was cleaned up
	_, exists := mcpFilter.sessionManager.Session(session.ID)
	if exists {
		t.Error("Session should be removed after context cancellation")
	}
}

func TestSendServerNotification(t *testing.T) {
	mcpFilter := createTestFilter(t)
	defer mcpFilter.sessionManager.Stop()

	// Create session with pipe
	session, _ := mcpFilter.sessionManager.EnsureSession("")
	pipeReader, pipeWriter := io.Pipe()
	session.PipeWriter = pipeWriter
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
	session, _ := mcpFilter.sessionManager.EnsureSession("")
	pipeReader, pipeWriter := io.Pipe()
	session.PipeWriter = pipeWriter
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
	cfg := &model.McpServerConfig{
		ServerInfo: model.ServerInfo{
			Name:    "Test Server",
			Version: "1.0.0",
		},
		Endpoint: "/mcp",
		Tools:    []model.ToolConfig{},
	}

	factory := &FilterFactory{cfg: cfg}
	if err := factory.Apply(); err != nil {
		t.Fatalf("Failed to apply filter factory: %v", err)
	}

	sessionManager := transport.NewSessionManager()
	sseHandler := transport.NewSSEHandler(sessionManager)
	contentNegotiator := transport.NewContentNegotiator()

	return &MCPServerFilter{
		cfg:               cfg,
		registry:          factory.registry,
		errorHandler:      NewErrorHandler(),
		responseBuilder:   NewResponseBuilder(),
		sessionManager:    sessionManager,
		sseHandler:        sseHandler,
		contentNegotiator: contentNegotiator,
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

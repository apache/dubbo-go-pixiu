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

package transport

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strings"
	"time"

	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	SSEKeepaliveInterval = 30 * time.Second
	SSEWriteTimeout      = 5 * time.Second
)

// SSEHandler handles Server-Sent Events for MCP
type SSEHandler struct {
	sessionManager *SessionManager
}

// NewSSEHandler creates a new SSE handler
func NewSSEHandler(sessionManager *SessionManager) *SSEHandler {
	return &SSEHandler{
		sessionManager: sessionManager,
	}
}

// EstablishSSEConnection establishes an SSE connection for the given session
func (h *SSEHandler) EstablishSSEConnection(w http.ResponseWriter, r *http.Request, session *MCPSession) error {
	// Validate that client accepts SSE
	acceptHeader := r.Header.Get(constant.HeaderKeyAccept)
	if !strings.Contains(acceptHeader, constant.HeaderValueTextEventStream) {
		return fmt.Errorf("client does not accept text/event-stream")
	}

	// Check if ResponseWriter supports flushing
	flusher, ok := w.(http.Flusher)
	if !ok {
		return fmt.Errorf("response writer does not support flushing")
	}

	// Set SSE headers
	h.setSSEHeaders(w, session.ID)

	// Store writer and flusher in session
	session.SSEWriter = w
	session.SSEFlusher = flusher

	// Start connection maintenance
	go h.maintainSSEConnection(session)

	logger.Infof("[dubbo-go-pixiu] mcp server established SSE connection for session: %s", session.ID)
	return nil
}

// SendSSEMessage sends a message through the SSE connection
func (h *SSEHandler) SendSSEMessage(session *MCPSession, message any) error {
	if session.SSEWriter == nil || session.SSEFlusher == nil {
		return fmt.Errorf("SSE connection not established")
	}

	// Marshal message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal SSE message: %w", err)
	}

	// Format as SSE event
	sseData := h.formatSSEMessage(string(messageJSON))

	// Write with timeout protection
	done := make(chan error, 1)
	go func() {
		if _, err := session.SSEWriter.Write([]byte(sseData)); err != nil {
			done <- err
			return
		}
		session.SSEFlusher.Flush()
		done <- nil
	}()

	select {
	case err := <-done:
		if err != nil {
			logger.Errorf("[dubbo-go-pixiu] mcp server failed to send SSE message: %v", err)
			return err
		}
	case <-time.After(SSEWriteTimeout):
		logger.Errorf("[dubbo-go-pixiu] mcp server SSE write timeout for session: %s", session.ID)
		return fmt.Errorf("SSE write timeout")
	}

	session.LastActivity = time.Now()
	logger.Debugf("[dubbo-go-pixiu] mcp server sent SSE message to session: %s", session.ID)
	return nil
}

// SendSSEEvent sends a custom SSE event
func (h *SSEHandler) SendSSEEvent(session *MCPSession, eventType, data string) error {
	if session.SSEWriter == nil || session.SSEFlusher == nil {
		return fmt.Errorf("SSE connection not established")
	}

	sseData := h.formatSSEEvent(eventType, data)

	// Write with timeout protection
	done := make(chan error, 1)
	go func() {
		if _, err := session.SSEWriter.Write([]byte(sseData)); err != nil {
			done <- err
			return
		}
		session.SSEFlusher.Flush()
		done <- nil
	}()

	select {
	case err := <-done:
		if err != nil {
			logger.Errorf("[dubbo-go-pixiu] mcp server failed to send SSE event: %v", err)
			return err
		}
	case <-time.After(SSEWriteTimeout):
		logger.Errorf("[dubbo-go-pixiu] mcp server SSE event write timeout for session: %s", session.ID)
		return fmt.Errorf("SSE event write timeout")
	}

	session.LastActivity = time.Now()
	return nil
}

// setSSEHeaders sets the necessary SSE response headers
func (h *SSEHandler) setSSEHeaders(w http.ResponseWriter, sessionID string) {
	w.Header().Set(constant.HeaderKeyContextType, constant.HeaderValueTextEventStream)
	w.Header().Set(constant.HeaderKeyCacheControl, constant.HeaderValueNoCache)
	w.Header().Set(constant.HeaderKeyConnection, constant.HeaderValueKeepAlive)
	w.Header().Set(constant.HeaderKeyMCPSessionId, sessionID)
	w.Header().Set(constant.HeaderKeyAccessControlAllowOrigin, constant.HeaderValueAll)
	w.Header().Set(constant.HeaderKeyAccessControlExposeHeaders, constant.HeaderKeyMCPSessionId)
	w.WriteHeader(http.StatusOK)
}

// formatSSEMessage formats a message as SSE data
func (h *SSEHandler) formatSSEMessage(messageJSON string) string {
	return fmt.Sprintf("data: %s\n\n", messageJSON)
}

// formatSSEEvent formats a custom SSE event
func (h *SSEHandler) formatSSEEvent(eventType, data string) string {
	if eventType != "" {
		return fmt.Sprintf("event: %s\ndata: %s\n\n", eventType, data)
	}
	return fmt.Sprintf("data: %s\n\n", data)
}

// maintainSSEConnection maintains the SSE connection with keepalive
func (h *SSEHandler) maintainSSEConnection(session *MCPSession) {
	ticker := time.NewTicker(SSEKeepaliveInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			// Send keepalive comment (ignored by SSE clients)
			if err := h.sendKeepalive(session); err != nil {
				logger.Errorf("[dubbo-go-pixiu] mcp server keepalive failed for session %s: %v", session.ID, err)
				h.sessionManager.RemoveSession(session.ID)
				return
			}

		case <-session.Done:
			logger.Debugf("[dubbo-go-pixiu] mcp server SSE connection closed for session: %s", session.ID)
			return
		}
	}
}

// sendKeepalive sends a keepalive comment to maintain the connection
func (h *SSEHandler) sendKeepalive(session *MCPSession) error {
	if session.SSEWriter == nil || session.SSEFlusher == nil {
		return fmt.Errorf("SSE connection not established")
	}

	keepaliveComment := fmt.Sprintf(": keepalive %d\n\n", time.Now().Unix())

	done := make(chan error, 1)
	go func() {
		if _, err := session.SSEWriter.Write([]byte(keepaliveComment)); err != nil {
			done <- err
			return
		}
		session.SSEFlusher.Flush()
		done <- nil
	}()

	select {
	case err := <-done:
		return err
	case <-time.After(SSEWriteTimeout):
		return fmt.Errorf("keepalive write timeout")
	}
}

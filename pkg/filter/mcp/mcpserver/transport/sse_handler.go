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
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// SSEHandler provides utility functions for SSE message formatting
type SSEHandler struct {
	sessionManager *SessionManager
}

// NewSSEHandler creates a new SSE handler
func NewSSEHandler(sessionManager *SessionManager) *SSEHandler {
	return &SSEHandler{
		sessionManager: sessionManager,
	}
}

// SendSSEMessage sends a message through the SSE pipe
func (h *SSEHandler) SendSSEMessage(session *MCPSession, message any) error {
	if session.PipeWriter == nil {
		return fmt.Errorf("SSE pipe not established")
	}

	// Marshal message to JSON
	messageJSON, err := json.Marshal(message)
	if err != nil {
		return fmt.Errorf("failed to marshal SSE message: %w", err)
	}

	// Format as SSE event
	sseData := h.FormatSSEMessage(string(messageJSON))

	// Write to pipe
	if _, err := session.PipeWriter.Write([]byte(sseData)); err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp server failed to send SSE message: %v", err)
		return fmt.Errorf("failed to write to SSE pipe: %w", err)
	}

	// Update LastActivity using session's mutex
	session.mu.Lock()
	session.LastActivity = time.Now()
	session.mu.Unlock()

	logger.Debugf("[dubbo-go-pixiu] mcp server sent SSE message to session: %s", session.ID)
	return nil
}

// FormatSSEMessage formats a message as SSE data
func (h *SSEHandler) FormatSSEMessage(messageJSON string) string {
	return fmt.Sprintf("data: %s\n\n", messageJSON)
}

// FormatSSEEvent formats a custom SSE event with event type
func (h *SSEHandler) FormatSSEEvent(eventType, data string) string {
	if eventType != "" {
		return fmt.Sprintf("event: %s\ndata: %s\n\n", eventType, data)
	}
	return fmt.Sprintf("data: %s\n\n", data)
}

// FormatSSEEventWithID formats a SSE event with ID for resumability
func (h *SSEHandler) FormatSSEEventWithID(eventID int64, data string) string {
	return fmt.Sprintf("id: %d\ndata: %s\n\n", eventID, data)
}

// FormatSSEKeepalive formats a SSE keepalive comment
func (h *SSEHandler) FormatSSEKeepalive(timestamp int64) string {
	return fmt.Sprintf(": keepalive %d\n\n", timestamp)
}

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
	"io"
	"testing"
	"time"
)

func TestNewSSEHandler(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()

	handler := NewSSEHandler(sm)
	if handler == nil {
		t.Fatal("NewSSEHandler returned nil")
	}
	if handler.sessionManager == nil {
		t.Error("sessionManager not initialized")
	}
}

func TestFormatSSEMessage(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()
	handler := NewSSEHandler(sm)

	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "simple message",
			input:    `{"jsonrpc":"2.0","method":"test"}`,
			expected: "data: {\"jsonrpc\":\"2.0\",\"method\":\"test\"}\n\n",
		},
		{
			name:     "empty message",
			input:    "",
			expected: "data: \n\n",
		},
		{
			name:     "message with newlines",
			input:    "line1\nline2",
			expected: "data: line1\nline2\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.FormatSSEMessage(tt.input)
			if result != tt.expected {
				t.Errorf("FormatSSEMessage() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatSSEEvent(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()
	handler := NewSSEHandler(sm)

	tests := []struct {
		name      string
		eventType string
		data      string
		expected  string
	}{
		{
			name:      "with event type",
			eventType: "message",
			data:      "test data",
			expected:  "event: message\ndata: test data\n\n",
		},
		{
			name:      "without event type",
			eventType: "",
			data:      "test data",
			expected:  "data: test data\n\n",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handler.FormatSSEEvent(tt.eventType, tt.data)
			if result != tt.expected {
				t.Errorf("FormatSSEEvent() = %q, want %q", result, tt.expected)
			}
		})
	}
}

func TestFormatSSEEventWithID(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()
	handler := NewSSEHandler(sm)

	result := handler.FormatSSEEventWithID(123, "test data")
	expected := "id: 123\ndata: test data\n\n"
	if result != expected {
		t.Errorf("FormatSSEEventWithID() = %q, want %q", result, expected)
	}
}

func TestFormatSSEKeepalive(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()
	handler := NewSSEHandler(sm)

	timestamp := time.Now().Unix()
	result := handler.FormatSSEKeepalive(timestamp)

	// Just check format is correct (contains keepalive and timestamp)
	if len(result) < 10 {
		t.Error("Keepalive format too short")
	}
	if result[0] != ':' {
		t.Error("Keepalive should start with ':'")
	}
	if result[len(result)-2:] != "\n\n" {
		t.Error("Keepalive should end with \\n\\n")
	}
}

func TestSendSSEMessage(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()
	handler := NewSSEHandler(sm)

	// Create session with pipe
	session, _ := sm.EnsureSession("")
	pipeReader, pipeWriter := io.Pipe()
	session.PipeWriter = pipeWriter

	// Send message in goroutine
	message := map[string]any{
		"jsonrpc": "2.0",
		"method":  "test",
	}

	errCh := make(chan error, 1)
	go func() {
		errCh <- handler.SendSSEMessage(session, message)
	}()

	// Read from pipe
	buf := make([]byte, 1024)
	n, err := pipeReader.Read(buf)
	if err != nil {
		t.Fatalf("Failed to read from pipe: %v", err)
	}

	received := string(buf[:n])
	if len(received) == 0 {
		t.Error("Should receive data")
	}
	if received[:5] != "data:" {
		t.Error("Should start with 'data:'")
	}
	if received[len(received)-2:] != "\n\n" {
		t.Error("Should end with \\n\\n")
	}

	// Check send error
	if err := <-errCh; err != nil {
		t.Errorf("SendSSEMessage failed: %v", err)
	}

	pipeWriter.Close()
	pipeReader.Close()
}

func TestSendSSEMessage_NoPipe(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()
	handler := NewSSEHandler(sm)

	// Create session without pipe
	session, _ := sm.EnsureSession("")

	message := map[string]any{"test": "data"}
	err := handler.SendSSEMessage(session, message)
	if err == nil {
		t.Error("Expected error when PipeWriter is nil")
	}
	if err.Error() != "SSE pipe not established" {
		t.Errorf("Unexpected error message: %v", err)
	}
}

func TestSendSSEMessage_InvalidJSON(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()
	handler := NewSSEHandler(sm)

	session, _ := sm.EnsureSession("")
	pipeReader, pipeWriter := io.Pipe()
	session.PipeWriter = pipeWriter
	defer pipeReader.Close()
	defer pipeWriter.Close()

	// Try to send invalid message (channel cannot be marshaled)
	invalidMessage := make(chan int)
	err := handler.SendSSEMessage(session, invalidMessage)
	if err == nil {
		t.Error("Expected error when marshaling invalid message")
	}
}

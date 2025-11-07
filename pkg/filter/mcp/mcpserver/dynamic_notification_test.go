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
	"strings"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestNotifyToolsListChanged(t *testing.T) {
	// Reset global state
	ResetGlobalState()
	defer ResetGlobalState()

	consumer := GetOrInitDynamicConsumer()
	sm := GetOrInitSessionManager()

	// Create session with SSE pipe
	session, _ := sm.EnsureSession("")
	pipeReader, pipeWriter := io.Pipe()
	session.PipeWriter = pipeWriter
	defer pipeReader.Close()
	defer pipeWriter.Close()

	// Apply config with tools (will trigger notification)
	config := createTestMcpServerConfig([]model.ToolConfig{
		{Name: "tool1", Cluster: "cluster1"},
	})

	// Read notification in background
	notificationCh := make(chan string, 1)
	go func() {
		buf := make([]byte, 1024)
		n, err := pipeReader.Read(buf)
		if err == nil && n > 0 {
			notificationCh <- string(buf[:n])
		}
	}()

	// Apply config
	consumer.ResetDebounceState()
	err := consumer.ApplyMcpServerConfigByServer("server1", config)
	assert.NoError(t, err)

	// Wait for notification
	select {
	case notification := <-notificationCh:
		// Verify notification format
		assert.True(t, strings.HasPrefix(notification, "data:"), "Should start with 'data:'")
		assert.True(t, strings.HasSuffix(notification, "\n\n"), "Should end with \\n\\n")
		assert.True(t, strings.Contains(notification, "notifications/tools/list_changed"), "Should contain method name")
		assert.True(t, strings.Contains(notification, "jsonrpc"), "Should contain jsonrpc field")
	case <-time.After(1 * time.Second):
		t.Fatal("Did not receive notification within timeout")
	}
}

func TestNotifyToolsListChanged_MultipleSessions(t *testing.T) {
	ResetGlobalState()
	defer ResetGlobalState()

	consumer := GetOrInitDynamicConsumer()
	sm := GetOrInitSessionManager()

	// Create multiple sessions with pipes
	numSessions := 3
	readers := make([]io.ReadCloser, numSessions)
	writers := make([]*io.PipeWriter, numSessions)
	notificationChs := make([]chan string, numSessions)

	for i := 0; i < numSessions; i++ {
		session, _ := sm.EnsureSession("")
		pipeReader, pipeWriter := io.Pipe()
		session.PipeWriter = pipeWriter
		readers[i] = pipeReader
		writers[i] = pipeWriter
		defer pipeReader.Close()
		defer pipeWriter.Close()

		// Start reader for each session
		notificationChs[i] = make(chan string, 1)
		go func(idx int, reader io.Reader, ch chan string) {
			buf := make([]byte, 1024)
			n, err := reader.Read(buf)
			if err == nil && n > 0 {
				ch <- string(buf[:n])
			}
		}(i, pipeReader, notificationChs[i])
	}

	// Apply config
	config := createTestMcpServerConfig([]model.ToolConfig{
		{Name: "tool1", Cluster: "cluster1"},
	})

	consumer.ResetDebounceState()
	err := consumer.ApplyMcpServerConfigByServer("server1", config)
	assert.NoError(t, err)

	// Verify all sessions received notification
	receivedCount := 0
	for i := 0; i < numSessions; i++ {
		select {
		case notification := <-notificationChs[i]:
			assert.Contains(t, notification, "notifications/tools/list_changed")
			receivedCount++
		case <-time.After(1 * time.Second):
			t.Logf("Session %d did not receive notification", i)
		}
	}

	assert.Equal(t, numSessions, receivedCount, "All sessions should receive notification")
}

func TestNotifyToolsListChanged_NoActiveSessions(t *testing.T) {
	ResetGlobalState()
	defer ResetGlobalState()

	consumer := GetOrInitDynamicConsumer()

	// Apply config without any active sessions
	config := createTestMcpServerConfig([]model.ToolConfig{
		{Name: "tool1", Cluster: "cluster1"},
	})

	consumer.ResetDebounceState()
	err := consumer.ApplyMcpServerConfigByServer("server1", config)
	assert.NoError(t, err)

	// Should not panic or error, just skip notification
	assert.Equal(t, 0, GetOrInitSessionManager().ActiveSessionCount())
}

func TestNotifyToolsListChanged_DisconnectedSession(t *testing.T) {
	ResetGlobalState()
	defer ResetGlobalState()

	consumer := GetOrInitDynamicConsumer()
	sm := GetOrInitSessionManager()

	// Create session but don't attach pipe
	_, _ = sm.EnsureSession("")
	// session.PipeWriter is nil

	// Apply config
	config := createTestMcpServerConfig([]model.ToolConfig{
		{Name: "tool1", Cluster: "cluster1"},
	})

	consumer.ResetDebounceState()
	err := consumer.ApplyMcpServerConfigByServer("server1", config)
	assert.NoError(t, err)

	// Should handle disconnected session gracefully (logged as warning)
	// No panic or error
}

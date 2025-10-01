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
	"testing"
	"time"
)

func TestNewSessionManager(t *testing.T) {
	sm := NewSessionManager()
	if sm == nil {
		t.Fatal("NewSessionManager returned nil")
	}
	if sm.sessions == nil {
		t.Error("sessions map not initialized")
	}
	if sm.stopCh == nil {
		t.Error("stopCh not initialized")
	}
	sm.Stop()
}

func TestEnsureSession_CreateNew(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()

	// Test creating new session with empty header
	session, isNew := sm.EnsureSession("")
	if !isNew {
		t.Error("Expected new session to be created")
	}
	if session == nil {
		t.Fatal("Session should not be nil")
	}
	if session.ID == "" {
		t.Error("Session ID should not be empty")
	}
	if session.Done == nil {
		t.Error("Session Done channel should be initialized")
	}
	if len(session.ID) != 32 {
		t.Errorf("Expected session ID length 32, got %d", len(session.ID))
	}
}

func TestEnsureSession_ReuseExisting(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()

	// Create first session
	session1, isNew1 := sm.EnsureSession("")
	if !isNew1 {
		t.Error("Expected new session")
	}

	// Try to get same session by ID
	session2, isNew2 := sm.EnsureSession(session1.ID)
	if isNew2 {
		t.Error("Expected existing session, not new")
	}
	if session1.ID != session2.ID {
		t.Error("Should return same session")
	}

	// Verify last activity was updated
	if session2.LastActivity.Before(session1.LastActivity) {
		t.Error("LastActivity should be updated")
	}
}

func TestSession(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()

	// Create session
	session1, _ := sm.EnsureSession("")
	sessionID := session1.ID

	// Retrieve session
	session2, exists := sm.Session(sessionID)
	if !exists {
		t.Error("Session should exist")
	}
	if session2.ID != sessionID {
		t.Error("Retrieved wrong session")
	}

	// Try non-existent session
	_, exists = sm.Session("non-existent-id")
	if exists {
		t.Error("Non-existent session should not be found")
	}
}

func TestRemoveSession(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()

	// Create session
	session, _ := sm.EnsureSession("")
	sessionID := session.ID

	// Verify session exists
	_, exists := sm.Session(sessionID)
	if !exists {
		t.Fatal("Session should exist before removal")
	}

	// Remove session
	sm.RemoveSession(sessionID)

	// Verify session removed
	_, exists = sm.Session(sessionID)
	if exists {
		t.Error("Session should be removed")
	}

	// Verify Done channel closed
	select {
	case <-session.Done:
		// Expected: channel should be closed
	default:
		t.Error("Done channel should be closed")
	}

	// Test removing non-existent session (should not panic)
	sm.RemoveSession("non-existent-id")
}

func TestSessionCleanup(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()

	// Create session
	session, _ := sm.EnsureSession("")
	sessionID := session.ID

	// Manually set old LastActivity to simulate timeout
	sm.mu.Lock()
	session.LastActivity = time.Now().Add(-SessionTimeout - 1*time.Minute)
	sm.mu.Unlock()

	// Trigger cleanup
	sm.cleanupExpiredSessions()

	// Verify session was cleaned up
	_, exists := sm.Session(sessionID)
	if exists {
		t.Error("Expired session should be cleaned up")
	}
}

func TestSessionManager_Stop(t *testing.T) {
	sm := NewSessionManager()

	// Create multiple sessions
	session1, _ := sm.EnsureSession("")
	session2, _ := sm.EnsureSession("")

	// Stop manager
	sm.Stop()

	// Verify all sessions removed
	_, exists1 := sm.Session(session1.ID)
	_, exists2 := sm.Session(session2.ID)
	if exists1 || exists2 {
		t.Error("All sessions should be removed on Stop")
	}

	// Verify Done channels closed
	select {
	case <-session1.Done:
		// Expected
	default:
		t.Error("Session1 Done should be closed")
	}

	select {
	case <-session2.Done:
		// Expected
	default:
		t.Error("Session2 Done should be closed")
	}

	// Verify stop channel closed
	select {
	case <-sm.stopCh:
		// Expected
	default:
		t.Error("stopCh should be closed")
	}
}

func TestGenerateSessionID(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()

	// Generate multiple session IDs
	ids := make(map[string]bool)
	for i := 0; i < 100; i++ {
		id := sm.generateSessionID()
		if id == "" {
			t.Error("Generated session ID should not be empty")
		}
		if len(id) != 32 {
			t.Errorf("Expected session ID length 32, got %d", len(id))
		}
		if ids[id] {
			t.Errorf("Duplicate session ID generated: %s", id)
		}
		ids[id] = true
	}
}

func TestConcurrentSessionAccess(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()

	// Create initial session
	session, _ := sm.EnsureSession("")
	sessionID := session.ID

	done := make(chan bool, 3)

	// Concurrent reads
	go func() {
		for i := 0; i < 100; i++ {
			sm.Session(sessionID)
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()

	// Concurrent writes (update LastActivity)
	go func() {
		for i := 0; i < 100; i++ {
			sm.Session(sessionID)
			time.Sleep(time.Millisecond)
		}
		done <- true
	}()

	// Concurrent session creation
	go func() {
		for i := 0; i < 10; i++ {
			sm.EnsureSession("")
			time.Sleep(10 * time.Millisecond)
		}
		done <- true
	}()

	// Wait for all goroutines
	for i := 0; i < 3; i++ {
		<-done
	}
}

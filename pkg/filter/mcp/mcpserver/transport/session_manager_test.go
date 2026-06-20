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
	"sync"
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

func TestCreateSession(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()

	session, err := sm.CreateSession()
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
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

func TestGetSessionReusesExisting(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()

	session1, err := sm.CreateSession()
	if err != nil {
		t.Fatalf("CreateSession failed: %v", err)
	}

	session2, exists := sm.GetSession(session1.ID)
	if !exists {
		t.Fatal("expected existing session")
	}
	if session1.ID != session2.ID {
		t.Error("Should return same session")
	}

	// Verify last activity was updated
	session1.mu.RLock()
	lastActivity1 := session1.LastActivity
	session1.mu.RUnlock()

	session2.mu.RLock()
	lastActivity2 := session2.LastActivity
	session2.mu.RUnlock()

	if lastActivity2.Before(lastActivity1) {
		t.Error("LastActivity should be updated")
	}
}

func TestGetSessionUnknownDoesNotCreate(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()

	_, exists := sm.GetSession("missing")
	if exists {
		t.Fatal("unknown session should not exist")
	}
	if sm.ActiveSessionCount() != 0 {
		t.Fatalf("unknown lookup created a session, count=%d", sm.ActiveSessionCount())
	}
}

func TestSession(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()

	// Create session
	session1, _ := sm.CreateSession()
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
	session, _ := sm.CreateSession()
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

func TestSessionRemovedHandlerCalledOutsideLock(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()

	session, _ := sm.CreateSession()
	done := make(chan struct{})
	sm.AddSessionRemovedHandler(func(string) {
		// This would deadlock if the callback ran while the manager lock was held.
		_ = sm.ActiveSessionCount()
		close(done)
	})

	sm.RemoveSession(session.ID)

	select {
	case <-done:
	case <-time.After(time.Second):
		t.Fatal("session removal callback did not complete")
	}
}

func TestSessionRemovedHandlerPanicDoesNotBreakCleanup(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()

	session, _ := sm.CreateSession()
	called := make(chan string, 1)
	sm.AddSessionRemovedHandler(func(string) {
		panic("boom")
	})
	sm.AddSessionRemovedHandler(func(id string) {
		called <- id
	})

	sm.RemoveSession(session.ID)

	select {
	case id := <-called:
		if id != session.ID {
			t.Fatalf("unexpected callback session id %s", id)
		}
	case <-time.After(time.Second):
		t.Fatal("second callback was not invoked after panic")
	}
}

func TestRemoveSessionClosesBlockedSSEWrite(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()

	session, _ := sm.CreateSession()
	reader, writer := io.Pipe()
	defer reader.Close()
	session.AttachStream(writer)

	writeDone := make(chan error, 1)
	go func() {
		writeDone <- session.WriteSSEData([]byte("data: blocked\n\n"), time.Now())
	}()

	time.Sleep(50 * time.Millisecond)

	removeDone := make(chan struct{})
	go func() {
		sm.RemoveSession(session.ID)
		close(removeDone)
	}()

	select {
	case <-removeDone:
	case <-time.After(time.Second):
		t.Fatal("RemoveSession blocked while an SSE write was pending")
	}

	select {
	case err := <-writeDone:
		if err == nil {
			t.Fatal("blocked write should fail after session close")
		}
	case <-time.After(time.Second):
		t.Fatal("blocked write did not unblock after session close")
	}
}

func TestStreamTokenCannotWriteReplacementStream(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()

	session, _ := sm.CreateSession()
	readerA, writerA := io.Pipe()
	defer readerA.Close()
	tokenA := session.AttachStream(writerA)

	readerB, writerB := io.Pipe()
	defer readerB.Close()
	defer writerB.Close()
	tokenB := session.AttachStream(writerB)

	readB := make(chan string, 1)
	go func() {
		buf := make([]byte, 64)
		n, err := readerB.Read(buf)
		if err == nil {
			readB <- string(buf[:n])
		}
	}()

	if err := session.WriteSSEDataForStream(tokenA, []byte("data: old\n\n"), time.Now()); err == nil {
		t.Fatal("old stream token should not write after replacement")
	}
	select {
	case got := <-readB:
		t.Fatalf("old stream token wrote to replacement stream: %q", got)
	case <-time.After(50 * time.Millisecond):
	}

	writeDone := make(chan error, 1)
	go func() {
		writeDone <- session.WriteSSEDataForStream(tokenB, []byte("data: new\n\n"), time.Now())
	}()

	select {
	case got := <-readB:
		if got != "data: new\n\n" {
			t.Fatalf("unexpected replacement stream data: %q", got)
		}
	case <-time.After(time.Second):
		t.Fatal("replacement stream did not receive token-owned write")
	}
	if err := <-writeDone; err != nil {
		t.Fatalf("replacement token write failed: %v", err)
	}
}

func TestSessionCleanup(t *testing.T) {
	sm := NewSessionManager()
	defer sm.Stop()

	// Create session
	session, _ := sm.CreateSession()
	sessionID := session.ID

	// Manually set old LastActivity to simulate timeout
	session.mu.Lock()
	session.LastActivity = time.Now().Add(-SessionTimeout - 1*time.Minute)
	session.mu.Unlock()

	// Trigger cleanup
	sm.cleanupExpiredSessions()

	// Verify session was cleaned up
	_, exists := sm.Session(sessionID)
	if exists {
		t.Error("Expired session should be cleaned up")
	}
}

func TestSessionExpiredLookupRemovesAndCallsHandler(t *testing.T) {
	var mu sync.Mutex
	now := time.Unix(100, 0)
	sm := NewSessionManagerWithNow(func() time.Time {
		mu.Lock()
		defer mu.Unlock()
		return now
	})
	defer sm.Stop()

	session, _ := sm.CreateSession()
	removed := make(chan string, 1)
	sm.AddSessionRemovedHandler(func(id string) {
		removed <- id
	})

	mu.Lock()
	now = now.Add(SessionTimeout + time.Nanosecond)
	mu.Unlock()

	_, exists := sm.Session(session.ID)
	if exists {
		t.Fatal("expired session should not be returned")
	}

	select {
	case id := <-removed:
		if id != session.ID {
			t.Fatalf("unexpected removed session id %s", id)
		}
	case <-time.After(time.Second):
		t.Fatal("expected removal callback for expired session")
	}
}

func TestSessionManager_Stop(t *testing.T) {
	sm := NewSessionManager()

	// Create multiple sessions
	session1, _ := sm.CreateSession()
	session2, _ := sm.CreateSession()

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

func TestSessionManagerStopCallsRemovalHandlers(t *testing.T) {
	sm := NewSessionManager()

	session1, _ := sm.CreateSession()
	session2, _ := sm.CreateSession()
	removed := make(chan string, 2)
	sm.AddSessionRemovedHandler(func(id string) {
		removed <- id
	})

	sm.Stop()
	sm.Stop()

	got := map[string]bool{}
	for i := 0; i < 2; i++ {
		select {
		case id := <-removed:
			got[id] = true
		case <-time.After(time.Second):
			t.Fatal("expected removal callback on Stop")
		}
	}
	if !got[session1.ID] || !got[session2.ID] {
		t.Fatalf("missing stopped sessions in callback set: %#v", got)
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
	session, _ := sm.CreateSession()
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
			sm.CreateSession()
			time.Sleep(10 * time.Millisecond)
		}
		done <- true
	}()

	// Wait for all goroutines
	for i := 0; i < 3; i++ {
		<-done
	}
}

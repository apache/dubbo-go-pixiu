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
	"crypto/rand"
	"encoding/hex"
	"fmt"
	"io"
	"sync"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	SessionTimeout    = 30 * time.Minute
	CleanupInterval   = 5 * time.Minute
	KeepaliveInterval = 30 * time.Second
)

// SessionRemovedHandler is invoked after a transport session is removed.
// Handlers are called outside the SessionManager lock.
type SessionRemovedHandler func(sessionID string)

// MCPSession represents an active MCP session
type MCPSession struct {
	ID                       string
	CreatedAt                time.Time
	LastActivity             time.Time
	PipeWriter               *io.PipeWriter // Pipe writer for sending SSE messages
	Done                     chan struct{}
	SupportsToolsListChanged bool
	PendingToolsListChanged  bool
	mu                       sync.RWMutex // Protects mutable session fields
	writeMu                  sync.Mutex
	closeOnce                sync.Once
}

// SessionManager manages MCP sessions for SSE connections
type SessionManager struct {
	sessions        map[string]*MCPSession
	removedHandlers []SessionRemovedHandler
	mu              sync.RWMutex
	stopCh          chan struct{}
	once            sync.Once
	now             func() time.Time
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	return NewSessionManagerWithNow(time.Now)
}

// NewSessionManagerWithNow creates a session manager with an injected clock.
func NewSessionManagerWithNow(now func() time.Time) *SessionManager {
	if now == nil {
		now = time.Now
	}
	sm := &SessionManager{
		sessions: make(map[string]*MCPSession),
		stopCh:   make(chan struct{}),
		now:      now,
	}
	go sm.startCleanupRoutine()
	return sm
}

// AddSessionRemovedHandler registers a callback for session removal.
func (sm *SessionManager) AddSessionRemovedHandler(handler SessionRemovedHandler) {
	if handler == nil {
		return
	}
	sm.mu.Lock()
	defer sm.mu.Unlock()
	sm.removedHandlers = append(sm.removedHandlers, handler)
}

// EnsureSession gets existing session or creates new one
func (sm *SessionManager) EnsureSession(sessionIDHeader string) (*MCPSession, bool) {
	sm.mu.Lock()
	now := sm.now()

	// Try to get existing session
	if sessionIDHeader != "" {
		if session, exists := sm.sessions[sessionIDHeader]; exists {
			if sm.sessionExpiredLocked(session, now) {
				sm.removeSessionLocked(sessionIDHeader)
				handlers := sm.handlersLocked()
				sm.mu.Unlock()
				sm.callRemovedHandlers(handlers, []string{sessionIDHeader})
				return sm.EnsureSession("")
			}
			session.touch(now)
			sm.mu.Unlock()
			return session, false // existing session
		}
	}

	// Create new session
	sessionID := sm.generateSessionID()
	session := &MCPSession{
		ID:           sessionID,
		CreatedAt:    now,
		LastActivity: now,
		Done:         make(chan struct{}),
	}

	sm.sessions[sessionID] = session
	sm.mu.Unlock()
	logger.Infof("[dubbo-go-pixiu] mcp server created new session: %s", sessionID)
	return session, true // new session
}

// Session retrieves a session by ID
func (sm *SessionManager) Session(sessionID string) (*MCPSession, bool) {
	if sessionID == "" {
		return nil, false
	}
	now := sm.now()
	sm.mu.Lock()
	session, exists := sm.sessions[sessionID]
	if !exists {
		sm.mu.Unlock()
		return nil, false
	}
	if sm.sessionExpiredLocked(session, now) {
		sm.removeSessionLocked(sessionID)
		handlers := sm.handlersLocked()
		sm.mu.Unlock()
		sm.callRemovedHandlers(handlers, []string{sessionID})
		return nil, false
	}
	session.touch(now)
	sm.mu.Unlock()
	return session, true
}

// RemoveSession removes a session and cleans up resources
func (sm *SessionManager) RemoveSession(sessionID string) {
	if sessionID == "" {
		return
	}
	sm.mu.Lock()
	if _, exists := sm.sessions[sessionID]; !exists {
		sm.mu.Unlock()
		return
	}
	sm.removeSessionLocked(sessionID)
	handlers := sm.handlersLocked()
	sm.mu.Unlock()
	sm.callRemovedHandlers(handlers, []string{sessionID})
	logger.Infof("[dubbo-go-pixiu] mcp server removed session: %s", sessionID)
}

// Stop stops the session manager
func (sm *SessionManager) Stop() {
	sm.once.Do(func() {
		close(sm.stopCh)

		sm.mu.Lock()
		removed := make([]string, 0, len(sm.sessions))
		for sessionID, session := range sm.sessions {
			session.close()
			delete(sm.sessions, sessionID)
			removed = append(removed, sessionID)
		}
		handlers := sm.handlersLocked()
		sm.mu.Unlock()
		sm.callRemovedHandlers(handlers, removed)
	})
}

// generateSessionID generates a unique session ID
func (sm *SessionManager) generateSessionID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID
		return hex.EncodeToString([]byte(time.Now().String()))
	}
	return hex.EncodeToString(bytes)
}

// startCleanupRoutine starts the session cleanup routine
func (sm *SessionManager) startCleanupRoutine() {
	ticker := time.NewTicker(CleanupInterval)
	defer ticker.Stop()

	for {
		select {
		case <-ticker.C:
			sm.cleanupExpiredSessions()
		case <-sm.stopCh:
			return
		}
	}
}

// cleanupExpiredSessions removes expired sessions
func (sm *SessionManager) cleanupExpiredSessions() {
	sm.mu.Lock()
	now := sm.now()
	var toRemove []string

	for sessionID, session := range sm.sessions {
		if sm.sessionExpiredLocked(session, now) {
			toRemove = append(toRemove, sessionID)
		}
	}

	for _, sessionID := range toRemove {
		sm.removeSessionLocked(sessionID)
		logger.Infof("[dubbo-go-pixiu] mcp server cleaned up expired session: %s", sessionID)
	}
	handlers := sm.handlersLocked()
	sm.mu.Unlock()
	sm.callRemovedHandlers(handlers, toRemove)

	if len(toRemove) > 0 {
		logger.Debugf("[dubbo-go-pixiu] mcp server cleaned up %d expired sessions", len(toRemove))
	}
}

// AllSessionIDs returns all active session IDs
func (sm *SessionManager) AllSessionIDs() []string {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	ids := make([]string, 0, len(sm.sessions))
	for id := range sm.sessions {
		ids = append(ids, id)
	}
	return ids
}

// ActiveSessionCount returns the number of active sessions
func (sm *SessionManager) ActiveSessionCount() int {
	sm.mu.RLock()
	defer sm.mu.RUnlock()
	return len(sm.sessions)
}

func (sm *SessionManager) sessionExpiredLocked(session *MCPSession, now time.Time) bool {
	return now.Sub(session.lastActivity()) > SessionTimeout
}

func (sm *SessionManager) removeSessionLocked(sessionID string) {
	session, exists := sm.sessions[sessionID]
	if !exists {
		return
	}
	session.close()
	delete(sm.sessions, sessionID)
}

func (sm *SessionManager) handlersLocked() []SessionRemovedHandler {
	if len(sm.removedHandlers) == 0 {
		return nil
	}
	handlers := make([]SessionRemovedHandler, len(sm.removedHandlers))
	copy(handlers, sm.removedHandlers)
	return handlers
}

func (sm *SessionManager) callRemovedHandlers(handlers []SessionRemovedHandler, sessionIDs []string) {
	for _, id := range sessionIDs {
		for _, handler := range handlers {
			func() {
				defer func() {
					if r := recover(); r != nil {
						logger.Warnf("[dubbo-go-pixiu] mcp session removal handler panicked for session %s: %v", id, r)
					}
				}()
				handler(id)
			}()
		}
	}
}

func (s *MCPSession) close() {
	s.closeOnce.Do(func() {
		close(s.Done)
		s.writeMu.Lock()
		defer s.writeMu.Unlock()
		if s.PipeWriter != nil {
			_ = s.PipeWriter.Close()
			s.PipeWriter = nil
		}
	})
}

func (s *MCPSession) touch(now time.Time) {
	s.mu.Lock()
	s.LastActivity = now
	s.mu.Unlock()
}

func (s *MCPSession) lastActivity() time.Time {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.LastActivity
}

// SetToolsListChangedSupported stores whether this client can receive
// notifications/tools/list_changed.
func (s *MCPSession) SetToolsListChangedSupported(supported bool) {
	s.mu.Lock()
	s.SupportsToolsListChanged = supported
	if !supported {
		s.PendingToolsListChanged = false
	}
	s.mu.Unlock()
}

// ToolsListChangedSupported reports whether notifications may be sent.
func (s *MCPSession) ToolsListChangedSupported() bool {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return s.SupportsToolsListChanged
}

// MarkToolsListChangedPending coalesces a pending list_changed notification.
func (s *MCPSession) MarkToolsListChangedPending() {
	s.mu.Lock()
	if s.SupportsToolsListChanged {
		s.PendingToolsListChanged = true
	}
	s.mu.Unlock()
}

// ConsumeToolsListChangedPending consumes one coalesced pending notification.
func (s *MCPSession) ConsumeToolsListChangedPending() bool {
	s.mu.Lock()
	defer s.mu.Unlock()
	if !s.SupportsToolsListChanged || !s.PendingToolsListChanged {
		return false
	}
	s.PendingToolsListChanged = false
	return true
}

// SetPipeWriter installs the current SSE pipe writer.
func (s *MCPSession) SetPipeWriter(writer *io.PipeWriter) {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	s.PipeWriter = writer
}

// HasPipeWriter reports whether an SSE stream is online.
func (s *MCPSession) HasPipeWriter() bool {
	s.writeMu.Lock()
	defer s.writeMu.Unlock()
	return s.PipeWriter != nil
}

// WriteSSEData writes a formatted SSE frame and updates last activity.
func (s *MCPSession) WriteSSEData(data []byte, now time.Time) error {
	s.writeMu.Lock()
	writer := s.PipeWriter
	s.writeMu.Unlock()
	if writer == nil {
		return fmt.Errorf("SSE pipe not established")
	}
	if _, err := writer.Write(data); err != nil {
		return err
	}
	s.touch(now)
	return nil
}

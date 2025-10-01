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

// MCPSession represents an active MCP session
type MCPSession struct {
	ID           string
	CreatedAt    time.Time
	LastActivity time.Time
	PipeWriter   *io.PipeWriter // Pipe writer for sending SSE messages
	Done         chan struct{}
}

// SessionManager manages MCP sessions for SSE connections
type SessionManager struct {
	sessions map[string]*MCPSession
	mu       sync.RWMutex
	stopCh   chan struct{}
	once     sync.Once
}

// NewSessionManager creates a new session manager
func NewSessionManager() *SessionManager {
	sm := &SessionManager{
		sessions: make(map[string]*MCPSession),
		stopCh:   make(chan struct{}),
	}
	go sm.startCleanupRoutine()
	return sm
}

// EnsureSession gets existing session or creates new one
func (sm *SessionManager) EnsureSession(sessionIDHeader string) (*MCPSession, bool) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	// Try to get existing session
	if sessionIDHeader != "" {
		if session, exists := sm.sessions[sessionIDHeader]; exists {
			session.LastActivity = time.Now()
			return session, false // existing session
		}
	}

	// Create new session
	sessionID := sm.generateSessionID()
	session := &MCPSession{
		ID:           sessionID,
		CreatedAt:    time.Now(),
		LastActivity: time.Now(),
		Done:         make(chan struct{}),
	}

	sm.sessions[sessionID] = session
	logger.Infof("[dubbo-go-pixiu] mcp server created new session: %s", sessionID)
	return session, true // new session
}

// Session retrieves a session by ID
func (sm *SessionManager) Session(sessionID string) (*MCPSession, bool) {
	sm.mu.RLock()
	defer sm.mu.RUnlock()

	session, exists := sm.sessions[sessionID]
	if exists {
		session.LastActivity = time.Now()
	}
	return session, exists
}

// RemoveSession removes a session and cleans up resources
func (sm *SessionManager) RemoveSession(sessionID string) {
	sm.mu.Lock()
	defer sm.mu.Unlock()

	if session, exists := sm.sessions[sessionID]; exists {
		// Close Done channel to signal goroutines
		close(session.Done)

		// Close PipeWriter to end the SSE stream
		if session.PipeWriter != nil {
			session.PipeWriter.Close()
		}

		delete(sm.sessions, sessionID)
		logger.Infof("[dubbo-go-pixiu] mcp server removed session: %s", sessionID)
	}
}

// Stop stops the session manager
func (sm *SessionManager) Stop() {
	sm.once.Do(func() {
		close(sm.stopCh)

		sm.mu.Lock()
		for sessionID, session := range sm.sessions {
			close(session.Done)
			delete(sm.sessions, sessionID)
		}
		sm.mu.Unlock()
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
	defer sm.mu.Unlock()

	now := time.Now()
	var toRemove []string

	for sessionID, session := range sm.sessions {
		if now.Sub(session.LastActivity) > SessionTimeout {
			toRemove = append(toRemove, sessionID)
		}
	}

	for _, sessionID := range toRemove {
		if session, exists := sm.sessions[sessionID]; exists {
			close(session.Done)
			delete(sm.sessions, sessionID)
			logger.Infof("[dubbo-go-pixiu] mcp server cleaned up expired session: %s", sessionID)
		}
	}

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

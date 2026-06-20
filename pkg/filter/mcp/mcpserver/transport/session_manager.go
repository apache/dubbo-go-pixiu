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

// StreamToken identifies one attached SSE stream generation.
type StreamToken uint64

type streamAttachment struct {
	token   StreamToken
	writer  *io.PipeWriter
	writeMu sync.Mutex
}

// MCPSession represents an active MCP session
type MCPSession struct {
	ID        string
	CreatedAt time.Time
	Done      chan struct{}

	mu                          sync.RWMutex
	LastActivity                time.Time
	toolListChangeVersion       uint64
	lastNotifiedToolListVersion uint64

	streamMu  sync.Mutex
	streamSeq StreamToken
	stream    *streamAttachment
	closeOnce sync.Once
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

// CreateSession creates a new MCP session. Client supplied session IDs are never
// accepted here; this prevents session fixation during initialize.
func (sm *SessionManager) CreateSession() (*MCPSession, error) {
	sm.mu.Lock()
	defer sm.mu.Unlock()
	now := sm.now()
	sessionID := sm.generateUniqueSessionIDLocked()
	session := &MCPSession{
		ID:           sessionID,
		CreatedAt:    now,
		LastActivity: now,
		Done:         make(chan struct{}),
	}
	sm.sessions[sessionID] = session
	logger.Infof("[dubbo-go-pixiu] mcp server created MCP session")
	return session, nil
}

// GetSession retrieves an existing MCP session. Unknown or expired IDs are not
// replaced with new sessions.
func (sm *SessionManager) GetSession(sessionID string) (*MCPSession, bool) {
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

// Session retrieves a session by ID.
func (sm *SessionManager) Session(sessionID string) (*MCPSession, bool) {
	return sm.GetSession(sessionID)
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
	logger.Infof("[dubbo-go-pixiu] mcp server removed MCP session")
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

// generateSessionID generates a session ID
func (sm *SessionManager) generateSessionID() string {
	bytes := make([]byte, 16)
	if _, err := rand.Read(bytes); err != nil {
		// Fallback to timestamp-based ID
		return hex.EncodeToString([]byte(time.Now().String()))
	}
	return hex.EncodeToString(bytes)
}

func (sm *SessionManager) generateUniqueSessionIDLocked() string {
	for {
		id := sm.generateSessionID()
		if _, exists := sm.sessions[id]; !exists {
			return id
		}
	}
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
		logger.Infof("[dubbo-go-pixiu] mcp server cleaned up expired MCP session")
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
						logger.Warnf("[dubbo-go-pixiu] mcp session removal handler panicked: %v", r)
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
		old := s.clearStream()
		closeStream(old)
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

// MarkToolsListChangedPending records that this session's visible tool set has
// changed. The version is monotonic and bounded to one counter per session.
func (s *MCPSession) MarkToolsListChangedPending() uint64 {
	s.mu.Lock()
	s.toolListChangeVersion++
	version := s.toolListChangeVersion
	s.mu.Unlock()
	return version
}

// PendingToolsListChangedVersion returns the latest unsent tool-list change.
func (s *MCPSession) PendingToolsListChangedVersion() (uint64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if s.toolListChangeVersion <= s.lastNotifiedToolListVersion {
		return 0, false
	}
	return s.toolListChangeVersion, true
}

// MarkToolsListChangedNotified marks one version as successfully sent. If a
// newer version was created while the send was in flight, it remains pending.
func (s *MCPSession) MarkToolsListChangedNotified(version uint64) {
	s.mu.Lock()
	if version > s.lastNotifiedToolListVersion {
		s.lastNotifiedToolListVersion = version
	}
	s.mu.Unlock()
}

// AttachStream installs a new SSE stream and closes the previous stream, if any.
func (s *MCPSession) AttachStream(writer *io.PipeWriter) StreamToken {
	s.streamMu.Lock()
	s.streamSeq++
	token := s.streamSeq
	old := s.stream
	s.stream = &streamAttachment{token: token, writer: writer}
	s.streamMu.Unlock()

	closeStream(old)
	return token
}

// DetachStream removes the current SSE stream only if the token still owns it.
func (s *MCPSession) DetachStream(token StreamToken) {
	s.streamMu.Lock()
	if s.stream == nil || s.stream.token != token {
		s.streamMu.Unlock()
		return
	}
	old := s.stream
	s.stream = nil
	s.streamMu.Unlock()

	closeStream(old)
}

func (s *MCPSession) clearStream() *streamAttachment {
	s.streamMu.Lock()
	old := s.stream
	s.stream = nil
	s.streamMu.Unlock()
	return old
}

// HasPipeWriter reports whether an SSE stream is online.
func (s *MCPSession) HasPipeWriter() bool {
	s.streamMu.Lock()
	defer s.streamMu.Unlock()
	return s.stream != nil
}

// WriteSSEData writes a formatted SSE frame and updates last activity.
func (s *MCPSession) WriteSSEData(data []byte, now time.Time) error {
	s.streamMu.Lock()
	stream := s.stream
	s.streamMu.Unlock()
	return s.writeSSEData(stream, data, now)
}

// WriteSSEDataForStream writes only when token still owns the active stream.
func (s *MCPSession) WriteSSEDataForStream(token StreamToken, data []byte, now time.Time) error {
	s.streamMu.Lock()
	stream := s.stream
	if stream == nil || stream.token != token {
		s.streamMu.Unlock()
		return fmt.Errorf("SSE stream no longer attached")
	}
	s.streamMu.Unlock()
	return s.writeSSEData(stream, data, now)
}

func (s *MCPSession) writeSSEData(stream *streamAttachment, data []byte, now time.Time) error {
	if stream == nil {
		return fmt.Errorf("SSE pipe not established")
	}

	stream.writeMu.Lock()
	_, err := stream.writer.Write(data)
	stream.writeMu.Unlock()
	if err != nil {
		return err
	}
	s.touch(now)
	return nil
}

func closeStream(stream *streamAttachment) {
	if stream == nil || stream.writer == nil {
		return
	}
	_ = stream.writer.Close()
}

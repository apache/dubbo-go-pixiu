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

package router

import (
	"sync"
	"time"
)

// sessionEntry bundles a session's current plan with its progressive-disclosure
// bookkeeping (successful call count and last-touch time for cleanup), plus the
// agent identifier for audit logs.
type sessionEntry struct {
	plan      *SelectionPlan
	callCount int
	agentID   string
	updatedAt time.Time
}

// SessionPlanStore is an in-process, concurrency-safe store of per-session
// selection plans keyed by Mcp-Session-Id. Plan lifetime is independent of the
// transport SessionManager, but a janitor evicts entries idle beyond TTL so the
// map cannot grow without bound when sessions vanish without cleanup.
type SessionPlanStore struct {
	mu      sync.RWMutex
	entries map[string]*sessionEntry
	ttl     time.Duration

	stopCh chan struct{}
	once   sync.Once
}

// DefaultPlanTTL is twice the transport session timeout (30m), so a plan
// outlives its session's normal lifecycle before being reclaimed.
const DefaultPlanTTL = 60 * time.Minute

// planCleanupInterval is how often the janitor scans for stale entries.
const planCleanupInterval = 5 * time.Minute

// NewSessionPlanStore creates a store with the default TTL and starts its
// background cleanup goroutine.
func NewSessionPlanStore() *SessionPlanStore {
	return NewSessionPlanStoreWithTTL(DefaultPlanTTL)
}

// NewSessionPlanStoreWithTTL creates a store with a custom TTL (used by tests).
func NewSessionPlanStoreWithTTL(ttl time.Duration) *SessionPlanStore {
	s := &SessionPlanStore{
		entries: make(map[string]*sessionEntry),
		ttl:     ttl,
		stopCh:  make(chan struct{}),
	}
	go s.cleanupLoop()
	return s
}

// Get returns the plan for a session, or (nil, false) if absent. An entry that
// exists only for agentID/callCount bookkeeping (no plan yet) is treated as
// absent so callers can rely on a non-nil plan when ok is true.
func (s *SessionPlanStore) Get(sessionID string) (*SelectionPlan, bool) {
	if sessionID == "" {
		return nil, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[sessionID]
	if !ok || e.plan == nil {
		return nil, false
	}
	return e.plan, true
}

// Set stores (or replaces) the plan for its session, preserving any existing
// call count so progressive disclosure survives plan recomputation.
func (s *SessionPlanStore) Set(plan *SelectionPlan) {
	if plan == nil || plan.SessionID == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[plan.SessionID]
	if !ok {
		e = &sessionEntry{}
		s.entries[plan.SessionID] = e
	}
	e.plan = plan
	e.updatedAt = time.Now()
}

// Delete removes a session's plan and bookkeeping.
func (s *SessionPlanStore) Delete(sessionID string) {
	s.mu.Lock()
	defer s.mu.Unlock()
	delete(s.entries, sessionID)
}

// IncrementCallCount records a successful tool call for the session and returns
// the new count. It is a no-op (returns 0) for unknown sessions.
func (s *SessionPlanStore) IncrementCallCount(sessionID string) int {
	if sessionID == "" {
		return 0
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[sessionID]
	if !ok {
		return 0
	}
	e.callCount++
	e.updatedAt = time.Now()
	return e.callCount
}

// CallCount returns the recorded successful-call count for a session.
func (s *SessionPlanStore) CallCount(sessionID string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if e, ok := s.entries[sessionID]; ok {
		return e.callCount
	}
	return 0
}

// Len returns the number of sessions with an active plan (for metrics/tests).
// Bookkeeping-only entries created by initialize are intentionally excluded.
func (s *SessionPlanStore) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	n := 0
	for _, e := range s.entries {
		if e.plan != nil {
			n++
		}
	}
	return n
}

// SetAgentID stores the agent identifier for a session. It is called during
// initialize to preserve the clientInfo.name for later decision logs.
func (s *SessionPlanStore) SetAgentID(sessionID, agentID string) {
	if sessionID == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	e, ok := s.entries[sessionID]
	if !ok {
		e = &sessionEntry{}
		s.entries[sessionID] = e
	}
	e.agentID = agentID
	e.updatedAt = time.Now()
}

// AgentID returns the stored agent identifier for a session, or "" if unknown.
func (s *SessionPlanStore) AgentID(sessionID string) string {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if e, ok := s.entries[sessionID]; ok {
		return e.agentID
	}
	return ""
}

// Stop terminates the cleanup goroutine. Safe to call multiple times.
func (s *SessionPlanStore) Stop() {
	s.once.Do(func() { close(s.stopCh) })
}

func (s *SessionPlanStore) cleanupLoop() {
	ticker := time.NewTicker(planCleanupInterval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.evictStale()
		case <-s.stopCh:
			return
		}
	}
}

// evictStale removes entries untouched for longer than the TTL.
func (s *SessionPlanStore) evictStale() {
	cutoff := time.Now().Add(-s.ttl)
	s.mu.Lock()
	defer s.mu.Unlock()
	for id, e := range s.entries {
		if e.updatedAt.Before(cutoff) {
			delete(s.entries, id)
		}
	}
}

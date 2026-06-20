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

type nowFunc func() time.Time

// PlanKey isolates cached plan state by router instance and MCP session.
type PlanKey struct {
	RouterID  string
	SessionID string
}

// NewPlanKey constructs a key for one router instance/session pair.
func NewPlanKey(routerID, sessionID string) PlanKey {
	return PlanKey{RouterID: routerID, SessionID: sessionID}
}

func (k PlanKey) valid() bool {
	return k.SessionID != ""
}

// sessionEntry bundles the immutable enforcement plan with the progressive
// state that is safe to retain for the session lifetime.
type sessionEntry struct {
	plan            *SelectionPlan
	callCount       int64
	expanded        bool
	identityHash    string
	progressiveHash string
	updatedAt       time.Time
}

// SessionPlanStoreOptions configures a SessionPlanStore.
type SessionPlanStoreOptions struct {
	TTL             time.Duration
	MaxEntries      int
	Now             nowFunc
	CleanupInterval time.Duration
}

// SessionPlanStore is an in-process, concurrency-safe store of selection plans
// keyed by router instance plus Mcp-Session-Id. Transport session removal
// deletes entries immediately through the mcpserver-registered hook; TTL is a
// safety net for stale entries that survive abnormal shutdown paths.
type SessionPlanStore struct {
	mu      sync.RWMutex
	entries map[PlanKey]*sessionEntry
	ttl     time.Duration
	max     int
	now     nowFunc

	stopCh chan struct{}
	once   sync.Once
}

// DefaultPlanTTL is the idle TTL safety net for cached plans.
const DefaultPlanTTL = 30 * time.Minute

// DefaultPlanMaxEntries caps the number of active session plans.
const DefaultPlanMaxEntries = 10000

// planCleanupInterval is how often the janitor scans for stale entries.
const planCleanupInterval = 5 * time.Minute

// NewSessionPlanStore creates a store with production defaults and starts its
// background cleanup goroutine.
func NewSessionPlanStore() *SessionPlanStore {
	return NewSessionPlanStoreWithOptions(SessionPlanStoreOptions{})
}

// NewSessionPlanStoreWithMaxEntries creates a store using the default TTL and
// the supplied capacity. A non-positive value means the production default.
func NewSessionPlanStoreWithMaxEntries(maxEntries int) *SessionPlanStore {
	return NewSessionPlanStoreWithOptions(SessionPlanStoreOptions{MaxEntries: maxEntries})
}

// NewSessionPlanStoreWithTTL creates a store with a custom TTL (used by tests).
func NewSessionPlanStoreWithTTL(ttl time.Duration) *SessionPlanStore {
	return NewSessionPlanStoreWithOptions(SessionPlanStoreOptions{TTL: ttl})
}

// NewSessionPlanStoreWithOptions creates a store with explicit options.
func NewSessionPlanStoreWithOptions(opts SessionPlanStoreOptions) *SessionPlanStore {
	ttl := opts.TTL
	if ttl <= 0 {
		ttl = DefaultPlanTTL
	}
	maxEntries := opts.MaxEntries
	if maxEntries <= 0 {
		maxEntries = DefaultPlanMaxEntries
	}
	now := opts.Now
	if now == nil {
		now = time.Now
	}
	cleanupInterval := opts.CleanupInterval
	if cleanupInterval <= 0 {
		cleanupInterval = planCleanupInterval
	}

	s := &SessionPlanStore{
		entries: make(map[PlanKey]*sessionEntry),
		ttl:     ttl,
		max:     maxEntries,
		now:     now,
		stopCh:  make(chan struct{}),
	}
	go s.cleanupLoop(cleanupInterval)
	return s
}

// Get returns an immutable copy of the plan for a session, or (nil, false) if
// absent. Callers cannot mutate store-owned slices or maps through the result.
func (s *SessionPlanStore) Get(key PlanKey) (*SelectionPlan, bool) {
	if !key.valid() {
		return nil, false
	}
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[key]
	if !ok || e.plan == nil {
		return nil, false
	}
	return clonePlan(e.plan, true), true
}

// Set stores or replaces the plan for its session. Successful-call state is
// preserved only while the verified identity and progressive config hashes are
// unchanged; identity/config changes reset progressive disclosure.
func (s *SessionPlanStore) Set(key PlanKey, plan *SelectionPlan) {
	if !key.valid() || plan == nil {
		return
	}
	now := s.now()
	stored := clonePlan(plan, false)
	stored.Reasons = nil
	stored.SessionID = key.SessionID

	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.entries[key]
	if !ok {
		s.evictForInsertLocked(now)
		e = &sessionEntry{}
		s.entries[key] = e
	} else if e.identityHash != stored.IdentityHash || e.progressiveHash != stored.ProgressiveHash {
		e.callCount = 0
		e.expanded = false
	}

	stored.Expanded = e.expanded || stored.Expanded
	e.plan = stored
	e.identityHash = stored.IdentityHash
	e.progressiveHash = stored.ProgressiveHash
	e.updatedAt = now
	s.publishActiveLocked()
}

// Delete removes one router instance's plan and progressive state.
func (s *SessionPlanStore) Delete(key PlanKey) {
	s.delete(key, "explicit")
}

func (s *SessionPlanStore) delete(key PlanKey, reason string) {
	if !key.valid() {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	if _, ok := s.entries[key]; !ok {
		return
	}
	delete(s.entries, key)
	recordPlanEvicted(reason)
	s.publishActiveLocked()
}

// DeleteSession removes all router-instance state for a transport session. It
// is the method registered with SessionManager teardown callbacks.
func (s *SessionPlanStore) DeleteSession(sessionID string) {
	if sessionID == "" {
		return
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	removed := 0
	for key := range s.entries {
		if key.SessionID == sessionID {
			delete(s.entries, key)
			removed++
		}
	}
	for i := 0; i < removed; i++ {
		recordPlanEvicted("session_end")
	}
	if removed > 0 {
		s.publishActiveLocked()
	}
}

// RecordCallSuccess records a successful tool call and returns whether this
// call crossed the progressive threshold. The threshold update and transition
// check happen under one lock, so concurrent calls can observe at most one
// transition.
func (s *SessionPlanStore) RecordCallSuccess(key PlanKey, requested string, expandAfter int) CallSuccessResult {
	if !key.valid() || requested == "" || expandAfter <= 0 {
		return CallSuccessResult{}
	}
	now := s.now()

	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.entries[key]
	if !ok || e.plan == nil || !e.plan.Contains(requested) {
		return CallSuccessResult{}
	}
	e.callCount++
	e.updatedAt = now

	result := CallSuccessResult{Count: e.callCount}
	if !e.expanded && e.callCount >= int64(expandAfter) {
		e.expanded = true
		result.Transitioned = true
		e.plan.Expanded = true
	}
	return result
}

// CallState returns progressive state for the current identity/config pair.
// Mismatched identity or progressive config returns the reset state without
// mutating the store; Set performs the reset atomically when the new plan is
// published.
func (s *SessionPlanStore) CallState(key PlanKey, identityHash, progressiveHash string) (int64, bool) {
	s.mu.RLock()
	defer s.mu.RUnlock()
	e, ok := s.entries[key]
	if !ok || e.identityHash != identityHash || e.progressiveHash != progressiveHash {
		return 0, false
	}
	return e.callCount, e.expanded
}

// CallCount returns the recorded successful-call count for a session.
func (s *SessionPlanStore) CallCount(key PlanKey) int64 {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if e, ok := s.entries[key]; ok {
		return e.callCount
	}
	return 0
}

// Len returns the number of active session plans.
func (s *SessionPlanStore) Len() int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	return len(s.entries)
}

// SetMaxEntries updates the store capacity. A non-positive value restores the
// production default. If the current size exceeds the new cap, oldest entries
// are evicted immediately.
func (s *SessionPlanStore) SetMaxEntries(maxEntries int) {
	if maxEntries <= 0 {
		maxEntries = DefaultPlanMaxEntries
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.max = maxEntries
	now := s.now()
	s.evictExpiredLocked(now)
	for len(s.entries) > s.max {
		s.evictOldestLocked("capacity")
	}
	s.publishActiveLocked()
}

// Stop terminates the cleanup goroutine and publishes an active-plan gauge of
// zero. Safe to call multiple times.
func (s *SessionPlanStore) Stop() {
	s.once.Do(func() {
		close(s.stopCh)
		s.mu.Lock()
		s.entries = make(map[PlanKey]*sessionEntry)
		s.publishActiveLocked()
		s.mu.Unlock()
	})
}

func (s *SessionPlanStore) cleanupLoop(interval time.Duration) {
	ticker := time.NewTicker(interval)
	defer ticker.Stop()
	for {
		select {
		case <-ticker.C:
			s.EvictExpired()
		case <-s.stopCh:
			return
		}
	}
}

// EvictExpired removes entries untouched for longer than the TTL.
func (s *SessionPlanStore) EvictExpired() {
	now := s.now()
	s.mu.Lock()
	defer s.mu.Unlock()
	s.evictExpiredLocked(now)
	s.publishActiveLocked()
}

func (s *SessionPlanStore) evictForInsertLocked(now time.Time) {
	s.evictExpiredLocked(now)
	for len(s.entries) >= s.max {
		s.evictOldestLocked("capacity")
	}
}

func (s *SessionPlanStore) evictExpiredLocked(now time.Time) {
	cutoff := now.Add(-s.ttl)
	for key, e := range s.entries {
		if e.updatedAt.Before(cutoff) {
			delete(s.entries, key)
			recordPlanEvicted("ttl")
		}
	}
}

func (s *SessionPlanStore) evictOldestLocked(reason string) {
	var oldestKey PlanKey
	var oldest time.Time
	for key, e := range s.entries {
		if !oldestKey.valid() || e.updatedAt.Before(oldest) || (e.updatedAt.Equal(oldest) && planKeyLess(key, oldestKey)) {
			oldestKey = key
			oldest = e.updatedAt
		}
	}
	if !oldestKey.valid() {
		return
	}
	delete(s.entries, oldestKey)
	recordPlanEvicted(reason)
}

func planKeyLess(a, b PlanKey) bool {
	if a.RouterID != b.RouterID {
		return a.RouterID < b.RouterID
	}
	return a.SessionID < b.SessionID
}

func (s *SessionPlanStore) publishActiveLocked() {
	setPlansActive(len(s.entries))
}

func clonePlan(plan *SelectionPlan, includeReasons bool) *SelectionPlan {
	if plan == nil {
		return nil
	}
	cp := *plan
	cp.ToolNames = copyStrings(plan.ToolNames)
	cp.VisibleToolNames = copyStrings(plan.VisibleToolNames)
	cp.StageCounts = copyStageCounts(plan.StageCounts)
	cp.toolSet = toolNameSet(cp.ToolNames)
	if includeReasons {
		cp.Reasons = copyDecisionTraces(plan.Reasons)
	} else {
		cp.Reasons = nil
	}
	return &cp
}

func copyStrings(values []string) []string {
	if values == nil {
		return nil
	}
	cp := make([]string, len(values))
	copy(cp, values)
	return cp
}

func copyDecisionTraces(values []DecisionTrace) []DecisionTrace {
	if values == nil {
		return nil
	}
	cp := make([]DecisionTrace, len(values))
	copy(cp, values)
	return cp
}

func copyStageCounts(values map[string]StageCount) map[string]StageCount {
	if values == nil {
		return nil
	}
	cp := make(map[string]StageCount, len(values))
	for k, v := range values {
		cp[k] = v
	}
	return cp
}

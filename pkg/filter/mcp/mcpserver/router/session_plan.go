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
	"errors"
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
	context         SelectionContext
	callCount       int64
	expanded        bool
	identityHash    string
	progressiveHash string
	configHash      string
	catalogVersion  string
	generation      uint64
	nextReceiptID   uint64
	countedReceipts map[uint64]struct{}
	updatedAt       time.Time
}

// SessionPlanStoreOptions configures a SessionPlanStore.
type SessionPlanStoreOptions struct {
	MaxEntries int
	Now        nowFunc
}

// SessionPlanStore is an in-process, concurrency-safe store of selection plans
// owned by one MCP filter instance. Transport session removal deletes entries
// immediately; plans do not expire independently from their sessions.
type SessionPlanStore struct {
	mu      sync.RWMutex
	entries map[PlanKey]*sessionEntry
	max     int
	now     nowFunc
}

// DefaultPlanMaxEntries caps the number of active session plans.
const DefaultPlanMaxEntries = 10000

// ErrPlanStoreFull is returned when a filter has reached its active session
// capacity. The caller must reject new sessions instead of evicting live plans.
var ErrPlanStoreFull = errors.New("mcp router session capacity reached")

// NewSessionPlanStore creates a store with production defaults.
func NewSessionPlanStore() *SessionPlanStore {
	return NewSessionPlanStoreWithOptions(SessionPlanStoreOptions{})
}

// NewSessionPlanStoreWithMaxEntries creates a store using the default TTL and
// the supplied capacity. A non-positive value means the production default.
func NewSessionPlanStoreWithMaxEntries(maxEntries int) *SessionPlanStore {
	return NewSessionPlanStoreWithOptions(SessionPlanStoreOptions{MaxEntries: maxEntries})
}

// NewSessionPlanStoreWithOptions creates a store with explicit options.
func NewSessionPlanStoreWithOptions(opts SessionPlanStoreOptions) *SessionPlanStore {
	maxEntries := opts.MaxEntries
	if maxEntries <= 0 {
		maxEntries = DefaultPlanMaxEntries
	}
	now := opts.Now
	if now == nil {
		now = time.Now
	}

	s := &SessionPlanStore{
		entries: make(map[PlanKey]*sessionEntry),
		max:     maxEntries,
		now:     now,
	}
	registerPlanStoreMetric(s)
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
func (s *SessionPlanStore) Set(key PlanKey, plan *SelectionPlan, sc SelectionContext) error {
	if !key.valid() || plan == nil {
		return nil
	}
	now := s.now()
	stored := clonePlan(plan, false)
	stored.Reasons = nil
	stored.SessionID = key.SessionID

	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.entries[key]
	if !ok {
		if len(s.entries) >= s.max {
			return ErrPlanStoreFull
		}
		e = &sessionEntry{}
		s.entries[key] = e
	} else if e.identityHash != stored.IdentityHash || e.progressiveHash != stored.ProgressiveHash {
		e.callCount = 0
		e.expanded = false
		e.countedReceipts = nil
	}

	e.generation++
	stored.Expanded = e.expanded || stored.Expanded
	stored.Generation = e.generation
	e.plan = stored
	e.context = cloneSelectionContext(sc)
	e.identityHash = stored.IdentityHash
	e.progressiveHash = stored.ProgressiveHash
	e.configHash = stored.ConfigHash
	e.catalogVersion = stored.CatalogVersion
	e.updatedAt = now
	s.publishActiveLocked()
	return nil
}

// SessionPlanContexts returns cloned plan/context pairs for management-plane
// recomputation without touching transport session activity.
func (s *SessionPlanStore) SessionPlanContexts() []SessionPlanContext {
	s.mu.RLock()
	defer s.mu.RUnlock()
	out := make([]SessionPlanContext, 0, len(s.entries))
	for key, entry := range s.entries {
		if entry.plan == nil {
			continue
		}
		out = append(out, SessionPlanContext{
			Key:     key,
			Plan:    clonePlan(entry.plan, false),
			Context: cloneSelectionContext(entry.context),
		})
	}
	return out
}

// SessionPlanContext is a cloned plan and the normalized selection attributes
// used to build it.
type SessionPlanContext struct {
	Key     PlanKey
	Plan    *SelectionPlan
	Context SelectionContext
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

// IssueReceipt returns a non-replayable authorization receipt for the current
// stored plan generation.
func (s *SessionPlanStore) IssueReceipt(key PlanKey, requested string, plan *SelectionPlan, routerID string) (AuthorizationReceipt, error) {
	if !key.valid() || requested == "" || plan == nil {
		return AuthorizationReceipt{}, ErrToolNotAuthorized
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.entries[key]
	if !ok || e.plan == nil ||
		e.plan.Version != plan.Version ||
		e.identityHash != plan.IdentityHash ||
		e.configHash != plan.ConfigHash ||
		e.catalogVersion != plan.CatalogVersion ||
		e.progressiveHash != plan.ProgressiveHash ||
		!e.plan.Contains(requested) {
		return AuthorizationReceipt{}, ErrToolNotAuthorized
	}
	e.nextReceiptID++
	return AuthorizationReceipt{
		RouterInstanceID: routerID,
		SessionID:        key.SessionID,
		ToolName:         requested,
		PlanGeneration:   e.plan.Generation,
		IdentityHash:     e.identityHash,
		ConfigHash:       e.configHash,
		CatalogVersion:   e.catalogVersion,
		ProgressiveHash:  e.progressiveHash,
		ReceiptID:        e.nextReceiptID,
	}, nil
}

// IssueReceiptForVersion signs a tools/call authorization without cloning the
// stored plan. It returns stale=true when the session has a plan, but its
// authorization inputs no longer match the caller's live inputs.
func (s *SessionPlanStore) IssueReceiptForVersion(key PlanKey, requested, expectedVersion, identityHash, configHash, catalogVersion, progressiveHash, routerID string) (AuthorizationReceipt, bool, error) {
	if !key.valid() || requested == "" || expectedVersion == "" {
		return AuthorizationReceipt{}, false, ErrToolNotAuthorized
	}
	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.entries[key]
	if !ok || e.plan == nil {
		return AuthorizationReceipt{}, false, ErrToolNotAuthorized
	}
	if e.plan.Version != expectedVersion ||
		e.identityHash != identityHash ||
		e.configHash != configHash ||
		e.catalogVersion != catalogVersion ||
		e.progressiveHash != progressiveHash {
		return AuthorizationReceipt{}, true, nil
	}
	if !e.plan.Contains(requested) {
		return AuthorizationReceipt{}, false, ErrToolNotAuthorized
	}
	e.nextReceiptID++
	return AuthorizationReceipt{
		RouterInstanceID: routerID,
		SessionID:        key.SessionID,
		ToolName:         requested,
		PlanGeneration:   e.plan.Generation,
		IdentityHash:     e.identityHash,
		ConfigHash:       e.configHash,
		CatalogVersion:   e.catalogVersion,
		ProgressiveHash:  e.progressiveHash,
		ReceiptID:        e.nextReceiptID,
	}, false, nil
}

// RecordCallSuccess records a successful tool call and returns whether this
// call crossed the progressive threshold. The threshold update and transition
// check happen under one lock, so concurrent calls can observe at most one
// transition.
func (s *SessionPlanStore) RecordCallSuccess(receipt AuthorizationReceipt, expandAfter int) CallSuccessResult {
	key := NewPlanKey(receipt.RouterInstanceID, receipt.SessionID)
	if !key.valid() || receipt.ToolName == "" || expandAfter <= 0 || receipt.ReceiptID == 0 {
		return CallSuccessResult{}
	}
	now := s.now()

	s.mu.Lock()
	defer s.mu.Unlock()

	e, ok := s.entries[key]
	if !ok || e.plan == nil ||
		e.generation != receipt.PlanGeneration ||
		e.identityHash != receipt.IdentityHash ||
		e.configHash != receipt.ConfigHash ||
		e.catalogVersion != receipt.CatalogVersion ||
		e.progressiveHash != receipt.ProgressiveHash ||
		!e.plan.Contains(receipt.ToolName) {
		return CallSuccessResult{}
	}
	if e.countedReceipts == nil {
		e.countedReceipts = make(map[uint64]struct{})
	}
	if _, counted := e.countedReceipts[receipt.ReceiptID]; counted {
		return CallSuccessResult{}
	}
	e.countedReceipts[receipt.ReceiptID] = struct{}{}
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
// production default. Existing sessions are never evicted to satisfy a smaller
// cap; new sessions will be rejected until usage drops below the limit.
func (s *SessionPlanStore) SetMaxEntries(maxEntries int) {
	if maxEntries <= 0 {
		maxEntries = DefaultPlanMaxEntries
	}
	s.mu.Lock()
	defer s.mu.Unlock()
	s.max = maxEntries
	s.publishActiveLocked()
}

// Stop clears all plans and publishes an active-plan gauge of zero.
func (s *SessionPlanStore) Stop() {
	s.mu.Lock()
	defer s.mu.Unlock()
	s.entries = make(map[PlanKey]*sessionEntry)
	s.publishActiveLocked()
	unregisterPlanStoreMetric(s)
}

func planKeyLess(a, b PlanKey) bool {
	if a.RouterID != b.RouterID {
		return a.RouterID < b.RouterID
	}
	return a.SessionID < b.SessionID
}

func (s *SessionPlanStore) publishActiveLocked() {
	setPlansActive(s, len(s.entries))
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

func cloneSelectionContext(sc SelectionContext) SelectionContext {
	cp := sc
	if sc.Claims != nil {
		cp.Claims = make(map[string]any, len(sc.Claims))
		for k, v := range sc.Claims {
			cp.Claims[k] = v
		}
	}
	return cp
}

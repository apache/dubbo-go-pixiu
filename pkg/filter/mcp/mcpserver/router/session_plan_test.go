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
	"fmt"
	"sync"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

type fakeClock struct {
	mu  sync.Mutex
	now time.Time
}

func newFakeClock() *fakeClock {
	return &fakeClock{now: time.Unix(100, 0)}
}

func (c *fakeClock) Now() time.Time {
	c.mu.Lock()
	defer c.mu.Unlock()
	return c.now
}

func (c *fakeClock) Advance(d time.Duration) {
	c.mu.Lock()
	c.now = c.now.Add(d)
	c.mu.Unlock()
}

func newTestPlanStore(clock *fakeClock, maxEntries int) *SessionPlanStore {
	return NewSessionPlanStoreWithOptions(SessionPlanStoreOptions{
		TTL:             time.Minute,
		MaxEntries:      maxEntries,
		Now:             clock.Now,
		CleanupInterval: time.Hour,
	})
}

func testPlanKey(sessionID string) PlanKey {
	return NewPlanKey("router-test", sessionID)
}

func TestSessionPlanStore_SetGet(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 10)
	defer s.Stop()

	plan := &SelectionPlan{
		SessionID:        "s1",
		ToolNames:        []string{"a"},
		VisibleToolNames: []string{"a"},
		Reasons:          []DecisionTrace{{Tool: "b"}},
		Version:          "v1",
	}
	s.Set(testPlanKey("s1"), plan)

	got, ok := s.Get(testPlanKey("s1"))
	require.True(t, ok)
	assert.Equal(t, []string{"a"}, got.ToolNames)
	assert.Equal(t, []string{"a"}, got.VisibleToolNames)
	assert.Nil(t, got.Reasons, "stored plans do not retain decision traces")

	_, ok = s.Get(testPlanKey("missing"))
	assert.False(t, ok)
}

func TestSessionPlanStore_SetIgnoresEmpty(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 10)
	defer s.Stop()

	s.Set(testPlanKey("ignored"), nil)
	s.Set(testPlanKey(""), &SelectionPlan{SessionID: ""})
	assert.Equal(t, 0, s.Len())
}

func TestSessionPlanStore_MutationIsolation(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 10)
	defer s.Stop()

	plan := &SelectionPlan{SessionID: "s1", ToolNames: []string{"a", "b"}, VisibleToolNames: []string{"a"}}
	s.Set(testPlanKey("s1"), plan)
	plan.ToolNames[0] = "mutated"
	plan.VisibleToolNames[0] = "mutated"

	got, ok := s.Get(testPlanKey("s1"))
	require.True(t, ok)
	assert.Equal(t, []string{"a", "b"}, got.ToolNames)
	assert.Equal(t, []string{"a"}, got.VisibleToolNames)

	got.ToolNames[0] = "caller-mutated"
	got.VisibleToolNames[0] = "caller-mutated"
	got2, ok := s.Get(testPlanKey("s1"))
	require.True(t, ok)
	assert.Equal(t, []string{"a", "b"}, got2.ToolNames)
	assert.Equal(t, []string{"a"}, got2.VisibleToolNames)
}

func TestSessionPlanStore_CallCountPreservedAndResetByIdentity(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 10)
	defer s.Stop()

	key := testPlanKey("s1")
	s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}, IdentityHash: "tenant-a", ProgressiveHash: "p1"})
	assert.Equal(t, CallSuccessResult{Count: 1}, s.RecordCallSuccess(key, "a", 3))
	assert.Equal(t, CallSuccessResult{Count: 2}, s.RecordCallSuccess(key, "a", 3))

	s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}, Version: "v2", IdentityHash: "tenant-a", ProgressiveHash: "p1"})
	assert.Equal(t, int64(2), s.CallCount(key))

	s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}, Version: "v3", IdentityHash: "tenant-b", ProgressiveHash: "p1"})
	assert.Equal(t, int64(0), s.CallCount(key))
	_, expanded := s.CallState(key, "tenant-b", "p1")
	assert.False(t, expanded)
}

func TestSessionPlanStore_RecordCallSuccessTransitionOnce(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 10)
	defer s.Stop()

	key := testPlanKey("s1")
	s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}})
	assert.Equal(t, CallSuccessResult{Count: 1}, s.RecordCallSuccess(key, "a", 2))
	assert.Equal(t, CallSuccessResult{Count: 2, Transitioned: true}, s.RecordCallSuccess(key, "a", 2))
	assert.Equal(t, CallSuccessResult{Count: 3}, s.RecordCallSuccess(key, "a", 2))

	got, ok := s.Get(key)
	require.True(t, ok)
	assert.True(t, got.Expanded)
}

func TestSessionPlanStore_RecordCallSuccessIgnoresUnknownOrOutsidePlan(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 10)
	defer s.Stop()

	assert.Equal(t, CallSuccessResult{}, s.RecordCallSuccess(testPlanKey("missing"), "a", 1))
	key := testPlanKey("s1")
	s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}})
	assert.Equal(t, CallSuccessResult{}, s.RecordCallSuccess(key, "ghost", 1))
	assert.Equal(t, int64(0), s.CallCount(key))
}

func TestSessionPlanStore_Delete(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 10)
	defer s.Stop()

	key := testPlanKey("s1")
	s.Set(key, &SelectionPlan{SessionID: "s1"})
	s.Delete(key)
	_, ok := s.Get(key)
	assert.False(t, ok)
}

func TestSessionPlanStore_IsolatesSameSessionAcrossRouters(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 10)
	defer s.Stop()

	keyA := NewPlanKey("router-a", "shared-session")
	keyB := NewPlanKey("router-b", "shared-session")
	s.Set(keyA, &SelectionPlan{SessionID: "shared-session", ToolNames: []string{"a"}})
	s.Set(keyB, &SelectionPlan{SessionID: "shared-session", ToolNames: []string{"b"}})

	gotA, ok := s.Get(keyA)
	require.True(t, ok)
	assert.Equal(t, []string{"a"}, gotA.ToolNames)

	gotB, ok := s.Get(keyB)
	require.True(t, ok)
	assert.Equal(t, []string{"b"}, gotB.ToolNames)

	assert.Equal(t, CallSuccessResult{Count: 1}, s.RecordCallSuccess(keyA, "a", 3))
	assert.Equal(t, int64(1), s.CallCount(keyA))
	assert.Equal(t, int64(0), s.CallCount(keyB))

	s.DeleteSession("shared-session")
	assert.Equal(t, 0, s.Len())
}

func TestSessionPlanStore_DeleteOneRouterKeepsOtherRouter(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 10)
	defer s.Stop()

	keyA := NewPlanKey("router-a", "shared-session")
	keyB := NewPlanKey("router-b", "shared-session")
	s.Set(keyA, &SelectionPlan{SessionID: "shared-session", ToolNames: []string{"a"}})
	s.Set(keyB, &SelectionPlan{SessionID: "shared-session", ToolNames: []string{"b"}})

	s.Delete(keyA)

	_, ok := s.Get(keyA)
	assert.False(t, ok)
	gotB, ok := s.Get(keyB)
	require.True(t, ok)
	assert.Equal(t, []string{"b"}, gotB.ToolNames)
}

func TestSessionPlanStore_EvictExpiredWithFakeClock(t *testing.T) {
	clock := newFakeClock()
	s := NewSessionPlanStoreWithOptions(SessionPlanStoreOptions{
		TTL:             time.Minute,
		MaxEntries:      10,
		Now:             clock.Now,
		CleanupInterval: time.Hour,
	})
	defer s.Stop()

	s.Set(testPlanKey("s1"), &SelectionPlan{SessionID: "s1"})
	clock.Advance(time.Minute)
	s.EvictExpired()
	assert.Equal(t, 1, s.Len(), "entry expires only after TTL, not exactly at boundary")

	clock.Advance(time.Nanosecond)
	s.EvictExpired()
	assert.Equal(t, 0, s.Len())
}

func TestSessionPlanStore_CapacityEvictsOldestDeterministically(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 2)
	defer s.Stop()

	s.Set(testPlanKey("s2"), &SelectionPlan{SessionID: "s2"})
	s.Set(testPlanKey("s1"), &SelectionPlan{SessionID: "s1"})
	s.Set(testPlanKey("s3"), &SelectionPlan{SessionID: "s3"})

	_, ok := s.Get(testPlanKey("s1"))
	assert.False(t, ok, "same timestamp ties evict lexicographically oldest session")
	_, ok = s.Get(testPlanKey("s2"))
	assert.True(t, ok)
	_, ok = s.Get(testPlanKey("s3"))
	assert.True(t, ok)
}

func TestSessionPlanStore_ExpiredEntriesEvictedBeforeCapacity(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 2)
	defer s.Stop()

	s.Set(testPlanKey("expired"), &SelectionPlan{SessionID: "expired"})
	clock.Advance(2 * time.Minute)
	s.Set(testPlanKey("fresh"), &SelectionPlan{SessionID: "fresh"})
	s.Set(testPlanKey("new"), &SelectionPlan{SessionID: "new"})

	_, ok := s.Get(testPlanKey("expired"))
	assert.False(t, ok)
	assert.Equal(t, 2, s.Len())
}

func TestSessionPlanStore_SetMaxEntries(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 5)
	defer s.Stop()

	for i := 0; i < 5; i++ {
		sessionID := fmt.Sprintf("s%d", i)
		s.Set(testPlanKey(sessionID), &SelectionPlan{SessionID: sessionID})
	}
	s.SetMaxEntries(2)
	assert.Equal(t, 2, s.Len())
}

func TestSessionPlanStore_ConcurrentAccess(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 20)
	defer s.Stop()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			id := fmt.Sprintf("s%d", n%10)
			key := testPlanKey(id)
			s.Set(key, &SelectionPlan{SessionID: id, ToolNames: []string{"a"}, Version: "v"})
			s.RecordCallSuccess(key, "a", 1000)
			s.Get(key)
			s.CallCount(key)
			if n%3 == 0 {
				s.Delete(key)
			}
		}(i)
	}
	wg.Wait()

	assert.LessOrEqual(t, s.Len(), 10)
}

func TestSessionPlanStore_StopIdempotentClearsEntries(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 10)
	s.Set(testPlanKey("s1"), &SelectionPlan{SessionID: "s1"})

	s.Stop()
	s.Stop()
	assert.Equal(t, 0, s.Len())
}

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

func TestSessionPlanStore_SetGet(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 10)
	defer s.Stop()

	plan := &SelectionPlan{
		SessionID:        "s1",
		ToolNames:        []string{"a"},
		VisibleToolNames: []string{"a"},
		Reasons:          []DecisionTrace{{Tool: "b", Kept: false}},
		Version:          "v1",
	}
	s.Set(plan)

	got, ok := s.Get("s1")
	require.True(t, ok)
	assert.Equal(t, []string{"a"}, got.ToolNames)
	assert.Equal(t, []string{"a"}, got.VisibleToolNames)
	assert.Nil(t, got.Reasons, "stored plans do not retain decision traces")

	_, ok = s.Get("missing")
	assert.False(t, ok)
}

func TestSessionPlanStore_SetIgnoresEmpty(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 10)
	defer s.Stop()

	s.Set(nil)
	s.Set(&SelectionPlan{SessionID: ""})
	assert.Equal(t, 0, s.Len())
}

func TestSessionPlanStore_MutationIsolation(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 10)
	defer s.Stop()

	plan := &SelectionPlan{SessionID: "s1", ToolNames: []string{"a", "b"}, VisibleToolNames: []string{"a"}}
	s.Set(plan)
	plan.ToolNames[0] = "mutated"
	plan.VisibleToolNames[0] = "mutated"

	got, ok := s.Get("s1")
	require.True(t, ok)
	assert.Equal(t, []string{"a", "b"}, got.ToolNames)
	assert.Equal(t, []string{"a"}, got.VisibleToolNames)

	got.ToolNames[0] = "caller-mutated"
	got.VisibleToolNames[0] = "caller-mutated"
	got2, ok := s.Get("s1")
	require.True(t, ok)
	assert.Equal(t, []string{"a", "b"}, got2.ToolNames)
	assert.Equal(t, []string{"a"}, got2.VisibleToolNames)
}

func TestSessionPlanStore_CallCountPreservedAndResetByIdentity(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 10)
	defer s.Stop()

	s.Set(&SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}, IdentityHash: "tenant-a", ProgressiveHash: "p1"})
	assert.Equal(t, CallSuccessResult{Count: 1}, s.RecordCallSuccess("s1", "a", 3))
	assert.Equal(t, CallSuccessResult{Count: 2}, s.RecordCallSuccess("s1", "a", 3))

	s.Set(&SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}, Version: "v2", IdentityHash: "tenant-a", ProgressiveHash: "p1"})
	assert.Equal(t, int64(2), s.CallCount("s1"))

	s.Set(&SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}, Version: "v3", IdentityHash: "tenant-b", ProgressiveHash: "p1"})
	assert.Equal(t, int64(0), s.CallCount("s1"))
	_, expanded := s.CallState("s1", "tenant-b", "p1")
	assert.False(t, expanded)
}

func TestSessionPlanStore_RecordCallSuccessTransitionOnce(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 10)
	defer s.Stop()

	s.Set(&SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}})
	assert.Equal(t, CallSuccessResult{Count: 1}, s.RecordCallSuccess("s1", "a", 2))
	assert.Equal(t, CallSuccessResult{Count: 2, Transitioned: true}, s.RecordCallSuccess("s1", "a", 2))
	assert.Equal(t, CallSuccessResult{Count: 3}, s.RecordCallSuccess("s1", "a", 2))

	got, ok := s.Get("s1")
	require.True(t, ok)
	assert.True(t, got.Expanded)
}

func TestSessionPlanStore_RecordCallSuccessIgnoresUnknownOrOutsidePlan(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 10)
	defer s.Stop()

	assert.Equal(t, CallSuccessResult{}, s.RecordCallSuccess("missing", "a", 1))
	s.Set(&SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}})
	assert.Equal(t, CallSuccessResult{}, s.RecordCallSuccess("s1", "ghost", 1))
	assert.Equal(t, int64(0), s.CallCount("s1"))
}

func TestSessionPlanStore_Delete(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 10)
	defer s.Stop()

	s.Set(&SelectionPlan{SessionID: "s1"})
	s.Delete("s1")
	_, ok := s.Get("s1")
	assert.False(t, ok)
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

	s.Set(&SelectionPlan{SessionID: "s1"})
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

	s.Set(&SelectionPlan{SessionID: "s2"})
	s.Set(&SelectionPlan{SessionID: "s1"})
	s.Set(&SelectionPlan{SessionID: "s3"})

	_, ok := s.Get("s1")
	assert.False(t, ok, "same timestamp ties evict lexicographically oldest session")
	_, ok = s.Get("s2")
	assert.True(t, ok)
	_, ok = s.Get("s3")
	assert.True(t, ok)
}

func TestSessionPlanStore_ExpiredEntriesEvictedBeforeCapacity(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 2)
	defer s.Stop()

	s.Set(&SelectionPlan{SessionID: "expired"})
	clock.Advance(2 * time.Minute)
	s.Set(&SelectionPlan{SessionID: "fresh"})
	s.Set(&SelectionPlan{SessionID: "new"})

	_, ok := s.Get("expired")
	assert.False(t, ok)
	assert.Equal(t, 2, s.Len())
}

func TestSessionPlanStore_SetMaxEntries(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 5)
	defer s.Stop()

	for i := 0; i < 5; i++ {
		s.Set(&SelectionPlan{SessionID: fmt.Sprintf("s%d", i)})
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
			s.Set(&SelectionPlan{SessionID: id, ToolNames: []string{"a"}, Version: "v"})
			s.RecordCallSuccess(id, "a", 1000)
			s.Get(id)
			s.CallCount(id)
			if n%3 == 0 {
				s.Delete(id)
			}
		}(i)
	}
	wg.Wait()

	assert.LessOrEqual(t, s.Len(), 10)
}

func TestSessionPlanStore_StopIdempotentClearsEntries(t *testing.T) {
	clock := newFakeClock()
	s := newTestPlanStore(clock, 10)
	s.Set(&SelectionPlan{SessionID: "s1"})

	s.Stop()
	s.Stop()
	assert.Equal(t, 0, s.Len())
}

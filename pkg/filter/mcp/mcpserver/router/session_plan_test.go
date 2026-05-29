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
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestSessionPlanStore_SetGet(t *testing.T) {
	s := NewSessionPlanStoreWithTTL(time.Minute)
	defer s.Stop()

	plan := &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}, Version: "v1"}
	s.Set(plan)

	got, ok := s.Get("s1")
	require.True(t, ok)
	assert.Equal(t, plan, got)

	_, ok = s.Get("missing")
	assert.False(t, ok)
}

func TestSessionPlanStore_SetIgnoresEmpty(t *testing.T) {
	s := NewSessionPlanStoreWithTTL(time.Minute)
	defer s.Stop()

	s.Set(nil)
	s.Set(&SelectionPlan{SessionID: ""})
	assert.Equal(t, 0, s.Len())
}

func TestSessionPlanStore_CallCountPreservedAcrossSet(t *testing.T) {
	s := NewSessionPlanStoreWithTTL(time.Minute)
	defer s.Stop()

	s.Set(&SelectionPlan{SessionID: "s1", Version: "v1"})
	assert.Equal(t, 1, s.IncrementCallCount("s1"))
	assert.Equal(t, 2, s.IncrementCallCount("s1"))

	// Recompute plan: call count must survive.
	s.Set(&SelectionPlan{SessionID: "s1", Version: "v2"})
	assert.Equal(t, 2, s.CallCount("s1"))
}

func TestSessionPlanStore_IncrementUnknownSession(t *testing.T) {
	s := NewSessionPlanStoreWithTTL(time.Minute)
	defer s.Stop()
	assert.Equal(t, 0, s.IncrementCallCount("ghost"))
}

func TestSessionPlanStore_Delete(t *testing.T) {
	s := NewSessionPlanStoreWithTTL(time.Minute)
	defer s.Stop()

	s.Set(&SelectionPlan{SessionID: "s1"})
	s.Delete("s1")
	_, ok := s.Get("s1")
	assert.False(t, ok)
}

func TestSessionPlanStore_EvictStale(t *testing.T) {
	s := NewSessionPlanStoreWithTTL(time.Millisecond)
	defer s.Stop()

	s.Set(&SelectionPlan{SessionID: "s1"})
	time.Sleep(5 * time.Millisecond)
	s.evictStale()
	assert.Equal(t, 0, s.Len())
}

func TestSessionPlanStore_ConcurrentAccess(t *testing.T) {
	s := NewSessionPlanStoreWithTTL(time.Minute)
	defer s.Stop()

	var wg sync.WaitGroup
	for i := 0; i < 50; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			id := "s" + string(rune('A'+n%5))
			s.Set(&SelectionPlan{SessionID: id, Version: "v"})
			s.IncrementCallCount(id)
			s.Get(id)
			s.CallCount(id)
		}(i)
	}
	wg.Wait()

	assert.LessOrEqual(t, s.Len(), 5)
}

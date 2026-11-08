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
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newTestPlanStore(maxEntries int) *SessionPlanStore {
	return NewSessionPlanStoreWithOptions(SessionPlanStoreOptions{MaxEntries: maxEntries})
}

func testPlanKey(sessionID string) PlanKey {
	return NewPlanKey("router-test", sessionID)
}

func issueReceiptForTest(t *testing.T, s *SessionPlanStore, key PlanKey, tool string) AuthorizationReceipt {
	t.Helper()
	plan, ok := s.Get(key)
	require.True(t, ok)
	receipt, err := s.IssueReceipt(key, tool, plan, key.RouterID)
	require.NoError(t, err)
	return receipt
}

func tryIssueReceiptForTest(s *SessionPlanStore, key PlanKey, tool string) AuthorizationReceipt {
	plan, ok := s.Get(key)
	if !ok {
		return AuthorizationReceipt{}
	}
	receipt, err := s.IssueReceipt(key, tool, plan, key.RouterID)
	if err != nil {
		return AuthorizationReceipt{}
	}
	return receipt
}

func activeReceiptCountForTest(s *SessionPlanStore, key PlanKey) int {
	s.mu.RLock()
	defer s.mu.RUnlock()
	if entry, ok := s.entries[key]; ok {
		return len(entry.activeReceipts)
	}
	return 0
}

func receiptVersionRequest(key PlanKey, requested, version string) ReceiptVersionRequest {
	return ReceiptVersionRequest{
		Key:             key,
		Requested:       requested,
		ExpectedVersion: version,
		IdentityHash:    "identity-a",
		ConfigHash:      "config-a",
		CatalogVersion:  "catalog-a",
		ProgressiveHash: "progressive-a",
		RouterID:        key.RouterID,
	}
}

func TestSessionPlanStore_SetGet(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	plan := &SelectionPlan{
		SessionID:        "s1",
		ToolNames:        []string{"a"},
		VisibleToolNames: []string{"a"},
		Reasons:          []DecisionTrace{{Tool: "b"}},
		Version:          "v1",
	}
	s.Set(testPlanKey("s1"), plan, SelectionContext{SessionID: testPlanKey("s1").SessionID})

	got, ok := s.Get(testPlanKey("s1"))
	require.True(t, ok)
	assert.Equal(t, []string{"a"}, got.ToolNames)
	assert.Equal(t, []string{"a"}, got.VisibleToolNames)
	assert.Nil(t, got.Reasons, "stored plans do not retain decision traces")

	_, ok = s.Get(testPlanKey("missing"))
	assert.False(t, ok)
}

func TestSessionPlanStore_CommitReturnsCurrentCommittedPlan(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	key := testPlanKey("s1")
	committedA, err := s.Commit(key, &SelectionPlan{
		SessionID:       "s1",
		ToolNames:       []string{"a"},
		Version:         "version-a",
		IdentityHash:    "identity-a",
		ConfigHash:      "config-a",
		CatalogVersion:  "catalog-a",
		ProgressiveHash: "progressive-a",
	}, SelectionContext{SessionID: key.SessionID})
	require.NoError(t, err)

	_, err = s.Commit(key, &SelectionPlan{
		SessionID:       "s1",
		ToolNames:       []string{"b"},
		Version:         "version-b",
		IdentityHash:    "identity-b",
		ConfigHash:      "config-a",
		CatalogVersion:  "catalog-a",
		ProgressiveHash: "progressive-a",
	}, SelectionContext{SessionID: key.SessionID})
	require.NoError(t, err)

	assert.Equal(t, []string{"a"}, committedA.ToolNames)
	assert.Equal(t, "identity-a", committedA.IdentityHash)
	got, ok := s.Get(key)
	require.True(t, ok)
	assert.Equal(t, []string{"b"}, got.ToolNames)
	assert.Equal(t, "identity-b", got.IdentityHash)
}

func TestSessionPlanStore_CommitAndIssueReceiptUsesCommittedPlan(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	key := testPlanKey("s1")
	committed, receipt, err := s.CommitAndIssueReceipt(key, &SelectionPlan{
		SessionID:       "s1",
		ToolNames:       []string{"a"},
		Version:         "version-a",
		IdentityHash:    "identity-a",
		ConfigHash:      "config-a",
		CatalogVersion:  "catalog-a",
		ProgressiveHash: "progressive-a",
	}, SelectionContext{SessionID: key.SessionID}, "a", key.RouterID)
	require.NoError(t, err)

	assert.Equal(t, []string{"a"}, committed.ToolNames)
	assert.Equal(t, committed.Generation, receipt.PlanGeneration)
	assert.Equal(t, "identity-a", receipt.IdentityHash)
	assert.Equal(t, 1, activeReceiptCountForTest(s, key))

	_, err = s.Commit(key, &SelectionPlan{
		SessionID:       "s1",
		ToolNames:       []string{"b"},
		Version:         "version-b",
		IdentityHash:    "identity-b",
		ConfigHash:      "config-a",
		CatalogVersion:  "catalog-a",
		ProgressiveHash: "progressive-a",
	}, SelectionContext{SessionID: key.SessionID})
	require.NoError(t, err)

	assert.Equal(t, CallSuccessResult{}, s.FinalizeReceipt(receipt, ReceiptSucceeded, 1))
	assert.Equal(t, int64(0), s.CallCount(key), "stale receipt from prior generation must not advance new plan")
}

func TestSessionPlanStore_IssueReceiptAllowsEquivalentGenerationRefresh(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	key := testPlanKey("s1")
	plan := &SelectionPlan{
		SessionID:       "s1",
		ToolNames:       []string{"a"},
		Version:         "v1",
		IdentityHash:    "identity-a",
		ConfigHash:      "config-a",
		CatalogVersion:  "catalog-a",
		ProgressiveHash: "progressive-a",
	}
	require.NoError(t, s.Set(key, plan, SelectionContext{SessionID: key.SessionID}))
	stale, ok := s.Get(key)
	require.True(t, ok)

	require.NoError(t, s.Set(key, plan, SelectionContext{SessionID: key.SessionID}))
	current, ok := s.Get(key)
	require.True(t, ok)
	require.Greater(t, current.Generation, stale.Generation)

	receipt, err := s.IssueReceipt(key, "a", stale, key.RouterID)
	require.NoError(t, err)
	assert.Equal(t, current.Generation, receipt.PlanGeneration)
}

func TestSessionPlanStore_IssueReceiptRejectsDifferentIdentityRefresh(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	key := testPlanKey("s1")
	require.NoError(t, s.Set(key, &SelectionPlan{
		SessionID:       "s1",
		ToolNames:       []string{"a"},
		Version:         "v1",
		IdentityHash:    "identity-a",
		ConfigHash:      "config-a",
		CatalogVersion:  "catalog-a",
		ProgressiveHash: "progressive-a",
	}, SelectionContext{SessionID: key.SessionID}))
	stale, ok := s.Get(key)
	require.True(t, ok)

	require.NoError(t, s.Set(key, &SelectionPlan{
		SessionID:       "s1",
		ToolNames:       []string{"a"},
		Version:         "v2",
		IdentityHash:    "identity-b",
		ConfigHash:      "config-a",
		CatalogVersion:  "catalog-a",
		ProgressiveHash: "progressive-a",
	}, SelectionContext{SessionID: key.SessionID}))

	_, err := s.IssueReceipt(key, "a", stale, key.RouterID)
	assert.ErrorIs(t, err, ErrToolNotAuthorized)
}

func TestSessionPlanStore_IssueReceiptForVersionReportsStale(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	key := testPlanKey("s1")
	require.NoError(t, s.Set(key, &SelectionPlan{
		SessionID:       "s1",
		ToolNames:       []string{"a"},
		Version:         "v1",
		IdentityHash:    "identity-a",
		ConfigHash:      "config-a",
		CatalogVersion:  "catalog-a",
		ProgressiveHash: "progressive-a",
	}, SelectionContext{SessionID: key.SessionID}))

	receipt, stale, err := s.IssueReceiptForVersion(receiptVersionRequest(key, "a", "v1"))
	require.NoError(t, err)
	assert.False(t, stale)
	assert.Equal(t, uint64(1), receipt.PlanGeneration)

	_, stale, err = s.IssueReceiptForVersion(receiptVersionRequest(key, "a", "v2"))
	require.NoError(t, err)
	assert.True(t, stale)

	_, stale, err = s.IssueReceiptForVersion(receiptVersionRequest(key, "b", "v1"))
	assert.False(t, stale)
	assert.ErrorIs(t, err, ErrToolNotAuthorized)
}

func TestSessionPlanStore_SetRejectsInvalidInput(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	assert.ErrorIs(t, s.Set(testPlanKey("ignored"), nil, SelectionContext{SessionID: testPlanKey("ignored").SessionID}), ErrInvalidSelectionPlan)
	assert.ErrorIs(t, s.Set(NewPlanKey("", "missing-router"), &SelectionPlan{SessionID: "missing-router"}, SelectionContext{SessionID: "missing-router"}), ErrInvalidPlanKey)
	assert.ErrorIs(t, s.Set(testPlanKey("s1"), &SelectionPlan{SessionID: "other"}, SelectionContext{SessionID: "s1"}), ErrInvalidSelectionPlan)
	assert.Equal(t, 0, s.Len())
}

func TestSessionPlanStore_MutationIsolation(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	plan := &SelectionPlan{SessionID: "s1", ToolNames: []string{"a", "b"}, VisibleToolNames: []string{"a"}}
	s.Set(testPlanKey("s1"), plan, SelectionContext{SessionID: testPlanKey("s1").SessionID})
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

func TestSessionPlanStore_ContextMutationIsolation(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	key := testPlanKey("s1")
	claims := map[string]any{"roles": []any{"reader"}}
	require.NoError(t, s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}}, SelectionContext{
		SessionID: key.SessionID,
		Claims:    claims,
	}))
	claims["roles"].([]any)[0] = "admin"

	contexts := s.SessionPlanContexts()
	require.Len(t, contexts, 1)
	contexts[0].Context.Claims["roles"].([]any)[0] = "mutated"

	contexts = s.SessionPlanContexts()
	assert.Equal(t, "reader", contexts[0].Context.Claims["roles"].([]any)[0])
}

func TestSessionPlanStore_SetRequiresActiveSessionGenerationWhenProvided(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	key := testPlanKey("s1")
	sc := SelectionContext{SessionID: key.SessionID, SessionGeneration: 7}
	assert.ErrorIs(t, s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}}, sc), ErrSessionPlanNotActive)

	require.NoError(t, s.ActivateSession(key.SessionID, 7))
	require.NoError(t, s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}}, sc))

	s.DeleteSession(key.SessionID)
	assert.ErrorIs(t, s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}}, sc), ErrSessionPlanNotActive)
}

func TestSessionPlanStore_CallCountPreservedAndResetByIdentity(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	key := testPlanKey("s1")
	s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}, IdentityHash: "tenant-a", ProgressiveHash: "p1"}, SelectionContext{SessionID: key.SessionID})
	assert.Equal(t, CallSuccessResult{Count: 1}, s.RecordCallSuccess(issueReceiptForTest(t, s, key, "a"), 3))
	assert.Equal(t, CallSuccessResult{Count: 2}, s.RecordCallSuccess(issueReceiptForTest(t, s, key, "a"), 3))

	s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}, Version: "v2", IdentityHash: "tenant-a", ProgressiveHash: "p1"}, SelectionContext{SessionID: key.SessionID})
	assert.Equal(t, int64(2), s.CallCount(key))

	s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}, Version: "v3", IdentityHash: "tenant-b", ProgressiveHash: "p1"}, SelectionContext{SessionID: key.SessionID})
	assert.Equal(t, int64(0), s.CallCount(key))
	_, expanded := s.CallState(key, "tenant-b", "p1")
	assert.False(t, expanded)
}

func TestSessionPlanStore_RecordCallSuccessTransitionOnce(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	key := testPlanKey("s1")
	s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}}, SelectionContext{SessionID: key.SessionID})
	assert.Equal(t, CallSuccessResult{Count: 1}, s.RecordCallSuccess(issueReceiptForTest(t, s, key, "a"), 2))
	assert.Equal(t, CallSuccessResult{Count: 2, Transitioned: true}, s.RecordCallSuccess(issueReceiptForTest(t, s, key, "a"), 2))
	assert.Equal(t, CallSuccessResult{Count: 3}, s.RecordCallSuccess(issueReceiptForTest(t, s, key, "a"), 2))

	got, ok := s.Get(key)
	require.True(t, ok)
	assert.True(t, got.Expanded)
}

func TestSessionPlanStore_ProgressiveTransitionInvalidatesOldPlanContext(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	key := testPlanKey("s1")
	sc := SelectionContext{SessionID: key.SessionID, SessionGeneration: 7}
	require.NoError(t, s.ActivateSession(key.SessionID, sc.SessionGeneration))
	require.NoError(t, s.Set(key, &SelectionPlan{
		SessionID:       "s1",
		ToolNames:       []string{"starter"},
		Version:         "v1",
		IdentityHash:    "identity-a",
		ConfigHash:      "config-a",
		CatalogVersion:  "catalog-a",
		ProgressiveHash: "progressive-a",
	}, sc))
	contexts := s.SessionPlanContexts()
	require.Len(t, contexts, 1)
	base := contexts[0]

	transitionReceipt := issueReceiptForTest(t, s, key, "starter")
	staleReceipt := issueReceiptForTest(t, s, key, "starter")
	assert.Equal(t, 2, activeReceiptCountForTest(s, key))

	assert.Equal(t, CallSuccessResult{Count: 1, Transitioned: true}, s.RecordCallSuccess(transitionReceipt, 1))
	assert.Equal(t, 0, activeReceiptCountForTest(s, key), "transition releases receipts from the superseded plan generation")
	assert.Equal(t, CallSuccessResult{}, s.RecordCallSuccess(staleReceipt, 1))

	got, ok := s.Get(key)
	require.True(t, ok)
	require.Greater(t, got.Generation, base.PlanGeneration)
	require.True(t, got.Expanded)

	committed, ok, err := s.SetIfCurrent(base, &SelectionPlan{
		SessionID:       "s1",
		ToolNames:       []string{"stale"},
		Version:         "v-stale",
		IdentityHash:    "identity-a",
		ConfigHash:      "config-a",
		CatalogVersion:  "catalog-a",
		ProgressiveHash: "progressive-a",
	})
	require.NoError(t, err)
	assert.False(t, ok)
	assert.Nil(t, committed)

	got, ok = s.Get(key)
	require.True(t, ok)
	assert.Equal(t, []string{"starter"}, got.ToolNames)
	assert.True(t, got.Expanded)
}

func TestSessionPlanStore_RecordCallSuccessIgnoresUnknownOrOutsidePlan(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	assert.Equal(t, CallSuccessResult{}, s.RecordCallSuccess(AuthorizationReceipt{}, 1))
	key := testPlanKey("s1")
	s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}}, SelectionContext{SessionID: key.SessionID})
	assert.Equal(t, CallSuccessResult{}, s.RecordCallSuccess(tryIssueReceiptForTest(s, key, "ghost"), 1))
	assert.Equal(t, int64(0), s.CallCount(key))
}

func TestSessionPlanStore_FinalizeReceiptAbortReleasesWithoutProgress(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	key := testPlanKey("s1")
	require.NoError(t, s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}}, SelectionContext{SessionID: key.SessionID}))
	receipt := issueReceiptForTest(t, s, key, "a")
	assert.Equal(t, 1, activeReceiptCountForTest(s, key))

	assert.Equal(t, CallSuccessResult{}, s.FinalizeReceipt(receipt, ReceiptAborted, 1))
	assert.Equal(t, 0, activeReceiptCountForTest(s, key))
	assert.Equal(t, int64(0), s.CallCount(key))
	assert.Equal(t, CallSuccessResult{}, s.RecordCallSuccess(receipt, 1), "aborted receipt must not be replayable")
}

func TestSessionPlanStore_FinalizeReceiptSuccessReleasesAndProgresses(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	key := testPlanKey("s1")
	require.NoError(t, s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}}, SelectionContext{SessionID: key.SessionID}))
	receipt := issueReceiptForTest(t, s, key, "a")

	assert.Equal(t, CallSuccessResult{Count: 1, Transitioned: true}, s.FinalizeReceipt(receipt, ReceiptSucceeded, 1))
	assert.Equal(t, 0, activeReceiptCountForTest(s, key))
	assert.Equal(t, int64(1), s.CallCount(key))
	assert.Equal(t, CallSuccessResult{}, s.FinalizeReceipt(receipt, ReceiptSucceeded, 1), "finalize must be idempotent for progress")
	assert.Equal(t, int64(1), s.CallCount(key))
}

func TestSessionPlanStore_FinalizeReceiptWithoutProgressiveStillReleases(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	key := testPlanKey("s1")
	require.NoError(t, s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}}, SelectionContext{SessionID: key.SessionID}))
	receipt := issueReceiptForTest(t, s, key, "a")

	assert.Equal(t, CallSuccessResult{}, s.FinalizeReceipt(receipt, ReceiptSucceeded, 0))
	assert.Equal(t, 0, activeReceiptCountForTest(s, key))
	assert.Equal(t, int64(0), s.CallCount(key))
}

func TestSessionPlanStore_FinalizeReceiptAbortManyFailuresDoesNotLeak(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	key := testPlanKey("s1")
	require.NoError(t, s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}}, SelectionContext{SessionID: key.SessionID}))
	for i := 0; i < 10000; i++ {
		receipt := issueReceiptForTest(t, s, key, "a")
		s.FinalizeReceipt(receipt, ReceiptAborted, 1)
	}

	assert.Equal(t, 0, activeReceiptCountForTest(s, key))
	assert.Equal(t, int64(0), s.CallCount(key))
}

func TestSessionPlanStore_ReceiptBookkeepingBoundedAfterManyCompletions(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	key := testPlanKey("s1")
	require.NoError(t, s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}}, SelectionContext{SessionID: key.SessionID}))
	for i := 0; i < 100000; i++ {
		result := s.RecordCallSuccess(issueReceiptForTest(t, s, key, "a"), 200000)
		require.Equal(t, int64(i+1), result.Count)
	}
	assert.Equal(t, 0, activeReceiptCountForTest(s, key))
	assert.Equal(t, int64(100000), s.CallCount(key))
}

func TestSessionPlanStore_ConcurrentReceiptCompletionCountsOnce(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	key := testPlanKey("s1")
	require.NoError(t, s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}}, SelectionContext{SessionID: key.SessionID}))
	receipt := issueReceiptForTest(t, s, key, "a")

	start := make(chan struct{})
	results := make(chan CallSuccessResult, 2)
	for i := 0; i < 2; i++ {
		go func() {
			<-start
			results <- s.RecordCallSuccess(receipt, 10)
		}()
	}
	close(start)

	first := <-results
	second := <-results
	assert.ElementsMatch(t, []CallSuccessResult{{Count: 1}, {}}, []CallSuccessResult{first, second})
	assert.Equal(t, int64(1), s.CallCount(key))
	assert.Equal(t, 0, activeReceiptCountForTest(s, key))
	assert.Equal(t, CallSuccessResult{}, s.RecordCallSuccess(receipt, 10))
}

func TestSessionPlanStore_AbortOneReceiptDoesNotConsumeAnother(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	key := testPlanKey("s1")
	require.NoError(t, s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}}, SelectionContext{SessionID: key.SessionID}))
	receipt := issueReceiptForTest(t, s, key, "a")
	assert.Equal(t, 1, activeReceiptCountForTest(s, key))

	other := issueReceiptForTest(t, s, key, "a")
	assert.Equal(t, 2, activeReceiptCountForTest(s, key))
	assert.Equal(t, CallSuccessResult{}, s.FinalizeReceipt(receipt, ReceiptAborted, 10))
	assert.Equal(t, 1, activeReceiptCountForTest(s, key), "aborting one receipt must not remove another in-flight call")

	assert.Equal(t, CallSuccessResult{Count: 1}, s.RecordCallSuccess(other, 10))
	assert.Equal(t, 0, activeReceiptCountForTest(s, key))
}

func TestSessionPlanStore_DeleteSessionReleasesReceipts(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	key := testPlanKey("s1")
	require.NoError(t, s.Set(key, &SelectionPlan{SessionID: "s1", ToolNames: []string{"a"}}, SelectionContext{SessionID: key.SessionID}))
	receipt := issueReceiptForTest(t, s, key, "a")
	assert.Equal(t, 1, activeReceiptCountForTest(s, key))

	s.DeleteSession(key.SessionID)
	assert.Equal(t, 0, s.Len())
	assert.Equal(t, 0, activeReceiptCountForTest(s, key))
	assert.Equal(t, CallSuccessResult{}, s.RecordCallSuccess(receipt, 10))
}

func TestSessionPlanStore_Delete(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	key := testPlanKey("s1")
	s.Set(key, &SelectionPlan{SessionID: "s1"}, SelectionContext{SessionID: key.SessionID})
	s.Delete(key)
	_, ok := s.Get(key)
	assert.False(t, ok)
}

func TestSessionPlanStore_IsolatesSameSessionAcrossRouters(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	keyA := NewPlanKey("router-a", "shared-session")
	keyB := NewPlanKey("router-b", "shared-session")
	s.Set(keyA, &SelectionPlan{SessionID: "shared-session", ToolNames: []string{"a"}}, SelectionContext{SessionID: keyA.SessionID})
	s.Set(keyB, &SelectionPlan{SessionID: "shared-session", ToolNames: []string{"b"}}, SelectionContext{SessionID: keyB.SessionID})

	gotA, ok := s.Get(keyA)
	require.True(t, ok)
	assert.Equal(t, []string{"a"}, gotA.ToolNames)

	gotB, ok := s.Get(keyB)
	require.True(t, ok)
	assert.Equal(t, []string{"b"}, gotB.ToolNames)

	assert.Equal(t, CallSuccessResult{Count: 1}, s.RecordCallSuccess(issueReceiptForTest(t, s, keyA, "a"), 3))
	assert.Equal(t, int64(1), s.CallCount(keyA))
	assert.Equal(t, int64(0), s.CallCount(keyB))

	s.DeleteSession("shared-session")
	assert.Equal(t, 0, s.Len())
}

func TestSessionPlanStore_DeleteOneRouterKeepsOtherRouter(t *testing.T) {
	s := newTestPlanStore(10)
	defer s.Stop()

	keyA := NewPlanKey("router-a", "shared-session")
	keyB := NewPlanKey("router-b", "shared-session")
	s.Set(keyA, &SelectionPlan{SessionID: "shared-session", ToolNames: []string{"a"}}, SelectionContext{SessionID: keyA.SessionID})
	s.Set(keyB, &SelectionPlan{SessionID: "shared-session", ToolNames: []string{"b"}}, SelectionContext{SessionID: keyB.SessionID})

	s.Delete(keyA)

	_, ok := s.Get(keyA)
	assert.False(t, ok)
	gotB, ok := s.Get(keyB)
	require.True(t, ok)
	assert.Equal(t, []string{"b"}, gotB.ToolNames)
}

func TestSessionPlanStore_CapacityRejectsNewPlan(t *testing.T) {
	s := newTestPlanStore(2)
	defer s.Stop()

	require.NoError(t, s.Set(testPlanKey("s2"), &SelectionPlan{SessionID: "s2"}, SelectionContext{SessionID: testPlanKey("s2").SessionID}))
	require.NoError(t, s.Set(testPlanKey("s1"), &SelectionPlan{SessionID: "s1"}, SelectionContext{SessionID: testPlanKey("s1").SessionID}))
	assert.ErrorIs(t, s.Set(testPlanKey("s3"), &SelectionPlan{SessionID: "s3"}, SelectionContext{SessionID: testPlanKey("s3").SessionID}), ErrPlanStoreFull)

	_, ok := s.Get(testPlanKey("s1"))
	assert.True(t, ok)
	_, ok = s.Get(testPlanKey("s2"))
	assert.True(t, ok)
	_, ok = s.Get(testPlanKey("s3"))
	assert.False(t, ok)
}

func TestSessionPlanStore_CapacityDoesNotExpireActivePlan(t *testing.T) {
	s := newTestPlanStore(2)
	defer s.Stop()

	require.NoError(t, s.Set(testPlanKey("active"), &SelectionPlan{SessionID: "active"}, SelectionContext{SessionID: testPlanKey("active").SessionID}))
	require.NoError(t, s.Set(testPlanKey("fresh"), &SelectionPlan{SessionID: "fresh"}, SelectionContext{SessionID: testPlanKey("fresh").SessionID}))
	assert.ErrorIs(t, s.Set(testPlanKey("new"), &SelectionPlan{SessionID: "new"}, SelectionContext{SessionID: testPlanKey("new").SessionID}), ErrPlanStoreFull)

	_, ok := s.Get(testPlanKey("active"))
	assert.True(t, ok)
	assert.Equal(t, 2, s.Len())
}

func TestSessionPlanStore_SetMaxEntries(t *testing.T) {
	s := newTestPlanStore(5)
	defer s.Stop()

	for i := 0; i < 5; i++ {
		sessionID := fmt.Sprintf("s%d", i)
		s.Set(testPlanKey(sessionID), &SelectionPlan{SessionID: sessionID}, SelectionContext{SessionID: testPlanKey(sessionID).SessionID})
	}
	s.SetMaxEntries(2)
	assert.Equal(t, 5, s.Len())
	assert.ErrorIs(t, s.Set(testPlanKey("new"), &SelectionPlan{SessionID: "new"}, SelectionContext{SessionID: "new"}), ErrPlanStoreFull)
}

func TestSessionPlanStore_ConcurrentAccess(t *testing.T) {
	s := newTestPlanStore(20)
	defer s.Stop()

	var wg sync.WaitGroup
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func(n int) {
			defer wg.Done()
			id := fmt.Sprintf("s%d", n%10)
			key := testPlanKey(id)
			s.Set(key, &SelectionPlan{SessionID: id, ToolNames: []string{"a"}, Version: "v"}, SelectionContext{SessionID: key.SessionID})
			s.RecordCallSuccess(tryIssueReceiptForTest(s, key, "a"), 1000)
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
	s := newTestPlanStore(10)
	s.Set(testPlanKey("s1"), &SelectionPlan{SessionID: "s1"}, SelectionContext{SessionID: testPlanKey("s1").SessionID})

	s.Stop()
	s.Stop()
	assert.Equal(t, 0, s.Len())
}

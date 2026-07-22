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
	"encoding/json"
	"strings"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestStageDropCounts(t *testing.T) {
	traces := []DecisionTrace{
		{Tool: "a", Stage: StagePolicy},
		{Tool: "b", Stage: StagePolicy},
		{Tool: "c", Stage: StageWorkflow},
	}

	counts := stageDropCounts(traces)
	assert.Equal(t, 2, counts[StagePolicy])
	assert.Equal(t, 1, counts[StageWorkflow])
}

func TestStageDropCounts_Empty(t *testing.T) {
	assert.Nil(t, stageDropCounts(nil))
}

func TestDeniedSamples_Capped(t *testing.T) {
	traces := make([]DecisionTrace, 0, 20)
	for i := 0; i < 20; i++ {
		traces = append(traces, DecisionTrace{Tool: "t", Stage: StagePolicy})
	}
	assert.Len(t, deniedSamples(traces), maxDeniedSamples)
}

func TestDeniedSamples_OnlyDropped(t *testing.T) {
	traces := []DecisionTrace{
		{Tool: "dropped"},
	}
	samples := deniedSamples(traces)
	assert.Equal(t, []string{"dropped"}, samples)
}

func TestDecisionLogger_NilSafe(t *testing.T) {
	var d *DecisionLogger
	// Must not panic on a nil logger or nil plan.
	d.Log(SelectionContext{}, nil, 0)
	NewDecisionLogger(1.0, false).Log(SelectionContext{SessionID: "s"}, nil, 0)
}

func TestDecisionLogger_ShouldSample(t *testing.T) {
	assert.False(t, NewDecisionLogger(0, false).shouldSample()) // 0 => off
	assert.True(t, NewDecisionLogger(1, false).shouldSample())  // 1 => always
}

func TestDecisionLogger_LogEmitsWithoutPanic(t *testing.T) {
	d := NewDecisionLogger(1.0, false)
	plan := &SelectionPlan{
		SessionID:      "s1",
		ToolNames:      []string{"a"},
		Mode:           ModeSelected,
		Version:        "identity-derived-v1",
		CatalogVersion: "catalog-v1",
		ConfigHash:     "config-v1",
		Reasons: []DecisionTrace{
			{Tool: "b", Stage: StagePolicy, Detail: "deny_tag:x"},
		},
	}
	d.Log(SelectionContext{SessionID: "s1", Tenant: "acme", Method: "tools/list"}, plan, 2)
}

func TestDecisionLogger_RecordOmitsDeniedSamplesUnlessDecisionDetailLogging(t *testing.T) {
	plan := &SelectionPlan{
		SessionID:      "s1",
		ToolNames:      []string{"a"},
		Mode:           ModeSelected,
		Version:        "identity-derived-v1",
		CatalogVersion: "catalog-v1",
		ConfigHash:     "config-v1",
		Reasons: []DecisionTrace{
			{Tool: "b", Stage: StagePolicy},
		},
	}

	normal := NewDecisionLogger(1.0, false).record(SelectionContext{SessionID: "s1"}, plan, 2)
	assert.Empty(t, normal.DeniedSamples)

	detailed := NewDecisionLogger(1.0, true).record(SelectionContext{SessionID: "s1"}, plan, 2)
	assert.Equal(t, []string{"b"}, detailed.DeniedSamples)
}

func TestDecisionLogger_RecordDoesNotSerializeIdentityOrPayload(t *testing.T) {
	plan := &SelectionPlan{
		SessionID:      "raw-session-id",
		ToolNames:      []string{"allowed"},
		Mode:           ModeSelected,
		Version:        "identity-fnv-abcd",
		CatalogVersion: "catalog-version",
		ConfigHash:     "config-version",
		IdentityHash:   "identity-fnv",
		Reasons: []DecisionTrace{
			{Tool: "hidden_tool", Stage: StagePolicy, Rule: "tenant-rule", Detail: "no_allow_tag"},
		},
	}
	sc := SelectionContext{
		SessionID: "raw-session-id",
		Tenant:    "tenant-acme",
		UserID:    "subject-123",
		Claims: map[string]any{
			"tenant": "tenant-acme",
			"sub":    "subject-123",
			"token":  "secret-token",
		},
		Method:    "tools/list",
		Requested: "hidden_tool",
	}

	normal := NewDecisionLogger(1.0, false).record(sc, plan, 2)
	data, err := json.Marshal(normal)
	assert.NoError(t, err)
	text := string(data)
	for _, forbidden := range []string{"raw-session-id", "tenant-acme", "subject-123", "agent-client", "secret-token", "hidden_tool", "tenant-rule", "identity-fnv", "identity-fnv-abcd", "plan_version"} {
		assert.False(t, strings.Contains(text, forbidden), "default decision log leaked %q: %s", forbidden, text)
	}
	assert.Contains(t, text, "catalog-version")
	assert.Contains(t, text, "config-version")

	detailed := NewDecisionLogger(1.0, true).record(sc, plan, 2)
	data, err = json.Marshal(detailed)
	assert.NoError(t, err)
	text = string(data)
	assert.Contains(t, text, "hidden_tool")
	for _, forbidden := range []string{"raw-session-id", "tenant-acme", "subject-123", "agent-client", "secret-token", "identity-fnv", "identity-fnv-abcd", "plan_version"} {
		assert.False(t, strings.Contains(text, forbidden), "decision detail logging leaked identity %q: %s", forbidden, text)
	}
}

func TestMetricHelpers_NilSafeBeforeInit(t *testing.T) {
	// These must be no-ops if initMetrics has not run (defensive).
	// We cannot un-init, so just ensure calling them after init is safe.
	initMetrics()
	recordSelection("ok", ModeSelected, 10, 3, 1.5)
	recordFallback(SelectionOutcomeNoMatch)
	recordCallDenied("not_in_plan")
	store := NewSessionPlanStore()
	setPlansActive(store, 5)
	store.Stop()
}

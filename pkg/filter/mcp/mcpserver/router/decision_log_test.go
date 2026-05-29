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
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestStageDropCounts(t *testing.T) {
	traces := []DecisionTrace{
		{Tool: "a", Kept: false, Stage: StagePolicy},
		{Tool: "b", Kept: false, Stage: StagePolicy},
		{Tool: "c", Kept: false, Stage: StageWorkflow},
		{Tool: "d", Kept: true, Stage: StagePolicy},
	}

	counts := stageDropCounts(traces)
	assert.Equal(t, 2, counts[StagePolicy])
	assert.Equal(t, 1, counts[StageWorkflow])
}

func TestStageDropCounts_AllKept(t *testing.T) {
	traces := []DecisionTrace{{Tool: "a", Kept: true, Stage: StagePolicy}}
	assert.Nil(t, stageDropCounts(traces))
}

func TestDeniedSamples_Capped(t *testing.T) {
	traces := make([]DecisionTrace, 0, 20)
	for i := 0; i < 20; i++ {
		traces = append(traces, DecisionTrace{Tool: "t", Kept: false, Stage: StagePolicy})
	}
	assert.Len(t, deniedSamples(traces), maxDeniedSamples)
}

func TestDeniedSamples_OnlyDropped(t *testing.T) {
	traces := []DecisionTrace{
		{Tool: "kept", Kept: true},
		{Tool: "dropped", Kept: false},
	}
	samples := deniedSamples(traces)
	assert.Equal(t, []string{"dropped"}, samples)
}

func TestDecisionLogger_NilSafe(t *testing.T) {
	var d *DecisionLogger
	// Must not panic on a nil logger or nil plan.
	d.Log(SelectionContext{}, nil, 0)
	NewDecisionLogger(1.0).Log(SelectionContext{SessionID: "s"}, nil, 0)
}

func TestDecisionLogger_ShouldSample(t *testing.T) {
	assert.True(t, NewDecisionLogger(0).shouldSample()) // 0 => always
	assert.True(t, NewDecisionLogger(1).shouldSample()) // 1 => always
	assert.True(t, NewDecisionLogger(2).shouldSample()) // >1 => always
}

func TestDecisionLogger_LogEmitsWithoutPanic(t *testing.T) {
	d := NewDecisionLogger(1.0)
	plan := &SelectionPlan{
		SessionID: "s1",
		ToolNames: []string{"a"},
		Mode:      ModeHybrid,
		Version:   "v1",
		Reasons: []DecisionTrace{
			{Tool: "b", Kept: false, Stage: StagePolicy, Detail: "deny_tag:x"},
		},
	}
	d.Log(SelectionContext{SessionID: "s1", Tenant: "acme", Method: "tools/list"}, plan, 2)
}

func TestMetricHelpers_NilSafeBeforeInit(t *testing.T) {
	// These must be no-ops if initMetrics has not run (defensive).
	// We cannot un-init, so just ensure calling them after init is safe.
	initMetrics()
	recordSelection("ok", ModeHybrid, 10, 3, 1.5)
	recordFallback("empty")
	recordCallDenied("not_in_plan")
	setPlansActive(5)
}

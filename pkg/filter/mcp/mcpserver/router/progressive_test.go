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
	"context"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestProgressiveGate_InitialBundleThenExpand(t *testing.T) {
	wf, err := NewWorkflowSelector([]model.WorkflowConfig{
		{Name: "starter", Tools: []string{"ping", "help"}},
	})
	require.NoError(t, err)

	gate := NewProgressiveGate(model.ProgressiveConfig{InitialBundle: "starter", ExpandAfterCalls: 1}, wf)
	store := NewSessionPlanStoreWithTTL(time.Minute)
	defer store.Stop()
	store.Set(&SelectionPlan{SessionID: "s1"})

	tools := testTools("ping", "help", "advanced1", "advanced2")

	// Before any calls: only initial bundle visible.
	out, _ := gate.Apply(tools, SelectionContext{SessionID: "s1"}, store)
	assert.ElementsMatch(t, []string{"ping", "help"}, keptNames(out))

	// After one successful call: full set visible.
	store.IncrementCallCount("s1")
	out2, _ := gate.Apply(tools, SelectionContext{SessionID: "s1"}, store)
	assert.Len(t, out2, 4)
}

func TestProgressiveGate_NoBundlePassthrough(t *testing.T) {
	wf, _ := NewWorkflowSelector(nil)
	gate := NewProgressiveGate(model.ProgressiveConfig{InitialBundle: "missing"}, wf)
	store := NewSessionPlanStoreWithTTL(time.Minute)
	defer store.Stop()
	store.Set(&SelectionPlan{SessionID: "s1"})

	tools := testTools("a", "b")
	out, _ := gate.Apply(tools, SelectionContext{SessionID: "s1"}, store)
	assert.Equal(t, tools, out)
}

func TestProgressiveGate_ViaCompositeExpands(t *testing.T) {
	cfg := &model.RouterConfig{
		Enabled: true,
		Stages:  model.RouterStages{Progressive: true},
		Workflows: []model.WorkflowConfig{
			{Name: "starter", Tools: []string{"t0"}},
		},
		Progressive: model.ProgressiveConfig{InitialBundle: "starter", ExpandAfterCalls: 1},
	}
	store := NewSessionPlanStoreWithTTL(time.Minute)
	sel, err := Build(cfg, store)
	require.NoError(t, err)
	cs := sel.(*CompositeSelector)
	defer store.Stop()

	tools := testTools("t0", "t1", "t2")

	// First tools/list: only starter bundle.
	p1, _ := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, tools)
	assert.Equal(t, []string{"t0"}, p1.ToolNames)

	// Authorize + call t0 to bump the counter.
	require.NoError(t, cs.AuthorizeCall(context.Background(), SelectionContext{SessionID: "s1", Requested: "t0"}))

	// Add a tool to force version change so the plan recomputes, then expand.
	p2, _ := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, append(tools, model.ToolConfig{Name: "t3"}))
	assert.Len(t, p2.ToolNames, 4)
}

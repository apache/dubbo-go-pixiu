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
	tools := testTools("ping", "help", "advanced1", "advanced2")

	// Before any calls: only initial bundle visible.
	out, _ := gate.Apply(tools, false)
	assert.ElementsMatch(t, []string{"ping", "help"}, keptNames(out))

	// After expansion: full set visible.
	out2, _ := gate.Apply(tools, true)
	assert.Len(t, out2, 4)
}

func TestProgressiveGate_MissingBundleFailsClosed(t *testing.T) {
	wf, _ := NewWorkflowSelector(nil)
	gate := NewProgressiveGate(model.ProgressiveConfig{InitialBundle: "missing"}, wf)

	tools := testTools("a", "b")
	out, traces := gate.Apply(tools, false)
	assert.Nil(t, out)
	require.Len(t, traces, 1)
	assert.Equal(t, "missing_initial_bundle", traces[0].Detail)
}

func TestProgressiveGate_ViaCompositeExpands(t *testing.T) {
	cfg := &model.RouterConfig{
		Enabled:  true,
		Fallback: FallbackFailClosed,
		Stages:   model.RouterStages{Progressive: true},
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

	// Authorization alone does not bump the counter; only a completed call does.
	require.NoError(t, cs.AuthorizeCall(context.Background(), SelectionContext{SessionID: "s1", Requested: "t0"}, tools))
	result, err := cs.RecordCallSuccess(context.Background(), SelectionContext{SessionID: "s1", Requested: "t0"})
	require.NoError(t, err)
	assert.True(t, result.Transitioned)

	// Second tools/list with same tools: crossing the threshold invalidates cache, full set revealed.
	p2, _ := cs.Select(context.Background(), SelectionContext{SessionID: "s1"}, tools)
	assert.Len(t, p2.ToolNames, 3, "after crossing expand threshold, all tools should be visible")
	assert.ElementsMatch(t, []string{"t0", "t1", "t2"}, p2.ToolNames)
}

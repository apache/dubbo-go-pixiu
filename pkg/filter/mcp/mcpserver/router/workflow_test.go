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
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestWorkflowSelector_MatchKeepsBundle(t *testing.T) {
	ws, err := NewWorkflowSelector([]model.WorkflowConfig{
		{
			Name:  "support",
			Tools: []string{"search_kb", "create_ticket"},
			When:  model.PolicyMatch{Claim: "role", Equals: "support"},
		},
	})
	require.NoError(t, err)

	tools := testTools("search_kb", "create_ticket", "admin_delete")
	out, traces := ws.Filter(tools, SelectionContext{Claims: map[string]any{"role": "support"}})

	assert.ElementsMatch(t, []string{"search_kb", "create_ticket"}, keptNames(out))
	require.Len(t, traces, 1)
	assert.Equal(t, "admin_delete", traces[0].Tool)
	assert.Equal(t, "support", traces[0].Rule)
}

func TestWorkflowSelector_NoMatchPassthrough(t *testing.T) {
	ws, err := NewWorkflowSelector([]model.WorkflowConfig{
		{Name: "support", Tools: []string{"a"}, When: model.PolicyMatch{Claim: "role", Equals: "support"}},
	})
	require.NoError(t, err)

	tools := testTools("a", "b", "c")
	out, traces := ws.Filter(tools, SelectionContext{Claims: map[string]any{"role": "analyst"}})

	assert.Equal(t, tools, out)
	assert.Nil(t, traces)
}

func TestWorkflowSelector_BundleWithoutWhenNotAutoMatched(t *testing.T) {
	ws, err := NewWorkflowSelector([]model.WorkflowConfig{
		{Name: "safe-minimal", Tools: []string{"ping"}}, // no When
	})
	require.NoError(t, err)

	tools := testTools("ping", "other")
	out, _ := ws.Filter(tools, SelectionContext{})

	// No When => never auto-matches => passthrough.
	assert.Equal(t, tools, out)

	// But it is addressable by name.
	bundle, ok := ws.bundleTools("safe-minimal")
	require.True(t, ok)
	_, hasPing := bundle["ping"]
	assert.True(t, hasPing)
}

func TestWorkflowSelector_FirstMatchWins(t *testing.T) {
	ws, err := NewWorkflowSelector([]model.WorkflowConfig{
		{Name: "first", Tools: []string{"a"}, When: model.PolicyMatch{Claim: "tenant", Equals: "acme"}},
		{Name: "second", Tools: []string{"b"}, When: model.PolicyMatch{Claim: "tenant", Equals: "acme"}},
	})
	require.NoError(t, err)

	tools := testTools("a", "b")
	out, _ := ws.Filter(tools, SelectionContext{Tenant: "acme"})

	assert.Equal(t, []string{"a"}, keptNames(out))
}

func TestWorkflowSelector_UnknownBundle(t *testing.T) {
	ws, err := NewWorkflowSelector(nil)
	require.NoError(t, err)
	_, ok := ws.bundleTools("ghost")
	assert.False(t, ok)
}

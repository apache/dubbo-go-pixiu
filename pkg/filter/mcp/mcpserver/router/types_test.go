/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License. You may obtain a copy of the License at
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

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func testTools(names ...string) []model.ToolConfig {
	tools := make([]model.ToolConfig, len(names))
	for i, n := range names {
		tools[i] = model.ToolConfig{Name: n, Cluster: "test"}
	}
	return tools
}

func TestSelectionPlan_Contains(t *testing.T) {
	plan := &SelectionPlan{ToolNames: []string{"x", "y"}}

	assert.True(t, plan.Contains("x"))
	assert.True(t, plan.Contains("y"))
	assert.False(t, plan.Contains("z"))

	var nilPlan *SelectionPlan
	assert.False(t, nilPlan.Contains("x"))
}

func TestSelectionPlan_VisibleNamesDefaultsToAuthorizedNames(t *testing.T) {
	plan := &SelectionPlan{ToolNames: []string{"x", "y"}}

	assert.Equal(t, []string{"x", "y"}, plan.VisibleNames())

	plan.VisibleToolNames = []string{"x"}
	assert.Equal(t, []string{"x"}, plan.VisibleNames())

	var nilPlan *SelectionPlan
	assert.Nil(t, nilPlan.VisibleNames())
}

func TestDefaultVisibleFingerprintUsesVisibleNamesOnly(t *testing.T) {
	a := model.ToolConfig{Name: "a", Description: "A"}
	b := model.ToolConfig{Name: "b", Description: "B"}
	hidden := false
	hiddenTool := model.ToolConfig{Name: "hidden", Meta: &model.ToolMeta{DiscoveryVisibility: &hidden}}

	assert.Equal(t,
		DefaultVisibleFingerprint([]model.ToolConfig{a, b}),
		DefaultVisibleFingerprint([]model.ToolConfig{a, b, hiddenTool}))
	assert.NotEqual(t,
		DefaultVisibleFingerprint([]model.ToolConfig{a, b}),
		DefaultVisibleFingerprint([]model.ToolConfig{b, a}))
}

func TestBuild_NilConfigRequiresStore(t *testing.T) {
	s, err := Build(nil, nil)
	assert.Error(t, err)
	assert.Nil(t, s)
}

func TestBuild_NilConfigReturnsComposite(t *testing.T) {
	store := NewSessionPlanStore()
	defer store.Stop()
	s, err := Build(nil, store)
	assert.NoError(t, err)
	assert.NotNil(t, s)
}

func TestBuild_EnabledReturnsComposite(t *testing.T) {
	store := NewSessionPlanStore()
	defer store.Stop()

	s, err := Build(&model.RouterConfig{Fallback: FallbackFailClosed}, store)

	assert.NoError(t, err)
	assert.NotNil(t, s)
}

func TestToolNames(t *testing.T) {
	assert.Equal(t, []string{"a", "b"}, toolNames(testTools("a", "b")))
	assert.Empty(t, toolNames(nil))
}

func TestVisibleToolNamesSkipsHiddenDiscovery(t *testing.T) {
	hidden := false
	tools := []model.ToolConfig{
		{Name: "visible"},
		{Name: "hidden", Meta: &model.ToolMeta{DiscoveryVisibility: &hidden}},
	}

	assert.Equal(t, []string{"visible"}, visibleToolNames(tools))
}

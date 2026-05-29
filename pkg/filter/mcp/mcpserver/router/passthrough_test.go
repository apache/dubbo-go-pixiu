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

func TestPassthroughSelector_SelectReturnsAllTools(t *testing.T) {
	s := NewPassthroughSelector()
	candidates := testTools("a", "b", "c")

	plan, err := s.Select(context.Background(), SelectionContext{SessionID: "sess-1"}, candidates)

	assert.NoError(t, err)
	assert.NotNil(t, plan)
	assert.Equal(t, "sess-1", plan.SessionID)
	assert.Equal(t, ModePassthrough, plan.Mode)
	assert.Equal(t, []string{"a", "b", "c"}, plan.ToolNames)
	assert.Positive(t, plan.CreatedAt)
}

func TestPassthroughSelector_SelectEmptyCandidates(t *testing.T) {
	s := NewPassthroughSelector()

	plan, err := s.Select(context.Background(), SelectionContext{SessionID: "sess-2"}, nil)

	assert.NoError(t, err)
	assert.NotNil(t, plan)
	assert.Empty(t, plan.ToolNames)
}

func TestPassthroughSelector_AuthorizeCallAlwaysAllows(t *testing.T) {
	s := NewPassthroughSelector()

	err := s.AuthorizeCall(context.Background(), SelectionContext{Requested: "anything"})

	assert.NoError(t, err)
}

func TestPassthroughSelector_OnInitializeNoop(t *testing.T) {
	s := NewPassthroughSelector()

	err := s.OnInitialize(context.Background(), SelectionContext{}, testTools("a"))

	assert.NoError(t, err)
}

func TestPassthroughSelector_Name(t *testing.T) {
	assert.Equal(t, ModePassthrough, NewPassthroughSelector().Name())
}

func TestSelectionPlan_Contains(t *testing.T) {
	plan := &SelectionPlan{ToolNames: []string{"x", "y"}}

	assert.True(t, plan.Contains("x"))
	assert.True(t, plan.Contains("y"))
	assert.False(t, plan.Contains("z"))

	var nilPlan *SelectionPlan
	assert.False(t, nilPlan.Contains("x"))
}

func TestBuild_NilConfigReturnsNilSelector(t *testing.T) {
	s, err := Build(nil, nil)
	assert.NoError(t, err)
	assert.Nil(t, s)
}

func TestBuild_DisabledReturnsNilSelector(t *testing.T) {
	s, err := Build(&model.RouterConfig{Enabled: false}, nil)
	assert.NoError(t, err)
	assert.Nil(t, s)
}

func TestBuild_EnabledReturnsComposite(t *testing.T) {
	store := NewSessionPlanStoreWithTTL(time.Minute)
	defer store.Stop()
	s, err := Build(&model.RouterConfig{Enabled: true}, store)
	assert.NoError(t, err)
	assert.NotNil(t, s)
	assert.Equal(t, "composite", s.Name())
}

func TestToolNames(t *testing.T) {
	assert.Equal(t, []string{"a", "b"}, toolNames(testTools("a", "b")))
	assert.Empty(t, toolNames(nil))
}

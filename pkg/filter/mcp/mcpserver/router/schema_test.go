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

func TestSchemaMatcher_RanksByRelevance(t *testing.T) {
	m := NewSchemaMatcher(model.SchemaConfig{})

	tools := []model.ToolConfig{
		toolWithMeta("unrelated", &model.ToolMeta{Tags: []string{"billing"}}),
		toolWithMeta("user_reader", &model.ToolMeta{Tags: []string{"user"}, Capabilities: []string{"user.read"}}),
	}

	out, _ := m.Rank(tools, SelectionContext{UserPrompt: "read the user profile"})

	// user_reader matches "user" + "read" and should rank first.
	require.Len(t, out, 2)
	assert.Equal(t, "user_reader", out[0].Name)
}

func TestSchemaMatcher_TopKTruncates(t *testing.T) {
	m := NewSchemaMatcher(model.SchemaConfig{TopK: 2})

	tools := testTools("a", "b", "c", "d")
	out, traces := m.Rank(tools, SelectionContext{UserPrompt: "anything"})

	assert.Len(t, out, 2)
	assert.Len(t, traces, 2)
	for _, tr := range traces {
		assert.Equal(t, "below_top_k", tr.Detail)
	}
}

func TestSchemaMatcher_NoPromptPreservesOrder(t *testing.T) {
	m := NewSchemaMatcher(model.SchemaConfig{})
	tools := testTools("a", "b", "c")

	out, _ := m.Rank(tools, SelectionContext{})

	assert.Equal(t, []string{"a", "b", "c"}, keptNames(out))
}

func TestSchemaMatcher_NoPromptWithTopK(t *testing.T) {
	m := NewSchemaMatcher(model.SchemaConfig{TopK: 1})
	tools := testTools("a", "b", "c")

	out, _ := m.Rank(tools, SelectionContext{})

	assert.Equal(t, []string{"a"}, keptNames(out))
}

func TestSchemaMatcher_DescriptionMatch(t *testing.T) {
	m := NewSchemaMatcher(model.SchemaConfig{})

	tools := []model.ToolConfig{
		{Name: "plain", Description: "does nothing special"},
		{Name: "invoicer", Description: "create an invoice for a customer"},
	}

	out, _ := m.Rank(tools, SelectionContext{UserPrompt: "create invoice"})
	assert.Equal(t, "invoicer", out[0].Name)
}

func TestTokenize(t *testing.T) {
	tokens := tokenize("Read the User-Profile, now!")
	assert.Contains(t, tokens, "read")
	assert.Contains(t, tokens, "user")
	assert.Contains(t, tokens, "profile")
	assert.NotContains(t, tokens, "")
}

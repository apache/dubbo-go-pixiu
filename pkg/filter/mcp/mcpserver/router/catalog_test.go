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

func TestToolCatalogViewPickOrderedMissingDuplicateAndDiscoverable(t *testing.T) {
	hidden := false
	tools := []model.ToolConfig{
		{Name: "first"},
		{Name: "dup", Description: "first duplicate wins for name lookup"},
		{Name: "hidden", Meta: &model.ToolMeta{DiscoveryVisibility: &hidden}},
		{Name: "dup", Description: "second duplicate stays in ordered scans"},
	}

	view := NewToolCatalogView(tools)

	assert.Equal(t, []string{"first", "dup", "hidden", "dup"}, view.Names())
	assert.Equal(t, []string{"first", "dup", "dup"}, view.DiscoverableNames())

	picked := view.PickOrdered([]string{"hidden", "missing", "dup"})
	require.Len(t, picked, 2)
	assert.Equal(t, []string{"hidden", "dup"}, toolNames(picked))
	assert.Equal(t, "first duplicate wins for name lookup", picked[1].Description)

	filtered := view.FilterSet(NewStringSet([]string{"dup", "first"}))
	assert.Equal(t, []string{"first", "dup", "dup"}, toolNames(filtered))
}

func TestValidatedStringSetReportsEmptyAndDuplicate(t *testing.T) {
	_, err := NewValidatedStringSet("workflow \"support\": tool", []string{"a", " ", "a"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "workflow \"support\": tool at index 1 is empty")

	_, err = NewValidatedStringSet("workflow \"support\": tool", []string{"a", "a"})
	require.Error(t, err)
	assert.Contains(t, err.Error(), "workflow \"support\": tool duplicate \"a\" at index 1")
}

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

package mcpserver

import (
	"encoding/json"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestMCPToolProjectionSharedByRegistryListAndFingerprint(t *testing.T) {
	tool := model.ToolConfig{
		Name:        "search",
		Description: "Search tickets",
		Args: []model.ArgConfig{{
			Name:        "query",
			Type:        "string",
			Description: "Query text",
			Required:    true,
			Enum:        []string{"open", "closed"},
			Default:     "open",
		}},
	}

	registry := NewToolRegistry()
	require.NoError(t, registry.ReplaceAllTools([]model.ToolConfig{tool}))

	registryTools, err := registry.ToMCPTools()
	require.NoError(t, err)
	require.Len(t, registryTools, 1)

	projected, err := MCPToolMap(BuildMCPTool(tool))
	require.NoError(t, err)
	assert.Equal(t, projected, registryTools[0])

	fromList, err := json.Marshal(BuildMCPTools([]model.ToolConfig{tool})[0])
	require.NoError(t, err)
	fromRegistry, err := json.Marshal(registryTools[0])
	require.NoError(t, err)
	assert.JSONEq(t, string(fromList), string(fromRegistry))

	assert.Equal(t, visibleToolsFingerprint([]model.ToolConfig{tool}), visibleToolsFingerprint(registry.ListTools()))
}

func TestVisibleFingerprintUsesDiscoverableProjectionOnly(t *testing.T) {
	hidden := false
	visible := model.ToolConfig{Name: "visible", Args: []model.ArgConfig{{Name: "q", Type: "string"}}}
	hiddenTool := model.ToolConfig{Name: "hidden", Args: []model.ArgConfig{{Name: "q", Type: "string"}}, Meta: &model.ToolMeta{DiscoveryVisibility: &hidden}}
	changedHidden := hiddenTool
	changedHidden.Description = "changed but still hidden"

	assert.Equal(t,
		visibleToolsFingerprint([]model.ToolConfig{visible, hiddenTool}),
		visibleToolsFingerprint([]model.ToolConfig{visible, changedHidden}))
	assert.NotEqual(t,
		visibleToolsFingerprint([]model.ToolConfig{visible}),
		visibleToolsFingerprint([]model.ToolConfig{{Name: "visible", Description: "changed", Args: visible.Args}}))
}

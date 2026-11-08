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

package model

import (
	"reflect"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"

	"gopkg.in/yaml.v3"
)

func TestMcpServerConfigRouterPresence(t *testing.T) {
	var omitted McpServerConfig
	require.NoError(t, yaml.Unmarshal([]byte("endpoint: /mcp\n"), &omitted))
	assert.Nil(t, omitted.Router)

	var empty McpServerConfig
	require.NoError(t, yaml.Unmarshal([]byte("endpoint: /mcp\nrouter: {}\n"), &empty))
	assert.NotNil(t, empty.Router)

	var configured McpServerConfig
	require.NoError(t, yaml.Unmarshal([]byte("endpoint: /mcp\nrouter:\n  fallback: fail_closed\n"), &configured))
	require.NotNil(t, configured.Router)
	assert.Equal(t, "fail_closed", configured.Router.Fallback)
}

func TestRouterConfigDoesNotExposeLegacySwitches(t *testing.T) {
	typ := reflect.TypeOf(RouterConfig{})
	_, hasEnabled := typ.FieldByName("Enabled")
	_, hasEnforceOnCall := typ.FieldByName("EnforceOnCall")

	assert.False(t, hasEnabled)
	assert.False(t, hasEnforceOnCall)
}

func TestMcpServerConfigDeepCopyRouterAndToolMetaIsolation(t *testing.T) {
	policyEnabled := true
	workflowEnabled := true
	visible := false
	cfg := &McpServerConfig{
		Tools: []ToolConfig{{
			Name: "hidden",
			Meta: &ToolMeta{
				Tags:                []string{"safe"},
				DiscoveryVisibility: &visible,
			},
		}},
		Router: &RouterConfig{
			Fallback: "fail_closed",
			Stages: RouterStages{
				Policy:   &policyEnabled,
				Workflow: &workflowEnabled,
			},
			Policy: PolicyConfig{Rules: []PolicyRule{{
				Name:      "tenant",
				AllowTags: []string{"safe"},
				DenyTags:  []string{"admin"},
				When:      PolicyMatch{Claim: "tenant", In: []string{"acme"}},
			}}},
			Workflows: []WorkflowConfig{{
				Name:  "safe-minimal",
				Tools: []string{"hidden"},
				When:  PolicyMatch{Claim: "role", In: []string{"support"}},
			}},
			Progressive: ProgressiveConfig{InitialBundle: "safe-minimal", ExpandAfterCalls: 1},
			Session:     RouterSessionConfig{MaxEntries: 7},
		},
	}

	cp := cfg.DeepCopy()
	require.NotSame(t, cfg.Router, cp.Router)
	require.NotSame(t, cfg.Tools[0].Meta, cp.Tools[0].Meta)
	require.NotSame(t, cfg.Tools[0].Meta.DiscoveryVisibility, cp.Tools[0].Meta.DiscoveryVisibility)
	require.NotSame(t, cfg.Router.Stages.Policy, cp.Router.Stages.Policy)
	require.NotSame(t, cfg.Router.Stages.Workflow, cp.Router.Stages.Workflow)

	*cfg.Tools[0].Meta.DiscoveryVisibility = true
	cfg.Tools[0].Meta.Tags[0] = "mutated"
	*cfg.Router.Stages.Policy = false
	cfg.Router.Policy.Rules[0].AllowTags[0] = "mutated"
	cfg.Router.Policy.Rules[0].When.In[0] = "mutated"
	cfg.Router.Workflows[0].Tools[0] = "mutated"
	cfg.Router.Workflows[0].When.In[0] = "mutated"
	cfg.Router.Session.MaxEntries = 99

	assert.False(t, *cp.Tools[0].Meta.DiscoveryVisibility)
	assert.Equal(t, []string{"safe"}, cp.Tools[0].Meta.Tags)
	assert.True(t, *cp.Router.Stages.Policy)
	assert.Equal(t, []string{"safe"}, cp.Router.Policy.Rules[0].AllowTags)
	assert.Equal(t, []string{"acme"}, cp.Router.Policy.Rules[0].When.In)
	assert.Equal(t, []string{"hidden"}, cp.Router.Workflows[0].Tools)
	assert.Equal(t, []string{"support"}, cp.Router.Workflows[0].When.In)
	assert.Equal(t, 7, cp.Router.Session.MaxEntries)
}

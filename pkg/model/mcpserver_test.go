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

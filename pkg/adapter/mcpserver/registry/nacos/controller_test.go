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

package nacos

import (
	"fmt"
	"testing"
)

import (
	"github.com/nacos-group/nacos-sdk-go/v2/model"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func newControllerTestClient(configs map[string]any) *NacosRegistryClient {
	return &NacosRegistryClient{
		configClient: &MockedNacosConfigClient{
			configs:           configs,
			configListenerMap: map[string][]func(string, string, string, string){},
		},
		namingClient: MockedNacosNamingClient{
			listenerMap: map[string][]func(services []model.Instance, err error){},
		},
		servers: map[string]*ServerContext{},
	}
}

func versionConfigEntry(id string) (string, string) {
	return fmt.Sprintf("%s-mcp-versions.json%smcp-server-versions", id, configKeySeparator),
		createMcpServerVersionConfig(id, id, "https", "mcp-streamable", testVersionLatest)
}

func TestMcpControllerServerRemovalPublishesTombstone(t *testing.T) {
	key, value := versionConfigEntry("server-b")
	client := newControllerTestClient(map[string]any{key: value})
	var tombstones []string
	controller := NewMcpController(client, func(serverId string, cfg *McpServerConfig) {
		if cfg == nil {
			tombstones = append(tombstones, serverId)
		}
	})
	controller.watched["server-a"] = true

	require.NoError(t, controller.reconcile())

	assert.Equal(t, []string{"server-a"}, tombstones)
	assert.False(t, controller.watched["server-a"])
}

func TestMcpControllerEmptyListThresholdPublishesTombstones(t *testing.T) {
	client := newControllerTestClient(map[string]any{})
	var tombstones []string
	controller := NewMcpController(client, func(serverId string, cfg *McpServerConfig) {
		if cfg == nil {
			tombstones = append(tombstones, serverId)
		}
	})
	controller.watched["server-a"] = true
	controller.watched["server-b"] = true

	require.NoError(t, controller.reconcile())
	assert.Empty(t, tombstones)
	require.NoError(t, controller.reconcile())
	assert.Empty(t, tombstones)
	require.NoError(t, controller.reconcile())

	assert.ElementsMatch(t, []string{"server-a", "server-b"}, tombstones)
	assert.Empty(t, controller.watched)
}

func TestMcpControllerClosePublishesTombstones(t *testing.T) {
	client := newControllerTestClient(map[string]any{})
	var tombstones []string
	controller := NewMcpController(client, func(serverId string, cfg *McpServerConfig) {
		if cfg == nil {
			tombstones = append(tombstones, serverId)
		}
	})
	controller.watched["server-a"] = true
	controller.watched["server-b"] = true

	require.NoError(t, controller.Close())

	assert.ElementsMatch(t, []string{"server-a", "server-b"}, tombstones)
	assert.Empty(t, controller.watched)
	require.NoError(t, controller.Close(), "repeated close must be idempotent")
	assert.Len(t, tombstones, 2)
}

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
	"sync"

	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

var (
	globalRegistry *ToolRegistry
	globalDynamic  *DynamicConsumer

	// sync.Once variables for thread-safe singleton initialization
	registryOnce sync.Once
	dynamicOnce  sync.Once
)

// GetOrInitRegistry returns a singleton ToolRegistry
func GetOrInitRegistry() *ToolRegistry {
	registryOnce.Do(func() {
		globalRegistry = NewToolRegistry()
	})
	return globalRegistry
}

// GetOrInitDynamic returns a singleton DynamicConsumer
func GetOrInitDynamic() *DynamicConsumer {
	dynamicOnce.Do(func() {
		globalDynamic = NewDynamicConsumer(GetOrInitRegistry())
	})
	return globalDynamic
}

// DynamicConsumer applies dynamic MCP configurations into the registry
type DynamicConsumer struct {
	registry *ToolRegistry
}

func NewDynamicConsumer(reg *ToolRegistry) *DynamicConsumer {
	return &DynamicConsumer{registry: reg}
}

// update tools from the remote config in nacos
func (d *DynamicConsumer) ApplyMcpServerConfig(cfg *model.McpServerConfig) error {
	if cfg == nil {
		return nil
	}

	// full sync tools
	d.registry.ReplaceAllTools(cfg.Tools)
	logger.Infof("[MCP Dynamic] applied config: tools synced=%d", len(cfg.Tools))
	return nil
}

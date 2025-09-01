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
	"os"
	"sync"

	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/adapter"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// TODO: Implement mcpserver/registry package
// "github.com/apache/dubbo-go-pixiu/pkg/adapter/mcpserver/registry"

func init() {
	adapter.RegisterAdapterPlugin(&Plugin{})
}

var (
	_ adapter.AdapterPlugin = new(Plugin)
	_ adapter.Adapter       = new(Adapter)
)

type (
	// Plugin to monitor mcp services on registry center
	Plugin struct{}

	// AdapterConfig holds configuration for multiple registries
	AdapterConfig struct {
		Registries map[string]model.Registry `yaml:"registries" json:"registries" mapstructure:"registries"`
	}

	// Adapter to monitor mcp services on registry center
	Adapter struct {
		id  string
		cfg *AdapterConfig
		// TODO: Use actual registry interface when implemented
		// registries map[string]registry.Registry
		mu sync.RWMutex
	}

	// McpServerInfo represents an MCP server instance from service discovery
	McpServerInfo struct {
		ServerID string            `json:"server_id"`
		Endpoint string            `json:"endpoint"`
		Protocol string            `json:"protocol"`
		Metadata map[string]string `json:"metadata"`
	}
)

// Kind returns the identifier of the plugin
func (p Plugin) Kind() string {
	return constant.McpServerAdapter
}

// CreateAdapter returns the mcp server adapter
func (p *Plugin) CreateAdapter(a *model.Adapter) (adapter.Adapter, error) {
	return &Adapter{
		id:  a.ID,
		cfg: &AdapterConfig{Registries: make(map[string]model.Registry)},
		// TODO: Initialize registries when implemented
	}, nil
}

// Start starts the adapter
func (a *Adapter) Start() {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// TODO: Start registries when implemented
	// for _, reg := range a.registries {
	//     if err := reg.Subscribe(); err != nil {
	//         logger.Errorf("MCP registry %s subscribe failed: %v", reg, err)
	//     }
	// }
	logger.Infof("MCP server adapter %s started successfully", a.id)
}

// Stop stops the adapter
func (a *Adapter) Stop() {
	a.mu.RLock()
	defer a.mu.RUnlock()

	// TODO: Stop registries when implemented
	// for _, reg := range a.registries {
	//     if err := reg.Unsubscribe(); err != nil {
	//         logger.Errorf("MCP registry %s unsubscribe failed: %v", reg, err)
	//     }
	// }
	logger.Infof("MCP server adapter %s stopped successfully", a.id)
}

// Apply inits the registries according to the configuration
func (a *Adapter) Apply() error {
	a.mu.Lock()
	defer a.mu.Unlock()

	// Support environment variable override for Nacos address
	nacosAddrFromEnv := os.Getenv(constant.EnvDubbogoPixiuNacosRegistryAddress)

	for k, registryConfig := range a.cfg.Registries {
		if nacosAddrFromEnv != "" && registryConfig.Protocol == constant.Nacos {
			registryConfig.Address = nacosAddrFromEnv
		}

		// TODO: Create registry when implemented
		// reg, err := registry.GetRegistry(k, registryConfig, a)
		// if err != nil {
		//     logger.Errorf("Create MCP registry %s failed: %v", k, err)
		//     return err
		// }
		// a.registries[k] = reg

		logger.Infof("MCP registry %s configured successfully", k)
	}

	return nil
}

// Config returns the config of the adapter
func (a *Adapter) Config() any {
	return a.cfg
}

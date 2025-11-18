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
	"context"
	"os"
	"strconv"
	"sync"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/mcpserver/common/util"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/mcpserver/registry"
	_ "github.com/apache/dubbo-go-pixiu/pkg/adapter/mcpserver/registry/nacos"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/adapter"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

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
		// single provider controller (provider-agnostic)
		controller registry.Controller
		ctx        context.Context
		cancel     context.CancelFunc
		mu         sync.RWMutex
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
func (p *Plugin) Kind() string {
	return constant.McpServerAdapter
}

// CreateAdapter returns the mcp server adapter
func (p *Plugin) CreateAdapter(a *model.Adapter) (adapter.Adapter, error) {
	return &Adapter{
		id:  a.ID,
		cfg: &AdapterConfig{Registries: make(map[string]model.Registry)},
	}, nil
}

// Start starts the adapter
func (a *Adapter) Start() {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.controller == nil {
		logger.Warnf("MCP server adapter %s start skipped: controller not initialized (call Apply first)", a.id)
		return
	}

	if a.cancel != nil {
		logger.Infof("MCP server adapter %s already running", a.id)
		return
	}

	a.ctx, a.cancel = context.WithCancel(context.Background())
	go func() {
		if err := a.controller.Run(a.ctx, 30*time.Second); err != nil {
			logger.Errorf("MCP server controller run error: %v", err)
		}
	}()

	logger.Infof("MCP server adapter %s started successfully", a.id)
}

// Stop stops the adapter
func (a *Adapter) Stop() {
	a.mu.RLock()
	defer a.mu.RUnlock()

	if a.cancel != nil {
		a.cancel()
		a.cancel = nil
	}

	if a.controller != nil {
		if err := a.controller.Close(); err != nil {
			logger.Errorf("MCP server controller close error: %v", err)
		}
	}
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
			// Validate environment variable address before overriding
			if err := util.ValidateNacosAddresses(nacosAddrFromEnv); err != nil {
				logger.Errorf("[dubbo-go-pixiu] mcp adapter invalid NACOS_ADDRESS environment variable: %v, keeping original config", err)
				// Continue with original configuration instead of failing
			} else {
				logger.Infof("[dubbo-go-pixiu] mcp adapter overriding nacos address with environment variable: %s", nacosAddrFromEnv)
				registryConfig.Address = nacosAddrFromEnv
			}
		}

		// only handle nacos for now
		if registryConfig.Protocol != constant.Nacos {
			logger.Infof("MCP registry %s skipped (protocol=%s)", k, registryConfig.Protocol)
			continue
		}

		onChange := func(serverId string, cfg *model.McpServerConfig) {
			if cfg == nil {
				return
			}

			if serverId == "" {
				serverId = "default"
			}

			// 1) apply tools dynamically to registry for filter usage
			if dc := mcpserver.GetOrInitDynamicConsumer(); dc != nil {
				if err := dc.ApplyMcpServerConfigByServer(serverId, cfg); err != nil {
					logger.Errorf("[dubbo-go-pixiu] mcp adapter apply server %s config error: %v", serverId, err)
				}
			} else {
				logger.Infof("[dubbo-go-pixiu] mcp adapter update received from server %s: tools=%d", serverId, len(cfg.Tools))
			}
			// 2) register endpoint for each tool using BackendURL (host:port) into cluster named by tool.Name
			for _, tool := range cfg.Tools {
				if tool.BackendURL == "" {
					continue
				}
				result, err := util.ParseHostPortFromURL(tool.BackendURL)
				if err != nil {
					logger.Errorf("[dubbo-go-pixiu] mcp adapter failed to parse BackendURL '%s' for tool '%s': %v",
						tool.BackendURL, tool.Name, err)
					continue
				}
				if result.UsedFallback {
					logger.Warnf("[dubbo-go-pixiu] mcp adapter using fallback for tool '%s' with BackendURL '%s': %s",
						tool.Name, tool.BackendURL, result.FallbackInfo)
				}
				endpointID := result.Host + ":" + strconv.Itoa(result.Port)
				server.GetClusterManager().SetEndpoint(tool.Cluster, &model.Endpoint{
					ID: endpointID,
					Address: model.SocketAddress{
						Address: result.Host,
						Port:    result.Port,
					},
				})
			}
		}

		// build controller via provider-agnostic factory
		ctrl, err := registry.BuildController(registryConfig, onChange)
		if err != nil {
			return err
		}
		a.controller = ctrl
		logger.Infof("MCP registry %s configured successfully (nacos)", k)
	}

	return nil
}

// Config returns the config of the adapter
func (a *Adapter) Config() any {
	return a.cfg
}

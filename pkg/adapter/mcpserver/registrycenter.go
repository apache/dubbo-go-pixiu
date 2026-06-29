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
)

func init() {
	adapter.RegisterAdapterPlugin(&Plugin{})
}

var (
	_ adapter.AdapterPlugin = new(Plugin)
	_ adapter.Adapter       = new(Adapter)

	serverPublicationSinkForSingleRuntime = mcpserver.ServerPublicationSinkForSingleRuntime
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
		id               string
		cfg              *AdapterConfig
		controllers      map[string]registry.Controller
		ctx              context.Context
		cancel           context.CancelFunc
		endpoints        *endpointReconciler
		sink             mcpserver.ServerPublicationSink
		publishedSources map[string]mcpserver.ServerSource
		mu               sync.RWMutex
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
		id:               a.ID,
		cfg:              &AdapterConfig{Registries: make(map[string]model.Registry)},
		controllers:      make(map[string]registry.Controller),
		endpoints:        newEndpointReconciler(clusterManagerEndpointSink{}),
		publishedSources: make(map[string]mcpserver.ServerSource),
	}, nil
}

// Start starts the adapter
func (a *Adapter) Start() {
	a.mu.Lock()
	if len(a.controllers) == 0 {
		a.mu.Unlock()
		logger.Warnf("MCP server adapter %s start skipped: controller not initialized (call Apply first)", a.id)
		return
	}

	if a.cancel != nil {
		a.mu.Unlock()
		logger.Infof("MCP server adapter %s already running", a.id)
		return
	}

	a.ctx, a.cancel = context.WithCancel(context.Background())
	ctx := a.ctx
	controllers := cloneControllers(a.controllers)
	a.mu.Unlock()

	for registryName, ctrl := range controllers {
		registryName, ctrl := registryName, ctrl
		go func() {
			if err := ctrl.Run(ctx, 30*time.Second); err != nil {
				logger.Errorf("MCP server controller %s run error: %v", registryName, err)
			}
		}()
	}

	logger.Infof("MCP server adapter %s started successfully with %d controller(s)", a.id, len(controllers))
}

// Stop stops the adapter
func (a *Adapter) Stop() {
	a.mu.Lock()
	cancel := a.cancel
	a.cancel = nil
	controllers := cloneControllers(a.controllers)
	a.ctx = nil
	a.mu.Unlock()

	if cancel != nil {
		cancel()
	}

	for registryName, ctrl := range controllers {
		if err := ctrl.Close(); err != nil {
			logger.Errorf("MCP server controller %s close error: %v", registryName, err)
		}
	}

	a.removeAllDynamicPublication()
	logger.Infof("MCP server adapter %s stopped successfully", a.id)
}

func cloneControllers(in map[string]registry.Controller) map[string]registry.Controller {
	out := make(map[string]registry.Controller, len(in))
	for k, v := range in {
		out[k] = v
	}
	return out
}

func (a *Adapter) runController(registryName string, ctrl registry.Controller, ctx context.Context) {
	go func() {
		if err := ctrl.Run(ctx, 30*time.Second); err != nil {
			logger.Errorf("MCP server controller %s run error: %v", registryName, err)
		}
	}()
}

// Apply inits the registries according to the configuration
func (a *Adapter) Apply() error {
	a.mu.Lock()
	if a.endpoints == nil {
		a.endpoints = newEndpointReconciler(clusterManagerEndpointSink{})
	}
	registries := make(map[string]model.Registry, len(a.cfg.Registries))
	for k, v := range a.cfg.Registries {
		registries[k] = v
	}
	running := a.cancel != nil
	a.mu.Unlock()

	// Support environment variable override for Nacos address
	nacosAddrFromEnv := os.Getenv(constant.EnvDubbogoPixiuNacosRegistryAddress)

	newControllers := make(map[string]registry.Controller)
	for k, registryConfig := range registries {
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

		registryName := k
		onChange := func(serverId string, cfg *model.McpServerConfig) {
			a.applyServerConfigEvent(registryName, serverId, cfg)
		}

		// build controller via provider-agnostic factory
		ctrl, err := registry.BuildController(registryConfig, onChange)
		if err != nil {
			closeControllers(newControllers)
			return err
		}
		newControllers[registryName] = ctrl
		logger.Infof("MCP registry %s configured successfully (nacos)", k)
	}

	a.mu.Lock()
	oldCancel := a.cancel
	oldControllers := cloneControllers(a.controllers)
	if oldCancel != nil {
		oldCancel()
	}
	a.controllers = newControllers
	a.cancel = nil
	a.ctx = nil
	if running && len(newControllers) > 0 {
		a.ctx, a.cancel = context.WithCancel(context.Background())
	}
	ctx := a.ctx
	a.mu.Unlock()

	closeControllers(oldControllers)
	a.removeAllDynamicPublication()

	if running && ctx != nil {
		for registryName, ctrl := range newControllers {
			a.runController(registryName, ctrl, ctx)
		}
	}

	return nil
}

func closeControllers(controllers map[string]registry.Controller) {
	for registryName, ctrl := range controllers {
		if err := ctrl.Close(); err != nil {
			logger.Errorf("MCP server controller %s close error: %v", registryName, err)
		}
	}
}

func (a *Adapter) applyServerConfigEvent(registryName, serverId string, cfg *model.McpServerConfig) {
	source := mcpserver.NewServerSource(registryName, serverId)
	reconciler := a.endpointReconciler()

	// Apply catalog and endpoints through one desired-state publication path. If
	// no runtime target is bound, skip the entire update so authorization catalog
	// and cluster endpoints cannot diverge.
	sink, err := a.bindPublicationSink()
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp adapter cannot bind runtime publication sink for source %s: %v", source, err)
		return
	}
	if sink == nil {
		logger.Errorf("[dubbo-go-pixiu] mcp adapter update received from source %s without a bound runtime publication sink", source)
		return
	}
	if sink.RuntimeID() == "" {
		logger.Errorf("[dubbo-go-pixiu] mcp adapter update received from source %s but publication sink has no runtime id", source)
		return
	}

	if cfg == nil {
		if err := sink.ApplyMcpServerConfigBySource(source, nil); err != nil {
			logger.Errorf("[dubbo-go-pixiu] mcp adapter remove source %s config error: %v", source, err)
			return
		}
		reconciler.RemoveSource(source)
		a.untrackPublishedSource(source)
		return
	}

	if err := reconciler.ValidateServerConfig(source, cfg); err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp adapter validate endpoints for source %s error: %v", source, err)
		return
	}
	if err := sink.ApplyMcpServerConfigBySource(source, cfg); err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp adapter apply source %s config error: %v", source, err)
		return
	}
	if err := reconciler.ApplyServerConfig(source, cfg); err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp adapter reconcile endpoints for source %s error: %v", source, err)
		return
	}
	a.trackPublishedSource(source)
}

func (a *Adapter) endpointReconciler() *endpointReconciler {
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.endpoints == nil {
		a.endpoints = newEndpointReconciler(clusterManagerEndpointSink{})
	}
	return a.endpoints
}

func (a *Adapter) bindPublicationSink() (mcpserver.ServerPublicationSink, error) {
	sink, err := serverPublicationSinkForSingleRuntime()
	a.mu.Lock()
	defer a.mu.Unlock()
	if err != nil {
		a.sink = nil
		return nil, err
	}
	if sink.RuntimeID() == "" {
		a.sink = nil
		return nil, mcpserver.ErrDynamicConsumerUnavailable
	}
	a.sink = sink
	return sink, nil
}

func (a *Adapter) trackPublishedSource(source mcpserver.ServerSource) {
	source = source.Normalize()
	a.mu.Lock()
	defer a.mu.Unlock()
	if a.publishedSources == nil {
		a.publishedSources = make(map[string]mcpserver.ServerSource)
	}
	a.publishedSources[source.Key()] = source
}

func (a *Adapter) untrackPublishedSource(source mcpserver.ServerSource) {
	source = source.Normalize()
	a.mu.Lock()
	defer a.mu.Unlock()
	delete(a.publishedSources, source.Key())
}

func (a *Adapter) removeAllDynamicPublication() {
	sink, err := a.bindPublicationSink()
	if err != nil {
		logger.Infof("[dubbo-go-pixiu] mcp adapter dynamic catalog cleanup skipped: %v", err)
	} else if err := sink.RemoveAllMcpServerConfigs(); err != nil {
		logger.Errorf("[dubbo-go-pixiu] mcp adapter dynamic catalog cleanup error: %v", err)
	}

	reconciler := a.endpointReconciler()
	reconciler.RemoveAll()

	a.mu.Lock()
	a.publishedSources = make(map[string]mcpserver.ServerSource)
	a.mu.Unlock()
}

// Config returns the config of the adapter
func (a *Adapter) Config() any {
	return a.cfg
}

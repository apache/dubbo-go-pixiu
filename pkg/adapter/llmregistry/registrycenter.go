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

package llmregistry

import (
	"os"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/llmregistry/common"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/llmregistry/registry"
	_ "github.com/apache/dubbo-go-pixiu/pkg/adapter/llmregistry/registry/nacos"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/adapter"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

const (
	Kind = constant.LLMRegistryCenterAdapter
)

func init() {
	adapter.RegisterAdapterPlugin(&Plugin{})
}

var (
	_ adapter.AdapterPlugin = new(Plugin)
	_ adapter.Adapter       = new(Adapter)
	// Explicitly declare that Adapter implements the listener interface.
	_ common.RegistryEventListener = new(Adapter)
)

type (
	// Plugin to monitor LLM services on a registry center.
	Plugin struct{}

	// AdapterConfig defines the configuration for the LLM registry adapter.
	AdapterConfig struct {
		Registries map[string]model.Registry `yaml:"registries" json:"registries" mapstructure:"registries"`
	}
)

// Kind returns the identifier of the plugin.
func (p *Plugin) Kind() string {
	return Kind
}

// CreateAdapter returns the LLM registry center adapter.
func (p *Plugin) CreateAdapter(a *model.Adapter) (adapter.Adapter, error) {
	adapter := &Adapter{
		id:         a.ID,
		registries: make(map[string]registry.Registry),
		cfg:        &AdapterConfig{Registries: make(map[string]model.Registry)},
	}
	return adapter, nil
}

// Adapter to monitor LLM services on a registry center.
type Adapter struct {
	id         string
	cfg        *AdapterConfig
	registries map[string]registry.Registry
}

// Start starts the adapter by subscribing to all configured registries.
func (a *Adapter) Start() {
	for name, reg := range a.registries {
		logger.Infof("Subscribing to LLM registry: %s", name)
		if err := reg.Subscribe(); err != nil {
			logger.Errorf("Failed to subscribe to registry %s: %s", name, err.Error())
		}
	}
}

// Stop stops the adapter by unsubscribing from all registries.
func (a *Adapter) Stop() {
	for name, reg := range a.registries {
		if err := reg.Unsubscribe(); err != nil {
			logger.Errorf("Failed to unsubscribe from registry %s: %s", name, err.Error())
		}
	}
}

// Apply initializes the registries according to the configuration.
func (a *Adapter) Apply() error {
	nacosAddrFromEnv := os.Getenv(constant.EnvDubbogoPixiuNacosRegistryAddress)

	for key, registryConfig := range a.cfg.Registries {
		var err error
		// Override address from environment variable if it's set for Nacos.
		if nacosAddrFromEnv != "" && registryConfig.Protocol == constant.Nacos {
			logger.Infof("Overriding Nacos address for registry '%s' with environment variable: %s", key, nacosAddrFromEnv)
			registryConfig.Address = nacosAddrFromEnv
		}

		a.registries[key], err = registry.GetRegistry(registryConfig, a)
		if err != nil {
			return err
		}
	}

	return nil
}

// Config returns the configuration of the adapter.
func (a *Adapter) Config() any {
	return a.cfg
}

// OnAddEndpoint is the callback that gets triggered when a new LLM endpoint is discovered.
func (a *Adapter) OnAddEndpoint(endpoint *model.Endpoint) error {
	// The endpoint metadata MUST contain a "cluster" key to identify the target cluster.
	clusterName, ok := endpoint.Metadata["cluster"]
	if !ok || clusterName == "" {
		logger.Warnf("Endpoint %s (ID: %s) is missing 'cluster' metadata, skipping.", endpoint.Name, endpoint.ID)
		return nil
	}
	logger.Infof("Adding endpoint %s to cluster %s", endpoint.ID+constant.PathParamIdentifier+endpoint.Name, clusterName)
	server.GetClusterManager().SetEndpoint(clusterName, endpoint)
	return nil
}

// OnRemoveEndpoint is the callback that gets triggered when an LLM endpoint is removed.
func (a *Adapter) OnRemoveEndpoint(endpoint *model.Endpoint) error {
	// The endpoint metadata MUST contain a "cluster" key.
	clusterName, ok := endpoint.Metadata["cluster"]
	if !ok || clusterName == "" {
		logger.Warnf("Endpoint %s (ID: %s) is missing 'cluster' metadata for removal, skipping.", endpoint.Name, endpoint.ID)
		return nil
	}
	logger.Infof("Removing endpoint ID %s from cluster %s", endpoint.ID, clusterName)
	server.GetClusterManager().DeleteEndpoint(clusterName, endpoint.ID)
	return nil
}

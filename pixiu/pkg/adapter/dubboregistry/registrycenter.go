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

package dubboregistry

import (
	"os"
)

import (
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/router"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/registry"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/registry/nacos"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/registry/zookeeper"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/adapter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server"
)

func init() {
	adapter.RegisterAdapterPlugin(&Plugin{})
}

var (
	_ adapter.AdapterPlugin = new(Plugin)
	_ adapter.Adapter       = new(Adapter)
)

type (
	// Plugin to monitor dubbo services on registry center
	Plugin struct {
	}

	AdaptorConfig struct {
		Registries map[string]model.Registry `yaml:"registries" json:"registries" mapstructure:"registries"`
	}
)

// Kind returns the identifier of the plugin
func (p Plugin) Kind() string {
	return constant.DubboRegistryCenterAdapter
}

// CreateAdapter returns the dubbo registry center adapter
func (p *Plugin) CreateAdapter(a *model.Adapter) (adapter.Adapter, error) {
	adapter := &Adapter{id: a.ID,
		registries: make(map[string]registry.Registry),
		cfg:        AdaptorConfig{Registries: make(map[string]model.Registry)}}
	return adapter, nil
}

// Adapter to monitor dubbo services on registry center
type Adapter struct {
	id         string
	cfg        AdaptorConfig
	registries map[string]registry.Registry
}

// Start starts the adaptor
func (a Adapter) Start() {
	for _, reg := range a.registries {
		if err := reg.Subscribe(); err != nil {
			logger.Errorf("Subscribe fail, error is {%s}", err.Error())
		}
	}
}

// Stop stops the adaptor
func (a *Adapter) Stop() {
	for _, reg := range a.registries {
		if err := reg.Unsubscribe(); err != nil {
			logger.Errorf("Unsubscribe fail, error is {%s}", err.Error())
		}
	}
}

// Apply inits the registries according to the configuration
func (a *Adapter) Apply() error {
	// create registry per config
	nacosAddrFromEnv := os.Getenv(constant.EnvDubbogoPixiuNacosRegistryAddress)
	for k, registryConfig := range a.cfg.Registries {
		var err error
		if nacosAddrFromEnv != "" && registryConfig.Protocol == constant.Nacos {
			registryConfig.Address = nacosAddrFromEnv
		}
		a.registries[k], err = registry.GetRegistry(k, registryConfig, a)
		if err != nil {
			return err
		}
	}

	return nil
}

// Config returns the config of the adaptor
func (a Adapter) Config() interface{} {
	return a.cfg
}

func (a *Adapter) OnAddAPI(r router.API) error {
	acm := server.GetApiConfigManager()
	return acm.AddAPI(a.id, r)
}

func (a *Adapter) OnRemoveAPI(r router.API) error {
	acm := server.GetApiConfigManager()
	return acm.RemoveAPI(a.id, r)
}

func (a *Adapter) OnDeleteRouter(r config.Resource) error {
	acm := server.GetApiConfigManager()
	return acm.DeleteRouter(a.id, r)
}

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

package sidecarprovider

import (

)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud/servicediscovery"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/adapter"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"sync"
	"time"
)

func init() {
	adapter.RegisterAdapterPlugin(&Plugin{})
}

type (
	// Plugin to monitor dubbo services on registry center
	Plugin struct {
	}

	// SidecarAdapter the adapter for spring cloud
	SidecarAdapter struct {
		cfg            *AdaptorConfig
		sd             servicediscovery.ServiceDiscovery
		currentService map[string]*Service
		mutex    sync.Mutex
		stopChan chan struct{}
	}

	// AdaptorConfig the config for SidecarAdapter
	AdaptorConfig struct {
		// Mode default 0 start backgroup fresh and watch, 1 only fresh 2 only watch
		Mode          int                 `yaml:"mode" json:"mode" default:"mode"`
		Registry      *model.RemoteConfig `yaml:"registry" json:"registry" default:"registry"`
		FreshInterval time.Duration       `yaml:"freshInterval" json:"freshInterval" default:"freshInterval"`
		Services      []string            `yaml:"services" json:"services" default:"services"`
		// SubscribePolicy subscribe config,
		// - adapting : if there is no any Services (App) names, fetch All services from registry center
		// - definitely : fetch services by the config Services (App) names
		SubscribePolicy string `yaml:"subscribe-policy" json:"subscribe-policy" default:"adapting"`
	}

	Service struct {
		Name string
	}

	SubscribePolicy int
)

// Kind returns the identifier of the plugin
func (p Plugin) Kind() string {
	return constant.SidecarProviderAdapter
}

// CreateAdapter returns the dubbo registry center adapter
func (p *Plugin) CreateAdapter(a *model.Adapter) (adapter.Adapter, error) {
	adapter := &SidecarAdapter{
		cfg:      &AdaptorConfig{},
		stopChan: make(chan struct{}),
		currentService: make(map[string]*Service),
	}
	return adapter,nil
}

// Start starts the adaptor
func (a SidecarAdapter) Start() {

	registerServiceInstance()
}

// Stop stops the adaptor
func (a *SidecarAdapter) Stop() {

}

// Apply init
func (a *SidecarAdapter) Apply() error {
	// create config

	return nil
}

// Config returns the config of the adaptor
func (a SidecarAdapter) Config() interface{} {
	return a.cfg
}

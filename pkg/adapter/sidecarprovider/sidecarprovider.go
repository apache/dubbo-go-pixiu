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
	"reflect"
	"strconv"
	"sync"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	dgconstant "dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	dg "dubbo.apache.org/dubbo-go/v3/config"
	"dubbo.apache.org/dubbo-go/v3/registry"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud/servicediscovery"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/adapter"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
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

func registerServiceInstance() {
	url := selectMetadataServiceExportedURL()
	if url == nil {
		return
	}
	instance, err := createInstance(url)
	if err != nil {
		panic(err)
	}
	p := extension.GetProtocol(dgconstant.RegistryKey)
	var rp registry.RegistryFactory
	var ok bool
	if rp, ok = p.(registry.RegistryFactory); !ok {
		panic("dubbo registry protocol{" + reflect.TypeOf(p).String() + "} is invalid")
	}
	rs := rp.GetRegistries()
	for _, r := range rs {
		var sdr registry.ServiceDiscoveryHolder
		if sdr, ok = r.(registry.ServiceDiscoveryHolder); !ok {
			continue
		}
		// publish app level data to registry
		err := sdr.GetServiceDiscovery().Register(instance)
		if err != nil {
			panic(err)
		}
	}
	// publish metadata to remote
	if dg.GetApplicationConfig().MetadataType == dgconstant.RemoteMetadataStorageType {
		if remoteMetadataService, err := extension.GetRemoteMetadataService(); err == nil && remoteMetadataService != nil {
			remoteMetadataService.PublishMetadata(dg.GetApplicationConfig().Name)
		}
	}
}

func createInstance(url *common.URL) (registry.ServiceInstance, error) {
	appConfig := dg.GetApplicationConfig()
	port, err := strconv.ParseInt(url.Port, 10, 32)
	if err != nil {
		return nil, perrors.WithMessage(err, "invalid port: "+url.Port)
	}

	host := url.Ip
	if len(host) == 0 {
		host = common.GetLocalIp()
	}

	// usually we will add more metadata
	metadata := make(map[string]string, 8)
	metadata[dgconstant.MetadataStorageTypePropertyName] = appConfig.MetadataType

	instance := &registry.DefaultServiceInstance{
		ServiceName: appConfig.Name,
		Host:        host,
		Port:        int(port),
		ID:          host + dgconstant.KeySeparator + url.Port,
		Enable:      true,
		Healthy:     true,
		Metadata:    metadata,
	}

	for _, cus := range extension.GetCustomizers() {
		cus.Customize(instance)
	}

	return instance, nil
}

// selectMetadataServiceExportedURL get already be exported url
func selectMetadataServiceExportedURL() *common.URL {
	var selectedUrl *common.URL
	metaDataService, err := extension.GetLocalMetadataService(dgconstant.DefaultKey)
	if err != nil {
		logger.Warnf("get metadata service exporter failed, pls check if you import _ \"dubbo.apache.org/dubbo-go/v3/metadata/service/local\"")
		return nil
	}
	urlList, err := metaDataService.GetExportedURLs(dgconstant.AnyValue, dgconstant.AnyValue, dgconstant.AnyValue, dgconstant.AnyValue)
	if err != nil {
		panic(err)
	}
	if len(urlList) == 0 {
		return nil
	}
	for _, url := range urlList {
		selectedUrl = url
		// rest first
		if url.Protocol == "rest" {
			break
		}
	}
	return selectedUrl
}
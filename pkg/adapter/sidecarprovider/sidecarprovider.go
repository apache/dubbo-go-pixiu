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
	"os"
	"reflect"
	"strconv"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	dgconstant "dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/common/extension"
	dg "dubbo.apache.org/dubbo-go/v3/config"
	dgregistry "dubbo.apache.org/dubbo-go/v3/registry"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/adapter"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/router"
)

const Kind = constant.SidecarProviderAdapter

func init() {
	adapter.RegisterAdapterPlugin(&Plugin{})
}

type (
	// Plugin the adapter plugin for sidecar provider
	Plugin struct {
	}

	// SidecarAdapter the adapter for sidecar register
	SidecarAdapter struct {
		id		  string
		cfg       *AdaptorConfig
		registries map[string]registry.Registry
		URL		  *common.URL
	}

	// AdaptorConfig the config for SidecarAdapter
	AdaptorConfig struct {
		Registries    map[string]model.Registry `yaml:"registry" json:"registry" default:"registry"`
		Interface	  string              `yaml:"interface" json:"interface" default:"interface"`
		Protocol	  *ProtocolConfig	  `yaml:"protocol" json:"protocol" default:"protocol"`
	}

	// ProtocolConfig is protocol configuration
	ProtocolConfig struct {
		Name   string      `default:"dubbo" validate:"required" yaml:"name" json:"name,omitempty" property:"name"`
		Ip     string      `yaml:"ip"  json:"ip,omitempty" property:"ip"`
		Port   string      `default:"20000" yaml:"port" json:"port,omitempty" property:"port"`
		Params interface{} `yaml:"params" json:"params,omitempty" property:"params"`
	}

)

// Kind returns the identifier of the plugin
func (p Plugin) Kind() string {
	return Kind
}

// CreateAdapter returns the dubbo registry center adapter
func (p *Plugin) CreateAdapter(a *model.Adapter) (adapter.Adapter, error) {
	adapter := &SidecarAdapter{
		cfg:      &AdaptorConfig{},
	}
	return adapter,nil
}

// Start starts the adaptor
func (a SidecarAdapter) Start() {
	//register sidecar
	for _, reg := range a.registries {
		if err := reg.Subscribe(); err != nil {
			logger.Errorf("Subscribe fail, error is {%s}", err.Error())
		}
	}
	registerServiceInstance(a.URL)
}

// Stop stops the adaptor
func (a *SidecarAdapter) Stop() {
	for _, reg := range a.registries {
		if err := reg.Unsubscribe(); err != nil {
			logger.Errorf("Unsubscribe fail, error is {%s}", err.Error())
		}
	}
}

// Apply init
func (a *SidecarAdapter) Apply() error {
	// create config
	a.URL = common.NewURLWithOptions(
		common.WithPath(a.cfg.Interface),
		common.WithProtocol(a.cfg.Protocol.Name),
		common.WithIp(a.cfg.Protocol.Ip),
		common.WithPort(a.cfg.Protocol.Port),
	)
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
func (a SidecarAdapter) Config() interface{} {
	return a.cfg
}

func (a *SidecarAdapter) OnAddAPI(r router.API) error {
	acm := server.GetApiConfigManager()
	return acm.AddAPI(a.id, r)
}

func (a *SidecarAdapter) OnRemoveAPI(r router.API) error {
	acm := server.GetApiConfigManager()
	return acm.RemoveAPI(a.id, r)
}

func (a *SidecarAdapter) OnDeleteRouter(r config.Resource) error {
	acm := server.GetApiConfigManager()
	return acm.DeleteRouter(a.id, r)
}

func registerServiceInstance(url *common.URL) {
	if url == nil {
		return
	}
	instance, err := createInstance(url)
	if err != nil {
		panic(err)
	}
	p := extension.GetProtocol(dgconstant.RegistryKey)
	var rp dgregistry.RegistryFactory
	var ok bool
	if rp, ok = p.(dgregistry.RegistryFactory); !ok {
		panic("dubbo registry protocol{" + reflect.TypeOf(p).String() + "} is invalid")
	}
	rs := rp.GetRegistries()
	for _, r := range rs {
		var sdr dgregistry.ServiceDiscoveryHolder
		if sdr, ok = r.(dgregistry.ServiceDiscoveryHolder); !ok {
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

func createInstance(url *common.URL) (dgregistry.ServiceInstance, error) {
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

	instance := &dgregistry.DefaultServiceInstance{
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
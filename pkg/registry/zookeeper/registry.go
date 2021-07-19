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

package zookeeper

import (
	"encoding/json"
	"fmt"
	"strings"
	"time"
)
import (
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	"github.com/apache/dubbo-go/common"
	ex "github.com/apache/dubbo-go/common/extension"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	"github.com/apache/dubbo-go/metadata/definition"
	_ "github.com/apache/dubbo-go/metadata/service/inmemory"
	_ "github.com/apache/dubbo-go/metadata/service/remote"
	dr "github.com/apache/dubbo-go/registry"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
	"github.com/apache/dubbo-go/remoting/zookeeper/curator_discovery"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/registry"
	"github.com/apache/dubbo-go-pixiu/pkg/registry/base"
	zk "github.com/apache/dubbo-go-pixiu/pkg/remoting/zookeeper"
)

const (
	// RegistryZkClient zk client name
	RegistryZkClient    = "zk registry"
	rootPath            = "/dubbo"
	defaultServicesPath = "/services"
	methodsRootPath     = "/dubbo/metadata"
)

func init() {
	extension.SetRegistry(constant.Zookeeper, newZKRegistry)
}

type ZKRegistry struct {
	*baseregistry.BaseRegistry
	client       *zk.ZooKeeperClient
	servicesPath string
}

func newZKRegistry(regConfig model.Registry) (registry.Registry, error) {
	var zkReg *ZKRegistry = &ZKRegistry{}
	baseReg := baseregistry.NewBaseRegistry(zkReg)
	timeout, err := time.ParseDuration(regConfig.Timeout)
	if err != nil {
		return nil, errors.Errorf("Incorrect timeout configuration: %s", regConfig.Timeout)
	}
	client, eventChan, err := zk.NewZooKeeperClient(RegistryZkClient, strings.Split(regConfig.Address, ","), timeout)
	if err != nil {
		return nil, errors.Errorf("Initialize zookeeper client failed: %s", err.Error())
	}
	client.RegisterHandler(eventChan)
	if regConfig.ServicesPath == "" {
		regConfig.ServicesPath = defaultServicesPath
	}
	return &ZKRegistry{
		BaseRegistry: baseReg,
		client:       client,
		servicesPath: regConfig.ServicesPath,
	}, nil
}

func (r *ZKRegistry) GetClient() *zk.ZooKeeperClient {
	return r.client
}

// LoadServices loads all the registered Dubbo services from registry
func (r *ZKRegistry) LoadServices() {
	r.LoadInterfaces()
	r.LoadApplications()
}

// LoadInterfaces load services registered before dubbo 2.7
func (r *ZKRegistry) LoadInterfaces() ([]router.API, []error) {
	subPaths, err := r.client.GetChildren(rootPath)
	if err != nil {
		return nil, []error{err}
	}
	if len(subPaths) == 0 {
		return nil, nil
	}
	errorStack := []error{}
	apis := []router.API{}
	for i := range subPaths {
		if subPaths[i] == "metadata" {
			continue
		}
		providerPath := strings.Join([]string{rootPath, subPaths[i], "providers"}, constant.PathSlash)
		providerString, err := r.client.GetChildren(providerPath)
		if err != nil {
			logger.Warnf("Get provider %s failed due to %s", providerPath, err.Error())
			errorStack = append(errorStack, errors.WithStack(err))
			continue
		}
		interfaceDetailString := providerString[0]
		bkConfig, methods, err := registry.ParseDubboString(interfaceDetailString)
		if len(bkConfig.ApplicationName) == 0 || len(bkConfig.Interface) == 0 {
			errorStack = append(errorStack, errors.Errorf("Path %s contains dubbo registration that interface or application not set", providerPath))
			continue
		}
		if err != nil {
			logger.Warnf("Parse dubbo interface provider %s failed; due to \n %s", interfaceDetailString, err.Error())
			errorStack = append(errorStack, errors.WithStack(err))
			continue
		}
		apiPattern := strings.Join([]string{"/" + bkConfig.ApplicationName, bkConfig.Interface, bkConfig.Version}, constant.PathSlash)
		mappingParams := []config.MappingParam{
			{
				Name:  "requestBody.values",
				MapTo: "opt.values",
			},
			{
				Name:  "requestBody.types",
				MapTo: "opt.types",
			},
		}
		for i := range methods {
			apis = append(apis, registry.CreateAPIConfig(apiPattern, bkConfig, methods[i], mappingParams))
		}
	}
	return apis, errorStack
}

// LoadApplications load services registered after dubbo 2.7
func (r *ZKRegistry) LoadApplications() ([]router.API, []error) {
	instances, errorStack := r.getInstances()
	apis := []router.API{}

	for _, instance := range instances {
		metaData := instance.GetMetadata()
		metadataStorageType, ok := metaData[constant.MetadataStorageTypeKey]
		if !ok {
			metadataStorageType = constant.DefaultMetadataStorageType
		}
		// get metadata service proxy factory according to the metadataStorageType
		proxyFactory := ex.GetMetadataServiceProxyFactory(metadataStorageType)
		if proxyFactory == nil {
			continue
		}
		metadataService := proxyFactory.GetProxy(instance)
		if metadataService == nil {
			continue
		}
		// call GetExportedURLs to get the exported urls
		urls, err := metadataService.GetExportedURLs(constant.AnyValue, constant.AnyValue, constant.AnyValue, constant.AnyValue)
		if err != nil {
			logger.Errorf("Get exported urls of instance %s failed; due to \n %s", instance, err.Error())
			continue
		}

		for _, url := range urls {
			bkConfig, _, err := registry.ParseDubboString(url.(string))
			if len(bkConfig.ApplicationName) == 0 || len(bkConfig.Interface) == 0 {
				errorStack = append(errorStack, errors.Errorf("Path %s contains dubbo registration that interface or application not set", metadataService))
				continue
			}
			if err != nil {
				logger.Warnf("Parse dubbo interface provider %s failed; due to \n %s", url.(string), err.Error())
				errorStack = append(errorStack, errors.WithStack(err))
				continue
			}
			methods, err := r.getMethods(bkConfig.Interface)
			if err != nil {
				logger.Warnf("Get methods of interface provider %s failed; due to \n %s", bkConfig.Interface, err.Error())
				errorStack = append(errorStack, errors.WithStack(err))
				continue
			}

			apiPattern := strings.Join([]string{"/" + bkConfig.ApplicationName, bkConfig.Interface, bkConfig.Version}, constant.PathSlash)
			mappingParams := []config.MappingParam{
				{
					Name:  "requestBody.values",
					MapTo: "opt.values",
				},
				{
					Name:  "requestBody.types",
					MapTo: "opt.types",
				},
			}
			for i := range methods {
				apis = append(apis, registry.CreateAPIConfig(apiPattern, bkConfig, methods[i], mappingParams))
			}
		}
	}
	return apis, errorStack
}

// getMethods will return the methods of a service
func (r *ZKRegistry) getMethods(in string) ([]string, error) {
	methods := []string{}

	path := strings.Join([]string{methodsRootPath, in}, constant.PathSlash)
	data, err := r.client.GetContent(path)
	if err != nil {
		return nil, err
	}

	sd := &definition.ServiceDefinition{}
	err = json.Unmarshal(data, sd)
	if err != nil {
		return nil, err
	}

	for _, m := range sd.Methods {
		methods = append(methods, m.Name)
	}
	return methods, nil
}

// getInstances will return the instances
func (r *ZKRegistry) getInstances() ([]dr.ServiceInstance, []error) {
	subPaths, err := r.client.GetChildren(r.servicesPath)
	if err != nil {
		return nil, []error{err}
	}
	if len(subPaths) == 0 {
		return nil, nil
	}

	errorStack := []error{}
	iss := []dr.ServiceInstance{}
	for _, subPath := range subPaths {
		serviceName := strings.Join([]string{r.servicesPath, subPath}, constant.PathSlash)
		ids, err := r.client.GetChildren(serviceName)
		if err != nil {
			logger.Warnf("Get serviceIDs %s failed due to %s", serviceName, err.Error())
			errorStack = append(errorStack, errors.WithStack(err))
			continue
		}
		for _, id := range ids {
			path := strings.Join([]string{serviceName, id}, constant.PathSlash)
			data, err := r.client.GetContent(path)
			if err != nil {
				logger.Warnf("Get service instance %s failed due to %s", path, err.Error())
				errorStack = append(errorStack, errors.WithStack(err))
				continue
			}
			instance := &curator_discovery.ServiceInstance{}
			err = json.Unmarshal(data, instance)
			if err != nil {
				logger.Warnf("Parse service instance %s failed due to %s", path, err.Error())
				errorStack = append(errorStack, errors.WithStack(err))
				continue
			}
			iss = append(iss, toZookeeperInstance(instance))
		}
	}
	return iss, errorStack
}

// toZookeeperInstance convert to registry's service instance
func toZookeeperInstance(cris *curator_discovery.ServiceInstance) dr.ServiceInstance {
	pl, ok := cris.Payload.(map[string]interface{})
	if !ok {
		logger.Errorf("toZookeeperInstance{%s} payload is not map[string]interface{}", cris.Id)
		return nil
	}
	mdi, ok := pl["metadata"].(map[string]interface{})
	if !ok {
		logger.Errorf("toZookeeperInstance{%s} metadata is not map[string]interface{}", cris.Id)
		return nil
	}
	md := make(map[string]string, len(mdi))
	for k, v := range mdi {
		md[k] = fmt.Sprint(v)
	}
	return &dr.DefaultServiceInstance{
		Id:          cris.Id,
		ServiceName: cris.Name,
		Host:        cris.Address,
		Port:        cris.Port,
		Enable:      true,
		Healthy:     true,
		Metadata:    md,
	}
}

// Subscribe monitors the target registry.
func (r *ZKRegistry) DoSubscribe(*common.URL) error {
	return nil
}

// Unsubscribe stops monitoring the target registry.
func (r *ZKRegistry) DoUnsubscribe(*common.URL) error {
	return nil
}

// CreateAPIFromRegistry creates the router from registry and save to local cache
func CreateAPIFromRegistry(api router.API) error {
	localAPIDiscSrv := extension.GetMustAPIDiscoveryService(constant.LocalMemoryApiDiscoveryService)
	err := localAPIDiscSrv.AddAPI(api)
	if err != nil {
		return err
	}
	return nil
}

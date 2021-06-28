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
	"strings"
	"time"
)
import (
	"github.com/apache/dubbo-go/common"
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
	zk "github.com/apache/dubbo-go-pixiu/pkg/remoting/zookeeper"
)

const (
	// RegistryZkClient zk client name
	RegistryZkClient = "zk registry"
	rootPath = "/dubbo"
)

func init() {
	extension.SetRegistry(constant.Zookeeper, newZKRegistry)
}

type ZKRegistry struct {
	*registry.BaseRegistry
	client *zk.ZookeeperClient
}

func newZKRegistry(regConfig model.Registry) (registry.Registry, error) {
	baseReg := registry.NewBaseRegistry()
	timeout, err := time.ParseDuration(regConfig.Timeout)
	if err != nil {
		return nil, errors.Errorf("Incorrect timeout configuration: %s", regConfig.Timeout)
	}
	client, err := zk.NewZookeeperClient(RegistryZkClient, strings.Split(regConfig.Address, ","), timeout)
	if err != nil {
		return nil, errors.Errorf("Initialize zookeeper client failed: %s", err.Error())
	}
	return &ZKRegistry{
		BaseRegistry: baseReg,
		client: client,
	}, nil
}

// LoadServices loads all the registered Dubbo services from registry
func (r *ZKRegistry)LoadServices() {
	r.LoadInterfaces()
	r.LoadApplications()
}

// LoadInterfaces load services registered before dubbo 2.7
func (r *ZKRegistry)LoadInterfaces() []error {
	subPaths, err := r.client.GetChildren(rootPath)
	if err != nil {
		return []error{err}
	}
	if len(subPaths) == 0 {
		return nil
	}
	errorStack := []error{}
	for i := range subPaths {
		if subPaths[i] == "metadata" {
			continue
		}
		providerPath := subPaths[i] + "/providers"
		providerString, err := r.client.GetChildren(providerPath)
		if err != nil {
			logger.Warnf("Get provider %s failed due to %s", providerPath, err.Error())
			errorStack = append(errorStack, errors.WithStack(err));
		}
		interfaceDetailString := providerString[0]
		bkConfig, methods, err := registry.ParseDubboString(interfaceDetailString)
		if err != nil {
			logger.Warnf("Parse dubbo interface provider %s failed; due to \n %s", interfaceDetailString, err.Error())
			errorStack = append(errorStack, errors.WithStack(err));
		}
		localAPIDiscSrv := extension.GetMustAPIDiscoveryService(constant.LocalMemoryApiDiscoveryService)
		apiPattern := bkConfig.ApplicationName + "/" + bkConfig.Interface + "/" + bkConfig.Version
		mappingParams := []config.MappingParam{
			{
				Name: "requestBody.values",
				MapTo: "opt.values",
			},
			{
				Name: "requestBody.types",
				MapTo: "opt.types",
			},
		}
		for i := range methods {
			err := localAPIDiscSrv.AddAPI(registry.CreateAPIConfig(apiPattern, bkConfig, methods[i], mappingParams))
			if err != nil {
				logger.Warnf("Add URL %s/method failed; due to \n %s", interfaceDetailString, err.Error())
				errorStack = append(errorStack, errors.WithStack(err));
			}
		}
	} 
	return errorStack
}

// LoadApplications load services registered after dubbo 2.7
func (r *ZKRegistry)LoadApplications() {}

// Subscribe monitors the target registry.
func (r *ZKRegistry) Subscribe(common.URL) error {
	return nil
}
// Unsubscribe stops monitoring the target registry.
func (r *ZKRegistry)Unsubscribe(common.URL) error {
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

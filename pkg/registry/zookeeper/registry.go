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
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
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

var (
	_ baseregistry.FacadeRegistry = new(ZKRegistry)
)

const (
	// RegistryZkClient zk client name
	RegistryZkClient = "zk registry"
	rootPath         = "/dubbo"
	providerCategory = "providers"
)

func init() {
	extension.SetRegistry(constant.Zookeeper, newZKRegistry)
}

type ZKRegistry struct {
	*baseregistry.BaseRegistry
	zkListeners  map[registry.RegisteredType]registry.Listener
	client *zk.ZooKeeperClient
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
	zkReg.BaseRegistry = baseReg
	zkReg.client = client
	initZKListeners(zkReg)
	return zkReg, nil
}

func initZKListeners(reg *ZKRegistry) {
	reg.zkListeners = make(map[registry.RegisteredType]registry.Listener)
	reg.zkListeners[registry.RegisteredTypeInterface] = newZKIntfListener(reg.client, reg)
	go reg.zkListeners[registry.RegisteredTypeInterface].WatchAndHandle()
	//reg.zkListeners[registry.RegisteredTypeApplication] = newZKIntfListener(reg.client, registry.RegisteredTypeApplication)
	//go reg.zkListeners[registry.RegisteredTypeInterface].Watch()
}

func (r *ZKRegistry) GetClient() *zk.ZooKeeperClient {
	return r.client
}

// LoadInterfaces load services registered before dubbo 2.7
func (r *ZKRegistry) LoadInterfaces() ([]registry.RegisteredAPI, []error) {
	subPaths, err := r.client.GetChildren(rootPath)
	if err != nil {
		return nil, []error{err}
	}
	if len(subPaths) == 0 {
		return nil, nil
	}
	var errorStack []error
	var apis []registry.RegisteredAPI
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
		if err != nil {
			logger.Warnf("Parse dubbo interface provider %s failed; due to \n %s", interfaceDetailString, err.Error())
			errorStack = append(errorStack, errors.WithStack(err))
			continue
		}
		if len(bkConfig.ApplicationName) == 0 || len(bkConfig.Interface) == 0 {
			errorStack = append(errorStack, errors.Errorf("Path %s contains dubbo registration that interface or application not set", providerPath))
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
			apis = append(apis, registry.RegisteredAPI{
				API:            registry.CreateAPIConfig(apiPattern, bkConfig, methods[i], mappingParams),
				RegisteredType: registry.RegisteredTypeInterface,
				RegisteredPath: providerPath,
			})
		}
	}
	return apis, errorStack
}

// LoadApplications load services registered after dubbo 2.7
func (r *ZKRegistry) LoadApplications() ([]registry.RegisteredAPI, []error) {
	return nil, nil
}

// DoSubscribe is the implementation of subscription on the target registry for interface level.
func (r *ZKRegistry) DoSubscribe(service registry.RegisteredAPI) error {
	intfListener, ok := r.zkListeners[registry.RegisteredTypeInterface]
	if !ok {
		return errors.New("Listener for interface level registration does not initialized")
	}
	go intfListener.WatchAndHandle()
	return nil
}

// DoSDSubscribe monitors the target registry for application level service discovery
func (r *ZKRegistry) DoSDSubscribe(service registry.RegisteredAPI) error {
	return nil
}

// DoUnsubscribe stops monitoring the target registry.
func (r *ZKRegistry) DoUnsubscribe(service registry.RegisteredAPI) error {
	return nil
}

// DoSDUnsubscribe monitors the target registry.
func (r *ZKRegistry) DoSDUnsubscribe(service registry.RegisteredAPI) error {
	return nil
}

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

package nacos

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"
	"sync"
)

import (
	dubboCommon "dubbo.apache.org/dubbo-go/v3/common"
	dubboConst "dubbo.apache.org/dubbo-go/v3/common/constant"
	dr "dubbo.apache.org/dubbo-go/v3/registry"
	"dubbo.apache.org/dubbo-go/v3/registry/servicediscovery"
	"dubbo.apache.org/dubbo-go/v3/remoting"

	"github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	nacosModel "github.com/nacos-group/nacos-sdk-go/model"
)

import (
	common2 "github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/common"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

var _ registry.Listener = new(appServiceListener)

type appServiceListener struct {
	client      naming_client.INamingClient
	instanceMap map[string]nacosModel.Instance
	cacheLock   sync.Mutex

	exit            chan struct{}
	wg              sync.WaitGroup
	adapterListener common2.RegistryEventListener
}

func newNacosAppSrvListener(client naming_client.INamingClient, adapterListener common2.RegistryEventListener) *appServiceListener {
	return &appServiceListener{
		client:          client,
		exit:            make(chan struct{}),
		adapterListener: adapterListener,
		instanceMap:     map[string]nacosModel.Instance{},
	}
}

func (l *appServiceListener) WatchAndHandle() {
	panic("implement me")
}

func (l *appServiceListener) Close() {
	close(l.exit)
	l.wg.Wait()
}

func (l *appServiceListener) Callback(services []nacosModel.SubscribeService, err error) {
	if err != nil {
		logger.Errorf("nacos subscribe callback error:%s", err.Error())
		return
	}

	addInstances := make([]nacosModel.Instance, 0, len(services))
	delInstances := make([]nacosModel.Instance, 0, len(services))
	updateInstances := make([]nacosModel.Instance, 0, len(services))
	newInstanceMap := make(map[string]nacosModel.Instance, len(services))

	l.cacheLock.Lock()
	defer l.cacheLock.Unlock()
	for i := range services {
		if !services[i].Enable {
			// instance is not available, so ignore it
			continue
		}
		host := services[i].Ip + ":" + strconv.Itoa(int(services[i].Port))
		services[i].ServiceName = handleServiceName(services[i].ServiceName)
		instance := generateInstance(services[i])
		newInstanceMap[host] = instance
		if old, ok := l.instanceMap[host]; ok {
			// instance does not exist in cache, add it to cache
			addInstances = append(addInstances, instance)
		} else {
			if !reflect.DeepEqual(old, instance) {
				// instance is not different from cache, update it to cache
				updateInstances = append(updateInstances, instance)
			}
		}
	}

	for host, inst := range l.instanceMap {
		if _, ok := newInstanceMap[host]; !ok {
			// cache instance does not exist in new instance list, remove it from cache
			delInstances = append(delInstances, inst)
		}
	}

	l.instanceMap = newInstanceMap
	for i := range addInstances {
		newURLs := l.getURLs(addInstances[i])
		for _, url := range newURLs {
			l.handle(url, remoting.EventTypeAdd)
		}
	}
	for i := range delInstances {
		newURLs := l.getURLs(delInstances[i])
		for _, url := range newURLs {
			l.handle(url, remoting.EventTypeDel)
		}
	}
	for i := range updateInstances {
		newURLs := l.getURLs(updateInstances[i])
		for _, url := range newURLs {
			l.handle(url, remoting.EventTypeUpdate)
		}
	}
}

func (l *appServiceListener) handle(url *dubboCommon.URL, action remoting.EventType) {
	logger.Infof("update begin, service event : %v %v", action, url)

	// NOTE: _ is methods, we can not get methods by application discovery
	bkConfig, _, location, err := registry.ParseDubboString(url.String())
	if err != nil {
		logger.Errorf("parse dubbo url error = %s", err)
		return
	}

	apiPattern := registry.GetAPIPattern(bkConfig)
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

	api := registry.CreateAPIConfig(apiPattern, location, bkConfig, constant.AnyValue, mappingParams)
	if action == remoting.EventTypeDel {
		if err := l.adapterListener.OnRemoveAPI(api); err != nil {
			logger.Errorf("Error={%s} happens when try to remove api %s", err.Error(), api.Path)
			return
		}
	} else {
		if err := l.adapterListener.OnAddAPI(api); err != nil {
			logger.Errorf("Error={%s} happens when try to add api %s", err.Error(), api.Path)
			return
		}
	}
}

func (l *appServiceListener) getURLs(nmis nacosModel.Instance) []*dubboCommon.URL {
	instance := toNacosInstance(nmis)
	metadata := instance.GetMetadata()
	metadataInfo, err := servicediscovery.GetMetadataInfo(instance.GetServiceName(), instance, metadata[dubboConst.ExportedServicesRevisionPropertyName])
	if err != nil {
		logger.Errorf("get instance metadata info error %v", err.Error())
		return nil
	}
	instance.SetServiceMetadata(metadataInfo)
	urls := make([]*dubboCommon.URL, 0, len(metadataInfo.Services))
	for _, service := range metadataInfo.Services {
		urls = append(urls, instance.ToURLs(service)...)
	}
	return urls
}

// toNacosInstance convert to registry's service instance
func toNacosInstance(nmis nacosModel.Instance) dr.ServiceInstance {
	md := make(map[string]string, len(nmis.Metadata))
	for k, v := range nmis.Metadata {
		md[k] = fmt.Sprint(v)
	}
	return &dr.DefaultServiceInstance{
		ID:          nmis.InstanceId,
		ServiceName: nmis.ServiceName,
		Host:        nmis.Ip,
		Port:        int(nmis.Port),
		Enable:      nmis.Enable,
		Healthy:     nmis.Healthy,
		Metadata:    md,
	}
}

// group@@serviceName convert to serviceName
func handleServiceName(serviceName string) string {
	parts := strings.Split(serviceName, "@@")
	if len(parts) > 1 {
		return parts[1]
	}
	return ""
}

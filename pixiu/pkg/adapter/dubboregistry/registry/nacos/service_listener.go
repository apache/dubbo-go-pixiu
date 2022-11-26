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
	"net/url"
	"reflect"
	"strconv"
	"sync"
)

import (
	dubboCommon "dubbo.apache.org/dubbo-go/v3/common"
	dubboRegistry "dubbo.apache.org/dubbo-go/v3/registry"
	_ "dubbo.apache.org/dubbo-go/v3/registry/nacos"
	"dubbo.apache.org/dubbo-go/v3/remoting"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	nacosModel "github.com/nacos-group/nacos-sdk-go/model"
)

import (
	common2 "github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/common"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/registry"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

// serviceListener normally monitors the /dubbo/[:url.service()]/providers
type serviceListener struct {
	url         *dubboCommon.URL
	client      naming_client.INamingClient
	instanceMap map[string]nacosModel.Instance
	cacheLock   sync.Mutex

	exit            chan struct{}
	wg              sync.WaitGroup
	adapterListener common2.RegistryEventListener
}

// WatchAndHandle todo WatchAndHandle is useless for service listener
func (z *serviceListener) WatchAndHandle() {
	panic("implement me")
}

// newNacosSrvListener creates a new zk service listener
func newNacosSrvListener(url *dubboCommon.URL, client naming_client.INamingClient, adapterListener common2.RegistryEventListener) *serviceListener {
	return &serviceListener{
		url:             url,
		client:          client,
		exit:            make(chan struct{}),
		adapterListener: adapterListener,
		instanceMap:     map[string]nacosModel.Instance{},
	}
}

func (z *serviceListener) Callback(services []nacosModel.SubscribeService, err error) {
	if err != nil {
		logger.Errorf("nacos subscribe callback error:%s", err.Error())
		return
	}

	addInstances := make([]nacosModel.Instance, 0, len(services))
	delInstances := make([]nacosModel.Instance, 0, len(services))
	updateInstances := make([]nacosModel.Instance, 0, len(services))
	newInstanceMap := make(map[string]nacosModel.Instance, len(services))

	z.cacheLock.Lock()
	defer z.cacheLock.Unlock()
	for i := range services {
		if !services[i].Enable {
			// instance is not available,so ignore it
			continue
		}
		host := services[i].Ip + ":" + strconv.Itoa(int(services[i].Port))
		instance := generateInstance(services[i])
		newInstanceMap[host] = instance
		if old, ok := z.instanceMap[host]; !ok {
			// instance does not exist in cache, add it to cache
			addInstances = append(addInstances, instance)
		} else {
			// instance is not different from cache, update it to cache
			if !reflect.DeepEqual(old, instance) {
				updateInstances = append(updateInstances, instance)
			}
		}
	}

	for host, inst := range z.instanceMap {
		if _, ok := newInstanceMap[host]; !ok {
			// cache instance does not exist in new instance list, remove it from cache
			delInstances = append(delInstances, inst)
		}
	}

	z.instanceMap = newInstanceMap
	for i := range addInstances {
		newUrl := generateURL(addInstances[i])
		if newUrl != nil {
			z.handle(newUrl, remoting.EventTypeAdd)
		}
	}
	for i := range delInstances {
		newUrl := generateURL(delInstances[i])
		if newUrl != nil {
			z.handle(newUrl, remoting.EventTypeDel)
		}
	}

	for i := range updateInstances {
		newUrl := generateURL(updateInstances[i])
		if newUrl != nil {
			z.handle(newUrl, remoting.EventTypeUpdate)
		}
	}
}

func (z *serviceListener) handle(url *dubboCommon.URL, action remoting.EventType) {

	logger.Infof("update begin, service event: %v %v", action, url)

	bkConfig, methods, location, err := registry.ParseDubboString(url.String())
	if err != nil {
		logger.Errorf("parse dubbo string error = %s", err)
		return
	}

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
	apiPattern := registry.GetAPIPattern(bkConfig)
	for i := range methods {
		api := registry.CreateAPIConfig(apiPattern, location, bkConfig, methods[i], mappingParams)
		if action == remoting.EventTypeDel {
			if err := z.adapterListener.OnRemoveAPI(api); err != nil {
				logger.Errorf("Error={%s} happens when try to remove api %s", err.Error(), api.Path)
				continue
			}
		} else {
			if err := z.adapterListener.OnAddAPI(api); err != nil {
				logger.Errorf("Error={%s} happens when try to add api %s", err.Error(), api.Path)
				continue
			}
		}

	}
}

func (z *serviceListener) NotifyAll(e []*dubboRegistry.ServiceEvent, f func()) {
}

// Close closes this listener
func (zkl *serviceListener) Close() {
	close(zkl.exit)
	zkl.wg.Wait()
}

func generateURL(instance nacosModel.Instance) *dubboCommon.URL {
	if instance.Metadata == nil {
		logger.Errorf("nacos instance metadata is empty,instance:%+v", instance)
		return nil
	}
	path := instance.Metadata["path"]
	myInterface := instance.Metadata["interface"]
	if len(path) == 0 && len(myInterface) == 0 {
		logger.Errorf("nacos instance metadata does not have  both path key and interface key,instance:%+v", instance)
		return nil
	}
	if len(path) == 0 && len(myInterface) != 0 {
		path = "/" + myInterface
	}
	protocol := instance.Metadata["protocol"]
	if len(protocol) == 0 {
		logger.Errorf("nacos instance metadata does not have protocol key,instance:%+v", instance)
		return nil
	}
	urlMap := url.Values{}
	for k, v := range instance.Metadata {
		urlMap.Set(k, v)
	}
	return dubboCommon.NewURLWithOptions(
		dubboCommon.WithIp(instance.Ip),
		dubboCommon.WithPort(strconv.Itoa(int(instance.Port))),
		dubboCommon.WithProtocol(protocol),
		dubboCommon.WithParams(urlMap),
		dubboCommon.WithPath(path),
	)
}

func generateInstance(ss nacosModel.SubscribeService) nacosModel.Instance {
	return nacosModel.Instance{
		InstanceId:  ss.InstanceId,
		Ip:          ss.Ip,
		Port:        ss.Port,
		ServiceName: ss.ServiceName,
		Valid:       ss.Valid,
		Enable:      ss.Enable,
		Weight:      ss.Weight,
		Metadata:    ss.Metadata,
		ClusterName: ss.ClusterName,
	}
}

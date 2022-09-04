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
	"strings"
	"sync"
)

import (
	"github.com/nacos-group/nacos-sdk-go/clients/cache"
	model2 "github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/util"
	"github.com/nacos-group/nacos-sdk-go/vo"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/springcloud/servicediscovery"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/remote/nacos"
)

type nacosServiceDiscovery struct {
	targetService []string
	//descriptor    string
	client      *nacos.NacosClient
	config      *model.RemoteConfig
	listener    servicediscovery.ServiceEventListener
	instanceMap map[string]servicediscovery.ServiceInstance

	cacheLock       sync.Mutex
	callbackFlagMap cache.ConcurrentMap
	//done      chan struct{}
}

func (n *nacosServiceDiscovery) Subscribe() error {
	if n.client == nil {
		return perrors.New("nacos naming client stopped")
	}

	serviceNames := n.listener.GetServiceNames()

	for _, serviceName := range serviceNames {
		subscribeParam := &vo.SubscribeParam{
			ServiceName:       serviceName,
			GroupName:         n.config.Group,
			Clusters:          []string{"DEFAULT"},
			SubscribeCallback: n.Callback,
		}
		key := util.GetServiceCacheKey(util.GetGroupName(subscribeParam.ServiceName, subscribeParam.GroupName), strings.Join(subscribeParam.Clusters, ","))
		if n.callbackFlagMap.Has(key) {
			continue
		}
		go func() {
			err := n.client.Subscribe(subscribeParam)
			if err != nil {
				logger.Warnf("[Pixiu-Nacos] Subscribe %s error %s", subscribeParam.ServiceName, err)
				return
			}
			n.callbackFlagMap.Set(key, true)
		}()
	}
	return nil
}

func (n *nacosServiceDiscovery) Unsubscribe() error {
	if n.client == nil {
		return perrors.New("nacos naming client stopped")
	}

	serviceNames := n.listener.GetServiceNames()

	for _, serviceName := range serviceNames {
		subscribeParam := &vo.SubscribeParam{
			ServiceName:       serviceName,
			GroupName:         n.config.Group,
			Clusters:          []string{"DEFAULT"},
			SubscribeCallback: n.Callback,
		}
		_ = n.client.Unsubscribe(subscribeParam)
	}
	return nil
}

func (n *nacosServiceDiscovery) Callback(services []model2.SubscribeService, err error) {

	addInstances := make([]servicediscovery.ServiceInstance, 0, len(services))
	delInstances := make([]servicediscovery.ServiceInstance, 0, len(services))
	updateInstances := make([]servicediscovery.ServiceInstance, 0, len(services))
	newInstanceMap := make(map[string]servicediscovery.ServiceInstance, len(services))

	n.cacheLock.Lock()

	for i := range services {
		service := services[i]
		if !service.Enable {
			// instance is not available,so ignore it
			continue
		}

		instance := fromSubscribeServiceToServiceInstance(service)
		key := instance.GetUniqKey()
		newInstanceMap[instance.GetUniqKey()] = instance
		if old, ok := n.instanceMap[key]; !ok {
			// instance does not exist in cache, add it to cache
			addInstances = append(addInstances, instance)
		} else {
			// instance is not different from cache, update it to cache
			if !reflect.DeepEqual(old, instance) {
				updateInstances = append(updateInstances, instance)
			}
		}
	}

	for host, inst := range n.instanceMap {
		if _, ok := newInstanceMap[host]; !ok {
			// cache instance does not exist in new instance list, remove it from cache
			delInstances = append(delInstances, inst)
		}
	}

	n.instanceMap = newInstanceMap
	n.cacheLock.Unlock()

	for _, instance := range addInstances {
		n.listener.OnAddServiceInstance(&instance)
	}
	for _, instance := range delInstances {
		n.listener.OnDeleteServiceInstance(&instance)
	}

	for _, instance := range updateInstances {
		n.listener.OnUpdateServiceInstance(&instance)
	}
}

func (n *nacosServiceDiscovery) QueryServicesByName(serviceNames []string) ([]servicediscovery.ServiceInstance, error) {
	res := make([]servicediscovery.ServiceInstance, 0, len(serviceNames))

	// need get all service instance api
	for _, serviceName := range serviceNames {

		instances, err := n.client.SelectInstances(vo.SelectInstancesParam{
			ServiceName: serviceName,
			GroupName:   n.config.Group,
			Clusters:    []string{"DEFAULT"},
			HealthyOnly: true,
		})

		if err != nil {
			logger.Warnf("QueryServices SelectInstances {key:%s} = error{%s}", serviceName, err)
			continue
		}

		for _, instance := range instances {
			si := fromInstanceToServiceInstance(serviceName, instance)
			res = append(res, si)
		}
	}

	n.cacheLock.Lock()
	defer n.cacheLock.Unlock()

	for _, instance := range res {
		n.instanceMap[instance.GetUniqKey()] = instance
	}

	return res, nil
}

func (n *nacosServiceDiscovery) QueryAllServices() ([]servicediscovery.ServiceInstance, error) {
	services, err := n.client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
		GroupName: n.config.Group,
		PageSize:  10,
	})

	if err != nil {
		return nil, err
	}

	return n.QueryServicesByName(services.Doms)
}

func (n *nacosServiceDiscovery) Register() error {
	panic("implement me")
}

func (n *nacosServiceDiscovery) UnRegister() error {
	panic("implement me")
}

func (n *nacosServiceDiscovery) Get(s string) []*servicediscovery.ServiceInstance {
	panic("implement me")
}

func (n *nacosServiceDiscovery) StartPeriodicalRefresh() error {
	panic("implement me")
}

func NewNacosServiceDiscovery(targetService []string, config *model.RemoteConfig, l servicediscovery.ServiceEventListener) (servicediscovery.ServiceDiscovery, error) {
	client, err := nacos.NewNacosClient(config)
	if err != nil {
		return nil, err
	}
	return &nacosServiceDiscovery{
		targetService:   targetService,
		client:          client,
		config:          config,
		listener:        l,
		instanceMap:     make(map[string]servicediscovery.ServiceInstance),
		callbackFlagMap: cache.NewConcurrentMap(),
	}, nil
}

func fromInstanceToServiceInstance(serviceName string, instance model2.Instance) servicediscovery.ServiceInstance {
	addr := instance.Ip + ":" + fmt.Sprint(instance.Port)

	return servicediscovery.ServiceInstance{
		// nacos sdk return empty instanceId, so use addr
		//ID: instance.InstanceId,
		ID:          addr,
		ServiceName: serviceName,
		Host:        instance.Ip,
		Port:        int(instance.Port),
		// SelectInstances default return all health instance, not unhealthy
		Healthy:     instance.Healthy,
		Enable:      instance.Enable,
		CLusterName: instance.ClusterName,
		Metadata:    instance.Metadata,
	}
}

func fromSubscribeServiceToServiceInstance(instance model2.SubscribeService) servicediscovery.ServiceInstance {
	addr := instance.Ip + ":" + fmt.Sprint(instance.Port)
	// because it value is DEFAULT_GROUP@@user-service, so split it with @@, and get service name
	serviceName := instance.ServiceName
	tmp := strings.Split(serviceName, "@@")
	if len(tmp) == 2 {
		serviceName = tmp[1]
	}

	return servicediscovery.ServiceInstance{
		// nacos sdk return empty instanceId, so use addr
		//ID: instance.InstanceId,
		ID:          addr,
		ServiceName: serviceName,
		Host:        instance.Ip,
		Port:        int(instance.Port),
		// subscribe callback service should be healthy
		Healthy:     true,
		Enable:      instance.Enable,
		CLusterName: instance.ClusterName,
		Metadata:    instance.Metadata,
	}
}

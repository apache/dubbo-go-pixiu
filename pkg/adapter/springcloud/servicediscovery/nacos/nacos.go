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
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud/servicediscovery"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/remote/nacos"
	"github.com/apache/dubbo-go/registry"
	model2 "github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

type nacosServiceDiscovery struct {
	descriptor        string
	client            *nacos.NacosClient
	config            *model.RemoteConfig
	registryInstances []servicediscovery.ServiceInstance
}

func (n *nacosServiceDiscovery) AddListener(listener servicediscovery.ServiceInstancesChangedListener) {
	for _, serviceName := range listener.GetServiceNames() {
		err := n.client.Subscribe(&vo.SubscribeParam{
			ServiceName: serviceName,
			SubscribeCallback: func(services []model2.SubscribeService, err error) {
				if err != nil {
					logger.Errorf("Could not handle the subscribe notification because the err is not nil."+
						" service name: %s, err: %v", serviceName, err)
				}

				instances := make([]registry.ServiceInstance, 0, len(services))
				for _, service := range services {
					// we won't use the nacos instance id here but use our instance id
					metadata := service.Metadata
					id := metadata[idKey]

					delete(metadata, idKey)

					instances = append(instances, &registry.DefaultServiceInstance{
						ID:          id,
						ServiceName: service.ServiceName,
						Host:        service.Ip,
						Port:        int(service.Port),
						Enable:      service.Enable,
						Healthy:     true,
						Metadata:    metadata,
					})
				}

				e := n.DispatchEventForInstances(serviceName, instances)
				if e != nil {
					logger.Errorf("Dispatching event got exception, service name: %s, err: %v", serviceName, err)
				}
			},
		})
		if err != nil {
			return err
		}
	}
	return nil
}

func (n *nacosServiceDiscovery) Stop() error {
	panic("implement me")
}

func (n *nacosServiceDiscovery) QueryServicesByName(serviceNames []string) ([]servicediscovery.ServiceInstance, error) {
	res := make([]servicediscovery.ServiceInstance, 0, len(serviceNames))

	// need get all service instance api
	for _, serviceName := range serviceNames {

		instances, err := n.client.SelectInstances(vo.SelectInstancesParam{
			ServiceName: serviceName,
			Clusters:    []string{"DEFAULT"},
			HealthyOnly: true,
		})

		if err != nil {
			logger.Warnf("QueryServices SelectInstances {key:%s} = error{%s}", serviceName, err)
			continue
		}

		for _, instance := range instances {
			addr := instance.Ip + ":" + fmt.Sprint(instance.Port)
			si := servicediscovery.ServiceInstance{
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
			res = append(res, si)
		}
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

func NewNacosServiceDiscovery(config *model.RemoteConfig) (servicediscovery.ServiceDiscovery, error) {
	client, err := nacos.NewNacosClient(config)
	if err != nil {
		return nil, err
	}
	return &nacosServiceDiscovery{client: client, config: config}, nil
}

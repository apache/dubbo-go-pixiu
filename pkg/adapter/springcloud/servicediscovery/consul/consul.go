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

package consul

import (
	"github.com/hashicorp/consul/api"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud/servicediscovery"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/remote/consul"
)

type consulServiceDiscovery struct {
	targetService     []string
	client            *api.Client
	config            *model.RemoteConfig
	//registryInstances []servicediscovery.ServiceInstance
	listener          servicediscovery.ServiceEventListener
}

func (c *consulServiceDiscovery) QueryAllServices() ([]servicediscovery.ServiceInstance, error) {
	panic("implement me")
}

func (c *consulServiceDiscovery) QueryServicesByName(serviceNames []string) ([]servicediscovery.ServiceInstance, error) {
	panic("implement me")
}

func (c *consulServiceDiscovery) Subscribe() error {
	panic("implement me")
}

func (c *consulServiceDiscovery) Unsubscribe() error {
	panic("implement me")
}

func (c *consulServiceDiscovery) AddListener(s string) {
	panic("implement me")
}

func (c *consulServiceDiscovery) QueryServices() ([]servicediscovery.ServiceInstance, error) {

	m, _, err := c.client.Catalog().Services(nil)

	if err != nil {
		return nil, err
	}
	res := make([]servicediscovery.ServiceInstance, 0, len(m))

	for serviceId, value := range m {

		logger.Debugf("serviceId : ", serviceId, " value : ", value)

		catalogService, _, _ := c.client.Catalog().Service(serviceId, "", nil)

		if len(catalogService) > 0 {
			for _, sever := range catalogService {
				si := servicediscovery.ServiceInstance{
					ID:          sever.ServiceID,
					ServiceName: sever.ServiceName,
					Host:        sever.Address,
					Port:        sever.ServicePort,
					Metadata:    sever.ServiceMeta,
				}
				res = append(res, si)
			}
		}
	}
	return res, nil
}

func (c *consulServiceDiscovery) Register() error {
	panic("implement me")
}

func (c *consulServiceDiscovery) UnRegister() error {
	panic("implement me")
}

// NewConsulServiceDiscovery
func NewConsulServiceDiscovery(targetService []string, config *model.RemoteConfig, l servicediscovery.ServiceEventListener) (servicediscovery.ServiceDiscovery, error) {
	client, err := consul.NewConsulClient(config)
	if err != nil {
		return nil, err
	}
	return &consulServiceDiscovery{targetService: targetService, client: client, config: config, listener: l}, nil
}

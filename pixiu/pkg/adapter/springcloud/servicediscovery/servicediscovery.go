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

package servicediscovery

import (
	"fmt"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

type (
	// ServiceInstance the service instance info fetched from registry such as nacos consul
	ServiceInstance struct {
		ID          string
		ServiceName string
		// host:port
		Host        string
		Port        int
		Healthy     bool
		CLusterName string
		Enable      bool
		// extra info such as label or other meta data
		Metadata map[string]string
	}

	ServiceEventListener interface {
		OnAddServiceInstance(r *ServiceInstance)
		OnDeleteServiceInstance(r *ServiceInstance)
		OnUpdateServiceInstance(r *ServiceInstance)
		GetServiceNames() []string
	}

	ServiceDiscovery interface {

		// QueryAllServices get all service from remote registry center
		QueryAllServices() ([]ServiceInstance, error)

		// QueryServicesByName get service by serviceName from remote registry center
		QueryServicesByName(serviceNames []string) ([]ServiceInstance, error)

		// Register register to remote registry center
		Register() error

		// UnRegister unregister to remote registry center
		UnRegister() error

		// Subscribe subscribe the service event from remote registry center
		Subscribe() error

		// Unsubscribe unsubscribe from remote registry center
		Unsubscribe() error
	}
)

func (i *ServiceInstance) GetUniqKey() string {
	return i.ServiceName + i.Host + fmt.Sprint(i.Port)
}

// ToEndpoint
func (i *ServiceInstance) ToEndpoint() *model.Endpoint {
	a := model.SocketAddress{Address: i.Host, Port: i.Port}
	return &model.Endpoint{ID: i.ID, Address: a, Name: i.ServiceName, Metadata: i.Metadata}
}

// ToRoute route ID is cluster name, so equal with endpoint name and routerMatch prefix is also service name
func (i *ServiceInstance) ToRoute() *model.Router {
	rm := model.NewRouterMatchPrefix(i.ServiceName)
	ra := model.RouteAction{Cluster: i.ServiceName}
	return &model.Router{ID: i.ServiceName, Match: rm, Route: ra}
}

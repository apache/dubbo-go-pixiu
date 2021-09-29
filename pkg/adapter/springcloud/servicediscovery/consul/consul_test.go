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
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud/servicediscovery"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"testing"
)

type DemoListener struct {
}

func (a *DemoListener) OnAddServiceInstance(instance *servicediscovery.ServiceInstance) {

}

func (a *DemoListener) OnDeleteServiceInstance(instance *servicediscovery.ServiceInstance) {
}

func (a *DemoListener) OnUpdateServiceInstance(instance *servicediscovery.ServiceInstance) {
}

func (a *DemoListener) GetServiceNames() []string {
	return nil
}

func TestNewConsulServiceDiscovery(t *testing.T) {

	config := &model.RemoteConfig{
		Address: "127.0.0.1:8500",
	}

	discovery, err := NewConsulServiceDiscovery(make([]string, 0), config, &DemoListener{})

	if err != nil {
		panic(err)
	}

	services, err := discovery.QueryAllServices()

	if err != nil {
		panic(err)
	}

	logger.Info("fetch services from consul : ", services)
}

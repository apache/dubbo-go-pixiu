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
package registry

import (
	"net/url"
	"strconv"
	"testing"
)

import (
	"github.com/apache/dubbo-go/common"
	"github.com/apache/dubbo-go/common/constant"
	"github.com/apache/dubbo-go/common/extension"
	_ "github.com/apache/dubbo-go/registry/consul"
	"github.com/apache/dubbo-go/remoting/consul"
	"github.com/stretchr/testify/assert"
)

var (
	registryHost = "localhost"
	registryPort = 8500
	providerHost = "localhost"
	providerPort = 8000
	consumerHost = "localhost"
	consumerPort = 8001
	service      = "HelloWorld"
	protocol     = "tcp"
	cluster      = "test_cluster"
)

func TestConsulRegistryLoad_GetCluster(t *testing.T) {
	loader, err := newConsulRegistryLoad(registryHost+":"+strconv.Itoa(registryPort), "test_cluster")
	assert.Nil(t, err)
	consulCluster, err := loader.GetCluster()
	assert.Nil(t, err)
	assert.Equal(t, cluster, consulCluster)
}

func TestConsulRegistryLoad_LoadAllServices(t *testing.T) {
	consulAgent := consul.NewConsulAgent(t, registryPort)
	defer consulAgent.Shutdown()
	registryUrl, _ := common.NewURL(protocol + "://" + providerHost + ":" + strconv.Itoa(providerPort) + "/" + service + "?anyhost=true&" +
		"application=BDTService&category=providers&default.timeout=10000&dubbo=dubbo-provider-golang-1.0.0&" +
		"environment=dev&interface=com.ikurento.user.UserProvider&ip=192.168.56.1&methods=GetUser%2C&" +
		"module=dubbogo+user-info+server&org=ikurento.com&owner=ZX&pid=1447&revision=0.0.1&" +
		"side=provider&timeout=3000&timestamp=1556509797245")

	registry, err := extension.GetRegistry("consul", common.NewURLWithOptions(common.WithParams(url.Values{}), common.WithIp("localhost"), common.WithPort("8500"), common.WithParamsValue(constant.ROLE_KEY, strconv.Itoa(common.PROVIDER))))
	assert.Nil(t, err)
	err = registry.Register(registryUrl)
	assert.Nil(t, err)
	defer registry.UnRegister(registryUrl)
	loader, err := newConsulRegistryLoad(registryHost+":"+strconv.Itoa(registryPort), "test_cluster")
	assert.Nil(t, err)
	services, err := loader.LoadAllServices()
	assert.Nil(t, err)
	assert.Len(t, services, 1)
	assert.Contains(t, services[0].Methods, "GetUser")
	assert.Equal(t, services[0].GetParams(), registryUrl.GetParams())
}

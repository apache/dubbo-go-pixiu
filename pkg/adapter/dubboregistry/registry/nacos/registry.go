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
	"strconv"
	"strings"
)

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	nacosConstant "github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/common"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry"
	baseRegistry "github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry/base"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func init() {
	registry.SetRegistry(constant.Nacos, newNacosRegistry)
}

type NacosRegistry struct {
	*baseRegistry.BaseRegistry
	nacosListeners map[registry.RegisteredType]registry.Listener
	client         naming_client.INamingClient
}

func (n *NacosRegistry) DoSubscribe() error {
	intfListener, ok := n.nacosListeners[registry.RegisteredTypeInterface]
	if !ok {
		return errors.New("Listener for interface level registration does not initialized")
	}
	go intfListener.WatchAndHandle()
	return nil
}

func (n *NacosRegistry) DoUnsubscribe() error {
	panic("implement me")
}

var _ registry.Registry = new(NacosRegistry)

func newNacosRegistry(regConfig model.Registry, adapterListener common.RegistryEventListener) (registry.Registry, error) {
	addrAndIP := strings.Split(regConfig.Address, ":")
	if len(addrAndIP) != 2 {
		return nil, errors.Errorf("nacos registry ip:port = %s is invalid, please check.", regConfig.Address)
	}
	port, err := strconv.Atoi(addrAndIP[1])
	if err != nil {
		return nil, errors.Errorf("nacos registry port = %s is not number, please check.", addrAndIP[1])
	}
	scs := []nacosConstant.ServerConfig{
		*nacosConstant.NewServerConfig(addrAndIP[0], uint64(port)),
	}
	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ServerConfigs: scs,
	})
	if err != nil {
		return nil, err
	}
	nacosRegistry := &NacosRegistry{
		client:         client,
		nacosListeners: make(map[registry.RegisteredType]registry.Listener),
	}
	nacosRegistry.nacosListeners[registry.RegisteredTypeInterface] = newNacosIntfListener(client, regConfig.Address, nacosRegistry, adapterListener)

	baseReg := baseRegistry.NewBaseRegistry(nacosRegistry, adapterListener)
	nacosRegistry.BaseRegistry = baseReg
	return baseReg, nil
}

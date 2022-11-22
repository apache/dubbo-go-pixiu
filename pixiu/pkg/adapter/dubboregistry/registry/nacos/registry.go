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
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	nacosConstant "github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/common"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/registry"
	baseRegistry "github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/registry/base"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/util/stringutil"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
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
	addrs, err := stringutil.GetIPAndPort(regConfig.Address)
	if err != nil {
		return nil, err
	}

	scs := make([]nacosConstant.ServerConfig, 0)
	for _, addr := range addrs {
		scs = append(scs, nacosConstant.ServerConfig{
			IpAddr: addr.IP.String(),
			Port:   uint64(addr.Port),
		})
	}

	ccs := nacosConstant.NewClientConfig(
		nacosConstant.WithNamespaceId(regConfig.Namespace),
		nacosConstant.WithUsername(regConfig.Username),
		nacosConstant.WithPassword(regConfig.Password),
		nacosConstant.WithNotLoadCacheAtStart(true),
		nacosConstant.WithUpdateCacheWhenEmpty(true))
	client, err := clients.NewNamingClient(vo.NacosClientParam{
		ServerConfigs: scs,
		ClientConfig:  ccs,
	})
	if err != nil {
		return nil, err
	}

	nacosRegistry := &NacosRegistry{
		client:         client,
		nacosListeners: make(map[registry.RegisteredType]registry.Listener),
	}
	nacosRegistry.nacosListeners[registry.RegisteredTypeInterface] = newNacosIntfListener(client, nacosRegistry, &regConfig, adapterListener)

	baseReg := baseRegistry.NewBaseRegistry(nacosRegistry, adapterListener)
	nacosRegistry.BaseRegistry = baseReg
	return baseReg, nil
}

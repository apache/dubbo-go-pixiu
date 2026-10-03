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
	"github.com/nacos-group/nacos-sdk-go/v2/clients"
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	nacosConstant "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/llmregistry/common"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/llmregistry/registry"
	baseRegistry "github.com/apache/dubbo-go-pixiu/pkg/adapter/llmregistry/registry/base"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/util/stringutil"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func init() {
	registry.SetRegistry(constant.Nacos, newNacosRegistry)
}

type NacosRegistry struct {
	*baseRegistry.BaseRegistry
	nacosListener *listener
	client        naming_client.INamingClient
}

func (n *NacosRegistry) DoSubscribe() error {
	go n.nacosListener.WatchAndHandle()
	return nil
}

func (n *NacosRegistry) DoUnsubscribe() error {
	// Stop the background listener first: it unsubscribes all services and
	// waits for the watch goroutine to exit so no callback races the close.
	n.nacosListener.Close()
	// v2 clients hold a gRPC connection and internal retry goroutines that
	// survive Unsubscribe; CloseClient() shuts them down (see MCP adapter,
	// pkg/adapter/mcpserver/registry/nacos/client.go). Without it, Stop/Apply
	// leaks the connection and goroutines after graceful shutdown.
	n.client.CloseClient()
	return nil
}

var _ registry.Registry = new(NacosRegistry)

func newNacosRegistry(regConfig model.Registry, adapterListener common.RegistryEventListener) (registry.Registry, error) {
	addrs, err := stringutil.GetIPAndPort(regConfig.Address)
	if err != nil {
		return nil, err
	}

	scs := make([]nacosConstant.ServerConfig, 0, len(addrs))
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

	return newNacosRegistryWithClient(regConfig, client, adapterListener)
}

// newNacosRegistryWithClient assembles a NacosRegistry around an existing
// naming client. It is split out from newNacosRegistry so tests can inject a
// mock client and assert on its lifecycle (e.g. that CloseClient is invoked
// on unsubscribe).
func newNacosRegistryWithClient(regConfig model.Registry, client naming_client.INamingClient, adapterListener common.RegistryEventListener) (*NacosRegistry, error) {
	nacosRegistry := &NacosRegistry{
		client: client,
	}
	nacosRegistry.BaseRegistry = baseRegistry.NewBaseRegistry(nacosRegistry, adapterListener)
	nacosRegistry.nacosListener = newNacosListener(client, &regConfig, adapterListener)

	return nacosRegistry, nil
}

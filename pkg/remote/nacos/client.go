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
	"context"
	"net"
	"strconv"
	"strings"
	"time"
)

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	model2 "github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"

	perrors "github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type NacosClient struct {
	namingClient naming_client.INamingClient
	configClient config_client.IConfigClient
}

func (client *NacosClient) GetAllServicesInfo(param vo.GetAllServiceInfoParam) (model2.ServiceList, error) {
	return client.namingClient.GetAllServicesInfo(param)
}

func (client *NacosClient) SelectInstances(param vo.SelectInstancesParam) ([]model2.Instance, error) {
	return client.namingClient.SelectInstances(param)
}

func (client *NacosClient) Subscribe(param *vo.SubscribeParam) error {
	return client.namingClient.Subscribe(param)
}

func (client *NacosClient) Unsubscribe(param *vo.SubscribeParam) error {
	return client.namingClient.Unsubscribe(param)
}

func NewNacosClient(config *model.RemoteConfig) (*NacosClient, error) {
	configMap := make(map[string]any, 2)
	addresses := strings.Split(config.Address, ",")
	serverConfigs := make([]constant.ServerConfig, 0, len(addresses))
	for _, addr := range addresses {
		ip, portStr, err := net.SplitHostPort(addr)
		if err != nil {
			return nil, perrors.WithMessagef(err, "split [%s] ", addr)
		}
		port, _ := strconv.Atoi(portStr)
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			IpAddr: ip,
			Port:   uint64(port),
		})
	}
	configMap["serverConfigs"] = serverConfigs

	duration, _ := time.ParseDuration(config.Timeout)
	naming, err := clients.NewNamingClient(
		vo.NacosClientParam{
			ClientConfig: constant.NewClientConfig(
				constant.WithTimeoutMs(uint64(duration.Milliseconds())),
				constant.WithNotLoadCacheAtStart(true),
				constant.WithUpdateCacheWhenEmpty(true),
			),
			ServerConfigs: serverConfigs,
		},
	)

	if err != nil {
		return nil, perrors.WithMessagef(err, "nacos naming client create error")
	}

	cfgClient, err := clients.NewConfigClient(vo.NacosClientParam{
		ClientConfig: constant.NewClientConfig(
			constant.WithTimeoutMs(uint64(duration.Milliseconds())),
			constant.WithNotLoadCacheAtStart(true),
			constant.WithUpdateCacheWhenEmpty(true),
		),
		ServerConfigs: serverConfigs,
	})
	if err != nil {
		return nil, perrors.WithMessagef(err, "nacos config client create error")
	}

	return &NacosClient{namingClient: naming, configClient: cfgClient}, nil
}

// Dynamic configuration APIs
// GetConfig retrieves configuration content by dataId and group
func (client *NacosClient) GetConfig(dataId, group string) (string, error) {
	if client.configClient == nil {
		return "", perrors.New("nacos config client not initialized")
	}
	return client.configClient.GetConfig(vo.ConfigParam{DataId: dataId, Group: group})
}

// ListenConfig listens for configuration changes; callback will be invoked with latest content
func (client *NacosClient) ListenConfig(ctx context.Context, dataId, group string, callback func(string)) error {
	if client.configClient == nil {
		return perrors.New("nacos config client not initialized")
	}
	return client.configClient.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			select {
			case <-ctx.Done():
				return
			default:
			}
			callback(data)
		},
	})
}

// CancelListenConfig stops listening for configuration changes
func (client *NacosClient) CancelListenConfig(dataId, group string) error {
	if client.configClient == nil {
		return perrors.New("nacos config client not initialized")
	}
	return client.configClient.CancelListenConfig(vo.ConfigParam{DataId: dataId, Group: group})
}

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

package configcenter

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

const (
	KeyDataId  = "dataId"
	KeyGroup   = "group"
	KeyContent = "content"
	KeyTag     = "tag"
	KeyAppName = "appName"
	KeyTenant  = "tenant"
)

const (
	DataId    = "pixiu.yaml"
	Group     = "DEFAULT_GROUP"
	Namespace = "dubbo-go-pixiu"

	IpAddr      = "localhost"
	ContextPath = "/nacos"
	Port        = 8848
	Scheme      = "http"
)

type (
	NacosConfig struct {
		client config_client.IConfigClient

		// todo not support now
		listenConfigCallback ListenConfig
	}
)

func NewNacosConfig(boot *model.Bootstrap) (configClient ConfigClient, err error) {

	var sc []constant.ServerConfig
	if len(boot.Nacos.ServerConfigs) == 0 {
		return nil, errors.New("no nacos server be setted")
	}
	for _, serveConfig := range boot.Nacos.ServerConfigs {
		sc = append(sc, constant.ServerConfig{
			Port:   serveConfig.Port,
			IpAddr: serveConfig.IpAddr,
		})
	}

	cc := constant.ClientConfig{
		NamespaceId:         boot.Nacos.ClientConfig.NamespaceId,
		TimeoutMs:           boot.Nacos.ClientConfig.TimeoutMs,
		NotLoadCacheAtStart: boot.Nacos.ClientConfig.NotLoadCacheAtStart,
		LogDir:              boot.Nacos.ClientConfig.LogDir,
		CacheDir:            boot.Nacos.ClientConfig.CacheDir,
		LogLevel:            boot.Nacos.ClientConfig.LogLevel,
	}

	pa := vo.NacosClientParam{
		ClientConfig:  &cc,
		ServerConfigs: sc,
	}
	nacos, err := clients.NewConfigClient(pa)
	if err != nil {
		return nil, err
	}
	configClient = &NacosConfig{
		client: nacos,
	}

	return configClient, nil
}

func (n *NacosConfig) LoadConfig(param map[string]interface{}) (string, error) {

	return n.client.GetConfig(vo.ConfigParam{
		DataId: getOrDefault(param[KeyDataId].(string), DataId),
		Group:  getOrDefault(param[KeyGroup].(string), Group),
	})
}

func getOrDefault(target string, quiet string) string {
	if len(target) == 0 {
		target = quiet
	}
	return target
}

func (n *NacosConfig) ListenConfig(param map[string]interface{}) (err error) {
	// todo noop, not support
	if true {
		return nil
	}
	listen := n.listen(getOrDefault(param[KeyDataId].(string), DataId), getOrDefault(param[KeyGroup].(string), Group))
	return listen()
}

func (n *NacosConfig) listen(dataId, group string) func() error {
	return func() error {
		return n.client.ListenConfig(vo.ConfigParam{
			DataId: dataId,
			Group:  group,
			OnChange: func(namespace, group, dataId, data string) {
				if len(data) == 0 {
					logger.Errorf("nacos listen callback data nil error ,  namespace : %s，group : %s , dataId : %s , data : %s")
					return
				}
				n.listenConfigCallback(data)
			},
		})
	}
}

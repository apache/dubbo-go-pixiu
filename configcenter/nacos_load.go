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
	"sync"
)

import (
	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// Constants for configuration keys.
const (
	KeyDataId  = "dataId"
	KeyGroup   = "group"
	KeyContent = "content"
	KeyTag     = "tag"
	KeyAppName = "appName"
	KeyTenant  = "tenant"
)

// Constants for Nacos configuration.
const (
	DataId    = "pixiu.yaml"
	Group     = "DEFAULT_GROUP"
	Namespace = "dubbo-go-pixiu"

	IpAddr      = "localhost"
	ContextPath = "/nacos"
	Port        = 8848
	Scheme      = "http"
)

// NacosConfig represents the Nacos configuration client and its state.
type NacosConfig struct {
	client       config_client.IConfigClient
	remoteConfig *model.Bootstrap
	mu           sync.Mutex
}

// NewNacosConfig creates a new NacosConfig instance.
// It returns an error if no Nacos server is configured or if the client cannot be created.
func NewNacosConfig(boot *model.Bootstrap) (ConfigClient, error) {
	if len(boot.Nacos.ServerConfigs) == 0 {
		return nil, errors.New("no Nacos server configured")
	}

	nacosClient, err := getNacosConfigClient(boot)
	if err != nil {
		return nil, err
	}

	return &NacosConfig{
		client: nacosClient,
	}, nil
}

// getNacosConfigClient initializes and returns a Nacos config client.
func getNacosConfigClient(boot *model.Bootstrap) (config_client.IConfigClient, error) {
	var serverConfigs []constant.ServerConfig
	for _, serverConfig := range boot.Nacos.ServerConfigs {
		serverConfigs = append(serverConfigs, constant.ServerConfig{
			Port:   serverConfig.Port,
			IpAddr: serverConfig.IpAddr,
		})
	}

	clientConfig := constant.ClientConfig{
		NamespaceId:         boot.Nacos.ClientConfig.NamespaceId,
		TimeoutMs:           boot.Nacos.ClientConfig.TimeoutMs,
		NotLoadCacheAtStart: boot.Nacos.ClientConfig.NotLoadCacheAtStart,
		LogDir:              boot.Nacos.ClientConfig.LogDir,
		CacheDir:            boot.Nacos.ClientConfig.CacheDir,
		LogLevel:            boot.Nacos.ClientConfig.LogLevel,
	}

	clientParam := vo.NacosClientParam{
		ClientConfig:  &clientConfig,
		ServerConfigs: serverConfigs,
	}

	return clients.NewConfigClient(clientParam)
}

// LoadConfig retrieves the configuration from Nacos based on the provided parameters.
func (n *NacosConfig) LoadConfig(param map[string]any) (string, error) {
	return n.client.GetConfig(vo.ConfigParam{
		DataId: getOrDefault(param[KeyDataId].(string), DataId),
		Group:  getOrDefault(param[KeyGroup].(string), Group),
	})
}

// getOrDefault returns the target value if it is not empty; otherwise, it returns the fallback value.
func getOrDefault(target, fallback string) string {
	if len(target) == 0 {
		return fallback
	}
	return target
}

// ListenConfig listens for configuration changes in Nacos.
func (n *NacosConfig) ListenConfig(param map[string]any) error {
	return n.client.ListenConfig(vo.ConfigParam{
		DataId:   getOrDefault(param[KeyDataId].(string), DataId),
		Group:    getOrDefault(param[KeyGroup].(string), Group),
		OnChange: n.onChange,
	})
}

// onChange is the callback function triggered when the configuration changes in Nacos.
func (n *NacosConfig) onChange(namespace, group, dataId, data string) {
	if len(data) == 0 {
		logger.Errorf("Nacos listen callback data is nil. Namespace: %s, Group: %s, DataId: %s", namespace, group, dataId)
		return
	}

	n.mu.Lock()
	defer n.mu.Unlock()

	var boot model.Bootstrap
	if err := Parsers[".yml"]([]byte(data), &boot); err != nil {
		logger.Errorf("Failed to parse the configuration loaded from the remote. Error: %v", err)
		return
	}

	n.remoteConfig = &boot
}

// ViewConfig returns the current remote configuration.
func (n *NacosConfig) ViewConfig() *model.Bootstrap {
	n.mu.Lock()
	defer n.mu.Unlock()
	return n.remoteConfig
}

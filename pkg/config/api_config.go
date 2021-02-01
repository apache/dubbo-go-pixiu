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

package config

import (
	"strings"
	"sync"
	"time"
)

import (
	"github.com/coreos/etcd/mvcc/mvccpb"
	fc "github.com/dubbogo/dubbo-go-proxy-filter/pkg/api/config"
	etcdv3 "github.com/dubbogo/dubbo-go-proxy/pkg/remoting/etcd3"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/yaml"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
)

var (
	apiConfig *fc.APIConfig
	once      sync.Once
	client    *etcdv3.Client
	listener  APIConfigListener
	lock      sync.RWMutex
)

// APIConfigListener defines api config listener interface
type APIConfigListener interface {
	APIConfigChange(apiConfig fc.APIConfig) bool //bool is return for interface implement is interesting
}

// LoadAPIConfigFromFile load the api config from file
func LoadAPIConfigFromFile(path string) (*fc.APIConfig, error) {
	if len(path) == 0 {
		return nil, perrors.Errorf("Config file not specified")
	}
	logger.Infof("Load API configuration file form %s", path)
	apiConf := &fc.APIConfig{}
	err := yaml.UnmarshalYMLConfig(path, apiConf)
	if err != nil {
		return nil, perrors.Errorf("unmarshalYmlConfig error %v", perrors.WithStack(err))
	}
	apiConfig = apiConf
	return apiConf, nil
}

// LoadAPIConfig load the api config from config center
func LoadAPIConfig(metaConfig *model.APIMetaConfig) (*fc.APIConfig, error) {

	client = etcdv3.NewConfigClient(
		etcdv3.WithName(etcdv3.RegistryETCDV3Client),
		etcdv3.WithTimeout(10*time.Second),
		etcdv3.WithEndpoints(strings.Split(metaConfig.Address, ",")...),
	)

	go listenAPIConfigNodeEvent(metaConfig.APIConfigPath)

	content, err := client.Get(metaConfig.APIConfigPath)

	if err != nil {
		return nil, perrors.Errorf("Get remote config fail error %v", err)
	}

	initAPIConfigFromString(content)

	return apiConfig, nil
}

func initAPIConfigFromString(content string) error{
	lock.Lock()
	defer lock.Unlock()

	apiConf := &fc.APIConfig{}
	if len(content) != 0 {
		err := yaml.UnmarshalYML([]byte(content), apiConf)
		if err != nil {
			return perrors.Errorf("unmarshalYmlConfig error %v", perrors.WithStack(err))
		}
		apiConfig = apiConf
	}
	return nil
}

func listenAPIConfigNodeEvent(key string) bool {
	for {
		wc, err := client.Watch(key)
		if err != nil {
			logger.Warnf("Watch api config {key:%s} = error{%v}", key, err)
			return false
		}

		select {

		// client stopped
		case <-client.Done():
			logger.Warnf("client stopped")
			return false
		// client ctx stop
		// handle etcd events
		case e, ok := <-wc:
			if !ok {
				logger.Warnf("watch-chan closed")
				return false
			}

			if e.Err() != nil {
				logger.Errorf("watch ERR {err: %s}", e.Err())
				continue
			}
			for _, event := range e.Events {
				switch event.Type {
				case mvccpb.PUT:
					initAPIConfigFromString(string(event.Kv.Value))
					listener.APIConfigChange(GetAPIConf())
				case mvccpb.DELETE:
					logger.Warnf("get event (key{%s}) = event{EventNodeDeleted}", event.Kv.Key)
					return true
				default:
					return false
				}
			}
		}
	}
}
// RegisterConfigListener register APIConfigListener
func RegisterConfigListener(li APIConfigListener) {
	listener = li
}

// GetAPIConf returns the initted api config
func GetAPIConf() fc.APIConfig {
	return *apiConfig
}

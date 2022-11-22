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
	"log"
	"os"
	"path/filepath"
	"strings"
	"sync"
)

import (
	"github.com/creasty/defaults"
	"github.com/ghodss/yaml"
	"github.com/goinggo/mapstructure"
	"github.com/imdario/mergo"
)

import (
	"github.com/apache/dubbo-go-pixiu/configcenter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

var (
	configPath     string
	config         *model.Bootstrap
	configLoadFunc LoadFunc = LoadYAMLConfig

	configCenterType map[string]interface{}
	once             sync.Once
)

// LoadFunc ConfigLoadFunc parse a input(usually file path) into a pixiu config
type LoadFunc func(path string) *model.Bootstrap

// GetBootstrap get config global, need a better name
func GetBootstrap() *model.Bootstrap {
	return config
}

// Load config file and parse
func Load(path string) *model.Bootstrap {
	logger.Infof("[dubbopixiu go] load path:%s", path)
	configPath, _ = filepath.Abs(path)
	if configPath != "" && CheckYamlFormat(configPath) {
		RegisterConfigLoadFunc(LoadYAMLConfig)
	}
	if cfg := configLoadFunc(configPath); cfg != nil {
		config = cfg
	}
	return config
}

// RegisterConfigLoadFunc can replace a new config load function instead of default
func RegisterConfigLoadFunc(f LoadFunc) {
	configLoadFunc = f
}

func CheckYamlFormat(path string) bool {
	ext := filepath.Ext(path)
	if ext == constant.YAML || ext == constant.YML {
		return true
	}
	return false
}

// LoadYAMLConfig YAMLConfigLoad config load yaml
func LoadYAMLConfig(path string) *model.Bootstrap {
	log.Println("load config in YAML format from : ", path)
	content, err := os.ReadFile(path)
	if err != nil {
		log.Fatalln("[config] [yaml load] load config failed, ", err)
	}
	cfg := &model.Bootstrap{}
	err = yaml.Unmarshal(content, cfg)
	if err != nil {
		log.Fatalln("[config] [yaml load] convert YAML to JSON failed, ", err)
	}
	if err = defaults.Set(cfg); err != nil {
		log.Fatalln("[config] [yaml load] initialize structs with default value failed, ", err)
	}
	err = Adapter(cfg)
	if err != nil {
		log.Fatalln("[config] [yaml load] yaml unmarshal config failed, ", err)
	}
	return cfg
}

func Adapter(cfg *model.Bootstrap) (err error) {
	if GetHttpConfig(cfg) != nil || GetProtocol(cfg) != nil ||
		GetLoadBalance(cfg) != nil || GetDiscoveryType(cfg) != nil {
		return err
	}
	return nil
}

func GetProtocol(cfg *model.Bootstrap) (err error) {
	if cfg == nil {
		logger.Error("Bootstrap configuration is null")
		return err
	}
	for _, l := range cfg.StaticResources.Listeners {
		if l.ProtocolStr == "" {
			l.ProtocolStr = constant.DefaultProtocolType
		}
		l.Protocol = model.ProtocolType(model.ProtocolTypeValue[l.ProtocolStr])
	}
	return nil
}

func GetHttpConfig(cfg *model.Bootstrap) (err error) {
	if cfg == nil {
		logger.Error("Bootstrap configuration is null")
		return err
	}
	for i, l := range cfg.StaticResources.Listeners {
		hc := &model.HttpConfig{}
		if l.Config != nil {
			if v, ok := l.Config.(map[string]interface{}); ok {
				logger.Info("http config:", v, ok)
				switch l.Name {
				case constant.DefaultHTTPType:
					if err := mapstructure.Decode(v, hc); err != nil {
						logger.Error(err)
					}
					cfg.StaticResources.Listeners[i].Config = *hc
					//cfg.StaticResources.Listeners[i].FilterChains
				}
			}
		}
	}
	return nil
}

func GetLoadBalance(cfg *model.Bootstrap) (err error) {
	if cfg == nil {
		logger.Error("Bootstrap configuration is null")
		return err
	}
	return nil
}

func GetDiscoveryType(cfg *model.Bootstrap) (err error) {
	if cfg == nil {
		logger.Error("Bootstrap configuration is null")
		return err
	}
	var discoveryType model.DiscoveryType
	for _, c := range cfg.StaticResources.Clusters {
		flag := true
		if c.TypeStr != "" {
			if t, ok := model.DiscoveryTypeValue[c.TypeStr]; ok {
				discoveryType = t
				flag = false
			}
		}
		if flag {
			c.TypeStr = constant.DefaultDiscoveryType
			discoveryType = model.DiscoveryTypeValue[c.TypeStr]
		}
		c.Type = discoveryType
	}
	return nil
}

type ConfigManager struct {
	path         string
	localConfig  *model.Bootstrap
	remoteConfig *model.Bootstrap
	load         configcenter.Load
}

func NewConfigManger() *ConfigManager {
	return &ConfigManager{}
}

func (m *ConfigManager) LoadBootConfig(path string) *model.Bootstrap {

	var configs *model.Bootstrap

	// load file
	configs = m.loadLocalBootConfigs(path)

	if m.localConfig != nil && m.localConfig.Config != nil {
		if strings.EqualFold(m.localConfig.Config.Enable, "true") {
			configs = m.loadRemoteBootConfigs()
		}
	}

	config = configs

	err := m.check()

	if err != nil {
		logger.Errorf("[Pixiu-Config] check bootstrap config fail : %s", err.Error())
		panic(err)
	}

	return configs
}

func (m *ConfigManager) loadLocalBootConfigs(path string) *model.Bootstrap {
	m.localConfig = Load(path)
	return m.localConfig
}

func (m *ConfigManager) loadRemoteBootConfigs() *model.Bootstrap {

	bootstrap := m.localConfig

	// load remote
	once.Do(func() {
		m.load = configcenter.NewConfigLoad(bootstrap)
	})

	configs, err := m.load.LoadConfigs(bootstrap, func(opt *configcenter.Options) {
		opt.Remote = true
	})

	if err != nil {
		panic(err)
	}

	err = mergo.Merge(configs, bootstrap, func(c *mergo.Config) {
		c.Overwrite = false
		c.AppendSlice = false
	})

	if err != nil {
		panic(err)
	}

	m.remoteConfig = configs

	return configs
}

func (m *ConfigManager) check() error {

	return Adapter(config)
}

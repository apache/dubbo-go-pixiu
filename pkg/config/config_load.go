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
	"encoding/json"
	"io/ioutil"
	"log"
	"path/filepath"
)

import (
	"github.com/ghodss/yaml"
	"github.com/goinggo/mapstructure"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

var (
	configPath     string
	config         *model.Bootstrap
	configLoadFunc ConfigLoadFunc = DefaultConfigLoad
)

// GetBootstrap get config global, need a better name
func GetBootstrap() *model.Bootstrap {
	return config
}

// Load config file and parse
func Load(path string) *model.Bootstrap {
	logger.Infof("[dubbopixiu go] load path:%s", path)
	configPath, _ = filepath.Abs(path)
	if yamlFormat(configPath) {
		RegisterConfigLoadFunc(YAMLConfigLoad)
	}
	if cfg := configLoadFunc(configPath); cfg != nil {
		config = cfg
	}
	return config
}

// ConfigLoadFunc parse a input(usually file path) into a pixiu config
type ConfigLoadFunc func(path string) *model.Bootstrap

// RegisterConfigLoadFunc can replace a new config load function instead of default
func RegisterConfigLoadFunc(f ConfigLoadFunc) {
	configLoadFunc = f
}

func yamlFormat(path string) bool {
	ext := filepath.Ext(path)
	if ext == ".yaml" || ext == ".yml" {
		return true
	}
	return false
}

// YAMLConfigLoad config load yaml
func YAMLConfigLoad(path string) *model.Bootstrap {
	log.Println("load config in YAML format from : ", path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("[config] [yaml load] load config failed, ", err)
	}
	cfg := &model.Bootstrap{}
	bytes, err := yaml.YAMLToJSON(content)
	if err != nil {
		log.Fatalln("[config] [yaml load] convert YAML to JSON failed, ", err)
	}
	err = json.Unmarshal(bytes, cfg)
	if err != nil {
		log.Fatalln("[config] [yaml load] yaml unmarshal config failed, ", err)
	}
	adapter(cfg)
	return cfg
}

// DefaultConfigLoad if not config, will load this
func DefaultConfigLoad(path string) *model.Bootstrap {
	logger.Debug("load config from : ", path)
	content, err := ioutil.ReadFile(path)
	if err != nil {
		log.Fatalln("[config] [default load] load config failed, ", err)
	}
	cfg := &model.Bootstrap{}
	// translate to lower case
	err = json.Unmarshal(content, cfg)
	if err != nil {
		log.Fatalln("[config] [default load] json unmarshal config failed, ", err)
	}
	return cfg
}

func adapter(cfg *model.Bootstrap) {
	FilterChain(cfg)
	HttpConfig(cfg)
	Protocol(cfg)
	LoadBalance(cfg)
	DiscoverType(cfg)
}

func Protocol(cfg *model.Bootstrap) {
	for _, l := range cfg.StaticResources.Listeners {
		if l.Address.SocketAddress.ProtocolStr == "" {
			l.Address.SocketAddress.ProtocolStr = "HTTP"
		}
		l.Address.SocketAddress.Protocol = model.ProtocolType(model.ProtocolTypeValue[l.Address.SocketAddress.ProtocolStr])
	}
}

func HttpConfig(cfg *model.Bootstrap) {
	for i, l := range cfg.StaticResources.Listeners {
		hc := &model.HttpConfig{}
		if l.Config != nil {
			if v, ok := l.Config.(map[string]interface{}); ok {
				switch l.Name {
				case "net/http":
					if err := mapstructure.Decode(v, hc); err != nil {
						logger.Error(err)
					}
					cfg.StaticResources.Listeners[i].Config = hc
				}
			}
		}
	}
}

func FilterChain(cfg *model.Bootstrap) {
	for _, l := range cfg.StaticResources.Listeners {
		for _, fc := range l.FilterChains {
			if fc.Filters != nil {
				for i, fcf := range fc.Filters {
					hcm := &model.HttpConnectionManager{}
					if fcf.Config != nil {
						switch fcf.Name {
						case "dgp.filters.http_connect_manager":
							if v, ok := fcf.Config.(map[string]interface{}); ok {
								if err := mapstructure.Decode(v, hcm); err != nil {
									logger.Error(err)
								}
								fc.Filters[i].Config = hcm
							}
						}
					}
				}
			}
		}
	}
}

func LoadBalance(cfg *model.Bootstrap) {
	for _, c := range cfg.StaticResources.Clusters {
		var lbPolicy int32
		if c.LbStr != "" {
			if lb, ok := model.LbPolicyValue[c.LbStr]; ok {
				lbPolicy = lb
			} else {
				c.LbStr = "RoundRobin"
				lbPolicy = model.LbPolicyValue[c.LbStr]
			}
		} else {
			c.LbStr = "RoundRobin"
			lbPolicy = model.LbPolicyValue[c.LbStr]
		}
		c.Lb = model.LbPolicy(lbPolicy)
	}
}

func DiscoverType(cfg *model.Bootstrap) {
	var discoverType int32
	for _, c := range cfg.StaticResources.Clusters {
		if c.TypeStr != "" {
			if t, ok := model.DiscoveryTypeValue[c.TypeStr]; ok {
				discoverType = t
			} else {
				c.TypeStr = "EDS"
				discoverType = model.DiscoveryTypeValue[c.TypeStr]
			}
		} else {
			c.TypeStr = "EDS"
			discoverType = model.DiscoveryTypeValue[c.TypeStr]
		}
		c.Type = model.DiscoveryType(discoverType)
	}
}

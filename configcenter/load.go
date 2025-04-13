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
	"strings"
)

import (
	"gopkg.in/yaml.v3"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	KEY_CONFIG_TYPE_NACOS = "nacos"
)

var Parsers = map[string]func(data []byte, v interface{}) error{
	".yaml": ParseYamlBytes,
	".yml":  ParseYamlBytes,
}

type (
	Load interface {
		LoadConfigs(boot *model.Bootstrap, opts ...Option) (v *model.Bootstrap, err error)
		ViewRemoteConfig() *model.Bootstrap
	}

	Option func(opt *Options)

	Options struct {
		Remote bool
		DataId string
		Group  string
		path   string // nolint:unused
		parser string // nolint:unused
	}
)

type DefaultConfigLoad struct {
	bootConfig   *model.Bootstrap
	configClient ConfigClient
}

func NewConfigLoad(bootConfig *model.Bootstrap) *DefaultConfigLoad {
	var configClient ConfigClient
	var err error
	// config center load
	if strings.EqualFold(bootConfig.Config.Type, KEY_CONFIG_TYPE_NACOS) {
		configClient, err = NewNacosConfig(bootConfig)
		if err != nil {
			logger.Errorf("Get new nacos config failed,err: %v", err)
		}
	}

	if configClient == nil {
		logger.Warnf("no remote config-center")
		return nil
	}

	return &DefaultConfigLoad{
		bootConfig:   bootConfig,
		configClient: configClient,
	}
}

func (d *DefaultConfigLoad) LoadConfigs(boot *model.Bootstrap, opts ...Option) (v *model.Bootstrap, err error) {
	var opt Options
	for _, o := range opts {
		o(&opt)
	}

	if !opt.Remote {
		return nil, nil
	}

	m := map[string]interface{}{}

	if strings.EqualFold(boot.Config.Type, KEY_CONFIG_TYPE_NACOS) {
		m["dataId"] = getOrDefault(opt.DataId, DataId)
		m["group"] = getOrDefault(opt.Group, Group)
	}

	if len(m) == 0 {
		logger.Errorf("no identify properties key when load from remote config center")
		return boot, nil
	}

	data, err := d.configClient.LoadConfig(m)
	if err != nil {
		return nil, err
	}

	if len(data) == 0 {
		logger.Errorf("the config data load from remote is nil, config center : %s", boot.Config.Type)
		return boot, err
	}

	if err = Parsers[".yml"]([]byte(data), boot); err != nil {
		logger.Errorf("failed to parse the configuration loaded from the remote,err: %v", err)
		return boot, err
	}

	if err = d.configClient.ListenConfig(m); err != nil {
		logger.Errorf("failed to listen the remote configcenter config,err: %v", err)
		return boot, err
	}

	return boot, err
}

// ViewRemoteConfig returns the current remote configuration.
func (d *DefaultConfigLoad) ViewRemoteConfig() *model.Bootstrap {
	return d.configClient.ViewConfig()
}

func ParseYamlBytes(content []byte, v interface{}) error {
	return yaml.Unmarshal(content, v)
}

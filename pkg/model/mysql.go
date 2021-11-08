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

package model

import (
	"github.com/mitchellh/mapstructure"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

type MysqlConfig struct {
	Salt          string            `yaml:"salt" json:"salt" mapstructure:"salt"`
	Users         map[string]string `yaml:"users" json:"users" mapstructure:"users"`
	ServerVersion string            `yaml:"server_version" json:"server_version" mapstructure:"server_version"`
}

func MapToMysqlConfig(cfg interface{}) *MysqlConfig {
	var hc *MysqlConfig
	if cfg != nil {
		if err := mapstructure.Decode(cfg, &hc); err != nil {
			logger.Error("Config error", err)
		}
	}
	return hc
}

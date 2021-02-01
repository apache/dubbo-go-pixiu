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
	fc "github.com/dubbogo/dubbo-go-proxy-filter/pkg/api/config"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/yaml"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
	"sync"
)

var (
	apiConfig *fc.APIConfig
	once      sync.Once
)

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
	once.Do(func() {
		apiConfig = apiConf
	})
	return apiConf, nil
}

// GetAPIConf returns the initted api config
func GetAPIConf() fc.APIConfig {
	return *apiConfig
}

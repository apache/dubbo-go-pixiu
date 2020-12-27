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
package plugins

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"github.com/stretchr/testify/assert"
	"testing"
)

func TestInitPluginsGroup(t *testing.T) {

	config, err := config.LoadAPIConfigFromFile("/Users/zhengxianle/dubbo-go-proxy/configs/api_config.yaml")
	assert.Empty(t, err)

	InitPluginsGroup(config.PluginsGroup, config.PluginFilePath)
}

func TestInitApiUrlWithFilterChain(t *testing.T) {
	config, err := config.LoadAPIConfigFromFile("/Users/zhengxianle/dubbo-go-proxy/configs/api_config.yaml")
	assert.Empty(t, err)

	InitPluginsGroup(config.PluginsGroup, config.PluginFilePath)
	InitApiUrlWithFilterChain(config.Resources)
}

func TestGetApiFilterFuncsWithApiUrl(t *testing.T) {
	config, err := config.LoadAPIConfigFromFile("/Users/zhengxianle/dubbo-go-proxy/configs/api_config.yaml")
	assert.Empty(t, err)

	InitPluginsGroup(config.PluginsGroup, config.PluginFilePath)
	InitApiUrlWithFilterChain(config.Resources)

	fltc := GetApiFilterFuncsWithApiUrl("/")

	assert.Equal(t, len(fltc), 3)
}

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

package config_test

import (
	_ "github.com/apache/dubbo-go/config_center/zookeeper"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
	"log"
	"testing"
)

import (
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
)

func TestLoadAPIConfigFromFile(t *testing.T) {
	apiC, err := config.LoadAPIConfigFromFile("")
	assert.Empty(t, apiC)
	assert.EqualError(t, err, "Config file not specified")

	apiC, err = config.LoadAPIConfigFromFile("./mock/api_config.yml")
	assert.Empty(t, err)
	assert.Equal(t, apiC.Name, "api name")
	bytes, _ := yaml.Marshal(apiC)
	log.Printf("%s", bytes)
}

func TestLoadAPIConfig(t *testing.T) {
	apiMeataConf := &model.ApiMetaConfig{"zookeeper", "127.0.0.1:2181", "proxy", "dubbo", "dubbo", "dubbo"}

	apiC, err := config.LoadAPIConfig(apiMeataConf)
	assert.Empty(t, err)
	assert.Equal(t, apiC.Name, "api name")
	bytes, _ := yaml.Marshal(apiC)
	log.Printf("%s", bytes)
}

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
	"log"
	"regexp"
	"testing"
)

import (
	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config"
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

func TestRegexSucc(t *testing.T) {
	et := regexp.MustCompile(".+/Resources/[^/]+/?$")
	ret := et.Match([]byte("/dis/sdf/Resources/service1/"))
	assert.True(t, ret, "not match")
}

func TestRegexExtractSucc(t *testing.T) {
	re := regexp.MustCompile("Resources/([^/]+)/Method")
	result := re.FindStringSubmatch("/Resources/user-service/Method/POST")
	assert.Equal(t, result[1], "user-service")
}

func TestRegexFail(t *testing.T) {
	et := regexp.MustCompile(".+/Resources/[^/]+/?$")
	ret := et.Match([]byte("/dis/sdf/Resources/service1/ddd"))
	assert.False(t, ret, "not match")
}

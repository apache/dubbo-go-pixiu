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
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
)

func TestLoadAPIConfigFromFile(t *testing.T) {
	apiC, err := config.LoadAPIConfigFromFile("")
	assert.Empty(t, apiC)
	assert.EqualError(t, err, "Config file not specified")

	apiC, err = config.LoadAPIConfigFromFile("./mock/api_config.yml")
	assert.Empty(t, err)
	assert.Equal(t, apiC.Name, "api name")
	bytes, _ := yaml.MarshalYML(apiC)
	log.Printf("%s", bytes)
}

func TestAPIConfigLoadsDirectGenericContractFields(t *testing.T) {
	raw := []byte(`
name: pixiu
resources:
  - path: /invoke
    type: restful
    methods:
      - httpVerb: POST
        enable: true
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: triple
          url: tri://127.0.0.1:20000
          protocol: tri
          interface: benchmark.BenchmarkService
          method: SayHello
          parameterTypes:
            - benchmark.HelloRequest
          serialization: protobuf
`)
	var cfg config.APIConfig
	err := yaml.UnmarshalYML(raw, &cfg)
	require.NoError(t, err)

	ir := cfg.Resources[0].Methods[0].IntegrationRequest
	require.Equal(t, []string{"benchmark.HelloRequest"}, ir.ParameterTypes)
	require.Equal(t, "protobuf", ir.Serialization)
}

func TestAPIConfigKeepsExplicitEmptyParameterTypes(t *testing.T) {
	raw := []byte(`
name: pixiu
resources:
  - path: /ping
    type: restful
    methods:
      - httpVerb: POST
        enable: true
        inboundRequest:
          requestType: http
        integrationRequest:
          requestType: dubbo
          url: dubbo://127.0.0.1:20880
          protocol: dubbo
          interface: demo.PingService
          method: Ping
          parameterTypes: []
          serialization: hessian2
`)
	var cfg config.APIConfig
	err := yaml.UnmarshalYML(raw, &cfg)
	require.NoError(t, err)

	ir := cfg.Resources[0].Methods[0].IntegrationRequest
	require.NotNil(t, ir.ParameterTypes)
	require.Len(t, ir.ParameterTypes, 0)
	require.Equal(t, "hessian2", ir.Serialization)
}

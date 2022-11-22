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

package circuitbreaker

import (
	stdHttp "net/http"
	"testing"
)

import (
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/mock"
	pkgs "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/sentinel"
)

func TestFilter(t *testing.T) {
	f := FilterFactory{cfg: &Config{}}

	mockYaml, err := yaml.MarshalYML(mockConfig())
	assert.Nil(t, err)

	assert.Nil(t, yaml.UnmarshalYML(mockYaml, f.Config()))

	assert.Nil(t, f.Apply())

	decoder := &Filter{cfg: f.cfg, matcher: f.matcher}
	request, _ := stdHttp.NewRequest(stdHttp.MethodGet, "https://www.dubbogopixiu.com/api/v1/test-dubbo/user/1111", nil)
	c := mock.GetMockHTTPContext(request)

	assert.Equal(t, decoder.Decode(c), filter.Continue)
}

func mockConfig() *Config {
	c := Config{
		Resources: []*pkgs.Resource{
			{
				Name: "test-dubbo",
				Items: []*pkgs.Item{
					{MatchStrategy: pkgs.EXACT, Pattern: "/api/v1/test-dubbo/user"},
					{MatchStrategy: pkgs.REGEX, Pattern: "/api/v1/test-dubbo/user/*"},
				},
			},
		},
		Rules: []*circuitbreaker.Rule{{
			Resource:         "test-dubbo",
			Strategy:         circuitbreaker.ErrorCount,
			RetryTimeoutMs:   3000,
			MinRequestAmount: 10,
			StatIntervalMs:   1000,
			Threshold:        1.0,
		}},
	}
	return &c
}

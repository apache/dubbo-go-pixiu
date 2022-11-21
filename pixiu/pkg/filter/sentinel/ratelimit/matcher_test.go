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

package ratelimit

import (
	"testing"
)

import (
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/sentinel"
)

func TestMatch(t *testing.T) {
	config := mockConfig()
	m := sentinel.NewMatcher()
	m.Load(config.Resources)

	tests := []struct {
		give    string
		matched bool
		res     string
	}{
		{
			give:    "/api/v1/test-dubbo/user",
			matched: true,
			res:     "test-dubbo",
		},
		{
			give:    "/api/v1/test-dubbo/user/getInfo",
			matched: true,
			res:     "test-dubbo",
		},
		{
			give:    "/api/v1/test-dubbo/A",
			matched: false,
			res:     "",
		},
		{
			give:    "/api/v1/http/foo",
			matched: true,
			res:     "test-http",
		},
		{
			give:    "/api/v1/http/bar/1.json",
			matched: true,
			res:     "test-http",
		},
		{
			give:    "/api/v1/http",
			matched: false,
			res:     "",
		},
	}
	for _, tt := range tests {
		t.Run(tt.give, func(t *testing.T) {
			res, ok := m.Match(tt.give)
			assert.Equal(t, res, tt.res)
			assert.Equal(t, ok, tt.matched)
		})
	}
}

func mockConfig() *Config {
	c := Config{
		Resources: []*sentinel.Resource{
			{
				Name: "test-dubbo",
				Items: []*sentinel.Item{
					{MatchStrategy: sentinel.EXACT, Pattern: "/api/v1/test-dubbo/user"},
					{MatchStrategy: sentinel.REGEX, Pattern: "/api/v1/test-dubbo/user/*"},
				},
			},
			{
				Name: "test-http",
				Items: []*sentinel.Item{
					{MatchStrategy: sentinel.EXACT, Pattern: "/api/v1/http/foo"},
					{MatchStrategy: sentinel.EXACT, Pattern: "/api/v1/http/bar"},

					{MatchStrategy: sentinel.REGEX, Pattern: "/api/v1/http/foo/*"},
					{MatchStrategy: sentinel.REGEX, Pattern: "/api/v1/http/bar/*"},
				},
			},
		},
		Rules: []*Rule{
			{
				Enable: true,
				FlowRule: flow.Rule{
					Threshold:        100,
					StatIntervalInMs: 1000,
				},
			},
		},
	}
	return &c
}

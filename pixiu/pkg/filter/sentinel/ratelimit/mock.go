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
	"github.com/alibaba/sentinel-golang/core/flow"
)

import (
	pkgs "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/sentinel"
)

func GetMockedRateLimitConfig() *Config {
	c := Config{
		Resources: []*pkgs.Resource{
			{
				Name: "test-dubbo",
				Items: []*pkgs.Item{
					{MatchStrategy: pkgs.EXACT, Pattern: "/api/v1/test-dubbo/user"},
					{MatchStrategy: pkgs.REGEX, Pattern: "/api/v1/test-dubbo/user/*"},
				},
			},
			{
				Name: "test-http",
				Items: []*pkgs.Item{
					{MatchStrategy: pkgs.EXACT, Pattern: "/api/v1/http/foo"},
					{MatchStrategy: pkgs.EXACT, Pattern: "/api/v1/http/bar"},

					{MatchStrategy: pkgs.REGEX, Pattern: "/api/v1/http/foo/*"},
					{MatchStrategy: pkgs.REGEX, Pattern: "/api/v1/http/bar/*"},
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

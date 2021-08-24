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

	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config/ratelimit"
)

func GetMockedRateLimitConfig() *ratelimit.Config {
	c := ratelimit.Config{
		Resources: []ratelimit.Resource{
			{
				Name: "test-dubbo",
				Items: []ratelimit.Item{
					{MatchStrategy: ratelimit.EXACT, Pattern: "/api/v1/test-dubbo/user"},
					{MatchStrategy: ratelimit.REGEX, Pattern: "/api/v1/test-dubbo/user/*"},
				},
			},
			{
				Name: "test-http",
				Items: []ratelimit.Item{
					{MatchStrategy: ratelimit.EXACT, Pattern: "/api/v1/http/foo"},
					{MatchStrategy: ratelimit.EXACT, Pattern: "/api/v1/http/bar"},

					{MatchStrategy: ratelimit.REGEX, Pattern: "/api/v1/http/foo/*"},
					{MatchStrategy: ratelimit.REGEX, Pattern: "/api/v1/http/bar/*"},
				},
			},
		},
		Rules: []ratelimit.Rule{
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

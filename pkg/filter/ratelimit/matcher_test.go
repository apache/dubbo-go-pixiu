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
	"github.com/apache/dubbo-go-pixiu/pkg/filter/ratelimit/matcher"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"testing"
)

func TestMatch(t *testing.T) {
	config, err := GetMockedRateLimitConfig()
	assert.Nil(t, err)

	matcher.Load(config.APIResources)

	res, _ := matcher.Match("/api/v1/test-dubbo/user")
	assert.Equal(t, "test-dubbo", res)
	res, _ = matcher.Match("/api/v1/test-dubbo/user/getInfo")
	assert.Equal(t, "test-dubbo", res)
	_, ok := matcher.Match("/api/v1/test-dubbo/A")
	assert.Equal(t, false, ok)

	res, _ = matcher.Match("/api/v1/http/foo")
	assert.Equal(t, "test-http", res)
	res, _ = matcher.Match("/api/v1/http/bar/1.json")
	assert.Equal(t, "test-http", res)
}

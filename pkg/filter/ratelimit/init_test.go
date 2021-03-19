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
	"encoding/json"
	"testing"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

func TestGetMockConfig(t *testing.T) {
	config, err := GetMockedRateLimitConfig()
	if err != nil {
		logger.Warn(err)
		t.Fail()
	}
	configStr, _ := json.Marshal(config)
	logger.Infof("rate limit config %s", configStr)
}

func TestInit(t *testing.T) {
	err := rateLimitInit()
	if err != nil {
		t.Fatal(err)
	}
}

func TestLoadRules(t *testing.T) {
	config, _ := GetMockedRateLimitConfig()
	loadRules(config.Rules)
}

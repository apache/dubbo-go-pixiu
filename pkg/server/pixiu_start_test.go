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

package server

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestStart_ReachesBlockPoint(t *testing.T) {
	// Replace blockForever to prevent the test from hanging.
	original := blockForever
	blockForever = func() {}
	defer func() { blockForever = original }()

	bs := &model.Bootstrap{}
	config.SetBootstrap(bs)
	defer config.SetBootstrap(nil)

	// Reset global server after test to avoid polluting other tests.
	defer func() { server = nil }()

	Start(bs)

	s := GetServer()
	assert.NotNil(t, s)
	assert.NotNil(t, s.GetListenerManager())
	assert.NotNil(t, s.GetClusterManager())
	assert.NotNil(t, s.GetRouterManager())
	assert.NotNil(t, s.GetApiConfigManager())
	assert.NotNil(t, s.GetTraceDriverManager())
}

func TestServerStart_WithMinimalConfig(t *testing.T) {
	bs := &model.Bootstrap{}
	config.SetBootstrap(bs)
	defer config.SetBootstrap(nil)

	s := NewServer()
	s.initialize(bs)
	s.Start()

	// Verify Start() returned without panicking.
	assert.NotNil(t, s.GetListenerManager())
	assert.NotNil(t, s.GetClusterManager())
}

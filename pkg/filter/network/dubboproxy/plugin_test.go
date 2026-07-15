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

package dubboproxy

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestPlugin_Kind(t *testing.T) {
	p := &Plugin{}
	assert.Equal(t, "dgp.filter.network.dubboconnectionmanager", p.Kind())
}

func TestPlugin_Config(t *testing.T) {
	p := &Plugin{}
	cfg := p.Config()
	_, ok := cfg.(*model.DubboProxyConnectionManagerConfig)
	assert.True(t, ok, "Config should return *model.DubboProxyConnectionManagerConfig")
}

func TestPlugin_CreateFilter_Success(t *testing.T) {
	p := &Plugin{}
	cfg := &model.DubboProxyConnectionManagerConfig{
		RouteConfig: model.RouteConfiguration{},
		TimeoutStr:  "3s",
	}

	filter, err := p.CreateFilter(cfg)
	assert.NoError(t, err)
	assert.NotNil(t, filter)

	// Verify timeout was set
	assert.Greater(t, cfg.Timeout.Seconds(), float64(0))
}

func TestPlugin_CreateFilter_InvalidConfigType(t *testing.T) {
	p := &Plugin{}

	// Test with nil config
	filter, err := p.CreateFilter(nil)
	assert.Error(t, err)
	assert.Nil(t, filter)
	assert.Contains(t, err.Error(), "invalid config type")
	assert.Contains(t, err.Error(), "expected *model.DubboProxyConnectionManagerConfig")

	// Test with wrong type - string
	filter, err = p.CreateFilter("invalid")
	assert.Error(t, err)
	assert.Nil(t, filter)
	assert.Contains(t, err.Error(), "invalid config type")

	// Test with wrong type - map
	filter, err = p.CreateFilter(map[string]interface{}{"key": "value"})
	assert.Error(t, err)
	assert.Nil(t, filter)
	assert.Contains(t, err.Error(), "invalid config type")

	// Test with wrong pointer type
	filter, err = p.CreateFilter(&model.HttpConnectionManagerConfig{})
	assert.Error(t, err)
	assert.Nil(t, filter)
	assert.Contains(t, err.Error(), "invalid config type")
}

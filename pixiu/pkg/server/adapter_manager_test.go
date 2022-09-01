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
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/adapter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

type (
	DemoAdapterPlugin struct{}

	DemoAdapter struct {
		cfg *Config
	}

	Config struct {
		Name string `yaml:"name" json:"name,omitempty"`
	}
)

func (d *DemoAdapterPlugin) Kind() string {
	return "test"
}

func (p *DemoAdapterPlugin) CreateAdapter(ad *model.Adapter) (adapter.Adapter, error) {
	return &DemoAdapter{cfg: &Config{}}, nil
}

func (a *DemoAdapter) Start() {

}

// Apply init
func (a *DemoAdapter) Apply() error {
	return nil
}

// Config get config for Adapter
func (a *DemoAdapter) Config() interface{} {
	return a.cfg
}

func (a *DemoAdapter) Stop() {

}

func TestAdapterManager(t *testing.T) {
	adapter.RegisterAdapterPlugin(&DemoAdapterPlugin{})

	bs := &model.Bootstrap{
		StaticResources: model.StaticResources{
			Adapters: []*model.Adapter{
				&model.Adapter{
					Name:   "test",
					Config: make(map[string]interface{}),
				},
			},
		},
	}

	am := CreateDefaultAdapterManager(nil, bs)

	assert.Equal(t, len(am.adapters), 1)

	am.Start()

	am.Stop()

}

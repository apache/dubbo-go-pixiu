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

package adapter

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

type (
	DemoAdapterPlugin struct{}

	DemoAdapter struct {
		cfg *Config
	}

	Config struct{}
)

func (d *DemoAdapterPlugin) Kind() string {
	return "test"
}

func (p *DemoAdapterPlugin) CreateAdapter(ad *model.Adapter) (Adapter, error) {
	return &DemoAdapter{cfg: &Config{}}, nil
}

// Apply init
func (a *DemoAdapter) Apply() error {
	return nil
}

// Config get config for Adapter
func (a *DemoAdapter) Config() interface{} {
	return a.cfg
}

func (a *DemoAdapter) Start() {

}

func (a *DemoAdapter) Stop() {

}

func TestRegisterAdapterPlugin(t *testing.T) {
	RegisterAdapterPlugin(&DemoAdapterPlugin{})
	_, err := GetAdapterPlugin("test")
	assert.NoError(t, err)
}

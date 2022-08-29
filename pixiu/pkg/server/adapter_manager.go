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
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/adapter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

type AdapterManager struct {
	configs  []*model.Adapter
	adapters []adapter.Adapter
}

func CreateDefaultAdapterManager(server *Server, bs *model.Bootstrap) *AdapterManager {
	am := &AdapterManager{configs: bs.StaticResources.Adapters}
	am.initAdapters()
	return am
}

func (am *AdapterManager) Start() {
	for _, a := range am.adapters {
		a.Start()
	}
}

func (am *AdapterManager) Stop() {
	for _, a := range am.adapters {
		a.Stop()
	}
}

func (am *AdapterManager) initAdapters() {

	var ads []adapter.Adapter

	for _, f := range am.configs {
		hp, err := adapter.GetAdapterPlugin(f.Name)
		if err != nil {
			logger.Error("initAdapters get plugin error %s", err)
		}

		if len(f.Enabled) > 0 && !strings.EqualFold(f.Enabled, "true") {
			logger.Warnf("the Adapter %s will stop starting, by config of Enabled : %s", f.Name, f.Enabled)
			return
		}

		hf, err := hp.CreateAdapter(f)
		if err != nil {
			logger.Error("initFilterIfNeed create adapter error %s", err)
		}

		cfg := hf.Config()
		if err := yaml.ParseConfig(cfg, f.Config); err != nil {
			logger.Error("initAdapters init config error %s", err)
		}

		err = hf.Apply()
		if err != nil {
			logger.Error("initFilterIfNeed apply adapter error %s", err)
		}
		ads = append(ads, hf)
	}
	am.adapters = ads
}

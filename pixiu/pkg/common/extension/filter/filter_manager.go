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

package filter

import (
	"sync"
)

import (
	"github.com/creasty/defaults"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

// FilterManager manage filters
type FilterManager struct {
	filters       map[string]HttpFilterFactory
	filtersArray  []*HttpFilterFactory
	filterConfigs []*model.HTTPFilter

	mu sync.RWMutex
}

// NewFilterManager create filter manager
func NewFilterManager(fs []*model.HTTPFilter) *FilterManager {
	fm := &FilterManager{filterConfigs: fs, filters: make(map[string]HttpFilterFactory)}
	return fm
}

// NewEmptyFilterManager create empty filter manager
func NewEmptyFilterManager() *FilterManager {
	return &FilterManager{filters: make(map[string]HttpFilterFactory)}
}

func (fm *FilterManager) CreateFilterChain(ctx *http.HttpContext) FilterChain {
	chain := NewDefaultFilterChain()

	for _, f := range fm.GetFactory() {
		_ = (*f).PrepareFilterChain(ctx, chain)
	}
	return chain
}

// GetFactory get all filter from manager
func (fm *FilterManager) GetFactory() []*HttpFilterFactory {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	return fm.filtersArray
}

// Load the filter from config
func (fm *FilterManager) Load() {
	fm.ReLoad(fm.filterConfigs)
}

// ReLoad filter configs
func (fm *FilterManager) ReLoad(filters []*model.HTTPFilter) {
	tmp := make(map[string]HttpFilterFactory)
	filtersArray := make([]*HttpFilterFactory, len(filters))
	for i, f := range filters {
		apply, err := fm.Apply(f.Name, f.Config)
		if err != nil {
			logger.Errorf("apply [%s] init fail, %s", f.Name, err.Error())
		}
		tmp[f.Name] = apply
		filtersArray[i] = &apply
	}
	// avoid filter inconsistency
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fm.filters = tmp
	fm.filtersArray = filtersArray
}

// Apply return a new filter factory by name & conf
func (fm *FilterManager) Apply(name string, conf map[string]interface{}) (HttpFilterFactory, error) {
	plugin, err := GetHttpFilterPlugin(name)
	if err != nil {
		return nil, errors.New("filter not found")
	}

	filter, err := plugin.CreateFilterFactory()

	if err != nil {
		return nil, errors.New("plugin create filter error")
	}

	factoryConf := filter.Config()
	if err := yaml.ParseConfig(factoryConf, conf); err != nil {
		return nil, errors.Wrap(err, "config error")
	}
	if err = defaults.Set(factoryConf); err != nil {
		return nil, err
	}
	err = filter.Apply()
	if err != nil {
		return nil, errors.Wrap(err, "create fail")
	}
	return filter, nil
}

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
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
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

// CreateFilterChain preserves the original best-effort public API: factories
// that fail are logged and skipped, while filters prepared successfully by
// other factories remain in the returned chain. Configuration publication uses
// CreateFilterChainChecked so errors can reject the complete resource.
func (fm *FilterManager) CreateFilterChain(ctx *http.HttpContext) FilterChain {
	chain := NewDefaultFilterChain()
	for index, f := range fm.GetFactory() {
		if f == nil || *f == nil {
			logger.Errorf("create HTTP filter chain: HTTP filter factory %d is nil", index)
			continue
		}
		if err := (*f).PrepareFilterChain(ctx, chain); err != nil {
			logger.Errorf("create HTTP filter chain: prepare HTTP filter %d: %v", index, err)
		}
	}
	return chain
}

func (fm *FilterManager) CreateFilterChainChecked(ctx *http.HttpContext) (FilterChain, error) {
	chain := NewDefaultFilterChain()

	for index, f := range fm.GetFactory() {
		if f == nil || *f == nil {
			return nil, errors.Errorf("HTTP filter factory %d is nil", index)
		}
		if err := (*f).PrepareFilterChain(ctx, chain); err != nil {
			return nil, errors.Wrapf(err, "prepare HTTP filter %d", index)
		}
	}
	return chain, nil
}

// GetFactory get all filter from manager
func (fm *FilterManager) GetFactory() []*HttpFilterFactory {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	return fm.filtersArray
}

// Load the filter from config
func (fm *FilterManager) Load() {
	if err := fm.LoadChecked(); err != nil {
		logger.Errorf("load HTTP filters: %v", err)
	}
}

func (fm *FilterManager) LoadChecked() error {
	return fm.ReLoadChecked(fm.filterConfigs)
}

// ReLoad filter configs
func (fm *FilterManager) ReLoad(filters []*model.HTTPFilter) {
	if err := fm.ReLoadChecked(filters); err != nil {
		logger.Errorf("reload HTTP filters: %v", err)
	}
}

func (fm *FilterManager) ReLoadChecked(filters []*model.HTTPFilter) error {
	tmp := make(map[string]HttpFilterFactory)
	filtersArray := make([]*HttpFilterFactory, len(filters))
	for i, f := range filters {
		if f == nil || f.Name == "" {
			return errors.Errorf("HTTP filter %d has an empty name", i)
		}
		apply, err := fm.Apply(f.Name, f.Config)
		if err != nil {
			return errors.Wrapf(err, "apply HTTP filter %q", f.Name)
		}
		tmp[f.Name] = apply
		filtersArray[i] = &apply
	}
	// avoid filter inconsistency
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fm.filters = tmp
	fm.filtersArray = filtersArray
	return nil
}

// Apply return a new filter factory by name & conf
func (fm *FilterManager) Apply(name string, conf map[string]any) (HttpFilterFactory, error) {
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

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

	mu          sync.RWMutex
	lifecycleMu sync.Mutex
	factoryRefs map[*HttpFilterFactory]int
	retired     map[*HttpFilterFactory]bool
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

	fm.mu.RLock()
	defer fm.mu.RUnlock()
	factories := append([]*HttpFilterFactory(nil), fm.filtersArray...)
	fm.leaseFactories(factories)
	chain.(*defaultFilterChain).setRelease(func() {
		fm.releaseFactories(factories)
	})
	for _, f := range factories {
		_ = (*f).PrepareFilterChain(ctx, chain)
	}
	return chain
}

// GetFactory get all filter from manager
func (fm *FilterManager) GetFactory() []*HttpFilterFactory {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	return append([]*HttpFilterFactory(nil), fm.filtersArray...)
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
	oldFilters := fm.filtersArray
	fm.filters = tmp
	fm.filtersArray = filtersArray
	ready := fm.retireFactories(oldFilters)
	fm.mu.Unlock()
	fm.closeAndLog(ready)
}

func (fm *FilterManager) leaseFactories(factories []*HttpFilterFactory) {
	fm.lifecycleMu.Lock()
	defer fm.lifecycleMu.Unlock()
	if fm.factoryRefs == nil {
		fm.factoryRefs = make(map[*HttpFilterFactory]int)
	}
	for _, factory := range factories {
		if factory != nil {
			fm.factoryRefs[factory]++
		}
	}
}

func (fm *FilterManager) retireFactories(factories []*HttpFilterFactory) []*HttpFilterFactory {
	fm.lifecycleMu.Lock()
	defer fm.lifecycleMu.Unlock()
	if fm.retired == nil {
		fm.retired = make(map[*HttpFilterFactory]bool)
	}
	var ready []*HttpFilterFactory
	for _, factory := range factories {
		if factory == nil {
			continue
		}
		fm.retired[factory] = true
		if fm.factoryRefs[factory] == 0 {
			delete(fm.retired, factory)
			ready = append(ready, factory)
		}
	}
	return ready
}

func (fm *FilterManager) releaseFactories(factories []*HttpFilterFactory) {
	fm.lifecycleMu.Lock()
	var ready []*HttpFilterFactory
	for _, factory := range factories {
		if factory == nil {
			continue
		}
		fm.factoryRefs[factory]--
		if fm.factoryRefs[factory] == 0 {
			delete(fm.factoryRefs, factory)
			if fm.retired[factory] {
				delete(fm.retired, factory)
				ready = append(ready, factory)
			}
		}
	}
	fm.lifecycleMu.Unlock()
	fm.closeAndLog(ready)
}

func (fm *FilterManager) closeAndLog(factories []*HttpFilterFactory) {
	if err := fm.closeFactories(factories); err != nil {
		logger.Warnf("failed to close retired HTTP filter factory: %v", err)
	}
}

func (fm *FilterManager) closeFactories(factories []*HttpFilterFactory) error {
	var firstErr error
	for _, factory := range factories {
		if factory == nil || *factory == nil {
			continue
		}
		closer, ok := (*factory).(interface{ Close() error })
		if !ok {
			continue
		}
		if err := closer.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
	}
	return firstErr
}

// Close releases resources held by HTTP filter factories that expose an
// optional Close method. Keeping this optional preserves compatibility with
// existing HTTP filter implementations.
func (fm *FilterManager) Close() error {
	fm.mu.Lock()
	factories := fm.filtersArray
	fm.filtersArray = nil
	ready := fm.retireFactories(factories)
	fm.mu.Unlock()
	return fm.closeFactories(ready)
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

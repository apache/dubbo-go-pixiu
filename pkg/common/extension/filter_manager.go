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
 * WITHOmanage.goUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package extension

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/filter"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
	"sync"
)

// FilterManager take over the entire life cycle of FilterChain,
// All filters are configurable and support dynamic refresh

// filterFactories all filter will be registered in filterFactories
// by call RegisterFilterFactory
var filterFactories = make(map[string]filterFactoryCreator, 16)

type filterFactoryCreator func() filter.Factory
type FilterChain []filter.Filter

type filterManager struct {
	filters FilterChain

	mu sync.RWMutex
}

func NewFilterManager() *filterManager {
	return &filterManager{filters: make([]filter.Filter, 0, 16)}
}

// RegisterFilterFactory will store the @filter factory and @name
func RegisterFilterFactory(name string, filter filterFactoryCreator) {
	filterFactories[name] = filter
}

func (fm *filterManager) GetFilters() []filter.Filter {
	fm.mu.RLock()
	defer fm.mu.RUnlock()

	return fm.filters
}

// Load init or reload filter configs
func (fm *filterManager) Load(filters []*config.Filter) {
	tmp := make([]filter.Filter, 0, len(filters))
	for _, f := range filters {
		apply, err := fm.Apply(f.Name, f.Config)
		if err != nil {
			logger.Errorf("apply [%s] init fail, %s", err)
		}
		tmp = append(tmp, apply)
	}
	// avoid filter inconsistency
	fm.mu.Lock()
	defer fm.mu.Unlock()

	fm.filters = tmp
}

// Apply return a new filter by name & conf
func (fm *filterManager) Apply(name string, conf map[string]interface{}) (filter.Filter, error) {
	creator, ok := filterFactories[name]
	if !ok {
		return nil, errors.New("filter not found")
	}
	factory := creator()
	factoryConf := factory.Config()
	if err := parseConfig(factoryConf, conf); err != nil {
		return nil, errors.Wrap(err, "config error")
	}
	apply, err := factory.Apply()
	if err != nil {
		return nil, errors.Wrap(err, "create fail")
	}
	return apply, nil
}

func parseConfig(factoryConfStruct interface{}, conf map[string]interface{}) error {
	// conf will be map, convert to yaml
	yamlBytes, err := yaml.Marshal(conf)
	if err != nil {
		return err
	}
	// Unmarshal yamlStr to factoryConf
	err = yaml.Unmarshal(yamlBytes, factoryConfStruct)
	if err != nil {
		return err
	}
	return nil
}

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

package extension

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/context"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/filter"
	"github.com/ghodss/yaml"
	"github.com/pkg/errors"
)

var filterFactories = make(map[string]filterFactoryCreator, 8)

type FilterManager struct {
	filters []context.Filter
}

type filterFactoryCreator func() filter.Factory

// RegisterFilterFactory will store the @filter factory and @name
func RegisterFilterFactory(name string, filter filterFactoryCreator) {
	filterFactories[name] = filter
}

func (fm *FilterManager) GetFilters() []context.Filter {
	return fm.filters
}

func (fm *FilterManager) Init(filters []config.Filter) {
	tmp := make([]context.Filter, len(filters))
	for _, f := range filters {
		apply, err := fm.GetFilter(f.Name, f.Config)
		if err != nil {
			logger.Errorf("apply [%s] init fail, %s", err)
		}
		tmp = append(tmp, apply)
	}
	//TODO lock and hotswap
	fm.filters = tmp
}

func (fm FilterManager) GetFilter(name string, conf interface{}) (context.Filter, error) {
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

func parseConfig(factoryConfStruct interface{}, conf interface{}) error {
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

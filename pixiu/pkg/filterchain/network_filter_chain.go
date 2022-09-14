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

package filterchain

import (
	"context"
	"net/http"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

type NetworkFilterChain struct {
	filtersArray []filter.NetworkFilter
	config       model.FilterChain
}

// ServeHTTP handle http request
func (fc NetworkFilterChain) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// todo: only one filter will exist for now, needs change when more than one
	for _, filter := range fc.filtersArray {
		filter.ServeHTTP(w, r)
	}
}

// OnDecode decode bytes received from getty listener
func (fc NetworkFilterChain) OnDecode(data []byte) (interface{}, int, error) {
	// todo: only one filter will exist for now, needs change when more than one
	for _, filter := range fc.filtersArray {
		return filter.OnDecode(data)
	}
	return nil, 0, errors.Errorf("filterChain don't have network filter")
}

// OnEncode encode struct to bytes sent to getty listener
func (fc NetworkFilterChain) OnEncode(p interface{}) ([]byte, error) {
	// todo: only one filter will exist for now, needs change when more than one
	for _, filter := range fc.filtersArray {
		return filter.OnEncode(p)
	}
	return nil, errors.Errorf("filterChain don't have network filter")
}

// OnData handle dubbo rpc invocation
func (fc NetworkFilterChain) OnData(data interface{}) (interface{}, error) {
	// todo: only one filter will exist for now, needs change when more than one
	for _, filter := range fc.filtersArray {
		return filter.OnData(data)
	}
	return nil, errors.Errorf("filterChain don't have network filter")
}

// OnTripleData handle triple rpc invocation
func (fc *NetworkFilterChain) OnTripleData(ctx context.Context, methodName string, arguments []interface{}) (interface{}, error) {
	// todo: only one filter will exist for now, needs change when more than one
	for _, filter := range fc.filtersArray {
		return filter.OnTripleData(ctx, methodName, arguments)
	}
	return nil, errors.Errorf("filterChain don't have network filter")
}

// CreateNetworkFilterChain create network filter chain
func CreateNetworkFilterChain(config model.FilterChain) *NetworkFilterChain {
	var filters []filter.NetworkFilter

	for _, f := range config.Filters {
		p, err := filter.GetNetworkFilterPlugin(f.Name)
		if err != nil {
			logger.Error("CreateNetworkFilterChain %s getNetworkFilterPlugin error %s", f.Name, err)
			continue
		}

		config := p.Config()
		if err := yaml.ParseConfig(config, f.Config); err != nil {
			logger.Error("CreateNetworkFilterChain %s parse config error %s", f.Name, err)
			continue
		}

		filter, err := p.CreateFilter(config)
		if err != nil {
			logger.Error("CreateNetworkFilterChain %s createFilter error %s", f.Name, err)
			continue
		}
		filters = append(filters, filter)
	}

	return &NetworkFilterChain{
		filtersArray: filters,
		config:       config,
	}
}

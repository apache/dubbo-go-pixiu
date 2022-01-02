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
	"net/http"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type NetworkFilterChain struct {
	filtersArray []filter.NetworkFilter
	config       model.FilterChain
}

func (fc NetworkFilterChain) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	// todo: only one filter will exist for now, needs change when more than one
	for _, filter := range fc.filtersArray {
		filter.ServeHTTP(w, r)
	}
}

func CreateNetworkFilterChain(config model.FilterChain, bs *model.Bootstrap) *NetworkFilterChain {
	filtersArray := make([]filter.NetworkFilter, len(config.Filters))
	// todo: split code block like http filter manager
	// todo: only one filter will exist for now, needs change when more than one
	for i, f := range config.Filters {
		if f.Name == constant.GRPCConnectManagerFilter {
			gcmc := &model.GRPCConnectionManagerConfig{}
			if err := yaml.ParseConfig(gcmc, f.Config); err != nil {
				logger.Error("CreateNetworkFilterChain %s parse config error %s", f.Name, err)
			}
			p, err := filter.GetNetworkFilterPlugin(constant.GRPCConnectManagerFilter)
			if err != nil {
				logger.Error("CreateNetworkFilterChain %s getNetworkFilterPlugin error %s", f.Name, err)
			}
			filter, err := p.CreateFilter(gcmc, bs)
			if err != nil {
				logger.Error("CreateNetworkFilterChain %s createFilter error %s", f.Name, err)
			}
			filtersArray[i] = filter
		} else if f.Name == constant.HTTPConnectManagerFilter {
			hcmc := &model.HttpConnectionManagerConfig{}
			if err := yaml.ParseConfig(hcmc, f.Config); err != nil {
				logger.Error("CreateNetworkFilterChain parse %s config error %s", f.Name, err)
			}
			p, err := filter.GetNetworkFilterPlugin(constant.HTTPConnectManagerFilter)
			if err != nil {
				logger.Error("CreateNetworkFilterChain %s getNetworkFilterPlugin error %s", f.Name, err)
			}
			filter, err := p.CreateFilter(hcmc, bs)
			if err != nil {
				logger.Error("CreateNetworkFilterChain %s createFilter error %s", f.Name, err)
			}
			filtersArray[i] = filter
		}
	}

	return &NetworkFilterChain{
		filtersArray: filtersArray,
		config:       config,
	}
}

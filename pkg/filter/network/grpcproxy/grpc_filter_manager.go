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

package grpcproxy

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// GrpcFilterManager manage gRPC filters, similar to DubboFilterManager
type GrpcFilterManager struct {
	filterConfigs []*model.GrpcFilter
	filters       []filter.GrpcFilter
}

// NewGrpcFilterManager create gRPC filter manager
func NewGrpcFilterManager(fs []*model.GrpcFilter) *GrpcFilterManager {
	filters := createGrpcFilter(fs)
	fm := &GrpcFilterManager{filterConfigs: fs, filters: filters}
	return fm
}

func createGrpcFilter(fs []*model.GrpcFilter) []filter.GrpcFilter {
	var filters []filter.GrpcFilter

	for _, f := range fs {
		p, err := filter.GetGrpcFilterPlugin(f.Name)
		if err != nil {
			logger.Error("createGrpcFilter %s getGrpcFilterPlugin error %s", f.Name, err)
			continue
		}

		config := p.Config()
		if err := yaml.ParseConfig(config, f.Config); err != nil {
			logger.Error("createGrpcFilter %s parse config error %s", f.Name, err)
			continue
		}

		grpcFilter, err := p.CreateFilter(config)
		if err != nil {
			logger.Error("createGrpcFilter %s createFilter error %s", f.Name, err)
			continue
		}
		filters = append(filters, grpcFilter)
	}
	return filters
}

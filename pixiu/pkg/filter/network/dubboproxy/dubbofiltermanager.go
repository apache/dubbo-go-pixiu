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

package dubboproxy

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

// DubboFilterManager manage filters
type DubboFilterManager struct {
	filterConfigs []*model.DubboFilter
	filters       []filter.DubboFilter
}

// NewDubboFilterManager create filter manager
func NewDubboFilterManager(fs []*model.DubboFilter) *DubboFilterManager {
	filters := createDubboFilter(fs)
	fm := &DubboFilterManager{filterConfigs: fs, filters: filters}
	return fm
}

func createDubboFilter(fs []*model.DubboFilter) []filter.DubboFilter {
	var filters []filter.DubboFilter

	for _, f := range fs {
		p, err := filter.GetDubboFilterPlugin(f.Name)
		if err != nil {
			logger.Error("createDubboFilter %s getNetworkFilterPlugin error %s", f.Name, err)
			continue
		}

		config := p.Config()
		if err := yaml.ParseConfig(config, f.Config); err != nil {
			logger.Error("createDubboFilter %s parse config error %s", f.Name, err)
			continue
		}

		filter, err := p.CreateFilter(config)
		if err != nil {
			logger.Error("createDubboFilter %s createFilter error %s", f.Name, err)
			continue
		}
		filters = append(filters, filter)
	}
	return filters
}

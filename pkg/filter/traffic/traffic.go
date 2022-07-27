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

package traffic

import (
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	canaryByHeader CanaryHeaders = "canary-by-header"
	canaryWeight   CanaryHeaders = "canary-weight"
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	CanaryHeaders string

	Plugin struct {
	}

	FilterFactory struct {
		cfg *Config
	}

	Filter struct {
		Rules map[CanaryHeaders]*Cluster
	}

	Config struct {
		Traffics []*Cluster `yaml:"traffics" json:"traffics" mapstructure:"traffics"`
	}

	Cluster struct {
		Name   string        `yaml:"name" json:"name" mapstructure:"name"`
		Router string        `yaml:"router" json:"router" mapstructure:"router"`
		Model  CanaryHeaders `yaml:"model" json:"model" mapstructure:"model"`    // canary model
		Canary string        `yaml:"canary" json:"canary" mapstructure:"canary"` // header key
		Value  string        `yaml:"value" json:"value" mapstructure:"value"`    // header value
	}
)

func (p *Plugin) Kind() string {
	return constant.HTTPTrafficFilter
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{
		cfg: &Config{},
	}, nil
}

func (factory *FilterFactory) Config() interface{} {
	return factory.cfg
}

func (factory *FilterFactory) Apply() error {
	initActions()
	return nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
	f := &Filter{}
	f.Rules = factory.rulesMatch(ctx.Request.RequestURI)
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(ctx *http.HttpContext) filter.FilterStatus {
	if f.Rules != nil {
		cluster := spilt(ctx.Request, f.Rules)
		if cluster != nil {
			ctx.Route.Cluster = cluster.Name
			logger.Debugf("[dubbo-go-pixiu] execute traffic split to cluster %s", cluster.Name)
		}
	} else {
		logger.Warnf("[dubbo-go-pixiu] execute traffic split fail because of empty rules.")
	}
	return filter.Continue
}

func (factory *FilterFactory) rulesMatch(path string) map[CanaryHeaders]*Cluster {
	clusters := factory.cfg.Traffics

	if clusters != nil {
		mp := make(map[CanaryHeaders]*Cluster)
		for i := range clusters {
			if strings.HasPrefix(clusters[i].Router, path) {
				mp[clusters[i].Model] = clusters[i]
			}
		}
		return mp
	}
	return nil
}

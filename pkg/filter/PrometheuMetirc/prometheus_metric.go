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

package metric

import ()

import ()

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

const (
	Kind = constant.HTTPPrometheusMetricFilter
)

type (
	Plugin struct {
	}
	// FilterFactory is http filter instance
	FilterFactory struct {
		cfg    *Config
		errMsg []byte
	}

	Filter struct {
		cfg    *Config
		errMsg []byte
	}

	// Config describe the config of FilterFactory
	Config struct {
		ErrMsg                string               `yaml:"err_msg" json:"err_msg" mapstructure:"err_msg"`
		Rules                 Rules                `yaml:"rules" json:"rules" mapstructure:"rules"`
		ApiStatsProviders     ApiStatsResponse     `yaml:"api_stats_providers" json:"api_stats_providers" mapstructure:"api_stats_providers"`
		ClusterStatsProviders ClusterStatsResponse `yaml:"cluster_stats_providers" json:"cluster_stats_providers" mapstructure:"cluster_stats_providers"`
		CommonStatsProviders  CommonStatsResponse  `yaml:"common_stats_providers" json:"common_stats_providers" mapstructure:"common_stats_providers"`
	}
)

func (p Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{
		cfg: &Config{},
	}, nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {

	f := &Filter{
		cfg:    factory.cfg,
		errMsg: factory.errMsg,
	}
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(ctx *http.HttpContext) filter.FilterStatus {

	return filter.Continue
}

func (factory *FilterFactory) Apply() error {
	return nil
}

func (factory *FilterFactory) Config() interface{} {
	return factory.cfg
}

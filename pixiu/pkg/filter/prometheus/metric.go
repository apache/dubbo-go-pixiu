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

package prometheus

import (
	"strconv"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	contextHttp "github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/metrics/collector"
	"github.com/apache/dubbo-go-pixiu/pkg/metrics/scrape"
	"github.com/apache/dubbo-go-pixiu/pkg/metrics/scrape/scrapeImpl"
)

const (
	Kind = constant.HTTPPrometheusMetricFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	Plugin struct {
	}

	FilterFactory struct {
		Cfg *MetricCollectConfiguration
	}

	Filter struct {
		Cfg *MetricCollectConfiguration
	}
)

func (p Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{Cfg: &MetricCollectConfiguration{}}, nil
}

func (factory *FilterFactory) Config() interface{} {
	return factory.Cfg
}

func (factory *FilterFactory) Apply() error {
	return nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *contextHttp.HttpContext, chain filter.FilterChain) error {
	f := &Filter{
		Cfg: factory.Cfg,
	}
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(ctx *contextHttp.HttpContext) filter.FilterStatus {

	cfg := collector.Config{
		Enable:       f.Cfg.Rules.Enable,
		MeticPath:    f.Cfg.Rules.MeticPath,
		ExporterPort: ":" + strconv.Itoa(f.Cfg.Rules.ExporterPort),
	}
	enabledScrapers := []scrape.Scraper{
		scrapeImpl.ApiScraper{},
		scrapeImpl.ClusterScraper{},
		scrapeImpl.FtontendScraper{},
	}
	go func() {
		collector.ExporterServerOnHttp(ctx, cfg, enabledScrapers)
	}()
	return filter.Continue
}

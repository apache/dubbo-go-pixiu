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
	stdHttp "net/http"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	contextHttp "github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	prom "github.com/apache/dubbo-go-pixiu/pkg/metrics/prometheus"
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
		Cfg  *MetricCollectConfiguration
		Prom *prom.Prometheus
	}

	Filter struct {
		Cfg  *MetricCollectConfiguration
		Prom *prom.Prometheus
	}
)

func (p Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {

	return &FilterFactory{
		Cfg:  &MetricCollectConfiguration{},
		Prom: prom.NewPrometheus(),
	}, nil
}

func (factory *FilterFactory) Config() interface{} {

	return factory.Cfg
}

func (factory *FilterFactory) Apply() error {

	return nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *contextHttp.HttpContext, chain filter.FilterChain) error {

	f := &Filter{
		Cfg:  factory.Cfg,
		Prom: factory.Prom,
	}
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(ctx *contextHttp.HttpContext) filter.FilterStatus {

	if f.Cfg == nil {
		logger.Errorf("Message:Filter Metric Collect Configuration is null")
		ctx.SendLocalReply(stdHttp.StatusForbidden, constant.Default403Body)
		return filter.Continue
	}
	if f.Prom == nil {
		logger.Errorf("Message:Prometheus Collector is not initialized")
		ctx.SendLocalReply(stdHttp.StatusForbidden, constant.Default403Body)
		return filter.Continue
	}

	f.Prom.SetMetricPath(f.Cfg.Rules.MetricPath)

	f.Prom.SetPushGatewayUrl(f.Cfg.Rules.PushGatewayURL, "/metrics", 3)
	start := f.Prom.HandlerFunc()
	err := start(ctx)
	if err != nil {
		logger.Errorf("Message:Context HandlerFunc error")
		ctx.SendLocalReply(stdHttp.StatusForbidden, constant.Default403Body)
	}
	return filter.Continue
}

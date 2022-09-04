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
	"time"
)
import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	contextHttp "github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	prom "github.com/apache/dubbo-go-pixiu/pkg/metrics/prometheus"
)

const (
	Kind = constant.HTTPPrometheusMetricFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{
		Cfg: &MetricCollectConfiguration{
			Rules: MetricCollectRule{
				Enable:              true,
				MeticPath:           "/metrics",
				PushGatewayURL:      "http://domain:port",
				PushIntervalSeconds: 3,
				PushJobName:         "prometheus",
			},
		},
	})
}

type (
	Plugin struct {
		Cfg *MetricCollectConfiguration
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

	return &FilterFactory{
		Cfg: p.Cfg,
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
		Cfg: factory.Cfg,
	}
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(ctx *contextHttp.HttpContext) filter.FilterStatus {
	p := prom.NewPrometheus("pixiu")

	p.SetMetricPath(f.Cfg.Rules.MeticPath)

	p.SetPushGateway(f.Cfg.Rules.PushGatewayURL, time.Duration(f.Cfg.Rules.PushIntervalSeconds), f.Cfg.Rules.PushJobName)

	start := p.HandlerFunc()
	start(ctx)
	return filter.Continue
}

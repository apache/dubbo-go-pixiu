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

import (
	"context"
	"sync/atomic"
	"time"
)

import (
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	Kind = constant.HTTPMetricFilter
)

var (
	totalElapsed int64
	totalCount   int64
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// FilterFactory is http filter plugin.
	Plugin struct {
	}
	// FilterFactory is http filter instance
	FilterFactory struct {
	}
	Filter struct {
		start time.Time
	}
	// Config describe the config of FilterFactory
	Config struct{}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	registerOtelMetric()
	return &FilterFactory{}, nil
}

func (factory *FilterFactory) Config() interface{} {
	return struct{}{}
}

func (factory *FilterFactory) Apply() error {
	// init
	registerOtelMetric()
	return nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
	f := &Filter{}
	chain.AppendDecodeFilters(f)
	chain.AppendEncodeFilters(f)
	return nil
}

func (f *Filter) Decode(c *http.HttpContext) filter.FilterStatus {
	f.start = time.Now()
	return filter.Continue
}

func (f *Filter) Encode(c *http.HttpContext) filter.FilterStatus {
	latency := time.Since(f.start)
	atomic.AddInt64(&totalElapsed, latency.Nanoseconds())
	atomic.AddInt64(&totalCount, 1)

	logger.Infof("[dubbo go server] [UPSTREAM] receive request | %d | %s | %s | %s | ", c.GetStatusCode(), latency, c.GetMethod(), c.GetUrl())
	return filter.Continue
}

func registerOtelMetric() {
	meter := global.GetMeterProvider().Meter("pixiu")
	observerElapsedCallback := func(_ context.Context, result metric.Int64ObserverResult) {
		result.Observe(totalElapsed)
	}
	_ = metric.Must(meter).NewInt64SumObserver("pixiu_request_elapsed", observerElapsedCallback,
		metric.WithDescription("request total elapsed in pixiu"),
	)

	observerCountCallback := func(_ context.Context, result metric.Int64ObserverResult) {
		result.Observe(totalCount)
	}
	_ = metric.Must(meter).NewInt64SumObserver("pixiu_request_count", observerCountCallback,
		metric.WithDescription("request total count in pixiu"),
	)
}

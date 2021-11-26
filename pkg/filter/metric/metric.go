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
	// Filter is http filter plugin.
	Plugin struct {
	}
	// Filter is http filter instance
	Filter struct {
	}
	// Config describe the config of Filter
	Config struct{}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter() (filter.HttpFilterFactory, error) {
	registerOtelMetric()
	return &Filter{}, nil
}

func (f *Filter) Config() interface{} {
	return nil
}

func (f *Filter) Apply() error {
	// init
	registerOtelMetric()
	return nil
}

func (f *Filter) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *Filter) Handle(c *http.HttpContext) {
	start := time.Now()

	c.Next()

	latency := time.Since(start)
	atomic.AddInt64(&totalElapsed, latency.Nanoseconds())
	atomic.AddInt64(&totalCount, 1)

	logger.Infof("[dubbo go server] [UPSTREAM] receive request | %d | %s | %s | %s | ", c.StatusCode(), latency, c.GetMethod(), c.GetUrl())
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

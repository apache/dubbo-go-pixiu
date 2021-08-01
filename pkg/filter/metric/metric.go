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
	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/context"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/filter"

	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/global"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	manager "github.com/apache/dubbo-go-pixiu/pkg/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

var (
	totalElapsed int64
	totalCount   int64
)

// nolint
func Init() {
	manager.RegisterFilterFactory(constant.MetricFilter, newMetricFilter)
}

// loggerFilter is a filter for simple metric.
type metricFilter struct {
}

// New create metric filter.
func newMetricFilter() filter.Factory {
	return &metricFilter{}
}

func (m *metricFilter) Config() interface{} {
	return nil
}

func (m *metricFilter) Apply() (filter.Filter, error) {
	// init
	registerOtelMetric()

	// Metric filter, record url and latency
	return func(c fc.Context) {
		start := time.Now()

		c.Next()

		latency := time.Now().Sub(start)
		atomic.AddInt64(&totalElapsed, latency.Nanoseconds())
		atomic.AddInt64(&totalCount, 1)

		logger.Infof("[dubbo go pixiu] [UPSTREAM] receive request | %d | %s | %s | %s | ", c.StatusCode(), latency, c.GetMethod(), c.GetUrl())
	}, nil
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

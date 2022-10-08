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
	"time"
)

import (
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

const (
	Kind = constant.HTTPMetricFilter
)

var (
	totalElapsed syncint64.Counter
	totalCount   syncint64.Counter
	totalError   syncint64.Counter
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
	return &FilterFactory{}, nil
}

func (factory *FilterFactory) Config() interface{} {
	return &struct{}{}
}

func (factory *FilterFactory) Apply() error {
	// init
	err := registerOtelMetric()
	return err
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
	totalCount.Add(c.Ctx, 1)
	totalElapsed.Add(c.Ctx, latency.Nanoseconds())
	if c.LocalReply() {
		totalError.Add(c.Ctx, 1)
	}

	logger.Debugf("[Metric] [UPSTREAM] receive request | %d | %s | %s | %s | ", c.GetStatusCode(), latency, c.GetMethod(), c.GetUrl())
	return filter.Continue
}

func registerOtelMetric() error {
	meter := global.MeterProvider().Meter("pixiu")

	elapsedCounter, err := meter.SyncInt64().Counter("pixiu_request_elapsed", instrument.WithDescription("request total elapsed in pixiu"))
	if err != nil {
		logger.Errorf("register pixiu_request_elapsed metric failed, err: %v", err)
		return err
	}
	totalElapsed = elapsedCounter

	count, err := meter.SyncInt64().Counter("pixiu_request_count", instrument.WithDescription("request total count in pixiu"))
	if err != nil {
		logger.Errorf("register pixiu_request_count metric failed, err: %v", err)
		return err
	}
	totalCount = count

	errorCounter, err := meter.SyncInt64().Counter("pixiu_request_error_count", instrument.WithDescription("request error total count in pixiu"))
	if err != nil {
		logger.Errorf("register pixiu_request_error_count metric failed, err: %v", err)
		return err
	}
	totalError = errorCounter
	return nil
}

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
	"fmt"
	stdhttp "net/http"
	"time"
)

import (
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	Kind = constant.HTTPMetricFilter
)

var (
	totalElapsed syncint64.Counter
	totalCount   syncint64.Counter
	totalError   syncint64.Counter

	sizeRequest  syncint64.Counter
	sizeResponse syncint64.Counter
	durationHist syncint64.Histogram
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

	commonAttrs := []attribute.KeyValue{
		attribute.String("code", fmt.Sprintf("%d", c.GetStatusCode())),
		attribute.String("method", c.Request.Method),
		attribute.String("url", c.GetUrl()),
		attribute.String("host", c.Request.Host),
	}

	latency := time.Since(f.start)
	totalCount.Add(c.Ctx, 1, commonAttrs...)
	latencyMilli := latency.Milliseconds()
	totalElapsed.Add(c.Ctx, latencyMilli, commonAttrs...)
	if c.LocalReply() {
		totalError.Add(c.Ctx, 1)
	}

	durationHist.Record(c.Ctx, latencyMilli, commonAttrs...)
	size, err := computeApproximateRequestSize(c.Request)
	if err != nil {
		logger.Warn("can not compute request size", err)
	} else {
		sizeRequest.Add(c.Ctx, int64(size), commonAttrs...)
	}

	size, err = computeApproximateResponseSize(c.TargetResp)
	if err != nil {
		logger.Warn("can not compute response size", err)
	} else {
		sizeResponse.Add(c.Ctx, int64(size), commonAttrs...)
	}

	logger.Debugf("[Metric] [UPSTREAM] receive request | %d | %s | %s | %s | ", c.GetStatusCode(), latency, c.GetMethod(), c.GetUrl())
	return filter.Continue
}

func computeApproximateResponseSize(res *client.Response) (int, error) {
	if res == nil {
		return 0, errors.New("client.Response is null pointer ")
	}
	return len(res.Data), nil
}

func computeApproximateRequestSize(r *stdhttp.Request) (int, error) {
	if r == nil {
		return 0, errors.New("http.Request is null pointer ")
	}
	s := 0
	if r.URL != nil {
		s = len(r.URL.Path)
	}
	s += len(r.Method)
	s += len(r.Proto)
	for name, values := range r.Header {
		s += len(name)
		for _, value := range values {
			s += len(value)
		}
	}
	s += len(r.Host)
	if r.ContentLength != -1 {
		s += int(r.ContentLength)
	}
	return s, nil
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

	sizeRequest, err = meter.SyncInt64().Counter("pixiu_request_content_length", instrument.WithDescription("request total content length in pixiu"))
	if err != nil {
		logger.Errorf("register pixiu_request_content_length metric failed, err: %v", err)
		return err
	}

	sizeResponse, err = meter.SyncInt64().Counter("pixiu_response_content_length", instrument.WithDescription("request total content length response in pixiu"))
	if err != nil {
		logger.Errorf("register pixiu_response_content_length metric failed, err: %v", err)
		return err
	}

	durationHist, err = meter.SyncInt64().Histogram(
		"pixiu_process_time_millicec",
		instrument.WithDescription("request process time response in pixiu"),
	)
	if err != nil {
		logger.Errorf("register pixiu_process_time_millisec metric failed, err: %v", err)
		return err
	}

	return nil
}

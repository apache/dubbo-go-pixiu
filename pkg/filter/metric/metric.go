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
	"sync"
	"time"
)

import (
	"github.com/pkg/errors"

	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/metric"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contextHttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	prom "github.com/apache/dubbo-go-pixiu/pkg/prometheus"
)

const (
	// Kind defines the filter kind
	Kind = constant.HTTPMetricFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is http filter plugin.
	Plugin struct{}

	// FilterFactory is http filter instance
	FilterFactory struct {
		cfg *Config
	}

	// Filter instance
	Filter struct {
		cfg *Config

		otelInstruments *OTelInstruments
		promCollector   *prom.Prometheus
		start           time.Time
	}
)

// Kind returns the filter kind.
func (p *Plugin) Kind() string {
	return Kind
}

// CreateFilterFactory creates a new filter factory.
func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{
		cfg: &Config{},
	}, nil
}

// Config returns the configuration.
func (factory *FilterFactory) Config() any {
	return factory.cfg
}

// Apply validates the configuration.
func (factory *FilterFactory) Apply() error {
	return factory.cfg.Validate()
}

var (
	globalOTelInstruments *OTelInstruments
	otelInitOnce          sync.Once
	otelInitErr           error
)

// initOTelInstruments initializes OpenTelemetry instruments (singleton).
func initOTelInstruments() (*OTelInstruments, error) {
	otelInitOnce.Do(func() {
		otelInitErr = doInitOTelInstruments()
	})
	return globalOTelInstruments, otelInitErr
}

func doInitOTelInstruments() error {
	meter := otel.GetMeterProvider().Meter("pixiu")

	instruments := &OTelInstruments{}

	elapsedCounter, err := meter.Int64Counter("pixiu_request_elapsed",
		metric.WithDescription("request total elapsed in pixiu"))
	if err != nil {
		return fmt.Errorf("register pixiu_request_elapsed metric failed: %w", err)
	}
	instruments.totalElapsed = elapsedCounter

	count, err := meter.Int64Counter("pixiu_request_count",
		metric.WithDescription("request total count in pixiu"))
	if err != nil {
		return fmt.Errorf("register pixiu_request_count metric failed: %w", err)
	}
	instruments.totalCount = count

	errorCounter, err := meter.Int64Counter("pixiu_request_error_count",
		metric.WithDescription("request error total count in pixiu"))
	if err != nil {
		return fmt.Errorf("register pixiu_request_error_count metric failed: %w", err)
	}
	instruments.totalError = errorCounter

	sizeRequest, err := meter.Int64Counter("pixiu_request_content_length",
		metric.WithDescription("request total content length in pixiu"))
	if err != nil {
		return fmt.Errorf("register pixiu_request_content_length metric failed: %w", err)
	}
	instruments.sizeRequest = sizeRequest

	sizeResponse, err := meter.Int64Counter("pixiu_response_content_length",
		metric.WithDescription("request total content length response in pixiu"))
	if err != nil {
		return fmt.Errorf("register pixiu_response_content_length metric failed: %w", err)
	}
	instruments.sizeResponse = sizeResponse

	durationHist, err := meter.Int64Histogram("pixiu_process_time_millisec",
		metric.WithDescription("request process time response in pixiu"))
	if err != nil {
		return fmt.Errorf("register pixiu_process_time_millisec metric failed: %w", err)
	}
	instruments.durationHist = durationHist

	globalOTelInstruments = instruments
	logger.Infof("[MetricReporter] OpenTelemetry instruments registered")
	return nil
}

// PrepareFilterChain prepares the filter chain.
func (factory *FilterFactory) PrepareFilterChain(ctx *contextHttp.HttpContext, chain filter.FilterChain) error {
	// Copy config to avoid sharing factory's pointer
	cfgCopy := *factory.cfg
	f := &Filter{cfg: &cfgCopy}

	// Initialize based on mode
	switch factory.cfg.Mode {
	case "pull":
		instruments, err := initOTelInstruments()
		if err != nil {
			return err
		}
		f.otelInstruments = instruments
		logger.Infof("[MetricReporter] Pull mode enabled")

	case "push":
		p := prom.NewPrometheus()
		p.SetPushGatewayUrl(factory.cfg.Push.GatewayURL, factory.cfg.Push.MetricPath)
		p.SetPushIntervalThreshold(true, factory.cfg.Push.PushInterval)
		p.SetPushGatewayJob(factory.cfg.Push.JobName)
		f.promCollector = p
		logger.Infof("[MetricReporter] Push mode enabled (gateway: %s, interval: %d)",
			factory.cfg.Push.GatewayURL, factory.cfg.Push.PushInterval)
	}

	// Both modes need decode and encode filters
	chain.AppendDecodeFilters(f)
	chain.AppendEncodeFilters(f)

	return nil
}

// Decode handles the decode phase - records start time for both modes.
func (f *Filter) Decode(ctx *contextHttp.HttpContext) filter.FilterStatus {
	// Record start time for latency calculation
	// Both pull and push modes report metrics in Encode phase
	f.start = time.Now()
	return filter.Continue
}

// Encode reports metrics for both modes.
func (f *Filter) Encode(ctx *contextHttp.HttpContext) filter.FilterStatus {
	switch f.cfg.Mode {
	case "pull":
		return f.reportWithOTel(ctx)
	case "push":
		return f.reportWithPrometheus(ctx)
	}

	return filter.Continue
}

// reportWithOTel reports metrics using OpenTelemetry.
func (f *Filter) reportWithOTel(ctx *contextHttp.HttpContext) filter.FilterStatus {
	if f.otelInstruments == nil {
		logger.Errorf("[MetricReporter] OpenTelemetry instruments not initialized")
		errResp := contextHttp.InternalError.New()
		ctx.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	// Report context metrics dynamically
	contextMetrics := ctx.GetAllMetrics()
	if len(contextMetrics) > 0 {
		meter := otel.GetMeterProvider().Meter("pixiu")

		for _, m := range contextMetrics {
			attrs := toOTelAttributes(m.Labels)

			opts := metric.WithAttributes(attrs...)

			switch m.Type {
			case "counter":
				counter, err := meter.Int64Counter(m.Name,
					metric.WithDescription(fmt.Sprintf("Context counter: %s", m.Name)))
				if err != nil {
					logger.Warnf("[MetricReporter] Failed to create counter %s: %v", m.Name, err)
					continue
				}
				counter.Add(ctx.Ctx, int64(m.Value), opts)

			case "histogram":
				histogram, err := meter.Float64Histogram(m.Name,
					metric.WithDescription(fmt.Sprintf("Context histogram: %s", m.Name)))
				if err != nil {
					logger.Warnf("[MetricReporter] Failed to create histogram %s: %v", m.Name, err)
					continue
				}
				histogram.Record(ctx.Ctx, m.Value, opts)

			case "gauge":
				gauge, err := meter.Int64UpDownCounter(m.Name,
					metric.WithDescription(fmt.Sprintf("Context gauge: %s", m.Name)))
				if err != nil {
					logger.Warnf("[MetricReporter] Failed to create gauge %s: %v", m.Name, err)
					continue
				}
				gauge.Add(ctx.Ctx, int64(m.Value), opts)
			}
		}
	}

	// Report built-in metrics
	commonAttrs := []attribute.KeyValue{
		attribute.String("code", fmt.Sprintf("%d", ctx.GetStatusCode())),
		attribute.String("method", ctx.Request.Method),
		attribute.String("url", ctx.GetUrl()),
		attribute.String("host", ctx.Request.Host),
	}
	commonOpts := metric.WithAttributes(commonAttrs...)

	latency := time.Since(f.start)
	f.otelInstruments.totalCount.Add(ctx.Ctx, 1, commonOpts)
	latencyMilli := latency.Milliseconds()
	f.otelInstruments.totalElapsed.Add(ctx.Ctx, latencyMilli, commonOpts)

	if ctx.LocalReply() {
		f.otelInstruments.totalError.Add(ctx.Ctx, 1)
	}

	f.otelInstruments.durationHist.Record(ctx.Ctx, latencyMilli, commonOpts)

	size, err := computeApproximateRequestSize(ctx.Request)
	if err != nil {
		logger.Warnf("[MetricReporter] Cannot compute request size: %v", err)
	} else {
		f.otelInstruments.sizeRequest.Add(ctx.Ctx, int64(size), commonOpts)
	}

	size, err = computeApproximateResponseSize(ctx.TargetResp)
	if err != nil {
		logger.Warnf("[MetricReporter] Cannot compute response size: %v", err)
	} else {
		f.otelInstruments.sizeResponse.Add(ctx.Ctx, int64(size), commonOpts)
	}

	logger.Debugf("[MetricReporter] [PULL] request | %d | %s | %s | %s |",
		ctx.GetStatusCode(), latency, ctx.GetMethod(), ctx.GetUrl())

	return filter.Continue
}

// reportWithPrometheus reports metrics using Prometheus.
func (f *Filter) reportWithPrometheus(ctx *contextHttp.HttpContext) filter.FilterStatus {
	if f.promCollector == nil {
		logger.Errorf("[MetricReporter] Prometheus collector not initialized")
		errResp := contextHttp.InternalError.New()
		ctx.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	// Process and report custom context metrics
	contextMetrics := ctx.GetAllMetrics()
	for _, m := range contextMetrics {
		if err := f.promCollector.RecordDynamicMetric(m.Name, m.Type, m.Value, m.Labels); err != nil {
			logger.Warnf("[MetricReporter] Failed to record dynamic metric %s: %v", m.Name, err)
		} else {
			logger.Debugf("[MetricReporter] Recorded custom metric: %s=%f (type: %s, labels: %v)",
				m.Name, m.Value, m.Type, m.Labels)
		}
	}

	// Report built-in Prometheus metrics
	handlerFunc := f.promCollector.HandlerFunc()
	if err := handlerFunc(ctx); err != nil {
		logger.Errorf("[MetricReporter] Prometheus handler error: %v", err)
	}

	logger.Debugf("[MetricReporter] [PUSH] request | %d | %s | %s |",
		ctx.GetStatusCode(), ctx.GetMethod(), ctx.GetUrl())

	return filter.Continue
}

// toOTelAttributes converts map[string]string to OpenTelemetry attributes.
func toOTelAttributes(labels map[string]string) []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, 0, len(labels))
	for k, v := range labels {
		attrs = append(attrs, attribute.String(k, v))
	}
	return attrs
}

// computeApproximateRequestSize computes the approximate size of an HTTP request.
func computeApproximateRequestSize(r *stdhttp.Request) (int, error) {
	if r == nil {
		return 0, errors.New("http.Request is null pointer")
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

// computeApproximateResponseSize computes the approximate size of an HTTP response.
func computeApproximateResponseSize(res any) (int, error) {
	if res == nil {
		return 0, errors.New("client response is nil")
	}
	if unaryResponse, ok := res.(*client.UnaryResponse); ok {
		return len(unaryResponse.Data), nil
	}
	return 0, errors.New("response is not of type client.UnaryResponse")
}

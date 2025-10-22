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

package metricreporter

import (
	"context"
	"fmt"
	"net/http"
	"strconv"
	"sync"
)

import (
	sdkprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"

	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric"
	"go.opentelemetry.io/otel/metric/instrument"
	"go.opentelemetry.io/otel/metric/instrument/syncfloat64"
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
	sdkmetric "go.opentelemetry.io/otel/sdk/metric"
)

import (
	contextHttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// OTelPullReporter handles pull mode using OpenTelemetry with Prometheus Exporter.
type OTelPullReporter struct {
	config   *PullConfig
	meter    metric.Meter
	exporter *prometheus.Exporter
	provider *sdkmetric.MeterProvider
	mu       sync.RWMutex

	// Cache for instruments (v0.32.x compatible)
	registeredCounters   map[string]syncint64.Counter
	registeredHistograms map[string]syncfloat64.Histogram
	registeredGauges     map[string]syncint64.UpDownCounter
}

// NewOTelPullReporter creates a new OpenTelemetry pull reporter.
func NewOTelPullReporter(config *PullConfig) *OTelPullReporter {
	return &OTelPullReporter{
		config:               config,
		registeredCounters:   make(map[string]syncint64.Counter),
		registeredHistograms: make(map[string]syncfloat64.Histogram),
		registeredGauges:     make(map[string]syncint64.UpDownCounter),
	}
}

// Start starts the pull reporter HTTP server with OpenTelemetry.
func (r *OTelPullReporter) Start() error {
	// Create Prometheus exporter (v0.32.x returns value, not pointer)
	exporter := prometheus.New()
	r.exporter = &exporter

	// Create MeterProvider
	provider := sdkmetric.NewMeterProvider(sdkmetric.WithReader(&exporter))
	r.provider = provider

	// Get Meter
	r.meter = provider.Meter("pixiu")

	// Setup HTTP handler
	registry := sdkprometheus.NewRegistry()
	if err := registry.Register(exporter.Collector); err != nil {
		logger.Warnf("[MetricReporter] Failed to register exporter collector: %v", err)
	}

	handler := promhttp.HandlerFor(registry, promhttp.HandlerOpts{
		EnableOpenMetrics: true,
	})

	mux := http.NewServeMux()
	mux.Handle(r.config.Path, handler)

	addr := ":" + strconv.Itoa(r.config.Port)

	go func() {
		logger.Infof("[MetricReporter] Starting OpenTelemetry pull mode HTTP server on %s%s", addr, r.config.Path)
		if err := http.ListenAndServe(addr, mux); err != nil {
			logger.Errorf("[MetricReporter] Pull mode HTTP server error: %v", err)
		}
	}()

	return nil
}

// Report updates metrics using OpenTelemetry API.
func (r *OTelPullReporter) Report(metrics []*contextHttp.MetricData) {
	r.mu.Lock()
	defer r.mu.Unlock()

	ctx := context.Background()

	for _, m := range metrics {
		// Convert labels to OpenTelemetry attributes
		attrs := toOTelAttributes(m.Labels)

		switch m.Type {
		case "counter":
			r.recordCounter(ctx, m.Name, int64(m.Value), attrs)

		case "histogram":
			r.recordHistogram(ctx, m.Name, m.Value, attrs)

		case "gauge":
			r.recordGauge(ctx, m.Name, int64(m.Value), attrs)
		}
	}
}

// recordCounter records a counter metric using OpenTelemetry v0.32.x API.
func (r *OTelPullReporter) recordCounter(ctx context.Context, name string, value int64, attrs []attribute.KeyValue) {
	// Get or create counter instrument
	counter, exists := r.registeredCounters[name]
	if !exists {
		var err error
		counter, err = r.meter.SyncInt64().Counter(
			name,
			instrument.WithDescription(fmt.Sprintf("Counter metric: %s", name)),
		)
		if err != nil {
			logger.Warnf("[MetricReporter] Failed to create counter %s: %v", name, err)
			return
		}
		r.registeredCounters[name] = counter
		logger.Debugf("[MetricReporter] Registered counter: %s", name)
	}

	// Record the value
	counter.Add(ctx, value, attrs...)
}

// recordHistogram records a histogram metric using OpenTelemetry v0.32.x API.
func (r *OTelPullReporter) recordHistogram(ctx context.Context, name string, value float64, attrs []attribute.KeyValue) {
	// Get or create histogram instrument
	histogram, exists := r.registeredHistograms[name]
	if !exists {
		var err error
		histogram, err = r.meter.SyncFloat64().Histogram(
			name,
			instrument.WithDescription(fmt.Sprintf("Histogram metric: %s", name)),
		)
		if err != nil {
			logger.Warnf("[MetricReporter] Failed to create histogram %s: %v", name, err)
			return
		}
		r.registeredHistograms[name] = histogram
		logger.Debugf("[MetricReporter] Registered histogram: %s", name)
	}

	// Record the value
	histogram.Record(ctx, value, attrs...)
}

// recordGauge records a gauge metric using OpenTelemetry v0.32.x API (UpDownCounter).
func (r *OTelPullReporter) recordGauge(ctx context.Context, name string, value int64, attrs []attribute.KeyValue) {
	// Get or create gauge instrument
	gauge, exists := r.registeredGauges[name]
	if !exists {
		var err error
		// In OpenTelemetry, Gauge is represented as UpDownCounter
		gauge, err = r.meter.SyncInt64().UpDownCounter(
			name,
			instrument.WithDescription(fmt.Sprintf("Gauge metric: %s", name)),
		)
		if err != nil {
			logger.Warnf("[MetricReporter] Failed to create gauge %s: %v", name, err)
			return
		}
		r.registeredGauges[name] = gauge
		logger.Debugf("[MetricReporter] Registered gauge: %s", name)
	}

	// Add the value
	gauge.Add(ctx, value, attrs...)
}

// toOTelAttributes converts map[string]string to OpenTelemetry attributes.
func toOTelAttributes(labels map[string]string) []attribute.KeyValue {
	attrs := make([]attribute.KeyValue, 0, len(labels))
	for k, v := range labels {
		attrs = append(attrs, attribute.String(k, v))
	}
	return attrs
}


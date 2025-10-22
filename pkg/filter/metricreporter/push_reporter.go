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
	"bytes"
	"fmt"
	"net/http"
	"os"
	"sync"
	"time"
)

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/common/expfmt"
)

import (
	contextHttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// PushReporter handles push mode (Push Gateway).
type PushReporter struct {
	config   *PushConfig
	registry *prometheus.Registry
	counter  int
	mu       sync.Mutex

	// Cache for metrics
	counters   map[string]*prometheus.CounterVec
	histograms map[string]*prometheus.HistogramVec
	gauges     map[string]*prometheus.GaugeVec
}

// NewPushReporter creates a new push reporter.
func NewPushReporter(config *PushConfig) *PushReporter {
	return &PushReporter{
		config:     config,
		registry:   prometheus.NewRegistry(),
		counters:   make(map[string]*prometheus.CounterVec),
		histograms: make(map[string]*prometheus.HistogramVec),
		gauges:     make(map[string]*prometheus.GaugeVec),
	}
}

// Report updates metrics and pushes to gateway when threshold is reached.
func (pr *PushReporter) Report(metrics []*contextHttp.MetricData) {
	pr.mu.Lock()
	defer pr.mu.Unlock()

	// Update metrics in registry
	for _, metric := range metrics {
		labelNames := make([]string, 0, len(metric.Labels))
		labelValues := make([]string, 0, len(metric.Labels))
		for k, v := range metric.Labels {
			labelNames = append(labelNames, k)
			labelValues = append(labelValues, v)
		}

		switch metric.Type {
		case "counter":
			counter := pr.getOrCreateCounter(metric.Name, labelNames)
			counter.WithLabelValues(labelValues...).Add(metric.Value)

		case "histogram":
			histogram := pr.getOrCreateHistogram(metric.Name, labelNames)
			histogram.WithLabelValues(labelValues...).Observe(metric.Value)

		case "gauge":
			gauge := pr.getOrCreateGauge(metric.Name, labelNames)
			gauge.WithLabelValues(labelValues...).Set(metric.Value)
		}
	}

	// Increment counter and check if we should push
	pr.counter++
	if pr.counter >= pr.config.PushInterval {
		go pr.push()
		pr.counter = 0
	}
}

// push sends metrics to Push Gateway.
func (pr *PushReporter) push() {
	// Gather metrics
	metricFamilies, err := pr.registry.Gather()
	if err != nil {
		logger.Errorf("[MetricReporter] Failed to gather metrics: %v", err)
		return
	}

	// Format metrics
	buf := &bytes.Buffer{}
	for _, mf := range metricFamilies {
		if _, err := expfmt.MetricFamilyToText(buf, mf); err != nil {
			logger.Errorf("[MetricReporter] Failed to format metric family: %v", err)
			return
		}
	}

	// Build push URL
	hostname, _ := os.Hostname()
	if hostname == "" {
		hostname = "unknown"
	}
	url := fmt.Sprintf("%s%s/job/%s/instance/%s",
		pr.config.GatewayURL,
		pr.config.MetricPath,
		pr.config.JobName,
		hostname,
	)

	// Send HTTP POST request
	req, err := http.NewRequest(http.MethodPost, url, buf)
	if err != nil {
		logger.Errorf("[MetricReporter] Failed to create push request: %v", err)
		return
	}
	req.Header.Set("Content-Type", string(expfmt.FmtText))

	client := &http.Client{Timeout: 5 * time.Second}
	resp, err := client.Do(req)
	if err != nil {
		logger.Errorf("[MetricReporter] Failed to push metrics to %s: %v", url, err)
		return
	}
	defer resp.Body.Close()

	if resp.StatusCode >= 200 && resp.StatusCode < 300 {
		logger.Debugf("[MetricReporter] Successfully pushed metrics to %s", url)
	} else {
		logger.Warnf("[MetricReporter] Push gateway returned status %d", resp.StatusCode)
	}
}

// getOrCreateCounter gets or creates a counter metric.
func (pr *PushReporter) getOrCreateCounter(name string, labelNames []string) *prometheus.CounterVec {
	if counter, exists := pr.counters[name]; exists {
		return counter
	}

	counter := prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: name,
			Help: fmt.Sprintf("Counter metric: %s", name),
		},
		labelNames,
	)

	if err := pr.registry.Register(counter); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			if existing, ok := are.ExistingCollector.(*prometheus.CounterVec); ok {
				pr.counters[name] = existing
				return existing
			}
		}
		logger.Warnf("[MetricReporter] Failed to register counter %s: %v", name, err)
	}

	pr.counters[name] = counter
	return counter
}

// getOrCreateHistogram gets or creates a histogram metric.
func (pr *PushReporter) getOrCreateHistogram(name string, labelNames []string) *prometheus.HistogramVec {
	if histogram, exists := pr.histograms[name]; exists {
		return histogram
	}

	histogram := prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    name,
			Help:    fmt.Sprintf("Histogram metric: %s", name),
			Buckets: prometheus.DefBuckets,
		},
		labelNames,
	)

	if err := pr.registry.Register(histogram); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			if existing, ok := are.ExistingCollector.(*prometheus.HistogramVec); ok {
				pr.histograms[name] = existing
				return existing
			}
		}
		logger.Warnf("[MetricReporter] Failed to register histogram %s: %v", name, err)
	}

	pr.histograms[name] = histogram
	return histogram
}

// getOrCreateGauge gets or creates a gauge metric.
func (pr *PushReporter) getOrCreateGauge(name string, labelNames []string) *prometheus.GaugeVec {
	if gauge, exists := pr.gauges[name]; exists {
		return gauge
	}

	gauge := prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Name: name,
			Help: fmt.Sprintf("Gauge metric: %s", name),
		},
		labelNames,
	)

	if err := pr.registry.Register(gauge); err != nil {
		if are, ok := err.(prometheus.AlreadyRegisteredError); ok {
			if existing, ok := are.ExistingCollector.(*prometheus.GaugeVec); ok {
				pr.gauges[name] = existing
				return existing
			}
		}
		logger.Warnf("[MetricReporter] Failed to register gauge %s: %v", name, err)
	}

	pr.gauges[name] = gauge
	return gauge
}


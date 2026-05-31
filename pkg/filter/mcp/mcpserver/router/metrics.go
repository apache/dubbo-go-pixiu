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

package router

import (
	"sync"
)

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Prometheus metrics for the MCP tool router. All series live under the
// pixiu_mcp_tool_router_* namespace and are registered once on first use so the
// router package can be imported without forcing metric registration.
var (
	metricsOnce sync.Once

	selectTotal      *prometheus.CounterVec
	selectionLatency *prometheus.HistogramVec
	candidatesCount  prometheus.Histogram
	selectedCount    prometheus.Histogram
	fallbackTotal    *prometheus.CounterVec
	callDeniedTotal  *prometheus.CounterVec
	plansActive      prometheus.Gauge
)

const (
	metricsNamespace = "pixiu"
	metricsSubsystem = "mcp_tool_router"
)

// initMetrics registers the metric collectors exactly once. promauto registers
// against the default registry; AlreadyRegisteredError is impossible here
// because of the sync.Once guard, but using promauto keeps it consistent with
// the rest of the codebase.
func initMetrics() {
	metricsOnce.Do(func() {
		selectTotal = promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Subsystem: metricsSubsystem,
			Name:      "select_total",
			Help:      "Total tool selections, partitioned by result and mode.",
		}, []string{"result", "mode"})

		selectionLatency = promauto.NewHistogramVec(prometheus.HistogramOpts{
			Namespace: metricsNamespace,
			Subsystem: metricsSubsystem,
			Name:      "selection_latency_ms",
			Help:      "Tool selection latency in milliseconds.",
			Buckets:   []float64{0.1, 0.5, 1, 2.5, 5, 10, 25, 50, 100},
		}, []string{"stage"})

		candidatesCount = promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: metricsNamespace,
			Subsystem: metricsSubsystem,
			Name:      "candidates_count",
			Help:      "Number of candidate tools before selection.",
			Buckets:   []float64{1, 5, 10, 25, 50, 100, 250, 500, 1000, 2500},
		})

		selectedCount = promauto.NewHistogram(prometheus.HistogramOpts{
			Namespace: metricsNamespace,
			Subsystem: metricsSubsystem,
			Name:      "selected_count",
			Help:      "Number of tools selected after the pipeline.",
			Buckets:   []float64{0, 1, 5, 10, 25, 50, 100, 250, 500, 1000},
		})

		fallbackTotal = promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Subsystem: metricsSubsystem,
			Name:      "fallback_total",
			Help:      "Total fallback activations, partitioned by reason.",
		}, []string{"reason"})

		callDeniedTotal = promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: metricsNamespace,
			Subsystem: metricsSubsystem,
			Name:      "call_denied_total",
			Help:      "Total tools/call denials, partitioned by reason.",
		}, []string{"reason"})

		plansActive = promauto.NewGauge(prometheus.GaugeOpts{
			Namespace: metricsNamespace,
			Subsystem: metricsSubsystem,
			Name:      "plans_active",
			Help:      "Current number of cached session plans.",
		})
	})
}

// recordSelection records the outcome of one Select call.
func recordSelection(result, mode string, candidates, selected int, totalMS float64) {
	if selectTotal == nil {
		return
	}
	selectTotal.WithLabelValues(result, mode).Inc()
	selectionLatency.WithLabelValues("total").Observe(totalMS)
	candidatesCount.Observe(float64(candidates))
	selectedCount.Observe(float64(selected))
}

// recordFallback records a fallback activation.
func recordFallback(reason string) {
	if fallbackTotal == nil {
		return
	}
	fallbackTotal.WithLabelValues(reason).Inc()
}

// recordCallDenied records a denied tools/call.
func recordCallDenied(reason string) {
	if callDeniedTotal == nil {
		return
	}
	callDeniedTotal.WithLabelValues(reason).Inc()
}

// setPlansActive publishes the current cached-plan count.
func setPlansActive(n int) {
	if plansActive == nil {
		return
	}
	plansActive.Set(float64(n))
}

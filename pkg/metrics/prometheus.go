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

package metrics

import (
	"context"
	"sync"
	"time"
)

import (
	"github.com/prometheus/client_golang/prometheus"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
)

const (
	registryName = "prometheus"

	MetricNamePrefix = "dubbo-go-proxy_"

	// service level.

	// MetricServicePrefix prefix of all service metric names.
	MetricServicePrefix       = MetricNamePrefix + "service_"
	serviceReqsTotalName      = MetricServicePrefix + "requests_total"
	serviceRetriesTotalName   = MetricServicePrefix + "retries_total"
	serviceRespsTimeMillsName = MetricServicePrefix + "responses_time"
)

var (
	promRegistry     *PromRegistry
	registryInitOnce sync.Once

	defaultBuckets = []float64{0.1, 0.3, 1.2, 5.0}
)

func init() {
	extension.SetRegistry(registryName, newPromRegistry())
}

type PromRegistry struct {
	serviceReqsTotalName      *prometheus.CounterVec
	serviceRetriesTotalName   *prometheus.CounterVec
	serviceRespsTimeMillsName *prometheus.HistogramVec
}

func (r *PromRegistry) Export(ctx context.Context, method config.Method, cost time.Duration) {
	r.serviceReqsTotalName.WithLabelValues(string(method.HTTPVerb)).Inc()
	r.serviceReqsTotalName.WithLabelValues(method.InboundRequest.RequestType).Inc()
	r.serviceReqsTotalName.WithLabelValues(method.IntegrationRequest.RequestType).Inc()

}

// newPromRegistry create new promRegistry
// it will register the metrics into prometheus
func newPromRegistry() Registry {
	// TODO
	config := config.Prometheus{}
	bucket := config.Buckets
	if bucket == nil {
		bucket = defaultBuckets
	}
	if promRegistry == nil {
		registryInitOnce.Do(func() {
			promRegistry = &PromRegistry{
				serviceReqsTotalName: newCounterForm(prometheus.CounterOpts{
					Name: serviceReqsTotalName,
					Help: "How many HTTP requests processed on a service, partitioned by status code, protocol, and method.",
				}, []string{"method", "protocol", "status", "dest_type", "service"}),
				serviceRetriesTotalName: newCounterForm(prometheus.CounterOpts{
					Name: serviceRetriesTotalName,
					Help: "How many request retries happened on a service.",
				}, []string{"method", "protocol", "dest_type", "service"}),
				serviceRespsTimeMillsName: newHistogramForm(prometheus.HistogramOpts{
					Name:    serviceRespsTimeMillsName,
					Help:    "http_response_time_milliseconds.",
					Buckets: bucket,
				}, []string{"method", "status", "dest_type", "service"}),
			}
		})

		doRegistry(promRegistry)
	}

	return promRegistry
}

func doRegistry(registry *PromRegistry) {
	var c []prometheus.Collector

	if registry.serviceReqsTotalName != nil {
		c = append(c, registry.serviceReqsTotalName)
	}
	if registry.serviceRetriesTotalName != nil {
		c = append(c, registry.serviceRetriesTotalName)
	}
	if registry.serviceRespsTimeMillsName != nil {
		c = append(c, registry.serviceRespsTimeMillsName)
	}

	prometheus.MustRegister(c...)
}

func newCounterForm(opts prometheus.CounterOpts, labelNames []string) *prometheus.CounterVec {
	return prometheus.NewCounterVec(opts, labelNames)
}

func newSummaryForm(opts prometheus.SummaryOpts, labelNames []string) *prometheus.SummaryVec {
	return prometheus.NewSummaryVec(opts, labelNames)
}

func newHistogramForm(opts prometheus.HistogramOpts, labelNames []string) *prometheus.HistogramVec {
	return prometheus.NewHistogramVec(opts, labelNames)
}

func newGuageForm(opts prometheus.GaugeOpts, labelNames []string) *prometheus.GaugeVec {
	return prometheus.NewGaugeVec(opts, labelNames)
}

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

package collector

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"path"
)

import (
	"github.com/prometheus/client_golang/prometheus"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

var (
	frontendLabelNames = []string{"frontend"}
	backendLabelNames  = []string{"backend"}
	serverLabelNames   = []string{"backend", "server"}
	serverMetrics      = []*metric{
		newServerMetric("current_queue", "Current number of queued requests assigned to this server.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("max_queue", "Maximum observed number of queued requests assigned to this server.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("current_sessions", "Current number of active sessions.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("max_sessions", "Maximum observed number of active sessions.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("limit_sessions", "Configured session limit.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("sessions_total", "Total number of sessions.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("bytes_in_total", "Current total of incoming bytes.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("bytes_out_total", "Current total of outgoing bytes.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("connection_errors_total", "Total of connection errors.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("response_errors_total", "Total of response errors.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("retry_warnings_total", "Total of retry warnings.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("redispatch_warnings_total", "Total of redispatch warnings.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("up", "Current health status of the server (1 = UP, 0 = DOWN).", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("weight", "Current weight of the server.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("check_failures_total", "Total number of failed health checks.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("downtime_seconds_total", "Total downtime in seconds.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("server_selected_total", "Total number of times a server was selected, either for new sessions, or when re-dispatching.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("current_session_rate", "Current number of sessions per second over last elapsed second.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("max_session_rate", "Maximum observed number of sessions per second.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("check_duration_seconds", "Previously run health check duration, in seconds", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("http_responses_total", "Total of HTTP responses.", prometheus.CounterValue, prometheus.Labels{"code": "1xx"}, labelValueFunc, typeValueFunc),
		newServerMetric("http_responses_total", "Total of HTTP responses.", prometheus.CounterValue, prometheus.Labels{"code": "2xx"}, labelValueFunc, typeValueFunc),
		newServerMetric("http_responses_total", "Total of HTTP responses.", prometheus.CounterValue, prometheus.Labels{"code": "3xx"}, labelValueFunc, typeValueFunc),
		newServerMetric("http_responses_total", "Total of HTTP responses.", prometheus.CounterValue, prometheus.Labels{"code": "4xx"}, labelValueFunc, typeValueFunc),
		newServerMetric("http_responses_total", "Total of HTTP responses.", prometheus.CounterValue, prometheus.Labels{"code": "5xx"}, labelValueFunc, typeValueFunc),
		newServerMetric("http_responses_total", "Total of HTTP responses.", prometheus.CounterValue, prometheus.Labels{"code": "other"}, labelValueFunc, typeValueFunc),
		newServerMetric("client_aborts_total", "Total number of data transfers aborted by the client.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("server_aborts_total", "Total number of data transfers aborted by the server.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("http_queue_time_average_seconds", "Avg. HTTP queue time for last 1024 successful connections.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("http_connect_time_average_seconds", "Avg. HTTP connect time for last 1024 successful connections.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("http_response_time_average_seconds", "Avg. HTTP response time for last 1024 successful connections.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newServerMetric("http_total_time_average_seconds", "Avg. HTTP total time for last 1024 successful connections.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
	}

	frontendMetrics = []*metric{
		newFrontendMetric("current_sessions", "Current number of active sessions.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newFrontendMetric("max_sessions", "Maximum observed number of active sessions.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newFrontendMetric("limit_sessions", "Configured session limit.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newFrontendMetric("sessions_total", "Total number of sessions.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newFrontendMetric("bytes_in_total", "Current total of incoming bytes.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newFrontendMetric("bytes_out_total", "Current total of outgoing bytes.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newFrontendMetric("requests_denied_total", "Total of requests denied for security.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newFrontendMetric("request_errors_total", "Total of request errors.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newFrontendMetric("current_session_rate", "Current number of sessions per second over last elapsed second.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newFrontendMetric("limit_session_rate", "Configured limit on new sessions per second.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newFrontendMetric("max_session_rate", "Maximum observed number of sessions per second.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newFrontendMetric("http_responses_total", "Total of HTTP responses.", prometheus.CounterValue, prometheus.Labels{"code": "1xx"}, labelValueFunc, typeValueFunc),
		newFrontendMetric("http_responses_total", "Total of HTTP responses.", prometheus.CounterValue, prometheus.Labels{"code": "2xx"}, labelValueFunc, typeValueFunc),
		newFrontendMetric("http_responses_total", "Total of HTTP responses.", prometheus.CounterValue, prometheus.Labels{"code": "3xx"}, labelValueFunc, typeValueFunc),
		newFrontendMetric("http_responses_total", "Total of HTTP responses.", prometheus.CounterValue, prometheus.Labels{"code": "4xx"}, labelValueFunc, typeValueFunc),
		newFrontendMetric("http_responses_total", "Total of HTTP responses.", prometheus.CounterValue, prometheus.Labels{"code": "5xx"}, labelValueFunc, typeValueFunc),
		newFrontendMetric("http_responses_total", "Total of HTTP responses.", prometheus.CounterValue, prometheus.Labels{"code": "other"}, labelValueFunc, typeValueFunc),
		newFrontendMetric("http_requests_total", "Total HTTP requests.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newFrontendMetric("compressor_bytes_in_total", "Number of HTTP response bytes fed to the compressor", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newFrontendMetric("compressor_bytes_out_total", "Number of HTTP response bytes emitted by the compressor", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newFrontendMetric("compressor_bytes_bypassed_total", "Number of bytes that bypassed the HTTP compressor", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newFrontendMetric("http_responses_compressed_total", "Number of HTTP responses that were compressed", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newFrontendMetric("connections_total", "Total number of connections", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
	}

	backendMetrics = []*metric{
		newBackendMetric("current_queue", "Current number of queued requests not assigned to any server.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("max_queue", "Maximum observed number of queued requests not assigned to any server.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("current_sessions", "Current number of active sessions.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("max_sessions", "Maximum observed number of active sessions.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("limit_sessions", "Configured session limit.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("sessions_total", "Total number of sessions.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("bytes_in_total", "Current total of incoming bytes.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("bytes_out_total", "Current total of outgoing bytes.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("connection_errors_total", "Total of connection errors.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("response_errors_total", "Total of response errors.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("retry_warnings_total", "Total of retry warnings.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("redispatch_warnings_total", "Total of redispatch warnings.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("up", "Current health status of the backend (1 = UP, 0 = DOWN).", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("weight", "Total weight of the servers in the backend.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("current_server", "Current number of active servers", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("server_selected_total", "Total number of times a server was selected, either for new sessions, or when re-dispatching.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("current_session_rate", "Current number of sessions per second over last elapsed second.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("max_session_rate", "Maximum number of sessions per second.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("http_responses_total", "Total of HTTP responses.", prometheus.CounterValue, prometheus.Labels{"code": "1xx"}, labelValueFunc, typeValueFunc),
		newBackendMetric("http_responses_total", "Total of HTTP responses.", prometheus.CounterValue, prometheus.Labels{"code": "2xx"}, labelValueFunc, typeValueFunc),
		newBackendMetric("http_responses_total", "Total of HTTP responses.", prometheus.CounterValue, prometheus.Labels{"code": "3xx"}, labelValueFunc, typeValueFunc),
		newBackendMetric("http_responses_total", "Total of HTTP responses.", prometheus.CounterValue, prometheus.Labels{"code": "4xx"}, labelValueFunc, typeValueFunc),
		newBackendMetric("http_responses_total", "Total of HTTP responses.", prometheus.CounterValue, prometheus.Labels{"code": "5xx"}, labelValueFunc, typeValueFunc),
		newBackendMetric("http_responses_total", "Total of HTTP responses.", prometheus.CounterValue, prometheus.Labels{"code": "other"}, labelValueFunc, typeValueFunc),
		newBackendMetric("client_aborts_total", "Total number of data transfers aborted by the client.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("server_aborts_total", "Total number of data transfers aborted by the server.", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("compressor_bytes_in_total", "Number of HTTP response bytes fed to the compressor", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("compressor_bytes_out_total", "Number of HTTP response bytes emitted by the compressor", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("compressor_bytes_bypassed_total", "Number of bytes that bypassed the HTTP compressor", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("http_responses_compressed_total", "Number of HTTP responses that were compressed", prometheus.CounterValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("http_queue_time_average_seconds", "Avg. HTTP queue time for last 1024 successful connections.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("http_connect_time_average_seconds", "Avg. HTTP connect time for last 1024 successful connections.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("http_response_time_average_seconds", "Avg. HTTP response time for last 1024 successful connections.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
		newBackendMetric("http_total_time_average_seconds", "Avg. HTTP total time for last 1024 successful connections.", prometheus.GaugeValue, nil, labelValueFunc, typeValueFunc),
	}
)

var typeValueFunc = func(stat interface{}) float64 {
	return float64(1)
}

var labelValueFunc = func(stat interface{}) []string {
	switch x := stat.(type) {
	case FrontendStat:
		return []string{

			x.Name,
		}
	case BackendStat:
		return []string{
			x.Name,
		}
	case ServerStat:
		return []string{
			x.Backend,
			x.Server,
		}
	default:
		return []string{}
	}
}

type LabelValuesFunc func(stat interface{}) []string

type TypeValueFunc = func(stat interface{}) float64

type metric struct {
	Type   prometheus.ValueType
	Desc   *prometheus.Desc
	Value  func(dataStreamStats interface{}) float64
	Labels func(dataStreamStats interface{}) []string
}

func newFrontendMetric(metricName string, helpString string, t prometheus.ValueType, constLabels prometheus.Labels, labelValueFunc LabelValuesFunc, typeValueFunc TypeValueFunc) *metric {
	return &metric{
		Type: t,
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "frontend", metricName),
			helpString,
			frontendLabelNames,
			constLabels,
		),
		Value:  typeValueFunc,
		Labels: labelValueFunc,
	}
}

func newBackendMetric(metricName string, helpString string, t prometheus.ValueType, constLabels prometheus.Labels, LabelValues LabelValuesFunc, typeValueFunc TypeValueFunc) *metric {
	return &metric{
		Type: t,
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "backend", metricName),
			helpString,
			backendLabelNames,
			constLabels,
		),
		Value:  typeValueFunc,
		Labels: labelValueFunc,
	}
}

func newServerMetric(metricName string, helpString string, t prometheus.ValueType, constLabels prometheus.Labels, labelValueFunc LabelValuesFunc, typeValueFunc TypeValueFunc) *metric {
	return &metric{
		Type: t,
		Desc: prometheus.NewDesc(
			prometheus.BuildFQName(Namespace, "server", metricName),
			helpString,
			serverLabelNames,
			constLabels,
		),
		Value:  typeValueFunc,
		Labels: labelValueFunc,
	}
}

type CommonMetricExporter struct {
	logger logger.Logger
	client *http.Client
	url    *url.URL

	up              prometheus.Gauge
	totalScrapes    prometheus.Counter
	parseFailures   prometheus.Counter
	serverMetrics   []*metric
	frontendMetrics []*metric
	backendMetrics  []*metric
}

func NewCommonMetricExporter(logger logger.Logger, client *http.Client, url *url.URL) prometheus.Collector {

	return &CommonMetricExporter{
		logger: logger,
		client: client,
		url:    url,

		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(Namespace, "common_exporter", "up"),
			Help: "Was the last scrape of Pixiu gatway successful.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(Namespace, "common_exporter", "total_scrapes"),
			Help: "Current total Pixiu  scrapes.",
		}),
		parseFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(Namespace, "common_exporter", "parse_failures"),
			Help: "Number of errors while parsing.",
		}),
		serverMetrics:   serverMetrics,
		frontendMetrics: frontendMetrics,
		backendMetrics:  backendMetrics,
	}
}

func (ce *CommonMetricExporter) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range ce.frontendMetrics {
		ch <- metric.Desc
	}
	for _, metric := range ce.serverMetrics {
		ch <- metric.Desc
	}
	for _, metric := range ce.backendMetrics {
		ch <- metric.Desc
	}

	ch <- ce.up.Desc()
	ch <- ce.totalScrapes.Desc()
	ch <- ce.parseFailures.Desc()

}

func (ce *CommonMetricExporter) Collect(ch chan<- prometheus.Metric) {
	ce.totalScrapes.Inc()
	defer func() {
		ch <- ce.up
		ch <- ce.totalScrapes
		ch <- ce.parseFailures
	}()

	statsResp, err := ce.FetchAndDecodeStats()
	if err != nil {
		ce.up.Set(0)
		ce.logger.Infof("watching (msg:{%s}) = error{%v}", "Failed to fetch and decode CommonMetric Stat", err.Error())
		return
	} else {
		ce.parseFailures.Inc()
	}
	ce.up.Set(1)

	for _, metric := range ce.frontendMetrics {
		ch <- prometheus.MustNewConstMetric(
			metric.Desc,
			metric.Type,
			metric.Value(statsResp.FrontedStats),
			metric.Labels(statsResp.FrontedStats)...,
		)

	}
	for _, metric := range ce.backendMetrics {
		ch <- prometheus.MustNewConstMetric(
			metric.Desc,
			metric.Type,
			metric.Value(statsResp.BackendStats),
			metric.Labels(statsResp.BackendStats)...,
		)

	}
	for _, metric := range ce.serverMetrics {
		ch <- prometheus.MustNewConstMetric(
			metric.Desc,
			metric.Type,
			metric.Value(statsResp.ServerStats),
			metric.Labels(statsResp.ServerStats)...,
		)
	}
}

type StatsResponse struct {
	BackendStats BackendStat  `json:"backend_stats"`
	FrontedStats FrontendStat `json:"fronted_stats"`
	ServerStats  ServerStat   `json:"server_stats"`
}

type FrontendStat struct {
	Name string `json:"frontname"`
}

type BackendStat struct {
	Name string `json:"backendname"`
}

type ServerStat struct {
	Backend string `json:"backend"`
	Server  string `json:"server"`
}

func (ce *CommonMetricExporter) FetchAndDecodeStats() (StatsResponse, error) {
	var sr StatsResponse
	u := *ce.url
	u.Path = path.Join(u.Path, "/_common/health")
	res, err := ce.client.Get(u.String())
	if err != nil {
		return sr, fmt.Errorf("failed to get stats  from %s://%s:%s%s: %s",
			u.Scheme, u.Hostname(), u.Port(), u.Path, err)
	}
	defer func() {
		err = res.Body.Close()
		if err != nil {
			ce.logger.Infof("watching (msg:{%s}) = error{%v}", "failed to close http.Client", err.Error())
		}
	}()

	if res.StatusCode != http.StatusOK {
		return sr, fmt.Errorf("HTTP Request failed with code %d", res.StatusCode)
	}

	bts, err := ioutil.ReadAll(res.Body)

	if err != nil {
		ce.parseFailures.Inc()
		return sr, err
	}
	if err := json.Unmarshal(bts, &sr); err != nil {
		ce.parseFailures.Inc()
		return sr, err
	}
	return sr, nil
}

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
	"github.com/apache/dubbo-go-pixiu/pkg/metrics/global"
)

type ApiStatsResponse struct {
	ApiStats []ApiStat `json:"api_stats"`
}

type ApiStatShare struct {
	Total      int64 `json:"total"`
	Successful int64 `json:"successful"`
	Failed     int64 `json:"failed"`
}

type ApiStat struct {
	ApiName            string `json:"api_name"`
	ApiRequests        int64  `json:"api_requests"`
	ApiRequestsLatency int64  `json:"api_requests_latency"`
}

type apiMetric struct {
	Type   prometheus.ValueType
	Desc   *prometheus.Desc
	Value  func(dataStreamStats ApiStat) float64
	Labels func(dataStreamStats ApiStat) []string
}

var (
	defaultApiMetricLabels      = []string{"api"}
	defaultApiMetricLabelValues = func(apiStat ApiStat) []string {
		return []string{apiStat.ApiName}
	}
)

type ApiHealth struct {
	logger logger.Logger
	client *http.Client
	url    *url.URL

	up                prometheus.Gauge
	totalScrapes      prometheus.Counter
	jsonParseFailures prometheus.Counter

	apiMetrics []*apiMetric
}

func NewApiHealth(logger logger.Logger, client *http.Client, url *url.URL) prometheus.Collector {
	return &ApiHealth{
		logger: logger,
		client: client,
		url:    url,

		up: prometheus.NewGauge(prometheus.GaugeOpts{
			Name: prometheus.BuildFQName(global.Namespace, "api_stats", "up"),
			Help: "Was the last scrape of the Pixiu Api Stat Data successful.",
		}),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(global.Namespace, "api_stats", "total_scrapes"),
			Help: "Current total Pixiu Api scrapes.",
		}),
		jsonParseFailures: prometheus.NewCounter(prometheus.CounterOpts{
			Name: prometheus.BuildFQName(global.Namespace, "api_stats", "json_parse_failures"),
			Help: "Number of errors while parsing JSON.",
		}),
		apiMetrics: []*apiMetric{
			{
				Type: prometheus.CounterValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(global.Namespace, "api_stats", "api_requests_total"),
					"Number of Api Requests",
					defaultApiMetricLabels, nil,
				),
				Value: func(apiStats ApiStat) float64 {
					return float64(apiStats.ApiRequests)
				},
				Labels: defaultApiMetricLabelValues,
			},
			{
				Type: prometheus.CounterValue,
				Desc: prometheus.NewDesc(
					prometheus.BuildFQName(global.Namespace, "api_stats", "api_requests_latency"),
					"Api Requests Latency ",
					defaultApiMetricLabels, nil,
				),
				Value: func(apiStats ApiStat) float64 {
					return float64(apiStats.ApiRequestsLatency)
				},
				Labels: defaultApiMetricLabelValues,
			},
		},
	}
}

func (ds *ApiHealth) Describe(ch chan<- *prometheus.Desc) {
	for _, metric := range ds.apiMetrics {
		ch <- metric.Desc
	}

	ch <- ds.up.Desc()
	ch <- ds.totalScrapes.Desc()
	ch <- ds.jsonParseFailures.Desc()
}

func (ds *ApiHealth) FetchAndDecodeStats() (ApiStatsResponse, error) {
	var dsr ApiStatsResponse

	u := *ds.url
	u.Path = path.Join(u.Path, "/_api/health")
	res, err := ds.client.Get(u.String())
	if err != nil {
		return dsr, fmt.Errorf("failed to get api stats health from %s://%s:%s%s: %s",
			u.Scheme, u.Hostname(), u.Port(), u.Path, err)
	}

	defer func() {
		err = res.Body.Close()
		if err != nil {
			ds.logger.Infof("watching (msg:{%s}) = error{%v}", "failed to close http.Client", err.Error())
		}
	}()

	if res.StatusCode != http.StatusOK {
		return dsr, fmt.Errorf("HTTP Request failed with code %d", res.StatusCode)
	}

	bts, err := ioutil.ReadAll(res.Body)
	if err != nil {
		ds.jsonParseFailures.Inc()
		return dsr, err
	}

	if err := json.Unmarshal(bts, &dsr); err != nil {
		ds.jsonParseFailures.Inc()
		return dsr, err
	}

	return dsr, nil
}

func (ds *ApiHealth) Collect(ch chan<- prometheus.Metric) {
	ds.totalScrapes.Inc()
	defer func() {
		ch <- ds.up
		ch <- ds.totalScrapes
		ch <- ds.jsonParseFailures
	}()

	statsResp, err := ds.FetchAndDecodeStats()
	if err != nil {
		ds.up.Set(0)
		ds.logger.Infof("watching (msg:{%s}) = error{%v}", "Failed to fetch and decode Api Stats", err.Error())
		return
	}
	ds.up.Set(1)

	for _, metric := range ds.apiMetrics {
		for _, v := range statsResp.ApiStats {
			ch <- prometheus.MustNewConstMetric(
				metric.Desc,
				metric.Type,
				metric.Value(v),
				metric.Labels(v)...,
			)

		}

	}
}

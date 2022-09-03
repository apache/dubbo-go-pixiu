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

package scrapeImpl

import (
	"encoding/json"
	"strings"
)

import (
	contextHttp "github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/metrics/global"
	"github.com/prometheus/client_golang/prometheus"
)

const (
	// Subsystem.
	apihealthscraper = "api_health_scraper"
)

var (
	apilabelValueFunc = func(apiStat ApiStat) []string {
		return []string{apiStat.ApiName}
	}
	apilabels  = []string{"api"}
	apimetrics = []*apiMetric{
		{
			Type: prometheus.CounterValue,
			Desc: prometheus.NewDesc(
				prometheus.BuildFQName(global.Namespace, "api_stats", "api_requests_total"),
				"Number of Api Requests",
				apilabels, nil,
			),
			Value: func(apiStats ApiStat) float64 {
				return float64(apiStats.ApiRequests)
			},
			Labels: apilabelValueFunc,
		},
	}
)

type apiMetric struct {
	Type   prometheus.ValueType
	Desc   *prometheus.Desc
	Value  func(dataStreamStats ApiStat) float64
	Labels func(dataStreamStats ApiStat) []string
}

type ApiStatsResponse struct {
	ApiStats []ApiStat `json:"api_stats"`
}

type ApiStat struct {
	ApiName     string `json:"api_name"`
	ApiRequests int64  `json:"api_requests"`
}

type ApiScraper struct{}

func (ApiScraper) Name() string {
	return apihealthscraper
}

func (ApiScraper) Help() string {
	return "The  Scraper for collect metric about API"
}

func (ApiScraper) Scrape(ctx *contextHttp.HttpContext, ch chan<- prometheus.Metric) error {
	var data ApiStatsResponse
	path := ctx.Request.URL.RequestURI()
	if strings.HasPrefix(path, "/_api/health") {
		_ = json.NewDecoder(ctx.Request.Body).Decode(&data)
		for _, metric := range apimetrics {
			for _, v := range data.ApiStats {
				ch <- prometheus.MustNewConstMetric(
					metric.Desc,
					metric.Type,
					metric.Value(v),
					metric.Labels(v)...,
				)

			}

		}
	}
	return nil
}

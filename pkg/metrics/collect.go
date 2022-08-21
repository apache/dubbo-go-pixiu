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
	"errors"
	"net/http"
	"net/url"
	"sync"
	"time"
)

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/exporter-toolkit/web"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

var ErrNoData = errors.New("collector returned no data")

var (
	scrapeDurationDesc = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "scrape", "duration_seconds"),
		"pixiu_exporter: Duration of a collector scrape.",
		[]string{"collector"},
		nil,
	)
	scrapeSuccessDesc = prometheus.NewDesc(
		prometheus.BuildFQName(config.Namespace, "scrape", "success"),
		"pixiu_exporter: Whether a collector succeeded.",
		[]string{"collector"},
		nil,
	)
)
var (
	factoryFuncs        = make(map[string]FactoryFunc)
	initiatedCollectors = make(map[string]prometheus.Collector)
	collectorState      = make(map[string]bool)
)

type FactoryFunc func(logger logger.Logger, client *http.Client, url *url.URL) prometheus.Collector

type PixiuMetricCollector struct {
	Collectors map[string]prometheus.Collector
	logger     logger.Logger
	pURL       *url.URL
	httpClient *http.Client
}

type Option func(*PixiuMetricCollector) error

func WithPixiuURL(URL *url.URL) Option {
	return func(e *PixiuMetricCollector) error {
		e.pURL = URL
		return nil
	}
}

func WithHTTPClient(hc *http.Client) Option {
	return func(e *PixiuMetricCollector) error {
		e.httpClient = hc
		return nil
	}
}

func NewPixiuMetricCollector(logger logger.Logger, filters []string, options ...Option) *PixiuMetricCollector {
	e := &PixiuMetricCollector{logger: logger}
	for _, o := range options {
		if err := o(e); err != nil {
			return nil
		}
	}

	for key, enabled := range collectorState {
		if !enabled {

			continue
		} else {
			initiatedCollectors[key] = factoryFuncs[key](logger, e.httpClient, e.pURL)

		}

	}
	e.Collectors = initiatedCollectors
	return e
}

func (e PixiuMetricCollector) Collect(ch chan<- prometheus.Metric) {
	wg := sync.WaitGroup{}
	wg.Add(len(e.Collectors))
	for name, c := range e.Collectors {
		go func(name string, c prometheus.Collector, ch chan<- prometheus.Metric, logger logger.Logger) {
			execute(name, c, ch, e.logger)
			wg.Done()
		}(name, c, ch, e.logger)
	}
	wg.Wait()
}

func (e PixiuMetricCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- scrapeDurationDesc
	ch <- scrapeSuccessDesc
}

func execute(name string, c prometheus.Collector, ch chan<- prometheus.Metric, logger logger.Logger) {
	begin := time.Now()
	c.Collect(ch)
	duration := time.Since(begin)
	var success float64 = 1
	ch <- prometheus.MustNewConstMetric(scrapeDurationDesc, prometheus.GaugeValue, duration.Seconds(), name)
	ch <- prometheus.MustNewConstMetric(scrapeSuccessDesc, prometheus.GaugeValue, success, name)
}

func IsNoDataError(err error) bool {
	return err == ErrNoData
}

func RegisterCollector(name string, isDefaultEnabled bool, createFunc FactoryFunc) {
	collectorState[name] = isDefaultEnabled
	factoryFuncs[name] = createFunc
}

func (e *PixiuMetricCollector) MustRegisterAllCollector() {
	prometheus.MustRegister(e)
	for _, v := range e.Collectors {
		prometheus.MustRegister(v)
	}
}

func NewHttpHandler(p *PixiuMetricCollector) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		registry := prometheus.NewRegistry()
		registry.MustRegister(p)
		gatherers := prometheus.Gatherers{
			registry,
		}
		h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	}
}

func ExporterMetrics() {
	httpClient := http.DefaultClient
	log := logger.GetLogger()
	url, _ := url.Parse(config.URI)
	if config.DefaultEnabled {
		RegisterCollector("api-health", true, NewApiHealth)
		RegisterCollector("cluster-health", true, NewClusterHealth)
		RegisterCollector("common-health", true, NewCommonMetricExporter)
		exporter := NewPixiuMetricCollector(
			log,
			[]string{},
			WithPixiuURL(url),
			WithHTTPClient(httpClient),
		)
		handlerFunc := NewHttpHandler(exporter)
		http.Handle(config.MetricsPath, promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, handlerFunc))
		srv := &http.Server{Addr: config.ListenAddress}
		web.ListenAndServe(srv, "", promlog.New(&promlog.Config{}))

	}
}

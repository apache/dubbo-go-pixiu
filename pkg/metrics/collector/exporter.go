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
	"net/http"
	"sync"
)

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/exporter-toolkit/web"
)

import (
	contextHttp "github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/metrics/global"
	"github.com/apache/dubbo-go-pixiu/pkg/metrics/scrape"
)

type Metrics struct {
	ExporterUp prometheus.Gauge
}

func NewMetrics() Metrics {
	return Metrics{
		ExporterUp: prometheus.NewGauge(prometheus.GaugeOpts{
			Namespace: global.Namespace,
			Name:      "up",
			Help:      "Whether the datacenter is up.",
		}),
	}
}

type Exporter struct {
	ctx      *contextHttp.HttpContext
	scrapers []scrape.Scraper
	metrics  Metrics
}

func New(ctx *contextHttp.HttpContext, metrics Metrics, scrapers []scrape.Scraper) *Exporter {
	return &Exporter{
		ctx:      ctx,
		scrapers: scrapers,
		metrics:  metrics,
	}
}

func (e *Exporter) Describe(ch chan<- *prometheus.Desc) {
	ch <- e.metrics.ExporterUp.Desc()
}

func (e *Exporter) Collect(ch chan<- prometheus.Metric) {
	e.scrape(e.ctx, ch)
	ch <- e.metrics.ExporterUp
}

func (e *Exporter) scrape(ctx *contextHttp.HttpContext, ch chan<- prometheus.Metric) {
	e.metrics.ExporterUp.Set(1)
	var wg sync.WaitGroup
	defer wg.Wait()

	for _, scraper := range e.scrapers {
		wg.Add(1)
		go func(scraper scrape.Scraper) {
			defer wg.Done()
			if err := scraper.Scrape(ctx, ch); err != nil {
				logger.Debugf("[PrometheusMetricError] watching (msg:{%s}) = error{%v}", "Failed To Scraper Data", err.Error())
			}
		}(scraper)
	}
}

type Config struct {
	Enable       bool
	MeticPath    string
	ExporterPort string
}

func ExporterServerOnHttp(ctx *contextHttp.HttpContext, cfg Config, scrapers []scrape.Scraper) {
	enabledScrapers := scrapers
	handlerFunc := newHandler(ctx, NewMetrics(), enabledScrapers)
	http.Handle(cfg.MeticPath, promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, handlerFunc))
	srv := &http.Server{Addr: cfg.ExporterPort}
	web.ListenAndServe(srv, "", promlog.New(&promlog.Config{}))
}

func newHandler(ctx *contextHttp.HttpContext, metrics Metrics, scrapers []scrape.Scraper) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		registry := prometheus.NewRegistry()
		registry.MustRegister(New(ctx, metrics, scrapers))
		gatherers := prometheus.Gatherers{
			registry,
		}
		h := promhttp.HandlerFor(gatherers, promhttp.HandlerOpts{})
		h.ServeHTTP(w, r)
	}
}

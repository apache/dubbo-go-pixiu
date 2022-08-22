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
	"strconv"
)

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"github.com/prometheus/common/promlog"
	"github.com/prometheus/exporter-toolkit/web"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	Namespace = "pixiu"
)

type MeticDriver struct {
	cfg model.Metric
	PMC *PixiuMetricCollector
}

func NewMeticDriver() *MeticDriver {
	return &MeticDriver{}
}

func InitDriver(bs *model.Bootstrap) *MeticDriver {
	driver := NewMeticDriver()
	var st model.Metric
	if bs.Metric == st {
		logger.Info("[dubbo-go-pixiu] no metric configuration in conf.yaml")
		return nil
	} else {
		driver.cfg = bs.Metric
	}
	return driver
}

func (d *MeticDriver) SetMetricCollector(hc *http.Client, loger logger.Logger) error {
	collector, err := NewMetricsExporter(hc, loger, &d.cfg)
	if err != nil {
		return err
	} else {
		d.PMC = collector
		return nil
	}
}

type MeticManager struct {
	driver    *MeticDriver
	bootstrap *model.Bootstrap
}

func CreateDefaultMeticManager(bs *model.Bootstrap) *MeticManager {
	manager := &MeticManager{
		bootstrap: bs,
	}
	manager.driver = InitDriver(bs)
	return manager
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

func (m *MeticManager) PrometheusMeticServerOnHttp() {
	handlerFunc := NewHttpHandler(m.driver.PMC)
	metricPath := m.driver.cfg.PrometheusMetricsPath
	http.Handle(metricPath, promhttp.InstrumentMetricHandler(prometheus.DefaultRegisterer, handlerFunc))
	listenAddres := ":" + strconv.Itoa(m.driver.cfg.PrometheusPort)
	srv := &http.Server{Addr: listenAddres}
	web.ListenAndServe(srv, "", promlog.New(&promlog.Config{}))
}

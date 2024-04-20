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

package server

import (
	"net/http"
	"strconv"
)

import (
	sdkprometheus "github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"go.opentelemetry.io/otel/exporters/prometheus"
	"go.opentelemetry.io/otel/metric/global"
	"go.opentelemetry.io/otel/sdk/metric"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func registerOtelMetricMeter(conf model.Metric) {
	if conf.Enable {
		exporter := prometheus.New()
		provider := metric.NewMeterProvider(metric.WithReader(exporter))

		registry := sdkprometheus.NewRegistry()
		err := registry.Register(exporter.Collector)
		if err != nil {
			logger.Errorf("register otel metric meter failed, err: %v", err)
			return
		}

		global.SetMeterProvider(provider)

		http.Handle("/", promhttp.HandlerFor(registry, promhttp.HandlerOpts{}))
		addr := ":" + strconv.Itoa(conf.PrometheusPort)
		go func() {
			logger.Info("Prometheus server running on " + addr)
			_ = http.ListenAndServe(addr, nil)
		}()
	}
}

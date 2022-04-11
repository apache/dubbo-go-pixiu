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
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync/atomic"
)

import (
	"go.opentelemetry.io/otel/exporters/prometheus"

	"go.opentelemetry.io/otel/metric/global"

	export "go.opentelemetry.io/otel/sdk/export/metric"
	"go.opentelemetry.io/otel/sdk/metric/aggregator/histogram"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	selector "go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"

	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/tracing"
)

func registerOtelMetricMeter(conf model.Metric) {
	if conf.Enable {
		config := prometheus.Config{}
		resources := resource.NewWithAttributes(
			semconv.SchemaURL,
		)
		c := controller.New(
			processor.New(
				selector.NewWithHistogramDistribution(
					histogram.WithExplicitBoundaries(config.DefaultHistogramBoundaries),
				),
				export.CumulativeExportKindSelector(),
				processor.WithMemory(true),
			),
			controller.WithResource(resources),
		)
		exporter, err := prometheus.New(config, c)
		if err != nil {
			logger.Error("failed to initialize prometheus exporter %v", err)
		}
		global.SetMeterProvider(exporter.MeterProvider())

		http.HandleFunc("/", exporter.ServeHTTP)
		addr := ":" + strconv.Itoa(conf.PrometheusPort)
		go func() {
			_ = http.ListenAndServe(addr, nil)
		}()

		logger.Info("Prometheus server running on " + addr)
	}
}

// NewTracer create tracer and need to be specified protocol.
func NewTracer(name tracing.ProtocolName) tracing.Trace {
	driver := GetTraceDriverManager().GetDriver()
	holder, ok := driver.Holders[name]
	if !ok {
		holder = &tracing.Holder{
			Tracers: make(map[string]tracing.Trace),
		}
		holder.ID = 0
		driver.Holders[name] = holder
	}

	builder := strings.Builder{}
	builder.WriteString(string(name))
	builder.WriteString("-" + fmt.Sprint(holder.ID))

	traceId := builder.String()
	tmp := driver.Tp.Tracer(traceId)
	tracer := &tracing.Tracer{
		ID:     traceId,
		Trace:  tmp,
		Holder: holder,
	}

	holder.Tracers[traceId] = tracing.TraceFactory[name](tracer)

	atomic.AddUint64(&holder.ID, 1)
	return holder.Tracers[traceId]
}

// GetTracer need to specify both the protocol and id.
func GetTracer(name tracing.ProtocolName, tracerID string) (tracing.Trace, error) {
	driver := GetTraceDriverManager().GetDriver()
	holder, ok := driver.Holders[name]
	if !ok {
		return nil, errors.New("can not find any tracer, please call NewTracer first\n")
	} else if _, ok = holder.Tracers[tracerID]; !ok {
		return nil, errors.New(fmt.Sprintf("can not find tracer %s with protocol %s\n", tracerID, name))
	} else {
		return holder.Tracers[tracerID], nil
	}
}

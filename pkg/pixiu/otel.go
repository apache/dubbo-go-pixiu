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

package pixiu

import (
	"context"
	"log"
	"time"
)

import (
	"go.opentelemetry.io/otel/exporters/stdout/stdoutmetric"
	"go.opentelemetry.io/otel/metric/global"
	controller "go.opentelemetry.io/otel/sdk/metric/controller/basic"
	processor "go.opentelemetry.io/otel/sdk/metric/processor/basic"
	"go.opentelemetry.io/otel/sdk/metric/selector/simple"
	"go.opentelemetry.io/otel/sdk/resource"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func registerOtelMetricMeter(conf model.Metric) {
	if conf.Enable {
		otelctx := context.Background()
		resources := resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("dubbo-go-pixiu"),
			semconv.ServiceVersionKey.String("0.3.0"),
		)

		metricExporter, err := stdoutmetric.New(
			stdoutmetric.WithPrettyPrint(),
		)
		if err != nil {
			log.Fatalf("failed to initialize stdoutmetric export pipeline: %v", err)
		}
		pusher := controller.New(
			processor.New(
				simple.NewWithExactDistribution(),
				metricExporter,
			),
			controller.WithResource(resources),
			controller.WithExporter(metricExporter),
			controller.WithCollectPeriod(conf.Interval*time.Second),
		)
		err = pusher.Start(otelctx)
		if err != nil {
			log.Fatalf("failed to initialize metric controller: %v", err)
		}
		global.SetMeterProvider(pusher.MeterProvider())
	}
}

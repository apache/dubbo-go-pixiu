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

package tracing

import (
	"context"
)

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/tracing/jaeger"
	"github.com/apache/dubbo-go-pixiu/pkg/tracing/otlp"
)

// Different Listeners listen to different protocol requests, and create the corresponding tracer
type ProtocolName string

const HTTPProtocol ProtocolName = "HTTP"

// Unique Name by making Id self-incrementing。
type Holder struct {
	Tracers map[string]Trace
	ID      uint64
}

// Tracers corresponding to the listening protocol are maintained by the holder
type TraceDriver struct {
	Holders map[ProtocolName]*Holder
	Tp      *sdktrace.TracerProvider
	context context.Context
}

func NewTraceDriver() *TraceDriver {
	return &TraceDriver{
		Holders: make(map[ProtocolName]*Holder),
	}
}

// Loading BootStrap configuration about trace
func (driver *TraceDriver) Init(bs *model.Bootstrap) *TraceDriver {
	config := bs.Trace
	ctx := context.Background()
	exp, err := newExporter(ctx, config)
	if err != nil {
		//TODO 错误处理
		return driver
	}
	provider := newTraceProvider(exp, config)

	otel.SetTracerProvider(provider)

	driver.Tp = provider
	return driver
}
func newExporter(ctx context.Context, cfg *model.TracerConfig) (sdktrace.SpanExporter, error) {
	// Your preferred exporter: console, jaeger, zipkin, OTLP, etc.
	switch cfg.Name {
	case "otlp":
		return otlp.NewOTLPExporter(ctx, cfg)
	default:
		return jaeger.NewJaegerExporter(cfg)
	}
}

func newTraceProvider(exp sdktrace.SpanExporter, cfg *model.TracerConfig) *sdktrace.TracerProvider {
	// The service.name attribute is required.
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String("dubbo-go-pixiu"),
	)

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource),
		sdktrace.WithSampler(newSampler(cfg.Sampler)),
	)
}

func newSampler(sample model.Sampler) sdktrace.Sampler {
	switch sample.Type {
	case "never":
		return sdktrace.NeverSample()
	case "ratio":
		return sdktrace.TraceIDRatioBased(sample.Param)
	default:
		return sdktrace.AlwaysSample()
	}
}

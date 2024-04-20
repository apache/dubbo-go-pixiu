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
	"errors"
	"fmt"
)

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/propagation"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
	"go.opentelemetry.io/otel/trace"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/tracing/jaeger"
	"github.com/apache/dubbo-go-pixiu/pkg/tracing/otlp"
)

// Different Listeners listen to different protocol requests, and create the corresponding tracer
type ProtocolName string

const HTTPProtocol ProtocolName = "HTTP"

const ServiceName = "dubbo-go-pixiu"

// Sampling Strategy
const (
	ALWAYS = "always"
	NEVER  = "never"
	RATIO  = "ratio"
)

// Exporter end
const (
	JAEGER = "jaeger"
	OTLP   = "otlp"
)

// Holder Unique Name by making Id self-incrementing。
type Holder struct {
	Tracers map[string]Trace
}

// Tracers corresponding to the listening protocol are maintained by the holder
type TraceDriver struct {
	trace.TracerProvider
}

// InitDriver loading BootStrap configuration about trace
func InitDriver(bs *model.Bootstrap) *TraceDriver {
	config := bs.Trace
	if config == nil {
		logger.Info("[dubbo-go-pixiu] no trace configuration in conf.yaml")
		return nil
	}
	ctx := context.Background()
	exp, err := newExporter(ctx, config)
	if err != nil {
		logger.Warnf("[dubbo-go-pixiu] create trace exporter failed: %v", err)
		return nil
	}
	provider := newTraceProvider(exp, config)
	otel.SetTracerProvider(provider)

	return &TraceDriver{TracerProvider: provider}
}

func newExporter(ctx context.Context, cfg *model.TracerConfig) (sdktrace.SpanExporter, error) {
	// You must specify exporter to collect traces, otherwise return nil.
	switch cfg.Name {
	case OTLP:
		return otlp.NewOTLPExporter(ctx, cfg)
	case JAEGER:
		return jaeger.NewJaegerExporter(cfg)
	default:
		return nil, errors.New("no exporter error\n")
	}
}

func newTraceProvider(exp sdktrace.SpanExporter, cfg *model.TracerConfig) *sdktrace.TracerProvider {
	// The service.name attribute is required.
	if cfg.ServiceName == "" {
		cfg.ServiceName = ServiceName
	}
	resource := resource.NewWithAttributes(
		semconv.SchemaURL,
		semconv.ServiceNameKey.String(cfg.ServiceName),
	)
	//otel.SetTextMapPropagator(propagation.TraceContext{})
	otel.SetTextMapPropagator(customJaegerInject{})
	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource),
		sdktrace.WithSampler(newSampler(cfg.Sampler)),
	)
}

type customJaegerInject struct {
	propagation.TraceContext
}

func (c customJaegerInject) Inject(ctx context.Context, carrier propagation.TextMapCarrier) {
	sc := trace.SpanContextFromContext(ctx)
	if !sc.IsValid() {
		return
	}

	// Clear all flags other than the trace-context supported sampling bit.
	flags := sc.TraceFlags() & trace.FlagsSampled

	h := fmt.Sprintf("%s:%s:%s:%s",
		sc.TraceID(),
		sc.SpanID(),
		sc.SpanID(),
		flags)
	carrier.Set("uber-trace-id", h)
}

func newSampler(sample model.Sampler) sdktrace.Sampler {
	// default sampling: always.
	switch sample.Type {
	case ALWAYS:
		return sdktrace.AlwaysSample()
	case NEVER:
		return sdktrace.NeverSample()
	case RATIO:
		return sdktrace.TraceIDRatioBased(sample.Param)
	default:
		return sdktrace.AlwaysSample()
	}
}

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
	"sync"
)

import (
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/tracing/jaeger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/tracing/otlp"
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

// Unique Name by making Id self-incrementing。
type Holder struct {
	Tracers map[string]Trace
}

// Tracers corresponding to the listening protocol are maintained by the holder
type TraceDriver struct {
	Holders sync.Map
	Tp      *sdktrace.TracerProvider
	context context.Context
}

func NewTraceDriver() *TraceDriver {
	return &TraceDriver{}
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
	driver := NewTraceDriver()
	provider := newTraceProvider(exp, config)

	otel.SetTracerProvider(provider)

	driver.Tp = provider
	return driver
}

// GetHolder return the holder of the corresponding protocol. If None, create new holder.
func (driver *TraceDriver) GetHolder(name ProtocolName) *Holder {
	holder, _ := driver.Holders.LoadOrStore(name, &Holder{
		Tracers: make(map[string]Trace),
	})
	return holder.(*Holder)
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

	return sdktrace.NewTracerProvider(
		sdktrace.WithBatcher(exp),
		sdktrace.WithResource(resource),
		sdktrace.WithSampler(newSampler(cfg.Sampler)),
	)
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

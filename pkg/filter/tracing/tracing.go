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
	"log"
	"net/http"
	"time"

	"go.opentelemetry.io/otel/trace"
)

import (
	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/context"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/filter"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/attribute"
	"go.opentelemetry.io/otel/exporters/jaeger"
	"go.opentelemetry.io/otel/sdk/resource"
	tracesdk "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.4.0"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

const (
	TracingType_Jaeger = "jaeger"

	traceName      = "http-server"
	spanNamePrefix = "HTTP SERVER"
	spanTagBody    = "body"

	jaegerTraceIDInHeader = "uber-trace-id"
)

// nolint
func Init() {
	extension.SetFilterFunc(constant.TracingFilter, tracerFilterFunc())
}

// tracerFilter is a filter for tracer
type tracerFilter struct {
	// global tracer
	waitTime time.Duration
}

func tracerFilterFunc() fc.FilterFunc {
	return New().Do()
}

func newTracerProvider(url string) (*tracesdk.TracerProvider, error) {
	// Create the Jaeger exporter
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(url)))
	if err != nil {
		return nil, err
	}
	tp := tracesdk.NewTracerProvider(
		// Always be sure to batch in production.
		tracesdk.WithBatcher(exp),
		// Record information about this application in an Resource.
		tracesdk.WithResource(resource.NewWithAttributes(
			semconv.SchemaURL,
			semconv.ServiceNameKey.String("pixiu"),
		)),
	)

	return tp, nil
}

// New create tracer filter.
func New() filter.Filter {
	tc := config.GetBootstrap().Tracing
	switch tc.Type {
	case "":
	case TracingType_Jaeger:
		tp, err := newTracerProvider(tc.URL)
		if err != nil {
			log.Fatal(err)
		}
		otel.SetTracerProvider(tp)
	default:
		panic("unsupported tracing")
	}

	return tracerFilter{}
}

// Do execute tracerFilter filter logic.
func (f tracerFilter) Do() fc.FilterFunc {
	return func(c fc.Context) {
		hc := c.(*contexthttp.HttpContext)
		spanName := spanNamePrefix + hc.Request.Method + " " + hc.Request.URL.Path
		tr := otel.Tracer(traceName)
		ctx := extractTraceCtxRequest(hc.Request)
		ctxWithTid, span := tr.Start(ctx, spanName)

		body := contexthttp.ExtractRequestBody(hc.Request)
		span.SetAttributes(attribute.Key(spanTagBody).String(string(body)))
		hc.Request = hc.Request.WithContext(ctxWithTid)
		hc.Next()
		span.End()
	}
}

func extractTraceCtxRequest(req *http.Request) context.Context {
	tidHex := req.Header.Get(jaegerTraceIDInHeader)
	if tidHex == "" {
		return req.Context()
	}
	tid, err := trace.TraceIDFromHex(tidHex)
	if err != nil {
		return req.Context()
	}
	sctx := trace.NewSpanContext(trace.SpanContextConfig{TraceID: tid})
	ctx := trace.ContextWithSpanContext(req.Context(), sctx)
	return ctx
}

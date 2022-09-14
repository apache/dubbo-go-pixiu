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
	"net/http"
)

import (
	"go.opentelemetry.io/otel/trace"
)

// Wrap the tracer provided by otel and be asked to implement the Trace interface
// to customize the Span implementation.
type Tracer struct {
	ID     string
	Trace  trace.Tracer
	Holder *Holder
}

// By inheritance we can override the default StartSpan and other methods.
type HTTPTracer struct {
	*Tracer
}

func (t *Tracer) GetID() string {
	return t.ID
}

func (t *Tracer) Close() {
	delete(t.Holder.Tracers, t.ID)
}

func (t *Tracer) StartSpan(name string, request interface{}) (context.Context, trace.Span) {
	return t.Trace.Start(request.(*http.Request).Context(), name)
}

func (t *Tracer) StartSpanFromContext(name string, ctx context.Context) (context.Context, trace.Span) {
	return t.Trace.Start(ctx, name)
}

func NewHTTPTracer(tracer Trace) Trace {
	return &HTTPTracer{
		tracer.(*Tracer),
	}
}

func (t *HTTPTracer) StartSpan(name string, request interface{}) (context.Context, trace.Span) {
	return t.Trace.Start(request.(*http.Request).Context(), name)
}

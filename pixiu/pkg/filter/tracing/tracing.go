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
	"go.opentelemetry.io/otel/trace"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/tracing"
)

const TraceIDInHeader = "dubbo-go-pixiu-trace-id"

// nolint
func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

// tracerFilter is a filter for tracer
type (
	Plugin struct {
	}

	TraceFilterFactory struct {
	}

	TraceFilter struct {
		span trace.Span
	}

	Config struct{}
)

func (ap *Plugin) Kind() string {
	return constant.TracingFilter
}

func (ap *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &TraceFilterFactory{}, nil
}

func (m *TraceFilterFactory) Config() interface{} {
	return &Config{}
}

func (m *TraceFilterFactory) Apply() error {
	return nil
}

func (mf *TraceFilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	t := &TraceFilter{}
	chain.AppendDecodeFilters(t)
	chain.AppendEncodeFilters(t)
	return nil
}

// Docode execute tracerFilter filter logic.
func (f *TraceFilter) Decode(hc *contexthttp.HttpContext) filter.FilterStatus {
	spanName := "HTTP-" + hc.Request.Method
	tr, _ := server.NewTracer(tracing.HTTPProtocol)
	ctxWithId, span := tr.StartSpanFromContext(spanName, hc.Ctx)

	hc.Request = hc.Request.WithContext(ctxWithId)
	f.span = span
	return filter.Continue
}

func (f *TraceFilter) Encode(hc *contexthttp.HttpContext) filter.FilterStatus {
	f.span.End()
	return filter.Continue
}

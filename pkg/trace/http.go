package trace

import (
	"context"
	"go.opentelemetry.io/otel/trace"
	"net/http"
)

type Tracer struct {
	Id string
	T  trace.Tracer
	H  *Holder
}

type HTTPTracer struct {
	*Tracer
}

func (t *Tracer) GetId() string {
	return t.Id
}

func (t *Tracer) Close() {
	delete(t.H.Tracers, t.Id)
}

func (t *Tracer) StartSpan(name string, request interface{}) (context.Context, trace.Span) {
	return t.T.Start(request.(*http.Request).Context(), name)
}

func (t *Tracer) StartSpanFromContext(name string, ctx context.Context) (context.Context, trace.Span) {
	return t.T.Start(ctx, name)
}

func NewHTTPTracer(tracer Trace) Trace {
	return &HTTPTracer{
		tracer.(*Tracer),
	}
}

func (t *HTTPTracer) StartSpan(name string, request interface{}) (context.Context, trace.Span) {
	return t.T.Start(request.(*http.Request).Context(), name)
}

func (t *HTTPTracer) StartSpanFromContext(name string, ctx context.Context) (context.Context, trace.Span) {
	return t.T.Start(ctx, name)
}

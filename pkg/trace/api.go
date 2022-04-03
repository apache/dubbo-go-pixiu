package trace

import (
	"context"
	"go.opentelemetry.io/otel/trace"
)

type Trace interface {
	GetId() string
	StartSpan(name string, request interface{}) (context.Context, trace.Span)
	StartSpanFromContext(name string, tx context.Context) (context.Context, trace.Span)
	Close()
}

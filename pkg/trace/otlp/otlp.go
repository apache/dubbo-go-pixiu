package otlp

import (
	"context"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace"
	"go.opentelemetry.io/otel/exporters/otlp/otlptrace/otlptracehttp"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

func NewOTLPExporter(ctx context.Context, cfg *model.TracerConfig) (sdktrace.SpanExporter, error) {
	client := otlptracehttp.NewClient()
	return otlptrace.New(ctx, client)
}

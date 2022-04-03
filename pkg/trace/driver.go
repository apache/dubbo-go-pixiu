package trace

import (
	"context"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/trace/jaeger"
	"github.com/apache/dubbo-go-pixiu/pkg/trace/otlp"
	"go.opentelemetry.io/otel"
	"go.opentelemetry.io/otel/sdk/resource"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
	semconv "go.opentelemetry.io/otel/semconv/v1.7.0"
)

// 不同listener监听的协议请求不同
type ProtocolName string

const HTTP ProtocolName = "HTTP"

type Holder struct {
	Tracers map[string]Trace
	Id      uint64
}
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
	if cfg.Name == "otlp" {
		return otlp.NewOTLPExporter(ctx, cfg)
	} else {
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

// 默认是fraction
func newSampler(sample model.Sampler) sdktrace.Sampler {
	if sample.Type == "never" {
		return sdktrace.NeverSample()
	} else if sample.Type == "always" {
		return sdktrace.AlwaysSample()
	} else {
		return sdktrace.TraceIDRatioBased(sample.Param)
	}
}

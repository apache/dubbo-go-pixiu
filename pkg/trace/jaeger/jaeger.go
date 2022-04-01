package jaeger

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/pkg/errors"
	"go.opentelemetry.io/otel/exporters/jaeger"
	sdktrace "go.opentelemetry.io/otel/sdk/trace"
)

type jaegerConfig struct {
	url string `yaml:"url" json:"url" mapstructure:"url"`
}

func NewJaegerExporter(cfg *model.TracerConfig) (sdktrace.SpanExporter, error) {
	config := &jaegerConfig{}
	if err := yaml.ParseConfig(config, cfg.Config); err != nil {
		return nil, errors.Wrap(err, "config error")
	}
	exp, err := jaeger.New(jaeger.WithCollectorEndpoint(jaeger.WithEndpoint(config.url)))
	return exp, err
}

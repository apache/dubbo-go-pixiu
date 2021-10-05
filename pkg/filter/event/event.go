package event

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.HTTPCorsFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is http filter plugin.
	Plugin struct {
	}

	// Filter is http filter instance
	Filter struct {
		cfg *Config
	}

	Config struct {
		mqType       MQType
		KafkaVersion string
		LoadBalance  string
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter() (filter.HttpFilter, error) {
	return &Filter{cfg: &Config{}}, nil
}

func (f *Filter) PrepareFilterChain(ctx *http.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *Filter) Handle(ctx *http.HttpContext) {
	// TODO handle request
}

func (f *Filter) Apply() error {
	// TODO init mq facade here
	return nil
}

func (f *Filter) Config() interface{} {
	return f.cfg
}

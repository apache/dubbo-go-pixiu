package proxyrewrite

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"regexp"
)

const (
	Kind = constant.HTTPProxyRewriteFilter
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
		cfg      *Config
		uriRegex *regexp.Regexp
	}
	Config struct {
		UriRegex []string          `yaml:"uri_regex" json:"uri_regex"`
		headers  map[string]string `yaml:"headers" json:"headers"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter() (filter.HttpFilter, error) {
	return &Filter{cfg: &Config{}}, nil
}

func (f *Filter) Config() interface{} {
	return f.cfg
}

func (f *Filter) Apply() error {
	if len(f.cfg.UriRegex) == 2 {
		f.uriRegex = regexp.MustCompile(f.cfg.UriRegex[0])
	}
	return nil
}

func (f *Filter) PrepareFilterChain(ctx *contexthttp.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *Filter) Handle(c *contexthttp.HttpContext) {
	url := c.GetUrl()
	if f.uriRegex != nil {
		newUrl := f.uriRegex.ReplaceAllString(url, f.cfg.UriRegex[1])
		logger.Infof("proxy rewrite filter change url from %s to %s", url, newUrl)
		c.SetUrl(newUrl)
	}
	if len(f.cfg.headers) > 0 {
		for k, v := range f.cfg.headers {
			logger.Infof("proxy rewrite filter add header  key %s and value %s", k, v)
			c.AddHeader(k, v)
		}
	}
	c.Next()
}

package api_config

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/http/api_config/api"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/pkg/errors"
	"net/http"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.ApiConfigFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	Plugin struct {
	}

	Filter struct {
		cfg *ApiConfigConfig
	}
)

func (ap *Plugin) Kind() string {
	return Kind
}

func (ap *Plugin) CreateFilter() (filter.HttpFilter, error) {
	return &Filter{cfg: &ApiConfigConfig{}}, nil
}

func (f *Filter) Config() interface{} {
	return f.cfg
}

func (f *Filter) Apply() error {
	initApiConfig(f.cfg)
	api.Init()
	return nil
}

func (f *Filter) PrepareFilterChain(ctx *contexthttp.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *Filter) Handle(ctx *contexthttp.HttpContext) {
	req := ctx.Request
	apiDiscSrv := api.GetMustAPIDiscoveryService(constant.LocalMemoryApiDiscoveryService)
	api, err := apiDiscSrv.GetAPI(req.URL.Path, fc.HTTPVerb(req.Method))
	if err != nil {
		ctx.WriteWithStatus(http.StatusNotFound, constant.Default404Body)
		ctx.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
		e := errors.Errorf("Requested URL %s not found", req.URL.Path)
		logger.Debug(e.Error())
		ctx.Abort()
		return
	}

	if !api.Method.OnAir {
		ctx.WriteWithStatus(http.StatusNotAcceptable, constant.Default406Body)
		ctx.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
		e := errors.Errorf("Requested API %s %s does not online", req.Method, req.URL.Path)
		logger.Debug(e.Error())
		ctx.Err = e
		ctx.Abort()
		return
	}
	ctx.API(api)
	ctx.Next()
}

// initApiConfig return value of the bool is for the judgment of whether is a api meta data error, a kind of silly (?)
func initApiConfig(cf *ApiConfigConfig) {
	initFromRemote := false
	if cf.APIMetaConfig != nil {
		if _, err := config.LoadAPIConfig(cf.APIMetaConfig); err != nil {
			logger.Warnf("load api config from etcd error:%+v", err)
		} else {
			initFromRemote = true
		}
	}

	if !initFromRemote {
		if _, err := config.LoadAPIConfigFromFile(cf.Path); err != nil {
			logger.Errorf("load api config error:%+v", err)
		}
	}
}

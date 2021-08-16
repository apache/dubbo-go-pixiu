package api_config

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/http/api_config/api"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/context"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/filter"
)

func Init() {
	extension.SetFilterFunc(constant.ApiConfigFilter, apiConfigFilterFunc())
}

func CreateFilterFactory(config interface{}, bs *model.Bootstrap) extension.FilterFactoryFunc {

	cf := config.(ApiConfigConfig)
	initApiConfig(&cf)
	api.Init()

	return func(hc *contexthttp.HttpContext) {
		hc.AppendFilterFunc(apiConfigFilterFunc())
	}
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

func apiConfigFilterFunc() fc.FilterFunc {
	return New().Do()
}

// loggerFilter is a filter for simple metric.
type ApiConfigFilter struct {
}

// New create ApiConfigFilter filter.
func New() filter.Filter {
	return &ApiConfigFilter{}
}

// Metric filter, record url and latency
func (f ApiConfigFilter) Do() fc.FilterFunc {
	return func(c fc.Context) {

		c.Next()
	}
}

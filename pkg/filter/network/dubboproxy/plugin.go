package dubboproxy

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	Kind = constant.DubboProxyFilter
)

func init() {
	filter.RegisterNetworkFilter(&Plugin{})
}

type (
	Plugin struct{}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (hp *Plugin) CreateFilter(config interface{}, bs *model.Bootstrap) (filter.NetworkFilter, error) {
	hcmc, ok := config.(*model.DubboProxyConnectionManagerConfig)
	if !ok {
		panic("CreateFilter occur some exception for the type is not suitable one.")
	}
	return CreateDubboProxyConnectionManager(hcmc, bs), nil
}

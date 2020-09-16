package api_load

import "github.com/dubbogo/dubbo-go-proxy/pkg/model"

// TODO
func init() {
	var _ ApiLoad = new(NacosApiLoader)
}

type NacosApiLoader struct {
	NacosAddress string
	ApiConfigs   []model.Api
}

type NacosApiLoaderOption func(*NacosApiLoader)

func WithNacosAddress(nacosAddress string) NacosApiLoaderOption {
	return func(opt *NacosApiLoader) {
		opt.NacosAddress = nacosAddress
	}
}

func NewNacosApiLoader(opts ...NacosApiLoaderOption) *NacosApiLoader {
	var NacosApiLoader = &NacosApiLoader{}
	for _, opt := range opts {
		opt(NacosApiLoader)
	}
	return NacosApiLoader
}

func (f *NacosApiLoader) InitLoad() (err error) {
	panic("")
}

func (f *NacosApiLoader) HotLoad() error {
	panic("")
}

func (f *NacosApiLoader) Clear() error {
	panic("")
}

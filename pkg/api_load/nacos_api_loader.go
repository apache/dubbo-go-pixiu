package api_load

import "github.com/dubbogo/dubbo-go-proxy/pkg/model"

// TODO
func init() {
	var _ ApiLoader = new(NacosApiLoader)
}

type NacosApiLoader struct {
	NacosAddress string
	ApiConfigs   []model.Api
	Prior        int
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

func (f *NacosApiLoader) GetPrior() int {
	return f.Prior
}

func (f *NacosApiLoader) GetLoadedApiConfigs() ([]model.Api, error) {
	return f.ApiConfigs, nil
}

func (f *NacosApiLoader) InitLoad() (err error) {
	panic("")
}

func (f *NacosApiLoader) HotLoad() (chan struct{}, error) {
	panic("")
}

func (f *NacosApiLoader) Clear() error {
	panic("")
}

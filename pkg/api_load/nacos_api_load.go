package api_load

import "github.com/dubbogo/dubbo-go-proxy/pkg/model"

// TODO
func init() {
	var _ ApiLoad = new(NacosApiLoader)
}

type NacosApiLoader struct {
	ApiConfigs []model.Api
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

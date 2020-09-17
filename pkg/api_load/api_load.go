package api_load

import "github.com/dubbogo/dubbo-go-proxy/pkg/model"

type ApiLoader interface {
	GetPrior() int
	GetLoadedApiConfigs() ([]model.Api, error)
	// 第一次初始化加载
	InitLoad() error
	// 后面动态更新加载
	HotLoad() (chan struct{}, error)
	// 清除所有加载的api
	Clear() error
}

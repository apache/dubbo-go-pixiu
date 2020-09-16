package api_load

import "github.com/dubbogo/dubbo-go-proxy/pkg/service"

type ApiLoad interface {
	// 第一次初始化加载
	InitLoad(service.ApiDiscoveryService) error
	// 后面动态更新加载
	HotLoad(service.ApiDiscoveryService) error
	// 清除所有加载的api
	Clear() error
}

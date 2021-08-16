package router

import "github.com/apache/dubbo-go-pixiu/pkg/model"

type RouterManager struct {
	config *model.RouteConfiguration
}

func CreateRouterManager(hcmc *model.HttpConnectionManager) *RouterManager {
	return &RouterManager{config: &hcmc.RouteConfig}
}

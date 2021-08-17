package router

import (
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type RouterManager struct {
	activeConfig *model.RouteConfiguration
}

func CreateRouterManager(hcmc *model.HttpConnectionManager) *RouterManager {
	return &RouterManager{activeConfig: &hcmc.RouteConfig}
}

func (rm *RouterManager) Route(hc *http.HttpContext) (*model.RouteAction, error) {
	return rm.activeConfig.Route(hc)
}

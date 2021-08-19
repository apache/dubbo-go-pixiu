package router

import (
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

type (
	RouterCoordinator struct {
		activeConfig *model.RouteConfiguration
	}
)

func CreateRouterCoordinator(hcmc *model.HttpConnectionManager) *RouterCoordinator {

	rc := &RouterCoordinator{activeConfig: &hcmc.RouteConfig}
	if hcmc.RouteConfig.Dynamic {
		server.GetRouterManager().AddRouterListener(rc)
	}
	return rc
}

func (rm *RouterCoordinator) Route(hc *http.HttpContext) (*model.RouteAction, error) {
	return rm.activeConfig.Route(hc)
}

func (rm *RouterCoordinator) OnAddRouter(r *model.Router) {

}

func (rm *RouterCoordinator) OnDeleteRouter(r *model.Router) {

}

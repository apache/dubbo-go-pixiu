package proxy

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
)

// DefaultHttpConnectionManager
func DefaultHttpConnectionManager() *model.HttpConnectionManager {
	return &model.HttpConnectionManager{
		RouteConfig: model.RouteConfiguration{
			Routes: []model.Router{
				{
					Match: model.RouterMatch{
						Prefix: "/api/v1",
					},
					Route: model.RouteAction{
						Cluster: constant.HeaderValueAll,
					},
				},
			},
		},
		HttpFilters: []model.HttpFilter{
			{
				Name: constant.HttpRouterFilter,
			},
		},
	}
}

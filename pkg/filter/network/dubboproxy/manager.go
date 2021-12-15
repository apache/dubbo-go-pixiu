package dubboproxy

import (
	router2 "github.com/apache/dubbo-go-pixiu/pkg/common/router"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// HttpConnectionManager network filter for http
type DubboProxyConnectionManager struct {
	config            *model.DubboProxyConnectionManagerConfig
	routerCoordinator *router2.RouterCoordinator
}

// CreateDubboProxyConnectionManager create dubbo proxy connection manager
func CreateDubboProxyConnectionManager(config *model.DubboProxyConnectionManagerConfig, bs *model.Bootstrap) *DubboProxyConnectionManager {
	hcm := &DubboProxyConnectionManager{config: config}
	hcm.routerCoordinator = router2.CreateRouterCoordinator(&config.RouteConfig)
	return hcm
}

func (dcm *DubboProxyConnectionManager) OnData(data []byte) (interface{}, int, error) {

	return nil, 0, nil
}

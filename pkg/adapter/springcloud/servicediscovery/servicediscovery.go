package servicediscovery

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type (
	// InstanceInfo the service instance info fetched from registry such as nacos consul
	InstanceInfo struct {
		ID          string
		ServiceName string
		// host:port
		Addr string
		// extra info such as label or other meta data
		Meta map[string]string
	}

	ServiceDiscovery interface {
		// 直接向远程注册中心查询所有服务实例
		QueryServices() ([]*InstanceInfo, error)

		// 注册自己
		Register() error

		// 取消注册自己
		UnRegister() error

		Stop() error
	}
)

// ToEndpoint
func (i *InstanceInfo) ToEndpoint() *model.Endpoint {
	a := model.SocketAddress{Address: i.Addr}
	return &model.Endpoint{ID: i.ID, Address: a, Name: i.ServiceName, Meta: i.Meta}
}

// ToRoute route ID is cluster name, so equal with endpoint name and routerMatch prefix is also service name
func (i *InstanceInfo) ToRoute() *model.Router {
	rm := model.RouterMatch{Prefix: i.ServiceName}
	ra := model.RouteAction{Cluster: i.ServiceName}
	return &model.Router{ID: i.ServiceName, Match: rm, Route: ra}
}

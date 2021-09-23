package servicediscovery

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go/common/observer"
)

type (
	// ServiceInstance the service instance info fetched from registry such as nacos consul
	ServiceInstance struct {
		ID          string
		ServiceName string
		// host:port
		Host        string
		Port        int
		Healthy     bool
		CLusterName string
		Enable      bool
		// extra info such as label or other meta data
		Metadata map[string]string
	}

	ServiceDiscovery interface {
		// 直接向远程注册中心查询所有服务实例
		QueryAllServices() ([]ServiceInstance, error)

		QueryServicesByName(serviceNames []string) ([]ServiceInstance, error)

		// 注册自己
		Register() error

		// 取消注册自己
		UnRegister() error

		AddListener(string)

		Stop() error
	}

	ServiceInstancesChangedListener interface {
		// OnEvent on ServiceInstancesChangedEvent the service instances change event
		OnEvent(e observer.Event) error
		GetServiceNames() []string
	}
)

// ToEndpoint
func (i *ServiceInstance) ToEndpoint() *model.Endpoint {
	a := model.SocketAddress{Address: i.Host, Port: i.Port}
	return &model.Endpoint{ID: i.ID, Address: a, Name: i.ServiceName, Metadata: i.Metadata}
}

// ToRoute route ID is cluster name, so equal with endpoint name and routerMatch prefix is also service name
func (i *ServiceInstance) ToRoute() *model.Router {
	prefix := "/" + i.ServiceName
	rm := model.RouterMatch{Prefix: prefix}
	ra := model.RouteAction{Cluster: i.ServiceName}
	return &model.Router{ID: i.ServiceName, Match: rm, Route: ra}
}

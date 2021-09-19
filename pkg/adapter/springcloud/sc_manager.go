package springcloud

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud/discovery"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

// SpringCloudConfig SpringCloud that life cycle configuration
type SpringCloudConfig struct {

	// runtime context : config properties
	boot            *model.Bootstrap
	discoveryClient discovery.DiscoveryClient
}

// SpringCloud runtime in Pixiu manager
type scManager struct {

	// runtime context : config properties
	//boot *model.Bootstrap

	scConfig *SpringCloudConfig

	Start func() *scManager

	// init all
	initAll func()

	// registryAutoConfiguration
	registryAutoConfiguration func()

	// PixiuRouter: init router , contain local config and dynamic config by remote registry
	initRouter func()

	// get service instance info by serviceId
	getServiceInfoById func(serviceId string)

	// hook
	postConstruct func()

	Stop func()
}

func (m *scManager) initCluster() {
	endpoint := &model.Endpoint{}
	endpoint.ID = "spring-cloud-producer" // 192.168.0.105:spring-cloud-producer:9000
	endpoint.Address = model.SocketAddress{
		Address: "192.168.0.105",
		Port:    9000,
	}
	cluster := &model.Cluster{}
	cluster.Name = "spring-cloud-producer"
	cluster.Lb = model.Rand
	cluster.Endpoints = []*model.Endpoint{
		endpoint,
	}
	// add cluster into manager
	cm := server.GetClusterManager()
	cm.AddCluster(cluster)

	// transform into route
	routeMatch := model.RouterMatch{
		Prefix: "/scp",
	}
	routeAction := model.RouteAction{
		Cluster: "spring-cloud-producer",
	}
	route := &model.Router{Match: routeMatch, Route: routeAction}

	server.GetRouterManager().AddRouter(route)
}

// load remote registry : nacos, consul...
func loadRouterByRemoteConfig() {

}

func loadRouterByLocalConfig() {

}

// SpringCloud runtime status
const (
	Initial = "Initial"
	Started = "Started"
	Running = "Running"
	Muting  = "Muting"
	Stopped = "Stopped"
)

var manager *scManager

var SCX = &sc{
	status: Initial,
}

type sc struct {
	status string `json:"status"`
}

func NewSpringCloudManager(config *SpringCloudConfig) *scManager {

	manager = &scManager{
		scConfig: config,
		Start: func() *scManager {
			manager.initAll()
			manager.initRouter()
			return manager
		},
		initAll: func() {

			// 1. 初始化 Register Client ：nacos，eureka，consul...
			config.discoveryClient, _ = discovery.NewEurekaClient("configs/eureka.json")

			// 1.1 开启定时自动刷新本地服务实例机制

			err := config.discoveryClient.StartPeriodicalRefresh()
			if err != nil {
				logger.Errorf("start discovery periodical refresh fail!", err)
			}
			// 1.2 注册自己，如果需要的话
			err = config.discoveryClient.Register()
			if err != nil {
				logger.Errorf("register self fail!", err)
			}
			// 2.

			//
			loadRouterByLocalConfig()
			loadRouterByRemoteConfig()
		},
		registryAutoConfiguration: func() {

			//bootstrap := manager.boot
			//resources := bootstrap.DynamicResources
			//resources
		},
		initRouter: func() {
			manager.initCluster()
		},
		getServiceInfoById: func(serviceId string) {

		},
		postConstruct: func() {

		},
	}

	// manager.Start()

	return manager
}

//type registry {
//
//}

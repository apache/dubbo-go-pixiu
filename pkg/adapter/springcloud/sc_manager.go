package springcloud

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud/discovery"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// SpringCloudConfig SpringCloud that life cycle configuration
type SpringCloudConfig struct {

	// runtime context : config properties
	boot *model.Bootstrap

	discoveryClient			discovery.Client


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
	Muting = "Muting"
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







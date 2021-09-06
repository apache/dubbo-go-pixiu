package springcloud

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// SpringCloud runtime in Pixiu manager
type scManager struct {

	// runtime context : config properties
	boot *model.Bootstrap

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

var SCX = &sc{
	status: Initial,
}

type sc struct {
	status string `json:"status"`
}

func SpringCloudManager(boot *model.Bootstrap) *scManager {
	
	return &scManager{
		boot:                      boot,
		initAll: func() {
			loadRouterByLocalConfig()
			loadRouterByRemoteConfig()
		},
		registryAutoConfiguration: func() {

		},
		initRouter: func() {

		},
		getServiceInfoById: func(serviceId string) {

		},
		postConstruct: func() {

		},
	}
}

//type registry {
//
//}







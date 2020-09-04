package extension

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/service"
)

var (
	apiDiscoveryServiceMap = map[string]service.ApiDiscoveryService{}
)

// SetApiDiscoveryService will store the @filter and @name
func SetApiDiscoveryService(name string, ads service.ApiDiscoveryService) {
	apiDiscoveryServiceMap[name] = ads
}

// GetMustApiDiscoveryService will return the service.ApiDiscoveryService
// if not found, it will panic
func GetMustApiDiscoveryService(name string) service.ApiDiscoveryService {
	if ds, ok := apiDiscoveryServiceMap[name]; ok {
		return ds
	}

	panic("api discovery service for " + name + " is not existing!")
}

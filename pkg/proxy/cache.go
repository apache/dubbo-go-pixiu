package proxy

var (
	apiCache = map[string]*Api{}

	filterFuncCache = map[string]func(Context){}
)

func GetMustFilterFunc(name string) FilterFunc {
	if ff, ok := filterFuncCache[name]; ok {
		return ff
	}

	panic("filter func for " + name + " is not existing!")
}

func AddFilterFunc(name string, ff FilterFunc) {
	filterFuncCache[name] = ff
}

var (
	apiDiscoveryService = map[string]*ApiDiscoveryService{}
)

func GetMustApiDiscoveryService(name string) *ApiDiscoveryService {
	if ff, ok := apiDiscoveryService[name]; ok {
		return ff
	}

	panic("api discovery service for " + name + " is not existing!")
}

func AddApiDiscoveryService(name string, ads *ApiDiscoveryService) {
	apiDiscoveryService[name] = ads
}

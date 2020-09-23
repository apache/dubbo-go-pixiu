package api_load

import "github.com/dubbogo/dubbo-go-proxy/pkg/model"

type ApiLoader interface {
	// every apiloader has a Priority, since remote api configurer such as nacos may cover fileapiloader.
	// and the larger priority number indicates it's apis can cover apis of lower priority number priority.
	GetPrior() int
	GetLoadedApiConfigs() ([]model.Api, error)
	// load all apis at first time
	InitLoad() error
	// watch apis when apis configured were changed.
	HotLoad() (chan struct{}, error)
	// clear all apis
	Clear() error
}

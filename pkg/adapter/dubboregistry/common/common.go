package common

import (
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"
)

type RegistryEventListener interface {
	OnAddAPI(r router.API) error
	OnDeleteRouter(r config.Resource) error
}

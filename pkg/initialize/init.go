package initialize

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter/api"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter/authority"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter/logger"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter/recovery"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter/remote"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter/response"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter/timeout"
	sa "github.com/dubbogo/dubbo-go-proxy/pkg/service/api"
)

// Run start init.
func Run() {
	filterInit()
	apiDiscoveryServiceInit()
}

func filterInit() {
	api.Init()
	authority.Init()
	logger.Init()
	recovery.Init()
	remote.Init()
	response.Init()
	timeout.Init()
}

func apiDiscoveryServiceInit() {
	sa.Init()
}

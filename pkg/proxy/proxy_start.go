package proxy

import (
	"encoding/json"
	"sync"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
	"github.com/dubbogo/dubbo-go-proxy/pkg/service"
)

type Proxy struct {
	startWG sync.WaitGroup
}

func (p *Proxy) Start() {
	conf := GetBootstrap()

	p.startWG.Add(1)

	defer func() {
		if re := recover(); re != nil {
			logger.Error(re)
			// TODO stop
		}
	}()

	p.beforeStart()

	listeners := conf.GetListeners()

	for _, s := range listeners {
		ls := ListenerService{Listener: &s}
		go ls.Start()
	}
}

func (p *Proxy) beforeStart() {
	NewDubboClient().Init()

	// TODO mock api register
	ads := GetMustApiDiscoveryService(pkg.ApiDiscoveryService_Dubbo)

	a1 := &Api{
		Name:     "/api/v1/test-dubbo/user",
		ITypeStr: "HTTP",
		OTypeStr: "DUBBO",
		Method:   "POST",
		Status:   1,
		Metadata: map[string]DubboMetadata{
			"dubbo": {
				ApplicationName: "BDTService",
				Group:           "test",
				Version:         "1.0.0",
				Interface:       "com.ikurento.user.UserProvider",
				Method:          "queryUser",
				Types: []string{
					"com.ikurento.user.User",
				},
			},
		},
	}
	a2 := &Api{
		Name:     "/api/v1/test-dubbo/getUserByName",
		ITypeStr: "HTTP",
		OTypeStr: "DUBBO",
		Method:   "POST",
		Status:   1,
		Metadata: map[string]DubboMetadata{
			"dubbo": {
				ApplicationName: "BDTService",
				Group:           "test",
				Version:         "1.0.0",
				Interface:       "com.ikurento.user.UserProvider",
				Method:          "GetUser",
				Types: []string{
					"java.lang.String",
				},
			},
		},
	}

	j1, _ := json.Marshal(a1)
	j2, _ := json.Marshal(a2)
	ads.AddApi(*service.NewDiscoveryRequest(j1))
	ads.AddApi(*service.NewDiscoveryRequest(j2))
}

func NewProxy() *Proxy {
	return &Proxy{
		startWG: sync.WaitGroup{},
	}
}

func Start(bs *model.Bootstrap) {
	logger.Infof("[dubboproxy go] start by config : %+v", bs)

	proxy := NewProxy()
	proxy.Start()

	proxy.startWG.Wait()
}

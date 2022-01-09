package triple

import (
	"context"
	"github.com/apache/dubbo-go-pixiu/pkg/filterchain"
	"github.com/apache/dubbo-go-pixiu/pkg/listener"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	tripleConstant "github.com/dubbogo/triple/pkg/common/constant"
	triConfig "github.com/dubbogo/triple/pkg/config"
	"github.com/dubbogo/triple/pkg/triple"
	"reflect"
	"sync"
)

func init() {
	listener.SetListenerServiceFactory(model.ProtocolTypeTriple, newTripleListenerService)
}

type (
	// ListenerService the facade of a listener
	TripleListenerService struct {
		listener.BaseListenerService
		server     *triple.TripleServer
		serviceMap *sync.Map
	}

	ProxyService struct {
		reqTypeMap sync.Map
		ls         *TripleListenerService
	}
)

func newTripleListenerService(lc *model.Listener, bs *model.Bootstrap) (listener.ListenerService, error) {

	fc := filterchain.CreateNetworkFilterChain(lc.FilterChain, bs)
	ls := &TripleListenerService{
		BaseListenerService: listener.BaseListenerService{
			Config:      lc,
			FilterChain: fc,
		},
	}

	opts := []triConfig.OptionFunction{
		triConfig.WithCodecType(tripleConstant.HessianCodecName),
		triConfig.WithLocation(lc.Address.SocketAddress.GetAddress()),
		triConfig.WithLogger(logger.GetLogger()),
		triConfig.WithProxyModeEnable(true),
	}

	triOption := triConfig.NewTripleOption(opts...)

	tripleService := &ProxyService{ls: ls}
	serviceMap := &sync.Map{}
	serviceMap.Store(tripleConstant.ProxyServiceKey, tripleService)
	server := triple.NewTripleServer(serviceMap, triOption)
	ls.serviceMap = serviceMap
	ls.server = server
	return ls, nil
}

func (ls *TripleListenerService) Start() error {
	ls.server.Start()
	return nil
}

func (d *ProxyService) setReqParamsTypes(methodName string, typ []reflect.Type) {
	d.reqTypeMap.Store(methodName, typ)
}

func (d *ProxyService) GetReqParamsInterfaces(methodName string) ([]interface{}, bool) {
	val, ok := d.reqTypeMap.Load(methodName)
	if !ok {
		return nil, false
	}
	typs := val.([]reflect.Type)
	reqParamsInterfaces := make([]interface{}, 0, len(typs))
	for _, typ := range typs {
		reqParamsInterfaces = append(reqParamsInterfaces, reflect.New(typ).Interface())
	}
	return reqParamsInterfaces, true
}

func (d *ProxyService) InvokeWithArgs(ctx context.Context, methodName string, arguments []interface{}) (interface{}, error) {
	return d.ls.FilterChain.OnTripleData(ctx, methodName, arguments)
}

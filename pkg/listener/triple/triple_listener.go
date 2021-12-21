package triple

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo3"
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
		server     triple.TripleServer
		serviceMap *sync.Map
	}

	UnaryService struct {
		reqTypeMap sync.Map
	}
)

func newTripleListenerService(lc *model.Listener, bs *model.Bootstrap) (listener.ListenerService, error) {

	opts := []triConfig.OptionFunction{
		triConfig.WithCodecType(tripleConstant.HessianCodecName),
		triConfig.WithLocation(lc.Address.SocketAddress.GetAddress()),
		triConfig.WithLogger(logger.GetLogger()),
	}

	triOption := triConfig.NewTripleOption(opts...)
	tripleService := &dubbo3.UnaryService{}
	serviceMap := &sync.Map{}
	serviceMap.Store(url.GetParam(constant.InterfaceKey, ""), tripleService)

	srv := triple.NewTripleServer(dp.serviceMap, triOption)

	fc := filterchain.CreateFilterChain(lc.FilterChain, bs)
	return &TripleListenerService{
		BaseListenerService: listener.BaseListenerService{
			Config:      lc,
			FilterChain: fc,
		},
	}, nil
}

func (ls *TripleListenerService) Start() error {

}

func (d *UnaryService) setReqParamsTypes(methodName string, typ []reflect.Type) {
	d.reqTypeMap.Store(methodName, typ)
}

func (d *UnaryService) GetReqParamsInterfaces(methodName string) ([]interface{}, bool) {
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

func (d *UnaryService) InvokeWithArgs(ctx context.Context, methodName string, arguments []interface{}) (interface{}, error) {

}

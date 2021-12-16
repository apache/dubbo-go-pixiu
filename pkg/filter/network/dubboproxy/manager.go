package dubboproxy

import (
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	"dubbo.apache.org/dubbo-go/v3/protocol/invocation"
	"dubbo.apache.org/dubbo-go/v3/remoting"
	hessian "github.com/apache/dubbo-go-hessian2"
	router2 "github.com/apache/dubbo-go-pixiu/pkg/common/router"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	perrors "github.com/pkg/errors"
	stdHttp "net/http"
	"reflect"
)

// HttpConnectionManager network filter for http
type DubboProxyConnectionManager struct {
	config            *model.DubboProxyConnectionManagerConfig
	routerCoordinator *router2.RouterCoordinator
	codec             *dubbo.DubboCodec
}

// CreateDubboProxyConnectionManager create dubbo proxy connection manager
func CreateDubboProxyConnectionManager(config *model.DubboProxyConnectionManagerConfig, bs *model.Bootstrap) *DubboProxyConnectionManager {
	hcm := &DubboProxyConnectionManager{config: config, codec: &dubbo.DubboCodec{}}
	hcm.routerCoordinator = router2.CreateRouterCoordinator(&config.RouteConfig)
	return hcm
}

func (dcm *DubboProxyConnectionManager) ServeHTTP(w stdHttp.ResponseWriter, r *stdHttp.Request) {
	panic("DubboProxyConnectionManager ServeHTTP shouldn't be called")
}

func (dcm *DubboProxyConnectionManager) OnDecode(data []byte) (interface{}, int, error) {

	resp, length, err := dcm.codec.Decode(data)
	// err := pkg.Unmarshal(buf, p.client)
	if err != nil {
		if perrors.Is(err, hessian.ErrHeaderNotEnough) || perrors.Is(err, hessian.ErrBodyNotEnough) {
			return nil, 0, nil
		}
		logger.Errorf("pkg.Unmarshal error:%+v", err)
		return nil, length, err
	}
	return resp, length, nil
}

func (dcm *DubboProxyConnectionManager) OnEncode(pkg interface{}) ([]byte, error) {
	res, ok := pkg.(*remoting.Response)
	if ok {
		buf, err := (dcm.codec).EncodeResponse(res)
		if err != nil {
			logger.Warnf("binary.Write(res{%#v}) = err{%#v}", res, perrors.WithStack(err))
			return nil, perrors.WithStack(err)
		}
		return buf.Bytes(), nil
	}

	req, ok := pkg.(*remoting.Request)
	if ok {
		buf, err := (dcm.codec).EncodeRequest(req)
		if err != nil {
			logger.Warnf("binary.Write(req{%#v}) = err{%#v}", res, perrors.WithStack(err))
			return nil, perrors.WithStack(err)
		}
		return buf.Bytes(), nil
	}

	logger.Errorf("illegal pkg:%+v\n, it is %+v", pkg, reflect.TypeOf(pkg))
	return nil, perrors.New("invalid rpc response")
}

func (dcm *DubboProxyConnectionManager) OnData(data interface{}) (interface{}, error) {
	_, ok := data.(*invocation.RPCInvocation)
	if !ok {
		panic("create invocation occur some exception for the type is not suitable one.")
	}

	result := protocol.RPCResult{}
	result.Rest = "success"
	return result, nil
}

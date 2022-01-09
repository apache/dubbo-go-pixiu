package dubboproxy

import (
	"context"
	dubboConstant "dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	"dubbo.apache.org/dubbo-go/v3/protocol/invocation"
	"dubbo.apache.org/dubbo-go/v3/remoting"
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	router2 "github.com/apache/dubbo-go-pixiu/pkg/common/router"
	dubbo2 "github.com/apache/dubbo-go-pixiu/pkg/context/dubbo"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/dubbogo/grpc-go/metadata"
	"github.com/go-errors/errors"
	perrors "github.com/pkg/errors"
	"reflect"
)

// HttpConnectionManager network filter for http
type DubboProxyConnectionManager struct {
	filter.EmptyNetworkFilter
	config            *model.DubboProxyConnectionManagerConfig
	routerCoordinator *router2.RouterCoordinator
	codec             *dubbo.DubboCodec
	filterManager     *DubboFilterManager
}

// CreateDubboProxyConnectionManager create dubbo proxy connection manager
func CreateDubboProxyConnectionManager(config *model.DubboProxyConnectionManagerConfig, bs *model.Bootstrap) *DubboProxyConnectionManager {
	filterManager := NewDubboFilterManager(config.DubboFilters)
	hcm := &DubboProxyConnectionManager{config: config, codec: &dubbo.DubboCodec{}, filterManager: filterManager}
	hcm.routerCoordinator = router2.CreateRouterCoordinator(&config.RouteConfig)
	return hcm
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

func (dcm *DubboProxyConnectionManager) OnTripleData(ctx context.Context, methodName string, arguments []interface{}) (interface{}, error) {
	dubboAttachment := make(map[string]interface{})
	md, ok := metadata.FromIncomingContext(ctx)
	if ok {
		for k := range md {
			dubboAttachment[k] = md.Get(k)[0]
		}
	}
	interfaceName := dubboAttachment[constant.InterfaceKey].(string)

	ra, err := dcm.routerCoordinator.RouteByPathAndName(interfaceName, methodName)

	if err != nil {
		return nil, errors.Errorf("Requested dubbo rpc invocation route not found")
	}

	len := len(arguments)
	inVArr := make([]reflect.Value, len)
	for i := 0; i < len; i++ {
		inVArr[i] = reflect.ValueOf(arguments[i])
	}
	// old invocation can't set parameterValues
	invoc := invocation.NewRPCInvocationWithOptions(invocation.WithMethodName(methodName),
		invocation.WithArguments(arguments),
		invocation.WithParameterValues(inVArr),
		invocation.WithAttachments(dubboAttachment))

	rpcContext := &dubbo2.RpcContext{}
	rpcContext.SetInvocation(invoc)
	rpcContext.SetRoute(ra)
	dcm.handleRpcInvocation(rpcContext)
	result := rpcContext.RpcResult
	return result, nil
}

func (dcm *DubboProxyConnectionManager) OnData(data interface{}) (interface{}, error) {
	old_invoc, ok := data.(*invocation.RPCInvocation)
	if !ok {
		panic("create invocation occur some exception for the type is not suitable one.")
	}
	// need reconstruct RPCInvocation ParameterValues witch is same with arguments. refer to dubbogo/common/proxy/proxy.makeDubboCallProxy
	arguments := old_invoc.Arguments()
	len := len(arguments)
	inVArr := make([]reflect.Value, len)
	for i := 0; i < len; i++ {
		inVArr[i] = reflect.ValueOf(arguments[i])
	}
	//old invocation can't set parameterValues
	invoc := invocation.NewRPCInvocationWithOptions(invocation.WithMethodName(old_invoc.MethodName()),
		invocation.WithArguments(old_invoc.Arguments()),
		invocation.WithCallBack(old_invoc.CallBack()), invocation.WithParameterValues(inVArr),
		invocation.WithAttachments(old_invoc.Attachments()))

	path := old_invoc.Attachment(dubboConstant.PathKey).(string)
	ra, err := dcm.routerCoordinator.RouteByPathAndName(path, old_invoc.MethodName())

	if err != nil {
		return nil, errors.Errorf("Requested dubbo rpc invocation %s %s route not found")
	}

	ctx := &dubbo2.RpcContext{}
	ctx.SetRoute(ra)
	ctx.SetInvocation(invoc)
	dcm.handleRpcInvocation(ctx)
	result := ctx.RpcResult
	return result, nil

	//var resp interface{}
	//invoc.SetReply(&resp)
	//
	//// todo: should use service key or path or interface
	//// argument check
	//path := invoc.Attachment(constant.PathKey).(string)
	//method := invoc.Arguments()[0].(string)
	//
	//ra, err := dcm.routerCoordinator.RouteByPathAndName(path, "GET")
	//
	//if err != nil {
	//	return nil, errors.Errorf("Requested dubbo rpc invocation %s %s route not found", path, method)
	//}
	//
	//clusterName := ra.Cluster
	//clusterManager := server.GetClusterManager()
	//endpoint := clusterManager.PickEndpoint(clusterName)
	//if endpoint == nil {
	//	return nil, errors.Errorf("Requested dubbo rpc invocation %s %s endpoint not found", path, method)
	//}
	//
	//// create URL from RpcInvocation
	//url, err := common.NewURL(endpoint.Address.GetAddress(),
	//	common.WithProtocol(tpconst.TRIPLE), common.WithParamsValue(constant.SerializationKey, constant.Hessian2Serialization),
	//	common.WithParamsValue(constant.GenericFilterKey, "true"),
	//	common.WithParamsValue(constant.AppVersionKey, "3.0.0"),
	//	common.WithParamsValue(constant.InterfaceKey, path),
	//	common.WithParamsValue(constant.ReferenceFilterKey, "generic,filter"),
	//	common.WithPath(path),
	//)
	//
	//if err != nil {
	//
	//}
	//
	//invoker, err := dubbo3.NewDubboInvoker(url)
	//
	//if err != nil {
	//
	//}
	//
	//invCtx := context.Background()
	//result := invoker.Invoke(invCtx, invoc)
	//result.SetAttachments(invoc.Attachments())
	//value := reflect.ValueOf(result.Result())
	//// todo: generic 是需要这样
	//result.SetResult(value.Elem().Interface())
	//
	//return result, nil
}

// handleRpcInvocation handle rpc request
func (dcm *DubboProxyConnectionManager) handleRpcInvocation(c *dubbo2.RpcContext) {
	filterChain := dcm.filterManager.filters

	// recover any err when filterChain run
	defer func() {
		if err := recover(); err != nil {
			logger.Warnf("[dubbopixiu go] Occur An Unexpected Err: %+v", err)
			c.SetError(errors.Errorf("Occur An Unexpected Err: %v", err))
		}
	}()

	for _, f := range filterChain {
		status := f.Handle(c)
		switch status {
		case filter.Continue:
			continue
		case filter.Stop:
			return
		}
	}
}

package dubboproxy

import (
	"context"
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo3"
	"dubbo.apache.org/dubbo-go/v3/protocol/invocation"
	"dubbo.apache.org/dubbo-go/v3/remoting"
	hessian "github.com/apache/dubbo-go-hessian2"
	router2 "github.com/apache/dubbo-go-pixiu/pkg/common/router"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
	tripleCommon "github.com/dubbogo/triple/pkg/common"
	tpconst "github.com/dubbogo/triple/pkg/common/constant"
	"github.com/go-errors/errors"
	perrors "github.com/pkg/errors"
	stdHttp "net/http"
	"reflect"
	"strings"
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

func (dcm *DubboProxyConnectionManager) OnTripleData(ctx context.Context, methodName string, arguments []interface{}) (interface{}, error) {
	dubboAttachment, _ := ctx.Value(tpconst.TripleAttachement).(tripleCommon.DubboAttachment)
	old_invoc := invocation.NewRPCInvocation(methodName, arguments, dubboAttachment)
	// /org.apache.dubbo.samples.UserProviderTriple/$invoke
	interfaceMethodName := ctx.Value("XXX_TRIPLE_GO_INTERFACE_NAME").(string)
	interfaceName := strings.Split(interfaceMethodName, "/")[1]

	path := interfaceName
	method := arguments[0] // 应该是第一个

	ra, err := dcm.routerCoordinator.RouteByPathAndName(path, "GET")

	if err != nil {
		return nil, errors.Errorf("Requested dubbo rpc invocation %s %s route not found", path, method)
	}

	clusterName := ra.Cluster
	clusterManager := server.GetClusterManager()
	endpoint := clusterManager.PickEndpoint(clusterName)
	if endpoint == nil {
		return nil, errors.Errorf("Requested dubbo rpc invocation %s %s endpoint not found", path, method)
	}

	// 需要设置hessian协议

	url, err := common.NewURL(endpoint.Address.GetAddress(),
		common.WithProtocol(dubbo.DUBBO), common.WithParamsValue(constant.SerializationKey, constant.Hessian2Serialization),
		common.WithParamsValue(constant.GenericFilterKey, "true"),
		common.WithParamsValue(constant.AppVersionKey, "2.0.2"),
		common.WithParamsValue(constant.InterfaceKey, path),
		common.WithParamsValue(constant.ReferenceFilterKey, "generic,filter"),
		common.WithPath(path),
	)

	dubboProtocol := dubbo.NewDubboProtocol()
	invoker := dubboProtocol.Refer(url)

	invCtx := context.Background()

	len := len(arguments)
	inVArr := make([]reflect.Value, len)
	for i := 0; i < len; i++ {
		inVArr[i] = reflect.ValueOf(arguments[i])
	}
	// old invocation can't set parameterValues
	invoc := invocation.NewRPCInvocationWithOptions(invocation.WithMethodName(old_invoc.MethodName()),
		invocation.WithArguments(arguments),
		invocation.WithCallBack(old_invoc.CallBack()), invocation.WithParameterValues(inVArr),
		invocation.WithAttachments(map[string]interface{}{
			"async":   "false",
			"generic": "true",
		}))
	var resp interface{}
	invoc.SetReply(&resp)

	result := invoker.Invoke(invCtx, invoc)
	result.SetAttachments(invoc.Attachments())
	value := reflect.ValueOf(result.Result())
	// todo: generic 是需要这样
	result.SetResult(value.Elem().Interface())
	result.SetAttachments(nil)
	return result, result.Error()
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
	// old invocation can't set parameterValues
	invoc := invocation.NewRPCInvocationWithOptions(invocation.WithMethodName(old_invoc.MethodName()),
		invocation.WithArguments(old_invoc.Arguments()),
		invocation.WithCallBack(old_invoc.CallBack()), invocation.WithParameterValues(inVArr),
		invocation.WithAttachments(old_invoc.Attachments()))
	// todo: 即使是nil也要设置，涉及指针问题，因为 pbwrapper的MarshalResponse会有 == nil 判断，否则判断不通过，导致返回值都是nil
	var resp interface{}
	invoc.SetReply(&resp)

	// todo: should use service key or path or interface
	// argument check
	path := invoc.Attachment(constant.PathKey).(string)
	method := invoc.Arguments()[0].(string)

	ra, err := dcm.routerCoordinator.RouteByPathAndName(path, "GET")

	if err != nil {
		return nil, errors.Errorf("Requested dubbo rpc invocation %s %s route not found", path, method)
	}

	clusterName := ra.Cluster
	clusterManager := server.GetClusterManager()
	endpoint := clusterManager.PickEndpoint(clusterName)
	if endpoint == nil {
		return nil, errors.Errorf("Requested dubbo rpc invocation %s %s endpoint not found", path, method)
	}

	// create URL from RpcInvocation
	url, err := common.NewURL(endpoint.Address.GetAddress(),
		common.WithProtocol(tpconst.TRIPLE), common.WithParamsValue(constant.SerializationKey, constant.Hessian2Serialization),
		common.WithParamsValue(constant.GenericFilterKey, "true"),
		common.WithParamsValue(constant.AppVersionKey, "3.0.0"),
		common.WithParamsValue(constant.InterfaceKey, path),
		common.WithParamsValue(constant.ReferenceFilterKey, "generic,filter"),
		common.WithPath(path),
	)

	if err != nil {

	}

	invoker, err := dubbo3.NewDubboInvoker(url)

	if err != nil {

	}

	invCtx := context.Background()
	result := invoker.Invoke(invCtx, invoc)
	result.SetAttachments(invoc.Attachments())
	value := reflect.ValueOf(result.Result())
	// todo: generic 是需要这样
	result.SetResult(value.Elem().Interface())

	return result, nil
}

/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to You under the Apache License, Version 2.0
 * (the "License"); you may not use this file except in compliance with
 * the License.  You may obtain a copy of the License at
 *
 *     http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package dubboproxy

import (
	"context"
	"reflect"
)

import (
	dubboConstant "dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	"dubbo.apache.org/dubbo-go/v3/protocol/invocation"
	"dubbo.apache.org/dubbo-go/v3/remoting"
	hessian "github.com/apache/dubbo-go-hessian2"
	"github.com/dubbogo/grpc-go/metadata"
	"github.com/go-errors/errors"
	perrors "github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	router2 "github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/router"
	dubbo2 "github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/dubbo"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

// DubboProxyConnectionManager network filter for dubbo
type DubboProxyConnectionManager struct {
	filter.EmptyNetworkFilter
	config            *model.DubboProxyConnectionManagerConfig
	routerCoordinator *router2.RouterCoordinator
	codec             *dubbo.DubboCodec
	filterManager     *DubboFilterManager
}

// CreateDubboProxyConnectionManager create dubbo proxy connection manager
func CreateDubboProxyConnectionManager(config *model.DubboProxyConnectionManagerConfig) *DubboProxyConnectionManager {
	filterManager := NewDubboFilterManager(config.DubboFilters)
	hcm := &DubboProxyConnectionManager{config: config, codec: &dubbo.DubboCodec{}, filterManager: filterManager}
	hcm.routerCoordinator = router2.CreateRouterCoordinator(&config.RouteConfig)
	return hcm
}

// OnDecode decode bytes to DecodeResult
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

// OnEncode encode Response to bytes
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

// OnTripleData handle triple rpc invocation
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
	return result, result.Error()
}

// OnData handle dubbo rpc invocation
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

	path, _ := old_invoc.GetAttachmentInterface(dubboConstant.PathKey).(string)
	ra, err := dcm.routerCoordinator.RouteByPathAndName(path, old_invoc.MethodName())

	if err != nil {
		return nil, errors.Errorf("Requested dubbo rpc invocation route not found")
	}

	ctx := &dubbo2.RpcContext{}
	ctx.SetRoute(ra)
	ctx.SetInvocation(invoc)
	dcm.handleRpcInvocation(ctx)
	result := ctx.RpcResult
	return result, nil
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

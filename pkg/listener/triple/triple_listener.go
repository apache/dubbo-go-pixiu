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

package triple

import (
	"context"
	"reflect"
	"sync"
)

import (
	tripleConstant "github.com/dubbogo/triple/pkg/common/constant"
	triConfig "github.com/dubbogo/triple/pkg/config"
	"github.com/dubbogo/triple/pkg/triple"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/filterchain"
	"github.com/apache/dubbo-go-pixiu/pkg/listener"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
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
	// ProxyService grpc proxy service definition
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

// Start start triple server
func (ls *TripleListenerService) Start() error {
	ls.server.Start()
	return nil
}

// GetReqParamsInterfaces get params
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

// InvokeWithArgs called when rpc invocation comes
func (d *ProxyService) InvokeWithArgs(ctx context.Context, methodName string, arguments []interface{}) (interface{}, error) {
	return d.ls.FilterChain.OnTripleData(ctx, methodName, arguments)
}

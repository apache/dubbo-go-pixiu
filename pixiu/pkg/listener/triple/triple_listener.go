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
	"time"
)

import (
	tripleConstant "github.com/dubbogo/triple/pkg/common/constant"
	triConfig "github.com/dubbogo/triple/pkg/config"
	"github.com/dubbogo/triple/pkg/triple"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/filterchain"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/listener"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

func init() {
	listener.SetListenerServiceFactory(model.ProtocolTypeTriple, newTripleListenerService)
}

type (
	// ListenerService the facade of a listener
	TripleListenerService struct {
		listener.BaseListenerService
		server          *triple.TripleServer
		serviceMap      *sync.Map
		gShutdownConfig *listener.ListenerGracefulShutdownConfig
	}
	// ProxyService grpc proxy service definition
	ProxyService struct {
		reqTypeMap sync.Map
		ls         *TripleListenerService
	}
)

func newTripleListenerService(lc *model.Listener, bs *model.Bootstrap) (listener.ListenerService, error) {

	fc := filterchain.CreateNetworkFilterChain(lc.FilterChain)
	ls := &TripleListenerService{
		BaseListenerService: listener.BaseListenerService{
			Config:      lc,
			FilterChain: fc,
		},
		gShutdownConfig: &listener.ListenerGracefulShutdownConfig{},
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

func (ls *TripleListenerService) Close() error {
	ls.server.Stop()
	return nil
}

func (ls *TripleListenerService) ShutDown(wg interface{}) error {
	timeout := config.GetBootstrap().GetShutdownConfig().GetTimeout()
	if timeout <= 0 {
		return nil
	}
	// stop accept request
	ls.gShutdownConfig.RejectRequest = true
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) && ls.gShutdownConfig.ActiveCount > 0 {
		// sleep 100 ms and check it again
		time.Sleep(100 * time.Millisecond)
		logger.Infof("waiting for active invocation count = %d", ls.gShutdownConfig.ActiveCount)
	}
	wg.(*sync.WaitGroup).Done()
	ls.server.Stop()
	return nil
}

func (ls *TripleListenerService) Refresh(c model.Listener) error {
	// There is no need to lock here for now, as there is at most one NetworkFilter
	fc := filterchain.CreateNetworkFilterChain(c.FilterChain)
	ls.FilterChain = fc
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
	d.ls.gShutdownConfig.AddActiveCount(1)
	defer d.ls.gShutdownConfig.AddActiveCount(-1)
	if d.ls.gShutdownConfig.RejectRequest {
		return nil, errors.Errorf("Pixiu is preparing to close, reject all new requests")
	}
	return d.ls.FilterChain.OnTripleData(ctx, methodName, arguments)
}

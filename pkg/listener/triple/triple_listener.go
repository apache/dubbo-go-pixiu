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
	"sync/atomic"
	"time"
)

import (
	tripleConstant "github.com/dubbogo/triple/pkg/common/constant"
	triConfig "github.com/dubbogo/triple/pkg/config"
	"github.com/dubbogo/triple/pkg/triple"

	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config"
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
		*listener.BaseListenerService
		server          *triple.TripleServer
		serviceMap      *sync.Map
		gShutdownConfig *listener.ListenerGracefulShutdownConfig
		started         atomic.Bool
	}
	// ProxyService grpc proxy service definition
	ProxyService struct {
		reqTypeMap sync.Map
		ls         *TripleListenerService
	}
)

func newTripleListenerService(lc *model.Listener, bs *model.Bootstrap) (listener.ListenerService, error) {

	fc, err := filterchain.BuildNetworkFilterChain(lc.FilterChain)
	if err != nil {
		return nil, err
	}
	ls := &TripleListenerService{
		BaseListenerService: listener.NewBaseListenerService(lc, fc),
		gShutdownConfig:     &listener.ListenerGracefulShutdownConfig{},
	}

	opts := []triConfig.OptionFunction{
		triConfig.WithCodecType(tripleConstant.HessianCodecName),
		triConfig.WithLocation(lc.Address.SocketAddress.GetAddress()),
		triConfig.WithLogger(logger.GetTripleLogger()),
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
func (ls *TripleListenerService) Start() (err error) {
	defer func() {
		if recovered := recover(); recovered != nil {
			err = errors.Errorf("start Triple listener: %v", recovered)
		}
	}()
	ls.server.Start()
	ls.started.Store(true)
	return nil
}

func (ls *TripleListenerService) Close() error {
	if ls.started.Swap(false) {
		ls.server.Stop()
	}
	return ls.CloseFilterChain()
}

func (ls *TripleListenerService) ShutDown(wg any) error {
	timeout := config.GetBootstrap().GetShutdownConfig().GetTimeout()
	if timeout <= 0 {
		return nil
	}
	// stop accept request
	ls.gShutdownConfig.SetRejectRequests(true)
	deadline := time.Now().Add(timeout)
	for time.Now().Before(deadline) && ls.gShutdownConfig.GetActiveCount() > 0 {
		// sleep 100 ms and check it again
		time.Sleep(100 * time.Millisecond)
		logger.Infof("waiting for active invocation count = %d", ls.gShutdownConfig.GetActiveCount())
	}
	wg.(*sync.WaitGroup).Done()
	ls.server.Stop()
	return nil
}

// GetReqParamsInterfaces get params
func (d *ProxyService) GetReqParamsInterfaces(methodName string) ([]any, bool) {
	val, ok := d.reqTypeMap.Load(methodName)
	if !ok {
		return nil, false
	}
	typs := val.([]reflect.Type)
	reqParamsInterfaces := make([]any, 0, len(typs))
	for _, typ := range typs {
		reqParamsInterfaces = append(reqParamsInterfaces, reflect.New(typ).Interface())
	}
	return reqParamsInterfaces, true
}

// InvokeWithArgs called when rpc invocation comes
func (d *ProxyService) InvokeWithArgs(ctx context.Context, methodName string, arguments []any) (any, error) {
	d.ls.gShutdownConfig.AddActiveCount(1)
	defer d.ls.gShutdownConfig.AddActiveCount(-1)
	if d.ls.gShutdownConfig.RejectRequests() {
		return nil, errors.Errorf("Pixiu is preparing to close, reject all new requests")
	}
	var result any
	err := d.ls.WithFilterChain(func(fc *filterchain.NetworkFilterChain) error {
		var invokeErr error
		result, invokeErr = fc.OnTripleData(ctx, methodName, arguments)
		return invokeErr
	})
	return result, err
}

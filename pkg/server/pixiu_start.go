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

package server

import (
	"errors"
	"fmt"
	"net/http"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/trace"
)

var server *Server

// PX is Pixiu start struct
type Server struct {
	startWG sync.WaitGroup

	listenerManager *ListenerManager
	clusterManager  *ClusterManager
	adapterManager  *AdapterManager
	// routerManager and apiConfigManager are duplicate, because route and dubbo-protocol api_config  are a bit repetitive
	routerManager         *RouterManager
	apiConfigManager      *ApiConfigManager
	dynamicResourceManger DynamicResourceManager
	traceDriverManager    *trace.TraceDriverManager
}

func (s *Server) initialize(bs *model.Bootstrap) {
	s.clusterManager = CreateDefaultClusterManager(bs)
	s.routerManager = CreateDefaultRouterManager(s, bs)
	s.apiConfigManager = CreateDefaultApiConfigManager(s, bs)
	s.adapterManager = CreateDefaultAdapterManager(s, bs)
	s.listenerManager = CreateDefaultListenerManager(bs)
	s.dynamicResourceManger = createDynamicResourceManger(bs)
	s.traceDriverManager = trace.CreateDefaultTraceDriverManager(bs)
}

func (s *Server) GetClusterManager() *ClusterManager {
	return s.clusterManager
}

func (s *Server) GetListenerManager() *ListenerManager {
	return s.listenerManager
}

func (s *Server) GetRouterManager() *RouterManager {
	return s.routerManager
}

func (s *Server) GetApiConfigManager() *ApiConfigManager {
	return s.apiConfigManager
}

func (s *Server) GetDynamicResourceManager() DynamicResourceManager {
	return s.dynamicResourceManger
}

func (s *Server) GetTraceDriverManager() *trace.TraceDriverManager {
	return s.traceDriverManager
}

// Start server start
func (s *Server) Start() {
	conf := config.GetBootstrap()

	s.startWG.Add(1)

	defer func() {
		if re := recover(); re != nil {
			logger.Error(re)
			// TODO stop
		}
	}()

	registerOtelMetricMeter(conf.Metric)
	s.listenerManager.StartListen()
	s.adapterManager.Start()

	if conf.GetPprof().Enable {
		addr := conf.GetPprof().Address.SocketAddress
		if len(addr.Address) == 0 {
			addr.Address = constant.PprofDefaultAddress
		}
		if addr.Port == 0 {
			addr.Port = constant.PprofDefaultPort
		}
		go http.ListenAndServe(addr.Address+":"+strconv.Itoa(addr.Port), nil)
		logger.Infof("[dubbopixiu go pprof] httpListener start by : %s", addr.Address+":"+strconv.Itoa(addr.Port))
	}
}

// NewServer create server
func NewServer() *Server {
	return &Server{
		startWG: sync.WaitGroup{},
	}
}

func Start(bs *model.Bootstrap) {
	logger.Infof("[dubbopixiu go] start by config : %+v", bs)
	// global variable
	server = NewServer()
	server.initialize(bs)
	server.Start()
	server.startWG.Wait()
}

func GetServer() *Server {
	return server
}

func GetClusterManager() *ClusterManager {
	return server.GetClusterManager()
}

func GetRouterManager() *RouterManager {
	return server.GetRouterManager()
}

func GetApiConfigManager() *ApiConfigManager {
	return server.GetApiConfigManager()
}

func GetDynamicResourceManager() DynamicResourceManager {
	return server.GetDynamicResourceManager()
}

func GetTraceDriverManager() *trace.TraceDriverManager {
	return server.traceDriverManager
}

func NewTracer(name trace.ProtocolName) trace.Trace {
	driver := GetTraceDriverManager().GetDriver()
	holder, ok := driver.Holders[name]
	if !ok {
		holder = &trace.Holder{
			Tracers: make(map[string]trace.Trace),
		}
		holder.Id = 0
		driver.Holders[name] = holder
	}
	// tarceId的生成，并且在协议接口唯一
	builder := strings.Builder{}
	builder.WriteString(string(name))
	builder.WriteString("-" + string(holder.Id))

	traceId := builder.String()
	tmp := driver.Tp.Tracer(traceId)
	tracer := &trace.Tracer{
		Id: traceId,
		T:  tmp,
		H:  holder,
	}

	holder.Tracers[traceId] = trace.TraceFactory[trace.HTTP](tracer)

	atomic.AddUint64(&holder.Id, 1)
	return holder.Tracers[traceId]
}

func GetTracer(name trace.ProtocolName, tracerId string) (trace.Trace, error) {
	driver := GetTraceDriverManager().GetDriver()
	holder, ok := driver.Holders[name]
	if !ok {
		return nil, errors.New("can not find any tracer, please call NewTracer first")
	} else if _, ok = holder.Tracers[tracerId]; !ok {
		return nil, errors.New(fmt.Sprintf("can not find tracer %s with protocol %s", tracerId, name))
	} else {
		return holder.Tracers[tracerId], nil
	}
}

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
	"net/http"
	"strconv"
	"sync"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	_ "github.com/apache/dubbo-go-pixiu/pkg/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

var server *Server

// PX is Pixiu start struct
type Server struct {
	startWG sync.WaitGroup

	listenerManager *ListenerManager
	clusterManager  *ClusterManager
}

func (s *Server) initialize(bs *model.Bootstrap) {
	s.listenerManager = CreateDefaultListenerManager(bs)
	s.clusterManager = CreateDefaultClusterManager(bs)
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

// NewPX create server
func NewServer() *Server {
	return &Server{
		startWG: sync.WaitGroup{},
	}
}

func Start(bs *model.Bootstrap) {
	logger.Infof("[dubbopixiu go] start by config : %+v", bs)

	server := NewServer()
	server.initialize(bs)
	server.Start()
	server.startWG.Wait()
}

func GetServer() *Server {
	return server
}

func GetClusterManager() *ClusterManager {
	return server.clusterManager
}

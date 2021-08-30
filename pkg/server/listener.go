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
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	h "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type (
	// ListenerService the facade of a listener
	ListenerService struct {
		cfg *model.Listener
		// TODO: just temporary because only one network filter
		nf  filter.NetworkFilter
		srv *http.Server
	}

	// DefaultHttpListener
	DefaultHttpWorker struct {
		pool sync.Pool
		ls   *ListenerService
	}
)

// NewListenerService create listener service
func CreateListenerService(lc *model.Listener, bs *model.Bootstrap) *ListenerService {
	hcm := createHttpManager(lc, bs)
	return &ListenerService{cfg: lc, nf: *hcm}
}

// Start start the listener
func (ls *ListenerService) Start() {
	sa := ls.cfg.Address.SocketAddress
	switch sa.Protocol {
	case model.HTTP:
		ls.httpListener()
	case model.HTTPS:
		ls.httpsListener()
	default:
		panic("unsupported protocol start: " + sa.ProtocolStr)
	}
}

func (ls *ListenerService) httpsListener() {
	hl := CreateDefaultHttpWorker(ls)
	hl.pool.New = func() interface{} {
		return ls.allocateContext()
	}
	// user customize http config
	var hc *model.HttpConfig
	hc = model.MapInStruct(ls.cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("/", hl.ServeHTTP)

	sa := ls.cfg.Address.SocketAddress
	ls.srv = &http.Server{
		Addr:           resolveAddress(sa.Address + ":" + strconv.Itoa(sa.Port)),
		Handler:        mux,
		ReadTimeout:    resolveStr2Time(hc.ReadTimeoutStr, 20*time.Second),
		WriteTimeout:   resolveStr2Time(hc.WriteTimeoutStr, 20*time.Second),
		IdleTimeout:    resolveStr2Time(hc.IdleTimeoutStr, 20*time.Second),
		MaxHeaderBytes: resolveInt2IntProp(hc.MaxHeaderBytes, 1<<20),
	}

	logger.Infof("[dubbo-go-server] httpsListener start at : %s", ls.srv.Addr)
	err := ls.srv.ListenAndServeTLS(hc.CertFile, hc.KeyFile)
	logger.Info("[dubbo-go-server] httpsListener result:", err)
}

func (ls *ListenerService) httpListener() {
	hl := CreateDefaultHttpWorker(ls)
	hl.pool.New = func() interface{} {
		return ls.allocateContext()
	}

	// user customize http config
	var hc *model.HttpConfig
	hc = model.MapInStruct(ls.cfg)

	mux := http.NewServeMux()
	mux.HandleFunc("/", hl.ServeHTTP)

	sa := ls.cfg.Address.SocketAddress
	ls.srv = &http.Server{
		Addr:           resolveAddress(sa.Address + ":" + strconv.Itoa(sa.Port)),
		Handler:        mux,
		ReadTimeout:    resolveStr2Time(hc.ReadTimeoutStr, 20*time.Second),
		WriteTimeout:   resolveStr2Time(hc.WriteTimeoutStr, 20*time.Second),
		IdleTimeout:    resolveStr2Time(hc.IdleTimeoutStr, 20*time.Second),
		MaxHeaderBytes: resolveInt2IntProp(hc.MaxHeaderBytes, 1<<20),
	}

	logger.Infof("[dubbo-go-server] httpListener start at : %s", ls.srv.Addr)

	log.Println(ls.srv.ListenAndServe())
}

func (ls *ListenerService) allocateContext() *h.HttpContext {
	return &h.HttpContext{
		Listener: ls.cfg,
	}
}

// NewDefaultHttpListener create http listener
func CreateDefaultHttpWorker(ls *ListenerService) *DefaultHttpWorker {
	return &DefaultHttpWorker{
		pool: sync.Pool{},
		ls:   ls,
	}
}

// ServeHTTP http request entrance.
func (s *DefaultHttpWorker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hc := s.pool.Get().(*h.HttpContext)
	defer s.pool.Put(hc)

	hc.Request = r
	hc.ResetWritermen(w)
	hc.Reset()

	// now only one filter http_connection_manager, so just get it and call
	err := s.ls.nf.OnData(hc)

	if err != nil {
		s.pool.Put(hc)
		logger.Errorf("ServeHTTP %s", err)
	}
}

func resolveInt2IntProp(currentV, defaultV int) int {
	if currentV == 0 {
		return defaultV
	}

	return currentV
}

func resolveStr2Time(currentV string, defaultV time.Duration) time.Duration {
	if currentV == "" {
		return defaultV
	} else {
		if duration, err := time.ParseDuration(currentV); err != nil {
			return 20 * time.Second
		} else {
			return duration
		}
	}
}

func resolveAddress(addr string) string {
	if addr == "" {
		logger.Debug("Addr is undefined. Using port :8080 by default")
		return ":8080"
	}

	return addr
}

func findHttpManager(l *model.Listener) *model.HttpConnectionManager {
	for _, fc := range l.FilterChains {
		for _, f := range fc.Filters {
			if f.Name == constant.HTTPConnectManagerFilter {
				hcmc := &model.HttpConnectionManager{}
				if err := yaml.ParseConfig(hcmc, f.Config); err != nil {
					return nil
				}

				return hcmc
			}
		}
	}

	return DefaultHttpConnectionManager()
}

func createHttpManager(lc *model.Listener, bs *model.Bootstrap) *filter.NetworkFilter {
	p, err := filter.GetNetworkFilterPlugin(constant.HTTPConnectManagerFilter)
	if err != nil {
		panic(err)
	}

	hcmc := findHttpManager(lc)
	hcm, err := p.CreateFilter(hcmc, bs)
	if err != nil {
		panic(err)
	}
	return &hcm
}

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

package http2

import (
	"net"
	"net/http"
	"strconv"
)

import (
	"golang.org/x/net/http2"
	"golang.org/x/net/http2/h2c"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/filterchain"
	"github.com/apache/dubbo-go-pixiu/pkg/listener"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func init() {
	listener.SetListenerServiceFactory(model.ProtocolTypeHTTP2, newHttp2ListenerService)
}

type (
	// Http2ListenerService the facade of a listener
	Http2ListenerService struct {
		listener.BaseListenerService
		listener net.Listener
		server   *http.Server
	}
)

type handleWrapper struct {
	fc *filterchain.NetworkFilterChain
}

type h2cWrapper struct {
	w *handleWrapper
	h http.Handler
}

// ServeHTTP call Handler to handle http request and response.
func (h *h2cWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.h.ServeHTTP(w, r)
}

// ServeHTTP call FilterChain to handle http request and response.
func (h *handleWrapper) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	h.fc.ServeHTTP(w, r)
}

func newHttp2ListenerService(lc *model.Listener, bs *model.Bootstrap) (listener.ListenerService, error) {
	fc := filterchain.CreateNetworkFilterChain(lc.FilterChain)
	return &Http2ListenerService{
		BaseListenerService: listener.BaseListenerService{
			Config:      lc,
			FilterChain: fc,
		},
		listener: nil,
		server:   nil,
	}, nil
}

// Start start listen
func (ls *Http2ListenerService) Start() error {

	sa := ls.Config.Address.SocketAddress
	addr := resolveAddress(sa.Address + ":" + strconv.Itoa(sa.Port))

	l, err := net.Listen("tcp", addr)
	if err != nil {
		return err
	}
	ls.listener = l

	handlerWrapper := &handleWrapper{ls.FilterChain}
	h2s := &http2.Server{}
	h := &h2cWrapper{
		w: handlerWrapper,
		h: h2c.NewHandler(handlerWrapper, h2s),
	}

	ls.server = &http.Server{
		Addr:    addr,
		Handler: h,
	}

	go func() {
		if err := ls.server.Serve(ls.listener); err != nil {
			if err == http.ErrServerClosed {
				logger.Infof("Listener %s closed", ls.Config.Name)
				return
			}
			logger.Errorf("Http2ListenerService.Serve: %v", err)
		}
	}()
	return nil
}

func (ls *Http2ListenerService) Close() error {
	return ls.server.Close()
}

func (ls *Http2ListenerService) ShutDown() error {
	//TODO implement me
	panic("implement me")
}

func (ls *Http2ListenerService) Refresh(c model.Listener) error {
	// There is no need to lock here for now, as there is at most one NetworkFilter
	fc := filterchain.CreateNetworkFilterChain(c.FilterChain)
	ls.FilterChain = fc
	return nil
}

func resolveAddress(addr string) string {
	if addr == "" {
		logger.Debug("Addr is undefined. Using port :8080 by default")
		return ":8080"
	}

	return addr
}

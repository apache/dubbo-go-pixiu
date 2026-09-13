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

package http

import (
	"context"
	"crypto/tls"
	"fmt"
	"net"
	"net/http"
	"strconv"
	"sync"
	"time"
)

import (
	"github.com/pkg/errors"

	"golang.org/x/crypto/acme/autocert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/filterchain"
	"github.com/apache/dubbo-go-pixiu/pkg/listener"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func init() {
	listener.SetListenerServiceFactory(model.ProtocolTypeHTTP, newHttpListenerService)
	listener.SetListenerServiceFactory(model.ProtocolTypeHTTPS, newHttpListenerService)
}

type (
	// ListenerService the facade of a listener
	HttpListenerService struct {
		*listener.BaseListenerService
		srv *http.Server
	}

	// DefaultHttpListener
	DefaultHttpWorker struct {
		ls *HttpListenerService
	}
)

func newHttpListenerService(lc *model.Listener, bs *model.Bootstrap) (listener.ListenerService, error) {
	fc, err := filterchain.BuildNetworkFilterChain(lc.FilterChain)
	if err != nil {
		return nil, errors.Wrap(err, "create HTTP listener filter chain")
	}
	return &HttpListenerService{
		BaseListenerService: listener.NewBaseListenerService(lc, fc),
		srv:                 nil,
	}, nil
}

// Start start the listener
func (ls *HttpListenerService) Start() error {
	switch ls.Config.Protocol {
	case model.ProtocolTypeHTTP:
		return ls.httpListener()
	case model.ProtocolTypeHTTPS:
		return ls.httpsListener()
	default:
		return fmt.Errorf("unsupported protocol start: %d", ls.Config.Protocol)
	}
}

func (ls *HttpListenerService) Close() error {
	var closeErr error
	if ls.srv != nil {
		closeErr = ls.srv.Close()
	}
	if filterErr := ls.CloseFilterChain(); closeErr == nil {
		closeErr = filterErr
	}
	return closeErr
}

func (ls *HttpListenerService) ShutDown(wg any) error {
	timeout := config.GetBootstrap().GetShutdownConfig().GetTimeout()
	if timeout <= 0 {
		return nil
	}
	ctx, cancel := context.WithTimeout(context.Background(), timeout)
	defer func() {
		cancel()
		wg.(*sync.WaitGroup).Done()
	}()
	return ls.srv.Shutdown(ctx)
}

func (ls *HttpListenerService) httpsListener() error {
	hl := createDefaultHttpWorker(ls)

	// user customize http config
	hc := model.MapInStruct(ls.Config.Config)

	mux := http.NewServeMux()
	mux.HandleFunc("/", hl.ServeHTTP)

	m := &autocert.Manager{
		Cache:      autocert.DirCache(ls.Config.Address.SocketAddress.CertsDir),
		Prompt:     autocert.AcceptTOS,
		HostPolicy: autocert.HostWhitelist(ls.Config.Address.SocketAddress.Domains...),
	}
	ls.srv = &http.Server{
		Addr:           ":https",
		Handler:        mux,
		ReadTimeout:    resolveStr2Time(hc.ReadTimeoutStr, 20*time.Second),
		WriteTimeout:   resolveStr2Time(hc.WriteTimeoutStr, 20*time.Second),
		IdleTimeout:    resolveStr2Time(hc.IdleTimeoutStr, 20*time.Second),
		MaxHeaderBytes: resolveInt2IntProp(hc.MaxHeaderBytes, 1<<20),
		TLSConfig:      m.TLSConfig(),
	}
	tcpListener, err := net.Listen("tcp", ls.srv.Addr)
	if err != nil {
		return errors.Wrapf(err, "bind HTTPS listener %s", ls.srv.Addr)
	}
	autoLs := tls.NewListener(tcpListener, m.TLSConfig())
	logger.Infof("[dubbo-go-server] httpsListener start at : %s", ls.srv.Addr)
	go func() {
		serveErr := ls.srv.Serve(autoLs)
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			logger.Errorf("[dubbo-go-server] httpsListener Serve error: %v", serveErr)
		}
	}()
	return nil
}

func (ls *HttpListenerService) httpListener() error {
	hl := createDefaultHttpWorker(ls)

	// user customize http config
	hc := model.MapInStruct(ls.Config.Config)

	mux := http.NewServeMux()
	mux.HandleFunc("/", hl.ServeHTTP)

	sa := ls.Config.Address.SocketAddress
	ls.srv = &http.Server{
		Addr:           resolveAddress(sa.Address + ":" + strconv.Itoa(sa.Port)),
		Handler:        mux,
		ReadTimeout:    resolveStr2Time(hc.ReadTimeoutStr, 20*time.Second),
		WriteTimeout:   resolveStr2Time(hc.WriteTimeoutStr, 20*time.Second),
		IdleTimeout:    resolveStr2Time(hc.IdleTimeoutStr, 20*time.Second),
		MaxHeaderBytes: resolveInt2IntProp(hc.MaxHeaderBytes, 1<<20),
	}

	logger.Infof("[dubbo-go-server] httpListener starting at %s with WriteTimeout: %v, IdleTimeout: %v, ReadTimeout: %v",
		ls.srv.Addr, ls.srv.WriteTimeout, ls.srv.IdleTimeout, ls.srv.ReadTimeout)

	netListener, err := net.Listen("tcp", ls.srv.Addr)
	if err != nil {
		return errors.Wrapf(err, "bind HTTP listener %s", ls.srv.Addr)
	}
	go func() {
		serveErr := ls.srv.Serve(netListener)
		if serveErr != nil && !errors.Is(serveErr, http.ErrServerClosed) {
			logger.Errorf("[dubbo-go-server] httpListener Serve error: %v", serveErr)
		} else {
			logger.Info("[dubbo-go-server] httpListener stopped gracefully.")
		}
	}()
	return nil
}

// createDefaultHttpWorker create http listener
func createDefaultHttpWorker(ls *HttpListenerService) *DefaultHttpWorker {
	return &DefaultHttpWorker{
		ls: ls,
	}
}

// ServeHTTP http request entrance.
func (s *DefaultHttpWorker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	if err := s.ls.WithFilterChain(func(fc *filterchain.NetworkFilterChain) error {
		fc.ServeHTTP(w, r)
		return nil
	}); err != nil {
		http.Error(w, err.Error(), http.StatusServiceUnavailable)
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
			logger.Errorf("Parse duration failed, err: %v", err)
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

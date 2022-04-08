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
	"fmt"
	"log"
	"net/http"
	"strconv"
	"time"
)

import (
	"github.com/pkg/errors"

	"golang.org/x/crypto/acme/autocert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/filterchain"
	"github.com/apache/dubbo-go-pixiu/pkg/listener"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
	"github.com/apache/dubbo-go-pixiu/pkg/tracing"
)

func init() {
	listener.SetListenerServiceFactory(model.ProtocolTypeHTTP, newHttpListenerService)
}

type (
	// ListenerService the facade of a listener
	HttpListenerService struct {
		listener.BaseListenerService
		srv *http.Server
	}

	// DefaultHttpListener
	DefaultHttpWorker struct {
		ls *HttpListenerService
	}
)

func newHttpListenerService(lc *model.Listener, bs *model.Bootstrap) (listener.ListenerService, error) {
	fc := filterchain.CreateNetworkFilterChain(lc.FilterChain, bs)
	return &HttpListenerService{
		BaseListenerService: listener.BaseListenerService{
			Config:      lc,
			FilterChain: fc,
		},
		srv: nil,
	}, nil
}

// Start start the listener
func (ls *HttpListenerService) Start() error {
	switch ls.Config.Protocol {
	case model.ProtocolTypeHTTP:
		ls.httpListener()
	case model.ProtocolTypeHTTPS:
		ls.httpsListener()
	default:
		return errors.New(fmt.Sprintf("unsupported protocol start: %d", ls.Config.Protocol))
	}
	return nil
}

func (ls *HttpListenerService) httpsListener() {
	hl := createDefaultHttpWorker(ls)

	// user customize http config
	var hc *model.HttpConfig
	hc = model.MapInStruct(ls.Config)

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
	autoLs := autocert.NewListener(ls.Config.Address.SocketAddress.Domains...)
	logger.Infof("[dubbo-go-server] httpsListener start at : %s", ls.srv.Addr)
	err := ls.srv.Serve(autoLs)
	logger.Info("[dubbo-go-server] httpsListener result:", err)
}

func (ls *HttpListenerService) httpListener() {
	hl := createDefaultHttpWorker(ls)

	// user customize http config
	var hc *model.HttpConfig
	hc = model.MapInStruct(ls.Config)

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

	logger.Infof("[dubbo-go-server] httpListener start at : %s", ls.srv.Addr)

	log.Println(ls.srv.ListenAndServe())
}

// createDefaultHttpWorker create http listener
func createDefaultHttpWorker(ls *HttpListenerService) *DefaultHttpWorker {
	return &DefaultHttpWorker{
		ls: ls,
	}
}

// ServeHTTP http request entrance.
func (s *DefaultHttpWorker) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	tracer := server.NewTracer(tracing.HTTPProtocol)
	defer tracer.Close()
	ctx, span := tracer.StartSpan("http_listener", r)
	defer span.End()
	r.Header.Set("tracing-id", tracer.GetID())

	s.ls.FilterChain.ServeHTTP(w, r.WithContext(ctx))
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

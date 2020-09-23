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

package proxy

import (
	"log"
	"net/http"
	"strconv"
	"sync"
	"time"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	h "github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
)

// ListenerService the facade of a listener
type ListenerService struct {
	*model.Listener
}

// Start start the listener
func (l *ListenerService) Start() {
	switch l.Address.SocketAddress.Protocol {
	case model.HTTP:
		l.httpListener()
	default:
		panic("un support protocol start: " + l.Address.SocketAddress.ProtocolStr)
	}
}

func (l *ListenerService) httpListener() {
	hl := NewDefaultHttpListener()
	hl.pool.New = func() interface{} {
		return l.allocateContext()
	}

	// user customize http config
	var hc model.HttpConfig
	if l.Config != nil {
		if c, ok := l.Config.(model.HttpConfig); ok {
			hc = c
		}
	}

	mux := http.NewServeMux()
	mux.HandleFunc("/", hl.ServeHTTP)

	srv := http.Server{
		Addr:           resolveAddress(l.Address.SocketAddress.Address + ":" + strconv.Itoa(l.Address.SocketAddress.Port)),
		Handler:        mux,
		ReadTimeout:    resolveStr2Time(hc.ReadTimeoutStr, 20*time.Second),
		WriteTimeout:   resolveStr2Time(hc.WriteTimeoutStr, 20*time.Second),
		IdleTimeout:    resolveStr2Time(hc.IdleTimeoutStr, 20*time.Second),
		MaxHeaderBytes: resolveInt2IntProp(hc.MaxHeaderBytes, 1<<20),
	}

	logger.Infof("[dubboproxy go] httpListener start by config : %+v", l)

	log.Println(srv.ListenAndServe())
}

func (l *ListenerService) allocateContext() *h.HttpContext {
	return &h.HttpContext{
		Listener:              l.Listener,
		FilterChains:          l.FilterChains,
		HttpConnectionManager: l.findHttpManager(),
		BaseContext:           context.NewBaseContext(),
	}
}

func (l *ListenerService) findHttpManager() model.HttpConnectionManager {
	for _, fc := range l.FilterChains {
		for _, f := range fc.Filters {
			if f.Name == constant.HttpConnectManagerFilter {
				return *f.Config.(*model.HttpConnectionManager)
			}
		}
	}

	return *DefaultHttpConnectionManager()
}

// DefaultHttpListener
type DefaultHttpListener struct {
	pool sync.Pool
}

// NewDefaultHttpListener create http listener
func NewDefaultHttpListener() *DefaultHttpListener {
	return &DefaultHttpListener{
		pool: sync.Pool{},
	}
}

// ServeHTTP
func (s *DefaultHttpListener) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	hc := s.pool.Get().(*h.HttpContext)
	hc.Request = r
	hc.ResetWritermen(w)
	hc.Reset()

	hc.AppendFilterFunc(
		extension.GetMustFilterFunc(constant.LoggerFilter),
		extension.GetMustFilterFunc(constant.RecoveryFilter),
	)

	hc.BuildFilters()

	s.handleHTTPRequest(hc)

	s.pool.Put(hc)
}

func (s *DefaultHttpListener) handleHTTPRequest(c *h.HttpContext) {
	if len(c.BaseContext.Filters) > 0 {
		c.Next()
		c.WriteHeaderNow()
		return
	}

	// TODO redirect
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

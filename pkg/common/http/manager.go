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
	"encoding/json"
	"fmt"
	"io"
	stdHttp "net/http"
	"sync"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/client/http"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	router2 "github.com/apache/dubbo-go-pixiu/pkg/common/router"
	"github.com/apache/dubbo-go-pixiu/pkg/common/util"
	pch "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// HttpConnectionManager network filter for http
type HttpConnectionManager struct {
	filter.EmptyNetworkFilter
	config            *model.HttpConnectionManagerConfig
	routerCoordinator *router2.RouterCoordinator
	filterManager     *filter.FilterManager
	pool              sync.Pool
}

// CreateHttpConnectionManager create http connection manager
func CreateHttpConnectionManager(hcmc *model.HttpConnectionManagerConfig) *HttpConnectionManager {
	hcm := &HttpConnectionManager{config: hcmc}
	hcm.pool.New = func() any {
		return hcm.allocateContext()
	}
	hcm.routerCoordinator = router2.CreateRouterCoordinator(&hcmc.RouteConfig)
	hcm.filterManager = filter.NewFilterManager(hcmc.HTTPFilters)
	hcm.filterManager.Load()
	return hcm
}

func (hcm *HttpConnectionManager) allocateContext() *pch.HttpContext {
	return &pch.HttpContext{
		Params: make(map[string]any),
	}
}

func (hcm *HttpConnectionManager) Handle(hc *pch.HttpContext) error {
	hc.Ctx = context.Background()
	err := hcm.findRoute(hc)
	if err != nil {
		return err
	}
	hcm.handleHTTPRequest(hc)
	return nil
}

func (hcm *HttpConnectionManager) ServeHTTP(w stdHttp.ResponseWriter, r *stdHttp.Request) {
	hc := hcm.pool.Get().(*pch.HttpContext)
	defer hcm.pool.Put(hc)

	hc.Writer = w
	hc.Request = r
	hc.Reset()
	hc.Timeout = hcm.config.Timeout
	err := hcm.Handle(hc)
	if err != nil {
		logger.Errorf("ServeHTTP %v", err)
	}
}

// handleHTTPRequest handle http request
func (hcm *HttpConnectionManager) handleHTTPRequest(c *pch.HttpContext) {
	filterChain := hcm.filterManager.CreateFilterChain(c)

	// recover any err when filterChain run
	defer func() {
		if err := recover(); err != nil {
			logger.Warnf("[dubbo-go-pixiu] Occur An Unexpected Err: %+v", err)
			c.SendLocalReply(stdHttp.StatusInternalServerError, []byte(fmt.Sprintf("Occur An Unexpected Err: %v", err)))
		}
	}()

	//todo timeout
	filterChain.OnDecode(c)
	hcm.buildTargetResponse(c)
	//todo: stream resp has to set HTTP Server's WriteTimeout to 0, need to check it
	filterChain.OnEncode(c)
	hcm.writeResponse(c)
}

func (hcm *HttpConnectionManager) writeResponse(c *pch.HttpContext) {
	if !c.LocalReply() {
		c.Writer.WriteHeader(c.GetStatusCode())
		if c.TargetResp != nil {
			switch res := c.TargetResp.(type) {
			case *client.UnaryResponse:
				_, err := c.Writer.Write(res.Data)
				if err != nil {
					logger.Errorf("Write response failed: %v", err)
				}
			case *client.StreamResponse:
				// create ctx helps goroutine exit
				ctx, cancel := context.WithCancel(c.Ctx)
				defer cancel()

				dataC := make(chan []byte)
				errC := make(chan error, 1)

				// goroutine read stream
				go func() {
					defer close(dataC)
					defer close(errC)
					buf := make([]byte, 1024) // 1KB buffer
					for {
						select {
						case <-ctx.Done():
							return
						default:
							n, err := res.Stream.Read(buf)
							if n > 0 {
								// copy data to prevent data cover
								data := make([]byte, n)
								copy(data, buf[:n])
								select {
								case dataC <- data:
								case <-ctx.Done():
									return
								}
							}
							if err != nil {
								if err != io.EOF {
									errC <- fmt.Errorf("stream read error: %w", err)
								} else {
									errC <- io.EOF
								}
								return
							}
						}
					}
				}()

				for {
					select {
					case <-ctx.Done():
						_ = res.Stream.Close()
						return
					case data, ok := <-dataC:
						if !ok {
							_ = res.Stream.Close()
							return
						}
						if _, err := c.Writer.Write(data); err != nil {
							cancel()
							_ = res.Stream.Close()
							return
						}

						if flusher, ok := c.Writer.(stdHttp.Flusher); ok {
							flusher.Flush()
						}

					case err := <-errC:
						if err != nil && err != io.EOF {
							logger.Errorf("Stream error: %v", err)
						}
						_ = res.Stream.Close()
						return
					}
				}
			default:
				logger.Errorf("Unknown response type: %T", c.TargetResp)
			}
		}
	}
}

func (hcm *HttpConnectionManager) buildTargetResponse(c *pch.HttpContext) {
	if c.LocalReply() {
		return
	}

	switch res := c.SourceResp.(type) {
	case *stdHttp.Response:
		//Merge header
		remoteHeader := res.Header
		for k := range remoteHeader {
			c.AddHeader(k, remoteHeader.Get(k))
		}
		//status code
		c.StatusCode(res.StatusCode)

		if http.IsStreamableResponse(res) {
			c.TargetResp = client.NewStreamResponse(res.Body, http.IsSSEStream(res))
		} else {
			body, err := io.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}
			//close body
			_ = res.Body.Close()
			c.TargetResp = client.NewUnaryResponse(body)
		}
	case []byte:
		c.StatusCode(stdHttp.StatusOK)
		if json.Valid(res) {
			c.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueApplicationJson)
		} else {
			c.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
		}
		c.TargetResp = client.NewUnaryResponse(res)
	default:
		//dubbo go generic invoke
		response := util.NewDubboResponse(res, false)
		c.StatusCode(stdHttp.StatusOK)
		c.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueJsonUtf8)
		c.TargetResp = response
	}
}

func (hcm *HttpConnectionManager) findRoute(hc *pch.HttpContext) error {
	ra, err := hcm.routerCoordinator.Route(hc)
	if err != nil {
		hc.SendLocalReply(stdHttp.StatusNotFound, constant.Default404Body)

		e := errors.Errorf("Requested URL %s not found", hc.GetUrl())
		logger.Debug(e.Error())
		return e
		// return 404
	}
	hc.RouteEntry(ra)
	return nil
}

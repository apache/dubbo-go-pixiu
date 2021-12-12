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
	"fmt"
	"io/ioutil"
	"net/http"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
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
	config            *model.HttpConnectionManagerConfig
	routerCoordinator *router2.RouterCoordinator
	filterManager     *filter.FilterManager
}

// CreateHttpConnectionManager create http connection manager
func CreateHttpConnectionManager(hcmc *model.HttpConnectionManagerConfig, bs *model.Bootstrap) *HttpConnectionManager {
	hcm := &HttpConnectionManager{config: hcmc}
	hcm.routerCoordinator = router2.CreateRouterCoordinator(hcmc)
	hcm.filterManager = filter.NewFilterManager(hcmc.HTTPFilters)
	hcm.filterManager.Load()
	return hcm
}

// OnData receive data from listener
func (hcm *HttpConnectionManager) OnData(hc *pch.HttpContext) error {
	hc.Ctx = context.Background()
	err := hcm.findRoute(hc)
	if err != nil {
		return err
	}
	hcm.handleHTTPRequest(hc)
	return nil
}

// handleHTTPRequest handle http request
func (hcm *HttpConnectionManager) handleHTTPRequest(c *pch.HttpContext) {
	filterChain := hcm.filterManager.CreateFilterChain(c)

	// recover any err when filterChain run
	defer func() {
		if err := recover(); err != nil {
			logger.Warnf("[dubbopixiu go] Occur An Unexpected Err: %+v", err)
			c.SendLocalReply(http.StatusInternalServerError, []byte(fmt.Sprintf("Occur An Unexpected Err: %v", err)))
		}
	}()

	//todo timeout
	filterChain.OnDecode(c)
	hcm.buildTargetResponse(c)
	filterChain.OnEncode(c)
	hcm.writeResponse(c)
}

func (hcm *HttpConnectionManager) writeResponse(c *pch.HttpContext) {
	if !c.LocalReply() {
		writer := c.Writer
		writer.WriteHeader(c.GetStatusCode())
		if _, err := writer.Write(c.TargetResp.Data); err != nil {
			panic(err)
		}
	}
}

func (hcm *HttpConnectionManager) buildTargetResponse(c *pch.HttpContext) {
	if !c.LocalReply() {
		switch res := c.SourceResp.(type) {
		case *http.Response:
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				panic(err)
			}
			//Merge header
			remoteHeader := res.Header
			for k := range remoteHeader {
				c.AddHeader(k, remoteHeader.Get(k))
			}
			//status code
			c.StatusCode(res.StatusCode)
			c.TargetResp = &client.Response{Data: body}
		case []byte:
			c.StatusCode(http.StatusOK)
			c.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
			c.TargetResp = &client.Response{Data: res}
		default:
			//dubbo go generic invoke
			response := util.NewDubboResponse(res, false)
			c.StatusCode(http.StatusOK)
			c.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueJsonUtf8)
			c.TargetResp = response
		}
	}
}

func (hcm *HttpConnectionManager) findRoute(hc *pch.HttpContext) error {
	ra, err := hcm.routerCoordinator.Route(hc)
	if err != nil {
		hc.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
		hc.SendLocalReply(http.StatusNotFound, constant.Default404Body)

		e := errors.Errorf("Requested URL %s not found", hc.GetUrl())
		logger.Debug(e.Error())
		return e
		// return 404
	}
	hc.RouteEntry(ra)
	return nil
}

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
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	router2 "github.com/apache/dubbo-go-pixiu/pkg/common/router"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

import (
	pch "github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

type HttpConnectionManager struct {
	config            *model.HttpConnectionManager
	routerCoordinator *router2.RouterCoordinator
	clusterManager    *server.ClusterManager
	filterManager     *filter.FilterManager
}

func CreateHttpConnectionManager(hcmc *model.HttpConnectionManager, bs *model.Bootstrap) *HttpConnectionManager {
	hcm := &HttpConnectionManager{config: hcmc}
	hcm.routerCoordinator = router2.CreateRouterCoordinator(hcmc)
	hcm.filterManager = filter.NewFilterManager(hcmc.HTTPFilters)
	hcm.filterManager.Load()
	return hcm
}

func (hcm *HttpConnectionManager) OnData(hc *pch.HttpContext) error {
	hc.Ctx = context.Background()
	hcm.findRoute(hc)
	hcm.addFilter(hc)
	hcm.handleHTTPRequest(hc)
	return nil
}

func (hcm *HttpConnectionManager) findRoute(hc *pch.HttpContext) {
	ra, err := hcm.routerCoordinator.Route(hc)
	if err != nil {
		// return 404
	}
	hc.RouteEntry(ra)
}

func (hcm *HttpConnectionManager) handleHTTPRequest(c *pch.HttpContext) {
	if len(c.Filters) > 0 {
		c.Next()
		c.WriteHeaderNow()
		return
	}

	// TODO redirect
}

func (hcm *HttpConnectionManager) addFilter(ctx *pch.HttpContext) {
	for _, f := range hcm.filterManager.GetFilters() {
		f.PrepareFilterChain(ctx)
	}
}

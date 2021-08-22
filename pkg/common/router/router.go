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

package router

import (
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

type (
	RouterCoordinator struct {
		activeConfig *model.RouteConfiguration
	}
)

func CreateRouterCoordinator(hcmc *model.HttpConnectionManager) *RouterCoordinator {

	rc := &RouterCoordinator{activeConfig: &hcmc.RouteConfig}
	if hcmc.RouteConfig.Dynamic {
		server.GetRouterManager().AddRouterListener(rc)
	}
	return rc
}

func (rm *RouterCoordinator) Route(hc *http.HttpContext) (*model.RouteAction, error) {
	return rm.activeConfig.Route(hc.Request)
}

func (rm *RouterCoordinator) OnAddRouter(r *model.Router) {

}

func (rm *RouterCoordinator) OnDeleteRouter(r *model.Router) {

}

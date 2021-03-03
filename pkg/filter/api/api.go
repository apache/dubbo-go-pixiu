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

package api

import (
	"net/http"
)

import (
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/context"
)

import (
	"github.com/dubbogo/dubbo-go-pixiu/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-pixiu/pkg/common/extension"
)

// nolint
func Init() {
	extension.SetFilterFunc(constant.HTTPApiFilter, apiFilterFunc())
}

// apiFilterFunc url match api
func apiFilterFunc() context.FilterFunc {
	return func(c context.Context) {
		url := c.GetUrl()
		method := c.GetMethod()
		// [williamfeng323]TO-DO: get the API details from router which saved in constant.LocalMemoryApiDiscoveryService
		if api, ok := api.EmptyApi.FindApi(url); ok {
			if !api.MatchMethod(method) {
				c.WriteWithStatus(http.StatusMethodNotAllowed, constant.Default405Body)
				c.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
				c.Abort()
				return
			}

			if !api.IsOk(api.Name) {
				c.WriteWithStatus(http.StatusNotAcceptable, constant.Default406Body)
				c.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
				c.Abort()
				return
			}
			// [williamfeng323]TO-DO: the c.Api method need to be updated to use the newest API definition
			c.Api(api)
			c.Next()
		} else {
			c.WriteWithStatus(http.StatusNotFound, constant.Default404Body)
			c.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
			c.Abort()
		}
	}
}

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
	"net/http/httptest"
)

import (
	ctxHttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func getTestContext() *ctxHttp.HttpContext {
	l := ListenerService{
		cfg: &model.Listener{
			Name: "test",
			Address: model.Address{
				SocketAddress: model.SocketAddress{
					Protocol: model.ProtocolTypeHTTP,
					Address:  "0.0.0.0",
					Port:     8888,
				},
			},
			FilterChains: []model.FilterChain{},
		},
	}

	hc := &ctxHttp.HttpContext{
		Listener: l.cfg,
		Filters:  []ctxHttp.FilterFunc{},
	}
	hc.ResetWritermen(httptest.NewRecorder())
	hc.Reset()
	return hc
}

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
	"bytes"
	"net/http"
	"net/http/httptest"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/mock"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	ctxHttp "github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
	"github.com/dubbogo/dubbo-go-proxy/pkg/router"
)

func getTestContext() *ctxHttp.HttpContext {
	l := ListenerService{
		Listener: &model.Listener{
			Name: "test",
			Address: model.Address{
				SocketAddress: model.SocketAddress{
					Protocol: model.HTTP,
					Address:  "0.0.0.0",
					Port:     8888,
				},
			},
			FilterChains: []model.FilterChain{},
		},
	}

	hc := &ctxHttp.HttpContext{
		Listener:              l.Listener,
		FilterChains:          l.FilterChains,
		HttpConnectionManager: l.findHttpManager(),
		BaseContext:           context.NewBaseContext(),
	}
	hc.ResetWritermen(httptest.NewRecorder())
	hc.Reset()
	return hc
}

func TestRouteRequest(t *testing.T) {
	mockAPI := mock.GetMockAPI(config.MethodPost, "/mock/test")
	mockAPI.Method.OnAir = false

	apiDiscoverySrv := extension.GetMustAPIDiscoveryService(constant.LocalMemoryApiDiscoveryService)
	apiDiscoverySrv.AddAPI(mockAPI)
	apiDiscoverySrv.AddAPI(mock.GetMockAPI(config.MethodGet, "/mock/test"))

	listener := NewDefaultHttpListener()
	listener.pool.New = func() interface{} {
		return getTestContext()
	}
	r := bytes.NewReader([]byte("test"))

	req, _ := http.NewRequest("GET", "/mock/test", r)
	ctx := listener.pool.Get().(*ctxHttp.HttpContext)
	api, err := listener.routeRequest(ctx, req)
	assert.Nil(t, err)
	assert.NotNil(t, api)
	assert.Equal(t, api.URLPattern, "/mock/test")
	assert.Equal(t, api.HTTPVerb, config.MethodGet)

	req, _ = http.NewRequest("GET", "/mock/test2", r)
	ctx = listener.pool.Get().(*ctxHttp.HttpContext)
	api, err = listener.routeRequest(ctx, req)
	assert.EqualError(t, err, "Requested URL /mock/test2 not found")
	assert.Equal(t, ctx.StatusCode(), 404)
	assert.Equal(t, api, router.API{})

	req, _ = http.NewRequest("POST", "/mock/test", r)
	ctx = listener.pool.Get().(*ctxHttp.HttpContext)
	api, err = listener.routeRequest(ctx, req)
	assert.EqualError(t, err, "Requested API POST /mock/test does not online")
	assert.Equal(t, ctx.StatusCode(), 406)
	assert.Equal(t, api, router.API{})
}

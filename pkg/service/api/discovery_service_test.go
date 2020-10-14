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
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"github.com/dubbogo/dubbo-go-proxy/pkg/router"
)

func getMockAPI(urlPattern string, verb config.HTTPVerb) router.API {
	inbound := config.InboundRequest{}
	integration := config.IntegrationRequest{}
	method := config.Method{
		OnAir:              true,
		HTTPVerb:           verb,
		InboundRequest:     inbound,
		IntegrationRequest: integration,
	}
	return router.API{
		URLPattern: urlPattern,
		Method:     method,
	}
}

func TestNewLocalMemoryAPIDiscoveryService(t *testing.T) {
	l := NewLocalMemoryAPIDiscoveryService()
	assert.NotNil(t, l)
	assert.NotNil(t, l.router)
}

func TestAddAPI(t *testing.T) {
	l := NewLocalMemoryAPIDiscoveryService()
	err := l.AddAPI(getMockAPI("/this/is/test", config.MethodPut))
	assert.Nil(t, err)
	_, found := l.router.FindAPI("/this/is/test", config.MethodPut)
	assert.True(t, found)
}

func TestGetAPI(t *testing.T) {
	l := NewLocalMemoryAPIDiscoveryService()
	err := l.AddAPI(getMockAPI("/this/is/test", config.MethodPut))
	assert.Nil(t, err)
	_, err = l.GetAPI("/this/is/test", config.MethodPut)
	assert.Nil(t, err)

	_, err = l.GetAPI("/this/is/test/or/else", config.MethodPut)
	assert.NotNil(t, err)
}

func TestLoadAPI(t *testing.T) {
	apiC, err := config.LoadAPIConfigFromFile("../../config/mock/api_config.yml")
	assert.Empty(t, err)
	err = InitAPIsFromConfig(*apiC)
	assert.Nil(t, err)
	apiDisSrv := extension.GetMustAPIDiscoveryService(constant.LocalMemoryApiDiscoveryService)
	rsp, err := apiDisSrv.GetAPI("/", config.MethodGet)
	assert.NotNil(t, rsp.URLPattern)
	rsp, err = apiDisSrv.GetAPI("/mockTest", config.MethodGet)
	assert.NotNil(t, rsp.URLPattern)
	rsp, err = apiDisSrv.GetAPI("/mockTest", config.MethodPost)
	assert.NotNil(t, rsp.URLPattern)
	rsp, err = apiDisSrv.GetAPI("/mockTest/12345", config.MethodGet)
	assert.NotNil(t, rsp.URLPattern)
}

func TestLoadAPIFromResource(t *testing.T) {
	apiDiscSrv := NewLocalMemoryAPIDiscoveryService()
	mockMethod1 := getMockAPI("", config.MethodPut).Method
	mockMethod2 := getMockAPI("", config.MethodPost).Method
	mockMethod3 := getMockAPI("", config.MethodGet).Method
	tempResources := []config.Resource{
		{
			Type:        "Restful",
			Path:        "/",
			Description: "test only",
			Methods: []config.Method{
				mockMethod1,
				mockMethod2,
				mockMethod3,
			},
			Resources: []config.Resource{
				{
					Type:        "Restful",
					Path:        "/mock",
					Description: "test only",
					Methods: []config.Method{
						mockMethod1,
						mockMethod2,
						mockMethod3,
					},
				},
				{
					Type:        "Restful",
					Path:        "/mock2",
					Description: "test only",
					Methods: []config.Method{
						mockMethod1,
					},
					Resources: []config.Resource{
						{
							Type:        "Restful",
							Path:        "/:id",
							Description: "test only",
							Methods: []config.Method{
								mockMethod1,
							},
						},
					},
				},
			},
		},
	}
	err := loadAPIFromResource("", tempResources, apiDiscSrv)
	assert.Nil(t, err)
	rsp, _ := apiDiscSrv.GetAPI("/", config.MethodPut)
	assert.Equal(t, rsp.URLPattern, "/")
	rsp, _ = apiDiscSrv.GetAPI("/", config.MethodGet)
	assert.Equal(t, rsp.URLPattern, "/")
	rsp, _ = apiDiscSrv.GetAPI("/mock", config.MethodGet)
	assert.Equal(t, rsp.URLPattern, "/mock")
	rsp, _ = apiDiscSrv.GetAPI("/mock2/12345", config.MethodPut)
	assert.Equal(t, rsp.URLPattern, "/mock2/:id")

	tempResources = []config.Resource{
		{
			Type:        "Restful",
			Path:        "/mock",
			Description: "test only",
			Methods: []config.Method{
				mockMethod1,
			},
			Resources: []config.Resource{
				{
					Type:        "Restful",
					Path:        ":id",
					Description: "test only",
					Methods: []config.Method{
						mockMethod1,
					},
				},
				{
					Type:        "Restful",
					Path:        ":ik",
					Description: "test only",
					Methods: []config.Method{
						mockMethod1,
					},
				},
			},
		},
	}
	apiDiscSrv = NewLocalMemoryAPIDiscoveryService()
	err = loadAPIFromResource("", tempResources, apiDiscSrv)
	assert.EqualError(t, err, "Path :id in /mock doesn't start with /; Path :ik in /mock doesn't start with /")
}

func TestLoadAPIFromMethods(t *testing.T) {
	tempMethods := []config.Method{
		getMockAPI("", config.MethodPut).Method,
		getMockAPI("", config.MethodGet).Method,
		getMockAPI("", config.MethodPut).Method,
	}
	apiDiscSrv := NewLocalMemoryAPIDiscoveryService()
	err := loadAPIFromMethods("/mock", tempMethods, apiDiscSrv)
	rsp, _ := apiDiscSrv.GetAPI("/mock", config.MethodPut)
	assert.Equal(t, rsp.URLPattern, "/mock")
	rsp, _ = apiDiscSrv.GetAPI("/mock", config.MethodGet)
	assert.Equal(t, rsp.URLPattern, "/mock")
	assert.EqualError(t, err, "Path: /mock, Method: PUT, error: Method PUT already exists in path /mock")
}

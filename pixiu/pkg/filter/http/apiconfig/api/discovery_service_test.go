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
	"strings"
	"testing"
)

import (
	fc "github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/mock"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/config"
)

func TestNewLocalMemoryAPIDiscoveryService(t *testing.T) {
	l := NewLocalMemoryAPIDiscoveryService()
	assert.NotNil(t, l)
	assert.NotNil(t, l.router)
}

func TestAddAPI(t *testing.T) {
	l := NewLocalMemoryAPIDiscoveryService()
	err := l.AddAPI(mock.GetMockAPI(fc.MethodPut, "/this/is/test"))
	assert.Nil(t, err)
	_, found := l.router.FindAPI("/this/is/test", fc.MethodPut)
	assert.True(t, found)
}

func TestGetAPI(t *testing.T) {
	l := NewLocalMemoryAPIDiscoveryService()
	err := l.AddAPI(mock.GetMockAPI(fc.MethodPut, "/this/is/test"))
	assert.Nil(t, err)
	_, err = l.GetAPI("/this/is/test", fc.MethodPut)
	assert.Nil(t, err)

	_, err = l.GetAPI("/this/is/test/or/else", fc.MethodPut)
	assert.NotNil(t, err)
}

func TestLoadAPI(t *testing.T) {
	apiDisSrv := NewLocalMemoryAPIDiscoveryService()
	apiC, err := config.LoadAPIConfigFromFile("../../../../config/mock/api_config.yml")
	assert.Empty(t, err)
	err = apiDisSrv.InitAPIsFromConfig(*apiC)
	assert.Nil(t, err)
	rsp, _ := apiDisSrv.GetAPI("/", fc.MethodGet)
	assert.NotNil(t, rsp.URLPattern)
	rsp, _ = apiDisSrv.GetAPI("/mockTest", fc.MethodGet)
	assert.NotNil(t, rsp.URLPattern)
	rsp, _ = apiDisSrv.GetAPI("/mockTest", fc.MethodPost)
	assert.NotNil(t, rsp.URLPattern)
	rsp, _ = apiDisSrv.GetAPI("/mockTest/12345", fc.MethodGet)
	assert.NotNil(t, rsp.URLPattern)
}

func TestLoadAPIFromResource(t *testing.T) {
	apiDiscSrv := NewLocalMemoryAPIDiscoveryService()
	mockMethod1 := mock.GetMockAPI(fc.MethodPut, "").Method
	mockMethod2 := mock.GetMockAPI(fc.MethodPost, "").Method
	mockMethod3 := mock.GetMockAPI(fc.MethodGet, "").Method
	tempResources := []fc.Resource{
		{
			Type:        "Restful",
			Path:        "/",
			Description: "test only",
			Methods: []fc.Method{
				mockMethod1,
				mockMethod2,
				mockMethod3,
			},
			Resources: []fc.Resource{
				{
					Type:        "Restful",
					Path:        "/mock",
					Description: "test only",
					Methods: []fc.Method{
						mockMethod1,
						mockMethod2,
						mockMethod3,
					},
				},
				{
					Type:        "Restful",
					Path:        "/mock2",
					Description: "test only",
					Methods: []fc.Method{
						mockMethod1,
					},
					Resources: []fc.Resource{
						{
							Type:        "Restful",
							Path:        "/:id",
							Description: "test only",
							Methods: []fc.Method{
								mockMethod1,
							},
						},
					},
				},
			},
		},
	}
	err := loadAPIFromResource("", tempResources, nil, apiDiscSrv)
	assert.Nil(t, err)
	rsp, _ := apiDiscSrv.GetAPI("/", fc.MethodPut)
	assert.Equal(t, rsp.URLPattern, "/")
	rsp, _ = apiDiscSrv.GetAPI("/", fc.MethodGet)
	assert.Equal(t, rsp.URLPattern, "/")
	rsp, _ = apiDiscSrv.GetAPI("/mock", fc.MethodGet)
	assert.Equal(t, rsp.URLPattern, "/mock")
	rsp, _ = apiDiscSrv.MatchAPI("/mock2/12345", fc.MethodPut)
	assert.Equal(t, rsp.URLPattern, "/mock2/:id")

	tempResources = []fc.Resource{
		{
			Type:        "Restful",
			Path:        "/mock",
			Description: "test only",
			Methods: []fc.Method{
				mockMethod1,
			},
			Resources: []fc.Resource{
				{
					Type:        "Restful",
					Path:        ":id",
					Description: "test only",
					Methods: []fc.Method{
						mockMethod1,
					},
				},
				{
					Type:        "Restful",
					Path:        ":ik",
					Description: "test only",
					Methods: []fc.Method{
						mockMethod1,
					},
				},
			},
		},
	}
	apiDiscSrv = NewLocalMemoryAPIDiscoveryService()
	err = loadAPIFromResource("", tempResources, nil, apiDiscSrv)
	assert.EqualError(t, err, "path :id in /mock doesn't start with /; path :ik in /mock doesn't start with /")
}

func TestLoadAPIFromMethods(t *testing.T) {
	mockPutAPIMethod := mock.GetMockAPI(fc.MethodPut, "").Method
	mockPutAPIMethod2 := mock.GetMockAPI(fc.MethodPut, "").Method
	mockPutAPIMethod.IntegrationRequest.URL = "localhost:8080"
	mockPutAPIMethod2.IntegrationRequest.URL = "localhost:8081"
	tempMethods := []fc.Method{
		mockPutAPIMethod,
		mockPutAPIMethod2,
		mock.GetMockAPI(fc.MethodGet, "").Method,
		mockPutAPIMethod,
	}
	apiDiscSrv := NewLocalMemoryAPIDiscoveryService()
	err := loadAPIFromMethods("/mock", tempMethods, nil, apiDiscSrv)
	rsp, _ := apiDiscSrv.GetAPI("/mock", fc.MethodPut)
	assert.Equal(t, rsp.URLPattern, "/mock")
	rsp, _ = apiDiscSrv.GetAPI("/mock", fc.MethodGet)
	assert.Equal(t, rsp.URLPattern, "/mock")
	assert.True(t, strings.Contains(err.Error(), "path: /mock, Method: PUT, error: Method PUT with address /mock already exists in path /mock"))
}

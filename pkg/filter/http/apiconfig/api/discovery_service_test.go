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
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/mock"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
)

func TestNewLocalMemoryAPIDiscoveryService(t *testing.T) {
	l := NewLocalMemoryAPIDiscoveryService()
	assert.NotNil(t, l)
	assert.NotNil(t, l.router)
}

func TestAddAPI(t *testing.T) {
	l := NewLocalMemoryAPIDiscoveryService()
	err := l.AddAPI(mock.GetMockAPI(constant.Put, "/this/is/test"))
	assert.Nil(t, err)
	_, found := l.router.FindAPI("/this/is/test", constant.Put)
	assert.True(t, found)
}

func TestGetAPI(t *testing.T) {
	l := NewLocalMemoryAPIDiscoveryService()
	err := l.AddAPI(mock.GetMockAPI(constant.Put, "/this/is/test"))
	assert.Nil(t, err)
	_, err = l.GetAPI("/this/is/test", constant.Put)
	assert.Nil(t, err)

	_, err = l.GetAPI("/this/is/test/or/else", constant.Put)
	assert.NotNil(t, err)
}

func TestLoadAPI(t *testing.T) {
	apiDisSrv := NewLocalMemoryAPIDiscoveryService()
	apiC, err := config.LoadAPIConfigFromFile("../../../../config/mock/api_config.yml")
	assert.Empty(t, err)
	err = apiDisSrv.InitAPIsFromConfig(*apiC)
	assert.Nil(t, err)
	rsp, _ := apiDisSrv.GetAPI("/", constant.Get)
	assert.NotNil(t, rsp.URLPattern)
	rsp, _ = apiDisSrv.GetAPI("/mockTest", constant.Get)
	assert.NotNil(t, rsp.URLPattern)
	rsp, _ = apiDisSrv.GetAPI("/mockTest", constant.Post)
	assert.NotNil(t, rsp.URLPattern)
	rsp, _ = apiDisSrv.GetAPI("/mockTest/12345", constant.Get)
	assert.NotNil(t, rsp.URLPattern)
}

func TestLoadAPIFromResource(t *testing.T) {
	apiDiscSrv := NewLocalMemoryAPIDiscoveryService()
	mockMethod1 := mock.GetMockAPI(constant.Put, "").Method
	mockMethod2 := mock.GetMockAPI(constant.Post, "").Method
	mockMethod3 := mock.GetMockAPI(constant.Get, "").Method
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
	err := loadAPIFromResource("", tempResources, nil, apiDiscSrv)
	assert.Nil(t, err)
	rsp, _ := apiDiscSrv.GetAPI("/", constant.Put)
	assert.Equal(t, rsp.URLPattern, "/")
	rsp, _ = apiDiscSrv.GetAPI("/", constant.Get)
	assert.Equal(t, rsp.URLPattern, "/")
	rsp, _ = apiDiscSrv.GetAPI("/mock", constant.Get)
	assert.Equal(t, rsp.URLPattern, "/mock")
	rsp, _ = apiDiscSrv.MatchAPI("/mock2/12345", constant.Put)
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
	err = loadAPIFromResource("", tempResources, nil, apiDiscSrv)
	assert.EqualError(t, err, "path :id in /mock doesn't start with /; path :ik in /mock doesn't start with /")
}

func TestLoadAPIFromMethods(t *testing.T) {
	mockPutAPIMethod := mock.GetMockAPI(constant.Put, "").Method
	mockPutAPIMethod2 := mock.GetMockAPI(constant.Put, "").Method
	mockPutAPIMethod.URL = "localhost:8080"
	mockPutAPIMethod2.URL = "localhost:8081"
	tempMethods := []config.Method{
		mockPutAPIMethod,
		mockPutAPIMethod2,
		mock.GetMockAPI(constant.Get, "").Method,
		mockPutAPIMethod,
	}
	apiDiscSrv := NewLocalMemoryAPIDiscoveryService()
	err := loadAPIFromMethods("/mock", tempMethods, nil, apiDiscSrv)
	rsp, _ := apiDiscSrv.GetAPI("/mock", constant.Put)
	assert.Equal(t, rsp.URLPattern, "/mock")
	rsp, _ = apiDiscSrv.GetAPI("/mock", constant.Get)
	assert.Equal(t, rsp.URLPattern, "/mock")
	assert.True(t, strings.Contains(err.Error(), "path: /mock, Method: PUT, error: Method PUT with address /mock already exists in path /mock"))
}

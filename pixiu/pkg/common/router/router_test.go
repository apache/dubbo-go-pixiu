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
	"bytes"
	"net/http"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/router/trie"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/mock"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

func TestCreateRouterCoordinator(t *testing.T) {
	hcmc := model.HttpConnectionManagerConfig{
		RouteConfig: model.RouteConfiguration{
			RouteTrie: trie.NewTrieWithDefault("POST/api/v1/**", model.RouteAction{
				Cluster:                     "test_dubbo",
				ClusterNotFoundResponseCode: 505,
			}),
			Dynamic: false,
		},
		HTTPFilters: []*model.HTTPFilter{
			{
				Name:   "test",
				Config: nil,
			},
		},
		ServerName:        "test_http_dubbo",
		GenerateRequestID: false,
		IdleTimeoutStr:    "100",
	}

	r := CreateRouterCoordinator(&hcmc.RouteConfig)
	request, err := http.NewRequest("POST", "http://www.dubbogopixiu.com/api/v1?name=tc", bytes.NewReader([]byte("{\"id\":\"12345\"}")))
	assert.NoError(t, err)
	c := mock.GetMockHTTPContext(request)
	a, err := r.Route(c)
	assert.NoError(t, err)
	assert.Equal(t, a.Cluster, "test_dubbo")

	router := &model.Router{
		ID: "1",
		Match: model.RouterMatch{
			Prefix: "/user",
		},
		Route: model.RouteAction{
			Cluster: "test",
		},
	}

	r.OnAddRouter(router)
	r.OnDeleteRouter(router)
}

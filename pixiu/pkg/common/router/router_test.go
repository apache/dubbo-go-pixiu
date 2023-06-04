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
			}, nil),
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

func TestRoute(t *testing.T) {
	const (
		Cluster1 = "test-cluster-1"
		Cluster2 = "test-cluster-2"
		Cluster3 = "test-cluster-3"
	)

	hcmc := model.HttpConnectionManagerConfig{
		RouteConfig: model.RouteConfiguration{
			Routes: []*model.Router{
				{
					ID: "1",
					Match: model.RouterMatch{
						Prefix: "/user",
						Headers: []model.HeaderMatcher{
							{
								Name:   "A",
								Values: []string{"1", "2", "3"},
							},
							{
								Name:   "A",
								Values: []string{"3", "4", "5"},
							},
							{
								Name:   "B",
								Values: []string{"1"},
							},
							{
								Name:   "normal-regex",
								Values: []string{"(k){2}"},
								Regex:  true,
							},
							{
								Name:   "broken-regex",
								Values: []string{"(t){2]]"},
								Regex:  true,
							},
						},
						Methods: []string{"GET", "POST"},
					},
					Route: model.RouteAction{
						Cluster:                     Cluster1,
						ClusterNotFoundResponseCode: 505,
					},
				},
				{
					ID: "2",
					Match: model.RouterMatch{
						Headers: []model.HeaderMatcher{
							{
								Name:   "B",
								Values: []string{"1", "9", "8"},
							},
							{
								Name:   "C",
								Values: []string{"7"},
							},
							{
								Name:   "normal-regex",
								Values: []string{"(t){2}"},
								Regex:  true,
							},
							{
								Name:   "forget-regex",
								Values: []string{"(f){4}"},
							},
						},
						Methods: []string{"GET", "PUT"},
					},
					Route: model.RouteAction{
						Cluster:                     Cluster2,
						ClusterNotFoundResponseCode: 505,
					},
				},
				{
					ID: "3",
					Match: model.RouterMatch{
						Headers: []model.HeaderMatcher{
							{
								Name:   "D",
								Values: []string{"7"},
							},
							{
								Name:   "normal-regex",
								Values: []string{"(n){2}"},
								Regex:  true,
							},
						},
						Methods: []string{"GET"},
					},
				},
				{
					ID: "4",
					Match: model.RouterMatch{
						Prefix: "/test",
					},
					Route: model.RouteAction{
						Cluster:                     Cluster3,
						ClusterNotFoundResponseCode: 505,
					},
				},
			},
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

	testCases := []struct {
		Name   string
		URL    string
		Method string
		Header map[string]string
		Expect string
	}{
		{
			Name:   "aim@ID3, no header but prefix match",
			URL:    "/test",
			Expect: Cluster3,
		},
		{
			Name: "aim@ID1, one override header",
			URL:  "/user",
			Header: map[string]string{
				"A": "1",
			},
			Expect: "prefix matched, but no headers matched.",
		},
		{
			Name: "aim@ID1, one header matched",
			URL:  "/user",
			Header: map[string]string{
				"A": "3",
			},
			Expect: Cluster1,
		},
		{
			Name: "aim@ID1, more header with one regex matched",
			URL:  "/user",
			Header: map[string]string{
				"A":            "5",
				"normal-regex": "kkkk",
			},
			Expect: Cluster1,
		},
		{
			Name:   "aim@ID1, one header but wrong method",
			URL:    "/user",
			Method: "PUT",
			Header: map[string]string{
				"A": "3",
			},
			Expect: "route failed for PUT/user, no rules matched.",
		},
		{
			Name: "aim@ID1, one header but no prefix match",
			URL:  "/error",
			Header: map[string]string{
				"A": "4",
			},
			Expect: "route failed for GET/error, no rules matched.",
		},
		{
			Name: "aim@ID1, one regex header mismatch",
			URL:  "/user",
			Header: map[string]string{
				"normal-regex": "tt",
			},
			Expect: Cluster2,
		},
		{
			Name: "aim@ID1, one broken regex header",
			URL:  "/user",
			Header: map[string]string{
				"broken-regex": "tt",
			},
			Expect: "prefix matched, but no headers matched.",
		},
		{
			Name: "aim@ID2, one matched header 1",
			Header: map[string]string{
				"B": "9",
			},
			Expect: Cluster2,
		},
		{
			Name: "aim@ID2, one matched header 2",
			Header: map[string]string{
				"B": "1",
			},
			Expect: Cluster2,
		},
		{
			Name:   "aim@ID2, only header but wrong method",
			Method: "DELETE",
			Header: map[string]string{
				"B": "1",
			},
			Expect: "route failed for DELETE, no rules matched.",
		},
		{
			Name: "aim@ID2, one header but mismatch",
			Header: map[string]string{
				"C": "4",
			},
			Expect: "prefix matched, but no headers matched.",
		},
		{
			Name: "aim@ID2, more matched header",
			Header: map[string]string{
				"C": "7",
				"B": "8",
			},
			Expect: Cluster2,
		},
		{
			Name: "aim@ID2, more header but mismatch",
			Header: map[string]string{
				"C": "4",
				"B": "5",
			},
			Expect: "prefix matched, but no headers matched.",
		},
		{
			Name: "aim@ID2, one matched regex header",
			Header: map[string]string{
				"normal-regex": "tt",
			},
			Expect: Cluster2,
		},
		{
			Name: "aim@ID2, but with wrong regex value",
			Header: map[string]string{
				"normal-regex": "gg",
			},
			Expect: "prefix matched, but no headers matched.",
		},
		{
			Name: "aim@ID2, intend regex match but forget to enable",
			Header: map[string]string{
				"forget-regex": "ffff",
			},
			Expect: "prefix matched, but no headers matched.",
		},
		{
			Name: "aim@ID3, but got no route action",
			Header: map[string]string{
				"D": "7",
			},
			Expect: "action is nil. please check your configuration.",
		},
		{
			Name: "aim@ID3, regex match but got no route action",
			Header: map[string]string{
				"normal-regex": "nn",
			},
			Expect: "action is nil. please check your configuration.",
		},
		{
			Name:   "aim@ID3, but wrong method",
			Method: "PUT",
			Header: map[string]string{
				"normal-regex": "nn",
			},
			Expect: "prefix matched, but no headers matched.",
		},
	}

	r := CreateRouterCoordinator(&hcmc.RouteConfig)

	for _, tc := range testCases {
		t.Run(tc.Name, func(t *testing.T) {
			method := "GET"
			if len(tc.Method) > 0 {
				method = tc.Method
			}
			request, err := http.NewRequest(method, tc.URL, nil)
			assert.NoError(t, err)

			if tc.Header != nil {
				for k, v := range tc.Header {
					request.Header.Set(k, v)
				}
			}
			c := mock.GetMockHTTPContext(request)

			a, err := r.Route(c)
			if err != nil {
				assert.Equal(t, tc.Expect, err.Error())
			} else {
				assert.Equal(t, tc.Expect, a.Cluster)
			}
		})
	}
}

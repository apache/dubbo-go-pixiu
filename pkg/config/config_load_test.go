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

package config

import (
	"encoding/json"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/dubbogo/dubbo-go-pixiu/pkg/model"
)

func TestLoad(t *testing.T) {
	conf := Load("conf_test.yaml")

	assert.Equal(t, 1, len(conf.StaticResources.Listeners))
	assert.Equal(t, 1, len(conf.StaticResources.Clusters))
}

func TestStruct2JSON(t *testing.T) {
	b := model.Bootstrap{
		StaticResources: model.StaticResources{
			Listeners: []model.Listener{
				{
					Name: "net/http",
					Address: model.Address{
						SocketAddress: model.SocketAddress{
							ProtocolStr: "HTTP",
							Address:     "127.0.0.0",
							Port:        8899,
						},
					},
					Config: model.HttpConfig{
						IdleTimeoutStr:  "5s",
						WriteTimeoutStr: "5s",
						ReadTimeoutStr:  "5s",
					},
					FilterChains: []model.FilterChain{
						{
							FilterChainMatch: model.FilterChainMatch{
								Domains: []string{
									"api.dubbo.com",
									"api.pixiu.com",
								},
							},
							Filters: []model.Filter{
								{
									Name: "dgp.filters.http_connect_manager",
									Config: model.HttpConnectionManager{
										RouteConfig: model.RouteConfiguration{
											Routes: []model.Router{
												{
													Match: model.RouterMatch{
														Prefix: "/api/v1",
													},
													Route: model.RouteAction{
														Cluster: "test_dubbo",
														Cors: model.CorsPolicy{
															AllowOrigin: []string{
																"*",
															},
														},
													},
												},
											},
										},
										HTTPFilters: []model.HTTPFilter{
											{
												Name: "dgp.filters.http.cors",
											},
											{
												Name: "dgp.filters.http.router",
											},
										},
									},
								},
							},
						},
					},
				},
			},
			Clusters: []model.Cluster{
				{
					Name:              "test_dubbo",
					TypeStr:           "EDS",
					LbStr:             "RoundRobin",
					ConnectTimeoutStr: "5s",
				},
			},
		},
	}

	if bytes, err := json.Marshal(b); err != nil {
		t.Fatal(err)
	} else {
		t.Log(string(bytes))
	}
}

package config

import (
	"encoding/json"
	"testing"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
)

func TestLoad(t *testing.T) {
	Load("/Users/tc/Documents/workspace_2020/aproxy/configs/conf.yaml")
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
									"api.proxy.com",
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
										HttpFilters: []model.HttpFilter{
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

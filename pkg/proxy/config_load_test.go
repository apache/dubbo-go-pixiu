package proxy

import (
	"encoding/json"
	"testing"
)

func TestLoad(t *testing.T) {
	Load("/Users/tc/Documents/workspace_2020/aproxy/configs/conf.yaml")
}

func TestStruct2JSON(t *testing.T) {
	b := Bootstrap{
		StaticResources: StaticResources{
			Listeners: []Listener{
				{
					Name: "net/http",
					Address: Address{
						SocketAddress: SocketAddress{
							ProtocolStr: "HTTP",
							Address:     "127.0.0.0",
							Port:        8899,
						},
					},
					Config: HttpConfig{
						IdleTimeoutStr:  "5s",
						WriteTimeoutStr: "5s",
						ReadTimeoutStr:  "5s",
					},
					FilterChains: []FilterChain{
						{
							FilterChainMatch: FilterChainMatch{
								Domains: []string{
									"api.dubbo.com",
									"api.proxy.com",
								},
							},
							Filters: []Filter{
								{
									Name: "dgp.filters.http_connect_manager",
									Config: HttpConnectionManager{
										RouteConfig: RouteConfiguration{
											Routes: []Router{
												{
													Match: RouterMatch{
														Prefix: "/api/v1",
													},
													Route: RouteAction{
														Cluster: "test_dubbo",
														Cors: CorsPolicy{
															AllowOrigin: []string{
																"*",
															},
														},
													},
												},
											},
										},
										HttpFilters: []HttpFilter{
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
			Clusters: []Cluster{
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

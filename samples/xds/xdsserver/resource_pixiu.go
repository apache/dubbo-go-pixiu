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

package main

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/xds"
	pixiupb "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/xds/model"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
	structpb2 "google.golang.org/protobuf/types/known/structpb"
)

var httpManagerConfigYaml = `
route_config:
  routes:
    - match:
        prefix: "/"
      route:
        cluster: "http_bin"
        cluster_not_found_response_code: 505
http_filters:
  - name: dgp.filter.http.httpproxy
    config:
`

func makeHttpFilter() *pixiupb.FilterChain {
	return &pixiupb.FilterChain{
		Filters: []*pixiupb.NetworkFilter{
			{
				Name: constant.HTTPConnectManagerFilter,
				//Config: &pixiupb.Filter_Yaml{Yaml: &pixiupb.Config{
				//	Content: httpManagerConfigYaml,
				//}},
				Config: &pixiupb.NetworkFilter_Struct{
					Struct: func() *structpb.Struct {
						v, err := structpb2.NewStruct(map[string]interface{}{
							"route_config": map[string]interface{}{
								"routes": []interface{}{
									map[string]interface{}{
										"match": map[string]interface{}{
											"prefix": "/",
										},
										"route": map[string]interface{}{
											"cluster":                         "http_bin",
											"cluster_not_found_response_code": 505,
										},
									},
								},
							},
							"http_filters": []interface{}{
								map[string]interface{}{
									"name":   "dgp.filter.http.httpproxy",
									"config": nil,
								},
							},
						})
						if err != nil {
							panic(err)
						}
						return v
					}(),
				},
			},
		},
	}
}
func makeListeners() *pixiupb.PixiuExtensionListeners {
	return &pixiupb.PixiuExtensionListeners{
		Listeners: []*pixiupb.Listener{
			{
				Name: "net/http",
				Address: &pixiupb.Address{
					SocketAddress: &pixiupb.SocketAddress{
						Address: "0.0.0.0",
						Port:    8888,
					},
					Name: "http_8888",
				},
				FilterChain: makeHttpFilter(),
			},
		},
	}
}

func makeClusters() *pixiupb.PixiuExtensionClusters {
	return &pixiupb.PixiuExtensionClusters{
		Clusters: []*pixiupb.Cluster{
			{
				Name:    "http_bin",
				TypeStr: "http",
				Endpoints: &pixiupb.Endpoint{
					Id: "backend",
					Address: &pixiupb.SocketAddress{
						Address: "httpbin.org",
						Port:    80,
					},
				},
				HealthChecks: []*pixiupb.HealthCheck{},
			},
		},
	}
}

func GenerateSnapshotPixiu() cache.Snapshot {
	ldsResource, _ := anypb.New(makeListeners())
	cdsResource, _ := anypb.New(makeClusters())
	snap, _ := cache.NewSnapshot("2",
		map[resource.Type][]types.Resource{
			resource.ExtensionConfigType: {
				&core.TypedExtensionConfig{
					Name:        xds.ClusterType,
					TypedConfig: cdsResource,
				},
				&core.TypedExtensionConfig{
					Name:        xds.ListenerType,
					TypedConfig: ldsResource,
				},
			},
		},
	)
	return snap
}

package main

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/model/xds"
	pixiupb "github.com/apache/dubbo-go-pixiu/pkg/model/xds/model"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

var httpManagerConfigYaml = `
route_config:
  routes:
    - match:
        prefix: "/"
      route:
        cluster: "http-baidu"
        cluster_not_found_response_code: 505
http_filters:
  - name: dgp.filter.http.httpproxy
    config:
  - name: dgp.filter.http.response
    config:
`

func makeHttpFilter() []*pixiupb.FilterChain {
	return []*pixiupb.FilterChain{
		{
			FilterChainMatch: &pixiupb.FilterChainMatch{
				Domains: []string{
					"api.dubbo.com",
					"api.pixiu.com",
				},
			},
			Filters: []*pixiupb.Filter{
				{
					Name: constant.HTTPConnectManagerFilter,
					Config: &pixiupb.Filter_Yaml{Yaml: &pixiupb.Config{
						Content: httpManagerConfigYaml,
					}},
					//Config: &pixiupb.Filter_Value{
					//	Value: func() *structpb2.Value {
					//		v, _ := structpb2.NewValue(nil)
					//		return v
					//	}(),
					//},
				},
			},
		},
	}
}
func makeListeners() *pixiupb.Listener {
	return &pixiupb.Listener{
		Name: "net/http",
		Address: &pixiupb.Address{
			SocketAddress: &pixiupb.SocketAddress{
				ProtocolStr: "http",
				Address:     "0.0.0.0",
				Port:        8888,
			},
			Name: "http_8888",
		},
		FilterChains: makeHttpFilter(),
	}
}

func makeClusters() *pixiupb.Cluster {
	return &pixiupb.Cluster{
		Name:    "http-baidu",
		TypeStr: "http",
		Endpoints: &pixiupb.Endpoint{
			Id: "backend",
			Address: &pixiupb.SocketAddress{
				ProtocolStr: "http",
				Address:     "httpbin.org",
				Port:        80,
			},
		},
	}
}

func GenerateSnapshotPixiu() cache.Snapshot {
	ldsResource, _ := anypb.New(proto.MessageV2(makeListeners()))
	cdsResource, _ := anypb.New(proto.MessageV2(makeClusters()))
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

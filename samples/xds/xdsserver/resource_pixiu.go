package main

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model/xds"
	pixiupb "github.com/apache/dubbo-go-pixiu/pkg/model/xds/model"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	"github.com/golang/protobuf/proto"
	"google.golang.org/protobuf/types/known/anypb"
)

func makeHttpFilter() []*pixiupb.FilterChain {
	return []*pixiupb.FilterChain{
		{
			FilterChainMatch: &pixiupb.FilterChainMatch{
				Domains: []string{
					"api.dubbo.com",
					"api.pixiu.com",
				},
			},
			Filters: nil,
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
		Name:    "http",
		TypeStr: "http",
		Endpoints: &pixiupb.Endpoint{
			Id: "backend",
			Address: &pixiupb.SocketAddress{
				ProtocolStr: "http",
				Address:     "127.0.0.1",
				Port:        8080,
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

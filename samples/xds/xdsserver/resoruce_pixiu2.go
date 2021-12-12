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

func makeListeners2() *pixiupb.Listener {
	return &pixiupb.Listener{
		Name: "net/http889",
		Address: &pixiupb.Address{
			SocketAddress: &pixiupb.SocketAddress{
				ProtocolStr: "http",
				Address:     "0.0.0.0",
				Port:        8889,
			},
			Name: "http_8889",
		},
		FilterChains: makeHttpFilter(),
	}
}

func GenerateSnapshotPixiu2() cache.Snapshot {
	ldsResource, _ := anypb.New(proto.MessageV2(makeListeners()))
	ldsResource2, _ := anypb.New(proto.MessageV2(makeListeners2()))
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
				&core.TypedExtensionConfig{
					Name:        xds.ListenerType,
					TypedConfig: ldsResource2,
				},
			},
		},
	)
	return snap
}

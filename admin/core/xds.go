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

package core

import (
	"context"
	"fmt"
	"net"
	"os"
	"strconv"
	"time"
)

import (
	fc "github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/xds"
	pixiupb "github.com/dubbo-go-pixiu/pixiu-api/pkg/xds/model"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	extensionpb "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	runtimeservice "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	secretservice "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
	envoyServer "github.com/envoyproxy/go-control-plane/pkg/server/v3"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
	structpb2 "google.golang.org/protobuf/types/known/structpb"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/admin/logic"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

var (
	port   = uint(18000)
	nodeID = "test-id"

	snaphost cache.SnapshotCache
)

const (
	grpcKeepaliveTime        = 30 * time.Second
	grpcKeepaliveTimeout     = 5 * time.Second
	grpcKeepaliveMinTime     = 30 * time.Second
	grpcMaxConcurrentStreams = 1000000
)

func registerServer(grpcServer *grpc.Server, server envoyServer.Server) {
	// register services
	discoverygrpc.RegisterAggregatedDiscoveryServiceServer(grpcServer, server)
	endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, server)
	clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, server)
	routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, server)
	listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, server)
	secretservice.RegisterSecretDiscoveryServiceServer(grpcServer, server)
	runtimeservice.RegisterRuntimeDiscoveryServiceServer(grpcServer, server)
	extensionpb.RegisterExtensionConfigDiscoveryServiceServer(grpcServer, server)
}

// StartxDsServer RunXDSServerWithCache starts an xDS server at the gi.ven port.
func StartxDsServer() error {
	// Create a snaphost
	snaphost = cache.NewSnapshotCache(false, cache.IDHash{}, logger.GetLogger())

	// Create the config that we'll serve to Envoy
	config := GenerateSnapshotPixiu()
	if err := config.Consistent(); err != nil {
		logger.Errorf("config inconsistency: %+v\n%+v", config, err)
		os.Exit(1)
	}

	// Add the config to the snaphost
	if err := snaphost.SetSnapshot(context.Background(), nodeID, config); err != nil {
		logger.Errorf("config error %q for %+v", err, config)
		os.Exit(1)
	}

	go watchConfigAndReload()

	// Run the xDS server
	ctx := context.Background()
	srv := envoyServer.NewServer(ctx, snaphost, nil)
	return runXDSServer(ctx, srv, port)
}

// runXDSServer starts an xDS server at the given port.
func runXDSServer(ctx context.Context, srv envoyServer.Server, port uint) error {
	// gRPC golang library sets a very small upper bound for the number gRPC/h2
	// streams over a single TCP connection. If a proxy multiplexes requests over
	// a single connection to the management server, then it might lead to
	// availability problems. Keepalive timeouts based on connection_keepalive parameter https://www.envoyproxy.io/docs/envoy/latest/configuration/overview/examples#dynamic
	var grpcOptions []grpc.ServerOption
	grpcOptions = append(grpcOptions,
		grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    grpcKeepaliveTime,
			Timeout: grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcKeepaliveMinTime,
			PermitWithoutStream: true,
		}),
	)
	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", port))
	if err != nil {
		return err
	}

	registerServer(grpcServer, srv)

	logger.Infof("management server listening on %d\n", port)
	if err = grpcServer.Serve(lis); err != nil {
		return nil
	}
	return nil
}

func watchConfigAndReload() {
	ch, err := config.Client.WatchWithPrefix(config.Bootstrap.EtcdConfig.Path)

	if err != nil {
		logger.Errorf("watch config error %q", err)
		panic(err)
	}

	for range ch {
		logger.Info("get etcd config change")
		// Create the config that we'll serve to Envoy
		config := GenerateSnapshotPixiu()
		if err := config.Consistent(); err != nil {
			logger.Errorf("config inconsistency: %+v\n%+v", config, err)
			os.Exit(1)
		}

		// Add the config to the snaphost
		if err := snaphost.SetSnapshot(context.Background(), nodeID, config); err != nil {
			logger.Errorf("config error %q for %+v", err, config)
			os.Exit(1)
		}
	}
}

// makeHTTPFilter returns a handler for the given resource.
func makeHTTPFilter(listener fc.Listener) *pixiupb.FilterChain {
	var filters, routes []any

	for _, f := range listener.HTTPFilters {
		filters = append(filters, map[string]any{
			"name":   f.Name,
			"config": f.Config,
		})
	}

	for _, r := range listener.RouteConfig.Routes {
		routes = append(filters, map[string]any{
			"match": map[string]any{
				"prefix": r.Match.Prefix,
			},
			"route": map[string]any{
				"cluster":                         r.Route.Cluster,
				"cluster_not_found_response_code": r.Route.ClusterNotFoundResponseCode,
			},
		})
	}

	return &pixiupb.FilterChain{
		Filters: []*pixiupb.NetworkFilter{
			{
				Name: constant.HTTPConnectManagerFilter,
				Config: &pixiupb.NetworkFilter_Struct{
					Struct: func() *structpb.Struct {
						v, err := structpb2.NewStruct(map[string]any{
							"route_config": map[string]any{
								"routes": routes,
							},
							"http_filters": filters,
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
	listeners, err := logic.BizGetListeners()
	if err != nil {
		logger.Errorf("get listeners error %q", err)
		return nil
	}

	if len(listeners) == 0 {
		return nil
	}

	pbListeners := &pixiupb.PixiuExtensionListeners{}
	for _, listener := range listeners {
		pbListeners.Listeners = append(pbListeners.Listeners, &pixiupb.Listener{
			Name: listener.Name,
			Address: &pixiupb.Address{
				SocketAddress: &pixiupb.SocketAddress{
					Address: listener.Address.SocketAddress.Address,
					Port:    int64(listener.Address.SocketAddress.Port),
				},
				Name: listener.Address.Name,
			},
			FilterChain: makeHTTPFilter(listener),
		})
	}
	return pbListeners
}

func makeClusters() *pixiupb.PixiuExtensionClusters {
	clusters, err := logic.BizGetClusters()
	if err != nil {
		logger.Errorf("get clusters error %q", err)
		return nil
	}

	if len(clusters) == 0 {
		return nil
	}

	pbCluster := &pixiupb.PixiuExtensionClusters{}

	for _, c := range clusters {
		pbCluster.Clusters = append(pbCluster.Clusters, &pixiupb.Cluster{
			Name:    c.Name,
			TypeStr: c.Type,
			Endpoints: []*pixiupb.Endpoint{
				{
					Id: c.Name + strconv.Itoa(c.ID),
					Address: &pixiupb.SocketAddress{
						Address: c.Address,
						Port:    int64(c.Port),
					},
				},
			},
		})
	}

	return pbCluster
}

// GenerateSnapshotPixiu returns a snapshot with a single cluster and endpoint.
func GenerateSnapshotPixiu() *cache.Snapshot {
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

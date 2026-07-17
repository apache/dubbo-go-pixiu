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
	"strconv"
	"time"
)

import (
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
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/admin/logic"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/config/xds/model"
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
// The server runs until ctx is canceled (e.g. when the admin HTTP server fails
// to start) or Serve returns an error; in either case it stops gracefully.
func StartxDsServer(ctx context.Context) error {
	// Create a snaphost
	snaphost = cache.NewSnapshotCache(false, cache.IDHash{}, logger.GetLogger())

	// Create the config that we'll serve to Envoy
	config := GenerateSnapshotPixiu()
	if err := config.Consistent(); err != nil {
		return fmt.Errorf("config inconsistency: %w", err)
	}

	// Add the config to the snaphost
	if err := snaphost.SetSnapshot(context.Background(), nodeID, config); err != nil {
		return fmt.Errorf("set snapshot error: %w", err)
	}

	go watchConfigAndReload()

	// Run the xDS server
	srv := envoyServer.NewServer(ctx, snaphost, nil)
	return runXDSServer(ctx, srv, port)
}

// runXDSServer starts an xDS server at the given port. It returns the Serve
// error, or stops the server gracefully when ctx is canceled.
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

	serveErr := make(chan error, 1)
	go func() {
		serveErr <- grpcServer.Serve(lis)
	}()

	select {
	case err := <-serveErr:
		return err
	case <-ctx.Done():
		grpcServer.GracefulStop()
		return <-serveErr
	}
}

func watchConfigAndReload() {
	const (
		maxRetries        = 5
		initialBackoff    = 1 * time.Second
		maxBackoff        = 30 * time.Second
		backoffMultiplier = 2.0
	)

	for {
		backoff := initialBackoff
		retries := 0

		// Try to establish watch with retry
		ch, err := adminconfig.Client.WatchWithPrefix(adminconfig.Bootstrap.EtcdConfig.Path)
		for err != nil && retries < maxRetries {
			retries++
			logger.Errorf("watch config error %q (retry %d/%d)", err, retries, maxRetries)

			// Wait with backoff before retrying
			time.Sleep(backoff)
			backoff = time.Duration(float64(backoff) * backoffMultiplier)
			if backoff > maxBackoff {
				backoff = maxBackoff
			}

			// Retry establishing watch
			ch, err = adminconfig.Client.WatchWithPrefix(adminconfig.Bootstrap.EtcdConfig.Path)
		}

		// Check if we failed to establish watch after max retries
		if err != nil {
			logger.Errorf("max retries reached for watch config, giving up")
			return
		}

		// Process watch events
		for range ch {
			logger.Info("get etcd config change")
			// Create the config that we'll serve to Envoy
			config := GenerateSnapshotPixiu()
			if err := config.Consistent(); err != nil {
				logger.Errorf("config inconsistency: %+v\n%+v", config, err)
				// Don't exit the process - continue running with previous valid config
				continue
			}

			// Add the config to the snaphost
			if err := snaphost.SetSnapshot(context.Background(), nodeID, config); err != nil {
				logger.Errorf("config error %q for %+v", err, config)
				// Don't exit the process - continue running with previous valid config
				continue
			}
		}

		// Channel closed, log and restart watch
		logger.Info("watch channel closed, restarting watch")
	}
}

// makeHTTPFilter returns a handler for the given resource.
func makeHTTPFilter(listener config.Listener) *model.FilterChain {
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

	return &model.FilterChain{
		Filters: []*model.NetworkFilter{
			{
				Name: constant.HTTPConnectManagerFilter,
				Config: &model.NetworkFilter_Struct{
					Struct: func() *structpb.Struct {
						v, err := structpb.NewStruct(map[string]any{
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

func makeListeners() *model.PixiuExtensionListeners {
	listeners, err := logic.BizGetListeners()
	if err != nil {
		logger.Errorf("get listeners error %q", err)
		return nil
	}

	if len(listeners) == 0 {
		return nil
	}

	pbListeners := &model.PixiuExtensionListeners{}
	for _, listener := range listeners {
		pbListeners.Listeners = append(pbListeners.Listeners, &model.Listener{
			Name: listener.Name,
			Address: &model.Address{
				SocketAddress: &model.SocketAddress{
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

func makeClusters() *model.PixiuExtensionClusters {
	clusters, err := logic.BizGetClusters()
	if err != nil {
		logger.Errorf("get clusters error %q", err)
		return nil
	}

	if len(clusters) == 0 {
		return nil
	}

	pbCluster := &model.PixiuExtensionClusters{}

	for _, c := range clusters {
		pbCluster.Clusters = append(pbCluster.Clusters, &model.Cluster{
			Name:    c.Name,
			TypeStr: c.Type,
			Endpoints: []*model.Endpoint{
				{
					Id: c.Name + strconv.Itoa(c.ID),
					Address: &model.SocketAddress{
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
					Name:        constant.ClusterType,
					TypedConfig: cdsResource,
				},
				&core.TypedExtensionConfig{
					Name:        constant.ListenerType,
					TypedConfig: ldsResource,
				},
			},
		},
	)
	return snap
}

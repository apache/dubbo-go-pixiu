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
	"time"
)

import (
	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discoverygrpc "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	endpointservice "github.com/envoyproxy/go-control-plane/envoy/service/endpoint/v3"
	extensionpb "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
	listenerservice "github.com/envoyproxy/go-control-plane/envoy/service/listener/v3"
	routeservice "github.com/envoyproxy/go-control-plane/envoy/service/route/v3"
	runtimeservice "github.com/envoyproxy/go-control-plane/envoy/service/runtime/v3"
	secretservice "github.com/envoyproxy/go-control-plane/envoy/service/secret/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	envoyServer "github.com/envoyproxy/go-control-plane/pkg/server/v3"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	adminxds "github.com/apache/dubbo-go-pixiu/admin/xds"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

var (
	snaphost        cache.SnapshotCache
	snapshotBuilder = adminxds.NewSnapshotBuilder(adminxds.LogicResourceLoader{})
)

const currentSnapshotVersion = "2"

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
	xdsConfig := adminconfig.Bootstrap.GetXDSConfig()
	adminxds.DefaultStatusStore.Reset(xdsConfig.NodeID)

	// Create a snaphost
	snaphost = cache.NewSnapshotCache(false, cache.IDHash{}, logger.GetLogger())

	// Create the config that we'll serve to Envoy
	result, err := GenerateSnapshotPixiu()
	if err != nil {
		adminxds.DefaultStatusStore.RecordError(err)
		logger.Errorf("generate xDS snapshot: %+v", err)
		os.Exit(1)
	}

	// Add the config to the snaphost
	if err := snaphost.SetSnapshot(context.Background(), xdsConfig.NodeID, result.Snapshot); err != nil {
		adminxds.DefaultStatusStore.RecordError(err)
		logger.Errorf("config error %q for %+v", err, result.Snapshot)
		os.Exit(1)
	}
	adminxds.DefaultStatusStore.RecordSuccess(result.Version, result.ListenerCount, result.ClusterCount)

	go watchConfigAndReload()

	// Run the xDS server
	ctx := context.Background()
	srv := envoyServer.NewServer(ctx, snaphost, nil)
	return runXDSServer(ctx, srv, xdsConfig.ListenPort)
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
	xdsConfig := adminconfig.Bootstrap.GetXDSConfig()
	ch, err := adminconfig.Client.WatchWithPrefix(adminconfig.Bootstrap.EtcdConfig.Path)

	if err != nil {
		logger.Errorf("watch config error %q", err)
		panic(err)
	}

	for range ch {
		logger.Info("get etcd config change")
		// Create the config that we'll serve to Envoy
		result, err := GenerateSnapshotPixiu()
		if err != nil {
			adminxds.DefaultStatusStore.RecordError(err)
			logger.Errorf("generate xDS snapshot: %+v", err)
			os.Exit(1)
		}

		// Add the config to the snaphost
		if err := snaphost.SetSnapshot(context.Background(), xdsConfig.NodeID, result.Snapshot); err != nil {
			adminxds.DefaultStatusStore.RecordError(err)
			logger.Errorf("config error %q for %+v", err, result.Snapshot)
			os.Exit(1)
		}
		adminxds.DefaultStatusStore.RecordSuccess(result.Version, result.ListenerCount, result.ClusterCount)
	}
}

// GenerateSnapshotPixiu builds and validates the current Admin resource view.
// Cache publication and version advancement remain the caller's responsibility.
func GenerateSnapshotPixiu() (*adminxds.SnapshotBuildResult, error) {
	return snapshotBuilder.Build(currentSnapshotVersion)
}

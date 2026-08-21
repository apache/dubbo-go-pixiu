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

	clientv3 "go.etcd.io/etcd/client/v3"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	adminxds "github.com/apache/dubbo-go-pixiu/admin/xds"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

var (
	snapshotCache   cache.SnapshotCache
	snapshotBuilder = adminxds.NewSnapshotBuilder(adminxds.LogicResourceLoader{})
)

const (
	grpcKeepaliveTime        = 30 * time.Second
	grpcKeepaliveTimeout     = 5 * time.Second
	grpcKeepaliveMinTime     = 30 * time.Second
	grpcMaxConcurrentStreams = 1000000
	configWatchRetryDelay    = time.Second
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
	ctx := context.Background()

	// Create a snapshot cache.
	snapshotCache = cache.NewSnapshotCache(false, cache.IDHash{}, logger.GetLogger())
	publisher := adminxds.NewSnapshotPublisher(
		xdsConfig.NodeID,
		snapshotBuilder,
		snapshotCache,
		adminxds.DefaultStatusStore,
	)

	// A failed initial candidate must not terminate Admin. The server and etcd
	// watch stay active so a later valid configuration can recover publication.
	if err := publisher.Publish(ctx); err != nil {
		logger.Errorf("initial xDS snapshot publication failed: %+v", err)
	}

	go watchConfigAndReload(ctx, publisher)

	// Run the xDS server
	srv := envoyServer.NewServer(ctx, snapshotCache, nil)
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
		return err
	}
	return nil
}

type snapshotPublisher interface {
	Publish(ctx context.Context) error
}

type configWatcher interface {
	WatchWithPrefix(key string) (clientv3.WatchChan, error)
}

func watchConfigAndReload(ctx context.Context, publisher snapshotPublisher) {
	if adminconfig.Client == nil {
		err := fmt.Errorf("watch xDS configuration: etcd client is not initialized")
		adminxds.DefaultStatusStore.RecordError(err)
		logger.Error(err)
		return
	}
	watchConfigWithRetry(ctx, adminconfig.Client, adminconfig.Bootstrap.EtcdConfig.Path, publisher, configWatchRetryDelay)
}

func watchConfigWithRetry(ctx context.Context, watcher configWatcher, path string, publisher snapshotPublisher, retryDelay time.Duration) {
	reconnecting := false
	for {
		ch, err := watcher.WatchWithPrefix(path)
		if err == nil && reconnecting {
			// Establish the new watch before rebuilding. Events that happen during
			// the rebuild remain queued on the new channel, while the rebuild
			// catches changes that occurred between the two watch sessions.
			if publishErr := publisher.Publish(ctx); publishErr != nil {
				logger.Errorf("resync xDS snapshot after watch reconnect failed: %+v", publishErr)
			}
		}
		if err == nil {
			err = consumeConfigWatch(ctx, ch, publisher)
		} else {
			err = fmt.Errorf("watch xDS configuration: %w", err)
		}
		if ctx.Err() != nil {
			return
		}

		adminxds.DefaultStatusStore.RecordError(err)
		logger.Error(err)
		reconnecting = true
		timer := time.NewTimer(retryDelay)
		select {
		case <-ctx.Done():
			if !timer.Stop() {
				<-timer.C
			}
			return
		case <-timer.C:
		}
	}
}

func consumeConfigWatch(ctx context.Context, ch clientv3.WatchChan, publisher snapshotPublisher) error {
	for {
		select {
		case <-ctx.Done():
			return nil
		case response, ok := <-ch:
			if !ok {
				if ctx.Err() != nil {
					return nil
				}
				return fmt.Errorf("watch xDS configuration: etcd watch channel closed")
			}
			if err := response.Err(); err != nil {
				return fmt.Errorf("watch xDS configuration: %w", err)
			}
			if response.Canceled {
				return fmt.Errorf("watch xDS configuration: etcd watch canceled")
			}
			if len(response.Events) == 0 {
				continue
			}
			logger.Info("get etcd config change")
			if err := publisher.Publish(ctx); err != nil {
				logger.Errorf("reload xDS snapshot failed: %+v", err)
			}
		}
	}
}

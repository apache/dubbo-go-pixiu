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

package xds

import (
	"context"
	"fmt"
	"net"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

import (
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	clusterservice "github.com/envoyproxy/go-control-plane/envoy/service/cluster/v3"
	discovery "github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
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

	"go.uber.org/zap"

	"google.golang.org/grpc"
	"google.golang.org/grpc/keepalive"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

import (
	adminmodel "github.com/apache/dubbo-go-pixiu/admin/internal/model"
	"github.com/apache/dubbo-go-pixiu/admin/internal/store"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config/xds/model"
	pkgmodel "github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	grpcKeepaliveTime        = 30 * time.Second
	grpcKeepaliveTimeout     = 5 * time.Second
	grpcKeepaliveMinTime     = 30 * time.Second
	grpcMaxConcurrentStreams = 1000000

	defaultNodeID = "test-id"
)

// Server represents an xDS server for serving configuration to Pixiu gateway.
type Server struct {
	etcd            *store.Etcd
	logger          *zap.Logger
	snapshot        cache.SnapshotCache
	nodeID          string
	port            uint
	instances       sync.Map // map[nodeID]*adminmodel.PixiuInstance
	snapshotVersion uint64   // atomic counter for snapshot versions
}

// NewServer creates a new xDS server.
func NewServer(etcd *store.Etcd, logger *zap.Logger, port uint) *Server {
	return &Server{
		etcd:   etcd,
		logger: logger,
		nodeID: defaultNodeID,
		port:   port,
	}
}

// Start starts the xDS server.
func (s *Server) Start(ctx context.Context) error {
	// Create a snapshot cache that supports any node ID
	s.snapshot = cache.NewSnapshotCache(false, cache.IDHash{}, &zapLogger{s.logger})

	// Create initial snapshot for default node
	snap := s.generateSnapshot()
	if err := snap.Consistent(); err != nil {
		s.logger.Error("config inconsistency", zap.Error(err))
		return err
	}

	// Add the config to the snapshot for default node
	if err := s.snapshot.SetSnapshot(ctx, s.nodeID, snap); err != nil {
		s.logger.Error("failed to set snapshot", zap.Error(err))
		return err
	}

	// Start watching for config changes
	go s.watchConfigAndReload(ctx)

	// Run the xDS server with callbacks
	srv := envoyServer.NewServer(ctx, s.snapshot, s.makeCallbacks(ctx))
	return s.runGRPCServer(ctx, srv)
}

// GetInstances returns all connected Pixiu instances.
func (s *Server) GetInstances() []adminmodel.PixiuInstance {
	var instances []adminmodel.PixiuInstance
	s.instances.Range(func(key, value any) bool {
		if inst, ok := value.(*adminmodel.PixiuInstance); ok {
			instances = append(instances, *inst)
		}
		return true
	})
	return instances
}

// GetInstanceCount returns the number of connected instances.
func (s *Server) GetInstanceCount() int {
	count := 0
	s.instances.Range(func(key, value any) bool {
		count++
		return true
	})
	return count
}

// runGRPCServer starts the gRPC server for xDS.
func (s *Server) runGRPCServer(ctx context.Context, srv envoyServer.Server) error {
	grpcOptions := []grpc.ServerOption{
		grpc.MaxConcurrentStreams(grpcMaxConcurrentStreams),
		grpc.KeepaliveParams(keepalive.ServerParameters{
			Time:    grpcKeepaliveTime,
			Timeout: grpcKeepaliveTimeout,
		}),
		grpc.KeepaliveEnforcementPolicy(keepalive.EnforcementPolicy{
			MinTime:             grpcKeepaliveMinTime,
			PermitWithoutStream: true,
		}),
	}
	grpcServer := grpc.NewServer(grpcOptions...)

	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", s.port))
	if err != nil {
		return err
	}

	// Register all xDS services
	discovery.RegisterAggregatedDiscoveryServiceServer(grpcServer, srv)
	endpointservice.RegisterEndpointDiscoveryServiceServer(grpcServer, srv)
	clusterservice.RegisterClusterDiscoveryServiceServer(grpcServer, srv)
	routeservice.RegisterRouteDiscoveryServiceServer(grpcServer, srv)
	listenerservice.RegisterListenerDiscoveryServiceServer(grpcServer, srv)
	secretservice.RegisterSecretDiscoveryServiceServer(grpcServer, srv)
	runtimeservice.RegisterRuntimeDiscoveryServiceServer(grpcServer, srv)
	extensionpb.RegisterExtensionConfigDiscoveryServiceServer(grpcServer, srv)

	s.logger.Info("xDS management server starting", zap.Uint("port", s.port))

	return grpcServer.Serve(lis)
}

// watchConfigAndReload watches etcd for config changes and reloads the snapshot.
func (s *Server) watchConfigAndReload(ctx context.Context) {
	ch, err := s.etcd.WatchConfig()
	if err != nil {
		s.logger.Error("failed to watch config", zap.Error(err))
		return
	}

	for {
		select {
		case <-ctx.Done():
			return
		case <-ch:
			s.logger.Info("etcd config changed, reloading snapshots for all nodes")

			snap := s.generateSnapshot()
			if err := snap.Consistent(); err != nil {
				s.logger.Error("config inconsistency on reload", zap.Error(err))
				continue
			}

			// Update snapshot for default node
			if err := s.snapshot.SetSnapshot(ctx, s.nodeID, snap); err != nil {
				s.logger.Error("failed to set snapshot for default node", zap.Error(err))
			}

			// Update snapshots for all connected instances
			s.instances.Range(func(key, value any) bool {
				nodeID := key.(string)
				if nodeID != s.nodeID {
					if err := s.snapshot.SetSnapshot(ctx, nodeID, snap); err != nil {
						s.logger.Error("failed to set snapshot for node", zap.String("nodeID", nodeID), zap.Error(err))
					}
				}
				return true
			})

			s.logger.Info("snapshots reloaded successfully", zap.Int("instances", s.GetInstanceCount()))
		}
	}
}

// nextVersion returns the next snapshot version number (atomically incremented).
func (s *Server) nextVersion() string {
	version := atomic.AddUint64(&s.snapshotVersion, 1)
	return strconv.FormatUint(version, 10)
}

// generateSnapshot generates a snapshot from etcd configuration.
func (s *Server) generateSnapshot() *cache.Snapshot {
	ldsResource, _ := anypb.New(s.makeListeners())
	cdsResource, _ := anypb.New(s.makeClusters())
	rdsResource, _ := anypb.New(s.makeRoutes())

	version := s.nextVersion()
	s.logger.Debug("generating snapshot", zap.String("version", version))

	snap, _ := cache.NewSnapshot(version,
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
				&core.TypedExtensionConfig{
					Name:        constant.RouterType,
					TypedConfig: rdsResource,
				},
			},
		},
	)
	return snap
}

// makeListeners creates the listeners configuration for xDS.
// Uses pkg/config.Listener format (compatible with old admin).
func (s *Server) makeListeners() *model.PixiuExtensionListeners {
	listeners, err := s.etcd.ListListeners()
	if err != nil {
		s.logger.Error("failed to get listeners", zap.Error(err))
		return nil
	}

	if len(listeners) == 0 {
		return nil
	}

	pbListeners := &model.PixiuExtensionListeners{}
	for _, listener := range listeners {
		// Convert protocol string to enum
		protocol := model.Listener_HTTP
		if p, ok := model.Listener_Protocols_value[listener.ProtocolStr]; ok {
			protocol = model.Listener_Protocols(p)
		}

		pbListeners.Listeners = append(pbListeners.Listeners, &model.Listener{
			Name: listener.Name,
			Address: &model.Address{
				SocketAddress: &model.SocketAddress{
					Address: listener.Address.SocketAddress.Address,
					Port:    int64(listener.Address.SocketAddress.Port),
				},
				Name: listener.Address.Name,
			},
			FilterChain: s.convertFilterChain(listener.FilterChain),
			Protocol:    protocol,
		})
	}
	return pbListeners
}

// convertFilterChain converts pkg/model.FilterChain to xDS model.FilterChain.
func (s *Server) convertFilterChain(fc pkgmodel.FilterChain) *model.FilterChain {
	var pbFilters []*model.NetworkFilter
	for _, f := range fc.Filters {
		pbFilters = append(pbFilters, &model.NetworkFilter{
			Name: f.Name,
			Config: &model.NetworkFilter_Struct{
				Struct: func() *structpb.Struct {
					if f.Config == nil {
						return nil
					}
					v, err := structpb.NewStruct(f.Config)
					if err != nil {
						s.logger.Error("failed to create struct", zap.Error(err))
						return nil
					}
					return v
				}(),
			},
		})
	}

	return &model.FilterChain{
		Filters: pbFilters,
	}
}

// makeClusters creates the clusters configuration for xDS using Pixiu's ClusterConfig format.
func (s *Server) makeClusters() *model.PixiuExtensionClusters {
	clusters, err := s.etcd.ListClusters()
	if err != nil {
		s.logger.Error("failed to get clusters", zap.Error(err))
		return nil
	}

	if len(clusters) == 0 {
		return nil
	}

	pbCluster := &model.PixiuExtensionClusters{}
	for _, c := range clusters {
		// Convert endpoints to xDS model.Endpoint
		var endpoints []*model.Endpoint
		for i, ep := range c.Endpoints {
			endpoints = append(endpoints, &model.Endpoint{
				Id: c.Name + strconv.Itoa(i),
				Address: &model.SocketAddress{
					Address: ep.Address.Address,
					Port:    int64(ep.Address.Port),
				},
			})
		}

		pbCluster.Clusters = append(pbCluster.Clusters, &model.Cluster{
			Name:      c.Name,
			TypeStr:   c.TypeStr,
			LbStr:     string(c.LbStr),
			Endpoints: endpoints,
		})
	}

	return pbCluster
}

// makeRoutes creates the routes configuration for xDS.
// Routes are extracted from filter_chains.filters[].config.route_config.routes in pkg/model.Listener.
func (s *Server) makeRoutes() *model.RouteConfiguration {
	listeners, err := s.etcd.ListListeners()
	if err != nil {
		s.logger.Error("failed to get listeners for routes", zap.Error(err))
		return nil
	}

	if len(listeners) == 0 {
		return nil
	}

	routeConfig := &model.RouteConfiguration{
		Dynamic: true,
	}

	// Extract routes from each listener's filter chain
	for _, listener := range listeners {
		for _, filter := range listener.FilterChain.Filters {
			// Look for httpconnectionmanager filter which contains routes
			if filter.Name != constant.HTTPConnectManagerFilter {
				continue
			}
			if filter.Config == nil {
				continue
			}

			// Extract route_config from filter config
			rcRaw, ok := filter.Config["route_config"]
			if !ok {
				continue
			}
			rc, ok := rcRaw.(map[string]any)
			if !ok {
				continue
			}

			// Extract routes array
			routesRaw, ok := rc["routes"]
			if !ok {
				continue
			}
			routes, ok := routesRaw.([]any)
			if !ok {
				continue
			}

			// Convert each route
			for _, routeRaw := range routes {
				route, ok := routeRaw.(map[string]any)
				if !ok {
					continue
				}

				// Extract match.prefix
				var prefix string
				if matchRaw, ok := route["match"].(map[string]any); ok {
					if p, ok := matchRaw["prefix"].(string); ok {
						prefix = p
					}
				}

				// Extract route.cluster
				var cluster string
				var clusterNotFoundCode int64
				if routeActionRaw, ok := route["route"].(map[string]any); ok {
					if c, ok := routeActionRaw["cluster"].(string); ok {
						cluster = c
					}
					if code, ok := routeActionRaw["cluster_not_found_response_code"].(float64); ok {
						clusterNotFoundCode = int64(code)
					}
				}

				if prefix == "" || cluster == "" {
					continue
				}

				routeConfig.Routes = append(routeConfig.Routes, &model.Router{
					Match: &model.RouterMatch{
						Prefix: prefix,
					},
					Route: &model.RouteAction{
						Cluster:                     cluster,
						ClusterNotFoundResponseCode: clusterNotFoundCode,
					},
				})
			}
		}
	}

	if len(routeConfig.Routes) == 0 {
		return nil
	}

	s.logger.Debug("extracted routes from listeners",
		zap.Int("routeCount", len(routeConfig.Routes)))

	return routeConfig
}

// zapLogger adapts zap.Logger for the go-control-plane cache.
type zapLogger struct {
	*zap.Logger
}

func (l *zapLogger) Debugf(format string, args ...any) {
	l.Sugar().Debugf(format, args...)
}

func (l *zapLogger) Infof(format string, args ...any) {
	l.Sugar().Infof(format, args...)
}

func (l *zapLogger) Warnf(format string, args ...any) {
	l.Sugar().Warnf(format, args...)
}

func (l *zapLogger) Errorf(format string, args ...any) {
	l.Sugar().Errorf(format, args...)
}

// makeCallbacks creates xDS server callbacks for tracking instance connections.
func (s *Server) makeCallbacks(ctx context.Context) *envoyServer.CallbackFuncs {
	return &envoyServer.CallbackFuncs{
		StreamClosedFunc: func(streamID int64, node *core.Node) {
			if node != nil && node.Id != "" {
				s.logger.Info("stream closed", zap.Int64("streamID", streamID), zap.String("nodeID", node.Id))
				if inst, ok := s.instances.Load(node.Id); ok {
					instance := inst.(*adminmodel.PixiuInstance)
					instance.Status = adminmodel.InstanceStatusDisconnected
					instance.LastSeen = time.Now()
				}
			}
		},
		StreamRequestFunc: func(streamID int64, req *discovery.DiscoveryRequest) error {
			if req.Node == nil || req.Node.Id == "" {
				return nil
			}
			s.registerOrUpdateInstance(ctx, req.Node)
			return nil
		},
		DeltaStreamClosedFunc: func(streamID int64, node *core.Node) {
			if node != nil && node.Id != "" {
				if inst, ok := s.instances.Load(node.Id); ok {
					instance := inst.(*adminmodel.PixiuInstance)
					instance.Status = adminmodel.InstanceStatusDisconnected
					instance.LastSeen = time.Now()
				}
			}
		},
	}
}

// registerOrUpdateInstance registers a new instance or updates an existing one.
func (s *Server) registerOrUpdateInstance(ctx context.Context, node *core.Node) {
	nodeID := node.Id
	now := time.Now()

	if _, ok := s.instances.Load(nodeID); !ok {
		// Register new instance
		instance := &adminmodel.PixiuInstance{
			NodeID:      nodeID,
			Cluster:     node.Cluster,
			Status:      adminmodel.InstanceStatusConnected,
			LastSeen:    now,
			ConnectedAt: now,
			Metadata:    make(map[string]string),
		}

		// Extract metadata
		if node.Metadata != nil {
			for k, v := range node.Metadata.Fields {
				if sv := v.GetStringValue(); sv != "" {
					instance.Metadata[k] = sv
				}
			}
		}

		// Extract address from locality
		if node.Locality != nil {
			instance.Address = node.Locality.Zone
		}

		// Extract version from user agent
		if node.UserAgentName != "" {
			instance.Version = node.UserAgentName
			if v, ok := node.UserAgentVersionType.(*core.Node_UserAgentVersion); ok {
				instance.Version += "/" + v.UserAgentVersion
			}
		}

		s.instances.Store(nodeID, instance)
		s.logger.Info("new instance connected",
			zap.String("nodeID", nodeID),
			zap.String("cluster", instance.Cluster),
			zap.String("version", instance.Version))

		// Create snapshot for this node
		snap := s.generateSnapshot()
		if err := s.snapshot.SetSnapshot(ctx, nodeID, snap); err != nil {
			s.logger.Error("failed to set snapshot for new node", zap.String("nodeID", nodeID), zap.Error(err))
		}
	} else {
		// Update last seen time
		if inst, ok := s.instances.Load(nodeID); ok {
			instance := inst.(*adminmodel.PixiuInstance)
			instance.LastSeen = now
			instance.Status = adminmodel.InstanceStatusConnected
		}
	}
}

/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to you under the Apache License, Version 2.0
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
	"fmt"
	"strconv"
)

import (
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/envoyproxy/go-control-plane/pkg/cache/types"
	cache "github.com/envoyproxy/go-control-plane/pkg/cache/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"

	"google.golang.org/protobuf/types/known/anypb"
	"google.golang.org/protobuf/types/known/structpb"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/admin/logic"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	xdsmodel "github.com/apache/dubbo-go-pixiu/pkg/config/xds/model"
)

// ResourceLoader reads one logical Admin resource view for snapshot building.
// A successful empty slice means that the resource type is intentionally
// empty; an error means no candidate snapshot may be published.
type ResourceLoader interface {
	LoadListeners() ([]config.Listener, error)
	LoadClusters() ([]config.Cluster, error)
}

// LogicResourceLoader reads the published listener and cluster configuration
// through the existing Admin logic layer.
type LogicResourceLoader struct{}

func (LogicResourceLoader) LoadListeners() ([]config.Listener, error) {
	if adminconfig.Client == nil {
		return nil, fmt.Errorf("admin etcd client is not initialized")
	}
	return logic.BizGetListeners()
}

func (LogicResourceLoader) LoadClusters() ([]config.Cluster, error) {
	if adminconfig.Client == nil {
		return nil, fmt.Errorf("admin etcd client is not initialized")
	}
	return logic.BizGetClusters()
}

// SnapshotBuildResult carries the validated snapshot and the metadata needed
// to update Admin diagnostics after successful cache publication.
type SnapshotBuildResult struct {
	Snapshot      *cache.Snapshot
	Version       string
	ListenerCount int
	ClusterCount  int
}

type SnapshotBuilder struct {
	loader ResourceLoader
}

func NewSnapshotBuilder(loader ResourceLoader) *SnapshotBuilder {
	return &SnapshotBuilder{loader: loader}
}

// Build loads, converts, packs, and validates a candidate snapshot. It does
// not mutate the SnapshotCache, version state, or last-good diagnostics.
func (b *SnapshotBuilder) Build(version string) (*SnapshotBuildResult, error) {
	if b == nil || b.loader == nil {
		return nil, fmt.Errorf("xDS snapshot resource loader is not configured")
	}
	if version == "" {
		return nil, fmt.Errorf("xDS snapshot version is empty")
	}

	listeners, err := b.loader.LoadListeners()
	if err != nil {
		return nil, fmt.Errorf("load listeners: %w", err)
	}
	clusters, err := b.loader.LoadClusters()
	if err != nil {
		return nil, fmt.Errorf("load clusters: %w", err)
	}

	listenerConfig, err := makeListeners(listeners)
	if err != nil {
		return nil, fmt.Errorf("convert listeners: %w", err)
	}
	clusterConfig := makeClusters(clusters)

	listenerResource, err := anypb.New(listenerConfig)
	if err != nil {
		return nil, fmt.Errorf("pack listener ExtensionConfig: %w", err)
	}
	clusterResource, err := anypb.New(clusterConfig)
	if err != nil {
		return nil, fmt.Errorf("pack cluster ExtensionConfig: %w", err)
	}

	snapshot, err := cache.NewSnapshot(version, map[resource.Type][]types.Resource{
		resource.ExtensionConfigType: {
			&corev3.TypedExtensionConfig{
				Name:        constant.ClusterType,
				TypedConfig: clusterResource,
			},
			&corev3.TypedExtensionConfig{
				Name:        constant.ListenerType,
				TypedConfig: listenerResource,
			},
		},
	})
	if err != nil {
		return nil, fmt.Errorf("construct xDS snapshot: %w", err)
	}
	if err := snapshot.Consistent(); err != nil {
		return nil, fmt.Errorf("validate xDS snapshot: %w", err)
	}

	return &SnapshotBuildResult{
		Snapshot:      snapshot,
		Version:       version,
		ListenerCount: len(listenerConfig.Listeners),
		ClusterCount:  len(clusterConfig.Clusters),
	}, nil
}

func makeHTTPFilter(listener config.Listener) (*xdsmodel.FilterChain, error) {
	filters := make([]any, 0, len(listener.HTTPFilters))
	routes := make([]any, 0, len(listener.RouteConfig.Routes))

	for _, filter := range listener.HTTPFilters {
		filters = append(filters, map[string]any{
			"name":   filter.Name,
			"config": filter.Config,
		})
	}

	for _, route := range listener.RouteConfig.Routes {
		routes = append(routes, map[string]any{
			"match": map[string]any{
				"prefix": route.Match.Prefix,
			},
			"route": map[string]any{
				"cluster":                         route.Route.Cluster,
				"cluster_not_found_response_code": route.Route.ClusterNotFoundResponseCode,
			},
		})
	}

	filterConfig, err := structpb.NewStruct(map[string]any{
		"route_config": map[string]any{
			"routes": routes,
		},
		"http_filters": filters,
	})
	if err != nil {
		return nil, fmt.Errorf("build HTTP connection manager config: %w", err)
	}

	return &xdsmodel.FilterChain{
		Filters: []*xdsmodel.NetworkFilter{
			{
				Name: constant.HTTPConnectManagerFilter,
				Config: &xdsmodel.NetworkFilter_Struct{
					Struct: filterConfig,
				},
			},
		},
	}, nil
}

func makeListeners(listeners []config.Listener) (*xdsmodel.PixiuExtensionListeners, error) {
	result := &xdsmodel.PixiuExtensionListeners{
		Listeners: make([]*xdsmodel.Listener, 0, len(listeners)),
	}

	for _, listener := range listeners {
		filterChain, err := makeHTTPFilter(listener)
		if err != nil {
			return nil, fmt.Errorf("listener %q: %w", listener.Name, err)
		}
		result.Listeners = append(result.Listeners, &xdsmodel.Listener{
			Name: listener.Name,
			Address: &xdsmodel.Address{
				SocketAddress: &xdsmodel.SocketAddress{
					Address: listener.Address.SocketAddress.Address,
					Port:    int64(listener.Address.SocketAddress.Port),
				},
				Name: listener.Address.Name,
			},
			FilterChain: filterChain,
		})
	}
	return result, nil
}

func makeClusters(clusters []config.Cluster) *xdsmodel.PixiuExtensionClusters {
	result := &xdsmodel.PixiuExtensionClusters{
		Clusters: make([]*xdsmodel.Cluster, 0, len(clusters)),
	}
	for _, cluster := range clusters {
		result.Clusters = append(result.Clusters, &xdsmodel.Cluster{
			Name:    cluster.Name,
			TypeStr: cluster.Type,
			Endpoints: []*xdsmodel.Endpoint{
				{
					Id: cluster.Name + strconv.Itoa(cluster.ID),
					Address: &xdsmodel.SocketAddress{
						Address: cluster.Address,
						Port:    int64(cluster.Port),
					},
				},
			},
		})
	}
	return result
}

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
	"errors"
	"strings"
	"testing"
)

import (
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	xdsmodel "github.com/apache/dubbo-go-pixiu/pkg/config/xds/model"
)

type fakeResourceLoader struct {
	listeners   []config.Listener
	clusters    []config.Cluster
	listenerErr error
	clusterErr  error
}

func (f fakeResourceLoader) LoadListeners() ([]config.Listener, error) {
	return f.listeners, f.listenerErr
}

func (f fakeResourceLoader) LoadClusters() ([]config.Cluster, error) {
	return f.clusters, f.clusterErr
}

func TestSnapshotBuilderBuildsValidEmptyResources(t *testing.T) {
	result, err := NewSnapshotBuilder(fakeResourceLoader{}).Build("7")
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if result.Version != "7" || result.ListenerCount != 0 || result.ClusterCount != 0 {
		t.Fatalf("unexpected build metadata: %+v", result)
	}
	if got := result.Snapshot.GetVersion(resource.ExtensionConfigType); got != "7" {
		t.Fatalf("snapshot version: want 7, got %q", got)
	}

	listeners := unpackListeners(t, result)
	clusters := unpackClusters(t, result)
	if len(listeners.Listeners) != 0 {
		t.Fatalf("empty listeners were not represented as an empty resource: %+v", listeners)
	}
	if len(clusters.Clusters) != 0 {
		t.Fatalf("empty clusters were not represented as an empty resource: %+v", clusters)
	}
}

func TestSnapshotBuilderConvertsRoutesAndFiltersSeparately(t *testing.T) {
	listener := config.Listener{Name: "http"}
	listener.Address.SocketAddress.Address = "0.0.0.0"
	listener.Address.SocketAddress.Port = 8080
	listener.HTTPFilters = config.HTTPFilters{{Name: "request-filter", Config: map[string]any{"enabled": true}}}
	listener.RouteConfig.Routes = make([]struct {
		Match struct {
			Prefix string `yaml:"prefix" json:"prefix"`
		} `yaml:"match" json:"match"`
		Route struct {
			Cluster                     string `yaml:"cluster" json:"cluster"`
			ClusterNotFoundResponseCode int    `yaml:"cluster_not_found_response_code" json:"cluster_not_found_response_code"`
		} `yaml:"route" json:"route"`
	}, 1)
	listener.RouteConfig.Routes[0].Match.Prefix = "/api"
	listener.RouteConfig.Routes[0].Route.Cluster = "backend"

	result, err := NewSnapshotBuilder(fakeResourceLoader{listeners: []config.Listener{listener}}).Build("8")
	if err != nil {
		t.Fatalf("Build: %v", err)
	}

	listeners := unpackListeners(t, result)
	if len(listeners.Listeners) != 1 {
		t.Fatalf("listener count: want 1, got %d", len(listeners.Listeners))
	}
	filterConfig := listeners.Listeners[0].FilterChain.Filters[0].GetStruct().AsMap()
	routes := filterConfig["route_config"].(map[string]any)["routes"].([]any)
	filters := filterConfig["http_filters"].([]any)
	if len(routes) != 1 || len(filters) != 1 {
		t.Fatalf("route/filter conversion mixed slices: routes=%v filters=%v", routes, filters)
	}
	if got := routes[0].(map[string]any)["route"].(map[string]any)["cluster"]; got != "backend" {
		t.Fatalf("route cluster: want backend, got %v", got)
	}
}

func TestSnapshotBuilderConvertsClusters(t *testing.T) {
	result, err := NewSnapshotBuilder(fakeResourceLoader{clusters: []config.Cluster{
		{
			Name:    "backend",
			Type:    "Static",
			Address: "127.0.0.1",
			Port:    20880,
			ID:      4,
		},
	}}).Build("9")
	if err != nil {
		t.Fatalf("Build: %v", err)
	}
	if result.ClusterCount != 1 {
		t.Fatalf("cluster count: want 1, got %d", result.ClusterCount)
	}

	clusters := unpackClusters(t, result)
	cluster := clusters.Clusters[0]
	if cluster.Name != "backend" || cluster.TypeStr != "Static" {
		t.Fatalf("unexpected cluster identity: %+v", cluster)
	}
	if len(cluster.Endpoints) != 1 || cluster.Endpoints[0].Id != "backend4" ||
		cluster.Endpoints[0].Address.Address != "127.0.0.1" ||
		cluster.Endpoints[0].Address.Port != 20880 {
		t.Fatalf("unexpected cluster endpoints: %+v", cluster.Endpoints)
	}
}

func TestSnapshotBuilderReturnsFilterConversionError(t *testing.T) {
	listener := config.Listener{Name: "invalid-filter"}
	listener.Address.SocketAddress.Address = "0.0.0.0"
	listener.Address.SocketAddress.Port = 8080
	listener.HTTPFilters = config.HTTPFilters{{
		Name:   "bad",
		Config: make(chan int),
	}}

	_, err := NewSnapshotBuilder(fakeResourceLoader{listeners: []config.Listener{listener}}).Build("9")
	if err == nil || !strings.Contains(err.Error(), `convert listeners: listener "invalid-filter"`) {
		t.Fatalf("Build error did not identify the invalid listener: %v", err)
	}
}

func TestSnapshotBuilderRejectsInvalidCustomResources(t *testing.T) {
	validListener := config.Listener{Name: "http"}
	validListener.Address.SocketAddress.Address = "0.0.0.0"
	validListener.Address.SocketAddress.Port = 8080
	validCluster := config.Cluster{Name: "backend", Type: "Static", Address: "127.0.0.1", Port: 20880}

	tests := []struct {
		name   string
		loader fakeResourceLoader
		want   string
	}{
		{
			name:   "empty cluster name",
			loader: fakeResourceLoader{clusters: []config.Cluster{{Type: "Static", Address: "127.0.0.1", Port: 20880}}},
			want:   "empty name",
		},
		{
			name: "duplicate cluster",
			loader: fakeResourceLoader{clusters: []config.Cluster{
				validCluster,
				{Name: "backend", Type: "Static", Address: "127.0.0.2", Port: 20880},
			}},
			want: "duplicate cluster name",
		},
		{
			name:   "invalid cluster endpoint",
			loader: fakeResourceLoader{clusters: []config.Cluster{{Name: "backend", Type: "Static"}}},
			want:   "invalid endpoint address",
		},
		{
			name:   "invalid cluster type",
			loader: fakeResourceLoader{clusters: []config.Cluster{{Name: "backend", Type: "", Address: "127.0.0.1", Port: 20880}}},
			want:   "unsupported type",
		},
		{
			name:   "empty listener name",
			loader: fakeResourceLoader{listeners: []config.Listener{{}}},
			want:   "empty name",
		},
		{
			name: "duplicate listener address",
			loader: fakeResourceLoader{listeners: []config.Listener{
				validListener,
				func() config.Listener {
					duplicate := validListener
					duplicate.Name = "other"
					return duplicate
				}(),
			}},
			want: "duplicate listener socket address",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSnapshotBuilder(tt.loader).Build("10")
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Build error: want %q, got %v", tt.want, err)
			}
		})
	}
}

func TestSnapshotBuilderPropagatesLoadErrors(t *testing.T) {
	tests := []struct {
		name   string
		loader fakeResourceLoader
		want   string
	}{
		{
			name:   "listeners",
			loader: fakeResourceLoader{listenerErr: errors.New("etcd unavailable")},
			want:   "load listeners: etcd unavailable",
		},
		{
			name:   "clusters",
			loader: fakeResourceLoader{clusterErr: errors.New("etcd unavailable")},
			want:   "load clusters: etcd unavailable",
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			_, err := NewSnapshotBuilder(tt.loader).Build("9")
			if err == nil || !strings.Contains(err.Error(), tt.want) {
				t.Fatalf("Build error: want %q, got %v", tt.want, err)
			}
		})
	}
}

func TestSnapshotBuilderRejectsMissingInputs(t *testing.T) {
	if _, err := NewSnapshotBuilder(nil).Build("1"); err == nil {
		t.Fatal("expected nil loader error")
	}
	if _, err := NewSnapshotBuilder(fakeResourceLoader{}).Build(""); err == nil {
		t.Fatal("expected empty version error")
	}
}

func TestLogicResourceLoaderRejectsMissingEtcdClient(t *testing.T) {
	previous := adminconfig.Client
	adminconfig.Client = nil
	t.Cleanup(func() { adminconfig.Client = previous })

	if _, err := (LogicResourceLoader{}).LoadListeners(); err == nil {
		t.Fatal("expected listener load to reject a missing etcd client")
	}
	if _, err := (LogicResourceLoader{}).LoadClusters(); err == nil {
		t.Fatal("expected cluster load to reject a missing etcd client")
	}
}

func unpackListeners(t *testing.T, result *SnapshotBuildResult) *xdsmodel.PixiuExtensionListeners {
	t.Helper()
	resourceMap := result.Snapshot.GetResources(resource.ExtensionConfigType)
	resourceValue, ok := resourceMap[constant.ListenerType]
	if !ok {
		t.Fatalf("listener ExtensionConfig is missing: %v", resourceMap)
	}
	typedConfig, ok := resourceValue.(*corev3.TypedExtensionConfig)
	if !ok {
		t.Fatalf("listener resource type: %T", resourceValue)
	}
	listeners := &xdsmodel.PixiuExtensionListeners{}
	if err := typedConfig.TypedConfig.UnmarshalTo(listeners); err != nil {
		t.Fatalf("unpack listeners: %v", err)
	}
	return listeners
}

func unpackClusters(t *testing.T, result *SnapshotBuildResult) *xdsmodel.PixiuExtensionClusters {
	t.Helper()
	resourceMap := result.Snapshot.GetResources(resource.ExtensionConfigType)
	resourceValue, ok := resourceMap[constant.ClusterType]
	if !ok {
		t.Fatalf("cluster ExtensionConfig is missing: %v", resourceMap)
	}
	typedConfig, ok := resourceValue.(*corev3.TypedExtensionConfig)
	if !ok {
		t.Fatalf("cluster resource type: %T", resourceValue)
	}
	clusters := &xdsmodel.PixiuExtensionClusters{}
	if err := typedConfig.TypedConfig.UnmarshalTo(clusters); err != nil {
		t.Fatalf("unpack clusters: %v", err)
	}
	return clusters
}

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
	"testing"
)

import (
	corev3 "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/envoyproxy/go-control-plane/pkg/resource/v3"

	"github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"
)

import (
	adminxds "github.com/apache/dubbo-go-pixiu/admin/xds"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/config/xds/apiclient"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server/controls/mocks"
)

type integrationResourceLoader struct {
	listeners []config.Listener
	clusters  []config.Cluster
}

func (l *integrationResourceLoader) LoadListeners() ([]config.Listener, error) {
	return l.listeners, nil
}

func (l *integrationResourceLoader) LoadClusters() ([]config.Cluster, error) {
	return l.clusters, nil
}

// TestAdminSnapshotComponentIntegration keeps conversion failures cheap to
// diagnose. The deployment-level path lives in
// TestAdminHTTPToRunningPixiuEndToEnd and must not be represented by this test.
func TestAdminSnapshotComponentIntegration(t *testing.T) {
	loader := &integrationResourceLoader{}
	clusterState := make(map[string]*model.ClusterConfig)
	clusterController := gomock.NewController(t)
	clusterManager := mocks.NewMockClusterManager(clusterController)
	replacingClusterManager := &replaceXDSClusterManager{ClusterManager: clusterManager, replace: func(clusters []*model.ClusterConfig) error {
		clear(clusterState)
		for _, cluster := range clusters {
			clusterState[cluster.Name] = cluster
		}
		return nil
	}}
	listenerManager := &mockListenerManager{m: make(map[string]*model.Listener)}
	cdsManager := &CdsManager{clusterMg: replacingClusterManager}
	ldsManager := &LdsManager{listenerMg: listenerManager}

	listener := config.Listener{Name: "gateway"}
	listener.Address.SocketAddress.Address = "0.0.0.0"
	listener.Address.SocketAddress.Port = 18080
	loader.listeners = []config.Listener{listener}
	loader.clusters = []config.Cluster{{
		Name:    "backend",
		Type:    "Static",
		Address: "10.0.0.1",
		Port:    20880,
		ID:      1,
	}}
	applyAdminSnapshotToPixiu(t, "1", loader, cdsManager, ldsManager)
	require.Equal(t, "10.0.0.1", clusterState["backend"].Endpoints[0].Address.Address)
	require.Contains(t, listenerManager.m, "0.0.0.0-18080-HTTP")

	loader.clusters[0].Address = "10.0.0.2"
	loader.listeners[0].Address.SocketAddress.Port = 18081
	applyAdminSnapshotToPixiu(t, "2", loader, cdsManager, ldsManager)
	require.Equal(t, "10.0.0.2", clusterState["backend"].Endpoints[0].Address.Address)
	require.NotContains(t, listenerManager.m, "0.0.0.0-18080-HTTP")
	require.Contains(t, listenerManager.m, "0.0.0.0-18081-HTTP")

	loader.clusters = nil
	loader.listeners = nil
	applyAdminSnapshotToPixiu(t, "3", loader, cdsManager, ldsManager)
	require.Empty(t, clusterState)
	require.Empty(t, listenerManager.m)
}

func applyAdminSnapshotToPixiu(
	t *testing.T,
	version string,
	loader *integrationResourceLoader,
	cdsManager *CdsManager,
	ldsManager *LdsManager,
) {
	t.Helper()
	result, err := adminxds.NewSnapshotBuilder(loader).Build(version)
	require.NoError(t, err)
	resources := result.Snapshot.GetResources(resource.ExtensionConfigType)

	clusterResource, ok := resources[constant.ClusterType].(*corev3.TypedExtensionConfig)
	require.True(t, ok)
	require.NoError(t, cdsManager.applyDelta(&apiclient.DeltaResources{
		NewResources: []*apiclient.ProtoAny{apiclient.NewProtoAny(clusterResource)},
	}))

	listenerResource, ok := resources[constant.ListenerType].(*corev3.TypedExtensionConfig)
	require.True(t, ok)
	require.NoError(t, ldsManager.applyDelta(&apiclient.DeltaResources{
		NewResources: []*apiclient.ProtoAny{apiclient.NewProtoAny(listenerResource)},
	}))
}

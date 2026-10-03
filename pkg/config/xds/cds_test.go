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
	stderr "errors"
	"testing"
)

import (
	"github.com/agiledragon/gomonkey/v2"

	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"

	"github.com/stretchr/testify/require"

	"go.uber.org/mock/gomock"

	"google.golang.org/protobuf/types/known/anypb"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config/xds/apiclient"
	xdsmodel "github.com/apache/dubbo-go-pixiu/pkg/config/xds/model"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server/controls"
	"github.com/apache/dubbo-go-pixiu/pkg/server/controls/mocks"
)

type replaceXDSClusterManager struct {
	controls.ClusterManager
	replace func([]*model.ClusterConfig) error
}

type legacyClusterManager struct {
	controls.ClusterManager
}

func (m *replaceXDSClusterManager) ReplaceXDSClusters(clusters []*model.ClusterConfig) error {
	return m.replace(clusters)
}

func makeClusters() *xdsmodel.PixiuExtensionClusters {
	return &xdsmodel.PixiuExtensionClusters{
		Clusters: []*xdsmodel.Cluster{
			{
				Name:    "http-baidu",
				TypeStr: "Static",
				Endpoints: []*xdsmodel.Endpoint{{
					Id: "backend",
					Address: &xdsmodel.SocketAddress{
						Address: "httpbin.org",
						Port:    80,
					},
				}},
			},
		},
	}
}

func getCdsConfig() *core.TypedExtensionConfig {
	makeClusters := func() *xdsmodel.PixiuExtensionClusters {
		return &xdsmodel.PixiuExtensionClusters{
			Clusters: []*xdsmodel.Cluster{
				{
					Name:    "http-baidu",
					TypeStr: "Static",
					Endpoints: []*xdsmodel.Endpoint{{
						Id: "backend",
						Address: &xdsmodel.SocketAddress{
							Address: "httpbin.org",
							Port:    80,
						},
					}},
				},
			},
		}
	}
	cdsResource, _ := anypb.New(makeClusters())
	return &core.TypedExtensionConfig{
		Name:        constant.ClusterType,
		TypedConfig: cdsResource,
	}
}

func TestCdsManager_Fetch(t *testing.T) {
	var fetchResult []*apiclient.ProtoAny
	var fetchError error
	var cluster = map[string]struct{}{}
	var updateCluster *model.ClusterConfig
	var addCluster *model.ClusterConfig
	xdsConfig := getCdsConfig()

	ctrl := gomock.NewController(t)
	clusterMg := mocks.NewMockClusterManager(ctrl)
	//var deltaResult chan *apiclient.DeltaResources
	//var deltaErr error
	patches := gomonkey.ApplyMethod(&apiclient.GrpcExtensionApiClient{}, "Fetch", func(_ *apiclient.GrpcExtensionApiClient, localVersion string) ([]*apiclient.ProtoAny, error) {
		return fetchResult, fetchError
	})
	defer patches.Reset()

	clusterMg.EXPECT().HasCluster(gomock.Any()).AnyTimes().DoAndReturn(func(clusterName string) bool {
		_, ok := cluster[clusterName]
		return ok
	})

	clusterMg.EXPECT().UpdateCluster(gomock.Any()).AnyTimes().Do(func(new *model.ClusterConfig) {
		updateCluster = new
	})
	clusterMg.EXPECT().AddCluster(gomock.Any()).AnyTimes().Do(func(c *model.ClusterConfig) {
		addCluster = c
	})
	clusterMg.EXPECT().RemoveCluster(gomock.Any()).AnyTimes()
	replacingClusterMg := &replaceXDSClusterManager{ClusterManager: clusterMg, replace: func(clusters []*model.ClusterConfig) error {
		if len(clusters) > 0 {
			addCluster = clusters[0]
		}
		return nil
	}}
	clusterMg.EXPECT().CloneXdsControlStore().AnyTimes().DoAndReturn(func() (controls.ClusterStore, error) {
		store := mocks.NewMockClusterStore(ctrl)
		store.EXPECT().Config().AnyTimes()
		return store, nil
	})
	//supermonkey.Patch((*server.ClusterManager).HasCluster, func(_ *server.ClusterManager, clusterName string) bool {
	//	_, ok := cluster[clusterName]
	//	return ok
	//})
	//supermonkey.Patch((*server.ClusterManager).UpdateCluster, func(_ *server.ClusterManager, new *model.Cluster) {
	//	updateCluster = new
	//})
	//supermonkey.Patch((*server.ClusterManager).AddCluster, func(_ *server.ClusterManager, c *model.Cluster) {
	//	addCluster = c
	//})
	//supermonkey.Patch((*server.ClusterManager).RemoveCluster, func(_ *server.ClusterManager, names []string) {
	//	//do nothing.
	//})
	//supermonkey.Patch((*server.ClusterManager).CloneStore, func(_ *server.ClusterManager) (*server.ClusterStore, error) {
	//	return &server.ClusterStore{}, nil
	//})
	//supermonkey.Patch((*apiclient.GrpcExtensionApiClient).Delta, func(_ *apiclient.GrpcExtensionApiClient) (chan *apiclient.DeltaResources, error) {
	//	return deltaResult, deltaErr
	//})

	tests := []struct {
		name              string
		mockResult        []*apiclient.ProtoAny
		mockError         error
		wantErr           bool
		wantNewCluster    bool
		wantUpdateCluster bool
	}{
		{"error", nil, stderr.New("error test"), true, false, false},
		{"simple", nil, nil, false, false, false},
		{"withValue", func() []*apiclient.ProtoAny {
			return []*apiclient.ProtoAny{
				apiclient.NewProtoAny(xdsConfig),
			}
		}(), nil, false, true, false},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CdsManager{
				DiscoverApi: &apiclient.GrpcExtensionApiClient{},
				clusterMg:   replacingClusterMg,
			}
			//reset context value.
			fetchError = tt.mockError
			fetchResult = tt.mockResult
			updateCluster = nil
			addCluster = nil

			err := c.Fetch()
			assert := require.New(t)
			if tt.wantErr {
				assert.Error(err)
				return
			}
			assert.NoError(err)
			if tt.wantUpdateCluster {
				assert.NotNil(updateCluster)
			} else {
				assert.Nil(updateCluster)
			}
			if tt.wantNewCluster {
				assert.NotNil(addCluster)
			} else {
				assert.Nil(addCluster)
			}
		})
	}
}

func TestCdsManager_makeCluster(t *testing.T) {
	c := &CdsManager{}
	cluster := makeClusters().Clusters[0]
	modelCluster, err := c.makeCluster(cluster)
	assert := require.New(t)
	assert.NoError(err)
	assert.Equal(cluster.Name, modelCluster.Name)
	assert.Equal(cluster.TypeStr, modelCluster.TypeStr)
	assert.Equal(cluster.Endpoints[0].Name, modelCluster.Endpoints[0].Name)
	assert.Equal(cluster.Endpoints[0].Address.Address, modelCluster.Endpoints[0].Address.Address)
	assert.Equal(cluster.Endpoints[0].Address.Port, int64(modelCluster.Endpoints[0].Address.Port))
}

func TestCdsManager_MakeEndpointsPreservesEDSHealth(t *testing.T) {
	manager := &CdsManager{}
	endpoints, err := manager.makeEndpoints([]*xdsmodel.Endpoint{{
		Id:       "endpoint-1",
		Address:  &xdsmodel.SocketAddress{Address: "10.0.0.1", Port: 20880},
		Metadata: map[string]string{endpointHealthMetadataKey: "true"},
	}})

	require.NoError(t, err)
	require.Len(t, endpoints, 1)
	require.True(t, endpoints[0].UnHealthy)
}

func TestCdsManagerRejectsInvalidExtensionClustersBeforePublication(t *testing.T) {
	tests := []struct {
		name    string
		cluster *xdsmodel.Cluster
		want    string
	}{
		{
			name: "unknown discovery type",
			cluster: &xdsmodel.Cluster{
				Name: "orders", TypeStr: "made-up",
			},
			want: "unsupported discovery type",
		},
		{
			name: "unknown load balancer",
			cluster: &xdsmodel.Cluster{
				Name: "orders", TypeStr: "Static", LbStr: "LeastRequest",
			},
			want: "unsupported load-balancing policy",
		},
		{
			name: "port above uint16",
			cluster: &xdsmodel.Cluster{
				Name: "orders", TypeStr: "Static",
				Endpoints: []*xdsmodel.Endpoint{{Id: "orders-1", Address: &xdsmodel.SocketAddress{Address: "127.0.0.1", Port: 65536}}},
			},
			want: "invalid socket address",
		},
		{
			name: "duplicate endpoint ID",
			cluster: &xdsmodel.Cluster{
				Name: "orders", TypeStr: "Static",
				Endpoints: []*xdsmodel.Endpoint{
					{Id: "duplicate", Address: &xdsmodel.SocketAddress{Address: "127.0.0.1", Port: 20880}},
					{Id: "duplicate", Address: &xdsmodel.SocketAddress{Address: "127.0.0.2", Port: 20880}},
				},
			},
			want: "duplicate endpoint ID",
		},
		{
			name: "unknown EDS API type",
			cluster: &xdsmodel.Cluster{
				Name: "orders", TypeStr: "EDS",
				EdsClusterConfig: &xdsmodel.EdsClusterConfig{EdsConfig: &xdsmodel.ConfigSource{
					ApiConfigSource: &xdsmodel.ApiConfigSource{APITypeStr: "made-up"},
				}},
			},
			want: "unsupported EDS API type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ctrl := gomock.NewController(t)
			clusterMg := mocks.NewMockClusterManager(ctrl)
			manager := &CdsManager{clusterMg: clusterMg}

			err := manager.setupCluster([]*xdsmodel.Cluster{tt.cluster})
			require.ErrorContains(t, err, tt.want)
		})
	}
}

func TestCdsManagerDefaultsOmittedLoadBalancerToRand(t *testing.T) {
	manager := &CdsManager{}
	cluster, err := manager.makeCluster(&xdsmodel.Cluster{Name: "orders", TypeStr: "Static"})
	require.NoError(t, err)
	require.Equal(t, model.LoadBalancerRand, cluster.LbStr)
}

func TestCdsManagerRejectsManagerWithoutTransactionalReplacement(t *testing.T) {
	manager := &CdsManager{clusterMg: &legacyClusterManager{}}
	err := manager.setupCluster([]*xdsmodel.Cluster{{
		Name:    "orders",
		TypeStr: "Static",
		Endpoints: []*xdsmodel.Endpoint{{
			Id:      "orders-1",
			Address: &xdsmodel.SocketAddress{Address: "127.0.0.1", Port: 20880},
		}},
	}})

	require.ErrorContains(t, err, "does not support transactional xDS replacement")
}

func TestCdsManager_StandardEDSClusterDoesNotRequireNestedConfigSource(t *testing.T) {
	ctrl := gomock.NewController(t)
	clusterMg := mocks.NewMockClusterManager(ctrl)
	replacingClusterMg := &replaceXDSClusterManager{ClusterManager: clusterMg, replace: func(clusters []*model.ClusterConfig) error {
		require.Len(t, clusters, 1)
		require.Equal(t, "orders", clusters[0].Name)
		require.Equal(t, "orders-eds", clusters[0].EdsClusterConfig.ServiceName)
		require.Equal(t, "10.0.0.1", clusters[0].Endpoints[0].Address.Address)
		return nil
	}}
	manager := &CdsManager{clusterMg: replacingClusterMg}

	require.NotPanics(t, func() {
		require.NoError(t, manager.setupCluster([]*xdsmodel.Cluster{{
			Name:             "orders",
			TypeStr:          "Static",
			LbStr:            "RoundRobin",
			EdsClusterConfig: &xdsmodel.EdsClusterConfig{ServiceName: "orders-eds"},
			Endpoints: []*xdsmodel.Endpoint{{
				Id:      "10.0.0.1:20880",
				Address: &xdsmodel.SocketAddress{Address: "10.0.0.1", Port: 20880},
			}},
		}}))
	})
}

func TestCdsManager_ApplyDelta(t *testing.T) {
	ctrl := gomock.NewController(t)
	clusterMg := mocks.NewMockClusterManager(ctrl)
	owned := map[string]*model.ClusterConfig{
		"old-dynamic": {Name: "old-dynamic"},
	}
	replacingClusterMg := &replaceXDSClusterManager{ClusterManager: clusterMg, replace: func(clusters []*model.ClusterConfig) error {
		clear(owned)
		for _, cluster := range clusters {
			owned[cluster.Name] = cluster
		}
		return nil
	}}
	manager := &CdsManager{clusterMg: replacingClusterMg}

	// An empty delta response is a no-op. Omission is not deletion in Delta xDS.
	require.NoError(t, manager.applyDelta(&apiclient.DeltaResources{}))
	require.Contains(t, owned, "old-dynamic")

	// The extension resource is an aggregate, so a new payload replaces the
	// xDS-owned aggregate while leaving static resources untouched.
	require.NoError(t, manager.applyDelta(&apiclient.DeltaResources{
		NewResources: []*apiclient.ProtoAny{apiclient.NewProtoAny(getCdsConfig())},
	}))
	require.NotContains(t, owned, "old-dynamic")
	require.Contains(t, owned, "http-baidu")

	// An explicitly delivered empty aggregate is authoritative and therefore
	// removes the clusters previously held by that aggregate.
	emptyClusters, err := anypb.New(&xdsmodel.PixiuExtensionClusters{})
	require.NoError(t, err)
	require.NoError(t, manager.applyDelta(&apiclient.DeltaResources{
		NewResources: []*apiclient.ProtoAny{apiclient.NewProtoAny(&core.TypedExtensionConfig{
			Name:        constant.ClusterType,
			TypedConfig: emptyClusters,
		})},
	}))
	require.NotContains(t, owned, "http-baidu")

	require.NoError(t, manager.applyDelta(&apiclient.DeltaResources{
		NewResources: []*apiclient.ProtoAny{apiclient.NewProtoAny(getCdsConfig())},
	}))
	require.Contains(t, owned, "http-baidu")

	// Removing the subscribed extension resource clears only xDS-owned state.
	require.NoError(t, manager.applyDelta(&apiclient.DeltaResources{
		RemovedResources: []string{constant.ClusterType},
	}))
	require.NotContains(t, owned, "http-baidu")
}

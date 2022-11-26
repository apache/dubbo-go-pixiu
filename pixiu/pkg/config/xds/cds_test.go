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
	"github.com/cch123/supermonkey"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/xds"
	pixiupb "github.com/dubbo-go-pixiu/pixiu-api/pkg/xds/model"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/config/xds/apiclient"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server/controls"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server/controls/mocks"
)

func makeClusters() *pixiupb.PixiuExtensionClusters {
	return &pixiupb.PixiuExtensionClusters{
		Clusters: []*pixiupb.Cluster{
			{
				Name:    "http-baidu",
				TypeStr: "http",
				Endpoints: []*pixiupb.Endpoint{{
					Id: "backend",
					Address: &pixiupb.SocketAddress{
						Address: "httpbin.org",
						Port:    80,
					},
				}},
			},
		},
	}
}

func getCdsConfig() *core.TypedExtensionConfig {
	makeClusters := func() *pixiupb.PixiuExtensionClusters {
		return &pixiupb.PixiuExtensionClusters{
			Clusters: []*pixiupb.Cluster{
				{
					Name:    "http-baidu",
					TypeStr: "http",
					Endpoints: []*pixiupb.Endpoint{{
						Id: "backend",
						Address: &pixiupb.SocketAddress{
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
		Name:        xds.ClusterType,
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
	supermonkey.Patch((*apiclient.GrpcExtensionApiClient).Fetch, func(_ *apiclient.GrpcExtensionApiClient, localVersion string) ([]*apiclient.ProtoAny, error) {
		return fetchResult, fetchError
	})
	//supermonkey.Patch(server.GetClusterManager, func() *server.ClusterManager {
	//	return nil
	//})
	clusterMg.EXPECT().HasCluster(gomock.Any()).DoAndReturn(func(clusterName string) bool {
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
	defer supermonkey.UnpatchAll()

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
				clusterMg:   clusterMg,
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
	modelCluster := c.makeCluster(cluster)
	assert := require.New(t)
	assert.Equal(cluster.Name, modelCluster.Name)
	assert.Equal(cluster.TypeStr, modelCluster.TypeStr)
	assert.Equal(cluster.Endpoints[0].Name, modelCluster.Endpoints[0].Name)
	assert.Equal(cluster.Endpoints[0].Address.Address, modelCluster.Endpoints[0].Address.Address)
	assert.Equal(cluster.Endpoints[0].Address.Port, int64(modelCluster.Endpoints[0].Address.Port))
}

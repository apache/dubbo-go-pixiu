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
	"testing"
)

import (
	"github.com/agiledragon/gomonkey/v2"

	"github.com/golang/mock/gomock"

	"github.com/stretchr/testify/require"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config/xds/apiclient"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server/controls"
	"github.com/apache/dubbo-go-pixiu/pkg/server/controls/mocks"
)

func TestAdapter_createApiManager(t *testing.T) {
	cluster := &model.ClusterConfig{
		Name:    "cluster-1",
		TypeStr: "GRPC",
		Endpoints: []*model.Endpoint{
			{
				Address: model.SocketAddress{
					Address: "localhost",
					Port:    18000,
				},
			},
		},
	}

	node := model.Node{
		Cluster: "test-cluster",
		Id:      "node-test-1",
	}

	apiConfig := model.ApiConfigSource{
		APIType:     model.ApiTypeGRPC,
		APITypeStr:  "GRPC",
		ClusterName: []string{"cluster-1"},
	}

	var state = connectivity.Ready
	gconn := &grpc.ClientConn{}
	patches := gomonkey.ApplyFunc(grpc.DialContext, func(ctx context.Context, target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
		fmt.Println("***** DialContext")
		return gconn, nil
	})
	defer patches.Reset()

	patches.ApplyMethod(&grpc.ClientConn{}, "Close", func(_ *grpc.ClientConn) error {
		return nil
	})
	patches.ApplyMethod(&grpc.ClientConn{}, "GetState", func(_ *grpc.ClientConn) connectivity.State {
		return state
	})

	patches.ApplyMethod(&apiclient.GrpcExtensionApiClient{}, "Fetch", func(_ *apiclient.GrpcExtensionApiClient, localVersion string) ([]*apiclient.ProtoAny, error) {
		return nil, nil
	})
	patches.ApplyMethod(&apiclient.GrpcExtensionApiClient{}, "Delta", func(_ *apiclient.GrpcExtensionApiClient) (chan *apiclient.DeltaResources, error) {
		ch := make(chan *apiclient.DeltaResources)
		close(ch)
		return ch, nil
	})

	//init cluster manager
	ctrl := gomock.NewController(t)
	clusterMg := mocks.NewMockClusterManager(ctrl)
	{
		clusterMg.EXPECT().
			HasCluster("cluster-1").AnyTimes().
			Return(true)
		clusterMg.EXPECT().CloneXdsControlStore().DoAndReturn(func() (controls.ClusterStore, error) {
			store := mocks.NewMockClusterStore(ctrl)
			store.EXPECT().Config().Return([]*model.ClusterConfig{cluster})
			return store, nil
		})
		// delete this stub by #https://github.com/golang/mock/issues/530
		//clusterMg.EXPECT().HasCluster("cluster-2").Return(false)
		apiclient.Init(clusterMg)
	}

	ada := Xds{
		clusterMg: clusterMg,
	}
	ada.Start()
	api := ada.createApiManager(&apiConfig, &node, constant.ClusterType)
	assert := require.New(t)
	assert.NotNil(api)
}

func TestAdapter_createApiManager_ErrorHandling(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clusterMg := mocks.NewMockClusterManager(ctrl)
	apiclient.Init(clusterMg)

	node := &model.Node{
		Cluster: "test-cluster",
		Id:      "node-test-1",
	}

	tests := []struct {
		name        string
		apiConfig   *model.ApiConfigSource
		expectNil   bool
		description string
	}{
		{
			name:        "nil config",
			apiConfig:   nil,
			expectNil:   true,
			description: "should return nil for nil config",
		},
		{
			name: "empty cluster name - GRPC type",
			apiConfig: &model.ApiConfigSource{
				APIType:     model.ApiTypeGRPC,
				APITypeStr:  "GRPC",
				ClusterName: []string{}, // Empty
			},
			expectNil:   true,
			description: "should return nil for empty cluster name in GRPC type",
		},
		{
			name: "unsupported API type",
			apiConfig: &model.ApiConfigSource{
				APIType:     model.ApiType(-1), // Invalid type
				APITypeStr:  "INVALID",
				ClusterName: []string{"cluster-1"},
			},
			expectNil:   true,
			description: "should return nil for unsupported API type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ada := &Xds{
				clusterMg: clusterMg,
				exitCh:    make(chan struct{}),
			}

			api := ada.createApiManager(tt.apiConfig, node, constant.ClusterType)
			assert := require.New(t)
			if tt.expectNil {
				assert.Nil(api, tt.description)
			} else {
				assert.NotNil(api, tt.description)
			}
		})
	}
}

func TestAdapter_Start_NilApiManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clusterMg := mocks.NewMockClusterManager(ctrl)
	drm := mocks.NewMockDynamicResourceManager(ctrl)

	apiclient.Init(clusterMg)

	// Setup mocks to trigger error scenarios
	// Start() will call GetLds() first at line 128
	drm.EXPECT().GetLds().Return(&model.ApiConfigSource{
		APIType:     model.ApiTypeGRPC,
		ClusterName: []string{}, // Empty - will cause createApiManager to return nil
	})
	// Then it calls GetLds() again at line 129 (inside createApiManager call)
	drm.EXPECT().GetLds().Return(&model.ApiConfigSource{
		APIType:     model.ApiTypeGRPC,
		ClusterName: []string{},
	})
	// Then it calls GetNode() at line 129
	drm.EXPECT().GetNode().Return(&model.Node{})
	// After discoverApi == nil, it returns early and doesn't call GetCds()

	ada := &Xds{
		clusterMg:         clusterMg,
		dynamicResourceMg: drm,
		exitCh:            make(chan struct{}),
	}

	// Start should handle nil DiscoverApi gracefully without panic
	ada.Start()

	// Verify that lds was not created because createApiManager returned nil
	assert := require.New(t)
	assert.Nil(ada.lds)
	assert.Nil(ada.cds) // Should also be nil because Start() returned early
}

func TestAdapter_Start_CdsNilApiManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clusterMg := mocks.NewMockClusterManager(ctrl)
	drm := mocks.NewMockDynamicResourceManager(ctrl)

	apiclient.Init(clusterMg)

	// Setup LDS to skip (nil config)
	drm.EXPECT().GetLds().Return(nil)

	// Setup CDS to fail
	// Start() will call GetCds() at line 143
	drm.EXPECT().GetCds().Return(&model.ApiConfigSource{
		APIType:     model.ApiTypeGRPC,
		ClusterName: []string{}, // Empty - will cause createApiManager to return nil
	})
	// Then it calls GetCds() again at line 144 (inside createApiManager call)
	drm.EXPECT().GetCds().Return(&model.ApiConfigSource{
		APIType:     model.ApiTypeGRPC,
		ClusterName: []string{},
	})
	// Then it calls GetNode() at line 144
	drm.EXPECT().GetNode().Return(&model.Node{})

	ada := &Xds{
		clusterMg:         clusterMg,
		dynamicResourceMg: drm,
		exitCh:            make(chan struct{}),
	}

	// Start should handle nil DiscoverApi gracefully without panic
	ada.Start()

	// Verify that cds was not created because createApiManager returned nil
	assert := require.New(t)
	assert.Nil(ada.lds) // Should be nil because GetLds() returned nil
	assert.Nil(ada.cds) // Should be nil because discoverApi was nil
}

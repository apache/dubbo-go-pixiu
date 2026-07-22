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
	require.NoError(t, ada.Start())
	api, err := ada.createApiManager(&apiConfig, &node, constant.ClusterType)
	assert := require.New(t)
	assert.NoError(err)
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
		expectErr   bool
		description string
	}{
		{
			name:        "nil config",
			apiConfig:   nil,
			expectErr:   false,
			description: "should return (nil, nil) for nil config (not configured is not an error)",
		},
		{
			name: "empty cluster name - GRPC type",
			apiConfig: &model.ApiConfigSource{
				APIType:     model.ApiTypeGRPC,
				APITypeStr:  "GRPC",
				ClusterName: []string{}, // Empty
			},
			expectErr:   true,
			description: "should return error for empty cluster name in GRPC type",
		},
		{
			name: "unsupported API type",
			apiConfig: &model.ApiConfigSource{
				APIType:     model.ApiType(-1), // Invalid type
				APITypeStr:  "INVALID",
				ClusterName: []string{"cluster-1"},
			},
			expectErr:   true,
			description: "should return error for unsupported API type",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			ada := &Xds{
				clusterMg: clusterMg,
				exitCh:    make(chan struct{}),
			}

			api, err := ada.createApiManager(tt.apiConfig, node, constant.ClusterType)
			assert := require.New(t)
			if tt.expectErr {
				assert.Error(err, tt.description)
				assert.Nil(api, tt.description)
			} else {
				assert.NoError(err, tt.description)
				assert.Nil(api, tt.description)
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

	// LDS configured but with an empty cluster name -> createApiManager returns an error.
	// Start() calls GetLds() once for the nil-check, then again inside createApiManager.
	gomock.InOrder(
		drm.EXPECT().GetLds().Return(&model.ApiConfigSource{
			APIType:     model.ApiTypeGRPC,
			ClusterName: []string{}, // Empty - will cause createApiManager to return an error
		}),
		drm.EXPECT().GetLds().Return(&model.ApiConfigSource{
			APIType:     model.ApiTypeGRPC,
			ClusterName: []string{},
		}),
		drm.EXPECT().GetNode().Return(&model.Node{}),
	)
	// CDS is not configured, so its block is skipped after the nil-check.
	drm.EXPECT().GetCds().Return(nil)

	ada := &Xds{
		clusterMg:         clusterMg,
		dynamicResourceMg: drm,
		exitCh:            make(chan struct{}),
	}

	// Start must surface the LDS init failure instead of silently returning (fail-fast).
	err := ada.Start()
	assert := require.New(t)
	assert.Error(err)
	assert.Contains(err.Error(), "LDS")
	assert.Nil(ada.lds) // not created because createApiManager failed
	assert.Nil(ada.cds) // not configured
}

func TestAdapter_Start_CdsNilApiManager(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clusterMg := mocks.NewMockClusterManager(ctrl)
	drm := mocks.NewMockDynamicResourceManager(ctrl)

	apiclient.Init(clusterMg)

	// LDS not configured.
	drm.EXPECT().GetLds().Return(nil)

	// CDS configured but with an empty cluster name -> createApiManager returns an error.
	// Start() calls GetCds() once for the nil-check, then again inside createApiManager.
	gomock.InOrder(
		drm.EXPECT().GetCds().Return(&model.ApiConfigSource{
			APIType:     model.ApiTypeGRPC,
			ClusterName: []string{}, // Empty - will cause createApiManager to return an error
		}),
		drm.EXPECT().GetCds().Return(&model.ApiConfigSource{
			APIType:     model.ApiTypeGRPC,
			ClusterName: []string{},
		}),
		drm.EXPECT().GetNode().Return(&model.Node{}),
	)

	ada := &Xds{
		clusterMg:         clusterMg,
		dynamicResourceMg: drm,
		exitCh:            make(chan struct{}),
	}

	// Start must surface the CDS init failure (fail-fast).
	err := ada.Start()
	assert := require.New(t)
	assert.Error(err)
	assert.Contains(err.Error(), "CDS")
	assert.Nil(ada.lds) // not configured
	assert.Nil(ada.cds) // not created because createApiManager failed
}

// TestAdapter_Start_LdsAndCdsIndependentErrors verifies the core P0 fix: an LDS init failure
// must NOT skip CDS initialization. Both are attempted, and the returned error reports both.
func TestAdapter_Start_LdsAndCdsIndependentErrors(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clusterMg := mocks.NewMockClusterManager(ctrl)
	drm := mocks.NewMockDynamicResourceManager(ctrl)

	apiclient.Init(clusterMg)

	// LDS configured with an empty cluster name -> fails.
	gomock.InOrder(
		drm.EXPECT().GetLds().Return(&model.ApiConfigSource{
			APIType:     model.ApiTypeGRPC,
			ClusterName: []string{},
		}),
		drm.EXPECT().GetLds().Return(&model.ApiConfigSource{
			APIType:     model.ApiTypeGRPC,
			ClusterName: []string{},
		}),
		drm.EXPECT().GetNode().Return(&model.Node{}),
	)
	// CDS also configured with an empty cluster name -> fails too, proving it was attempted
	// despite the LDS failure.
	gomock.InOrder(
		drm.EXPECT().GetCds().Return(&model.ApiConfigSource{
			APIType:     model.ApiTypeGRPC,
			ClusterName: []string{},
		}),
		drm.EXPECT().GetCds().Return(&model.ApiConfigSource{
			APIType:     model.ApiTypeGRPC,
			ClusterName: []string{},
		}),
		drm.EXPECT().GetNode().Return(&model.Node{}),
	)

	ada := &Xds{
		clusterMg:         clusterMg,
		dynamicResourceMg: drm,
		exitCh:            make(chan struct{}),
	}

	err := ada.Start()
	assert := require.New(t)
	assert.Error(err)
	// Both LDS and CDS failures are present in the joined error.
	assert.Contains(err.Error(), "LDS")
	assert.Contains(err.Error(), "CDS")
	assert.Nil(ada.lds)
	assert.Nil(ada.cds)
}

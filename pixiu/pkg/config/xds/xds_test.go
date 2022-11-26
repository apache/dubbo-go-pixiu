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
	monkey "github.com/cch123/supermonkey"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/xds"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/config/xds/apiclient"
	_ "github.com/apache/dubbo-go-pixiu/pixiu/pkg/config/xds/apiclient"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server/controls"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server/controls/mocks"
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
	monkey.Patch(grpc.DialContext, func(ctx context.Context, target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
		fmt.Println("***** DialContext")
		return gconn, nil
	})
	monkey.Patch((*grpc.ClientConn).Close, func(_ *grpc.ClientConn) error {
		return nil
	})
	monkey.Patch((*grpc.ClientConn).GetState, func(_ *grpc.ClientConn) connectivity.State {
		return state
	})
	//monkey.Patch(server.GetDynamicResourceManager, func() server.DynamicResourceManager {
	//	return &server.DynamicResourceManagerImpl{}
	//})
	//monkey.Patch((*server.DynamicResourceManagerImpl).GetLds, func(_ *server.DynamicResourceManagerImpl) *model.ApiConfigSource {
	//	return &apiConfig
	//})
	//monkey.Patch((*server.DynamicResourceManagerImpl).GetCds, func(_ *server.DynamicResourceManagerImpl) *model.ApiConfigSource {
	//	return &apiConfig
	//})
	//monkey.Patch(server.GetClusterManager, func() *server.ClusterManager {
	//	return nil
	//})
	//monkey.Patch((*server.ClusterManager).CloneStore, func(_ *server.ClusterManager) (*server.ClusterStore, error) {
	//	return &server.ClusterStore{
	//		Config:  []*model.Cluster{cluster},
	//		Version: 1,
	//	}, nil
	//})

	monkey.Patch((*apiclient.GrpcExtensionApiClient).Fetch, func(_ *apiclient.GrpcExtensionApiClient, localVersion string) ([]*apiclient.ProtoAny, error) {
		return nil, nil
	})
	monkey.Patch((*apiclient.GrpcExtensionApiClient).Delta, func(_ *apiclient.GrpcExtensionApiClient) (chan *apiclient.DeltaResources, error) {
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

	defer monkey.UnpatchAll()

	ada := Xds{
		clusterMg: clusterMg,
	}
	ada.Start()
	api := ada.createApiManager(&apiConfig, &node, xds.ClusterType)
	assert := require.New(t)
	assert.NotNil(api)
}

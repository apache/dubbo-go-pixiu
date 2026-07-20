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

package apiclient

import (
	"context"
	"sync"
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
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server/controls"
	"github.com/apache/dubbo-go-pixiu/pkg/server/controls/mocks"
)

func TestGRPCClusterManager_GetGrpcCluster(t *testing.T) {
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
	type fields struct {
		clusters  *sync.Map
		clusterMg controls.ClusterManager
	}
	type args struct {
		name string
	}
	ctrl := gomock.NewController(t)
	clusterMg := mocks.NewMockClusterManager(ctrl)
	clusterMg.EXPECT().
		HasCluster("cluster-1").AnyTimes().
		Return(true)
	clusterMg.EXPECT().CloneXdsControlStore().DoAndReturn(func() (controls.ClusterStore, error) {
		store := mocks.NewMockClusterStore(ctrl)
		store.EXPECT().Config().Return([]*model.ClusterConfig{cluster})
		return store, nil
	})
	clusterMg.EXPECT().HasCluster("cluster-2").Return(false)
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{"test-simple", fields{
			clusters:  &sync.Map{},
			clusterMg: clusterMg,
		}, args{name: "cluster-1"}, nil},
		{"test-not-exist", fields{
			clusters:  &sync.Map{},
			clusterMg: clusterMg,
		}, args{name: "cluster-2"}, ErrClusterNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GRPCClusterManager{
				clusters:  tt.fields.clusters,
				clusterMg: tt.fields.clusterMg,
			}
			got, err := g.GetGrpcCluster(tt.args.name)
			assert := require.New(t)
			if tt.wantErr != nil {
				assert.ErrorIs(err, tt.wantErr)
				return
			}
			assert.NotNil(got)
			//run two times.
			got, err = g.GetGrpcCluster(tt.args.name)
			assert.NoError(err)
			assert.NotNil(got)

		})
	}
}

func TestGRPCCluster_GetConnect(t *testing.T) {
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
	g := GRPCCluster{
		name:   "name",
		config: cluster,
		once:   sync.Once{},
		conn:   nil,
	}

	gconn := &grpc.ClientConn{}
	var state = connectivity.Ready
	patches := gomonkey.ApplyFunc(grpc.DialContext, func(ctx context.Context, target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
		return gconn, nil
	})
	defer patches.Reset()
	patches.ApplyMethod(&grpc.ClientConn{}, "Close", func(_ *grpc.ClientConn) error {
		return nil
	})
	patches.ApplyMethod(&grpc.ClientConn{}, "GetState", func(_ *grpc.ClientConn) connectivity.State {
		return state
	})

	assert := require.New(t)
	conn, err := g.GetConnection()
	assert.NoError(err)
	assert.NotNil(conn)
	assert.NotNil(g.conn)
	assert.True(g.IsAlive())
	err = g.Close()
	assert.NoError(err)
	state = connectivity.Shutdown
	assert.False(g.IsAlive())
}

func TestCreateEnvoyGrpcApiClient_EmptyClusterName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test with empty cluster name - should return error without panic
	config := &model.ApiConfigSource{
		ClusterName: []string{}, // Empty cluster name
	}
	node := &model.Node{}
	exitCh := make(chan struct{})

	client, err := CreateEnvoyGrpcApiClient(config, node, exitCh, constant.ListenerType)
	assert := require.New(t)
	assert.Error(err)
	assert.Contains(err.Error(), "cluster name is required")
	assert.Nil(client)
}

func TestCreateEnvoyGrpcApiClient_ClusterNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clusterMg := mocks.NewMockClusterManager(ctrl)
	clusterMg.EXPECT().HasCluster("non-existent-cluster").Return(false)

	// Initialize the global grpc manager with mock
	Init(clusterMg)

	config := &model.ApiConfigSource{
		ClusterName: []string{"non-existent-cluster"},
	}
	node := &model.Node{}
	exitCh := make(chan struct{})

	client, err := CreateEnvoyGrpcApiClient(config, node, exitCh, constant.ClusterType)
	assert := require.New(t)
	assert.Error(err)
	assert.Contains(err.Error(), "get cluster for init error")
	assert.Nil(client)
}

func TestCreateGrpExtensionApiClient_EmptyClusterName(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	// Test with empty cluster name - should return error without panic
	config := &model.ApiConfigSource{
		ClusterName: []string{}, // Empty cluster name
	}
	node := &model.Node{}
	exitCh := make(chan struct{})

	client, err := CreateGrpExtensionApiClient(config, node, exitCh, constant.ListenerType)
	assert := require.New(t)
	assert.Error(err)
	assert.Contains(err.Error(), "cluster name is required")
	assert.Nil(client)
}

func TestCreateGrpExtensionApiClient_ClusterNotFound(t *testing.T) {
	ctrl := gomock.NewController(t)
	defer ctrl.Finish()

	clusterMg := mocks.NewMockClusterManager(ctrl)
	clusterMg.EXPECT().HasCluster("non-existent-cluster").Return(false)

	// Initialize the global grpc manager with mock
	Init(clusterMg)

	config := &model.ApiConfigSource{
		ClusterName: []string{"non-existent-cluster"},
	}
	node := &model.Node{}
	exitCh := make(chan struct{})

	client, err := CreateGrpExtensionApiClient(config, node, exitCh, constant.ClusterType)
	assert := require.New(t)
	assert.Error(err)
	assert.Contains(err.Error(), "get cluster for init error")
	assert.Nil(client)
}

// TestGRPCCluster_GetConnection_PersistInitError verifies that when LDS and CDS
// share an xDS cluster and the first connection fails, the second call returns
// the same error instead of (nil, nil), preventing panic from calling methods
// on a nil ClientConn.
func TestGRPCCluster_GetConnection_PersistInitError(t *testing.T) {
	cluster := &model.ClusterConfig{
		Name:    "cluster-unreachable",
		TypeStr: "GRPC",
		Endpoints: []*model.Endpoint{
			{
				Address: model.SocketAddress{
					Address: "localhost",
					Port:    19999, // Unreachable port
				},
			},
		},
	}

	g := GRPCCluster{
		name:   "cluster-unreachable",
		config: cluster,
		once:   sync.Once{},
		conn:   nil,
	}

	assert := require.New(t)

	// First call should fail to connect
	conn1, err1 := g.GetConnection()
	assert.Error(err1)
	assert.Contains(err1.Error(), "failed")
	assert.Nil(conn1)

	// Second call should return the same error, not (nil, nil)
	conn2, err2 := g.GetConnection()
	assert.Error(err2)
	assert.Equal(err1, err2, "second call should return the same error as first")
	assert.Nil(conn2)

	// Verify that initErr was persisted
	assert.Equal(g.initErr, err1)
}

// TestGRPCCluster_GetConnection_NoEndpoints verifies that GetConnection returns
// an error when cluster has no endpoints configured, without panic.
func TestGRPCCluster_GetConnection_NoEndpoints(t *testing.T) {
	cluster := &model.ClusterConfig{
		Name:      "cluster-no-endpoints",
		TypeStr:   "GRPC",
		Endpoints: nil, // No endpoints
	}

	g := GRPCCluster{
		name:   "cluster-no-endpoints",
		config: cluster,
		once:   sync.Once{},
		conn:   nil,
	}

	assert := require.New(t)

	conn, err := g.GetConnection()
	assert.Error(err)
	assert.Contains(err.Error(), "no endpoints configured")
	assert.Nil(conn)

	// Second call should return the same error
	conn2, err2 := g.GetConnection()
	assert.Error(err2)
	assert.Equal(err, err2)
	assert.Nil(conn2)
}

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
	"github.com/cch123/supermonkey"
	"github.com/golang/mock/gomock"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server/controls"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server/controls/mocks"
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
	supermonkey.Patch(grpc.DialContext, func(ctx context.Context, target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
		return gconn, nil
	})
	supermonkey.Patch((*grpc.ClientConn).Close, func(_ *grpc.ClientConn) error {
		return nil
	})
	supermonkey.Patch((*grpc.ClientConn).GetState, func(_ *grpc.ClientConn) connectivity.State {
		return state
	})

	defer supermonkey.UnpatchAll()
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

package xds

import (
	"context"
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/xds/apiclient"
	_ "github.com/apache/dubbo-go-pixiu/pkg/adapter/xds/apiclient"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
	monkey "github.com/cch123/supermonkey"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/xds"
	"github.com/stretchr/testify/require"
	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"sync"
	"testing"
)

func TestGrpcClusterManager_GetGrpcCluster(t *testing.T) {
	cluster := &model.Cluster{
		Name:    "cluster-1",
		TypeStr: "GRPC",
		Endpoints: []*model.Endpoint{
			{
				Address: model.SocketAddress{
					ProtocolStr: "http",
					Address:     "localhost",
					Port:        18000,
				},
			},
		},
	}
	type fields struct {
		clusters *sync.Map
		store    *server.ClusterStore
	}
	type args struct {
		name string
	}
	tests := []struct {
		name    string
		fields  fields
		args    args
		wantErr error
	}{
		{"test-simple", fields{
			clusters: &sync.Map{},
			store: &server.ClusterStore{
				Config:  []*model.Cluster{cluster},
				Version: 1,
			},
		}, args{name: "cluster-1"}, nil},
		{"test-not-exist", fields{
			clusters: &sync.Map{},
			store: &server.ClusterStore{
				Config:  []*model.Cluster{cluster},
				Version: 1,
			},
		}, args{name: "cluster-2"}, ErrClusterNotFound},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			g := &GrpcClusterManager{
				clusters: tt.fields.clusters,
				store:    tt.fields.store,
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
			assert.NotNil(got)

		})
	}
}

func TestGrpcCluster_GetConnect(t *testing.T) {
	cluster := &model.Cluster{
		Name:    "cluster-1",
		TypeStr: "GRPC",
		Endpoints: []*model.Endpoint{
			{
				Address: model.SocketAddress{
					ProtocolStr: "http",
					Address:     "localhost",
					Port:        18000,
				},
			},
		},
	}
	g := GrpcCluster{
		name:   "name",
		config: cluster,
		once:   sync.Once{},
		conn:   nil,
	}

	gconn := &grpc.ClientConn{}
	var state connectivity.State = connectivity.Ready
	monkey.Patch(grpc.DialContext, func(ctx context.Context, target string, opts ...grpc.DialOption) (conn *grpc.ClientConn, err error) {
		return gconn, nil
	})
	monkey.Patch((*grpc.ClientConn).Close, func(_ *grpc.ClientConn) error {
		return nil
	})
	monkey.Patch((*grpc.ClientConn).GetState, func(_ *grpc.ClientConn) connectivity.State {
		return state
	})

	defer monkey.UnpatchAll()
	conn := g.GetConnect()
	assert := require.New(t)
	assert.NotNil(conn)
	assert.NotNil(g.conn)
	assert.True(g.IsAlive())
	err := g.Close()
	assert.NoError(err)
	state = connectivity.Shutdown
	assert.False(g.IsAlive())
}

func TestAdapter_createApiManager(t *testing.T) {
	cluster := &model.Cluster{
		Name:    "cluster-1",
		TypeStr: "GRPC",
		Endpoints: []*model.Endpoint{
			{
				Address: model.SocketAddress{
					ProtocolStr: "http",
					Address:     "localhost",
					Port:        18000,
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
	monkey.Patch(server.GetDynamicResourceManager, func() server.DynamicResourceManager {
		return &server.DynamicResourceManagerImpl{}
	})
	monkey.Patch((*server.DynamicResourceManagerImpl).GetLds, func(_ *server.DynamicResourceManagerImpl) *model.ApiConfigSource {
		return &apiConfig
	})
	monkey.Patch((*server.DynamicResourceManagerImpl).GetCds, func(_ *server.DynamicResourceManagerImpl) *model.ApiConfigSource {
		return &apiConfig
	})
	monkey.Patch(server.GetClusterManager, func() *server.ClusterManager {
		return nil
	})
	monkey.Patch((*server.ClusterManager).CloneStore, func(_ *server.ClusterManager) (*server.ClusterStore, error) {
		return &server.ClusterStore{
			Config:  []*model.Cluster{cluster},
			Version: 1,
		}, nil
	})

	monkey.Patch((*apiclient.GrpcApiClient).Fetch, func(_ *apiclient.GrpcApiClient, localVersion string) ([]*apiclient.ProtoAny, error) {
		return nil, nil
	})
	monkey.Patch((*apiclient.GrpcApiClient).Delta, func(_ *apiclient.GrpcApiClient) (chan *apiclient.DeltaResources, error) {
		ch := make(chan *apiclient.DeltaResources)
		close(ch)
		return ch, nil
	})

	defer monkey.UnpatchAll()

	ada := Adapter{}
	ada.Start()
	api := ada.createApiManager(&apiConfig, &node, xds.ClusterType)
	assert := require.New(t)
	assert.NotNil(api)
}

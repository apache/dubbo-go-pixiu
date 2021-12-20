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
	"github.com/envoyproxy/go-control-plane/envoy/service/discovery/v3"
	extensionpb "github.com/envoyproxy/go-control-plane/envoy/service/extension/v3"
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

type xdsClient struct {
}

func (x *xdsClient) StreamExtensionConfigs(ctx context.Context, opts ...grpc.CallOption) (extensionpb.ExtensionConfigDiscoveryService_StreamExtensionConfigsClient, error) {
	panic("implement me")
}

func (x *xdsClient) DeltaExtensionConfigs(ctx context.Context, opts ...grpc.CallOption) (extensionpb.ExtensionConfigDiscoveryService_DeltaExtensionConfigsClient, error) {
	panic("implement me")
}

func (x *xdsClient) FetchExtensionConfigs(ctx context.Context, in *envoy_service_discovery_v3.DiscoveryRequest, opts ...grpc.CallOption) (*envoy_service_discovery_v3.DiscoveryResponse, error) {
	panic("implement me")
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

	grpcClusterManager = GrpcClusterManager{
		clusters: &sync.Map{},
		store: &server.ClusterStore{
			Config:  []*model.Cluster{cluster},
			Version: 1,
		},
	}

	defer func() {
		grpcClusterManager = GrpcClusterManager{}
	}()

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

	var xdsClient extensionpb.ExtensionConfigDiscoveryServiceClient = &xdsClient{}
	monkey.Patch((*apiclient.GrpcApiClient).MakeClient, func(_ *apiclient.GrpcApiClient, cluster apiclient.GrpcCluster) extensionpb.ExtensionConfigDiscoveryServiceClient {
		fmt.Println("========================= get ExtensionConfigDiscoveryServiceClient")
		return xdsClient
	})

	//defer monkey.UnpatchAll()

	ada := Adapter{}
	api := ada.createApiManager(&apiConfig, &node, xds.ClusterType)
	assert := require.New(t)
	assert.NotNil(api)
	r, err := api.Fetch("")
	assert.NoError(err)
	assert.NotNil(r)

}

//func TestGrpcCluster_GetConnect(t *testing.T) {
//	type fields struct {
//		name   string
//		config *model.Cluster
//		once   sync.Once
//		conn   *grpc.ClientConn
//	}
//	tests := []struct {
//		name   string
//		fields fields
//		want   *grpc.ClientConn
//	}{
//		// TODO: Add test cases.
//	}
//	for _, tt := range tests {
//		t.Run(tt.name, func(t *testing.T) {
//			g := &GrpcCluster{
//				name:   tt.fields.name,
//				config: tt.fields.config,
//				once:   tt.fields.once,
//				conn:   tt.fields.conn,
//			}
//			if got := g.GetConnect(); !reflect.DeepEqual(got, tt.want) {
//				t.Errorf("GetConnect() = %v, want %v", got, tt.want)
//			}
//		})
//	}
//}

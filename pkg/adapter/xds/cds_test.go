package xds

import (
	"errors"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/xds/apiclient"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
	"github.com/cch123/supermonkey"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/xds"
	pixiupb "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/xds/model"
	core "github.com/envoyproxy/go-control-plane/envoy/config/core/v3"
	"github.com/golang/protobuf/proto"
	"github.com/stretchr/testify/require"
	"google.golang.org/protobuf/types/known/anypb"
	"testing"
)

func makeClusters() *pixiupb.PixiuExtensionClusters {
	return &pixiupb.PixiuExtensionClusters{
		Clusters: []*pixiupb.Cluster{
			{
				Name:    "http-baidu",
				TypeStr: "http",
				Endpoints: &pixiupb.Endpoint{
					Id: "backend",
					Address: &pixiupb.SocketAddress{
						ProtocolStr: "http",
						Address:     "httpbin.org",
						Port:        80,
					},
				},
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
					Endpoints: &pixiupb.Endpoint{
						Id: "backend",
						Address: &pixiupb.SocketAddress{
							ProtocolStr: "http",
							Address:     "httpbin.org",
							Port:        80,
						},
					},
				},
			},
		}
	}
	cdsResource, _ := anypb.New(proto.MessageV2(makeClusters()))
	return &core.TypedExtensionConfig{
		Name:        xds.ClusterType,
		TypedConfig: cdsResource,
	}
}

func TestCdsManager_Fetch(t *testing.T) {
	var fetchResult []*apiclient.ProtoAny
	var fetchError error
	var cluster = map[string]struct{}{}
	var updateCluster *model.Cluster
	var addCluster *model.Cluster
	xdsConfig := getCdsConfig()
	//var deltaResult chan *apiclient.DeltaResources
	//var deltaErr error
	supermonkey.Patch((*apiclient.GrpcApiClient).Fetch, func(_ *apiclient.GrpcApiClient, localVersion string) ([]*apiclient.ProtoAny, error) {
		return fetchResult, fetchError
	})
	supermonkey.Patch(server.GetClusterManager, func() *server.ClusterManager {
		return nil
	})
	supermonkey.Patch((*server.ClusterManager).HasCluster, func(_ *server.ClusterManager, clusterName string) bool {
		_, ok := cluster[clusterName]
		return ok
	})
	supermonkey.Patch((*server.ClusterManager).UpdateCluster, func(_ *server.ClusterManager, new *model.Cluster) {
		updateCluster = new
	})
	supermonkey.Patch((*server.ClusterManager).AddCluster, func(_ *server.ClusterManager, c *model.Cluster) {
		addCluster = c
	})
	//supermonkey.Patch((*apiclient.GrpcApiClient).Delta, func(_ *apiclient.GrpcApiClient) (chan *apiclient.DeltaResources, error) {
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
		{"error", nil, errors.New("error test"), true, false, false},
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
				DiscoverApi: &apiclient.GrpcApiClient{},
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
	assert.Equal(cluster.Endpoints.Name, modelCluster.Endpoints[0].Name)
	assert.Equal(cluster.Endpoints.Address.Address, modelCluster.Endpoints[0].Address.Address)
	assert.Equal(cluster.Endpoints.Address.Port, int64(modelCluster.Endpoints[0].Address.Port))
	assert.Equal(cluster.Endpoints.Address.ProtocolStr, modelCluster.Endpoints[0].Address.ProtocolStr)
}

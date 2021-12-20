package xds

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/xds/apiclient"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
	"github.com/cch123/supermonkey"
	"github.com/stretchr/testify/require"
	"testing"
)

func TestCdsManager_Fetch(t *testing.T) {
	var fetchResult []*apiclient.ProtoAny
	var fetchError error
	var cluster = map[string]struct{}{}
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
		return
	})
	supermonkey.Patch((*server.ClusterManager).AddCluster, func(_ *server.ClusterManager, c *model.Cluster) {
		return
	})
	//supermonkey.Patch((*apiclient.GrpcApiClient).Delta, func(_ *apiclient.GrpcApiClient) (chan *apiclient.DeltaResources, error) {
	//	return deltaResult, deltaErr
	//})
	defer supermonkey.UnpatchAll()

	tests := []struct {
		name    string
		wantErr bool
	}{
		// TODO: Add test cases.
		{},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			c := &CdsManager{
				DiscoverApi: &apiclient.GrpcApiClient{},
			}
			err := c.Fetch()
			assert := require.New(t)
			assert.NoError(err)
		})
	}
}

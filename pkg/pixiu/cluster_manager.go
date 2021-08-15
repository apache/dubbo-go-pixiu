package pixiu

import "github.com/apache/dubbo-go-pixiu/pkg/model"

type ClusterManager struct {
	cConfig []*model.Cluster
	caConfig []*model.ClusterAdapter

}

type ClusterAdapter interface {

}

func CreateDefaultClusterManager(bs *model.Bootstrap) *ClusterManager {
	return &ClusterManager{bs.StaticResources.Clusters, bs.StaticResources.ClusterAdapters}
}


func (* ClusterManager) startAdapter() {

}

func (* ClusterManager) addOrUpdateCluster() {

}

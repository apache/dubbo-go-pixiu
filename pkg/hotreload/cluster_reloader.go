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

package hotreload

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

type ClusterReloader struct{}

// CheckUpdate compares the old and new cluster configurations to determine if a reload is needed.
func (c *ClusterReloader) CheckUpdate(oldConfig, newConfig *model.Bootstrap) bool {
	oldClusters := oldConfig.StaticResources.Clusters
	newClusters := newConfig.StaticResources.Clusters

	if len(oldClusters) != len(newClusters) {
		logger.Infof("Cluster count changed: old=%d, new=%d", len(oldClusters), len(newClusters))
		return true
	}

	oldClusterMap := make(map[string]*model.ClusterConfig)
	for _, cluster := range oldClusters {
		oldClusterMap[cluster.Name] = cluster
	}

	for _, newCluster := range newClusters {
		oldCluster, exists := oldClusterMap[newCluster.Name]
		if !exists {
			logger.Infof("New cluster found: %s", newCluster.Name)
			return true
		}

		if len(oldCluster.Endpoints) != len(newCluster.Endpoints) {
			logger.Infof("Cluster %s endpoint count changed: old=%d, new=%d",
				newCluster.Name, len(oldCluster.Endpoints), len(newCluster.Endpoints))
			return true
		}

		for i := range newCluster.Endpoints {
			if i >= len(oldCluster.Endpoints) {
				return true
			}
			oldEp := oldCluster.Endpoints[i]
			newEp := newCluster.Endpoints[i]
			if oldEp.Address.Address != newEp.Address.Address ||
				oldEp.Address.Port != newEp.Address.Port {
				logger.Infof("Cluster %s endpoint changed: old=%s:%d, new=%s:%d",
					newCluster.Name,
					oldEp.Address.Address, oldEp.Address.Port,
					newEp.Address.Address, newEp.Address.Port)
				return true
			}
		}

		if oldCluster.LbStr != newCluster.LbStr {
			logger.Infof("Cluster %s lb policy changed: old=%s, new=%s",
				newCluster.Name, oldCluster.LbStr, newCluster.LbStr)
			return true
		}
	}

	return false
}

// HotReload applies the new cluster configuration.
func (c *ClusterReloader) HotReload(oldConfig, newConfig *model.Bootstrap) error {
	logger.Info("Starting cluster hot reload")

	srv := server.GetServer()
	if srv == nil {
		logger.Error("Server instance is nil")
		return errors.New("server instance is nil")
	}

	clusterManager := srv.GetClusterManager()
	if clusterManager == nil {
		logger.Error("Cluster manager is nil")
		return errors.New("cluster manager is nil")
	}

	oldClusterMap := make(map[string]*model.ClusterConfig)
	for _, cluster := range oldConfig.StaticResources.Clusters {
		oldClusterMap[cluster.Name] = cluster
	}

	newClusterMap := make(map[string]*model.ClusterConfig)
	for _, cluster := range newConfig.StaticResources.Clusters {
		newClusterMap[cluster.Name] = cluster
	}

	clustersToRemove := []string{}
	for name := range oldClusterMap {
		if _, exists := newClusterMap[name]; !exists {
			clustersToRemove = append(clustersToRemove, name)
		}
	}
	if len(clustersToRemove) > 0 {
		logger.Infof("Removing %d cluster(s): %v", len(clustersToRemove), clustersToRemove)
		clusterManager.RemoveCluster(clustersToRemove)
	}

	updated := 0
	added := 0
	for _, newCluster := range newConfig.StaticResources.Clusters {
		if oldCluster, exists := oldClusterMap[newCluster.Name]; exists {
			if c.clusterChanged(oldCluster, newCluster) {
				logger.Infof("Updating cluster: %s (endpoints: %d)", newCluster.Name, len(newCluster.Endpoints))
				clusterManager.UpdateCluster(newCluster)
				updated++
			}
		} else {
			logger.Infof("Adding new cluster: %s (endpoints: %d)", newCluster.Name, len(newCluster.Endpoints))
			clusterManager.AddCluster(newCluster)
			added++
		}
	}

	logger.Infof("Cluster hot reload completed: added=%d, updated=%d, removed=%d",
		added, updated, len(clustersToRemove))
	logger.Info("Cluster hot reload completed successfully")
	return nil
}

// clusterChanged checks if a cluster configuration has changed
func (c *ClusterReloader) clusterChanged(old, new *model.ClusterConfig) bool {
	if len(old.Endpoints) != len(new.Endpoints) {
		return true
	}

	for i := range new.Endpoints {
		if i >= len(old.Endpoints) {
			return true
		}
		oldEp := old.Endpoints[i]
		newEp := new.Endpoints[i]
		if oldEp.Address.Address != newEp.Address.Address ||
			oldEp.Address.Port != newEp.Address.Port {
			return true
		}
	}

	return old.LbStr != new.LbStr
}

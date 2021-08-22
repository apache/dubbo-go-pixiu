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

package server

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"sync"
)

type ClusterManager struct {
	rw      sync.RWMutex
	cConfig []*model.Cluster
}

type ClusterAdapter interface {
}

func CreateDefaultClusterManager(bs *model.Bootstrap) *ClusterManager {
	return &ClusterManager{cConfig: bs.StaticResources.Clusters}
}

func (cm *ClusterManager) AddCluster(c *model.Cluster) {
	cm.rw.Lock()
	defer cm.rw.Unlock()
	cm.cConfig = append(cm.cConfig, c)
}

func (cm *ClusterManager) UpdateCluster(new *model.Cluster) {
	cm.rw.Lock()
	defer cm.rw.Unlock()

	for i, c := range cm.cConfig {
		if c.Name == new.Name {
			cm.cConfig[i] = new
			break
		}
	}
	logger.Warnf("not found modified cluster %s", new.Name)
}

func (cm *ClusterManager) AddEndpoint(clusterName string, endpoint *model.Endpoint) {
	cm.rw.Lock()
	defer cm.rw.Unlock()

	for _, c := range cm.cConfig {
		if c.Name == clusterName {
			c.Endpoints = append(c.Endpoints, endpoint)
			break
		}
	}

	logger.Warnf("not found  cluster %s", clusterName)
}

func (cm *ClusterManager) DeleteEndpoint(clusterName string, endpointID string) {
	cm.rw.Lock()
	defer cm.rw.Unlock()

	for _, c := range cm.cConfig {
		if c.Name == clusterName {
			for i, e := range c.Endpoints {
				if e.ID == endpointID {
					c.Endpoints = append(c.Endpoints[:i], c.Endpoints[i+1:]...)
					break
				}
			}
			logger.Warnf("not found endpoint %s", endpointID)
			break
		}
	}
	logger.Warnf("not found  cluster %s", clusterName)

}

func (cm *ClusterManager) PickEndpoint(clusterName string) *model.Endpoint {
	cm.rw.RLock()
	defer cm.rw.RUnlock()

	for _, cluster := range cm.cConfig {
		if cluster.Name == clusterName {
			// according to lb to choose one endpoint, now only random
			return cluster.PickOneEndpoint()
		}
	}
	return nil
}

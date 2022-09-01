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
	"fmt"
	"sync"
	"sync/atomic"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/cluster"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/cluster/loadbalancer"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server/controls"
)

// generate cluster name for unnamed cluster
var clusterIndex int32 = 1

type (
	ClusterManager struct {
		rw sync.RWMutex

		store *ClusterStore
		//cConfig []*model.ClusterConfig
	}

	// ClusterStore store for cluster array
	ClusterStore struct {
		Config      []*model.ClusterConfig `yaml:"config" json:"config"`
		Version     int32                  `yaml:"version" json:"version"`
		clustersMap map[string]*cluster.Cluster
	}

	// xdsControlStore help convert ClusterStore to controls.ClusterStore interface
	xdsControlStore struct {
		*ClusterStore
	}
)

func (x *xdsControlStore) Config() []*model.ClusterConfig {
	return x.ClusterStore.Config
}

// CloneXdsControlStore clone cluster store for xds
func (cm *ClusterManager) CloneXdsControlStore() (controls.ClusterStore, error) {
	store, err := cm.CloneStore()
	return &xdsControlStore{store}, err
}

func CreateDefaultClusterManager(bs *model.Bootstrap) *ClusterManager {
	return &ClusterManager{store: newClusterStore(bs)}
}

func newClusterStore(bs *model.Bootstrap) *ClusterStore {
	store := &ClusterStore{
		clustersMap: map[string]*cluster.Cluster{},
	}
	for _, c := range bs.StaticResources.Clusters {
		store.AddCluster(c)
	}
	return store
}

func (cm *ClusterManager) AddCluster(c *model.ClusterConfig) {
	cm.rw.Lock()
	defer cm.rw.Unlock()

	cm.store.IncreaseVersion()
	cm.store.AddCluster(c)
}

func (cm *ClusterManager) UpdateCluster(new *model.ClusterConfig) {
	cm.rw.Lock()
	defer cm.rw.Unlock()

	cm.store.IncreaseVersion()
	cm.store.UpdateCluster(new)
}

func (cm *ClusterManager) SetEndpoint(clusterName string, endpoint *model.Endpoint) {
	cm.rw.Lock()
	defer cm.rw.Unlock()

	cm.store.IncreaseVersion()
	cm.store.SetEndpoint(clusterName, endpoint)
}

func (cm *ClusterManager) DeleteEndpoint(clusterName string, endpointID string) {
	cm.rw.Lock()
	defer cm.rw.Unlock()

	cm.store.IncreaseVersion()
	cm.store.DeleteEndpoint(clusterName, endpointID)
}

func (cm *ClusterManager) CloneStore() (*ClusterStore, error) {
	cm.rw.Lock()
	defer cm.rw.Unlock()

	b, err := yaml.MarshalYML(cm.store)
	if err != nil {
		return nil, err
	}

	c := &ClusterStore{
		clustersMap: map[string]*cluster.Cluster{},
	}
	if err := yaml.UnmarshalYML(b, c); err != nil {
		return nil, err
	}
	return c, nil
}

func (cm *ClusterManager) NewStore(version int32) *ClusterStore {
	cm.rw.Lock()
	defer cm.rw.Unlock()

	return &ClusterStore{Version: version, clustersMap: map[string]*cluster.Cluster{}}
}

func (cm *ClusterManager) CompareAndSetStore(store *ClusterStore) bool {
	cm.rw.Lock()
	defer cm.rw.Unlock()

	if store.Version != cm.store.Version {
		return false
	}

	cm.store = store
	return true
}

func (cm *ClusterManager) PickEndpoint(clusterName string) *model.Endpoint {
	cm.rw.RLock()
	defer cm.rw.RUnlock()

	for _, c := range cm.store.Config {
		if c.Name == clusterName {
			return cm.pickOneEndpoint(c)
		}
	}
	return nil
}

func (cm *ClusterManager) pickOneEndpoint(c *model.ClusterConfig) *model.Endpoint {
	if c.Endpoints == nil || len(c.Endpoints) == 0 {
		return nil
	}

	if len(c.Endpoints) == 1 {
		if !c.Endpoints[0].UnHealthy {
			return c.Endpoints[0]
		}
		return nil
	}

	loadBalancer, ok := loadbalancer.LoadBalancerStrategy[c.LbStr]
	if ok {
		return loadBalancer.Handler(c)
	}
	return loadbalancer.LoadBalancerStrategy[model.LoadBalancerRand].Handler(c)
}

func (cm *ClusterManager) RemoveCluster(namesToDel []string) {
	cm.rw.Lock()
	defer cm.rw.Unlock()

	for i, c := range cm.store.Config {
		if c == nil {
			continue
		}
		for _, name := range namesToDel { // suppose resource to remove and clusters is few
			if name == c.Name {
				removed := cm.store.Config[i]
				cm.store.clustersMap[removed.Name].Stop()
				cm.store.Config[i] = nil
				delete(cm.store.clustersMap, removed.Name)
			}
		}
	}
	//re-construct cm.store.Config remove nil element
	for i := 0; i < len(cm.store.Config); {
		if cm.store.Config[i] != nil {
			i++
			continue
		}
		cm.store.Config = append(cm.store.Config[:i], cm.store.Config[i+1:]...)
	}
	cm.store.IncreaseVersion()
}

func (cm *ClusterManager) HasCluster(clusterName string) bool {
	cm.rw.Lock()
	defer cm.rw.Unlock()
	return cm.store.HasCluster(clusterName)
}

func (s *ClusterStore) AddCluster(c *model.ClusterConfig) {
	if c.Name == "" {
		index := atomic.AddInt32(&clusterIndex, 1)
		c.Name = fmt.Sprintf("cluster%d", index)
	}
	s.Config = append(s.Config, c)
	s.clustersMap[c.Name] = cluster.NewCluster(c)
	c.CreateConsistentHash()
}

func (s *ClusterStore) UpdateCluster(new *model.ClusterConfig) {

	for i, c := range s.Config {
		if c.Name == new.Name {
			s.Config[i] = new
			return
		}
	}
	logger.Warnf("not found modified cluster %s", new.Name)
}

func (s *ClusterStore) SetEndpoint(clusterName string, endpoint *model.Endpoint) {
	cluster := s.clustersMap[clusterName]
	if cluster == nil {
		c := &model.ClusterConfig{Name: clusterName, LbStr: model.LoadBalancerRoundRobin, Endpoints: []*model.Endpoint{}}
		s.AddCluster(c)
		cluster = s.clustersMap[clusterName]
	}

	for _, c := range s.Config {
		if c.Name == clusterName {
			for _, e := range c.Endpoints {
				// endpoint update
				if e.ID == endpoint.ID {
					cluster.RemoveEndpoint(e)
					e.Name = endpoint.Name
					e.Metadata = endpoint.Metadata
					e.Address = endpoint.Address
					cluster.AddEndpoint(e)
					return
				}
			}
			// endpoint create
			c.Endpoints = append(c.Endpoints, endpoint)
			cluster.AddEndpoint(endpoint)
			if c.Hash.ConsistentHash != nil {
				c.Hash.ConsistentHash.Add(endpoint.Address.Address)
			}
			return
		}
	}
}

func (s *ClusterStore) DeleteEndpoint(clusterName string, endpointID string) {
	cluster := s.clustersMap[clusterName]
	if cluster == nil {
		return
	}
	for _, c := range s.Config {
		if c.Name == clusterName {
			for i, e := range c.Endpoints {
				if e.ID == endpointID {
					cluster.RemoveEndpoint(e)
					c.Endpoints = append(c.Endpoints[:i], c.Endpoints[i+1:]...)
					if c.Hash.ConsistentHash != nil {
						c.Hash.ConsistentHash.Remove(e.Address.Address)
					}
					return
				}
			}
			logger.Warnf("not found endpoint %s", endpointID)
			return
		}
	}
	logger.Warnf("not found  cluster %s", clusterName)
}

func (s *ClusterStore) HasCluster(clusterName string) bool {
	for _, c := range s.Config {
		if c.Name == clusterName {
			return true
		}
	}
	return false
}

func (s *ClusterStore) IncreaseVersion() {
	atomic.AddInt32(&s.Version, 1)
}

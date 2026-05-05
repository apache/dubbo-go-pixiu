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
	"github.com/hashicorp/go-uuid"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cluster"
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server/controls"
)

// generate cluster name for unnamed cluster
var (
	clusterIndex int32 = 1
)

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

// CompareAndSetStore swaps the store only when versions match.
// Version mismatch must leave both stores and runtime clusters untouched.
func (cm *ClusterManager) CompareAndSetStore(store *ClusterStore) bool {
	swapped, replacedClusters := cm.compareAndSetStore(store)
	if !swapped {
		return false
	}

	// Stop old runtime after publishing the swap; Stop may touch timers/goroutines.
	stopClusters(replacedClusters)
	return true
}

func (cm *ClusterManager) compareAndSetStore(store *ClusterStore) (bool, []*cluster.Cluster) {
	cm.rw.Lock()
	defer cm.rw.Unlock()

	if store.Version != cm.store.Version {
		return false, nil
	}

	currentStore := cm.store
	replacedClusters := store.ensureRuntimeClusters()
	store.carryOverRuntimeStateFrom(currentStore)
	cm.store = store
	if store != currentStore {
		replacedClusters = append(replacedClusters, currentStore.runtimeClustersNotIn(store)...)
	}
	return true, replacedClusters
}

// PickEndpoint picks an endpoint from the cluster by its name and load balancing policy.
func (cm *ClusterManager) PickEndpoint(clusterName string, policy model.LbPolicy) *model.Endpoint {
	cm.rw.RLock()
	defer cm.rw.RUnlock()
	c := cm.getCluster(clusterName)
	if c == nil {
		logger.Warnf("[dubbo-go-pixiu] cluster %s not found", clusterName)
		return nil
	}
	return cm.pickOneEndpoint(c, policy)
}

// PickNextEndpoint picks the next endpoint in the cluster after the current endpoint ID.
func (cm *ClusterManager) PickNextEndpoint(clusterName string, curEndpointID string) *model.Endpoint {
	cm.rw.RLock()
	defer cm.rw.RUnlock()

	c := cm.getCluster(clusterName)
	if c == nil {
		logger.Warnf("[dubbo-go-pixiu] cluster %s not found", clusterName)
		return nil
	}

	for i, endpoint := range c.Endpoints {
		if endpoint.ID == curEndpointID {
			// pick next endpoint
			if i < len(c.Endpoints)-1 {
				return c.Endpoints[i+1]
			}
			return nil // have tried all endpoints
		}
	}

	return nil
}

// GetEndpointByID returns the endpoint by ID in the given cluster.
func (cm *ClusterManager) GetEndpointByID(clusterName string, endpointID string) *model.Endpoint {
	cm.rw.RLock()
	defer cm.rw.RUnlock()

	c := cm.getCluster(clusterName)
	if c == nil {
		return nil
	}
	for _, endpoint := range c.Endpoints {
		if endpoint.ID == endpointID && !endpoint.UnHealthy {
			return endpoint
		}
	}
	return nil
}

// getCluster returns the cluster configuration by its name.
func (cm *ClusterManager) getCluster(clusterName string) *model.ClusterConfig {
	for _, c := range cm.store.Config {
		if c.Name == clusterName {
			return c
		}
	}
	return nil
}

func (cm *ClusterManager) pickOneEndpoint(c *model.ClusterConfig, policy model.LbPolicy) *model.Endpoint {
	if len(c.Endpoints) == 0 {
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
		return loadBalancer.Handler(c, policy)
	}
	return loadbalancer.LoadBalancerStrategy[model.LoadBalancerRand].Handler(c, policy)
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
				stopClusters([]*cluster.Cluster{cm.store.clustersMap[removed.Name]})
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
		c.Name = fmt.Sprintf("cluster-%d", clusterIndex)
		atomic.AddInt32(&clusterIndex, 1)
	}

	s.prepareClusterConfig(c)

	s.Config = append(s.Config, c)
	stopClusters([]*cluster.Cluster{s.replaceClusterRuntime(c.Name, c)})
}

// prepareClusterConfig rebuilds endpoint defaults and hash from current endpoints.
func (s *ClusterStore) prepareClusterConfig(c *model.ClusterConfig) {
	s.assembleClusterEndpoints(c)
	c.CreateConsistentHash()
}

// assembleClusterEndpoints assembles the cluster endpoints
// by formatting the ID, name and domains for each endpoint
// If endpoint.LLMMeta is not nil, the assimilation of name and domain is based on
// the LLM provider denoted in the endpoint LLMMeta.
func (s *ClusterStore) assembleClusterEndpoints(c *model.ClusterConfig) {
	if c == nil {
		return
	}

	for i, endpoint := range c.Endpoints {
		// If the endpoint ID is not set, set it to the index + 1
		if endpoint.ID == "" {
			endpoint.ID, _ = uuid.GenerateUUID()
		}

		// If the endpoint has no name, set a default name
		if endpoint.Name == "" && endpoint.LLMMeta != nil {
			endpoint.Name = fmt.Sprintf("endpoint-%d#%s", i+1, endpoint.LLMMeta.Provider)
		} else if endpoint.Name == "" && endpoint.LLMMeta == nil {
			endpoint.Name = fmt.Sprintf("endpoint-%d", i+1)
		}
	}
}

// replaceClusterRuntime returns the old runtime so callers decide when to stop it.
func (s *ClusterStore) replaceClusterRuntime(name string, config *model.ClusterConfig) *cluster.Cluster {
	if s.clustersMap == nil {
		s.clustersMap = map[string]*cluster.Cluster{}
	}

	oldRuntime := s.clustersMap[name]
	s.clustersMap[name] = cluster.NewCluster(config)
	return oldRuntime
}

// ensureRuntimeClusters repairs clustersMap to match Config by name and pointer.
func (s *ClusterStore) ensureRuntimeClusters() []*cluster.Cluster {
	if s == nil {
		return nil
	}
	if s.clustersMap == nil {
		s.clustersMap = map[string]*cluster.Cluster{}
	}

	replacedClusters := make([]*cluster.Cluster, 0)
	configsByName := make(map[string]*model.ClusterConfig, len(s.Config))
	for _, clusterConfig := range s.Config {
		if clusterConfig == nil {
			continue
		}
		s.prepareClusterConfig(clusterConfig)
		configsByName[clusterConfig.Name] = clusterConfig

		runtimeCluster := s.clustersMap[clusterConfig.Name]
		if runtimeCluster == nil || runtimeCluster.Config != clusterConfig {
			if oldRuntime := s.replaceClusterRuntime(clusterConfig.Name, clusterConfig); oldRuntime != nil {
				replacedClusters = append(replacedClusters, oldRuntime)
			}
		}
	}

	for name, runtimeCluster := range s.clustersMap {
		if _, ok := configsByName[name]; !ok {
			replacedClusters = append(replacedClusters, runtimeCluster)
			delete(s.clustersMap, name)
		}
	}
	return replacedClusters
}

// runtimeClustersNotIn finds old runtime clusters dropped by the next store.
func (s *ClusterStore) runtimeClustersNotIn(next *ClusterStore) []*cluster.Cluster {
	if s == nil {
		return nil
	}

	nextClusters := make(map[*cluster.Cluster]struct{})
	if next != nil {
		for _, runtimeCluster := range next.clustersMap {
			if runtimeCluster != nil {
				nextClusters[runtimeCluster] = struct{}{}
			}
		}
	}

	replacedClusters := make([]*cluster.Cluster, 0)
	for _, runtimeCluster := range s.clustersMap {
		if runtimeCluster == nil {
			continue
		}
		if _, ok := nextClusters[runtimeCluster]; !ok {
			replacedClusters = append(replacedClusters, runtimeCluster)
		}
	}
	return replacedClusters
}

// stopClusters is nil-safe and avoids stopping the same runtime twice.
func stopClusters(clusters []*cluster.Cluster) {
	stopped := make(map[*cluster.Cluster]struct{}, len(clusters))
	for _, runtimeCluster := range clusters {
		if runtimeCluster == nil {
			continue
		}
		if _, ok := stopped[runtimeCluster]; ok {
			continue
		}
		stopped[runtimeCluster] = struct{}{}
		runtimeCluster.Stop()
	}
}

// UpdateCluster replaces config/runtime together while preserving the RR cursor.
func (s *ClusterStore) UpdateCluster(new *model.ClusterConfig) {
	for i, c := range s.Config {
		if c == nil {
			continue
		}
		if c.Name == new.Name {
			s.prepareClusterConfig(new)
			atomic.StoreUint32(
				&new.PrePickEndpointIndex,
				atomic.LoadUint32(&c.PrePickEndpointIndex),
			)
			s.Config[i] = new
			stopClusters([]*cluster.Cluster{s.replaceClusterRuntime(new.Name, new)})
			return
		}
	}
	logger.Warnf("not found modified cluster %s", new.Name)
}

func (s *ClusterStore) SetEndpoint(clusterName string, endpoint *model.Endpoint) {
	clusterConfig := s.findClusterConfig(clusterName)
	if clusterConfig == nil {
		c := &model.ClusterConfig{Name: clusterName, LbStr: model.LoadBalancerRoundRobin, Endpoints: []*model.Endpoint{}}
		s.AddCluster(c)
		clusterConfig = c
	}

	runtimeCluster := s.clustersMap[clusterName]
	if runtimeCluster == nil || runtimeCluster.Config != clusterConfig {
		stopClusters([]*cluster.Cluster{s.replaceClusterRuntime(clusterName, clusterConfig)})
		runtimeCluster = s.clustersMap[clusterName]
	}

	for _, e := range clusterConfig.Endpoints {
		if e.ID == endpoint.ID {
			// Remove before mutating address because healthcheck keys by address.
			runtimeCluster.RemoveEndpoint(e)
			e.Name = endpoint.Name
			e.Metadata = endpoint.Metadata
			e.Address = endpoint.Address
			s.prepareClusterConfig(clusterConfig)
			runtimeCluster.AddEndpoint(e)
			return
		}
	}
	clusterConfig.Endpoints = append(clusterConfig.Endpoints, endpoint)
	s.prepareClusterConfig(clusterConfig)
	runtimeCluster.AddEndpoint(endpoint)
}

func (s *ClusterStore) DeleteEndpoint(clusterName string, endpointID string) {
	clusterConfig := s.findClusterConfig(clusterName)
	if clusterConfig == nil {
		logger.Warnf("not found cluster %s", clusterName)
		return
	}

	runtimeCluster := s.clustersMap[clusterName]
	if runtimeCluster == nil || runtimeCluster.Config != clusterConfig {
		stopClusters([]*cluster.Cluster{s.replaceClusterRuntime(clusterName, clusterConfig)})
		runtimeCluster = s.clustersMap[clusterName]
	}

	for i, e := range clusterConfig.Endpoints {
		if e.ID == endpointID {
			runtimeCluster.RemoveEndpoint(e)
			clusterConfig.Endpoints = append(clusterConfig.Endpoints[:i], clusterConfig.Endpoints[i+1:]...)
			s.prepareClusterConfig(clusterConfig)
			return
		}
	}
	logger.Warnf("not found endpoint %s", endpointID)
}

func (s *ClusterStore) findClusterConfig(clusterName string) *model.ClusterConfig {
	for _, c := range s.Config {
		if c != nil && c.Name == clusterName {
			return c
		}
	}
	return nil
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

func (s *ClusterStore) carryOverRuntimeStateFrom(old *ClusterStore) {
	if s == nil || old == nil {
		return
	}

	oldConfigsByName := make(map[string]*model.ClusterConfig, len(old.Config))
	for _, clusterConfig := range old.Config {
		if clusterConfig != nil {
			oldConfigsByName[clusterConfig.Name] = clusterConfig
		}
	}

	// Preserve runtime-only load-balancer state when a rebuilt store is swapped in.
	for _, clusterConfig := range s.Config {
		if clusterConfig == nil {
			continue
		}
		if oldConfig := oldConfigsByName[clusterConfig.Name]; oldConfig != nil {
			atomic.StoreUint32(
				&clusterConfig.PrePickEndpointIndex,
				atomic.LoadUint32(&oldConfig.PrePickEndpointIndex),
			)
		}
	}
}

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
	"reflect"
	"sync"
	"sync/atomic"
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

// SetEndpoint registers or refreshes a single endpoint in the named
// cluster, creating the cluster on first use.
//
// Dedup contract — follows the same collision rules as
// assembleClusterEndpoints (ID match → hash match → suffix on content
// mismatch) so static and dynamic registrations agree on outcomes, with
// one deliberate divergence around Name documented below:
//
//   - When the incoming endpoint matches an existing slot's routing
//     identity AND its content (Address, Metadata, LLMMeta — see
//     endpointContentEqualForSet), the call is idempotent. The existing
//     slot stays in place; this is the safety net for replayed Nacos
//     instance events.
//   - When the incoming endpoint collides on explicit ID or generated hash
//     with an existing slot but differs in content, the new endpoint is
//     appended with a -2/-3/... suffix (a WARN is emitted, mirroring
//     assembleClusterEndpoints).
//   - Otherwise the endpoint is appended with its explicit or generated ID.
//
// Name is intentionally NOT part of the content equality check (see
// endpointContentEqualForSet's godoc). As a consequence, Name-only updates
// (same Address/LLMMeta/Metadata, different Name) are absorbed by the
// idempotent path and the cluster keeps the original Name. Callers that
// need to rename an endpoint must DeleteEndpoint + SetEndpoint, or replace
// the whole cluster config via AddCluster/UpdateCluster.
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
	var replacedClusters []*cluster.Cluster
	if store == currentStore {
		replacedClusters = store.ensureRuntimeClusters()
	} else {
		replacedClusters = store.ensureRuntimeClustersFrom(currentStore)
	}
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
	runtimeCluster := cm.getRuntimeCluster(clusterName)
	if runtimeCluster == nil {
		logger.Warnf("[dubbo-go-pixiu] cluster %s not found", clusterName)
		return nil
	}
	return cm.pickOneEndpoint(runtimeCluster, policy)
}

// PickNextEndpoint picks the next endpoint in the cluster after the current endpoint ID.
func (cm *ClusterManager) PickNextEndpoint(clusterName, curEndpointID string) *model.Endpoint {
	cm.rw.RLock()
	defer cm.rw.RUnlock()

	runtimeCluster := cm.getRuntimeCluster(clusterName)
	if runtimeCluster == nil {
		logger.Warnf("[dubbo-go-pixiu] cluster %s not found", clusterName)
		return nil
	}

	return runtimeCluster.EndpointSnapshot().NextHealthyEndpoint(curEndpointID)
}

// GetEndpointByID returns the healthy runtime endpoint by ID in the given
// cluster.
//
// Deprecated: kept as an alias for backward compatibility with v1.x callers.
// New code should call GetHealthyEndpointByID directly to make the
// healthy-only semantic explicit, or GetAnyEndpointByID when the caller also
// wants to observe runtime-unhealthy endpoints.
func (cm *ClusterManager) GetEndpointByID(clusterName, endpointID string) *model.Endpoint {
	return cm.GetHealthyEndpointByID(clusterName, endpointID)
}

// GetAnyEndpointByID returns the runtime endpoint by ID regardless of its
// current health in the runtime snapshot.
func (cm *ClusterManager) GetAnyEndpointByID(clusterName, endpointID string) *model.Endpoint {
	cm.rw.RLock()
	defer cm.rw.RUnlock()

	runtimeCluster := cm.getRuntimeCluster(clusterName)
	if runtimeCluster == nil {
		return nil
	}
	return runtimeCluster.EndpointSnapshot().EndpointByID(endpointID)
}

// GetHealthyEndpointByID returns the runtime endpoint by ID only when it is
// healthy in the current runtime snapshot.
func (cm *ClusterManager) GetHealthyEndpointByID(clusterName, endpointID string) *model.Endpoint {
	cm.rw.RLock()
	defer cm.rw.RUnlock()

	runtimeCluster := cm.getRuntimeCluster(clusterName)
	if runtimeCluster == nil {
		return nil
	}
	return runtimeCluster.EndpointSnapshot().HealthyEndpointByID(endpointID)
}

// getCluster returns the cluster configuration by its name.
func (cm *ClusterManager) getCluster(clusterName string) *model.ClusterConfig {
	runtimeCluster := cm.getRuntimeCluster(clusterName)
	if runtimeCluster == nil {
		return nil
	}
	return runtimeCluster.Config
}

func (cm *ClusterManager) getRuntimeCluster(clusterName string) *cluster.Cluster {
	return cm.store.clustersMap[clusterName]
}

func (cm *ClusterManager) pickOneEndpoint(runtimeCluster *cluster.Cluster, policy model.LbPolicy) *model.Endpoint {
	snapshot := runtimeCluster.EndpointSnapshot()
	if snapshot.HealthyEndpointCount() == 0 {
		return nil
	}
	healthyEndpoints := snapshot.HealthyEndpointsForPick()

	c := runtimeCluster.Config
	loadBalancer, ok := loadbalancer.LoadBalancerStrategy[c.LbStr]
	if !ok {
		loadBalancer = loadbalancer.LoadBalancerStrategy[model.LoadBalancerRand]
	}

	var allEndpoints []*model.Endpoint
	if loadbalancer.NeedsAllEndpoints(loadBalancer) {
		allEndpoints = snapshot.AllEndpointsForPick()
	}

	return loadbalancer.PickEndpoint(loadBalancer, loadbalancer.PickContext{
		Config:                c,
		HealthyConsistentHash: snapshot.HealthyConsistentHash(),
		AllEndpoints:          allEndpoints,
		HealthyEndpoints:      healthyEndpoints,
	}, policy)
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

// assembleClusterEndpoints assembles the cluster endpoints by formatting the
// ID, name and domains for each endpoint. If endpoint.LLMMeta is not nil, the
// assimilation of name and domain is based on the LLM provider denoted in the
// endpoint LLMMeta. The store first deep-clones c.Endpoints, so ID/name
// defaulting never mutates operator-supplied *model.Endpoint values.
func (s *ClusterStore) assembleClusterEndpoints(c *model.ClusterConfig) {
	if c == nil {
		return
	}

	c.Endpoints = model.CloneEndpoints(c.Endpoints)
	endpointIDs := make(map[string]struct{}, len(c.Endpoints))
	for i, endpoint := range c.Endpoints {
		if endpoint == nil {
			continue
		}
		// Endpoint IDs are runtime health keys, so keep them unique per cluster.
		if endpoint.ID == "" {
			endpoint.ID = nextStableEndpointID(c.Name, endpoint, endpointIDs)
		} else if _, exists := endpointIDs[endpoint.ID]; exists {
			duplicateID := endpoint.ID
			endpoint.ID = nextStableEndpointID(c.Name, endpoint, endpointIDs)
			logger.Warnf(
				"[dubbo-go-pixiu] duplicate endpoint ID %s in cluster %s, assigned endpoint ID %s",
				duplicateID,
				c.Name,
				endpoint.ID,
			)
		}
		endpointIDs[endpoint.ID] = struct{}{}

		// If the endpoint has no name, set a default name
		if endpoint.Name == "" && endpoint.LLMMeta != nil {
			endpoint.Name = fmt.Sprintf("endpoint-%d#%s", i+1, endpoint.LLMMeta.Provider)
		} else if endpoint.Name == "" && endpoint.LLMMeta == nil {
			endpoint.Name = fmt.Sprintf("endpoint-%d", i+1)
		}
	}
}

// nextStableEndpointID returns a unique endpoint ID for the cluster's dedup
// set. If endpoint.ID is set (operator-supplied) and only collides with a
// sibling, it appends -2, -3, ... to preserve the operator's choice. If
// endpoint.ID is empty, it derives a deterministic generated-* base via
// model.GenerateEndpointID and suffixes that on collision. This matches
// uniqueSnapshotEndpointID in pkg/cluster and keeps the same dashboard/log
// identity post-rebuild.
func nextStableEndpointID(clusterName string, endpoint *model.Endpoint, endpointIDs map[string]struct{}) string {
	baseID := ""
	if endpoint != nil {
		baseID = endpoint.ID
	}
	if baseID == "" {
		baseID = model.GenerateEndpointID(clusterName, endpoint)
	}
	if _, exists := endpointIDs[baseID]; !exists {
		return baseID
	}
	for suffix := 2; ; suffix++ {
		candidate := fmt.Sprintf("%s-%d", baseID, suffix)
		if _, exists := endpointIDs[candidate]; !exists {
			return candidate
		}
	}
}

// replaceClusterRuntime returns the old runtime so callers decide when to stop it.
func (s *ClusterStore) replaceClusterRuntime(name string, config *model.ClusterConfig) *cluster.Cluster {
	var previous *cluster.EndpointSnapshot
	if oldRuntime := s.clustersMap[name]; oldRuntime != nil {
		previous = oldRuntime.SnapshotForRuntimeReplacement()
	}
	return s.replaceClusterRuntimeWithSnapshot(name, config, previous)
}

func (s *ClusterStore) replaceClusterRuntimeWithSnapshot(
	name string,
	config *model.ClusterConfig,
	previous *cluster.EndpointSnapshot,
) *cluster.Cluster {
	s.ensureRuntimeClusterMap()

	oldRuntime := s.clustersMap[name]
	s.clustersMap[name] = cluster.NewClusterWithEndpointSnapshot(config, previous)
	return oldRuntime
}

func (s *ClusterStore) ensureRuntimeClusterMap() {
	if s.clustersMap == nil {
		s.clustersMap = map[string]*cluster.Cluster{}
	}
}

// ensureRuntimeClusters repairs clustersMap to match Config by name and pointer.
func (s *ClusterStore) ensureRuntimeClusters() []*cluster.Cluster {
	if s == nil {
		return nil
	}
	s.ensureRuntimeClusterMap()

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

	return append(replacedClusters, s.removeRuntimeClustersNotIn(configsByName)...)
}

// ensureRuntimeClustersFrom repairs a candidate store using the currently
// published runtime snapshots before the candidate runtime starts health checks.
func (s *ClusterStore) ensureRuntimeClustersFrom(old *ClusterStore) []*cluster.Cluster {
	if s == nil {
		return nil
	}
	if old == nil || s == old {
		return s.ensureRuntimeClusters()
	}
	s.ensureRuntimeClusterMap()

	replacedClusters := make([]*cluster.Cluster, 0)
	configsByName := make(map[string]*model.ClusterConfig, len(s.Config))
	for _, clusterConfig := range s.Config {
		clusterName, oldRuntime := s.repairRuntimeClusterFromPrevious(clusterConfig, old)
		if clusterName == "" {
			continue
		}
		configsByName[clusterName] = clusterConfig
		if oldRuntime != nil {
			replacedClusters = append(replacedClusters, oldRuntime)
		}
	}

	return append(replacedClusters, s.removeRuntimeClustersNotIn(configsByName)...)
}

func (s *ClusterStore) repairRuntimeClusterFromPrevious(
	clusterConfig *model.ClusterConfig,
	old *ClusterStore,
) (string, *cluster.Cluster) {
	if clusterConfig == nil {
		return "", nil
	}
	s.prepareClusterConfig(clusterConfig)
	previous := snapshotForRuntimeReplacement(old, clusterConfig.Name)
	runtimeCluster := s.clustersMap[clusterConfig.Name]
	if runtimeCluster != nil && runtimeCluster.Config == clusterConfig && previous == nil {
		return clusterConfig.Name, nil
	}
	return clusterConfig.Name, s.replaceClusterRuntimeWithSnapshot(clusterConfig.Name, clusterConfig, previous)
}

func snapshotForRuntimeReplacement(old *ClusterStore, clusterName string) *cluster.EndpointSnapshot {
	if old == nil {
		return nil
	}
	if oldRuntime := old.clustersMap[clusterName]; oldRuntime != nil {
		return oldRuntime.SnapshotForRuntimeReplacement()
	}
	return nil
}

func (s *ClusterStore) removeRuntimeClustersNotIn(configsByName map[string]*model.ClusterConfig) []*cluster.Cluster {
	replacedClusters := make([]*cluster.Cluster, 0)
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
	endpoint = model.CloneEndpoint(endpoint)
	if endpoint == nil {
		return
	}

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

	targetID, idempotent := resolveSetEndpointSlot(clusterName, endpoint, clusterConfig.Endpoints)
	endpoint.ID = targetID
	if idempotent {
		// A content-equal endpoint already occupies this slot, so the cluster
		// state, consistent hash, and runtime health-check registration are
		// already correct. Returning here keeps re-registration idempotent —
		// the LLM/Nacos path can replay the same instance event without
		// growing the endpoint slice.
		return
	}

	clusterConfig.Endpoints = append(clusterConfig.Endpoints, endpoint)
	s.prepareClusterConfig(clusterConfig)
	runtimeCluster.RefreshEndpoints()
	runtimeCluster.AddEndpoint(endpoint)
}

// resolveSetEndpointSlot enforces the same dedup invariant as
// assembleClusterEndpoints on the dynamic SetEndpoint path: an endpoint ID
// identifies exactly one runtime endpoint at a time, and two registrations
// that differ in routing-relevant content must both survive (suffixed -2,
// -3, ...) instead of one silently overwriting the other.
//
// It returns the ID the incoming endpoint should be stored under and whether
// the call is an idempotent re-registration (an existing slot already has
// matching content). When idempotent is false the caller should append a
// fresh endpoint with the returned ID.
//
// Decision tree:
//
//  1. Explicit ID matches an existing slot → idempotent iff content equal,
//     otherwise return a -2/-3/... suffix off the explicit ID.
//  2. Empty ID, hash matches an existing slot → idempotent iff content
//     equal, otherwise return a suffix off the generated hash.
//  3. No collision → use the explicit ID, or the generated hash for empty
//     IDs.
//
// Share this helper across new entry points (admin API, xDS, gRPC
// reflection, ...) to keep the dedup contract uniform; reimplementing the
// branches in each caller is how the static and dynamic paths drifted
// apart in the first place.
func resolveSetEndpointSlot(clusterName string, incoming *model.Endpoint, existing []*model.Endpoint) (string, bool) {
	if incoming == nil {
		return "", false
	}
	if incoming.ID != "" {
		return resolveSetEndpointSlotByID(clusterName, incoming, existing)
	}
	return resolveSetEndpointSlotByHash(clusterName, incoming, existing)
}

// resolveSetEndpointSlotByID handles the explicit-ID branch of the decision
// tree documented on resolveSetEndpointSlot. Probing the existing slice by
// ID equality keeps the operator's pinned ID readable in the suffix on a
// content mismatch (foo → foo-2, not foo → generated-...-2).
func resolveSetEndpointSlotByID(clusterName string, incoming *model.Endpoint, existing []*model.Endpoint) (string, bool) {
	for _, e := range existing {
		if e == nil {
			continue
		}
		if e.ID != incoming.ID {
			continue
		}
		if endpointContentEqualForSet(e, incoming) {
			return incoming.ID, true
		}
		suffixedID := nextStableEndpointID(clusterName, incoming, existingEndpointIDs(existing))
		logSetEndpointSuffix(clusterName, incoming.ID, suffixedID)
		return suffixedID, false
	}
	return incoming.ID, false
}

// resolveSetEndpointSlotByHash handles the empty-ID branch of the decision
// tree documented on resolveSetEndpointSlot. Probing by generated hash —
// not by ID — is what lets an empty-ID call reuse an operator-pinned slot
// (truth-table row 3) instead of inventing a parallel generated-* entry.
func resolveSetEndpointSlotByHash(clusterName string, incoming *model.Endpoint, existing []*model.Endpoint) (string, bool) {
	incomingHash := model.GenerateEndpointID(clusterName, incoming)
	for _, e := range existing {
		if e == nil {
			continue
		}
		if model.GenerateEndpointID(clusterName, e) != incomingHash {
			continue
		}
		if endpointContentEqualForSet(e, incoming) {
			return e.ID, true
		}
		suffixedID := nextStableEndpointID(clusterName, incoming, existingEndpointIDs(existing))
		logSetEndpointSuffix(clusterName, incomingHash, suffixedID)
		return suffixedID, false
	}
	return incomingHash, false
}

// logSetEndpointSuffix uses the same WARN format as
// assembleClusterEndpoints (pkg/server/cluster_manager.go) so operators
// tail one log line whether a duplicate came from static YAML or from a
// runtime SetEndpoint call. Note: the dynamic path additionally logs
// empty-ID hash collisions (when two SetEndpoint calls hash-collide but
// differ in content), which assemble currently leaves silent — operator
// visibility was the explicit motivation for adding this log on the
// dynamic path.
//
// This log is intentionally NOT pinned by a test. The contractual floor
// is that suffix decisions actually happen (locked by the
// *AppendsSuffix tests); the WARN is operator-visibility supplement.
// Pinning the exact log line via stderr capture would couple the tests
// to log formatting and offer little benefit. If a refactor accidentally
// drops the call site, the suffix tests still pass — the regression
// shows up as operators not seeing collision warnings they were
// previously relying on, which is a separate signal worth treating as
// such.
func logSetEndpointSuffix(clusterName, collidingID, assignedID string) {
	logger.Warnf(
		"[dubbo-go-pixiu] duplicate endpoint ID %s in cluster %s, assigned endpoint ID %s",
		collidingID,
		clusterName,
		assignedID,
	)
}

// endpointContentEqualForSet reports whether the incoming SetEndpoint
// payload represents the same endpoint as an existing slot for the purposes
// of dedup. It compares the routing-relevant content — Address (full
// SocketAddress including Domains), Metadata, and LLMMeta — but
// deliberately excludes ID (the dedup contract uses the ID as the collision
// key, not as a content field), Name (assembleClusterEndpoints defaults
// missing names to "endpoint-<index>" between calls, so comparing Name
// would treat re-registration with no name as a duplicate-content case),
// and runtime-only state.
//
// nil and zero-length maps/slices are treated as equal: callers building
// endpoints from different code paths (static YAML vs. Nacos metadata) often
// produce one or the other for the same logical "no metadata" state.
func endpointContentEqualForSet(a, b *model.Endpoint) bool {
	if a == nil || b == nil {
		return a == b
	}
	if !socketAddressContentEqual(a.Address, b.Address) {
		return false
	}
	if !stringMapEqualForSet(a.Metadata, b.Metadata) {
		return false
	}
	return llmMetaEqualForSet(a.LLMMeta, b.LLMMeta)
}

// socketAddressContentEqual is a deliberately stricter equality oracle than
// model.SocketAddress.Equal. SocketAddress.Equal (see its godoc) is the
// identity oracle used by load-balancer routing, which only compares
// Address+Port (or Domains[0]) and ignores ResolverName, CertsDir, and
// Domains beyond the first element. That is the right semantic for "pick
// an endpoint by address" but too lax for SetEndpoint dedup: two endpoints
// that differ only in CertsDir or Domains list ordering are distinct
// configurations and must not be silently treated as one. So we compare
// every field here, with the nil/empty-Domains normalization callers
// expect.
func socketAddressContentEqual(a, b model.SocketAddress) bool {
	if a.Address != b.Address || a.Port != b.Port {
		return false
	}
	if a.ResolverName != b.ResolverName || a.CertsDir != b.CertsDir {
		return false
	}
	if len(a.Domains) != len(b.Domains) {
		return false
	}
	for i := range a.Domains {
		if a.Domains[i] != b.Domains[i] {
			return false
		}
	}
	return true
}

func stringMapEqualForSet(a, b map[string]string) bool {
	if len(a) != len(b) {
		return false
	}
	for k, v := range a {
		bv, ok := b[k]
		if !ok || bv != v {
			return false
		}
	}
	return true
}

func llmMetaEqualForSet(a, b *model.LLMMeta) bool {
	if a == nil || b == nil {
		return a == b
	}
	return reflect.DeepEqual(*a, *b)
}

func existingEndpointIDs(endpoints []*model.Endpoint) map[string]struct{} {
	ids := make(map[string]struct{}, len(endpoints))
	for _, endpoint := range endpoints {
		if endpoint == nil || endpoint.ID == "" {
			continue
		}
		ids[endpoint.ID] = struct{}{}
	}
	return ids
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
			runtimeCluster.RefreshEndpoints()
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

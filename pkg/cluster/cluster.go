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

package cluster

import (
	"fmt"
	"sync"
	"sync/atomic"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/healthcheck"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type Cluster struct {
	HealthCheck *healthcheck.HealthChecker
	// Config is the desired cluster configuration. Runtime picks read the
	// published EndpointSnapshot, so direct edits to Config.Endpoints or
	// Endpoint.UnHealthy are not observed by PickEndpoint immediately. Publish
	// membership/address changes with RefreshEndpoints; publish runtime health
	// changes with UpdateEndpointHealth, or RefreshEndpoints for clusters
	// without health checks.
	Config             *model.ClusterConfig
	healthMu           sync.Mutex
	legacyPickMu       sync.Mutex
	acceptHealthEvents bool
	endpoints          atomic.Pointer[EndpointSnapshot]
}

func NewCluster(clusterConfig *model.ClusterConfig) *Cluster {
	return NewClusterWithEndpointSnapshot(clusterConfig, nil)
}

func NewClusterWithEndpointSnapshot(clusterConfig *model.ClusterConfig, previous *EndpointSnapshot) *Cluster {
	c := &Cluster{
		Config:             clusterConfig,
		acceptHealthEvents: true,
	}
	c.RefreshEndpointsFrom(previous)

	// only handle one health checker
	if len(c.Config.HealthChecks) != 0 {
		c.HealthCheck = healthcheck.CreateHealthCheckWithCallback(
			clusterConfig,
			c.Config.HealthChecks[0],
			c.handleEndpointHealth,
		)
		c.HealthCheck.Start()
	}
	return c
}

func (c *Cluster) Stop() {
	if c.HealthCheck != nil {
		c.HealthCheck.Stop()
	}
}

func (c *Cluster) RemoveEndpoint(endpoint *model.Endpoint) {
	if c.HealthCheck != nil {
		c.HealthCheck.StopOne(endpoint)
	}
}

func (c *Cluster) AddEndpoint(endpoint *model.Endpoint) {
	if c.HealthCheck != nil {
		c.HealthCheck.StartOne(endpoint)
	}
}

// LegacyPickLock returns the runtime-scoped lock for legacy load balancer picks.
func (c *Cluster) LegacyPickLock() sync.Locker {
	if c == nil {
		return nil
	}
	return &c.legacyPickMu
}

func (c *Cluster) EndpointSnapshot() *EndpointSnapshot {
	snapshot := c.endpoints.Load()
	if snapshot == nil {
		return emptyEndpointSnapshot
	}
	return snapshot
}

func (c *Cluster) RefreshEndpoints() {
	c.RefreshEndpointsFrom(c.endpoints.Load())
}

func (c *Cluster) RefreshEndpointsFrom(previous *EndpointSnapshot) {
	for {
		current := c.endpoints.Load()
		source := previous
		if current != nil && current != previous {
			source = current
		}
		next := newEndpointSnapshot(c.Config, source, len(c.Config.HealthChecks) != 0)
		if c.endpoints.CompareAndSwap(current, next) {
			return
		}
	}
}

func (c *Cluster) UpdateEndpointHealth(endpointID, endpointAddress string, healthy bool) bool {
	c.healthMu.Lock()
	defer c.healthMu.Unlock()
	if !c.acceptHealthEvents {
		return false
	}

	for {
		current := c.endpoints.Load()
		if current == nil {
			return false
		}
		next, ok := current.withEndpointHealth(endpointID, endpointAddress, healthy)
		if !ok {
			return false
		}
		if next == current {
			return true
		}
		if c.endpoints.CompareAndSwap(current, next) {
			return true
		}
	}
}

// UpdateEndpointAddressHealth updates every endpoint that shares an address.
// Address-keyed health checkers use this path because one probe represents the
// reachability of all endpoint identities at the same network address.
func (c *Cluster) UpdateEndpointAddressHealth(endpointAddress string, healthy bool) bool {
	c.healthMu.Lock()
	defer c.healthMu.Unlock()
	if !c.acceptHealthEvents {
		return false
	}

	for {
		current := c.endpoints.Load()
		if current == nil {
			return false
		}
		next, ok := current.withEndpointAddressHealth(endpointAddress, healthy)
		if !ok {
			return false
		}
		if next == current {
			return true
		}
		if c.endpoints.CompareAndSwap(current, next) {
			return true
		}
	}
}

// SnapshotForRuntimeReplacement freezes health updates on this runtime and
// returns the final snapshot that a replacement runtime should inherit.
func (c *Cluster) SnapshotForRuntimeReplacement() *EndpointSnapshot {
	if c == nil {
		return nil
	}
	c.healthMu.Lock()
	defer c.healthMu.Unlock()
	c.acceptHealthEvents = false
	return c.EndpointSnapshot()
}

func (c *Cluster) handleEndpointHealth(event healthcheck.EndpointHealthEvent) {
	c.UpdateEndpointAddressHealth(event.EndpointAddress, event.Healthy)
}

// EndpointSnapshot endpoint membership and health indexes are immutable after
// publication. Snapshot endpoints are cloned from config endpoints when the
// snapshot is built, so request-path health changes do not mutate config
// Endpoint.UnHealthy or config metadata.
type EndpointSnapshot struct {
	all                  []*model.Endpoint
	healthy              []*model.Endpoint
	endpointByID         map[string]*model.Endpoint
	healthyEndpointByID  map[string]*model.Endpoint
	addressByID          map[string]string
	healthyByID          map[string]bool
	healthyByAddress     map[string]bool
	lbPolicy             model.LbPolicyType
	consistentHashOnce   sync.Once
	consistentHash       model.LbConsistentHashView
	consistentHashConfig model.ConsistentHash
}

var emptyEndpointSnapshot = &EndpointSnapshot{
	all:                 []*model.Endpoint{},
	healthy:             []*model.Endpoint{},
	endpointByID:        map[string]*model.Endpoint{},
	healthyEndpointByID: map[string]*model.Endpoint{},
	addressByID:         map[string]string{},
	healthyByID:         map[string]bool{},
	healthyByAddress:    map[string]bool{},
}

func newEndpointSnapshot(config *model.ClusterConfig, previous *EndpointSnapshot, inheritRuntimeHealth bool) *EndpointSnapshot {
	var endpoints []*model.Endpoint
	clusterName := ""
	if config != nil {
		clusterName = config.Name
		endpoints = config.Endpoints
	}
	snapshot := newEndpointSnapshotIndex(len(endpoints))
	if config != nil {
		snapshot.lbPolicy = config.LbStr
		snapshot.consistentHashConfig = config.ConsistentHash
		snapshot.consistentHashConfig.Hash = nil
	}
	endpointIDs := make(map[string]struct{}, len(endpoints))
	for _, endpoint := range endpoints {
		if endpoint == nil {
			snapshot.all = append(snapshot.all, nil)
			continue
		}
		snapshotEndpoint := model.CloneEndpoint(endpoint)
		snapshotEndpoint.ID = uniqueSnapshotEndpointID(clusterName, snapshotEndpoint, endpointIDs)
		endpointIDs[snapshotEndpoint.ID] = struct{}{}
		address := snapshotEndpoint.Address.GetAddress()
		healthy := endpointSnapshotHealth(snapshotEndpoint, address, previous, inheritRuntimeHealth)
		snapshot.addEndpoint(snapshotEndpoint, address, healthy)
	}
	return snapshot
}

func newEndpointSnapshotIndex(endpointCount int) *EndpointSnapshot {
	return &EndpointSnapshot{
		all:                 make([]*model.Endpoint, 0, endpointCount),
		healthy:             make([]*model.Endpoint, 0, endpointCount),
		endpointByID:        make(map[string]*model.Endpoint, endpointCount),
		healthyEndpointByID: make(map[string]*model.Endpoint, endpointCount),
		addressByID:         make(map[string]string, endpointCount),
		healthyByID:         make(map[string]bool, endpointCount),
		healthyByAddress:    make(map[string]bool, endpointCount),
	}
}

func uniqueSnapshotEndpointID(clusterName string, endpoint *model.Endpoint, endpointIDs map[string]struct{}) string {
	id := ""
	if endpoint != nil {
		id = endpoint.ID
	}
	if id == "" {
		id = model.GeneratedEndpointID(clusterName, endpoint)
	}
	if _, exists := endpointIDs[id]; !exists {
		return id
	}
	baseID := model.GeneratedEndpointID(clusterName, endpoint)
	for suffix := 2; ; suffix++ {
		candidate := fmt.Sprintf("%s-%d", baseID, suffix)
		if _, exists := endpointIDs[candidate]; !exists {
			return candidate
		}
	}
}

func endpointSnapshotHealth(
	endpoint *model.Endpoint,
	address string,
	previous *EndpointSnapshot,
	inheritRuntimeHealth bool,
) bool {
	healthy := !endpoint.UnHealthy
	if previous == nil {
		return healthy
	}
	// Carry health only while this runtime still has a health checker that can
	// later correct it; otherwise seed health from config.
	if inheritRuntimeHealth {
		if previousAddress, ok := previous.addressByID[endpoint.ID]; ok && previousAddress == address {
			return previous.healthyByID[endpoint.ID]
		}
		if addressHealthy, ok := previous.healthyByAddress[address]; ok {
			return addressHealthy
		}
	}
	return healthy
}

func (s *EndpointSnapshot) addEndpoint(
	endpoint *model.Endpoint,
	address string,
	healthy bool,
) {
	endpoint.UnHealthy = !healthy
	s.all = append(s.all, endpoint)
	s.endpointByID[endpoint.ID] = endpoint
	s.addressByID[endpoint.ID] = address
	s.healthyByID[endpoint.ID] = healthy
	s.addAddressHealth(address, healthy)
	if healthy {
		s.healthy = append(s.healthy, endpoint)
		s.healthyEndpointByID[endpoint.ID] = endpoint
	}
}

func (s *EndpointSnapshot) addAddressHealth(address string, healthy bool) {
	if s == nil {
		return
	}
	addressHealthy, ok := s.healthyByAddress[address]
	if !ok {
		s.healthyByAddress[address] = healthy
		return
	}
	s.healthyByAddress[address] = addressHealthy && healthy
}

func (s *EndpointSnapshot) AllEndpoints() []*model.Endpoint {
	if s == nil {
		return nil
	}
	return model.CloneEndpoints(s.all)
}

func (s *EndpointSnapshot) EndpointCount() int {
	if s == nil {
		return 0
	}
	return len(s.all)
}

func (s *EndpointSnapshot) HealthyEndpoints() []*model.Endpoint {
	if s == nil {
		return nil
	}
	return model.CloneEndpoints(s.healthy)
}

func (s *EndpointSnapshot) HealthyEndpointCount() int {
	if s == nil {
		return 0
	}
	return len(s.healthy)
}

func (s *EndpointSnapshot) HealthyConsistentHash() model.LbConsistentHashView {
	if s == nil {
		return nil
	}
	s.consistentHashOnce.Do(s.rebuildConsistentHash)
	return s.consistentHash
}

// PickHealthyEndpoint gives request-path selectors a read-only view of healthy
// endpoints and clones only the selected endpoint before returning it. The pick
// callback must not mutate or retain the endpoint view.
func (s *EndpointSnapshot) PickHealthyEndpoint(pick func([]*model.Endpoint) *model.Endpoint) *model.Endpoint {
	if s == nil || len(s.healthy) == 0 || pick == nil {
		return nil
	}
	return model.CloneEndpoint(pick(s.healthy))
}

func (s *EndpointSnapshot) NextHealthyEndpoint(curEndpointID string) *model.Endpoint {
	if s == nil {
		return nil
	}
	start := s.nextEndpointStartIndex(curEndpointID)
	if start < 0 {
		return nil
	}
	for _, endpoint := range s.all[start:] {
		if endpoint == nil {
			continue
		}
		if s.healthyByID[endpoint.ID] {
			return model.CloneEndpoint(endpoint)
		}
	}
	return nil
}

func (s *EndpointSnapshot) nextEndpointStartIndex(curEndpointID string) int {
	for i, endpoint := range s.all {
		if endpoint != nil && endpoint.ID == curEndpointID {
			return i + 1
		}
	}
	return -1
}

func (s *EndpointSnapshot) EndpointByID(endpointID string) *model.Endpoint {
	if s == nil {
		return nil
	}
	return model.CloneEndpoint(s.endpointByID[endpointID])
}

func (s *EndpointSnapshot) HealthyEndpointByID(endpointID string) *model.Endpoint {
	if s == nil {
		return nil
	}
	return model.CloneEndpoint(s.healthyEndpointByID[endpointID])
}

func (s *EndpointSnapshot) withEndpointHealth(
	endpointID, endpointAddress string,
	healthy bool,
) (*EndpointSnapshot, bool) {
	if s == nil {
		return nil, false
	}
	address, ok := s.addressByID[endpointID]
	if !ok || address != endpointAddress {
		return nil, false
	}
	return s.withEndpointHealthForIDs(map[string]struct{}{endpointID: {}}, healthy)
}

func (s *EndpointSnapshot) withEndpointAddressHealth(
	endpointAddress string,
	healthy bool,
) (*EndpointSnapshot, bool) {
	if s == nil {
		return nil, false
	}
	endpointIDs := make(map[string]struct{})
	for endpointID, address := range s.addressByID {
		if address == endpointAddress {
			endpointIDs[endpointID] = struct{}{}
		}
	}
	return s.withEndpointHealthForIDs(endpointIDs, healthy)
}

func (s *EndpointSnapshot) withEndpointHealthForIDs(
	endpointIDs map[string]struct{},
	healthy bool,
) (*EndpointSnapshot, bool) {
	if len(endpointIDs) == 0 {
		return nil, false
	}

	changed := false
	for endpointID := range endpointIDs {
		if s.healthyByID[endpointID] != healthy {
			changed = true
			break
		}
	}
	if !changed {
		return s, true
	}

	// Reuse unchanged endpoint/address views and rebuild the endpoint clone
	// whose runtime health flag changed.
	next := &EndpointSnapshot{
		all:                  make([]*model.Endpoint, 0, len(s.all)),
		healthy:              make([]*model.Endpoint, 0, len(s.all)),
		endpointByID:         make(map[string]*model.Endpoint, len(s.endpointByID)),
		healthyEndpointByID:  make(map[string]*model.Endpoint, len(s.endpointByID)),
		addressByID:          s.addressByID,
		healthyByID:          make(map[string]bool, len(s.healthyByID)),
		healthyByAddress:     make(map[string]bool, len(s.healthyByAddress)),
		lbPolicy:             s.lbPolicy,
		consistentHashConfig: s.consistentHashConfig,
	}

	for id, wasHealthy := range s.healthyByID {
		nextHealthy := wasHealthy
		if _, ok := endpointIDs[id]; ok {
			nextHealthy = healthy
		}
		next.healthyByID[id] = nextHealthy
		next.addAddressHealth(s.addressByID[id], nextHealthy)
	}

	for _, endpoint := range s.all {
		if endpoint == nil {
			next.all = append(next.all, nil)
			continue
		}
		id := endpoint.ID
		nextEndpoint := endpoint
		if _, ok := endpointIDs[id]; ok {
			nextEndpoint = model.CloneEndpoint(endpoint)
			nextEndpoint.UnHealthy = !healthy
		}
		next.all = append(next.all, nextEndpoint)
		next.endpointByID[id] = nextEndpoint
		if next.healthyByID[id] {
			next.healthy = append(next.healthy, nextEndpoint)
			next.healthyEndpointByID[id] = nextEndpoint
		}
	}

	return next, true
}

func (s *EndpointSnapshot) rebuildConsistentHash() {
	if s == nil || len(s.healthy) == 0 {
		return
	}
	newConsistentHash, ok := model.ConsistentHashInitMap[s.lbPolicy]
	if !ok {
		return
	}
	s.consistentHash = model.ReadOnlyConsistentHash(newConsistentHash(s.consistentHashConfig, s.healthy))
}

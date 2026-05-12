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
		next := newEndpointSnapshot(c.Config.Endpoints, source, len(c.Config.HealthChecks) != 0)
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
	c.UpdateEndpointHealth(event.EndpointID, event.EndpointAddress, event.Healthy)
}

// EndpointSnapshot endpoint membership and health indexes are immutable after
// publication. Snapshot endpoints are cloned from config endpoints when the
// snapshot is built, so request-path health changes do not mutate config
// Endpoint.UnHealthy or config metadata.
type EndpointSnapshot struct {
	all                 []*model.Endpoint
	healthy             []*model.Endpoint
	endpointByID        map[string]*model.Endpoint
	healthyEndpointByID map[string]*model.Endpoint
	addressByID         map[string]string
	healthyByID         map[string]bool
}

var emptyEndpointSnapshot = &EndpointSnapshot{
	all:                 []*model.Endpoint{},
	healthy:             []*model.Endpoint{},
	endpointByID:        map[string]*model.Endpoint{},
	healthyEndpointByID: map[string]*model.Endpoint{},
	addressByID:         map[string]string{},
	healthyByID:         map[string]bool{},
}

func newEndpointSnapshot(endpoints []*model.Endpoint, previous *EndpointSnapshot, inheritRuntimeHealth bool) *EndpointSnapshot {
	snapshot := newEndpointSnapshotIndex(len(endpoints))
	for _, endpoint := range endpoints {
		if endpoint == nil {
			snapshot.all = append(snapshot.all, nil)
			continue
		}
		snapshotEndpoint := cloneEndpoint(endpoint)
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
	previousAddress, ok := previous.addressByID[endpoint.ID]
	if !ok || previousAddress != address {
		return healthy
	}

	// Carry health only while this runtime still has a health checker that can
	// later correct it; otherwise seed health from config.
	if inheritRuntimeHealth {
		healthy = previous.healthyByID[endpoint.ID]
	}
	return healthy
}

func (s *EndpointSnapshot) addEndpoint(
	endpoint *model.Endpoint,
	address string,
	healthy bool,
) {
	s.all = append(s.all, endpoint)
	s.endpointByID[endpoint.ID] = endpoint
	s.addressByID[endpoint.ID] = address
	s.healthyByID[endpoint.ID] = healthy
	if healthy {
		s.healthy = append(s.healthy, endpoint)
		s.healthyEndpointByID[endpoint.ID] = endpoint
	}
}

func (s *EndpointSnapshot) AllEndpoints() []*model.Endpoint {
	if s == nil {
		return nil
	}
	return cloneEndpoints(s.all)
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
	return cloneEndpoints(s.healthy)
}

func (s *EndpointSnapshot) EndpointByID(endpointID string) *model.Endpoint {
	if s == nil {
		return nil
	}
	return s.endpointByID[endpointID]
}

func (s *EndpointSnapshot) HealthyEndpointByID(endpointID string) *model.Endpoint {
	if s == nil {
		return nil
	}
	return s.healthyEndpointByID[endpointID]
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
	if s.healthyByID[endpointID] == healthy {
		return s, true
	}

	// Reuse the immutable endpoint/address views and rebuild only the health
	// views that change for this event.
	next := &EndpointSnapshot{
		all:                 s.all,
		healthy:             make([]*model.Endpoint, 0, len(s.all)),
		endpointByID:        s.endpointByID,
		healthyEndpointByID: make(map[string]*model.Endpoint, len(s.endpointByID)),
		addressByID:         s.addressByID,
		healthyByID:         make(map[string]bool, len(s.healthyByID)),
	}

	for id, wasHealthy := range s.healthyByID {
		nextHealthy := wasHealthy
		if id == endpointID {
			nextHealthy = healthy
		}
		next.healthyByID[id] = nextHealthy
	}

	for _, endpoint := range s.all {
		if endpoint == nil {
			continue
		}
		id := endpoint.ID
		if next.healthyByID[id] {
			next.healthy = append(next.healthy, endpoint)
			next.healthyEndpointByID[id] = endpoint
		}
	}

	return next, true
}

func cloneEndpoints(endpoints []*model.Endpoint) []*model.Endpoint {
	if endpoints == nil {
		return nil
	}
	cloned := make([]*model.Endpoint, len(endpoints))
	copy(cloned, endpoints)
	return cloned
}

func cloneEndpoint(endpoint *model.Endpoint) *model.Endpoint {
	if endpoint == nil {
		return nil
	}
	cloned := *endpoint
	cloned.Address = cloneSocketAddress(endpoint.Address)
	cloned.Metadata = cloneMetadata(endpoint.Metadata)
	cloned.LLMMeta = cloneLLMMeta(endpoint.LLMMeta)
	return &cloned
}

func cloneSocketAddress(address model.SocketAddress) model.SocketAddress {
	cloned := address
	if address.Domains != nil {
		cloned.Domains = append([]string(nil), address.Domains...)
	}
	return cloned
}

func cloneMetadata(metadata map[string]string) map[string]string {
	if metadata == nil {
		return nil
	}
	cloned := make(map[string]string, len(metadata))
	for key, value := range metadata {
		cloned[key] = value
	}
	return cloned
}

func cloneLLMMeta(meta *model.LLMMeta) *model.LLMMeta {
	if meta == nil {
		return nil
	}
	cloned := *meta
	if meta.RetryPolicy.Config != nil {
		cloned.RetryPolicy.Config = make(map[string]any, len(meta.RetryPolicy.Config))
		for key, value := range meta.RetryPolicy.Config {
			cloned.RetryPolicy.Config[key] = value
		}
	}
	return &cloned
}

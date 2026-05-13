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

package loadbalancer

import (
	"sync"
	"sync/atomic"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type PickContext struct {
	// Config carries cluster-level load-balancer configuration and runtime
	// cursor state. Snapshot-aware balancers should not reread Config.Endpoints
	// for health filtering.
	Config *model.ClusterConfig
	// HealthyConsistentHash is built from HealthyEndpoints for the same
	// immutable runtime snapshot and intentionally exposes lookup methods only.
	// Consistent-hash balancers should prefer this view over
	// Config.ConsistentHash.Hash, which is mutable and can include
	// runtime-unhealthy endpoints.
	HealthyConsistentHash model.LbConsistentHashView
	// AllEndpoints is the current runtime snapshot, including endpoints marked
	// unhealthy by runtime health checks.
	AllEndpoints []*model.Endpoint
	// HealthyEndpoints is already filtered from the current runtime snapshot.
	// Snapshot-aware balancers must treat endpoints as read-only and return the
	// chosen endpoint without mutating or retaining it.
	HealthyEndpoints []*model.Endpoint
}

type LoadBalancer interface {
	Handler(c *model.ClusterConfig, policy model.LbPolicy) *model.Endpoint
}

// SnapshotLoadBalancer is optional for balancers that consume the current
// runtime snapshot. Implementations should pick only from
// PickContext.HealthyEndpoints; Config.Endpoints may include unhealthy or stale
// config entries.
type SnapshotLoadBalancer interface {
	HandlerWithSnapshot(c PickContext, policy model.LbPolicy) *model.Endpoint
}

// HealthyOnlySnapshotLoadBalancer marks snapshot-aware balancers that do not
// need PickContext.AllEndpoints. Unmarked snapshot balancers keep receiving the
// full snapshot for compatibility with custom implementations.
type HealthyOnlySnapshotLoadBalancer interface {
	UseHealthyEndpointsOnly() bool
}

// ZeroCopySnapshotLoadBalancer marks trusted balancers that never mutate or
// retain snapshot endpoints. Other snapshot balancers receive defensive copies.
type ZeroCopySnapshotLoadBalancer interface {
	UseZeroCopySnapshot() bool
}

// ClusterScopedLegacyLoadBalancer lets a legacy load balancer opt in to
// runtime-cluster scoped serialization. Legacy balancers that do not implement
// this interface keep the package-level compatibility lock because strategy
// instances are shared globally. Implementations must ensure the same balancer
// instance can run Handler concurrently across different clusters.
type ClusterScopedLegacyLoadBalancer interface {
	UseClusterScopedLegacyLock() bool
}

// LoadBalancerStrategy load balancer strategy mode
var LoadBalancerStrategy = map[model.LbPolicyType]LoadBalancer{}

var legacyPickMu sync.Mutex

func RegisterLoadBalancer(name model.LbPolicyType, balancer LoadBalancer) {
	if _, ok := LoadBalancerStrategy[name]; ok {
		panic("load balancer register fail " + name)
	}
	LoadBalancerStrategy[name] = balancer
}

func PickEndpoint(balancer LoadBalancer, context PickContext, policy model.LbPolicy) *model.Endpoint {
	return pickEndpointWithLegacyLock(balancer, nil, context, policy)
}

// PickEndpointWithLegacyLock serializes legacy balancers. The caller-provided
// runtime lock is used only when the balancer explicitly opts in to scoped
// serialization; other legacy balancers keep the package-level compatibility
// lock.
func PickEndpointWithLegacyLock(balancer LoadBalancer, legacyPickLock sync.Locker, context PickContext, policy model.LbPolicy) *model.Endpoint {
	return pickEndpointWithLegacyLock(balancer, legacyPickLock, context, policy)
}

// NeedsAllEndpoints reports whether a snapshot-aware balancer should receive
// PickContext.AllEndpoints on the request path.
func NeedsAllEndpoints(balancer LoadBalancer) bool {
	healthyOnly, ok := balancer.(HealthyOnlySnapshotLoadBalancer)
	return !ok || !healthyOnly.UseHealthyEndpointsOnly()
}

// ConsistentHashForHealthyEndpoints returns a consistent hash view that only
// contains the healthy endpoints visible to this pick.
func ConsistentHashForHealthyEndpoints(context PickContext) model.LbConsistentHashView {
	if context.HealthyConsistentHash != nil {
		return context.HealthyConsistentHash
	}
	if context.Config == nil || len(context.HealthyEndpoints) == 0 {
		return nil
	}
	newConsistentHash, ok := model.ConsistentHashInitMap[context.Config.LbStr]
	if ok {
		return model.ReadOnlyConsistentHash(newConsistentHash(context.Config.ConsistentHash, context.HealthyEndpoints))
	}
	return model.ReadOnlyConsistentHash(context.Config.ConsistentHash.Hash)
}

func pickEndpointWithLegacyLock(balancer LoadBalancer, legacyPickLock sync.Locker, context PickContext, policy model.LbPolicy) *model.Endpoint {
	if balancer == nil || context.Config == nil {
		return nil
	}
	if snapshotBalancer, ok := balancer.(SnapshotLoadBalancer); ok {
		snapshotContext := context
		zeroCopy, ok := balancer.(ZeroCopySnapshotLoadBalancer)
		if !ok || !zeroCopy.UseZeroCopySnapshot() {
			snapshotContext = defensiveSnapshotPickContext(context)
		}
		endpoint := snapshotBalancer.HandlerWithSnapshot(snapshotContext, policy)
		return healthyEndpointFromSnapshot(endpoint, context.HealthyEndpoints)
	}

	// Legacy balancers only understand ClusterConfig. Serialize this
	// compatibility path so cursor-style state reconciles predictably; custom
	// mutable state should move to SnapshotLoadBalancer instead.
	lock := legacyPickLockFor(balancer, legacyPickLock)
	lock.Lock()
	defer lock.Unlock()

	allEndpoints := context.AllEndpoints
	if allEndpoints == nil {
		allEndpoints = context.HealthyEndpoints
	}
	config := *context.Config
	config.Endpoints = model.CloneEndpoints(allEndpoints)
	cursorBefore := atomic.LoadUint32(&context.Config.PrePickEndpointIndex)
	atomic.StoreUint32(&config.PrePickEndpointIndex, cursorBefore)
	endpoint := balancer.Handler(&config, policy)
	cursorAfter := atomic.LoadUint32(&config.PrePickEndpointIndex)
	if cursorAfter != cursorBefore {
		atomic.AddUint32(&context.Config.PrePickEndpointIndex, cursorAfter-cursorBefore)
	}
	return healthyEndpointFromSnapshot(endpoint, context.HealthyEndpoints)
}

func defensiveSnapshotPickContext(context PickContext) PickContext {
	defensive := context
	defensive.AllEndpoints = model.CloneEndpoints(context.AllEndpoints)
	defensive.HealthyEndpoints = model.CloneEndpoints(context.HealthyEndpoints)
	return defensive
}

func legacyPickLockFor(balancer LoadBalancer, legacyPickLock sync.Locker) sync.Locker {
	if scoped, ok := balancer.(ClusterScopedLegacyLoadBalancer); ok && scoped.UseClusterScopedLegacyLock() && legacyPickLock != nil {
		return legacyPickLock
	}
	return &legacyPickMu
}

func healthyEndpointFromSnapshot(endpoint *model.Endpoint, healthyEndpoints []*model.Endpoint) *model.Endpoint {
	if endpoint == nil {
		return nil
	}
	for _, candidate := range healthyEndpoints {
		if sameEndpointIdentity(candidate, endpoint) {
			return model.CloneEndpoint(candidate)
		}
	}
	return nil
}

func sameEndpointIdentity(candidate, endpoint *model.Endpoint) bool {
	if candidate == nil || endpoint == nil {
		return false
	}
	if endpoint.ID != "" || candidate.ID != "" {
		if candidate.ID != endpoint.ID {
			return false
		}
		endpointAddress := endpoint.Address.GetAddress()
		return endpointAddress == "" || candidate.Address.GetAddress() == endpointAddress
	}
	return candidate.Address.GetAddress() == endpoint.Address.GetAddress()
}

func RegisterConsistentHashInit(name model.LbPolicyType, function model.ConsistentHashInitFunc) {
	if _, ok := model.ConsistentHashInitMap[name]; ok {
		panic("consistent hash load balancer register fail " + name)
	}
	model.ConsistentHashInitMap[name] = function
}

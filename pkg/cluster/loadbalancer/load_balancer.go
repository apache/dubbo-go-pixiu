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

// PickContext bundles the snapshot view a load balancer needs to make a
// single pick. The runtime publishes an immutable EndpointSnapshot per
// cluster and derives PickContext from it, so balancers must treat every
// field as read-only.
type PickContext struct {
	// Config carries cluster-level load-balancer configuration and runtime
	// cursor state. Snapshot-aware balancers should not reread
	// Config.Endpoints for health filtering — use HealthyEndpoints instead.
	Config *model.ClusterConfig
	// HealthyConsistentHash is built from HealthyEndpoints for the same
	// immutable runtime snapshot and intentionally exposes lookup methods only.
	// Consistent-hash balancers should prefer this view over
	// Config.ConsistentHash.Hash, which is mutable and can include
	// runtime-unhealthy endpoints.
	HealthyConsistentHash model.LbConsistentHashView
	// AllEndpoints is the current runtime snapshot, including endpoints marked
	// unhealthy by runtime health checks. Provided only when the balancer
	// requests it via NeedsAllEndpoints.
	AllEndpoints []*model.Endpoint
	// HealthyEndpoints is already filtered from the current runtime snapshot.
	// Snapshot-aware balancers must treat endpoints as read-only and return
	// the chosen endpoint without mutating or retaining it.
	HealthyEndpoints []*model.Endpoint
}

type LoadBalancer interface {
	Handler(c *model.ClusterConfig, policy model.LbPolicy) *model.Endpoint
}

// SnapshotLoadBalancer is optional for balancers that consume the current
// runtime snapshot. Implementations should pick only from
// PickContext.HealthyEndpoints; Config.Endpoints may include unhealthy or
// stale config entries.
type SnapshotLoadBalancer interface {
	HandlerWithSnapshot(c PickContext, policy model.LbPolicy) *model.Endpoint
}

// HealthyOnlySnapshotLoadBalancer marks snapshot-aware balancers that do not
// need PickContext.AllEndpoints. Unmarked snapshot balancers keep receiving
// the full snapshot for compatibility with custom implementations.
type HealthyOnlySnapshotLoadBalancer interface {
	UseHealthyEndpointsOnly() bool
}

// ZeroCopySnapshotLoadBalancer marks trusted balancers that never mutate or
// retain snapshot endpoints. Other snapshot balancers receive defensive
// copies.
type ZeroCopySnapshotLoadBalancer interface {
	UseZeroCopySnapshot() bool
}

// LoadBalancerStrategy load balancer strategy mode
var LoadBalancerStrategy = map[model.LbPolicyType]LoadBalancer{}

// legacyPickMu serializes pre-snapshot balancers. They share strategy
// instances globally and may carry mutable cursor state, so concurrent
// Handler calls across clusters are unsafe by default.
var legacyPickMu sync.Mutex

func RegisterLoadBalancer(name model.LbPolicyType, balancer LoadBalancer) {
	if _, ok := LoadBalancerStrategy[name]; ok {
		panic("load balancer register fail " + name)
	}
	LoadBalancerStrategy[name] = balancer
}

// PickEndpoint picks an endpoint for the supplied snapshot context. Use this
// from snapshot-published pick paths. Legacy balancers fall back to the
// package-level compatibility lock automatically.
func PickEndpoint(balancer LoadBalancer, context PickContext, policy model.LbPolicy) *model.Endpoint {
	return pickEndpoint(balancer, context, policy)
}

// NeedsAllEndpoints reports whether a snapshot-aware balancer should receive
// PickContext.AllEndpoints on the request path.
func NeedsAllEndpoints(balancer LoadBalancer) bool {
	healthyOnly, ok := balancer.(HealthyOnlySnapshotLoadBalancer)
	return !ok || !healthyOnly.UseHealthyEndpointsOnly()
}

// ConsistentHashForHealthyEndpoints returns a consistent hash view that only
// contains the healthy endpoints visible to this pick. Returns nil if the
// context has no healthy endpoints or no consistent-hash factory registered.
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

func pickEndpoint(balancer LoadBalancer, context PickContext, policy model.LbPolicy) *model.Endpoint {
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
	// compatibility path so cursor-style state reconciles predictably;
	// custom mutable state should move to SnapshotLoadBalancer instead.
	//
	// RR cursor fairness caveat (issue #905): the cursor delta-apply
	// below uses atomic.AddUint32 on context.Config.PrePickEndpointIndex.
	// If a snapshot-path RoundRobin picks against the SAME *ClusterConfig
	// concurrently (also via atomic.AddUint32), the two paths compete on
	// the same counter and the legacy delta may overshoot or undershoot
	// the snapshot increment. This is not a data race (both are atomic)
	// and not a correctness regression (the picked endpoint is still
	// valid), but the visit order may skip / repeat one slot during the
	// race window. In practice a single cluster registers a single
	// LbStr, so legacy and snapshot RR do not coexist; if a future
	// deployment mixes them, accept the fairness skew or migrate the
	// legacy plugin to SnapshotLoadBalancer.
	//
	// Shared-pointer hazard: config := *context.Config is a shallow copy.
	// config.Endpoints is replaced with a clone below, but other
	// pointer-bearing fields on ClusterConfig (ConsistentHash.Hash,
	// operator-supplied Metadata, etc.) stay shared with context.Config.
	// The package lock above serializes legacy picks, but it does not
	// serialize legacy picks against concurrent ClusterStore mutations that
	// touch those fields. In-tree this is safe because all ClusterStore
	// mutators hold ClusterManager.rw and PickEndpoint takes
	// ClusterManager.rw.RLock; external legacy plugins must observe the same
	// rule.
	legacyPickMu.Lock()
	defer legacyPickMu.Unlock()

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

func healthyEndpointFromSnapshot(endpoint *model.Endpoint, healthyEndpoints []*model.Endpoint) *model.Endpoint {
	if endpoint == nil {
		return nil
	}
	for _, candidate := range healthyEndpoints {
		if candidate == endpoint {
			return model.CloneEndpoint(candidate)
		}
	}
	for _, candidate := range healthyEndpoints {
		if sameEndpointIdentity(candidate, endpoint) {
			return model.CloneEndpoint(candidate)
		}
	}
	return nil
}

// sameEndpointIdentity reports whether candidate (from the snapshot's healthy
// set) matches endpoint (returned by a balancer pick). Used to confirm the
// balancer chose a still-healthy entry before handing it back to the request
// path.
//
// Parameter convention (load-bearing): the FIRST argument is the snapshot
// healthy-set member, the SECOND is the balancer return. The wildcard in
// rule (2) is keyed off the SNAPSHOT side because that is where legacy
// resolvers produce placeholder endpoints. Crossing the arguments will
// silently invert rule (2) and reject valid picks that should match by ID.
// Call-site convention enforced by healthyEndpointFromSnapshot:
//
//	sameEndpointIdentity(candidate /* snapshot */, endpoint /* return */)
//
// Rules:
//
//  1. When either side carries a non-empty ID, IDs must match. ID is the
//     authoritative identity since PR-2 (deterministic generation) and is
//     immune to address drift caused by DNS or service-discovery refreshes.
//  2. With IDs matched, when the SNAPSHOT endpoint's (candidate) address is
//     fully empty (Domains == [""] AND Address == "" AND Port == 0 — i.e.
//     SocketAddress is at its zero value bar the single blank domain
//     marker), any balancer return address is accepted. This is a narrow
//     wildcard for legacy resolvers that publish placeholder endpoints
//     whose address is intentionally absent and is reconciled via ID; the
//     balancer (or a downstream resolver) supplies the resolved address.
//     Any non-zero address field on the snapshot side forfeits the
//     wildcard so an operator cannot accidentally widen the trust
//     boundary by typing one blank domain entry next to a real port.
//  3. Otherwise (or when neither side has an ID), addresses must compare
//     equal via SocketAddress.Equal — no string formatting, no allocation.
//
// Risk: rule (2) means anyone with the right ID matches the placeholder
// endpoint regardless of where the balancer routed them. The ID is
// therefore treated as a trust boundary and must remain operator-controlled
// or system-generated. The fully-empty-address gate keeps the wildcard
// from catching real snapshot addresses that happen to share an ID via
// misconfiguration.
//
// Locked by TestSameEndpointIdentityBlankDomainWildcard (unit) and
// TestHealthyEndpointFromSnapshotAcceptsResolvedAddressForBlankPlaceholder
// (integration through the real pick path) in load_balancer_test.go.
func sameEndpointIdentity(candidate, endpoint *model.Endpoint) bool {
	if candidate == nil || endpoint == nil {
		return false
	}
	if endpoint.ID != "" || candidate.ID != "" {
		if candidate.ID != endpoint.ID {
			return false
		}
		if isBlankDomainPlaceholderAddress(candidate.Address) {
			return true
		}
		return candidate.Address.Equal(endpoint.Address)
	}
	return candidate.Address.Equal(endpoint.Address)
}

// isBlankDomainPlaceholderAddress reports whether addr is the narrow
// "address-absent placeholder" form: exactly one blank domain entry and
// zero values everywhere else. Used by sameEndpointIdentity to bound the
// blank-domain wildcard.
func isBlankDomainPlaceholderAddress(addr model.SocketAddress) bool {
	return len(addr.Domains) == 1 &&
		addr.Domains[0] == "" &&
		addr.Address == "" &&
		addr.Port == 0
}

func RegisterConsistentHashInit(name model.LbPolicyType, function model.ConsistentHashInitFunc) {
	if _, ok := model.ConsistentHashInitMap[name]; ok {
		panic("consistent hash load balancer register fail " + name)
	}
	model.ConsistentHashInitMap[name] = function
}

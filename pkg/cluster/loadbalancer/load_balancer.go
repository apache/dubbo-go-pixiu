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
	// HealthyEndpoints is already filtered from the current runtime snapshot.
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
	if balancer == nil || context.Config == nil {
		return nil
	}
	if snapshotBalancer, ok := balancer.(SnapshotLoadBalancer); ok {
		return snapshotBalancer.HandlerWithSnapshot(context, policy)
	}

	// Legacy balancers only understand ClusterConfig. Serialize this
	// compatibility path so cursor-style state reconciles predictably; custom
	// mutable state should move to SnapshotLoadBalancer instead.
	legacyPickMu.Lock()
	defer legacyPickMu.Unlock()

	config := *context.Config
	config.Endpoints = cloneEndpoints(context.HealthyEndpoints)
	cursorBefore := atomic.LoadUint32(&context.Config.PrePickEndpointIndex)
	atomic.StoreUint32(&config.PrePickEndpointIndex, cursorBefore)
	endpoint := balancer.Handler(&config, policy)
	cursorAfter := atomic.LoadUint32(&config.PrePickEndpointIndex)
	if cursorAfter != cursorBefore {
		atomic.AddUint32(&context.Config.PrePickEndpointIndex, cursorAfter-cursorBefore)
	}
	return endpoint
}

func cloneEndpoints(endpoints []*model.Endpoint) []*model.Endpoint {
	if endpoints == nil {
		return nil
	}
	cloned := make([]*model.Endpoint, len(endpoints))
	copy(cloned, endpoints)
	return cloned
}

func RegisterConsistentHashInit(name model.LbPolicyType, function model.ConsistentHashInitFunc) {
	if _, ok := model.ConsistentHashInitMap[name]; ok {
		panic("consistent hash load balancer register fail " + name)
	}
	model.ConsistentHashInitMap[name] = function
}

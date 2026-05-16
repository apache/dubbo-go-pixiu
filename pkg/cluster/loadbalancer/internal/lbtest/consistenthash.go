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

// Package lbtest holds load-balancer test helpers shared between
// pkg/cluster/loadbalancer sub-packages.
package lbtest

import (
	"testing"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	testLoopbackAddress = "127.0.0.1"
	pickedWantMsgFormat = "picked %q, want %q"
)

// StaticHashPolicy returns a fixed hash key, satisfying model.LbPolicy.
type StaticHashPolicy string

func (p StaticHashPolicy) GenerateHash() string {
	return string(p)
}

// FixedConsistentHash always resolves to Host, satisfying both
// model.LbConsistentHash and model.LbConsistentHashView for tests.
type FixedConsistentHash struct {
	Host string
}

func (h FixedConsistentHash) Hash(string) uint32 { return 0 }

// Add intentionally does nothing: the fixture exposes a frozen single-host
// view and rejects runtime membership changes by design.
func (h FixedConsistentHash) Add(string) {
	// no-op: see method doc above.
}

func (h FixedConsistentHash) Get(string) (string, error)     { return h.Host, nil }
func (h FixedConsistentHash) GetHash(uint32) (string, error) { return h.Host, nil }
func (h FixedConsistentHash) Remove(string) bool             { return false }

// ConsistentHashBalancer is the union of legacy and snapshot pick entry points
// that consistent-hash balancers (Maglev, RingHash, ...) implement.
type ConsistentHashBalancer interface {
	loadbalancer.LoadBalancer
	loadbalancer.SnapshotLoadBalancer
}

// RunConsistentHashHandlerSuite exercises the configured-hash, unhealthy
// fallback, and snapshot-aware paths for a consistent-hash balancer.
// namePrefix labels the cluster configs so failure output stays readable.
func RunConsistentHashHandlerSuite(t *testing.T, balancer ConsistentHashBalancer, namePrefix string) {
	t.Helper()

	t.Run("HandlerUsesConfiguredHashWithoutClusterLBPolicy", func(t *testing.T) {
		first := &model.Endpoint{
			ID:      "first",
			Address: model.SocketAddress{Address: testLoopbackAddress, Port: 18080},
		}
		second := &model.Endpoint{
			ID:      "second",
			Address: model.SocketAddress{Address: testLoopbackAddress, Port: 18081},
		}
		cluster := &model.ClusterConfig{
			Name:      namePrefix + "-direct",
			Endpoints: []*model.Endpoint{first, second},
			ConsistentHash: model.ConsistentHash{
				Hash: FixedConsistentHash{Host: second.GetHost()},
			},
		}

		got := balancer.Handler(cluster, StaticHashPolicy("request-key"))
		if got == nil {
			t.Fatal("expected endpoint, got nil")
		}
		if got.ID != second.ID {
			t.Fatalf(pickedWantMsgFormat, got.ID, second.ID)
		}
	})

	t.Run("HandlerFallsBackWhenConfiguredHashHitsUnhealthyEndpoint", func(t *testing.T) {
		healthy := &model.Endpoint{
			ID:      "healthy",
			Address: model.SocketAddress{Address: testLoopbackAddress, Port: 18080},
		}
		unhealthy := &model.Endpoint{
			ID:        "unhealthy",
			Address:   model.SocketAddress{Address: testLoopbackAddress, Port: 18081},
			UnHealthy: true,
		}
		cluster := &model.ClusterConfig{
			Name:      namePrefix + "-direct-unhealthy-hit",
			Endpoints: []*model.Endpoint{healthy, unhealthy},
			ConsistentHash: model.ConsistentHash{
				Hash: FixedConsistentHash{Host: unhealthy.GetHost()},
			},
		}

		got := balancer.Handler(cluster, StaticHashPolicy("request-key"))
		if got == nil {
			t.Fatal("expected fallback endpoint, got nil")
		}
		if got.ID != healthy.ID {
			t.Fatalf(pickedWantMsgFormat, got.ID, healthy.ID)
		}
	})

	t.Run("UsesHealthyConsistentHashSnapshot", func(t *testing.T) {
		first := &model.Endpoint{
			ID:      "first",
			Address: model.SocketAddress{Address: testLoopbackAddress, Port: 18080},
		}
		unhealthy := &model.Endpoint{
			ID:        "unhealthy",
			Address:   model.SocketAddress{Address: testLoopbackAddress, Port: 18081},
			UnHealthy: true,
		}
		second := &model.Endpoint{
			ID:      "second",
			Address: model.SocketAddress{Address: testLoopbackAddress, Port: 18082},
		}
		cluster := &model.ClusterConfig{
			Name: namePrefix + "-healthy-snapshot",
			ConsistentHash: model.ConsistentHash{
				Hash: FixedConsistentHash{Host: unhealthy.GetHost()},
			},
		}

		got := balancer.HandlerWithSnapshot(loadbalancer.PickContext{
			Config:                cluster,
			HealthyConsistentHash: FixedConsistentHash{Host: second.GetHost()},
			HealthyEndpoints:      []*model.Endpoint{first, second},
		}, StaticHashPolicy("key-for-unhealthy-slot"))

		if got == nil {
			t.Fatal("expected healthy endpoint, got nil")
		}
		if got.ID != second.ID {
			t.Fatalf(pickedWantMsgFormat, got.ID, second.ID)
		}
	})
}

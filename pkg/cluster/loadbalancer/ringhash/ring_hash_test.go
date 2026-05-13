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

package ringhash

import (
	"fmt"
	stdHttp "net/http"
	"strconv"
	"testing"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type staticHashPolicy string

func (p staticHashPolicy) GenerateHash() string {
	return string(p)
}

type fixedConsistentHash struct {
	host string
}

func (h fixedConsistentHash) Hash(string) uint32 {
	return 0
}

func (h fixedConsistentHash) Add(string) {}

func (h fixedConsistentHash) Get(string) (string, error) {
	return h.host, nil
}

func (h fixedConsistentHash) GetHash(uint32) (string, error) {
	return h.host, nil
}

func (h fixedConsistentHash) Remove(string) bool {
	return false
}

func TestRingHashHandlerUsesConfiguredHashWithoutClusterLBPolicy(t *testing.T) {
	first := &model.Endpoint{
		ID:      "first",
		Address: model.SocketAddress{Address: "127.0.0.1", Port: 18080},
	}
	second := &model.Endpoint{
		ID:      "second",
		Address: model.SocketAddress{Address: "127.0.0.1", Port: 18081},
	}
	cluster := &model.ClusterConfig{
		Name:      "ring-direct",
		Endpoints: []*model.Endpoint{first, second},
		ConsistentHash: model.ConsistentHash{
			Hash: fixedConsistentHash{host: second.GetHost()},
		},
	}

	got := RingHashing{}.Handler(cluster, staticHashPolicy("request-key"))
	if got == nil {
		t.Fatal("expected endpoint, got nil")
	}
	if got.ID != second.ID {
		t.Fatalf("RingHashing picked %q, want %q", got.ID, second.ID)
	}
}

func TestRingHashHandlerFallsBackWhenConfiguredHashHitsUnhealthyEndpoint(t *testing.T) {
	healthy := &model.Endpoint{
		ID:      "healthy",
		Address: model.SocketAddress{Address: "127.0.0.1", Port: 18080},
	}
	unhealthy := &model.Endpoint{
		ID:        "unhealthy",
		Address:   model.SocketAddress{Address: "127.0.0.1", Port: 18081},
		UnHealthy: true,
	}
	cluster := &model.ClusterConfig{
		Name:      "ring-direct-unhealthy-hit",
		Endpoints: []*model.Endpoint{healthy, unhealthy},
		ConsistentHash: model.ConsistentHash{
			Hash: fixedConsistentHash{host: unhealthy.GetHost()},
		},
	}

	got := RingHashing{}.Handler(cluster, staticHashPolicy("request-key"))
	if got == nil {
		t.Fatal("expected fallback endpoint, got nil")
	}
	if got.ID != healthy.ID {
		t.Fatalf("RingHashing picked %q, want %q", got.ID, healthy.ID)
	}
}

func TestHashRing(t *testing.T) {

	nodeCount := 5

	nodes := make([]*model.Endpoint, 0, nodeCount)

	for i := 1; i <= nodeCount; i++ {
		name := strconv.Itoa(i)
		nodes = append(nodes, &model.Endpoint{ID: name, Name: name,
			Address: model.SocketAddress{Address: "192.168.1." + name, Port: 1000 + i}})
	}

	cluster := &model.ClusterConfig{
		Name:           "cluster1",
		Endpoints:      nodes,
		LbStr:          model.LoadBalancerRingHashing,
		ConsistentHash: model.ConsistentHash{ReplicaNum: 10, MaxVnodeNum: 1023},
	}
	cluster.CreateConsistentHash()

	var (
		hashing = RingHashing{}
		path    string
	)

	for i := 1; i <= 20; i++ {
		path = fmt.Sprintf("/pixiu?total=%d", i)
		t.Log(hashing.HandlerWithSnapshot(loadbalancer.PickContext{
			Config:           cluster,
			HealthyEndpoints: cluster.GetEndpoint(true),
		}, &http.HttpContext{Request: &stdHttp.Request{Method: stdHttp.MethodGet, RequestURI: path}}))
	}

}

func TestRingHashUsesHealthyConsistentHashSnapshot(t *testing.T) {
	first := &model.Endpoint{
		ID:      "first",
		Address: model.SocketAddress{Address: "127.0.0.1", Port: 18080},
	}
	unhealthy := &model.Endpoint{
		ID:        "unhealthy",
		Address:   model.SocketAddress{Address: "127.0.0.1", Port: 18081},
		UnHealthy: true,
	}
	second := &model.Endpoint{
		ID:      "second",
		Address: model.SocketAddress{Address: "127.0.0.1", Port: 18082},
	}
	cluster := &model.ClusterConfig{
		Name: "ring-healthy-snapshot",
		ConsistentHash: model.ConsistentHash{
			Hash: fixedConsistentHash{host: unhealthy.GetHost()},
		},
	}

	got := RingHashing{}.HandlerWithSnapshot(loadbalancer.PickContext{
		Config:                cluster,
		HealthyConsistentHash: fixedConsistentHash{host: second.GetHost()},
		HealthyEndpoints:      []*model.Endpoint{first, second},
	}, staticHashPolicy("key-for-unhealthy-slot"))

	if got == nil {
		t.Fatal("expected healthy endpoint, got nil")
	}
	if got.ID != second.ID {
		t.Fatalf("RingHashing picked %q, want %q", got.ID, second.ID)
	}
}

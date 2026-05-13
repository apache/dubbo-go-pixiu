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

package maglev

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

func TestMaglevHashHandlerUsesConfiguredHashWithoutClusterLBPolicy(t *testing.T) {
	first := &model.Endpoint{
		ID:      "first",
		Address: model.SocketAddress{Address: "127.0.0.1", Port: 18080},
	}
	second := &model.Endpoint{
		ID:      "second",
		Address: model.SocketAddress{Address: "127.0.0.1", Port: 18081},
	}
	cluster := &model.ClusterConfig{
		Name:      "maglev-direct",
		Endpoints: []*model.Endpoint{first, second},
		ConsistentHash: model.ConsistentHash{
			Hash: fixedConsistentHash{host: second.GetHost()},
		},
	}

	got := MaglevHash{}.Handler(cluster, staticHashPolicy("request-key"))
	if got == nil {
		t.Fatal("expected endpoint, got nil")
	}
	if got.ID != second.ID {
		t.Fatalf("MaglevHash picked %q, want %q", got.ID, second.ID)
	}
}

func TestMaglevHashHandlerFallsBackWhenConfiguredHashHitsUnhealthyEndpoint(t *testing.T) {
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
		Name:      "maglev-direct-unhealthy-hit",
		Endpoints: []*model.Endpoint{healthy, unhealthy},
		ConsistentHash: model.ConsistentHash{
			Hash: fixedConsistentHash{host: unhealthy.GetHost()},
		},
	}

	got := MaglevHash{}.Handler(cluster, staticHashPolicy("request-key"))
	if got == nil {
		t.Fatal("expected fallback endpoint, got nil")
	}
	if got.ID != healthy.ID {
		t.Fatalf("MaglevHash picked %q, want %q", got.ID, healthy.ID)
	}
}

func TestMaglevHash(t *testing.T) {

	nodeCount := 5

	nodes := make([]*model.Endpoint, 0, nodeCount)

	for i := 1; i <= nodeCount; i++ {
		name := strconv.Itoa(i)
		nodes = append(nodes, &model.Endpoint{ID: name, Name: name,
			Address: model.SocketAddress{Address: "192.168.1." + name, Port: 1000 + i}})
	}

	cluster := &model.ClusterConfig{
		Name:           "test-cluster",
		Endpoints:      nodes,
		LbStr:          model.LoadBalancerMaglevHashing,
		ConsistentHash: model.ConsistentHash{MaglevTableSize: 521},
	}
	cluster.CreateConsistentHash()

	var (
		hashing = MaglevHash{}
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

func TestLookUpTableHashIsStableForSameKey(t *testing.T) {
	table, err := NewLookUpTable(521, []string{"127.0.0.1:18080", "127.0.0.1:18081"})
	if err != nil {
		t.Fatalf("NewLookUpTable() error = %v", err)
	}

	first := table.Hash("same-request-key")
	for i := 0; i < 20; i++ {
		if got := table.Hash("same-request-key"); got != first {
			t.Fatalf("Hash() = %d, want stable %d", got, first)
		}
	}
}

func TestMaglevHashUsesHealthyConsistentHashSnapshot(t *testing.T) {
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
		Name: "maglev-healthy-snapshot",
		ConsistentHash: model.ConsistentHash{
			Hash: fixedConsistentHash{host: unhealthy.GetHost()},
		},
	}

	got := MaglevHash{}.HandlerWithSnapshot(loadbalancer.PickContext{
		Config:                cluster,
		HealthyConsistentHash: fixedConsistentHash{host: second.GetHost()},
		HealthyEndpoints:      []*model.Endpoint{first, second},
	}, staticHashPolicy("key-for-unhealthy-slot"))

	if got == nil {
		t.Fatal("expected healthy endpoint, got nil")
	}
	if got.ID != second.ID {
		t.Fatalf("MaglevHash picked %q, want %q", got.ID, second.ID)
	}
}

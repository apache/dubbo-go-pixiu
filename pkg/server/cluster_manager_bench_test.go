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
	"sync/atomic"
	"testing"
)

import (
	_ "github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/maglev"     // Register Maglev for benchmark coverage.
	_ "github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/rand"       // Register Rand for benchmark coverage.
	_ "github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/ringhash"   // Register RingHash for benchmark coverage.
	_ "github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/roundrobin" // Register RoundRobin for benchmark coverage.
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

var (
	// Keep benchmark results live so the compiler cannot optimize the hot path away.
	benchmarkEndpointSink  *model.Endpoint
	benchmarkEndpointsSink []*model.Endpoint
	benchmarkClusterSink   *model.ClusterConfig
)

type benchmarkHashPolicy string

func (p benchmarkHashPolicy) GenerateHash() string {
	return string(p)
}

func BenchmarkClusterPickEndpointSerial(b *testing.B) {
	for _, lbType := range []model.LbPolicyType{model.LoadBalancerRand, model.LoadBalancerRoundRobin} {
		for _, clusterCount := range []int{1, 32, 256, 1024} {
			b.Run(fmt.Sprintf("%s/clusters=%d", lbType, clusterCount), func(b *testing.B) {
				cm, names := benchmarkClusterManager(clusterCount, 4, lbType)

				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					benchmarkEndpointSink = cm.PickEndpoint(names[i%len(names)], nil)
				}
			})
		}
	}
}

func BenchmarkClusterLookupSerial(b *testing.B) {
	for _, clusterCount := range []int{1, 32, 256, 1024} {
		b.Run(fmt.Sprintf("clusters=%d", clusterCount), func(b *testing.B) {
			cm, names := benchmarkClusterManager(clusterCount, 4, model.LoadBalancerRoundRobin)

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				cm.rw.RLock()
				benchmarkClusterSink = cm.getCluster(names[i%len(names)])
				cm.rw.RUnlock()
			}
		})
	}
}

func BenchmarkClusterPickEndpointParallel(b *testing.B) {
	for _, lbType := range []model.LbPolicyType{model.LoadBalancerRand, model.LoadBalancerRoundRobin} {
		b.Run(string(lbType), func(b *testing.B) {
			cm, names := benchmarkClusterManager(256, 4, lbType)
			var workerCounter uint64

			b.ReportAllocs()
			b.ResetTimer()
			b.RunParallel(func(pb *testing.PB) {
				// Give each worker a stable starting offset once so bookkeeping stays off the hot path.
				idx := int(atomic.AddUint64(&workerCounter, 1)-1) % len(names)
				var endpoint *model.Endpoint
				for pb.Next() {
					endpoint = cm.PickEndpoint(names[idx], nil)
					idx++
					if idx == len(names) {
						idx = 0
					}
				}
				benchmarkEndpointSink = endpoint
			})
		})
	}
}

func BenchmarkClusterLoadBalancerHotPathSerial(b *testing.B) {
	for _, lbType := range []model.LbPolicyType{model.LoadBalancerRand, model.LoadBalancerRoundRobin} {
		for _, endpointCount := range []int{4, 64, 512} {
			b.Run(fmt.Sprintf("%s/endpoints=%d", lbType, endpointCount), func(b *testing.B) {
				cm := &ClusterManager{}
				cluster := benchmarkClusterConfig("lb-hot-path", lbType, endpointCount, 0)

				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					benchmarkEndpointSink = cm.pickOneEndpoint(cluster, nil)
				}
			})
		}
	}
}

func BenchmarkClusterHealthyFilterCost(b *testing.B) {
	for _, endpointCount := range []int{8, 64, 512} {
		for _, healthyRatio := range []int{100, 50, 0} {
			b.Run(fmt.Sprintf("endpoints=%d/healthy=%d", endpointCount, healthyRatio), func(b *testing.B) {
				cluster := benchmarkClusterConfig("healthy-filter", model.LoadBalancerRoundRobin, endpointCount, 0)
				healthyCount := endpointCount * healthyRatio / 100
				for i := healthyCount; i < len(cluster.Endpoints); i++ {
					cluster.Endpoints[i].UnHealthy = true
				}

				b.ReportAllocs()
				b.ResetTimer()
				for i := 0; i < b.N; i++ {
					benchmarkEndpointsSink = cluster.GetEndpoint(true)
				}
			})
		}
	}
}

func BenchmarkClusterCompareAndSetStoreMixed(b *testing.B) {
	cm, names := benchmarkClusterManager(128, 4, model.LoadBalancerRoundRobin)
	readIndex := 0

	b.ReportAllocs()
	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		if i%20 == 19 {
			newStore, err := benchmarkFreshStore(cm, names, model.LoadBalancerRoundRobin, 4)
			if err != nil {
				b.Fatal(err)
			}
			if !cm.CompareAndSetStore(newStore) {
				b.Fatal("CompareAndSetStore returned false")
			}
			continue
		}

		benchmarkEndpointSink = cm.PickEndpoint(names[readIndex%len(names)], nil)
		readIndex++
	}
}

func BenchmarkClusterConsistentHashResolve(b *testing.B) {
	keys := make([]benchmarkHashPolicy, 1024)
	for i := range keys {
		keys[i] = benchmarkHashPolicy(fmt.Sprintf("hash-key-%d", i))
	}

	for _, lbType := range []model.LbPolicyType{model.LoadBalancerRingHashing, model.LoadBalancerMaglevHashing} {
		b.Run(string(lbType), func(b *testing.B) {
			cm, names := benchmarkClusterManager(1, 64, lbType)
			clusterName := names[0]

			b.ReportAllocs()
			b.ResetTimer()
			for i := 0; i < b.N; i++ {
				benchmarkEndpointSink = cm.PickEndpoint(clusterName, keys[i%len(keys)])
			}
		})
	}
}

func benchmarkClusterManager(clusterCount int, endpointCount int, lbType model.LbPolicyType) (*ClusterManager, []string) {
	clusters := make([]*model.ClusterConfig, 0, clusterCount)
	names := make([]string, 0, clusterCount)
	for i := 0; i < clusterCount; i++ {
		name := fmt.Sprintf("cluster-%d", i)
		names = append(names, name)
		clusters = append(clusters, benchmarkClusterConfig(name, lbType, endpointCount, i*endpointCount))
	}
	return testClusterManager(clusters...), names
}

func benchmarkClusterConfig(name string, lbType model.LbPolicyType, endpointCount int, offset int) *model.ClusterConfig {
	endpoints := make([]*model.Endpoint, 0, endpointCount)
	for i := 0; i < endpointCount; i++ {
		endpoints = append(endpoints, testEndpoint(fmt.Sprintf("%s-ep-%d", name, i), "127.0.0.1", 20000+offset+i))
	}

	cluster := testCluster(name, lbType, endpoints)
	switch lbType {
	case model.LoadBalancerRingHashing:
		cluster.ConsistentHash = model.ConsistentHash{
			ReplicaNum:  128,
			MaxVnodeNum: 4099,
		}
	case model.LoadBalancerMaglevHashing:
		cluster.ConsistentHash = model.ConsistentHash{}
	}
	return cluster
}

func benchmarkFreshStore(cm *ClusterManager, names []string, lbType model.LbPolicyType, endpointCount int) (*ClusterStore, error) {
	oldStore, err := cm.CloneStore()
	if err != nil {
		return nil, err
	}

	newStore := cm.NewStore(oldStore.Version)
	for i, name := range names {
		newStore.AddCluster(benchmarkClusterConfig(name, lbType, endpointCount, i*endpointCount))
	}
	return newStore, nil
}

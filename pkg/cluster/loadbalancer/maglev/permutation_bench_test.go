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
	"crypto/sha3"
	"encoding/binary"
	"fmt"
	"hash/maphash"
	"os"
	"testing"
)

var benchmarkHashSink uint32
var benchmarkEndpointSink string

func BenchmarkLookUpTableHash(b *testing.B) {
	key := "GET./pixiu?total=100&user=benchmark"

	b.Run("maphash-random-seed", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchmarkHashSink = benchmarkMapHash(key)
		}
	})

	b.Run("sha3-512", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchmarkHashSink = benchmarkSHA3Hash(key)
		}
	})

	b.Run("xxhash", func(b *testing.B) {
		for i := 0; i < b.N; i++ {
			benchmarkHashSink = requestHash(key)
		}
	})
}

func BenchmarkLookUpTableGet(b *testing.B) {
	table := createTableWithNodes(10007, 100)
	table.Populate()

	keys := []string{
		"GET./pixiu?total=1",
		"GET./pixiu?total=2",
		"POST./api/v1/orders?id=100",
		"GET./grpc.service.Method",
		"GET./dubbo.method.provider",
	}

	b.ResetTimer()
	for i := 0; i < b.N; i++ {
		endpoint, err := table.Get(keys[i%len(keys)])
		if err != nil {
			b.Fatal(err)
		}
		benchmarkEndpointSink = endpoint
	}
}

func benchmarkMapHash(key string) uint32 {
	var h maphash.Hash
	h.SetSeed(maphash.MakeSeed())
	h.WriteString(key)
	return uint32(h.Sum64())
}

func benchmarkSHA3Hash(key string) uint32 {
	out := sha3.Sum512([]byte(key))
	return binary.LittleEndian.Uint32(out[:])
}

var benchmarkTableSink *LookUpTable

// BenchmarkLookUpTableFullBuild includes allocation, permutation generation,
// and slot population. Large cases are opt-in to keep routine benchmark runs short.
func BenchmarkLookUpTableFullBuild(b *testing.B) {
	for _, hostCount := range []int{1, 32, 256, 1024} {
		for _, tableSize := range []int{1009, 10007, 131071} {
			b.Run(fmt.Sprintf("hosts=%d/table=%d", hostCount, tableSize), func(b *testing.B) {
				if tableSize == 131071 && os.Getenv("PIXIU_MAGLEV_LARGE_BENCH") != "1" {
					b.Skip("set PIXIU_MAGLEV_LARGE_BENCH=1 to run large-memory cases")
				}
				hosts := make([]string, hostCount)
				for i := range hosts {
					hosts[i] = fmt.Sprintf("127.0.0.1:%d", 8000+i)
				}
				for _, ordered := range []bool{false, true} {
					name := "Populate"
					if ordered {
						name = "FixedOrder"
					}
					b.Run(name, func(b *testing.B) {
						b.ReportAllocs()
						for b.Loop() {
							table, err := NewLookUpTable(tableSize, hosts)
							if err != nil {
								b.Fatal(err)
							}
							if ordered {
								// Match Populate's work and locking, but use input order
								// instead of map iteration for reproducible comparisons.
								table.Lock()
								table.permutations = make([]*permutation, 0, len(hosts))
								for i, host := range hosts {
									table.generatePerm(host, i)
								}
								table.populate()
								table.Unlock()
							} else {
								table.Populate()
							}
							benchmarkTableSink = table
						}
					})
				}
			})
		}
	}
}

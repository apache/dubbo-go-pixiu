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
	"hash/maphash"
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

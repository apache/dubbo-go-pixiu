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
	"encoding/binary"
	"fmt"
	"math"
	"sync"

	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"

	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

type permutation struct {
	pos   []uint64
	next  int
	index int
	hit   int
}

type LookUpTable struct {
	slot         []string
	permutations []*permutation
	buckets      map[int]string
	endpointNum  int
	size         int
	sync.RWMutex
}

func NewLookUpTable(factor int, endpoints []*model.Endpoint) *LookUpTable {
	if factor < 10 {
		logger.Debugf("[dubbo-go-pixiu] The factor of Maglev load balancing should be an integer greater than 10,"+
			" but got %d instead. Setting factor to 10 by default.", factor)
		factor = 10
	}

	buckets := make(map[int]string)
	for i, endpoint := range endpoints {
		address := endpoint.Address
		buckets[i] = fmt.Sprintf("%s:%d", address.Address, address.Port)
	}
	n := len(buckets)
	m := factor * n

	return &LookUpTable{
		slot:         make([]string, m),
		permutations: make([]*permutation, n),
		buckets:      buckets,
		endpointNum:  n,
		size:         m,
	}
}

// Populate Magelev hashing look up table.
func (t *LookUpTable) Populate() {
	t.Lock()
	defer t.Unlock()

	t.genPermutations()

	full, miss := 0, 0
	for miss == t.endpointNum && t.endpointNum != 0 || full != t.size {
		for _, p := range t.permutations {
			if p.next == t.size {
				continue
			}
			start := p.next
			for start < t.size && len(t.slot[p.pos[start]]) > 0 {
				start++
			}
			if start < t.size {
				t.slot[p.pos[start]] = t.buckets[p.index]
				p.hit++
				full++
			} else {
				miss++
			}
			p.next = start
		}
	}

	// Fill the empty slot with the least placed Endpoint.
	if full != t.size {
		t.fillMissingSlot()
	}
}

func (t *LookUpTable) genPermutations() {
	var offset, skip, j uint64
	m := uint64(t.size)
	for i, B := range t.buckets {
		pos := make([]uint64, m)
		offset = _hash1(B) % m
		skip = _hash2(B)%m*(m-1) + 1
		for j = 0; j < m; j++ {
			pos[j] = (offset + j*skip) % m
		}
		t.permutations[i] = &permutation{pos, 0, i, 0}
	}
}

func (t *LookUpTable) fillMissingSlot() {
	// todo: what to do when hit == 0
	var minP *permutation
	minHit := math.MaxInt
	for _, p := range t.permutations {
		if p.hit < minHit {
			minHit = p.hit
			minP = p
		}
	}
	for i, s := range t.slot {
		if len(s) == 0 {
			t.slot[i] = t.buckets[minP.index]
			minP.hit++
		}
	}
}

func _hash1(key string) uint64 {
	out := sha3.Sum512([]byte(key))
	return binary.LittleEndian.Uint64(out[:])
}

func _hash2(key string) uint64 {
	out := blake2b.Sum512([]byte(key))
	return binary.LittleEndian.Uint64(out[:])
}

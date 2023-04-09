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
	"hash/maphash"
	"math"
	"math/big"
	"sync"
)

import (
	"github.com/pkg/errors"
	"golang.org/x/crypto/blake2b"
	"golang.org/x/crypto/sha3"
)

// Default table size support maximum 10,000 endpoints
var defaultTableSize = []int{
	307, 503, 1511, 2503, 5099, 7877,
	10007, 14207, 20903, 40423, 88721,
	123911, 164821, 199999, 218077, 299777,
	344941, 399989, 463031, 514399, 604613,
	686669, 726337, 799789, 857011, 999983,
}

type permutation struct {
	pos   []uint32
	next  int
	index int
	hit   int
}

type LookUpTable struct {
	slots        []string
	permutations []*permutation
	buckets      map[int]string
	endpointNum  int
	size         int
	sync.RWMutex
}

func NewLookUpTable(tableSize int, hosts []string) (*LookUpTable, error) {
	expectN := len(hosts) * 100
	if tableSize == 0 {
		// find closet table size
		for _, size := range defaultTableSize {
			if expectN < size {
				tableSize = size
				break
			}
		}
	}

	if tableSize < expectN && !big.NewInt(0).SetUint64(uint64(tableSize)).ProbablyPrime(1) {
		return nil, errors.Errorf("table size should be a prime number that greater than at least "+
			"100 times of endpoints size, but got %d instead", tableSize)
	}

	buckets := make(map[int]string, len(hosts))
	for i, host := range hosts {
		buckets[i] = host
	}
	n := len(buckets)

	return &LookUpTable{
		buckets:     buckets,
		endpointNum: n,
		size:        tableSize,
	}, nil
}

// Populate Magelev hashing look up table.
func (t *LookUpTable) Populate() {
	t.Lock()
	defer t.Unlock()

	t.generatePerms()
	t.populate()
}

func (t *LookUpTable) populate() {
	t.slots = make([]string, t.size)

	full, miss := 0, 0
	for miss < t.endpointNum && full < t.size {
		for _, p := range t.permutations {
			if p.next == t.size {
				continue
			}
			start := p.next
			for start < t.size && len(t.slots[p.pos[start]]) > 0 {
				start++
			}
			if start < t.size {
				t.slots[p.pos[start]] = t.buckets[p.index]
				p.hit++
				full++
			} else {
				miss++
			}
			p.next = start
		}
	}

	// Fill the empty slots with the least placed Endpoint.
	if full != t.size && miss > 0 {
		t.fillMissingSlots()
	}
}

func (t *LookUpTable) fillMissingSlots() {
	var minP *permutation
	minHit := math.MaxInt
	for _, p := range t.permutations {
		if p.hit < minHit {
			minHit = p.hit
			minP = p
		}
	}
	for i, s := range t.slots {
		if len(s) == 0 {
			t.slots[i] = t.buckets[minP.index]
			minP.hit++
		}
	}
}

func (t *LookUpTable) generatePerms() {
	t.permutations = make([]*permutation, 0, t.endpointNum)

	for i, b := range t.buckets {
		t.generatePerm(b, i)
	}
}

func (t *LookUpTable) generatePerm(bucket string, i int) {
	var offs, skip, j uint32

	m := uint32(t.size)
	pos := make([]uint32, m)
	offs = _hash1(bucket) % m
	skip = _hash2(bucket)%(m-1) + 1
	for j = 0; j < m; j++ {
		pos[j] = (offs + j*skip) % m
	}
	t.permutations = append(t.permutations, &permutation{pos, 0, i, 0})
}

func (t *LookUpTable) resetPerms() {
	for _, p := range t.permutations {
		p.next = 0
		p.hit = 0
	}
}

func (t *LookUpTable) removePerm(dst int) {
	del := -1
	for i, p := range t.permutations {
		if p.index == dst {
			del = i
			break
		}
	}
	if del != -1 {
		t.permutations[del] = nil
		t.permutations = append(t.permutations[:del], t.permutations[del+1:]...)
	}
}

// Hash the input key.
func (t *LookUpTable) Hash(key string) uint32 {
	var h maphash.Hash
	h.SetSeed(maphash.MakeSeed())
	h.WriteString(key)
	return uint32(h.Sum64())
}

// Get a slot by hashing the input key.
func (t *LookUpTable) Get(key string) (string, error) {
	t.RLock()
	defer t.RUnlock()

	if len(t.slots) == 0 {
		return "", errors.New("no slot initialized, please call Populate() before Get()")
	}

	if t.size == 0 || t.endpointNum == 0 {
		return "", errors.New("no host added")
	}

	dst := t.Hash(key) % uint32(t.size)
	return t.slots[dst], nil
}

// GetHash a slot by a hashed key.
func (t *LookUpTable) GetHash(key uint32) (string, error) {
	t.RLock()
	defer t.RUnlock()

	if len(t.slots) == 0 {
		return "", errors.New("no slot initialized, please call Populate() before Get()")
	}

	if t.size == 0 || t.endpointNum == 0 {
		return "", errors.New("no host added")
	}

	return t.slots[key], nil
}

// Add one endpoint into lookup table.
func (t *LookUpTable) Add(host string) {
	t.Lock()
	defer t.Unlock()

	t.add(host)
}

func (t *LookUpTable) add(host string) {
	dst := 0
	for i, bucket := range t.buckets {
		if bucket == host {
			return
		}
		if i > dst {
			dst = i
		}
	}
	t.buckets[dst+1] = host
	t.endpointNum++
	t.resetPerms()
	t.generatePerm(host, dst+1)
	t.populate()
}

// Remove one endpoint from lookup table.
func (t *LookUpTable) Remove(host string) bool {
	t.Lock()
	defer t.Unlock()

	return t.remove(host)
}

func (t *LookUpTable) remove(host string) bool {
	if t.endpointNum <= 0 {
		return false
	}

	for i, bucket := range t.buckets {
		if bucket == host {
			delete(t.buckets, i)
			t.endpointNum--
			t.removePerm(i)
			t.resetPerms()
			t.populate()
			return true
		}
	}

	return false
}

func _hash1(key string) uint32 {
	out := sha3.Sum512([]byte(key))
	return binary.LittleEndian.Uint32(out[:])
}

func _hash2(key string) uint32 {
	out := blake2b.Sum512([]byte(key))
	return binary.LittleEndian.Uint32(out[:])
}

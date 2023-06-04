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

package model

// LbPolicyType the load balance policy enum
type LbPolicyType string

const (
	LoadBalancerRand          LbPolicyType = "Rand"
	LoadBalancerRoundRobin    LbPolicyType = "RoundRobin"
	LoadBalancerRingHashing   LbPolicyType = "RingHashing"
	LoadBalancerMaglevHashing LbPolicyType = "MaglevHashing"
)

var LbPolicyTypeValue = map[string]LbPolicyType{
	"Rand":          LoadBalancerRand,
	"RoundRobin":    LoadBalancerRoundRobin,
	"RingHashing":   LoadBalancerRingHashing,
	"MaglevHashing": LoadBalancerMaglevHashing,
}

type LbPolicy interface {
	GenerateHash() string
}

// LbConsistentHash supports consistent hash load balancing
type LbConsistentHash interface {
	Hash(key string) uint32
	Add(host string)
	Get(key string) (string, error)
	GetHash(key uint32) (string, error)
	Remove(host string) bool
}

type ConsistentHashInitFunc = func(ConsistentHash, []*Endpoint) LbConsistentHash

// ConsistentHashInitMap stores the Init functions for consistent hash load balancing
var ConsistentHashInitMap = map[LbPolicyType]ConsistentHashInitFunc{}

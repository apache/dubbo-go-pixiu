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
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer"
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer/ringhash"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func init() {
	loadbalancer.RegisterLoadBalancer(model.LoadBalancerMaglevHashing, MaglevHash{})
	loadbalancer.RegisterConsistentHashInit(model.LoadBalancerMaglevHashing, NewMaglevHash)
}

func NewMaglevHash(config model.ConsistentHash, endpoints []*model.Endpoint) model.LbConsistentHash {
	hosts := make([]string, len(endpoints))
	for i, endpoint := range endpoints {
		hosts[i] = endpoint.GetHost()
	}

	h, err := NewLookUpTable(config.MaglevTableSize, hosts)
	if err == nil {
		h.Populate()
		return h
	}

	logger.Warnf("[dubbo-go-pixiu] maglev hash load balancing fail: %v, using ring hash instead", err)
	if config.ReplicaNum == 0 {
		config.ReplicaNum = 2 * len(endpoints)
	}
	return ringhash.NewRingHash(config, endpoints)
}

type MaglevHash struct{}

func (MaglevHash) UseHealthyEndpointsOnly() bool {
	return true
}

func (MaglevHash) UseZeroCopySnapshot() bool {
	return true
}

func (m MaglevHash) Handler(c *model.ClusterConfig, policy model.LbPolicy) *model.Endpoint {
	if c == nil {
		return nil
	}
	endpoints := c.GetEndpoint(true)
	hashView := model.ReadOnlyConsistentHash(c.ConsistentHash.Hash)
	if hashView == nil && len(endpoints) > 0 {
		hashView = model.ReadOnlyConsistentHash(NewMaglevHash(c.ConsistentHash, endpoints))
	}
	return m.HandlerWithSnapshot(loadbalancer.PickContext{
		Config:                c,
		HealthyConsistentHash: hashView,
		HealthyEndpoints:      endpoints,
	}, policy)
}

func (m MaglevHash) HandlerWithSnapshot(c loadbalancer.PickContext, policy model.LbPolicy) *model.Endpoint {
	endpoints := c.HealthyEndpoints
	hashView := loadbalancer.ConsistentHashForHealthyEndpoints(c)
	if len(endpoints) == 0 || hashView == nil || policy == nil {
		return nil
	}

	dst, err := hashView.Get(policy.GenerateHash())
	if err != nil {
		logger.Warnf("[dubbo-go-pixiu] error of getting from maglev hash: %v", err)
		return nil
	}

	for _, endpoint := range endpoints {
		if endpoint.GetHost() == dst {
			return endpoint
		}
	}
	return endpoints[0]
}

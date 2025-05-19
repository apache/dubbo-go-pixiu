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
	"math"
)

import (
	"github.com/dubbogo/gost/hash/consistent"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func init() {
	loadbalancer.RegisterLoadBalancer(model.LoadBalancerRingHashing, RingHashing{})
	loadbalancer.RegisterConsistentHashInit(model.LoadBalancerRingHashing, NewRingHash)
}

func NewRingHash(config model.ConsistentHash, endpoints []*model.Endpoint) model.LbConsistentHash {
	var ops []consistent.Option

	if config.ReplicaNum != 0 {
		ops = append(ops, consistent.WithReplicaNum(config.ReplicaNum))
	}

	if config.MaxVnodeNum != 0 {
		ops = append(ops, consistent.WithMaxVnodeNum(int(config.MaxVnodeNum)))
	} else {
		config.MaxVnodeNum = math.MinInt32
	}

	h := consistent.NewConsistentHash(ops...)
	for _, endpoint := range endpoints {
		h.Add(endpoint.GetHost())
	}
	return h
}

type RingHashing struct{}

func (r RingHashing) Handler(c *model.ClusterConfig, policy model.LbPolicy) *model.Endpoint {
	u := c.ConsistentHash.Hash.Hash(policy.GenerateHash())
	hash, err := c.ConsistentHash.Hash.GetHash(u)
	if err != nil {
		logger.Warnf("[dubbo-go-pixiu] error of getting from ring hash: %v", err)
		return nil
	}

	endpoints := c.GetEndpoint(true)

	for _, endpoint := range endpoints {
		if endpoint.GetHost() == hash {
			return endpoint
		}
	}

	if len(endpoints) == 0 {
		return nil
	}

	return endpoints[0]
}

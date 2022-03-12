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

package consistenthashing

import (
	"math/rand"
)

import (
	"github.com/dubbogo/gost/hash/consistent"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func init() {
	loadbalancer.RegisterLoadBalancer(model.LoadBalanceConsistentHashing, ConsistentHashing{})
}

var clusterMap = map[string]*consistent.Consistent{}

type ConsistentHashing struct{}

func (ConsistentHashing) Handler(c *model.Cluster) *model.Endpoint {
	data, ok := clusterMap[c.Name]
	if ok {

		// todo random ?
		hash, err := data.GetHash(uint32(rand.Int31n(1023)))
		if err != nil {
			return nil
		}

		for _, endpoint := range c.Endpoints {
			if endpoint.Address.Address == hash {
				return endpoint
			}
		}

		return nil
	}

	// todo replicaNum parameters or yaml config ?
	h := consistent.NewConsistentHash(consistent.WithReplicaNum(10), consistent.WithMaxVnodeNum(1023))
	for _, endpoint := range c.Endpoints {
		h.Add(endpoint.Address.Address)
	}
	clusterMap[c.Name] = h
	return c.Endpoints[rand.Intn(len(c.Endpoints))]
}

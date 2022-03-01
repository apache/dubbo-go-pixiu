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
	"strconv"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func init() {
	loadbalancer.RegisterLoadBalancer(model.LoadBalanceConsistentHashing, ConsistentHashing{})
}

// todo wait cluster nodes judge logic
// todo add read write lock or mutex or other lock?
var clusterMap = map[string]*HashRing{}

type ConsistentHashing struct{}

func (ConsistentHashing) Handler(c *model.Cluster) *model.Endpoint {
	data, ok := clusterMap[c.Name]
	if ok {
		return data.getNode(strconv.Itoa(rand.Intn(10)))
	}
	// todo replicaNum parameters or yaml config ?
	h := NewHashRing(c.Endpoints, 10)
	clusterMap[c.Name] = h
	return h.getNode(strconv.Itoa(rand.Intn(10)))
}

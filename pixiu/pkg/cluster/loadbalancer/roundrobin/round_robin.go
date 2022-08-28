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

package roundrobin

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/cluster/loadbalancer"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

func init() {
	loadbalancer.RegisterLoadBalancer(model.LoadBalancerRoundRobin, RoundRobin{})
}

type RoundRobin struct{}

func (RoundRobin) Handler(c *model.ClusterConfig) *model.Endpoint {
	endpoints := c.GetEndpoint(true)
	lens := len(endpoints)
	if c.PrePickEndpointIndex >= lens {
		c.PrePickEndpointIndex = 0
	}
	e := endpoints[c.PrePickEndpointIndex]
	c.PrePickEndpointIndex = (c.PrePickEndpointIndex + 1) % lens
	return e
}

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

package weightrandom

import (
	"math/rand"
	"strconv"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/loadbalancer"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func init() {
	loadbalancer.RegisterLoadBalancer(model.LoadBalancerWeightRandom, WeightRandom{})
}

type weightedEndpoint struct {
	endpoint *model.Endpoint
	weight   int
}

// WeightRandom is a load balancing strategy that selects an endpoint based on weighted random selection.
// It assigns weights to endpoints and uses these weights to influence the probability of selection.
type WeightRandom struct{}

func (WeightRandom) Handler(c *model.ClusterConfig, _ model.LbPolicy) *model.Endpoint {
	endpoints := c.GetEndpoint(true)

	if len(endpoints) == 0 {
		return nil
	}

	var (
		weightedEndpoints = make([]*weightedEndpoint, 0, len(endpoints))
		totalWeight       int
	)

	for _, endpoint := range endpoints {
		weight := 0 // default weight
		if weightStr, ok := endpoint.Metadata["weight"]; ok {
			if w, err := strconv.Atoi(weightStr); err == nil && w >= 0 {
				weight = w
			}
		}
		totalWeight += weight
		weightedEndpoints = append(weightedEndpoints, &weightedEndpoint{endpoint: endpoint, weight: weight})
	}

	if totalWeight <= 0 {
		// if the sum of weights is 0 or negative, return a random endpoint
		randomIndex := rand.Intn(len(endpoints)) // NOSONAR
		return endpoints[randomIndex]
	}

	randomNumber := rand.Intn(totalWeight) // NOSONAR

	// iterate through the weighted endpoints
	// find the one that corresponds to the random number
	currentWeightSum := 0
	for _, we := range weightedEndpoints {
		currentWeightSum += we.weight
		if randomNumber < currentWeightSum {
			return we.endpoint
		}
	}

	// this line should not be reached, just for safety
	return endpoints[0]
}

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

package exit

import (
	"github.com/apache/dubbo-go-pixiu/pkg/cluster/retryer"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func init() {
	retryer.RegisterRetryer(model.RetryerTryNext, TryNext{})
}

type TryNext struct{}

func (TryNext) Handler(c *model.ClusterConfig, policy model.Policy) *model.Endpoint {
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

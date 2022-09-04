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

package cluster

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/cluster/healthcheck"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

type Cluster struct {
	HealthCheck *healthcheck.HealthChecker
	Config      *model.ClusterConfig
}

func NewCluster(clusterConfig *model.ClusterConfig) *Cluster {
	c := &Cluster{
		Config: clusterConfig,
	}

	// only handle one health checker
	if len(c.Config.HealthChecks) != 0 {
		c.HealthCheck = healthcheck.CreateHealthCheck(clusterConfig, c.Config.HealthChecks[0])
		c.HealthCheck.Start()
	}
	return c
}

func (c *Cluster) Stop() {
	if c.HealthCheck != nil {
		c.HealthCheck.Stop()
	}
}

func (c *Cluster) RemoveEndpoint(endpoint *model.Endpoint) {
	if c.HealthCheck != nil {
		c.HealthCheck.StopOne(endpoint)
	}
}

func (c *Cluster) AddEndpoint(endpoint *model.Endpoint) {
	if c.HealthCheck != nil {
		c.HealthCheck.StartOne(endpoint)
	}
}

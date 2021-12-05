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

package xds

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	xdspb "github.com/apache/dubbo-go-pixiu/pkg/model/xds/model"
)

type CdsManager struct {
	DiscoverApi
	clusters []*xdspb.Cluster
}

// Fetch overwrite DiscoverApi.Fetch.
func (c *CdsManager) Fetch() error {
	r, err := c.DiscoverApi.Fetch("") //todo use local version
	if err != nil {
		logger.Error("can not fetch lds", err)
		return err
	}
	if len(c.clusters) == 0 {
		c.clusters = make([]*xdspb.Cluster, 0, len(r))
	}
	for _, one := range r {
		_cluster := &xdspb.Cluster{}
		if err := one.To(_cluster); err != nil {
			logger.Errorf("unknown resource of %s, expect Listener", one.GetName())
			continue
		}
		logger.Infof("clusters from xds server %v", _cluster)
		c.clusters = append(c.clusters, _cluster)
	}
	return nil
}

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

package server

import (
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

func TestClusterManager(t *testing.T) {
	bs := &model.Bootstrap{
		StaticResources: model.StaticResources{
			Clusters: []*model.ClusterConfig{
				{
					Name: "test",
					Endpoints: []*model.Endpoint{
						{
							Address: model.SocketAddress{},
							ID:      "1",
						},
					},
				},
			},
		},
	}

	cm := CreateDefaultClusterManager(bs)
	assert.Equal(t, len(cm.store.Config), 1)

	cm.AddCluster(&model.ClusterConfig{
		Name: "test2",
		Endpoints: []*model.Endpoint{
			{
				Address: model.SocketAddress{},
				ID:      "1",
			},
		},
	})

	assert.Equal(t, len(cm.store.Config), 2)

	cm.SetEndpoint("test2", &model.Endpoint{
		Address: model.SocketAddress{},
		ID:      "2",
	})
	assert.Equal(t, cm.PickEndpoint("test").ID, "1")
	cm.DeleteEndpoint("test2", "1")
}

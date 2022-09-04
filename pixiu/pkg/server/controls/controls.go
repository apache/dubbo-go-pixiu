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

package controls

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

type (
	ClusterManager interface {
		RemoveCluster(names []string)
		HasCluster(name string) bool
		UpdateCluster(cluster *model.ClusterConfig)
		AddCluster(cluster *model.ClusterConfig)
		CloneXdsControlStore() (ClusterStore, error)
	}

	ListenerManager interface {
		AddListener(m *model.Listener) error
		UpdateListener(m *model.Listener) error
		RemoveListener(names []string)
		HasListener(name string) bool
		CloneXdsControlListener() ([]*model.Listener, error)
	}

	DynamicResourceManager interface {
		GetLds() *model.ApiConfigSource
		GetCds() *model.ApiConfigSource
		GetNode() *model.Node
	}

	ClusterStore interface {
		Config() []*model.ClusterConfig
		HasCluster(name string) bool
	}
)

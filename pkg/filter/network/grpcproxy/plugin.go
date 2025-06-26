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

package grpcproxy

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	// Kind filter plugin kind
	Kind = constant.GRPCProxyConnectionFilter
)

func init() {
	filter.RegisterNetworkFilterPlugin(&Plugin{})
}

// Plugin gRPC connection manager plugin, similar to dubboproxy plugin
type Plugin struct{}

// Kind returns the unique kind name to represent itself.
func (p *Plugin) Kind() string {
	return Kind
}

// CreateFilter return the filter instance
func (p *Plugin) CreateFilter(config any) (filter.NetworkFilter, error) {
	return CreateGrpcProxyConnectionManager(config.(*model.GRPCConnectionManagerConfig)), nil
}

// Config return the config struct
func (p *Plugin) Config() any {
	return &model.GRPCConnectionManagerConfig{}
}

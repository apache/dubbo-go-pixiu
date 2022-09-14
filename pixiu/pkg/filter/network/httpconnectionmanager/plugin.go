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

package httpconnectionmanager

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/util/stringutil"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

const (
	Kind = constant.HTTPConnectManagerFilter
)

func init() {
	filter.RegisterNetworkFilterPlugin(&Plugin{})
}

type (
	Plugin struct{}
)

// Kind
func (p *Plugin) Kind() string {
	return Kind
}

// CreateFilter create http network filter
func (p *Plugin) CreateFilter(config interface{}) (filter.NetworkFilter, error) {
	hcmc := config.(*model.HttpConnectionManagerConfig)
	hcmc.Timeout = stringutil.ResolveTimeStr2Time(hcmc.TimeoutStr, constant.DefaultReqTimeout)
	return http.CreateHttpConnectionManager(hcmc), nil
}

// Config return HttpConnectionManagerConfig
func (p *Plugin) Config() interface{} {
	return &model.HttpConnectionManagerConfig{}
}

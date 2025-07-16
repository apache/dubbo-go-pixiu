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

package mcpserver

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
)

const (
	// Kind 是 MCP Server Filter 的类型标识
	Kind = constant.MCPServerFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

// Plugin 实现 filter.HttpFilterPlugin 接口
type Plugin struct{}

// Kind 返回插件类型
func (p *Plugin) Kind() string {
	return Kind
}

// CreateFilterFactory 创建 FilterFactory
func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

// Config 返回配置结构体
func (p *Plugin) Config() any {
	return &Config{}
}

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

package a2a

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

// Plugin implements filter.HttpFilterPlugin interface for A2A Server Filter
type Plugin struct{}

// Kind returns the plugin type identifier
func (p *Plugin) Kind() string {
	return Kind
}

// CreateFilterFactory creates a new FilterFactory instance
func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

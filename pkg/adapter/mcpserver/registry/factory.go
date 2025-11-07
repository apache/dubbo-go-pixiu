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

package registry

import (
	"fmt"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

var providers = map[string]BuildFunc{}

// RegisterProvider registers a provider-specific controller builder.
func RegisterProvider(protocol string, fn BuildFunc) {
	providers[protocol] = fn
}

// BuildController builds a provider controller based on registry protocol.
func BuildController(reg model.Registry, onChange func(serverId string, cfg *model.McpServerConfig)) (Controller, error) {
	if fn, ok := providers[reg.Protocol]; ok {
		return fn(reg, onChange)
	}
	return nil, fmt.Errorf("no provider for protocol: %s", reg.Protocol)
}

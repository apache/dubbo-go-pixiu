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

package nacos

import (
	"fmt"
	"strings"
)

import (
	nacosconstant "github.com/nacos-group/nacos-sdk-go/v2/common/constant"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/mcpserver/common/util"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/mcpserver/registry"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	defaultNacosTimeoutMs = 5000
)

// provider self registration: only register the builder
func init() {
	registry.RegisterProvider(constant.Nacos, BuildController)
}

// BuildController builds a Nacos MCP registry controller
func BuildController(reg model.Registry, onChange func(serverId string, cfg *model.McpServerConfig)) (registry.Controller, error) {
	// Pre-validate addresses using util function
	if reg.Address != "" {
		if err := util.ValidateNacosAddresses(reg.Address); err != nil {
			return nil, fmt.Errorf("[dubbo-go-pixiu] nacos registry address validation failed: %v", err)
		}
	}

	// build server configs from comma-separated addresses
	serverCfgs := []nacosconstant.ServerConfig{}

	if reg.Address != "" {
		for _, part := range strings.Split(reg.Address, ",") {
			addr := strings.TrimSpace(part)
			if addr == "" {
				continue
			}
			result, err := util.ParseHostPortFromURL(addr)
			if err != nil {
				// This should not happen after validation, but keep for safety
				logger.Errorf("[dubbo-go-pixiu] nacos registry unexpected parse error for address '%s': %v", addr, err)
				continue
			}
			if result.UsedFallback {
				logger.Warnf("[dubbo-go-pixiu] nacos registry using fallback for address '%s': %s", addr, result.FallbackInfo)
			}
			serverCfgs = append(serverCfgs, nacosconstant.ServerConfig{IpAddr: result.Host, Port: uint64(result.Port)})
		}
	}

	logger.Infof("[dubbo-go-pixiu] nacos registry initialized with %d server(s)", len(serverCfgs))

	clientCfg := nacosconstant.ClientConfig{
		TimeoutMs:   defaultNacosTimeoutMs,
		NamespaceId: reg.Namespace,
		Username:    reg.Username,
		Password:    reg.Password,
	}

	client, err := NewMcpRegistryClient(&clientCfg, serverCfgs, reg.Namespace)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] nacos registry create Nacos MCP client failed: %v", err)
		return nil, err
	}

	controller := NewMcpController(client, func(serverId string, cfg *McpServerConfig) {
		if cfg == nil {
			onChange(serverId, nil)
			return
		}

		mc := &model.McpServerConfig{
			Tools: cfg.ToolConfigs,
		}
		onChange(serverId, mc)
	})

	return controller, nil
}

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

package dubbo

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// DubboProxyConfig the config for dubbo proxy
type DubboProxyConfig struct {
	// Registries such as zk,nacos or etcd
	Registries map[string]model.Registry `yaml:"registries" json:"registries"`
	// Timeout
	Timeout *model.TimeoutConfig `yaml:"timeout_config" json:"timeout_config"`
	// IsDefaultMap whether to use DefaultMap role
	IsDefaultMap bool
	// AutoResolve whether to resolve api config from request
	AutoResolve bool `yaml:"auto_resolve" json:"auto_resolve,omitempty"`
	// Protoset path to load protoset files
	Protoset []string `yaml:"protoset" json:"protoset,omitempty"`
	// Load balance
	LoadBalance string `yaml:"load_balance"  json:"load_balance,omitempty"`
	// Retries number of retries
	Retries string `yaml:"retries" json:"retries,omitempty"`

	// Cluster strategy for dubbo client.
	// Valid values (case-sensitive): "failover", "failfast", "failsafe", "failback", "forking", "broadcast"
	// Source: dubbo-go v3.3.1 cluster constants
	// Default: "failover"
	Cluster string `yaml:"cluster,omitempty" json:"cluster,omitempty"`

	// Check whether to check provider availability on startup.
	// Uses pointer type to distinguish "not configured" (nil) from "explicitly false".
	// nil = let dubbo-go SDK use its default behavior
	Check *bool `yaml:"check,omitempty" json:"check,omitempty"`

	// Protocol type for dubbo client.
	// Valid values (case-sensitive): "dubbo", "tri", "rest"
	// ⚠️ Important: Triple protocol value is "tri", NOT "triple" (from triple.TRIPLE constant)
	// Source: dubbo-go v3.3.1 protocol names
	// Default: "dubbo"
	Protocol string `yaml:"protocol,omitempty" json:"protocol,omitempty"`
}

// GetCluster returns cluster strategy with default value "failover".
func (dpc *DubboProxyConfig) GetCluster() string {
	if dpc.Cluster == "" {
		return "failover"
	}
	return dpc.Cluster
}

// GetProtocol returns protocol with default value "dubbo".
func (dpc *DubboProxyConfig) GetProtocol() string {
	if dpc.Protocol == "" {
		return "dubbo"
	}
	return dpc.Protocol
}

// GetCheck returns check pointer (nil means not configured, let dubbo-go use its default).
// Does not provide default value to maintain backward compatibility.
func (dpc *DubboProxyConfig) GetCheck() *bool {
	return dpc.Check
}

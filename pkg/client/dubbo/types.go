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
	"context"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// DubboClient is the protocol-specific client contract for outbound Dubbo calls.
type DubboClient interface {
	Apply() error
	Close() error
	Call(ctx context.Context, req *DubboOutboundRequest) (any, error)
}

// DubboOutboundRequest is an immutable contract for a single Dubbo outbound invocation.
type DubboOutboundRequest struct {
	Service string
	Method  string
	Group   string
	Version string

	Address string

	Protocol      string
	Serialization string

	Arguments   []any
	ParamTypes  []string
	Attachments map[string]any

	Timeout time.Duration
}

// DubboProxyConfig the config for dubbo proxy
type DubboProxyConfig struct {
	// Deprecated: AutoResolve is no longer supported. Remove auto_resolve from your
	// dubboProxyConfig and configure integrationRequest explicitly in the API definition.
	AutoResolve *bool `yaml:"auto_resolve" json:"auto_resolve,omitempty"`
	// Registries such as zk,nacos or etcd
	Registries map[string]model.Registry `yaml:"registries" json:"registries"`
	// Timeout
	Timeout *model.TimeoutConfig `yaml:"timeout_config" json:"timeout_config"`
	// Protoset path to load protoset files
	Protoset []string `yaml:"protoset" json:"protoset,omitempty"`
	// Load balance
	LoadBalance string `yaml:"load_balance"  json:"load_balance,omitempty"`
	// Retries number of retries
	Retries string `yaml:"retries" json:"retries,omitempty"`
}

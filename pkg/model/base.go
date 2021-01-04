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

package model

type Metadata struct {
	Info map[string]MetadataValue
}

type MetadataValue interface {
}

// Status is the components status
type Status int32

const (
	Down    Status = 0
	Up      Status = 1
	Unknown Status = 2
)

var StatusName = map[int32]string{
	0: "Down",
	1: "Up",
	2: "Unknown",
}

var StatusValue = map[string]int32{
	"Down":    0,
	"Up":      1,
	"Unknown": 2,
}

// ProtocolType
type ProtocolType int32

const (
	HTTP ProtocolType = 0 + iota // support for 1.0
	TCP
	UDP
)

// ProtocolTypeName
var ProtocolTypeName = map[int32]string{
	0: "HTTP",
	1: "TCP",
	2: "UDP",
}

// ProtocolTypeValue
var ProtocolTypeValue = map[string]int32{
	"HTTP": 0,
	"TCP":  1,
	"UDP":  2,
}

// Address the address
type Address struct {
	SocketAddress SocketAddress `yaml:"socket_address" json:"socket_address" mapstructure:"socket_address"`
	Name          string        `yaml:"name" json:"name" mapstructure:"name"`
}

// Address specify either a logical or physical address and port, which are
// used to tell proxy where to bind/listen, connect to upstream and find
// management servers
type SocketAddress struct {
	ProtocolStr  string       `yaml:"protocol_type" json:"protocol_type" mapstructure:"protocol_type"`
	Protocol     ProtocolType `yaml:"omitempty" json:"omitempty"`
	Address      string       `yaml:"address" json:"address" mapstructure:"address"`
	Port         int          `yaml:"port" json:"port" mapstructure:"port"`
	ResolverName string       `yaml:"resolver_name" json:"resolver_name" mapstructure:"resolver_name"`
}

// ConfigSource
type ConfigSource struct {
	Path            string          `yaml:"path" json:"path" mapstructure:"path"`
	ApiConfigSource ApiConfigSource `yaml:"api_config_source" json:"api_config_source" mapstructure:"api_config_source"`
}

// ApiConfigSource
type ApiConfigSource struct {
	ApiType     ApiType  `yaml:"omitempty" json:"omitempty"`
	ApiTypeStr  string   `yaml:"api_type" json:"api_type" mapstructure:"api_type"`
	ClusterName []string `yaml:"cluster_name" json:"cluster_name" mapstructure:"cluster_name"`
}

type ApiType int32

const (
	REST_VALUE  = "REST"
	GRPC_VALUE  = "GRPC"
	DUBBO_VALUE = "DUBBO"
)

const (
	REST ApiType = 0 + iota // support for 1.0
	GRPC
	DUBBO
)

var ApiTypeName = map[int32]string{
	0: REST_VALUE,
	1: GRPC_VALUE,
	2: DUBBO_VALUE,
}

var ApiTypeValue = map[string]int32{
	REST_VALUE:  0,
	GRPC_VALUE:  1,
	DUBBO_VALUE: 2,
}

// HeaderValueOption
type HeaderValueOption struct {
	Header []HeaderValue `yaml:"header" json:"header" mapstructure:"header"`
	Append []bool        `yaml:"append" json:"append" mapstructure:"append"`
}

// HeaderValue
type HeaderValue struct {
	Key   string `yaml:"key" json:"key" mapstructure:"key"`
	Value string `yaml:"value" json:"value" mapstructure:"value"`
}

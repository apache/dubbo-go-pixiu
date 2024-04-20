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

import (
	"fmt"
)

import (
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/api"
)

type Metadata struct {
	Info map[string]MetadataValue
}

type MetadataValue interface{}

// Status is the components status

const (
	// Down
	Down api.Status = 0
	// Up
	Up      api.Status = 1
	Unknown api.Status = 2
)

const (
	ApiTypeREST api.ApiType = 0 + iota // support for 1.0
	ApiTypeGRPC
	ApiTypeDUBBO
	ApiTypeIstioGRPC
)

var (
	StatusName = map[int32]string{
		0: "Down",
		1: "Up",
		2: "Unknown",
	}

	StatusValue = map[string]int32{
		"Down":    0,
		"Up":      1,
		"Unknown": 2,
	}

	ApiTypeName = map[int32]string{
		0: REST_VALUE,
		1: GRPC_VALUE,
		2: DUBBO_VALUE,
	}

	ApiTypeValue = map[string]int32{
		REST_VALUE:      0,
		GRPC_VALUE:      1,
		DUBBO_VALUE:     2,
		ISTIOGRPC_VALUE: 3,
	}
)

type (

	// Address the address
	Address struct {
		SocketAddress SocketAddress `yaml:"socket_address" json:"socket_address" mapstructure:"socket_address"`
		Name          string        `yaml:"name" json:"name" mapstructure:"name"`
	}

	// SocketAddress specify either a logical or physical address and port, which are
	// used to tell server where to bind/listen, connect to upstream and find
	// management servers
	SocketAddress struct {
		Address      string   `default:"0.0.0.0" yaml:"address" json:"address" mapstructure:"address"`
		Port         int      `default:"8881" yaml:"port" json:"port" mapstructure:"port"`
		ResolverName string   `yaml:"resolver_name" json:"resolver_name" mapstructure:"resolver_name"`
		Domains      []string `yaml:"domains" json:"domains" mapstructure:"domains"`
		CertsDir     string   `yaml:"certs_dir" json:"certs_dir" mapstructure:"certs_dir"`
	}

	// ConfigSource todo remove un-used
	ConfigSource struct {
		Path            string          `yaml:"path" json:"path" mapstructure:"path"`
		ApiConfigSource ApiConfigSource `yaml:"api_config_source" json:"api_config_source" mapstructure:"api_config_source"`
	}

	// ApiConfigSource config the api info. compatible with envoy xDS API
	//	{
	//  "api_type": "...",
	//  "transport_api_version": "...",
	//  "cluster_names": [],
	//  "grpc_services": [],
	//  "refresh_delay": "{...}",
	//  "request_timeout": "{...}",
	//  "rate_limit_settings": "{...}",
	//  "set_node_on_first_message_only": "..."
	//	}
	ApiConfigSource struct {
		APIType        api.ApiType   `yaml:"omitempty" json:"omitempty"`
		APITypeStr     string        `yaml:"api_type" json:"api_type" mapstructure:"api_type"`
		ClusterName    []string      `yaml:"cluster_name" json:"cluster_name" mapstructure:"cluster_name"`
		RefreshDelay   string        `yaml:"refresh_delay" json:"refresh_delay" mapstructure:"refresh_delay"`
		RequestTimeout string        `yaml:"request_timeout" json:"request_timeout" mapstructure:"request_timeout"`
		GrpcServices   []GrpcService `yaml:"grpc_services" json:"grpc_services" mapstructure:"grpc_services"`
		//rate_limit_settings todo implement
		//set_node_on_first_message_only todo implement

	}

	// GrpcService define grpc service context
	GrpcService struct {
		Timeout         string        `yaml:"timeout" json:"timeout" mapstructure:"timeout"`
		InitialMetadata []HeaderValue `yaml:"initial_metadata" json:"initial_metadata" mapstructure:"initial_metadata"`
	}

	// HeaderValueOption
	HeaderValueOption struct {
		Header []HeaderValue `yaml:"header" json:"header" mapstructure:"header"`
		Append []bool        `yaml:"append" json:"append" mapstructure:"append"`
	}

	// HeaderValue
	HeaderValue struct {
		Key   string `yaml:"key" json:"key" mapstructure:"key"`
		Value string `yaml:"value" json:"value" mapstructure:"value"`
	}
)

func (a SocketAddress) GetAddress() string {
	return fmt.Sprintf("%s:%v", a.Address, a.Port)
}

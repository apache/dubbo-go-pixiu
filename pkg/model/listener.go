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
	"github.com/creasty/defaults"

	"github.com/mitchellh/mapstructure"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	ProtocolTypeHTTP ProtocolType = 0 + iota // support for 1.0
	ProtocolTypeTCP
	ProtocolTypeUDP
	ProtocolTypeHTTPS
	ProtocolTypeGRPC
	ProtocolTypeHTTP2
	ProtocolTypeTriple
)

const (
	REST_VALUE      = "REST"
	GRPC_VALUE      = "GRPC"
	DUBBO_VALUE     = "DUBBO"
	ISTIOGRPC_VALUE = "ISTIO"
)

var (
	// ProtocolTypeName enum seq to protocol type name
	ProtocolTypeName = map[int32]string{
		0: "HTTP",
		1: "TCP",
		2: "UDP",
		3: "HTTPS",
		4: "GRPC",
		5: "HTTP2",
		6: "TRIPLE",
	}

	// ProtocolTypeValue protocol type name to enum seq
	ProtocolTypeValue = map[string]int32{
		"HTTP":   0,
		"TCP":    1,
		"UDP":    2,
		"HTTPS":  3,
		"GRPC":   4,
		"HTTP2":  5,
		"TRIPLE": 6,
	}
)

type (
	// ProtocolType protocol type enum
	ProtocolType int32

	// Listener is a server, listener a port
	Listener struct {
		Name        string       `yaml:"name" json:"name" mapstructure:"name"`
		Address     Address      `yaml:"address" json:"address" mapstructure:"address"`
		ProtocolStr string       `default:"http" yaml:"protocol_type" json:"protocol_type" mapstructure:"protocol_type"`
		Protocol    ProtocolType `default:"http" yaml:"omitempty" json:"omitempty"`
		FilterChain FilterChain  `yaml:"filter_chains" json:"filter_chains" mapstructure:"filter_chains"`
		Config      any          `yaml:"config" json:"config" mapstructure:"config"`
	}

	// GrpcConfig gRPC listener specific configuration
	GrpcConfig struct {
		MaxReceiveMessageSize int        `default:"4194304" yaml:"max_receive_message_size" json:"max_receive_message_size" mapstructure:"max_receive_message_size"`
		MaxSendMessageSize    int        `default:"4194304" yaml:"max_send_message_size" json:"max_send_message_size" mapstructure:"max_send_message_size"`
		EnableCompression     bool       `yaml:"enable_compression" json:"enable_compression" mapstructure:"enable_compression"`
		IdleTimeout           string     `yaml:"idle_timeout" json:"idle_timeout" mapstructure:"idle_timeout"`
		MaxConnectionAge      string     `yaml:"max_connection_age" json:"max_connection_age" mapstructure:"max_connection_age"`
		EnableTLS             bool       `yaml:"enable_tls" json:"enable_tls" mapstructure:"enable_tls"`
		TLS                   *TLSConfig `yaml:"tls,omitempty" json:"tls,omitempty" mapstructure:"tls,omitempty"`
	}

	// TLSConfig TLS configuration for gRPC (optional)
	TLSConfig struct {
		CertFile string `yaml:"cert_file" json:"cert_file" mapstructure:"cert_file"`
		KeyFile  string `yaml:"key_file" json:"key_file" mapstructure:"key_file"`
	}
)

// MapInGrpcStruct maps any config to GrpcConfig struct
func MapInGrpcStruct(cfg any) *GrpcConfig {
	var gc GrpcConfig
	if cfg != nil {
		if err := mapstructure.Decode(cfg, &gc); err != nil {
			logger.Error("gRPC Config error", err)
		}
	}
	if err := defaults.Set(&gc); err != nil {
		logger.Errorf("set grpc config default error %v", err)
	}
	return &gc
}

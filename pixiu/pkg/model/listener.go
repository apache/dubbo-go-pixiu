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
		Config      interface{}  `yaml:"config" json:"config" mapstructure:"config"`
	}
)

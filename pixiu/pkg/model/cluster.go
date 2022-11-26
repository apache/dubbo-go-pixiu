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
	"github.com/dubbogo/gost/hash/consistent"
)

const (
	Static DiscoveryType = iota
	StrictDNS
	LogicalDns
	EDS
	OriginalDst
)

var (
	// DiscoveryTypeName
	DiscoveryTypeName = map[DiscoveryType]string{
		Static:      "Static",
		StrictDNS:   "StrictDNS",
		LogicalDns:  "LogicalDns",
		EDS:         "EDS",
		OriginalDst: "OriginalDst",
	}

	// DiscoveryTypeValue
	DiscoveryTypeValue = map[string]DiscoveryType{
		"Static":      Static,
		"StrictDNS":   StrictDNS,
		"LogicalDns":  LogicalDns,
		"EDS":         EDS,
		"OriginalDst": OriginalDst,
	}
)

type (
	// ClusterConfig a single upstream cluster
	ClusterConfig struct {
		Name                 string              `yaml:"name" json:"name"` // Name the cluster unique name
		TypeStr              string              `yaml:"type" json:"type"` // Type the cluster discovery type string value
		Type                 DiscoveryType       `yaml:"-" json:"-"`       // Type the cluster discovery type
		EdsClusterConfig     EdsClusterConfig    `yaml:"eds_cluster_config" json:"eds_cluster_config" mapstructure:"eds_cluster_config"`
		LbStr                LbPolicyType        `yaml:"lb_policy" json:"lb_policy"`   // Lb the cluster select node used loadBalance policy
		Hash                 Hash                `yaml:"consistent" json:"consistent"` // Consistent hash config info
		HealthChecks         []HealthCheckConfig `yaml:"health_checks" json:"health_checks"`
		Endpoints            []*Endpoint         `yaml:"endpoints" json:"endpoints"`
		PrePickEndpointIndex int
	}

	// EdsClusterConfig todo remove un-used EdsClusterConfig
	EdsClusterConfig struct {
		EdsConfig   ConfigSource `yaml:"eds_config" json:"eds_config" mapstructure:"eds_config"`
		ServiceName string       `yaml:"service_name" json:"service_name" mapstructure:"service_name"`
	}

	// Registry remote registry where dubbo apis are registered.
	// Here comes a problem, dubbo protocol proxy does not use the same registry as pixiu,
	// so any modification to the config, should apply to both `pkg/client/dubbo/dubbo.go`
	// and `pkg\adapter\dubboregistry\registry`
	Registry struct {
		Protocol  string `default:"zookeeper" yaml:"protocol" json:"protocol"`
		Timeout   string `yaml:"timeout" json:"timeout"`
		Address   string `yaml:"address" json:"address"`
		Username  string `yaml:"username" json:"username"`
		Password  string `yaml:"password" json:"password"`
		Group     string `default:"DEFAULT_GROUP"  yaml:"group" json:"group"`
		Namespace string `yaml:"namespace" json:"namespace"`
	}

	// DiscoveryType
	DiscoveryType int32

	// Endpoint
	Endpoint struct {
		ID        string            `yaml:"ID" json:"ID"`                                                       // ID indicate one endpoint
		Name      string            `yaml:"name" json:"name"`                                                   // Name the cluster unique name
		Address   SocketAddress     `yaml:"socket_address" json:"socket_address" mapstructure:"socket_address"` // Address socket address
		Metadata  map[string]string `yaml:"meta" json:"meta"`                                                   // Metadata extra info such as label or other meta data
		UnHealthy bool
	}

	Hash struct {
		ReplicaNum     int   `yaml:"replica_num" json:"replica_num"`
		MaxVnodeNum    int32 `yaml:"max_vnode_num" json:"max_vnode_num"`
		ConsistentHash *consistent.Consistent
	}
)

func (c *ClusterConfig) GetEndpoint(mustHealth bool) []*Endpoint {
	var endpoints []*Endpoint
	for _, e := range c.Endpoints {
		// select all endpoint or endpoint is health
		if !mustHealth || !e.UnHealthy {
			endpoints = append(endpoints, e)
		}
	}
	return endpoints
}

func (c *ClusterConfig) CreateConsistentHash() {
	if c.LbStr == LoadBalanceConsistentHashing {
		h := consistent.NewConsistentHash(consistent.WithReplicaNum(c.Hash.ReplicaNum),
			consistent.WithMaxVnodeNum(int(c.Hash.MaxVnodeNum)))
		for _, endpoint := range c.Endpoints {
			h.Add(endpoint.Address.Address)
		}
		c.Hash.ConsistentHash = h
	}
}

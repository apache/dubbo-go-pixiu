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

// Cluster a single upstream cluster
type Cluster struct {
	Name             string              `yaml:"name" json:"name"`             // Name the cluster unique name
	TypeStr          string              `yaml:"type" json:"type"`             // Type the cluster discovery type string value
	Type             DiscoveryType       `yaml:",omitempty" json:",omitempty"` // Type the cluster discovery type
	EdsClusterConfig EdsClusterConfig    `yaml:"eds_cluster_config" json:"eds_cluster_config" mapstructure:"eds_cluster_config"`
	LbStr            string              `yaml:"lb_policy" json:"lb_policy"`   // Lb the cluster select node used loadBalance policy
	Lb               LbPolicy            `yaml:",omitempty" json:",omitempty"` // Lb the cluster select node used loadBalance policy
	HealthChecks     []HealthCheck       `yaml:"health_checks" json:"health_checks"`
	Hosts            []Address           `yaml:"hosts" json:"hosts"` // Hosts whe discovery type is Static, StrictDNS or LogicalDns, this need config
	Registries       map[string]Registry `yaml:"registries" json:"registries"`
}

// DiscoveryType
type DiscoveryType int32

const (
	Static DiscoveryType = 0 + iota
	StrictDNS
	LogicalDns
	EDS
	OriginalDst
)

// DiscoveryTypeName
var DiscoveryTypeName = map[int32]string{
	0: "Static",
	1: "StrictDNS",
	2: "LogicalDns",
	3: "EDS",
	4: "OriginalDst",
}

// DiscoveryTypeValue
var DiscoveryTypeValue = map[string]int32{
	"Static":      0,
	"StrictDNS":   1,
	"LogicalDns":  2,
	"EDS":         3,
	"OriginalDst": 4,
}

// EdsClusterConfig
type EdsClusterConfig struct {
	EdsConfig   ConfigSource `yaml:"eds_config" json:"eds_config" mapstructure:"eds_config"`
	ServiceName string       `yaml:"service_name" json:"service_name" mapstructure:"service_name"`
}

// Registry remote registry where dubbo apis are registered.
type Registry struct {
	Protocol     string `yaml:"protocol" json:"protocol" default:"zookeeper"`
	Timeout      string `yaml:"timeout" json:"timeout"`
	Address      string `yaml:"address" json:"address"`
	Username     string `yaml:"username" json:"username"`
	Password     string `yaml:"password" json:"password"`
	ServicesPath string `yaml:"services_path,omitempty" json:"services_path,omitempty"`
}

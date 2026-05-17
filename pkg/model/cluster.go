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
	"crypto/sha256"
	"fmt"
)

const generatedEndpointIDPrefix = "pixiu-generated-endpoint-"

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
		ConsistentHash       ConsistentHash      `yaml:"consistent" json:"consistent"` // Consistent hash config info
		HealthChecks         []HealthCheckConfig `yaml:"health_checks" json:"health_checks"`
		Endpoints            []*Endpoint         `yaml:"endpoints" json:"endpoints"`
		PrePickEndpointIndex uint32              `yaml:"-" json:"-"` // runtime-only round-robin cursor state
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
		Protocol     string `default:"zookeeper" yaml:"protocol" json:"protocol"`
		Timeout      string `yaml:"timeout" json:"timeout"`
		Address      string `yaml:"address" json:"address"`
		Username     string `yaml:"username" json:"username"`
		Password     string `yaml:"password" json:"password"`
		Group        string `default:"DEFAULT_GROUP"  yaml:"group" json:"group"`
		Namespace    string `yaml:"namespace" json:"namespace"`
		RegistryType string `default:"interface" yaml:"registry_type" json:"registry_type"` // "application", "interface"
	}

	// DiscoveryType
	DiscoveryType int32

	// Endpoint
	Endpoint struct {
		ID        string            `yaml:"ID" json:"ID"`                                                       // ID indicate one endpoint
		Name      string            `yaml:"name" json:"name"`                                                   // Name the endpoint unique name
		Address   SocketAddress     `yaml:"socket_address" json:"socket_address" mapstructure:"socket_address"` // Address socket address
		Metadata  map[string]string `yaml:"meta" json:"meta" mapstructure:"meta"`                               // Metadata extra info such as label or other meta data
		UnHealthy bool

		LLMMeta *LLMMeta `yaml:"llm_meta" json:"llm_meta" mapstructure:"llm_meta"` // LLMMeta extra info such as label or other meta data
	}

	// ConsistentHash methods include: RingHash, MaglevHash
	ConsistentHash struct {
		ReplicaNum      int   `yaml:"replica_num" json:"replica_num"`
		MaxVnodeNum     int32 `yaml:"max_vnode_num" json:"max_vnode_num"`
		MaglevTableSize int   `yaml:"maglev_table_size" json:"maglev_table_size"`
		Hash            LbConsistentHash
	}
)

func (c *ClusterConfig) GetEndpoint(mustHealth bool) []*Endpoint {
	var endpoints = make([]*Endpoint, 0, len(c.Endpoints))
	for _, e := range c.Endpoints {
		// select all endpoint or endpoint is health
		if !mustHealth || !e.UnHealthy {
			endpoints = append(endpoints, e)
		}
	}
	return endpoints
}

// CreateConsistentHash creates consistent hashing algorithms for Load Balance.
func (c *ClusterConfig) CreateConsistentHash() {
	if newConsistentHash, ok := ConsistentHashInitMap[c.LbStr]; ok {
		c.ConsistentHash.Hash = newConsistentHash(c.ConsistentHash, c.Endpoints)
	}
}

func (e Endpoint) GetHost() string {
	return fmt.Sprintf("%s:%d", e.Address.Address, e.Address.Port)
}

// GenerateEndpointID returns a deterministic runtime identity for endpoints
// that do not provide an explicit ID. The hash material includes cluster name,
// endpoint address, and (when present) LLM provider + API key, so endpoints
// that differ only by credential do not collide and the same endpoint hashes
// to the same ID across process restarts.
//
// Design notes:
//   - The output is sha256(material) truncated to the first 8 bytes (64 bits).
//     Birthday-collision probability becomes meaningful only around 2^32
//     endpoints, far above any realistic per-cluster scale.
//   - The output uses the pixiu-generated-endpoint- prefix so generated IDs
//     are recognizable in logs and dashboards without implying they only come
//     from the LLM registry path. Static-config endpoints use this helper too.
//   - clusterName is part of the material so endpoints in different clusters
//     never alias. Callers from the Nacos LLM path supply
//     instance.Metadata["cluster"]; if that metadata is missing the value is
//     the empty string, and identity is then determined by address +
//     credential only.
//   - The LLM API key is included on purpose so two endpoints that share an
//     address but use different credentials never alias to the same ID. The
//     output is one-way (sha256), so the raw key never appears in the ID.
func GenerateEndpointID(clusterName string, endpoint *Endpoint) string {
	sum := sha256.Sum256([]byte(endpointIDMaterial(clusterName, endpoint)))
	return fmt.Sprintf("%s%x", generatedEndpointIDPrefix, sum[:8])
}

// endpointIDMaterial builds the byte string fed into the hash inside
// GenerateEndpointID. Each component is tagged and length-prefixed so field
// boundaries remain unambiguous even if values contain punctuation or newlines.
//
// Contract: this function MUST NOT depend on endpoint.Name. Callers
// (notably ClusterStore.assembleClusterEndpoints) rely on being able to
// derive the ID before assigning a default Name. Adding Name into the
// material would also break the rename-invariance guarantee asserted by
// TestGenerateEndpointIDIgnoresEndpointName.
func endpointIDMaterial(clusterName string, endpoint *Endpoint) string {
	address := ""
	provider := ""
	apiKey := ""
	if endpoint != nil {
		address = endpoint.Address.GetAddress()
	}
	if endpoint != nil && endpoint.LLMMeta != nil {
		provider = endpoint.LLMMeta.Provider
		apiKey = endpoint.LLMMeta.APIKey
	}
	return endpointIDMaterialField("cluster", clusterName) +
		endpointIDMaterialField("address", address) +
		endpointIDMaterialField("llm_provider", provider) +
		endpointIDMaterialField("llm_api_key", apiKey)
}

func endpointIDMaterialField(name, value string) string {
	return fmt.Sprintf("%s:%d:%s\n", name, len(value), value)
}

// CloneEndpoints returns a deep copy of endpoints. Nil input is preserved.
func CloneEndpoints(endpoints []*Endpoint) []*Endpoint {
	if endpoints == nil {
		return nil
	}
	cloned := make([]*Endpoint, len(endpoints))
	for i, endpoint := range endpoints {
		cloned[i] = CloneEndpoint(endpoint)
	}
	return cloned
}

// CloneEndpoint returns a deep copy of endpoint with independent Address,
// Metadata, and LLMMeta. Returns nil for nil input. Snapshot consumers clone
// before handing endpoints to callers so downstream mutation does not leak
// back into the runtime snapshot.
//
// Cost: O(depth) — LLMMeta.RetryPolicy.Config is recursively cloned via
// cloneAnyMap, which allocates per nested map/slice. The request path
// should clone at most once per pick (typically when returning the chosen
// endpoint to the caller); avoid CloneEndpoint inside per-iteration loops
// over a snapshot's endpoint slice. Use HealthyEndpointsForPick to scan
// without cloning, then clone the single chosen endpoint.
func CloneEndpoint(endpoint *Endpoint) *Endpoint {
	if endpoint == nil {
		return nil
	}
	cloned := *endpoint
	cloned.Address = cloneSocketAddress(endpoint.Address)
	cloned.Metadata = cloneMetadata(endpoint.Metadata)
	cloned.LLMMeta = cloneLLMMeta(endpoint.LLMMeta)
	return &cloned
}

func cloneSocketAddress(address SocketAddress) SocketAddress {
	cloned := address
	if address.Domains != nil {
		cloned.Domains = append([]string(nil), address.Domains...)
	}
	return cloned
}

func cloneMetadata(metadata map[string]string) map[string]string {
	if metadata == nil {
		return nil
	}
	cloned := make(map[string]string, len(metadata))
	for key, value := range metadata {
		cloned[key] = value
	}
	return cloned
}

func cloneLLMMeta(meta *LLMMeta) *LLMMeta {
	if meta == nil {
		return nil
	}
	cloned := *meta
	cloned.RetryPolicy.Config = cloneAnyMap(meta.RetryPolicy.Config)
	return &cloned
}

func cloneAnyMap(input map[string]any) map[string]any {
	if input == nil {
		return nil
	}
	cloned := make(map[string]any, len(input))
	for key, value := range input {
		cloned[key] = cloneAnyValue(value)
	}
	return cloned
}

func cloneAnyValue(value any) any {
	switch typed := value.(type) {
	case map[string]any:
		return cloneAnyMap(typed)
	case []any:
		cloned := make([]any, len(typed))
		for i, item := range typed {
			cloned[i] = cloneAnyValue(item)
		}
		return cloned
	case []string:
		return append([]string(nil), typed...)
	case []int:
		return append([]int(nil), typed...)
	case []int64:
		return append([]int64(nil), typed...)
	case []float64:
		return append([]float64(nil), typed...)
	case []bool:
		return append([]bool(nil), typed...)
	case map[string]string:
		return cloneMetadata(typed)
	default:
		return value
	}
}

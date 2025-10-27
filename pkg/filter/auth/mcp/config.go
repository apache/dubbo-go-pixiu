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

package mcp

import (
	"fmt"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/filter/auth/mcp/internal/validator"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// Config defines the MCP auth filter configuration.
// It wires resource metadata (RFC9728), JWT providers, and path-based auth rules.
type Config struct {
	// ResourceMetadata controls /.well-known/oauth-protected-resource exposure and values.
	ResourceMetadata ResourceMetadata `yaml:"resource_metadata" json:"resource_metadata" mapstructure:"resource_metadata"`

	// Providers declares JWT validation providers reused by rules.
	Providers []validator.Provider `yaml:"providers" json:"providers" mapstructure:"providers"`

	// Rules binds request paths to a provider and required scopes.
	Rules []Rule `yaml:"rules" json:"rules" mapstructure:"rules"`
}

// ResourceMetadata represents OAuth 2.0 Protected Resource Metadata (RFC9728)
// that MCP clients discover via /.well-known/oauth-protected-resource
type ResourceMetadata struct {
	// Path is the well-known endpoint path to serve metadata from.
	// Default: "/.well-known/oauth-protected-resource"
	Path string `yaml:"path" json:"path" mapstructure:"path"`

	// Resource is the canonical resource identifier (RFC8707) that clients
	// should request tokens for (e.g. "https://mcp.example.com").
	Resource string `yaml:"resource" json:"resource" mapstructure:"resource"`

	// AuthorizationServers lists candidate Authorization Server metadata endpoints
	// (e.g. "https://auth.example.com/.well-known/oauth-authorization-server").
	AuthorizationServers []string `yaml:"authorization_servers" json:"authorization_servers" mapstructure:"authorization_servers"`
}

// Rule describes how to protect requests under a given path prefix.
type Rule struct {
	// Cluster is the route cluster name matched by the framework router.
	// The MCP filter will protect routes that resolve to this cluster.
	Cluster string `yaml:"cluster" json:"cluster" mapstructure:"cluster"`
}

// Validate performs basic semantic checks on the configuration.
func (c *Config) Validate() error {
	// Resource metadata
	if c.ResourceMetadata.Path == "" {
		c.ResourceMetadata.Path = "/.well-known/oauth-protected-resource"
		logger.Warnf("[dubbo-go-pixiu] resource_metadata.path is not set, using default value: %s", c.ResourceMetadata.Path)
	}
	if c.ResourceMetadata.Resource == "" {
		return fmt.Errorf("resource_metadata.resource must be set to the canonical MCP server URI")
	}
	if len(c.ResourceMetadata.AuthorizationServers) == 0 {
		return fmt.Errorf("resource_metadata.authorization_servers must not be empty")
	}

	// Providers presence
	if len(c.Providers) == 0 {
		return fmt.Errorf("providers must not be empty")
	}

	// Validate provider entries and index names to detect duplicates
	providerNames := make(map[string]struct{}, len(c.Providers))
	for _, p := range c.Providers {
		if p.Name == "" {
			return fmt.Errorf("provider name must not be empty")
		}
		if p.Audience == "" {
			p.Audience = c.ResourceMetadata.Resource
			logger.Warnf("[dubbo-go-pixiu] provider '%s' has no audience; defaulting to resource_metadata.resource '%s'", p.Name, c.ResourceMetadata.Resource)
		}
		if p.Issuer == "" {
			return fmt.Errorf("provider '%s': issuer must not be empty", p.Name)
		}
		if p.JWKS == "" {
			return fmt.Errorf("provider '%s': jwks must not be empty", p.Name)
		}
		if _, exists := providerNames[p.Name]; exists {
			return fmt.Errorf("duplicated provider name '%s'", p.Name)
		}
		providerNames[p.Name] = struct{}{}
	}

	// Rules
	for idx, r := range c.Rules {
		if r.Cluster == "" {
			return fmt.Errorf("rules[%d].cluster must not be empty", idx)
		}
	}

	return nil
}

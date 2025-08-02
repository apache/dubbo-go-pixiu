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
	"github.com/apache/dubbo-go-pixiu/pkg/filter/auth/mcp/internal/validator"
)

type (
	// Config defines the configuration for the MCP Auth filter.
	Config struct {
		// Providers lists the trusted JWT providers for this filter.
		Providers []validator.Provider `yaml:"providers" json:"providers"`
		// ResourceMetadata configures the /.well-known/oauth-protected-resource endpoint.
		ResourceMetadata ResourceMetadata `yaml:"resource_metadata" json:"resource_metadata"`
		// Rules define which paths require authentication and what scopes they need.
		Rules []Rule `yaml:"rules" json:"rules"`
	}

	// ResourceMetadata defines the content for the metadata endpoint.
	ResourceMetadata struct {
		// Enabled controls whether to serve the metadata endpoint.
		Enabled bool `yaml:"enabled" json:"enabled"`
		// ServerHost is the public-facing hostname of the gateway, used to build the
		// WWW-Authenticate header. e.g., "https://api.example.com"
		ServerHost string `yaml:"server_host" json:"server_host"`
		// AuthorizationServers lists the authorization servers to be included in the metadata response.
		AuthorizationServers []AuthorizationServer `yaml:"authorization_servers" json:"authorization_servers"`
	}

	// AuthorizationServer corresponds to an entry in the "authorization_servers" array.
	AuthorizationServer struct {
		Issuer string `yaml:"issuer" json:"issuer"`
	}

	// Rule defines an authentication/authorization rule for a specific path.
	Rule struct {
		Match Match `yaml:"match" json:"match"`
		// RequiredScopes lists the scopes required to access this path.
		// The token must contain ALL of these scopes.
		RequiredScopes []string `yaml:"required_scopes" json:"required_scopes"`
		// ProviderName specifies which provider to use for validating the token for this rule.
		ProviderName string `yaml:"provider_name" json:"provider_name"`
	}

	// Match defines the path matching rule.
	Match struct {
		Prefix string `yaml:"prefix" json:"prefix"`
	}
)

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

package validator

// InternalValidatorConfig represents the internal configuration for the JWT validator
type InternalValidatorConfig struct {
	// Providers is the list of JWT providers (using external Provider type)
	Providers []Provider `yaml:"providers" json:"providers"`
}

// Provider represents a JWT provider configuration (internal use)
type Provider struct {
	// Name is the unique identifier for this provider
	Name string `yaml:"name" json:"name" mapstructure:"name"`

	// Issuer is the JWT issuer identifier
	Issuer string `yaml:"issuer" json:"issuer" mapstructure:"issuer"`

	// Audience is the single valid audience value
	Audience string `yaml:"audience" json:"audience" mapstructure:"audience"`

	// JWKSSource defines how to obtain the JWKS
	JWKSSource JWKSSource `yaml:"jwks_source" json:"jwks_source" mapstructure:"jwks_source"`
}

// JWKSSource defines the source of JWKS
type JWKSSource struct {
	// Remote defines remote JWKS configuration
	Remote *RemoteJWKS `yaml:"remote" json:"remote" mapstructure:"remote"`

	// Local defines local JWKS configuration
	Local *LocalJWKS `yaml:"local" json:"local" mapstructure:"local"`
}

// RemoteJWKS defines remote JWKS configuration with jwx v3 support
type RemoteJWKS struct {
	// URI is the JWKS endpoint URL
	URI string `yaml:"uri" json:"uri" mapstructure:"uri"`

	// RefreshInterval is the interval for automatic JWKS refresh (jwx v3)
	RefreshInterval string `yaml:"refresh_interval" json:"refresh_interval" mapstructure:"refresh_interval" default:"15m"`

	// CacheTTL is the cache time-to-live for JWKS (jwx v3)
	CacheTTL string `yaml:"cache_ttl" json:"cache_ttl" mapstructure:"cache_ttl" default:"1h"`

	// Timeout is the HTTP request timeout
	Timeout string `yaml:"timeout" json:"timeout" mapstructure:"timeout" default:"5s"`
}

// LocalJWKS defines local JWKS configuration
type LocalJWKS struct {
	// InlineString is the JWKS JSON string
	InlineString string `yaml:"inline_string" json:"inline_string" mapstructure:"inline_string"`

	// FilePath is the path to JWKS file
	FilePath string `yaml:"file_path" json:"file_path" mapstructure:"file_path"`
}

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

package saml

import (
	"fmt"
)

// Config describes the configuration for SAML authentication filter
type Config struct {
	// Service Provider (SP) configuration
	EntityID             string `yaml:"entity_id" json:"entity_id" mapstructure:"entity_id"`          // SP entity ID
	AssertionConsumerURL string `yaml:"acs_url" json:"acs_url" mapstructure:"acs_url"`                // Assertion Consumer Service URL
	MetadataURL          string `yaml:"metadata_url" json:"metadata_url" mapstructure:"metadata_url"` // SP metadata endpoint

	// Identity Provider (IdP) configuration
	IdPMetadataURL  string `yaml:"idp_metadata_url" json:"idp_metadata_url" mapstructure:"idp_metadata_url"`    // IdP metadata URL
	IdPMetadataFile string `yaml:"idp_metadata_file" json:"idp_metadata_file" mapstructure:"idp_metadata_file"` // IdP metadata file path

	// Certificate configuration
	CertFile string `yaml:"cert_file" json:"cert_file" mapstructure:"cert_file"` // SP certificate file
	KeyFile  string `yaml:"key_file" json:"key_file" mapstructure:"key_file"`    // SP private key file

	// Routing rules
	Rules []Rule `yaml:"rules" json:"rules" mapstructure:"rules"` // URL matching rules

	// AllowIDPInitiated skips InResponseTo validation. Required for HTTP (non-TLS)
	// testing because SAML request tracking cookies need Secure + SameSite=None
	// which browsers reject over plain HTTP. In production with HTTPS, set to false.
	AllowIDPInitiated bool `yaml:"allow_idp_initiated" json:"allow_idp_initiated" mapstructure:"allow_idp_initiated"`

	// Attribute forwarding to backend services
	ForwardAttributes []ForwardAttribute `yaml:"forward_attributes" json:"forward_attributes" mapstructure:"forward_attributes"`

	// Error message
	ErrMsg string `yaml:"err_msg" json:"err_msg" mapstructure:"err_msg"` // Custom error message
}

// Rule defines URL matching rules for SAML authentication
type Rule struct {
	Match Match `yaml:"match" json:"match" mapstructure:"match"` // URL matching pattern
}

// Match defines the URL pattern to match
type Match struct {
	Prefix string `yaml:"prefix" json:"prefix" mapstructure:"prefix"` // URL prefix to match
}

// ForwardAttribute maps a SAML assertion attribute to an HTTP request header.
type ForwardAttribute struct {
	SAMLAttribute string `yaml:"saml_attribute" json:"saml_attribute" mapstructure:"saml_attribute"` // SAML attribute name
	Header        string `yaml:"header" json:"header" mapstructure:"header"`                         // HTTP header name
}

func (cfg *Config) Validate() error {
	if cfg.EntityID == "" {
		return fmt.Errorf("entity_id is required")
	}
	if cfg.AssertionConsumerURL == "" {
		return fmt.Errorf("acs_url is required")
	}
	if cfg.MetadataURL == "" {
		return fmt.Errorf("metadata_url is required")
	}
	if cfg.IdPMetadataURL == "" && cfg.IdPMetadataFile == "" {
		return fmt.Errorf("either idp_metadata_url or idp_metadata_file is required")
	}
	if cfg.CertFile == "" {
		return fmt.Errorf("cert_file is required")
	}
	if cfg.KeyFile == "" {
		return fmt.Errorf("key_file is required")
	}
	return nil
}

// DeepCopy returns a new independent copy of Config.
func (cfg *Config) DeepCopy() *Config {
	if cfg == nil {
		return nil
	}
	cp := *cfg

	if cfg.Rules != nil {
		cp.Rules = make([]Rule, len(cfg.Rules))
		copy(cp.Rules, cfg.Rules)
	}

	if cfg.ForwardAttributes != nil {
		cp.ForwardAttributes = make([]ForwardAttribute, len(cfg.ForwardAttributes))
		copy(cp.ForwardAttributes, cfg.ForwardAttributes)
	}

	return &cp
}

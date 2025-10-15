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

// Config represents the configuration for the JWT validator
type Config struct {
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

	// JWKS is a single URI-like string that specifies how to obtain JWKS
	// Supported schemes:
	//  - http(s)://...  (remote JWKS, uses default timeout)
	//  - file:///abs/path/jwks.json  (local file)
	JWKS string `yaml:"jwks" json:"jwks" mapstructure:"jwks"`
}

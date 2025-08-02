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

type (
	// Provider represents a single trusted JWT issuer.
	Provider struct {
		// Name is a unique identifier for this provider.
		Name string `yaml:"name" json:"name"`
		// Issuer is the required "iss" claim for the JWT.
		Issuer string `yaml:"issuer" json:"issuer"`
		// Audiences is a list of required "aud" claims. Validation will pass
		// if the token's "aud" claim contains ANY of these values.
		Audiences []string `yaml:"audiences" json:"audiences"`
		// JwksSource specifies the location of the JSON Web Key Set.
		JwksSource JwksSource `yaml:"jwks_source" json:"jwks_source"`
	}

	// JwksSource defines where to fetch the JWKS from.
	// Only one of its fields should be set.
	JwksSource struct {
		// Remote specifies a remote JWKS endpoint.
		Remote *Remote `yaml:"remote,omitempty" json:"remote,omitempty"`
		// Local specifies an inline JWKS.
		Local *Local `yaml:"local,omitempty" json:"local,omitempty"`
	}

	// Remote defines configuration for a remote JWKS endpoint.
	Remote struct {
		// URI is the URL for the remote JWKS endpoint.
		URI string `yaml:"uri" json:"uri"`
	}

	// Local defines an inline JWKS.
	Local struct {
		// Inline is the raw JSON string of the JWKS.
		Inline string `yaml:"inline" json:"inline"`
	}
)

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

package jwt

import (
	"github.com/MicahParks/keyfunc"
)

type (
	// FromHeaders Get the token from a field in the header，default Authorization: Bearer <token>
	FromHeaders struct {
		Name        string `default:"Authorization" yaml:"name" json:"name" mapstructure:"name"`                   // header key
		ValuePrefix string `default:"Bearer " yaml:"value_prefix" json:"value_prefix" mapstructure:"value_prefix"` // header value
	}

	// Rules router match
	Rules struct {
		Match    Match    `yaml:"match" json:"match" mapstructure:"match"`          // router
		Requires Requires `yaml:"requires" json:"requires" mapstructure:"requires"` // router requires jwks check
	}

	Match struct {
		Prefix string `yaml:"prefix" json:"prefix" mapstructure:"prefix"` // url
	}

	Requires struct {
		RequiresAny Requirement   `yaml:"requires_any" json:"requires_any" mapstructure:"requires_any"` // single token check
		RequiresAll []Requirement `yaml:"requires_all" json:"requires_all" mapstructure:"requires_all"` // multiple token check,will only verify one, if the token fails, continue to the next check
	}

	Requirement struct {
		ProviderName string `yaml:"provider_name" json:"provider_name" mapstructure:"provider_name"` // jwks providers name
	}

	Providers struct {
		Name                 string      `yaml:"name" json:"name" mapstructure:"name"`                                                       // jwt name
		ForwardPayloadHeader string      `yaml:"forward_payload_header" json:"forward_payload_header" mapstructure:"forward_payload_header"` // header add issuer
		FromHeaders          FromHeaders `yaml:"from_headers" json:"from_headers" mapstructure:"from_headers"`                               // from header get token
		Issuer               string      `yaml:"issuer" json:"issuer" mapstructure:"issuer"`                                                 // jwt issuer
		Local                *Local      `yaml:"local_jwks" json:"local_jwks" mapstructure:"local_jwks"`                                     // local jwks
		Remote               *Remote     `yaml:"remote_jwks" json:"remote_jwks" mapstructure:"remote_jwks"`                                  // remote jwks
	}

	Local struct {
		InlineString string `yaml:"inline_string" json:"inline_string" mapstructure:"inline_string"` // local jwks public key
	}

	Remote struct {
		HttpURI HttpURI `yaml:"http_uri" json:"http_uri" mapstructure:"http_uri"` // remote jwks public key
	}

	HttpURI struct {
		Uri     string `yaml:"uri" json:"uri" mapstructure:"uri"`
		Cluster string `yaml:"cluster" json:"cluster" mapstructure:"cluster"`
		TimeOut string `default:"5s" yaml:"timeout" json:"timeout" mapstructure:"timeout"`
	}
)

type Provider struct {
	jwk                  *keyfunc.JWKS
	issuer               string
	forwardPayloadHeader string
	headers              FromHeaders
}

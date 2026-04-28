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
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

func TestConfigValidate(t *testing.T) {
	tests := []struct {
		name    string
		cfg     Config
		wantErr string
	}{
		{
			name: "missing entity id",
			cfg: func() Config {
				cfg := validConfig()
				cfg.EntityID = ""
				return cfg
			}(),
			wantErr: "entity_id is required",
		},
		{
			name: "missing acs url",
			cfg: func() Config {
				cfg := validConfig()
				cfg.AssertionConsumerURL = ""
				return cfg
			}(),
			wantErr: "acs_url is required",
		},
		{
			name: "missing metadata url",
			cfg: func() Config {
				cfg := validConfig()
				cfg.MetadataURL = ""
				return cfg
			}(),
			wantErr: "metadata_url is required",
		},
		{
			name: "missing idp metadata source",
			cfg: func() Config {
				cfg := validConfig()
				cfg.IdPMetadataURL = ""
				cfg.IdPMetadataFile = ""
				return cfg
			}(),
			wantErr: "either idp_metadata_url or idp_metadata_file is required",
		},
		{
			name: "missing cert file",
			cfg: func() Config {
				cfg := validConfig()
				cfg.CertFile = ""
				return cfg
			}(),
			wantErr: "cert_file is required",
		},
		{
			name: "missing key file",
			cfg: func() Config {
				cfg := validConfig()
				cfg.KeyFile = ""
				return cfg
			}(),
			wantErr: "key_file is required",
		},
		{
			name: "valid config",
			cfg:  validConfig(),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := tt.cfg.Validate()
			if tt.wantErr == "" {
				assert.NoError(t, err)
				return
			}
			assert.EqualError(t, err, tt.wantErr)
		})
	}
}

func validConfig() Config {
	return Config{
		EntityID:             "pixiu-sp",
		AssertionConsumerURL: "http://localhost:8888/saml/acs",
		MetadataURL:          "http://localhost:8888/saml/metadata",
		IdPMetadataURL:       "http://localhost:9000/metadata",
		CertFile:             "sp.crt",
		KeyFile:              "sp.key",
	}
}

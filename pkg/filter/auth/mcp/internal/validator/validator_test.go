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

import (
	"os"
	"path/filepath"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func tempJWKSFileURL(t *testing.T) string {
	t.Helper()
	dir := t.TempDir()
	p := filepath.Join(dir, "jwks.json")
	if err := os.WriteFile(p, []byte(`{"keys":[]}`), 0644); err != nil {
		t.Fatalf("write temp jwks: %v", err)
	}
	return "file://" + p
}

func TestNewValidator(t *testing.T) {
	tests := []struct {
		name    string
		config  Config
		wantErr bool
	}{
		{
			name:    "empty providers",
			config:  Config{Providers: []Provider{}},
			wantErr: true,
		},
		{
			name: "valid provider with local JWKS",
			config: Config{
				Providers: []Provider{
					{
						Name:     "test-provider",
						Issuer:   "https://test.issuer.com",
						Audience: "test-audience",
						JWKS:     tempJWKSFileURL(t),
					},
				},
			},
			wantErr: false,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			validator, err := NewValidator(tt.config)
			if tt.wantErr {
				assert.Error(t, err)
				assert.Nil(t, validator)
			} else {
				assert.NoError(t, err)
				assert.NotNil(t, validator)
			}
		})
	}
}

func TestValidator_ListProviders(t *testing.T) {
	config := Config{
		Providers: []Provider{
			{
				Name:     "provider1",
				Issuer:   "https://issuer1.com",
				Audience: "audience1",
				JWKS:     tempJWKSFileURL(t),
			},
			{
				Name:     "provider2",
				Issuer:   "https://issuer2.com",
				Audience: "audience2",
				JWKS:     tempJWKSFileURL(t),
			},
		},
	}

	validator, err := NewValidator(config)
	require.NoError(t, err)

	providers := validator.Providers()
	assert.Len(t, providers, 2)
	assert.Contains(t, providers, "provider1")
	assert.Contains(t, providers, "provider2")
}

func TestValidator_GetProvider(t *testing.T) {
	config := Config{
		Providers: []Provider{
			{
				Name:     "test-provider",
				Issuer:   "https://test.issuer.com",
				Audience: "test-audience",
				JWKS:     tempJWKSFileURL(t),
			},
		},
	}

	validator, err := NewValidator(config)
	require.NoError(t, err)

	// Test existing provider
	provider, exists := validator.Provider("test-provider")
	assert.True(t, exists)
	assert.Equal(t, "test-provider", provider.Name)
	assert.Equal(t, "https://test.issuer.com", provider.Issuer)

	// Test non-existing provider
	provider, exists = validator.Provider("non-existing")
	assert.False(t, exists)
	assert.Nil(t, provider)
}

func TestValidationError_Error(t *testing.T) {
	err := ValidationError{
		Code:    "invalid_token",
		Message: "token is expired",
	}

	expected := "invalid_token: token is expired"
	assert.Equal(t, expected, err.Error())
}

func TestProvider_Configuration(t *testing.T) {
	config := Config{
		Providers: []Provider{
			{
				Name:     "test-provider",
				Issuer:   "https://test.issuer.com",
				Audience: "test-audience",
				JWKS:     tempJWKSFileURL(t),
			},
		},
	}

	validator, err := NewValidator(config)
	require.NoError(t, err)

	provider, exists := validator.Provider("test-provider")
	assert.True(t, exists)
	assert.Equal(t, "test-provider", provider.Name)
	assert.Equal(t, "https://test.issuer.com", provider.Issuer)
	assert.Equal(t, "test-audience", provider.Audience)
}

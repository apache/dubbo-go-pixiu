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

package util

import (
	"strings"
	"testing"
)

func TestParseHostPortFromURL(t *testing.T) {
	tests := []struct {
		name             string
		input            string
		expectedHost     string
		expectedPort     int
		expectedFallback bool
		expectError      bool
	}{
		{
			name:             "complete HTTP URL",
			input:            "http://localhost:8080/path",
			expectedHost:     "localhost",
			expectedPort:     8080,
			expectedFallback: false,
			expectError:      false,
		},
		{
			name:             "HTTPS URL with fallback",
			input:            "https://example.com",
			expectedHost:     "example.com",
			expectedPort:     443,
			expectedFallback: true,
			expectError:      false,
		},
		{
			name:             "HTTP URL with fallback",
			input:            "http://test.com",
			expectedHost:     "test.com",
			expectedPort:     80,
			expectedFallback: true,
			expectError:      false,
		},
		{
			name:             "host:port format",
			input:            "example.com:9090",
			expectedHost:     "example.com",
			expectedPort:     9090,
			expectedFallback: false,
			expectError:      false,
		},
		{
			name:             "host only with fallback",
			input:            "nacos-server",
			expectedHost:     "nacos-server",
			expectedPort:     8080,
			expectedFallback: true,
			expectError:      false,
		},
		{
			name:        "empty string",
			input:       "",
			expectError: true,
		},
		{
			name:        "invalid port",
			input:       "localhost:abc",
			expectError: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := ParseHostPortFromURL(tt.input)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				return
			}

			if err != nil {
				t.Errorf("unexpected error: %v", err)
				return
			}

			if result.Host != tt.expectedHost {
				t.Errorf("expected host %s, got %s", tt.expectedHost, result.Host)
			}

			if result.Port != tt.expectedPort {
				t.Errorf("expected port %d, got %d", tt.expectedPort, result.Port)
			}

			if result.UsedFallback != tt.expectedFallback {
				t.Errorf("expected fallback %v, got %v", tt.expectedFallback, result.UsedFallback)
			}

			if tt.expectedFallback && result.FallbackInfo == "" {
				t.Errorf("expected fallback info but got empty string")
			}
		})
	}
}

func TestExtractPathFromURL(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "complete URL with path",
			input:    "http://localhost:8080/api/v1/test",
			expected: "/api/v1/test",
		},
		{
			name:     "URL without path",
			input:    "http://localhost:8080",
			expected: "/",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "/",
		},
		{
			name:     "path only",
			input:    "/api/test",
			expected: "/api/test",
		},
		{
			name:     "host:port/path format",
			input:    "localhost:8080/service",
			expected: "/service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractPathFromURL(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestReplaceGoTemplateArgsInPath(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "template with args",
			input:    "/api/{{.args.service}}/{{.args.method}}",
			expected: "/api/{service}/{method}",
		},
		{
			name:     "no template",
			input:    "/api/v1/test",
			expected: "/api/v1/test",
		},
		{
			name:     "empty path",
			input:    "",
			expected: "/",
		},
		{
			name:     "single template arg",
			input:    "/{{.args.id}}",
			expected: "/{id}",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ReplaceGoTemplateArgsInPath(tt.input)
			if result != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, result)
			}
		})
	}
}

func TestValidateNacosAddresses(t *testing.T) {
	tests := []struct {
		name        string
		addresses   string
		expectError bool
		errorMsg    string
	}{
		{
			name:        "valid single address",
			addresses:   "localhost:8848",
			expectError: false,
		},
		{
			name:        "valid multiple addresses",
			addresses:   "nacos1:8848,nacos2:8848,nacos3:8848",
			expectError: false,
		},
		{
			name:        "empty string",
			addresses:   "",
			expectError: true,
			errorMsg:    "nacos addresses cannot be empty",
		},
		{
			name:        "all invalid addresses",
			addresses:   "invalid1:abc,invalid2:,invalid3:invalid",
			expectError: true,
			errorMsg:    "no valid nacos addresses found",
		},
		{
			name:        "mixed valid and invalid",
			addresses:   "valid:8848,invalid:abc,another:8849",
			expectError: false, // Should succeed with warnings
		},
		{
			name:        "whitespace only",
			addresses:   "   ",
			expectError: true,
		},
		{
			name:        "valid address with URL format",
			addresses:   "http://nacos:8848",
			expectError: false,
		},
		{
			name:        "invalid port number",
			addresses:   "nacos:invalid_port",
			expectError: true,
			errorMsg:    "no valid nacos addresses found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := ValidateNacosAddresses(tt.addresses)

			if tt.expectError {
				if err == nil {
					t.Errorf("expected error but got none")
				}
				if tt.errorMsg != "" && !strings.Contains(err.Error(), tt.errorMsg) {
					t.Errorf("expected error message to contain '%s', got '%s'", tt.errorMsg, err.Error())
				}
			} else {
				if err != nil {
					t.Errorf("unexpected error: %v", err)
				}
			}
		})
	}
}

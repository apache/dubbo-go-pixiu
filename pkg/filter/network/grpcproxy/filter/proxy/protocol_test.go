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

package proxy

import (
	"testing"
)

func TestExtractTripleMetadata(t *testing.T) {
	tests := []struct {
		name        string
		attachments map[string]any
		expected    map[string]string
	}{
		{
			name: "with triple headers",
			attachments: map[string]any{
				"tri-service-version": "1.0.0",
				"tri-service-group":   "production",
				"tri-req-id":          "12345",
				"other-header":        "value",
			},
			expected: map[string]string{
				"tri-service-version": "1.0.0",
				"tri-service-group":   "production",
				"tri-req-id":          "12345",
			},
		},
		{
			name: "no triple headers",
			attachments: map[string]any{
				"content-type":  "application/grpc",
				"authorization": "Bearer token",
			},
			expected: map[string]string{},
		},
		{
			name:        "nil attachments",
			attachments: nil,
			expected:    map[string]string{},
		},
		{
			name:        "empty attachments",
			attachments: map[string]any{},
			expected:    map[string]string{},
		},
		{
			name: "mixed case headers",
			attachments: map[string]any{
				"Tri-Service-Version": "2.0.0",
				"TRI-SERVICE-GROUP":   "staging",
			},
			expected: map[string]string{
				"tri-service-version": "2.0.0",
				"tri-service-group":   "staging",
			},
		},
		{
			name: "non-string values ignored",
			attachments: map[string]any{
				"tri-service-version": "1.0.0",
				"tri-count":           123, // non-string, should be ignored
			},
			expected: map[string]string{
				"tri-service-version": "1.0.0",
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ExtractTripleMetadata(tt.attachments)
			if len(result) != len(tt.expected) {
				t.Errorf("ExtractTripleMetadata() returned %d items, want %d", len(result), len(tt.expected))
			}
			for k, v := range tt.expected {
				if result[k] != v {
					t.Errorf("ExtractTripleMetadata()[%q] = %q, want %q", k, result[k], v)
				}
			}
		})
	}
}

func TestIsTripleHeader(t *testing.T) {
	tests := []struct {
		name     string
		key      string
		expected bool
	}{
		{"lowercase tri prefix", "tri-service-version", true},
		{"uppercase TRI prefix", "TRI-SERVICE-VERSION", true},
		{"mixed case", "Tri-Service-Group", true},
		{"no tri prefix", "content-type", false},
		{"partial match", "triangle", false},
		{"empty string", "", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsTripleHeader(tt.key)
			if result != tt.expected {
				t.Errorf("IsTripleHeader(%q) = %v, want %v", tt.key, result, tt.expected)
			}
		})
	}
}

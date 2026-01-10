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
	"time"
)

func TestParseDurationWithDefault(t *testing.T) {
	tests := []struct {
		name        string
		durationStr string
		defaultVal  time.Duration
		expected    time.Duration
	}{
		{
			name:        "empty string uses default",
			durationStr: "",
			defaultVal:  30 * time.Second,
			expected:    30 * time.Second,
		},
		{
			name:        "valid duration seconds",
			durationStr: "10s",
			defaultVal:  30 * time.Second,
			expected:    10 * time.Second,
		},
		{
			name:        "valid duration minutes",
			durationStr: "5m",
			defaultVal:  30 * time.Second,
			expected:    5 * time.Minute,
		},
		{
			name:        "valid duration hours",
			durationStr: "2h",
			defaultVal:  30 * time.Second,
			expected:    2 * time.Hour,
		},
		{
			name:        "valid duration milliseconds",
			durationStr: "500ms",
			defaultVal:  1 * time.Second,
			expected:    500 * time.Millisecond,
		},
		{
			name:        "valid complex duration",
			durationStr: "1h30m",
			defaultVal:  30 * time.Second,
			expected:    1*time.Hour + 30*time.Minute,
		},
		{
			name:        "invalid duration uses default",
			durationStr: "invalid",
			defaultVal:  45 * time.Second,
			expected:    45 * time.Second,
		},
		{
			name:        "number without unit uses default",
			durationStr: "100",
			defaultVal:  20 * time.Second,
			expected:    20 * time.Second,
		},
		{
			name:        "negative duration",
			durationStr: "-5s",
			defaultVal:  10 * time.Second,
			expected:    -5 * time.Second,
		},
		{
			name:        "zero duration",
			durationStr: "0s",
			defaultVal:  10 * time.Second,
			expected:    0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := parseDurationWithDefault(tt.durationStr, tt.defaultVal)
			if result != tt.expected {
				t.Errorf("parseDurationWithDefault(%q, %v) = %v, want %v",
					tt.durationStr, tt.defaultVal, result, tt.expected)
			}
		})
	}
}

func TestPlugin_Kind(t *testing.T) {
	plugin := Plugin{}
	expected := Kind

	if got := plugin.Kind(); got != expected {
		t.Errorf("Kind() = %q, want %q", got, expected)
	}
}

func TestReflectionModeConstants(t *testing.T) {
	// Verify reflection mode constants are defined correctly
	tests := []struct {
		name     string
		mode     string
		expected string
	}{
		{"passthrough mode", ReflectionModePassthrough, "passthrough"},
		{"reflection mode", ReflectionModeReflection, "reflection"},
		{"hybrid mode", ReflectionModeHybrid, "hybrid"},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if tt.mode != tt.expected {
				t.Errorf("%s = %q, want %q", tt.name, tt.mode, tt.expected)
			}
		})
	}
}

func TestConfig_DefaultValues(t *testing.T) {
	// Test that config can be created with default values
	config := &Config{}

	// Verify default empty values
	if config.ReflectionMode != "" {
		t.Errorf("default ReflectionMode = %q, want empty", config.ReflectionMode)
	}
	if config.DescriptorCacheTTLStr != "" {
		t.Errorf("default DescriptorCacheTTLStr = %q, want empty", config.DescriptorCacheTTLStr)
	}
	if config.ExtractTripleMetadata {
		t.Error("default ExtractTripleMetadata should be false")
	}
}

func TestConfig_WithReflectionSettings(t *testing.T) {
	config := &Config{
		ReflectionMode:        ReflectionModeReflection,
		DescriptorCacheTTLStr: "10m",
		ExtractTripleMetadata: true,
	}

	if config.ReflectionMode != "reflection" {
		t.Errorf("ReflectionMode = %q, want %q", config.ReflectionMode, "reflection")
	}

	// Parse the TTL
	ttl := parseDurationWithDefault(config.DescriptorCacheTTLStr, 5*time.Minute)
	if ttl != 10*time.Minute {
		t.Errorf("parsed TTL = %v, want %v", ttl, 10*time.Minute)
	}

	if !config.ExtractTripleMetadata {
		t.Error("ExtractTripleMetadata should be true")
	}
}

func TestFilter_Close_NilReflectionManager(t *testing.T) {
	// Test that Close handles nil reflection manager gracefully
	filter := &Filter{
		Config:            &Config{},
		reflectionManager: nil,
	}

	// Should not panic
	filter.Close()
}

func TestFilter_Close_WithReflectionManager(t *testing.T) {
	rm := NewReflectionManager(5 * time.Minute)
	filter := &Filter{
		Config:            &Config{},
		reflectionManager: rm,
	}

	// Should not panic and should close the manager
	filter.Close()
}

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

package config

import (
	"encoding/json"
	"testing"
	"time"
)

func TestNewDefaultConfig(t *testing.T) {
	cfg := NewDefaultConfig()

	if cfg.LogLevel != DefaultLogLevel {
		t.Errorf("expected LogLevel %s, got %s", DefaultLogLevel, cfg.LogLevel)
	}

	if cfg.ControllerName != DefaultControllerName {
		t.Errorf("expected ControllerName %s, got %s", DefaultControllerName, cfg.ControllerName)
	}

	if cfg.LeaderElectionID != DefaultLeaderElectionID {
		t.Errorf("expected LeaderElectionID %s, got %s", DefaultLeaderElectionID, cfg.LeaderElectionID)
	}

	if cfg.Gateway == nil {
		t.Fatal("expected Gateway config to be non-nil")
	}

	if cfg.Gateway.Image != DefaultGatewayImage {
		t.Errorf("expected Gateway.Image %s, got %s", DefaultGatewayImage, cfg.Gateway.Image)
	}

	if cfg.Gateway.ImagePullPolicy != DefaultImagePullPolicy {
		t.Errorf("expected Gateway.ImagePullPolicy %s, got %s", DefaultImagePullPolicy, cfg.Gateway.ImagePullPolicy)
	}
}

func TestNewDefaultGatewayConfig(t *testing.T) {
	gw := NewDefaultGatewayConfig()

	if gw.Image != DefaultGatewayImage {
		t.Errorf("expected Image %s, got %s", DefaultGatewayImage, gw.Image)
	}

	if gw.ImagePullPolicy != DefaultImagePullPolicy {
		t.Errorf("expected ImagePullPolicy %s, got %s", DefaultImagePullPolicy, gw.ImagePullPolicy)
	}
}

func TestGatewayConfigJSON(t *testing.T) {
	gw := &GatewayConfig{
		Image:           "test/image:v1",
		ImagePullPolicy: "Always",
	}

	data, err := json.Marshal(gw)
	if err != nil {
		t.Fatalf("failed to marshal GatewayConfig: %v", err)
	}

	var parsed GatewayConfig
	if err := json.Unmarshal(data, &parsed); err != nil {
		t.Fatalf("failed to unmarshal GatewayConfig: %v", err)
	}

	if parsed.Image != gw.Image {
		t.Errorf("expected Image %s, got %s", gw.Image, parsed.Image)
	}

	if parsed.ImagePullPolicy != gw.ImagePullPolicy {
		t.Errorf("expected ImagePullPolicy %s, got %s", gw.ImagePullPolicy, parsed.ImagePullPolicy)
	}
}

func TestSetControllerConfig(t *testing.T) {
	// Note: This test modifies global state, so it cannot run in parallel
	// Save and restore to minimize impact on other tests
	original := ControllerConfig
	defer SetControllerConfig(original)

	newCfg := &Config{
		LogLevel:       "debug",
		ControllerName: "test-controller",
		Gateway: &GatewayConfig{
			Image:           "custom/image:test",
			ImagePullPolicy: "Never",
		},
	}

	SetControllerConfig(newCfg)

	if ControllerConfig.LogLevel != "debug" {
		t.Errorf("expected LogLevel debug, got %s", ControllerConfig.LogLevel)
	}

	if ControllerConfig.Gateway.Image != "custom/image:test" {
		t.Errorf("expected Gateway.Image custom/image:test, got %s", ControllerConfig.Gateway.Image)
	}
}

func TestTimeDurationJSON(t *testing.T) {
	tests := []struct {
		name     string
		duration time.Duration
		expected string
	}{
		{"seconds", 30 * time.Second, `"30s"`},
		{"minutes", 5 * time.Minute, `"5m0s"`},
		{"hours", 2 * time.Hour, `"2h0m0s"`},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			td := TimeDuration{Duration: tt.duration}
			data, err := json.Marshal(&td)
			if err != nil {
				t.Fatalf("failed to marshal: %v", err)
			}

			if string(data) != tt.expected {
				t.Errorf("expected %s, got %s", tt.expected, string(data))
			}

			var parsed TimeDuration
			if err := json.Unmarshal(data, &parsed); err != nil {
				t.Fatalf("failed to unmarshal: %v", err)
			}

			if parsed.Duration != tt.duration {
				t.Errorf("expected duration %v, got %v", tt.duration, parsed.Duration)
			}
		})
	}
}

func TestTimeDurationUnmarshalNumeric(t *testing.T) {
	// Test unmarshaling numeric value (nanoseconds)
	data := []byte(`1000000000`)
	var td TimeDuration
	if err := json.Unmarshal(data, &td); err != nil {
		t.Fatalf("failed to unmarshal numeric: %v", err)
	}

	if td.Duration != time.Second {
		t.Errorf("expected 1s, got %v", td.Duration)
	}
}

func TestNewLeaderElection(t *testing.T) {
	le := NewLeaderElection()

	if le.LeaseDuration.Duration != 30*time.Second {
		t.Errorf("expected LeaseDuration 30s, got %v", le.LeaseDuration.Duration)
	}

	if le.RenewDeadline.Duration != 20*time.Second {
		t.Errorf("expected RenewDeadline 20s, got %v", le.RenewDeadline.Duration)
	}

	if le.RetryPeriod.Duration != 2*time.Second {
		t.Errorf("expected RetryPeriod 2s, got %v", le.RetryPeriod.Duration)
	}

	if le.Disable != false {
		t.Error("expected Disable to be false")
	}
}

func TestValidateImagePullPolicy(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{"valid Always", "Always", "Always"},
		{"valid IfNotPresent", "IfNotPresent", "IfNotPresent"},
		{"valid Never", "Never", "Never"},
		{"empty string returns default", "", DefaultImagePullPolicy},
		{"invalid policy returns default", "InvalidPolicy", DefaultImagePullPolicy},
		{"lowercase invalid", "always", DefaultImagePullPolicy},
		{"mixed case invalid", "ALWAYS", DefaultImagePullPolicy},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := ValidateImagePullPolicy(tt.input)
			if result != tt.expected {
				t.Errorf("ValidateImagePullPolicy(%q) = %q, want %q", tt.input, result, tt.expected)
			}
		})
	}
}

func TestIsValidImagePullPolicy(t *testing.T) {
	tests := []struct {
		name     string
		input    string
		expected bool
	}{
		{"valid Always", "Always", true},
		{"valid IfNotPresent", "IfNotPresent", true},
		{"valid Never", "Never", true},
		{"empty string", "", false},
		{"invalid policy", "InvalidPolicy", false},
		{"lowercase always", "always", false},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := IsValidImagePullPolicy(tt.input)
			if result != tt.expected {
				t.Errorf("IsValidImagePullPolicy(%q) = %v, want %v", tt.input, result, tt.expected)
			}
		})
	}
}

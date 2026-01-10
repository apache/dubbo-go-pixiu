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

package dubbo

import (
	"testing"
)

func TestDubboProxyConfig_GetCluster(t *testing.T) {
	tests := []struct {
		name     string
		cluster  string
		expected string
	}{
		{
			name:     "empty - use default",
			cluster:  "",
			expected: "failover",
		},
		{
			name:     "explicit value",
			cluster:  "failfast",
			expected: "failfast",
		},
		{
			name:     "case sensitive - lowercase accepted",
			cluster:  "failover",
			expected: "failover",
		},
		{
			name:     "failsafe strategy",
			cluster:  "failsafe",
			expected: "failsafe",
		},
		{
			name:     "broadcast strategy",
			cluster:  "broadcast",
			expected: "broadcast",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &DubboProxyConfig{Cluster: tt.cluster}
			if got := config.GetCluster(); got != tt.expected {
				t.Errorf("GetCluster() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDubboProxyConfig_GetProtocol(t *testing.T) {
	tests := []struct {
		name     string
		protocol string
		expected string
	}{
		{
			name:     "empty - use default",
			protocol: "",
			expected: "dubbo",
		},
		{
			name:     "triple protocol - tri",
			protocol: "tri",
			expected: "tri",
		},
		{
			name:     "dubbo protocol",
			protocol: "dubbo",
			expected: "dubbo",
		},
		{
			name:     "rest protocol",
			protocol: "rest",
			expected: "rest",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &DubboProxyConfig{Protocol: tt.protocol}
			if got := config.GetProtocol(); got != tt.expected {
				t.Errorf("GetProtocol() = %v, want %v", got, tt.expected)
			}
		})
	}
}

func TestDubboProxyConfig_GetCheck(t *testing.T) {
	tests := []struct {
		name     string
		check    *bool
		expected *bool
	}{
		{
			name:     "nil - not configured",
			check:    nil,
			expected: nil,
		},
		{
			name:     "explicit false",
			check:    boolPtr(false),
			expected: boolPtr(false),
		},
		{
			name:     "explicit true",
			check:    boolPtr(true),
			expected: boolPtr(true),
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			config := &DubboProxyConfig{Check: tt.check}
			got := config.GetCheck()

			if (got == nil) != (tt.expected == nil) {
				t.Errorf("GetCheck() nil mismatch: got %v, want %v", got, tt.expected)
				return
			}

			if got != nil && tt.expected != nil && *got != *tt.expected {
				t.Errorf("GetCheck() = %v, want %v", *got, *tt.expected)
			}
		})
	}
}

func boolPtr(b bool) *bool {
	return &b
}

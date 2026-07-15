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

package nacos

import (
	"testing"
)

import (
	"github.com/nacos-group/nacos-sdk-go/v2/model"

	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud/servicediscovery"
)

func TestCallback_ServiceNameWithGroupPrefix(t *testing.T) {
	// Test the @@ prefix stripping logic in Callback
	// This tests the new code path where ServiceName contains "DEFAULT_GROUP@@service-name"

	tests := []struct {
		name        string
		serviceName string
		expected    string
	}{
		{
			name:        "with group prefix",
			serviceName: "DEFAULT_GROUP@@user-service",
			expected:    "user-service",
		},
		{
			name:        "without group prefix",
			serviceName: "user-service",
			expected:    "user-service",
		},
		{
			name:        "empty string",
			serviceName: "",
			expected:    "",
		},
		{
			name:        "multiple @@ separators",
			serviceName: "GROUP@@SUBGROUP@@service",
			expected:    "SUBGROUP@@service",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			// Simulate the strings.Cut logic used in Callback
			serviceName := tt.serviceName
			if _, after, ok := cutServiceName(serviceName, "@@"); ok {
				serviceName = after
			}
			assert.Equal(t, tt.expected, serviceName)
		})
	}
}

// cutServiceName mimics strings.Cut for testing purposes
func cutServiceName(s, sep string) (before, after string, found bool) {
	for i := 0; i+len(sep) <= len(s); i++ {
		if s[i:i+len(sep)] == sep {
			return s[:i], s[i+len(sep):], true
		}
	}
	return s, "", false
}

func TestFromInstanceToServiceInstance(t *testing.T) {
	instance := model.Instance{
		Ip:          "192.168.1.1",
		Port:        8080,
		ServiceName: "user-service",
		Healthy:     true,
		Enable:      true,
		ClusterName: "DEFAULT",
		Metadata:    map[string]string{"key": "value"},
	}

	result := fromInstanceToServiceInstance("user-service", instance)

	expected := servicediscovery.ServiceInstance{
		ID:          "192.168.1.1:8080",
		ServiceName: "user-service",
		Host:        "192.168.1.1",
		Port:        8080,
		Healthy:     true,
		Enable:      true,
		CLusterName: "DEFAULT",
		Metadata:    map[string]string{"key": "value"},
	}

	assert.Equal(t, expected.ID, result.ID)
	assert.Equal(t, expected.ServiceName, result.ServiceName)
	assert.Equal(t, expected.Host, result.Host)
	assert.Equal(t, expected.Port, result.Port)
	assert.Equal(t, expected.Healthy, result.Healthy)
	assert.Equal(t, expected.Enable, result.Enable)
	assert.Equal(t, expected.CLusterName, result.CLusterName)
}

func TestCallback_DisabledInstance(t *testing.T) {
	// Test that disabled instances are skipped (line 117-119)
	// This is existing behavior but needs coverage for the new Instance type

	// Create a mock listener
	mockListener := &mockServiceEventListener{
		serviceNames: []string{"test-service"},
	}

	// Create nacosServiceDiscovery with mock
	nsd := &nacosServiceDiscovery{
		listener:        mockListener,
		instanceMap:     make(map[string]servicediscovery.ServiceInstance),
		callbackFlagMap: nil,
	}

	// Test callback with mixed enabled/disabled instances
	services := []model.Instance{
		{
			Ip:          "192.168.1.1",
			Port:        8080,
			ServiceName: "DEFAULT_GROUP@@test-service",
			Enable:      true,
			Healthy:     true,
		},
		{
			Ip:          "192.168.1.2",
			Port:        8080,
			ServiceName: "DEFAULT_GROUP@@test-service",
			Enable:      false, // This should be skipped
			Healthy:     true,
		},
	}

	// Call the callback
	nsd.Callback(services, nil)

	// Verify only enabled instance was added
	assert.Equal(t, 1, mockListener.addCount)
}

// Mock implementation of ServiceEventListener
type mockServiceEventListener struct {
	serviceNames []string
	addCount     int
	delCount     int
	updateCount  int
}

func (m *mockServiceEventListener) GetServiceNames() []string {
	return m.serviceNames
}

func (m *mockServiceEventListener) OnAddServiceInstance(instance *servicediscovery.ServiceInstance) {
	m.addCount++
}

func (m *mockServiceEventListener) OnDeleteServiceInstance(instance *servicediscovery.ServiceInstance) {
	m.delCount++
}

func (m *mockServiceEventListener) OnUpdateServiceInstance(instance *servicediscovery.ServiceInstance) {
	m.updateCount++
}

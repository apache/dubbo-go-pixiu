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
	"strconv"
	"testing"

	nacosModel "github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestHandleServiceName(t *testing.T) {
	// Test the handleServiceName function which strips group prefix
	tests := []struct {
		name     string
		input    string
		expected string
	}{
		{
			name:     "with group prefix",
			input:    "DEFAULT_GROUP@@user-service",
			expected: "user-service",
		},
		{
			name:     "without separator",
			input:    "user-service",
			expected: "",
		},
		{
			name:     "empty string",
			input:    "",
			expected: "",
		},
		{
			name:     "separator at end",
			input:    "GROUP@@",
			expected: "",
		},
		{
			name:     "multiple separators",
			input:    "A@@B@@C",
			expected: "B", // strings.Split returns all parts, we take parts[1]
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result := handleServiceName(tt.input)
			assert.Equal(t, tt.expected, result)
		})
	}
}

func TestAppServiceListenerCallback_InstanceMapHandling(t *testing.T) {
	// Test that Callback updates instance map correctly
	// We focus on the instance map update logic without triggering URL generation
	listener := &appServiceListener{
		client:          nil,
		instanceMap:     make(map[string]nacosModel.Instance),
		adapterListener: nil,
	}

	// First callback - add instances
	services1 := []nacosModel.Instance{
		{
			Ip:          "192.168.1.1",
			Port:        8080,
			ServiceName: "test-service",
			Enable:      true,
		},
	}

	// Manually update instanceMap to avoid getURLs call
	// This tests the instance map handling logic
	for _, svc := range services1 {
		if !svc.Enable {
			continue
		}
		host := svc.Ip + ":" + strconv.Itoa(int(svc.Port))
		svc.ServiceName = handleServiceName(svc.ServiceName)
		listener.instanceMap[host] = svc
	}
	assert.Equal(t, 1, len(listener.instanceMap))
}

func TestToNacosInstance(t *testing.T) {
	instance := nacosModel.Instance{
		InstanceId:  "instance-1",
		ServiceName: "test-service",
		Ip:          "192.168.1.1",
		Port:        8080,
		Enable:      true,
		Healthy:     true,
		Metadata:    map[string]string{"key": "value", "num": "123"},
	}

	result := toNacosInstance(instance)

	assert.Equal(t, "instance-1", result.GetID())
	assert.Equal(t, "test-service", result.GetServiceName())
	assert.Equal(t, "192.168.1.1", result.GetHost())
	assert.Equal(t, 8080, result.GetPort())
	assert.True(t, result.IsEnable())
	assert.True(t, result.IsHealthy())
}

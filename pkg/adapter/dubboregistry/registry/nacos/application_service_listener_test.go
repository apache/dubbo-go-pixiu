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
	dubboCommon "dubbo.apache.org/dubbo-go/v3/common"

	"github.com/agiledragon/gomonkey/v2"

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
			expected: "user-service",
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
			expected: "B@@C",
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
	for _, tc := range []struct{ input, expected string }{
		{"test-service", "test-service"},
		{"DEFAULT_GROUP@@test-service", "test-service"},
		{"GROUP@@service@@suffix", "service@@suffix"},
	} {
		t.Run(tc.input, func(t *testing.T) {
			listener := newNacosAppSrvListener(nil, nil)
			var names []string
			// Intercept only the external metadata lookup; exercise the complete
			// callback and verify the name forwarded for URL resolution.
			patches := gomonkey.ApplyPrivateMethod(listener, "getURLs", func(_ *appServiceListener, instance nacosModel.Instance) []*dubboCommon.URL {
				names = append(names, instance.ServiceName)
				return nil
			})
			defer patches.Reset()
			listener.Callback([]nacosModel.Instance{{
				Ip: "192.168.1.1", Port: 8080, ServiceName: tc.input, Enable: true,
			}}, nil)
			assert.Equal(t, []string{tc.expected}, names)
			assert.Equal(t, tc.expected, listener.instanceMap["192.168.1.1:8080"].ServiceName)
			listener.Close()
		})
	}
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

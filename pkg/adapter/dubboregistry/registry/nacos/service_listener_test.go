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

	nacosModel "github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/stretchr/testify/assert"
)

func TestServiceListenerCallback_SkipsDisabledInstances(t *testing.T) {
	// Test that Callback skips instances where Enable is false (line 86-89)
	// This tests the direct assignment on line 91: instance := services[i]

	mockListener := &mockRegistryEventListener{}

	listener := &serviceListener{
		client:          nil,
		instanceMap:     make(map[string]nacosModel.Instance),
		adapterListener: mockListener,
	}

	services := []nacosModel.Instance{
		{
			Ip:          "192.168.1.1",
			Port:        8080,
			ServiceName: "test-service",
			Enable:      true,
			Metadata:    map[string]string{"protocol": "dubbo", "path": "/test"},
		},
		{
			Ip:          "192.168.1.2",
			Port:        8080,
			ServiceName: "test-service",
			Enable:      false, // Should be skipped
		},
	}

	// Call callback - tests the Enable check and direct instance assignment
	listener.Callback(services, nil)

	// Verify only enabled instance was processed (the one with metadata)
	// Since the enabled instance has valid metadata, it should be processed
	assert.Equal(t, 1, mockListener.addCount)     // One instance added
	assert.Equal(t, 1, len(listener.instanceMap)) // One instance in map
}

func TestServiceListenerCallback_InstanceMapHandling(t *testing.T) {
	// Test the instance map update logic (lines 93-102)
	mockListener := &mockRegistryEventListener{}

	listener := &serviceListener{
		client:          nil,
		instanceMap:     make(map[string]nacosModel.Instance),
		adapterListener: mockListener,
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

	listener.Callback(services1, nil)
	assert.Equal(t, 1, len(listener.instanceMap))

	// Second callback - update with same instances
	listener.Callback(services1, nil)
	assert.Equal(t, 1, len(listener.instanceMap))

	// Third callback - remove instances (empty list)
	listener.Callback([]nacosModel.Instance{}, nil)
	assert.Equal(t, 0, len(listener.instanceMap))
}

func TestGenerateURL_WithMetadata(t *testing.T) {
	// Test generateURL function
	instance := nacosModel.Instance{
		Ip:   "192.168.1.1",
		Port: 8080,
		Metadata: map[string]string{
			"protocol": "dubbo",
			"path":     "/test/api",
			"version":  "1.0.0",
		},
	}

	url := generateURL(instance)
	assert.NotNil(t, url)
	assert.Equal(t, "192.168.1.1", url.Ip)
	assert.Equal(t, "8080", url.Port)
	assert.Equal(t, "dubbo", url.Protocol)
	assert.Equal(t, "/test/api", url.Path)
}

func TestGenerateURL_EmptyMetadata(t *testing.T) {
	// Test generateURL with empty metadata
	instance := nacosModel.Instance{
		Ip:       "192.168.1.1",
		Port:     8080,
		Metadata: nil,
	}

	url := generateURL(instance)
	assert.Nil(t, url)
}

func TestGenerateURL_MissingProtocol(t *testing.T) {
	// Test generateURL with missing protocol
	instance := nacosModel.Instance{
		Ip:   "192.168.1.1",
		Port: 8080,
		Metadata: map[string]string{
			"path": "/test/api",
		},
	}

	url := generateURL(instance)
	assert.Nil(t, url)
}

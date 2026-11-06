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
	"sync"
	"testing"
)

import (
	"github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"

	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud/servicediscovery"
	pixiumodel "github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestCallback_ServiceNameWithGroupPrefix(t *testing.T) {
	// Drive the production Callback (not a copy of the separator logic) and
	// assert the @@ group prefix is stripped from ServiceName before the
	// instance reaches the listener. The mock records the ServiceInstance it
	// receives so this fails if the prefix handling is removed, broken, or
	// never wired through to the listener.

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
			mockListener := &mockServiceEventListener{}
			nsd := &nacosServiceDiscovery{
				listener:    mockListener,
				instanceMap: make(map[string]servicediscovery.ServiceInstance),
			}

			services := []model.Instance{
				{
					Ip:          "192.168.1.1",
					Port:        8080,
					ServiceName: tt.serviceName,
					Enable:      true,
					Healthy:     true,
				},
			}

			nsd.Callback(services, nil)

			names := mockListener.addedServiceNames()
			if assert.Len(t, names, 1, "Callback should forward exactly one instance") {
				assert.Equal(t, tt.expected, names[0],
					"ServiceName should have the @@ group prefix stripped by Callback")
			}
		})
	}
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
	serviceNames   []string
	addCount       int
	delCount       int
	updateCount    int
	mu             sync.Mutex
	addedInstances []servicediscovery.ServiceInstance
}

func (m *mockServiceEventListener) GetServiceNames() []string {
	return m.serviceNames
}

func (m *mockServiceEventListener) OnAddServiceInstance(instance *servicediscovery.ServiceInstance) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.addCount++
	m.addedInstances = append(m.addedInstances, *instance)
}

func (m *mockServiceEventListener) OnDeleteServiceInstance(instance *servicediscovery.ServiceInstance) {
	m.delCount++
}

func (m *mockServiceEventListener) OnUpdateServiceInstance(instance *servicediscovery.ServiceInstance) {
	m.updateCount++
}

// addedServiceNames returns the ServiceName of each instance reported via
// OnAddServiceInstance, in arrival order.
func (m *mockServiceEventListener) addedServiceNames() []string {
	m.mu.Lock()
	defer m.mu.Unlock()
	names := make([]string, 0, len(m.addedInstances))
	for _, inst := range m.addedInstances {
		names = append(names, inst.ServiceName)
	}
	return names
}

// mockNamingClient is a stand-in for *nacos.NacosClient satisfying the
// nacosNamingClient interface, so tests can assert on lifecycle without a real
// gRPC connection.
type mockNamingClient struct {
	mu             sync.Mutex
	subscribed     map[string]struct{}
	unsubscribed   map[string]struct{}
	closeClientCnt int
}

func newMockNamingClient() *mockNamingClient {
	return &mockNamingClient{
		subscribed:   make(map[string]struct{}),
		unsubscribed: make(map[string]struct{}),
	}
}

func (m *mockNamingClient) GetAllServicesInfo(param vo.GetAllServiceInfoParam) (model.ServiceList, error) {
	return model.ServiceList{}, nil
}

func (m *mockNamingClient) SelectInstances(param vo.SelectInstancesParam) ([]model.Instance, error) {
	return nil, nil
}

func (m *mockNamingClient) Subscribe(param *vo.SubscribeParam) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subscribed[param.ServiceName] = struct{}{}
	return nil
}

func (m *mockNamingClient) Unsubscribe(param *vo.SubscribeParam) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.unsubscribed[param.ServiceName] = struct{}{}
	return nil
}

func (m *mockNamingClient) Close() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closeClientCnt++
}

func (m *mockNamingClient) closeCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.closeClientCnt
}

// TestUnsubscribeClosesClient locks in the v2 close path: Unsubscribe must
// close the naming client in addition to unsubscribing, otherwise the gRPC
// connection and internal retry goroutines outlive shutdown. See AlexStocks'
// [P1] review on PR #982.
func TestUnsubscribeClosesClient(t *testing.T) {
	client := newMockNamingClient()
	listener := &mockServiceEventListener{serviceNames: []string{"service-A"}}
	nsd := &nacosServiceDiscovery{
		client:      client,
		config:      &pixiumodel.RemoteConfig{Group: "DEFAULT_GROUP"},
		listener:    listener,
		instanceMap: make(map[string]servicediscovery.ServiceInstance),
	}

	assert.Equal(t, 0, client.closeCalls(), "client must not be closed before unsubscribe")

	err := nsd.Unsubscribe()
	assert.NoError(t, err)

	assert.Equal(t, 1, client.closeCalls(),
		"Unsubscribe must close the naming client to release the v2 gRPC connection")
	assert.Contains(t, client.unsubscribed, "service-A", "Unsubscribe must unsubscribe configured services first")
}

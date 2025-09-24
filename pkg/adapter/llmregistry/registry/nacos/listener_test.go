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
	"errors"
	"sync"
	"testing"
	"time"
)

import (
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	nacosModel "github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/llmregistry/common"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type mockNacosClient struct {
	naming_client.INamingClient

	mu                   sync.Mutex
	servicesToReturn     nacosModel.ServiceList
	servicesToReturnErr  error
	subscribeCallback    func(services []nacosModel.SubscribeService, err error)
	subscribedServices   map[string]struct{}
	unsubscribedServices map[string]struct{}
}

func newMockNacosClient() *mockNacosClient {
	return &mockNacosClient{
		subscribedServices:   make(map[string]struct{}),
		unsubscribedServices: make(map[string]struct{}),
	}
}

func (m *mockNacosClient) GetAllServicesInfo(param vo.GetAllServiceInfoParam) (nacosModel.ServiceList, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.servicesToReturn, m.servicesToReturnErr
}

func (m *mockNacosClient) Subscribe(param *vo.SubscribeParam) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subscribedServices[param.ServiceName] = struct{}{}
	m.subscribeCallback = param.SubscribeCallback
	return nil
}

func (m *mockNacosClient) Unsubscribe(param *vo.SubscribeParam) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.unsubscribedServices[param.ServiceName] = struct{}{}
	return nil
}

type mockAdapterListener struct {
	mu               sync.Mutex
	addedEndpoints   map[string]*model.Endpoint
	removedEndpoints map[string]*model.Endpoint
}

var _ common.RegistryEventListener = (*mockAdapterListener)(nil)

func newMockAdapterListener() *mockAdapterListener {
	return &mockAdapterListener{
		addedEndpoints:   make(map[string]*model.Endpoint),
		removedEndpoints: make(map[string]*model.Endpoint),
	}
}

func (m *mockAdapterListener) OnAddEndpoint(endpoint *model.Endpoint) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.addedEndpoints[endpoint.ID] = endpoint
	return nil
}

func (m *mockAdapterListener) OnRemoveEndpoint(endpoint *model.Endpoint) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.removedEndpoints[endpoint.ID] = endpoint
	return nil
}

func (m *mockAdapterListener) reset() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.addedEndpoints = make(map[string]*model.Endpoint)
	m.removedEndpoints = make(map[string]*model.Endpoint)
}

func testSetup() (*listener, *mockNacosClient, *mockAdapterListener) {
	client := newMockNacosClient()
	adapterListener := newMockAdapterListener()
	regConf := &model.Registry{Group: "test_group", Namespace: "test_namespace"}
	nacosListener := newNacosListener(client, regConf, adapterListener)
	return nacosListener, client, adapterListener
}

func TestGenerateEndpoint(t *testing.T) {
	t.Run("Full valid metadata", func(t *testing.T) {
		instance := nacosModel.Instance{
			Metadata: map[string]string{
				"id":                         "ep-123",
				"name":                       "my-llm",
				"ip":                         "127.0.0.1",
				"port":                       "8080",
				"address":                    "openai.com,openai1.com",
				"llm-meta.retry_policy.name": "ExponentialBackoff",
				"llm-meta.fallback":          "true",
			},
		}

		endpoint := generateEndpoint(instance)
		assert.NotNil(t, endpoint)
		assert.Equal(t, "ep-123", endpoint.ID)
		assert.Equal(t, "my-llm", endpoint.Name)
		assert.Equal(t, "127.0.0.1", endpoint.Address.Address)
		assert.Equal(t, 8080, endpoint.Address.Port)
		assert.Equal(t, 2, len(endpoint.Address.Domains))
		assert.Equal(t, "openai.com", endpoint.Address.Domains[0])
		assert.Equal(t, "openai1.com", endpoint.Address.Domains[1])
		assert.Equal(t, model.RetryerExponentialBackoff, endpoint.LLMMeta.RetryPolicy.Name)
		assert.True(t, endpoint.LLMMeta.Fallback)
	})

	t.Run("Nil metadata", func(t *testing.T) {
		instance := nacosModel.Instance{Metadata: nil}
		endpoint := generateEndpoint(instance)
		assert.Nil(t, endpoint)
	})

	t.Run("Invalid port", func(t *testing.T) {
		instance := nacosModel.Instance{Metadata: map[string]string{"port": "not-a-number"}}
		endpoint := generateEndpoint(instance)
		assert.NotNil(t, endpoint)
		assert.Equal(t, 0, endpoint.Address.Port)
	})
}

func TestDiscoverAndSubscribe(t *testing.T) {
	l, client, _ := testSetup()

	t.Run("Discover and subscribe to a new service", func(t *testing.T) {
		// CHANGE THIS LINE:
		client.servicesToReturn = nacosModel.ServiceList{Doms: []string{"service-A"}} // Was nacosModel.Service
		l.discoverAndSubscribe()

		assert.Equal(t, struct{}{}, client.subscribedServices["service-A"], "Should subscribe to service-A")
		_, loaded := l.subscribedServices.Load("service-A")
		assert.True(t, loaded, "service-A should be in the subscribedServices map")
	})

	t.Run("Unsubscribe from a removed service", func(t *testing.T) {
		// Ensure service-A is already subscribed for the test setup
		l.subscribedServices.Store("service-A", true)
		client.servicesToReturn = nacosModel.ServiceList{Doms: []string{}} // Nacos now returns an empty list

		l.discoverAndSubscribe()

		assert.Equal(t, struct{}{}, client.unsubscribedServices["service-A"], "Should unsubscribe from service-A")
		_, loaded := l.subscribedServices.Load("service-A")
		assert.False(t, loaded, "service-A should be removed from subscribedServices map")
	})

	t.Run("Handle Nacos API error", func(t *testing.T) {
		client.subscribedServices = make(map[string]struct{})
		l.subscribedServices.Store("stale-service", true)

		client.servicesToReturnErr = errors.New("Nacos unavailable")
		l.discoverAndSubscribe()

		_, loaded := l.subscribedServices.Load("stale-service")
		assert.True(t, loaded, "Should not change subscriptions on API error")
		assert.Empty(t, client.subscribedServices, "Should not attempt to subscribe on API error")
	})
}

func TestServiceCallback(t *testing.T) {
	l, client, adapterListener := testSetup()

	_ = client.Subscribe(&vo.SubscribeParam{
		ServiceName:       "service-A",
		SubscribeCallback: l.serviceCallback,
	})

	instance1 := nacosModel.SubscribeService{
		InstanceId: "ep-1", ServiceName: "service-A", Enable: true, Healthy: true,
		Metadata: map[string]string{"id": "ep-1", "name": "inst-1"},
	}
	instance2 := nacosModel.SubscribeService{
		InstanceId: "ep-2", ServiceName: "service-A", Enable: true, Healthy: true,
		Metadata: map[string]string{"id": "ep-2", "name": "inst-2"},
	}

	t.Run("Initial instance registration", func(t *testing.T) {
		adapterListener.reset()
		client.subscribeCallback([]nacosModel.SubscribeService{instance1, instance2}, nil)

		assert.Len(t, adapterListener.addedEndpoints, 2, "Should add 2 endpoints")
		assert.Contains(t, adapterListener.addedEndpoints, "ep-1")
		assert.Contains(t, adapterListener.addedEndpoints, "ep-2")
		assert.Empty(t, adapterListener.removedEndpoints, "Should not remove any endpoints")
	})

	t.Run("One instance is removed", func(t *testing.T) {
		adapterListener.reset()
		client.subscribeCallback([]nacosModel.SubscribeService{instance1}, nil)

		assert.Empty(t, adapterListener.addedEndpoints, "Should not add any new endpoints")
		assert.Len(t, adapterListener.removedEndpoints, 1, "Should remove 1 endpoint")
		assert.Contains(t, adapterListener.removedEndpoints, "ep-2")
	})

	t.Run("One instance is updated", func(t *testing.T) {
		adapterListener.reset()
		updatedInstance1 := instance1
		updatedInstance1.Metadata = map[string]string{"id": "ep-1", "name": "inst-1-updated"}

		client.subscribeCallback([]nacosModel.SubscribeService{updatedInstance1}, nil)

		assert.Len(t, adapterListener.addedEndpoints, 1, "Should fire an add/update event for 1 endpoint")
		assert.Contains(t, adapterListener.addedEndpoints, "ep-1")
		assert.Equal(t, "inst-1-updated", adapterListener.addedEndpoints["ep-1"].Name)
		assert.Empty(t, adapterListener.removedEndpoints, "Should not remove any endpoints")
	})

	t.Run("Filter unhealthy or disabled instances", func(t *testing.T) {
		adapterListener.reset()
		unhealthyInstance := instance1
		unhealthyInstance.Healthy = false
		disabledInstance := instance2
		disabledInstance.Enable = false

		client.subscribeCallback([]nacosModel.SubscribeService{unhealthyInstance, disabledInstance}, nil)

		assert.Empty(t, adapterListener.addedEndpoints, "Should not add unhealthy/disabled endpoints")
		assert.Len(t, adapterListener.removedEndpoints, 1, "Should remove the previously active endpoints")
	})

	t.Run("No changes in instances", func(t *testing.T) {
		client.subscribeCallback([]nacosModel.SubscribeService{instance1}, nil)

		adapterListener.reset()

		client.subscribeCallback([]nacosModel.SubscribeService{instance1}, nil)

		assert.Empty(t, adapterListener.addedEndpoints, "Should not trigger add for unchanged instance")
		assert.Empty(t, adapterListener.removedEndpoints, "Should not trigger remove for unchanged instance")
	})
}

func TestLifecycle(t *testing.T) {
	l, _, _ := testSetup()

	l.WatchAndHandle()

	time.Sleep(100 * time.Millisecond)

	// test that Close works without panic
	assert.NotPanics(t, func() {
		l.Close()
	})
}

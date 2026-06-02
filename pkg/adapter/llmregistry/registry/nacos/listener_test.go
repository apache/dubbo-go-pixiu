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
	"strings"
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
		instance := nacosModel.Instance{
			Port:     8080,
			Metadata: map[string]string{"port": "not-a-number"},
		}
		endpoint := generateEndpoint(instance)
		assert.NotNil(t, endpoint)
		assert.Equal(t, 8080, endpoint.Address.Port)
	})

	t.Run("Missing metadata ID uses Nacos instance identity", func(t *testing.T) {
		instance := nacosModel.Instance{
			InstanceId: "nacos-instance-1",
			Ip:         "10.0.0.1",
			Port:       18080,
			Metadata: map[string]string{
				"cluster": "llm-cluster",
				"name":    "shared-llm",
			},
		}

		endpoint := generateEndpoint(instance)
		assert.NotNil(t, endpoint)
		assert.Equal(t, "nacos-instance-1", endpoint.ID)
		assert.Equal(t, "10.0.0.1", endpoint.Address.Address)
		assert.Equal(t, 18080, endpoint.Address.Port)
	})

	t.Run("Missing ID uses stable generated endpoint ID", func(t *testing.T) {
		instance := nacosModel.Instance{
			Metadata: map[string]string{
				"cluster":          "llm-cluster",
				"name":             "shared-llm",
				"ip":               "127.0.0.1",
				"port":             "8080",
				"llm-meta.api_key": "key-a",
			},
		}

		first := generateEndpoint(instance)
		second := generateEndpoint(instance)
		changedCredential := instance
		changedCredential.Metadata = map[string]string{
			"cluster":          "llm-cluster",
			"name":             "shared-llm",
			"ip":               "127.0.0.1",
			"port":             "8080",
			"llm-meta.api_key": "key-b",
		}

		assert.NotNil(t, first)
		assert.Equal(t, first.ID, second.ID)
		assert.NotEqual(t, first.ID, generateEndpoint(changedCredential).ID)
		assert.Contains(t, first.ID, "pixiu-generated-endpoint-")
		assert.NotContains(t, first.ID, "key-a")
	})

	t.Run("Missing metadata cluster falls back to empty cluster hash", func(t *testing.T) {
		build := func(clusterName string) nacosModel.Instance {
			return nacosModel.Instance{
				Ip:          "127.0.0.1",
				Port:        8080,
				ClusterName: clusterName,
				Metadata: map[string]string{
					"name":             "shared-llm",
					"llm-meta.api_key": "key-a",
				},
			}
		}

		clusterA := generateEndpoint(build("cluster-a"))
		clusterB := generateEndpoint(build("cluster-b"))

		assert.NotNil(t, clusterA)
		assert.NotNil(t, clusterB)
		assert.Contains(t, clusterA.ID, "pixiu-generated-endpoint-")
		assert.Equal(t, clusterA.ID, clusterB.ID)
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

func TestServiceCallbackUsesStableGeneratedEndpointID(t *testing.T) {
	l, client, adapterListener := testSetup()

	_ = client.Subscribe(&vo.SubscribeParam{
		ServiceName:       "service-generated",
		SubscribeCallback: l.serviceCallback,
	})

	instance1 := nacosModel.SubscribeService{
		ServiceName: "service-generated",
		Enable:      true,
		Healthy:     true,
		Metadata: map[string]string{
			"cluster":          "llm-cluster",
			"name":             "shared-llm-a",
			"ip":               "127.0.0.1",
			"port":             "8080",
			"llm-meta.api_key": "key-a",
		},
	}
	instance2 := nacosModel.SubscribeService{
		ServiceName: "service-generated",
		Enable:      true,
		Healthy:     true,
		Metadata: map[string]string{
			"cluster":          "llm-cluster",
			"name":             "shared-llm-b",
			"ip":               "127.0.0.1",
			"port":             "8080",
			"llm-meta.api_key": "key-b",
		},
	}

	client.subscribeCallback([]nacosModel.SubscribeService{instance1, instance2}, nil)
	assert.Len(t, adapterListener.addedEndpoints, 2)

	var removedID string
	for id, endpoint := range adapterListener.addedEndpoints {
		if endpoint.Name == "shared-llm-a" {
			removedID = id
		}
	}
	assert.NotEmpty(t, removedID)

	adapterListener.reset()
	client.subscribeCallback([]nacosModel.SubscribeService{instance2}, nil)

	assert.Empty(t, adapterListener.addedEndpoints)
	assert.Contains(t, adapterListener.removedEndpoints, removedID)
}

func TestServiceCallbackKeepsGeneratedEndpointIDWhenNameChanges(t *testing.T) {
	l, client, adapterListener := testSetup()

	_ = client.Subscribe(&vo.SubscribeParam{
		ServiceName:       "service-rename",
		SubscribeCallback: l.serviceCallback,
	})

	instance := nacosModel.SubscribeService{
		ServiceName: "service-rename",
		Ip:          "127.0.0.1",
		Port:        8080,
		Enable:      true,
		Healthy:     true,
		Metadata: map[string]string{
			"cluster":          "llm-cluster",
			"name":             "shared-llm-old",
			"llm-meta.api_key": "key-a",
		},
	}

	client.subscribeCallback([]nacosModel.SubscribeService{instance}, nil)
	assert.Len(t, adapterListener.addedEndpoints, 1)

	var endpointID string
	for id := range adapterListener.addedEndpoints {
		endpointID = id
	}
	assert.NotEmpty(t, endpointID)

	renamed := instance
	renamed.Metadata = map[string]string{
		"cluster":          "llm-cluster",
		"name":             "shared-llm-new",
		"llm-meta.api_key": "key-a",
	}

	adapterListener.reset()
	client.subscribeCallback([]nacosModel.SubscribeService{renamed}, nil)

	assert.Empty(t, adapterListener.removedEndpoints)
	if assert.Contains(t, adapterListener.addedEndpoints, endpointID) {
		assert.Equal(t, "shared-llm-new", adapterListener.addedEndpoints[endpointID].Name)
	}
}

func TestServiceCallbackKeepsNacosInstancesWithoutMetadataIDDistinct(t *testing.T) {
	l, client, adapterListener := testSetup()

	_ = client.Subscribe(&vo.SubscribeParam{
		ServiceName:       "service-instance-id",
		SubscribeCallback: l.serviceCallback,
	})

	instance1 := nacosModel.SubscribeService{
		InstanceId:  "nacos-instance-a",
		ServiceName: "service-instance-id",
		Ip:          "10.0.0.1",
		Port:        18080,
		Enable:      true,
		Healthy:     true,
		Metadata: map[string]string{
			"cluster":          "llm-cluster",
			"name":             "shared-llm",
			"llm-meta.api_key": "same-key",
		},
	}
	instance2 := nacosModel.SubscribeService{
		InstanceId:  "nacos-instance-b",
		ServiceName: "service-instance-id",
		Ip:          "10.0.0.2",
		Port:        18081,
		Enable:      true,
		Healthy:     true,
		Metadata: map[string]string{
			"cluster":          "llm-cluster",
			"name":             "shared-llm",
			"llm-meta.api_key": "same-key",
		},
	}

	client.subscribeCallback([]nacosModel.SubscribeService{instance1, instance2}, nil)

	assert.Len(t, adapterListener.addedEndpoints, 2)
	assert.Contains(t, adapterListener.addedEndpoints, "nacos-instance-a")
	assert.Contains(t, adapterListener.addedEndpoints, "nacos-instance-b")

	adapterListener.reset()
	client.subscribeCallback([]nacosModel.SubscribeService{instance2}, nil)

	assert.Empty(t, adapterListener.addedEndpoints)
	assert.Contains(t, adapterListener.removedEndpoints, "nacos-instance-a")
}

// TestNacosEndpointIDMissingClusterFallsBackToEmptyClusterHash ensures that
// when a nacos instance has no metadata["id"], no InstanceId, and no
// metadata["cluster"], nacosEndpointID returns a generated- ID derived
// with an empty cluster name. The downstream LLM registry adapter
// (Adapter.OnAddEndpoint) drops such endpoints, so we do not attempt to
// disambiguate cross-service collisions here. The contract this test
// locks: the function does not invent a synthesized fallback prefix and
// does not pretend an unreachable endpoint will be admitted.
func TestNacosEndpointIDMissingClusterFallsBackToEmptyClusterHash(t *testing.T) {
	instance := nacosModel.Instance{
		// InstanceId intentionally empty; metadata has neither "id" nor "cluster".
		Ip:          "10.0.0.1",
		Port:        18080,
		ServiceName: "alpha",
		ClusterName: "DEFAULT",
		Metadata:    map[string]string{"llm-meta.api_key": "key-shared"},
	}

	endpoint := generateEndpoint(instance)
	assert.True(t, strings.HasPrefix(endpoint.ID, "pixiu-generated-endpoint-"),
		"missing metadata[\"cluster\"] falls through to model.GenerateEndpointID with empty cluster")

	// Deterministic: re-generating from the same instance returns the same ID.
	assert.Equal(t, endpoint.ID, generateEndpoint(instance).ID)

	// Acknowledged limitation: at this code level we cannot tell two
	// instances apart that share address+credential and both lack
	// metadata["cluster"]. They alias to the same generated- ID. The LLM
	// registry adapter skips them before the alias has any runtime effect.
	collidingInstance := nacosModel.Instance{
		Ip:          instance.Ip,
		Port:        instance.Port,
		ServiceName: "bravo", // different service, but the adapter ignores ServiceName for ID
		ClusterName: instance.ClusterName,
		Metadata:    instance.Metadata,
	}
	assert.Equal(t, endpoint.ID, generateEndpoint(collidingInstance).ID,
		"two instances missing metadata[\"cluster\"] at the same address will alias here; "+
			"disambiguation is the adapter's job (currently: drop)")
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

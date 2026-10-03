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
	"time"
)

import (
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	nacosModel "github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	registryapi "github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry"
	registrybase "github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry/base"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type lifecycleNamingClient struct {
	naming_client.INamingClient

	mu           sync.Mutex
	services     nacosModel.ServiceList
	subscribed   []*vo.SubscribeParam
	unsubscribed []*vo.SubscribeParam
}

func (m *lifecycleNamingClient) GetAllServicesInfo(vo.GetAllServiceInfoParam) (nacosModel.ServiceList, error) {
	return m.services, nil
}

func (m *lifecycleNamingClient) Subscribe(param *vo.SubscribeParam) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.subscribed = append(m.subscribed, param)
	return nil
}

func (m *lifecycleNamingClient) Unsubscribe(param *vo.SubscribeParam) error {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.unsubscribed = append(m.unsubscribed, param)
	return nil
}

func (m *lifecycleNamingClient) subscriptionCounts() (int, int) {
	m.mu.Lock()
	defer m.mu.Unlock()
	return len(m.subscribed), len(m.unsubscribed)
}

func TestNacosRegistry_DoUnsubscribeRequiresInitialization(t *testing.T) {
	n := &NacosRegistry{}

	err := n.DoUnsubscribe()

	require.Error(t, err)
	assert.Contains(t, err.Error(), "base registry is not initialized")
}

func TestNacosRegistry_ApplicationLifecycle(t *testing.T) {
	client := &lifecycleNamingClient{
		services: nacosModel.ServiceList{Doms: []string{"orders-app"}},
	}
	regConfig := &model.Registry{Group: "test-group", Namespace: "test-namespace"}
	listener := newNacosAppListener(client, nil, regConfig, nil)
	n := newLifecycleRegistry(registryapi.RegisteredTypeApplication, listener)

	require.NoError(t, n.DoSubscribe())
	require.Eventually(t, func() bool {
		subscribed, _ := client.subscriptionCounts()
		return subscribed == 1
	}, time.Second, 10*time.Millisecond)

	var wg sync.WaitGroup
	for range 8 {
		wg.Add(1)
		go func() {
			defer wg.Done()
			assert.NoError(t, n.DoUnsubscribe())
		}()
	}
	wg.Wait()

	subscribed, unsubscribed := client.subscriptionCounts()
	assert.Equal(t, 1, subscribed)
	assert.Equal(t, 1, unsubscribed)
	client.mu.Lock()
	defer client.mu.Unlock()
	assert.Same(t, client.subscribed[0], client.unsubscribed[0])
}

func TestNacosRegistry_InterfaceLifecycle(t *testing.T) {
	client := &lifecycleNamingClient{
		services: nacosModel.ServiceList{Doms: []string{"providers:org.example.OrderService:1.0.0:test-group"}},
	}
	regConfig := &model.Registry{Group: "test-group", Namespace: "test-namespace"}
	listener := newNacosIntfListener(client, nil, regConfig, nil)
	n := newLifecycleRegistry(registryapi.RegisteredTypeInterface, listener)

	require.NoError(t, n.DoSubscribe())
	require.Eventually(t, func() bool {
		subscribed, _ := client.subscriptionCounts()
		return subscribed == 1
	}, time.Second, 10*time.Millisecond)
	require.NoError(t, n.DoUnsubscribe())

	client.mu.Lock()
	defer client.mu.Unlock()
	require.Len(t, client.subscribed, 1)
	require.Len(t, client.unsubscribed, 1)
	assert.Equal(t, client.subscribed[0].ServiceName, client.unsubscribed[0].ServiceName)
	assert.Same(t, client.subscribed[0], client.unsubscribed[0])
}

func newLifecycleRegistry(registeredType registryapi.RegisteredType, listener registryapi.Listener) *NacosRegistry {
	return &NacosRegistry{
		BaseRegistry: &registrybase.BaseRegistry{RegisteredType: registeredType},
		nacosListeners: map[registryapi.RegisteredType]registryapi.Listener{
			registeredType: listener,
		},
	}
}

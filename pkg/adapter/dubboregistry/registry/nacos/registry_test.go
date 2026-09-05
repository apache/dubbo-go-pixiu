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
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"

	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry"
	baseRegistry "github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry/base"
)

// mockNamingClient embeds the SDK interface and only overrides CloseClient so
// the shutdown path can be asserted without a real gRPC connection.
type mockNamingClient struct {
	naming_client.INamingClient

	mu             sync.Mutex
	closeClientCnt int
}

func (m *mockNamingClient) CloseClient() {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.closeClientCnt++
}

func (m *mockNamingClient) closeCalls() int {
	m.mu.Lock()
	defer m.mu.Unlock()
	return m.closeClientCnt
}

// closeTrackingListener is a minimal registry.Listener that records Close.
type closeTrackingListener struct {
	closed bool
	mu     sync.Mutex
}

func (l *closeTrackingListener) Close() {
	l.mu.Lock()
	defer l.mu.Unlock()
	l.closed = true
}

func (l *closeTrackingListener) WatchAndHandle() {}

// newTestNacosRegistry builds a NacosRegistry wired with a mock naming client
// and a listener registered under the given registered type, without going
// through the factory (which would create a real gRPC client).
func newTestNacosRegistry(client naming_client.INamingClient, rt registry.RegisteredType, listener registry.Listener) *NacosRegistry {
	reg := &NacosRegistry{
		client:         client,
		nacosListeners: map[registry.RegisteredType]registry.Listener{rt: listener},
	}
	reg.BaseRegistry = baseRegistry.NewBaseRegistry(reg, &mockRegistryEventListener{}, rt)
	return reg
}

// TestDoUnsubscribeClosesClient locks in the v2 close path: DoUnsubscribe must
// close the naming client (CloseClient) in addition to stopping the listener,
// otherwise the gRPC connection and internal retry goroutines outlive shutdown.
// See AlexStocks' [P1] review on PR #982.
func TestDoUnsubscribeClosesClient(t *testing.T) {
	client := &mockNamingClient{}
	listener := &closeTrackingListener{}
	reg := newTestNacosRegistry(client, registry.RegisteredTypeInterface, listener)

	// Register a service-level listener so the close loop over GetAllSvcListener
	// is also exercised (it must Close and remove each one).
	svcListener := &closeTrackingListener{}
	reg.SetSvcListener("svc-1", svcListener)

	assert.Equal(t, 0, client.closeCalls(), "client must not be closed before unsubscribe")
	assert.False(t, listener.closed, "listener must not be closed before unsubscribe")

	err := reg.DoUnsubscribe()
	assert.NoError(t, err)

	assert.True(t, listener.closed, "DoUnsubscribe must close the listener first")
	assert.True(t, svcListener.closed, "DoUnsubscribe must close service-level listeners")
	assert.Nil(t, reg.GetSvcListener("svc-1"), "DoUnsubscribe must remove closed service-level listeners")
	assert.Equal(t, 1, client.closeCalls(),
		"DoUnsubscribe must close the naming client to release the v2 gRPC connection")
}

// TestDoUnsubscribe_NoListener ensures DoUnsubscribe returns an error when no
// listener is registered for the configured registered type (instead of the
// previous panic("implement me") behavior).
func TestDoUnsubscribe_NoListener(t *testing.T) {
	client := &mockNamingClient{}
	reg := newTestNacosRegistry(client, registry.RegisteredTypeInterface, nil)
	// Overwrite the listener map so the configured type has no listener.
	reg.nacosListeners = map[registry.RegisteredType]registry.Listener{}

	err := reg.DoUnsubscribe()
	assert.Error(t, err)
	assert.Equal(t, 0, client.closeCalls(), "client must not be closed when there is no listener")
}

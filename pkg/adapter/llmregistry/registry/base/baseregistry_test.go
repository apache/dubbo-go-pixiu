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

package baseregistry

import (
	"errors"
	"sync"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/llmregistry/registry"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// mockFacadeRegistry is a mock implementation of the FacadeRegistry interface for testing.
type mockFacadeRegistry struct {
	subscribeCalled   bool
	unsubscribeCalled bool
	subscribeErr      error
	unsubscribeErr    error
}

func (m *mockFacadeRegistry) DoSubscribe() error {
	m.subscribeCalled = true
	return m.subscribeErr
}

func (m *mockFacadeRegistry) DoUnsubscribe() error {
	m.unsubscribeCalled = true
	return m.unsubscribeErr
}

// mockListener is a mock implementation of the registry.Listener interface.
// Since BaseRegistry only stores and retrieves it, we don't need a complex implementation.
type mockListener struct{}

func (m *mockListener) WatchAndHandle() {
	panic("implement me") // NOSONAR
}

// Close is a mock method.
func (m *mockListener) Close() {} // NOSONAR

// mockAdapterListener is a mock implementation of the common.RegistryEventListener interface.
type mockAdapterListener struct{}

func (m *mockAdapterListener) OnAddEndpoint(r *model.Endpoint) error {
	panic("implement me") // NOSONAR
}

func (m *mockAdapterListener) OnRemoveEndpoint(r *model.Endpoint) error {
	panic("implement me") // NOSONAR
}

func TestSvcListeners(t *testing.T) {
	// Initialization
	svcListeners := &SvcListeners{
		listeners: make(map[string]registry.Listener),
	}
	mockL := &mockListener{}
	id := "test-service-1"

	// 1. Test SetListener and GetListener
	t.Run("SetAndGet", func(t *testing.T) {
		// Try to get a non-existent listener
		listener := svcListeners.GetListener(id)
		assert.Nil(t, listener, "Getting a non-existent listener should return nil")

		// Set and then get an existing listener
		svcListeners.SetListener(id, mockL)
		listener = svcListeners.GetListener(id)
		assert.NotNil(t, listener, "Getting an existing listener should not return nil")
		assert.Equal(t, mockL, listener, "The retrieved listener should be the same one that was set")
	})

	// 2. Test GetAllListener
	t.Run("GetAll", func(t *testing.T) {
		id2 := "test-service-2"
		mockL2 := &mockListener{}
		svcListeners.SetListener(id2, mockL2)

		allListeners := svcListeners.GetAllListener()
		assert.Len(t, allListeners, 2, "GetAllListener should return all set listeners")
		assert.Contains(t, allListeners, id, "The returned map should contain the first listener")
		assert.Contains(t, allListeners, id2, "The returned map should contain the second listener")
	})

	// 3. Test RemoveListener
	t.Run("Remove", func(t *testing.T) {
		svcListeners.RemoveListener(id)
		listener := svcListeners.GetListener(id)
		assert.Nil(t, listener, "Getting a removed listener should return nil")

		allListeners := svcListeners.GetAllListener()
		assert.Len(t, allListeners, 1, "The total count of listeners should decrease after removal")
		assert.NotContains(t, allListeners, id, "The returned map should no longer contain the removed listener")
	})
}

// Test the concurrency safety of SvcListeners.
func TestSvcListeners_concurrency(t *testing.T) {
	svcListeners := &SvcListeners{
		listeners: make(map[string]registry.Listener),
	}
	id := "concurrent-service"
	mockL := &mockListener{}

	var wg sync.WaitGroup
	// Start multiple goroutines to read and write concurrently
	for i := 0; i < 100; i++ {
		wg.Add(1)
		go func() {
			defer wg.Done()
			svcListeners.SetListener(id, mockL)
			_ = svcListeners.GetListener(id)
			svcListeners.RemoveListener(id)
		}()
	}

	wg.Wait()
	// The main purpose of this test is to detect race conditions.
	// Running tests with the `go test -race` command can expose such issues.
	// If the test completes without errors, it indicates the lock mechanism is working correctly.
}

// --- Tests for BaseRegistry ---

// Test the constructor NewBaseRegistry.
func TestNewBaseRegistry(t *testing.T) {
	mockFacade := &mockFacadeRegistry{}
	mockAdapter := &mockAdapterListener{}

	br := NewBaseRegistry(mockFacade, mockAdapter)

	assert.NotNil(t, br, "NewBaseRegistry should not return nil")
	assert.Equal(t, mockFacade, br.facadeRegistry, "facadeRegistry should be initialized correctly")
	assert.Equal(t, mockAdapter, br.AdapterListener, "AdapterListener should be initialized correctly")
	assert.NotNil(t, br.svcListeners, "svcListeners should be initialized")
	assert.Empty(t, br.svcListeners.listeners, "The svcListeners map should be empty initially")
}

// Test the listener management methods of BaseRegistry.
func TestBaseRegistry_listenerMethods(t *testing.T) {
	br := NewBaseRegistry(&mockFacadeRegistry{}, &mockAdapterListener{})
	mockL := &mockListener{}
	id := "test-service"

	// Test SetSvcListener and GetSvcListener
	l := br.GetSvcListener(id)
	assert.Nil(t, l)

	br.SetSvcListener(id, mockL)
	l = br.GetSvcListener(id)
	assert.NotNil(t, l)
	assert.Equal(t, mockL, l)

	// Test GetAllSvcListener
	allListeners := br.GetAllSvcListener()
	assert.Len(t, allListeners, 1)
	assert.Equal(t, mockL, allListeners[id])

	// Test RemoveSvcListener
	br.RemoveSvcListener(id)
	l = br.GetSvcListener(id)
	assert.Nil(t, l)
	allListeners = br.GetAllSvcListener()
	assert.Empty(t, allListeners)
}

// Test the Subscribe and Unsubscribe methods.
func TestBaseRegistry_subscribeUnsubscribe(t *testing.T) {
	t.Run("Success", func(t *testing.T) {
		mockFacade := &mockFacadeRegistry{}
		br := NewBaseRegistry(mockFacade, &mockAdapterListener{})

		// Test Subscribe
		err := br.Subscribe()
		assert.NoError(t, err)
		assert.True(t, mockFacade.subscribeCalled, "facadeRegistry.DoSubscribe should have been called")

		// Test Unsubscribe
		err = br.Unsubscribe()
		assert.NoError(t, err)
		assert.True(t, mockFacade.unsubscribeCalled, "facadeRegistry.DoUnsubscribe should have been called")
	})

	t.Run("Error", func(t *testing.T) {
		subscribeError := errors.New("failed to subscribe")
		unsubscribeError := errors.New("failed to unsubscribe")

		mockFacade := &mockFacadeRegistry{
			subscribeErr:   subscribeError,
			unsubscribeErr: unsubscribeError,
		}
		br := NewBaseRegistry(mockFacade, &mockAdapterListener{})

		// Test Subscribe error propagation
		err := br.Subscribe()
		assert.Error(t, err)
		assert.Equal(t, subscribeError, err, "Subscribe should propagate the error from facadeRegistry")

		// Test Unsubscribe error propagation
		err = br.Unsubscribe()
		assert.Error(t, err)
		assert.Equal(t, unsubscribeError, err, "Unsubscribe should propagate the error from facadeRegistry")
	})
}

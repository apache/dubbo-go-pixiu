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

package server

import (
	"errors"
	"sync"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type transactionListenerService struct {
	config   model.Listener
	failName string
}

func (s *transactionListenerService) Start() error       { return nil }
func (s *transactionListenerService) Close() error       { return nil }
func (s *transactionListenerService) ShutDown(any) error { return nil }
func (s *transactionListenerService) Refresh(config model.Listener) error {
	if config.Name == s.failName {
		return errors.New("refresh rejected")
	}
	s.config = config
	return nil
}

func TestListenerManager_XDSOwnershipPreservesStaticListener(t *testing.T) {
	staticListener := &model.Listener{
		Name:        "static",
		ProtocolStr: "HTTP",
		Address: model.Address{SocketAddress: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    18080,
		}},
	}
	key := resolveListenerName(staticListener)
	lm := &ListenerManager{
		activeListenerService: map[string]*wrapListenerService{
			key: {config: staticListener},
		},
		xdsManaged: make(map[string]struct{}),
		rwLock:     &sync.RWMutex{},
	}

	require.Error(t, lm.UpsertXDSListener(staticListener))
	lm.RemoveXDSListeners([]string{key})

	assert.True(t, lm.HasListener(key))
	assert.Empty(t, lm.XDSListenerNames())
}

func TestListenerManager_ReplaceXDSListenersKeepsLastGoodOnCreateFailure(t *testing.T) {
	oldListener := &model.Listener{
		Name:        "old-dynamic",
		ProtocolStr: "HTTP",
		Protocol:    model.ProtocolTypeHTTP,
		Address: model.Address{SocketAddress: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    18080,
		}},
	}
	oldKey := resolveListenerName(oldListener)
	lm := &ListenerManager{
		activeListenerService: map[string]*wrapListenerService{
			oldKey: {config: oldListener},
		},
		xdsManaged: map[string]struct{}{oldKey: {}},
		rwLock:     &sync.RWMutex{},
	}

	err := lm.ReplaceXDSListeners([]*model.Listener{{
		Name:        "invalid-new",
		ProtocolStr: "UNSUPPORTED",
		Protocol:    model.ProtocolType(999),
		Address: model.Address{SocketAddress: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    18081,
		}},
	}})

	require.ErrorContains(t, err, "does not support yet")
	assert.True(t, lm.HasListener(oldKey))
	assert.Equal(t, []string{oldKey}, lm.XDSListenerNames())
}

func TestListenerManager_ReplaceXDSListenersRollsBackEarlierRefresh(t *testing.T) {
	oldOne := &model.Listener{Name: "old-one", ProtocolStr: "HTTP", Protocol: model.ProtocolTypeHTTP}
	oldOne.Address.SocketAddress = model.SocketAddress{Address: "127.0.0.1", Port: 18080}
	oldTwo := &model.Listener{Name: "old-two", ProtocolStr: "HTTP", Protocol: model.ProtocolTypeHTTP}
	oldTwo.Address.SocketAddress = model.SocketAddress{Address: "127.0.0.1", Port: 18081}
	keyOne := resolveListenerName(oldOne)
	keyTwo := resolveListenerName(oldTwo)
	serviceOne := &transactionListenerService{config: *oldOne}
	serviceTwo := &transactionListenerService{config: *oldTwo, failName: "new-two"}
	lm := &ListenerManager{
		activeListenerService: map[string]*wrapListenerService{
			keyOne: {ListenerService: serviceOne, config: oldOne},
			keyTwo: {ListenerService: serviceTwo, config: oldTwo},
		},
		xdsManaged: map[string]struct{}{keyOne: {}, keyTwo: {}},
		rwLock:     &sync.RWMutex{},
	}
	newOne := *oldOne
	newOne.Name = "new-one"
	newTwo := *oldTwo
	newTwo.Name = "new-two"

	err := lm.ReplaceXDSListeners([]*model.Listener{&newOne, &newTwo})

	require.ErrorContains(t, err, "refresh rejected")
	assert.Equal(t, "old-one", serviceOne.config.Name)
	assert.Equal(t, "old-two", serviceTwo.config.Name)
	assert.Same(t, oldOne, lm.activeListenerService[keyOne].config)
	assert.Same(t, oldTwo, lm.activeListenerService[keyTwo].config)
}

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

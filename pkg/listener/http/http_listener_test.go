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

package http

import (
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/filterchain"
	listenerpkg "github.com/apache/dubbo-go-pixiu/pkg/listener"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestCloseFilterChainAfterShutdownDoesNotWaitForActiveRequest(t *testing.T) {
	listener := &HttpListenerService{
		BaseListenerService: listenerpkg.BaseListenerService{
			FilterChain: &filterchain.NetworkFilterChain{},
		},
	}
	listener.filterState = newHTTPFilterChainState(listener.FilterChain)
	state := listener.filterState
	require.True(t, state.acquire())

	finished := make(chan struct{})
	go func() {
		_ = listener.closeFilterChainAfterShutdown()
		close(finished)
	}()

	select {
	case <-finished:
	case <-time.After(100 * time.Millisecond):
		t.Fatal("shutdown filter cleanup waited for an active request")
	}

	state.release()
	select {
	case <-state.done:
	case <-time.After(time.Second):
		t.Fatal("deferred filter cleanup did not run after the request released its chain lease")
	}
}

func TestRefreshSwapsChainWithoutWaitingForActiveRequest(t *testing.T) {
	listener := &HttpListenerService{
		BaseListenerService: listenerpkg.BaseListenerService{
			FilterChain: &filterchain.NetworkFilterChain{},
		},
	}
	listener.filterState = newHTTPFilterChainState(listener.FilterChain)
	oldState := listener.filterState
	require.True(t, oldState.acquire())

	refreshed := make(chan error, 1)
	go func() {
		refreshed <- listener.Refresh(model.Listener{})
	}()

	select {
	case err := <-refreshed:
		require.NoError(t, err)
	case <-time.After(100 * time.Millisecond):
		oldState.release()
		t.Fatal("refresh waited for an active request on the old filter chain")
	}
	require.NotSame(t, oldState, listener.filterState)
	oldState.release()
	select {
	case <-oldState.done:
	case <-time.After(time.Second):
		t.Fatal("old filter chain was not closed after the active request released")
	}
}

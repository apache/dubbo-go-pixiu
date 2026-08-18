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
	"github.com/apache/dubbo-go-pixiu/pkg/filterchain"
	listenerpkg "github.com/apache/dubbo-go-pixiu/pkg/listener"
)

func TestCloseFilterChainAfterShutdownDoesNotWaitForActiveRequest(t *testing.T) {
	listener := &HttpListenerService{
		BaseListenerService: listenerpkg.BaseListenerService{
			FilterChain: &filterchain.NetworkFilterChain{},
		},
	}

	listener.filterMu.RLock()
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

	listener.filterMu.RUnlock()
	deadline := time.Now().Add(time.Second)
	for time.Now().Before(deadline) {
		listener.filterMu.RLock()
		removed := listener.FilterChain == nil
		listener.filterMu.RUnlock()
		if removed {
			return
		}
		time.Sleep(time.Millisecond)
	}
	t.Fatal("deferred filter cleanup did not run after the request released its read lock")
}

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
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func TestListenersCloseDuringRetryBackoff(t *testing.T) {
	factories := map[string]func(*mockNamingClient) registry.Listener{
		"interface": func(c *mockNamingClient) registry.Listener {
			return newNacosIntfListener(c, nil, &model.Registry{}, nil)
		},
		"application": func(c *mockNamingClient) registry.Listener {
			return newNacosAppListener(c, nil, &model.Registry{}, nil)
		},
	}
	for name, newListener := range factories {
		t.Run(name, func(t *testing.T) {
			client := &mockNamingClient{pollStarted: make(chan struct{}), pollErr: errors.New("Nacos is unavailable")}
			listener := newListener(client)
			listener.WatchAndHandle()
			select {
			case <-client.pollStarted:
			case <-time.After(time.Second):
				t.Fatal("listener did not start polling")
			}
			done := make(chan struct{})
			go func() { listener.Close(); close(done) }()
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("Close did not interrupt retry backoff")
			}
		})
	}
}

func TestClosedListenersDoNotStartPolling(t *testing.T) {
	client := &mockNamingClient{pollStarted: make(chan struct{})}
	reg := newTestNacosRegistry(client, registry.RegisteredTypeInterface, nil)
	app := newNacosAppListener(client, reg, &model.Registry{}, &mockRegistryEventListener{}).(*nacosAppListener)
	intf := newNacosIntfListener(client, reg, &model.Registry{}, &mockRegistryEventListener{}).(*nacosIntfListener)
	for _, watch := range []struct {
		close func()
		start func()
		wait  *sync.WaitGroup
	}{
		{app.Close, app.WatchAndHandle, &app.wg},
		{intf.Close, intf.WatchAndHandle, &intf.wg},
	} {
		watch.close()
		watch.start()
		done := make(chan struct{})
		go func() { watch.wait.Wait(); close(done) }()
		select {
		case <-done:
		case <-time.After(time.Second):
			t.Fatal("closed listener did not exit")
		}
	}
	select {
	case <-client.pollStarted:
		t.Fatal("listener polled Nacos after it was closed")
	default:
	}
}

func TestNacosPollDelayCanExpireOrBeCanceled(t *testing.T) {
	assert.True(t, waitForNacosPoll(make(chan struct{}), 0), "an elapsed delay must allow the next discovery poll")
	exit := make(chan struct{})
	close(exit)
	assert.False(t, waitForNacosPoll(exit, time.Hour), "closing must interrupt a long retry delay")
}

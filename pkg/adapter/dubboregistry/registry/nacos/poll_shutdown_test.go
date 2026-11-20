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
	"sync/atomic"
	"testing"
	"time"
)

import (
	"github.com/nacos-group/nacos-sdk-go/v2/clients/naming_client"
	nacosModel "github.com/nacos-group/nacos-sdk-go/v2/model"
	"github.com/nacos-group/nacos-sdk-go/v2/vo"
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

// A temporary discovery failure must resume polling after backoff, and the
// listener must remain stoppable after recovery. Use lifecycle entry points so
// the test also applies when the internal polling implementation changes.
type recoveringNamingClient struct {
	*mockNamingClient
	calls     atomic.Int32
	recovered chan struct{}
}

func (c *recoveringNamingClient) GetAllServicesInfo(vo.GetAllServiceInfoParam) (nacosModel.ServiceList, error) {
	switch c.calls.Add(1) {
	case 1:
		return nacosModel.ServiceList{}, errors.New("temporary Nacos outage")
	case 2:
		close(c.recovered)
	}
	return nacosModel.ServiceList{}, nil
}

func TestListenersResumePollingAfterTemporaryFailure(t *testing.T) {
	factories := map[string]func(naming_client.INamingClient) registry.Listener{
		"interface": func(c naming_client.INamingClient) registry.Listener {
			return newNacosIntfListener(c, nil, &model.Registry{}, nil)
		},
		"application": func(c naming_client.INamingClient) registry.Listener {
			return newNacosAppListener(c, nil, &model.Registry{}, nil)
		},
	}
	for name, newListener := range factories {
		t.Run(name, func(t *testing.T) {
			client := &recoveringNamingClient{mockNamingClient: &mockNamingClient{}, recovered: make(chan struct{})}
			listener := newListener(client)
			listener.WatchAndHandle()
			select {
			case <-client.recovered:
			case <-time.After(10 * time.Second):
				t.Fatal("listener did not resume discovery after a temporary failure")
			}
			done := make(chan struct{})
			go func() { listener.Close(); close(done) }()
			select {
			case <-done:
			case <-time.After(time.Second):
				t.Fatal("listener did not stop after recovery")
			}
		})
	}
}

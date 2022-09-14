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

package listener

import (
	"sync/atomic"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/filterchain"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

var factoryMap = make(map[model.ProtocolType]func(lc *model.Listener, bs *model.Bootstrap) (ListenerService, error), 8)

type (
	ListenerService interface {
		// Start the listener service
		Start() error
		// Close the listener service forcefully
		Close() error
		// ShutDown gracefully shuts down the listener.
		ShutDown(interface{}) error
		// Refresh config
		Refresh(model.Listener) error
	}

	BaseListenerService struct {
		Config      *model.Listener
		FilterChain *filterchain.NetworkFilterChain
	}

	ListenerGracefulShutdownConfig struct {
		ActiveCount   int32
		RejectRequest bool
	}
)

// SetListenerServiceFactory will store the listenerService factory by name
func SetListenerServiceFactory(protocol model.ProtocolType, newRegFunc func(lc *model.Listener, bs *model.Bootstrap) (ListenerService, error)) {
	factoryMap[protocol] = newRegFunc
}

// CreateListenerService create listener service
func CreateListenerService(lc *model.Listener, bs *model.Bootstrap) (ListenerService, error) {
	if registry, ok := factoryMap[lc.Protocol]; ok {
		reg, err := registry(lc, bs)
		if err != nil {
			panic("Initialize ListenerService " + lc.Name + "failed due to: " + err.Error())
		}
		return reg, nil
	}
	return nil, errors.New("Registry " + lc.ProtocolStr + " does not support yet")
}

func (lgsc *ListenerGracefulShutdownConfig) AddActiveCount(num int32) {
	atomic.AddInt32(&lgsc.ActiveCount, num)
}

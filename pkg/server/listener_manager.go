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
	"runtime/debug"
	"strconv"
	"sync"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/listener"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

// wrapListenerService wrap listener service and its configuration.
type wrapListenerService struct {
	cfg *model.Listener
	listener.ListenerService
}

// ListenerManager the listener manager
type ListenerManager struct {
	bootstrap *model.Bootstrap

	// name(host-port-protocol) -> wrapListenerService
	activeListenerService map[string]*wrapListenerService
	//readWriteLock
	rwLock *sync.RWMutex
}

// CreateDefaultListenerManager create listener manager from config
func CreateDefaultListenerManager(bs *model.Bootstrap) *ListenerManager {
	listeners := map[string]*wrapListenerService{}
	sl := bs.GetStaticListeners()

	for _, lsCof := range sl {
		ls, err := listener.CreateListenerService(lsCof, bs)
		if err != nil {
			logger.Error("CreateDefaultListenerManager %s error: %v", lsCof.Name, err)
		}
		listeners[resolveListenerName(lsCof)] = &wrapListenerService{
			cfg:             lsCof,
			ListenerService: ls,
		}
	}

	return &ListenerManager{
		activeListenerService: listeners,
		bootstrap:             bs,
		rwLock:                &sync.RWMutex{},
	}
}

func resolveListenerName(c *model.Listener) string {
	return c.Address.SocketAddress.Address + "-" + strconv.Itoa(c.Address.SocketAddress.Port) + "-" + c.ProtocolStr
}

func (lm *ListenerManager) AddListener(lsConf *model.Listener) error {
	ls, err := listener.CreateListenerService(lsConf, lm.bootstrap)
	if err != nil {
		return err
	}
	lm.startListenerServiceAsync(ls)
	lm.addListenerService(ls, lsConf)
	return nil
}

func (lm *ListenerManager) UpdateListener(m *model.Listener) error {
	//TODO implement me
	logger.Infof("UpdateListener %s", m.Name)
	service := lm.GetListenerService(m.Name)
	if service != nil {
		return service.Refresh(*m)
	}
	return nil
}

func (lm *ListenerManager) HasListener(name string) bool {
	lm.rwLock.RLock()
	defer lm.rwLock.RUnlock()
	_, ok := lm.activeListenerService[name]
	return ok
}

func (lm *ListenerManager) CloneXdsControlListener() ([]*model.Listener, error) {
	lm.rwLock.RLock()
	defer lm.rwLock.RUnlock()
	var listeners []*model.Listener
	for _, ls := range lm.activeListenerService {
		listeners = append(listeners, ls.cfg)
	}
	return listeners, nil
}

func (lm *ListenerManager) getListener(name string) *model.Listener {
	service := lm.activeListenerService[name]
	if service != nil {
		return service.cfg
	}
	return nil
}

func (lm *ListenerManager) StartListen() {
	for _, s := range lm.activeListenerService {
		lm.startListenerServiceAsync(s)
	}
}

func (lm *ListenerManager) startListenerServiceAsync(s listener.ListenerService) chan<- struct{} {
	done := make(chan struct{})
	go func() {
		defer func() {
			panicErr := recover()
			if panicErr != nil {
				logger.Errorf("recover from panic %v", panicErr)
				debug.PrintStack()
			}
			close(done)
		}()
		err := s.Start()
		if err != nil {
			logger.Error("start listener service error.  %v", err)
		}
	}()
	return done
}

func (lm *ListenerManager) addListenerService(ls listener.ListenerService, lsConf *model.Listener) {
	lm.rwLock.Lock()
	defer lm.rwLock.Unlock()
	lm.activeListenerService[resolveListenerName(lsConf)] = &wrapListenerService{
		cfg:             lsConf,
		ListenerService: ls,
	}
}

func (lm *ListenerManager) GetListenerService(name string) listener.ListenerService {
	lm.rwLock.RLock()
	defer lm.rwLock.RUnlock()
	return lm.activeListenerService[name]
}

func (lm *ListenerManager) RemoveListener(names []string) {
	//close ListenerService
	for _, name := range names {
		logger.Infof("listener %s closing", name)
		ls := lm.GetListenerService(name)
		if ls == nil {
			logger.Warnf("listener %s not found", name)
			continue
		}
		if err := ls.Close(); err != nil {
			logger.Error("close listener service error.  %name", err)
			continue
		}
		logger.Infof("listener %s closed", name)
	}
	lm.rwLock.Lock()
	defer lm.rwLock.Unlock()
	//remove from activeListenerService
	for _, name := range names {
		delete(lm.activeListenerService, name)
	}
}

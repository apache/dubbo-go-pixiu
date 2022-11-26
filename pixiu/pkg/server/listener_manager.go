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
	"os"
	"os/signal"
	"runtime/debug"
	"strconv"
	"sync"
	"time"
)

import (
	"github.com/pkg/errors"
	"gopkg.in/yaml.v2"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/shutdown"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/listener"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

// wrapListenerService wrap listener service and its configuration.
type wrapListenerService struct {
	listener.ListenerService

	config *model.Listener
}

// ListenerManager the listener manager
type ListenerManager struct {
	bootstrap *model.Bootstrap

	// name(host-port-protocol) -> wrapListenerService
	activeListenerService map[string]*wrapListenerService
	//readWriteLock
	rwLock *sync.RWMutex
	//shutdownWaitGroup
	shutdownWG *sync.WaitGroup
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
			config:          lsCof,
			ListenerService: ls,
		}
	}

	lm := &ListenerManager{
		activeListenerService: listeners,
		bootstrap:             bs,
		rwLock:                &sync.RWMutex{},
		shutdownWG:            &sync.WaitGroup{},
	}
	lm.gracefulShutdownInit()

	return lm
}

func (lm *ListenerManager) gracefulShutdownInit() {
	sdc := lm.bootstrap.GetShutdownConfig()
	timeout := sdc.GetTimeout()
	if timeout <= 0 {
		return
	}
	signals := make(chan os.Signal, 1)
	signal.Notify(signals, shutdown.ShutdownSignals...)

	go func() {
		sig := <-signals
		logger.Infof("get signal %s, dubbo-go-pixiu will start shutdown.", sig)

		time.AfterFunc(timeout, func() {
			logger.Warn("Shutdown gracefully timeout, listeners will shutdown immediately. ")
			os.Exit(0)
		})

		for _, listener := range lm.activeListenerService {
			lm.shutdownWG.Add(1)
			go func(listener *wrapListenerService) {
				err := listener.ListenerService.ShutDown(lm.shutdownWG)
				if err != nil {
					logger.Errorf("Shutdown Error: %+v", err)
					os.Exit(0)
				}
			}(listener)
		}
		lm.shutdownWG.Wait()

		// those signals' original behavior is exit with dump ths stack, so we try to keep the behavior
		for _, dumpSignal := range shutdown.DumpHeapShutdownSignals {
			if sig == dumpSignal {
				debug.WriteHeapDump(os.Stdout.Fd())
			}
		}
		os.Exit(0)
	}()
}

func resolveListenerName(c *model.Listener) string {
	return c.Address.SocketAddress.Address + "-" + strconv.Itoa(c.Address.SocketAddress.Port) + "-" + c.ProtocolStr
}

func (lm *ListenerManager) AddListener(lsConf *model.Listener) error {
	logger.Infof("Add Listener %s", lsConf.Name)
	ls, err := listener.CreateListenerService(lsConf, lm.bootstrap)
	if err != nil {
		return err
	}
	lm.addListenerService(ls, lsConf)
	lm.startListenerServiceAsync(ls)
	return nil
}

func (lm *ListenerManager) UpdateListener(m *model.Listener) error {
	// lock
	lm.rwLock.Lock()
	defer lm.rwLock.Unlock()
	ls, ok := lm.activeListenerService[m.Name]
	if !ok {
		return errors.New("ListenerManager UpdateListener error: listener not found")
	}
	logger.Infof("Update Listener %s", m.Name)
	ls.config = m
	err := ls.Refresh(*m)
	if err != nil {
		logger.Warnf("Update Listener %s error: %s", m.Name, err)
		return err
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
		listeners = append(listeners, ls.config)
	}
	//deep copy
	bytes, err := yaml.Marshal(listeners)
	if err != nil {
		return nil, err
	}
	var cloneListeners []*model.Listener
	if err = yaml.Unmarshal(bytes, &cloneListeners); err != nil {
		return nil, err
	}
	return cloneListeners, nil
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
			logger.Errorf("start listener service error.  %v", err)
		}
	}()
	return done
}

func (lm *ListenerManager) addListenerService(ls listener.ListenerService, lsConf *model.Listener) {
	lm.rwLock.Lock()
	defer lm.rwLock.Unlock()
	lm.activeListenerService[resolveListenerName(lsConf)] = &wrapListenerService{
		config:          lsConf,
		ListenerService: ls,
	}
}

func (lm *ListenerManager) GetListenerService(name string) listener.ListenerService {
	lm.rwLock.RLock()
	defer lm.rwLock.RUnlock()

	ls, ok := lm.activeListenerService[name]
	if ok {
		return ls
	}
	return nil
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
			logger.Errorf("close listener %s service error.  %s", name, err)
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

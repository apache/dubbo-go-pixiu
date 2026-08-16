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
	"slices"
	"strconv"
	"sync"
	"time"
)

import (
	"github.com/pkg/errors"

	"gopkg.in/yaml.v3"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/shutdown"
	"github.com/apache/dubbo-go-pixiu/pkg/listener"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
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
	// xdsManaged contains only listener keys created through the xDS API.
	xdsManaged map[string]struct{}
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
			logger.Errorf("CreateDefaultListenerManager %s error: %v", lsCof.Name, err)
		}
		listeners[resolveListenerName(lsCof)] = &wrapListenerService{
			config:          lsCof,
			ListenerService: ls,
		}
	}

	lm := &ListenerManager{
		activeListenerService: listeners,
		xdsManaged:            make(map[string]struct{}),
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
				err := listener.ShutDown(lm.shutdownWG)
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
	if m == nil {
		return errors.New("UpdateListener error: provided listener config is nil")
	}
	// lock
	lm.rwLock.Lock()
	defer lm.rwLock.Unlock()

	// Use resolveListenerName to get the correct key
	listenerKey := resolveListenerName(m)
	ls, ok := lm.activeListenerService[listenerKey]
	if !ok {
		return errors.Errorf("ListenerManager UpdateListener error: listener not found with key %s", listenerKey)
	}
	logger.Infof("Update Listener %s (key: %s)", m.Name, listenerKey)
	ls.config = m
	err := ls.Refresh(*m)
	if err != nil {
		logger.Warnf("Update Listener %s error: %s", m.Name, err)
		return err
	}
	return nil
}

// UpsertXDSListener adds or updates a listener owned by xDS. An xDS resource
// is not allowed to replace a listener created from static configuration.
func (lm *ListenerManager) UpsertXDSListener(m *model.Listener) error {
	if m == nil {
		return errors.New("xDS listener config is nil")
	}

	listenerKey := resolveListenerName(m)
	lm.rwLock.Lock()
	if active, exists := lm.activeListenerService[listenerKey]; exists {
		if _, owned := lm.xdsManaged[listenerKey]; !owned {
			lm.rwLock.Unlock()
			return errors.Errorf("xDS listener %q conflicts with a non-xDS listener", listenerKey)
		}
		logger.Infof("Update xDS Listener %s (key: %s)", m.Name, listenerKey)
		if err := active.Refresh(*m); err != nil {
			lm.rwLock.Unlock()
			return errors.Wrapf(err, "refresh xDS listener %q", listenerKey)
		}
		active.config = m
		lm.rwLock.Unlock()
		return nil
	}

	ls, err := listener.CreateListenerService(m, lm.bootstrap)
	if err != nil {
		lm.rwLock.Unlock()
		return err
	}
	lm.activeListenerService[listenerKey] = &wrapListenerService{
		config:          m,
		ListenerService: ls,
	}
	lm.xdsManaged[listenerKey] = struct{}{}
	lm.rwLock.Unlock()

	logger.Infof("Add xDS Listener %s (key: %s)", m.Name, listenerKey)
	lm.startListenerServiceAsync(ls)
	return nil
}

// XDSListenerNames returns the listener keys currently owned by xDS.
func (lm *ListenerManager) XDSListenerNames() []string {
	lm.rwLock.RLock()
	defer lm.rwLock.RUnlock()

	names := make([]string, 0, len(lm.xdsManaged))
	for name := range lm.xdsManaged {
		names = append(names, name)
	}
	slices.Sort(names)
	return names
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
				debug.PrintStack() // NOSONAR
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

// RemoveXDSListeners removes only listeners owned by xDS. Names belonging to
// static or other runtime configuration sources are ignored.
func (lm *ListenerManager) RemoveXDSListeners(names []string) {
	lm.rwLock.Lock()
	listeners := make(map[string]*wrapListenerService, len(names))
	for _, name := range names {
		if _, owned := lm.xdsManaged[name]; !owned {
			continue
		}
		if active := lm.activeListenerService[name]; active != nil {
			listeners[name] = active
		}
		delete(lm.activeListenerService, name)
		delete(lm.xdsManaged, name)
	}
	lm.rwLock.Unlock()

	for name, active := range listeners {
		logger.Infof("xDS listener %s closing", name)
		if err := active.Close(); err != nil {
			logger.Errorf("close xDS listener %s service error: %s", name, err)
		}
	}
}

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
	"sync/atomic"
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

// ShutdownFunc is the signature for a listener shutdown function.
// It returns an error if shutdown failed.
type ShutdownFunc func() error

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
			logger.Errorf("CreateDefaultListenerManager %s error: %v", lsCof.Name, err)
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

		// Handle shutdown signal (extracted for testability)
		shutdownErrors, timedOut := lm.handleShutdownSignal(sig, timeout)

		// Exit with appropriate code based on shutdown errors
		if len(shutdownErrors) > 0 {
			logger.Errorf("Shutdown completed with %d errors", len(shutdownErrors))
			os.Exit(1)
		}
		if timedOut {
			logger.Warn("Shutdown gracefully timeout, some listeners may not have shut down cleanly")
		}
		os.Exit(0)
	}()
}

// handleShutdownSignal processes a shutdown signal by coordinating listener shutdowns.
// It returns a slice of errors from failed shutdowns and a boolean indicating timeout.
// This function is extracted from gracefulShutdownInit for testability.
func (lm *ListenerManager) handleShutdownSignal(sig os.Signal, timeout time.Duration) ([]error, bool) {
	// Build shutdown functions for all listeners
	shutdownFuncs := make([]ShutdownFunc, 0, len(lm.activeListenerService))
	for _, listener := range lm.activeListenerService {
		shutdownFuncs = append(shutdownFuncs, func() error {
			lm.shutdownWG.Add(1)
			// Note: listener.ShutDown() internally calls wg.Done()
			return listener.ShutDown(lm.shutdownWG)
		})
	}

	// Execute shutdown coordination
	shutdownErrors, timedOut := shutdownListeners(shutdownFuncs, timeout)

	// those signals' original behavior is exit with dump ths stack, so we try to keep the behavior
	for _, dumpSignal := range shutdown.DumpHeapShutdownSignals {
		if sig == dumpSignal {
			debug.WriteHeapDump(os.Stdout.Fd())
		}
	}

	return shutdownErrors, timedOut
}

// shutdownListeners coordinates the shutdown of multiple listeners.
// It returns a slice of errors from failed shutdowns and a boolean indicating
// whether the shutdown timed out before all listeners completed.
//
// The overall timeout is the sole bound on how long we wait for listeners: it
// must honor the user-configured graceful-shutdown budget. We deliberately do
// NOT impose a shorter per-listener deadline. A fixed per-listener timeout
// would let a single slow listener be considered "done" early and cause the
// caller to os.Exit(0) and interrupt the graceful shutdown of listeners that
// are still within their allowed time (P0 review feedback on PR #993).
func shutdownListeners(shutdownFuncs []ShutdownFunc, timeout time.Duration) ([]error, bool) {
	if len(shutdownFuncs) == 0 {
		return nil, false
	}

	// Error channel is buffered so a slow send never blocks the worker goroutine.
	errCh := make(chan error, len(shutdownFuncs))
	// doneCh is signalled by each worker AFTER its error (if any) has been sent,
	// so receiving all doneCh signals guarantees every error is already queued.
	doneCh := make(chan struct{}, len(shutdownFuncs))

	var completed int32
	// Start shutdown for all listeners
	for _, shutdownFunc := range shutdownFuncs {
		go func(fn ShutdownFunc) {
			err := fn()
			if err != nil {
				logger.Errorf("Shutdown Error: %+v", err)
				errCh <- err
			}
			atomic.AddInt32(&completed, 1)
			// Signal that this goroutine has completed (after the error send).
			doneCh <- struct{}{}
		}(shutdownFunc)
	}

	// Wait for every listener to finish, bounded only by the overall timeout.
	// There is no per-listener deadline: a listener that is slow but still
	// within the configured budget must be allowed to run to completion.
	allDone := make(chan struct{})
	go func() {
		for i := 0; i < len(shutdownFuncs); i++ {
			<-doneCh
		}
		close(allDone)
	}()

	var timedOut bool
	select {
	case <-allDone:
		// All listener goroutines completed
		logger.Info("All listeners shut down gracefully")
	case <-time.After(timeout):
		timedOut = true
		unfinished := int32(len(shutdownFuncs)) - atomic.LoadInt32(&completed)
		logger.Warnf("Shutdown timed out after %s; %d of %d listener(s) did not finish gracefully",
			timeout, unfinished, len(shutdownFuncs))
	}

	// Drain any queued errors (non-blocking). After allDone this is complete
	// and safe (all sends are done); after a timeout it is best-effort.
	var shutdownErrors []error
drainErrors:
	for {
		select {
		case err := <-errCh:
			shutdownErrors = append(shutdownErrors, err)
		default:
			break drainErrors
		}
	}

	return shutdownErrors, timedOut
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

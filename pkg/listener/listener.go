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
	"sync"
	"sync/atomic"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/filterchain"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

var factoryMap = make(map[model.ProtocolType]func(lc *model.Listener, bs *model.Bootstrap) (ListenerService, error), 8)

type (
	ListenerService interface {
		// Start the listener service
		Start() error
		// Close the listener service forcefully
		Close() error
		// ShutDown gracefully shuts down the listener.
		ShutDown(any) error
		// Refresh config
		Refresh(model.Listener) error
	}

	BaseListenerService struct {
		Config      *model.Listener
		FilterChain *filterchain.NetworkFilterChain

		filterChain    atomic.Pointer[filterChainSnapshot]
		active         atomic.Bool
		updateGate     *sync.RWMutex
		updateGateOnce sync.Once
	}

	PreparedUpdate struct {
		snapshot *filterChainSnapshot
	}

	RetiredUpdate struct {
		snapshot *filterChainSnapshot
	}

	filterChainSnapshot struct {
		filterChain *filterchain.NetworkFilterChain
		state       atomic.Uint64
		closeOnce   sync.Once
		closeErr    error
	}

	ListenerGracefulShutdownConfig struct {
		// ActiveCount and RejectRequest remain exported for source compatibility.
		// Listener implementations in this repository use the atomic helpers below.
		ActiveCount   int32
		RejectRequest bool

		rejectRequests atomic.Bool
	}
)

const filterChainRetired = uint64(1) << 63

// SetListenerServiceFactory will store the listenerService factory by name
func SetListenerServiceFactory(protocol model.ProtocolType, newRegFunc func(lc *model.Listener, bs *model.Bootstrap) (ListenerService, error)) {
	factoryMap[protocol] = newRegFunc
}

// CreateListenerService creates a listener service and preserves the original
// public factory signature for out-of-tree listener implementations.
func CreateListenerService(lc *model.Listener, bs *model.Bootstrap) (ListenerService, error) {
	return createListenerService(lc, bs, nil)
}

// CreateListenerServiceWithUpdateGate creates a listener whose BaseListenerService
// participates in a manager-wide, short critical section during xDS publication.
func CreateListenerServiceWithUpdateGate(lc *model.Listener, bs *model.Bootstrap, updateGate *sync.RWMutex) (ListenerService, error) {
	return createListenerService(lc, bs, updateGate)
}

func createListenerService(lc *model.Listener, bs *model.Bootstrap, updateGate *sync.RWMutex) (ListenerService, error) {
	if registry, ok := factoryMap[lc.Protocol]; ok {
		reg, err := registry(lc, bs)
		if err != nil {
			return nil, errors.Wrapf(err, "initialize listener service %q", lc.Name)
		}
		if reg == nil {
			return nil, errors.Errorf("listener service factory for %q returned nil", lc.Name)
		}
		if updateGate != nil {
			if setter, ok := reg.(interface{ setUpdateGate(*sync.RWMutex) }); ok {
				setter.setUpdateGate(updateGate)
			}
		}
		return reg, nil
	}
	return nil, errors.New("Registry " + lc.ProtocolStr + " does not support yet")
}

// NewBaseListenerService initializes the race-safe filter-chain holder shared
// by all listener implementations.
func NewBaseListenerService(lc *model.Listener, fc *filterchain.NetworkFilterChain) *BaseListenerService {
	base := &BaseListenerService{Config: lc, FilterChain: fc, updateGate: &sync.RWMutex{}}
	base.filterChain.Store(newFilterChainSnapshot(fc))
	base.active.Store(true)
	return base
}

func (ls *BaseListenerService) setUpdateGate(gate *sync.RWMutex) {
	if gate != nil {
		ls.updateGate = gate
	}
}

func (ls *BaseListenerService) gate() *sync.RWMutex {
	ls.updateGateOnce.Do(func() {
		if ls.updateGate == nil {
			ls.updateGate = &sync.RWMutex{}
		}
	})
	return ls.updateGate
}

// SetActive controls whether requests may enter the filter chain. Listener
// managers keep newly bound sockets inactive until the complete transaction
// has been committed.
func (ls *BaseListenerService) SetActive(active bool) {
	ls.active.Store(active)
}

// WithFilterChain acquires a chain snapshot while briefly holding the shared
// publication gate. Network I/O runs without that gate; the snapshot's
// reference count delays closing a retired chain until the request completes.
func (ls *BaseListenerService) WithFilterChain(fn func(*filterchain.NetworkFilterChain) error) error {
	gate := ls.gate()
	gate.RLock()
	if !ls.active.Load() {
		gate.RUnlock()
		return errors.New("listener is not active")
	}
	snapshot := ls.currentSnapshot()
	if snapshot == nil || !snapshot.acquire() {
		gate.RUnlock()
		return errors.New("listener filter chain is not initialized")
	}
	gate.RUnlock()
	defer snapshot.release()
	return fn(snapshot.filterChain)
}

func (ls *BaseListenerService) PrepareRefresh(config model.Listener) (*PreparedUpdate, error) {
	fc, err := filterchain.BuildNetworkFilterChain(config.FilterChain)
	if err != nil {
		return nil, err
	}
	return &PreparedUpdate{snapshot: newFilterChainSnapshot(fc)}, nil
}

func (ls *BaseListenerService) CommitRefresh(update *PreparedUpdate) *RetiredUpdate {
	if update == nil || update.snapshot == nil {
		return nil
	}
	next := update.snapshot
	update.snapshot = nil
	old := ls.filterChain.Swap(next)
	ls.FilterChain = next.filterChain
	return &RetiredUpdate{snapshot: old}
}

// ClosePreparedUpdate releases filters created for a transaction that was
// rejected before publication.
func ClosePreparedUpdate(update *PreparedUpdate) error {
	if update == nil || update.snapshot == nil {
		return nil
	}
	snapshot := update.snapshot
	update.snapshot = nil
	return snapshot.retire()
}

// Close retires the old chain. Active requests keep their snapshot alive and
// release it after their network operation finishes.
func (update *RetiredUpdate) Close() error {
	if update == nil || update.snapshot == nil {
		return nil
	}
	snapshot := update.snapshot
	update.snapshot = nil
	return snapshot.retire()
}

func (ls *BaseListenerService) Refresh(config model.Listener) error {
	update, err := ls.PrepareRefresh(config)
	if err != nil {
		return err
	}
	gate := ls.gate()
	gate.Lock()
	old := ls.CommitRefresh(update)
	gate.Unlock()
	if old != nil {
		return old.Close()
	}
	return nil
}

func (ls *BaseListenerService) CloseFilterChain() error {
	gate := ls.gate()
	gate.Lock()
	ls.active.Store(false)
	old := ls.filterChain.Swap(nil)
	ls.FilterChain = nil
	gate.Unlock()
	if old != nil {
		return old.retire()
	}
	return nil
}

func newFilterChainSnapshot(fc *filterchain.NetworkFilterChain) *filterChainSnapshot {
	if fc == nil {
		return nil
	}
	return &filterChainSnapshot{filterChain: fc}
}

func (ls *BaseListenerService) currentSnapshot() *filterChainSnapshot {
	if snapshot := ls.filterChain.Load(); snapshot != nil {
		return snapshot
	}
	if ls.FilterChain == nil {
		return nil
	}
	candidate := newFilterChainSnapshot(ls.FilterChain)
	if ls.filterChain.CompareAndSwap(nil, candidate) {
		return candidate
	}
	return ls.filterChain.Load()
}

func (snapshot *filterChainSnapshot) acquire() bool {
	for snapshot != nil {
		state := snapshot.state.Load()
		if state&filterChainRetired != 0 || state == filterChainRetired-1 {
			return false
		}
		if snapshot.state.CompareAndSwap(state, state+1) {
			return true
		}
	}
	return false
}

func (snapshot *filterChainSnapshot) release() {
	if snapshot == nil {
		return
	}
	state := snapshot.state.Add(^uint64(0))
	if state == filterChainRetired {
		if err := snapshot.close(); err != nil {
			logger.Warnf("close retired listener filter chain: %v", err)
		}
	}
}

func (snapshot *filterChainSnapshot) retire() error {
	if snapshot == nil {
		return nil
	}
	for {
		state := snapshot.state.Load()
		if state&filterChainRetired != 0 {
			if state == filterChainRetired {
				return snapshot.close()
			}
			return nil
		}
		if snapshot.state.CompareAndSwap(state, state|filterChainRetired) {
			if state == 0 {
				return snapshot.close()
			}
			return nil
		}
	}
}

func (snapshot *filterChainSnapshot) close() error {
	snapshot.closeOnce.Do(func() {
		snapshot.closeErr = snapshot.filterChain.Close()
	})
	return snapshot.closeErr
}

func (lgsc *ListenerGracefulShutdownConfig) AddActiveCount(num int32) {
	atomic.AddInt32(&lgsc.ActiveCount, num)
}

func (lgsc *ListenerGracefulShutdownConfig) GetActiveCount() int32 {
	return atomic.LoadInt32(&lgsc.ActiveCount)
}

func (lgsc *ListenerGracefulShutdownConfig) SetRejectRequests(reject bool) {
	lgsc.RejectRequest = reject
	lgsc.rejectRequests.Store(reject)
}

func (lgsc *ListenerGracefulShutdownConfig) RejectRequests() bool {
	return lgsc.rejectRequests.Load()
}

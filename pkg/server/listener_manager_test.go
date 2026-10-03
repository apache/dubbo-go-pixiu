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
	"errors"
	"net"
	"os"
	"sync"
	"syscall"
	"testing"
	"time"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/listener"
	_ "github.com/apache/dubbo-go-pixiu/pkg/listener/http"
	_ "github.com/apache/dubbo-go-pixiu/pkg/listener/http2"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

type transactionListenerService struct {
	listener.BaseListenerService
	config   model.Listener
	pending  model.Listener
	failName string
}

type legacyListenerService struct {
	config model.Listener
}

type closeTrackingListenerService struct {
	closed int
}

type handoffTrackingListenerService struct {
	starts int
	stops  int
}

var _ listener.ListenerService = (*legacyListenerService)(nil)

func (s *legacyListenerService) Start() error       { return nil }
func (s *legacyListenerService) Close() error       { return nil }
func (s *legacyListenerService) ShutDown(any) error { return nil }
func (s *legacyListenerService) Refresh(config model.Listener) error {
	s.config = config
	return nil
}

func (s *closeTrackingListenerService) Start() error                 { return nil }
func (s *closeTrackingListenerService) Close() error                 { s.closed++; return nil }
func (s *closeTrackingListenerService) ShutDown(any) error           { return nil }
func (s *closeTrackingListenerService) Refresh(model.Listener) error { return nil }

func (s *handoffTrackingListenerService) Start() error {
	s.starts++
	return nil
}
func (s *handoffTrackingListenerService) StopAccepting() error {
	s.stops++
	return nil
}
func (s *handoffTrackingListenerService) Close() error                 { return nil }
func (s *handoffTrackingListenerService) ShutDown(any) error           { return nil }
func (s *handoffTrackingListenerService) Refresh(model.Listener) error { return nil }

func (s *transactionListenerService) Start() error       { return nil }
func (s *transactionListenerService) Close() error       { return nil }
func (s *transactionListenerService) ShutDown(any) error { return nil }
func (s *transactionListenerService) Refresh(config model.Listener) error {
	if config.Name == s.failName {
		return errors.New("refresh rejected")
	}
	s.config = config
	return nil
}
func (s *transactionListenerService) PrepareRefresh(config model.Listener) (*listener.PreparedUpdate, error) {
	if config.Name == s.failName {
		return nil, errors.New("refresh rejected")
	}
	s.pending = config
	return s.BaseListenerService.PrepareRefresh(config)
}
func (s *transactionListenerService) CommitRefresh(update *listener.PreparedUpdate) *listener.RetiredUpdate {
	s.config = s.pending
	return s.BaseListenerService.CommitRefresh(update)
}

func TestListenerManager_XDSOwnershipPreservesStaticListener(t *testing.T) {
	staticListener := &model.Listener{
		Name:        "static",
		ProtocolStr: "HTTP",
		Address: model.Address{SocketAddress: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    18080,
		}},
	}
	key := resolveListenerName(staticListener)
	lm := &ListenerManager{
		activeListenerService: map[string]*wrapListenerService{
			key: {config: staticListener},
		},
		xdsManaged: make(map[string]struct{}),
		rwLock:     &sync.RWMutex{},
	}

	require.Error(t, lm.UpsertXDSListener(staticListener))
	lm.RemoveXDSListeners([]string{key})

	assert.True(t, lm.HasListener(key))
	assert.Empty(t, lm.XDSListenerNames())
}

func TestCreateDefaultListenerManagerSkipsInvalidStaticListener(t *testing.T) {
	config := &model.Listener{
		Name:        "invalid-static",
		ProtocolStr: "HTTP",
		Protocol:    model.ProtocolTypeHTTP,
		Address: model.Address{SocketAddress: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    18080,
		}},
		FilterChain: model.FilterChain{Filters: []model.NetworkFilter{{
			Name: "missing.network.filter.plugin",
		}}},
	}
	manager := CreateDefaultListenerManager(&model.Bootstrap{StaticResources: model.StaticResources{
		Listeners: []*model.Listener{config},
	}})

	require.False(t, manager.HasListener(resolveListenerName(config)))
}

func TestReplaceXDSListenersPreservesLegacyListenerWhenTransactionalRefreshIsUnavailable(t *testing.T) {
	lastGood := &model.Listener{
		Name:        "last-good",
		ProtocolStr: "HTTP",
		Protocol:    model.ProtocolTypeHTTP,
		Address: model.Address{SocketAddress: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    18080,
		}},
	}
	key := resolveListenerName(lastGood)
	service := &legacyListenerService{config: *lastGood}
	manager := &ListenerManager{
		activeListenerService: map[string]*wrapListenerService{
			key: {config: lastGood, ListenerService: service},
		},
		xdsManaged: map[string]struct{}{key: {}},
		rwLock:     &sync.RWMutex{},
	}

	next := *lastGood
	next.Name = "next"
	err := manager.ReplaceXDSListeners([]*model.Listener{&next})

	require.ErrorContains(t, err, "does not support transactional xDS refresh")
	require.Equal(t, "last-good", service.config.Name)
	require.Same(t, lastGood, manager.activeListenerService[key].config)
}

func TestListenerManager_ReplaceXDSListenersKeepsLastGoodOnCreateFailure(t *testing.T) {
	oldListener := &model.Listener{
		Name:        "old-dynamic",
		ProtocolStr: "HTTP",
		Protocol:    model.ProtocolTypeHTTP,
		Address: model.Address{SocketAddress: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    18080,
		}},
	}
	oldKey := resolveListenerName(oldListener)
	lm := &ListenerManager{
		activeListenerService: map[string]*wrapListenerService{
			oldKey: {config: oldListener},
		},
		xdsManaged: map[string]struct{}{oldKey: {}},
		rwLock:     &sync.RWMutex{},
	}

	err := lm.ReplaceXDSListeners([]*model.Listener{{
		Name:        "invalid-new",
		ProtocolStr: "UNSUPPORTED",
		Protocol:    model.ProtocolType(999),
		Address: model.Address{SocketAddress: model.SocketAddress{
			Address: "127.0.0.1",
			Port:    18081,
		}},
	}})

	require.ErrorContains(t, err, "does not support yet")
	assert.True(t, lm.HasListener(oldKey))
	assert.Equal(t, []string{oldKey}, lm.XDSListenerNames())
}

func TestListenerManager_ReplaceXDSListenersNACKsOccupiedPortBeforeCommit(t *testing.T) {
	occupied, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = occupied.Close() })
	port := occupied.Addr().(*net.TCPAddr).Port

	oldListener := &model.Listener{Name: "last-good", ProtocolStr: "HTTP", Protocol: model.ProtocolTypeHTTP}
	oldListener.Address.SocketAddress = model.SocketAddress{Address: "127.0.0.1", Port: port + 1}
	oldKey := resolveListenerName(oldListener)
	lm := &ListenerManager{
		activeListenerService: map[string]*wrapListenerService{
			oldKey: {ListenerService: &transactionListenerService{config: *oldListener}, config: oldListener},
		},
		xdsManaged: map[string]struct{}{oldKey: {}},
		rwLock:     &sync.RWMutex{},
		updateGate: &sync.RWMutex{},
	}

	candidate := &model.Listener{Name: "candidate", ProtocolStr: "HTTP", Protocol: model.ProtocolTypeHTTP}
	candidate.Address.SocketAddress = model.SocketAddress{Address: "127.0.0.1", Port: port}
	err = lm.ReplaceXDSListeners([]*model.Listener{candidate})

	require.ErrorContains(t, err, "bind HTTP listener")
	require.True(t, lm.HasListener(oldKey))
	require.Equal(t, []string{oldKey}, lm.XDSListenerNames())
	require.False(t, lm.HasListener(resolveListenerName(candidate)))
}

func TestListenerManager_ReplaceXDSListenersHandsOffHTTPPort(t *testing.T) {
	for _, tt := range []struct {
		name        string
		oldProtocol model.ProtocolType
		oldName     string
		newProtocol model.ProtocolType
		newName     string
	}{
		{name: "HTTP to HTTP2", oldProtocol: model.ProtocolTypeHTTP, oldName: "HTTP", newProtocol: model.ProtocolTypeHTTP2, newName: "HTTP2"},
		{name: "HTTP2 to HTTP", oldProtocol: model.ProtocolTypeHTTP2, oldName: "HTTP2", newProtocol: model.ProtocolTypeHTTP, newName: "HTTP"},
	} {
		t.Run(tt.name, func(t *testing.T) {
			reservation, err := net.Listen("tcp", "127.0.0.1:0")
			require.NoError(t, err)
			port := reservation.Addr().(*net.TCPAddr).Port
			require.NoError(t, reservation.Close())

			manager := &ListenerManager{
				bootstrap:             &model.Bootstrap{},
				activeListenerService: make(map[string]*wrapListenerService),
				xdsManaged:            make(map[string]struct{}),
				rwLock:                &sync.RWMutex{},
				updateGate:            &sync.RWMutex{},
			}
			t.Cleanup(func() { require.NoError(t, manager.ReplaceXDSListeners(nil)) })
			old := &model.Listener{Name: "before", Protocol: tt.oldProtocol, ProtocolStr: tt.oldName}
			old.Address.SocketAddress = model.SocketAddress{Address: "127.0.0.1", Port: port}
			require.NoError(t, manager.ReplaceXDSListeners([]*model.Listener{old}))

			next := &model.Listener{Name: "after", Protocol: tt.newProtocol, ProtocolStr: tt.newName}
			next.Address.SocketAddress = old.Address.SocketAddress
			require.NoError(t, manager.ReplaceXDSListeners([]*model.Listener{next}))
			require.False(t, manager.HasListener(resolveListenerName(old)))
			require.True(t, manager.HasListener(resolveListenerName(next)))
			connection, err := net.DialTimeout("tcp", reservation.Addr().String(), time.Second)
			require.NoError(t, err)
			require.NoError(t, connection.Close())
		})
	}
}

func TestListenerManager_ReplaceXDSListenersRestoresOldListenerOnHandoffBindFailure(t *testing.T) {
	occupied, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	t.Cleanup(func() { _ = occupied.Close() })
	port := occupied.Addr().(*net.TCPAddr).Port

	old := &model.Listener{Name: "old", Protocol: model.ProtocolTypeHTTP, ProtocolStr: "HTTP"}
	old.Address.SocketAddress = model.SocketAddress{Address: "127.0.0.1", Port: port}
	oldKey := resolveListenerName(old)
	service := &handoffTrackingListenerService{}
	manager := &ListenerManager{
		bootstrap: &model.Bootstrap{},
		activeListenerService: map[string]*wrapListenerService{
			oldKey: {ListenerService: service, config: old},
		},
		xdsManaged: map[string]struct{}{oldKey: {}},
		rwLock:     &sync.RWMutex{},
		updateGate: &sync.RWMutex{},
	}
	next := &model.Listener{Name: "next", Protocol: model.ProtocolTypeHTTP2, ProtocolStr: "HTTP2"}
	next.Address.SocketAddress = old.Address.SocketAddress

	err = manager.ReplaceXDSListeners([]*model.Listener{next})

	require.ErrorContains(t, err, "address already in use")
	require.Equal(t, 1, service.stops)
	require.Equal(t, 1, service.starts)
	require.True(t, manager.HasListener(oldKey))
	require.False(t, manager.HasListener(resolveListenerName(next)))
	require.Equal(t, []string{oldKey}, manager.XDSListenerNames())
}

func TestListenerManager_ReplaceXDSListenersRestoresEarlierHandoffAfterLaterFailure(t *testing.T) {
	firstPort, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	secondPort, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	ports := []net.Listener{firstPort, secondPort}
	if firstPort.Addr().String() > secondPort.Addr().String() {
		ports[0], ports[1] = ports[1], ports[0]
	}
	require.NoError(t, ports[0].Close())
	t.Cleanup(func() { _ = ports[1].Close() })

	manager := &ListenerManager{
		bootstrap:             &model.Bootstrap{},
		activeListenerService: make(map[string]*wrapListenerService),
		xdsManaged:            make(map[string]struct{}),
		rwLock:                &sync.RWMutex{},
		updateGate:            &sync.RWMutex{},
	}
	t.Cleanup(func() { require.NoError(t, manager.ReplaceXDSListeners(nil)) })
	first := &model.Listener{Name: "first", Protocol: model.ProtocolTypeHTTP, ProtocolStr: "HTTP"}
	first.Address.SocketAddress = model.SocketAddress{Address: "127.0.0.1", Port: ports[0].Addr().(*net.TCPAddr).Port}
	require.NoError(t, manager.ReplaceXDSListeners([]*model.Listener{first}))
	second := &model.Listener{Name: "second", Protocol: model.ProtocolTypeHTTP, ProtocolStr: "HTTP"}
	second.Address.SocketAddress = model.SocketAddress{Address: "127.0.0.1", Port: ports[1].Addr().(*net.TCPAddr).Port}
	secondKey := resolveListenerName(second)
	secondService := &handoffTrackingListenerService{}
	manager.activeListenerService[secondKey] = &wrapListenerService{config: second, ListenerService: secondService}
	manager.xdsManaged[secondKey] = struct{}{}
	firstNext := &model.Listener{Name: "first-next", Protocol: model.ProtocolTypeHTTP2, ProtocolStr: "HTTP2", Address: first.Address}
	secondNext := &model.Listener{Name: "second-next", Protocol: model.ProtocolTypeHTTP2, ProtocolStr: "HTTP2", Address: second.Address}

	err = manager.ReplaceXDSListeners([]*model.Listener{firstNext, secondNext})

	require.ErrorContains(t, err, "address already in use")
	require.Equal(t, 1, secondService.starts)
	require.Equal(t, 1, secondService.stops)
	require.Equal(t, []string{resolveListenerName(first), secondKey}, manager.XDSListenerNames())
	require.False(t, manager.HasListener(resolveListenerName(firstNext)))
	connection, err := net.DialTimeout("tcp", ports[0].Addr().String(), time.Second)
	require.NoError(t, err)
	require.NoError(t, connection.Close())
}

func TestListenerManager_ReplaceXDSListenersRejectsInvalidFilter(t *testing.T) {
	oldListener := &model.Listener{Name: "last-good", ProtocolStr: "HTTP", Protocol: model.ProtocolTypeHTTP}
	oldListener.Address.SocketAddress = model.SocketAddress{Address: "127.0.0.1", Port: 18080}
	oldKey := resolveListenerName(oldListener)
	lm := &ListenerManager{
		activeListenerService: map[string]*wrapListenerService{
			oldKey: {ListenerService: &transactionListenerService{config: *oldListener}, config: oldListener},
		},
		xdsManaged: map[string]struct{}{oldKey: {}},
		rwLock:     &sync.RWMutex{},
		updateGate: &sync.RWMutex{},
	}
	candidate := *oldListener
	candidate.Name = "invalid-filter-update"
	candidate.FilterChain.Filters = []model.NetworkFilter{{
		Name:   "missing.network.filter.plugin",
		Config: map[string]any{},
	}}

	err := lm.ReplaceXDSListeners([]*model.Listener{&candidate})
	require.ErrorContains(t, err, "missing.network.filter.plugin")
	require.Equal(t, []string{oldKey}, lm.XDSListenerNames())
	require.Same(t, oldListener, lm.activeListenerService[oldKey].config)
}

func TestListenerManager_ReplaceXDSListenersRollsBackEarlierRefresh(t *testing.T) {
	oldOne := &model.Listener{Name: "old-one", ProtocolStr: "HTTP", Protocol: model.ProtocolTypeHTTP}
	oldOne.Address.SocketAddress = model.SocketAddress{Address: "127.0.0.1", Port: 18080}
	oldTwo := &model.Listener{Name: "old-two", ProtocolStr: "HTTP", Protocol: model.ProtocolTypeHTTP}
	oldTwo.Address.SocketAddress = model.SocketAddress{Address: "127.0.0.1", Port: 18081}
	keyOne := resolveListenerName(oldOne)
	keyTwo := resolveListenerName(oldTwo)
	serviceOne := &transactionListenerService{config: *oldOne}
	serviceTwo := &transactionListenerService{config: *oldTwo, failName: "new-two"}
	lm := &ListenerManager{
		activeListenerService: map[string]*wrapListenerService{
			keyOne: {ListenerService: serviceOne, config: oldOne},
			keyTwo: {ListenerService: serviceTwo, config: oldTwo},
		},
		xdsManaged: map[string]struct{}{keyOne: {}, keyTwo: {}},
		rwLock:     &sync.RWMutex{},
	}
	newOne := *oldOne
	newOne.Name = "new-one"
	newTwo := *oldTwo
	newTwo.Name = "new-two"

	err := lm.ReplaceXDSListeners([]*model.Listener{&newOne, &newTwo})

	require.ErrorContains(t, err, "refresh rejected")
	assert.Equal(t, "old-one", serviceOne.config.Name)
	assert.Equal(t, "old-two", serviceTwo.config.Name)
	assert.Same(t, oldOne, lm.activeListenerService[keyOne].config)
	assert.Same(t, oldTwo, lm.activeListenerService[keyTwo].config)
}

func TestListenerManager_ReplaceXDSListenersClosesStagedOnLaterPrepareFailure(t *testing.T) {
	const trackingProtocol model.ProtocolType = 1000
	var stagedService *closeTrackingListenerService
	listener.SetListenerServiceFactory(trackingProtocol, func(*model.Listener, *model.Bootstrap) (listener.ListenerService, error) {
		stagedService = &closeTrackingListenerService{}
		return stagedService, nil
	})

	oldListener := &model.Listener{Name: "last-good", ProtocolStr: "HTTP", Protocol: model.ProtocolTypeHTTP}
	oldListener.Address.SocketAddress = model.SocketAddress{Address: "127.0.0.1", Port: 18080}
	oldKey := resolveListenerName(oldListener)
	lm := &ListenerManager{
		bootstrap: &model.Bootstrap{},
		activeListenerService: map[string]*wrapListenerService{
			oldKey: {ListenerService: &transactionListenerService{config: *oldListener, failName: "rejected-refresh"}, config: oldListener},
		},
		xdsManaged: map[string]struct{}{oldKey: {}},
		rwLock:     &sync.RWMutex{},
		updateGate: &sync.RWMutex{},
	}
	stagedListener := &model.Listener{Name: "staged", ProtocolStr: "TRACKING", Protocol: trackingProtocol}
	stagedListener.Address.SocketAddress = model.SocketAddress{Address: "127.0.0.1", Port: 18081}
	rejectedRefresh := *oldListener
	rejectedRefresh.Name = "rejected-refresh"

	err := lm.ReplaceXDSListeners([]*model.Listener{stagedListener, &rejectedRefresh})

	require.ErrorContains(t, err, "refresh rejected")
	require.NotNil(t, stagedService)
	require.Equal(t, 1, stagedService.closed)
	require.False(t, lm.HasListener(resolveListenerName(stagedListener)))
	require.Equal(t, []string{oldKey}, lm.XDSListenerNames())
}

func TestListenerManager_ReplaceXDSListenersClosesStagedOnValidationFailure(t *testing.T) {
	const trackingProtocol model.ProtocolType = 1001
	tests := []struct {
		name      string
		configure func(*ListenerManager, *model.Listener) []*model.Listener
		errorText string
	}{
		{
			name: "nil listener",
			configure: func(_ *ListenerManager, staged *model.Listener) []*model.Listener {
				return []*model.Listener{staged, nil}
			},
			errorText: "config is nil",
		},
		{
			name: "duplicate key",
			configure: func(_ *ListenerManager, staged *model.Listener) []*model.Listener {
				return []*model.Listener{staged, staged}
			},
			errorText: "duplicate xDS listener",
		},
		{
			name: "static conflict",
			configure: func(manager *ListenerManager, staged *model.Listener) []*model.Listener {
				staticListener := &model.Listener{Name: "static", ProtocolStr: "HTTP", Protocol: model.ProtocolTypeHTTP}
				staticListener.Address.SocketAddress = model.SocketAddress{Address: "127.0.0.1", Port: 18082}
				manager.activeListenerService[resolveListenerName(staticListener)] = &wrapListenerService{
					ListenerService: &legacyListenerService{config: *staticListener},
					config:          staticListener,
				}
				return []*model.Listener{staged, staticListener}
			},
			errorText: "conflicts with a non-xDS listener",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			var stagedService *closeTrackingListenerService
			listener.SetListenerServiceFactory(trackingProtocol, func(*model.Listener, *model.Bootstrap) (listener.ListenerService, error) {
				stagedService = &closeTrackingListenerService{}
				return stagedService, nil
			})
			manager := &ListenerManager{
				bootstrap:             &model.Bootstrap{},
				activeListenerService: make(map[string]*wrapListenerService),
				xdsManaged:            make(map[string]struct{}),
				rwLock:                &sync.RWMutex{},
				updateGate:            &sync.RWMutex{},
			}
			staged := &model.Listener{Name: "staged", ProtocolStr: "TRACKING", Protocol: trackingProtocol}
			staged.Address.SocketAddress = model.SocketAddress{Address: "127.0.0.1", Port: 18081}

			err := manager.ReplaceXDSListeners(tt.configure(manager, staged))

			require.ErrorContains(t, err, tt.errorText)
			require.NotNil(t, stagedService)
			require.Equal(t, 1, stagedService.closed)
			require.False(t, manager.HasListener(resolveListenerName(staged)))
		})
	}
}

type shutdownListenerService struct {
	err     error
	release <-chan struct{}
}

func (s *shutdownListenerService) Start() error { return nil }
func (s *shutdownListenerService) Close() error { return nil }
func (s *shutdownListenerService) ShutDown(wg any) error {
	waitGroup := wg.(*sync.WaitGroup)
	defer waitGroup.Done()
	if s.release != nil {
		<-s.release
	}
	return s.err
}
func (s *shutdownListenerService) Refresh(model.Listener) error { return nil }

func TestShutdownListeners(t *testing.T) {
	testErr := errors.New("shutdown failed")

	t.Run("all listeners complete", func(t *testing.T) {
		errs, timedOut := shutdownListeners([]shutdownFunc{
			func() error { return nil },
			func() error { return nil },
		}, time.Second)

		require.Empty(t, errs)
		require.False(t, timedOut)
	})

	t.Run("errors are collected", func(t *testing.T) {
		errs, timedOut := shutdownListeners([]shutdownFunc{
			func() error { return testErr },
			func() error { return nil },
		}, time.Second)

		require.Equal(t, []error{testErr}, errs)
		require.False(t, timedOut)
	})

	t.Run("overall timeout bounds the wait", func(t *testing.T) {
		release := make(chan struct{})
		t.Cleanup(func() { close(release) })

		errs, timedOut := shutdownListeners([]shutdownFunc{
			func() error {
				<-release
				return nil
			},
		}, 20*time.Millisecond)

		require.Empty(t, errs)
		require.True(t, timedOut)
	})

	t.Run("slow listener within overall budget completes", func(t *testing.T) {
		errs, timedOut := shutdownListeners([]shutdownFunc{
			func() error {
				time.Sleep(20 * time.Millisecond)
				return nil
			},
		}, time.Second)

		require.Empty(t, errs)
		require.False(t, timedOut)
	})

	t.Run("empty listener set completes", func(t *testing.T) {
		errs, timedOut := shutdownListeners(nil, time.Second)

		require.Empty(t, errs)
		require.False(t, timedOut)
	})
}

func TestListenerManagerHandleShutdownSignal(t *testing.T) {
	testErr := errors.New("shutdown failed")
	lm := &ListenerManager{
		activeListenerService: map[string]*wrapListenerService{
			"success": {ListenerService: &shutdownListenerService{}},
			"failure": {ListenerService: &shutdownListenerService{err: testErr}},
		},
		rwLock: &sync.RWMutex{},
	}

	errs, timedOut := lm.handleShutdownSignal(os.Interrupt, time.Second)

	require.Equal(t, []error{testErr}, errs)
	require.False(t, timedOut)
}

func TestListenerManagerHandleShutdownSignalTimeout(t *testing.T) {
	release := make(chan struct{})
	t.Cleanup(func() { close(release) })
	lm := &ListenerManager{
		activeListenerService: map[string]*wrapListenerService{
			"blocked": {ListenerService: &shutdownListenerService{release: release}},
		},
		rwLock: &sync.RWMutex{},
	}

	errs, timedOut := lm.handleShutdownSignal(os.Interrupt, 20*time.Millisecond)

	require.Empty(t, errs)
	require.True(t, timedOut)
}

func TestShutdownExitCode(t *testing.T) {
	testErr := errors.New("shutdown failed")
	tests := []struct {
		name     string
		errs     []error
		timedOut bool
		want     int
	}{
		{name: "success", want: 0},
		{name: "listener error", errs: []error{testErr}, want: 1},
		{name: "timeout", timedOut: true, want: 1},
		{name: "error and timeout", errs: []error{testErr}, timedOut: true, want: 1},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			require.Equal(t, tt.want, shutdownExitCode(tt.errs, tt.timedOut))
		})
	}
}

func TestShouldDumpHeap(t *testing.T) {
	require.True(t, shouldDumpHeap(syscall.SIGQUIT))
	require.False(t, shouldDumpHeap(os.Interrupt))
}

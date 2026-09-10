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

package grpcproxy

import (
	"container/list"
	"context"
	"fmt"
	"sync"
	"time"
)

import (
	"golang.org/x/sync/singleflight"

	"google.golang.org/grpc"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
)

const (
	defaultGRPCDialTimeout = 5 * time.Second
	// Endpoint tombstones only need to cover recently delivered lifecycle
	// events. Active connections and in-flight dials are pinned separately and
	// are never evicted by this bound.
	maxEndpointTombstones = 1024
)

type grpcConnectionDialer func(context.Context, string) (*grpc.ClientConn, error)

// grpcConnectionManager owns long-lived backend connections for the HTTP gRPC
// proxy. A grpc.ClientConn is safe for concurrent use and multiplexes calls
// over HTTP/2, so a sync.Pool is both unnecessary and incorrect here.
type grpcConnectionManager struct {
	connections sync.Map
	creates     singleflight.Group
	dial        grpcConnectionDialer
	dialTimeout time.Duration
	onRemove    func(*grpc.ClientConn)

	mu                  sync.Mutex
	closed              bool
	endpointGenerations map[string]uint64
	endpointEventVers   map[string]uint64
	endpointRefs        map[string]int
	endpointRemoved     map[string]bool
	endpointTombstones  map[string]*list.Element
	tombstoneOrder      *list.List
	// endpointPresent is an optional authoritative snapshot check used after
	// bounded tombstone metadata has been evicted. It is intentionally kept
	// outside the manager so direct-endpoint users do not need a cluster store.
	endpointPresent func(key, endpoint string) bool
}

func newGRPCConnectionManager() *grpcConnectionManager {
	return &grpcConnectionManager{
		dial:                dialGRPCConnection,
		dialTimeout:         defaultGRPCDialTimeout,
		endpointGenerations: make(map[string]uint64),
		endpointEventVers:   make(map[string]uint64),
		endpointRefs:        make(map[string]int),
		endpointRemoved:     make(map[string]bool),
		endpointTombstones:  make(map[string]*list.Element),
		tombstoneOrder:      list.New(),
	}
}

func (m *grpcConnectionManager) initEndpointStateLocked() {
	if m.endpointGenerations == nil {
		m.endpointGenerations = make(map[string]uint64)
	}
	if m.endpointEventVers == nil {
		m.endpointEventVers = make(map[string]uint64)
	}
	if m.endpointRefs == nil {
		m.endpointRefs = make(map[string]int)
	}
	if m.endpointRemoved == nil {
		m.endpointRemoved = make(map[string]bool)
	}
	if m.endpointTombstones == nil {
		m.endpointTombstones = make(map[string]*list.Element)
	}
	if m.tombstoneOrder == nil {
		m.tombstoneOrder = list.New()
	}
}

func (m *grpcConnectionManager) discardEndpointTombstoneLocked(key string) {
	if element, ok := m.endpointTombstones[key]; ok {
		m.tombstoneOrder.Remove(element)
		delete(m.endpointTombstones, key)
	}
}

func (m *grpcConnectionManager) discardEndpointStateLocked(key string) {
	m.discardEndpointTombstoneLocked(key)
	delete(m.endpointGenerations, key)
	delete(m.endpointEventVers, key)
	delete(m.endpointRemoved, key)
}

func (m *grpcConnectionManager) rememberEndpointTombstoneLocked(key string) {
	m.initEndpointStateLocked()
	if element, ok := m.endpointTombstones[key]; ok {
		m.tombstoneOrder.MoveToFront(element)
		return
	}
	element := m.tombstoneOrder.PushFront(key)
	m.endpointTombstones[key] = element

	for len(m.endpointTombstones) > maxEndpointTombstones {
		var evict *list.Element
		for element := m.tombstoneOrder.Back(); element != nil; element = element.Prev() {
			candidate := element.Value.(string)
			if m.endpointRefs[candidate] == 0 && !m.hasConnection(candidate) {
				evict = element
				break
			}
		}
		if evict == nil {
			return
		}
		candidate := evict.Value.(string)
		m.tombstoneOrder.Remove(evict)
		delete(m.endpointTombstones, candidate)
		delete(m.endpointGenerations, candidate)
		delete(m.endpointEventVers, candidate)
		delete(m.endpointRemoved, candidate)
	}
}

func (m *grpcConnectionManager) hasConnection(key string) bool {
	_, ok := m.connections.Load(key)
	return ok
}

func (m *grpcConnectionManager) pinEndpoint(key string) (uint64, error) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.initEndpointStateLocked()
	if m.closed {
		return 0, fmt.Errorf("grpc connection manager is closed")
	}
	if m.endpointRemoved[key] {
		return 0, fmt.Errorf("grpc endpoint was removed")
	}
	m.endpointRefs[key]++
	return m.endpointGenerations[key], nil
}

func (m *grpcConnectionManager) isEndpointPresent(key, endpoint string) bool {
	if m.endpointPresent == nil {
		return true
	}
	return m.endpointPresent(key, endpoint)
}

func (m *grpcConnectionManager) unpinEndpoint(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.endpointRefs[key] > 1 {
		m.endpointRefs[key]--
		return
	}
	delete(m.endpointRefs, key)
	if m.endpointRemoved[key] && !m.hasConnection(key) {
		m.rememberEndpointTombstoneLocked(key)
	} else if !m.endpointRemoved[key] && m.endpointEventVers[key] == 0 && !m.hasConnection(key) {
		// Get may create a temporary generation entry for a direct endpoint
		// before the cluster manager has delivered a lifecycle event. Do not
		// retain that request-only state after a failed or canceled dial.
		m.discardEndpointStateLocked(key)
	}
}

func (m *grpcConnectionManager) finalizeRemovedEndpoint(key string) {
	m.mu.Lock()
	defer m.mu.Unlock()
	if m.endpointRemoved[key] && m.endpointRefs[key] == 0 && !m.hasConnection(key) {
		m.rememberEndpointTombstoneLocked(key)
	}
}

func dialGRPCConnection(ctx context.Context, endpoint string) (*grpc.ClientConn, error) {
	return grpc.DialContext( //nolint:staticcheck // SA1019: the context is required to enforce the dial timeout.
		ctx,
		endpoint,
		grpc.WithTransportCredentials(insecure.NewCredentials()),
	)
}

func (m *grpcConnectionManager) Get(ctx context.Context, key, endpoint string) (*grpc.ClientConn, error) {
	if key == "" || endpoint == "" {
		return nil, fmt.Errorf("grpc connection key and endpoint must not be empty")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	if !m.isEndpointPresent(key, endpoint) {
		return nil, fmt.Errorf("grpc endpoint was removed")
	}

	if conn, ok := m.loadHealthy(key); ok {
		if !m.isEndpointPresent(key, endpoint) {
			m.remove(key, conn)
			return nil, fmt.Errorf("grpc endpoint was removed")
		}
		return conn, nil
	}

	endpointGeneration, err := m.pinEndpoint(key)
	if err != nil {
		return nil, err
	}
	createKey := fmt.Sprintf("%s\x00%d", key, endpointGeneration)
	result := m.creates.DoChan(createKey, func() (any, error) {
		if !m.isEndpointPresent(key, endpoint) {
			return nil, fmt.Errorf("grpc endpoint was removed")
		}
		m.mu.Lock()
		currentGeneration := m.endpointGenerations[key]
		removed := m.endpointRemoved[key]
		m.mu.Unlock()
		if removed || currentGeneration != endpointGeneration {
			return nil, fmt.Errorf("grpc endpoint was removed while connecting")
		}
		if conn, ok := m.loadHealthy(key); ok {
			if !m.isEndpointPresent(key, endpoint) {
				m.remove(key, conn)
				return nil, fmt.Errorf("grpc endpoint was removed")
			}
			return conn, nil
		}

		m.mu.Lock()
		if m.closed {
			m.mu.Unlock()
			return nil, fmt.Errorf("grpc connection manager is closed")
		}
		dial := m.dial
		dialTimeout := m.dialTimeout
		m.mu.Unlock()

		dialCtx, cancel := context.WithTimeout(context.Background(), dialTimeout)
		defer cancel()
		conn, err := dial(dialCtx, endpoint)
		if err != nil {
			return nil, err
		}

		m.mu.Lock()
		closed := m.closed
		removed = m.endpointRemoved[key] || m.endpointGenerations[key] != endpointGeneration
		if !closed && !removed {
			m.connections.Store(key, conn)
		}
		m.mu.Unlock()
		if closed || removed || !m.isEndpointPresent(key, endpoint) {
			removedFromManager := false
			if !closed && !removed {
				removedFromManager = m.remove(key, conn)
			}
			if !removedFromManager {
				_ = conn.Close()
			}
			if closed {
				return nil, fmt.Errorf("grpc connection manager closed while dialing")
			}
			return nil, fmt.Errorf("grpc endpoint was removed while connecting")
		}

		return conn, nil
	})

	select {
	case <-ctx.Done():
		go func() {
			<-result
			m.unpinEndpoint(key)
		}()
		return nil, ctx.Err()
	case result := <-result:
		m.unpinEndpoint(key)
		if result.Err != nil {
			return nil, result.Err
		}
		return result.Val.(*grpc.ClientConn), nil
	}
}

func (m *grpcConnectionManager) loadHealthy(key string) (*grpc.ClientConn, bool) {
	m.mu.Lock()
	removed := m.endpointRemoved[key]
	m.mu.Unlock()
	if removed {
		return nil, false
	}
	value, ok := m.connections.Load(key)
	if !ok {
		return nil, false
	}

	conn, ok := value.(*grpc.ClientConn)
	if ok && m.isHealthy(conn) {
		return conn, true
	}
	if ok {
		m.remove(key, conn)
	}
	return nil, false
}

func (m *grpcConnectionManager) isHealthy(conn *grpc.ClientConn) bool {
	if conn == nil {
		return false
	}
	state := conn.GetState()
	return state != connectivity.Shutdown
}

// Invalidate removes a connection only when its transport is known to be
// unusable. Application-level RPC errors must not cause healthy connections to
// churn.
func (m *grpcConnectionManager) Invalidate(key string, conn *grpc.ClientConn) {
	if conn == nil || m.isHealthy(conn) {
		return
	}
	m.remove(key, conn)
}

// RemoveEndpoint closes and removes the connection for a deleted endpoint.
func (m *grpcConnectionManager) RemoveEndpoint(clusterName, endpoint string) {
	key := grpcConnectionKey(clusterName, endpoint)
	m.mu.Lock()
	m.initEndpointStateLocked()
	if m.closed {
		m.mu.Unlock()
		return
	}
	m.endpointGenerations[key]++
	m.endpointRemoved[key] = true
	m.mu.Unlock()
	value, ok := m.connections.Load(key)
	if !ok {
		m.finalizeRemovedEndpoint(key)
		return
	}
	conn, ok := value.(*grpc.ClientConn)
	if ok {
		m.remove(key, conn)
	}
	m.finalizeRemovedEndpoint(key)
}

// UpdateEndpointState applies an ordered endpoint lifecycle event. Additions
// advance the same generation as removals so a delayed removal cannot delete
// a connection created after the endpoint was re-added.
func (m *grpcConnectionManager) UpdateEndpointState(clusterName, endpoint string, present bool, eventVersion uint64) {
	key := grpcConnectionKey(clusterName, endpoint)
	// Lifecycle callbacks are delivered outside the cluster manager lock and can
	// therefore arrive out of order. Once a bounded tombstone has been evicted,
	// its version watermark is no longer available to reject an old removal.
	// The cluster snapshot is authoritative in that case: do not let a delayed
	// removal turn an endpoint that is currently present back into a tombstone.
	if !present && m.isEndpointPresent(key, endpoint) {
		return
	}
	m.mu.Lock()
	m.initEndpointStateLocked()
	if m.closed {
		m.mu.Unlock()
		return
	}
	if eventVersion <= m.endpointEventVers[key] {
		m.mu.Unlock()
		return
	}
	m.endpointEventVers[key] = eventVersion
	m.endpointGenerations[key]++
	m.endpointRemoved[key] = !present
	if present {
		m.discardEndpointTombstoneLocked(key)
	}
	var conn *grpc.ClientConn
	if value, ok := m.connections.Load(key); ok {
		conn, _ = value.(*grpc.ClientConn)
	}
	m.mu.Unlock()
	if conn != nil {
		// A present event starts a new endpoint incarnation. Remove any
		// connection left over from the previous incarnation before a new Get
		// can reuse it.
		m.remove(key, conn)
	}
	if !present {
		m.finalizeRemovedEndpoint(key)
	}
}

func (m *grpcConnectionManager) remove(key string, expected *grpc.ClientConn) bool {
	value, ok := m.connections.Load(key)
	if !ok || value != expected {
		return false
	}
	if m.connections.CompareAndDelete(key, expected) {
		_ = expected.Close()
		if m.onRemove != nil {
			m.onRemove(expected)
		}
		return true
	}
	return false
}

func (m *grpcConnectionManager) Close() error {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return nil
	}
	m.closed = true
	m.endpointGenerations = nil
	m.endpointEventVers = nil
	m.endpointRefs = nil
	m.endpointRemoved = nil
	m.endpointTombstones = nil
	m.tombstoneOrder = nil
	var connections []*grpc.ClientConn
	m.connections.Range(func(key, value any) bool {
		m.connections.Delete(key)
		if conn, ok := value.(*grpc.ClientConn); ok {
			connections = append(connections, conn)
		}
		return true
	})
	m.mu.Unlock()

	var firstErr error
	for _, conn := range connections {
		if err := conn.Close(); err != nil && firstErr == nil {
			firstErr = err
		}
		if m.onRemove != nil {
			m.onRemove(conn)
		}
	}
	return firstErr
}

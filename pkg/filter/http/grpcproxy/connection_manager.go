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

const defaultGRPCDialTimeout = 5 * time.Second

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
}

func newGRPCConnectionManager() *grpcConnectionManager {
	return &grpcConnectionManager{
		dial:                dialGRPCConnection,
		dialTimeout:         defaultGRPCDialTimeout,
		endpointGenerations: make(map[string]uint64),
		endpointEventVers:   make(map[string]uint64),
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

	if conn, ok := m.loadHealthy(key); ok {
		return conn, nil
	}

	m.mu.Lock()
	if m.endpointGenerations == nil {
		m.endpointGenerations = make(map[string]uint64)
	}
	endpointGeneration := m.endpointGenerations[key]
	m.mu.Unlock()
	createKey := fmt.Sprintf("%s\x00%d", key, endpointGeneration)
	result := m.creates.DoChan(createKey, func() (any, error) {
		m.mu.Lock()
		currentGeneration := m.endpointGenerations[key]
		m.mu.Unlock()
		if currentGeneration != endpointGeneration {
			return nil, fmt.Errorf("grpc endpoint was removed while connecting")
		}
		if conn, ok := m.loadHealthy(key); ok {
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
		removed := m.endpointGenerations[key] != endpointGeneration
		if !closed && !removed {
			m.connections.Store(key, conn)
		}
		m.mu.Unlock()
		if closed || removed {
			_ = conn.Close()
			if closed {
				return nil, fmt.Errorf("grpc connection manager closed while dialing")
			}
			return nil, fmt.Errorf("grpc endpoint was removed while connecting")
		}

		return conn, nil
	})

	select {
	case <-ctx.Done():
		return nil, ctx.Err()
	case result := <-result:
		if result.Err != nil {
			return nil, result.Err
		}
		return result.Val.(*grpc.ClientConn), nil
	}
}

func (m *grpcConnectionManager) loadHealthy(key string) (*grpc.ClientConn, bool) {
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
	return state != connectivity.Shutdown && state != connectivity.TransientFailure
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
	if m.endpointGenerations == nil {
		m.endpointGenerations = make(map[string]uint64)
	}
	m.endpointGenerations[key]++
	m.mu.Unlock()
	value, ok := m.connections.Load(key)
	if !ok {
		return
	}
	conn, ok := value.(*grpc.ClientConn)
	if ok {
		m.remove(key, conn)
	}
}

// UpdateEndpointState applies an ordered endpoint lifecycle event. Additions
// advance the same generation as removals so a delayed removal cannot delete
// a connection created after the endpoint was re-added.
func (m *grpcConnectionManager) UpdateEndpointState(clusterName, endpoint string, present bool, eventVersion uint64) {
	key := grpcConnectionKey(clusterName, endpoint)
	m.mu.Lock()
	if m.endpointGenerations == nil {
		m.endpointGenerations = make(map[string]uint64)
	}
	if m.endpointEventVers == nil {
		m.endpointEventVers = make(map[string]uint64)
	}
	if eventVersion <= m.endpointEventVers[key] {
		m.mu.Unlock()
		return
	}
	m.endpointEventVers[key] = eventVersion
	m.endpointGenerations[key]++
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
}

func (m *grpcConnectionManager) remove(key string, expected *grpc.ClientConn) {
	value, ok := m.connections.Load(key)
	if !ok || value != expected {
		return
	}
	if m.connections.CompareAndDelete(key, expected) {
		_ = expected.Close()
		if m.onRemove != nil {
			m.onRemove(expected)
		}
	}
}

func (m *grpcConnectionManager) Close() error {
	m.mu.Lock()
	if m.closed {
		m.mu.Unlock()
		return nil
	}
	m.closed = true
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

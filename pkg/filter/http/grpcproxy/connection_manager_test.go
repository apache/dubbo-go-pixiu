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
	"net"
	"sync"
	"sync/atomic"
	"testing"
	"time"
)

import (
	"github.com/jhump/protoreflect/desc"            //nolint:staticcheck // legacy descriptor API used by grpcproxy.
	"github.com/jhump/protoreflect/desc/protoparse" //nolint:staticcheck // legacy parser used by grpcproxy.

	"github.com/stretchr/testify/require"

	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/connectivity"
	"google.golang.org/grpc/credentials/insecure"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"google.golang.org/grpc/status"
)

import (
	ct "github.com/apache/dubbo-go-pixiu/pkg/context"
)

func startTestGRPCServer(t *testing.T) string {
	t.Helper()

	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	server := grpc.NewServer()
	go func() {
		_ = server.Serve(listener)
	}()
	t.Cleanup(func() {
		server.Stop()
		_ = listener.Close()
	})
	return listener.Addr().String()
}

func testConnectionManager(t *testing.T, calls *atomic.Int32) *grpcConnectionManager {
	t.Helper()
	return &grpcConnectionManager{
		dial: func(ctx context.Context, endpoint string) (*grpc.ClientConn, error) {
			calls.Add(1)
			return grpc.DialContext(ctx, endpoint, grpc.WithTransportCredentials(insecure.NewCredentials())) //nolint:staticcheck // the test verifies context-bounded dialing.
		},
		dialTimeout: time.Second,
	}
}

func TestGRPCConnectionManagerSharesConcurrentConnection(t *testing.T) {
	endpoint := startTestGRPCServer(t)
	var dialCalls atomic.Int32
	manager := testConnectionManager(t, &dialCalls)
	t.Cleanup(func() { require.NoError(t, manager.Close()) })

	const requests = 64
	connections := make([]*grpc.ClientConn, requests)
	errs := make(chan error, requests)
	var wg sync.WaitGroup
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			conn, err := manager.Get(context.Background(), grpcConnectionKey("cluster", endpoint), endpoint)
			if err != nil {
				errs <- err
				return
			}
			connections[index] = conn
		}(i)
	}
	wg.Wait()
	close(errs)
	for err := range errs {
		require.NoError(t, err)
	}

	require.Equal(t, int32(1), dialCalls.Load())
	for _, conn := range connections[1:] {
		require.Same(t, connections[0], conn)
	}
}

func TestGRPCConnectionManagerRecreatesUnhealthyConnection(t *testing.T) {
	endpoint := startTestGRPCServer(t)
	var dialCalls atomic.Int32
	manager := testConnectionManager(t, &dialCalls)
	t.Cleanup(func() { require.NoError(t, manager.Close()) })

	key := grpcConnectionKey("cluster", endpoint)
	first, err := manager.Get(context.Background(), key, endpoint)
	require.NoError(t, err)
	require.NoError(t, first.Close())
	manager.Invalidate(key, first)

	second, err := manager.Get(context.Background(), key, endpoint)
	require.NoError(t, err)
	require.NotSame(t, first, second)
	require.Equal(t, int32(2), dialCalls.Load())
}

func TestGRPCConnectionManagerRetainsTransientFailureConnection(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)
	endpoint := listener.Addr().String()
	require.NoError(t, listener.Close())

	conn, err := grpc.NewClient("passthrough:///"+endpoint, grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	conn.Connect()
	require.Eventually(t, func() bool {
		return conn.GetState() == connectivity.TransientFailure
	}, time.Second, time.Millisecond)

	manager := newGRPCConnectionManager()
	key := grpcConnectionKey("cluster", endpoint)
	manager.connections.Store(key, conn)
	t.Cleanup(func() { require.NoError(t, manager.Close()) })

	got, ok := manager.loadHealthy(key)
	require.True(t, ok)
	require.Same(t, conn, got)
	manager.Invalidate(key, conn)
	_, stillCached := manager.connections.Load(key)
	require.True(t, stillCached)
}

func TestGRPCConnectionManagerHonorsCallerTimeoutWhileCreating(t *testing.T) {
	var dialCalls atomic.Int32
	manager := &grpcConnectionManager{
		dial: func(ctx context.Context, _ string) (*grpc.ClientConn, error) {
			dialCalls.Add(1)
			<-ctx.Done()
			return nil, ctx.Err()
		},
		dialTimeout: time.Second,
	}

	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Millisecond)
	defer cancel()
	_, err := manager.Get(ctx, "cluster\x00endpoint", "endpoint")
	require.ErrorIs(t, err, context.DeadlineExceeded)
	require.Equal(t, int32(1), dialCalls.Load())
	require.NoError(t, manager.Close())
}

func TestGRPCConnectionManagerDoesNotPublishRemovedEndpointAfterDial(t *testing.T) {
	endpoint := startTestGRPCServer(t)
	dialStarted := make(chan struct{})
	releaseDial := make(chan struct{})
	manager := &grpcConnectionManager{
		dial: func(ctx context.Context, endpoint string) (*grpc.ClientConn, error) {
			close(dialStarted)
			<-releaseDial
			return grpc.DialContext(ctx, endpoint, grpc.WithTransportCredentials(insecure.NewCredentials())) //nolint:staticcheck // the test verifies endpoint lifecycle.
		},
		dialTimeout: time.Second,
	}
	t.Cleanup(func() { require.NoError(t, manager.Close()) })

	key := grpcConnectionKey("cluster", endpoint)
	result := make(chan error, 1)
	go func() {
		_, err := manager.Get(context.Background(), key, endpoint)
		result <- err
	}()
	<-dialStarted
	manager.RemoveEndpoint("cluster", endpoint)
	close(releaseDial)

	require.Error(t, <-result)
	_, ok := manager.connections.Load(key)
	require.False(t, ok, "a removed endpoint must not be published after dialing")
}

func TestGRPCConnectionManagerBoundsEndpointTombstones(t *testing.T) {
	manager := newGRPCConnectionManager()
	t.Cleanup(func() { require.NoError(t, manager.Close()) })

	const churn = maxEndpointTombstones * 4
	for i := 0; i < churn; i++ {
		manager.UpdateEndpointState(
			"cluster",
			fmt.Sprintf("127.0.0.1:%d", 20000+i),
			false,
			uint64(i+1),
		)
	}

	manager.mu.Lock()
	defer manager.mu.Unlock()
	require.LessOrEqual(t, len(manager.endpointGenerations), maxEndpointTombstones)
	require.LessOrEqual(t, len(manager.endpointEventVers), maxEndpointTombstones)
	require.LessOrEqual(t, len(manager.endpointTombstones), maxEndpointTombstones)
}

func TestGRPCConnectionManagerRejectsEvictedRemovedEndpointFromSnapshot(t *testing.T) {
	endpoint := "127.0.0.1:20000"
	current := make(map[string]bool)
	var currentMu sync.Mutex
	var dialCalls atomic.Int32
	manager := &grpcConnectionManager{
		dial: func(context.Context, string) (*grpc.ClientConn, error) {
			dialCalls.Add(1)
			return nil, fmt.Errorf("unexpected dial")
		},
		dialTimeout: time.Second,
		endpointPresent: func(_, address string) bool {
			currentMu.Lock()
			defer currentMu.Unlock()
			return current[address]
		},
	}
	t.Cleanup(func() { require.NoError(t, manager.Close()) })

	current[endpoint] = true
	manager.UpdateEndpointState("cluster", endpoint, true, 1)
	current[endpoint] = false
	manager.UpdateEndpointState("cluster", endpoint, false, 2)
	for i := 0; i < maxEndpointTombstones+1; i++ {
		address := fmt.Sprintf("127.0.0.1:%d", 21000+i)
		manager.UpdateEndpointState("cluster", address, false, uint64(i+3))
	}

	_, err := manager.Get(context.Background(), grpcConnectionKey("cluster", endpoint), endpoint)
	require.EqualError(t, err, "grpc endpoint was removed")
	require.Zero(t, dialCalls.Load())
}

func TestGRPCConnectionManagerIgnoresStaleRemovalAfterTombstoneEviction(t *testing.T) {
	endpoint := startTestGRPCServer(t)
	current := make(map[string]bool)
	var currentMu sync.Mutex
	var dialCalls atomic.Int32
	manager := testConnectionManager(t, &dialCalls)
	manager.endpointPresent = func(_, address string) bool {
		currentMu.Lock()
		defer currentMu.Unlock()
		return current[address]
	}
	t.Cleanup(func() { require.NoError(t, manager.Close()) })

	setPresent := func(address string, present bool) {
		currentMu.Lock()
		current[address] = present
		currentMu.Unlock()
	}

	setPresent(endpoint, true)
	manager.UpdateEndpointState("cluster", endpoint, true, 1)
	setPresent(endpoint, false)
	manager.UpdateEndpointState("cluster", endpoint, false, 100)
	for i := 0; i < maxEndpointTombstones+1; i++ {
		manager.UpdateEndpointState(
			"cluster",
			fmt.Sprintf("127.0.0.1:%d", 21000+i),
			false,
			uint64(i+101),
		)
	}

	// The authoritative snapshot has re-added the endpoint, but the matching
	// present callback has not arrived yet. A delayed removal must not turn the
	// current endpoint back into a tombstone after its old version was evicted.
	setPresent(endpoint, true)
	manager.UpdateEndpointState("cluster", endpoint, false, 50)

	conn, err := manager.Get(context.Background(), grpcConnectionKey("cluster", endpoint), endpoint)
	require.NoError(t, err)
	require.NotNil(t, conn)
	require.Equal(t, int32(1), dialCalls.Load())
}

func TestGRPCConnectionManagerReclaimsRequestOnlyEndpointState(t *testing.T) {
	manager := &grpcConnectionManager{
		dial: func(context.Context, string) (*grpc.ClientConn, error) {
			return nil, fmt.Errorf("dial failed")
		},
		dialTimeout: time.Second,
	}
	t.Cleanup(func() { require.NoError(t, manager.Close()) })

	_, err := manager.Get(context.Background(), "cluster\x00endpoint", "endpoint")
	require.EqualError(t, err, "dial failed")

	manager.mu.Lock()
	defer manager.mu.Unlock()
	require.Empty(t, manager.endpointGenerations)
	require.Empty(t, manager.endpointEventVers)
}

func TestGRPCConnectionManagerClosePreventsNewConnections(t *testing.T) {
	var dialCalls atomic.Int32
	manager := testConnectionManager(t, &dialCalls)
	require.NoError(t, manager.Close())

	_, err := manager.Get(context.Background(), grpcConnectionKey("cluster", "endpoint"), "endpoint")
	require.Error(t, err)
	require.Equal(t, int32(0), dialCalls.Load())
}

func TestDescriptorSourceHonorsRequestTimeout(t *testing.T) {
	listener, err := net.Listen("tcp", "127.0.0.1:0")
	require.NoError(t, err)

	server := grpc.NewServer()
	reflectpb.RegisterServerReflectionServer(server, &blockingReflectionServer{})
	go func() {
		_ = server.Serve(listener)
	}()
	t.Cleanup(func() {
		server.Stop()
		_ = listener.Close()
	})

	conn, err := grpc.NewClient(listener.Addr().String(), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close()) })

	descriptor := &Descriptor{}
	ctx, cancel := context.WithTimeout(context.Background(), 100*time.Millisecond)
	defer cancel()
	source, err := descriptor.getServerDescriptorSourceCtx(
		context.WithValue(ctx, ct.ContextKey(GrpcClientConnKey), conn),
		&Config{},
	)
	require.NoError(t, err)

	_, err = source.FindSymbol("test.SlowService")
	require.Error(t, err)
	if st, ok := status.FromError(err); ok {
		require.Equal(t, codes.DeadlineExceeded, st.Code())
	} else {
		require.ErrorIs(t, err, context.DeadlineExceeded)
	}
}

type blockingReflectionServer struct{}

func (blockingReflectionServer) ServerReflectionInfo(stream reflectpb.ServerReflection_ServerReflectionInfoServer) error {
	<-stream.Context().Done()
	return stream.Context().Err()
}

type countingDescriptorSource struct {
	descriptor desc.Descriptor
	findCalls  atomic.Int32
}

func (s *countingDescriptorSource) ListServices() ([]string, error) {
	return []string{"test.Greeter"}, nil
}

func (s *countingDescriptorSource) FindSymbol(string) (desc.Descriptor, error) {
	s.findCalls.Add(1)
	return s.descriptor, nil
}

func (s *countingDescriptorSource) AllExtensionsForType(string) ([]*desc.FieldDescriptor, error) {
	return nil, nil
}

func TestDescriptorCachesMethodLookupPerConnection(t *testing.T) {
	files, err := (protoparse.Parser{
		Accessor: protoparse.FileContentsFromMap(map[string]string{
			"test.proto": `syntax = "proto3";
package test;

service Greeter {
  rpc Hello(Request) returns (Response);
}

message Request {}
message Response {}
`,
		}),
	}).ParseFiles("test.proto")
	require.NoError(t, err)

	source := &countingDescriptorSource{descriptor: files[0].FindSymbol("test.Greeter")}
	descriptor := &Descriptor{}
	conn := &grpc.ClientConn{}
	_, err = descriptor.getMethodDescriptor(source, conn, "test.Greeter", "Hello")
	require.NoError(t, err)

	const requests = 32
	methods := make([]*desc.MethodDescriptor, requests)
	errs := make([]error, requests)
	var wg sync.WaitGroup
	for i := 0; i < requests; i++ {
		wg.Add(1)
		go func(index int) {
			defer wg.Done()
			methods[index], errs[index] = descriptor.getMethodDescriptor(source, conn, "test.Greeter", "Hello")
		}(i)
	}
	wg.Wait()

	for index, method := range methods {
		require.NoError(t, errs[index])
		require.NotNil(t, method)
		require.Equal(t, "Hello", method.GetName())
	}
	require.Equal(t, int32(1), source.findCalls.Load())
}

func TestDescriptorLookupDoesNotShareRequestBoundSource(t *testing.T) {
	files, err := (protoparse.Parser{
		Accessor: protoparse.FileContentsFromMap(map[string]string{
			"test.proto": `syntax = "proto3";
package test;

service Greeter {
  rpc Hello(Request) returns (Response);
}

message Request {}
message Response {}
`,
		}),
	}).ParseFiles("test.proto")
	require.NoError(t, err)

	firstSource := &blockingDescriptorSource{
		descriptor: files[0].FindSymbol("test.Greeter"),
		started:    make(chan struct{}),
		release:    make(chan struct{}),
	}
	secondSource := &countingDescriptorSource{descriptor: files[0].FindSymbol("test.Greeter")}
	descriptor := &Descriptor{}
	conn := &grpc.ClientConn{}

	firstResult := make(chan error, 1)
	go func() {
		_, lookupErr := descriptor.getMethodDescriptor(firstSource, conn, "test.Greeter", "Hello")
		firstResult <- lookupErr
	}()
	<-firstSource.started

	secondResult := make(chan error, 1)
	go func() {
		_, lookupErr := descriptor.getMethodDescriptor(secondSource, conn, "test.Greeter", "Hello")
		secondResult <- lookupErr
	}()

	select {
	case lookupErr := <-secondResult:
		require.NoError(t, lookupErr)
	case <-time.After(200 * time.Millisecond):
		close(firstSource.release)
		t.Fatal("request-bound descriptor lookup was shared with a canceled/slow caller")
	}

	close(firstSource.release)
	require.NoError(t, <-firstResult)
}

type blockingDescriptorSource struct {
	descriptor desc.Descriptor
	started    chan struct{}
	release    chan struct{}
}

func (s *blockingDescriptorSource) ListServices() ([]string, error) {
	return []string{"test.Greeter"}, nil
}

func (s *blockingDescriptorSource) FindSymbol(string) (desc.Descriptor, error) {
	close(s.started)
	<-s.release
	return s.descriptor, nil
}

func (s *blockingDescriptorSource) AllExtensionsForType(string) ([]*desc.FieldDescriptor, error) {
	return nil, nil
}

func TestDescriptorRemovalForOneConnectionDoesNotInvalidateAnother(t *testing.T) {
	files, err := (protoparse.Parser{
		Accessor: protoparse.FileContentsFromMap(map[string]string{
			"test.proto": `syntax = "proto3";
package test;

service Greeter {
  rpc Hello(Request) returns (Response);
}

message Request {}
message Response {}
`,
		}),
	}).ParseFiles("test.proto")
	require.NoError(t, err)

	source := &blockingDescriptorSource{
		descriptor: files[0].FindSymbol("test.Greeter"),
		started:    make(chan struct{}),
		release:    make(chan struct{}),
	}
	descriptor := &Descriptor{}
	connA := &grpc.ClientConn{}
	connB := &grpc.ClientConn{}
	_, err = descriptor.getMethodDescriptor(
		&countingDescriptorSource{descriptor: files[0].FindSymbol("test.Greeter")},
		connB, "test.Greeter", "Hello",
	)
	require.NoError(t, err)

	result := make(chan error, 1)
	go func() {
		_, lookupErr := descriptor.getMethodDescriptor(source, connA, "test.Greeter", "Hello")
		result <- lookupErr
	}()
	<-source.started

	descriptor.removeConnection(connB)
	close(source.release)

	require.NoError(t, <-result)
}

func TestDescriptorRemovalSeparatesNewLookupOnSameConnection(t *testing.T) {
	files, err := (protoparse.Parser{
		Accessor: protoparse.FileContentsFromMap(map[string]string{
			"test.proto": `syntax = "proto3";
package test;

service Greeter {
  rpc Hello(Request) returns (Response);
}

message Request {}
message Response {}
`,
		}),
	}).ParseFiles("test.proto")
	require.NoError(t, err)

	oldSource := &blockingDescriptorSource{
		descriptor: files[0].FindSymbol("test.Greeter"),
		started:    make(chan struct{}),
		release:    make(chan struct{}),
	}
	newSource := &countingDescriptorSource{descriptor: files[0].FindSymbol("test.Greeter")}
	descriptor := &Descriptor{}
	conn := &grpc.ClientConn{}

	oldResult := make(chan error, 1)
	go func() {
		_, lookupErr := descriptor.getMethodDescriptor(oldSource, conn, "test.Greeter", "Hello")
		oldResult <- lookupErr
	}()
	<-oldSource.started
	descriptor.removeConnection(conn)

	newResult := make(chan error, 1)
	go func() {
		_, lookupErr := descriptor.getMethodDescriptor(newSource, conn, "test.Greeter", "Hello")
		newResult <- lookupErr
	}()
	select {
	case lookupErr := <-newResult:
		require.NoError(t, lookupErr)
	case <-time.After(200 * time.Millisecond):
		close(oldSource.release)
		t.Fatal("new lookup joined the removed connection's in-flight lookup")
	}

	close(oldSource.release)
	require.Error(t, <-oldResult)
}

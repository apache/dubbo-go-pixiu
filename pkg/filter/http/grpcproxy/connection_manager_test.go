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

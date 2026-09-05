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
	"os"
	"path/filepath"
	"testing"
	"time"
)

import (
	"github.com/jhump/protoreflect/desc"            //nolint:staticcheck // legacy descriptor API used by grpcproxy.
	"github.com/jhump/protoreflect/desc/protoparse" //nolint:staticcheck // legacy parser used by grpcproxy.

	"github.com/stretchr/testify/require"

	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
)

import (
	ct "github.com/apache/dubbo-go-pixiu/pkg/context"
)

const testGreeterProto = `syntax = "proto3";
package test;

service Greeter {
  rpc Hello(Request) returns (Response);
}

message Request {}
message Response {}
`

// writeTestProtoDir writes a proto file to a temporary directory usable as Config.Path.
func writeTestProtoDir(t *testing.T) string {
	t.Helper()

	dir := t.TempDir()
	require.NoError(t, os.WriteFile(filepath.Join(dir, "test.proto"), []byte(testGreeterProto), 0o600))
	return dir
}

// TestDescriptorAutoStrategyFallsBackToLocalProtoFiles covers the AUTO strategy, which is
// the default one, against a backend that does not implement the server reflection API.
// The documented `file + reflection` fallback has to resolve the method from the local
// proto files instead of dereferencing an uninitialized file source.
func TestDescriptorAutoStrategyFallsBackToLocalProtoFiles(t *testing.T) {
	cfg := &Config{DescriptorSourceStrategy: AUTO, Path: writeTestProtoDir(t)}
	descriptor := (&Descriptor{}).initDescriptorSource(cfg)
	require.NotNil(t, descriptor.getFileSource())

	conn, err := grpc.NewClient(startTestGRPCServer(t), grpc.WithTransportCredentials(insecure.NewCredentials()))
	require.NoError(t, err)
	t.Cleanup(func() { require.NoError(t, conn.Close()) })

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	source, err := descriptor.getDescriptorSource(context.WithValue(ctx, ct.ContextKey(GrpcClientConnKey), conn), cfg)
	require.NoError(t, err)

	method, err := descriptor.getMethodDescriptor(source, conn, "test.Greeter", "Hello")
	require.NoError(t, err)
	require.Equal(t, "Hello", method.GetName())
}

// TestDescriptorLocalStrategyReportsUnloadableProtoFiles keeps a misconfigured `path` from
// producing a descriptor source that only fails once it is dereferenced.
func TestDescriptorLocalStrategyReportsUnloadableProtoFiles(t *testing.T) {
	cfg := &Config{DescriptorSourceStrategy: LOCAL, Path: filepath.Join(t.TempDir(), "missing")}
	descriptor := (&Descriptor{}).initDescriptorSource(cfg)

	source, err := descriptor.getDescriptorSource(context.Background(), cfg)
	require.Error(t, err)
	require.Nil(t, source)
}

// TestCompositeSourceFindSymbolWithoutFileSource keeps the fallback safe when the local
// proto files are unavailable altogether.
func TestCompositeSourceFindSymbolWithoutFileSource(t *testing.T) {
	cs := &compositeSource{}
	_, err := cs.FindSymbol("test.Greeter")
	require.Error(t, err)
}

// TestCompositeSourceAllExtensionsForTypeWithoutReflection asserts the extensions found in
// the local proto files are returned instead of being computed and then discarded.
func TestCompositeSourceAllExtensionsForTypeWithoutReflection(t *testing.T) {
	files, err := (protoparse.Parser{
		Accessor: protoparse.FileContentsFromMap(map[string]string{
			"extension.proto": `syntax = "proto2";
package test;

message Request {
  extensions 100 to 200;
}

extend Request {
  optional string note = 100;
}
`,
		}),
	}).ParseFiles("extension.proto")
	require.NoError(t, err)

	cs := &compositeSource{file: &fileSource{files: map[string]*desc.FileDescriptor{"extension.proto": files[0]}}}
	exts, err := cs.AllExtensionsForType("test.Request")
	require.NoError(t, err)
	require.Len(t, exts, 1)
	require.Equal(t, "test.note", exts[0].GetFullyQualifiedName())
}

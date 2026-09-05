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
	"os"
	"path/filepath"
	"strings"
	"sync"
)

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/grpcreflect"

	"github.com/pkg/errors"

	"google.golang.org/grpc"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
)

import (
	ct "github.com/apache/dubbo-go-pixiu/pkg/context"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

type Descriptor struct {
	fileSource  *fileSource
	methodMu    sync.RWMutex
	methodDescs map[*grpc.ClientConn]map[string]*desc.MethodDescriptor
	// connectionStates invalidates only lookups for the connection that
	// was removed. closeGeneration invalidates all in-flight lookups on close.
	connectionStates  map[*grpc.ClientConn]*descriptorConnectionState
	closeGeneration   uint64
	nextStateSequence uint64
	closed            bool
}

type descriptorConnectionState struct {
	generation uint64
}

type serviceNotExposedError struct {
	service string
}

func (e *serviceNotExposedError) Error() string {
	return fmt.Sprintf("service not exposed: %s", e.service)
}

func (dr *Descriptor) GetCurrentDescriptorSource(ctx context.Context) (DescriptorSource, error) {

	value := ctx.Value(ct.ContextKey(DescriptorSourceKey))

	switch t := value.(type) {
	case *DescriptorSource:
		return value.(DescriptorSource), nil
	case nil:
		return nil, errors.New("the descriptor source not found!")
	default:
		return nil, errors.Errorf("found a value of type %s, which is not DescriptorSource, ", t)
	}
}

func (dr *Descriptor) getDescriptorSource(ctx context.Context, cfg *Config) (DescriptorSource, error) {

	var ds DescriptorSource
	var err error

	switch strings.ToLower(cfg.DescriptorSourceStrategy.String()) {
	case LOCAL:
		// file only
		ds, err = dr.getFileDescriptorCompose(ctx, cfg)
	case REMOTE:
		// server reflection only
		ds, err = dr.getServerDescriptorSourceCtx(ctx, cfg)
	case AUTO:
		// file + reflection
		ds, err = dr.getDescriptorCompose(ctx, cfg)
	case NONE:
		// nope
		logger.Warnf("%s grpc descriptor source is none config , check descriptor_source_strategy %s ", loggerHeader, cfg.DescriptorSourceStrategy.String())
	default:
		err = errors.Errorf("grpc descriptor source not initialized cause the config of `descriptor_source_strategy` is %s, maybe set it `AUTO`", cfg.DescriptorSourceStrategy)
	}

	return ds, err
}

func (dr *Descriptor) getDescriptorCompose(ctx context.Context, cfg *Config) (DescriptorSource, error) {

	var err error

	cs := &compositeSource{}
	cs.reflection, err = dr.getServerDescriptorSourceCtx(ctx, cfg)
	// Never store a nil *fileSource: it would leave cs.file a non-nil interface
	// holding a nil pointer, so the fallback in compositeSource would dereference
	// it instead of reporting that the local proto files are unavailable.
	if fs := dr.getFileSource(); fs != nil {
		cs.file = fs
	}

	return cs, err
}

func (dr *Descriptor) initDescriptorSource(cfg *Config) *Descriptor {

	switch strings.ToLower(cfg.DescriptorSourceStrategy.String()) {
	case LOCAL, AUTO:
		// AUTO is `file + reflection`, so the local proto files have to be loaded
		// as well, otherwise the fallback of the reflection lookup has no source.
		dr.initFileDescriptorSource(cfg)
	}

	return dr
}

func (dr *Descriptor) getServerDescriptorSourceCtx(refCtx context.Context, cfg *Config) (DescriptorSource, error) {
	var (
		err error
		cc  *grpc.ClientConn
	)
	switch t := refCtx.Value(ct.ContextKey(GrpcClientConnKey)).(type) {
	case *grpc.ClientConn:
		cc = t
	case nil:
		err = errors.New("the descriptor source not found!")
	default:
		err = errors.Errorf("found a value of type %s, which is not *grpc.ClientConn, ", t)
	}
	if err != nil {
		return nil, err
	}

	// The reflection client is created per lookup and bound to the request
	// context so every remote reflection RPC honors the request timeout.
	// It must not be cached connection-scoped: grpcreflect reuses the root
	// context for every RPC, and a cached client would lose the deadline and
	// keep the per-request timeout from applying. The method descriptor
	// cache in getMethodDescriptor below is what avoids repeating the
	// reflection RPC after the first lookup.
	return &serverSource{client: grpcreflect.NewClientV1Alpha(refCtx, reflectpb.NewServerReflectionClient(cc))}, nil
}

// nolint
func (dr *Descriptor) getServerDescriptorSource(refCtx context.Context, cc *grpc.ClientConn) DescriptorSource {
	return &serverSource{client: grpcreflect.NewClientV1Alpha(refCtx, reflectpb.NewServerReflectionClient(cc))}
}

func (dr *Descriptor) removeConnection(cc *grpc.ClientConn) {
	if cc == nil {
		return
	}
	dr.methodMu.Lock()
	delete(dr.methodDescs, cc)
	state := dr.connectionStates[cc]
	if state != nil {
		dr.nextStateSequence++
		state.generation = dr.nextStateSequence
		// Keep the state alive for any in-flight lookup, but do not retain the
		// removed connection in the descriptor's long-lived map.
		delete(dr.connectionStates, cc)
	}
	dr.methodMu.Unlock()
}

func (dr *Descriptor) Close() {
	dr.methodMu.Lock()
	dr.closed = true
	dr.methodDescs = nil
	dr.connectionStates = nil
	dr.closeGeneration++
	dr.methodMu.Unlock()
}

func (dr *Descriptor) getMethodDescriptor(source DescriptorSource, cc *grpc.ClientConn, service, method string) (*desc.MethodDescriptor, error) {
	key := service + "\x00" + method
	dr.methodMu.Lock()
	if dr.closed {
		dr.methodMu.Unlock()
		return nil, errors.New("descriptor is closed")
	}
	if methods := dr.methodDescs[cc]; methods != nil {
		if descriptor, ok := methods[key]; ok {
			dr.methodMu.Unlock()
			return descriptor, nil
		}
	}
	state := dr.connectionStates[cc]
	if state == nil {
		if dr.connectionStates == nil {
			dr.connectionStates = make(map[*grpc.ClientConn]*descriptorConnectionState)
		}
		dr.nextStateSequence++
		state = &descriptorConnectionState{generation: dr.nextStateSequence}
		dr.connectionStates[cc] = state
	}
	connectionGeneration := state.generation
	closeGeneration := dr.closeGeneration
	dr.methodMu.Unlock()
	// Do not singleflight this lookup across requests. DescriptorSource carries
	// the request context used by server reflection, so sharing the first
	// caller's source would let its timeout or cancellation fail other callers.
	// The per-connection method cache below still removes repeated reflection
	// calls after the first successful lookup.
	dscp, err := source.FindSymbol(service)
	if err != nil {
		return nil, err
	}
	svcDesc, ok := dscp.(*desc.ServiceDescriptor)
	if !ok {
		return nil, &serviceNotExposedError{service: service}
	}
	descriptor := svcDesc.FindMethodByName(method)
	if descriptor == nil {
		return nil, fmt.Errorf("method not found: %s/%s", service, method)
	}

	dr.methodMu.Lock()
	defer dr.methodMu.Unlock()
	if dr.closed || dr.closeGeneration != closeGeneration || state.generation != connectionGeneration {
		return nil, errors.New("descriptor cache invalidated")
	}
	if methods := dr.methodDescs[cc]; methods != nil {
		if cached, ok := methods[key]; ok {
			return cached, nil
		}
	}
	if dr.methodDescs == nil {
		dr.methodDescs = make(map[*grpc.ClientConn]map[string]*desc.MethodDescriptor)
	}
	if dr.methodDescs[cc] == nil {
		dr.methodDescs[cc] = make(map[string]*desc.MethodDescriptor)
	}
	dr.methodDescs[cc][key] = descriptor
	return descriptor, nil
}

func (dr *Descriptor) getFileDescriptorCompose(ctx context.Context, cfg *Config) (DescriptorSource, error) {
	dr.initFileDescriptorSource(cfg)
	fs := dr.getFileSource()
	if fs == nil {
		return nil, errors.New("the local proto file descriptor source is not available")
	}
	return fs, nil
}

func (dr *Descriptor) initFileDescriptorSource(cfg *Config) *Descriptor {
	if dr.getFileSource() != nil {
		return dr
	}

	descriptor, err := loadFileSource(cfg)

	if err != nil {
		logger.Errorf("%s init gRPC descriptor by local file error: %v", loggerHeader, err)
		return dr
	}

	dr.methodMu.Lock()
	if dr.fileSource == nil {
		dr.fileSource = descriptor
	}
	dr.methodMu.Unlock()

	return dr
}

func (dr *Descriptor) getFileSource() *fileSource {
	dr.methodMu.RLock()
	defer dr.methodMu.RUnlock()
	return dr.fileSource
}

func loadFileSource(gc *Config) (*fileSource, error) {

	var fsrc fileSource

	cur := gc.Path
	if !filepath.IsAbs(cur) {
		ex, err := os.Executable()
		if err != nil {
			return nil, err
		}
		cur = filepath.Dir(ex) + string(os.PathSeparator) + gc.Path
	}

	logger.Infof("%s load proto files from %s", loggerHeader, cur)

	fileLists := make([]string, 0)
	items, err := os.ReadDir(cur)
	if err != nil {
		return nil, err
	}

	for _, item := range items {
		if !item.IsDir() {
			sp := strings.Split(item.Name(), ".")
			length := len(sp)
			if length >= 2 && sp[length-1] == "proto" {
				fileLists = append(fileLists, item.Name())
			}
		}
	}

	if err != nil {
		return nil, err
	}

	importPaths := []string{gc.Path}

	fileNames, err := protoparse.ResolveFilenames(importPaths, fileLists...)
	if err != nil {
		return nil, err
	}
	p := protoparse.Parser{
		ImportPaths:           importPaths,
		InferImportPaths:      len(importPaths) == 0,
		IncludeSourceCodeInfo: true,
	}
	fds, err := p.ParseFiles(fileNames...)
	if err != nil {
		return nil, fmt.Errorf("could not parse given files: %v", err)
	}

	fsrc.files = make(map[string]*desc.FileDescriptor)
	for _, fd := range fds {
		name := fd.GetName()
		fsrc.files[name] = fd
	}

	return &fsrc, nil
}

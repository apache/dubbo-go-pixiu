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

	"golang.org/x/sync/singleflight"

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
	methodLoads singleflight.Group
	generation  uint64
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
	cs.file = dr.fileSource

	return cs, err
}

func (dr *Descriptor) initDescriptorSource(cfg *Config) *Descriptor {

	if cfg.DescriptorSourceStrategy.String() == LOCAL {
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
	dr.generation++
	dr.methodMu.Unlock()
}

func (dr *Descriptor) Close() {
	dr.methodMu.Lock()
	dr.methodDescs = nil
	dr.generation++
	dr.methodMu.Unlock()
}

func (dr *Descriptor) getMethodDescriptor(source DescriptorSource, cc *grpc.ClientConn, service, method string) (*desc.MethodDescriptor, error) {
	key := service + "\x00" + method
	dr.methodMu.RLock()
	if methods := dr.methodDescs[cc]; methods != nil {
		if descriptor, ok := methods[key]; ok {
			dr.methodMu.RUnlock()
			return descriptor, nil
		}
	}
	dr.methodMu.RUnlock()

	dr.methodMu.RLock()
	generation := dr.generation
	dr.methodMu.RUnlock()
	loadKey := fmt.Sprintf("%p:%d:%s", cc, generation, key)
	result, err, _ := dr.methodLoads.Do(loadKey, func() (any, error) {
		dr.methodMu.RLock()
		if methods := dr.methodDescs[cc]; methods != nil {
			if descriptor, ok := methods[key]; ok {
				dr.methodMu.RUnlock()
				return descriptor, nil
			}
		}
		dr.methodMu.RUnlock()

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
		if dr.generation != generation {
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
	})
	if err != nil {
		return nil, err
	}
	return result.(*desc.MethodDescriptor), nil
}

func (dr *Descriptor) getFileDescriptorCompose(ctx context.Context, cfg *Config) (DescriptorSource, error) {
	if dr.fileSource == nil {
		dr.initFileDescriptorSource(cfg)
	}
	return dr.fileSource, nil
}

func (dr *Descriptor) initFileDescriptorSource(cfg *Config) *Descriptor {

	if dr.fileSource != nil {
		return dr
	}

	descriptor, err := loadFileSource(cfg)

	if err != nil {
		logger.Errorf("%s init gRPC descriptor by local file error: %v", loggerHeader, err)
		return dr
	}

	dr.fileSource = descriptor

	return dr
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

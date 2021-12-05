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
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/desc/protoparse"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/pkg/errors"
	"google.golang.org/grpc"
	reflectpb "google.golang.org/grpc/reflection/grpc_reflection_v1alpha"
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"
)

type Descriptor struct {

	//descSource DescriptorSource
	//reflSource *serverSource

	fileSource *fileSource
}

func (dr *Descriptor) GetCurrentDescriptor(ctx context.Context) (DescriptorSource, error) {

	value := ctx.Value(DescriptorSourceKey)

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

	switch cfg.DescriptorSourceStrategy {
	case LOCAL:
		// file only
		if dr.fileSource == nil {
			dr.initFileDescriptorSource(cfg)
		}
		ds = dr.fileSource
	case REMOTE:
		// server reflection only
		var cc *grpc.ClientConn
		gconn := ctx.Value(GrpcClientConnKey)
		switch t := gconn.(type) {
		case *grpc.ClientConn:
			cc = gconn.(*grpc.ClientConn)
		case nil:
			err = errors.New("the descriptor source not found!")
		default:
			err = errors.Errorf("found a value of type %s, which is not *grpc.ClientConn, ", t)
		}
		ds = dr.getServerDescriptorSource(ctx, cc)
	case AUTO:
		// file + reflection
		cs := &compositeSource{}
		cs.file = dr.fileSource
		cs.reflection = dr.getDescriptorCompose(ctx, nil)
	case NONE:
		// nope
		logger.Warnf("%s grpc descriptor source is no need initialize cause the strategy is %s ", loggerHeader, cfg.DescriptorSourceStrategy)
	default:
		err = errors.Errorf("grpc descriptor source not initialized cause the config of `descriptor_source_strategy` is %s, maybe set it `AUTO`", cfg.DescriptorSourceStrategy)
	}

	if err == nil {
		context.WithValue(ctx, DescriptorSourceKey, ds)
	}

	return ds, err
}

func (dr *Descriptor) getDescriptorCompose(ctx context.Context, cc *grpc.ClientConn) DescriptorSource {

	cs := &compositeSource{}
	cs.reflection = dr.getServerDescriptorSource(ctx, cc)
	cs.file = dr.fileSource

	return cs
}

func (dr *Descriptor) initDescriptorSource(cfg *Config) *Descriptor {

	if cfg.DescriptorSourceStrategy == LOCAL {
		dr.initFileDescriptorSource(cfg)
	}

	return dr
}

func (dr *Descriptor) getServerDescriptorSource(refCtx context.Context, cc *grpc.ClientConn) DescriptorSource {
	return &serverSource{client: grpcreflect.NewClient(refCtx, reflectpb.NewServerReflectionClient(cc))}
}

func (dr *Descriptor) initFileDescriptorSource(cfg *Config) *Descriptor {

	if dr.fileSource != nil {
		return dr
	}

	descriptor, err := loadFileSource(cfg)

	if err != nil {
		logger.Errorf("%s init gRPC descriptor by local file error, ", loggerHeader, err)
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
	items, err := ioutil.ReadDir(cur)
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

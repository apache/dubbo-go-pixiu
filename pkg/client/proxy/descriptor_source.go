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

// This file contains code that is copied and modified from:
// https://github.com/fullstorydev/grpcurl/blob/v1.8.7/desc_source.go

package proxy

import (
	"os"
	"sync"
)

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/pkg/errors"
	"google.golang.org/protobuf/proto"
	"google.golang.org/protobuf/types/descriptorpb"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/filter/http/grpcproxy"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

var (
	sourceOnce     sync.Once
	protosetSource grpcproxy.DescriptorSource
)

func InitProtosetSource(protoset []string) {
	sourceOnce.Do(func() {
		var err error
		protosetSource, err = DescriptorSourceFromProtoset(protoset)
		if err != nil {
			logger.Infof("[dubbo-go-pixiu] could not load protoset files: %v", err)
		}
	})
}

// DescriptorSourceFromProtoset creates a DescriptorSource that is backed by the named files,
// whose contents are Protocol Buffer source files. The given importPaths are used to locate
// any imported files.
func DescriptorSourceFromProtoset(filenames []string) (grpcproxy.DescriptorSource, error) {
	if len(filenames) < 1 {
		return nil, errors.New("no protoset files provided")
	}
	files := &descriptorpb.FileDescriptorSet{}
	for _, filename := range filenames {
		b, err := os.ReadFile(filename)
		if err != nil {
			return nil, errors.Errorf("wrong path to load protoset file %q: %v", filename, err)
		}
		var fs descriptorpb.FileDescriptorSet
		err = proto.Unmarshal(b, &fs)
		if err != nil {
			return nil, errors.Errorf("could not parse contents of protoset file %q: %v", filename, err)
		}
		files.File = append(files.File, fs.File...)
	}
	return DescriptorSourceFromFileDescriptorSet(files)
}

// DescriptorSourceFromFileDescriptorSet creates a DescriptorSource that is backed by the FileDescriptorSet.
func DescriptorSourceFromFileDescriptorSet(files *descriptorpb.FileDescriptorSet) (grpcproxy.DescriptorSource, error) {
	unresolved := map[string]*descriptorpb.FileDescriptorProto{}
	for _, fd := range files.File {
		unresolved[fd.GetName()] = fd
	}
	resolved := map[string]*desc.FileDescriptor{}
	for _, fd := range files.File {
		_, err := resolveFileDescriptor(unresolved, resolved, fd.GetName())
		if err != nil {
			return nil, err
		}
	}
	return &protosetFileSource{files: resolved}, nil
}

func resolveFileDescriptor(unresolved map[string]*descriptorpb.FileDescriptorProto, resolved map[string]*desc.FileDescriptor, filename string) (*desc.FileDescriptor, error) {
	if r, ok := resolved[filename]; ok {
		return r, nil
	}
	fd, ok := unresolved[filename]
	if !ok {
		return nil, errors.Errorf("no descriptor found for %q", filename)
	}
	deps := make([]*desc.FileDescriptor, 0, len(fd.GetDependency()))
	for _, dep := range fd.GetDependency() {
		depFd, err := resolveFileDescriptor(unresolved, resolved, dep)
		if err != nil {
			return nil, err
		}
		deps = append(deps, depFd)
	}
	result, err := desc.CreateFileDescriptor(fd, deps...)
	if err != nil {
		return nil, err
	}
	resolved[filename] = result
	return result, nil
}

type protosetFileSource struct {
	files  map[string]*desc.FileDescriptor
	er     *dynamic.ExtensionRegistry
	erInit sync.Once
}

func (fs *protosetFileSource) ListServices() ([]string, error) {
	set := map[string]bool{}
	for _, fd := range fs.files {
		for _, svc := range fd.GetServices() {
			set[svc.GetFullyQualifiedName()] = true
		}
	}
	sl := make([]string, 0, len(set))
	for svc := range set {
		sl = append(sl, svc)
	}
	return sl, nil
}

func (fs *protosetFileSource) FindSymbol(fullyQualifiedName string) (desc.Descriptor, error) {
	for _, fd := range fs.files {
		if dsc := fd.FindSymbol(fullyQualifiedName); dsc != nil {
			return dsc, nil
		}
	}
	return nil, errors.Errorf("Symbol not found: %s", fullyQualifiedName)
}

func (fs *protosetFileSource) AllExtensionsForType(typeName string) ([]*desc.FieldDescriptor, error) {
	fs.erInit.Do(func() {
		fs.er = &dynamic.ExtensionRegistry{}
		for _, fd := range fs.files {
			fs.er.AddExtensionsFromFile(fd)
		}
	})
	return fs.er.AllExtensionsForType(typeName), nil
}

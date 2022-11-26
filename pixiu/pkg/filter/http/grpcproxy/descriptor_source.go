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
	"fmt"
	"sync"
)

import (
	"github.com/jhump/protoreflect/desc"
	"github.com/jhump/protoreflect/dynamic"
	"github.com/jhump/protoreflect/grpcreflect"
	"github.com/pkg/errors"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

type DescriptorSource interface {
	// ListServices returns a list of fully-qualified service names. It will be all services in a set of
	// descriptor files or the set of all services exposed by a gRPC server.
	ListServices() ([]string, error)
	// FindSymbol returns a descriptor for the given fully-qualified symbol name.
	FindSymbol(fullyQualifiedName string) (desc.Descriptor, error)
	// AllExtensionsForType returns all known extension fields that extend the given message type name.
	AllExtensionsForType(typeName string) ([]*desc.FieldDescriptor, error)
}

var ErrReflectionNotSupported = errors.New("server does not support the reflection API")

// serverSource by gRPC server reflection
type serverSource struct {
	client *grpcreflect.Client
}

func (s *serverSource) ListServices() ([]string, error) {
	svcs, err := s.client.ListServices()
	return svcs, reflectionSupport(err)
}

func (s *serverSource) FindSymbol(fullyQualifiedName string) (desc.Descriptor, error) {
	file, err := s.client.FileContainingSymbol(fullyQualifiedName)
	if err != nil {
		return nil, reflectionSupport(err)
	}
	d := file.FindSymbol(fullyQualifiedName)
	if d == nil {
		return nil, errors.New(fmt.Sprintf("%s not found: %s", "Symbol", fullyQualifiedName))
	}
	return d, nil
}

func (s *serverSource) AllExtensionsForType(typeName string) ([]*desc.FieldDescriptor, error) {
	var exts []*desc.FieldDescriptor
	nums, err := s.client.AllExtensionNumbersForType(typeName)
	if err != nil {
		return nil, reflectionSupport(err)
	}
	for _, fieldNum := range nums {
		ext, err := s.client.ResolveExtension(typeName, fieldNum)
		if err != nil {
			return nil, reflectionSupport(err)
		}
		exts = append(exts, ext)
	}
	return exts, nil
}

func reflectionSupport(err error) error {
	if err == nil {
		return nil
	}
	if stat, ok := status.FromError(err); ok && stat.Code() == codes.Unimplemented {
		return ErrReflectionNotSupported
	}
	return err
}

// fileSource by gRPC proto file
type fileSource struct {
	files  map[string]*desc.FileDescriptor
	er     *dynamic.ExtensionRegistry
	erInit sync.Once
}

func (fs *fileSource) ListServices() ([]string, error) {
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

func (fs *fileSource) FindSymbol(fullyQualifiedName string) (desc.Descriptor, error) {
	for _, fd := range fs.files {
		if dsc := fd.FindSymbol(fullyQualifiedName); dsc != nil {
			return dsc, nil
		}
	}
	return nil, fmt.Errorf("could not found symbol %v", fullyQualifiedName)
}

func (fs *fileSource) AllExtensionsForType(typeName string) ([]*desc.FieldDescriptor, error) {
	fs.erInit.Do(func() {
		fs.er = &dynamic.ExtensionRegistry{}
		for _, fd := range fs.files {
			fs.er.AddExtensionsFromFile(fd)
		}
	})
	return fs.er.AllExtensionsForType(typeName), nil
}

// compositeSource by fileSource and serverSource
type compositeSource struct {
	reflection DescriptorSource
	file       DescriptorSource
}

func (cs *compositeSource) ListServices() ([]string, error) {
	if cs.reflection == nil {
		return nil, ErrReflectionNotSupported
	}
	return cs.reflection.ListServices()
}

func (cs *compositeSource) FindSymbol(fullyQualifiedName string) (desc.Descriptor, error) {
	if cs.reflection != nil {
		descriptor, err := cs.reflection.FindSymbol(fullyQualifiedName)
		if err == nil {
			logger.Debugf("%s find symbol by reflection : %v", loggerHeader, descriptor)
			return descriptor, nil
		}
	}

	return cs.file.FindSymbol(fullyQualifiedName)
}

func (cs *compositeSource) AllExtensionsForType(typeName string) ([]*desc.FieldDescriptor, error) {

	if cs.reflection == nil {
		fileExts, err := cs.file.AllExtensionsForType(typeName)
		if err != nil {
			return fileExts, nil
		}
	} else {
		exts, err := cs.reflection.AllExtensionsForType(typeName)
		if err != nil {
			return cs.file.AllExtensionsForType(typeName)
		}
		tags := make(map[int32]bool)
		for _, ext := range exts {
			tags[ext.GetNumber()] = true
		}

		fileExts, err := cs.file.AllExtensionsForType(typeName)
		if err != nil {
			return exts, nil
		}
		for _, ext := range fileExts {
			if !tags[ext.GetNumber()] {
				exts = append(exts, ext)
			}
		}
		return exts, nil
	}
	return nil, nil
}

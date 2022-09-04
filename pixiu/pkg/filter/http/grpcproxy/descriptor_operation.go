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
)

import (
	"github.com/jhump/protoreflect/desc"
)

func DescriptorSourceFromFileDescriptors(files ...*desc.FileDescriptor) (*fileSource, error) {
	fds := map[string]*desc.FileDescriptor{}
	for _, fd := range files {
		if err := addFile(fd, fds); err != nil {
			return nil, err
		}
	}
	return &fileSource{files: fds}, nil
}

func addFile(fd *desc.FileDescriptor, fds map[string]*desc.FileDescriptor) error {
	name := fd.GetName()
	if existing, ok := fds[name]; ok {
		// already added this file
		if existing != fd {
			// doh! duplicate files provided
			return fmt.Errorf("given files include multiple copies of %q", name)
		}
		return nil
	}
	fds[name] = fd
	for _, dep := range fd.GetDependencies() {
		if err := addFile(dep, fds); err != nil {
			return err
		}
	}
	return nil
}

func (fs *fileSource) GetAllFiles() ([]*desc.FileDescriptor, error) {
	files := make([]*desc.FileDescriptor, len(fs.files))
	i := 0
	for _, fd := range fs.files {
		files[i] = fd
		i++
	}
	return files, nil
}

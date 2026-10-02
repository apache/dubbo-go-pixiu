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

package tcp

import (
	"github.com/apache/dubbo-getty"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/filterchain"
)

type PackageHandler struct {
	ls *TcpListenerService
}

func NewPackageHandler(ls *TcpListenerService) *PackageHandler {
	return &PackageHandler{ls}
}

func (h *PackageHandler) Read(ss getty.Session, data []byte) (any, int, error) {
	var result any
	var length int
	err := h.ls.WithFilterChain(func(fc *filterchain.NetworkFilterChain) error {
		var decodeErr error
		result, length, decodeErr = fc.OnDecode(data)
		return decodeErr
	})
	return result, length, err
}

func (h *PackageHandler) Write(ss getty.Session, p any) ([]byte, error) {
	var result []byte
	err := h.ls.WithFilterChain(func(fc *filterchain.NetworkFilterChain) error {
		var encodeErr error
		result, encodeErr = fc.OnEncode(p)
		return encodeErr
	})
	return result, err
}

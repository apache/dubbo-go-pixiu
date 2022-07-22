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

package wasm

import (
	"fmt"
	"net/http"
)

import (
	"mosn.io/proxy-wasm-go-host/common"

	"mosn.io/proxy-wasm-go-host/proxywasm"
)

type importHandler struct {
	reqHeader common.HeaderMap
	proxywasm.DefaultImportsHandler
}

type myHeaderMap struct {
	realMap http.Header
}

func (im *importHandler) GetHttpRequestHeader() common.HeaderMap {
	return im.reqHeader
}

func (im *importHandler) Log(level proxywasm.LogLevel, msg string) proxywasm.WasmResult {
	fmt.Println(msg)
	return proxywasm.WasmResultOk
}

func (m *myHeaderMap) Get(key string) (string, bool) {
	return m.realMap.Get(key), true
}

func (m *myHeaderMap) Set(key, value string) {
	m.realMap.Set(key, value)
}

func (m *myHeaderMap) Add(key, value string) { panic("implemented") }

func (m *myHeaderMap) Del(key string) { panic("implemented") }

func (m *myHeaderMap) Range(f func(key string, value string) bool) {
	for k, _ := range m.realMap {
		v := m.realMap.Get(k)
		f(k, v)
	}
}

func (m *myHeaderMap) Clone() common.HeaderMap { panic("implemented") }

func (m *myHeaderMap) ByteSize() uint64 { panic("implemented") }

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
	"net/http"
)

import (
	"mosn.io/proxy-wasm-go-host/common"
	"mosn.io/proxy-wasm-go-host/proxywasm"
)

import (
	contexthttp "github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

type (
	headerWrapper struct {
		header http.Header
	}

	importHandler struct {
		ctx       *contexthttp.HttpContext
		reqHeader common.HeaderMap
		proxywasm.DefaultImportsHandler
	}
)

func (h *headerWrapper) Get(key string) (string, bool) {
	value := h.header.Get(key)
	if value == "" {
		return "", false
	}
	return value, true
}

func (h *headerWrapper) Set(key, value string) {
	h.header.Set(key, value)
}

func (h *headerWrapper) Add(key, value string) {
	h.header.Add(key, value)
}

func (h *headerWrapper) Del(key string) {
	h.header.Del(key)
}

func (h *headerWrapper) Range(f func(key string, value string) bool) {
	for k := range h.header {
		v := h.header.Get(k)
		f(k, v)
	}
}

func (h *headerWrapper) Clone() common.HeaderMap {
	copy := &headerWrapper{
		header: h.header.Clone(),
	}
	return copy
}

func (h *headerWrapper) ByteSize() uint64 {
	var size uint64

	for k := range h.header {
		v := h.header.Get(k)
		size += uint64(len(k) + len(v))
	}
	return size
}

func (im *importHandler) GetHttpRequestHeader() common.HeaderMap {
	return im.reqHeader
}

func (im *importHandler) Log(level proxywasm.LogLevel, msg string) proxywasm.WasmResult {
	switch level {
	case proxywasm.LogLevelDebug:
		logger.Debugf(msg)
	case proxywasm.LogLevelInfo:
		logger.Infof(msg)
	case proxywasm.LogLevelWarn:
		logger.Warnf(msg)
	case proxywasm.LogLevelError:
		logger.Errorf(msg)
	default:
		logger.Infof(msg)
	}
	return proxywasm.WasmResultOk
}

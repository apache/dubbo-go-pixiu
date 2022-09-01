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
	"math"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

import (
	"github.com/mitchellh/mapstructure"
	"mosn.io/proxy-wasm-go-host/common"
	"mosn.io/proxy-wasm-go-host/proxywasm"
	"mosn.io/proxy-wasm-go-host/wasmer"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

// WasmService manages a collection of WasmInstances of the same type.
type WasmService struct {
	rootContextID      int32
	contextIDGenerator int32
	name               string
	path               string
	mutex              sync.Mutex
	model              model.WasmService
	instancePool       sync.Pool
}

// WasmConfig contains config about every wasmFile.
type WasmConfig struct {
	Path string `yaml:"path" json:"path,omitempty"`
}

// ABIContextWrapper is a request wrapper which indicates request was handled in the wasmInstance.
type ABIContextWrapper struct {
	ContextID int32
	name      string
	Context   *proxywasm.ABIContext
}

func createWasmService(service model.WasmService) (*WasmService, error) {
	var cfg WasmConfig
	if err := mapstructure.Decode(service.Config, &cfg); err != nil {
		return nil, err
	}

	wasmService := &WasmService{
		contextIDGenerator: 0,
		name:               service.Name,
		path:               cfg.Path,
		mutex:              sync.Mutex{},
		model:              service,
	}
	wasmService.rootContextID = atomic.AddInt32(&wasmService.contextIDGenerator, 1)

	wasmService.instancePool = sync.Pool{
		New: func() interface{} {
			pwd, _ := os.Getwd()
			path := filepath.Join(pwd, cfg.Path)
			if _, err := os.Stat(path); err != nil {
				logger.Warnf("[dubbo-go-pixiu] wasmFile path:%v error: %v", path, err)
				return nil
			}
			instance := wasmer.NewWasmerInstanceFromFile(path)
			proxywasm.RegisterImports(instance)
			_ = instance.Start()
			return instance
		},
	}
	return wasmService, nil
}

func (ws *WasmService) getWasmInstance() common.WasmInstance {
	return ws.instancePool.Get().(common.WasmInstance)
}

func (ws *WasmService) putWasmInstance(instance common.WasmInstance) error {
	ws.instancePool.Put(instance)
	return nil
}

func (ws *WasmService) createABIContextWrapper(ctx *http.HttpContext) *ABIContextWrapper {

	contextWrapper := new(ABIContextWrapper)
	contextWrapper.ContextID = atomic.AddInt32(&ws.contextIDGenerator, 1)
	atomic.CompareAndSwapInt32(&ws.contextIDGenerator, math.MaxInt32, 1)

	contextWrapper.name = ws.name
	contextWrapper.Context = &proxywasm.ABIContext{
		Imports: &importHandler{
			ctx: ctx,
			reqHeader: &headerWrapper{
				header: ctx.Request.Header,
			},
		},
		Instance: ws.getWasmInstance(),
	}
	_ = contextWrapper.Context.GetExports().ProxyOnContextCreate(ws.rootContextID, 0)

	return contextWrapper
}

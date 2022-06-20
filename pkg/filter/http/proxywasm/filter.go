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

package proxywasm

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"mosn.io/proxy-wasm-go-host/proxywasm/common"
	proxywasm "mosn.io/proxy-wasm-go-host/proxywasm/v1"
	"mosn.io/proxy-wasm-go-host/wasmer"
	"os"
	"path/filepath"
	"sync"
	"sync/atomic"
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type importHandler struct {
	reqHeader common.HeaderMap
	proxywasm.DefaultImportsHandler
}

type (
	// Plugin manages wasm instances by pool.
	Plugin struct {
		instanceMap sync.Pool
	}

	WasmFilterFactory struct {
		plugin             *Plugin
		cfg                *Config
		contextIDGenerator int32
		rootContextID      int32
	}

	WasmFilter struct {
		instance      common.WasmInstance
		ctx           *proxywasm.ABIContext
		contextID     int32
		rootContextID int32
	}

	Config struct {
		Path string `yaml:"path" json:"path,omitempty"`
	}
)

func (w *WasmFilter) Encode(ctx *http.HttpContext) filter.FilterStatus {
	_ = w.ctx.GetExports().ProxyOnContextCreate(w.contextID, w.rootContextID)

	_, _ = w.ctx.GetExports().ProxyOnRequestHeaders(w.contextID, int32(len(ctx.Request.Header)), 1)

	_ = w.ctx.GetExports().ProxyOnDelete(w.contextID)

	return filter.Continue
}

func (factory *WasmFilterFactory) Config() interface{} {
	return factory.cfg
}

func (factory *WasmFilterFactory) Apply() error {
	factory.rootContextID = atomic.AddInt32(&factory.contextIDGenerator, 1)

	factory.plugin.instanceMap = sync.Pool{
		New: func() interface{} {
			pwd, _ := os.Getwd()
			instance := wasmer.NewWasmerInstanceFromFile(filepath.Join(pwd, factory.cfg.Path))
			proxywasm.RegisterImports(instance)
			// Call Start() here is OK?
			_ = instance.Start()
			return instance
		},
	}
	return nil
}

func (factory *WasmFilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
	filter := &WasmFilter{
		rootContextID: factory.rootContextID,
	}
	filter.contextID = atomic.AddInt32(&factory.contextIDGenerator, 1)
	filter.instance = factory.getWasmInstance()
	filter.ctx = &proxywasm.ABIContext{
		Imports:  &importHandler{reqHeader: &myHeaderMap{ctx.Request.Header}},
		Instance: filter.instance,
	}
	chain.AppendEncodeFilters(filter)
	var once sync.Once
	once.Do(func() {
		_ = filter.ctx.GetExports().ProxyOnContextCreate(filter.rootContextID, 0)
	})
	return nil
}

func (factory *WasmFilterFactory) getWasmInstance() common.WasmInstance {
	return factory.plugin.instanceMap.Get().(common.WasmInstance)
}

func (p *Plugin) Kind() string {
	return constant.HTTPWasmFilter
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &WasmFilterFactory{
		cfg:                &Config{},
		contextIDGenerator: 0,
		plugin:             p,
	}, nil
}

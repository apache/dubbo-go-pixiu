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
	"sync"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

type WasmManager struct {
	model      *model.WasmConfig
	serviceMap map[string]*WasmService
}

var manager *WasmManager
var once sync.Once

// InitWasmManager loads config in conf.yaml
func InitWasmManager(model *model.WasmConfig) {
	manager = &WasmManager{
		model:      model,
		serviceMap: make(map[string]*WasmService),
	}

	for _, service := range model.Services {
		if !contains(service.Name) {
			logger.Warnf("[dubbo-go-pixiu] wasm config error: service name isn't in pixiu keys.")
			continue
		}
		wasmService, err := createWasmService(service)
		if err != nil {
			logger.Warnf("[dubbo-go-pixiu] wasm config error: %v", err)
			continue
		}
		manager.serviceMap[service.Name] = wasmService
	}
}

// GetWasmServiceByName get WasmService By Name
func getWasmServiceByName(name string) *WasmService {
	once.Do(func() {
		bs := config.GetBootstrap()
		InitWasmManager(bs.Wasm)
	})
	if wasmService, exist := manager.serviceMap[name]; exist {
		return wasmService
	}
	logger.Warnf("[dubbo-go-pixiu] no wasmService named %v, check pixiu conf.yaml", name)

	return nil
}

func GetServiceRootID(name string) int32 {
	wasmService := getWasmServiceByName(name)
	return wasmService.rootContextID
}

// CreateAbIContextByName create abiContext according to every request
func CreateABIContextByName(name string, ctx *http.HttpContext) *ABIContextWrapper {
	wasmService := getWasmServiceByName(name)

	if wasmService != nil {
		return wasmService.createABIContextWrapper(ctx)
	}
	return nil
}

// ContextDone put wasmInstance back to Pool.
func ContextDone(wrapper *ABIContextWrapper) error {
	wasmService := getWasmServiceByName(wrapper.name)
	return wasmService.putWasmInstance(wrapper.Context.Instance)
}

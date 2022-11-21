//go:build windows

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
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	Plugin struct {
	}

	WasmFilterFactory struct {
		cfg *Config
	}

	WasmFilter struct {
	}

	Config struct {
	}
)

func (w *WasmFilter) Decode(ctx *http.HttpContext) filter.FilterStatus {
	return filter.Continue
}

func (w *WasmFilter) Encode(ctx *http.HttpContext) filter.FilterStatus {
	return filter.Continue
}

func (factory *WasmFilterFactory) Config() interface{} {
	return factory.cfg
}

func (factory *WasmFilterFactory) Apply() error {
	return nil
}

func (factory *WasmFilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
	logger.Warnf("Wasm not support on windows now")
	return nil
}

func (p *Plugin) Kind() string {
	return constant.HTTPWasmFilter
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &WasmFilterFactory{
		cfg: &Config{},
	}, nil
}

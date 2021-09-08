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

package cors

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.HTTPCorsFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is http filter plugin.
	Plugin struct {
	}
	// Filter is http filter instance
	Filter struct {
		cfg *Config
	}
	// Config describe the config of Filter
	Config struct{}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter() (filter.HttpFilter, error) {
	return &Filter{cfg: &Config{}}, nil
}

func (f *Filter) PrepareFilterChain(ctx *http.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *Filter) Handle(ctx *http.HttpContext) {
	f.handleCors(ctx)

	ctx.Next()
}

func (f *Filter) handleCors(ctx *http.HttpContext) {
	cp := ctx.HttpConnectionManager.CorsPolicy
	if cp == nil || !cp.Enabled {
		return
	}

	domains := cp.AllowOrigin
	if len(domains) != 0 {
		for _, domain := range domains {
			if ctx.GetHeader("Host") == domain {
				ctx.AddHeader(constant.HeaderKeyAccessControlAllowOrigin, domain)
			}
		}
	}

	if cp.AllowHeaders != "" {
		ctx.AddHeader(constant.HeaderKeyAccessControlExposeHeaders, cp.AllowHeaders)
	}

	if cp.AllowMethods != "" {
		ctx.AddHeader(constant.HeaderKeyAccessControlAllowMethods, cp.AllowMethods)
	}

	if cp.MaxAge != "" {
		ctx.AddHeader(constant.HeaderKeyAccessControlMaxAge, cp.MaxAge)
	}

	if cp.AllowCredentials {
		ctx.AddHeader(constant.HeaderKeyAccessControlAllowCredentials, "true")
	}
}

func (f *Filter) Apply() error {
	return nil
}

func (f *Filter) Config() interface{} {
	return f.cfg
}

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
	stdHttp "net/http"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
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

	// FilterFactory is http filter instance
	FilterFactory struct {
		cfg *Config
	}
	Filter struct {
		cfg *Config
	}

	// Config describe the config of FilterFactory
	Config struct {
		AllowOrigin []string `yaml:"allow_origin" json:"allow_origin" mapstructure:"allow_origin"`
		// AllowMethods access-control-allow-methods
		AllowMethods string `yaml:"allow_methods" json:"allow_methods" mapstructure:"allow_methods"`
		// AllowHeaders access-control-allow-headers
		AllowHeaders string `yaml:"allow_headers" json:"allow_headers" mapstructure:"allow_headers"`
		// ExposeHeaders access-control-expose-headers
		ExposeHeaders string `yaml:"expose_headers" json:"expose_headers" mapstructure:"expose_headers"`
		// MaxAge access-control-max-age
		MaxAge           string `yaml:"max_age" json:"max_age" mapstructure:"max_age"`
		AllowCredentials bool   `yaml:"allow_credentials" json:"allow_credentials" mapstructure:"allow_credentials"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
	f := &Filter{cfg: factory.cfg}
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(ctx *http.HttpContext) filter.FilterStatus {
	f.handleCors(ctx)
	return filter.Continue
}

func (f *Filter) handleCors(ctx *http.HttpContext) {
	c := f.cfg
	if c == nil {
		return
	}

	domains := c.AllowOrigin
	if len(domains) != 0 {
		for _, domain := range domains {
			if ctx.Request.Host == domain || ctx.Request.URL.Host == domain ||
				ctx.GetHeader("Host") == domain || ctx.GetHeader("host") == domain {
				ctx.SourceResp.(*stdHttp.Response).Header.Add(constant.HeaderKeyAccessControlAllowOrigin, domain)
			}
		}
	}

	if c.AllowHeaders != "" {
		ctx.SourceResp.(*stdHttp.Response).Header.Add(constant.HeaderKeyAccessControlExposeHeaders, c.AllowHeaders)
	}

	if c.AllowMethods != "" {
		ctx.SourceResp.(*stdHttp.Response).Header.Add(constant.HeaderKeyAccessControlAllowMethods, c.AllowMethods)
	}

	if c.MaxAge != "" {
		ctx.SourceResp.(*stdHttp.Response).Header.Add(constant.HeaderKeyAccessControlMaxAge, c.MaxAge)
	}

	if c.AllowCredentials {
		ctx.SourceResp.(*stdHttp.Response).Header.Add(constant.HeaderKeyAccessControlAllowCredentials, "true")
	}
}

func (factory *FilterFactory) Apply() error {
	return nil
}

func (factory *FilterFactory) Config() interface{} {
	return factory.cfg
}

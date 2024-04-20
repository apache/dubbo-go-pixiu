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
	writer := ctx.Writer
	c := f.cfg
	if c == nil {
		return filter.Continue
	}
	if ctx.GetHeader(constant.Origin) == "" {
		// not a cors request
		return filter.Continue
	}

	domains := c.AllowOrigin
	if len(domains) != 0 {
		for _, domain := range domains {
			if domain == "*" || ctx.GetHeader("Origin") == domain {
				writer.Header().Add(constant.HeaderKeyAccessControlAllowOrigin, domain)
				continue
			}
		}
	}

	if c.AllowHeaders != "" {
		writer.Header().Add(constant.HeaderKeyAccessControlAllowHeaders, c.AllowHeaders)
	}

	if c.ExposeHeaders != "" {
		writer.Header().Add(constant.HeaderKeyAccessControlExposeHeaders, c.ExposeHeaders)
	}

	if c.AllowMethods != "" {
		writer.Header().Add(constant.HeaderKeyAccessControlAllowMethods, c.AllowMethods)
	}

	if c.MaxAge != "" {
		writer.Header().Add(constant.HeaderKeyAccessControlMaxAge, c.MaxAge)
	}

	if c.AllowCredentials {
		writer.Header().Add(constant.HeaderKeyAccessControlAllowCredentials, "true")
	}
	if ctx.Request.Method == stdHttp.MethodOptions {
		ctx.SendLocalReply(stdHttp.StatusOK, nil)
		return filter.Stop
	}
	return filter.Continue
}

func (factory *FilterFactory) Apply() error {
	return nil
}

func (factory *FilterFactory) Config() interface{} {
	return factory.cfg
}

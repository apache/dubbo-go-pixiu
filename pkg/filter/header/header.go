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

package header

import (
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

const (
	Kind = constant.HTTPHeaderFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is http filter plugin.
	Plugin struct {
	}
	// HeaderFilter is http filter instance
	HeaderFilter struct {
		cfg *Config
	}
	// Config describe the config of HeaderFilter
	Config struct{}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter() (filter.HttpFilter, error) {
	return &HeaderFilter{}, nil
}

func (f *HeaderFilter) Config() interface{} {
	return f.cfg
}

func (f *HeaderFilter) Apply() error {
	return nil
}

func (hf *HeaderFilter) PrepareFilterChain(ctx *http.HttpContext) error {
	ctx.AppendFilterFunc(hf.Handle)
	return nil
}

func (hf *HeaderFilter) Handle(hc *http.HttpContext) {
	api := hc.GetAPI()
	headers := api.Headers
	if len(headers) <= 0 {
		hc.Next()
		return
	}

	urlHeaders := hc.AllHeaders()
	if len(urlHeaders) <= 0 {
		hc.Abort()
		return
	}

	for headerName, headerValue := range headers {
		urlHeaderValues := urlHeaders.Values(strings.ToLower(headerName))
		if urlHeaderValues == nil {
			hc.Abort()
			return
		}
		for _, urlHeaderValue := range urlHeaderValues {
			if urlHeaderValue == headerValue {
				goto FOUND
			}
		}
		hc.Abort()
	FOUND:
		continue
	}
}

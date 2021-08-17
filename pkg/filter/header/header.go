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
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	http2 "github.com/apache/dubbo-go-pixiu/pkg/common/http"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"strings"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.HTTPHeaderFilter
)

func init() {
	extension.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is http filter plugin.
	Plugin struct {
	}
	// AccessFilter is http filter instance
	HeaderFilter struct {
		cfg *Config
		alw *model.AccessLogWriter
		alc *model.AccessLogConfig
	}
	// Config describe the config of AccessFilter
	Config struct{}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter(hcm *http2.HttpConnectionManager, config interface{}, bs *model.Bootstrap) (extension.HttpFilter, error) {
	alc := bs.StaticResources.AccessLogConfig
	if !alc.Enable {
		return nil, errors.Errorf("AccessPlugin CreateFilter error the access_log config not enable")
	}

	accessLogWriter := &model.AccessLogWriter{AccessLogDataChan: make(chan model.AccessLogData, constant.LogDataBuffer)}
	specConfig := config.(AuthorityConfiguration)
	return &Filter{cfg: &specConfig, alw: accessLogWriter, alc: &alc}, nil
}

type headerFilter struct{}

// nolint.
func New() *headerFilter {
	return &headerFilter{}
}

func (h *headerFilter) Do() context.FilterFunc {
	return func(c context.Context) {
		api := c.GetAPI()
		headers := api.Headers
		if len(headers) <= 0 {
			c.Next()
			return
		}
		switch c.(type) {
		case *http.HttpContext:
			hc := c.(*http.HttpContext)
			urlHeaders := hc.AllHeaders()
			if len(urlHeaders) <= 0 {
				c.Abort()
				return
			}

			for headerName, headerValue := range headers {
				urlHeaderValues := urlHeaders.Values(strings.ToLower(headerName))
				if urlHeaderValues == nil {
					c.Abort()
					return
				}
				for _, urlHeaderValue := range urlHeaderValues {
					if urlHeaderValue == headerValue {
						goto FOUND
					}
				}
				c.Abort()
			FOUND:
				continue
			}
			break
		}
	}
}

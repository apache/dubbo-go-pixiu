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

package proxyrewrite

import (
	"regexp"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	Kind = constant.HTTPProxyRewriteFilter
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
		cfg      *Config
		uriRegex *regexp.Regexp
	}
	Config struct {
		UriRegex []string          `yaml:"uri_regex" json:"uri_regex"`
		Headers  map[string]string `yaml:"headers" json:"headers"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter() (filter.HttpFilter, error) {
	return &Filter{cfg: &Config{}}, nil
}

func (f *Filter) Config() interface{} {
	return f.cfg
}

func (f *Filter) Apply() error {
	if len(f.cfg.UriRegex) == 2 {
		f.uriRegex = regexp.MustCompile(f.cfg.UriRegex[0])
	}
	return nil
}

func (f *Filter) PrepareFilterChain(ctx *contexthttp.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *Filter) Handle(c *contexthttp.HttpContext) {
	url := c.GetUrl()
	if f.uriRegex != nil {
		newUrl := f.uriRegex.ReplaceAllString(url, f.cfg.UriRegex[1])
		logger.Infof("proxy rewrite filter change url from %s to %s", url, newUrl)
		c.SetUrl(newUrl)
	}
	if len(f.cfg.Headers) > 0 {
		for k, v := range f.cfg.Headers {
			logger.Infof("proxy rewrite filter add header  key %s and value %s", k, v)
			c.AddHeader(k, v)
		}
	}
	c.Next()
}

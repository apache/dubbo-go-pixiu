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
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
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

	// FilterFactory is http filter instance
	FilterFactory struct {
		cfg *Config
	}
	//Filter
	Filter struct {
		Headers map[string]string

		uriRegex *regexp.Regexp
		replace  string
	}
	//Config
	Config struct {
		UriRegex []string          `yaml:"uri_regex" json:"uri_regex"`
		Headers  map[string]string `yaml:"headers" json:"headers"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

func (factory *FilterFactory) Config() interface{} {
	return factory.cfg
}

func (factory *FilterFactory) Apply() error {
	if len(factory.cfg.UriRegex) != 2 {
		return errors.Errorf("UriRegex len must == 2, %v", factory.cfg.UriRegex)
	}
	return nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	cfg := factory.cfg

	headers := make(map[string]string, len(cfg.Headers))
	for k, v := range cfg.Headers {
		headers[k] = v
	}
	f := &Filter{uriRegex: regexp.MustCompile(cfg.UriRegex[0]), replace: cfg.UriRegex[1], Headers: headers}

	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(c *contexthttp.HttpContext) filter.FilterStatus {
	url := c.GetUrl()

	newUrl := f.uriRegex.ReplaceAllString(url, f.replace)
	logger.Infof("proxy rewrite filter change url from %s to %s", url, newUrl)
	c.SetUrl(newUrl)

	if len(f.Headers) > 0 {
		for k, v := range f.Headers {
			logger.Infof("proxy rewrite filter add header  key %s and value %s", k, v)
			c.AddHeader(k, v)
		}
	}
	return filter.Continue
}

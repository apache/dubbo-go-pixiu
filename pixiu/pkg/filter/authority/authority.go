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

package authority

import (
	nh "net/http"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
)

const (
	Kind = constant.HTTPAuthorityFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// AuthorityPlugin is http filter plugin.
	Plugin struct {
	}
	// FilterFactory is http filter instance
	FilterFactory struct {
		cfg *AuthorityConfiguration
	}
	Filter struct {
		cfg *AuthorityConfiguration
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &AuthorityConfiguration{}}, nil
}

func (factory *FilterFactory) Config() interface{} {
	return factory.cfg
}

func (factory *FilterFactory) Apply() error {
	return nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
	f := &Filter{cfg: factory.cfg}
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(c *http.HttpContext) filter.FilterStatus {
	for _, r := range f.cfg.Rules {
		item := c.GetClientIP()
		if r.Limit == App {
			item = c.GetApplicationName()
		}

		result := passCheck(item, r)
		if !result {
			c.SendLocalReply(nh.StatusForbidden, constant.Default403Body)
			return filter.Stop
		}
	}
	return filter.Continue
}

func passCheck(item string, rule AuthorityRule) bool {
	result := false
	for _, it := range rule.Items {
		if it == item {
			result = true
			break
		}
	}

	if (rule.Strategy == Blacklist && result) || (rule.Strategy == Whitelist && !result) {
		return false
	}

	return true
}

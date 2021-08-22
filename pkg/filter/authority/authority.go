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
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.HTTPAuthorityFilter
)

func init() {
	extension.RegisterHttpFilter(&Plugin{})
}

type (
	// AuthorityPlugin is http filter plugin.
	Plugin struct {
	}
	// Filter is http filter instance
	Filter struct {
		cfg *AuthorityConfiguration
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter() (extension.HttpFilter, error) {
	return &Filter{cfg: &AuthorityConfiguration{}}, nil
}

func (f *Filter) Config() interface{} {
	return f.cfg
}

func (f *Filter) Apply() error {
	return nil
}

func (f *Filter) PrepareFilterChain(ctx *http.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *Filter) Handle(c *http.HttpContext) {
	for _, r := range f.cfg.Rules {
		item := c.GetClientIP()
		if r.Limit == App {
			item = c.GetApplicationName()
		}

		result := passCheck(item, r)
		if !result {
			c.WriteWithStatus(nh.StatusForbidden, constant.Default403Body)
			c.Abort()
			return
		}
	}

	c.Next()
}


func passCheck(item string, rule AuthorityRule) bool {
	result := false
	for _, it := range rule.Items {
		if it == item {
			result = true
			break
		}
	}

	if (rule.Strategy == Blacklist && result == true) || (rule.Strategy == Whitelist && result == false) {
		return false
	}

	return true
}

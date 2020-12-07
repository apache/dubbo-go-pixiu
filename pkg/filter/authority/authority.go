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
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
)

func Init() {
	extension.SetFilterFunc(constant.HTTPAuthorityFilter, authorityFilterFunc())
}

func authorityFilterFunc() context.FilterFunc {
	return New().Do()
}

// authorityFilter is a filter for blacklist/whitelist.
type authorityFilter struct {
}

// New create blacklist/whitelist filter.
func New() filter.Filter {
	return &authorityFilter{}
}

// Do execute blacklist/whitelist filter logic.
func (f authorityFilter) Do() context.FilterFunc {
	return func(c context.Context) {
		f.doAuthorityFilter(c.(*http.HttpContext))
	}
}

func (f authorityFilter) doAuthorityFilter(c *http.HttpContext) {
	for _, r := range c.HttpConnectionManager.AuthorityConfig.Rules {
		item := c.GetClientIP()
		if r.Limit == model.App {
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

func passCheck(item string, rule model.AuthorityRule) bool {
	result := false
	for _, it := range rule.Items {
		if it == item {
			result = true
			break
		}
	}

	if (rule.Strategy == model.Blacklist && result == true) || (rule.Strategy == model.Whitelist && result == false) {
		return false
	}

	return true
}

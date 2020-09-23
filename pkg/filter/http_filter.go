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

package filter

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
)

func init() {
	extension.SetFilterFunc(constant.HttpDomainFilter, Domain())
}

// Domain
// https :authority
// http Host
func Domain() context.FilterFunc {
	return func(c context.Context) {
		if MatchDomainFilter(c.(*http.HttpContext)) {
			c.Next()
		}
	}
}

// MatchDomainFilter
func MatchDomainFilter(c *http.HttpContext) bool {
	for _, v := range c.Listener.FilterChains {
		for _, d := range v.FilterChainMatch.Domains {
			if d == c.GetHeader(constant.Http1HeaderKeyHost) {
				return true
			}
		}
	}

	return false
}

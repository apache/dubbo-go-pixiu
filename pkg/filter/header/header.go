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
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/context"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
)

type headerFilter struct {
}

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

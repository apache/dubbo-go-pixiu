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

package host

import (
	selfcontext "github.com/dubbogo/dubbo-go-proxy/pkg/context"
	contexthttp "github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter"
)

// hostFilter is a filter for host.
type hostFilter struct {
	host string
}

// New create host filter.
func New(host string) filter.Filter {
	return &hostFilter{host: host}
}

// // Do execute hostFilter filter logic.
func (f *hostFilter) Do() selfcontext.FilterFunc {
	return func(c selfcontext.Context) {
		f.doHostFilter(c.(*contexthttp.HttpContext))
	}
}

func (f *hostFilter) doHostFilter(c *contexthttp.HttpContext) {
	c.Request.Host = f.host
	c.Next()
}

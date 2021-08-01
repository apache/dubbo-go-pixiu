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
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/context"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/filter"
)

import (
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	manager "github.com/apache/dubbo-go-pixiu/pkg/filter"
)

func Init() {
	manager.RegisterFilterFactory(constant.HostFilter, newHostFilter)
}

// hostFilter is a filter for host.
type hostFilter struct {
}

// newHostFilter create host filter.
func newHostFilter() filter.Factory {
	return &hostFilter{}
}

func (f *hostFilter) Config() interface{} {
	return nil
}

func (f *hostFilter) Apply() (filter.Filter, error) {
	return func(c fc.Context) {
		f.doHostFilter(c.(*contexthttp.HttpContext))
	}, nil
}

func (f hostFilter) doHostFilter(c *contexthttp.HttpContext) {
	c.Request.Host = c.GetAPI().Method.IntegrationRequest.Host
	c.Next()
}

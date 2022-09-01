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

package context

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

type (
	ContextKey string

	// Context run context
	Context interface {
		Next()
		Abort()
		AbortWithError(string, error)

		Status(code int)
		StatusCode() int
		WriteWithStatus(int, []byte) (int, error)
		Write([]byte) (int, error)
		AddHeader(k, v string)
		GetHeader(k string) string
		GetUrl() string
		GetMethod() string

		BuildFilters()

		GetRouteEntry() *model.RouteAction
		SetRouteEntry(ra *model.RouteAction)
		GetClientIP() string
		GetApplicationName() string

		WriteErr(p interface{})

		Request()
		Response()
	}
)

func (c ContextKey) String() string {
	return string(c)
}

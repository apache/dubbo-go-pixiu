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

package extension

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
)

var (
	filterFuncCacheMap = make(map[string]func(ctx context.Context), 4)
)

// SetFilterFunc will store the @filter and @name
func SetFilterFunc(name string, filter context.FilterFunc) {
	filterFuncCacheMap[name] = filter
}

// GetMustFilterFunc will return the proxy.FilterFunc
// if not found, it will panic
func GetMustFilterFunc(name string) context.FilterFunc {
	if filter, ok := filterFuncCacheMap[name]; ok {
		return filter
	}

	panic("filter func for " + name + " is not existing!")
}

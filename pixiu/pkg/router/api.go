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

package router

import (
	"net/url"
	"strings"
)

import (
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/router"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
)

// GetURIParams returns the values retrieved from the rawURL
func GetURIParams(api *router.API, rawURL url.URL) url.Values {
	return wildcardMatch(api.URLPattern, rawURL.Path)
}

// IsWildCardBackendPath checks whether the configured path of
// the upstream restful service contains parameters
func IsWildCardBackendPath(api *router.API) bool {
	if len(api.IntegrationRequest.Path) == 0 {
		return false
	}
	return strings.Contains(api.IntegrationRequest.Path, constant.PathParamIdentifier)
}

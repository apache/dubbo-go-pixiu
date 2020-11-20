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

package client

import (
	"context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"net/http"
	"strings"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/router"
)

// Request request for endpoint
type Request struct {
	Context context.Context

	IngressRequest *http.Request
	API            router.API
}

// NewReq create a request
func NewReq(ctx context.Context, request *http.Request, api router.API) *Request {
	return &Request{
		Context:        ctx,
		IngressRequest: request,
		API:            api,
	}
}

// GetURL new url
func (r *Request) GetURL() string {
	ir := r.API.IntegrationRequest
	if ir.RequestType == config.HTTPRequest {
		if len(ir.HTTPBackendConfig.URL) != 0 {
			return ir.HTTPBackendConfig.URL
		}

		// now only support http.
		scheme := "http"
		host := r.IngressRequest.URL.Host
		if len(ir.HTTPBackendConfig.Host) != 0 {
			host = ir.HTTPBackendConfig.Host
		}
		path := r.IngressRequest.URL.Path
		if len(ir.HTTPBackendConfig.Path) != 0 {
			path = ir.HTTPBackendConfig.Path
		}
		if !strings.HasPrefix(path, "/") {
			path = "/" + path
		}

		return scheme + "://" + host + path
	}

	return ""
}

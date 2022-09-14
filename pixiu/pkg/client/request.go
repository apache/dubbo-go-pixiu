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
	"net/http"
	"time"
)

import (
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/router"
)

// Request request for endpoint
type Request struct {
	Context        context.Context
	IngressRequest *http.Request
	API            router.API
	Timeout        time.Duration
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

		return scheme + "://" + r.IngressRequest.Host + r.IngressRequest.URL.Path
	}

	return ""
}

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

package api_config

import (
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"
)

import (
	pc "github.com/apache/dubbo-go-pixiu/pkg/config"
)

// DiscoveryRequest a request for discovery
type DiscoveryRequest struct {
	Body []byte
}

// NewDiscoveryRequest return a DiscoveryRequest with body
func NewDiscoveryRequest(b []byte) *DiscoveryRequest {
	return &DiscoveryRequest{
		Body: b,
	}
}

// DiscoveryResponse a response for discovery
type DiscoveryResponse struct {
	Success bool
	Data    interface{}
}

// NewDiscoveryResponseWithSuccess return a DiscoveryResponse with success
func NewDiscoveryResponseWithSuccess(b bool) *DiscoveryResponse {
	return &DiscoveryResponse{
		Success: b,
	}
}

// NewDiscoveryResponse return a DiscoveryResponse with Data and success true
func NewDiscoveryResponse(d interface{}) *DiscoveryResponse {
	return &DiscoveryResponse{
		Success: true,
		Data:    d,
	}
}

var EmptyDiscoveryResponse = &DiscoveryResponse{}

// APIDiscoveryService api discovery service interface
type APIDiscoveryService interface {
	pc.APIConfigResourceListener
	AddAPI(router.API) error
	ClearAPI() error
	GetAPI(string, config.HTTPVerb) (router.API, error)
	RemoveAPIByPath(deleted config.Resource) error
	RemoveAPI(fullPath string, method config.Method) error
}

// DiscoveryService is come from envoy, it can used for admin

// ListenerDiscoveryService
type ListenerDiscoveryService interface {
	AddListeners(request DiscoveryRequest) (DiscoveryResponse, error)
	GetListeners(request DiscoveryRequest) (DiscoveryResponse, error)
}

// RouteDiscoveryService
type RouteDiscoveryService interface {
	AddRoutes(r DiscoveryRequest) (DiscoveryResponse, error)
	GetRoutes(r DiscoveryRequest) (DiscoveryResponse, error)
}

// ClusterDiscoveryService
type ClusterDiscoveryService interface {
	AddClusters(r DiscoveryRequest) (DiscoveryResponse, error)
	GetClusters(r DiscoveryRequest) (DiscoveryResponse, error)
}

// EndpointDiscoveryService
type EndpointDiscoveryService interface {
	AddEndpoints(r DiscoveryRequest) (DiscoveryResponse, error)
	GetEndpoints(r DiscoveryRequest) (DiscoveryResponse, error)
}

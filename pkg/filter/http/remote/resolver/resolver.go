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

package resolver

import (
	"errors"
	"fmt"
	"net/http"
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/router"
)

// Resolver defines the interface for resolving an HTTP request to a specific API configuration.
// It allows for multiple resolution strategies to be implemented and tried in sequence.
type Resolver interface {
	// Resolve attempts to create an API configuration from the given HTTP context.
	// If the resolver can handle the request, it returns a configured *router.API and a nil error.
	// If the resolver is not applicable to this request, it should return (nil, error).
	// If an actual error occurs during processing, it should return (nil, error).
	Resolve(ctx *contexthttp.HttpContext) (*router.API, error)
}

// BaseResolver contains common logic for checking request prerequisites.
type BaseResolver struct{}

func (b *BaseResolver) PreCheck(req *http.Request) error {
	// 1. Method must be POST.
	// 2. Header must have x-dubbo-http1.1-dubbo-version.
	// 3. Path must be in {application}/{service}/{method} format.
	if req.Method != http.MethodPost || req.Header.Get(constant.DubboHttpDubboVersion) == "" {
		return errors.New("http request must be POST and have x-dubbo-http1.1-dubbo-version header")
	}

	rawPath := strings.Trim(req.URL.Path, "/")
	if len(strings.Split(rawPath, "/")) != 3 {
		return errors.New("http request path must be in {application}/{service}/{method} format")
	}

	return nil
}

func (b *BaseResolver) BuildAPI(req *http.Request, mappingParams []config.MappingParam) (*router.API, error) {
	integrationRequest := config.IntegrationRequest{}
	resolveProtocol := req.Header.Get(constant.DubboServiceProtocol)
	switch resolveProtocol {
	case string(constant.HTTPRequest):
		integrationRequest.RequestType = resolveProtocol
	case string(constant.DubboRequest):
		integrationRequest.RequestType = resolveProtocol
	case "triple":
		integrationRequest.RequestType = resolveProtocol
	default:
		// If the protocol is specified but unknown, it's an error.
		if resolveProtocol != "" {
			return nil, errors.New("http request has unknown protocol in x-dubbo-service-protocol")
		}
		// Default to dubbo if not specified
		integrationRequest.RequestType = constant.DubboRequest
	}

	dubboBackendConfig := config.DubboBackendConfig{
		Version: req.Header.Get(constant.DubboServiceVersion),
		Group:   req.Header.Get(constant.DubboGroup),
	}
	integrationRequest.DubboBackendConfig = dubboBackendConfig
	integrationRequest.MappingParams = mappingParams

	method := config.Method{
		Enable:             true,
		HTTPVerb:           http.MethodPost,
		IntegrationRequest: integrationRequest,
		InboundRequest:     config.InboundRequest{RequestType: constant.HTTPRequest},
	}

	api := router.API{
		URLPattern: "/:application/:interface/:method",
		Method:     method,
	}

	return &api, nil
}

// resolverRegistry holds all available resolver factory.
var resolverRegistry = make(map[string]Resolver)

// RegisterResolver register resolver factory to registry.
// This function is called from init() functions in files that define a resolver.
func RegisterResolver(name string, r Resolver) {
	name = strings.ToLower(name)
	if _, exists := resolverRegistry[name]; exists {
		logger.Warnf("retry policy type '%s' is being overwritten", name)
	}
	resolverRegistry[name] = r
}

// GetResolver dynamically creates a resolver.
func GetResolver(name string) (Resolver, error) {
	r, exists := resolverRegistry[strings.ToLower(name)]
	if !exists {
		return nil, fmt.Errorf("unknown resolver type '%s'", name)
	}

	return r, nil
}

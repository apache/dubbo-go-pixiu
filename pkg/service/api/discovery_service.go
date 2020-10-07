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

package api

import (
	"encoding/json"
)

import (
	"github.com/goinggo/mapstructure"
	"github.com/pkg/errors"
)

import (
	"fmt"
	"github.com/dubbogo/dubbo-go-proxy/pkg/client/dubbo"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
	"github.com/dubbogo/dubbo-go-proxy/pkg/router"
	"github.com/dubbogo/dubbo-go-proxy/pkg/service"
	"strings"
)

func init() {
	extension.SetApiDiscoveryService(constant.LocalMemoryApiDiscoveryService, NewLocalMemoryAPIDiscoveryService())
}

// LocalMemoryAPIDiscoveryService is the local cached API discovery service
type LocalMemoryAPIDiscoveryService struct {
	router *router.Route
}

// NewLocalMemoryAPIDiscoveryService creates a new LocalMemoryApiDiscoveryService instance
func NewLocalMemoryAPIDiscoveryService() *LocalMemoryAPIDiscoveryService {
	return &LocalMemoryAPIDiscoveryService{
		router: router.NewRoute(),
	}
}

// AddApi adds a method to the router tree
func (ads *LocalMemoryAPIDiscoveryService) AddApi(request service.DiscoveryRequest) (service.DiscoveryResponse, error) {
	aj := model.NewApi()
	if err := json.Unmarshal(request.Body, aj); err != nil {
		return *service.EmptyDiscoveryResponse, err
	}

	if _, loaded := model.CacheApi.LoadOrStore(aj.Name, aj); loaded {
		// loaded
		logger.Warnf("api:%s is exist", aj)
	} else {
		// store
		if aj.Metadata == nil {

		} else {
			if v, ok := aj.Metadata.(map[string]interface{}); ok {
				if d, ok := v["dubbo"]; ok {
					dm := &dubbo.DubboMetadata{}
					if err := mapstructure.Decode(d, dm); err != nil {
						return *service.EmptyDiscoveryResponse, err
					}
					aj.Metadata = dm
				}
			}

			aj.RequestMethod = model.RequestMethod(model.RequestMethodValue[aj.Method])
		}
	}

	return *service.NewDiscoveryResponseWithSuccess(true), nil
}

func (ads *LocalMemoryAPIDiscoveryService) GetApi(request service.DiscoveryRequest) (service.DiscoveryResponse, error) {
	n := string(request.Body)

	if a, ok := model.CacheApi.Load(n); ok {
		return *service.NewDiscoveryResponse(a), nil
	}

	return *service.EmptyDiscoveryResponse, errors.New("not found")
}

// AddAPI adds a method to the router tree
func (ads *LocalMemoryAPIDiscoveryService) AddAPI(api service.API) error {
	return ads.router.Put(api.URLPattern, api.Method)
}

// GetAPI returns the method to the caller
func (ads *LocalMemoryAPIDiscoveryService) GetAPI(url string, httpVerb config.HTTPVerb) (service.DiscoveryResponse, error) {
	if method, ok := ads.router.FindMethod(url, httpVerb); ok {
		return *service.NewDiscoveryResponse(method), nil
	}

	return *service.EmptyDiscoveryResponse, errors.New("not found")
}

// InitAPIsFromConfig inits the router from API config and to local cache
func InitAPIsFromConfig(apiConfig *config.APIConfig) error {
	localAPIDiscSrv := extension.GetMustApiDiscoveryService(constant.LocalMemoryApiDiscoveryService)
	if len(apiConfig.Resources) == 0 {
		return nil
	}
	return loadAPIFromResource("", apiConfig.Resources, localAPIDiscSrv)
}

func loadAPIFromResource(parrentPath string, resources []config.Resource, localSrv service.ApiDiscoveryService) error {
	errStack := []string{}
	if len(resources) == 0 {
		return nil
	}
	groupPath := parrentPath
	if parrentPath == constant.PathSlash {
		groupPath = ""
	}
	for _, resource := range resources {
		fullPath := groupPath + resource.Path
		if !strings.HasPrefix(resource.Path, constant.PathSlash) {
			errStack = append(errStack, fmt.Sprintf("Path %s in %s doesn't start with /", resource.Path, parrentPath))
			continue
		}
		if len(resource.Resources) > 0 {
			if err := loadAPIFromResource(resource.Path, resource.Resources, localSrv); err != nil {
				errStack = append(errStack, err.Error())
			}
		}

		if err := loadAPIFromMethods(fullPath, resource.Methods, localSrv); err != nil {
			errStack = append(errStack, err.Error())
		}
	}
	if len(errStack) > 0 {
		return errors.New(strings.Join(errStack, "; "))
	}
	return nil
}

func loadAPIFromMethods(fullPath string, methods []config.Method, localSrv service.ApiDiscoveryService) error {
	errStack := []string{}
	for _, method := range methods {
		api := service.API{
			URLPattern: fullPath,
			Method:     method,
		}
		if err := localSrv.AddAPI(api); err != nil {
			errStack = append(errStack, fmt.Sprintf("Path: %s, Method: %s, error: %s", fullPath, method.HTTPVerb, err.Error()))
		}
	}
	if len(errStack) > 0 {
		return errors.New(strings.Join(errStack, "\n"))
	}
	return nil
}

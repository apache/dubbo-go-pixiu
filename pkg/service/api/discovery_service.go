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
	"errors"
	"fmt"
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	pc "github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/plugins"
	"github.com/apache/dubbo-go-pixiu/pkg/router"
	"github.com/apache/dubbo-go-pixiu/pkg/service"
)

import (
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	fr "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"
)

// Init set api discovery local_memory service.
func Init() {
	extension.SetAPIDiscoveryService(constant.LocalMemoryApiDiscoveryService, NewLocalMemoryAPIDiscoveryService())
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

// AddAPI adds a method to the router tree
func (ads *LocalMemoryAPIDiscoveryService) AddAPI(api fr.API) error {
	return ads.router.PutAPI(api)
}

// GetAPI returns the method to the caller
func (ads *LocalMemoryAPIDiscoveryService) GetAPI(url string, httpVerb config.HTTPVerb) (fr.API, error) {
	if api, ok := ads.router.FindAPI(url, httpVerb); ok {
		return *api, nil
	}

	return fr.API{}, errors.New("not found")
}

// ClearAPI clear all api
func (ads *LocalMemoryAPIDiscoveryService) ClearAPI() error {
	ads.router.ClearAPI()
	return nil
}

// APIConfigChange to response to api config change
func (ads *LocalMemoryAPIDiscoveryService) APIConfigChange(apiConfig config.APIConfig) bool {
	ads.ClearAPI()
	// load pluginsGroup
	plugins.InitPluginsGroup(apiConfig.PluginsGroup, apiConfig.PluginFilePath)
	// init plugins from resource
	plugins.InitAPIURLWithFilterChain(apiConfig.Resources)
	loadAPIFromResource("", apiConfig.Resources, nil, ads)
	return true
}

// InitAPIsFromConfig inits the router from API config and to local cache
func InitAPIsFromConfig(apiConfig config.APIConfig) error {
	localAPIDiscSrv := extension.GetMustAPIDiscoveryService(constant.LocalMemoryApiDiscoveryService)
	if len(apiConfig.Resources) == 0 {
		return nil
	}
	// register config change listener
	pc.RegisterConfigListener(localAPIDiscSrv)
	return loadAPIFromResource("", apiConfig.Resources, nil, localAPIDiscSrv)
}

// RefreshAPIsFromConfig fresh the router from API config and to local cache
func RefreshAPIsFromConfig(apiConfig config.APIConfig) error {
	localAPIDiscSrv := NewLocalMemoryAPIDiscoveryService()
	if len(apiConfig.Resources) == 0 {
		return nil
	}
	error := loadAPIFromResource("", apiConfig.Resources, nil, localAPIDiscSrv)
	if error == nil {
		extension.SetAPIDiscoveryService(constant.LocalMemoryApiDiscoveryService, localAPIDiscSrv)
	}
	return error
}

func loadAPIFromResource(parrentPath string, resources []config.Resource, parentHeaders map[string]string, localSrv service.APIDiscoveryService) error {
	errStack := []string{}
	if len(resources) == 0 {
		return nil
	}
	groupPath := parrentPath
	if parrentPath == constant.PathSlash {
		groupPath = ""
	}
	fullHeaders := parentHeaders
	if fullHeaders == nil {
		fullHeaders = make(map[string]string, 9)
	}
	for _, resource := range resources {
		fullPath := groupPath + resource.Path
		if !strings.HasPrefix(resource.Path, constant.PathSlash) {
			errStack = append(errStack, fmt.Sprintf("Path %s in %s doesn't start with /", resource.Path, parrentPath))
			continue
		}
		for headerName, headerValue := range resource.Headers {
			fullHeaders[headerName] = headerValue
		}
		if len(resource.Resources) > 0 {
			if err := loadAPIFromResource(resource.Path, resource.Resources, fullHeaders, localSrv); err != nil {
				errStack = append(errStack, err.Error())
			}
		}

		if err := loadAPIFromMethods(fullPath, resource.Methods, fullHeaders, localSrv); err != nil {
			errStack = append(errStack, err.Error())
		}
	}
	if len(errStack) > 0 {
		return errors.New(strings.Join(errStack, "; "))
	}
	return nil
}

func loadAPIFromMethods(fullPath string, methods []config.Method, headers map[string]string, localSrv service.APIDiscoveryService) error {
	errStack := []string{}
	for _, method := range methods {
		api := fr.API{
			URLPattern: fullPath,
			Method:     method,
			Headers:    headers,
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

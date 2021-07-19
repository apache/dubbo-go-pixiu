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

package baseregistry

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	"github.com/apache/dubbo-go-pixiu/pkg/registry"
	"github.com/apache/dubbo-go/common"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"
)

type FacadeRegistry interface {
	// LoadInterfaces loads the dubbo services from interfaces level
	LoadInterfaces() ([]router.API, []error)
	// LoadApplication loads the dubbo services from application level
	LoadApplications() ([]router.API, []error)
	// DoSubscribe subscribes the registries to monitor the changes.
	DoSubscribe(*common.URL) error
	// DoUnsubscribe unsubscribes the registries.
	DoUnsubscribe(*common.URL) error
}

type BaseRegistry struct {
	listeners      []registry.Listener
	facadeRegistry FacadeRegistry
}

func NewBaseRegistry(facade FacadeRegistry) *BaseRegistry {
	return &BaseRegistry{
		listeners:      []registry.Listener{},
		facadeRegistry: facade,
	}
}

// LoadServices loads all the registered Dubbo services from registry
func (r *BaseRegistry) LoadServices() error {
	interfaceAPIs, _ := r.facadeRegistry.LoadInterfaces()
	applicationAPIs, _ := r.facadeRegistry.LoadApplications()
	apis := r.deduplication(append(interfaceAPIs, applicationAPIs...))
	localAPIDiscSrv := extension.GetMustAPIDiscoveryService(constant.LocalMemoryApiDiscoveryService)
	for i := range apis {
		localAPIDiscSrv.AddAPI(apis[i])
		// r.facadeRegistry.DoSubscribe()
	}
	return nil
}

func (r *BaseRegistry) deduplication(apis []router.API) []router.API {
	urlPatternMap := make(map[string]router.API)
	for i := range apis {
		urlPatternMap[apis[i].URLPattern] = apis[i]
	}
	rstAPIs := []router.API{}
	for i := range urlPatternMap {
		rstAPIs = append(rstAPIs, urlPatternMap[i])
	}
	return rstAPIs
}

// Subscribe monitors the target registry.
func (r *BaseRegistry) Subscribe(url *common.URL) error {
	return r.facadeRegistry.DoSubscribe(url)
}

// Unsubscribe stops monitoring the target registry.
func (r *BaseRegistry) Unsubscribe(url *common.URL) error {
	return r.facadeRegistry.DoUnsubscribe(url)
}

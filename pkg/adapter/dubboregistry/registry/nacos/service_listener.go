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

package nacos

import (
	"sync"
)

import (
	dubboCommon "dubbo.apache.org/dubbo-go/v3/common"
	dubboRegistry "dubbo.apache.org/dubbo-go/v3/registry"
	_ "dubbo.apache.org/dubbo-go/v3/registry/nacos"
	"dubbo.apache.org/dubbo-go/v3/remoting"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
)

import (
	common2 "github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/common"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// serviceListener normally monitors the /dubbo/[:url.service()]/providers
type serviceListener struct {
	url    *dubboCommon.URL
	client naming_client.INamingClient

	exit            chan struct{}
	wg              sync.WaitGroup
	adapterListener common2.RegistryEventListener
}

// WatchAndHandle todo WatchAndHandle is useless for service listener
func (z *serviceListener) WatchAndHandle() {
	panic("implement me")
}

// newNacosSrvListener creates a new zk service listener
func newNacosSrvListener(url *dubboCommon.URL, client naming_client.INamingClient, adapterListener common2.RegistryEventListener) *serviceListener {
	return &serviceListener{
		url:             url,
		client:          client,
		exit:            make(chan struct{}),
		adapterListener: adapterListener,
	}
}

func (z *serviceListener) Notify(e *dubboRegistry.ServiceEvent) {
	bkConfig, methods, location, err := registry.ParseDubboString(e.Service.String())
	if err != nil {
		logger.Errorf("parse dubbo string error = %s", err)
		return
	}

	mappingParams := []config.MappingParam{
		{
			Name:  "requestBody.values",
			MapTo: "opt.values",
		},
		{
			Name:  "requestBody.types",
			MapTo: "opt.types",
		},
	}
	apiPattern := registry.GetAPIPattern(bkConfig)
	for i := range methods {
		api := registry.CreateAPIConfig(apiPattern, location, bkConfig, methods[i], mappingParams)
		if e.Action == remoting.EventTypeDel {
			if err := z.adapterListener.OnRemoveAPI(api); err != nil {
				logger.Errorf("Error={%s} happens when try to remove api %s", err.Error(), api.Path)
				continue
			}
		} else {
			if err := z.adapterListener.OnAddAPI(api); err != nil {
				logger.Errorf("Error={%s} happens when try to add api %s", err.Error(), api.Path)
				continue
			}
		}

	}
}

func (z *serviceListener) NotifyAll(e []*dubboRegistry.ServiceEvent, f func()) {
}

// Close closes this listener
func (zkl *serviceListener) Close() {
	close(zkl.exit)
	zkl.wg.Wait()
}

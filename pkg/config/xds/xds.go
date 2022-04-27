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

package xds

import (
	"sync"
)

import (
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/xds"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/config/xds/apiclient"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server/controls"
)

type (
	DiscoverApi interface {
		Fetch(localVersion string) ([]*apiclient.ProtoAny, error)
		Delta() (chan *apiclient.DeltaResources, error)
	}

	AdapterConfig struct {
	}

	Xds struct {
		//ads    DiscoverApi //aggregate discover service manager todo to implement
		cds               *CdsManager //cluster discover service manager
		lds               *LdsManager //listener discover service manager
		exitCh            chan struct{}
		listenerMg        controls.ListenerManager
		clusterMg         controls.ClusterManager
		dynamicResourceMg controls.DynamicResourceManager
	}
)

func (a *Xds) createApiManager(config *model.ApiConfigSource,
	node *model.Node,
	resourceTypes ...apiclient.ResourceTypeName) DiscoverApi {
	if config == nil {
		return nil
	}

	switch config.APIType {
	case model.ApiTypeGRPC:
		return apiclient.CreateGrpcApiClient(config, node, a.exitCh, resourceTypes...)
	default:
		logger.Errorf("un-support the api type %s", config.APITypeStr)
		return nil
	}
}

func (a *Xds) Start() {
	if a.dynamicResourceMg == nil { // if dm is nil, then config not initialized.
		logger.Infof("can not get dynamic resource manager. maybe the config has not initialized")
		return
	}
	apiclient.Init(a.clusterMg)

	// lds fetch just run on init phase.
	if a.dynamicResourceMg.GetLds() != nil {
		a.lds = &LdsManager{
			DiscoverApi: a.createApiManager(a.dynamicResourceMg.GetLds(), a.dynamicResourceMg.GetNode(), xds.ListenerType),
			listenerMg:  a.listenerMg,
		}
		if err := a.lds.Delta(); err != nil {
			logger.Errorf("can not fetch lds err is %+v", err)
		}
	}
	// catch the ongoing cds config change.
	if a.dynamicResourceMg.GetCds() != nil {
		a.cds = &CdsManager{
			DiscoverApi: a.createApiManager(a.dynamicResourceMg.GetCds(), a.dynamicResourceMg.GetNode(), xds.ClusterType),
			clusterMg:   a.clusterMg,
		}
		if err := a.cds.Delta(); err != nil {
			logger.Errorf("can not fetch lds")
		}
	}

}

func (a *Xds) Stop() {
	apiclient.Stop()
	close(a.exitCh)
}

var (
	client Client
	once   sync.Once
)

// Client xds client
type Client interface {
	Stop()
}

// StartXdsClient create XdsClient and run. only one xds client create at first(singleton)
func StartXdsClient(listenerMg controls.ListenerManager, clusterMg controls.ClusterManager, drm controls.DynamicResourceManager) Client {
	once.Do(func() {
		xdsClient := &Xds{
			listenerMg:        listenerMg,
			clusterMg:         clusterMg,
			dynamicResourceMg: drm,
			exitCh:            make(chan struct{}),
		}
		xdsClient.Start()
		client = xdsClient
	})

	return client
}

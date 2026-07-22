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
	stderr "errors"
	"sync"
)

import (
	"github.com/mitchellh/mapstructure"

	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
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
	resourceType apiclient.ResourceTypeName) (DiscoverApi, error) {
	if config == nil {
		return nil, nil
	}

	switch config.APIType {
	case model.ApiTypeGRPC:
		client, err := apiclient.CreateGrpExtensionApiClient(config, node, a.exitCh, resourceType)
		if err != nil {
			return nil, errors.Wrap(err, "create grpc extension api client")
		}
		return client, nil
	case model.ApiTypeIstioGRPC:
		dubboServices, err := a.readDubboServiceFromListener()
		if err != nil {
			return nil, errors.Wrap(err, "read dubbo service from listener")
		}
		client, err := apiclient.CreateEnvoyGrpcApiClient(config, node, a.exitCh, resourceType, apiclient.WithIstioService(dubboServices...))
		if err != nil {
			return nil, errors.Wrap(err, "create envoy grpc api client")
		}
		return client, nil
	default:
		return nil, errors.Errorf("un-support the api type %s", config.APITypeStr)
	}
}

func (a *Xds) readDubboServiceFromListener() ([]string, error) {
	dubboServices := make([]string, 0)
	listeners, err := a.listenerMg.CloneXdsControlListener()
	if err != nil {
		return nil, err
	}

	for _, l := range listeners {
		for _, filter := range l.FilterChain.Filters {
			if filter.Name != constant.HTTPConnectManagerFilter {
				continue
			}
			var cfg *model.HttpConnectionManagerConfig
			if filter.Config != nil {
				if err := mapstructure.Decode(filter.Config, &cfg); err != nil {
					logger.Error("read listener config error when init xds", err)
					continue
				}
			}
			for _, httpFilter := range cfg.HTTPFilters {
				if httpFilter.Name == constant.HTTPDirectDubboProxyFilter {
					for _, route := range cfg.RouteConfig.Routes {
						dubboServices = append(dubboServices, route.Route.Cluster)
					}
				}
			}
		}
	}
	return dubboServices, nil
}

func (a *Xds) Start() error {
	if a.dynamicResourceMg == nil { // if dm is nil, then config not initialized.
		logger.Infof("can not get dynamic resource manager. maybe the config has not initialized")
		return nil
	}
	apiclient.Init(a.clusterMg)

	// lds and cds are independent dynamic resource watches: a failure in one must not skip the
	// other's initialization. Collect each error and fail startup if any occurred, so a missing
	// config or xDS connection failure surfaces to the caller instead of leaving the process in a
	// false-healthy state with parts of the dynamic config unavailable.
	var errs []error

	// lds fetch just run on init phase.
	if a.dynamicResourceMg.GetLds() != nil {
		discoverApi, err := a.createApiManager(a.dynamicResourceMg.GetLds(), a.dynamicResourceMg.GetNode(), constant.ListenerType)
		if err != nil {
			errs = append(errs, errors.Wrap(err, "create LDS api manager"))
		} else {
			a.lds = &LdsManager{
				DiscoverApi: discoverApi,
				listenerMg:  a.listenerMg,
			}
			if err := a.lds.Delta(); err != nil {
				errs = append(errs, errors.Wrap(err, "fetch lds"))
			}
		}
	}
	// catch the ongoing cds config change.
	if a.dynamicResourceMg.GetCds() != nil {
		discoverApi, err := a.createApiManager(a.dynamicResourceMg.GetCds(), a.dynamicResourceMg.GetNode(), constant.ClusterType)
		if err != nil {
			errs = append(errs, errors.Wrap(err, "create CDS api manager"))
		} else {
			a.cds = &CdsManager{
				DiscoverApi: discoverApi,
				clusterMg:   a.clusterMg,
			}
			if err := a.cds.Delta(); err != nil {
				errs = append(errs, errors.Wrap(err, "fetch cds"))
			}
		}
	}

	return stderr.Join(errs...)
}

func (a *Xds) Stop() {
	apiclient.Stop()
	close(a.exitCh)
}

var (
	client   Client
	startErr error
	once     sync.Once
)

// Client xds client
type Client interface {
	Stop()
}

// StartXdsClient create XdsClient and run. only one xds client create at first(singleton)
func StartXdsClient(listenerMg controls.ListenerManager, clusterMg controls.ClusterManager, drm controls.DynamicResourceManager) (Client, error) {
	// Note: on first-call failure once.Do still completes, leaving client == nil, so subsequent
	// calls return (nil, nil). This is acceptable because the caller fails the whole startup on
	// the first error and never reaches a second call.
	once.Do(func() {
		xdsClient := &Xds{
			listenerMg:        listenerMg,
			clusterMg:         clusterMg,
			dynamicResourceMg: drm,
			exitCh:            make(chan struct{}),
		}
		if err := xdsClient.Start(); err != nil {
			client = nil
			startErr = err
			return
		}
		client = xdsClient
	})

	if client == nil {
		return nil, startErr
	}
	return client, nil
}

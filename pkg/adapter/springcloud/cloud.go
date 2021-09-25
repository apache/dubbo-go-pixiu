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

package springcloud

import (
	"sync"
	"time"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud/servicediscovery"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud/servicediscovery/consul"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud/servicediscovery/nacos"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/adapter"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

const (
	// Kind is the kind of Adapter Plugin.
	Kind = constant.SpringCloudAdapter
)

func init() {
	adapter.RegisterAdapterPlugin(&CloudPlugin{})
}

type (
	// CloudPlugin the adapter plugin for spring cloud
	CloudPlugin struct {
	}

	// CloudAdapter the adapter for spring cloud
	CloudAdapter struct {
		cfg *Config
		sd  servicediscovery.ServiceDiscovery

		mutex    sync.Mutex
		stopChan chan struct{}
	}

	// Config the config for CloudAdapter
	Config struct {
		// Mode default 0 start backgroup fresh and watch, 1 only fresh 2 only watch
		Mode          int                 `yaml:"mode" json:"mode" default:"mode"`
		Registry      *model.RemoteConfig `yaml:"registry" json:"registry" default:"registry"`
		FreshInterval time.Duration       `yaml:"freshInterval" json:"freshInterval" default:"freshInterval"`
	}
)

// Kind return plugin kind
func (p *CloudPlugin) Kind() string {
	return Kind
}

// CreateAdapter create adapter
func (p *CloudPlugin) CreateAdapter(ad *model.Adapter) (adapter.Adapter, error) {
	return &CloudAdapter{cfg: &Config{}, stopChan: make(chan struct{})}, nil
}

// Start start the adapter
func (a *CloudAdapter) Start(adapter *model.Adapter) {
	// init get all service instance
	a.firstFetch()
	// backgroup sync service instance from remote
	// do not block the main goroutine
	a.backgroupSyncPeriod()
	//a.watch()
}

// Stop stop the adapter
func (a *CloudAdapter) Stop() {
	a.stop()
}

// Apply init
func (a *CloudAdapter) Apply() error {
	//registryUsed := ad.Config["registry"].(map[string]interface{})
	switch a.cfg.Registry.Protocol {
	case "nacos":
		sd, err := nacos.NewNacosServiceDiscovery(a.cfg.Registry)
		if err != nil {
			logger.Errorf("Apply NewNacosServiceDiscovery", err.Error())
			return err
		}
		a.sd = sd
	case "consul":
		sd, err := consul.NewConsulServiceDiscovery(a.cfg.Registry)
		if err != nil {
			logger.Errorf("new consul client fail : ", err.Error())
			return err
		}
		a.sd = sd
	case "zookeeper":
	default:
		return errors.New("adapter init error registry not recognise")
	}
	return nil
}

// Config get config for Adapter
func (a *CloudAdapter) Config() interface{} {
	return a.cfg
}

func (a *CloudAdapter) firstFetch() error {
	instances, err := a.sd.QueryServices()
	if err != nil {
		logger.Errorf("start query all service error ", err.Error())
		return err
	}
	// manage cluster and route
	cm := server.GetClusterManager()
	rm := server.GetRouterManager()
	for _, instance := range instances {
		endpoint := instance.ToEndpoint()
		// todo: maybe instance service name not equal with cluster name ?
		cm.AddEndpoint(endpoint.Name, endpoint)
		// route ID is cluster name, so equal with endpoint name
		route := instance.ToRoute()
		rm.AddRouter(route)
	}
	return nil
}

func (a *CloudAdapter) fetchCompareAndSet() {
	instances, err := a.sd.QueryServices()
	if err != nil {
		logger.Warnf("fetchCompareAndSet all service error ", err.Error())
	}
	// manage cluster and route
	cm := server.GetClusterManager()
	rm := server.GetRouterManager()

	a.mutex.Lock()
	defer a.mutex.Unlock()

	oldStore, err := cm.CloneStore()

	if err != nil {
		logger.Warnf("fetchCompareAndSet clone store error ", err.Error())
	}

	newStore := cm.NewStore(oldStore.Version)

	for _, instance := range instances {
		endpoint := instance.ToEndpoint()
		// endpoint name should equal with cluster name
		newStore.AddEndpoint(endpoint.Name, endpoint)
	}

	// maximize reduction the interval of down state
	// first remove the router for removed cluster
	for _, c := range oldStore.Config {
		if !newStore.HasCluster(c.Name) {
			delete := &model.Router{ID: c.Name}
			rm.DeleteRouter(delete)
		}
	}
	// second set cluster
	ret := cm.CompareAndSetStore(newStore)

	if !ret {
		// fast fail the delete route at first phase shouldn't recover
		return
	}

	// third add new router
	for _, c := range newStore.Config {
		if !oldStore.HasCluster(c.Name) {
			match := model.RouterMatch{Prefix: c.Name}
			route := model.RouteAction{Cluster: c.Name}
			added := &model.Router{ID: c.Name, Match: match, Route: route}
			rm.AddRouter(added)
		}
	}
}

func (a *CloudAdapter) backgroupSyncPeriod() error {
	timer := time.NewTicker(a.cfg.FreshInterval)
	go func() {
		defer timer.Stop()
		for {
			select {
			case <-timer.C:
				a.fetchCompareAndSet()
				break
			case <-a.stopChan:
				logger.Info("stop the adapter")
				return
			}
		}
	}()
	return nil
}

func (a *CloudAdapter) watch() error {
	return nil
}

func (a *CloudAdapter) stop() error {
	a.sd.Stop()
	close(a.stopChan)
	return nil
}

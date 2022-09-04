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
	"strings"
	"sync"
	"time"
)

import (
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/springcloud/servicediscovery"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/springcloud/servicediscovery/nacos"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/springcloud/servicediscovery/zookeeper"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/adapter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server"
)

const (
	// Kind is the kind of Adapter Plugin.
	Kind = constant.SpringCloudAdapter

	Nacos     = "nacos"
	Zookeeper = "zookeeper"
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
		cfg            *Config
		sd             servicediscovery.ServiceDiscovery
		currentService map[string]*Service

		mutex    sync.Mutex
		stopChan chan struct{}
	}

	// Config the config for CloudAdapter
	Config struct {
		// Mode default 0 start backgroup fresh and watch, 1 only fresh 2 only watch
		Mode          int                 `yaml:"mode" json:"mode" default:"mode"`
		Registry      *model.RemoteConfig `yaml:"registry" json:"registry" default:"registry"`
		FreshInterval time.Duration       `yaml:"freshInterval" json:"freshInterval" default:"freshInterval"`
		Services      []string            `yaml:"services" json:"services" default:"services"`
		// SubscribePolicy subscribe config,
		// - adapting : if there is no any Services (App) names, fetch All services from registry center
		// - definitely : fetch services by the config Services (App) names
		SubscribePolicy string `yaml:"subscribe-policy" json:"subscribe-policy" default:"adapting"`
	}

	Service struct {
		Name string
	}

	SubscribePolicy int
)

const (
	Adapting SubscribePolicy = iota
	Definitely
)

func (sp SubscribePolicy) String() string {
	return [...]string{"adapting", "definitely"}[sp]
}

// Kind return plugin kind
func (p *CloudPlugin) Kind() string {
	return Kind
}

// CreateAdapter create adapter
func (p *CloudPlugin) CreateAdapter(ad *model.Adapter) (adapter.Adapter, error) {
	return &CloudAdapter{cfg: &Config{}, stopChan: make(chan struct{}), currentService: make(map[string]*Service)}, nil
}

// Start start the adapter
func (a *CloudAdapter) Start() {

	// do not block the main goroutine
	// init get all service instance
	err := a.firstFetch()
	if err != nil {
		logger.Errorf("init fetch service fail", err.Error())
	}

	// background sync service instance from remote
	err = a.backgroundSyncPeriod()
	if err != nil {
		logger.Errorf("init periodicity fetch service task fail", err.Error())
	}

	// watch then fetch is more safety for consistent but there is background fresh mechanism
	err = a.watch()
	if err != nil {
		logger.Errorf("init watch the register fail", err.Error())
	}
}

// Stop stop the adapter
func (a *CloudAdapter) Stop() {
	err := a.stop()
	if err != nil {
		logger.Errorf("stop the adapter fail", err.Error())
		return
	}
}

// Apply init
func (a *CloudAdapter) Apply() error {
	//registryUsed := ad.Config["registry"].(map[string]interface{})
	var (
		sd  servicediscovery.ServiceDiscovery
		err error
	)
	switch strings.ToLower(a.cfg.Registry.Protocol) {
	case Nacos:
		sd, err = nacos.NewNacosServiceDiscovery(a.cfg.Services, a.cfg.Registry, a)
	case Zookeeper:
		sd, err = zookeeper.NewZKServiceDiscovery(a.cfg.Services, a.cfg.Registry, a)
	default:
		return errors.New("adapter init error registry not recognise")
	}
	if err != nil {
		logger.Errorf("Apply NewServiceDiscovery %s ", a.cfg.Registry.Protocol, err.Error())
		return err
	}
	a.sd = sd
	return nil
}

func (a *CloudAdapter) OnAddServiceInstance(instance *servicediscovery.ServiceInstance) {
	cm := server.GetClusterManager()
	endpoint := instance.ToEndpoint()

	a.mutex.Lock()
	defer a.mutex.Unlock()

	// endpoint name should equal with cluster name
	cm.SetEndpoint(endpoint.Name, endpoint)
	if a.checkHasExistService(endpoint.Name) {
		return
	}
	// new service, so add route and into CurrentService map
	a.addNewService(instance)
	// route ID is cluster name, so equal with endpoint name
	rm := server.GetRouterManager()
	route := instance.ToRoute()
	rm.AddRouter(route)
}

func (a *CloudAdapter) OnDeleteServiceInstance(instance *servicediscovery.ServiceInstance) {
	cm := server.GetClusterManager()
	endpoint := instance.ToEndpoint()

	cm.DeleteEndpoint(endpoint.Name, endpoint.ID)

}

func (a *CloudAdapter) OnUpdateServiceInstance(instance *servicediscovery.ServiceInstance) {
	cm := server.GetClusterManager()
	endpoint := instance.ToEndpoint()
	a.mutex.Lock()
	defer a.mutex.Unlock()
	cm.SetEndpoint(endpoint.Name, endpoint)
}

func (a *CloudAdapter) GetServiceNames() []string {
	a.mutex.Lock()
	defer a.mutex.Unlock()

	res := make([]string, 0, len(a.currentService))

	for k := range a.currentService {
		res = append(res, k)
	}

	return res
}

// Config get config for Adapter
func (a *CloudAdapter) Config() interface{} {
	return a.cfg
}

func (a *CloudAdapter) fetchServiceByConfig() ([]servicediscovery.ServiceInstance, error) {
	var instances []servicediscovery.ServiceInstance
	var err error
	// if configure specific services, then fetch those service instance only
	if a.subscribeServiceDefinitely() {
		if len(a.cfg.Services) > 0 {
			instances, err = a.sd.QueryServicesByName(a.cfg.Services)
		} else {
			logger.Warnf("No any Service(App) need Subscribe, config the Service(App) Names or make the `subscribe-policy: adapting` pls.")
		}
	} else {
		if len(a.cfg.Services) > 0 {
			instances, err = a.sd.QueryServicesByName(a.cfg.Services)
		} else {
			instances, err = a.sd.QueryAllServices()
		}
	}

	if err != nil {
		logger.Errorf("fetchServiceByConfig error ", err.Error())
		return instances, err
	}
	return instances, nil
}

func (a *CloudAdapter) firstFetch() error {

	instances, err := a.fetchServiceByConfig()

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
		cm.SetEndpoint(endpoint.Name, endpoint)
		// route ID is cluster name, so equal with endpoint name
		route := instance.ToRoute()
		rm.AddRouter(route)
	}
	a.clearAndSetCurrentService(instances)
	return nil
}

func (a *CloudAdapter) clearAndSetCurrentService(instances []servicediscovery.ServiceInstance) {
	a.currentService = make(map[string]*Service)

	for _, instance := range instances {
		if _, ok := a.currentService[instance.ServiceName]; ok {
			continue
		}
		a.currentService[instance.ServiceName] = &Service{Name: instance.ServiceName}
	}
}

func (a *CloudAdapter) checkHasExistService(name string) bool {
	_, ok := a.currentService[name]
	return ok
}

func (a *CloudAdapter) addNewService(instance *servicediscovery.ServiceInstance) {
	a.currentService[instance.ServiceName] = &Service{Name: instance.ServiceName}
}

func (a *CloudAdapter) fetchCompareAndSet() {
	instances, err := a.fetchServiceByConfig()
	if err != nil {
		logger.Warnf("fetchCompareAndSet all service error ", err.Error())
		return
	}
	_ = a.watch()
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
		newStore.SetEndpoint(endpoint.Name, endpoint)
	}

	// maximize reduction the interval of down state
	// first remove the router for removed cluster
	for _, c := range oldStore.Config {
		if !newStore.HasCluster(c.Name) {
			rm.DeleteRouter(&model.Router{ID: c.Name})
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
			match := model.NewRouterMatchPrefix(c.Name)
			route := model.RouteAction{Cluster: c.Name}
			added := &model.Router{ID: c.Name, Match: match, Route: route}
			rm.AddRouter(added)
		}
	}
}

func (a *CloudAdapter) backgroundSyncPeriod() error {
	if a.cfg.FreshInterval <= 0 {
		return nil
	}
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
	return a.sd.Subscribe()
}

func (a *CloudAdapter) stop() error {
	err := a.sd.Unsubscribe()
	if err != nil {
		logger.Errorf("unsubscribe registry fail ", err.Error())
	}
	close(a.stopChan)
	return nil
}

func (a *CloudAdapter) subscribeServiceDefinitely() bool {
	return strings.EqualFold(a.cfg.SubscribePolicy, Definitely.String())
}

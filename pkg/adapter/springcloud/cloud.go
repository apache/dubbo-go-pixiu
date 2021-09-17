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
	"fmt"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/adapter"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/registry"
	"github.com/apache/dubbo-go-pixiu/pkg/registry/nacos"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
	"net"
	"strconv"
	"strings"
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
		cfg        *Config
		LoaderUsed registry.Loader
	}

	// Config the config for CloudAdapter
	Config struct {
		Registry *Registry `yaml:"registrys" json:"registrys" default:"registrys"`
	}

	// Registry remote registry where spring cloud apis are registered.
	Registry struct {
		Protocol string `yaml:"protocol" json:"protocol" default:"zookeeper"`
		Timeout  string `yaml:"timeout" json:"timeout"`
		Address  string `yaml:"address" json:"address"`
		Username string `yaml:"username" json:"username"`
		Password string `yaml:"password" json:"password"`
	}
)

// Kind return plugin kind
func (p *CloudPlugin) Kind() string {
	return Kind
}

// CreateAdapter create adapter
func (p *CloudPlugin) CreateAdapter(config interface{}, ad *model.Adapter) (adapter.Adapter, error) {
	//registryUsed := ad.Config["registry"].(map[string]interface{})
	registryUsed := config.(map[string]interface{})
	fmt.Println(registryUsed)
	var loader registry.Loader
	reg := registryUsed["registries"].(map[string]interface{})
	for key, _ := range reg {
		switch key {
		case "Nacos":
			loader = new(nacos.NacosRegistry)
		case "consul":
			loader = new(registry.ConsulRegistryLoad)
		case "zookeeper":
			loader = new(registry.ZookeeperRegistryLoad)
		}
	}
	return &CloudAdapter{cfg: &Config{}, LoaderUsed: loader}, nil
}

// Start start the adapter
func (a *CloudAdapter) Start(adapter *model.Adapter) {
	reg := a.LoaderUsed
	RegistryLoader, err := reg.NewRegistryLoader(adapter)
	if err != nil {
		logger.Info("fail to get registyloader")
	}
	// do not block the main goroutine
	go func() {
		//init registry client

		// init SpringCloud Manager for control initialize
		cloudManager := NewSpringCloudManager(&SpringCloudConfig{boot: a.boot})

		cloudManager.Start()

		// fetch service instance from consul

		// transform into endpoint and cluster
		//endpoint := &model.Endpoint{}
		//endpoint.ID = "user-mock-service"
		//endpoint.Address = model.SocketAddress{
		//	Address: "127.0.0.1",
		//	Port:    8080,
		//}
		//cluster := &model.Cluster{}
		//cluster.Name = "userservice"
		//cluster.Lb = model.Rand
		//cluster.Endpoints = []*model.Endpoint{
		//	endpoint,
		//}
		//// add cluster into manager
		//cm := server.GetClusterManager()
		//cm.AddCluster(cluster)
		//
		//// transform into route
		//routeMatch := model.RouterMatch{
		//	Prefix: "/user",
		//}
		//routeAction := model.RouteAction{
		//	Cluster: "userservice",
		//}
		//route := &model.Router{Match: routeMatch, Route: routeAction}
		//
		//server.GetRouterManager().AddRouter(route)

		urls, err := RegistryLoader.LoadAllServices()
		if err != nil {
			logger.Error("can't get service for %s re")
		}
		var endpoints []*model.Endpoint
		for _, url := range urls {
			endpoint := url.GetParam("endpoint", "")
			tmp := strings.Split(endpoint, ",")
			for _, path := range tmp {
				ep := &model.Endpoint{}
				ip, port, err := net.SplitHostPort(path)
				porti, err := strconv.Atoi(port)
				if err != nil {
					logger.Info("split ip & port failed")
				}
				ep.ID = url.GetParam("Name", "")
				ep.Address = model.SocketAddress{
					Address: ip,
					Port:    porti,
				}
				endpoints = append(endpoints, ep)
			}
			// get cluster
			cluster := &model.Cluster{}
			cluster.Name = url.GetParam(constant.NameKey, "")
			lb := url.GetParam("loadbalance", "")
			switch lb {
			case "RoundRobin":
				cluster.Lb = model.RoundRobin
			case "IPHash":
				cluster.Lb = model.IPHash
			case "WightRobin":
				cluster.Lb = model.WightRobin
			case "Rand":
				cluster.Lb = model.Rand
			}
			cluster.Endpoints = endpoints
			cm := server.GetClusterManager()
			cm.AddCluster(cluster)
			// transform into route
			routeMatch := model.RouterMatch{
				Prefix: "/user",
			}
			routeAction := model.RouteAction{
				Cluster: "userservice",
			}
			route := &model.Router{Match: routeMatch, Route: routeAction}

			server.GetRouterManager().AddRouter(route)
		}
	}()
}

// Stop stop the adapter
func (a *CloudAdapter) Stop() {

}

// Apply init
func (a *CloudAdapter) Apply() error {

	return nil
}

// Config get config for Adapter
func (a *CloudAdapter) Config() interface{} {
	return a.cfg
}

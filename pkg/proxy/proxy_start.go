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

package proxy

import (
	"encoding/json"
	"github.com/dubbogo/dubbo-go-proxy/pkg/api_load"
	"sync"
	"time"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client/dubbo"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	_ "github.com/dubbogo/dubbo-go-proxy/pkg/filter"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
	"github.com/dubbogo/dubbo-go-proxy/pkg/model"
	"github.com/dubbogo/dubbo-go-proxy/pkg/service"
	_ "github.com/dubbogo/dubbo-go-proxy/pkg/service/api"
)

// Proxy
type Proxy struct {
	startWG sync.WaitGroup
	bs      *model.Bootstrap
}

// Start proxy start
func (p *Proxy) Start() {
	conf := config.GetBootstrap()

	p.startWG.Add(1)

	defer func() {
		if re := recover(); re != nil {
			logger.Error(re)
			// TODO stop
		}
	}()

	p.beforeStart()

	listeners := conf.GetListeners()

	for _, s := range listeners {
		ls := ListenerService{Listener: &s}
		go ls.Start()
	}
}

func (p *Proxy) beforeStart() error {
	dubbo.SingleDubboClient().Init()

	// TODO mock api register
	ads := extension.GetMustApiDiscoveryService(constant.LocalMemoryApiDiscoveryService)

	apiLoader := api_load.NewApiManager(time.Second, ads)
	apiLoader.AddApiLoad(p.bs.DynamicResources.ApiConfig)
	err := apiLoader.StartLoadApi()
	if err != nil {
		logger.Errorf("error load api:%v", err)
		return err
	}
	err = apiLoader.SelectMergeApiTask()
	if err != nil {
		logger.Errorf("error select merge api task :%v", err)
	}
	a1 := &model.Api{
		Name:     "/api/v1/test-dubbo/user",
		ITypeStr: "HTTP",
		OTypeStr: "DUBBO",
		Method:   "POST",
		Status:   1,
		Metadata: map[string]dubbo.DubboMetadata{
			"dubbo": {
				ApplicationName: "BDTService",
				Group:           "test",
				Version:         "1.0.0",
				Interface:       "com.ikurento.user.UserProvider",
				Method:          "queryUser",
				Types: []string{
					"com.ikurento.user.User",
				},
				ClusterName: "test_dubbo",
			},
		},
	}
	a2 := &model.Api{
		Name:     "/api/v1/test-dubbo/getUserByName",
		ITypeStr: "HTTP",
		OTypeStr: "DUBBO",
		Method:   "POST",
		Status:   1,
		Metadata: map[string]dubbo.DubboMetadata{
			"dubbo": {
				ApplicationName: "BDTService",
				Group:           "test",
				Version:         "1.0.0",
				Interface:       "com.ikurento.user.UserProvider",
				Method:          "GetUser",
				Types: []string{
					"java.lang.String",
				},
				ClusterName: "test_dubbo",
			},
		},
	}

	j1, _ := json.Marshal(a1)
	j2, _ := json.Marshal(a2)
	ads.AddApi(*service.NewDiscoveryRequest(j1))
	ads.AddApi(*service.NewDiscoveryRequest(j2))
}

// NewProxy create proxy
func NewProxy(bs *model.Bootstrap) *Proxy {
	return &Proxy{
		startWG: sync.WaitGroup{},
		bs:      bs,
	}
}

func Start(bs *model.Bootstrap) {
	logger.Infof("[dubboproxy go] start by config : %+v", bs)

	proxy := NewProxy(bs)
	proxy.Start()

	proxy.startWG.Wait()
}

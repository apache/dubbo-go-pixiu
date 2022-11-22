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

package apiconfig

import (
	"net/http"
)

import (
	fc "github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/router"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/config"
	contexthttp "github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/http/apiconfig/api"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.HTTPApiConfigFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	Plugin struct {
	}

	FilterFactory struct {
		cfg        *ApiConfigConfig
		apiService api.APIDiscoveryService
	}
	Filter struct {
		apiService api.APIDiscoveryService
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &ApiConfigConfig{}}, nil
}

func (factory *FilterFactory) Config() interface{} {
	return factory.cfg
}

func (factory *FilterFactory) Apply() error {
	factory.apiService = api.NewLocalMemoryAPIDiscoveryService()

	if factory.cfg.Dynamic {
		server.GetApiConfigManager().AddApiConfigListener(factory.cfg.DynamicAdapter, factory)
		return nil
	}

	config, err := initApiConfig(factory.cfg)
	if err != nil {
		logger.Errorf("Get ApiConfig fail: %v", err)
	}
	if err := factory.apiService.InitAPIsFromConfig(*config); err != nil {
		logger.Errorf("InitAPIsFromConfig fail: %v", err)
	}

	return nil
}

func (factory *FilterFactory) OnAddAPI(r router.API) error {
	return factory.apiService.AddOrUpdateAPI(r)
}
func (factory *FilterFactory) OnRemoveAPI(r router.API) error {
	return factory.apiService.RemoveAPIByIntance(r)
}

func (factory *FilterFactory) OnDeleteRouter(r fc.Resource) error {
	return factory.apiService.RemoveAPIByPath(r)
}

func (factory *FilterFactory) GetAPIService() api.APIDiscoveryService {
	return factory.apiService
}

func (factory *FilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	f := &Filter{apiService: factory.apiService}
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(ctx *contexthttp.HttpContext) filter.FilterStatus {
	req := ctx.Request
	v, err := f.apiService.MatchAPI(req.URL.Path, fc.HTTPVerb(req.Method))
	if err != nil {
		ctx.SendLocalReply(http.StatusNotFound, constant.Default404Body)
		e := errors.Errorf("Requested URL %s not found", req.URL.Path)
		logger.Debug(e.Error())
		return filter.Stop
	}

	if !v.Method.Enable {
		ctx.SendLocalReply(http.StatusNotAcceptable, constant.Default406Body)
		e := errors.Errorf("Requested API %s %s does not online", req.Method, req.URL.Path)
		logger.Debug(e.Error())
		return filter.Stop
	}
	ctx.API(v)
	return filter.Continue
}

func (factory *FilterFactory) GetApiService() api.APIDiscoveryService {
	return factory.apiService
}

// initApiConfig return value of the bool is for the judgment of whether is a api meta data error, a kind of silly (?)
func initApiConfig(cf *ApiConfigConfig) (*fc.APIConfig, error) {
	if cf.APIMetaConfig != nil {
		a, err := config.LoadAPIConfig(cf.APIMetaConfig)
		if err != nil {
			logger.Warnf("load api config from etcd error:%+v", err)
			return nil, err
		}
		return a, nil
	}

	a, err := config.LoadAPIConfigFromFile(cf.Path)
	if err != nil {
		logger.Errorf("load api config error:%+v", err)
		return nil, err
	}
	return a, nil
}

var _ filter.HttpFilterFactory = new(FilterFactory)

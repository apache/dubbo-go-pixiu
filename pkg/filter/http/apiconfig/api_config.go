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
	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"

	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/http/apiconfig/api"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
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
		filterInstance *Filter
	}

	Filter struct {
		cfg        *ApiConfigConfig
		apiService api.APIDiscoveryService
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &Filter{cfg: &ApiConfigConfig{}}, nil
}

func (f *Plugin) GetInstance() *Filter {
	return f.filterInstance
}

func (f *Filter) Config() interface{} {
	return f.cfg
}

func (f *Filter) Apply() error {
	f.apiService = api.NewLocalMemoryAPIDiscoveryService()

	if f.cfg.Dynamic {
		server.GetApiConfigManager().AddApiConfigListener(f.cfg.DynamicAdapter, f)
		return nil
	}

	config, err := initApiConfig(f.cfg)
	if err != nil {
		logger.Errorf("Get ApiConfig fail: %v", err)
	}
	if err := f.apiService.InitAPIsFromConfig(*config); err != nil {
		logger.Errorf("InitAPIsFromConfig fail: %v", err)
	}

	return nil
}

func (f *Filter) OnAddAPI(r router.API) error {
	return f.apiService.AddAPI(r)
}
func (f *Filter) OnRemoveAPI(r router.API) error {
	return f.apiService.RemoveAPIByIntance(r)
}

func (f *Filter) OnDeleteRouter(r fc.Resource) error {
	return f.apiService.RemoveAPIByPath(r)
}

func (f *Filter) GetAPIService() api.APIDiscoveryService {
	return f.apiService
}

func (f *Filter) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *Filter) Handle(ctx *contexthttp.HttpContext) {
	req := ctx.Request
	api, err := f.apiService.GetAPI(req.URL.Path, fc.HTTPVerb(req.Method))
	if err != nil {
		if _, err := ctx.WriteWithStatus(http.StatusNotFound, constant.Default404Body); err != nil {
			logger.Errorf("WriteWithStatus fail: %v", err)
		}
		ctx.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
		e := errors.Errorf("Requested URL %s not found", req.URL.Path)
		logger.Debug(e.Error())
		ctx.Abort()
		return
	}

	if !api.Method.Enable {
		if _, err := ctx.WriteWithStatus(http.StatusNotAcceptable, constant.Default406Body); err != nil {
			logger.Errorf("WriteWithStatus fail: %v", err)
		}
		ctx.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
		e := errors.Errorf("Requested API %s %s does not online", req.Method, req.URL.Path)
		logger.Debug(e.Error())
		ctx.Err = e
		ctx.Abort()
		return
	}
	ctx.API(api)
	ctx.Next()
}

func (f *Filter) GetApiService() api.APIDiscoveryService {
	return f.apiService
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

var _ filter.HttpFilterFactory = new(Filter)

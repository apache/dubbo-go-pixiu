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
	"context"
	dubboConstant "dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo3"
	"net/http"
	"reflect"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	tpconst "github.com/dubbogo/triple/pkg/common/constant"
	tripleConstant "github.com/dubbogo/triple/pkg/common/constant"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	dubbo2 "github.com/apache/dubbo-go-pixiu/pkg/context/dubbo"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/network/dubboproxy/filter/apiconfig/api"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
	fc "github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/router"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.DubboApiConfigFilter
)

func init() {
	filter.RegisterDubboFilterPlugin(&Plugin{})
}

type (
	Plugin struct {
	}

	// Config
	Config struct {
		Protocol string `default:"protocol" yaml:"protocol" json:"protocol" mapstructure:"protocol"`
	}

	FilterFactory struct {
		cfg        *ApiConfigConfig
		apiService api.APIDiscoveryService
	}

	Filter struct {
		apiService api.APIDiscoveryService
		Config *Config
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

// CreateFilter return the filter callback
func (p Plugin) CreateFilter(config interface{}) (filter.DubboFilter, error) {
	cfg := config.(*Config)
	return Filter{Config: cfg}, nil
}

// Config Expose the config so that Filter Manger can inject it, so it must be a pointer
func (p Plugin) Config() interface{} {
	return &Config{}
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

// Handle handle rpc invocation
func (f Filter) Handle(ctx *dubbo2.RpcContext) filter.FilterStatus {
	switch f.Config.Protocol {
	case dubbo.DUBBO:
		return f.sendDubboRequest(ctx)
	case tripleConstant.TRIPLE:
		return f.sendTripleRequest(ctx)
	default:
		return filter.Stop
	}
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

func (f Filter) sendDubboRequest(ctx *dubbo2.RpcContext) filter.FilterStatus {

	ra := ctx.Route
	clusterName := ra.Cluster
	clusterManager := server.GetClusterManager()
	endpoint := clusterManager.PickEndpoint(clusterName)
	if endpoint == nil {
		ctx.SetError(errors.Errorf("Requested dubbo rpc invocation endpoint not found"))
		return filter.Stop
	}

	invoc := ctx.RpcInvocation
	interfaceKey, _ := invoc.GetAttachment(dubboConstant.InterfaceKey)
	groupKey, _ := invoc.GetAttachment(dubboConstant.GroupKey)
	versionKey, _ := invoc.GetAttachment(dubboConstant.VersionKey)

	url, err := common.NewURL(endpoint.Address.GetAddress(),
		common.WithProtocol(dubbo.DUBBO), common.WithParamsValue(dubboConstant.SerializationKey, dubboConstant.Hessian2Serialization),
		common.WithParamsValue(dubboConstant.GenericFilterKey, "true"),
		common.WithParamsValue(dubboConstant.InterfaceKey, interfaceKey),
		common.WithParamsValue(dubboConstant.ReferenceFilterKey, "generic,filter"),
		// dubboAttachment must contains group and version info
		common.WithParamsValue(dubboConstant.GroupKey, groupKey),
		common.WithParamsValue(dubboConstant.VersionKey, versionKey),
		common.WithPath(interfaceKey),
	)
	if err != nil {
		ctx.SetError(err)
		return filter.Stop
	}

	dubboProtocol := dubbo.NewDubboProtocol()

	// TODO: will print many Error when failed to connect server
	invoker := dubboProtocol.Refer(url)
	if invoker == nil {
		ctx.SetError(errors.Errorf("can't connect to upstream server %s with address %s", endpoint.Name, endpoint.Address.GetAddress()))
		return filter.Stop
	}
	var resp interface{}
	invoc.SetReply(&resp)

	invCtx := context.Background()
	result := invoker.Invoke(invCtx, invoc)
	result.SetAttachments(invoc.Attachments())
	if result.Error() != nil {
		ctx.SetError(result.Error())
		return filter.Stop
	}

	value := reflect.ValueOf(result.Result())
	result.SetResult(value.Elem().Interface())
	ctx.SetResult(result)
	return filter.Continue
}

func (f Filter) sendTripleRequest(ctx *dubbo2.RpcContext) filter.FilterStatus {
	ra := ctx.Route
	clusterName := ra.Cluster
	clusterManager := server.GetClusterManager()
	endpoint := clusterManager.PickEndpoint(clusterName)
	if endpoint == nil {
		ctx.SetError(errors.Errorf("Requested dubbo rpc invocation endpoint not found"))
		return filter.Stop
	}

	invoc := ctx.RpcInvocation
	path := invoc.GetAttachmentInterface(dubboConstant.PathKey).(string)
	// create URL from RpcInvocation
	url, err := common.NewURL(endpoint.Address.GetAddress(),
		common.WithProtocol(tpconst.TRIPLE), common.WithParamsValue(dubboConstant.SerializationKey, dubboConstant.Hessian2Serialization),
		common.WithParamsValue(dubboConstant.GenericFilterKey, "true"),
		common.WithParamsValue(dubboConstant.AppVersionKey, "3.0.0"),
		common.WithParamsValue(dubboConstant.InterfaceKey, path),
		common.WithParamsValue(dubboConstant.ReferenceFilterKey, "generic,filter"),
		common.WithPath(path),
	)
	if err != nil {
		ctx.SetError(err)
		return filter.Stop
	}

	invoker, err := dubbo3.NewDubboInvoker(url)
	if err != nil {
		ctx.SetError(err)
		return filter.Stop
	}

	var resp interface{}
	invoc.SetReply(&resp)
	invCtx := context.Background()
	result := invoker.Invoke(invCtx, invoc)

	if result.Error() != nil {
		ctx.SetError(result.Error())
		return filter.Stop
	}

	// when upstream server down, the result and error in result are both nil
	if result.Result() == nil {
		ctx.SetError(errors.New("result from upstream server is nil"))
		return filter.Stop
	}

	result.SetAttachments(invoc.Attachments())
	value := reflect.ValueOf(result.Result())
	result.SetResult(value.Elem().Interface())
	ctx.SetResult(result)
	return filter.Continue
}
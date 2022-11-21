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
	"context"
	"reflect"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	dubboConstant "dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo3"
	tpconst "github.com/dubbogo/triple/pkg/common/constant"
	tripleConstant "github.com/dubbogo/triple/pkg/common/constant"
	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	dubbo2 "github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/dubbo"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server"
)

const (
	Kind = constant.DubboProxyFilter
)

func init() {
	filter.RegisterDubboFilterPlugin(&Plugin{})
}

type (
	// Plugin dubbo triple plugin
	Plugin struct {
	}

	// Config
	Config struct {
		Protocol string `default:"protocol" yaml:"protocol" json:"protocol" mapstructure:"protocol"`
	}

	// Filter dubbo triple filter
	Filter struct {
		Config *Config
	}
)

// Kind the filter kind
func (p Plugin) Kind() string {
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

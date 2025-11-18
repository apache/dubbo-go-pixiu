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

package dubboproxy

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"os"
	"reflect"
	"strings"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	dubboConstant "dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"
	"dubbo.apache.org/dubbo-go/v3/protocol/invocation"

	hessian "github.com/apache/dubbo-go-hessian2"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	pixiuHttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

const (
	// Kind is the kind of plugin.
	Kind = constant.HTTPDirectDubboProxyFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (

	// Plugin is http to dubbo with direct generic call filter plugin.
	Plugin struct{}

	// FilterFactory is http to dubbo with direct generic call filter instance
	FilterFactory struct {
		cfg *Config
	}

	// Filter http to dubbo with direct generic call
	Filter struct{}

	// Config http to dubbo with direct generic call config
	Config struct{}
)

// Kind return plugin kind
func (p *Plugin) Kind() string {
	return Kind
}

// CreateFilterFactory create filter factory instance
func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

// Config return filter facotry config, now is empty
func (factory *FilterFactory) Config() any {
	return factory.cfg
}

// Apply init filter factory, now is empty
func (factory *FilterFactory) Apply() error {
	return nil
}

// PrepareFilterChain prepare filter chain
func (factory *FilterFactory) PrepareFilterChain(ctx *pixiuHttp.HttpContext, chain filter.FilterChain) error {
	f := &Filter{}
	chain.AppendDecodeFilters(f)
	return nil
}

// Decode handle http request to dubbo direct generic call and return http response
func (f *Filter) Decode(hc *pixiuHttp.HttpContext) filter.FilterStatus {
	rEntry := hc.GetRouteEntry()
	if rEntry == nil {
		logger.Info("[dubbo-go-pixiu] http not match route")
		errResp := pixiuHttp.RouteNotFound.New()
		hc.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}
	logger.Debugf("[dubbo-go-pixiu] client choose endpoint from cluster :%v", rEntry.Cluster)

	clusterName := rEntry.Cluster
	clusterManager := server.GetClusterManager()
	endpoint := clusterManager.PickEndpoint(clusterName, hc)
	if endpoint == nil {
		logger.Info("[dubbo-go-pixiu] cluster not found endpoint")
		errResp := pixiuHttp.ServiceUnavailable.WithError(errors.New("endpoint not found"))
		hc.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	// http://host/{application}/{service}/{method} or https://host/{application}/{service}/{method}
	rawPath := hc.Request.URL.Path
	rawPath = strings.Trim(rawPath, "/")
	splits := strings.Split(rawPath, "/")

	if len(splits) != 3 {
		logger.Info("[dubbo-go-pixiu] http path pattern error. path pattern should be http://127.0.0.1/{application}/{service}/{method}")
		errResp := pixiuHttp.BadRequest.WithError(errors.New("http path pattern error"))
		hc.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	service := splits[1]
	method := splits[2]
	interfaceKey := service

	groupKey := hc.Request.Header.Get(constant.DubboGroup)
	versionKey := hc.Request.Header.Get(constant.DubboServiceVersion)
	types := hc.Request.Header.Get(constant.DubboServiceMethodTypes)

	rawBody, err := io.ReadAll(hc.Request.Body)
	if err != nil {
		logger.Infof("[dubbo-go-pixiu] read request body error %v", err)
		errResp := pixiuHttp.BadRequest.WithError(fmt.Errorf("read request body: %w", err))
		hc.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	var body any
	if err := json.Unmarshal(rawBody, &body); err != nil {
		logger.Infof("[dubbo-go-pixiu] unmarshal request body error %v", err)
		errResp := pixiuHttp.BadRequest.WithError(fmt.Errorf("unmarshal request body: %w", err))
		hc.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	inIArr := make([]any, 3)
	inVArr := make([]reflect.Value, 3)
	inIArr[0] = method

	var (
		typesList  []string
		valuesList []hessian.Object
	)

	if types != "" {
		typesList = strings.Split(types, ",")
	}

	values := body
	if _, ok := values.([]any); ok {
		for _, v := range values.([]any) {
			valuesList = append(valuesList, v)
		}
	} else {
		valuesList = append(valuesList, values)
	}

	inIArr[1] = typesList
	inIArr[2] = valuesList

	inVArr[0] = reflect.ValueOf(inIArr[0])
	inVArr[1] = reflect.ValueOf(inIArr[1])
	inVArr[2] = reflect.ValueOf(inIArr[2])

	invoc := invocation.NewRPCInvocationWithOptions(invocation.WithMethodName("$invoke"),
		invocation.WithArguments(inIArr),
		invocation.WithParameterValues(inVArr))

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
		logger.Infof("[dubbo-go-pixiu] newURL error %v", err)
		errResp := pixiuHttp.BadGateway.WithError(fmt.Errorf("newURL error: %w", err))
		hc.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	dubboProtocol := dubbo.NewDubboProtocol()

	// TODO: will print many Error when failed to connect server
	invoker := dubboProtocol.Refer(url)
	if invoker == nil {
		logger.Info("[dubbo-go-pixiu] dubbo protocol refer error")
		errResp := pixiuHttp.BadGateway.WithError(fmt.Errorf("upstream service error"))
		hc.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	var resp any
	invoc.SetReply(&resp)

	invCtx, cancel := context.WithTimeout(context.Background(), hc.Timeout)
	defer cancel()
	result := invoker.Invoke(invCtx, invoc)
	result.SetAttachments(invoc.Attachments())

	if result.Error() != nil {
		err := result.Error()
		logger.Debugf("[dubbo-go-pixiu] invoke result error %v", err)
		// Prefer reliable timeout detection over substring matching
		if errors.Is(err, context.DeadlineExceeded) || os.IsTimeout(err) {
			errResp := pixiuHttp.GatewayTimeout.WithError(fmt.Errorf("upstream timeout: %w", err))
			hc.SendLocalReply(errResp.Status, errResp.ToJSON())
			return filter.Stop
		}
		errResp := pixiuHttp.BadGateway.WithError(fmt.Errorf("invoke error: %w", err))
		hc.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	value := reflect.ValueOf(result.Result())
	result.SetResult(value.Elem().Interface())
	hc.SourceResp = resp
	invoker.Destroy()
	// response write in hcm
	return filter.Continue
}

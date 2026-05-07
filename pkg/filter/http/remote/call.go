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

package remote

import (
	"errors"
	"fmt"
	"os"
	"strconv"
	"strings"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/client/dubbo"
	clienthttp "github.com/apache/dubbo-go-pixiu/pkg/client/http"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/router"
)

const (
	OPEN = iota
	CLOSE
	ALL
)

const (
	Kind = constant.HTTPDubboProxyFilter
)

var (
	initDubboClient = dubbo.InitDefaultDubboClient
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	mockLevel int8

	Plugin struct {
	}

	FilterFactory struct {
		conf *filterConfig
	}

	Filter struct {
		conf        filterConfig
		dubboClient dubbo.DubboClient
	}

	filterConfig struct {
		Level            mockLevel               `yaml:"level,omitempty" json:"level,omitempty"`
		DubboProxyConfig *dubbo.DubboProxyConfig `yaml:"dubboProxyConfig,omitempty" json:"dubboProxyConfig,omitempty"`
	}

	mockResponse struct {
		Message string `json:"message"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{conf: &filterConfig{}}, nil
}

func (factory *FilterFactory) Config() any {
	return factory.conf
}

func (factory *FilterFactory) Apply() error {
	mock := 1
	mockStr := os.Getenv(constant.EnvMock)
	if len(mockStr) > 0 {
		i, err := strconv.Atoi(mockStr)
		if err == nil {
			mock = i
		}
	}
	level := mockLevel(mock)
	if level < 0 || level > 2 {
		level = CLOSE
	}
	factory.conf.Level = level
	// must init it at apply function
	if factory.conf.DubboProxyConfig == nil {
		return errors.New("expect the dubboProxyConfig config the registries")
	}
	if factory.conf.DubboProxyConfig.AutoResolve != nil {
		return errors.New("dubboProxyConfig.auto_resolve is no longer supported; remove it and configure integrationRequest explicitly in the API definition")
	}
	initDubboClient(factory.conf.DubboProxyConfig)
	return nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	f := &Filter{
		conf:        *factory.conf,
		dubboClient: dubbo.SingletonDubboClient(),
	}
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(c *contexthttp.HttpContext) filter.FilterStatus {
	api := c.GetAPI()

	if (f.conf.Level == OPEN && api.Mock) || (f.conf.Level == ALL) {
		c.SourceResp = &mockResponse{Message: "mock success"}
		return filter.Continue
	}

	typ := api.IntegrationRequest.RequestType
	switch strings.ToLower(typ) {
	case constant.DubboRequest, constant.TripleRequest:
		return f.callDubbo(c, *api)
	case constant.HTTPRequest:
		return f.callHTTP(c, *api)
	default:
		panic(errors.New("not support"))
	}
}

func (f *Filter) callHTTP(c *contexthttp.HttpContext, api router.API) filter.FilterStatus {
	cli, err := f.matchHTTPClient(api.IntegrationRequest.RequestType)
	if err != nil {
		panic(err)
	}

	req := client.NewReq(c.Request.Context(), c.Request, api)
	req.Timeout = c.Timeout
	resp, err := cli.Call(req)
	if err != nil {
		return f.handleClientError(c, err)
	}

	logger.Debugf("[dubbo-go-pixiu] client call resp: %v", resp)

	c.SourceResp = resp
	return filter.Continue
}

func (f *Filter) callDubbo(c *contexthttp.HttpContext, api router.API) filter.FilterStatus {
	// BuildOutbound keeps HTTP mapping details out of the Dubbo client.
	outbound, err := (&DubboHandler{}).BuildOutbound(c.Request, api)
	if err != nil {
		return f.handleClientError(c, err)
	}
	outbound.Timeout = c.Timeout

	resp, err := f.dubboClient.Call(c.Request.Context(), outbound)
	if err != nil {
		return f.handleClientError(c, err)
	}

	logger.Debugf("[dubbo-go-pixiu] client call resp: %v", resp)

	c.SourceResp = resp
	return filter.Continue
}

func (f *Filter) handleClientError(c *contexthttp.HttpContext, err error) filter.FilterStatus {
	logger.Errorf("[dubbo-go-pixiu] client call err: %v!", err)
	if strings.Contains(strings.ToLower(err.Error()), "timeout") {
		errResp := contexthttp.GatewayTimeout.WithError(fmt.Errorf("client timeout: %w", err))
		c.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}
	errResp := contexthttp.InternalError.WithError(fmt.Errorf("client call error: %w", err))
	c.SendLocalReply(errResp.Status, errResp.ToJSON())
	return filter.Stop
}

func (f *Filter) matchHTTPClient(typ string) (client.Client, error) {
	switch strings.ToLower(typ) {
	case constant.HTTPRequest:
		return clienthttp.SingletonHTTPClient(), nil
	default:
		return nil, errors.New("not support")
	}
}

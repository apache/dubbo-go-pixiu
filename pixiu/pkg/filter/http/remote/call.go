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
	"net/http"
	"os"
	"strconv"
	"strings"
)

import (
	apiConf "github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/router"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/client/dubbo"
	clienthttp "github.com/apache/dubbo-go-pixiu/pixiu/pkg/client/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/client/triple"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

const (
	OPEN = iota
	CLOSE
	ALL
)

const (
	Kind = constant.HTTPDubboProxyFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	mockLevel int8

	Plugin struct {
	}

	FilterFactory struct {
		conf *config
	}

	Filter struct {
		conf config
	}

	config struct {
		Level mockLevel               `yaml:"level,omitempty" json:"level,omitempty"`
		Dpc   *dubbo.DubboProxyConfig `yaml:"dubboProxyConfig,omitempty" json:"dubboProxyConfig,omitempty"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{conf: &config{}}, nil
}

func (factory *FilterFactory) Config() interface{} {
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
	dubbo.InitDefaultDubboClient(factory.conf.Dpc)
	triple.InitDefaultTripleClient()
	return nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	f := &Filter{conf: *factory.conf}
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(c *contexthttp.HttpContext) filter.FilterStatus {
	if f.conf.Dpc != nil && f.conf.Dpc.AutoResolve {
		if err := f.resolve(c); err != nil {
			c.SendLocalReply(http.StatusInternalServerError, []byte(fmt.Sprintf("auto resolve err: %s", err)))
			return filter.Stop
		}
	}

	api := c.GetAPI()

	if (f.conf.Level == OPEN && api.Mock) || (f.conf.Level == ALL) {
		c.SourceResp = &contexthttp.ErrResponse{
			Message: "mock success",
		}
		return filter.Continue
	}

	typ := api.Method.IntegrationRequest.RequestType

	cli, err := matchClient(typ)
	if err != nil {
		panic(err)
	}

	req := client.NewReq(c.Request.Context(), c.Request, *api)
	req.Timeout = c.Timeout
	resp, err := cli.Call(req)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] client call err:%v!", err)
		c.SendLocalReply(http.StatusInternalServerError, []byte(fmt.Sprintf("client call err: %s", err)))
		return filter.Stop
	}

	logger.Debugf("[dubbo-go-pixiu] client call resp:%v", resp)

	c.SourceResp = resp
	return filter.Continue
}

func matchClient(typ apiConf.RequestType) (client.Client, error) {
	switch strings.ToLower(string(typ)) {
	case string(apiConf.DubboRequest):
		return dubbo.SingletonDubboClient(), nil
	// todo @(laurence) add triple to apiConf
	case "triple":
		return triple.SingletonTripleClient(), nil
	case string(apiConf.HTTPRequest):
		return clienthttp.SingletonHTTPClient(), nil
	default:
		return nil, errors.New("not support")
	}
}

func (f *Filter) resolve(ctx *contexthttp.HttpContext) error {
	// method must be post
	req := ctx.Request
	if req.Method != http.MethodPost {
		return errors.New("http request method must be post when trying to auto resolve")
	}
	// header must has x-dubbo-http1.1-dubbo-version to declare using auto resolve rule
	version := req.Header.Get(constant.DubboHttpDubboVersion)
	if version == "" {
		return errors.New("http request must has x-dubbo-http1.1-dubbo-version header when trying to auto resolve")
	}

	// http://host/{application}/{service}/{method} or https://host/{application}/{service}/{method}
	rawPath := req.URL.Path
	rawPath = strings.Trim(rawPath, "/")
	splits := strings.Split(rawPath, "/")
	if len(splits) != 3 {
		return errors.New("http request path must meet {application}/{service}/{method} format when trying to auto resolve")
	}

	integrationRequest := apiConf.IntegrationRequest{}
	resolveProtocol := req.Header.Get(constant.DubboServiceProtocol)
	if resolveProtocol == string(apiConf.HTTPRequest) {
		integrationRequest.RequestType = apiConf.HTTPRequest
	} else if resolveProtocol == string(apiConf.DubboRequest) {
		integrationRequest.RequestType = apiConf.DubboRequest
	} else if resolveProtocol == "triple" {
		integrationRequest.RequestType = "triple"
	} else {
		return errors.New("http request has unknown protocol in x-dubbo-service-protocol when trying to auto resolve")
	}

	dubboBackendConfig := apiConf.DubboBackendConfig{}
	dubboBackendConfig.Version = req.Header.Get(constant.DubboServiceVersion)
	dubboBackendConfig.Group = req.Header.Get(constant.DubboGroup)
	integrationRequest.DubboBackendConfig = dubboBackendConfig

	defaultMappingParams := []apiConf.MappingParam{
		{
			Name:  "requestBody.values",
			MapTo: "opt.values",
		}, {
			Name:  "requestBody.types",
			MapTo: "opt.types",
		}, {
			Name:  "uri.application",
			MapTo: "opt.application",
		}, {
			Name:  "uri.interface",
			MapTo: "opt.interface",
		}, {
			Name:  "uri.method",
			MapTo: "opt.method",
		},
	}
	integrationRequest.MappingParams = defaultMappingParams

	method := apiConf.Method{
		Enable:   true,
		Mock:     false,
		HTTPVerb: http.MethodPost,
	}
	method.IntegrationRequest = integrationRequest

	inboundRequest := apiConf.InboundRequest{}
	inboundRequest.RequestType = apiConf.HTTPRequest
	method.InboundRequest = inboundRequest

	api := router.API{}
	api.URLPattern = "/:application/:interface/:method"
	api.Method = method
	ctx.API(api)
	return nil
}

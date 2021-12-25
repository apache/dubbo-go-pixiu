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
	"net/http"
	"os"
	"strconv"
	"strings"
)

import (
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
	apiConf "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	fconf "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/router"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/client/dubbo"
	clienthttp "github.com/apache/dubbo-go-pixiu/pkg/client/http"
	"github.com/apache/dubbo-go-pixiu/pkg/client/triple"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	open = iota
	close
	all
)

const (
	Kind = constant.HTTPDubboProxyFilter
)

const (
	x_dubbo_http_dubbo_version = "x-dubbo-http1.1-dubbo-version"
	x_dubbo_service_protocol   = "x-dubbo-service-protocol"
	x_dubbo_service_version    = "x-dubbo-service-version"
	x_dubbo_group              = "x-dubbo-service-group"
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	mockLevel int8

	Plugin struct {
	}

	Filter struct {
		conf *config
	}

	config struct {
		Level mockLevel               `yaml:"level,omitempty" json:"level,omitempty"`
		Dpc   *dubbo.DubboProxyConfig `yaml:"dubboProxyConfig,omitempty" json:"dubboProxyConfig,omitempty"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter() (filter.HttpFilter, error) {
	return &Filter{conf: &config{}}, nil
}

func (f *Filter) Config() interface{} {
	return f.conf
}

func (f *Filter) Apply() error {
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
		level = close
	}
	f.conf.Level = level
	// must init it at apply function
	dubbo.InitDefaultDubboClient(f.conf.Dpc)
	triple.InitDefaultTripleClient()
	return nil
}

func (f *Filter) PrepareFilterChain(ctx *contexthttp.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func NewDefaultAPI() router.API {
	api := router.API{}
	return api
}

func NewDefaultMethod() fconf.Method {
	method := fconf.Method{}
	method.Enable = true
	method.Mock = false
	method.HTTPVerb = http.MethodPost
	return method
}

func (f *Filter) resolve(ctx *contexthttp.HttpContext) error {
	// method must be post
	req := ctx.Request
	if req.Method != http.MethodPost {
		return errors.New("http request method must be post when trying to auto resolve")
	}
	// header must has x-dubbo-http1.1-dubbo-version to declare using auto resolve rule
	version := req.Header.Get(x_dubbo_http_dubbo_version)
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

	integrationRequest := fconf.IntegrationRequest{}
	resolve_protocol := req.Header.Get(x_dubbo_service_protocol)
	if resolve_protocol == string(fconf.HTTPRequest) {
		integrationRequest.RequestType = fconf.HTTPRequest
	} else if resolve_protocol == string(fconf.DubboRequest) {
		integrationRequest.RequestType = fconf.DubboRequest
	} else if resolve_protocol == "triple" {
		integrationRequest.RequestType = "triple"
	} else {
		return errors.New("http request has unknown protocol in x-dubbo-protocol when trying to auto resolve")
	}

	dubboBackendConfig := fconf.DubboBackendConfig{}
	dubboBackendConfig.Version = req.Header.Get(x_dubbo_service_version)
	dubboBackendConfig.Group = req.Header.Get(x_dubbo_group)
	integrationRequest.DubboBackendConfig = dubboBackendConfig

	defaultMappingParams := []fconf.MappingParam{
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

	method := NewDefaultMethod()
	method.IntegrationRequest = integrationRequest

	inboundRequest := fconf.InboundRequest{}
	inboundRequest.RequestType = fconf.HTTPRequest
	method.InboundRequest = inboundRequest

	api := NewDefaultAPI()
	api.URLPattern = "/:application/:interface/:method"
	api.Method = method
	ctx.API(api)
	return nil
}

func (f *Filter) Handle(c *contexthttp.HttpContext) {
	if f.conf.Dpc.AutoResolve {
		if err := f.resolve(c); err != nil {
			c.Err = err
			return
		}
	}

	api := c.GetAPI()

	if (f.conf.Level == open && api.Mock) || (f.conf.Level == all) {
		c.SourceResp = &contexthttp.ErrResponse{
			Message: "mock success",
		}
		c.Next()
		return
	}

	typ := api.Method.IntegrationRequest.RequestType

	cli, err := matchClient(typ)
	if err != nil {
		c.Err = err
		c.Next()
		return
	}

	req := client.NewReq(c.Request.Context(), c.Request, *api)
	resp, err := cli.Call(req)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] client call err:%v!", err)
		c.Err = err
		c.Next()
		return
	}

	logger.Debugf("[dubbo-go-pixiu] client call resp:%v", resp)

	c.SourceResp = resp
	// response write in response filter.
	c.Next()
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

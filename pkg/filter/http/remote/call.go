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
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/client/dubbo"
	clienthttp "github.com/apache/dubbo-go-pixiu/pkg/client/http"
	"github.com/apache/dubbo-go-pixiu/pkg/client/triple"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/http/remote/resolver"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
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
		conf     config
		resolver resolver.Resolver
	}

	config struct {
		Level            mockLevel               `yaml:"level,omitempty" json:"level,omitempty"`
		DubboProxyConfig *dubbo.DubboProxyConfig `yaml:"dubboProxyConfig,omitempty" json:"dubboProxyConfig,omitempty"`
		// Resolver is the Resolver to resolve HTTP requests to Dubbo services.
		Resolver string `yaml:"resolver,omitempty" json:"resolver,omitempty" default:"StandardDubboResolver"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{conf: &config{}}, nil
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
	dubbo.InitDefaultDubboClient(factory.conf.DubboProxyConfig)
	triple.InitDefaultTripleClient(factory.conf.DubboProxyConfig.Protoset)
	return nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	r, err := resolver.GetResolver(factory.conf.Resolver)

	if err != nil {
		logger.Errorf("get resolver fail %s", err.Error())
	}

	f := &Filter{
		conf:     *factory.conf,
		resolver: r,
	}
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(c *contexthttp.HttpContext) filter.FilterStatus {
	if f.conf.DubboProxyConfig != nil && f.conf.DubboProxyConfig.AutoResolve {
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

	cli, err := f.matchClient(typ)
	if err != nil {
		panic(err)
	}

	req := client.NewReq(c.Request.Context(), c.Request, *api)
	req.Timeout = c.Timeout
	resp, err := cli.Call(req)
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] client call err: %v!", err)
		if strings.Contains(strings.ToLower(err.Error()), "timeout") {
			c.SendLocalReply(http.StatusGatewayTimeout, []byte(fmt.Sprintf("client call timeout err: %s", err)))
			return filter.Stop
		}
		c.SendLocalReply(http.StatusInternalServerError, []byte(fmt.Sprintf("client call err: %s", err)))
		return filter.Stop
	}

	logger.Debugf("[dubbo-go-pixiu] client call resp: %v", resp)

	c.SourceResp = resp
	return filter.Continue
}

func (f *Filter) matchClient(typ apiConf.RequestType) (client.Client, error) {
	switch strings.ToLower(string(typ)) {
	case string(apiConf.DubboRequest):
		return dubbo.SingletonDubboClient(), nil
	// todo @(laurence) add triple to apiConf
	case "triple":
		return triple.SingletonTripleClient(f.conf.DubboProxyConfig.Protoset), nil
	case string(apiConf.HTTPRequest):
		return clienthttp.SingletonHTTPClient(), nil
	default:
		return nil, errors.New("not support")
	}
}

// Resolve is the function calls resolver.Resolve.
func (f *Filter) resolve(ctx *contexthttp.HttpContext) error {
	api, err := f.resolver.Resolve(ctx)
	if err != nil {
		logger.Warnf("[dubbo-go-pixiu] resolver err: %v", err)
		return err
	}
	if api != nil {
		// Resolver successfully processed the request.
		ctx.API(*api)
		return nil
	}

	// If no resolver could handle the request, return a generic error.
	// This maintains the original behavior of failing if auto-resolve conditions aren't met.
	err = errors.New("http request cannot be resolved to a Dubbo service")
	logger.Errorf("[dubbo-go-pixiu] resolver err: %v", err)
	return err
}

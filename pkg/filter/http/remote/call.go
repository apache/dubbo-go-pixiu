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
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"

	apiConf "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config"
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
		level = close
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
	api := c.GetAPI()

	if (f.conf.Level == open && api.Mock) || (f.conf.Level == all) {
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

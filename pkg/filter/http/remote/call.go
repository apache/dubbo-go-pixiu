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
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
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
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)


const (
	open = iota
	close
	all
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.MetricFilter
)

func init() {
	extension.RegisterHttpFilter(&Plugin{})
}

type (
	mockLevel int8

	Plugin struct {

	}

	clientFilter struct {
		conf *config
	}

	config struct {
		Level mockLevel `yaml:"level,omitempty" json:"level,omitempty"`
		dpc * DubboProxyConfig `yaml:"dubboProxyConfig,omitempty" json:"dubboProxyConfig,omitempty"`
	}
)

func (ap *Plugin) Kind() string {
	return Kind
}

func (ap *Plugin) CreateFilter() (extension.HttpFilter, error) {
	return &clientFilter{conf: &config{}},nil
}

func (f *clientFilter) Config() interface{} {
	return f.conf
}

func (f *clientFilter) Apply() error {
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
	dubbo.InitDefaultDubboClient(f.conf.dpc)
	return nil
}

func (f *clientFilter) PrepareFilterChain(ctx *contexthttp.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *clientFilter) Handle(c *contexthttp.HttpContext) {
	api := c.GetAPI()

	if (f.conf.Level == open && api.Mock) || (f.conf.Level == all) {
		c.SourceResp = &extension.ErrResponse{
			Message: "mock success",
		}
		c.Next()
		return
	}

	typ := api.Method.IntegrationRequest.RequestType

	cli, err := matchClient(typ)
	if err != nil {
		c.Err = err
		return
	}

	resp, err := cli.Call(client.NewReq(c.Ctx, c.Request, *api))
	if err != nil {
		logger.Errorf("[dubbo-go-pixiu] client call err:%v!", err)
		c.Err = err
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
	case string(apiConf.HTTPRequest):
		return clienthttp.SingletonHTTPClient(), nil
	default:
		return nil, errors.New("not support")
	}
}

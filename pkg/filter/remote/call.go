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
)

import (
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client"
	"github.com/dubbogo/dubbo-go-proxy/pkg/client/dubbo"
	clienthttp "github.com/dubbogo/dubbo-go-proxy/pkg/client/http"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	selfcontext "github.com/dubbogo/dubbo-go-proxy/pkg/context"
	contexthttp "github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
)

func init() {
	extension.SetFilterFunc(constant.RemoteCallFilter, New(false).Do())
}

// clientFilter is a filter for recover.
type clientFilter struct {
	mock bool
}

// New create timeout filter.
func New(mock bool) filter.Filter {
	return &clientFilter{
		mock: mock,
	}
}

// Do execute clientFilter filter logic
// support: 1 http 2 dubbo 2 http 2 http
func (f *clientFilter) Do() selfcontext.FilterFunc {
	return func(c selfcontext.Context) {
		f.doRemoteCall(c.(*contexthttp.HttpContext))
	}
}

func (f *clientFilter) doRemoteCall(c *contexthttp.HttpContext) {
	api := c.GetAPI()

	if f.mock {
		return
	}

	typ := api.Method.IntegrationRequest.RequestType

	cli, err := matchClient(typ)
	if err != nil {
		c.WriteWithStatus(http.StatusServiceUnavailable, constant.Default503Body)
		c.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueTextPlain)
		return
	}

	resp, err := cli.Call(client.NewReq(c.Ctx, c.Request, *api))

	if err != nil {
		logger.Errorf("[dubbo-go-proxy go] client do err:%v!", err)
		c.WriteErr(err)
		return
	}

	logger.Debugf("[dubbo-go-proxy go] : %v", resp)

	// response write in response filter.
	c.Next()
}

func matchClient(typ config.RequestType) (client.Client, error) {
	switch typ {
	case config.DubboRequest:
		return dubbo.SingletonDubboClient(), nil
	case config.HTTPRequest:
		return clienthttp.SingletonHTTPClient(), nil
	default:
		return nil, errors.New("not support")
	}
}

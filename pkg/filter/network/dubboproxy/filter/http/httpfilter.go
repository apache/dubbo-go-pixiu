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

package http

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	http3 "net/http"
	stdHttp "net/http"
	"net/url"
	"strings"
)

import (
	"dubbo.apache.org/dubbo-go/v3/protocol"
	"dubbo.apache.org/dubbo-go/v3/protocol/dubbo"

	"github.com/pkg/errors"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	dubbo2 "github.com/apache/dubbo-go-pixiu/pkg/context/dubbo"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

const (
	x_dubbo_http_dubbo_version = "x-dubbo-http1.1-dubbo-version"
	x_dubbo_service_protocol   = "x-dubbo-service-protocol"
	x_dubbo_service_version    = "x-dubbo-service-version"
	x_dubbo_group              = "x-dubbo-service-group"
)

const (
	Kind = constant.DubboHttpFilter
)

func init() {
	filter.RegisterDubboFilterPlugin(&Plugin{})
}

type (
	// Plugin dubbo to http transform plugin
	Plugin struct {
	}

	// Config config
	Config struct {
	}

	// Filter dubbo to http transform filter
	Filter struct {
		Config *Config
	}
)

// Kind the filter kind
func (p Plugin) Kind() string {
	return Kind
}

// CreateFilter return the filter
func (p Plugin) CreateFilter(config interface{}) (filter.DubboFilter, error) {
	cfg := config.(*Config)
	return Filter{Config: cfg}, nil
}

// Config Expose the config so that Filter Manger can inject it, so it must be a pointer
func (p Plugin) Config() interface{} {
	return &Config{}
}

// Handle transform rpcInvocation to http
func (f Filter) Handle(ctx *dubbo2.RpcContext) filter.FilterStatus {

	ra := ctx.Route
	clusterName := ra.Cluster
	clusterManager := server.GetClusterManager()
	endpoint := clusterManager.PickEndpoint(clusterName)
	if endpoint == nil {
		ctx.SetError(errors.Errorf("Requested dubbo rpc invocation endpoint not found"))
		return filter.Stop
	}

	var (
		req *http3.Request
		err error
	)

	invoc := ctx.RpcInvocation
	// path's format /{service}/{method}
	interfaceKey := invoc.AttachmentsByKey(constant.InterfaceKey, "")
	// work when invocation is generic
	// when invocation is generic, there are three value in arguments. first is methodName, second is types, third is values
	methodName := invoc.Arguments()[0].(string)
	path := interfaceKey + "/" + methodName

	parsedURL := url.URL{
		Host:   endpoint.Address.GetAddress(),
		Scheme: "http",
		Path:   path,
	}

	body := invoc.Arguments()[2]
	b, _ := json.Marshal(body)
	req, err = http3.NewRequest(stdHttp.MethodPost, parsedURL.String(), strings.NewReader(string(b)))
	if err != nil {
		err := errors.New(fmt.Sprintf("create new request failed: %v", err))
		ctx.SetError(err)
		return filter.Stop
	}

	req.Header.Set(x_dubbo_http_dubbo_version, "1.0.0")
	req.Header.Set(x_dubbo_service_protocol, dubbo.DUBBO)
	req.Header.Set(x_dubbo_service_version, invoc.AttachmentsByKey(constant.VersionKey, ""))
	req.Header.Set(x_dubbo_group, invoc.AttachmentsByKey(constant.GroupKey, ""))

	resp, err := (&http3.Client{}).Do(req)
	if err != nil {
		ctx.SetError(err)
		return filter.Stop
	}

	if resp.StatusCode != stdHttp.StatusOK {
		ctx.SetError(errors.New(fmt.Sprintf("upstream http response status code %d", resp.StatusCode)))
		return filter.Stop
	}

	s, _ := ioutil.ReadAll(resp.Body)
	result := &protocol.RPCResult{}
	result.Rest = string(s)
	ctx.SetResult(result)
	return filter.Continue
}

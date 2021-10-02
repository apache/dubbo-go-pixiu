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

package httpproxy

import (
	"encoding/json"
	"fmt"
	http3 "net/http"
	"net/url"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.HTTPProxyFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is http filter plugin.
	Plugin struct {
	}
	// Filter is http filter instance
	Filter struct {
		cfg       *Config
		transport http3.RoundTripper
	}
	// Config describe the config of Filter
	Config struct{}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter() (filter.HttpFilter, error) {
	return &Filter{cfg: &Config{}, transport: &http3.Transport{}}, nil
}

func (f *Filter) Config() interface{} {
	return f.cfg
}

func (f *Filter) Apply() error {
	return nil
}

func (f *Filter) PrepareFilterChain(ctx *http.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *Filter) Handle(hc *http.HttpContext) {
	rEntry := hc.GetRouteEntry()
	if rEntry == nil {
		panic("no route entry")
	}
	logger.Debugf("[dubbo-go-pixiu] client choose endpoint from cluster :%v", rEntry.Cluster)

	clusterName := rEntry.Cluster
	clusterManager := server.GetClusterManager()
	endpoint := clusterManager.PickEndpoint(clusterName)
	if endpoint == nil {
		bt, _ := json.Marshal(http.ErrResponse{Message: "cluster not found endpoint"})
		hc.SourceResp = bt
		hc.TargetResp = &client.Response{Data: bt}
		hc.WriteJSONWithStatus(http3.StatusServiceUnavailable, bt)
		hc.Abort()
		return
	}

	logger.Debugf("[dubbo-go-pixiu] client choose endpoint :%v", endpoint.Address.GetAddress())
	r := hc.Request

	var errPrefix string
	defer func() {
		if err := recover(); err != nil {
			bt, _ := json.Marshal(http.ErrResponse{Message: fmt.Sprintf("remoteFilterErr: %s: %v", errPrefix, err)})
			hc.SourceResp = bt
			hc.TargetResp = &client.Response{Data: bt}
			hc.WriteJSONWithStatus(http3.StatusServiceUnavailable, bt)
			hc.Abort()
			return
		}
	}()

	var (
		req *http3.Request
		err error
	)

	parsedURL := url.URL{
		Host:     endpoint.Address.GetAddress(),
		Scheme:   "http",
		Path:     r.URL.Path,
		RawQuery: r.URL.RawQuery,
	}

	req, err = http3.NewRequest(r.Method, parsedURL.String(), r.Body)
	if err != nil {
		bt, _ := json.Marshal(http.ErrResponse{Message: fmt.Sprintf("BUG: new request failed: %v", err)})
		hc.SourceResp = bt
		hc.TargetResp = &client.Response{Data: bt}
		hc.WriteJSONWithStatus(http3.StatusInternalServerError, bt)
		hc.Abort()
		return
	}
	req.Header = r.Header

	errPrefix = "do request"
	resp, err := (&http3.Client{Transport: f.transport}).Do(req)
	if err != nil {
		panic(err)
	}
	logger.Debugf("[dubbo-go-pixiu] client call resp:%v", resp)

	hc.SourceResp = resp
	// response write in response filter.
	hc.Next()
}

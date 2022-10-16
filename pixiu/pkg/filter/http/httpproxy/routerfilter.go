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
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/server"
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
	// FilterFactory is http filter instance
	FilterFactory struct {
		cfg *Config
	}
	//Filter
	Filter struct {
		transport http3.RoundTripper
	}
	// Config describe the config of FilterFactory
	Config struct{}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

func (factory *FilterFactory) Config() interface{} {
	return factory.cfg
}

func (factory *FilterFactory) Apply() error {
	return nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
	f := &Filter{transport: &http3.Transport{}}
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(hc *http.HttpContext) filter.FilterStatus {
	rEntry := hc.GetRouteEntry()
	if rEntry == nil {
		panic("no route entry")
	}
	logger.Debugf("[dubbo-go-pixiu] client choose endpoint from cluster :%v", rEntry.Cluster)

	clusterName := rEntry.Cluster
	clusterManager := server.GetClusterManager()
	endpoint := clusterManager.PickEndpoint(clusterName)
	if endpoint == nil {
		logger.Debugf("[dubbo-go-pixiu] cluster not found endpoint")
		bt, _ := json.Marshal(http.ErrResponse{Message: "cluster not found endpoint"})
		hc.SendLocalReply(http3.StatusServiceUnavailable, bt)
		return filter.Stop
	}

	logger.Debugf("[dubbo-go-pixiu] client choose endpoint :%v", endpoint.Address.GetAddress())
	r := hc.Request

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
		hc.SendLocalReply(http3.StatusInternalServerError, bt)
		return filter.Stop
	}
	req.Header = r.Header

	cli := &http3.Client{
		Transport: f.transport,
		Timeout:   hc.Timeout,
	}

	resp, err := cli.Do(req)
	if err != nil {
		hc.SendLocalReply(http3.StatusServiceUnavailable, []byte(err.Error()))
		return filter.Stop
	}
	logger.Debugf("[dubbo-go-pixiu] client call resp:%v", resp)

	hc.SourceResp = resp
	// response write in hcm
	return filter.Continue
}

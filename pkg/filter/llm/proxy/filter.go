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

package proxy

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"net/url"
	"strings"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/common/util"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pkg/server"
)

const (
	Kind = constant.LLMProxyFilter
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
		cfg    *Config
		client http.Client
	}
	Filter struct {
		client http.Client
		scheme string
	}
	// Config describe the config of FilterFactory
	Config struct {
		Timeout             time.Duration `yaml:"timeout" json:"timeout,omitempty"`
		MaxIdleConns        int           `yaml:"maxIdleConns" json:"maxIdleConns,omitempty"`
		MaxIdleConnsPerHost int           `yaml:"maxIdleConnsPerHost" json:"maxIdleConnsPerHost,omitempty"`
		MaxConnsPerHost     int           `yaml:"maxConnsPerHost" json:"maxConnsPerHost,omitempty"`
		Scheme              string        `yaml:"scheme" json:"scheme,omitempty" default:"http"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

func (factory *FilterFactory) Config() any {
	return factory.cfg
}

func (factory *FilterFactory) Apply() error {
	scheme := strings.TrimSpace(strings.ToLower(factory.cfg.Scheme))

	if scheme != "http" && scheme != "https" {
		return fmt.Errorf("%s: scheme must be http or https", Kind)
	}

	factory.cfg.Scheme = scheme

	cfg := factory.cfg
	client := http.Client{
		Timeout: cfg.Timeout,
		Transport: http.RoundTripper(&http.Transport{
			MaxIdleConns:        cfg.MaxIdleConns,
			MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
			MaxConnsPerHost:     cfg.MaxConnsPerHost,
		}),
	}
	factory.client = client
	return nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	//reuse http client
	f := &Filter{factory.client, factory.cfg.Scheme}
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(hc *contexthttp.HttpContext) filter.FilterStatus {
	rEntry := hc.GetRouteEntry()
	if rEntry == nil {
		bt, _ := json.Marshal(contexthttp.ErrResponse{Message: "no route entry"})
		hc.SendLocalReply(http.StatusBadRequest, bt)
		return filter.Stop
	}

	logger.Debugf("[dubbo-go-pixiu] client choose endpoint from cluster: %v", rEntry.Cluster)

	var (
		clusterName    = rEntry.Cluster
		clusterManager = server.GetClusterManager()
		endpoint       = clusterManager.PickEndpoint(clusterName, hc)
	)

	if endpoint == nil {
		logger.Debugf("[dubbo-go-pixiu] cluster not found endpoint")
		bt, _ := json.Marshal(contexthttp.ErrResponse{Message: "cluster not found endpoint"})
		hc.SendLocalReply(http.StatusServiceUnavailable, bt)
		return filter.Stop
	}

	r := hc.Request
	defer r.Body.Close()

	var (
		req  *http.Request
		resp *http.Response
		err  error
	)

	if hc.Request.Body != nil && hc.Request.GetBody == nil {
		bodyBytes, err := io.ReadAll(hc.Request.Body)
		hc.Request.Body.Close()

		if err != nil {
			bt, _ := json.Marshal(contexthttp.ErrResponse{Message: fmt.Sprintf("failed to read request body: %v", err)})
			hc.SendLocalReply(http.StatusInternalServerError, bt)
			return filter.Stop
		}

		hc.Request.Body = io.NopCloser(bytes.NewReader(bodyBytes))
		hc.Request.GetBody = func() (io.ReadCloser, error) {
			return io.NopCloser(bytes.NewReader(bodyBytes)), nil
		}
	}

	logger.Debugf("[dubbo-go-pixiu] client choose endpoint [%s: %v]", endpoint.ID, endpoint.Address.GetAddress())

	// make request
FALLBACK:
	for {
	RETRY:
		for retry := uint(0); retry <= endpoint.LLMMeta.RetryTimes; retry++ {
			req, err = f.assembleRequest(endpoint, r)
			if err != nil {
				logger.Warnf("[dubbo-go-pixiu] client assemble request failed: %v", err)
				break RETRY
			}

			resp, err = f.client.Do(req)
			if err != nil {
				logger.Warnf("[dubbo-go-pixiu] client call endpoint [%s: %v] failed: %v", endpoint.ID, endpoint.Address.GetAddress(), err)
				break RETRY
			}
			if util.IsHTTPRespSuccessful(resp.StatusCode) {
				// If the response is successful, we can break out of the fallback loop.
				break FALLBACK
			}
			// If the response is not successful, we will retry with the next endpoint.
			logger.Debugf("[dubbo-go-pixiu] client retry endpoint [%s: %v]", endpoint.ID, endpoint.Address.GetAddress())
		}

		if !endpoint.LLMMeta.Fallback {
			// If fallback is not enabled, we will break out of the fallback loop.
			break FALLBACK
		}

		endpoint = clusterManager.PickNextEndpoint(clusterName, endpoint.ID)
		if endpoint == nil {
			break FALLBACK
		}

		// If we have a next endpoint, we will retry with the next endpoint.
		logger.Debugf("[dubbo-go-pixiu] client fallback to endpoint [%s: %v]", endpoint.ID, endpoint.Address.GetAddress())
	}

	if err != nil {
		var urlErr *url.Error
		ok := errors.As(err, &urlErr)
		if ok && urlErr.Timeout() {
			hc.SendLocalReply(http.StatusGatewayTimeout, []byte(err.Error()))
			return filter.Stop
		}
		hc.SendLocalReply(http.StatusServiceUnavailable, []byte(err.Error()))
		return filter.Stop
	}

	logger.Debugf("[dubbo-go-pixiu] client call resp:%v", resp)
	hc.SourceResp = resp
	// response write in hcm
	return filter.Continue

}

func (f *Filter) assembleRequest(endpoint *model.Endpoint, r *http.Request) (*http.Request, error) {
	parsedURL := url.URL{
		Host:     endpoint.Address.GetAddress(),
		Scheme:   f.scheme,
		Path:     r.URL.Path,
		RawQuery: r.URL.RawQuery,
	}

	req, err := http.NewRequest(r.Method, parsedURL.String(), r.Body)
	if err != nil {
		return nil, err
	}
	req.Header = r.Header

	return req, nil
}

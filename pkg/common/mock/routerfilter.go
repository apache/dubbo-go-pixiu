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

package mock

import (
	"encoding/json"
	"fmt"
	"net/http"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	// Kind is the kind of Fallback.
	Kind = "dgp.filter.http.sse.httpproxy"
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
	//Filter
	Filter struct {
		client http.Client
	}
	// Config describe the config of FilterFactory
	Config struct {
		Timeout             time.Duration `yaml:"timeout" json:"timeout,omitempty"`
		MaxIdleConns        int           `yaml:"maxIdleConns" json:"maxIdleConns,omitempty"`
		MaxIdleConnsPerHost int           `yaml:"maxIdleConnsPerHost" json:"maxIdleConnsPerHost,omitempty"`
		MaxConnsPerHost     int           `yaml:"maxConnsPerHost" json:"maxConnsPerHost,omitempty"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

func (ff *FilterFactory) Config() any {
	return ff.cfg
}

func (ff *FilterFactory) Apply() error {
	cfg := ff.cfg
	client := http.Client{
		Timeout: cfg.Timeout,
		Transport: http.RoundTripper(&http.Transport{
			MaxIdleConns:        cfg.MaxIdleConns,
			MaxIdleConnsPerHost: cfg.MaxIdleConnsPerHost,
			MaxConnsPerHost:     cfg.MaxConnsPerHost,
		}),
	}
	ff.client = client
	return nil
}

func (ff *FilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	//reuse http client
	f := &Filter{ff.client}
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(hc *contexthttp.HttpContext) filter.FilterStatus {
	r := hc.Request

	var (
		req *http.Request
		err error
	)

	req, err = http.NewRequest(r.Method, r.URL.String(), r.Body)
	if err != nil {
		bt, _ := json.Marshal(contexthttp.ErrResponse{Message: fmt.Sprintf("BUG: new request failed: %v", err)})
		hc.SendLocalReply(http.StatusInternalServerError, bt)
		return filter.Stop
	}
	req.Header = r.Header

	resp, err := f.client.Do(req)

	if err != nil {
		hc.SendLocalReply(http.StatusServiceUnavailable, []byte(err.Error()))
		return filter.Stop
	}

	hc.SourceResp = resp

	logger.Debugf("[dubbo-go-pixiu] client call resp:%v", resp)
	// response write in hcm
	return filter.Continue
}

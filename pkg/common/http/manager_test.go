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
	"bytes"
	"fmt"
	"net/http"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/context/mock"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	DEMO = "dgp.filters.http.demo"
	// Kind is the kind of plugin.
	Kind = DEMO
)

type (
	// Plugin is http filter plugin.
	Plugin struct {
	}
	// HeaderFilter is http filter instance
	DemoFilter struct {
		conf *Config
		str  string
	}
	// Config describe the config of ResponseFilter
	Config struct {
		Foo string `json:"foo,omitempty" yaml:"foo,omitempty"`
		Bar string `json:"bar,omitempty" yaml:"bar,omitempty"`
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter() (filter.HttpFilter, error) {
	return &DemoFilter{conf: &Config{Foo: "default foo", Bar: "default bar"}}, nil
}

func (f *DemoFilter) PrepareFilterChain(ctx *contexthttp.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *DemoFilter) Handle(c *contexthttp.HttpContext) {
	logger.Info(f.str)
}

func (f *DemoFilter) Config() interface{} {
	return f.conf
}

func (f *DemoFilter) Apply() error {
	c := f.conf
	f.str = fmt.Sprintf("%s is drinking in the %s", c.Foo, c.Bar)
	//return the filter func
	return nil
}

func TestCreateHttpConnectionManager(t *testing.T) {
	filter.RegisterHttpFilter(&Plugin{})

	hcmc := model.HttpConnectionManager{
		RouteConfig: model.RouteConfiguration{
			Routes: []*model.Router{
				{
					ID: "1",
					Match: model.RouterMatch{
						Prefix: "/api/v1",
						Methods: []string{
							"POST",
						},
						Path:  "",
						Regex: "",
						Headers: []model.HeaderMatcher{
							{Name: "X-DGP-WAY",
								Values: []string{"Dubbo"},
								Regex:  false,
							},
						},
					},
					Route: model.RouteAction{
						Cluster:                     "test_dubbo",
						ClusterNotFoundResponseCode: 505,
					},
				},
			},
			Dynamic: false,
		},
		HTTPFilters: []*model.HTTPFilter{
			{
				Name:   DEMO,
				Config: nil,
			},
		},
		ServerName:        "test_http_dubbo",
		GenerateRequestID: false,
		IdleTimeoutStr:    "100",
	}

	hcm := CreateHttpConnectionManager(&hcmc, nil)
	assert.Equal(t, len(hcm.filterManager.GetFilters()), 1)
	request, err := http.NewRequest("POST", "http://www.dubbogopixiu.com/api/v1?name=tc", bytes.NewReader([]byte("{\"id\":\"12345\"}")))
	assert.NoError(t, err)
	c := mock.GetMockHTTPContext(request)
	hcm.addFilter(c)
	assert.Equal(t, len(c.Filters), 1)
	err = hcm.findRoute(c)
	assert.NoError(t, err)
	err = hcm.OnData(c)
	assert.NoError(t, err)
}

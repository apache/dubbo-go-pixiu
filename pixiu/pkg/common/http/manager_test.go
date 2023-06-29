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
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/router/trie"
	contexthttp "github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/mock"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
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
	DemoFilterFactory struct {
		conf *Config
	}
	DemoFilter struct {
		str string
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

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &DemoFilterFactory{conf: &Config{Foo: "default foo", Bar: "default bar"}}, nil
}

func (f *DemoFilter) Decode(ctx *contexthttp.HttpContext) filter.FilterStatus {
	logger.Info("decode phase: ", f.str)

	runes := []rune(f.str)
	for i := 0; i < len(runes)/2; i += 1 {
		runes[i], runes[len(runes)-1-i] = runes[len(runes)-1-i], runes[i]
	}
	f.str = string(runes)

	return filter.Continue
}

func (f *DemoFilter) Encode(ctx *contexthttp.HttpContext) filter.FilterStatus {
	logger.Info("encode phase: ", f.str)
	return filter.Continue
}

func (f *DemoFilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	c := f.conf
	str := fmt.Sprintf("%s is drinking in the %s", c.Foo, c.Bar)
	filter := &DemoFilter{str: str}

	chain.AppendDecodeFilters(filter)
	chain.AppendEncodeFilters(filter)
	return nil
}

func (f *DemoFilterFactory) Config() interface{} {
	return f.conf
}

func (f *DemoFilterFactory) Apply() error {
	return nil
}

func TestCreateHttpConnectionManager(t *testing.T) {
	filter.RegisterHttpFilter(&Plugin{})

	hcmc := model.HttpConnectionManagerConfig{
		RouteConfig: model.RouteConfiguration{
			RouteTrie: trie.NewTrieWithDefault("POST/api/v1/**", model.RouteAction{
				Cluster:                     "test_dubbo",
				ClusterNotFoundResponseCode: 505,
			}),
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

	hcm := CreateHttpConnectionManager(&hcmc)
	assert.Equal(t, len(hcm.filterManager.GetFactory()), 1)
	request, err := http.NewRequest("POST", "http://www.dubbogopixiu.com/api/v1?name=tc", bytes.NewReader([]byte("{\"id\":\"12345\"}")))
	assert.NoError(t, err)
	request.Header = map[string][]string{
		"X-Dgp-Way": []string{"Dubbo"},
	}
	assert.NoError(t, err)
	c := mock.GetMockHTTPContext(request)
	err = hcm.findRoute(c)
	assert.NoError(t, err)
	err = hcm.Handle(c)
	assert.NoError(t, err)
}

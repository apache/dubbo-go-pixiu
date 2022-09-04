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

package proxyrewrite

import (
	"net/http"
	"testing"
)

import (
	"github.com/go-playground/assert/v2"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/mock"
)

func TestDecode(t *testing.T) {
	url := "/user-service/query"

	request, _ := http.NewRequest("Get", url, nil)
	ctx := mock.GetMockHTTPContext(request)

	factory := &FilterFactory{cfg: &Config{UriRegex: []string{"^/([^/]*)/(.*)$", "/$2"}}}
	chain := filter.NewDefaultFilterChain()
	_ = factory.PrepareFilterChain(ctx, chain)
	chain.OnDecode(ctx)

	assert.Equal(t, ctx.GetUrl(), "/query")
}

func TestConfigUpdate(t *testing.T) {
	url := "/user-service/query"

	request, _ := http.NewRequest("Get", url, nil)
	ctx := mock.GetMockHTTPContext(request)

	cfg := &Config{UriRegex: []string{"^/([^/]*)/(.*)$", "/$2"}}

	factory := &FilterFactory{cfg: cfg}
	chain := filter.NewDefaultFilterChain()
	_ = factory.PrepareFilterChain(ctx, chain)

	cfg.UriRegex = []string{}

	chain.OnDecode(ctx)

	assert.Equal(t, ctx.GetUrl(), "/query")
}

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

package ratelimit

import (
	"encoding/json"

	"net/http"
)

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config/ratelimit"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	http2 "github.com/apache/dubbo-go-pixiu/pkg/common/http"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/ratelimit/matcher"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	// Kind is the kind of Fallback.
	Kind = constant.RateLimitFilter
)

func init() {
	extension.RegisterHttpFilter(&Plugin{})
}

type (
	// Plugin is http filter plugin.
	Plugin struct {
	}

	// Filter is http filter instance
	Filter struct {
		cfg *ratelimit.Config
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter(hcm *http2.HttpConnectionManager, config interface{}, bs *model.Bootstrap) (extension.HttpFilter, error) {
	specConfig := config.(ratelimit.Config)

	if err := rateLimitInit(&specConfig); err != nil {
		logger.Errorf("rate limit init fail: %s", err)
	}

	return &Filter{cfg: &specConfig}, nil
}

func (f *Filter) PrepareFilterChain(ctx *contexthttp.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *Filter) Handle(hc *contexthttp.HttpContext) {

	path := hc.GetUrl()
	resourceName, ok := matcher.Match(path)
	//if not exists, just skip it.
	if !ok {
		return
	}

	entry, blockErr := sentinel.Entry(resourceName, sentinel.WithResourceType(base.ResTypeAPIGateway), sentinel.WithTrafficType(base.Inbound))

	//if blockErr not nil, indicates the request was blocked by Sentinel
	if blockErr != nil {
		bt, _ := json.Marshal(extension.ErrResponse{Message: "blocked by rate limit"})
		hc.SourceResp = bt
		hc.TargetResp = &client.Response{Data: bt}
		hc.WriteJSONWithStatus(http.StatusTooManyRequests, bt)
		hc.Abort()
		return
	}
	defer entry.Exit()
}

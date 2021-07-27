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
	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/context"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/filter"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/ratelimit/matcher"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// Init cache the filter func & init sentinel
func Init(config *ratelimit.Config) {
	if err := rateLimitInit(config); err != nil {
		logger.Errorf("rate limit init fail: %s", err)

		//if sentinel init fail, just return a empty filter func to avoid error.
		extension.SetFilterFunc(constant.RateLimitFilter, func(context fc.Context) {})
		return
	}

	extension.SetFilterFunc(constant.RateLimitFilter, New().Do())
}

// rateLimit
type rateLimit struct {
}

// New create the rate limit filter
func New() filter.Filter {
	return &rateLimit{}
}

// Do extract the url target & pass or block it.
func (r *rateLimit) Do() fc.FilterFunc {
	return func(ctx fc.Context) {
		hc := ctx.(*contexthttp.HttpContext)

		path := hc.GetAPI().URLPattern
		resourceName, ok := matcher.Match(path)
		//if not exists, just skip it.
		if !ok {
			return
		}

		entry, blockErr := sentinel.Entry(resourceName, sentinel.WithResourceType(base.ResTypeAPIGateway), sentinel.WithTrafficType(base.Inbound))

		//if blockErr not nil, indicates the request was blocked by Sentinel
		if blockErr != nil {
			bt, _ := json.Marshal(filter.ErrResponse{Message: "blocked by rate limit"})
			hc.SourceResp = bt
			hc.TargetResp = &client.Response{Data: bt}
			hc.WriteJSONWithStatus(http.StatusTooManyRequests, bt)
			hc.Abort()
			return
		}
		defer entry.Exit()
	}
}

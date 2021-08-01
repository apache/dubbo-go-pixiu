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
	sc "github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"

	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/context"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/filter"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	fm "github.com/apache/dubbo-go-pixiu/pkg/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// Init cache the filter func & init sentinel
func Init() {
	fm.RegisterFilterFactory(constant.RateLimitFilter, newFactory)
}

// rateLimit
type rateLimit struct {
	conf *Config
}

// newFactory create the rate limit filter
func newFactory() filter.Factory {
	return &rateLimit{conf: &Config{}}
}

func (r *rateLimit) Config() interface{} {
	return r.conf
}

func (r *rateLimit) Apply() (filter.Filter, error) {
	// init matcher
	matcher := newMatcher()
	conf := r.conf
	matcher.load(conf.Resources)

	// init sentinel
	sentinelConf := sc.NewDefaultConfig()
	if len(conf.LogPath) > 0 {
		sentinelConf.Sentinel.Log.Dir = conf.LogPath
	}
	_ = logging.ResetGlobalLogger(getWrappedLogger())

	if err := sentinel.InitWithConfig(sentinelConf); err != nil {
		return nil, err
	}
	OnRulesUpdate(conf.Rules)

	// filter func
	return filterFunc(matcher), nil
}

func filterFunc(m *Matcher) filter.Filter {
	return func(ctx fc.Context) {
		hc := ctx.(*contexthttp.HttpContext)

		path := hc.GetAPI().URLPattern
		resourceName, ok := m.match(path)
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

// OnRulesUpdate update rule
func OnRulesUpdate(rules []*Rule) {
	var enableRules []*flow.Rule
	for _, v := range rules {
		if v.Enable {
			enableRules = append(enableRules, &v.FlowRule)
		}
	}

	if _, err := flow.LoadRules(enableRules); err != nil {
		logger.Warnf("rate limit load rules err: %v", err)
	}
}

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
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	contexthttp "github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
	pkgs "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/sentinel"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

const (
	Kind = constant.HTTPRateLimitFilter
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
		conf    *Config
		matcher *pkgs.Matcher
	}

	// Filter is http filter instance
	Filter struct {
		conf    *Config
		matcher *pkgs.Matcher
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{conf: &Config{}}, nil
}

func (factory *FilterFactory) PrepareFilterChain(ctx *contexthttp.HttpContext, chain filter.FilterChain) error {
	f := &Filter{conf: factory.conf, matcher: factory.matcher}
	chain.AppendDecodeFilters(f)
	return nil
}

func (f *Filter) Decode(hc *contexthttp.HttpContext) filter.FilterStatus {

	path := hc.GetUrl()
	resourceName, ok := f.matcher.Match(path)
	//if not exists, just skip it.
	if !ok {
		return filter.Continue
	}

	entry, blockErr := sentinel.Entry(resourceName, sentinel.WithResourceType(base.ResTypeAPIGateway), sentinel.WithTrafficType(base.Inbound))

	//if blockErr not nil, indicates the request was blocked by Sentinel
	if blockErr != nil {
		bt, _ := json.Marshal(contexthttp.ErrResponse{Message: "blocked by rate limit"})
		hc.SendLocalReply(http.StatusTooManyRequests, bt)
		return filter.Stop
	}
	defer entry.Exit()
	return filter.Continue
}

func (factory *FilterFactory) Config() interface{} {
	return factory.conf
}

func (factory *FilterFactory) Apply() error {
	// init matcher
	factory.matcher = pkgs.NewMatcher()
	conf := factory.conf
	factory.matcher.Load(conf.Resources)

	// init sentinel
	sentinelConf := sc.NewDefaultConfig()
	if len(conf.LogPath) > 0 {
		sentinelConf.Sentinel.Log.Dir = conf.LogPath
	}
	_ = logging.ResetGlobalLogger(pkgs.GetWrappedLogger())

	if err := sentinel.InitWithConfig(sentinelConf); err != nil {
		return err
	}
	OnRulesUpdate(conf.Rules)
	return nil
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

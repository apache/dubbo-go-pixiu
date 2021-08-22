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
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
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
		conf *Config
		matcher *Matcher
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilter() (extension.HttpFilter, error) {
	return &Filter{ conf: &Config{}}, nil
}

func (f *Filter) PrepareFilterChain(ctx *contexthttp.HttpContext) error {
	ctx.AppendFilterFunc(f.Handle)
	return nil
}

func (f *Filter) Handle(hc *contexthttp.HttpContext) {

	path := hc.GetUrl()
	resourceName, ok := f.matcher.match(path)
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


func (r *Filter) Config() interface{} {
	return r.conf
}

func (r *Filter) Apply() error {
	// init matcher
	r.matcher = newMatcher()
	conf := r.conf
	r.matcher.load(conf.Resources)

	// init sentinel
	sentinelConf := sc.NewDefaultConfig()
	if len(conf.LogPath) > 0 {
		sentinelConf.Sentinel.Log.Dir = conf.LogPath
	}
	_ = logging.ResetGlobalLogger(getWrappedLogger())

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

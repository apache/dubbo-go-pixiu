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

package failinject

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/extension/filter"
	contextHttp "github.com/apache/dubbo-go-pixiu/pixiu/pkg/context/http"
)

const (
	Kind = constant.HTTPFailInjectFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

type (
	Plugin struct {
	}

	FilterFactory struct {
		cfg *Config
	}
	Filter struct {
		cfg *Config
	}
)

func (factory *FilterFactory) Apply() error {
	return nil
}

func (factory *FilterFactory) Config() interface{} {
	return factory.cfg
}

func (factory *FilterFactory) PrepareFilterChain(ctx *contextHttp.HttpContext, chain filter.FilterChain) error {
	fmt.Println(factory.cfg)
	f := Filter{
		cfg: factory.cfg,
	}
	chain.AppendDecodeFilters(f)
	return nil
}

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{
		cfg: &Config{
			Rules: make(map[URI]*Rule),
		},
	}, nil
}

func (f Filter) Decode(ctx *contextHttp.HttpContext) filter.FilterStatus {
	uriCfg := &Rule{}
	if v, ok := f.cfg.Rules[URI(ctx.Request.RequestURI)]; ok {
		uriCfg = v
	}
	if uriCfg == nil {
		return filter.Continue
	}

	matched := f.match(uriCfg)
	if !matched {
		return filter.Continue
	}

	if uriCfg.Type == TypeAbort {
		ctx.SendLocalReply(uriCfg.StatusCode, []byte(uriCfg.Body))
		return filter.Stop
	}

	if uriCfg.Type == TypeDelay {
		time.Sleep(uriCfg.Delay)
		ctx.SendLocalReply(uriCfg.StatusCode, []byte(uriCfg.Body))
		return filter.Stop
	}

	return filter.Continue
}

// 当前请求是否命中判断
func (f Filter) match(rule *Rule) (matched bool) {
	if rule == nil {
		return false
	}
	if rule.TriggerType == TriggerTypeAlways {
		return true
	}
	if rule.TriggerType == TriggerTypePercentage {
		return percentage(rule.Odds)
	}
	if rule.TriggerType == TriggerTypeRandom {
		return random()
	}
	return false
}

// random if the current request is matched
func random() bool {
	return (rand.Intn(1000)+1)%2 == 0
}

// percentage if the current request is matched, odds is the percentage
func percentage(odds int) bool {
	if odds <= 0 {
		return false
	}
	// generate rand number in 1-100
	num := rand.Intn(100) + 1
	return num <= odds
}

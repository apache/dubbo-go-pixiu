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
	"math/rand"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contextHttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
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
		return filter.Stop
	}

	return filter.Continue
}

// match determines whether a given request hits the error injection based on the rule's trigger type.
//
// - If the rule is nil, the request doesn't hit the error injection.
// - If the rule's trigger type is "Always", the request always hits the error injection.
// - If the rule's trigger type is "Percentage", the request hits the error injection based on a specified percentage chance.
// - If the rule's trigger type is "Random", the request hits the error injection based on a random chance.
//
// Returns:
//   - true if the request hits the error injection
//   - false otherwise
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

// percentage checks if the current request matches based on a given odds percentage.
// The function returns true with a probability equal to the odds percentage.
//
// Parameters:
//   - odds: Percentage chance (0-100) of the function returning true.
//
// Returns:
//   - true if a generated random number between 1 and 100 (inclusive) is less than or equal to the odds.
//   - false otherwise.
func percentage(odds int) bool {
	if odds <= 0 {
		return false
	}
	// generate rand number in 1-100
	num := rand.Intn(100) + 1
	return num <= odds
}

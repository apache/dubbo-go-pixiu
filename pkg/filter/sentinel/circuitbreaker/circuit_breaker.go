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

package circuitbreaker

import (
	"fmt"
	"strings"
)

import (
	sentinel "github.com/alibaba/sentinel-golang/api"
	"github.com/alibaba/sentinel-golang/core/base"
	"github.com/alibaba/sentinel-golang/core/circuitbreaker"
	sc "github.com/alibaba/sentinel-golang/core/config"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	pkgs "github.com/apache/dubbo-go-pixiu/pkg/filter/sentinel"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	Kind         = constant.HTTPCircuitBreakerFilter
	Segmentation = "@"
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
		cfg     *Config
		matcher *pkgs.Matcher
	}

	// Filter is http filter instance
	Filter struct {
		cfg     *Config
		matcher *pkgs.Matcher
	}

	// Config describe the config of FilterFactory
	Config struct {
		Resources []*pkgs.Resource       `json:"resources,omitempty" yaml:"resources,omitempty"`
		Rules     []*circuitbreaker.Rule `json:"rules" yaml:"rules"` // circuit breaker base config info
	}
)

func (p *Plugin) Kind() string {
	return Kind
}

func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{cfg: &Config{}}, nil
}

// Deep copy config to avoid pointer sharing (factory.cfg may change at runtime)
func (factory *FilterFactory) PrepareFilterChain(ctx *http.HttpContext, chain filter.FilterChain) error {
	f := &Filter{cfg: factory.cfg.DeepCopy(), matcher: factory.matcher}
	chain.AppendDecodeFilters(f)
	chain.AppendEncodeFilters(f)
	return nil
}

func (f *Filter) Decode(ctx *http.HttpContext) filter.FilterStatus {

	url := ctx.GetUrl()
	resourceName, ok := f.matcher.Match(url)

	if !ok {
		return filter.Continue
	}

	// entry, blockErr := sentinel.Entry(resourceName, sentinel.WithResourceType(base.ResTypeAPIGateway), sentinel.WithTrafficType(base.Inbound))
	entry, blockErr := sentinel.Entry(resourceName)

	// if blockErr not nil, indicates the request was blocked by Sentinel
	if blockErr != nil {
		logger.Warnf("circuit breaker request blocked for resource %s: %v", resourceName, blockErr)
		errResp := http.ServiceUnavailable.WithError(fmt.Errorf("circuit breaker open for resource: %s", resourceName))
		ctx.SendLocalReply(errResp.Status, errResp.ToJSON())
		return filter.Stop
	}

	// Store entry in context for later use in Encode phase
	if ctx.Params == nil {
		ctx.Params = make(map[string]any)
	}
	ctx.Params[constant.SentinelEntryKey] = entry

	return filter.Continue
}

// Encode processes the response and reports statistics to Sentinel
func (f *Filter) Encode(ctx *http.HttpContext) filter.FilterStatus {
	entryVal, ok := ctx.Params[constant.SentinelEntryKey]
	if !ok {
		// No entry in context, skip
		return filter.Continue
	}

	entry, ok := entryVal.(*base.SentinelEntry)
	if !ok || entry == nil {
		logger.Warnf("circuit breaker invalid sentinel entry type in context")
		return filter.Continue
	}

	// Ensure entry.Exit() is called
	defer entry.Exit()

	// Report error to Sentinel if response indicates failure
	// Consider 5xx status codes as errors for circuit breaker
	statusCode := ctx.GetStatusCode()
	if statusCode >= 500 && statusCode < 600 {
		// Create detailed error with status code and request context
		err := fmt.Errorf("backend returned HTTP %d for %s %s",
			statusCode, ctx.GetMethod(), ctx.GetUrl())
		entry.SetError(err)
		logger.Debugf("circuit breaker reported error to Sentinel: %v", err)
	}

	return filter.Continue
}

func (factory *FilterFactory) Config() any {
	return factory.cfg
}

func (factory *FilterFactory) Apply() error {
	if len(factory.cfg.Rules) == 0 || len(factory.cfg.Resources) == 0 {
		return fmt.Errorf("circuit breaker router or resources is null")
	}

	factory.resetResources()

	// init matcher
	factory.matcher = pkgs.NewMatcher()
	factory.matcher.Load(factory.cfg.Resources)

	resourcesMap := make(map[string]circuitbreaker.Rule, len(factory.cfg.Rules))
	for _, rule := range factory.cfg.Rules {
		resourcesMap[rule.Resource] = *rule
	}

	rules := make([]*circuitbreaker.Rule, 0, len(factory.cfg.Rules))
	for _, rule := range factory.cfg.Resources {
		c, ok := resourcesMap[strings.Split(rule.Name, Segmentation)[0]]
		if !ok {
			logger.Warn("circuit breaker resource does not exist")
			continue
		}

		c.Resource = rule.Name
		rules = append(rules, &c)
	}

	return loadRules(rules)
}

func (factory *FilterFactory) resetResources() {
	resources := make([]*pkgs.Resource, 0, len(factory.cfg.Resources))

	for _, rule := range factory.cfg.Resources {
		for _, item := range rule.Items {
			resources = append(resources, &pkgs.Resource{Name: rule.Name + Segmentation + item.Pattern, Items: []*pkgs.Item{item}})
		}
	}
	factory.cfg.Resources = resources
}

func loadRules(rules []*circuitbreaker.Rule) error {
	conf := sc.NewDefaultConfig()
	conf.Sentinel.Log.Logger = pkgs.GetWrappedLogger()
	err := sentinel.InitWithConfig(conf)
	if err != nil {
		logger.Error("circuit breaker init fail ", err)
		return err
	}

	if _, err = circuitbreaker.LoadRules(rules); err != nil {
		logger.Error("circuit breaker load rules fail ", err)
		return err
	}
	return nil
}

// DeepCopy returns a new independent copy of Config
// Deep copy slices/maps to avoid sharing pointers with the factory
func (config *Config) DeepCopy() *Config {
	if config == nil {
		return nil
	}
	cpConfig := *config

	if config.Resources != nil {
		cpConfig.Resources = make([]*pkgs.Resource, len(config.Resources))
		for index, resource := range config.Resources {
			if resource != nil {
				cpConfig.Resources[index] = resource.DeepCopy()
			}
		}
	}

	if config.Rules != nil {
		cpConfig.Rules = make([]*circuitbreaker.Rule, len(config.Rules))
		for index, rule := range config.Rules {
			if rule != nil {
				cpRule := *rule
				cpConfig.Rules[index] = &cpRule
			}
		}
	}

	return &cpConfig
}

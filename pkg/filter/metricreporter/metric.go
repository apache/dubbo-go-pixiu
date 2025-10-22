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

package metricreporter

import (
	"fmt"
	"sync"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension/filter"
	contextHttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	// Kind defines the filter kind
	Kind = constant.HTTPMetricReporterFilter
)

func init() {
	filter.RegisterHttpFilter(&Plugin{})
}

// Plugin is the plugin for metric reporter filter.
type Plugin struct{}

// Kind returns the filter kind.
func (p *Plugin) Kind() string {
	return Kind
}

// CreateFilterFactory creates a new filter factory.
func (p *Plugin) CreateFilterFactory() (filter.HttpFilterFactory, error) {
	return &FilterFactory{
		cfg: &Config{},
	}, nil
}

// FilterFactory is the factory for metric reporter filter.
type FilterFactory struct {
	cfg          *Config
	pullReporter *OTelPullReporter
	pushReporter *PushReporter
	initOnce     sync.Once
}

// Config returns the configuration.
func (factory *FilterFactory) Config() any {
	return factory.cfg
}

// Apply initializes the filter factory.
func (factory *FilterFactory) Apply() error {
	var initErr error
	factory.initOnce.Do(func() {
		initErr = factory.initialize()
	})
	return initErr
}

// initialize performs the actual initialization of reporters.
func (factory *FilterFactory) initialize() error {
	// Validate configuration
	if err := factory.cfg.Validate(); err != nil {
		return err
	}

	logger.Infof("[MetricReporter] Initializing with mode: %s", factory.cfg.Mode)

	// Initialize pull reporter if enabled
	if factory.cfg.PullConfig.Enabled {
		if err := factory.initPullReporter(); err != nil {
			return err
		}
	}

	// Initialize push reporter if enabled
	if factory.cfg.PushConfig.Enabled {
		factory.initPushReporter()
	}

	return nil
}

// initPullReporter initializes the pull mode reporter using OpenTelemetry.
func (factory *FilterFactory) initPullReporter() error {
	factory.pullReporter = NewOTelPullReporter(&factory.cfg.PullConfig)
	if err := factory.pullReporter.Start(); err != nil {
		return fmt.Errorf("failed to start pull reporter: %w", err)
	}
	logger.Infof("[MetricReporter] OpenTelemetry pull mode enabled on port %d, path %s",
		factory.cfg.PullConfig.Port, factory.cfg.PullConfig.Path)
	return nil
}

// initPushReporter initializes the push mode reporter.
func (factory *FilterFactory) initPushReporter() {
	factory.pushReporter = NewPushReporter(&factory.cfg.PushConfig)
	logger.Infof("[MetricReporter] Push mode enabled, gateway: %s, job: %s, interval: %d",
		factory.cfg.PushConfig.GatewayURL,
		factory.cfg.PushConfig.JobName,
		factory.cfg.PushConfig.PushInterval)
}

// PrepareFilterChain prepares the filter chain.
func (factory *FilterFactory) PrepareFilterChain(ctx *contextHttp.HttpContext, chain filter.FilterChain) error {
	f := &Filter{
		factory:      factory,
		pullReporter: factory.pullReporter,
		pushReporter: factory.pushReporter,
	}
	// Append to encode filters (runs after response)
	chain.AppendEncodeFilters(f)
	return nil
}

// Filter is the metric reporter filter instance.
type Filter struct {
	factory      *FilterFactory
	pullReporter *OTelPullReporter
	pushReporter *PushReporter
}

// Decode is not used in this filter.
func (f *Filter) Decode(ctx *contextHttp.HttpContext) filter.FilterStatus {
	return filter.Continue
}

// Encode collects metrics from context and reports them.
func (f *Filter) Encode(ctx *contextHttp.HttpContext) filter.FilterStatus {
	// Get all metrics from context
	metrics := ctx.GetAllMetrics()
	if len(metrics) == 0 {
		return filter.Continue
	}

	// Report to pull reporter (update in-memory registry)
	if f.pullReporter != nil {
		f.pullReporter.Report(metrics)
	}

	// Report to push reporter (push to gateway)
	if f.pushReporter != nil {
		f.pushReporter.Report(metrics)
	}

	return filter.Continue
}


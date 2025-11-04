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
)

import (
	"go.opentelemetry.io/otel/metric/instrument/syncint64"
)

// Config defines the configuration for the unified metric reporter filter.
type Config struct {
	// Mode defines the metric reporting mode: "pull" or "push"
	Mode string `yaml:"mode" json:"mode"`

	// PushConfig configuration for push mode (Push Gateway)
	// Note: Pull mode uses global metric configuration (metric.enable, metric.prometheus_port)
	PushConfig PushConfig `yaml:"push_config" json:"push_config"`
}

// PushConfig defines the configuration for push mode.
type PushConfig struct {
	// GatewayURL is the Push Gateway URL (e.g., http://localhost:9091)
	GatewayURL string `yaml:"gateway_url" json:"gateway_url"`

	// JobName is the job name for Push Gateway
	JobName string `yaml:"job_name" json:"job_name"`

	// PushInterval defines how many requests to process before pushing metrics
	PushInterval int `yaml:"push_interval" json:"push_interval"`

	// MetricPath is the path to push metrics to Push Gateway
	MetricPath string `yaml:"metric_path" json:"metric_path"`
}

type OTelInstruments struct {
	totalElapsed syncint64.Counter
	totalCount   syncint64.Counter
	totalError   syncint64.Counter
	sizeRequest  syncint64.Counter
	sizeResponse syncint64.Counter
	durationHist syncint64.Histogram
}

// Validate validates the configuration based on mode.
func (c *Config) Validate() error {
	// Validate mode
	if c.Mode != "pull" && c.Mode != "push" {
		return fmt.Errorf("invalid mode '%s', must be 'pull' or 'push'", c.Mode)
	}

	// Apply defaults and validate push config if in push mode
	// Pull mode has no filter-level configuration (uses global metric config)
	if c.Mode == "push" {
		c.PushConfig.ApplyDefaults()
		return c.PushConfig.Validate()
	}

	return nil
}

// ApplyDefaults applies default values to empty fields.
func (c *PushConfig) ApplyDefaults() {
	if c.GatewayURL == "" {
		c.GatewayURL = "http://localhost:9091"
	}

	if c.JobName == "" {
		c.JobName = "pixiu"
	}

	if c.PushInterval <= 0 {
		c.PushInterval = 100
	}

	if c.MetricPath == "" {
		c.MetricPath = "/metrics"
	}
}

// Validate validates push mode configuration.
func (c *PushConfig) Validate() error {
	if c.GatewayURL == "" {
		return fmt.Errorf("push gateway_url cannot be empty")
	}

	if c.JobName == "" {
		return fmt.Errorf("push job_name cannot be empty")
	}

	if c.PushInterval <= 0 {
		return fmt.Errorf("push interval %d must be greater than 0", c.PushInterval)
	}

	if c.MetricPath == "" {
		return fmt.Errorf("push metric_path cannot be empty")
	}

	return nil
}

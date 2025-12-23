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

package metric

import (
	"fmt"
)

import (
	"go.opentelemetry.io/otel/metric"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
)

// Config defines the configuration for the unified metric reporter filter.
type Config struct {
	// Mode defines the metric reporting mode: "pull" or "push"
	Mode string `yaml:"mode" json:"mode"`

	// Push configuration for push mode (Push Gateway)
	// Note: Pull mode uses global metric configuration (metric.enable, metric.prometheus_port)
	Push PushConfig `yaml:"push_config" json:"push_config"`
}

// PushConfig defines the configuration for push mode.
type PushConfig struct {
	// GatewayURL is the Push Gateway URL (default: http://localhost:9091)
	GatewayURL string `yaml:"gateway_url" json:"gateway_url"`

	// JobName is the job name for Push Gateway (default: pixiu)
	JobName string `yaml:"job_name" json:"job_name"`

	// PushInterval defines how many requests to process before pushing metrics (default: 100)
	PushInterval int `yaml:"push_interval" json:"push_interval"`

	// MetricPath is the path to push metrics to Push Gateway (default: /metrics)
	MetricPath string `yaml:"metric_path" json:"metric_path"`
}

type OTelInstruments struct {
	totalElapsed metric.Int64Counter
	totalCount   metric.Int64Counter
	totalError   metric.Int64Counter
	sizeRequest  metric.Int64Counter
	sizeResponse metric.Int64Counter
	durationHist metric.Int64Histogram
}

// Validate validates the configuration based on mode.
func (c *Config) Validate() error {
	// Apply default mode if not specified
	if c.Mode == "" {
		c.Mode = constant.DefaultMetricMode
	}

	// Validate mode
	if c.Mode != "pull" && c.Mode != "push" {
		return fmt.Errorf("invalid mode '%s', must be 'pull' or 'push'", c.Mode)
	}

	// Validate push config if in push mode
	// Pull mode has no filter-level configuration (uses global metric config)
	if c.Mode == "push" {
		return c.Push.Validate()
	}

	return nil
}

// Validate validates push mode configuration and applies defaults for empty fields.
func (c *PushConfig) Validate() error {
	// Apply defaults for empty fields
	if c.GatewayURL == "" {
		c.GatewayURL = constant.DefaultMetricPushGatewayURL
	}

	if c.JobName == "" {
		c.JobName = constant.DefaultMetricPushJobName
	}

	if c.PushInterval <= 0 {
		c.PushInterval = constant.DefaultMetricPushInterval
	}

	if c.MetricPath == "" {
		c.MetricPath = constant.DefaultMetricPushPath
	}

	// All fields now have values (either user-provided or defaults)
	return nil
}

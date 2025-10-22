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

// Config defines the configuration for the unified metric reporter filter.
type Config struct {
	// Mode defines the metric reporting mode: "pull", "push", or "both"
	Mode string `yaml:"mode" json:"mode" default:"pull"`

	// PullConfig configuration for pull mode (Prometheus scraping)
	PullConfig PullConfig `yaml:"pull_config" json:"pull_config"`

	// PushConfig configuration for push mode (Push Gateway)
	PushConfig PushConfig `yaml:"push_config" json:"push_config"`
}

// PullConfig defines the configuration for pull mode.
type PullConfig struct {
	// Enabled enables pull mode
	Enabled bool `yaml:"enabled" json:"enabled" default:"true"`

	// Port is the port to expose metrics endpoint
	Port int `yaml:"port" json:"port" default:"9090"`

	// Path is the HTTP path for metrics endpoint
	Path string `yaml:"path" json:"path" default:"/metrics"`
}

// PushConfig defines the configuration for push mode.
type PushConfig struct {
	// Enabled enables push mode
	Enabled bool `yaml:"enabled" json:"enabled" default:"false"`

	// GatewayURL is the Push Gateway URL (e.g., http://localhost:9091)
	GatewayURL string `yaml:"gateway_url" json:"gateway_url" default:"http://127.0.0.1:9091"`

	// JobName is the job name for Push Gateway
	JobName string `yaml:"job_name" json:"job_name" default:"pixiu"`

	// PushInterval defines how many requests to process before pushing metrics
	// If set to 100, metrics will be pushed every 100 requests
	PushInterval int `yaml:"push_interval" json:"push_interval" default:"100"`

	// MetricPath is the path to push metrics to Push Gateway
	MetricPath string `yaml:"metric_path" json:"metric_path" default:"/metrics"`
}

// Validate validates the configuration.
func (c *Config) Validate() error {
	if c.Mode != "pull" && c.Mode != "push" && c.Mode != "both" {
		c.Mode = "pull" // default to pull mode
	}

	// Set defaults for pull config
	if c.PullConfig.Port == 0 {
		c.PullConfig.Port = 9090
	}
	if c.PullConfig.Path == "" {
		c.PullConfig.Path = "/metrics"
	}

	// Set defaults for push config
	if c.PushConfig.GatewayURL == "" {
		c.PushConfig.GatewayURL = "http://127.0.0.1:9091"
	}
	if c.PushConfig.JobName == "" {
		c.PushConfig.JobName = "pixiu"
	}
	if c.PushConfig.PushInterval <= 0 {
		c.PushConfig.PushInterval = 100
	}
	if c.PushConfig.MetricPath == "" {
		c.PushConfig.MetricPath = "/metrics"
	}

	// Enable based on mode
	switch c.Mode {
	case "pull":
		c.PullConfig.Enabled = true
		c.PushConfig.Enabled = false
	case "push":
		c.PullConfig.Enabled = false
		c.PushConfig.Enabled = true
	case "both":
		c.PullConfig.Enabled = true
		c.PushConfig.Enabled = true
	}

	return nil
}


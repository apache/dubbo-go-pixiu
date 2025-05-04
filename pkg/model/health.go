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

package model

// HealthCheck
type HealthCheckConfig struct {
	Protocol            string                 `yaml:"protocol" json:"protocol,omitempty" mapstructure:"protocol"`
	TimeoutConfig       string                 `yaml:"timeout" json:"timeout,omitempty" mapstructure:"timeout"`
	IntervalConfig      string                 `yaml:"interval" json:"interval,omitempty" mapstructure:"interval"`
	InitialDelaySeconds string                 `yaml:"initial_delay_seconds" json:"initial_delay_seconds,omitempty" mapstructure:"initial_delay_seconds"`
	HealthyThreshold    uint32                 `yaml:"healthy_threshold" json:"healthy_threshold,omitempty" mapstructure:"healthy_threshold"`
	UnhealthyThreshold  uint32                 `yaml:"unhealthy_threshold" json:"unhealthy_threshold,omitempty" mapstructure:"unhealthy_threshold"`
	ServiceName         string                 `yaml:"service_name" json:"service_name,omitempty" mapstructure:"service_name"`
	SessionConfig       map[string]interface{} `yaml:"check_config" json:"check_config,omitempty" mapstructure:"check_config"`
	CommonCallbacks     []string               `yaml:"common_callbacks" json:"common_callbacks,omitempty" mapstructure:"common_callbacks"`
}

// HttpHealthCheck
type HttpHealthCheck struct {
	HealthCheckConfig
	Host             string
	Path             string
	UseHttp2         bool
	ExpectedStatuses int64
}

// GrpcHealthCheck
type GrpcHealthCheck struct {
	HealthCheckConfig
	ServiceName string
	Authority   string
}

// CustomHealthCheck
type CustomHealthCheck struct {
	HealthCheckConfig
	Name   string
	Config interface{}
}

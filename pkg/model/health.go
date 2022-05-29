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
	Protocol            string                 `json:"protocol,omitempty"`
	TimeoutConfig       string                 `json:"timeout,omitempty"`
	IntervalConfig      string                 `json:"interval,omitempty"`
	InitialDelaySeconds string                 `json:"initial_delay_seconds,omitempty"`
	HealthyThreshold    uint32                 `json:"healthy_threshold,omitempty"`
	UnhealthyThreshold  uint32                 `json:"unhealthy_threshold,omitempty"`
	ServiceName         string                 `json:"service_name,omitempty"`
	SessionConfig       map[string]interface{} `json:"check_config,omitempty"`
	CommonCallbacks     []string               `json:"common_callbacks,omitempty"`
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

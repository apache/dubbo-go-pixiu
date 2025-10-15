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

type (
	// LLMMeta LLM metadata for llm call
	LLMMeta struct {
		Provider            string      `yaml:"provider" json:"provider"`                                                                               // Provider the cluster unique name
		APIKey              string      `yaml:"api_key" json:"api_key" mapstructure:"api_key"`                                                          // APIKey the cluster unique name
		RetryPolicy         RetryPolicy `yaml:"retry_policy" json:"retry_policy" mapstructure:"retry_policy"`                                           // RetryPolicy key
		Fallback            bool        `yaml:"fallback" json:"fallback" mapstructure:"fallback"`                                                       // Fallback to the next provider if failed
		HealthCheckInterval int64       `yaml:"health_check_interval" json:"health_check_interval" mapstructure:"health_check_interval" default:"5000"` // HealthCheckInterval the interval for health check
	}
)

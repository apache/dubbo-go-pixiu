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

package prometheus

type (
	MetricCollectConfiguration struct {
		Rules MetricCollectRule `yaml:"metric_collect_rules" json:"metric_collect_rules"`
	}

	MetricCollectRule struct {
		Enable     bool   `json:"enbale,omitempty" yaml:"enable,omitempty"`
		MetricPath string `json:"metric_path,omitempty" yaml:"metric_path,omitempty"`
		// Push Gateway URL in format http://domain:port
		// where JOBNAME can be any string of your choice
		PushGatewayURL string `default:"http://127.0.0.1:9091" json:"push_gateway_url,omitempty" yaml:"push_gateway_url,omitempty"`
		// Push interval in seconds
		// lint:ignore ST1011 renaming would be breaking change
		PushIntervalSeconds int    `json:"push_interval_seconds,omitempty" yaml:"push_interval_seconds,omitempty"`
		PushJobName         string `json:"push_job_name,omitempty" yaml:"push_job_name,omitempty"`
		// Pushgateway job name, defaults to "prometheus"
	}
)

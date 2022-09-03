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
	// AuthorityConfiguration blacklist/whitelist config
	MetricCollectConfiguration struct {
		Rules MetricCollectRule `yaml:"metric_collect_rules" json:"metric_collect_rules"` // Rules the authority rule list
	}
	// AuthorityRule blacklist/whitelist rule
	MetricCollectRule struct {
		Enable       bool   `json:"enbale,omitempty" yaml:"enable,omitempty"`
		MeticPath    string `json:"metric_path,omitempty" yaml:"metric_path,omitempty"`
		ExporterPort int    `json:"exporter_port,omitempty" yaml:"exporter_port,omitempty"` // Items the authority rule items

	}
)

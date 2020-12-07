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

// AccessLog
type AccessLog struct {
	Name   string          `yaml:"name" json:"name" mapstructure:"name"`
	Filter AccessLogFilter `yaml:"filter" json:"filter" mapstructure:"filter"`
	Config interface{}     `yaml:"config" json:"config" mapstructure:"config"`
}

type AccessLogConfig struct {
	Enable     bool   `yaml:"enable" json:"enable" mapstructure:"enable"`
	OutPutPath string `yaml:"outPutPath" json:"outPutPath" mapstructure:"outPutPath"`
}

// AccessLogFilter
type AccessLogFilter struct {
	StatusCodeFilter StatusCodeFilter `yaml:"status_code_filter" json:"status_code_filter" mapstructure:"status_code_filter"`
	DurationFilter   DurationFilter   `yaml:"duration_filter" json:"duration_filter" mapstructure:"duration_filter"`
}

type StatusCodeFilter struct {
}

type DurationFilter struct {
}

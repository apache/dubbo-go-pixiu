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

package config

import (
	"github.com/alibaba/sentinel-golang/core/flow"
)

// Config rate limit config
type Config struct {
	APIResources []APIResource `json:"resources" yaml:"resources"`
	Rules        []Rule        `json:"rules" yaml:"rules"`
	LogPath      string        `json:"logPath" yaml:"logPath"`
}

// APIResource API group for rate limit, all API in group is considered to be the same resource
type APIResource struct {
	Name  string    `json:"name" yaml:"name"`
	Items []APIItem `json:"items" yaml:"items"`
}

// APIItem API item for group
type APIItem struct {
	MatchStrategy MatchStrategy `json:"matchStrategy" yaml:"matchStrategy"`
	Pattern       string        `json:"pattern,omitempty" yaml:"pattern,omitempty"`
}

// MatchStrategy API match strategy
type MatchStrategy int32

const (
	EXACT MatchStrategy = iota
	REGEX
	ANT_PATH
)

type Rule struct {
	flow.Rule `yaml:",inline"`
	Enable    bool `json:"enable" yaml:"enable"`
}

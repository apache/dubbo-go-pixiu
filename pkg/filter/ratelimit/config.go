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

package ratelimit

import (
	"github.com/alibaba/sentinel-golang/core/flow"
)

// Config rate limit config
type Config struct {
	Resources []*Resource `json:"resources,omitempty" yaml:"resources,omitempty"`
	Rules     []*Rule     `json:"rules,omitempty" yaml:"rules,omitempty"`
	LogPath   string      `json:"logPath,omitempty" yaml:"logPath,omitempty"`
}

// Resource API group for rate limit, all API in group is considered to be the same resource
type Resource struct {
	ID    int64   `json:"id,omitempty" yaml:"id,omitempty"`
	Name  string  `json:"name,omitempty" yaml:"name,omitempty"`
	Items []*Item `json:"items" yaml:"items"`
}

// Item API item for group
type Item struct {
	MatchStrategy MatchStrategy `json:"matchStrategy,omitempty" yaml:"matchStrategy,omitempty"`
	Pattern       string        `json:"pattern,omitempty" yaml:"pattern"`
}

// Rule api group 's rate-limit rule
type Rule struct {
	ID       int64     `json:"id,omitempty" yaml:"id,omitempty"`
	FlowRule flow.Rule `json:"flowRule,omitempty" yaml:"flowRule,omitempty"`
	Enable   bool      `json:"enable,omitempty" yaml:"enable,omitempty"`
}

// MatchStrategy API match strategy
type MatchStrategy int32

const (
	EXACT    MatchStrategy = 0
	REGEX    MatchStrategy = 1
	ANT_PATH MatchStrategy = 2
)

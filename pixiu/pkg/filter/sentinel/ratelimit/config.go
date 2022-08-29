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

import (
	pkgs "github.com/apache/dubbo-go-pixiu/pixiu/pkg/filter/sentinel"
)

type (
	// Config rate limit config
	Config struct {
		Resources []*pkgs.Resource `json:"resources,omitempty" yaml:"resources,omitempty"`
		Rules     []*Rule          `json:"rules,omitempty" yaml:"rules,omitempty"`
		LogPath   string           `json:"logPath,omitempty" yaml:"logPath,omitempty"`
	}

	// Rule api group 's rate-limit rule
	Rule struct {
		ID       int64     `json:"id,omitempty" yaml:"id,omitempty"`
		FlowRule flow.Rule `json:"flowRule,omitempty" yaml:"flowRule,omitempty"`
		Enable   bool      `json:"enable,omitempty" yaml:"enable,omitempty"`
	}
)

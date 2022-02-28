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

package sentinel

const (
	EXACT MatchStrategy = iota
	REGEX
	AntPath
)

type (
	// Resource API group for rate sentinel, all API in group is considered to be the same resource
	Resource struct {
		ID    int64   `json:"id,omitempty" yaml:"id,omitempty"`
		Name  string  `json:"name,omitempty" yaml:"name,omitempty"`
		Items []*Item `json:"items" yaml:"items"`
	}

	// Item API item for group
	Item struct {
		MatchStrategy MatchStrategy `json:"matchStrategy,omitempty" yaml:"matchStrategy,omitempty"`
		Pattern       string        `json:"pattern,omitempty" yaml:"pattern"`
	}

	// MatchStrategy API match strategy
	MatchStrategy int32
)

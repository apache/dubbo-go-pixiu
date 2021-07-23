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

// StringMatcher matcher string
type StringMatcher struct {
	Matcher MatcherType
}

// Match
func (sm *StringMatcher) Match() (bool, error) {
	return true, nil
}

// MatcherType matcher type
type MatcherType int32

const (
	Exact MatcherType = 0 + iota
	Prefix
	Suffix
	Regex
)

var MatcherTypeName = map[int32]string{
	0: "Exact",
	1: "Prefix",
	2: "Suffix",
	3: "Regex",
}

var MatcherTypeValue = map[string]int32{
	"Exact":  0,
	"Prefix": 1,
	"Suffix": 2,
	"Regex":  3,
}

// HeaderMatcher header matcher struct
// Name header key, Value header value, Regex header value is regex
type HeaderMatcher struct {
	Name  string `yaml:"name" json:"name" mapstructure:"name"`
	Value string `yaml:"value" json:"value" mapstructure:"value"`
	Regex bool   `yaml:"regex" json:"regex" mapstructure:"regex"`
}

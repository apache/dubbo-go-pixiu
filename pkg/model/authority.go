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

// AuthorityConfiguration blacklist/whitelist config
type AuthorityConfiguration struct {
	Rules []AuthorityRule `yaml:"authority_rules" json:"authority_rules"` //Rules the authority rule list
}

// AuthorityRule blacklist/whitelist rule
type AuthorityRule struct {
	Strategy StrategyType `yaml:"strategy" json:"strategy"` //Strategy the authority rule strategy
	Limit    LimitType    `yaml:"limit" json:"limit"`       //Limit the authority rule limit
	Items    []string     `yaml:"items" json:"items"`       //Items the authority rule items
}

// StrategyType the authority rule strategy enum
type StrategyType int32

// StrategyType strategy type const
const (
	Whitelist StrategyType = 0
	Blacklist StrategyType = 1
)

// StrategyTypeName key int32 for StrategyType, value string
var StrategyTypeName = map[int32]string{
	0: "Whitelist",
	1: "Blacklist",
}

// StrategyTypeValue key string, value int32 for StrategyType
var StrategyTypeValue = map[string]int32{
	"Whitelist": 0,
	"Blacklist": 1,
}

// LimitType the authority rule limit enum
type LimitType int32

// LimitType limit type const
const (
	IP  LimitType = 0
	App LimitType = 1
)

// LimitTypeName key int32 for LimitType, value string
var LimitTypeName = map[int32]string{
	0: "IP",
	1: "App",
}

// LimitTypeValue key string, value int32 for LimitType
var LimitTypeValue = map[string]int32{
	"IP":  0,
	"App": 1,
}

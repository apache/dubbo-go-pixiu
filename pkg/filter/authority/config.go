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

package authority

// StrategyType strategy type const
const (
	Whitelist StrategyType = 0
	Blacklist StrategyType = 1
)

// LimitType limit type const
const (
	IP  LimitType = 0
	App LimitType = 1
)

var (
	// StrategyTypeName key int32 for StrategyType, value string
	StrategyTypeName = map[int32]string{
		0: "Whitelist",
		1: "Blacklist",
	}

	// StrategyTypeValue key string, value int32 for StrategyType
	StrategyTypeValue = map[string]int32{
		"Whitelist": 0,
		"Blacklist": 1,
	}

	// LimitTypeName key int32 for LimitType, value string
	LimitTypeName = map[int32]string{
		0: "IP",
		1: "App",
	}

	// LimitTypeValue key string, value int32 for LimitType
	LimitTypeValue = map[string]int32{
		"IP":  0,
		"App": 1,
	}
)

type (
	// AuthorityConfiguration blacklist/whitelist config
	AuthorityConfiguration struct {
		Rules []AuthorityRule `yaml:"authority_rules" json:"authority_rules"` // Rules the authority rule list
	}

	// AuthorityRule blacklist/whitelist rule
	AuthorityRule struct {
		Strategy StrategyType `yaml:"strategy" json:"strategy"` // Strategy the authority rule strategy
		Limit    LimitType    `yaml:"limit" json:"limit"`       // Limit the authority rule limit
		Items    []string     `yaml:"items" json:"items"`       // Items the authority rule items
	}

	// StrategyType the authority rule strategy enum
	StrategyType int32
	// LimitType the authority rule limit enum
	LimitType int32
)

// DeepCopy returns a new independent copy of Config
// Deep copy slices/maps to avoid sharing pointers with the factory
func (rule *AuthorityRule) DeepCopy() *AuthorityRule {
	if rule == nil {
		return nil
	}

	cp := *rule

	if rule.Items != nil {
		cp.Items = make([]string, len(rule.Items))
		copy(cp.Items, rule.Items)
	} else {
		cp.Items = nil
	}

	return &cp
}

// DeepCopy returns a new independent copy of Config
// Deep copy slices/maps to avoid sharing pointers with the factory
func (ac *AuthorityConfiguration) DeepCopy() *AuthorityConfiguration {
	if ac == nil {
		return nil
	}

	cp := *ac

	if ac.Rules != nil {
		cp.Rules = make([]AuthorityRule, len(ac.Rules))
		for i := range ac.Rules {
			cp.Rules[i] = *ac.Rules[i].DeepCopy()
		}
	} else {
		cp.Rules = nil
	}

	return &cp
}

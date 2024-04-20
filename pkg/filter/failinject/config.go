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

package failinject

import (
	"time"
)

type URI string

const (
	TypeAbort = "abort"
	TypeDelay = "delay"

	TriggerTypeRandom     = "random"
	TriggerTypeAlways     = "always"
	TriggerTypePercentage = "percentage"
)

type (
	Config struct {
		Rules map[URI]*Rule `yaml:"fail_inject_rules" json:"fail_inject_rules"`
	}
	Rule struct {
		StatusCode  int           `yaml:"status_code" json:"status_code"`
		Body        string        `yaml:"body" json:"body"`
		Type        string        `yaml:"type" json:"type"`                 // abort, delay
		TriggerType string        `yaml:"trigger_type" json:"trigger_type"` // random, always, percentage
		Odds        int           `yaml:"odds" json:"odds"`                 // 0-100
		Delay       time.Duration `yaml:"delay" json:"delay"`
	}
)

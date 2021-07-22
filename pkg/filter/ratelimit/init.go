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
	sentinel "github.com/alibaba/sentinel-golang/api"
	sc "github.com/alibaba/sentinel-golang/core/config"
	"github.com/alibaba/sentinel-golang/core/flow"
	"github.com/alibaba/sentinel-golang/logging"

	"github.com/pkg/errors"

	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/api/config/ratelimit"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/filter/ratelimit/matcher"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

func rateLimitInit(c *ratelimit.Config) error {
	sentinelConf := sc.NewDefaultConfig()
	if len(c.LogPath) > 0 {
		sentinelConf.Sentinel.Log.Dir = c.LogPath
	}
	_ = logging.ResetGlobalLogger(getWrappedLogger())

	if err := sentinel.InitWithConfig(sentinelConf); err != nil {
		return errors.Wrap(err, "rate limit init fail")
	}
	matcher.Init()

	OnUpdate(c)
	return nil
}

// OnUpdate update api & rule
func OnUpdate(c *ratelimit.Config) {
	OnResourcesUpdate(c.Resources)
	OnRulesUpdate(c.Rules)
}

// OnRulesUpdate update rule
func OnRulesUpdate(rules []ratelimit.Rule) {
	var enableRules []*flow.Rule
	for _, v := range rules {
		if v.Enable {
			enableRules = append(enableRules, &v.FlowRule)
		}
	}

	if _, err := flow.LoadRules(enableRules); err != nil {
		logger.Warnf("rate limit load rules err: %v", err)
	}
}

// OnResourcesUpdate update matcher for resources
func OnResourcesUpdate(apis []ratelimit.Resource) {
	matcher.Load(apis)
}

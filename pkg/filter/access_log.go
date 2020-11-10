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

package filter

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
	"time"
)

func init() {
	extension.SetFilterFunc(constant.AccessLogFilter, AccessLogFilter())
}

func AccessLogFilter() context.FilterFunc {
	return func(c context.Context) {
		filters := c.GetAPI().Filters
		headers := c.GetAPI().Headers
		if hasFilterConf(filters) || hasHeader(headers) {
			doLog(c)
		}
	}
}

func doLog(c context.Context) {
	start := time.Now()
	c.Next()
	latency := time.Now().Sub(start)
	logger.Infof("[dubboproxy go] [UPSTREAM] receive request | %d | %s | %s | %s | ", c.StatusCode(), latency, c.GetMethod(), c.GetUrl())
}

func hasHeader(headers []config.Params) bool {
	if len(headers) > 0 {
		for _, h := range headers {
			if h.Name == constant.AccessLogFilter {
				return true
			}
		}
	}
	return false
}

func hasFilterConf(filters []string) bool {
	if len(filters) > 0 {
		for _, f := range filters {
			if f == constant.AccessLogFilter {
				return true
			}
		}
	}
	return false
}

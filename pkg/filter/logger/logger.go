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

package logger

import (
	"time"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
)

// nolint
func Init() {
	extension.SetFilterFunc(constant.LoggerFilter, loggerFilterFunc())
}

func loggerFilterFunc() context.FilterFunc {
	return New().Do()
}

// loggerFilter is a filter for simple logger.
type loggerFilter struct {
}

// New create logger filter.
func New() filter.Filter {
	return &loggerFilter{}
}

// Logger logger filter, print url and latency
func (f loggerFilter) Do() context.FilterFunc {
	return func(c context.Context) {
		start := time.Now()

		c.Next()

		latency := time.Now().Sub(start)

		logger.Infof("[dubboproxy go] [UPSTREAM] receive request | %d | %s | %s | %s | ", c.StatusCode(), latency, c.GetMethod(), c.GetUrl())
	}
}

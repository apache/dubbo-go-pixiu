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

package recovery

import (
	"github.com/dubbogo/dubbo-go-proxy-filter/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy-filter/pkg/filter"
)

import (
	"github.com/dubbogo/dubbo-go-pixiu/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-pixiu/pkg/common/extension"
	"github.com/dubbogo/dubbo-go-pixiu/pkg/logger"
)

// nolint
func Init() {
	extension.SetFilterFunc(constant.RecoveryFilter, recoveryFilterFunc())
}

func recoveryFilterFunc() context.FilterFunc {
	return New().Do()
}

// recoveryFilter is a filter for recover.
type recoveryFilter struct {
}

// New create timeout filter.
func New() filter.Filter {
	return &recoveryFilter{}
}

// Recovery execute recoveryFilter filter logic, if recover happen, print log or do other things.
func (f recoveryFilter) Do() context.FilterFunc {
	return func(c context.Context) {
		defer func() {
			if err := recover(); err != nil {
				logger.Warnf("[dubbopixiu go] error:%+v", err)

				c.WriteErr(err)
			}
		}()
		c.Next()
	}
}

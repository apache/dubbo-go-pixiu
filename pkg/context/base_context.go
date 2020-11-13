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

package context

import (
	"context"
	"math"
	"time"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
)

const abortIndex int8 = math.MaxInt8 / 2

// BaseContext
type BaseContext struct {
	Context
	Index   int8
	Filters FilterChain
	Timeout time.Duration
	Ctx     context.Context
}

// NewBaseContext create base context.
func NewBaseContext() *BaseContext {
	return &BaseContext{Index: -1}
}

// Next should be used only inside middleware.
// It executes the pending handlers in the chain inside the calling handler.
// See example in GitHub.
func (c *BaseContext) Next() {
	c.Index++
	for c.Index < int8(len(c.Filters)) {
		c.Filters[c.Index](c)
		c.Index++
	}
}

// Abort  filter chain break , filter after the current filter will not executed.
func (c *BaseContext) Abort() {
	c.Index = abortIndex
}

// AbortWithError  filter chain break , filter after the current filter will not executed. And log will print.
func (c *BaseContext) AbortWithError(message string, err error) {
	c.Index = abortIndex
	logger.Errorf("abort with err : %v", err)
}

// AppendFilterFunc  append filter func.
func (c *BaseContext) AppendFilterFunc(ff ...FilterFunc) {
	for _, v := range ff {
		c.Filters = append(c.Filters, v)
	}
}

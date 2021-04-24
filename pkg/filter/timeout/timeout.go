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

package timeout

import (
	"context"
	"encoding/json"
	"net/http"
	"time"
)

import (
	fc "github.com/dubbogo/dubbo-go-pixiu-filter/pkg/context"
	"github.com/dubbogo/dubbo-go-pixiu-filter/pkg/filter"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/client"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/common/extension"
	contexthttp "github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// nolint
func init() {
	extension.SetFilterFunc(constant.TimeoutFilter, timeoutFilterFunc(0))
}

func timeoutFilterFunc(duration time.Duration) fc.FilterFunc {
	return New(duration).Do()
}

// timeoutFilter is a filter for control request time out.
type timeoutFilter struct {
	// global timeout
	waitTime time.Duration
}

// New create timeout filter.
func New(t time.Duration) filter.Filter {
	if t <= 0 {
		t = constant.DefaultTimeout
	}
	return &timeoutFilter{
		waitTime: t,
	}
}

// Do execute timeoutFilter filter logic.
func (f timeoutFilter) Do() fc.FilterFunc {
	return func(c fc.Context) {
		hc := c.(*contexthttp.HttpContext)

		ctx, cancel := context.WithTimeout(hc.Ctx, f.getTimeout(hc.Timeout))
		defer cancel()
		// Channel capacity must be greater than 0.
		// Otherwise, if the parent coroutine quit due to timeout,
		// the child coroutine may never be able to quit.
		finishChan := make(chan struct{}, 1)
		go func() {
			// panic by recovery
			c.Next()
			finishChan <- struct{}{}
		}()

		select {
		// timeout do.
		case <-ctx.Done():
			logger.Warnf("api:%s request timeout", hc.GetAPI().URLPattern)
			bt, _ := json.Marshal(filter.ErrResponse{Message: http.ErrHandlerTimeout.Error()})
			hc.SourceResp = bt
			hc.TargetResp = &client.Response{Data: bt}
			hc.WriteJSONWithStatus(http.StatusGatewayTimeout, bt)
			c.Abort()
		case <-finishChan:
			// finish call do something.
		}
	}
}

func (f timeoutFilter) getTimeout(t time.Duration) time.Duration {
	if t <= 0 {
		return f.waitTime
	}

	return t
}

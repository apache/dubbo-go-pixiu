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
	"github.com/dubbogo/dubbo-go-proxy/pkg/client"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	selfcontext "github.com/dubbogo/dubbo-go-proxy/pkg/context"
	contexthttp "github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
)

func init() {
	extension.SetFilterFunc(constant.TimeoutFilter, New(0).Do())
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
func (f *timeoutFilter) Do() selfcontext.FilterFunc {
	return func(c selfcontext.Context) {
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
			hc.WriteFailWithCode(http.StatusGatewayTimeout, bt)
			c.Abort()
		case <-finishChan:
			// finish call do something.
		}
	}
}

func (f *timeoutFilter) getTimeout(t time.Duration) time.Duration {
	if t <= 0 {
		return f.waitTime
	}

	return t
}

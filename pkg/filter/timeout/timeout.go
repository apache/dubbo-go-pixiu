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
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
	"net/http"
	"time"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	selfcontext "github.com/dubbogo/dubbo-go-proxy/pkg/context"
	contexthttp "github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
)

const (
	// TimeoutError timeout code
	TimeoutError = "PROXY005"
)

func init() {
	extension.SetFilterFunc(constant.TimeoutFilter, NewTimeoutFilter(0).Do())
}

// timeoutFilter is a filter for control request time out.
type timeoutFilter struct {
	// global timeout
	waitTime time.Duration
}

// NewTimeoutFilter create timeout filter.
func NewTimeoutFilter(t time.Duration) *timeoutFilter {
	if t <= 0 {
		t = constant.DefaultTimeout
	}
	return &timeoutFilter{
		waitTime: t,
	}
}

// Do timeoutFilter execute filter logic.
func (f *timeoutFilter) Do() selfcontext.FilterFunc {
	return func(c selfcontext.Context) {
		hc := c.(*contexthttp.HttpContext)

		ctx, cancel := context.WithTimeout(hc.Ctx, hc.GetTimeout(hc.Timeout))
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
			hc.Lock.Lock()
			defer hc.Lock.Unlock()
			logger.Warnf("api %v request timeout", hc.GetAPI())
			bt, _ := json.Marshal(errResponse{Code: TimeoutError,
				Message: http.ErrHandlerTimeout.Error()})
			hc.WriteWithStatus(http.StatusServiceUnavailable, bt)
			hc.AddHeader(constant.HeaderKeyContextType, constant.HeaderValueJsonUtf8)
			hc.Abort()
		case <-finishChan:
			hc.Lock.Lock()
			defer hc.Lock.Unlock()
			// finish call do something.
		}
	}
}

type errResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

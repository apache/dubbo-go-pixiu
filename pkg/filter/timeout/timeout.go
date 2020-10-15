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
	"fmt"
	"net/http"
	"time"

	"github.com/dubbogo/dubbo-go-proxy/pkg/common/constant"
	"github.com/dubbogo/dubbo-go-proxy/pkg/common/extension"
	selfcontext "github.com/dubbogo/dubbo-go-proxy/pkg/context"
	contexthttp "github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
)

const (
	// timeout code
	HandlerFuncTimeout = "E509"
	// unknown code
	ErrUnknownError = "E003"
)

// TODO 改初始化的时候注入
func init() {
	extension.SetFilterFunc(constant.TimeoutFilter, NewTimeoutFilter().Run())
}

// timeoutFilter is a filter for control request time out.
type timeoutFilter struct {
	waitTime time.Duration
}

// NewTimeoutFilter create timeout filter.
func NewTimeoutFilter() *timeoutFilter {
	return &timeoutFilter{
		waitTime: time.Minute,
	}
}

// DoFilter do filter.
func (f *timeoutFilter) Run() selfcontext.FilterFunc {
	return func(c selfcontext.Context) {
		hc := c.(*contexthttp.HttpContext)

		ctx, cancel := context.WithTimeout(context.Background(), hc.GetTimeout(hc.Timeout))

		// Channel capacity must be greater than 0.
		// Otherwise, if the parent coroutine quit due to timeout,
		// the child coroutine may never be able to quit.
		finish := make(chan struct{}, 1)
		panicChan := make(chan interface{}, 1)
		go func() {
			// 变动点1: 增加子协程的recover
			defer func() {
				if p := recover(); p != nil {
					panicChan <- p
				}
			}()
			c.Next()
			finish <- struct{}{}
		}()

		select {
		case p := <-panicChan:
			bt, _ := json.Marshal(errResponse{Code: ErrUnknownError,
				Msg: fmt.Sprintf("unknow internal error, %v", p)})
			c.WriteWithStatus(http.StatusInternalServerError, bt)
			c.Abort()
		case <-ctx.Done():
			hc.Lock.Lock()
			defer hc.Lock.Unlock()

			bt, _ := json.Marshal(errResponse{Code: HandlerFuncTimeout,
				Msg: http.ErrHandlerTimeout.Error()})
			c.WriteWithStatus(http.StatusInternalServerError, bt)
			c.Abort()
			cancel()
			// If timeout happen, the buffer cannot be cleared actively,
			// but wait for the GC to recycle.
		case <-finish:
			hc.Lock.Lock()
			defer hc.Lock.Unlock()
		}
	}
}

type errResponse struct {
	Code string `json:"code"`
	Msg  string `json:"msg"`
}

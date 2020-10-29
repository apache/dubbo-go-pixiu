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
	"fmt"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter/recovery"
	"net/http"
	"testing"
	"time"
)

import (
	selfcontext "github.com/dubbogo/dubbo-go-proxy/pkg/context"
	selfhttp "github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
)

func TestPanic(t *testing.T) {
	c := MockHttpContext(testPanicFilter)
	c.Next()
	// print
	// 500
	// "timeout filter test panic"
}

var testPanicFilter = func(c selfcontext.Context) {
	time.Sleep(time.Millisecond * 100)
	panic("timeout filter test panic")
}

func TestTimeout(t *testing.T) {
	c := MockHttpContext(testTimeoutFilter)
	c.Next()
	// print
	// 503
	// {"code":"S005","message":"http: Handler timeout"}
}

var testTimeoutFilter = func(c selfcontext.Context) {
	time.Sleep(time.Second * 3)
}

func TestNormal(t *testing.T) {
	c := MockHttpContext(testNormalFilter)
	c.Next()
}

var testNormalFilter = func(c selfcontext.Context) {
	time.Sleep(time.Millisecond * 200)
	fmt.Println("normal call")
}

func MockHttpContext(fc ...selfcontext.FilterFunc) *selfhttp.HttpContext {
	result := &selfhttp.HttpContext{
		BaseContext: &selfcontext.BaseContext{
			Index: -1,
			Ctx:   context.Background(),
		},
	}

	w := MockWriter{header: map[string][]string{}}
	result.ResetWritermen(&w)
	result.Reset()

	result.Filters = append(result.Filters, NewTimeoutFilter().Do(), recovery.NewRecoveryFilter().Do())
	for i := range fc {
		result.Filters = append(result.Filters, fc[i])
	}

	return result
}

type MockWriter struct {
	header http.Header
}

func (w *MockWriter) Header() http.Header {
	return w.header
}

func (w *MockWriter) Write(b []byte) (int, error) {
	fmt.Println(string(b))
	return -1, nil
}

func (w *MockWriter) WriteHeader(statusCode int) {
	fmt.Println(statusCode)
}

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
	"context"
	"fmt"
	"net/http"
	"testing"
	time "time"
)

import (
	selfcontext "github.com/dubbogo/dubbo-go-proxy/pkg/context"
	selfhttp "github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
)

func TestRecovery(t *testing.T) {
	c := MockHTTPContext(func(c selfcontext.Context) {
		time.Sleep(time.Millisecond * 100)
		// panic
		var m map[string]string
		m["name"] = "1"
	})
	c.Next()
	// print
	// 500
	// "assignment to entry in nil map"
}

func MockHTTPContext(fc ...selfcontext.FilterFunc) *selfhttp.HttpContext {
	result := &selfhttp.HttpContext{
		BaseContext: &selfcontext.BaseContext{
			Index: -1,
			Ctx:   context.Background(),
		},
	}

	w := MockWriter{header: map[string][]string{}}
	result.ResetWritermen(&w)
	result.Reset()

	result.Filters = append(result.Filters, NewRecoveryFilter().Do())
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

func TestName(t *testing.T) {
	tt, _ := time.ParseDuration("111")
	t.Log(tt)
	if tt == 0 {
		t.Log("0")
	}
}

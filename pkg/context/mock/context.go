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

package mock

import (
	"context"
	"fmt"
	"net/http"
)

import (
	pkgcontext "github.com/dubbogo/dubbo-go-proxy/pkg/context"
	contexthttp "github.com/dubbogo/dubbo-go-proxy/pkg/context/http"
)

func GetMockHTTPContext(r *http.Request, fc ...pkgcontext.FilterFunc) *contexthttp.HttpContext {
	result := &contexthttp.HttpContext{
		BaseContext: &pkgcontext.BaseContext{
			Index: -1,
			Ctx:   context.Background(),
		},
		Request: r,
	}

	w := mockWriter{header: map[string][]string{}}
	result.ResetWritermen(&w)
	result.Reset()

	for i := range fc {
		result.Filters = append(result.Filters, fc[i])
	}

	return result
}

type mockWriter struct {
	header http.Header
}

func (w *mockWriter) Header() http.Header {
	return w.header
}

func (w *mockWriter) Write(b []byte) (int, error) {
	fmt.Println(string(b))
	return -1, nil
}

func (w *mockWriter) WriteHeader(statusCode int) {
	fmt.Println(statusCode)
}

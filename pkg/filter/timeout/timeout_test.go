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
	"fmt"
	"testing"
	"time"
)

import (
	pkgcontext "github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context/mock"
	"github.com/dubbogo/dubbo-go-proxy/pkg/filter/recovery"
)

func TestPanic(t *testing.T) {
	c := mock.GetMockHTTPContext(nil, timeoutFilterFunc(0), recovery.New().Do(), testPanicFilter)
	c.Next()
	// print
	// 500
	// "timeout filter test panic"
}

var testPanicFilter = func(c pkgcontext.Context) {
	time.Sleep(time.Millisecond * 100)
	panic("timeout filter test panic")
}

func TestTimeout(t *testing.T) {
	c := mock.GetMockHTTPContext(nil, timeoutFilterFunc(0), testTimeoutFilter)
	c.Next()
	// print
	// 503
	// {"code":"S005","message":"http: Handler timeout"}
}

var testTimeoutFilter = func(c pkgcontext.Context) {
	time.Sleep(time.Second * 3)
}

func TestNormal(t *testing.T) {
	c := mock.GetMockHTTPContext(nil, testNormalFilter)
	c.Next()
}

var testNormalFilter = func(c pkgcontext.Context) {
	time.Sleep(time.Millisecond * 200)
	fmt.Println("normal call")
}

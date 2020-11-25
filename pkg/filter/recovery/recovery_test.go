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
	"testing"
	"time"
)

import (
	selfcontext "github.com/dubbogo/dubbo-go-proxy/pkg/context"
	"github.com/dubbogo/dubbo-go-proxy/pkg/context/mock"
)

func TestRecovery(t *testing.T) {
	c := mock.GetMockHTTPContext(nil, func(c selfcontext.Context) {
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

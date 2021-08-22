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
	"github.com/apache/dubbo-go-pixiu/pkg/context/http"
	"testing"
	"time"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/context/mock"
)

// nolint
func TestRecovery(t *testing.T) {
	filter := GetMock()
	sleepFilter := &SleepFilter{}
	c := mock.GetMockHTTPContext(nil, filter, sleepFilter)
	c.Next()
	// print
	// 500
	// "assignment to entry in nil map"
}

type SleepFilter struct {
}

func (rf *SleepFilter) PrepareFilterChain(ctx *http.HttpContext) error {
	ctx.AppendFilterFunc(rf.Handle)
	return nil
}

func (rf *SleepFilter) Handle(c *http.HttpContext) {
	time.Sleep(time.Millisecond * 100)
	// panic
	var m map[string]string
	m["name"] = "1"
}

func (f *SleepFilter) Config() interface{} {
	return nil
}

func (f *SleepFilter) Apply() error {
	return nil
}

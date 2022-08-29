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

package tracing

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
)

type traceGenerator func(Trace) Trace

// One protocol for one Trace.
var TraceFactory map[ProtocolName]traceGenerator

func init() {
	TraceFactory = make(map[ProtocolName]traceGenerator)
	TraceFactory[HTTPProtocol] = NewHTTPTracer
}

// Driver maintains all tracers and the provider.
type TraceDriverManager struct {
	driver    *TraceDriver
	bootstrap *model.Bootstrap
}

func CreateDefaultTraceDriverManager(bs *model.Bootstrap) *TraceDriverManager {
	manager := &TraceDriverManager{
		bootstrap: bs,
	}
	manager.driver = InitDriver(bs)
	return manager
}

func (manager *TraceDriverManager) GetDriver() *TraceDriver {
	return manager.driver
}

func (manager *TraceDriverManager) Close() {
	driver := manager.driver
	_ = driver.Tp.Shutdown(driver.context)
}

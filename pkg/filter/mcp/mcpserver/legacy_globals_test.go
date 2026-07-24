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

package mcpserver

import (
	"sync"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/transport"
)

var testGlobals struct {
	sync.Mutex
	registry *ToolRegistry
	sessions *transport.SessionManager
	dynamic  *DynamicConsumer
}

func init() {
	resetTestGlobals = func() {
		testGlobals.Lock()
		defer testGlobals.Unlock()
		if testGlobals.sessions != nil {
			testGlobals.sessions.Stop()
		}
		testGlobals.registry = nil
		testGlobals.sessions = nil
		testGlobals.dynamic = nil
	}
}

func GetOrInitRegistry() *ToolRegistry {
	testGlobals.Lock()
	defer testGlobals.Unlock()
	if testGlobals.registry == nil {
		testGlobals.registry = NewToolRegistry()
	}
	return testGlobals.registry
}

func GetOrInitSessionManager() *transport.SessionManager {
	testGlobals.Lock()
	defer testGlobals.Unlock()
	if testGlobals.sessions == nil {
		testGlobals.sessions = transport.NewSessionManager()
	}
	return testGlobals.sessions
}

func getOrInitTestDynamicConsumer() *DynamicConsumer {
	testGlobals.Lock()
	defer testGlobals.Unlock()
	if testGlobals.dynamic == nil {
		if testGlobals.registry == nil {
			testGlobals.registry = NewToolRegistry()
		}
		if testGlobals.sessions == nil {
			testGlobals.sessions = transport.NewSessionManager()
		}
		testGlobals.dynamic = NewDynamicConsumer(testGlobals.registry, testGlobals.sessions, transport.NewSSEHandler(testGlobals.sessions))
	}
	return testGlobals.dynamic
}

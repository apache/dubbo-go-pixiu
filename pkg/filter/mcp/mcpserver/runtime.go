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
	"fmt"
	"sync"
	"sync/atomic"
)

var (
	runtimeIndexMu   sync.RWMutex
	runtimeIndex     = map[string]*RuntimeState{}
	runtimeSeq       uint64
	resetTestGlobals func()
)

func registerRuntime(runtime *RuntimeState) string {
	if runtime == nil {
		return ""
	}
	id := fmt.Sprintf("mcp-runtime-%d", atomic.AddUint64(&runtimeSeq, 1))
	runtimeIndexMu.Lock()
	runtimeIndex[id] = runtime
	runtimeIndexMu.Unlock()
	return id
}

func unregisterRuntime(id string) {
	if id == "" {
		return
	}
	runtimeIndexMu.Lock()
	delete(runtimeIndex, id)
	runtimeIndexMu.Unlock()
}

// Stop releases all runtime-owned state. It does not affect other MCP filters.
func (r *RuntimeState) Stop() {
	if r == nil {
		return
	}
	unregisterRuntime(r.id)
	if r.sessionManager != nil {
		r.sessionManager.Stop()
	}
	if r.plans != nil {
		r.plans.Stop()
	}
}

// GetOrInitDynamicConsumer preserves the legacy registry-center entry point
// only when a single MCP filter instance is registered. Multiple instances
// require an explicit runtime binding, so this function refuses to guess.
func GetOrInitDynamicConsumer() *DynamicConsumer {
	runtimeIndexMu.RLock()
	defer runtimeIndexMu.RUnlock()
	if len(runtimeIndex) != 1 {
		return nil
	}
	for _, runtime := range runtimeIndex {
		return runtime.dynamic
	}
	return nil
}

// ResetGlobalState resets runtime references created by tests.
func ResetGlobalState() {
	runtimeIndexMu.Lock()
	runtimes := make([]*RuntimeState, 0, len(runtimeIndex))
	for _, runtime := range runtimeIndex {
		runtimes = append(runtimes, runtime)
	}
	runtimeIndex = map[string]*RuntimeState{}
	runtimeSeq = 0
	runtimeIndexMu.Unlock()

	for _, runtime := range runtimes {
		runtime.Stop()
	}
	if resetTestGlobals != nil {
		resetTestGlobals()
	}
}

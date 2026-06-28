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
	"errors"
	"fmt"
	"sync"
	"sync/atomic"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

var (
	runtimeIndexMu   sync.RWMutex
	runtimeIndex     = map[string]*RuntimeState{}
	runtimeSeq       uint64
	resetTestGlobals func()
)

var ErrDynamicConsumerUnavailable = errors.New("mcp dynamic consumer unavailable")
var ErrDynamicConsumerAmbiguous = errors.New("multiple mcp dynamic consumers registered")

// ServerPublicationSink is the dynamic registry publication target owned by one
// MCP runtime. Registry adapters bind to this sink once instead of discovering a
// runtime on every event.
type ServerPublicationSink interface {
	RuntimeID() string
	ApplyMcpServerConfigByServer(serverId string, cfg *model.McpServerConfig) error
}

func registerRuntime(runtime *RuntimeState) string {
	if runtime == nil {
		return ""
	}
	id := fmt.Sprintf("mcp-runtime-%d", atomic.AddUint64(&runtimeSeq, 1))
	runtimeIndexMu.Lock()
	runtimeIndex[id] = runtime
	runtimeIndexMu.Unlock()
	if runtime.dynamic != nil {
		runtime.dynamic.setRuntimeID(id)
	}
	return id
}

func unregisterRuntime(id string) {
	if id == "" {
		return
	}
	runtimeIndexMu.Lock()
	runtime := runtimeIndex[id]
	delete(runtimeIndex, id)
	runtimeIndexMu.Unlock()
	if runtime != nil && runtime.dynamic != nil {
		runtime.dynamic.setRuntimeID("")
	}
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

// DynamicConsumerForSingleRuntime returns the only registered dynamic consumer.
// Registry-center callbacks are process-global today, so multiple MCP filters
// must be treated as ambiguous instead of guessing which runtime should update.
func DynamicConsumerForSingleRuntime() (*DynamicConsumer, error) {
	runtimeIndexMu.RLock()
	defer runtimeIndexMu.RUnlock()
	switch len(runtimeIndex) {
	case 0:
		return nil, ErrDynamicConsumerUnavailable
	case 1:
		for _, runtime := range runtimeIndex {
			if runtime.dynamic == nil {
				return nil, ErrDynamicConsumerUnavailable
			}
			return runtime.dynamic, nil
		}
	default:
		return nil, ErrDynamicConsumerAmbiguous
	}
	return nil, ErrDynamicConsumerUnavailable
}

// ServerPublicationSinkForSingleRuntime returns the publication sink for the
// only registered runtime. It is a compatibility bridge for the process-global
// registry adapter; callers must treat ambiguous runtime state as no target.
func ServerPublicationSinkForSingleRuntime() (ServerPublicationSink, error) {
	return DynamicConsumerForSingleRuntime()
}

// GetOrInitDynamicConsumer preserves the legacy registry-center entry point.
//
// Deprecated: use DynamicConsumerForSingleRuntime to distinguish unavailable
// and ambiguous runtime state.
func GetOrInitDynamicConsumer() *DynamicConsumer {
	consumer, err := DynamicConsumerForSingleRuntime()
	if err != nil {
		return nil
	}
	return consumer
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

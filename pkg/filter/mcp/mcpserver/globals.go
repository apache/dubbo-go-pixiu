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
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/router"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/transport"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

// Global singletons for MCP server components
var (
	globalRegistry       *ToolRegistry
	globalDynamic        *DynamicConsumer
	globalSessionManager *transport.SessionManager
	globalPlanStore      *router.SessionPlanStore

	// sync.Once variables for thread-safe singleton initialization
	registryOnce       sync.Once
	dynamicOnce        sync.Once
	sessionManagerOnce sync.Once
	planStoreOnce      sync.Once
	planStoreHookOnce  sync.Once
)

// GetOrInitRegistry returns the global tool registry singleton
func GetOrInitRegistry() *ToolRegistry {
	registryOnce.Do(func() {
		globalRegistry = NewToolRegistry()
		logger.Infof("[dubbo-go-pixiu] mcp server initialized global tool registry")
	})
	return globalRegistry
}

// GetOrInitDynamicConsumer returns the global dynamic consumer singleton
func GetOrInitDynamicConsumer() *DynamicConsumer {
	dynamicOnce.Do(func() {
		reg := GetOrInitRegistry()
		sm := GetOrInitSessionManager()
		sseHandler := transport.NewSSEHandler(sm)
		globalDynamic = NewDynamicConsumer(reg, sm, sseHandler)
		logger.Infof("[dubbo-go-pixiu] mcp server initialized global dynamic consumer")
	})
	return globalDynamic
}

// GetOrInitSessionManager returns the global session manager singleton
func GetOrInitSessionManager() *transport.SessionManager {
	sessionManagerOnce.Do(func() {
		globalSessionManager = transport.NewSessionManager()
		logger.Infof("[dubbo-go-pixiu] mcp server initialized global session manager")
	})
	return globalSessionManager
}

// GetOrInitPlanStore returns the global router session plan store singleton.
func GetOrInitPlanStore() *router.SessionPlanStore {
	return GetOrInitPlanStoreWithMaxEntries(router.DefaultPlanMaxEntries)
}

// GetOrInitPlanStoreWithMaxEntries returns the global router session plan
// store singleton and applies the configured capacity.
func GetOrInitPlanStoreWithMaxEntries(maxEntries int) *router.SessionPlanStore {
	planStoreOnce.Do(func() {
		globalPlanStore = router.NewSessionPlanStoreWithMaxEntries(maxEntries)
		sm := GetOrInitSessionManager()
		planStoreHookOnce.Do(func() {
			sm.AddSessionRemovedHandler(globalPlanStore.Delete)
		})
		logger.Infof("[dubbo-go-pixiu] mcp server initialized global router plan store")
	})
	if globalPlanStore != nil {
		globalPlanStore.SetMaxEntries(maxEntries)
	}
	return globalPlanStore
}

// ResetGlobalState resets all global singletons (for testing)
func ResetGlobalState() {
	globalRegistry = nil
	globalDynamic = nil
	if globalSessionManager != nil {
		globalSessionManager.Stop()
	}
	globalSessionManager = nil
	if globalPlanStore != nil {
		globalPlanStore.Stop()
	}
	globalPlanStore = nil

	registryOnce = sync.Once{}
	dynamicOnce = sync.Once{}
	sessionManagerOnce = sync.Once{}
	planStoreOnce = sync.Once{}
	planStoreHookOnce = sync.Once{}

	logger.Debugf("[dubbo-go-pixiu] mcp server global state reset")
}

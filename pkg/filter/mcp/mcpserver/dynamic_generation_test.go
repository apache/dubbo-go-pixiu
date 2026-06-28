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
	"context"
	"testing"
)

import (
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/router"
	"github.com/apache/dubbo-go-pixiu/pkg/filter/mcp/mcpserver/transport"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

func newDynamicRefreshFixture(t *testing.T, cfg *model.RouterConfig, tools []model.ToolConfig) (*ToolRegistry, *transport.SessionManager, *router.SessionPlanStore, router.ToolSelector, *transport.MCPSession) {
	t.Helper()

	registry := NewToolRegistry()
	require.NoError(t, registry.ReplaceAllTools(tools))
	sm := transport.NewSessionManager()
	store := router.NewSessionPlanStore()
	selector, err := router.BuildWithOptions(cfg, store, router.BuildOptions{VisibleFP: visibleToolsFingerprint})
	require.NoError(t, err)
	session, err := sm.CreateSession()
	require.NoError(t, err)
	require.NoError(t, store.ActivateSession(session.ID, session.Generation))
	sm.AddSessionRemovedHandler(store.DeleteSession)
	return registry, sm, store, selector, session
}

func selectPlanForSession(t *testing.T, registry *ToolRegistry, selector router.ToolSelector, session *transport.MCPSession, claims map[string]any) *router.SelectionPlan {
	t.Helper()
	snapshot := registry.toolCatalogSnapshotUnsafe()
	plan, err := selector.Select(context.Background(), router.SelectionContext{
		SessionID:         session.ID,
		SessionGeneration: session.Generation,
		Claims:            claims,
		CatalogVersion:    snapshot.Version,
	}, snapshot.orderedToolsUnsafe())
	require.NoError(t, err)
	return plan
}

func singlePlanContext(t *testing.T, store *router.SessionPlanStore) router.SessionPlanContext {
	t.Helper()
	contexts := store.SessionPlanContexts()
	require.Len(t, contexts, 1)
	return contexts[0]
}

func TestDynamicRefreshDoesNotReviveDeletedSessionPlan(t *testing.T) {
	registry, sm, store, selector, session := newDynamicRefreshFixture(t, &model.RouterConfig{}, []model.ToolConfig{
		createTestToolConfig("old", "old"),
	})
	defer sm.Stop()
	defer store.Stop()

	selectPlanForSession(t, registry, selector, session, nil)
	base := singlePlanContext(t, store)
	sm.RemoveSession(session.ID)

	require.NoError(t, registry.ReplaceAllTools([]model.ToolConfig{createTestToolConfig("new", "new")}))
	snapshot := registry.toolCatalogSnapshotUnsafe()
	refresher := selector.(router.PlanRefreshSelector)
	plan, committed, err := refresher.RefreshPlan(context.Background(), base, snapshot.orderedToolsUnsafe())

	require.NoError(t, err)
	assert.False(t, committed)
	assert.Nil(t, plan)
	assert.Equal(t, 0, store.Len())
}

func TestDynamicRefreshDoesNotOverwriteNewerToolsListPlan(t *testing.T) {
	registry, sm, store, selector, session := newDynamicRefreshFixture(t, &model.RouterConfig{}, []model.ToolConfig{
		createTestToolConfig("old", "old"),
	})
	defer sm.Stop()
	defer store.Stop()

	oldPlan := selectPlanForSession(t, registry, selector, session, nil)
	require.Equal(t, []string{"old"}, oldPlan.ToolNames)
	base := singlePlanContext(t, store)

	require.NoError(t, registry.ReplaceAllTools([]model.ToolConfig{createTestToolConfig("new", "new")}))
	newPlan := selectPlanForSession(t, registry, selector, session, nil)
	require.Equal(t, []string{"new"}, newPlan.ToolNames)

	refresher := selector.(router.PlanRefreshSelector)
	plan, committed, err := refresher.RefreshPlan(context.Background(), base, []model.ToolConfig{createTestToolConfig("stale", "stale")})

	require.NoError(t, err)
	assert.False(t, committed)
	assert.Nil(t, plan)
	got, ok := store.Get(base.Key)
	require.True(t, ok)
	assert.Equal(t, []string{"new"}, got.ToolNames)
}

func TestDynamicRefreshRejectsStaleIdentityGeneration(t *testing.T) {
	registry, sm, store, selector, session := newDynamicRefreshFixture(t, &model.RouterConfig{}, []model.ToolConfig{
		createTestToolConfig("tool", "tool"),
	})
	defer sm.Stop()
	defer store.Stop()

	selectPlanForSession(t, registry, selector, session, map[string]any{"tenant": "acme"})
	base := singlePlanContext(t, store)
	selectPlanForSession(t, registry, selector, session, map[string]any{"tenant": "globex"})

	refresher := selector.(router.PlanRefreshSelector)
	plan, committed, err := refresher.RefreshPlan(context.Background(), base, registry.toolCatalogSnapshotUnsafe().orderedToolsUnsafe())

	require.NoError(t, err)
	assert.False(t, committed)
	assert.Nil(t, plan)
	got, ok := store.Get(base.Key)
	require.True(t, ok)
	assert.NotEqual(t, base.Plan.IdentityHash, got.IdentityHash)
}

func TestDynamicSessionDeletionCallbacksAreIdempotent(t *testing.T) {
	registry, sm, store, selector, session := newDynamicRefreshFixture(t, &model.RouterConfig{}, []model.ToolConfig{
		createTestToolConfig("tool", "tool"),
	})
	defer sm.Stop()
	defer store.Stop()

	selectPlanForSession(t, registry, selector, session, nil)
	store.DeleteSession(session.ID)
	store.DeleteSession(session.ID)
	sm.RemoveSession(session.ID)
	sm.RemoveSession(session.ID)
	assert.Equal(t, 0, store.Len())
}

func TestDynamicChangedVisibleDefinitionDetectsDescriptionChange(t *testing.T) {
	tool := createTestToolConfig("tool", "old description")
	registry, sm, store, selector, session := newDynamicRefreshFixture(t, &model.RouterConfig{}, []model.ToolConfig{tool})
	defer sm.Stop()
	defer store.Stop()

	oldPlan := selectPlanForSession(t, registry, selector, session, nil)
	require.NotEmpty(t, oldPlan.VisibleFingerprint)

	consumer := NewDynamicConsumer(registry, sm, transport.NewSSEHandler(sm))
	consumer.SetGovernance(selector, store)
	consumer.SetDebounceTime(0)

	tool.Description = "new description"
	_, err := consumer.applyServerTools("server-a", []model.ToolConfig{tool})
	require.NoError(t, err)

	changed := consumer.sessionsWithChangedVisibleSet()
	_, ok := changed[session.ID]
	assert.True(t, ok, "same tool name with changed tools/list definition must notify")
}

func TestDynamicChangedVisibleDefinitionIgnoresBackendOnlyChange(t *testing.T) {
	tool := createTestToolConfig("tool", "description")
	tool.BackendURL = "http://127.0.0.1:8080"
	registry, sm, store, selector, session := newDynamicRefreshFixture(t, &model.RouterConfig{}, []model.ToolConfig{tool})
	defer sm.Stop()
	defer store.Stop()

	selectPlanForSession(t, registry, selector, session, nil)

	consumer := NewDynamicConsumer(registry, sm, transport.NewSSEHandler(sm))
	consumer.SetGovernance(selector, store)
	consumer.SetDebounceTime(0)

	tool.BackendURL = "http://127.0.0.1:9090"
	_, err := consumer.applyServerTools("server-a", []model.ToolConfig{tool})
	require.NoError(t, err)

	changed := consumer.sessionsWithChangedVisibleSet()
	_, ok := changed[session.ID]
	assert.False(t, ok, "backend-only changes must not trigger tools/list_changed")
}

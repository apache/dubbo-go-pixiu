/*
 * Licensed to the Apache Software Foundation (ASF) under one or more
 * contributor license agreements.  See the NOTICE file distributed with
 * this work for additional information regarding copyright ownership.
 * The ASF licenses this file to you under the Apache License, Version 2.0
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

package logic

import (
	"context"
	"errors"
	"testing"
)

import (
	"github.com/stretchr/testify/require"

	"gopkg.in/yaml.v3"
)

import (
	legacyconfig "github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema"
)

func runtimeClusterAt(t *testing.T, kv *fakeRouteBindingKV, key string) legacyconfig.Cluster {
	t.Helper()
	var cluster legacyconfig.Cluster
	require.NoError(t, yaml.Unmarshal(kv.fakeValue(key), &cluster))
	return cluster
}

func TestClusterPublishHistoryRollbackAndDraftDeletion(t *testing.T) {
	store, kv := newTestClusterDraftStore(t)
	ctx := context.Background()
	invalid, err := store.CreateDraft(ctx, schema.AdminObject{Kind: schema.KindAdminCluster})
	require.NoError(t, err)
	_, err = store.Publish(ctx, invalid.ID, invalid.Revision)
	var issues schema.ValidationErrors
	require.ErrorAs(t, err, &issues)
	require.False(t, kv.fakeHasPrefix(store.root+"/"+Clusters+"/"))
	require.False(t, kv.fakeHasPrefix(store.root+"/"+clusterHistorySegment+"/"))

	draft, err := store.CreateDraft(ctx, testClusterObject("backend", "127.0.0.1", 20880))
	require.NoError(t, err)
	first, err := store.Publish(ctx, draft.ID, draft.Revision)
	require.NoError(t, err)
	require.Equal(t, 1, first.RuntimeID)
	require.Equal(t, "127.0.0.1", runtimeClusterAt(t, kv, store.runtimeKey(first.RuntimeID)).Address)
	status, err := store.Status(ctx, draft.ID)
	require.NoError(t, err)
	require.True(t, status.PublishedExists)
	require.False(t, status.Dirty)
	diff, err := store.Diff(ctx, draft.ID)
	require.NoError(t, err)
	require.Empty(t, diff.Changes)

	current, err := store.GetDraft(ctx, draft.ID)
	require.NoError(t, err)
	changed := testClusterObject("backend", "127.0.0.2", 20881)
	current, err = store.UpdateDraft(ctx, draft.ID, changed, current.Revision)
	require.NoError(t, err)
	_, err = store.Publish(ctx, draft.ID, draft.Revision)
	require.ErrorIs(t, err, ErrClusterPublishConflict)
	second, err := store.Publish(ctx, draft.ID, current.Revision)
	require.NoError(t, err)
	require.Equal(t, first.RuntimeID, second.RuntimeID)
	require.Equal(t, "127.0.0.2", runtimeClusterAt(t, kv, store.runtimeKey(first.RuntimeID)).Address)
	history, err := store.History(ctx, draft.ID)
	require.NoError(t, err)
	require.Len(t, history, 2)
	require.Equal(t, "127.0.0.1", history[0].Object.Spec["endpoint"].(map[string]any)["address"])

	_, err = store.Rollback(ctx, draft.ID, history[0].Revision, current.Revision)
	require.ErrorIs(t, err, ErrClusterPublishConflict)
	rolled, err := store.Rollback(ctx, draft.ID, history[0].Revision, second.DraftRevision)
	require.NoError(t, err)
	require.Equal(t, first.RuntimeID, rolled.RuntimeID)
	require.Equal(t, "127.0.0.1", runtimeClusterAt(t, kv, store.runtimeKey(first.RuntimeID)).Address)
	history, err = store.History(ctx, draft.ID)
	require.NoError(t, err)
	require.Len(t, history, 3)
	require.Greater(t, history[2].Revision, history[1].Revision)
	require.NoError(t, store.DeleteDraft(ctx, draft.ID, rolled.DraftRevision))
	require.True(t, kv.fakeHasKey(store.runtimeKey(first.RuntimeID)), "draft deletion must preserve the published cluster")
	reopened, err := store.CreateFromPublished(ctx, first.RuntimeID)
	require.NoError(t, err)
	require.NotEqual(t, draft.ID, reopened.ID)
	history, err = store.History(ctx, reopened.ID)
	require.NoError(t, err)
	require.Len(t, history, 3, "a new draft for the same runtime Cluster must retain its history")
	republished, err := store.Rollback(ctx, reopened.ID, history[1].Revision, reopened.Revision)
	require.NoError(t, err)
	require.Equal(t, first.RuntimeID, republished.RuntimeID)
	require.Equal(t, "127.0.0.2", runtimeClusterAt(t, kv, store.runtimeKey(first.RuntimeID)).Address)
}

func TestClusterPublishRejectsNameAndLegacyWriteConflicts(t *testing.T) {
	store, kv := newTestClusterDraftStore(t)
	ctx := context.Background()
	legacy := &legacyconfig.Cluster{Name: "taken", Type: "Static", Address: "127.0.0.1", Port: 20880}
	require.NoError(t, store.CreateRuntime(ctx, legacy))
	require.Equal(t, 1, legacy.ID)
	draft, err := store.CreateDraft(ctx, testClusterObject("taken", "127.0.0.2", 20880))
	require.NoError(t, err)
	_, err = store.Publish(ctx, draft.ID, draft.Revision)
	require.ErrorIs(t, err, ErrClusterNameConflict)
	require.False(t, kv.fakeHasPrefix(store.root+"/"+clusterHistorySegment+"/"))
	spaced, err := store.CreateDraft(ctx, testClusterObject(" taken ", "127.0.0.2", 20880))
	require.NoError(t, err)
	_, err = store.Publish(ctx, spaced.ID, spaced.Revision)
	require.ErrorIs(t, err, ErrClusterNameConflict)
	require.False(t, kv.fakeHasKey(store.runtimeKey(legacy.ID+1)), "a whitespace-equivalent name must not reach runtime")
	require.False(t, kv.fakeHasPrefix(store.root+"/"+clusterHistorySegment+"/"))

	fromPublished, err := store.CreateFromPublished(ctx, legacy.ID)
	require.NoError(t, err)
	legacy.Address = "127.0.0.3"
	require.NoError(t, store.UpdateRuntime(ctx, legacy))
	_, err = store.Publish(ctx, fromPublished.ID, fromPublished.Revision)
	require.ErrorIs(t, err, ErrClusterPublishConflict)
	require.Equal(t, "127.0.0.3", runtimeClusterAt(t, kv, store.runtimeKey(legacy.ID)).Address)
}

func TestClusterPreviewNameAvailabilityUsesDraftSourceIdentity(t *testing.T) {
	store, _ := newTestClusterDraftStore(t)
	ctx := context.Background()
	legacy := &legacyconfig.Cluster{Name: "taken", Type: "Static", Address: "127.0.0.1", Port: 20880}
	require.NoError(t, store.CreateRuntime(ctx, legacy))
	checkDuplicate := func(id string) {
		t.Helper()
		err := store.ValidateNameAvailability(ctx, id, "taken")
		var issues schema.ValidationErrors
		require.ErrorAs(t, err, &issues)
		require.Equal(t, "metadata.name", issues[0].Path)
	}
	checkDuplicate("")
	var spacedIssues schema.ValidationErrors
	require.ErrorAs(t, store.ValidateNameAvailability(ctx, "", " taken "), &spacedIssues)
	unlinked, err := store.CreateDraft(ctx, testClusterObject("taken", "127.0.0.2", 20881))
	require.NoError(t, err)
	checkDuplicate(unlinked.ID)
	linked, err := store.CreateFromPublished(ctx, legacy.ID)
	require.NoError(t, err)
	require.NoError(t, store.ValidateNameAvailability(ctx, linked.ID, "taken"))
	require.NoError(t, store.ValidateNameAvailability(ctx, linked.ID, "other"))
}

func TestClusterDeleteRejectsReferencedRuntimeAndConcurrentListener(t *testing.T) {
	store, kv := newTestClusterDraftStore(t)
	ctx := context.Background()
	cluster := &legacyconfig.Cluster{Name: "backend", Type: "Static", Address: "127.0.0.1", Port: 20880}
	require.NoError(t, store.CreateRuntime(ctx, cluster))
	var listener legacyconfig.Listener
	require.NoError(t, yaml.Unmarshal([]byte("name: gateway\nroute_config:\n  routes:\n    - route:\n        cluster: backend\n"), &listener))
	require.NoError(t, store.WriteListener(ctx, &listener, true))
	rename := *cluster
	rename.Name = " backend "
	require.ErrorIs(t, store.UpdateRuntime(ctx, &rename), ErrClusterReferenced)
	require.ErrorIs(t, store.DeleteRuntime(ctx, cluster.ID), ErrClusterReferenced)
	require.True(t, kv.fakeHasKey(store.runtimeKey(cluster.ID)))
	require.NoError(t, store.DeleteListener(ctx, listener.Name))

	triggered := false
	guarded := *store
	guarded.kv = &hookedRouteBindingKV{routeBindingKV: kv, afterGet: func(key string) {
		if triggered || key != store.root+"/"+Listeners+"/" {
			return
		}
		triggered = true
		require.NoError(t, store.WriteListener(ctx, &listener, true))
	}}
	err := guarded.DeleteRuntime(ctx, cluster.ID)
	require.True(t, triggered)
	require.True(t, errors.Is(err, ErrClusterPublishConflict) || errors.Is(err, ErrClusterReferenced))
	require.True(t, kv.fakeHasKey(store.runtimeKey(cluster.ID)))
	require.NoError(t, store.DeleteListener(ctx, listener.Name))
	require.NoError(t, store.DeleteRuntime(ctx, cluster.ID))
	require.False(t, kv.fakeHasKey(store.runtimeKey(cluster.ID)))
}

func TestClusterRollbackRestoresDeletedRuntimeWithoutChangingItsID(t *testing.T) {
	store, kv := newTestClusterDraftStore(t)
	ctx := context.Background()
	draft, err := store.CreateDraft(ctx, testClusterObject("backend", "127.0.0.1", 20880))
	require.NoError(t, err)
	published, err := store.Publish(ctx, draft.ID, draft.Revision)
	require.NoError(t, err)
	require.NoError(t, store.DeleteRuntime(ctx, published.RuntimeID))
	_, err = store.Publish(ctx, draft.ID, published.DraftRevision)
	require.ErrorIs(t, err, ErrClusterPublishConflict)
	require.False(t, kv.fakeHasKey(store.runtimeKey(published.RuntimeID)))

	restored, err := store.Rollback(ctx, draft.ID, published.HistoryRevision, published.DraftRevision)
	require.NoError(t, err)
	require.Equal(t, published.RuntimeID, restored.RuntimeID)
	require.Equal(t, "backend", runtimeClusterAt(t, kv, store.runtimeKey(restored.RuntimeID)).Name)
	history, err := store.History(ctx, draft.ID)
	require.NoError(t, err)
	require.Len(t, history, 2)
}

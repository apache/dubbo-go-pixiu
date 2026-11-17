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

func newTestClusterDraftStore(t *testing.T) (*ClusterDraftStore, *fakeRouteBindingKV) {
	t.Helper()
	kv := newFakeRouteBindingKV()
	store, err := NewClusterDraftStore(kv, "/pixiu/config/api")
	require.NoError(t, err)
	return store, kv
}

func testClusterObject(name, address string, port any) schema.AdminObject {
	return schema.AdminObject{Kind: schema.KindAdminCluster, Metadata: schema.ObjectMetadata{Name: name}, Spec: map[string]any{
		"type": "Static", "endpoint": map[string]any{"address": address, "port": port},
	}}
}

func TestClusterDraftSavesInvalidObjectAndRejectsPreview(t *testing.T) {
	store, kv := newTestClusterDraftStore(t)
	object := schema.AdminObject{Kind: schema.KindAdminCluster, Spec: map[string]any{"type": "invalid"}}
	created, err := store.CreateDraft(context.Background(), object)
	require.NoError(t, err)
	require.Len(t, created.ID, 32)
	require.Positive(t, created.Revision)
	require.True(t, kv.fakeHasKey(store.draftKey(created.ID)))
	require.False(t, kv.fakeHasPrefix(store.root+"/"+Clusters+"/"))
	read, err := store.GetDraft(context.Background(), created.ID)
	require.NoError(t, err)
	require.Equal(t, object, read.Object)
	listed, err := store.ListDrafts(context.Background())
	require.NoError(t, err)
	require.Equal(t, []ClusterDraft{read}, listed)

	registry, err := schema.NewBuiltinRegistry()
	require.NoError(t, err)
	_, err = registry.Normalize(read.Object)
	var issues schema.ValidationErrors
	require.ErrorAs(t, err, &issues)
	require.Contains(t, err.Error(), "metadata.name")
	_, err = schema.CompileAdminCluster(registry, read.Object)
	require.ErrorAs(t, err, &issues)
	status, err := store.Status(context.Background(), created.ID)
	require.NoError(t, err)
	require.True(t, status.DraftExists)
	require.False(t, status.PublishedExists)
	require.True(t, status.Dirty)
	diff, err := store.Diff(context.Background(), created.ID)
	require.NoError(t, err)
	require.Nil(t, diff.Published)
	require.NotEmpty(t, diff.Changes)
}

func TestClusterDraftFromPublishedStatusDiffAndRevision(t *testing.T) {
	store, kv := newTestClusterDraftStore(t)
	published := legacyconfig.Cluster{Name: "backend", Type: "Static", Address: "127.0.0.1", Port: 20880, ID: 17}
	value, err := yaml.Marshal(published)
	require.NoError(t, err)
	kv.fakePut(store.runtimeKey(17), value)
	draft, err := store.CreateFromPublished(context.Background(), 17)
	require.NoError(t, err)
	require.Equal(t, 17, draft.SourceID)
	status, err := store.Status(context.Background(), draft.ID)
	require.NoError(t, err)
	require.True(t, status.PublishedExists)
	require.False(t, status.Dirty)
	require.Positive(t, status.PublishedRevision)
	diff, err := store.Diff(context.Background(), draft.ID)
	require.NoError(t, err)
	require.NotNil(t, diff.Published)
	require.Empty(t, diff.Changes)

	modified := testClusterObject("backend", "127.0.0.2", 20880)
	updated, err := store.UpdateDraft(context.Background(), draft.ID, modified, draft.Revision)
	require.NoError(t, err)
	require.Greater(t, updated.Revision, draft.Revision)
	require.Equal(t, 17, updated.SourceID)
	status, err = store.Status(context.Background(), draft.ID)
	require.NoError(t, err)
	require.True(t, status.Dirty)
	diff, err = store.Diff(context.Background(), draft.ID)
	require.NoError(t, err)
	require.Equal(t, []ClusterDiffChange{{Path: "spec.endpoint.address", Before: "127.0.0.1", After: "127.0.0.2"}}, diff.Changes)
	require.Equal(t, value, kv.fakeValue(store.runtimeKey(17)), "draft writes must leave runtime keys unchanged")

	_, err = store.UpdateDraft(context.Background(), draft.ID, modified, draft.Revision)
	require.ErrorIs(t, err, ErrClusterDraftConflict)
	_, err = store.UpdateDraft(context.Background(), draft.ID, modified, 0)
	require.ErrorContains(t, err, "expectedRevision")

	// A direct legacy write remains visible in status and diff by source ID.
	changedPublished := published
	changedPublished.Address = "127.0.0.3"
	value, err = yaml.Marshal(changedPublished)
	require.NoError(t, err)
	kv.fakePut(store.runtimeKey(17), value)
	diff, err = store.Diff(context.Background(), draft.ID)
	require.NoError(t, err)
	require.Equal(t, "127.0.0.3", diff.Changes[0].Before)
}

func TestClusterDraftUpdateRejectsConcurrentWrite(t *testing.T) {
	store, kv := newTestClusterDraftStore(t)
	created, err := store.CreateDraft(context.Background(), testClusterObject("backend", "127.0.0.1", 20880))
	require.NoError(t, err)
	triggered := false
	guarded := *store
	guarded.kv = &hookedRouteBindingKV{routeBindingKV: kv, afterGet: func(key string) {
		if triggered || key != store.draftKey(created.ID) {
			return
		}
		triggered = true
		_, updateErr := store.UpdateDraft(context.Background(), created.ID, testClusterObject("backend", "new", 20880), created.Revision)
		require.NoError(t, updateErr)
	}}
	_, err = guarded.UpdateDraft(context.Background(), created.ID, testClusterObject("backend", "stale", 20880), created.Revision)
	require.True(t, triggered)
	require.ErrorIs(t, err, ErrClusterDraftConflict)
	current, err := store.GetDraft(context.Background(), created.ID)
	require.NoError(t, err)
	require.Equal(t, "new", current.Object.Spec["endpoint"].(map[string]any)["address"])
}

func TestClusterDraftIDCollisionAndInvalidInputs(t *testing.T) {
	store, kv := newTestClusterDraftStore(t)
	firstID := "aaaaaaaaaaaaaaaaaaaaaaaaaaaaaaaa"
	secondID := "bbbbbbbbbbbbbbbbbbbbbbbbbbbbbbbb"
	kv.fakePut(store.draftKey(firstID), []byte(`{"object":{"kind":"AdminCluster"}}`))
	ids := []string{firstID, secondID}
	store.newID = func() (string, error) { id := ids[0]; ids = ids[1:]; return id, nil }
	created, err := store.CreateDraft(context.Background(), schema.AdminObject{})
	require.NoError(t, err)
	require.Equal(t, secondID, created.ID)
	require.Equal(t, schema.KindAdminCluster, created.Object.Kind)
	_, err = store.CreateDraft(context.Background(), schema.AdminObject{Kind: schema.KindAdminRouteBinding})
	require.ErrorContains(t, err, "kind")
	_, err = store.GetDraft(context.Background(), "../clusters/17")
	require.ErrorContains(t, err, "id is invalid")
	_, err = store.GetDraft(context.Background(), "cccccccccccccccccccccccccccccccc")
	require.ErrorIs(t, err, ErrClusterDraftNotFound)
	_, err = store.CreateFromPublished(context.Background(), 100)
	require.ErrorIs(t, err, ErrClusterRuntimeNotFound)
	store.newID = func() (string, error) { return "", errors.New("entropy unavailable") }
	_, err = store.CreateDraft(context.Background(), schema.AdminObject{})
	require.ErrorContains(t, err, "entropy unavailable")
}

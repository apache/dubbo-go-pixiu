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
	"crypto/sha256"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"sort"
	"strconv"
	"strings"
	"time"

	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
	"gopkg.in/yaml.v3"

	legacyconfig "github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema"
)

const clusterHistorySegment = "history/clusters"

var (
	ErrClusterPublishConflict = errors.New("cluster publish conflict")
	ErrClusterNameConflict    = errors.New("runtime cluster name already exists")
	ErrClusterReferenced      = errors.New("cluster is referenced by a published listener")
)

type ClusterPublishResult struct {
	ID                string `json:"id"`
	RuntimeID         int    `json:"runtimeId"`
	DraftRevision     int64  `json:"draftRevision"`
	PublishedRevision int64  `json:"publishedRevision"`
	HistoryRevision   int64  `json:"historyRevision"`
}

type ClusterHistoryEntry struct {
	Revision    int64              `json:"revision"`
	PublishedAt time.Time          `json:"publishedAt"`
	RuntimeID   int                `json:"runtimeId"`
	Object      schema.AdminObject `json:"object"`
}

type clusterRuntimeEntry struct {
	cluster legacyconfig.Cluster
	kv      *mvccpb.KeyValue
}

func clusterRuntimeDigest(cluster legacyconfig.Cluster) string {
	value, _ := yaml.Marshal(cluster)
	sum := sha256.Sum256(value)
	return hex.EncodeToString(sum[:])
}

func (s *ClusterDraftStore) lifecycleKey() string { return s.root + "/unpublished/cluster-generation" }
func (s *ClusterDraftStore) counterKey() string   { return s.root + "/" + ClusterID }
func (s *ClusterDraftStore) historyPrefix(runtimeID int) string {
	return s.root + "/" + clusterHistorySegment + "/" + strconv.Itoa(runtimeID) + "/"
}
func (s *ClusterDraftStore) historyKey(runtimeID int, revision int64) string {
	return s.historyPrefix(runtimeID) + fmt.Sprintf("%020d", revision)
}

func (s *ClusterDraftStore) readNumber(ctx context.Context, key string) (int64, int64, error) {
	response, err := s.kv.Get(ctx, key)
	if err != nil {
		return 0, 0, fmt.Errorf("read %q: %w", key, err)
	}
	if len(response.Kvs) == 0 {
		return 0, 0, nil
	}
	value, err := strconv.ParseInt(string(response.Kvs[0].Value), 10, 64)
	if err != nil || value < 0 {
		return 0, 0, fmt.Errorf("invalid counter at %q", key)
	}
	return value, response.Kvs[0].ModRevision, nil
}

func compareRevision(key string, revision int64) clientv3.Cmp {
	if revision == 0 {
		return clientv3.Compare(clientv3.CreateRevision(key), "=", 0)
	}
	return clientv3.Compare(clientv3.ModRevision(key), "=", revision)
}

func (s *ClusterDraftStore) runtimeEntries(ctx context.Context) ([]clusterRuntimeEntry, error) {
	prefix := s.root + "/" + Clusters + "/"
	response, err := s.kv.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("list runtime clusters: %w", err)
	}
	entries := make([]clusterRuntimeEntry, 0, len(response.Kvs))
	for _, kv := range response.Kvs {
		idPart := strings.TrimPrefix(string(kv.Key), prefix)
		id, parseErr := strconv.Atoi(idPart)
		if parseErr != nil || id <= 0 {
			continue
		}
		var cluster legacyconfig.Cluster
		if err := yaml.Unmarshal(kv.Value, &cluster); err != nil {
			return nil, fmt.Errorf("decode runtime cluster %d: %w", id, err)
		}
		if cluster.ID != id {
			return nil, fmt.Errorf("runtime cluster key %d contains ID %d", id, cluster.ID)
		}
		entries = append(entries, clusterRuntimeEntry{cluster: cluster, kv: kv})
	}
	return entries, nil
}

func (s *ClusterDraftStore) listenerReferences(ctx context.Context, clusterName string) ([]clientv3.Cmp, bool, error) {
	prefix := s.root + "/" + Listeners + "/"
	response, err := s.kv.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return nil, false, fmt.Errorf("list published listeners: %w", err)
	}
	comparisons := make([]clientv3.Cmp, 0, len(response.Kvs))
	referenced := false
	for _, kv := range response.Kvs {
		if strings.Contains(strings.TrimPrefix(string(kv.Key), prefix), "/") {
			continue
		}
		var listener legacyconfig.Listener
		if err := yaml.Unmarshal(kv.Value, &listener); err != nil {
			return nil, false, fmt.Errorf("decode listener %q: %w", string(kv.Key), err)
		}
		comparisons = append(comparisons, compareRevision(string(kv.Key), kv.ModRevision))
		for _, route := range listener.RouteConfig.Routes {
			if strings.TrimSpace(route.Route.Cluster) == clusterName {
				referenced = true
			}
		}
	}
	return comparisons, referenced, nil
}

// ValidateNameAvailability checks the current runtime set before a preview.
// The optional draft ID resolves the server-owned source identity, so an edit
// may keep its existing name without trusting a caller-supplied runtime ID.
func (s *ClusterDraftStore) ValidateNameAvailability(ctx context.Context, draftID, name string) error {
	sourceID := 0
	if draftID != "" {
		draft, err := s.GetDraft(ctx, draftID)
		if err != nil {
			return err
		}
		sourceID = draft.SourceID
	}
	entries, err := s.runtimeEntries(routeBindingContext(ctx))
	if err != nil {
		return err
	}
	for _, entry := range entries {
		if entry.cluster.ID != sourceID && strings.TrimSpace(entry.cluster.Name) == strings.TrimSpace(name) {
			return schema.ValidationErrors{{Path: "metadata.name", Code: "duplicate", Message: "runtime cluster name already exists"}}
		}
	}
	return nil
}

func (s *ClusterDraftStore) Publish(ctx context.Context, id string, expectedRevision int64) (ClusterPublishResult, error) {
	draft, err := s.GetDraft(ctx, id)
	if err != nil {
		return ClusterPublishResult{}, err
	}
	return s.publish(ctx, draft, draft.Object, expectedRevision, false)
}

func (s *ClusterDraftStore) publish(ctx context.Context, draft ClusterDraft, object schema.AdminObject, expectedRevision int64, restoreMissing bool) (ClusterPublishResult, error) {
	if expectedRevision <= 0 {
		return ClusterPublishResult{}, errors.New("positive expectedRevision is required")
	}
	if draft.Revision != expectedRevision {
		return ClusterPublishResult{}, fmt.Errorf("%w: draft revision changed", ErrClusterPublishConflict)
	}
	registry, err := schema.NewBuiltinRegistry()
	if err != nil {
		return ClusterPublishResult{}, err
	}
	compiled, err := schema.CompileAdminCluster(registry, object)
	if err != nil {
		return ClusterPublishResult{}, err
	}
	ctx = routeBindingContext(ctx)
	generation, generationRevision, err := s.readNumber(ctx, s.lifecycleKey())
	if err != nil {
		return ClusterPublishResult{}, err
	}
	entries, err := s.runtimeEntries(ctx)
	if err != nil {
		return ClusterPublishResult{}, err
	}
	var source *clusterRuntimeEntry
	maxID := 0
	comparisons := []clientv3.Cmp{compareRevision(s.lifecycleKey(), generationRevision), compareRevision(s.draftKey(draft.ID), expectedRevision)}
	for i := range entries {
		entry := &entries[i]
		if entry.cluster.ID > maxID {
			maxID = entry.cluster.ID
		}
		comparisons = append(comparisons, compareRevision(string(entry.kv.Key), entry.kv.ModRevision))
		if entry.cluster.ID == draft.SourceID {
			source = entry
		}
		if strings.TrimSpace(entry.cluster.Name) == strings.TrimSpace(compiled.Cluster.Name) && entry.cluster.ID != draft.SourceID {
			return ClusterPublishResult{}, fmt.Errorf("%w: %q", ErrClusterNameConflict, compiled.Cluster.Name)
		}
	}
	runtimeID := draft.SourceID
	operations := make([]clientv3.Op, 0, 5)
	if runtimeID > 0 {
		if source == nil && !restoreMissing {
			return ClusterPublishResult{}, fmt.Errorf("%w: source cluster %d was removed", ErrClusterPublishConflict, runtimeID)
		}
		if source != nil && draft.sourceDigest != "" && clusterRuntimeDigest(source.cluster) != draft.sourceDigest {
			return ClusterPublishResult{}, fmt.Errorf("%w: source cluster %d changed", ErrClusterPublishConflict, runtimeID)
		}
		if source == nil {
			comparisons = append(comparisons, compareRevision(s.runtimeKey(runtimeID), 0))
		}
	} else {
		counter, counterRevision, err := s.readNumber(ctx, s.counterKey())
		if err != nil {
			return ClusterPublishResult{}, err
		}
		if counter > int64(maxID) {
			maxID = int(counter)
		}
		runtimeID = maxID + 1
		comparisons = append(comparisons, compareRevision(s.counterKey(), counterRevision), compareRevision(s.runtimeKey(runtimeID), 0))
		operations = append(operations, clientv3.OpPut(s.counterKey(), strconv.Itoa(runtimeID)))
	}
	if source != nil && source.cluster.Name != compiled.Cluster.Name {
		listenerComparisons, referenced, err := s.listenerReferences(ctx, strings.TrimSpace(source.cluster.Name))
		if err != nil {
			return ClusterPublishResult{}, err
		}
		if referenced {
			return ClusterPublishResult{}, fmt.Errorf("%w: %q", ErrClusterReferenced, source.cluster.Name)
		}
		comparisons = append(comparisons, listenerComparisons...)
	}
	runtime := compiled.RuntimeCluster(runtimeID)
	runtimeYAML, err := yaml.Marshal(runtime)
	if err != nil {
		return ClusterPublishResult{}, err
	}
	historyRevision := generation + 1
	history := ClusterHistoryEntry{Revision: historyRevision, PublishedAt: time.Now().UTC(), RuntimeID: runtimeID, Object: compiled.Source}
	historyValue, err := json.Marshal(history)
	if err != nil {
		return ClusterPublishResult{}, err
	}
	draftValue, err := json.Marshal(clusterDraftRecord{Object: compiled.Source, SourceID: runtimeID, SourceDigest: clusterRuntimeDigest(runtime)})
	if err != nil {
		return ClusterPublishResult{}, err
	}
	operations = append(operations,
		clientv3.OpPut(s.runtimeKey(runtimeID), string(runtimeYAML)),
		clientv3.OpPut(s.draftKey(draft.ID), string(draftValue)),
		clientv3.OpPut(s.historyKey(runtimeID, historyRevision), string(historyValue)),
		clientv3.OpPut(s.lifecycleKey(), strconv.FormatInt(historyRevision, 10)))
	response, err := s.kv.Txn(ctx).If(comparisons...).Then(operations...).Commit()
	if err != nil {
		return ClusterPublishResult{}, fmt.Errorf("publish cluster: %w", err)
	}
	if !response.Succeeded {
		return ClusterPublishResult{}, fmt.Errorf("%w: draft, runtime, or listener changed", ErrClusterPublishConflict)
	}
	return ClusterPublishResult{ID: draft.ID, RuntimeID: runtimeID, DraftRevision: response.Header.Revision, PublishedRevision: response.Header.Revision, HistoryRevision: historyRevision}, nil
}

func (s *ClusterDraftStore) History(ctx context.Context, id string) ([]ClusterHistoryEntry, error) {
	draft, err := s.GetDraft(ctx, id)
	if err != nil {
		return nil, err
	}
	if draft.SourceID == 0 {
		return []ClusterHistoryEntry{}, nil
	}
	response, err := s.kv.Get(routeBindingContext(ctx), s.historyPrefix(draft.SourceID), clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("list cluster history: %w", err)
	}
	entries := make([]ClusterHistoryEntry, 0, len(response.Kvs))
	for _, kv := range response.Kvs {
		var entry ClusterHistoryEntry
		if err := json.Unmarshal(kv.Value, &entry); err != nil {
			return nil, fmt.Errorf("decode cluster history: %w", err)
		}
		entries = append(entries, entry)
	}
	sort.Slice(entries, func(i, j int) bool { return entries[i].Revision < entries[j].Revision })
	return entries, nil
}

func (s *ClusterDraftStore) Rollback(ctx context.Context, id string, historyRevision, expectedDraftRevision int64) (ClusterPublishResult, error) {
	if historyRevision <= 0 {
		return ClusterPublishResult{}, errors.New("positive history revision is required")
	}
	draft, err := s.GetDraft(ctx, id)
	if err != nil {
		return ClusterPublishResult{}, err
	}
	response, err := s.kv.Get(routeBindingContext(ctx), s.historyKey(draft.SourceID, historyRevision))
	if err != nil {
		return ClusterPublishResult{}, fmt.Errorf("get cluster history: %w", err)
	}
	if len(response.Kvs) == 0 {
		return ClusterPublishResult{}, errors.New("cluster history revision not found")
	}
	var entry ClusterHistoryEntry
	if err := json.Unmarshal(response.Kvs[0].Value, &entry); err != nil {
		return ClusterPublishResult{}, fmt.Errorf("decode cluster history: %w", err)
	}
	if entry.Revision != historyRevision || entry.RuntimeID != draft.SourceID {
		return ClusterPublishResult{}, fmt.Errorf("%w: history does not match draft source", ErrClusterPublishConflict)
	}
	return s.publish(ctx, draft, entry.Object, expectedDraftRevision, true)
}

func (s *ClusterDraftStore) DeleteDraft(ctx context.Context, id string, expectedRevision int64) error {
	if !validClusterDraftID(id) || expectedRevision <= 0 {
		return errors.New("valid draft id and positive expectedRevision are required")
	}
	response, err := s.kv.Txn(routeBindingContext(ctx)).If(compareRevision(s.draftKey(id), expectedRevision)).Then(clientv3.OpDelete(s.draftKey(id))).Commit()
	if err != nil {
		return fmt.Errorf("delete cluster draft: %w", err)
	}
	if !response.Succeeded {
		return fmt.Errorf("%w: draft revision changed", ErrClusterDraftConflict)
	}
	return nil
}

// Legacy CRUD writes use the same generation guard as lifecycle publication.
// This lets draft publication detect a direct Admin API edit before overwriting it.
func (s *ClusterDraftStore) CreateRuntime(ctx context.Context, cluster *legacyconfig.Cluster) error {
	if cluster == nil || strings.TrimSpace(cluster.Name) == "" {
		return errors.New("cluster name is required")
	}
	ctx = routeBindingContext(ctx)
	generation, generationRevision, err := s.readNumber(ctx, s.lifecycleKey())
	if err != nil {
		return err
	}
	entries, err := s.runtimeEntries(ctx)
	if err != nil {
		return err
	}
	maxID := 0
	comparisons := []clientv3.Cmp{compareRevision(s.lifecycleKey(), generationRevision)}
	for _, entry := range entries {
		if strings.TrimSpace(entry.cluster.Name) == strings.TrimSpace(cluster.Name) {
			return fmt.Errorf("%w: %q", ErrClusterNameConflict, cluster.Name)
		}
		if entry.cluster.ID > maxID {
			maxID = entry.cluster.ID
		}
		comparisons = append(comparisons, compareRevision(string(entry.kv.Key), entry.kv.ModRevision))
	}
	counter, counterRevision, err := s.readNumber(ctx, s.counterKey())
	if err != nil {
		return err
	}
	if counter > int64(maxID) {
		maxID = int(counter)
	}
	newID := maxID + 1
	clusterValue := *cluster
	clusterValue.ID = newID
	value, err := yaml.Marshal(clusterValue)
	if err != nil {
		return err
	}
	comparisons = append(comparisons, compareRevision(s.counterKey(), counterRevision), compareRevision(s.runtimeKey(newID), 0))
	response, err := s.kv.Txn(ctx).If(comparisons...).Then(
		clientv3.OpPut(s.counterKey(), strconv.Itoa(newID)),
		clientv3.OpPut(s.runtimeKey(newID), string(value)),
		clientv3.OpPut(s.lifecycleKey(), strconv.FormatInt(generation+1, 10)),
	).Commit()
	if err != nil {
		return fmt.Errorf("create runtime cluster: %w", err)
	}
	if !response.Succeeded {
		return ErrClusterPublishConflict
	}
	cluster.ID = newID
	return nil
}

func (s *ClusterDraftStore) UpdateRuntime(ctx context.Context, cluster *legacyconfig.Cluster) error {
	if cluster == nil || cluster.ID <= 0 || strings.TrimSpace(cluster.Name) == "" {
		return errors.New("valid cluster ID and name are required")
	}
	ctx = routeBindingContext(ctx)
	generation, generationRevision, err := s.readNumber(ctx, s.lifecycleKey())
	if err != nil {
		return err
	}
	entries, err := s.runtimeEntries(ctx)
	if err != nil {
		return err
	}
	var source *clusterRuntimeEntry
	comparisons := []clientv3.Cmp{compareRevision(s.lifecycleKey(), generationRevision)}
	for i := range entries {
		entry := &entries[i]
		comparisons = append(comparisons, compareRevision(string(entry.kv.Key), entry.kv.ModRevision))
		if entry.cluster.ID == cluster.ID {
			source = entry
		}
		if entry.cluster.ID != cluster.ID && strings.TrimSpace(entry.cluster.Name) == strings.TrimSpace(cluster.Name) {
			return fmt.Errorf("%w: %q", ErrClusterNameConflict, cluster.Name)
		}
	}
	if source == nil {
		return fmt.Errorf("%w: %d", ErrClusterRuntimeNotFound, cluster.ID)
	}
	if source.cluster.Name != cluster.Name {
		listenerComparisons, referenced, err := s.listenerReferences(ctx, strings.TrimSpace(source.cluster.Name))
		if err != nil {
			return err
		}
		if referenced {
			return fmt.Errorf("%w: %q", ErrClusterReferenced, source.cluster.Name)
		}
		comparisons = append(comparisons, listenerComparisons...)
	}
	value, err := yaml.Marshal(cluster)
	if err != nil {
		return err
	}
	response, err := s.kv.Txn(ctx).If(comparisons...).Then(
		clientv3.OpPut(s.runtimeKey(cluster.ID), string(value)),
		clientv3.OpPut(s.lifecycleKey(), strconv.FormatInt(generation+1, 10)),
	).Commit()
	if err != nil {
		return fmt.Errorf("update runtime cluster: %w", err)
	}
	if !response.Succeeded {
		return ErrClusterPublishConflict
	}
	return nil
}

func (s *ClusterDraftStore) DeleteRuntime(ctx context.Context, id int) error {
	if id <= 0 {
		return errors.New("valid cluster ID is required")
	}
	ctx = routeBindingContext(ctx)
	generation, generationRevision, err := s.readNumber(ctx, s.lifecycleKey())
	if err != nil {
		return err
	}
	published, revision, exists, err := s.getPublished(ctx, id)
	if err != nil {
		return err
	}
	if !exists {
		return fmt.Errorf("%w: %d", ErrClusterRuntimeNotFound, id)
	}
	listenerComparisons, referenced, err := s.listenerReferences(ctx, strings.TrimSpace(published.Name))
	if err != nil {
		return err
	}
	if referenced {
		return fmt.Errorf("%w: %q", ErrClusterReferenced, published.Name)
	}
	comparisons := append([]clientv3.Cmp{compareRevision(s.lifecycleKey(), generationRevision), compareRevision(s.runtimeKey(id), revision)}, listenerComparisons...)
	response, err := s.kv.Txn(ctx).If(comparisons...).Then(
		clientv3.OpDelete(s.runtimeKey(id)),
		clientv3.OpPut(s.lifecycleKey(), strconv.FormatInt(generation+1, 10)),
	).Commit()
	if err != nil {
		return fmt.Errorf("delete runtime cluster: %w", err)
	}
	if !response.Succeeded {
		return ErrClusterPublishConflict
	}
	return nil
}

// Listener CRUD also advances the shared guard, so a new reference cannot
// appear between the deletion check and the transaction commit.
func (s *ClusterDraftStore) WriteListener(ctx context.Context, listener *legacyconfig.Listener, create bool) error {
	if listener == nil || strings.TrimSpace(listener.Name) == "" || strings.Contains(listener.Name, "/") {
		return errors.New("valid listener name is required")
	}
	ctx = routeBindingContext(ctx)
	generation, generationRevision, err := s.readNumber(ctx, s.lifecycleKey())
	if err != nil {
		return err
	}
	key := s.root + "/" + Listeners + "/" + listener.Name
	response, err := s.kv.Get(ctx, key)
	if err != nil {
		return err
	}
	revision := int64(0)
	if len(response.Kvs) > 0 {
		revision = response.Kvs[0].ModRevision
	}
	if create && revision > 0 {
		return errors.New("listener already exists")
	}
	if !create && revision == 0 {
		return errors.New("listener not found")
	}
	value, err := yaml.Marshal(listener)
	if err != nil {
		return err
	}
	result, err := s.kv.Txn(ctx).If(compareRevision(s.lifecycleKey(), generationRevision), compareRevision(key, revision)).Then(
		clientv3.OpPut(key, string(value)),
		clientv3.OpPut(s.lifecycleKey(), strconv.FormatInt(generation+1, 10)),
	).Commit()
	if err != nil {
		return err
	}
	if !result.Succeeded {
		return ErrClusterPublishConflict
	}
	return nil
}

func (s *ClusterDraftStore) DeleteListener(ctx context.Context, name string) error {
	if strings.TrimSpace(name) == "" || strings.Contains(name, "/") {
		return errors.New("valid listener name is required")
	}
	ctx = routeBindingContext(ctx)
	generation, generationRevision, err := s.readNumber(ctx, s.lifecycleKey())
	if err != nil {
		return err
	}
	key := s.root + "/" + Listeners + "/" + name
	response, err := s.kv.Get(ctx, key)
	if err != nil {
		return err
	}
	if len(response.Kvs) == 0 {
		return nil
	}
	result, err := s.kv.Txn(ctx).If(compareRevision(s.lifecycleKey(), generationRevision), compareRevision(key, response.Kvs[0].ModRevision)).Then(
		clientv3.OpDelete(key),
		clientv3.OpPut(s.lifecycleKey(), strconv.FormatInt(generation+1, 10)),
	).Commit()
	if err != nil {
		return err
	}
	if !result.Succeeded {
		return ErrClusterPublishConflict
	}
	return nil
}

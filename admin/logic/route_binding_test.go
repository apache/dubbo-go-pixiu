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
	"sort"
	"strings"
	"sync"
	"testing"
	"time"
)

import (
	pb "go.etcd.io/etcd/api/v3/etcdserverpb"
	"go.etcd.io/etcd/api/v3/mvccpb"
	clientv3 "go.etcd.io/etcd/client/v3"
)

import (
	commonyaml "github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	legacyconfig "github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema"
)

func TestRouteBindingStoreSavePublishAndDelete(t *testing.T) {
	store := newTestRouteBindingStore(t)

	saved, err := store.SaveDraft(testRouteBinding("user-get", "/api/users/:id", "GET"), true, 0)
	if err != nil {
		t.Fatalf("SaveDraft: %v", err)
	}
	if saved.ResourceID != 1 || saved.MethodID != 1 {
		t.Fatalf("runtime identity: want 1/1, got %d/%d", saved.ResourceID, saved.MethodID)
	}
	if saved.Revision == 0 {
		t.Fatal("SaveDraft returned zero key revision")
	}

	status, err := store.PublishStatus()
	if err != nil {
		t.Fatalf("PublishStatus: %v", err)
	}
	result, err := store.PublishAll(status.DraftRevision)
	if err != nil {
		t.Fatalf("PublishAll: %v", err)
	}
	if result.PublishedCount != 1 || result.DeletedCount != 0 {
		t.Fatalf("publish result: %+v", result)
	}
	if result.Revision == 0 || result.DraftRevision != status.DraftRevision {
		t.Fatalf("publish revisions: %+v, status=%+v", result, status)
	}

	fake := fakeRouteBindingStoreKV(t, store)
	resourceValue := fake.fakeValue(store.runtimeResourceKey(saved.ResourceID))
	methodValue := fake.fakeValue(store.runtimeMethodKey(saved.ResourceID, saved.MethodID))
	var resource legacyconfig.Resource
	if err := commonyaml.UnmarshalYML(resourceValue, &resource); err != nil {
		t.Fatalf("decode generated Resource: %v\n%s", err, resourceValue)
	}
	var method legacyconfig.Method
	if err := commonyaml.UnmarshalYML(methodValue, &method); err != nil {
		t.Fatalf("decode generated Method: %v\n%s", err, methodValue)
	}
	if resource.ID != saved.ResourceID || resource.Path != "/api/users/:id" || resource.Timeout != time.Second {
		t.Fatalf("generated resource: %+v", resource)
	}
	if method.ID != saved.MethodID || method.ResourcePath != resource.Path || method.HTTPVerb != "GET" || method.Timeout != time.Second {
		t.Fatalf("generated method: %+v", method)
	}
	if method.IntegrationRequest.Interface != "com.example.UserService" || method.IntegrationRequest.Method != "GetUser" {
		t.Fatalf("generated integration request: %+v", method.IntegrationRequest)
	}

	published, err := store.List(false)
	if err != nil {
		t.Fatalf("List published: %v", err)
	}
	if len(published) != 1 || published[0].Object.Metadata.Name != "user-get" {
		t.Fatalf("published bindings: %+v", published)
	}

	if err := store.DeleteDraft("user-get", saved.Revision); err != nil {
		t.Fatalf("DeleteDraft: %v", err)
	}
	status, err = store.PublishStatus()
	if err != nil {
		t.Fatalf("PublishStatus after delete: %v", err)
	}
	result, err = store.PublishAll(status.DraftRevision)
	if err != nil {
		t.Fatalf("PublishAll after delete: %v", err)
	}
	if result.PublishedCount != 0 || result.DeletedCount != 1 {
		t.Fatalf("delete publish result: %+v", result)
	}
	if fake.fakeHasKey(store.runtimeResourceKey(saved.ResourceID)) || fake.fakeHasKey(store.runtimeMethodKey(saved.ResourceID, saved.MethodID)) {
		t.Fatal("published runtime keys remain after atomic deletion")
	}
	if _, err := store.Get("user-get", false); err == nil {
		t.Fatal("deleted binding still exists in published namespace")
	}
}

func TestRouteBindingStorePublishOneKeepsOtherRoutesDraft(t *testing.T) {
	store := newTestRouteBindingStore(t)
	first, err := store.SaveDraft(testRouteBinding("user-get", "/api/users", "GET"), true, 0)
	if err != nil {
		t.Fatalf("first SaveDraft: %v", err)
	}
	second, err := store.SaveDraft(testRouteBinding("order-get", "/api/orders", "GET"), true, 0)
	if err != nil {
		t.Fatalf("second SaveDraft: %v", err)
	}

	result, err := store.Publish("user-get", first.Revision)
	if err != nil {
		t.Fatalf("Publish one route: %v", err)
	}
	if result.Name != "user-get" || result.PublishedCount != 1 || result.DeletedCount != 0 {
		t.Fatalf("single-route publish result: %+v", result)
	}
	if _, err := store.Get("order-get", false); err == nil {
		t.Fatal("publishing user-get unexpectedly published order-get")
	}
	if _, err := store.Get("order-get", true); err != nil {
		t.Fatalf("order-get draft disappeared: %v", err)
	}

	status, err := store.Status("user-get")
	if err != nil {
		t.Fatalf("user-get Status: %v", err)
	}
	if !status.DraftExists || !status.PublishedExists || status.Dirty {
		t.Fatalf("user-get status: %+v", status)
	}
	status, err = store.Status("order-get")
	if err != nil {
		t.Fatalf("order-get Status: %v", err)
	}
	if !status.DraftExists || status.PublishedExists || !status.Dirty {
		t.Fatalf("order-get status: %+v", status)
	}

	diff, err := store.Diff("order-get")
	if err != nil {
		t.Fatalf("order-get Diff: %v", err)
	}
	if diff.Published != nil || len(diff.Changes) == 0 {
		t.Fatalf("order-get diff: %+v", diff)
	}

	result, err = store.Publish("order-get", second.Revision)
	if err != nil {
		t.Fatalf("Publish second route: %v", err)
	}
	if result.Name != "order-get" || result.PublishedCount != 1 {
		t.Fatalf("second single-route publish result: %+v", result)
	}
	published, err := store.List(false)
	if err != nil {
		t.Fatalf("List published: %v", err)
	}
	if len(published) != 2 {
		t.Fatalf("published route count: got %d, want 2", len(published))
	}
}

func TestRouteBindingStorePublishRejectsConcurrentDraftChange(t *testing.T) {
	store := newTestRouteBindingStore(t)
	saved, err := store.SaveDraft(testRouteBinding("user-get", "/api/users/:id", "GET"), true, 0)
	if err != nil {
		t.Fatalf("initial SaveDraft: %v", err)
	}
	status, err := store.PublishStatus()
	if err != nil {
		t.Fatalf("PublishStatus: %v", err)
	}
	if _, err := store.SaveDraft(testRouteBinding("order-get", "/api/orders/:id", "GET"), true, 0); err != nil {
		t.Fatalf("concurrent SaveDraft: %v", err)
	}
	if _, err := store.PublishAll(status.DraftRevision); !strings.Contains(err.Error(), ErrRouteBindingPublishConflict.Error()) {
		t.Fatalf("stale publish error: want %q, got %v", ErrRouteBindingPublishConflict, err)
	}
	if fakeRouteBindingStoreKV(t, store).fakeHasKey(store.runtimeResourceKey(saved.ResourceID)) {
		t.Fatal("stale publish changed runtime keys")
	}
	if _, err := store.Get("order-get", false); err == nil {
		t.Fatal("stale publish wrote a published binding")
	}
}

func TestRouteBindingStoreRejectsDuplicatePublishedRoute(t *testing.T) {
	store := newTestRouteBindingStore(t)
	if _, err := store.SaveDraft(testRouteBinding("user-get-a", "/api/users", "GET"), true, 0); err != nil {
		t.Fatalf("first SaveDraft: %v", err)
	}
	if _, err := store.SaveDraft(testRouteBinding("user-get-b", "/api/users", "GET"), true, 0); err != nil {
		t.Fatalf("second SaveDraft: %v", err)
	}
	if _, err := store.PublishAll(0); err == nil || !strings.Contains(err.Error(), "conflicts") {
		t.Fatalf("duplicate publish error: %v", err)
	}
	fake := fakeRouteBindingStoreKV(t, store)
	if fake.fakeHasPrefix(store.bindingPrefix(false)) || fake.fakeHasPrefix(store.runtimeResourceKey(1)) {
		t.Fatal("duplicate publish wrote a partial published snapshot")
	}
}

func TestRouteBindingStoreSaveUsesPublishedRuntimeIdentity(t *testing.T) {
	store := newTestRouteBindingStore(t)
	saved, err := store.SaveDraft(testRouteBinding("user-get", "/api/users", "GET"), true, 0)
	if err != nil {
		t.Fatalf("SaveDraft: %v", err)
	}
	if _, err := store.PublishAll(0); err != nil {
		t.Fatalf("PublishAll: %v", err)
	}
	if err := store.DeleteDraft("user-get", saved.Revision); err != nil {
		t.Fatalf("DeleteDraft: %v", err)
	}
	updated := testRouteBinding("user-get", "/api/users/:id", "GET")
	updated.Spec["target"].(map[string]any)["method"] = "GetUserByID"
	restored, err := store.SaveDraft(updated, false, 0)
	if err != nil {
		t.Fatalf("SaveDraft from published identity: %v", err)
	}
	if restored.ResourceID != saved.ResourceID || restored.MethodID != saved.MethodID {
		t.Fatalf("runtime identity changed: old=%+v new=%+v", saved, restored)
	}
}

func newTestRouteBindingStore(t *testing.T) *RouteBindingStore {
	t.Helper()
	registry, err := schema.NewBuiltinRegistry()
	if err != nil {
		t.Fatalf("NewBuiltinRegistry: %v", err)
	}
	store, err := NewRouteBindingStore(newFakeRouteBindingKV(), context.Background(), "/pixiu/config/api", registry)
	if err != nil {
		t.Fatalf("NewRouteBindingStore: %v", err)
	}
	return store
}

func fakeRouteBindingStoreKV(t *testing.T, store *RouteBindingStore) *fakeRouteBindingKV {
	t.Helper()
	fake, ok := store.kv.(*fakeRouteBindingKV)
	if !ok {
		t.Fatalf("store is not using fake route binding kv: %T", store.kv)
	}
	return fake
}

func testRouteBinding(name, path, method string) schema.AdminObject {
	return schema.AdminObject{
		Kind: schema.KindAdminRouteBinding,
		Metadata: schema.ObjectMetadata{
			Name: name,
		},
		Spec: map[string]any{
			"entry": map[string]any{
				"path":   path,
				"method": method,
			},
			"target": map[string]any{
				"application": "UserProvider",
				"interface":   "com.example.UserService",
				"method":      "GetUser",
				"cluster":     "user-dubbo",
			},
			"params": []any{
				map[string]any{
					"from": "uri.id",
					"to":   0,
					"type": "java.lang.String",
				},
			},
		},
	}
}

type fakeRouteBindingValue struct {
	value     []byte
	createRev int64
	modifyRev int64
}

type fakeRouteBindingKV struct {
	mu       sync.Mutex
	revision int64
	values   map[string]fakeRouteBindingValue
}

func newFakeRouteBindingKV() *fakeRouteBindingKV {
	return &fakeRouteBindingKV{values: make(map[string]fakeRouteBindingValue)}
}

func (f *fakeRouteBindingKV) Get(_ context.Context, key string, options ...clientv3.OpOption) (*clientv3.GetResponse, error) {
	op := clientv3.OpGet(key, options...)
	f.mu.Lock()
	defer f.mu.Unlock()
	keys := f.matchingKeys(string(op.KeyBytes()), op.RangeBytes())
	response := &clientv3.GetResponse{Header: &pb.ResponseHeader{Revision: f.revision}}
	for _, matchedKey := range keys {
		value := f.values[matchedKey]
		response.Kvs = append(response.Kvs, &mvccpb.KeyValue{
			Key:            []byte(matchedKey),
			Value:          append([]byte(nil), value.value...),
			CreateRevision: value.createRev,
			ModRevision:    value.modifyRev,
		})
	}
	return response, nil
}

func (f *fakeRouteBindingKV) Txn(_ context.Context) clientv3.Txn {
	return &fakeRouteBindingTxn{store: f}
}

func (f *fakeRouteBindingKV) matchingKeys(key string, end []byte) []string {
	keys := make([]string, 0)
	for candidate := range f.values {
		if len(end) == 0 {
			if candidate == key {
				keys = append(keys, candidate)
			}
			continue
		}
		if candidate >= key && candidate < string(end) {
			keys = append(keys, candidate)
		}
	}
	sort.Strings(keys)
	return keys
}

func (f *fakeRouteBindingKV) fakeValue(key string) []byte {
	f.mu.Lock()
	defer f.mu.Unlock()
	return append([]byte(nil), f.values[key].value...)
}

func (f *fakeRouteBindingKV) fakeHasKey(key string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	_, exists := f.values[key]
	return exists
}

func (f *fakeRouteBindingKV) fakeHasPrefix(prefix string) bool {
	f.mu.Lock()
	defer f.mu.Unlock()
	for key := range f.values {
		if strings.HasPrefix(key, prefix) {
			return true
		}
	}
	return false
}

type fakeRouteBindingTxn struct {
	store       *fakeRouteBindingKV
	comparisons []clientv3.Cmp
	thenOps     []clientv3.Op
	elseOps     []clientv3.Op
}

func (t *fakeRouteBindingTxn) If(comparisons ...clientv3.Cmp) clientv3.Txn {
	t.comparisons = append(t.comparisons, comparisons...)
	return t
}

func (t *fakeRouteBindingTxn) Then(operations ...clientv3.Op) clientv3.Txn {
	t.thenOps = append(t.thenOps, operations...)
	return t
}

func (t *fakeRouteBindingTxn) Else(operations ...clientv3.Op) clientv3.Txn {
	t.elseOps = append(t.elseOps, operations...)
	return t
}

func (t *fakeRouteBindingTxn) Commit() (*clientv3.TxnResponse, error) {
	t.store.mu.Lock()
	defer t.store.mu.Unlock()
	succeeded := true
	for _, comparison := range t.comparisons {
		if !t.store.evaluateComparison(comparison) {
			succeeded = false
			break
		}
	}
	operations := t.thenOps
	if !succeeded {
		operations = t.elseOps
	}
	if len(operations) > 0 {
		t.store.revision++
		for _, operation := range operations {
			t.store.applyOperation(operation)
		}
	}
	return &clientv3.TxnResponse{
		Header:    &pb.ResponseHeader{Revision: t.store.revision},
		Succeeded: succeeded,
	}, nil
}

func (f *fakeRouteBindingKV) evaluateComparison(comparison clientv3.Cmp) bool {
	key := string(comparison.KeyBytes())
	value, exists := f.values[key]
	actual := int64(0)
	switch comparison.Target {
	case pb.Compare_CREATE:
		actual = value.createRev
	case pb.Compare_MOD:
		actual = value.modifyRev
	default:
		return false
	}
	var expected int64
	switch target := comparison.TargetUnion.(type) {
	case *pb.Compare_CreateRevision:
		expected = target.CreateRevision
	case *pb.Compare_ModRevision:
		expected = target.ModRevision
	default:
		return false
	}
	if comparison.Result != pb.Compare_EQUAL {
		return false
	}
	return actual == expected || (!exists && expected == 0)
}

func (f *fakeRouteBindingKV) applyOperation(operation clientv3.Op) {
	if operation.IsPut() {
		key := string(operation.KeyBytes())
		value := append([]byte(nil), operation.ValueBytes()...)
		entry := f.values[key]
		if entry.createRev == 0 {
			entry.createRev = f.revision
		}
		entry.modifyRev = f.revision
		entry.value = value
		f.values[key] = entry
		return
	}
	if !operation.IsDelete() {
		return
	}
	key := string(operation.KeyBytes())
	end := operation.RangeBytes()
	for candidate := range f.values {
		if len(end) == 0 {
			if candidate == key {
				delete(f.values, candidate)
			}
			continue
		}
		if candidate >= key && candidate < string(end) {
			delete(f.values, candidate)
		}
	}
}

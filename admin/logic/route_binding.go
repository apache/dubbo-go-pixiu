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
	"encoding/json"
	"errors"
	"fmt"
	"net/url"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	commonyaml "github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	legacyconfig "github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema"
)

const (
	routeBindingPathSegment        = "route-bindings"
	routeBindingRevisionSegment    = "route-bindings-revision"
	routeBindingRuntimeResourceDir = "resources"
	maxRouteBindingIDRetries       = 16
)

var (
	ErrRouteBindingNotFound        = errors.New("admin route binding not found")
	ErrRouteBindingAlreadyExists   = errors.New("admin route binding already exists")
	ErrRouteBindingConflict        = errors.New("admin route binding was modified concurrently")
	ErrRouteBindingPublishConflict = errors.New("admin route binding publish conflict")
)

// routeBindingKV is the small part of clientv3.KV needed by the route store.
// Keeping this boundary narrow makes the planning and transaction logic easy
// to exercise without coupling it to the gost wrapper.
type routeBindingKV interface {
	Get(context.Context, string, ...clientv3.OpOption) (*clientv3.GetResponse, error)
	Txn(context.Context) clientv3.Txn
}

// RouteBinding is the Admin-facing response shape. ResourceID and MethodID
// are stable runtime identities used by Pixiu's legacy etcd watcher; they are
// not part of the user-editable AdminRouteBinding object.
type RouteBinding struct {
	Object     schema.AdminObject `json:"object"`
	ResourceID int                `json:"resourceId"`
	MethodID   int                `json:"methodId"`
	Revision   int64              `json:"revision"`
}

// RouteBindingPublishResult describes one atomic publish transaction.
type RouteBindingPublishResult struct {
	Name              string `json:"name"`
	Revision          int64  `json:"revision"`
	DraftRevision     int64  `json:"draftRevision"`
	PublishedRevision int64  `json:"publishedRevision"`
	PublishedCount    int    `json:"publishedCount"`
	DeletedCount      int    `json:"deletedCount"`
}

// RouteBindingPublishStatus exposes one route's draft and published state.
// Revisions are the etcd key mod revisions for that route, not a global
// configuration revision.
type RouteBindingPublishStatus struct {
	Name              string `json:"name"`
	DraftRevision     int64  `json:"draftRevision"`
	PublishedRevision int64  `json:"publishedRevision"`
	DraftExists       bool   `json:"draftExists"`
	PublishedExists   bool   `json:"publishedExists"`
	Dirty             bool   `json:"dirty"`
}

type RouteBindingDiffChange struct {
	Path   string `json:"path"`
	Before any    `json:"before"`
	After  any    `json:"after"`
}

// RouteBindingDiff compares one route's draft object with its published
// object. A nil side represents a route that only exists in the other scope.
type RouteBindingDiff struct {
	Name      string                   `json:"name"`
	Draft     *RouteBinding            `json:"draft,omitempty"`
	Published *RouteBinding            `json:"published,omitempty"`
	Changes   []RouteBindingDiffChange `json:"changes"`
}

type routeBindingRecord struct {
	Object     schema.AdminObject `json:"object"`
	ResourceID int                `json:"resourceId"`
	MethodID   int                `json:"methodId"`
}

type routeBindingEntry struct {
	record   routeBindingRecord
	key      string
	revision int64
}

// RouteBindingStore owns the AdminRouteBinding persistence and publish
// boundary. Published runtime data is deliberately written in the existing
// Resource/Method layout so the current Pixiu watcher can consume it.
type RouteBindingStore struct {
	kv       routeBindingKV
	ctx      context.Context
	root     string
	registry *schema.Registry
}

// NewRouteBindingStore creates a store over an etcd v3 client.
func NewRouteBindingStore(kv routeBindingKV, ctx context.Context, root string, registry *schema.Registry) (*RouteBindingStore, error) {
	if kv == nil {
		return nil, errors.New("route binding etcd client is nil")
	}
	if ctx == nil {
		ctx = context.Background()
	}
	root = strings.TrimRight(strings.TrimSpace(root), "/")
	if root == "" {
		return nil, errors.New("route binding etcd root path is empty")
	}
	if registry == nil {
		var err error
		registry, err = schema.NewBuiltinRegistry()
		if err != nil {
			return nil, fmt.Errorf("create route binding schema registry: %w", err)
		}
	}
	return &RouteBindingStore{kv: kv, ctx: ctx, root: root, registry: registry}, nil
}

// NewAdminRouteBindingStore creates the production store from the Admin
// process' configured etcd client and root path.
func NewAdminRouteBindingStore() (*RouteBindingStore, error) {
	if adminconfig.Bootstrap == nil {
		return nil, errors.New("admin bootstrap is nil")
	}
	if adminconfig.Client == nil {
		return nil, errors.New("admin etcd client is nil")
	}
	rawClient := adminconfig.Client.GetRawClient()
	if rawClient == nil {
		return nil, errors.New("admin raw etcd client is nil")
	}
	return NewRouteBindingStore(rawClient, adminconfig.Client.GetCtx(), adminconfig.Bootstrap.GetPath(), nil)
}

// Registry returns the schema registry used by this store.
func (s *RouteBindingStore) Registry() *schema.Registry {
	return s.registry
}

// List returns high-level route bindings from either the draft or published
// namespace. Empty namespaces are returned as an empty slice.
func (s *RouteBindingStore) List(unpublished bool) ([]RouteBinding, error) {
	entries, err := s.listEntries(s.bindingPrefix(unpublished))
	if err != nil {
		return nil, err
	}
	result := make([]RouteBinding, 0, len(entries))
	for _, entry := range entries {
		result = append(result, routeBindingView(entry))
	}
	return result, nil
}

// Get returns one high-level route binding from either the draft or published
// namespace.
func (s *RouteBindingStore) Get(name string, unpublished bool) (RouteBinding, error) {
	name, err := normalizeRouteBindingName(name)
	if err != nil {
		return RouteBinding{}, err
	}
	entry, exists, _, err := s.getEntry(s.bindingPrefix(unpublished), name)
	if err != nil {
		return RouteBinding{}, err
	}
	if !exists {
		return RouteBinding{}, fmt.Errorf("%w: %s", ErrRouteBindingNotFound, name)
	}
	return routeBindingView(entry), nil
}

// Normalize validates an object and applies all registered schema defaults.
// It is useful to the controller's validate and preview endpoints and does
// not require an etcd write.
func (s *RouteBindingStore) Normalize(object schema.AdminObject) (schema.AdminObject, error) {
	normalized, err := s.registry.Normalize(object)
	if err != nil {
		return normalized, err
	}
	name, err := normalizeRouteBindingName(normalized.Metadata.Name)
	if err != nil {
		return normalized, err
	}
	normalized.Metadata.Name = name
	return normalized, nil
}

// Preview validates and compiles an Admin object without writing it.
func (s *RouteBindingStore) Preview(object schema.AdminObject) (schema.AdminObject, []byte, error) {
	normalized, err := s.Normalize(object)
	if err != nil {
		return normalized, nil, err
	}
	compiled, err := schema.CompileAdminRouteBinding(s.registry, normalized)
	if err != nil {
		return normalized, nil, err
	}
	preview, err := compiled.PreviewYAML()
	if err != nil {
		return normalized, nil, err
	}
	return normalized, preview, nil
}

// SaveDraft validates and atomically stores one AdminRouteBinding in the
// draft namespace. POST callers use create=true; PUT callers use create=false.
// expectedRevision is the draft key's mod revision and may be zero to allow
// an unconditional write.
func (s *RouteBindingStore) SaveDraft(object schema.AdminObject, create bool, expectedRevision int64) (RouteBinding, error) {
	if expectedRevision < 0 {
		return RouteBinding{}, errors.New("expected revision must not be negative")
	}
	if create && expectedRevision > 0 {
		return RouteBinding{}, errors.New("create route binding does not accept an expected revision")
	}
	normalized, err := s.Normalize(object)
	if err != nil {
		return RouteBinding{}, err
	}
	name := normalized.Metadata.Name
	draftPrefix := s.bindingPrefix(true)
	draftEntry, draftExists, draftRevision, err := s.getEntry(draftPrefix, name)
	if err != nil {
		return RouteBinding{}, err
	}
	if create && draftExists {
		return RouteBinding{}, fmt.Errorf("%w: %s", ErrRouteBindingAlreadyExists, name)
	}
	if expectedRevision > 0 && (!draftExists || draftRevision != expectedRevision) {
		return RouteBinding{}, fmt.Errorf("%w: %s", ErrRouteBindingConflict, name)
	}

	resourceID, methodID := 0, 0
	if draftExists {
		resourceID, methodID = draftEntry.record.ResourceID, draftEntry.record.MethodID
	} else {
		publishedEntry, publishedExists, _, getErr := s.getEntry(s.bindingPrefix(false), name)
		if getErr != nil {
			return RouteBinding{}, getErr
		}
		if publishedExists {
			resourceID, methodID = publishedEntry.record.ResourceID, publishedEntry.record.MethodID
		}
	}
	if resourceID <= 0 {
		resourceID, err = s.allocateRuntimeID()
		if err != nil {
			return RouteBinding{}, err
		}
	}
	if methodID <= 0 {
		// One AdminRouteBinding currently compiles to exactly one Method. Method
		// IDs are scoped below the resource, so sharing the stable route ID is
		// safe and avoids a second allocation transaction.
		methodID = resourceID
	}

	record := routeBindingRecord{
		Object:     normalized,
		ResourceID: resourceID,
		MethodID:   methodID,
	}
	value, err := json.Marshal(record)
	if err != nil {
		return RouteBinding{}, fmt.Errorf("encode route binding %q: %w", name, err)
	}

	comparisons := make([]clientv3.Cmp, 0, 1)
	if create {
		comparisons = append(comparisons, clientv3.Compare(clientv3.CreateRevision(s.bindingKey(draftPrefix, name)), "=", 0))
	} else if expectedRevision > 0 {
		comparisons = append(comparisons, clientv3.Compare(clientv3.ModRevision(s.bindingKey(draftPrefix, name)), "=", expectedRevision))
	}
	operations := []clientv3.Op{
		clientv3.OpPut(s.bindingKey(draftPrefix, name), string(value)),
		clientv3.OpPut(s.draftRevisionKey(), strconv.FormatInt(time.Now().UnixNano(), 10)),
	}
	transaction := s.kv.Txn(s.ctx)
	if len(comparisons) > 0 {
		transaction = transaction.If(comparisons...)
	}
	response, err := transaction.Then(operations...).Commit()
	if err != nil {
		return RouteBinding{}, fmt.Errorf("save route binding %q: %w", name, err)
	}
	if !response.Succeeded {
		if create {
			return RouteBinding{}, fmt.Errorf("%w: %s", ErrRouteBindingAlreadyExists, name)
		}
		return RouteBinding{}, fmt.Errorf("%w: %s", ErrRouteBindingConflict, name)
	}

	saved, exists, revision, err := s.getEntry(draftPrefix, name)
	if err != nil {
		return RouteBinding{}, err
	}
	if !exists {
		return RouteBinding{}, fmt.Errorf("route binding %q disappeared after save", name)
	}
	saved.revision = revision
	return routeBindingView(saved), nil
}

// DeleteDraft removes a draft. If only a published binding exists, the
// deletion is still recorded by advancing the draft revision; the next
// PublishAll transaction will remove that published route atomically.
func (s *RouteBindingStore) DeleteDraft(name string, expectedRevision int64) error {
	if expectedRevision < 0 {
		return errors.New("expected revision must not be negative")
	}
	name, err := normalizeRouteBindingName(name)
	if err != nil {
		return err
	}
	draftPrefix := s.bindingPrefix(true)
	_, exists, revision, err := s.getEntry(draftPrefix, name)
	if err != nil {
		return err
	}
	_, publishedExists, _, err := s.getEntry(s.bindingPrefix(false), name)
	if err != nil {
		return err
	}
	if !exists && !publishedExists {
		return fmt.Errorf("%w: %s", ErrRouteBindingNotFound, name)
	}
	if expectedRevision > 0 && (!exists || revision != expectedRevision) {
		return fmt.Errorf("%w: %s", ErrRouteBindingConflict, name)
	}
	draftMarkerRevision, draftMarkerExists, err := s.markerRevision(s.draftRevisionKey())
	if err != nil {
		return err
	}

	comparisons := make([]clientv3.Cmp, 0, 2)
	if exists {
		observedRevision := revision
		if expectedRevision > 0 {
			observedRevision = expectedRevision
		}
		comparisons = append(comparisons, clientv3.Compare(clientv3.ModRevision(s.bindingKey(draftPrefix, name)), "=", observedRevision))
	}
	if draftMarkerExists {
		comparisons = append(comparisons, clientv3.Compare(clientv3.ModRevision(s.draftRevisionKey()), "=", draftMarkerRevision))
	} else {
		comparisons = append(comparisons, clientv3.Compare(clientv3.CreateRevision(s.draftRevisionKey()), "=", 0))
	}
	operations := make([]clientv3.Op, 0, 2)
	if exists {
		operations = append(operations, clientv3.OpDelete(s.bindingKey(draftPrefix, name)))
	}
	operations = append(operations, clientv3.OpPut(s.draftRevisionKey(), strconv.FormatInt(time.Now().UnixNano(), 10)))
	transaction := s.kv.Txn(s.ctx)
	if len(comparisons) > 0 {
		transaction = transaction.If(comparisons...)
	}
	response, err := transaction.Then(operations...).Commit()
	if err != nil {
		return fmt.Errorf("delete route binding %q: %w", name, err)
	}
	if len(comparisons) > 0 && !response.Succeeded {
		return fmt.Errorf("%w: %s", ErrRouteBindingConflict, name)
	}
	return nil
}

// Publish publishes exactly one route binding. The generated legacy resource
// and method, the published high-level binding, and the publish marker are
// written or deleted in one etcd transaction. Other routes' draft and
// published values are left untouched.
func (s *RouteBindingStore) Publish(name string, expectedDraftRevision int64) (RouteBindingPublishResult, error) {
	if expectedDraftRevision < 0 {
		return RouteBindingPublishResult{}, errors.New("expected draft revision must not be negative")
	}
	name, err := normalizeRouteBindingName(name)
	if err != nil {
		return RouteBindingPublishResult{}, err
	}

	draftPrefix := s.bindingPrefix(true)
	publishedPrefix := s.bindingPrefix(false)
	draft, draftExists, draftRevision, err := s.getEntry(draftPrefix, name)
	if err != nil {
		return RouteBindingPublishResult{}, err
	}
	published, publishedExists, publishedRevision, err := s.getEntry(publishedPrefix, name)
	if err != nil {
		return RouteBindingPublishResult{}, err
	}
	if !draftExists && !publishedExists {
		return RouteBindingPublishResult{}, fmt.Errorf("%w: %s", ErrRouteBindingNotFound, name)
	}
	if expectedDraftRevision > 0 && (!draftExists || draftRevision != expectedDraftRevision) {
		return RouteBindingPublishResult{}, fmt.Errorf("%w: %s", ErrRouteBindingPublishConflict, name)
	}

	publishedEntries, err := s.listEntries(publishedPrefix)
	if err != nil {
		return RouteBindingPublishResult{}, err
	}

	var candidate *compiledRouteEntry
	if draftExists {
		entries := make([]routeBindingEntry, 0, len(publishedEntries)+1)
		for _, entry := range publishedEntries {
			if entry.record.Object.Metadata.Name != name {
				entries = append(entries, entry)
			}
		}
		entries = append(entries, draft)
		compiled, compileErr := s.compileSnapshot(entries)
		if compileErr != nil {
			return RouteBindingPublishResult{}, compileErr
		}
		for index := range compiled {
			if compiled[index].object.Metadata.Name == name {
				candidate = &compiled[index]
				break
			}
		}
		if candidate == nil {
			return RouteBindingPublishResult{}, fmt.Errorf("route binding %q was not compiled", name)
		}
	}

	comparisons := make([]clientv3.Cmp, 0, 2)
	if draftExists {
		// Always compare the value observed above, even when the caller omitted
		// expectedDraftRevision. This closes the read/modify/write race.
		comparisons = append(comparisons,
			clientv3.Compare(clientv3.ModRevision(s.bindingKey(draftPrefix, name)), "=", draftRevision))
	} else {
		comparisons = append(comparisons,
			clientv3.Compare(clientv3.ModRevision(s.bindingKey(publishedPrefix, name)), "=", publishedRevision))
	}
	publishedMarkerRevision, publishedMarkerExists, err := s.markerRevision(s.publishedRevisionKey())
	if err != nil {
		return RouteBindingPublishResult{}, err
	}
	if publishedMarkerExists {
		comparisons = append(comparisons,
			clientv3.Compare(clientv3.ModRevision(s.publishedRevisionKey()), "=", publishedMarkerRevision))
	} else {
		comparisons = append(comparisons,
			clientv3.Compare(clientv3.CreateRevision(s.publishedRevisionKey()), "=", 0))
	}

	operations := make([]clientv3.Op, 0, 4)
	publishedCount, deletedCount := 0, 0
	if candidate != nil {
		compiled := candidate.legacy
		if publishedExists && published.record.ResourceID != compiled.Resource.ID {
			operations = append(operations,
				clientv3.OpDelete(s.runtimeResourceKey(published.record.ResourceID), clientv3.WithPrefix()))
		} else if publishedExists && published.record.MethodID != compiled.Method.ID {
			operations = append(operations,
				clientv3.OpDelete(s.runtimeMethodKey(published.record.ResourceID, published.record.MethodID)))
		}
		resourceValue, methodValue, marshalErr := marshalCompiledRuntimeRoute(compiled)
		if marshalErr != nil {
			return RouteBindingPublishResult{}, fmt.Errorf("encode route binding %q for runtime: %w", name, marshalErr)
		}
		recordValue, marshalErr := json.Marshal(routeBindingRecord{
			Object:     candidate.object,
			ResourceID: compiled.Resource.ID,
			MethodID:   compiled.Method.ID,
		})
		if marshalErr != nil {
			return RouteBindingPublishResult{}, fmt.Errorf("encode published route binding %q: %w", name, marshalErr)
		}
		operations = append(operations,
			clientv3.OpPut(s.runtimeResourceKey(compiled.Resource.ID), string(resourceValue)),
			clientv3.OpPut(s.runtimeMethodKey(compiled.Resource.ID, compiled.Method.ID), string(methodValue)),
			clientv3.OpPut(s.bindingKey(publishedPrefix, name), string(recordValue)))
		publishedCount = 1
	} else {
		operations = append(operations,
			clientv3.OpDelete(s.runtimeResourceKey(published.record.ResourceID), clientv3.WithPrefix()),
			clientv3.OpDelete(s.bindingKey(publishedPrefix, name)))
		deletedCount = 1
	}
	operations = append(operations,
		clientv3.OpPut(s.publishedRevisionKey(), strconv.FormatInt(time.Now().UnixNano(), 10)))

	response, err := s.kv.Txn(s.ctx).If(comparisons...).Then(operations...).Commit()
	if err != nil {
		return RouteBindingPublishResult{}, fmt.Errorf("publish route binding %q: %w", name, err)
	}
	if !response.Succeeded {
		return RouteBindingPublishResult{}, fmt.Errorf("%w: %s", ErrRouteBindingPublishConflict, name)
	}

	publishedRevision = 0
	if candidate != nil {
		_, exists, revision, getErr := s.getEntry(publishedPrefix, name)
		if getErr != nil {
			return RouteBindingPublishResult{}, getErr
		}
		if !exists {
			return RouteBindingPublishResult{}, fmt.Errorf("route binding %q disappeared after publish", name)
		}
		publishedRevision = revision
	}
	return RouteBindingPublishResult{
		Name:              name,
		Revision:          response.Header.Revision,
		DraftRevision:     draftRevision,
		PublishedRevision: publishedRevision,
		PublishedCount:    publishedCount,
		DeletedCount:      deletedCount,
	}, nil
}

// Status returns draft and published revisions for one route. The revisions
// are independent of other route bindings in the same config namespace.
func (s *RouteBindingStore) Status(name string) (RouteBindingPublishStatus, error) {
	name, err := normalizeRouteBindingName(name)
	if err != nil {
		return RouteBindingPublishStatus{}, err
	}
	draft, draftExists, draftRevision, err := s.getEntry(s.bindingPrefix(true), name)
	if err != nil {
		return RouteBindingPublishStatus{}, err
	}
	published, publishedExists, publishedRevision, err := s.getEntry(s.bindingPrefix(false), name)
	if err != nil {
		return RouteBindingPublishStatus{}, err
	}
	if !draftExists && !publishedExists {
		return RouteBindingPublishStatus{}, fmt.Errorf("%w: %s", ErrRouteBindingNotFound, name)
	}
	dirty := draftExists != publishedExists
	if draftExists && publishedExists {
		dirty = !reflect.DeepEqual(draft.record.Object, published.record.Object)
	}
	return RouteBindingPublishStatus{
		Name:              name,
		DraftRevision:     draftRevision,
		PublishedRevision: publishedRevision,
		DraftExists:       draftExists,
		PublishedExists:   publishedExists,
		Dirty:             dirty,
	}, nil
}

// Diff compares one route's draft and published objects. It does not read or
// mutate any other route's draft state.
func (s *RouteBindingStore) Diff(name string) (RouteBindingDiff, error) {
	name, err := normalizeRouteBindingName(name)
	if err != nil {
		return RouteBindingDiff{}, err
	}
	draft, draftExists, _, err := s.getEntry(s.bindingPrefix(true), name)
	if err != nil {
		return RouteBindingDiff{}, err
	}
	published, publishedExists, _, err := s.getEntry(s.bindingPrefix(false), name)
	if err != nil {
		return RouteBindingDiff{}, err
	}
	if !draftExists && !publishedExists {
		return RouteBindingDiff{}, fmt.Errorf("%w: %s", ErrRouteBindingNotFound, name)
	}

	var draftView, publishedView *RouteBinding
	if draftExists {
		view := routeBindingView(draft)
		draftView = &view
	}
	if publishedExists {
		view := routeBindingView(published)
		publishedView = &view
	}
	changes, err := diffRouteBindingObjects(
		func() schema.AdminObject {
			if publishedExists {
				return published.record.Object
			}
			return schema.AdminObject{}
		}(),
		func() schema.AdminObject {
			if draftExists {
				return draft.record.Object
			}
			return schema.AdminObject{}
		}(),
		publishedExists,
		draftExists,
	)
	if err != nil {
		return RouteBindingDiff{}, err
	}
	return RouteBindingDiff{
		Name:      name,
		Draft:     draftView,
		Published: publishedView,
		Changes:   changes,
	}, nil
}

func diffRouteBindingObjects(before, after schema.AdminObject, beforeExists, afterExists bool) ([]RouteBindingDiffChange, error) {
	beforeValues := make(map[string]any)
	afterValues := make(map[string]any)
	if beforeExists {
		if err := flattenRouteBindingValue("", before, beforeValues); err != nil {
			return nil, err
		}
	}
	if afterExists {
		if err := flattenRouteBindingValue("", after, afterValues); err != nil {
			return nil, err
		}
	}
	keys := make(map[string]struct{}, len(beforeValues)+len(afterValues))
	for key := range beforeValues {
		keys[key] = struct{}{}
	}
	for key := range afterValues {
		keys[key] = struct{}{}
	}
	paths := make([]string, 0, len(keys))
	for key := range keys {
		paths = append(paths, key)
	}
	sort.Strings(paths)
	changes := make([]RouteBindingDiffChange, 0, len(paths))
	for _, path := range paths {
		beforeValue, beforeOK := beforeValues[path]
		afterValue, afterOK := afterValues[path]
		if beforeOK && afterOK && reflect.DeepEqual(beforeValue, afterValue) {
			continue
		}
		changes = append(changes, RouteBindingDiffChange{
			Path:   path,
			Before: valueOrNil(beforeValue, beforeOK),
			After:  valueOrNil(afterValue, afterOK),
		})
	}
	return changes, nil
}

func flattenRouteBindingValue(path string, value any, result map[string]any) error {
	switch typed := value.(type) {
	case schema.AdminObject:
		return flattenRouteBindingValue(path, map[string]any{
			"kind":     typed.Kind,
			"metadata": typed.Metadata,
			"spec":     typed.Spec,
		}, result)
	case schema.ObjectMetadata:
		return flattenRouteBindingValue(path, map[string]any{"name": typed.Name}, result)
	case map[string]any:
		if len(typed) == 0 {
			result[path] = typed
			return nil
		}
		keys := make([]string, 0, len(typed))
		for key := range typed {
			keys = append(keys, key)
		}
		sort.Strings(keys)
		for _, key := range keys {
			childPath := key
			if path != "" {
				childPath = path + "." + key
			}
			if err := flattenRouteBindingValue(childPath, typed[key], result); err != nil {
				return err
			}
		}
		return nil
	case []any:
		if len(typed) == 0 {
			result[path] = typed
			return nil
		}
		for index, item := range typed {
			childPath := fmt.Sprintf("%s[%d]", path, index)
			if err := flattenRouteBindingValue(childPath, item, result); err != nil {
				return err
			}
		}
		return nil
	default:
		result[path] = typed
		return nil
	}
}

func valueOrNil(value any, exists bool) any {
	if !exists {
		return nil
	}
	return value
}

// PublishAll compiles the complete draft snapshot and publishes it with one
// etcd transaction. The transaction updates the high-level published records,
// generated legacy resources/methods, stale-route deletions, and the published
// revision marker together. A draft marker compare prevents a concurrent save
// or delete from being silently omitted from the snapshot.
func (s *RouteBindingStore) PublishAll(expectedDraftRevision int64) (RouteBindingPublishResult, error) {
	if expectedDraftRevision < 0 {
		return RouteBindingPublishResult{}, errors.New("expected draft revision must not be negative")
	}
	draftEntries, err := s.listEntries(s.bindingPrefix(true))
	if err != nil {
		return RouteBindingPublishResult{}, err
	}
	publishedEntries, err := s.listEntries(s.bindingPrefix(false))
	if err != nil {
		return RouteBindingPublishResult{}, err
	}

	draftRevision, draftMarkerExists, err := s.markerRevision(s.draftRevisionKey())
	if err != nil {
		return RouteBindingPublishResult{}, err
	}
	if expectedDraftRevision > 0 && (!draftMarkerExists || draftRevision != expectedDraftRevision) {
		return RouteBindingPublishResult{}, fmt.Errorf("%w: draft revision changed", ErrRouteBindingPublishConflict)
	}

	compiled, err := s.compileSnapshot(draftEntries)
	if err != nil {
		return RouteBindingPublishResult{}, err
	}
	operations, publishedCount, deletedCount, err := s.buildPublishOperations(compiled, publishedEntries)
	if err != nil {
		return RouteBindingPublishResult{}, err
	}

	comparisons := make([]clientv3.Cmp, 0, 1)
	if draftMarkerExists {
		comparisons = append(comparisons, clientv3.Compare(clientv3.ModRevision(s.draftRevisionKey()), "=", draftRevision))
	} else {
		comparisons = append(comparisons, clientv3.Compare(clientv3.CreateRevision(s.draftRevisionKey()), "=", 0))
	}
	operations = append(operations,
		clientv3.OpPut(s.publishedRevisionKey(), strconv.FormatInt(draftRevision, 10)),
	)
	response, err := s.kv.Txn(s.ctx).If(comparisons...).Then(operations...).Commit()
	if err != nil {
		return RouteBindingPublishResult{}, fmt.Errorf("publish route bindings: %w", err)
	}
	if !response.Succeeded {
		return RouteBindingPublishResult{}, fmt.Errorf("%w: draft revision changed", ErrRouteBindingPublishConflict)
	}
	return RouteBindingPublishResult{
		Revision:       response.Header.Revision,
		DraftRevision:  draftRevision,
		PublishedCount: publishedCount,
		DeletedCount:   deletedCount,
	}, nil
}

// PublishStatus returns the current draft and published marker revisions.
func (s *RouteBindingStore) PublishStatus() (RouteBindingPublishStatus, error) {
	draftRevision, _, err := s.markerRevision(s.draftRevisionKey())
	if err != nil {
		return RouteBindingPublishStatus{}, err
	}
	publishedRevision, _, err := s.markerRevision(s.publishedRevisionKey())
	if err != nil {
		return RouteBindingPublishStatus{}, err
	}
	return RouteBindingPublishStatus{
		DraftRevision:     draftRevision,
		PublishedRevision: publishedRevision,
	}, nil
}

type compiledRouteEntry struct {
	entry  routeBindingEntry
	object schema.AdminObject
	legacy schema.CompiledRoute
}

func (s *RouteBindingStore) compileSnapshot(entries []routeBindingEntry) ([]compiledRouteEntry, error) {
	result := make([]compiledRouteEntry, 0, len(entries))
	paths := make(map[string]string, len(entries))
	resourceIDs := make(map[int]string, len(entries))
	for _, entry := range entries {
		if entry.record.ResourceID <= 0 || entry.record.MethodID <= 0 {
			return nil, fmt.Errorf("route binding %q has invalid runtime identity", entry.record.Object.Metadata.Name)
		}
		normalized, err := s.Normalize(entry.record.Object)
		if err != nil {
			return nil, fmt.Errorf("route binding %q: %w", entry.record.Object.Metadata.Name, err)
		}
		legacy, err := schema.CompileAdminRouteBinding(s.registry, normalized)
		if err != nil {
			return nil, fmt.Errorf("route binding %q: %w", entry.record.Object.Metadata.Name, err)
		}
		legacy.Resource.ID = entry.record.ResourceID
		legacy.Method.ID = entry.record.MethodID
		routeKey := legacy.Resource.Path + "\x00" + legacy.Method.HTTPVerb
		if previous, exists := paths[routeKey]; exists {
			return nil, fmt.Errorf("route binding %q conflicts with route binding %q for %s %s", entry.record.Object.Metadata.Name, previous, legacy.Method.HTTPVerb, legacy.Resource.Path)
		}
		if previous, exists := resourceIDs[legacy.Resource.ID]; exists && previous != entry.record.Object.Metadata.Name {
			return nil, fmt.Errorf("route binding %q reuses runtime resource id %d from %q", entry.record.Object.Metadata.Name, legacy.Resource.ID, previous)
		}
		paths[routeKey] = entry.record.Object.Metadata.Name
		resourceIDs[legacy.Resource.ID] = entry.record.Object.Metadata.Name
		result = append(result, compiledRouteEntry{entry: entry, object: normalized, legacy: legacy})
	}
	return result, nil
}

func (s *RouteBindingStore) buildPublishOperations(compiled []compiledRouteEntry, published []routeBindingEntry) ([]clientv3.Op, int, int, error) {
	oldByName := make(map[string]routeBindingEntry, len(published))
	oldIDs := make(map[int]string, len(published))
	for _, entry := range published {
		if entry.record.ResourceID <= 0 || entry.record.MethodID <= 0 {
			return nil, 0, 0, fmt.Errorf("published route binding %q has invalid runtime identity", entry.record.Object.Metadata.Name)
		}
		name := entry.record.Object.Metadata.Name
		oldByName[name] = entry
		if previous, exists := oldIDs[entry.record.ResourceID]; exists && previous != name {
			return nil, 0, 0, fmt.Errorf("published route bindings %q and %q share runtime resource id %d", previous, name, entry.record.ResourceID)
		}
		oldIDs[entry.record.ResourceID] = name
	}

	newByName := make(map[string]compiledRouteEntry, len(compiled))
	newIDs := make(map[int]string, len(compiled))
	for _, entry := range compiled {
		name := entry.object.Metadata.Name
		newByName[name] = entry
		newIDs[entry.legacy.Resource.ID] = name
	}

	operations := make([]clientv3.Op, 0, len(compiled)*3+len(published)*2)
	deletedCount := 0
	for name, old := range oldByName {
		_, exists := newByName[name]
		if exists {
			continue
		}
		if owner, reused := newIDs[old.record.ResourceID]; reused {
			return nil, 0, 0, fmt.Errorf("cannot delete published route binding %q because runtime resource id %d is reused by %q", name, old.record.ResourceID, owner)
		}
		operations = append(operations,
			clientv3.OpDelete(s.runtimeResourceKey(old.record.ResourceID), clientv3.WithPrefix()),
			clientv3.OpDelete(s.bindingKey(s.bindingPrefix(false), name)),
		)
		deletedCount++
	}

	for _, entry := range compiled {
		name := entry.object.Metadata.Name
		old, existed := oldByName[name]
		if existed && old.record.ResourceID != entry.legacy.Resource.ID {
			// The old identity was deleted above only when the name is present in
			// the new snapshot, so delete it here before writing the new identity.
			operations = append(operations, clientv3.OpDelete(s.runtimeResourceKey(old.record.ResourceID), clientv3.WithPrefix()))
		} else if existed && old.record.MethodID != entry.legacy.Method.ID {
			operations = append(operations, clientv3.OpDelete(s.runtimeMethodKey(old.record.ResourceID, old.record.MethodID)))
		}
		resourceValue, methodValue, err := marshalCompiledRuntimeRoute(entry.legacy)
		if err != nil {
			return nil, 0, 0, fmt.Errorf("encode route binding %q for runtime: %w", name, err)
		}
		recordValue, err := json.Marshal(routeBindingRecord{
			Object:     entry.object,
			ResourceID: entry.legacy.Resource.ID,
			MethodID:   entry.legacy.Method.ID,
		})
		if err != nil {
			return nil, 0, 0, fmt.Errorf("encode published route binding %q: %w", name, err)
		}
		operations = append(operations,
			clientv3.OpPut(s.runtimeResourceKey(entry.legacy.Resource.ID), string(resourceValue)),
			clientv3.OpPut(s.runtimeMethodKey(entry.legacy.Resource.ID, entry.legacy.Method.ID), string(methodValue)),
			clientv3.OpPut(s.bindingKey(s.bindingPrefix(false), name), string(recordValue)),
		)
	}
	return operations, len(compiled), deletedCount, nil
}

func marshalCompiledRuntimeRoute(compiled schema.CompiledRoute) ([]byte, []byte, error) {
	resource := routeBindingResourceYAML{
		ID:      compiled.Resource.ID,
		Type:    compiled.Resource.Type,
		Path:    compiled.Resource.Path,
		Timeout: compiled.Resource.Timeout.String(),
	}
	method := routeBindingMethodYAML{
		ID:                 compiled.Method.ID,
		ResourcePath:       compiled.Method.ResourcePath,
		Enable:             compiled.Method.Enable,
		Timeout:            compiled.Method.Timeout.String(),
		Mock:               compiled.Method.Mock,
		Filters:            compiled.Method.Filters,
		HTTPVerb:           compiled.Method.HTTPVerb,
		InboundRequest:     compiled.Method.InboundRequest,
		IntegrationRequest: routeBindingIntegrationRequestYAMLFrom(compiled.Method.IntegrationRequest),
	}
	resourceValue, err := commonyaml.MarshalYML(resource)
	if err != nil {
		return nil, nil, err
	}
	methodValue, err := commonyaml.MarshalYML(method)
	if err != nil {
		return nil, nil, err
	}
	return resourceValue, methodValue, nil
}

type routeBindingResourceYAML struct {
	ID      int    `yaml:"id,omitempty"`
	Type    string `yaml:"type"`
	Path    string `yaml:"path"`
	Timeout string `yaml:"timeout"`
}

type routeBindingMethodYAML struct {
	ID                 int                                `yaml:"id,omitempty"`
	ResourcePath       string                             `yaml:"resourcePath"`
	Enable             bool                               `yaml:"enable"`
	Timeout            string                             `yaml:"timeout"`
	Mock               bool                               `yaml:"mock"`
	Filters            []legacyconfig.Filter              `yaml:"filters,omitempty"`
	HTTPVerb           string                             `yaml:"httpVerb"`
	InboundRequest     legacyconfig.InboundRequest        `yaml:"inboundRequest"`
	IntegrationRequest routeBindingIntegrationRequestYAML `yaml:"integrationRequest"`
}

type routeBindingIntegrationRequestYAML struct {
	RequestType     string                      `yaml:"requestType"`
	ClusterName     string                      `yaml:"clusterName"`
	ApplicationName string                      `yaml:"applicationName"`
	Protocol        string                      `yaml:"protocol,omitempty"`
	Group           string                      `yaml:"group,omitempty"`
	Version         string                      `yaml:"version,omitempty"`
	Interface       string                      `yaml:"interface"`
	Method          string                      `yaml:"method"`
	ParameterTypes  []string                    `yaml:"parameterTypes,omitempty"`
	Serialization   string                      `yaml:"serialization,omitempty"`
	Retries         string                      `yaml:"retries,omitempty"`
	MappingParams   []legacyconfig.MappingParam `yaml:"mappingParams,omitempty"`
}

func routeBindingIntegrationRequestYAMLFrom(request legacyconfig.IntegrationRequest) routeBindingIntegrationRequestYAML {
	return routeBindingIntegrationRequestYAML{
		RequestType:     request.RequestType,
		ClusterName:     request.ClusterName,
		ApplicationName: request.ApplicationName,
		Protocol:        request.Protocol,
		Group:           request.Group,
		Version:         request.Version,
		Interface:       request.Interface,
		Method:          request.Method,
		ParameterTypes:  append([]string(nil), request.ParameterTypes...),
		Serialization:   request.Serialization,
		Retries:         request.Retries,
		MappingParams:   append([]legacyconfig.MappingParam(nil), request.MappingParams...),
	}
}

func (s *RouteBindingStore) listEntries(prefix string) ([]routeBindingEntry, error) {
	response, err := s.kv.Get(s.ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
	if err != nil {
		return nil, fmt.Errorf("list route bindings under %q: %w", prefix, err)
	}
	entries := make([]routeBindingEntry, 0, len(response.Kvs))
	for _, kv := range response.Kvs {
		name, err := routeBindingNameFromKey(prefix, string(kv.Key))
		if err != nil {
			return nil, err
		}
		var record routeBindingRecord
		if err := json.Unmarshal(kv.Value, &record); err != nil {
			return nil, fmt.Errorf("decode route binding %q: %w", name, err)
		}
		if record.Object.Metadata.Name != name {
			return nil, fmt.Errorf("route binding key %q does not match metadata.name %q", name, record.Object.Metadata.Name)
		}
		entries = append(entries, routeBindingEntry{record: record, key: string(kv.Key), revision: kv.ModRevision})
	}
	return entries, nil
}

func (s *RouteBindingStore) getEntry(prefix, name string) (routeBindingEntry, bool, int64, error) {
	key := s.bindingKey(prefix, name)
	response, err := s.kv.Get(s.ctx, key)
	if err != nil {
		return routeBindingEntry{}, false, 0, fmt.Errorf("get route binding %q: %w", name, err)
	}
	if len(response.Kvs) == 0 {
		return routeBindingEntry{}, false, 0, nil
	}
	var record routeBindingRecord
	if err := json.Unmarshal(response.Kvs[0].Value, &record); err != nil {
		return routeBindingEntry{}, false, 0, fmt.Errorf("decode route binding %q: %w", name, err)
	}
	if record.Object.Metadata.Name != name {
		return routeBindingEntry{}, false, 0, fmt.Errorf("route binding key %q does not match metadata.name %q", name, record.Object.Metadata.Name)
	}
	entry := routeBindingEntry{record: record, key: key, revision: response.Kvs[0].ModRevision}
	return entry, true, entry.revision, nil
}

func (s *RouteBindingStore) markerRevision(key string) (int64, bool, error) {
	response, err := s.kv.Get(s.ctx, key)
	if err != nil {
		return 0, false, fmt.Errorf("get route binding revision marker %q: %w", key, err)
	}
	if len(response.Kvs) == 0 {
		return 0, false, nil
	}
	return response.Kvs[0].ModRevision, true, nil
}

func (s *RouteBindingStore) allocateRuntimeID() (int, error) {
	counterKey := s.root + "/" + ResourceID
	resourcePrefix := s.root + "/" + routeBindingRuntimeResourceDir + "/"
	for attempt := 0; attempt < maxRouteBindingIDRetries; attempt++ {
		counterResponse, err := s.kv.Get(s.ctx, counterKey)
		if err != nil {
			return 0, fmt.Errorf("read route binding resource id counter: %w", err)
		}
		current, counterRevision := 0, int64(0)
		if len(counterResponse.Kvs) > 0 {
			current, err = strconv.Atoi(string(counterResponse.Kvs[0].Value))
			if err != nil || current < 0 {
				return 0, fmt.Errorf("invalid resource id counter %q", string(counterResponse.Kvs[0].Value))
			}
			counterRevision = counterResponse.Kvs[0].ModRevision
		}

		maxExisting := current
		resourceResponse, err := s.kv.Get(s.ctx, resourcePrefix, clientv3.WithPrefix())
		if err != nil {
			return 0, fmt.Errorf("scan route binding resource ids: %w", err)
		}
		for _, kv := range resourceResponse.Kvs {
			relative := strings.TrimPrefix(string(kv.Key), resourcePrefix)
			if strings.Contains(relative, "/") {
				continue
			}
			id, parseErr := strconv.Atoi(relative)
			if parseErr == nil && id > maxExisting {
				maxExisting = id
			}
		}

		next := maxExisting + 1
		comparison := clientv3.Compare(clientv3.CreateRevision(counterKey), "=", 0)
		if counterRevision > 0 {
			comparison = clientv3.Compare(clientv3.ModRevision(counterKey), "=", counterRevision)
		}
		response, err := s.kv.Txn(s.ctx).
			If(comparison).
			Then(clientv3.OpPut(counterKey, strconv.Itoa(next))).
			Commit()
		if err != nil {
			return 0, fmt.Errorf("allocate route binding resource id: %w", err)
		}
		if response.Succeeded {
			return next, nil
		}
	}
	return 0, fmt.Errorf("%w: resource id allocation retries exhausted", ErrRouteBindingConflict)
}

func routeBindingView(entry routeBindingEntry) RouteBinding {
	return RouteBinding{
		Object:     entry.record.Object.Clone(),
		ResourceID: entry.record.ResourceID,
		MethodID:   entry.record.MethodID,
		Revision:   entry.revision,
	}
}

func normalizeRouteBindingName(name string) (string, error) {
	name = strings.TrimSpace(name)
	if name == "" {
		return "", errors.New("route binding metadata.name is required")
	}
	if name == "." || name == ".." || strings.ContainsAny(name, `/\\`) {
		return "", fmt.Errorf("route binding metadata.name %q is not a safe key", name)
	}
	if len(name) > 128 {
		return "", errors.New("route binding metadata.name must be at most 128 characters")
	}
	return name, nil
}

func routeBindingNameFromKey(prefix, key string) (string, error) {
	if !strings.HasPrefix(key, prefix) {
		return "", fmt.Errorf("route binding key %q is outside prefix %q", key, prefix)
	}
	encodedName := strings.TrimPrefix(key, prefix)
	if encodedName == "" || strings.Contains(encodedName, "/") {
		return "", fmt.Errorf("route binding key %q is not a direct binding key", key)
	}
	name, err := url.PathUnescape(encodedName)
	if err != nil {
		return "", fmt.Errorf("decode route binding key %q: %w", key, err)
	}
	return normalizeRouteBindingName(name)
}

func (s *RouteBindingStore) bindingPrefix(unpublished bool) string {
	if unpublished {
		return s.root + "/" + Unpublished + "/" + routeBindingPathSegment + "/"
	}
	return s.root + "/" + routeBindingPathSegment + "/"
}

func (s *RouteBindingStore) bindingKey(prefix, name string) string {
	return prefix + url.PathEscape(name)
}

func (s *RouteBindingStore) draftRevisionKey() string {
	return s.root + "/" + Unpublished + "/" + routeBindingRevisionSegment
}

func (s *RouteBindingStore) publishedRevisionKey() string {
	return s.root + "/" + routeBindingRevisionSegment
}

func (s *RouteBindingStore) runtimeResourceKey(resourceID int) string {
	return s.root + "/" + routeBindingRuntimeResourceDir + "/" + strconv.Itoa(resourceID)
}

func (s *RouteBindingStore) runtimeMethodKey(resourceID, methodID int) string {
	return s.runtimeResourceKey(resourceID) + "/" + Method + "/" + strconv.Itoa(methodID)
}

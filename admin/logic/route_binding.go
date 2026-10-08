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
	ErrRouteBindingNameImmutable   = errors.New("admin route binding metadata.name is immutable")
	ErrRouteBindingConflict        = errors.New("admin route binding was modified concurrently")
	ErrRouteBindingPublishConflict = errors.New("admin route binding publish conflict")
	ErrRouteBindingRuntimeConflict = errors.New("admin route binding conflicts with an existing runtime route")
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
	Object        schema.AdminObject         `json:"object"`
	ResourceID    int                        `json:"resourceId"`
	MethodID      int                        `json:"methodId"`
	Revision      int64                      `json:"revision"`
	PublishStatus *RouteBindingPublishStatus `json:"publishStatus,omitempty"`
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
	root     string
	registry *schema.Registry
}

// NewRouteBindingStore creates a store over an etcd v3 client.
func NewRouteBindingStore(kv routeBindingKV, root string, registry *schema.Registry) (*RouteBindingStore, error) {
	if kv == nil {
		return nil, errors.New("route binding etcd client is nil")
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
	return &RouteBindingStore{kv: kv, root: root, registry: registry}, nil
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
	registry, err := schema.DefaultRegistry()
	if err != nil {
		return nil, fmt.Errorf("get Admin route binding schema registry: %w", err)
	}
	return NewRouteBindingStore(rawClient, adminconfig.Bootstrap.GetPath(), registry)
}

func routeBindingContext(ctx context.Context) context.Context {
	if ctx == nil {
		return context.Background()
	}
	return ctx
}

// Registry returns the schema registry used by this store.
func (s *RouteBindingStore) Registry() *schema.Registry {
	return s.registry
}

// List returns high-level route bindings from either the draft or published
// namespace. Empty namespaces are returned as an empty slice.
func (s *RouteBindingStore) List(ctx context.Context, unpublished bool) ([]RouteBinding, error) {
	ctx = routeBindingContext(ctx)
	entries, err := s.listEntries(ctx, s.bindingPrefix(unpublished))
	if err != nil {
		return nil, err
	}
	var statuses map[string]RouteBindingPublishStatus
	if unpublished {
		publishedEntries, err := s.listEntries(ctx, s.bindingPrefix(false))
		if err != nil {
			return nil, err
		}
		statuses = routeBindingStatuses(entries, publishedEntries)
	}
	result := make([]RouteBinding, 0, len(entries))
	for _, entry := range entries {
		view := routeBindingView(entry)
		if status, exists := statuses[entry.record.Object.Metadata.Name]; exists {
			view.PublishStatus = &status
		}
		result = append(result, view)
	}
	return result, nil
}

func routeBindingStatuses(draftEntries, publishedEntries []routeBindingEntry) map[string]RouteBindingPublishStatus {
	publishedByName := make(map[string]routeBindingEntry, len(publishedEntries))
	for _, entry := range publishedEntries {
		publishedByName[entry.record.Object.Metadata.Name] = entry
	}

	statuses := make(map[string]RouteBindingPublishStatus, len(draftEntries))
	for _, draft := range draftEntries {
		name := draft.record.Object.Metadata.Name
		published, publishedExists := publishedByName[name]
		status := RouteBindingPublishStatus{
			Name:              name,
			DraftRevision:     draft.revision,
			DraftExists:       true,
			PublishedExists:   publishedExists,
			PublishedRevision: published.revision,
			Dirty:             !publishedExists || !reflect.DeepEqual(draft.record.Object, published.record.Object),
		}
		statuses[name] = status
	}
	return statuses
}

// Get returns one high-level route binding from either the draft or published
// namespace.
func (s *RouteBindingStore) Get(ctx context.Context, name string, unpublished bool) (RouteBinding, error) {
	ctx = routeBindingContext(ctx)
	name, err := normalizeRouteBindingName(name)
	if err != nil {
		return RouteBinding{}, err
	}
	entry, exists, _, err := s.getEntry(ctx, s.bindingPrefix(unpublished), name)
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
func (s *RouteBindingStore) SaveDraft(ctx context.Context, object schema.AdminObject, create bool, expectedRevision int64) (RouteBinding, error) {
	ctx = routeBindingContext(ctx)
	if err := validateSaveDraftRequest(create, expectedRevision); err != nil {
		return RouteBinding{}, err
	}
	normalized, err := s.Normalize(object)
	if err != nil {
		return RouteBinding{}, err
	}
	name := normalized.Metadata.Name
	draftPrefix := s.bindingPrefix(true)
	draftEntry, draftExists, draftRevision, err := s.getEntry(ctx, draftPrefix, name)
	if err != nil {
		return RouteBinding{}, err
	}
	if err := validateDraftSaveState(name, create, expectedRevision, draftExists, draftRevision); err != nil {
		return RouteBinding{}, err
	}

	resourceID, methodID, err := s.resolveRuntimeIdentity(ctx, name, draftEntry, draftExists)
	if err != nil {
		return RouteBinding{}, err
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
	if err := s.commitDraft(ctx, name, draftPrefix, value, create, expectedRevision); err != nil {
		return RouteBinding{}, err
	}

	saved, exists, revision, err := s.getEntry(ctx, draftPrefix, name)
	if err != nil {
		return RouteBinding{}, err
	}
	if !exists {
		return RouteBinding{}, fmt.Errorf("route binding %q disappeared after save", name)
	}
	saved.revision = revision
	return routeBindingView(saved), nil
}

// UpdateDraft updates one existing route binding draft without allowing its
// metadata.name to change. The original name is the stable identity supplied
// by the caller, while the object name is user-editable in the YAML editor.
func (s *RouteBindingStore) UpdateDraft(ctx context.Context, originalName string, object schema.AdminObject, expectedRevision int64) (RouteBinding, error) {
	originalName, err := normalizeRouteBindingName(originalName)
	if err != nil {
		return RouteBinding{}, err
	}
	normalized, err := s.Normalize(object)
	if err != nil {
		return RouteBinding{}, err
	}
	if normalized.Metadata.Name != originalName {
		return RouteBinding{}, fmt.Errorf("%w: %q cannot be changed to %q", ErrRouteBindingNameImmutable, originalName, normalized.Metadata.Name)
	}
	return s.SaveDraft(ctx, normalized, false, expectedRevision)
}

func validateSaveDraftRequest(create bool, expectedRevision int64) error {
	if expectedRevision < 0 {
		return errors.New("expected revision must not be negative")
	}
	if create && expectedRevision > 0 {
		return errors.New("create route binding does not accept an expected revision")
	}
	return nil
}

func validateDraftSaveState(name string, create bool, expectedRevision int64, exists bool, revision int64) error {
	if create && exists {
		return fmt.Errorf("%w: %s", ErrRouteBindingAlreadyExists, name)
	}
	if expectedRevision > 0 && (!exists || revision != expectedRevision) {
		return fmt.Errorf("%w: %s", ErrRouteBindingConflict, name)
	}
	return nil
}

func (s *RouteBindingStore) resolveRuntimeIdentity(ctx context.Context, name string, draft routeBindingEntry, draftExists bool) (int, int, error) {
	resourceID, methodID := draft.record.ResourceID, draft.record.MethodID
	if !draftExists {
		published, exists, _, err := s.getEntry(ctx, s.bindingPrefix(false), name)
		if err != nil {
			return 0, 0, err
		}
		if exists {
			resourceID, methodID = published.record.ResourceID, published.record.MethodID
		}
	}
	if resourceID <= 0 {
		var err error
		resourceID, err = s.allocateRuntimeID(ctx)
		if err != nil {
			return 0, 0, err
		}
	}
	if methodID <= 0 {
		// One AdminRouteBinding currently compiles to exactly one Method. Method
		// IDs are scoped below the resource, so sharing the stable route ID is
		// safe and avoids a second allocation transaction.
		methodID = resourceID
	}
	return resourceID, methodID, nil
}

func (s *RouteBindingStore) commitDraft(ctx context.Context, name, draftPrefix string, value []byte, create bool, expectedRevision int64) error {
	comparisons := make([]clientv3.Cmp, 0, 1)
	if create {
		comparisons = append(comparisons, clientv3.Compare(clientv3.CreateRevision(s.bindingKey(draftPrefix, name)), "=", 0))
	} else if expectedRevision > 0 {
		comparisons = append(comparisons, clientv3.Compare(clientv3.ModRevision(s.bindingKey(draftPrefix, name)), "=", expectedRevision))
	}
	operations := []clientv3.Op{
		clientv3.OpPut(s.bindingKey(draftPrefix, name), string(value)),
	}
	transaction := s.kv.Txn(ctx)
	if len(comparisons) > 0 {
		transaction = transaction.If(comparisons...)
	}
	response, err := transaction.Then(operations...).Commit()
	if err != nil {
		return fmt.Errorf("save route binding %q: %w", name, err)
	}
	if response.Succeeded {
		return nil
	}
	if create {
		return fmt.Errorf("%w: %s", ErrRouteBindingAlreadyExists, name)
	}
	return fmt.Errorf("%w: %s", ErrRouteBindingConflict, name)
}

// DeleteAndPublish removes one route draft and its published runtime state in
// the same etcd transaction. This keeps the delete flow atomic for callers:
// there is no intermediate state where the draft is gone but the route still
// remains active in Pixiu.
func (s *RouteBindingStore) DeleteAndPublish(ctx context.Context, name string, expectedRevision int64) (RouteBindingPublishResult, error) {
	ctx = routeBindingContext(ctx)
	if expectedRevision < 0 {
		return RouteBindingPublishResult{}, errors.New("expected revision must not be negative")
	}
	name, err := normalizeRouteBindingName(name)
	if err != nil {
		return RouteBindingPublishResult{}, err
	}
	state, err := s.readPublishState(ctx, name)
	if err != nil {
		return RouteBindingPublishResult{}, err
	}
	if expectedRevision > 0 && (!state.draftExists || state.draftRevision != expectedRevision) {
		return RouteBindingPublishResult{}, fmt.Errorf("%w: %s", ErrRouteBindingConflict, name)
	}
	comparisons, err := s.buildPublishComparisons(ctx, state)
	if err != nil {
		return RouteBindingPublishResult{}, err
	}
	operations := make([]clientv3.Op, 0, 4)
	if state.draftExists {
		operations = append(operations, clientv3.OpDelete(s.bindingKey(state.draftPrefix, name)))
	}
	deletedCount := 0
	if state.publishedExists {
		if state.published.record.ResourceID <= 0 || state.published.record.MethodID <= 0 {
			return RouteBindingPublishResult{}, fmt.Errorf("published route binding %q has invalid runtime identity", name)
		}
		operations = append(operations,
			clientv3.OpDelete(s.runtimeResourceKey(state.published.record.ResourceID), clientv3.WithPrefix()),
			clientv3.OpDelete(s.bindingKey(state.publishedPrefix, name)),
		)
		deletedCount = 1
	}
	operations = append(operations,
		clientv3.OpPut(s.publishedRevisionKey(), strconv.FormatInt(time.Now().UnixNano(), 10)))
	response, err := s.kv.Txn(ctx).If(comparisons...).Then(operations...).Commit()
	if err != nil {
		return RouteBindingPublishResult{}, fmt.Errorf("delete route binding %q: %w", name, err)
	}
	if !response.Succeeded {
		return RouteBindingPublishResult{}, fmt.Errorf("%w: %s", ErrRouteBindingPublishConflict, name)
	}
	return RouteBindingPublishResult{
		Name:              name,
		Revision:          response.Header.Revision,
		DraftRevision:     state.draftRevision,
		PublishedRevision: 0,
		DeletedCount:      deletedCount,
	}, nil
}

// Publish publishes exactly one route binding. The generated legacy resource
// and method, the published high-level binding, and the internal generation
// guard are written or deleted in one etcd transaction. Other routes' draft
// and published values are left untouched.
func (s *RouteBindingStore) Publish(ctx context.Context, name string, expectedDraftRevision int64) (RouteBindingPublishResult, error) {
	ctx = routeBindingContext(ctx)
	if expectedDraftRevision < 0 {
		return RouteBindingPublishResult{}, errors.New("expected draft revision must not be negative")
	}
	name, err := normalizeRouteBindingName(name)
	if err != nil {
		return RouteBindingPublishResult{}, err
	}
	state, err := s.readPublishState(ctx, name)
	if err != nil {
		return RouteBindingPublishResult{}, err
	}
	if expectedDraftRevision > 0 && (!state.draftExists || state.draftRevision != expectedDraftRevision) {
		return RouteBindingPublishResult{}, fmt.Errorf("%w: %s", ErrRouteBindingPublishConflict, name)
	}
	publishedEntries, err := s.listEntries(ctx, state.publishedPrefix)
	if err != nil {
		return RouteBindingPublishResult{}, err
	}
	candidate, err := s.compilePublishCandidate(name, state.draft, state.draftExists, publishedEntries)
	if err != nil {
		return RouteBindingPublishResult{}, err
	}
	if err := s.checkRuntimeRouteConflict(ctx, candidate, state); err != nil {
		return RouteBindingPublishResult{}, err
	}
	comparisons, err := s.buildPublishComparisons(ctx, state)
	if err != nil {
		return RouteBindingPublishResult{}, err
	}
	operations, publishedCount, deletedCount, err := s.buildSinglePublishOperations(state, candidate)
	if err != nil {
		return RouteBindingPublishResult{}, err
	}
	operations = append(operations,
		clientv3.OpPut(s.publishedRevisionKey(), strconv.FormatInt(time.Now().UnixNano(), 10)))
	response, err := s.kv.Txn(ctx).If(comparisons...).Then(operations...).Commit()
	if err != nil {
		return RouteBindingPublishResult{}, fmt.Errorf("publish route binding %q: %w", name, err)
	}
	if !response.Succeeded {
		return RouteBindingPublishResult{}, fmt.Errorf("%w: %s", ErrRouteBindingPublishConflict, name)
	}
	publishedRevision, err := s.publishedRouteRevision(ctx, state.publishedPrefix, name, candidate)
	if err != nil {
		return RouteBindingPublishResult{}, err
	}
	return RouteBindingPublishResult{
		Name:              name,
		Revision:          response.Header.Revision,
		DraftRevision:     state.draftRevision,
		PublishedRevision: publishedRevision,
		PublishedCount:    publishedCount,
		DeletedCount:      deletedCount,
	}, nil
}

type routeBindingPublishState struct {
	draftPrefix       string
	publishedPrefix   string
	draft             routeBindingEntry
	draftExists       bool
	draftRevision     int64
	published         routeBindingEntry
	publishedExists   bool
	publishedRevision int64
}

func (s *RouteBindingStore) readPublishState(ctx context.Context, name string) (routeBindingPublishState, error) {
	state := routeBindingPublishState{
		draftPrefix:     s.bindingPrefix(true),
		publishedPrefix: s.bindingPrefix(false),
	}
	var err error
	state.draft, state.draftExists, state.draftRevision, err = s.getEntry(ctx, state.draftPrefix, name)
	if err != nil {
		return routeBindingPublishState{}, err
	}
	state.published, state.publishedExists, state.publishedRevision, err = s.getEntry(ctx, state.publishedPrefix, name)
	if err != nil {
		return routeBindingPublishState{}, err
	}
	if !state.draftExists && !state.publishedExists {
		return routeBindingPublishState{}, fmt.Errorf("%w: %s", ErrRouteBindingNotFound, name)
	}
	return state, nil
}

func (s *RouteBindingStore) compilePublishCandidate(name string, draft routeBindingEntry, draftExists bool, published []routeBindingEntry) (*compiledRouteEntry, error) {
	if !draftExists {
		return nil, nil
	}
	entries := make([]routeBindingEntry, 0, len(published)+1)
	for _, entry := range published {
		if entry.record.Object.Metadata.Name != name {
			entries = append(entries, entry)
		}
	}
	entries = append(entries, draft)
	compiled, err := s.compileSnapshot(entries)
	if err != nil {
		return nil, err
	}
	for index := range compiled {
		if compiled[index].object.Metadata.Name == name {
			return &compiled[index], nil
		}
	}
	return nil, fmt.Errorf("route binding %q was not compiled", name)
}

func (s *RouteBindingStore) buildPublishComparisons(ctx context.Context, state routeBindingPublishState) ([]clientv3.Cmp, error) {
	comparisons := make([]clientv3.Cmp, 0, 2)
	if state.draftExists {
		// Always compare the value observed above, even when the caller omitted
		// expectedDraftRevision. This closes the read/modify/write race.
		comparisons = append(comparisons,
			clientv3.Compare(clientv3.ModRevision(s.bindingKey(state.draftPrefix, state.draft.record.Object.Metadata.Name)), "=", state.draftRevision))
		if state.publishedExists {
			comparisons = append(comparisons,
				clientv3.Compare(clientv3.ModRevision(s.bindingKey(state.publishedPrefix, state.published.record.Object.Metadata.Name)), "=", state.publishedRevision))
		} else {
			comparisons = append(comparisons,
				clientv3.Compare(clientv3.CreateRevision(s.bindingKey(state.publishedPrefix, state.draft.record.Object.Metadata.Name)), "=", 0))
		}
	} else {
		comparisons = append(comparisons,
			clientv3.Compare(clientv3.ModRevision(s.bindingKey(state.publishedPrefix, state.published.record.Object.Metadata.Name)), "=", state.publishedRevision))
	}
	publishedMarkerRevision, publishedMarkerExists, err := s.markerRevision(ctx, s.publishedRevisionKey())
	if err != nil {
		return nil, err
	}
	if publishedMarkerExists {
		comparisons = append(comparisons,
			clientv3.Compare(clientv3.ModRevision(s.publishedRevisionKey()), "=", publishedMarkerRevision))
	} else {
		comparisons = append(comparisons,
			clientv3.Compare(clientv3.CreateRevision(s.publishedRevisionKey()), "=", 0))
	}
	return comparisons, nil
}

func (s *RouteBindingStore) buildSinglePublishOperations(state routeBindingPublishState, candidate *compiledRouteEntry) ([]clientv3.Op, int, int, error) {
	if candidate == nil {
		return []clientv3.Op{
			clientv3.OpDelete(s.runtimeResourceKey(state.published.record.ResourceID), clientv3.WithPrefix()),
			clientv3.OpDelete(s.bindingKey(state.publishedPrefix, state.published.record.Object.Metadata.Name)),
		}, 0, 1, nil
	}
	compiled := candidate.legacy
	operations := make([]clientv3.Op, 0, 4)
	if state.publishedExists && state.published.record.ResourceID != compiled.Resource.ID {
		operations = append(operations,
			clientv3.OpDelete(s.runtimeResourceKey(state.published.record.ResourceID), clientv3.WithPrefix()))
	} else if state.publishedExists && state.published.record.MethodID != compiled.Method.ID {
		operations = append(operations,
			clientv3.OpDelete(s.runtimeMethodKey(state.published.record.ResourceID, state.published.record.MethodID)))
	}
	resourceValue, methodValue, err := marshalCompiledRuntimeRoute(compiled)
	if err != nil {
		return nil, 0, 0, fmt.Errorf("encode route binding %q for runtime: %w", candidate.object.Metadata.Name, err)
	}
	recordValue, err := json.Marshal(routeBindingRecord{
		Object:     candidate.object,
		ResourceID: compiled.Resource.ID,
		MethodID:   compiled.Method.ID,
	})
	if err != nil {
		return nil, 0, 0, fmt.Errorf("encode published route binding %q: %w", candidate.object.Metadata.Name, err)
	}
	operations = append(operations,
		clientv3.OpPut(s.runtimeResourceKey(compiled.Resource.ID), string(resourceValue)),
		clientv3.OpPut(s.runtimeMethodKey(compiled.Resource.ID, compiled.Method.ID), string(methodValue)),
		clientv3.OpPut(s.bindingKey(state.publishedPrefix, candidate.object.Metadata.Name), string(recordValue)))
	return operations, 1, 0, nil
}

func (s *RouteBindingStore) publishedRouteRevision(ctx context.Context, prefix, name string, candidate *compiledRouteEntry) (int64, error) {
	if candidate == nil {
		return 0, nil
	}
	_, exists, revision, err := s.getEntry(ctx, prefix, name)
	if err != nil {
		return 0, err
	}
	if !exists {
		return 0, fmt.Errorf("route binding %q disappeared after publish", name)
	}
	return revision, nil
}

// checkRuntimeRouteConflict protects the breaking migration boundary between
// the legacy Resource/Method keys and AdminRouteBinding. Existing runtime
// method keys are still consumed by Pixiu's watcher, even when they are not
// visible through the new Admin API. Publishing a new route with the same
// HTTP method and path would otherwise commit successfully and leave the
// watcher with duplicate runtime routes.
func (s *RouteBindingStore) checkRuntimeRouteConflict(ctx context.Context, candidate *compiledRouteEntry, state routeBindingPublishState) error {
	if candidate == nil {
		return nil
	}
	prefix := s.runtimeResourcePrefix()
	response, err := s.kv.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("list existing runtime resources and methods: %w", err)
	}
	resources := make(map[int]legacyconfig.Resource)
	for _, kv := range response.Kvs {
		resourceID, ok := runtimeResourceID(prefix, string(kv.Key))
		if !ok {
			continue
		}
		var resource legacyconfig.Resource
		if err := commonyaml.UnmarshalYML(kv.Value, &resource); err != nil {
			return fmt.Errorf("decode existing runtime resource %q: %w", string(kv.Key), err)
		}
		resources[resourceID] = resource
	}

	resourcePaths := make(map[int]string)
	runtimeRoutes := make([]runtimeRouteBinding, 0, len(response.Kvs))
	for resourceID, resource := range resources {
		collectRuntimeResourceRoutes(resource, "", resourceID, resourcePaths, &runtimeRoutes)
	}
	for _, kv := range response.Kvs {
		resourceID, methodID, ok, err := runtimeMethodIDs(prefix, string(kv.Key))
		if err != nil {
			return err
		}
		if !ok {
			continue
		}
		// A republish of the same published AdminRouteBinding owns this runtime
		// identity; a new route must not bypass the check by reusing an orphaned
		// legacy method key with the same numeric IDs.
		if state.publishedExists && resourceID == state.published.record.ResourceID && methodID == state.published.record.MethodID {
			continue
		}
		var method legacyconfig.Method
		if err := commonyaml.UnmarshalYML(kv.Value, &method); err != nil {
			return fmt.Errorf("decode existing runtime method %q: %w", string(kv.Key), err)
		}
		path := strings.TrimSpace(method.ResourcePath)
		if path == "" {
			path = resourcePaths[resourceID]
		}
		if path != "" {
			runtimeRoutes = append(runtimeRoutes, runtimeRouteBinding{path: path, httpVerb: method.HTTPVerb, key: string(kv.Key)})
		}
	}

	candidatePath := strings.ToLower(candidate.legacy.Resource.Path)
	candidateVerb := strings.TrimSpace(candidate.legacy.Method.HTTPVerb)
	for _, route := range runtimeRoutes {
		if strings.ToLower(route.path) == candidatePath && strings.EqualFold(route.httpVerb, candidateVerb) {
			return fmt.Errorf("%w: %s %s already exists at %q", ErrRouteBindingRuntimeConflict, candidateVerb, candidate.legacy.Resource.Path, route.key)
		}
	}
	return nil
}

type runtimeRouteBinding struct {
	path     string
	httpVerb string
	key      string
}

func collectRuntimeResourceRoutes(resource legacyconfig.Resource, parentPath string, keyResourceID int, resourcePaths map[int]string, routes *[]runtimeRouteBinding) {
	groupPath := parentPath
	if groupPath == "/" {
		groupPath = ""
	}
	if !strings.HasPrefix(resource.Path, "/") {
		return
	}
	fullPath := groupPath + resource.Path
	if keyResourceID > 0 {
		resourcePaths[keyResourceID] = fullPath
	}
	if resource.ID > 0 {
		resourcePaths[resource.ID] = fullPath
	}
	resourceKey := fmt.Sprintf("resource %d", keyResourceID)
	for _, method := range resource.Methods {
		*routes = append(*routes, runtimeRouteBinding{path: fullPath, httpVerb: method.HTTPVerb, key: resourceKey})
	}
	for _, nested := range resource.Resources {
		// Match addAPIFromResource: nested resources are visited with their
		// parent's declared path, and the router lowercases paths at insertion.
		collectRuntimeResourceRoutes(nested, resource.Path, nested.ID, resourcePaths, routes)
	}
}

// Status returns draft and published revisions for one route. The revisions
// are independent of other route bindings in the same config namespace.
func (s *RouteBindingStore) Status(ctx context.Context, name string) (RouteBindingPublishStatus, error) {
	ctx = routeBindingContext(ctx)
	name, err := normalizeRouteBindingName(name)
	if err != nil {
		return RouteBindingPublishStatus{}, err
	}
	draft, draftExists, draftRevision, err := s.getEntry(ctx, s.bindingPrefix(true), name)
	if err != nil {
		return RouteBindingPublishStatus{}, err
	}
	published, publishedExists, publishedRevision, err := s.getEntry(ctx, s.bindingPrefix(false), name)
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
func (s *RouteBindingStore) Diff(ctx context.Context, name string) (RouteBindingDiff, error) {
	ctx = routeBindingContext(ctx)
	name, err := normalizeRouteBindingName(name)
	if err != nil {
		return RouteBindingDiff{}, err
	}
	draft, draftExists, _, err := s.getEntry(ctx, s.bindingPrefix(true), name)
	if err != nil {
		return RouteBindingDiff{}, err
	}
	published, publishedExists, _, err := s.getEntry(ctx, s.bindingPrefix(false), name)
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

func (s *RouteBindingStore) listEntries(ctx context.Context, prefix string) ([]routeBindingEntry, error) {
	response, err := s.kv.Get(ctx, prefix, clientv3.WithPrefix(), clientv3.WithSort(clientv3.SortByKey, clientv3.SortAscend))
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

func (s *RouteBindingStore) getEntry(ctx context.Context, prefix, name string) (routeBindingEntry, bool, int64, error) {
	key := s.bindingKey(prefix, name)
	response, err := s.kv.Get(ctx, key)
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

// markerRevision reads the internal published generation used to protect a
// single-route publish from publishing against a stale published snapshot.
// It is intentionally not exposed as a global publish status API.
func (s *RouteBindingStore) markerRevision(ctx context.Context, key string) (int64, bool, error) {
	response, err := s.kv.Get(ctx, key)
	if err != nil {
		return 0, false, fmt.Errorf("get route binding revision marker %q: %w", key, err)
	}
	if len(response.Kvs) == 0 {
		return 0, false, nil
	}
	return response.Kvs[0].ModRevision, true, nil
}

type routeBindingIDCounter struct {
	value    int
	revision int64
}

func (s *RouteBindingStore) allocateRuntimeID(ctx context.Context) (int, error) {
	counterKey := s.root + "/" + ResourceID
	resourcePrefix := s.root + "/" + routeBindingRuntimeResourceDir + "/"
	for attempt := 0; attempt < maxRouteBindingIDRetries; attempt++ {
		id, allocated, err := s.tryAllocateRuntimeID(ctx, counterKey, resourcePrefix)
		if err != nil {
			return 0, err
		}
		if allocated {
			return id, nil
		}
	}
	return 0, fmt.Errorf("%w: resource id allocation retries exhausted", ErrRouteBindingConflict)
}

func (s *RouteBindingStore) tryAllocateRuntimeID(ctx context.Context, counterKey, resourcePrefix string) (int, bool, error) {
	counter, err := s.readRuntimeIDCounter(ctx, counterKey)
	if err != nil {
		return 0, false, err
	}
	resourceResponse, err := s.kv.Get(ctx, resourcePrefix, clientv3.WithPrefix())
	if err != nil {
		return 0, false, fmt.Errorf("scan route binding resource ids: %w", err)
	}
	next := maxRuntimeResourceID(resourceResponse, resourcePrefix, counter.value) + 1
	comparison := clientv3.Compare(clientv3.CreateRevision(counterKey), "=", 0)
	if counter.revision > 0 {
		comparison = clientv3.Compare(clientv3.ModRevision(counterKey), "=", counter.revision)
	}
	response, err := s.kv.Txn(ctx).
		If(comparison).
		Then(clientv3.OpPut(counterKey, strconv.Itoa(next))).
		Commit()
	if err != nil {
		return 0, false, fmt.Errorf("allocate route binding resource id: %w", err)
	}
	return next, response.Succeeded, nil
}

func (s *RouteBindingStore) readRuntimeIDCounter(ctx context.Context, counterKey string) (routeBindingIDCounter, error) {
	response, err := s.kv.Get(ctx, counterKey)
	if err != nil {
		return routeBindingIDCounter{}, fmt.Errorf("read route binding resource id counter: %w", err)
	}
	counter := routeBindingIDCounter{}
	if len(response.Kvs) == 0 {
		return counter, nil
	}
	value := string(response.Kvs[0].Value)
	current, err := strconv.Atoi(value)
	if err != nil || current < 0 {
		return routeBindingIDCounter{}, fmt.Errorf("invalid resource id counter %q", value)
	}
	counter.value = current
	counter.revision = response.Kvs[0].ModRevision
	return counter, nil
}

func maxRuntimeResourceID(response *clientv3.GetResponse, resourcePrefix string, current int) int {
	maxExisting := current
	for _, kv := range response.Kvs {
		relative := strings.TrimPrefix(string(kv.Key), resourcePrefix)
		if strings.Contains(relative, "/") {
			continue
		}
		id, parseErr := strconv.Atoi(relative)
		if parseErr == nil && id > maxExisting {
			maxExisting = id
		}
	}
	return maxExisting
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

func (s *RouteBindingStore) publishedRevisionKey() string {
	return s.root + "/" + routeBindingRevisionSegment
}

func (s *RouteBindingStore) runtimeResourceKey(resourceID int) string {
	return s.root + "/" + routeBindingRuntimeResourceDir + "/" + strconv.Itoa(resourceID)
}

func (s *RouteBindingStore) runtimeResourcePrefix() string {
	return s.root + "/" + routeBindingRuntimeResourceDir + "/"
}

func (s *RouteBindingStore) runtimeMethodKey(resourceID, methodID int) string {
	return s.runtimeResourceKey(resourceID) + "/" + Method + "/" + strconv.Itoa(methodID)
}

func runtimeResourceID(prefix, key string) (int, bool) {
	if !strings.HasPrefix(key, prefix) {
		return 0, false
	}
	relative := strings.Trim(strings.TrimPrefix(key, prefix), "/")
	if relative == "" || strings.Contains(relative, "/") {
		return 0, false
	}
	resourceID, err := strconv.Atoi(relative)
	if err != nil {
		return 0, false
	}
	return resourceID, true
}

func runtimeMethodIDs(prefix, key string) (int, int, bool, error) {
	if !strings.HasPrefix(key, prefix) {
		return 0, 0, false, nil
	}
	parts := strings.Split(strings.Trim(strings.TrimPrefix(key, prefix), "/"), "/")
	if len(parts) != 3 || parts[1] != Method {
		return 0, 0, false, nil
	}
	resourceID, err := strconv.Atoi(parts[0])
	if err != nil {
		return 0, 0, true, fmt.Errorf("invalid runtime resource key %q: %w", key, err)
	}
	methodID, err := strconv.Atoi(parts[2])
	if err != nil {
		return 0, 0, true, fmt.Errorf("invalid runtime method key %q: %w", key, err)
	}
	return resourceID, methodID, true, nil
}

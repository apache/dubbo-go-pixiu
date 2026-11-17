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
	"bytes"
	"context"
	"crypto/rand"
	"encoding/hex"
	"encoding/json"
	"errors"
	"fmt"
	"reflect"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

import (
	clientv3 "go.etcd.io/etcd/client/v3"

	"gopkg.in/yaml.v3"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	legacyconfig "github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/config/schema"
)

const clusterDraftSegment = "unpublished/cluster-drafts"

var (
	ErrClusterDraftNotFound   = errors.New("cluster draft not found")
	ErrClusterDraftConflict   = errors.New("cluster draft was modified concurrently")
	ErrClusterRuntimeNotFound = errors.New("runtime cluster not found")
	clusterDraftIDPattern     = regexp.MustCompile(`^[0-9a-f]{32}$`)
)

type clusterDraftKV interface {
	Get(context.Context, string, ...clientv3.OpOption) (*clientv3.GetResponse, error)
	Txn(context.Context) clientv3.Txn
}

type clusterDraftRecord struct {
	Object schema.AdminObject `json:"object"`
	// SourceID is assigned only when a draft is created from a published key.
	// It never comes from the user-editable AdminObject.
	SourceID     int    `json:"sourceId,omitempty"`
	SourceDigest string `json:"sourceDigest,omitempty"`
}

type ClusterDraft struct {
	ID           string             `json:"id"`
	Object       schema.AdminObject `json:"object"`
	SourceID     int                `json:"sourceId,omitempty"`
	Revision     int64              `json:"revision"`
	sourceDigest string
}

type ClusterDraftStatus struct {
	ID                string `json:"id"`
	SourceID          int    `json:"sourceId,omitempty"`
	DraftRevision     int64  `json:"draftRevision"`
	PublishedRevision int64  `json:"publishedRevision"`
	DraftExists       bool   `json:"draftExists"`
	PublishedExists   bool   `json:"publishedExists"`
	Dirty             bool   `json:"dirty"`
}

type ClusterDiffChange struct {
	Path   string `json:"path"`
	Before any    `json:"before"`
	After  any    `json:"after"`
}

type ClusterDraftDiff struct {
	ID        string              `json:"id"`
	Draft     ClusterDraft        `json:"draft"`
	Published *schema.AdminObject `json:"published,omitempty"`
	Changes   []ClusterDiffChange `json:"changes"`
}

type ClusterDraftStore struct {
	kv    clusterDraftKV
	root  string
	newID func() (string, error)
}

func NewClusterDraftStore(kv clusterDraftKV, root string) (*ClusterDraftStore, error) {
	if kv == nil {
		return nil, errors.New("cluster draft etcd client is nil")
	}
	root = strings.TrimRight(strings.TrimSpace(root), "/")
	if root == "" {
		return nil, errors.New("cluster draft etcd root path is empty")
	}
	return &ClusterDraftStore{kv: kv, root: root, newID: randomClusterDraftID}, nil
}

func NewAdminClusterDraftStore() (*ClusterDraftStore, error) {
	if adminconfig.Client == nil || adminconfig.Bootstrap == nil || adminconfig.Client.GetRawClient() == nil {
		return nil, errors.New("admin cluster draft etcd client is not initialized")
	}
	return NewClusterDraftStore(adminconfig.Client.GetRawClient(), adminconfig.Bootstrap.GetPath())
}

func randomClusterDraftID() (string, error) {
	var bytes [16]byte
	if _, err := rand.Read(bytes[:]); err != nil {
		return "", fmt.Errorf("generate cluster draft id: %w", err)
	}
	return hex.EncodeToString(bytes[:]), nil
}

func clusterDraftObject(object schema.AdminObject) (schema.AdminObject, error) {
	if object.Kind == "" {
		object.Kind = schema.KindAdminCluster
	}
	if object.Kind != schema.KindAdminCluster {
		return schema.AdminObject{}, fmt.Errorf("cluster draft kind must be %q", schema.KindAdminCluster)
	}
	return object, nil
}

func (s *ClusterDraftStore) draftKey(id string) string {
	return s.root + "/" + clusterDraftSegment + "/" + id
}
func (s *ClusterDraftStore) draftPrefix() string { return s.root + "/" + clusterDraftSegment + "/" }
func (s *ClusterDraftStore) runtimeKey(id int) string {
	return s.root + "/" + Clusters + "/" + strconv.Itoa(id)
}

func validClusterDraftID(id string) bool { return clusterDraftIDPattern.MatchString(id) }

// CreateDraft stores an object without business validation. An empty name and
// incomplete spec are valid drafts; Validate/Preview apply strict rules.
func (s *ClusterDraftStore) CreateDraft(ctx context.Context, object schema.AdminObject) (ClusterDraft, error) {
	object, err := clusterDraftObject(object)
	if err != nil {
		return ClusterDraft{}, err
	}
	return s.create(ctx, clusterDraftRecord{Object: object})
}

// CreateFromPublished seeds a draft from an existing runtime Cluster. The
// runtime key, rather than a caller-supplied form field, binds the identity.
func (s *ClusterDraftStore) CreateFromPublished(ctx context.Context, sourceID int) (ClusterDraft, error) {
	if sourceID <= 0 {
		return ClusterDraft{}, errors.New("sourceId must be positive")
	}
	published, _, exists, err := s.getPublished(ctx, sourceID)
	if err != nil {
		return ClusterDraft{}, err
	}
	if !exists {
		return ClusterDraft{}, fmt.Errorf("%w: %d", ErrClusterRuntimeNotFound, sourceID)
	}
	return s.create(ctx, clusterDraftRecord{Object: schema.AdminClusterFromLegacy(published), SourceID: sourceID, SourceDigest: clusterRuntimeDigest(published)})
}

func (s *ClusterDraftStore) create(ctx context.Context, record clusterDraftRecord) (ClusterDraft, error) {
	value, err := json.Marshal(record)
	if err != nil {
		return ClusterDraft{}, fmt.Errorf("encode cluster draft: %w", err)
	}
	ctx = routeBindingContext(ctx)
	for attempt := 0; attempt < 8; attempt++ {
		id, err := s.newID()
		if err != nil {
			return ClusterDraft{}, err
		}
		if !validClusterDraftID(id) {
			return ClusterDraft{}, errors.New("generated cluster draft id is invalid")
		}
		key := s.draftKey(id)
		response, err := s.kv.Txn(ctx).If(clientv3.Compare(clientv3.CreateRevision(key), "=", 0)).Then(clientv3.OpPut(key, string(value))).Commit()
		if err != nil {
			return ClusterDraft{}, fmt.Errorf("create cluster draft: %w", err)
		}
		if response.Succeeded {
			return ClusterDraft{ID: id, Object: record.Object, SourceID: record.SourceID, Revision: response.Header.Revision, sourceDigest: record.SourceDigest}, nil
		}
	}
	return ClusterDraft{}, fmt.Errorf("%w: draft id allocation retries exhausted", ErrClusterDraftConflict)
}

// UpdateDraft requires the revision returned by the last read or write.
func (s *ClusterDraftStore) UpdateDraft(ctx context.Context, id string, object schema.AdminObject, expectedRevision int64) (ClusterDraft, error) {
	if !validClusterDraftID(id) {
		return ClusterDraft{}, errors.New("cluster draft id is invalid")
	}
	if expectedRevision <= 0 {
		return ClusterDraft{}, errors.New("positive expectedRevision is required")
	}
	object, err := clusterDraftObject(object)
	if err != nil {
		return ClusterDraft{}, err
	}
	current, err := s.GetDraft(ctx, id)
	if err != nil {
		return ClusterDraft{}, err
	}
	if current.Revision != expectedRevision {
		return ClusterDraft{}, fmt.Errorf("%w: %s", ErrClusterDraftConflict, id)
	}
	value, err := json.Marshal(clusterDraftRecord{Object: object, SourceID: current.SourceID, SourceDigest: current.sourceDigest})
	if err != nil {
		return ClusterDraft{}, fmt.Errorf("encode cluster draft: %w", err)
	}
	response, err := s.kv.Txn(routeBindingContext(ctx)).If(clientv3.Compare(clientv3.ModRevision(s.draftKey(id)), "=", expectedRevision)).Then(clientv3.OpPut(s.draftKey(id), string(value))).Commit()
	if err != nil {
		return ClusterDraft{}, fmt.Errorf("update cluster draft %q: %w", id, err)
	}
	if !response.Succeeded {
		return ClusterDraft{}, fmt.Errorf("%w: %s", ErrClusterDraftConflict, id)
	}
	return ClusterDraft{ID: id, Object: object, SourceID: current.SourceID, Revision: response.Header.Revision, sourceDigest: current.sourceDigest}, nil
}

func (s *ClusterDraftStore) GetDraft(ctx context.Context, id string) (ClusterDraft, error) {
	if !validClusterDraftID(id) {
		return ClusterDraft{}, errors.New("cluster draft id is invalid")
	}
	response, err := s.kv.Get(routeBindingContext(ctx), s.draftKey(id))
	if err != nil {
		return ClusterDraft{}, fmt.Errorf("get cluster draft %q: %w", id, err)
	}
	if len(response.Kvs) == 0 {
		return ClusterDraft{}, fmt.Errorf("%w: %s", ErrClusterDraftNotFound, id)
	}
	var record clusterDraftRecord
	if err := json.Unmarshal(response.Kvs[0].Value, &record); err != nil {
		return ClusterDraft{}, fmt.Errorf("decode cluster draft %q: %w", id, err)
	}
	return ClusterDraft{ID: id, Object: record.Object, SourceID: record.SourceID, Revision: response.Kvs[0].ModRevision, sourceDigest: record.SourceDigest}, nil
}

func (s *ClusterDraftStore) ListDrafts(ctx context.Context) ([]ClusterDraft, error) {
	response, err := s.kv.Get(routeBindingContext(ctx), s.draftPrefix(), clientv3.WithPrefix())
	if err != nil {
		return nil, fmt.Errorf("list cluster drafts: %w", err)
	}
	result := make([]ClusterDraft, 0, len(response.Kvs))
	for _, kv := range response.Kvs {
		id := strings.TrimPrefix(string(kv.Key), s.draftPrefix())
		if !validClusterDraftID(id) {
			continue
		}
		var record clusterDraftRecord
		if err := json.Unmarshal(kv.Value, &record); err != nil {
			return nil, fmt.Errorf("decode cluster draft %q: %w", id, err)
		}
		result = append(result, ClusterDraft{ID: id, Object: record.Object, SourceID: record.SourceID, Revision: kv.ModRevision, sourceDigest: record.SourceDigest})
	}
	sort.Slice(result, func(i, j int) bool { return result[i].ID < result[j].ID })
	return result, nil
}

func (s *ClusterDraftStore) getPublished(ctx context.Context, sourceID int) (legacyconfig.Cluster, int64, bool, error) {
	response, err := s.kv.Get(routeBindingContext(ctx), s.runtimeKey(sourceID))
	if err != nil {
		return legacyconfig.Cluster{}, 0, false, fmt.Errorf("get runtime cluster %d: %w", sourceID, err)
	}
	if len(response.Kvs) == 0 {
		return legacyconfig.Cluster{}, 0, false, nil
	}
	var cluster legacyconfig.Cluster
	if err := yaml.Unmarshal(response.Kvs[0].Value, &cluster); err != nil {
		return legacyconfig.Cluster{}, 0, false, fmt.Errorf("decode runtime cluster %d: %w", sourceID, err)
	}
	return cluster, response.Kvs[0].ModRevision, true, nil
}

func (s *ClusterDraftStore) Status(ctx context.Context, id string) (ClusterDraftStatus, error) {
	draft, err := s.GetDraft(ctx, id)
	if err != nil {
		return ClusterDraftStatus{}, err
	}
	status := ClusterDraftStatus{ID: id, SourceID: draft.SourceID, DraftRevision: draft.Revision, DraftExists: true, Dirty: true}
	if draft.SourceID == 0 {
		return status, nil
	}
	published, revision, exists, err := s.getPublished(ctx, draft.SourceID)
	if err != nil {
		return ClusterDraftStatus{}, err
	}
	status.PublishedExists = exists
	status.PublishedRevision = revision
	if exists {
		status.Dirty = !clusterObjectsEqual(schema.AdminClusterFromLegacy(published), draft.Object)
	}
	return status, nil
}

func (s *ClusterDraftStore) Diff(ctx context.Context, id string) (ClusterDraftDiff, error) {
	draft, err := s.GetDraft(ctx, id)
	if err != nil {
		return ClusterDraftDiff{}, err
	}
	var published *schema.AdminObject
	if draft.SourceID > 0 {
		cluster, _, exists, err := s.getPublished(ctx, draft.SourceID)
		if err != nil {
			return ClusterDraftDiff{}, err
		}
		if exists {
			object := schema.AdminClusterFromLegacy(cluster)
			published = &object
		}
	}
	changes, err := diffClusterObjects(published, draft.Object)
	if err != nil {
		return ClusterDraftDiff{}, err
	}
	return ClusterDraftDiff{ID: id, Draft: draft, Published: published, Changes: changes}, nil
}

func clusterObjectsEqual(before, after schema.AdminObject) bool {
	left, leftErr := json.Marshal(before)
	right, rightErr := json.Marshal(after)
	return leftErr == nil && rightErr == nil && bytes.Equal(left, right)
}

func diffClusterObjects(before *schema.AdminObject, after schema.AdminObject) ([]ClusterDiffChange, error) {
	beforeValues := make(map[string]any)
	if before != nil {
		if err := flattenClusterObject(*before, beforeValues); err != nil {
			return nil, err
		}
	}
	afterValues := make(map[string]any)
	if err := flattenClusterObject(after, afterValues); err != nil {
		return nil, err
	}
	paths := make([]string, 0, len(beforeValues)+len(afterValues))
	for path := range beforeValues {
		paths = append(paths, path)
	}
	for path := range afterValues {
		if _, exists := beforeValues[path]; !exists {
			paths = append(paths, path)
		}
	}
	sort.Strings(paths)
	changes := make([]ClusterDiffChange, 0)
	for _, path := range paths {
		beforeValue, beforeExists := beforeValues[path]
		afterValue, afterExists := afterValues[path]
		if beforeExists && afterExists && reflect.DeepEqual(beforeValue, afterValue) {
			continue
		}
		changes = append(changes, ClusterDiffChange{Path: path, Before: valueOrNil(beforeValue, beforeExists), After: valueOrNil(afterValue, afterExists)})
	}
	return changes, nil
}

func flattenClusterObject(object schema.AdminObject, values map[string]any) error {
	encoded, err := json.Marshal(object)
	if err != nil {
		return fmt.Errorf("encode cluster diff object: %w", err)
	}
	var canonical map[string]any
	if err := json.Unmarshal(encoded, &canonical); err != nil {
		return fmt.Errorf("decode cluster diff object: %w", err)
	}
	return flattenAdminObjectValue("", canonical, values)
}

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

package logic

import (
	"context"
	"errors"
	"fmt"
	"regexp"
	"sort"
	"strconv"
	"strings"
)

import (
	clientv3 "go.etcd.io/etcd/client/v3"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	commonyaml "github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
)

const legacyResources = "resources"

func legacyResourceKey(id string, unpublished bool) string {
	if unpublished {
		return getUnpublishedRootPath(legacyResources) + "/" + id
	}
	return getRootPath(legacyResources) + "/" + id
}

func legacyMethodPrefix(resourceID string, unpublished bool) string {
	return legacyResourceKey(resourceID, unpublished) + "/" + Method
}

func legacyMethodKey(resourceID, methodID string, unpublished bool) string {
	return legacyMethodPrefix(resourceID, unpublished) + "/" + methodID
}

func legacyResourceKeyPattern(unpublished bool) *regexp.Regexp {
	root := getRootPath(legacyResources)
	if unpublished {
		root = getUnpublishedRootPath(legacyResources)
	}
	return regexp.MustCompile("^" + regexp.QuoteMeta(root) + `/[^/]+/?$`)
}

// BizGetResourceList restores the original Resource API over the runtime key
// layout. RouteBinding records are not included because they use a separate
// namespace and are compiled to these keys only when published.
func BizGetResourceList(unpublished bool) ([]config.Resource, error) {
	keys, values, err := adminconfig.Client.GetChildrenKVList(func() string {
		if unpublished {
			return getUnpublishedRootPath(legacyResources)
		}
		return getRootPath(legacyResources)
	}())
	if err != nil {
		return nil, err
	}
	if len(keys) != len(values) {
		return nil, fmt.Errorf("legacy resource key/value count mismatch: %d != %d", len(keys), len(values))
	}
	resources := make([]config.Resource, 0, len(values))
	pattern := legacyResourceKeyPattern(unpublished)
	for index, key := range keys {
		if !pattern.MatchString(key) {
			continue
		}
		var resource config.Resource
		if err := commonyaml.UnmarshalYML([]byte(values[index]), &resource); err != nil {
			return nil, fmt.Errorf("decode legacy resource %q: %w", key, err)
		}
		resources = append(resources, resource)
	}
	return resources, nil
}

func BizGetResourceDetail(id string, unpublished bool) (string, error) {
	return adminconfig.Client.Get(legacyResourceKey(id, unpublished))
}

func BizGetResourceDetailWithRevision(id string, unpublished bool) (string, int64, error) {
	client, err := legacyResourceClient()
	if err != nil {
		return "", 0, err
	}
	response, err := client.Get(adminconfig.Client.GetCtx(), legacyResourceKey(id, unpublished))
	if err != nil {
		return "", 0, fmt.Errorf("read legacy resource %q: %w", id, err)
	}
	if len(response.Kvs) == 0 {
		return "", 0, fmt.Errorf("legacy resource %q not found", id)
	}
	return string(response.Kvs[0].Value), response.Kvs[0].ModRevision, nil
}

// BizResourceExists checks an exact legacy Resource key without conflating a
// missing record with an etcd read failure.
func BizResourceExists(id string, unpublished bool) (bool, error) {
	id, err := normalizePositiveLegacyID(id, ResourceID)
	if err != nil {
		return false, err
	}
	if adminconfig.Client == nil {
		return false, errors.New("admin etcd client is unavailable")
	}
	client := adminconfig.Client.GetRawClient()
	if client == nil {
		return false, errors.New("admin etcd client is unavailable")
	}
	response, err := client.Get(adminconfig.Client.GetCtx(), legacyResourceKey(id, unpublished))
	if err != nil {
		return false, fmt.Errorf("read legacy resource %q: %w", id, err)
	}
	return len(response.Kvs) > 0, nil
}

func BizSetResourceInfo(resource *config.Resource, created, unpublished bool) error {
	storedResource, methods, resourceData, err := prepareLegacyResourceInfo(resource, created)
	if err != nil {
		return err
	}
	if unpublished {
		return saveLegacyResourceDraft(resource, storedResource, methods, resourceData, created)
	}
	return publishLegacyResourceInfo(resource, storedResource, methods, resourceData, created)
}

func prepareLegacyResourceInfo(resource *config.Resource, created bool) (config.Resource, []config.Method, []byte, error) {
	if resource == nil {
		return config.Resource{}, nil, nil, errors.New("legacy resource is nil")
	}
	methods := append([]config.Method(nil), resource.Methods...)
	storedResource := *resource
	storedResource.Methods = nil
	if created && storedResource.ID <= 0 {
		storedResource.ID = nextLegacyResourceID()
		if storedResource.ID <= 0 {
			return config.Resource{}, nil, nil, errors.New("could not allocate legacy resource id")
		}
	} else if storedResource.ID <= 0 {
		return config.Resource{}, nil, nil, errors.New("legacy resource id must be positive")
	}
	resource.ID = storedResource.ID
	data, err := commonyaml.MarshalYML(storedResource)
	if err != nil {
		return config.Resource{}, nil, nil, fmt.Errorf("encode legacy resource: %w", err)
	}
	return storedResource, methods, data, nil
}

func saveLegacyResourceDraft(resource *config.Resource, storedResource config.Resource, methods []config.Method, data []byte, created bool) error {
	resourceID := strconv.Itoa(storedResource.ID)
	key := legacyResourceKey(resourceID, true)
	var err error
	if created {
		err = adminconfig.Client.Create(key, string(data))
	} else {
		err = adminconfig.Client.Update(key, string(data))
	}
	if err != nil {
		return fmt.Errorf("save legacy resource: %w", err)
	}
	if !created {
		return nil
	}
	for index := range methods {
		methods[index].ResourcePath = storedResource.Path
		if err := BizSetResourceMethod(resourceID, &methods[index], true, true); err != nil {
			return err
		}
		resource.Methods[index] = methods[index]
	}
	return nil
}

func publishLegacyResourceInfo(resource *config.Resource, storedResource config.Resource, methods []config.Method, data []byte, created bool) error {
	key := legacyResourceKey(strconv.Itoa(storedResource.ID), false)
	operations := []clientv3.Op{clientv3.OpPut(key, string(data))}
	comparisons := make([]clientv3.Cmp, 0, len(methods)+1)
	if created {
		comparisons = append(comparisons, clientv3.Compare(clientv3.CreateRevision(key), "=", 0))
		for index := range methods {
			if methods[index].ID <= 0 {
				methods[index].ID = nextLegacyMethodID()
				if methods[index].ID <= 0 {
					return errors.New("could not allocate legacy method id")
				}
			}
			methods[index].ResourcePath = storedResource.Path
			methodData, err := commonyaml.MarshalYML(methods[index])
			if err != nil {
				return fmt.Errorf("encode legacy method: %w", err)
			}
			methodKey := legacyMethodKey(strconv.Itoa(storedResource.ID), strconv.Itoa(methods[index].ID), false)
			operations = append(operations, clientv3.OpPut(methodKey, string(methodData)))
			comparisons = append(comparisons, clientv3.Compare(clientv3.CreateRevision(methodKey), "=", 0))
			resource.Methods[index] = methods[index]
		}
	}
	if err := publishLegacyResourceOperations(operations, comparisons...); err != nil {
		return fmt.Errorf("save legacy resource: %w", err)
	}
	return nil
}

// BizCreatePublishedResourceInfo writes a legacy Resource and its Methods to
// both the published and draft namespaces in the same transaction. The
// published route checks and generation update remain part of that transaction.
func BizCreatePublishedResourceInfo(resource *config.Resource) error {
	if resource == nil {
		return errors.New("legacy resource is nil")
	}
	if adminconfig.Client == nil {
		return errors.New("admin etcd client is unavailable")
	}
	if resource.ID <= 0 {
		resource.ID = nextLegacyResourceID()
		if resource.ID <= 0 {
			return errors.New("could not allocate legacy resource id")
		}
	}
	resourceID := strconv.Itoa(resource.ID)
	methods := append([]config.Method(nil), resource.Methods...)
	for index := range methods {
		if methods[index].ID <= 0 {
			methods[index].ID = nextLegacyMethodID()
			if methods[index].ID <= 0 {
				return errors.New("could not allocate legacy method id")
			}
		}
		methods[index].ResourcePath = resource.Path
		resource.Methods[index] = methods[index]
	}

	storedResource := *resource
	storedResource.Methods = nil
	resourceData, err := commonyaml.MarshalYML(storedResource)
	if err != nil {
		return fmt.Errorf("encode legacy resource: %w", err)
	}
	runtimeOps := []clientv3.Op{clientv3.OpPut(legacyResourceKey(resourceID, false), string(resourceData))}
	draftOps := []clientv3.Op{clientv3.OpPut(legacyResourceKey(resourceID, true), string(resourceData))}
	comparisons := []clientv3.Cmp{
		clientv3.Compare(clientv3.CreateRevision(legacyResourceKey(resourceID, false)), "=", 0),
		clientv3.Compare(clientv3.CreateRevision(legacyResourceKey(resourceID, true)), "=", 0),
	}
	for _, method := range methods {
		methodData, err := commonyaml.MarshalYML(method)
		if err != nil {
			return fmt.Errorf("encode legacy method %d: %w", method.ID, err)
		}
		runtimeKey := legacyMethodKey(resourceID, strconv.Itoa(method.ID), false)
		draftKey := legacyMethodKey(resourceID, strconv.Itoa(method.ID), true)
		runtimeOps = append(runtimeOps, clientv3.OpPut(runtimeKey, string(methodData)))
		draftOps = append(draftOps, clientv3.OpPut(draftKey, string(methodData)))
		comparisons = append(comparisons,
			clientv3.Compare(clientv3.CreateRevision(runtimeKey), "=", 0),
			clientv3.Compare(clientv3.CreateRevision(draftKey), "=", 0),
		)
	}
	client, err := legacyResourceClient()
	if err != nil {
		return err
	}
	store, err := legacyRouteBindingStore(client)
	if err != nil {
		return err
	}
	return store.publishLegacyResourceOperationsWithAdditional(adminconfig.Client.GetCtx(), runtimeOps, draftOps, comparisons)
}

func BizDeleteResourceInfo(id string, unpublished bool) error {
	id, err := normalizePositiveLegacyID(id, ResourceID)
	if err != nil {
		return err
	}
	if adminconfig.Client == nil {
		return errors.New("admin etcd client is unavailable")
	}
	client := adminconfig.Client.GetRawClient()
	if client == nil {
		return errors.New("admin etcd client is unavailable")
	}
	key := legacyResourceKey(id, unpublished)
	if unpublished {
		return deleteLegacyResource(adminconfig.Client.GetCtx(), client, key)
	}
	store, err := legacyRouteBindingStore(client)
	if err != nil {
		return err
	}
	return store.deletePublishedLegacyResource(adminconfig.Client.GetCtx(), id)
}

func deleteLegacyResource(ctx context.Context, kv routeBindingKV, key string) error {
	_, err := kv.Txn(ctx).Then(
		clientv3.OpDelete(key),
		clientv3.OpDelete(key+"/", clientv3.WithPrefix()),
	).Commit()
	return err
}

func BizGetMethodList(resourceID string, unpublished bool) ([]config.Method, error) {
	_, values, err := adminconfig.Client.GetChildrenKVList(legacyMethodPrefix(resourceID, unpublished))
	if err != nil {
		return nil, err
	}
	methods := make([]config.Method, 0, len(values))
	for _, value := range values {
		var method config.Method
		if err := commonyaml.UnmarshalYML([]byte(value), &method); err != nil {
			return nil, fmt.Errorf("decode legacy method: %w", err)
		}
		methods = append(methods, method)
	}
	return methods, nil
}

type LegacyKeyRevision struct {
	Key      string
	Revision int64
	Exists   bool
}

// LegacyMethodSnapshot holds method records and the revisions read with them
// in one etcd transaction.
type LegacyMethodSnapshot struct {
	Methods          []config.Method
	Revisions        []LegacyKeyRevision
	SnapshotRevision LegacyKeyRevision
}

func BizGetMethodListWithRevisions(resourceID string, unpublished bool) (LegacyMethodSnapshot, error) {
	client, err := legacyResourceClient()
	if err != nil {
		return LegacyMethodSnapshot{}, err
	}
	prefix := legacyMethodPrefix(resourceID, unpublished) + "/"
	store, err := legacyRouteBindingStore(client)
	if err != nil {
		return LegacyMethodSnapshot{}, err
	}
	revisionKey := store.legacyDraftMethodRevisionKey(resourceID)
	if !unpublished {
		revisionKey = store.publishedRevisionKey()
	}
	response, err := client.Txn(adminconfig.Client.GetCtx()).Then(
		clientv3.OpGet(prefix, clientv3.WithPrefix()),
		clientv3.OpGet(revisionKey),
	).Commit()
	if err != nil {
		return LegacyMethodSnapshot{}, fmt.Errorf("read legacy methods for resource %q: %w", resourceID, err)
	}
	methodResponse := response.Responses[0].GetResponseRange()
	revisionResponse := response.Responses[1].GetResponseRange()
	methods := make([]config.Method, 0, len(methodResponse.Kvs))
	revisions := make([]LegacyKeyRevision, 0, len(methodResponse.Kvs))
	for _, item := range methodResponse.Kvs {
		relativeKey := strings.TrimPrefix(string(item.Key), prefix)
		if relativeKey == "" || strings.Contains(relativeKey, "/") {
			continue
		}
		var method config.Method
		if err := commonyaml.UnmarshalYML(item.Value, &method); err != nil {
			return LegacyMethodSnapshot{}, fmt.Errorf("decode legacy method %q: %w", item.Key, err)
		}
		keyID, err := normalizePositiveLegacyID(relativeKey, MethodID)
		if err != nil {
			return LegacyMethodSnapshot{}, fmt.Errorf("invalid legacy method key %q: %w", item.Key, err)
		}
		keyIDNumber, _ := strconv.Atoi(keyID)
		if method.ID != keyIDNumber {
			return LegacyMethodSnapshot{}, fmt.Errorf("legacy method key %q contains ID %d", item.Key, method.ID)
		}
		methods = append(methods, method)
		revisions = append(revisions, LegacyKeyRevision{Key: string(item.Key), Revision: item.ModRevision})
	}
	methodRevision := LegacyKeyRevision{Key: revisionKey}
	if len(revisionResponse.Kvs) > 0 {
		methodRevision.Revision = revisionResponse.Kvs[0].ModRevision
		methodRevision.Exists = true
	}
	return LegacyMethodSnapshot{
		Methods:          methods,
		Revisions:        revisions,
		SnapshotRevision: methodRevision,
	}, nil
}

// BizUpdateResourceAndMethods updates a legacy Resource and its Methods in one
// etcd transaction. This keeps Method.ResourcePath aligned with the Resource
// when a legacy client changes the route path.
func BizUpdateResourceAndMethods(resource *config.Resource, methodSnapshot LegacyMethodSnapshot, unpublished bool, resourceRevision int64) error {
	if resource == nil {
		return errors.New("legacy resource is nil")
	}
	resourceID, err := normalizePositiveLegacyID(strconv.Itoa(resource.ID), ResourceID)
	if err != nil {
		return err
	}
	client, err := legacyResourceClient()
	if err != nil {
		return err
	}
	if resourceRevision <= 0 {
		return errors.New("legacy resource revision must be positive")
	}
	resourceData, err := marshalLegacyResourceForPathUpdate(resource)
	if err != nil {
		return err
	}
	store, err := legacyRouteBindingStore(client)
	if err != nil {
		return err
	}
	operations, comparisons, err := buildLegacyResourceAndMethodUpdate(resource, resourceID, methodSnapshot, unpublished, resourceRevision, resourceData, store)
	if err != nil {
		return err
	}
	if err := store.updateLegacyResourceAndMethods(adminconfig.Client.GetCtx(), operations, comparisons, unpublished); err != nil {
		return fmt.Errorf("save legacy resource and methods atomically: %w", err)
	}
	return nil
}

func marshalLegacyResourceForPathUpdate(resource *config.Resource) ([]byte, error) {
	storedResource := *resource
	storedResource.Methods = nil
	resourceData, err := commonyaml.MarshalYML(storedResource)
	if err != nil {
		return nil, fmt.Errorf("encode legacy resource: %w", err)
	}
	return resourceData, nil
}

func buildLegacyResourceAndMethodUpdate(
	resource *config.Resource,
	resourceID string,
	methodSnapshot LegacyMethodSnapshot,
	unpublished bool,
	resourceRevision int64,
	resourceData []byte,
	store *RouteBindingStore,
) ([]clientv3.Op, []clientv3.Cmp, error) {
	operations := []clientv3.Op{clientv3.OpPut(legacyResourceKey(resourceID, unpublished), string(resourceData))}
	comparisons := []clientv3.Cmp{clientv3.Compare(clientv3.ModRevision(legacyResourceKey(resourceID, unpublished)), "=", resourceRevision)}
	revisionByKey, snapshotComparisons, err := buildLegacyMethodSnapshotComparisons(resourceID, unpublished, store, methodSnapshot.SnapshotRevision, methodSnapshot.Revisions)
	if err != nil {
		return nil, nil, err
	}
	comparisons = append(comparisons, snapshotComparisons...)
	methodOperations, err := buildLegacyMethodUpdateOperations(resource, resourceID, methodSnapshot.Methods, unpublished, revisionByKey)
	if err != nil {
		return nil, nil, err
	}
	operations = append(operations, methodOperations...)
	return operations, comparisons, nil
}

func buildLegacyMethodUpdateOperations(
	resource *config.Resource,
	resourceID string,
	methods []config.Method,
	unpublished bool,
	revisionByKey map[string]int64,
) ([]clientv3.Op, error) {
	operations := make([]clientv3.Op, 0, len(methods))
	for index := range methods {
		if methods[index].ID <= 0 {
			return nil, fmt.Errorf("legacy method id must be positive, got %d", methods[index].ID)
		}
		methodKey := legacyMethodKey(resourceID, strconv.Itoa(methods[index].ID), unpublished)
		if _, exists := revisionByKey[methodKey]; !exists {
			return nil, fmt.Errorf("legacy method %d was not included in the path update snapshot", methods[index].ID)
		}
		methods[index].ResourcePath = resource.Path
		methodData, err := commonyaml.MarshalYML(methods[index])
		if err != nil {
			return nil, fmt.Errorf("encode legacy method %d: %w", methods[index].ID, err)
		}
		operations = append(operations, clientv3.OpPut(methodKey, string(methodData)))
	}
	return operations, nil
}

func buildLegacyMethodSnapshotComparisons(
	resourceID string,
	unpublished bool,
	store *RouteBindingStore,
	snapshotRevision LegacyKeyRevision,
	methodRevisions []LegacyKeyRevision,
) (map[string]int64, []clientv3.Cmp, error) {
	expectedSnapshotKey := store.publishedRevisionKey()
	if unpublished {
		expectedSnapshotKey = store.legacyDraftMethodRevisionKey(resourceID)
	}
	if snapshotRevision.Key != expectedSnapshotKey ||
		(snapshotRevision.Exists && snapshotRevision.Revision <= 0) ||
		(!snapshotRevision.Exists && snapshotRevision.Revision != 0) {
		return nil, nil, fmt.Errorf("invalid legacy method snapshot revision for resource %q", resourceID)
	}
	revisionByKey := make(map[string]int64, len(methodRevisions))
	for _, revision := range methodRevisions {
		if revision.Revision <= 0 || !strings.HasPrefix(revision.Key, legacyMethodPrefix(resourceID, unpublished)+"/") {
			return nil, nil, fmt.Errorf("invalid legacy method revision for resource %q", resourceID)
		}
		revisionByKey[revision.Key] = revision.Revision
	}
	comparisons := make([]clientv3.Cmp, 0, len(revisionByKey)+1)
	if snapshotRevision.Exists {
		comparisons = append(comparisons, clientv3.Compare(clientv3.ModRevision(snapshotRevision.Key), "=", snapshotRevision.Revision))
	} else {
		comparisons = append(comparisons, clientv3.Compare(clientv3.CreateRevision(snapshotRevision.Key), "=", 0))
	}
	for key, revision := range revisionByKey {
		comparisons = append(comparisons, clientv3.Compare(clientv3.ModRevision(key), "=", revision))
	}
	return revisionByKey, comparisons, nil
}

func BizGetMethodDetail(resourceID, methodID string, unpublished bool) (string, error) {
	return adminconfig.Client.Get(legacyMethodKey(resourceID, methodID, unpublished))
}

func BizSetResourceMethod(resourceID string, method *config.Method, created, unpublished bool) error {
	if method == nil {
		return errors.New("legacy method is nil")
	}
	if created && method.ID <= 0 {
		method.ID = nextLegacyMethodID()
		if method.ID <= 0 {
			return errors.New("could not allocate legacy method id")
		}
	} else if method.ID <= 0 {
		return errors.New("legacy method id must be positive")
	}
	resourceID, err := normalizePositiveLegacyID(resourceID, ResourceID)
	if err != nil {
		return err
	}
	client, err := legacyResourceClient()
	if err != nil {
		return err
	}
	resourceKey := legacyResourceKey(resourceID, unpublished)
	resourceResponse, err := client.Get(adminconfig.Client.GetCtx(), resourceKey)
	if err != nil {
		return fmt.Errorf("read legacy resource %q before saving Method: %w", resourceID, err)
	}
	if len(resourceResponse.Kvs) == 0 {
		return fmt.Errorf("legacy resource %q not found", resourceID)
	}
	var resource config.Resource
	if err := commonyaml.UnmarshalYML(resourceResponse.Kvs[0].Value, &resource); err != nil {
		return fmt.Errorf("decode legacy resource %q before saving Method: %w", resourceID, err)
	}
	method.ResourcePath = resource.Path
	data, err := commonyaml.MarshalYML(method)
	if err != nil {
		return fmt.Errorf("encode legacy method: %w", err)
	}
	key := legacyMethodKey(resourceID, strconv.Itoa(method.ID), unpublished)
	if unpublished {
		store, err := legacyRouteBindingStore(client)
		if err != nil {
			return err
		}
		return store.saveLegacyDraftMethod(adminconfig.Client.GetCtx(), resourceID, key, string(data), resourceResponse.Kvs[0].ModRevision, created)
	}
	comparisons := []clientv3.Cmp{clientv3.Compare(clientv3.ModRevision(resourceKey), "=", resourceResponse.Kvs[0].ModRevision)}
	if created {
		comparisons = append(comparisons, clientv3.Compare(clientv3.CreateRevision(key), "=", 0))
	}
	if err := publishLegacyResourceOperations([]clientv3.Op{clientv3.OpPut(key, string(data))}, comparisons...); err != nil {
		return fmt.Errorf("save legacy method: %w", err)
	}
	return nil
}

func (s *RouteBindingStore) saveLegacyDraftMethod(ctx context.Context, resourceID, methodKey, value string, resourceRevision int64, created bool) error {
	resourceKey := s.legacyDraftResourceKey(resourceID)
	revisionKey := s.legacyDraftMethodRevisionKey(resourceID)
	revisionResponse, err := s.kv.Get(routeBindingContext(ctx), revisionKey)
	if err != nil {
		return fmt.Errorf("read draft legacy method revision: %w", err)
	}
	methodResponse, err := s.kv.Get(routeBindingContext(ctx), methodKey)
	if err != nil {
		return fmt.Errorf("read draft legacy Method: %w", err)
	}
	comparisons := []clientv3.Cmp{clientv3.Compare(clientv3.ModRevision(resourceKey), "=", resourceRevision)}
	if len(revisionResponse.Kvs) > 0 {
		comparisons = append(comparisons, clientv3.Compare(clientv3.ModRevision(revisionKey), "=", revisionResponse.Kvs[0].ModRevision))
	} else {
		comparisons = append(comparisons, clientv3.Compare(clientv3.CreateRevision(revisionKey), "=", 0))
	}
	if created {
		comparisons = append(comparisons, clientv3.Compare(clientv3.CreateRevision(methodKey), "=", 0))
	} else if len(methodResponse.Kvs) > 0 {
		comparisons = append(comparisons, clientv3.Compare(clientv3.ModRevision(methodKey), "=", methodResponse.Kvs[0].ModRevision))
	} else {
		return fmt.Errorf("draft legacy Method %q not found", methodKey)
	}
	response, err := s.kv.Txn(routeBindingContext(ctx)).If(comparisons...).Then(
		clientv3.OpPut(methodKey, value),
		clientv3.OpPut(revisionKey, "1"),
	).Commit()
	if err != nil {
		return fmt.Errorf("save draft legacy Method: %w", err)
	}
	if !response.Succeeded {
		return fmt.Errorf("%w: draft legacy Method changed during update", ErrRouteBindingPublishConflict)
	}
	return nil
}

func (s *RouteBindingStore) deleteLegacyDraftMethod(ctx context.Context, resourceID, methodKey string) error {
	resourceKey := s.legacyDraftResourceKey(resourceID)
	revisionKey := s.legacyDraftMethodRevisionKey(resourceID)
	resourceResponse, err := s.kv.Get(routeBindingContext(ctx), resourceKey)
	if err != nil {
		return fmt.Errorf("read draft legacy Resource before deleting Method: %w", err)
	}
	methodResponse, err := s.kv.Get(routeBindingContext(ctx), methodKey)
	if err != nil {
		return fmt.Errorf("read draft legacy Method before deletion: %w", err)
	}
	revisionResponse, err := s.kv.Get(routeBindingContext(ctx), revisionKey)
	if err != nil {
		return fmt.Errorf("read draft legacy method revision before deletion: %w", err)
	}
	comparisons := make([]clientv3.Cmp, 0, 3)
	if len(resourceResponse.Kvs) > 0 {
		comparisons = append(comparisons, clientv3.Compare(clientv3.ModRevision(resourceKey), "=", resourceResponse.Kvs[0].ModRevision))
	} else {
		comparisons = append(comparisons, clientv3.Compare(clientv3.CreateRevision(resourceKey), "=", 0))
	}
	if len(methodResponse.Kvs) > 0 {
		comparisons = append(comparisons, clientv3.Compare(clientv3.ModRevision(methodKey), "=", methodResponse.Kvs[0].ModRevision))
	} else {
		comparisons = append(comparisons, clientv3.Compare(clientv3.CreateRevision(methodKey), "=", 0))
	}
	if len(revisionResponse.Kvs) > 0 {
		comparisons = append(comparisons, clientv3.Compare(clientv3.ModRevision(revisionKey), "=", revisionResponse.Kvs[0].ModRevision))
	} else {
		comparisons = append(comparisons, clientv3.Compare(clientv3.CreateRevision(revisionKey), "=", 0))
	}
	operations := []clientv3.Op{clientv3.OpDelete(methodKey), clientv3.OpPut(revisionKey, "1")}
	response, err := s.kv.Txn(routeBindingContext(ctx)).If(comparisons...).Then(operations...).Commit()
	if err != nil {
		return fmt.Errorf("delete draft legacy Method: %w", err)
	}
	if !response.Succeeded {
		return fmt.Errorf("%w: draft legacy Method changed during deletion", ErrRouteBindingPublishConflict)
	}
	return nil
}

func BizDeleteMethodInfo(resourceID, methodID string, unpublished bool) error {
	var err error
	resourceID, err = normalizePositiveLegacyID(resourceID, ResourceID)
	if err != nil {
		return err
	}
	methodID, err = normalizePositiveLegacyID(methodID, MethodID)
	if err != nil {
		return err
	}
	if adminconfig.Client == nil {
		return errors.New("admin etcd client is unavailable")
	}
	if unpublished {
		client, err := legacyResourceClient()
		if err != nil {
			return err
		}
		store, err := legacyRouteBindingStore(client)
		if err != nil {
			return err
		}
		return store.deleteLegacyDraftMethod(adminconfig.Client.GetCtx(), resourceID, legacyMethodKey(resourceID, methodID, true))
	}
	client, err := legacyResourceClient()
	if err != nil {
		return err
	}
	store, err := legacyRouteBindingStore(client)
	if err != nil {
		return err
	}
	resourceNumber, _ := strconv.Atoi(resourceID)
	methodNumber, _ := strconv.Atoi(methodID)
	return store.deletePublishedLegacyMethod(adminconfig.Client.GetCtx(), resourceNumber, methodNumber)
}

// BizMethodExists checks an exact legacy Method key without conflating a
// missing record with an etcd read failure.
func BizMethodExists(resourceID, methodID string, unpublished bool) (bool, error) {
	var err error
	resourceID, err = normalizePositiveLegacyID(resourceID, ResourceID)
	if err != nil {
		return false, err
	}
	methodID, err = normalizePositiveLegacyID(methodID, MethodID)
	if err != nil {
		return false, err
	}
	if adminconfig.Client == nil {
		return false, errors.New("admin etcd client is unavailable")
	}
	client := adminconfig.Client.GetRawClient()
	if client == nil {
		return false, errors.New("admin raw etcd client is unavailable")
	}
	response, err := client.Get(adminconfig.Client.GetCtx(), legacyMethodKey(resourceID, methodID, unpublished))
	if err != nil {
		return false, fmt.Errorf("read legacy method %q/%q: %w", resourceID, methodID, err)
	}
	return len(response.Kvs) > 0, nil
}

// BizBatchReleaseResources publishes each staged legacy Resource and its
// Method records in one etcd transaction. The legacy API can address either
// Resource or Method publication, so both endpoints share route conflict and
// generation checks with AdminRouteBinding publishing.
func BizBatchReleaseResources() error {
	client, err := legacyResourceClient()
	if err != nil {
		return err
	}
	fromPrefix := getUnpublishedRootPath(legacyResources) + "/"
	response, err := client.Get(adminconfig.Client.GetCtx(), fromPrefix, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("read staged legacy resources: %w", err)
	}
	grouped := make(map[string][]clientv3.Op)
	for _, item := range response.Kvs {
		resourceID, operation, err := legacyResourcePublishOperation(item.Key, item.Value, fromPrefix)
		if err != nil {
			return err
		}
		grouped[resourceID] = append(grouped[resourceID], operation)
	}

	return publishLegacyResourceGroups(grouped)
}

func legacyResourceClient() (*clientv3.Client, error) {
	if adminconfig.Client == nil {
		return nil, errors.New("admin etcd client is unavailable")
	}
	client := adminconfig.Client.GetRawClient()
	if client == nil {
		return nil, errors.New("admin etcd client is unavailable")
	}
	return client, nil
}

func legacyRouteBindingStore(client *clientv3.Client) (*RouteBindingStore, error) {
	if client == nil {
		return nil, errors.New("admin etcd client is unavailable")
	}
	if adminconfig.Bootstrap == nil {
		return nil, errors.New("admin bootstrap is unavailable")
	}
	root := strings.TrimRight(strings.TrimSpace(adminconfig.Bootstrap.GetPath()), "/")
	if root == "" {
		return nil, errors.New("admin etcd root path is empty")
	}
	return &RouteBindingStore{kv: client, root: root}, nil
}

func (s *RouteBindingStore) legacyDraftResourceKey(resourceID string) string {
	return s.root + "/" + Unpublished + "/" + legacyResources + "/" + resourceID
}

func (s *RouteBindingStore) legacyDraftMethodRevisionKey(resourceID string) string {
	return s.root + "/" + Unpublished + "/legacy-method-revisions/" + resourceID
}

func legacyResourcePublishOperation(key, value []byte, fromPrefix string) (string, clientv3.Op, error) {
	keyString := string(key)
	relative := strings.TrimPrefix(keyString, fromPrefix)
	parts := strings.SplitN(relative, "/", 2)
	if relative == "" || strings.HasPrefix(relative, "/") {
		return "", clientv3.Op{}, fmt.Errorf("invalid staged legacy resource key %q", keyString)
	}
	resourceID, err := normalizePositiveLegacyID(parts[0], ResourceID)
	if err != nil {
		return "", clientv3.Op{}, fmt.Errorf("invalid staged legacy resource key %q: %w", keyString, err)
	}
	if len(parts) == 2 {
		methodParts := strings.Split(parts[1], "/")
		if len(methodParts) != 2 || methodParts[0] != Method {
			return "", clientv3.Op{}, fmt.Errorf("invalid staged legacy resource key %q", keyString)
		}
		if _, err := normalizePositiveLegacyID(methodParts[1], MethodID); err != nil {
			return "", clientv3.Op{}, fmt.Errorf("invalid staged legacy resource key %q: %w", keyString, err)
		}
	}
	operation := clientv3.OpPut(getRootPath(legacyResources)+"/"+relative, string(value))
	return resourceID, operation, nil
}

func publishLegacyResourceGroups(grouped map[string][]clientv3.Op) error {
	resourceIDs := make([]string, 0, len(grouped))
	for resourceID := range grouped {
		resourceIDs = append(resourceIDs, resourceID)
	}
	sort.Strings(resourceIDs)
	for _, resourceID := range resourceIDs {
		ops := grouped[resourceID]
		sort.Slice(ops, func(i, j int) bool { return string(ops[i].KeyBytes()) < string(ops[j].KeyBytes()) })
		if err := publishLegacyResourceOperations(ops); err != nil {
			return fmt.Errorf("publish legacy resource %q atomically: %w", resourceID, err)
		}
	}
	return nil
}

func publishLegacyResourceOperations(operations []clientv3.Op, comparisons ...clientv3.Cmp) error {
	client, err := legacyResourceClient()
	if err != nil {
		return err
	}
	store, err := legacyRouteBindingStore(client)
	if err != nil {
		return err
	}
	return store.publishLegacyResourceOperations(adminconfig.Client.GetCtx(), operations, comparisons)
}

func (s *RouteBindingStore) publishLegacyResourceOperations(ctx context.Context, operations []clientv3.Op, comparisons []clientv3.Cmp) error {
	return s.publishLegacyResourceOperationsWithAdditional(ctx, operations, nil, comparisons)
}

func (s *RouteBindingStore) updateLegacyResourceAndMethods(ctx context.Context, operations []clientv3.Op, comparisons []clientv3.Cmp, unpublished bool) error {
	if !unpublished {
		return s.publishLegacyResourceOperations(ctx, operations, comparisons)
	}
	response, err := s.kv.Txn(routeBindingContext(ctx)).If(comparisons...).Then(operations...).Commit()
	if err != nil {
		return fmt.Errorf("save legacy resource and methods atomically: %w", err)
	}
	if !response.Succeeded {
		return fmt.Errorf("%w: legacy resource changed during path update", ErrRouteBindingPublishConflict)
	}
	return nil
}

func (s *RouteBindingStore) publishLegacyResourceOperationsWithAdditional(ctx context.Context, operations, additionalOperations []clientv3.Op, comparisons []clientv3.Cmp) error {
	if len(operations) == 0 {
		return nil
	}
	ctx = routeBindingContext(ctx)
	markerRevision, markerExists, err := s.markerRevision(ctx, s.publishedRevisionKey())
	if err != nil {
		return err
	}
	generation, err := s.nextPublishedGeneration(ctx)
	if err != nil {
		return err
	}
	prefix := s.runtimeResourcePrefix()
	runtimeResponse, err := s.kv.Get(ctx, prefix, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("list existing runtime routes before legacy publish: %w", err)
	}
	resources, err := runtimeResourcesFromResponse(prefix, runtimeResponse)
	if err != nil {
		return err
	}
	resourcePaths, runtimeRoutes := indexRuntimeResourceRoutes(resources, len(runtimeResponse.Kvs))
	methodRoutes, err := runtimeMethodRoutesFromResponse(prefix, runtimeResponse, resourcePaths, routeBindingPublishState{})
	if err != nil {
		return err
	}
	runtimeRoutes = append(runtimeRoutes, methodRoutes...)
	candidateRoutes, err := legacyResourcePublishRoutes(operations, runtimeResponse, prefix, resourcePaths)
	if err != nil {
		return err
	}
	if err := s.rejectLegacyAdminRouteBindingWrites(ctx, candidateRoutes); err != nil {
		return err
	}
	if err := findLegacyRuntimeRouteConflict(candidateRoutes, runtimeRoutes); err != nil {
		return err
	}
	if markerExists {
		comparisons = append(comparisons, clientv3.Compare(clientv3.ModRevision(s.publishedRevisionKey()), "=", markerRevision))
	} else {
		comparisons = append(comparisons, clientv3.Compare(clientv3.CreateRevision(s.publishedRevisionKey()), "=", 0))
	}
	allOperations := append(append([]clientv3.Op(nil), operations...), additionalOperations...)
	allOperations = append(allOperations, clientv3.OpPut(s.publishedRevisionKey(), strconv.FormatInt(generation, 10)))
	response, err := s.kv.Txn(ctx).If(comparisons...).Then(allOperations...).Commit()
	if err != nil {
		return fmt.Errorf("publish legacy resources atomically: %w", err)
	}
	if !response.Succeeded {
		return fmt.Errorf("%w: legacy resource publish", ErrRouteBindingPublishConflict)
	}
	return nil
}

func (s *RouteBindingStore) deletePublishedLegacyResource(ctx context.Context, resourceID string) error {
	normalizedID, err := normalizePositiveLegacyID(resourceID, ResourceID)
	if err != nil {
		return err
	}
	resourceNumber, _ := strconv.Atoi(normalizedID)
	return s.deletePublishedLegacyRuntime(ctx, resourceNumber, nil)
}

func (s *RouteBindingStore) deletePublishedLegacyMethod(ctx context.Context, resourceID, methodID int) error {
	if resourceID <= 0 || methodID <= 0 {
		return errors.New("resourceId and methodId must be positive integers")
	}
	return s.deletePublishedLegacyRuntime(ctx, resourceID, &methodID)
}

func (s *RouteBindingStore) deletePublishedLegacyRuntime(ctx context.Context, resourceID int, methodID *int) error {
	ctx = routeBindingContext(ctx)
	markerRevision, markerExists, err := s.markerRevision(ctx, s.publishedRevisionKey())
	if err != nil {
		return err
	}
	generation, err := s.nextPublishedGeneration(ctx)
	if err != nil {
		return err
	}
	publishedBindings, err := s.listEntries(ctx, s.bindingPrefix(false))
	if err != nil {
		return err
	}
	for _, binding := range publishedBindings {
		if binding.record.ResourceID != resourceID {
			continue
		}
		if methodID == nil || binding.record.MethodID == *methodID {
			return fmt.Errorf("%w: runtime identity %d/%d is owned by AdminRouteBinding %q", ErrRouteBindingRuntimeConflict, resourceID, binding.record.MethodID, binding.record.Object.Metadata.Name)
		}
	}

	comparisons := make([]clientv3.Cmp, 0, 1)
	if markerExists {
		comparisons = append(comparisons, clientv3.Compare(clientv3.ModRevision(s.publishedRevisionKey()), "=", markerRevision))
	} else {
		comparisons = append(comparisons, clientv3.Compare(clientv3.CreateRevision(s.publishedRevisionKey()), "=", 0))
	}
	resourceKey := s.runtimeResourceKey(resourceID)
	operations := []clientv3.Op{}
	if methodID == nil {
		operations = append(operations,
			clientv3.OpDelete(resourceKey),
			clientv3.OpDelete(resourceKey+"/", clientv3.WithPrefix()),
		)
	} else {
		operations = append(operations, clientv3.OpDelete(s.runtimeMethodKey(resourceID, *methodID)))
	}
	operations = append(operations, clientv3.OpPut(s.publishedRevisionKey(), strconv.FormatInt(generation, 10)))
	response, err := s.kv.Txn(ctx).If(comparisons...).Then(operations...).Commit()
	if err != nil {
		return fmt.Errorf("delete legacy runtime configuration: %w", err)
	}
	if !response.Succeeded {
		return fmt.Errorf("%w: legacy runtime delete", ErrRouteBindingPublishConflict)
	}
	return nil
}

func (s *RouteBindingStore) rejectLegacyAdminRouteBindingWrites(ctx context.Context, candidates []runtimeRouteBinding) error {
	publishedBindings, err := s.listEntries(ctx, s.bindingPrefix(false))
	if err != nil {
		return err
	}
	methodOwners := make(map[[2]int]string, len(publishedBindings))
	for _, binding := range publishedBindings {
		methodOwners[[2]int{binding.record.ResourceID, binding.record.MethodID}] = binding.record.Object.Metadata.Name
	}
	for _, candidate := range candidates {
		resourceID, methodID, ok, err := runtimeMethodIDs(s.runtimeResourcePrefix(), candidate.key)
		if err != nil {
			return err
		}
		if !ok {
			continue
		}
		if name, exists := methodOwners[[2]int{resourceID, methodID}]; exists {
			return fmt.Errorf("%w: runtime method %q is owned by AdminRouteBinding %q", ErrRouteBindingRuntimeConflict, candidate.key, name)
		}
	}
	return nil
}

func legacyResourcePublishRoutes(
	operations []clientv3.Op,
	runtimeResponse *clientv3.GetResponse,
	prefix string,
	resourcePaths map[int]string,
) ([]runtimeRouteBinding, error) {
	affectedResourceIDs, stagedMethodValues, routes, err := collectLegacyResourcePublishInputs(operations, prefix, resourcePaths)
	if err != nil {
		return nil, err
	}
	methodValues, err := collectLegacyResourcePublishMethodValues(runtimeResponse, prefix, affectedResourceIDs, stagedMethodValues)
	if err != nil {
		return nil, err
	}
	methodRoutes, err := buildLegacyResourcePublishMethodRoutes(methodValues, prefix, resourcePaths)
	if err != nil {
		return nil, err
	}
	return append(routes, methodRoutes...), nil
}

func collectLegacyResourcePublishInputs(
	operations []clientv3.Op,
	prefix string,
	resourcePaths map[int]string,
) (map[int]struct{}, map[string][]byte, []runtimeRouteBinding, error) {
	affectedResourceIDs := make(map[int]struct{})
	stagedMethodValues := make(map[string][]byte)
	routes := make([]runtimeRouteBinding, 0)
	for _, operation := range operations {
		key := string(operation.KeyBytes())
		if resourceID, ok := runtimeResourceID(prefix, key); ok {
			var resource config.Resource
			if err := commonyaml.UnmarshalYML(operation.ValueBytes(), &resource); err != nil {
				return nil, nil, nil, fmt.Errorf("decode staged legacy resource %q: %w", key, err)
			}
			collectLegacyResourceIDs(resource, resourceID, affectedResourceIDs)
			collectRuntimeResourceRoutes(resource, "", resourceID, resourcePaths, &routes)
			continue
		}
		resourceID, _, ok, err := runtimeMethodIDs(prefix, key)
		if err != nil {
			return nil, nil, nil, err
		}
		if !ok {
			return nil, nil, nil, fmt.Errorf("invalid legacy runtime publish key %q", key)
		}
		affectedResourceIDs[resourceID] = struct{}{}
		stagedMethodValues[key] = operation.ValueBytes()
	}
	return affectedResourceIDs, stagedMethodValues, routes, nil
}

func collectLegacyResourcePublishMethodValues(
	runtimeResponse *clientv3.GetResponse,
	prefix string,
	affectedResourceIDs map[int]struct{},
	stagedMethodValues map[string][]byte,
) (map[string][]byte, error) {
	methodValues := make(map[string][]byte)
	for _, item := range runtimeResponse.Kvs {
		resourceID, _, ok, err := runtimeMethodIDs(prefix, string(item.Key))
		if err != nil {
			return nil, err
		}
		if _, affected := affectedResourceIDs[resourceID]; ok && affected {
			methodValues[string(item.Key)] = item.Value
		}
	}
	for key, value := range stagedMethodValues {
		methodValues[key] = value
	}
	return methodValues, nil
}

func buildLegacyResourcePublishMethodRoutes(
	methodValues map[string][]byte,
	prefix string,
	resourcePaths map[int]string,
) ([]runtimeRouteBinding, error) {
	methodKeys := make([]string, 0, len(methodValues))
	for key := range methodValues {
		methodKeys = append(methodKeys, key)
	}
	sort.Strings(methodKeys)
	routes := make([]runtimeRouteBinding, 0, len(methodKeys))
	for _, key := range methodKeys {
		resourceID, _, ok, err := runtimeMethodIDs(prefix, key)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("invalid legacy runtime method key %q", key)
		}
		var method config.Method
		if err := commonyaml.UnmarshalYML(methodValues[key], &method); err != nil {
			return nil, fmt.Errorf("decode staged legacy method %q: %w", key, err)
		}
		path := strings.TrimSpace(method.ResourcePath)
		if path == "" {
			path = resourcePaths[resourceID]
		}
		if path == "" {
			continue
		}
		routes = append(routes, runtimeRouteBinding{path: path, httpVerb: method.HTTPVerb, key: key})
	}
	return routes, nil
}

func collectLegacyResourceIDs(resource config.Resource, keyResourceID int, ids map[int]struct{}) {
	if keyResourceID > 0 {
		ids[keyResourceID] = struct{}{}
	}
	if resource.ID > 0 {
		ids[resource.ID] = struct{}{}
	}
	for _, nested := range resource.Resources {
		collectLegacyResourceIDs(nested, nested.ID, ids)
	}
}

func findLegacyRuntimeRouteConflict(candidates, existing []runtimeRouteBinding) error {
	seen := append([]runtimeRouteBinding(nil), existing...)
	for _, candidate := range candidates {
		for _, route := range seen {
			if route.key == candidate.key {
				continue
			}
			if runtimeRouteConflicts(candidate, route) {
				return fmt.Errorf("%w: %s %s already exists at %q", ErrRouteBindingRuntimeConflict, strings.TrimSpace(candidate.httpVerb), candidate.path, route.key)
			}
		}
		seen = append(seen, candidate)
	}
	return nil
}

func normalizePositiveLegacyID(value, field string) (string, error) {
	id, err := strconv.Atoi(strings.TrimSpace(value))
	if err != nil || id <= 0 {
		return "", fmt.Errorf("%s must be a positive integer", field)
	}
	return strconv.Itoa(id), nil
}

func nextLegacyResourceID() int {
	return nextLegacyID(getRootPath(ResourceID))
}

func nextLegacyMethodID() int {
	return nextLegacyID(getRootPath(MethodID))
}

func nextLegacyID(key string) int {
	client := adminconfig.Client.GetRawClient()
	if client == nil {
		return ErrID
	}
	ctx := adminconfig.Client.GetCtx()
	for attempt := 0; attempt < maxRouteBindingIDRetries; attempt++ {
		response, err := client.Get(ctx, key)
		if err != nil {
			return ErrID
		}
		current := int64(0)
		var comparison clientv3.Cmp
		if len(response.Kvs) == 0 {
			comparison = clientv3.Compare(clientv3.CreateRevision(key), "=", 0)
		} else {
			current, err = strconv.ParseInt(string(response.Kvs[0].Value), 10, 64)
			if err != nil || current < 0 || current >= int64(^uint(0)>>1) {
				return ErrID
			}
			comparison = clientv3.Compare(clientv3.ModRevision(key), "=", response.Kvs[0].ModRevision)
		}
		newValue := current + 1
		result, err := client.Txn(ctx).If(comparison).Then(clientv3.OpPut(key, strconv.FormatInt(newValue, 10))).Commit()
		if err != nil {
			return ErrID
		}
		if result.Succeeded {
			return int(newValue)
		}
	}
	return ErrID
}

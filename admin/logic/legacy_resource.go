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
	if resource == nil {
		return errors.New("legacy resource is nil")
	}
	methods := append([]config.Method(nil), resource.Methods...)
	storedResource := *resource
	storedResource.Methods = nil
	if created && storedResource.ID <= 0 {
		storedResource.ID = nextLegacyResourceID()
		if storedResource.ID <= 0 {
			return errors.New("could not allocate legacy resource id")
		}
	} else if storedResource.ID <= 0 {
		return errors.New("legacy resource id must be positive")
	}
	resource.ID = storedResource.ID
	data, err := commonyaml.MarshalYML(storedResource)
	if err != nil {
		return fmt.Errorf("encode legacy resource: %w", err)
	}
	key := legacyResourceKey(strconv.Itoa(storedResource.ID), unpublished)
	if unpublished {
		if created {
			err = adminconfig.Client.Create(key, string(data))
		} else {
			err = adminconfig.Client.Update(key, string(data))
		}
		if err != nil {
			return fmt.Errorf("save legacy resource: %w", err)
		}
		if created {
			for index := range methods {
				methods[index].ResourcePath = storedResource.Path
				if err := BizSetResourceMethod(strconv.Itoa(storedResource.ID), &methods[index], true, unpublished); err != nil {
					return err
				}
				resource.Methods[index] = methods[index]
			}
		}
		return nil
	}

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
	return deleteLegacyResource(adminconfig.Client.GetCtx(), client, key)
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

// BizUpdateResourceAndMethods updates a legacy Resource and its Methods in one
// etcd transaction. This keeps Method.ResourcePath aligned with the Resource
// when a legacy client changes the route path.
func BizUpdateResourceAndMethods(resource *config.Resource, methods []config.Method, unpublished bool) error {
	if resource == nil {
		return errors.New("legacy resource is nil")
	}
	resourceID, err := normalizePositiveLegacyID(strconv.Itoa(resource.ID), ResourceID)
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

	storedResource := *resource
	storedResource.Methods = nil
	resourceData, err := commonyaml.MarshalYML(storedResource)
	if err != nil {
		return fmt.Errorf("encode legacy resource: %w", err)
	}
	operations := []clientv3.Op{clientv3.OpPut(
		legacyResourceKey(resourceID, unpublished), string(resourceData),
	)}
	for index := range methods {
		if methods[index].ID <= 0 {
			return fmt.Errorf("legacy method id must be positive, got %d", methods[index].ID)
		}
		methods[index].ResourcePath = resource.Path
		methodData, err := commonyaml.MarshalYML(methods[index])
		if err != nil {
			return fmt.Errorf("encode legacy method %d: %w", methods[index].ID, err)
		}
		operations = append(operations, clientv3.OpPut(
			legacyMethodKey(resourceID, strconv.Itoa(methods[index].ID), unpublished), string(methodData),
		))
	}
	if unpublished {
		if _, err := client.Txn(adminconfig.Client.GetCtx()).Then(operations...).Commit(); err != nil {
			return fmt.Errorf("save legacy resource and methods atomically: %w", err)
		}
		return nil
	}
	if err := publishLegacyResourceOperations(operations); err != nil {
		return fmt.Errorf("save legacy resource and methods atomically: %w", err)
	}
	return nil
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
	data, err := commonyaml.MarshalYML(method)
	if err != nil {
		return fmt.Errorf("encode legacy method: %w", err)
	}
	key := legacyMethodKey(resourceID, strconv.Itoa(method.ID), unpublished)
	if unpublished {
		if created {
			err = adminconfig.Client.Create(key, string(data))
		} else {
			err = adminconfig.Client.Update(key, string(data))
		}
		if err != nil {
			return fmt.Errorf("save legacy method: %w", err)
		}
		return nil
	}
	comparisons := make([]clientv3.Cmp, 0, 1)
	if created {
		comparisons = append(comparisons, clientv3.Compare(clientv3.CreateRevision(key), "=", 0))
	}
	if err := publishLegacyResourceOperations([]clientv3.Op{clientv3.OpPut(key, string(data))}, comparisons...); err != nil {
		return fmt.Errorf("save legacy method: %w", err)
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
	return adminconfig.Client.Delete(legacyMethodKey(resourceID, methodID, unpublished))
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
	client, err := legacyResourcePublishClient()
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

func legacyResourcePublishClient() (*clientv3.Client, error) {
	if adminconfig.Client == nil {
		return nil, errors.New("admin etcd client is unavailable")
	}
	client := adminconfig.Client.GetRawClient()
	if client == nil {
		return nil, errors.New("admin etcd client is unavailable")
	}
	return client, nil
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
	client, err := legacyResourcePublishClient()
	if err != nil {
		return err
	}
	if adminconfig.Bootstrap == nil {
		return errors.New("admin bootstrap is unavailable")
	}
	root := strings.TrimRight(strings.TrimSpace(adminconfig.Bootstrap.GetPath()), "/")
	if root == "" {
		return errors.New("admin etcd root path is empty")
	}
	store := &RouteBindingStore{kv: client, root: root}
	return store.publishLegacyResourceOperations(adminconfig.Client.GetCtx(), operations, comparisons)
}

func (s *RouteBindingStore) publishLegacyResourceOperations(ctx context.Context, operations []clientv3.Op, comparisons []clientv3.Cmp) error {
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
	operations = append(append([]clientv3.Op(nil), operations...), clientv3.OpPut(s.publishedRevisionKey(), strconv.FormatInt(generation, 10)))
	response, err := s.kv.Txn(ctx).If(comparisons...).Then(operations...).Commit()
	if err != nil {
		return fmt.Errorf("publish legacy resources atomically: %w", err)
	}
	if !response.Succeeded {
		return fmt.Errorf("%w: legacy resource publish", ErrRouteBindingPublishConflict)
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
	affectedResourceIDs := make(map[int]struct{})
	stagedMethodValues := make(map[string][]byte)
	routes := make([]runtimeRouteBinding, 0)
	for _, operation := range operations {
		key := string(operation.KeyBytes())
		if resourceID, ok := runtimeResourceID(prefix, key); ok {
			var resource config.Resource
			if err := commonyaml.UnmarshalYML(operation.ValueBytes(), &resource); err != nil {
				return nil, fmt.Errorf("decode staged legacy resource %q: %w", key, err)
			}
			collectLegacyResourceIDs(resource, resourceID, affectedResourceIDs)
			collectRuntimeResourceRoutes(resource, "", resourceID, resourcePaths, &routes)
			continue
		}
		resourceID, _, ok, err := runtimeMethodIDs(prefix, key)
		if err != nil {
			return nil, err
		}
		if !ok {
			return nil, fmt.Errorf("invalid legacy runtime publish key %q", key)
		}
		affectedResourceIDs[resourceID] = struct{}{}
		stagedMethodValues[key] = operation.ValueBytes()
	}

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
	methodKeys := make([]string, 0, len(methodValues))
	for key := range methodValues {
		methodKeys = append(methodKeys, key)
	}
	sort.Strings(methodKeys)
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

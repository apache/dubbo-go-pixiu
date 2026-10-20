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
	"errors"
	"fmt"
	"regexp"
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

func BizDeleteResourceInfo(id string, unpublished bool) error {
	_, err := adminconfig.Client.GetRawClient().Delete(adminconfig.Client.GetCtx(), legacyResourceKey(id, unpublished), clientv3.WithPrefix())
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

func BizDeleteMethodInfo(resourceID, methodID string, unpublished bool) error {
	return adminconfig.Client.Delete(legacyMethodKey(resourceID, methodID, unpublished))
}

// BizBatchReleaseResources copies the legacy staged runtime tree to the
// published tree. The old endpoints historically used a snapshot copy rather
// than route bindings; keep this compatibility path isolated from the new
// per-route atomic publisher.
func BizBatchReleaseResources() error {
	client := adminconfig.Client.GetRawClient()
	if client == nil {
		return errors.New("admin etcd client is unavailable")
	}
	fromPrefix := getUnpublishedRootPath(legacyResources)
	response, err := client.Get(adminconfig.Client.GetCtx(), fromPrefix, clientv3.WithPrefix())
	if err != nil {
		return fmt.Errorf("read staged legacy resources: %w", err)
	}
	for _, item := range response.Kvs {
		relative := strings.TrimPrefix(string(item.Key), fromPrefix)
		if relative == "" || strings.HasPrefix(relative, "/") {
			continue
		}
		if err := adminconfig.Client.Put(getRootPath(legacyResources)+"/"+relative, string(item.Value)); err != nil {
			return fmt.Errorf("publish legacy resource %q: %w", relative, err)
		}
	}
	return nil
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

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
	"regexp"
	"strconv"
	"strings"
)

import (
	fc "github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"

	gxetcd "github.com/dubbogo/gost/database/kv/etcd/v3"

	perrors "github.com/pkg/errors"

	clientv3 "go.etcd.io/etcd/client/v3"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	Base        = "base"
	Resources   = "resources"
	Method      = "method"
	ResourceID  = "resourceId"
	ClusterID   = "clusterId"
	Listener    = "listener"
	MethodID    = "methodId"
	PluginGroup = "pluginGroup"
	Plugin      = "plugin"
	Filter      = "filter"
	Ratelimit   = "ratelimit"
	Clusters    = "clusters"
	Listeners   = "listeners"
	Unpublished = "unpublished"

	ErrID = -1
)

// BizGetBaseInfo get base info
func BizGetBaseInfo() (*config.BaseInfo, error) {
	content, err := config.Client.Get(getRootPath(Base))
	if err != nil {
		logger.Errorf("BizGetBaseInfo err, %v\n", err)
		return nil, perrors.WithMessage(err, "BizGetBaseInfo error")
	}
	data := &config.BaseInfo{}
	_ = yaml.UnmarshalYML([]byte(content), data)

	return data, nil
}

// BizSetBaseInfo create or modify base info
func BizSetBaseInfo(info *config.BaseInfo, created bool) error {
	// validate the api config

	data, _ := yaml.MarshalYML(info)

	if created {
		setErr := config.Client.Put(getRootPath(Base), string(data))
		if setErr != nil {
			logger.Warnf("BizSetBaseInfo create error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetBaseInfo error")
		}
	} else {
		setErr := config.Client.Update(getRootPath(Base), string(data))
		if setErr != nil {
			logger.Warnf("BizSetBaseInfo update error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetBaseInfo error")
		}
	}

	return nil
}

// BizGetResourceList get resource list
func BizGetResourceList(unpublished bool) ([]fc.Resource, error) {
	var kList, vList []string
	var err error
	if unpublished {
		kList, vList, err = config.Client.GetChildrenKVList(getUnpublishedRootPath(Resources))
	} else {
		kList, vList, err = config.Client.GetChildrenKVList(getRootPath(Resources))
	}

	if err != nil {
		logger.Errorf("BizGetResourceList err, %v\n", err)
		return nil, perrors.WithMessage(err, "BizGetResourceList error")
	}
	var ret []fc.Resource

	for i, k := range kList {
		// only handle resource, filter method
		re := getCheckResourceRegexp()
		if m := re.Match([]byte(k)); !m {
			continue
		}
		v := vList[i]
		res := &fc.Resource{}
		err := yaml.UnmarshalYML([]byte(v), res)
		if err != nil {
			logger.Errorf("UnmarshalYML err, %v\n", err)
		}
		ret = append(ret, *res)
	}

	return ret, nil
}

// BizGetResourceDetail get resource detail
func BizGetResourceDetail(id string, unpublished bool) (string, error) {
	key := getResourceKey(id, unpublished)
	detail, err := config.Client.Get(key)
	if err != nil {
		logger.Errorf("BizGetResourceDetail err, %v\n", err)
		return "", perrors.WithMessage(err, "BizGetResourceDetail error")
	}
	return detail, nil
}

// BizSetResourceInfo create resource
func BizSetResourceInfo(res *fc.Resource, created, unpublished bool) error {

	if created {
		// backups method
		methods := res.Methods
		res.Methods = nil

		res.ID = getResourceId()
		if res.ID == ErrID {
			logger.Warnf("can't get id from etcd")
			return perrors.New("BizSetResourceInfo error can't get id from etcd")
		}
		data, _ := yaml.MarshalYML(res)

		setErr := config.Client.Create(getResourceKey(strconv.Itoa(res.ID), unpublished), string(data))
		if setErr != nil {
			logger.Warnf("Create etcd error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetResourceInfo error")
		}

		for i := range methods {
			methods[i].ResourcePath = res.Path
		}
		// create methods
		_ = BizBatchCreateResourceMethod(strconv.Itoa(res.ID), methods, unpublished)
	} else {
		key := getResourceKey(strconv.Itoa(res.ID), unpublished)
		data, _ := yaml.MarshalYML(res)

		// should set method in this situation
		res.Methods = nil
		setErr := config.Client.Update(key, string(data))
		if setErr != nil {
			logger.Warnf("update etcd error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetResourceInfo error")
		}
	}
	return nil
}

// BizDeleteResourceInfo delete resource
func BizDeleteResourceInfo(id string, unpublished bool) error {
	key := getResourceKey(id, unpublished)
	// delete all key with prefix to delete method key
	_, err := config.Client.GetRawClient().Delete(config.Client.GetCtx(), key, clientv3.WithPrefix())
	if err != nil {
		logger.Warnf("BizDeleteResourceInfo, %v\n", err)
		return perrors.WithMessage(err, "BizDeleteResourceInfo error")
	}
	return nil
}

// BizGetMethodList get method list
func BizGetMethodList(resourceId string, unpublished bool) ([]fc.Method, error) {
	key := getResourceMethodPrefixKey(resourceId, unpublished)

	_, vList, err := config.Client.GetChildrenKVList(key)
	if err != nil {
		logger.Errorf("BizGetMethodList err, %v\n", err)
		return nil, perrors.WithMessage(err, "BizGetMethodList error")
	}

	var ret []fc.Method
	for _, v := range vList {
		res := &fc.Method{}
		err := yaml.UnmarshalYML([]byte(v), res)
		if err != nil {
			logger.Errorf("UnmarshalYML err, %v\n", err)
		}
		ret = append(ret, *res)
	}

	return ret, nil
}

// BizGetMethodDetail get method detail
func BizGetMethodDetail(resourceId string, methodId string, unpublished bool) (string, error) {
	key := getMethodKey(resourceId, methodId, unpublished)
	detail, err := config.Client.Get(key)
	if err != nil {
		logger.Errorf("BizGetMethodDetail err, %v\n", err)
		return "", perrors.WithMessage(err, "BizGetMethodDetail error")
	}
	return detail, nil
}

// BizBatchCreateResourceMethod batch create method below one resource
func BizBatchCreateResourceMethod(resourceId string, methods []fc.Method, unpublished bool) error {

	if len(methods) == 0 {
		return nil
	}

	var kList, vList []string

	for _, method := range methods {
		method.ID = getMethodId()
		if method.ID == ErrID {
			logger.Warnf("can't get id from etcd")
			continue
		}
		kList = append(kList, getMethodKey(resourceId, strconv.Itoa(method.ID), unpublished))
		data, _ := yaml.MarshalYML(method)
		vList = append(vList, string(data))
	}

	err := config.Client.BatchCreate(kList, vList)
	if err != nil {
		logger.Warnf("update etcd error, %v\n", err)
		return perrors.WithMessage(err, "BizBatchCreateResourceMethod error")
	}
	return nil
}

// BizSetResourceMethod create or update method below specific path
func BizSetResourceMethod(resourceId string, method *fc.Method, created, unpublished bool) error {

	if created {
		method.ID = getMethodId()
		key := getMethodKey(resourceId, strconv.Itoa(method.ID), unpublished)

		if method.ID == ErrID {
			logger.Warnf("can't get id from etcd")
			return perrors.New("BizSetResourceMethod error can't get id from etcd")
		}
		data, _ := yaml.MarshalYML(method)

		err := config.Client.Create(key, string(data))
		if err != nil {
			logger.Warnf("BizSetResourceMethod etcd error, %v\n", err)
			return perrors.WithMessage(err, "BizSetResourceMethod error")
		}
	} else {
		data, _ := yaml.MarshalYML(method)
		key := getMethodKey(resourceId, strconv.Itoa(method.ID), unpublished)
		err := config.Client.Update(key, string(data))
		if err != nil {
			logger.Warnf("BizSetResourceMethod etcd error, %v\n", err)
			return perrors.WithMessage(err, "BizSetResourceMethod error")
		}
	}

	return nil
}

// BizDeleteMethodInfo delete method
func BizDeleteMethodInfo(resourceId string, methodId string, unpublished bool) error {
	key := getMethodKey(resourceId, methodId, unpublished)
	err := config.Client.Delete(key)
	if err != nil {
		logger.Warnf("BizDeleteMethodInfo, %v\n", err)
		return perrors.WithMessage(err, "BizDeleteMethodInfo error")
	}
	return nil
}

// BRGetResourceList GetResourceList
func BRGetResourceList(unpublished bool) ([]string, []string, error) {
	if unpublished {
		return config.Client.GetChildrenKVList(getUnpublishedRootPath(Resources))
	} else {
		return config.Client.GetChildrenKVList(getRootPath(Resources))
	}
}

// BRGetMethodList GetMethodList
func BRGetMethodList(resourceId string, unpublished bool) ([]string, []string, error) {
	key := getResourceMethodPrefixKey(resourceId, unpublished)
	return config.Client.GetChildrenKVList(key)
}

// BRGetPluginGroupList GetPluginGroupList
func BRGetPluginGroupList(unpublished bool) ([]string, []string, error) {
	key := getPluginGroupPrefixKey(unpublished)
	return config.Client.GetChildrenKVList(key)
}

// BizGetClusters get clusters
func BizGetClusters() ([]fc.Cluster, error) {
	var (
		kList, vList []string
		err          error
	)
	kList, vList, err = config.Client.GetChildrenKVList(getRootPath(Clusters))
	if err != nil {
		logger.Debugf("get clusters error from etcd, %+v, %+v, %s", kList, vList, err)
		return nil, perrors.WithMessage(err, "get clusters error")
	}

	var ret []fc.Cluster

	for i, k := range kList {
		// only handle resource, filter method
		re := getCheckClusterRegexp()
		if m := re.Match([]byte(k)); !m {
			continue
		}
		v := vList[i]
		res := &fc.Cluster{}
		err := yaml.UnmarshalYML([]byte(v), res)
		if err != nil {
			logger.Errorf("UnmarshalYML err, %v\n", err)
		}
		ret = append(ret, *res)
	}

	return ret, nil
}

// BizCreateCluster create cluster
func BizCreateCluster(res *fc.Cluster) error {
	res.ID = getClusterId()
	if res.ID == ErrID {
		logger.Warnf("can't get id from etcd")
		return perrors.New("BizSetCluster error can't get id from etcd")
	}
	data, _ := yaml.MarshalYML(res)
	setErr := config.Client.Create(getClusterKey(strconv.Itoa(res.ID)), string(data))

	if setErr != nil {
		logger.Warnf("Create etcd error, %v\n", setErr)
		return perrors.WithMessage(setErr, "BizSetCluster error")
	}

	return nil
}

// BizUpdateCluster create cluster
func BizUpdateCluster(res *fc.Cluster) error {
	if res.ID <= 0 {
		logger.Warnf("invalid cluster id, %d", res.ID)
		return perrors.New("invalid cluster id")
	}
	data, _ := yaml.MarshalYML(res)
	setErr := config.Client.Update(getClusterKey(strconv.Itoa(res.ID)), string(data))

	if setErr != nil {
		logger.Warnf("Update etcd error, %v\n", setErr)
		return perrors.WithMessage(setErr, "BizSetCluster error")
	}

	return nil
}

// BizDeleteCluster delete cluster
func BizDeleteCluster(id string) error {
	key := getClusterKey(id)
	// delete all key with prefix to delete method key
	_, err := config.Client.GetRawClient().Delete(config.Client.GetCtx(), key, clientv3.WithPrefix())
	if err != nil {
		logger.Warnf("BizDeleteCluster, %v\n", err)
		return perrors.WithMessage(err, "BizDeleteCluster error")
	}
	return nil
}

// BizGetCluster get cluster
func BizGetCluster(id string) (string, error) {
	key := getClusterKey(id)
	detail, err := config.Client.Get(key)
	if err != nil {
		logger.Errorf("BizGetClusterDetail error, %v\n", err)
		return "", perrors.WithMessage(err, "BizGetClusterDetail error")
	}
	return detail, nil
}

// BizGetListeners get Listeners
func BizGetListeners() ([]fc.Listener, error) {
	kList, vList, err := config.Client.GetChildrenKVList(getRootPath(Listeners))
	if err != nil {
		logger.Debugf("get listeners error from etcd, %+v, %+v, %s", kList, vList, err)
		return nil, perrors.WithMessage(err, "get listeners error")
	}

	var ret []fc.Listener

	for i, k := range kList {
		// only handle resource, filter method
		re := getCheckListenerRegexp()
		if m := re.Match([]byte(k)); !m {
			continue
		}
		v := vList[i]
		res := &fc.Listener{}
		err := yaml.UnmarshalYML([]byte(v), res)
		if err != nil {
			logger.Errorf("UnmarshalYML err, %v\n", err)
		}
		ret = append(ret, *res)
	}

	return ret, nil
}

// BizCreateListener create Listener
func BizCreateListener(res *fc.Listener) error {
	if strings.TrimSpace(res.Name) == "" {
		logger.Warnf("invalid listener name, %s", res.Name)
		return perrors.New("invalid listener id")
	}
	data, _ := yaml.MarshalYML(res)
	setErr := config.Client.Create(getListenerKey(res.Name), string(data))

	if setErr != nil {
		logger.Warnf("Create etcd error, %v\n", setErr)
		return perrors.WithMessage(setErr, "BizSetCluster error")
	}

	return nil
}

// BizUpdateListener create Listener
func BizUpdateListener(res *fc.Listener) error {
	if strings.TrimSpace(res.Name) == "" {
		logger.Warnf("invalid listener name, %s", res.Name)
		return perrors.New("invalid listener name")
	}
	data, _ := yaml.MarshalYML(res)
	setErr := config.Client.Update(getListenerKey(res.Name), string(data))

	if setErr != nil {
		logger.Warnf("Update etcd error, %v\n", setErr)
		return perrors.WithMessage(setErr, "BizSetListener error")
	}

	return nil
}

// BizDeleteListener delete Listener
func BizDeleteListener(name string) error {
	key := getListenerKey(name)
	// delete all key with prefix to delete listener key
	_, err := config.Client.GetRawClient().Delete(config.Client.GetCtx(), key, clientv3.WithPrefix())
	if err != nil {
		logger.Warnf("BizDeleteListener, %v\n", err)
		return perrors.WithMessage(err, "BizDeleteListener error")
	}
	return nil
}

// BizGetListener get Listener
func BizGetListener(name string) (string, error) {
	key := getListenerKey(name)
	detail, err := config.Client.Get(key)
	if err != nil {
		logger.Errorf("BizGetListenerDetail error, %v\n", err)
		return "", perrors.WithMessage(err, "BizGetListenerDetail error")
	}
	return detail, nil
}

// BRUpdate
func BRUpdate(key, value string) error {
	return config.Client.Update(key, value)
}

func BRCreate(key, value, configType string) error {
	if strings.EqualFold(configType, Resources) {
		return config.Client.Create(getResourceKey(key, false), value)
	} else if strings.EqualFold(configType, Method) {

	} else if strings.EqualFold(configType, PluginGroup) {
		return config.Client.Create(getPluginGroupKey(key, false), value)
	} else {
		return config.Client.Create(getPluginRatelimitKey(false), value)
	}
	return errors.New("")
}

func getResourceKey(path string, unpublished bool) string {
	if unpublished {
		return getUnpublishedRootPath(Resources) + "/" + path
	}
	return getRootPath(Resources) + "/" + path
}

func getClusterKey(path string) string {
	return getRootPath(Clusters) + "/" + path
}

func getListenerKey(path string) string {
	return getRootPath(Listeners) + "/" + path
}

func getPluginRatelimitKey(unpublished bool) string {
	return getFilterPrefixKey(unpublished) + "/" + Ratelimit
}

func getPluginGroupKey(name string, unpublished bool) string {
	return getPluginGroupPrefixKey(unpublished) + "/" + name
}

func getPluginGroupPrefixKey(unpublished bool) string {
	if unpublished {
		return getUnpublishedRootPath(PluginGroup)
	}
	return getRootPath(PluginGroup)
}

func getFilterPrefixKey(unpublished bool) string {
	if unpublished {
		return getUnpublishedRootPath(Filter)
	}
	return getRootPath(Filter)
}

func getResourceMethodPrefixKey(path string, unpublished bool) string {
	return getResourceKey(path, unpublished) + "/" + Method
}

func getMethodKey(path string, method string, unpublished bool) string {
	return getResourceMethodPrefixKey(path, unpublished) + "/" + method
}

// create method, No need to judge whether to publish or not
func getResourceId() int {
	return loopGetId(getRootPath(ResourceID))
}

func getClusterId() int {
	return loopGetId(getRootPath(ClusterID))
}

func getMethodId() int {
	return loopGetId(getRootPath(MethodID))
}

func loopGetId(k string) int {

	for true {

		rawClient := config.Client.GetRawClient()
		if rawClient == nil {
			logger.Error("GetId etcd client is null")
			return ErrID
		}

		resp, err := rawClient.Get(config.Client.GetCtx(), k)
		if err != nil {
			return ErrID
		}

		var val string
		var rev int64
		if len(resp.Kvs) != 0 {
			val = string(resp.Kvs[0].Value)
			rev = resp.Kvs[0].ModRevision
		} else {
			val = "0"
			rev = 0
		}
		id, err := strconv.Atoi(val)
		if err != nil {
			logger.Error("GetId Atoi error, %v\n", err)
			return ErrID
		}

		id += 1

		if rev == 0 {
			err = config.Client.Create(k, strconv.Itoa(id))
		} else {
			err = config.Client.UpdateWithRev(k, strconv.Itoa(id), rev)
		}

		if err != nil {
			if !errors.Is(err, gxetcd.ErrCompareFail) {
				logger.Error("GetId UpdateWithRev error, %v\n", err)
				return ErrID
			}
			logger.Info("retry get id")
		} else {
			return id
		}
	}
	return ErrID
}

func getRootPath(key string) string {
	return config.Bootstrap.GetPath() + "/" + key
}

func getUnpublishedRootPath(key string) string {
	return config.Bootstrap.GetPath() + "/" + Unpublished + "/" + key
}

func getCheckResourceRegexp() *regexp.Regexp {
	return regexp.MustCompile(".+/resources/[^/]+/?$")
}

func getCheckClusterRegexp() *regexp.Regexp {
	return regexp.MustCompile(".+/clusters/[^/]+/?$")
}

func getCheckListenerRegexp() *regexp.Regexp {
	return regexp.MustCompile(".+/listeners/[^/]+/?$")
}

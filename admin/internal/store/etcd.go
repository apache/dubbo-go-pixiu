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

package store

import (
	"regexp"
	"strconv"
	"strings"
)

import (
	gxetcd "github.com/dubbogo/gost/database/kv/etcd/v3"

	"github.com/pkg/errors"

	clientv3 "go.etcd.io/etcd/client/v3"
)

import (
	"github.com/apache/dubbo-go-pixiu/admin/internal/model"
	"github.com/apache/dubbo-go-pixiu/admin/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	pkgconfig "github.com/apache/dubbo-go-pixiu/pkg/config"
	pkgmodel "github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	keyBase        = "base"
	keyResources   = "resources"
	keyMethod      = "method"
	keyResourceID  = "resourceId"
	keyClusterID   = "clusterId"
	keyMethodID    = "methodId"
	keyClusters    = "clusters"
	keyListeners   = "listeners"
	keyPluginGroup = "pluginGroup"
	keyUnpublished = "unpublished"

	errID = -1
)

type Etcd struct {
	client   *gxetcd.Client
	basePath string
}

func NewEtcd(client *gxetcd.Client, basePath string) *Etcd {
	return &Etcd{client: client, basePath: basePath}
}

func (e *Etcd) Close() {
	if e.client != nil {
		e.client.Close()
	}
}

func (e *Etcd) Client() *gxetcd.Client { return e.client }
func (e *Etcd) BasePath() string       { return e.basePath }

func (e *Etcd) GetBaseInfo() (*config.BaseInfo, error) {
	content, err := e.client.Get(e.path(keyBase))
	if err != nil {
		return nil, errors.Wrap(err, "failed to get base info")
	}
	info := &config.BaseInfo{}
	if err := yaml.UnmarshalYML([]byte(content), info); err != nil {
		return nil, errors.Wrap(err, "failed to parse base info")
	}
	return info, nil
}

func (e *Etcd) SetBaseInfo(info *config.BaseInfo, create bool) error {
	data, err := yaml.MarshalYML(info)
	if err != nil {
		return errors.Wrap(err, "failed to marshal base info")
	}
	if create {
		err = e.client.Put(e.path(keyBase), string(data))
	} else {
		err = e.client.Update(e.path(keyBase), string(data))
	}
	if err != nil {
		return errors.Wrap(err, "failed to save base info")
	}
	return nil
}

func (e *Etcd) ListClusters() ([]pkgmodel.ClusterConfig, error) {
	kList, vList, err := e.client.GetChildrenKVList(e.path(keyClusters))
	if err != nil {
		return []pkgmodel.ClusterConfig{}, nil
	}
	re := regexp.MustCompile(`.+/clusters/[^/]+/?$`)
	var clusters []pkgmodel.ClusterConfig
	for i, k := range kList {
		if !re.Match([]byte(k)) {
			continue
		}
		var cluster pkgmodel.ClusterConfig
		if err := yaml.UnmarshalYML([]byte(vList[i]), &cluster); err != nil {
			continue
		}
		clusters = append(clusters, cluster)
	}
	return clusters, nil
}

func (e *Etcd) GetCluster(name string) (string, error) {
	return e.client.Get(e.clusterPath(name))
}

func (e *Etcd) CreateCluster(cluster *pkgmodel.ClusterConfig) error {
	if strings.TrimSpace(cluster.Name) == "" {
		return errors.New("invalid cluster name")
	}
	data, _ := yaml.MarshalYML(cluster)
	return e.client.Create(e.clusterPath(cluster.Name), string(data))
}

func (e *Etcd) UpdateCluster(cluster *pkgmodel.ClusterConfig) error {
	if strings.TrimSpace(cluster.Name) == "" {
		return errors.New("invalid cluster name")
	}
	data, _ := yaml.MarshalYML(cluster)
	return e.client.Update(e.clusterPath(cluster.Name), string(data))
}

func (e *Etcd) DeleteCluster(name string) error {
	_, err := e.client.GetRawClient().Delete(e.client.GetCtx(), e.clusterPath(name), clientv3.WithPrefix())
	return err
}

func (e *Etcd) ListListeners() ([]*pkgmodel.Listener, error) {
	kList, vList, err := e.client.GetChildrenKVList(e.path(keyListeners))
	if err != nil {
		return []*pkgmodel.Listener{}, nil
	}
	re := regexp.MustCompile(`.+/listeners/[^/]+/?$`)
	var listeners []*pkgmodel.Listener
	for i, k := range kList {
		if !re.Match([]byte(k)) {
			continue
		}
		var listener pkgmodel.Listener
		if err := yaml.UnmarshalYML([]byte(vList[i]), &listener); err != nil {
			continue
		}
		listeners = append(listeners, &listener)
	}
	return listeners, nil
}

func (e *Etcd) GetListener(name string) (string, error) {
	return e.client.Get(e.listenerPath(name))
}

func (e *Etcd) CreateListener(listener *pkgmodel.Listener) error {
	if strings.TrimSpace(listener.Name) == "" {
		return errors.New("invalid listener name")
	}
	data, _ := yaml.MarshalYML(listener)
	return e.client.Create(e.listenerPath(listener.Name), string(data))
}

func (e *Etcd) UpdateListener(listener *pkgmodel.Listener) error {
	if strings.TrimSpace(listener.Name) == "" {
		return errors.New("invalid listener name")
	}
	data, _ := yaml.MarshalYML(listener)
	return e.client.Update(e.listenerPath(listener.Name), string(data))
}

func (e *Etcd) DeleteListener(name string) error {
	_, err := e.client.GetRawClient().Delete(e.client.GetCtx(), e.listenerPath(name), clientv3.WithPrefix())
	return err
}

func (e *Etcd) ListResources(unpublished bool) ([]pkgconfig.Resource, error) {
	path := e.resourcesPath(unpublished)
	kList, vList, err := e.client.GetChildrenKVList(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list resources")
	}
	re := regexp.MustCompile(`.+/resources/([^/]+)/?$`)
	var resources []pkgconfig.Resource
	for i, k := range kList {
		matches := re.FindStringSubmatch(k)
		if matches == nil {
			continue
		}
		var res pkgconfig.Resource
		if err := yaml.UnmarshalYML([]byte(vList[i]), &res); err != nil {
			continue
		}
		if res.ID == 0 && len(matches) > 1 {
			if id, err := strconv.Atoi(matches[1]); err == nil {
				res.ID = id
			}
		}
		resources = append(resources, res)
	}
	return resources, nil
}

func (e *Etcd) GetResource(id string, unpublished bool) (string, error) {
	return e.client.Get(e.resourcePath(id, unpublished))
}

func (e *Etcd) CreateResource(res *pkgconfig.Resource, unpublished bool) error {
	methods := res.Methods
	res.Methods = nil
	res.ID = e.nextID(keyResourceID)
	if res.ID == errID {
		return errors.New("failed to generate resource ID")
	}
	data, _ := yaml.MarshalYML(res)
	if err := e.client.Create(e.resourcePath(strconv.Itoa(res.ID), unpublished), string(data)); err != nil {
		return err
	}
	for i := range methods {
		methods[i].ResourcePath = res.Path
	}
	return e.batchCreateMethods(strconv.Itoa(res.ID), methods, unpublished)
}

func (e *Etcd) UpdateResource(res *pkgconfig.Resource, unpublished bool) error {
	res.Methods = nil
	data, _ := yaml.MarshalYML(res)
	return e.client.Update(e.resourcePath(strconv.Itoa(res.ID), unpublished), string(data))
}

func (e *Etcd) DeleteResource(id string, unpublished bool) error {
	_, err := e.client.GetRawClient().Delete(e.client.GetCtx(), e.resourcePath(id, unpublished), clientv3.WithPrefix())
	return err
}

func (e *Etcd) ListMethods(resourceID string, unpublished bool) ([]pkgconfig.Method, error) {
	path := e.methodsPath(resourceID, unpublished)
	_, vList, err := e.client.GetChildrenKVList(path)
	if err != nil {
		return nil, errors.Wrap(err, "failed to list methods")
	}
	var methods []pkgconfig.Method
	for _, v := range vList {
		var m pkgconfig.Method
		if err := yaml.UnmarshalYML([]byte(v), &m); err != nil {
			continue
		}
		methods = append(methods, m)
	}
	return methods, nil
}

func (e *Etcd) GetMethod(resourceID, methodID string, unpublished bool) (string, error) {
	return e.client.Get(e.methodPath(resourceID, methodID, unpublished))
}

func (e *Etcd) CreateMethod(resourceID string, method *pkgconfig.Method, unpublished bool) error {
	method.ID = e.nextID(keyMethodID)
	if method.ID == errID {
		return errors.New("failed to generate method ID")
	}
	data, _ := yaml.MarshalYML(method)
	return e.client.Create(e.methodPath(resourceID, strconv.Itoa(method.ID), unpublished), string(data))
}

func (e *Etcd) UpdateMethod(resourceID string, method *pkgconfig.Method, unpublished bool) error {
	data, _ := yaml.MarshalYML(method)
	return e.client.Update(e.methodPath(resourceID, strconv.Itoa(method.ID), unpublished), string(data))
}

func (e *Etcd) DeleteMethod(resourceID, methodID string, unpublished bool) error {
	return e.client.Delete(e.methodPath(resourceID, methodID, unpublished))
}

func (e *Etcd) batchCreateMethods(resourceID string, methods []pkgconfig.Method, unpublished bool) error {
	if len(methods) == 0 {
		return nil
	}
	var keys, values []string
	for _, m := range methods {
		m.ID = e.nextID(keyMethodID)
		if m.ID == errID {
			continue
		}
		keys = append(keys, e.methodPath(resourceID, strconv.Itoa(m.ID), unpublished))
		data, _ := yaml.MarshalYML(m)
		values = append(values, string(data))
	}
	return e.client.BatchCreate(keys, values)
}

func (e *Etcd) ListPluginGroups(unpublished bool) ([]model.PluginGroup, error) {
	path := e.pluginGroupsPath(unpublished)
	_, vList, err := e.client.GetChildrenKVList(path)
	if err != nil {
		return []model.PluginGroup{}, nil
	}
	var groups []model.PluginGroup
	for _, v := range vList {
		var g model.PluginGroup
		if err := yaml.UnmarshalYML([]byte(v), &g); err != nil {
			continue
		}
		groups = append(groups, g)
	}
	return groups, nil
}

func (e *Etcd) GetPluginGroup(name string, unpublished bool) (string, error) {
	return e.client.Get(e.pluginGroupPath(name, unpublished))
}

func (e *Etcd) CreatePluginGroup(group *model.PluginGroup, unpublished bool) error {
	if strings.TrimSpace(group.GroupName) == "" {
		return errors.New("invalid plugin group name")
	}
	data, _ := yaml.MarshalYML(group)
	return e.client.Create(e.pluginGroupPath(group.GroupName, unpublished), string(data))
}

func (e *Etcd) UpdatePluginGroup(group *model.PluginGroup, unpublished bool) error {
	if strings.TrimSpace(group.GroupName) == "" {
		return errors.New("invalid plugin group name")
	}
	data, _ := yaml.MarshalYML(group)
	return e.client.Update(e.pluginGroupPath(group.GroupName, unpublished), string(data))
}

func (e *Etcd) DeletePluginGroup(name string, unpublished bool) error {
	_, err := e.client.GetRawClient().Delete(e.client.GetCtx(), e.pluginGroupPath(name, unpublished), clientv3.WithPrefix())
	return err
}

func (e *Etcd) GetResourcesKV(unpublished bool) ([]string, []string, error) {
	return e.client.GetChildrenKVList(e.resourcesPath(unpublished))
}

func (e *Etcd) GetMethodsKV(resourceID string, unpublished bool) ([]string, []string, error) {
	return e.client.GetChildrenKVList(e.methodsPath(resourceID, unpublished))
}

func (e *Etcd) GetPluginGroupsKV(unpublished bool) ([]string, []string, error) {
	return e.client.GetChildrenKVList(e.pluginGroupsPath(unpublished))
}

func (e *Etcd) Update(key, value string) error { return e.client.Update(key, value) }
func (e *Etcd) Create(key, value string) error { return e.client.Create(key, value) }

func (e *Etcd) path(key string) string { return e.basePath + "/" + key }
func (e *Etcd) unpublishedPath(key string) string {
	return e.basePath + "/" + keyUnpublished + "/" + key
}
func (e *Etcd) clusterPath(id string) string    { return e.path(keyClusters) + "/" + id }
func (e *Etcd) listenerPath(name string) string { return e.path(keyListeners) + "/" + name }

func (e *Etcd) resourcesPath(unpublished bool) string {
	if unpublished {
		return e.unpublishedPath(keyResources)
	}
	return e.path(keyResources)
}

func (e *Etcd) resourcePath(id string, unpublished bool) string {
	return e.resourcesPath(unpublished) + "/" + id
}

func (e *Etcd) methodsPath(resourceID string, unpublished bool) string {
	return e.resourcePath(resourceID, unpublished) + "/" + keyMethod
}

func (e *Etcd) methodPath(resourceID, methodID string, unpublished bool) string {
	return e.methodsPath(resourceID, unpublished) + "/" + methodID
}

func (e *Etcd) pluginGroupsPath(unpublished bool) string {
	if unpublished {
		return e.unpublishedPath(keyPluginGroup)
	}
	return e.path(keyPluginGroup)
}

func (e *Etcd) pluginGroupPath(name string, unpublished bool) string {
	return e.pluginGroupsPath(unpublished) + "/" + name
}

func (e *Etcd) WatchConfig() (<-chan bool, error) {
	ch, err := e.client.WatchWithPrefix(e.basePath)
	if err != nil {
		return nil, errors.Wrap(err, "failed to watch config")
	}
	out := make(chan bool)
	go func() {
		for range ch {
			out <- true
		}
		close(out)
	}()
	return out, nil
}

func (e *Etcd) nextID(keyType string) int {
	key := e.path(keyType)
	for {
		rawClient := e.client.GetRawClient()
		if rawClient == nil {
			return errID
		}
		resp, err := rawClient.Get(e.client.GetCtx(), key)
		if err != nil {
			return errID
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
			return errID
		}
		id++
		if rev == 0 {
			err = e.client.Create(key, strconv.Itoa(id))
		} else {
			err = e.client.UpdateWithRev(key, strconv.Itoa(id), rev)
		}
		if err == nil {
			return id
		}
		if !errors.Is(err, gxetcd.ErrCompareFail) {
			return errID
		}
	}
}

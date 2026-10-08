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
	"bytes"
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"io"
	"net/http"
	"regexp"
	"strconv"
	"strings"
	"time"
)

import (
	gxetcd "github.com/dubbogo/gost/database/kv/etcd/v3"

	perrors "github.com/pkg/errors"

	clientv3 "go.etcd.io/etcd/client/v3"
)

import (
	adminconfig "github.com/apache/dubbo-go-pixiu/admin/config"
	"github.com/apache/dubbo-go-pixiu/pkg/common/yaml"
	"github.com/apache/dubbo-go-pixiu/pkg/config"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	Base        = "base"
	Method      = "method"
	ResourceID  = "resourceId"
	ClusterID   = "clusterId"
	Listener    = "listener"
	Plugin      = "plugin"
	Filter      = "filter"
	Ratelimit   = "ratelimit"
	OPA         = "opa"
	Clusters    = "clusters"
	Listeners   = "listeners"
	PluginGroup = "pluginGroup"
	Unpublished = "unpublished"

	ErrID = -1
)

// Use a shared client to enable HTTP keep-alive and prevent connection exhaustion.
// Timeout is enforced per request so it can honor late-loaded config.
// Set a large fallback timeout as a last-resort guard.
const opaHTTPClientFallbackTimeout = 30 * time.Second

var opaHTTPClient = &http.Client{
	Timeout: opaHTTPClientFallbackTimeout,
}

func getOPATimeout() time.Duration {
	if adminconfig.Bootstrap != nil && adminconfig.Bootstrap.OPA.RequestTimeout > 0 {
		return adminconfig.Bootstrap.OPA.RequestTimeout
	}
	return adminconfig.DefaultOPAPolicyTimeout
}

// BizGetBaseInfo get base info
func BizGetBaseInfo() (*adminconfig.BaseInfo, error) {
	content, err := adminconfig.Client.Get(getRootPath(Base))
	if err != nil {
		logger.Errorf("BizGetBaseInfo err, %v\n", err)
		return nil, perrors.WithMessage(err, "BizGetBaseInfo error")
	}
	data := &adminconfig.BaseInfo{}
	_ = yaml.UnmarshalYML([]byte(content), data)

	return data, nil
}

// BizSetBaseInfo create or modify base info
func BizSetBaseInfo(info *adminconfig.BaseInfo, created bool) error {
	// validate the api config

	data, _ := yaml.MarshalYML(info)

	if created {
		setErr := adminconfig.Client.Put(getRootPath(Base), string(data))
		if setErr != nil {
			logger.Warnf("BizSetBaseInfo create error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetBaseInfo error")
		}
	} else {
		setErr := adminconfig.Client.Update(getRootPath(Base), string(data))
		if setErr != nil {
			logger.Warnf("BizSetBaseInfo update error, %v\n", setErr)
			return perrors.WithMessage(setErr, "BizSetBaseInfo error")
		}
	}

	return nil
}

// BizGetClusters get clusters
func BizGetClusters() ([]config.Cluster, error) {
	var (
		kList, vList []string
		err          error
	)
	kList, vList, err = adminconfig.Client.GetChildrenKVList(getRootPath(Clusters))
	if err != nil {
		if errors.Is(err, gxetcd.ErrKVPairNotFound) {
			return []config.Cluster{}, nil
		}
		logger.Debugf("get clusters error from etcd, %+v, %+v, %s", kList, vList, err)
		return nil, perrors.WithMessage(err, "get clusters error")
	}

	return decodeClusters(kList, vList)
}

func decodeClusters(kList, vList []string) ([]config.Cluster, error) {
	if len(kList) != len(vList) {
		return nil, perrors.Errorf("cluster key/value length mismatch: %d != %d", len(kList), len(vList))
	}
	ret := make([]config.Cluster, 0, len(kList))
	for i, k := range kList {
		// only handle resource, filter method
		re := getCheckClusterRegexp()
		if m := re.Match([]byte(k)); !m {
			continue
		}
		v := vList[i]
		res := &config.Cluster{}
		if err := yaml.UnmarshalYML([]byte(v), res); err != nil {
			return nil, perrors.Wrapf(err, "decode cluster configuration at %q", k)
		}
		ret = append(ret, *res)
	}

	return ret, nil
}

// BizCreateCluster create cluster
func BizCreateCluster(res *config.Cluster) error {
	res.ID = getClusterId()
	if res.ID == ErrID {
		logger.Warnf("can't get id from etcd")
		return perrors.New("BizSetCluster error can't get id from etcd")
	}
	data, _ := yaml.MarshalYML(res)
	setErr := adminconfig.Client.Create(getClusterKey(strconv.Itoa(res.ID)), string(data))

	if setErr != nil {
		logger.Warnf("Create etcd error, %v\n", setErr)
		return perrors.WithMessage(setErr, "BizSetCluster error")
	}

	return nil
}

// BizUpdateCluster create cluster
func BizUpdateCluster(res *config.Cluster) error {
	if res.ID <= 0 {
		logger.Warnf("invalid cluster id, %d", res.ID)
		return perrors.New("invalid cluster id")
	}
	data, _ := yaml.MarshalYML(res)
	setErr := adminconfig.Client.Update(getClusterKey(strconv.Itoa(res.ID)), string(data))

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
	_, err := adminconfig.Client.GetRawClient().Delete(adminconfig.Client.GetCtx(), key, clientv3.WithPrefix())
	if err != nil {
		logger.Warnf("BizDeleteCluster, %v\n", err)
		return perrors.WithMessage(err, "BizDeleteCluster error")
	}
	return nil
}

// BizGetCluster get cluster
func BizGetCluster(id string) (string, error) {
	key := getClusterKey(id)
	detail, err := adminconfig.Client.Get(key)
	if err != nil {
		logger.Errorf("BizGetClusterDetail error, %v\n", err)
		return "", perrors.WithMessage(err, "BizGetClusterDetail error")
	}
	return detail, nil
}

// BizGetListeners get Listeners
func BizGetListeners() ([]config.Listener, error) {
	kList, vList, err := adminconfig.Client.GetChildrenKVList(getRootPath(Listeners))
	if err != nil {
		if errors.Is(err, gxetcd.ErrKVPairNotFound) {
			return []config.Listener{}, nil
		}
		logger.Debugf("get listeners error from etcd, %+v, %+v, %s", kList, vList, err)
		return nil, perrors.WithMessage(err, "get listeners error")
	}

	return decodeListeners(kList, vList)
}

func decodeListeners(kList, vList []string) ([]config.Listener, error) {
	if len(kList) != len(vList) {
		return nil, perrors.Errorf("listener key/value length mismatch: %d != %d", len(kList), len(vList))
	}
	ret := make([]config.Listener, 0, len(kList))
	for i, k := range kList {
		// only handle resource, filter method
		re := getCheckListenerRegexp()
		if m := re.Match([]byte(k)); !m {
			continue
		}
		v := vList[i]
		res := &config.Listener{}
		if err := yaml.UnmarshalYML([]byte(v), res); err != nil {
			return nil, perrors.Wrapf(err, "decode listener configuration at %q", k)
		}
		ret = append(ret, *res)
	}

	return ret, nil
}

// BizCreateListener create Listener
func BizCreateListener(res *config.Listener) error {
	if strings.TrimSpace(res.Name) == "" {
		logger.Warnf("invalid listener name, %s", res.Name)
		return perrors.New("invalid listener id")
	}
	data, _ := yaml.MarshalYML(res)
	setErr := adminconfig.Client.Create(getListenerKey(res.Name), string(data))

	if setErr != nil {
		logger.Warnf("Create etcd error, %v\n", setErr)
		return perrors.WithMessage(setErr, "BizSetCluster error")
	}

	return nil
}

// BizUpdateListener create Listener
func BizUpdateListener(res *config.Listener) error {
	if strings.TrimSpace(res.Name) == "" {
		logger.Warnf("invalid listener name, %s", res.Name)
		return perrors.New("invalid listener name")
	}
	data, _ := yaml.MarshalYML(res)
	setErr := adminconfig.Client.Update(getListenerKey(res.Name), string(data))

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
	_, err := adminconfig.Client.GetRawClient().Delete(adminconfig.Client.GetCtx(), key, clientv3.WithPrefix())
	if err != nil {
		logger.Warnf("BizDeleteListener, %v\n", err)
		return perrors.WithMessage(err, "BizDeleteListener error")
	}
	return nil
}

// BizGetListener get Listener
func BizGetListener(name string) (string, error) {
	key := getListenerKey(name)
	detail, err := adminconfig.Client.Get(key)
	if err != nil {
		logger.Errorf("BizGetListenerDetail error, %v\n", err)
		return "", perrors.WithMessage(err, "BizGetListenerDetail error")
	}
	return detail, nil
}

// BizPublishPluginGroupConfig copies the first staged PluginGroup config into
// the published namespace, creating it if it does not already exist.
func BizPublishPluginGroupConfig() error {
	if adminconfig.Client == nil {
		return perrors.New("admin etcd client is not initialized")
	}
	if adminconfig.Bootstrap == nil {
		return perrors.New("admin bootstrap is not initialized")
	}

	draftKeys, draftValues, err := adminconfig.Client.GetChildrenKVList(getPluginGroupPrefixKey(true))
	if err != nil {
		return perrors.WithMessage(err, "get unpublished PluginGroup config")
	}
	if len(draftKeys) == 0 || len(draftKeys) != len(draftValues) {
		return perrors.New("unpublished PluginGroup config is empty or invalid")
	}

	publishedKeys, publishedValues, err := adminconfig.Client.GetChildrenKVList(getPluginGroupPrefixKey(false))
	if err != nil && !errors.Is(err, gxetcd.ErrKVPairNotFound) {
		return perrors.WithMessage(err, "get published PluginGroup config")
	}
	if len(publishedKeys) == 0 {
		name := strings.TrimPrefix(draftKeys[0], getPluginGroupPrefixKey(true)+"/")
		if name == draftKeys[0] || name == "" || strings.Contains(name, "/") {
			return perrors.Errorf("invalid unpublished PluginGroup key %q", draftKeys[0])
		}
		key := getPluginGroupPrefixKey(false) + "/" + name
		return adminconfig.Client.Create(key, draftValues[0])
	}
	if len(publishedKeys) != len(publishedValues) {
		return perrors.New("published PluginGroup config is invalid")
	}
	if strings.EqualFold(draftValues[0], publishedValues[0]) {
		return nil
	}
	return adminconfig.Client.Update(publishedKeys[0], draftValues[0])
}

// BizGetOPAPolicy fetches the policy raw text. Returns empty string if not found.
func BizGetOPAPolicy(serverURL, policyID, bearerToken string) (string, error) {
	url, err := buildOPAPolicyURL(serverURL, policyID)
	if err != nil {
		return "", err
	}

	status, body, err := doOPARequestWithStatus(http.MethodGet, url, bearerToken, "", nil)
	if err != nil {
		return "", err
	}

	// Handle the "initial state" where the policy doesn't exist yet
	if status == http.StatusNotFound {
		return "", nil
	}

	var decoded adminconfig.OPAPolicyGetResponse
	if err := json.Unmarshal(body, &decoded); err != nil {
		return "", perrors.WithMessage(err, "failed to decode OPA response")
	}

	return decoded.Result.Raw, nil
}

// BizPutOPAPolicy updates or creates a policy.
func BizPutOPAPolicy(serverURL, policyID, bearerToken, policy string) error {
	url, err := buildOPAPolicyURL(serverURL, policyID)
	if err != nil {
		return err
	}

	if strings.TrimSpace(policy) == "" {
		return perrors.New("policy content is required")
	}

	normalized := strings.ReplaceAll(policy, "\r\n", "\n")
	_, _, err = doOPARequestWithStatus(http.MethodPut, url, bearerToken, "text/plain", []byte(normalized))
	return err
}

// BizDeleteOPAPolicy removes a policy. Returns nil if policy is already gone.
func BizDeleteOPAPolicy(serverURL, policyID, bearerToken string) error {
	url, err := buildOPAPolicyURL(serverURL, policyID)
	if err != nil {
		return err
	}

	_, _, err = doOPARequestWithStatus(http.MethodDelete, url, bearerToken, "", nil)
	return err
}

func buildOPAPolicyURL(serverURL, policyID string) (string, error) {
	serverURL = strings.TrimSpace(serverURL)
	policyID = strings.TrimSpace(policyID)

	if serverURL == "" || policyID == "" {
		return "", perrors.New("server_url and policy_id are required")
	}

	base := strings.TrimRight(serverURL, "/")
	return fmt.Sprintf("%s/v1/policies/%s", base, policyID), nil
}

// doOPARequestWithStatus is the core proxy function
func doOPARequestWithStatus(method, url, bearerToken, contentType string, body []byte) (int, []byte, error) {
	var reader io.Reader
	if body != nil {
		reader = bytes.NewReader(body)
	}

	ctx := context.Background()
	if adminconfig.Client != nil {
		ctx = adminconfig.Client.GetCtx()
	}
	if timeout := getOPATimeout(); timeout > 0 {
		var cancel context.CancelFunc
		ctx, cancel = context.WithTimeout(ctx, timeout)
		defer cancel()
	}
	req, err := http.NewRequestWithContext(ctx, method, url, reader)
	if err != nil {
		return 0, nil, perrors.Wrap(err, "failed to create OPA request")
	}

	if contentType != "" {
		req.Header.Set("Content-Type", contentType)
	}
	if token := strings.TrimSpace(bearerToken); token != "" {
		req.Header.Set("Authorization", "Bearer "+token)
	}

	resp, err := opaHTTPClient.Do(req)
	if err != nil {
		return 0, nil, perrors.Wrap(err, "OPA connection failed")
	}
	defer resp.Body.Close()

	respBody, err := io.ReadAll(resp.Body)
	if err != nil {
		logger.Warnf("failed to read OPA response body: %v", err)
	}

	switch resp.StatusCode {
	case http.StatusOK:
		return resp.StatusCode, respBody, nil

	case http.StatusNotFound:
		if method == http.MethodGet || method == http.MethodDelete {
			return resp.StatusCode, nil, nil
		}
		return resp.StatusCode, nil, perrors.Errorf("OPA endpoint not found: %s", url)

	default:
		if resp.StatusCode < http.StatusOK || resp.StatusCode >= http.StatusMultipleChoices {
			return resp.StatusCode, nil, perrors.Errorf("OPA status %d: %s", resp.StatusCode, strings.TrimSpace(string(respBody)))
		}
		return resp.StatusCode, respBody, nil
	}
}

func getClusterKey(path string) string {
	return getRootPath(Clusters) + "/" + path
}

func getListenerKey(path string) string {
	return getRootPath(Listeners) + "/" + path
}

func getPluginGroupPrefixKey(unpublished bool) string {
	if unpublished {
		return getUnpublishedRootPath(PluginGroup)
	}
	return getRootPath(PluginGroup)
}

func getPluginRatelimitKey(unpublished bool) string {
	return getFilterPrefixKey(unpublished) + "/" + Ratelimit
}

func getFilterPrefixKey(unpublished bool) string {
	if unpublished {
		return getUnpublishedRootPath(Filter)
	}
	return getRootPath(Filter)
}

func getClusterId() int {
	return loopGetId(getRootPath(ClusterID))
}

func loopGetId(k string) int {
	for {

		rawClient := adminconfig.Client.GetRawClient()
		if rawClient == nil {
			logger.Error("GetId etcd client is null")
			return ErrID
		}

		resp, err := rawClient.Get(adminconfig.Client.GetCtx(), k)
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
			err = adminconfig.Client.Create(k, strconv.Itoa(id))
		} else {
			err = adminconfig.Client.UpdateWithRev(k, strconv.Itoa(id), rev)
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
}

func getRootPath(key string) string {
	return adminconfig.Bootstrap.GetPath() + "/" + key
}

func getUnpublishedRootPath(key string) string {
	return adminconfig.Bootstrap.GetPath() + "/" + Unpublished + "/" + key
}

func getCheckClusterRegexp() *regexp.Regexp {
	return regexp.MustCompile(".+/clusters/[^/]+/?$")
}

func getCheckListenerRegexp() *regexp.Regexp {
	return regexp.MustCompile(".+/listeners/[^/]+/?$")
}

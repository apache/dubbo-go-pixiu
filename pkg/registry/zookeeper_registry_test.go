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
package registry

import (
	"fmt"
	"github.com/apache/dubbo-go/common"
	"github.com/apache/dubbo-go/common/constant"
	zkRegistry "github.com/apache/dubbo-go/registry"
	"github.com/apache/dubbo-go/remoting/zookeeper"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
	"github.com/dubbogo/go-zookeeper/zk"
	perrors "github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

func newMockZkRegistry(providerUrl common.URL, opts ...zookeeper.Option) (*zk.TestCluster, error) {
	var (
		err      error
		c        *zk.TestCluster
		client   *zookeeper.ZookeeperClient
		registry *zkRegistry.BaseRegistry
		//event <-chan zk.Event
	)
	c, client, _, err = zookeeper.NewMockZookeeperClient("test", 15*time.Second, opts...)
	if err != nil {
		return nil, err
	}
	var (
		//revision   string
		params     url.Values
		rawURL     string
		zkPath     string
		encodedURL string
		dubboPath  string
		//conf       config.URL
	)
	params = url.Values{}

	providerUrl.RangeParams(func(key, value string) bool {
		params.Add(key, value)
		return true
	})

	params.Add("pid", "")
	params.Add("ip", "")
	//params.Add("timeout", fmt.Sprintf("%d", int64(r.Timeout)/1e6))

	role, _ := strconv.Atoi(registry.URL.GetParam(constant.ROLE_KEY, ""))
	switch role {

	case common.PROVIDER:
		dubboPath = fmt.Sprintf("/dubbo/%s/%s", url.QueryEscape(providerUrl.Service()), common.DubboNodes[common.PROVIDER])
		client.CreateWithValue(dubboPath, []byte(""))

	}
	var host string
	if providerUrl.Ip == "" {
		host = ""
	} else {
		host = providerUrl.Ip
	}
	rawURL = fmt.Sprintf("%s://%s%s?%s", providerUrl.Protocol, host, c.Path, params.Encode())

	encodedURL = url.QueryEscape(rawURL)
	dubboPath = strings.ReplaceAll(dubboPath, "$", "%24")

	err = client.Create(dubboPath)
	if err != nil {
		return nil, perrors.WithStack(err)
	}

	// try to register the node
	zkPath, err = client.RegisterTemp(dubboPath, encodedURL)
	if err != nil {
		logger.Errorf("Register temp node(root{%s}, node{%s}) = error{%v}", dubboPath, encodedURL, perrors.WithStack(err))
		if perrors.Cause(err) == zk.ErrNodeExists {
			// should delete the old node
			logger.Info("Register temp node failed, try to delete the old and recreate  (root{%s}, node{%s}) , ignore!", dubboPath, encodedURL)
			if err = client.Delete(zkPath); err == nil {
				_, err = client.RegisterTemp(dubboPath, encodedURL)
			}
			if err != nil {
				logger.Errorf("Recreate the temp node failed, (root{%s}, node{%s}) = error{%v}", dubboPath, encodedURL, perrors.WithStack(err))
			}
		}
		return nil, perrors.WithMessagef(err, "RegisterTempNode(root{%s}, node{%s})", dubboPath, encodedURL)
	}
	logger.Debugf("Create a zookeeper node:%s", zkPath)
	defer c.Stop()
	return c, nil
}

func registerProvider() (*zk.TestCluster, error) {
	providerUrl, _ := common.NewURL("dubbo://127.0.0.1:20000/com.ikurento.user.UserProvider", common.WithParamsValue(constant.CLUSTER_KEY, "mock"), common.WithParamsValue("serviceid", "soa.mock"), common.WithMethods([]string{"GetUser", "AddUser"}))
	return newMockZkRegistry(providerUrl)
}

func TestZookeeperRegistryLoad_GetCluster(t *testing.T) {
	testCluster, err := registerProvider()
	assert.Nil(t, err)
	defer testCluster.Stop()
	registryLoad, err := newZookeeperRegistryLoad("127.0.0.1:1111", "test-cluster")
	assert.Nil(t, err)
	assert.NotNil(t, registryLoad)
	clusterName, err := registryLoad.GetCluster()
	assert.NotNil(t, err)
	assert.Equal(t, "test-cluster", clusterName)
}

func TestZookeeperRegistryLoad_LoadAllServices(t *testing.T) {
}

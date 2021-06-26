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
	"net/url"
	"strconv"
	"strings"
	"testing"
	"time"
)

import (
	"github.com/apache/dubbo-go/common"
	"github.com/apache/dubbo-go/common/constant"
	"github.com/dubbogo/go-zookeeper/zk"
	gxzookeeper "github.com/dubbogo/gost/database/kv/zk"
	"github.com/stretchr/testify/assert"
)

const providerUrlstr = "dubbo://127.0.0.1:20000/com.ikurento.user.UserProvider?methods.GetUser.retries=1"

func newMockZkRegistry(t *testing.T, providerURL *common.URL, opts ...gxzookeeper.Option) (*zk.TestCluster, error) {
	var (
		err    error
		c      *zk.TestCluster
		client *gxzookeeper.ZookeeperClient
	)
	c, client, _, err = gxzookeeper.NewMockZookeeperClient("test", 15*time.Second, opts...)
	if err != nil {
		return nil, err
	}
	var (
		params     url.Values
		rawURL     string
		encodedURL string
		dubboPath  string
	)
	params = url.Values{}

	providerURL.RangeParams(func(key, value string) bool {
		params.Add(key, value)
		return true
	})

	dubboPath = fmt.Sprintf("/dubbo/%s/%s", url.QueryEscape(providerURL.Service()), common.DubboNodes[common.PROVIDER])
	err = client.CreateWithValue(dubboPath, []byte(""))
	assert.Nil(t, err)

	if len(providerURL.Methods) > 0 {
		params.Add(constant.METHODS_KEY, strings.Join(providerURL.Methods, ","))
	}
	var host string
	if providerURL.Ip == "" {
		host = ""
	} else {
		host = providerURL.Ip
	}
	host += ":" + providerURL.Port
	rawURL = fmt.Sprintf("%s://%s%s?%s", providerURL.Protocol, host, providerURL.Path, params.Encode())

	encodedURL = url.QueryEscape(rawURL)
	dubboPath = strings.ReplaceAll(dubboPath, "$", "%24")

	err = client.Create(dubboPath)
	assert.Nil(t, err)

	// to register the node
	_, err = client.RegisterTemp(dubboPath, encodedURL)
	assert.Nil(t, err)
	return c, nil
}

func registerProvider(providerURL *common.URL, t *testing.T) (*zk.TestCluster, error) {
	return newMockZkRegistry(t, providerURL)
}

func TestZookeeperRegistryLoad_GetCluster(t *testing.T) {
	providerURL, _ := common.NewURL(providerUrlstr)
	testCluster, err := registerProvider(providerURL, t)
	assert.Nil(t, err)
	defer testCluster.Stop()
	registryLoad, err := newZookeeperRegistryLoad("127.0.0.1:1111", "test-cluster")
	assert.Nil(t, err)
	assert.NotNil(t, registryLoad)
	clusterName, err := registryLoad.GetCluster()
	assert.Nil(t, err)
	assert.Equal(t, "test-cluster", clusterName)
}

func TestZookeeperRegistryLoad_LoadAllServices(t *testing.T) {
	providerURL, _ := common.NewURL(providerUrlstr)
	testCluster, err := registerProvider(providerURL, t)
	assert.Nil(t, err)
	defer testCluster.Stop()
	registryLoad, err := newZookeeperRegistryLoad(fmt.Sprintf("%s:%s", "127.0.0.1", strconv.Itoa(testCluster.Servers[0].Port)), "test-cluster")
	assert.Nil(t, err)
	assert.NotNil(t, registryLoad)
	services, err := registryLoad.LoadAllServices()
	assert.Nil(t, err)
	assert.GreaterOrEqual(t, len(services), 1)
	assert.Equal(t, services[0].Protocol, "dubbo")
	assert.Equal(t, services[0].Location, "127.0.0.1:20000")
	assert.Equal(t, services[0].Path, "/com.ikurento.user.UserProvider")
	assert.Equal(t, services[0].GetMethodParam("GetUser", "retries", ""), "1")
}

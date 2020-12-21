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
	"path"
	"strings"
	"time"
)

import (
	"github.com/apache/dubbo-go/common"
	"github.com/apache/dubbo-go/remoting/zookeeper"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
)

const (
	rootPath = "/dubbo"
)

func init() {
	var _ Loader = new(ZookeeperRegistryLoad)
}

// ZookeeperRegistryLoad load dubbo apis from zookeeper registry
type ZookeeperRegistryLoad struct {
	zkName  string
	client  *zookeeper.ZookeeperClient
	Address string
	cluster string
}

func newZookeeperRegistryLoad(address, cluster string) (Loader, error) {
	newClient, err := zookeeper.NewZookeeperClient("zkClient", strings.Split(address, ","), 15*time.Second)
	if err != nil {
		logger.Warnf("newZookeeperClient error:%v", err)
		return nil, err
	}

	r := &ZookeeperRegistryLoad{
		Address: address,
		client:  newClient,
		cluster: cluster,
	}

	return r, nil
}

// nolint
func (crl *ZookeeperRegistryLoad) GetCluster() (string, error) {
	return crl.cluster, nil
}

// LoadAllServices load all services from zookeeper registry
func (crl *ZookeeperRegistryLoad) LoadAllServices() ([]common.URL, error) {
	children, err := crl.client.GetChildren(rootPath)
	if err != nil {
		logger.Errorf("[zookeeper registry] get zk children error:%v", err)
		return nil, err
	}
	var urls []common.URL
	for _, _interface := range children {
		providerStr := path.Join(rootPath, "/", _interface, "/", "providers")
		urlStrs, err := crl.client.GetChildren(providerStr)
		if err != nil {
			logger.Errorf("[zookeeper registry] get zk children \"%s\" error:%v", providerStr, err)
			return nil, err
		}
		for _, url := range urlStrs {
			dubboURL, err := common.NewURL(url)
			if err != nil {
				logger.Warnf("[zookeeper registry] transfer zk info to url error:%v", err)
				continue
			}
			urls = append(urls, *dubboURL)
		}
	}
	return urls, nil
}

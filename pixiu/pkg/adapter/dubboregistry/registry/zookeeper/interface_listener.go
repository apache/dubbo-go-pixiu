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

package zookeeper

import (
	"path"
	"sync"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/common"
	"github.com/dubbogo/go-zookeeper/zk"
)

import (
	common2 "github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/common"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/registry"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/remoting/zookeeper"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

const (
	rootPath         = "/dubbo"
	providerCategory = "providers"
)

var _ registry.Listener = new(zkIntfListener)

type zkIntfListener struct {
	path            string
	exit            chan struct{}
	client          *zookeeper.ZooKeeperClient
	reg             *ZKRegistry
	wg              sync.WaitGroup
	adapterListener common2.RegistryEventListener
}

// newZKIntfListener returns a new zkIntfListener with pre-defined path according to the registered type.
func newZKIntfListener(client *zookeeper.ZooKeeperClient, reg *ZKRegistry, adapterListener common2.RegistryEventListener) registry.Listener {
	p := rootPath
	return &zkIntfListener{
		path:            p,
		exit:            make(chan struct{}),
		client:          client,
		reg:             reg,
		adapterListener: adapterListener,
	}
}

func (z *zkIntfListener) Close() {
	close(z.exit)
	z.wg.Wait()
}

func (z *zkIntfListener) WatchAndHandle() {
	z.wg.Add(1)
	go z.watch()
}

func (z *zkIntfListener) watch() {
	defer z.wg.Done()
	var (
		failTimes  int64 = 0
		delayTimer       = time.NewTimer(ConnDelay * time.Duration(failTimes))
	)
	defer delayTimer.Stop()
	for {
		_, e, err := z.client.GetChildrenW(z.path)
		// error handling
		if err != nil {
			failTimes++
			logger.Infof("watching (path{%s}) = error{%v}", z.path, err)
			// Exit the watch if root node is in error
			if err == zookeeper.ErrNilNode {
				logger.Errorf("watching (path{%s}) got errNilNode,so exit listen", z.path)
				return
			}
			if failTimes > MaxFailTimes {
				logger.Errorf("Error happens on (path{%s}) exceed max fail times: %s,so exit listen",
					z.path, MaxFailTimes)
				return
			}
			delayTimer.Reset(ConnDelay * time.Duration(failTimes))
			<-delayTimer.C
			continue
		}
		failTimes = 0
		if continueLoop := z.waitEventAndHandlePeriod(z.path, e); !continueLoop {
			return
		}
	}
}

func (z *zkIntfListener) waitEventAndHandlePeriod(path string, e <-chan zk.Event) bool {
	tickerTTL := defaultTTL
	ticker := time.NewTicker(tickerTTL)
	defer ticker.Stop()
	z.handleEvent(z.path)
	for {
		select {
		case <-ticker.C:
			z.handleEvent(z.path)
		case zkEvent := <-e:
			logger.Warnf("get a zookeeper e{type:%s, server:%s, path:%s, state:%d-%s, err:%s}",
				zkEvent.Type.String(), zkEvent.Server, zkEvent.Path, zkEvent.State, zookeeper.StateToString(zkEvent.State), zkEvent.Err)
			if zkEvent.Type != zk.EventNodeChildrenChanged {
				return true
			}
			z.handleEvent(zkEvent.Path)
			return true
		case <-z.exit:
			logger.Warnf("listen(path{%s}) goroutine exit now...", z.path)
			return false
		}
	}
}

func (z *zkIntfListener) handleEvent(basePath string) {
	newChildren, err := z.client.GetChildren(basePath)
	if err != nil {
		logger.Errorf("Error when retrieving newChildren in path: %s, Error:%s", basePath, err.Error())
	}
	for i := range newChildren {
		if newChildren[i] == "metadata" {
			continue
		}
		providerPath := path.Join(basePath, newChildren[i], providerCategory)
		// TO-DO: modify here to only handle child that changed
		providers, err := z.client.GetChildren(providerPath)
		if err != nil {
			logger.Warnf("Get provider %s failed due to %s", providerPath, err.Error())
			continue
		}
		srvUrl, err := common.NewURL(providers[0])
		if err != nil {
			logger.Warnf("Parse provider service url %s failed due to %s", providers[0], err.Error())
			continue
		}
		if z.reg.GetSvcListener(srvUrl.ServiceKey()) != nil {
			continue
		}
		l := newZkSrvListener(srvUrl, providerPath, z.client, z.adapterListener)
		l.wg.Add(1)
		go l.WatchAndHandle()
		z.reg.SetSvcListener(srvUrl.ServiceKey(), l)
	}
}

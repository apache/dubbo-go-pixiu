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
	"strings"
	"sync"
	"time"
)

import (
	"github.com/dubbogo/go-zookeeper/zk"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/common"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/registry"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/remoting/zookeeper"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

const (
	defaultServicesPath = "/services"
	methodsRootPath     = "/dubbo/metadata"
)

var _ registry.Listener = new(zkAppListener)

type zkAppListener struct {
	servicesPath    string
	exit            chan struct{}
	client          *zookeeper.ZooKeeperClient
	reg             *ZKRegistry
	wg              sync.WaitGroup
	adapterListener common.RegistryEventListener
}

// newZkAppListener returns a new newZkAppListener with pre-defined servicesPath according to the registered type.
func newZkAppListener(client *zookeeper.ZooKeeperClient, reg *ZKRegistry, adapterListener common.RegistryEventListener) registry.Listener {
	p := defaultServicesPath
	return &zkAppListener{
		servicesPath:    p,
		exit:            make(chan struct{}),
		client:          client,
		reg:             reg,
		adapterListener: adapterListener,
	}
}

func (z *zkAppListener) Close() {
	close(z.exit)
	z.wg.Wait()
}

func (z *zkAppListener) WatchAndHandle() {
	z.wg.Add(1)
	go z.watch()
}

func (z *zkAppListener) watch() {
	defer z.wg.Done()

	var (
		failTimes  int64 = 0
		delayTimer       = time.NewTimer(ConnDelay * time.Duration(failTimes))
	)
	defer delayTimer.Stop()
	for {
		children, e, err := z.client.GetChildrenW(z.servicesPath)
		// error handling
		if err != nil {
			failTimes++
			logger.Infof("watching (path{%s}) = error{%v}", z.servicesPath, err)
			// Exit the watch if root node is in error
			if err == zookeeper.ErrNilNode {
				logger.Errorf("watching (path{%s}) got errNilNode,so exit listen", z.servicesPath)
				return
			}
			if failTimes > MaxFailTimes {
				logger.Errorf("Error happens on (path{%s}) exceed max fail times: %s,so exit listen",
					z.servicesPath, MaxFailTimes)
				return
			}
			delayTimer.Reset(ConnDelay * time.Duration(failTimes))
			<-delayTimer.C
			continue
		}
		failTimes = 0
		if continueLoop := z.waitEventAndHandlePeriod(children, e); !continueLoop {
			return
		}
	}
}

func (z *zkAppListener) waitEventAndHandlePeriod(children []string, e <-chan zk.Event) bool {
	tickerTTL := defaultTTL
	ticker := time.NewTicker(tickerTTL)
	defer ticker.Stop()
	z.handleEvent(children)
	for {
		select {
		case <-ticker.C:
			z.handleEvent(children)
		case zkEvent := <-e:
			logger.Warnf("get a zookeeper e{type:%s, server:%s, path:%s, state:%d-%s, err:%s}",
				zkEvent.Type.String(), zkEvent.Server, zkEvent.Path, zkEvent.State, zookeeper.StateToString(zkEvent.State), zkEvent.Err)
			if zkEvent.Type != zk.EventNodeChildrenChanged {
				return true
			}
			z.handleEvent(children)
			return true
		case <-z.exit:
			logger.Warnf("listen(path{%s}) goroutine exit now...", z.servicesPath)
			return false
		}
	}
}

func (z *zkAppListener) handleEvent(children []string) {
	fetchChildren, err := z.client.GetChildren(z.servicesPath)
	if err != nil {
		logger.Warnf("Error when retrieving newChildren in path: %s, Error:%s", z.servicesPath, err.Error())
	}
	for _, path := range fetchChildren {
		serviceName := strings.Join([]string{z.servicesPath, path}, constant.PathSlash)
		if z.reg.GetSvcListener(serviceName) != nil {
			continue
		}
		l := newApplicationServiceListener(serviceName, z.client, z.adapterListener)
		l.wg.Add(1)
		go l.WatchAndHandle()
		z.reg.SetSvcListener(serviceName, l)
	}
}

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
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/registry"
	"github.com/apache/dubbo-go-pixiu/pkg/remoting/zookeeper"
	_ "github.com/apache/dubbo-go/cluster/cluster_impl"
	_ "github.com/apache/dubbo-go/cluster/loadbalance"
	_ "github.com/apache/dubbo-go/common/proxy/proxy_factory"
	_ "github.com/apache/dubbo-go/filter/filter_impl"
	_ "github.com/apache/dubbo-go/metadata/service/inmemory"
	_ "github.com/apache/dubbo-go/metadata/service/remote"
	_ "github.com/apache/dubbo-go/registry/protocol"
	_ "github.com/apache/dubbo-go/registry/zookeeper"
	"github.com/dubbogo/go-zookeeper/zk"
	"strings"
	"sync"
	"time"
)

const (
	defaultServicesPath = "/services"
	methodsRootPath     = "/dubbo/metadata"
)

var _ registry.Listener = new(zkAppListener)

type zkAppListener struct {
	servicesPath string
	exit         chan struct{}
	client       *zookeeper.ZooKeeperClient
	reg          *ZKRegistry
	wg           sync.WaitGroup
}

// newZkAppListener returns a new newZkAppListener with pre-defined servicesPath according to the registered type.
func newZkAppListener(client *zookeeper.ZooKeeperClient, reg *ZKRegistry) registry.Listener {
	p := defaultServicesPath
	return &zkAppListener{
		servicesPath: p,
		exit:         make(chan struct{}),
		client:       client,
		reg:          reg,
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
		_, e, err := z.client.GetChildrenW(z.servicesPath)
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
		tickerTTL := defaultTTL
		ticker := time.NewTicker(tickerTTL)
		z.handleEvent(z.servicesPath)
	WATCH:
		for {
			select {
			case <-ticker.C:
				z.handleEvent(z.servicesPath)
			case zkEvent := <-e:
				logger.Warnf("get a zookeeper e{type:%s, server:%s, path:%s, state:%d-%s, err:%s}",
					zkEvent.Type.String(), zkEvent.Server, zkEvent.Path, zkEvent.State, zookeeper.StateToString(zkEvent.State), zkEvent.Err)
				ticker.Stop()
				if zkEvent.Type != zk.EventNodeChildrenChanged {
					break WATCH
				}
				z.handleEvent(zkEvent.Path)
				break WATCH
			case <-z.exit:
				logger.Warnf("listen(path{%s}) goroutine exit now...", z.servicesPath)
				ticker.Stop()
				return
			}
		}

	}
}

func (z *zkAppListener) handleEvent(basePath string) {
	paths := z.getIssPaths(basePath)
	for _, path := range paths {
		if z.reg.GetSvcListener(path) != nil {
			continue
		}
		l := newApplicationServiceListener(path, z.client)
		l.wg.Add(1)
		go l.WatchAndHandle()
		z.reg.SetSvcListener(path, l)
	}
}

// getIssPaths will return the paths of instances
func (z *zkAppListener) getIssPaths(basePath string) []string {
	subPaths, err := z.client.GetChildren(basePath)
	if err != nil {
		logger.Errorf("Error when retrieving newChildren in path: %s, Error:%s", basePath, err.Error())
		return nil
	}
	if len(subPaths) == 0 {
		return nil
	}

	paths := []string{}
	for _, subPath := range subPaths {
		serviceName := strings.Join([]string{z.servicesPath, subPath}, constant.PathSlash)
		ids, err := z.client.GetChildren(serviceName)
		if err != nil {
			logger.Warnf("Get serviceIDs %s failed due to %s", serviceName, err.Error())
			continue
		}
		for _, id := range ids {
			insPath := strings.Join([]string{serviceName, id}, constant.PathSlash)
			paths = append(paths, insPath)
		}
	}
	return paths
}

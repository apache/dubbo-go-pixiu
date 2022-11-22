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
	gzk "github.com/dubbogo/gost/database/kv/zk"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/springcloud/common"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/remote/zookeeper"
)

type zkAppListener struct {
	servicesPath string
	exit         chan struct{}
	svcListeners *SvcListeners
	wg           sync.WaitGroup

	ds *zookeeperDiscovery
}

func newZkAppListener(ds *zookeeperDiscovery) zookeeper.Listener {
	return &zkAppListener{
		servicesPath: ds.basePath,
		exit:         make(chan struct{}),
		svcListeners: &SvcListeners{listeners: make(map[string]zookeeper.Listener), listenerLock: sync.Mutex{}},
		ds:           ds,
	}
}

func (z *zkAppListener) Close() {
	for k, listener := range z.svcListeners.GetAllListener() {
		z.svcListeners.RemoveListener(k)
		listener.Close()
	}
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
		children, e, err := z.ds.getClient().GetChildrenW(z.servicesPath)
		if err != nil {
			failTimes++
			logger.Infof("watching (path{%s}) = error{%v}", z.servicesPath, err)
			if err == zk.ErrNoNode {
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
		if continueLoop := z.watchEventHandle(children, e); !continueLoop {
			return
		}
	}
}

func (z *zkAppListener) watchEventHandle(children []string, e <-chan zk.Event) bool {
	tickerTTL := defaultTTL
	ticker := time.NewTicker(tickerTTL)
	defer ticker.Stop()
	z.handleEvent(children)
	for {
		select {
		case <-ticker.C:
			z.handleEvent(children)
		case zkEvent := <-e:
			logger.Debugf("get a zookeeper e{type:%s, server:%s, path:%s, state:%d-%s, err:%s}",
				zkEvent.Type.String(), zkEvent.Server, zkEvent.Path, zkEvent.State, gzk.StateToString(zkEvent.State), zkEvent.Err)
			z.handleEvent(children)
			return true
		case <-z.exit:
			logger.Warnf("listen(path{%s}) goroutine exit now...", z.servicesPath)
			return false
		}
	}
}

func (z *zkAppListener) handleEvent(children []string) {

	fetchChildren, err := z.ds.getClient().GetChildren(z.servicesPath)

	if err != nil {
		// todo refactor gost zk, make it return the definite err
		if strings.Contains(err.Error(), "none children") {
			logger.Debugf("%s get nodes from zookeeper fail: %s", common.ZKLogDiscovery, err.Error())
		} else {
			logger.Warnf("Error when retrieving newChildren in path: %s, Error:%s", z.servicesPath, err.Error())
		}
	}

	discovery := z.ds
	serviceMap := discovery.getServiceMap()

	// del services
	for sn := range serviceMap {

		if !contains(fetchChildren, sn) {

			// service zk event listener
			serviceNodePath := strings.Join([]string{z.servicesPath, sn}, constant.PathSlash)
			z.svcListeners.RemoveListener(serviceNodePath)

			// service cluster
			for _, instance := range serviceMap[sn] {
				_, _ = discovery.delServiceInstance(instance)
			}
		}
	}

	for _, serviceName := range fetchChildren {
		serviceNodePath := strings.Join([]string{z.servicesPath, serviceName}, constant.PathSlash)
		instances := serviceMap[serviceName]
		if len(instances) == 0 {
			z.svcListeners.RemoveListener(serviceNodePath)
		}
		if z.svcListeners.GetListener(serviceNodePath) != nil {
			continue
		}
		l := newApplicationServiceListener(serviceNodePath, serviceName, discovery)
		l.wg.Add(1)
		go l.WatchAndHandle()
		z.svcListeners.SetListener(serviceNodePath, l)
	}
}

type SvcListeners struct {
	listeners    map[string]zookeeper.Listener
	listenerLock sync.Mutex
}

func (s *SvcListeners) GetListener(id string) zookeeper.Listener {
	s.listenerLock.Lock()
	defer s.listenerLock.Unlock()
	listener, ok := s.listeners[id]
	if !ok {
		return nil
	}
	return listener
}

func (s *SvcListeners) SetListener(id string, listener zookeeper.Listener) {
	s.listenerLock.Lock()
	defer s.listenerLock.Unlock()
	s.listeners[id] = listener
}

func (s *SvcListeners) RemoveListener(id string) {
	s.listenerLock.Lock()
	defer s.listenerLock.Unlock()
	delete(s.listeners, id)
}

func (s *SvcListeners) GetAllListener() map[string]zookeeper.Listener {
	s.listenerLock.Lock()
	defer s.listenerLock.Unlock()
	return s.listeners
}

func contains(elems []string, v string) bool {
	for _, s := range elems {
		if v == s {
			return true
		}
	}
	return false
}

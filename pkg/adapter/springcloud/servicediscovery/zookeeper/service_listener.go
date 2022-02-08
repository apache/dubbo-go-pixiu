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
	"sync"
	"time"
)

import (
	"github.com/dubbogo/go-zookeeper/zk"
	gzk "github.com/dubbogo/gost/database/kv/zk"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/springcloud/common"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

type applicationServiceListener struct {
	servicePath string
	serviceName string
	exit        chan struct{}
	wg          sync.WaitGroup

	ds *zookeeperDiscovery
}

func newApplicationServiceListener(path string, serviceName string, ds *zookeeperDiscovery) *applicationServiceListener {
	return &applicationServiceListener{
		servicePath: path,
		exit:        make(chan struct{}),
		ds:          ds,
		serviceName: serviceName,
	}
}

func (asl *applicationServiceListener) WatchAndHandle() {
	defer asl.wg.Done()

	var (
		failTimes  int64 = 0
		delayTimer       = time.NewTimer(ConnDelay * time.Duration(failTimes))
	)
	defer delayTimer.Stop()
	for {
		children, e, err := asl.ds.getClient().GetChildrenW(asl.servicePath)
		if err != nil {
			failTimes++
			logger.Infof("watching (path{%s}) = error{%v}", asl.servicePath, err)
			if err == gzk.ErrNilChildren {
				return
			}
			if err == gzk.ErrNilNode {
				logger.Errorf("watching (path{%s}) got errNilNode,so exit listen", asl.servicePath)
				return
			}
			if failTimes > MaxFailTimes {
				logger.Errorf("Error happens on (path{%s}) exceed max fail times: %v,so exit listen", asl.servicePath, MaxFailTimes)
				return
			}
			delayTimer.Reset(ConnDelay * time.Duration(failTimes))
			<-delayTimer.C
			continue
		}
		failTimes = 0
		if continueLoop := asl.watchEventHandle(children, e); !continueLoop {
			return
		}

	}
}

func (asl *applicationServiceListener) watchEventHandle(children []string, e <-chan zk.Event) bool {
	tickerTTL := defaultTTL
	ticker := time.NewTicker(tickerTTL)
	defer ticker.Stop()
	asl.handleEvent(children)
	for {
		select {
		case <-ticker.C:
			asl.handleEvent(children)
		case zkEvent := <-e:
			logger.Warnf("get a zookeeper e{type:%s, server:%s, path:%s, state:%d-%s, err:%s}",
				zkEvent.Type.String(), zkEvent.Server, zkEvent.Path, zkEvent.State, gzk.StateToString(zkEvent.State), zkEvent.Err)
			if zkEvent.Type != zk.EventNodeChildrenChanged {
				return true
			}
			asl.handleEvent(children)
			return true
		case <-asl.exit:
			logger.Warnf("listen(path{%s}) goroutine exit now...", asl.servicePath)
			return false
		}
	}
}

func (asl *applicationServiceListener) handleEvent(children []string) {

	fetchChildren, err := asl.ds.getClient().GetChildren(asl.servicePath)
	if err != nil {
		logger.Warnf("%s Error when retrieving service node [%s] in path: %s, Error:%s", common.ZKLogDiscovery, asl.serviceName, asl.servicePath, err.Error())
		return
	}
	discovery := asl.ds
	instanceMap := discovery.instanceMap
	addf := func() {
		if addInstanceIds := Diff(fetchChildren, children); addInstanceIds != nil {
			for _, id := range addInstanceIds {
				discovery.addServiceInstance(instanceMap[id])
			}
		}
	}
	addf()

	delf := func() {
		if delInstanceIds := Diff(children, fetchChildren); delInstanceIds != nil {
			for _, id := range delInstanceIds {
				discovery.delServiceInstance(instanceMap[id])
			}
		}
	}
	delf()
}

// Close closes this listener
func (asl *applicationServiceListener) Close() {
	close(asl.exit)
	asl.wg.Wait()
}

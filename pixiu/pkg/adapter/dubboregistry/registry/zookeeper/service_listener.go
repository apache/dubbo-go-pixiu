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
	"dubbo.apache.org/dubbo-go/v3/common"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/dubbogo/go-zookeeper/zk"
)

import (
	common2 "github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/common"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/registry"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/dubboregistry/remoting/zookeeper"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
)

var _ registry.Listener = new(serviceListener)

// serviceListener normally monitors the /dubbo/[:url.service()]/providers
type serviceListener struct {
	url    *common.URL
	path   string
	client *zookeeper.ZooKeeperClient

	exit            chan struct{}
	wg              sync.WaitGroup
	adapterListener common2.RegistryEventListener
	registryMethod  map[string]*config.Method
	mutex           sync.Mutex
}

// newZkSrvListener creates a new zk service listener
func newZkSrvListener(url *common.URL, path string, client *zookeeper.ZooKeeperClient, adapterListener common2.RegistryEventListener) *serviceListener {
	return &serviceListener{
		url:             url,
		path:            path,
		client:          client,
		exit:            make(chan struct{}),
		adapterListener: adapterListener,
		registryMethod:  make(map[string]*config.Method),
	}
}

func (zkl *serviceListener) WatchAndHandle() {
	defer zkl.wg.Done()

	var (
		failTimes  int64 = 0
		delayTimer       = time.NewTimer(ConnDelay * time.Duration(failTimes))
	)

	for {
		children, e, err := zkl.client.GetChildrenW(zkl.path)
		// error handling
		if err != nil {
			failTimes++
			logger.Infof("watching (path{%s}) = error{%v}", zkl.path, err)
			// Exit the watch if root node is in error
			if err == zookeeper.ErrNilNode {
				logger.Errorf("watching (path{%s}) got errNilNode,so exit listen", zkl.path)
				return
			}
			if failTimes > MaxFailTimes {
				logger.Errorf("Error happens on (path{%s}) exceed max fail times: %s,so exit listen",
					zkl.path, MaxFailTimes)
				return
			}
			delayTimer.Reset(ConnDelay * time.Duration(failTimes))
			<-delayTimer.C
			continue
		}
		failTimes = 0
		if continueLoop := zkl.waitEventAndHandlePeriod(children, e); !continueLoop {
			return
		}
	}
}

func (zkl *serviceListener) waitEventAndHandlePeriod(children []string, e <-chan zk.Event) bool {
	tickerTTL := defaultTTL
	ticker := time.NewTicker(tickerTTL)
	defer ticker.Stop()
	zkl.handleEvent()

	for {
		select {
		case <-ticker.C:
			zkl.handleEvent()
		case zkEvent := <-e:
			logger.Warnf("get a zookeeper childEventCh{type:%s, server:%s, path:%s, state:%d-%s, err:%s}",
				zkEvent.Type.String(), zkEvent.Server, zkEvent.Path, zkEvent.State, zookeeper.StateToString(zkEvent.State), zkEvent.Err)
			ticker.Stop()
			if zkEvent.Type != zk.EventNodeChildrenChanged {
				return true
			}
			zkl.handleEvent()
			return true
		case <-zkl.exit:
			logger.Warnf("listen(path{%s}) goroutine exit now...", zkl.path)
			ticker.Stop()
			return false
		}
	}
}

// whenever it is called, the children node changed and refresh the api configuration.
func (zkl *serviceListener) handleEvent() {
	children, err := zkl.client.GetChildren(zkl.path)
	if err != nil {
		// disable the API
		bkConf, methods, _, _ := registry.ParseDubboString(zkl.url.String())
		apiPattern := registry.GetAPIPattern(bkConf)
		for i := range methods {
			path := strings.Join([]string{apiPattern, methods[i]}, constant.PathSlash)
			if err := zkl.adapterListener.OnDeleteRouter(config.Resource{Path: path}); err != nil {
				logger.Errorf("Error={%s} when try to remove API by path: %s", err.Error(), path)
			}
		}
		return
	}
	zkl.url, err = common.NewURL(children[0])
	if err != nil {
		logger.Warnf("Parse service path failed: %s", children[0])
	}
	bkConfig, methods, location, err := registry.ParseDubboString(children[0])
	if err != nil {
		logger.Warnf("Parse dubbo interface provider %s failed; due to \n %s", children[0], err.Error())
		return
	}
	if len(bkConfig.ApplicationName) == 0 || len(bkConfig.Interface) == 0 {
		return
	}

	mappingParams := []config.MappingParam{
		{
			Name:  "requestBody.values",
			MapTo: "opt.values",
		},
		{
			Name:  "requestBody.types",
			MapTo: "opt.types",
		},
	}
	apiPattern := registry.GetAPIPattern(bkConfig)
	zkl.mutex.Lock()
	defer zkl.mutex.Unlock()
	for i := range methods {
		api := registry.CreateAPIConfig(apiPattern, location, bkConfig, methods[i], mappingParams)
		key := api.URLPattern + ":" + string(api.Method.HTTPVerb)
		if _, ok := zkl.registryMethod[key]; ok {
			return
		}
		if err := zkl.adapterListener.OnAddAPI(api); err != nil {
			logger.Errorf("Error={%s} happens when try to add api %s", err.Error(), api.Path)
		} else {
			zkl.registryMethod[key] = &api.Method
		}
	}
}

// Close closes this listener
func (zkl *serviceListener) Close() {
	close(zkl.exit)
	zkl.wg.Wait()
}

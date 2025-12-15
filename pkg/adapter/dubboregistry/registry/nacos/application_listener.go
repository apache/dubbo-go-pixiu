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

package nacos

import (
	"strings"
	"sync"
	"time"
)

import (
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

import (
	common2 "github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/common"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/remoting/zookeeper"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

var _ registry.Listener = new(nacosAppListener)

type nacosAppListener struct {
	exit            chan struct{}
	client          naming_client.INamingClient
	regConf         *model.Registry
	reg             *NacosRegistry
	wg              sync.WaitGroup
	addr            string
	adapterListener common2.RegistryEventListener
	appInfoMap      map[string]*applicationInfo
}

// newNacosAppListener returns a new nacosAppListener with pre-defined path according to the registered type.
func newNacosAppListener(client naming_client.INamingClient, reg *NacosRegistry, regConf *model.Registry, adapterListener common2.RegistryEventListener) registry.Listener {
	return &nacosAppListener{
		exit:            make(chan struct{}),
		client:          client,
		regConf:         regConf,
		reg:             reg,
		addr:            regConf.Address,
		adapterListener: adapterListener,
		appInfoMap:      map[string]*applicationInfo{},
	}
}

func (n *nacosAppListener) Close() {
	close(n.exit)
	n.wg.Wait()
}

func (n *nacosAppListener) WatchAndHandle() {
	n.wg.Add(1)
	go n.watch()
}

func (n *nacosAppListener) watch() {
	defer n.wg.Done()
	var (
		failTimes  int64 = 0
		delayTimer       = time.NewTimer(ConnDelay * time.Duration(failTimes))
	)
	defer delayTimer.Stop()
	
	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()
	
	for {
		// Check for exit signal before processing
		select {
		case <-n.exit:
			logger.Info("nacosAppListener watch goroutine received exit signal, shutting down gracefully")
			return
		default:
		}
		
		serviceList, err := n.client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
			GroupName: n.regConf.Group,
			NameSpace: n.regConf.Namespace,
			PageSize:  100,
		})
		if err != nil {
			failTimes++
			logger.Infof("watching nacos interface with error{%v}", err)
			// Exit the watch if root node is in error
			// TODO: do not use zookeeper error
			if err == zookeeper.ErrNilNode {
				logger.Errorf("watching nacos services got errNilNode,so exit listen")
				return
			}
			if failTimes > MaxFailTimes {
				logger.Errorf("Error happens on nacos exceed max fail times: %s,so exit listen", MaxFailTimes)
				return
			}
			delayTimer.Reset(ConnDelay * time.Duration(failTimes))
			select {
			case <-n.exit:
				logger.Info("nacosAppListener watch goroutine received exit signal during error backoff")
				return
			case <-delayTimer.C:
			}
			continue
		}
		failTimes = 0
		if err := n.updateServiceList(serviceList.Doms); err != nil {
			logger.Errorf("update service list failed %s", err)
		}
		
		// Wait for next tick or exit signal
		select {
		case <-n.exit:
			logger.Info("nacosAppListener watch goroutine received exit signal, shutting down gracefully")
			return
		case <-ticker.C:
		}
	}
}

type applicationInfo struct {
	appName  string
	listener *appServiceListener
}

func (a *applicationInfo) String() string {
	return a.appName
}

func fromServiceKey(serviceKey string) *applicationInfo {
	// if serviceKey contains ":" means it is a interface registry
	// we should ignore it
	if strings.Contains(serviceKey, ":") {
		return nil
	}
	return &applicationInfo{
		appName: serviceKey,
	}
}

func (n *nacosAppListener) updateServiceList(serviceList []string) error {
	// add new service info and watch
	newServiceMap := make(map[string]struct{}, len(serviceList))

	for _, v := range serviceList {
		appInfo := fromServiceKey(v)
		if appInfo == nil {
			// ignore interface registry
			continue
		}
		key := appInfo.String()
		newServiceMap[key] = struct{}{}
		if _, ok := n.appInfoMap[key]; !ok {
			l := newNacosAppSrvListener(n.client, n.adapterListener)
			l.wg.Add(1)
			defer l.wg.Done() // Ensure wg.Done() is always called, even on error paths

			appInfo.listener = l
			n.appInfoMap[key] = appInfo

			sub := &vo.SubscribeParam{
				ServiceName:       appInfo.appName,
				SubscribeCallback: l.Callback,
				GroupName:         n.regConf.Group,
			}

			// Retry subscription with exponential backoff
			var subscribeErr error
			for retry := 0; retry < MaxSubscribeRetry; retry++ {
				subscribeErr = n.client.Subscribe(sub)
				if subscribeErr == nil {
					logger.Infof("successfully subscribed to application %s", key)
					break
				}
				
				// Check if it's a "hosts is empty" error, which is expected during startup
				if strings.Contains(subscribeErr.Error(), "hosts is empty") {
					logger.Warnf("subscribe attempt %d/%d for application %s: hosts is empty, will retry after %v", 
						retry+1, MaxSubscribeRetry, key, SubscribeRetryWait)
					if retry < MaxSubscribeRetry-1 {
						time.Sleep(SubscribeRetryWait)
						continue
					}
					// On last retry, log as warning instead of error since service might register later
					logger.Warnf("subscribe to application %s still has no hosts after %d retries, will continue monitoring", 
						key, MaxSubscribeRetry)
					subscribeErr = nil // Don't treat as fatal error
					break
				}
				
				// For other errors, log and retry
				logger.Warnf("subscribe attempt %d/%d for application %s failed: %s", 
					retry+1, MaxSubscribeRetry, key, subscribeErr)
				if retry < MaxSubscribeRetry-1 {
					time.Sleep(SubscribeRetryWait)
				}
			}
			
			if subscribeErr != nil {
				logger.Errorf("failed to subscribe to application %s after %d retries: %s", 
					key, MaxSubscribeRetry, subscribeErr)
				// Clean up orphaned entry to prevent resource leak
				// Note: wg.Done() will be called by defer before Close() to avoid deadlock
				delete(n.appInfoMap, key)
				l.Close()
			}
		}
	}

	// handle deleted service
	for k, v := range n.appInfoMap {
		if _, ok := newServiceMap[k]; !ok {
			delete(n.appInfoMap, k)
			v.listener.Close()
		}
	}

	return nil
}

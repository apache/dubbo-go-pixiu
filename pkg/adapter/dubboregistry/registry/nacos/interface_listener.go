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
	"bytes"
	"fmt"
	"strings"
	"sync"
	"time"
)

import (
	dubboCommon "dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

import (
	common2 "github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/common"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	MaxFailTimes       = 2
	ConnDelay          = 3 * time.Second
	MaxSubscribeRetry  = 5               // Increased retries to allow more time for provider registration
	SubscribeRetryWait = 3 * time.Second // Increased wait time between retries
)

var _ registry.Listener = new(nacosIntfListener)

type nacosIntfListener struct {
	exit            chan struct{}
	client          naming_client.INamingClient
	regConf         *model.Registry
	reg             *NacosRegistry
	wg              sync.WaitGroup
	addr            string
	adapterListener common2.RegistryEventListener
	serviceInfoMap  map[string]*serviceInfo
}

// newNacosIntfListener returns a new nacosIntfListener with pre-defined path according to the registered type.
func newNacosIntfListener(client naming_client.INamingClient, reg *NacosRegistry, regConf *model.Registry, adapterListener common2.RegistryEventListener) registry.Listener {
	return &nacosIntfListener{
		exit:            make(chan struct{}),
		client:          client,
		regConf:         regConf,
		reg:             reg,
		addr:            regConf.Address,
		adapterListener: adapterListener,
		serviceInfoMap:  map[string]*serviceInfo{},
	}
}

func (n *nacosIntfListener) Close() {
	close(n.exit)
	n.wg.Wait()
	// Cleanup all subscribed service listeners to prevent resource leaks
	for _, v := range n.serviceInfoMap {
		if v.listener != nil {
			v.listener.Close()
		}
	}
}

func (n *nacosIntfListener) WatchAndHandle() {
	n.wg.Add(1)
	go n.watch()
}

func (n *nacosIntfListener) watch() {
	defer n.wg.Done()
	var failTimes int64 = 0

	// Initial wait to allow Nacos and providers time to initialize
	logger.Info("nacosIntfListener waiting for initial service registration...")
	select {
	case <-n.exit:
		logger.Info("nacosIntfListener received exit signal during initial wait")
		return
	case <-time.After(2 * time.Second):
	}

	ticker := time.NewTicker(time.Second * 5)
	defer ticker.Stop()

	for {
		// Check for exit signal before processing
		select {
		case <-n.exit:
			logger.Info("nacosIntfListener watch goroutine received exit signal, shutting down gracefully")
			return
		default:
		}

		serviceList, err := n.client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
			GroupName: n.regConf.Group,
			NameSpace: n.regConf.Namespace,
			PageSize:  100,
		})
		// error handling
		if err != nil {
			failTimes++
			logger.Infof("watching nacos interface with error{%v}", err)
			if failTimes > MaxFailTimes {
				logger.Errorf("Error happens on nacos exceed max fail times: %d,so exit listen", MaxFailTimes)
				return
			}
			// Create timer only when needed for backoff (avoids 0-duration timer bug)
			delayTimer := time.NewTimer(ConnDelay * time.Duration(failTimes))
			select {
			case <-n.exit:
				delayTimer.Stop()
				logger.Info("nacosIntfListener watch goroutine received exit signal during error backoff")
				return
			case <-delayTimer.C:
				delayTimer.Stop()
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
			logger.Info("nacosIntfListener watch goroutine received exit signal, shutting down gracefully")
			return
		case <-ticker.C:
		}
	}
}

type serviceInfo struct {
	interfaceName string
	version       string
	group         string
	listener      *serviceListener
}

func (s *serviceInfo) String() string {
	return fmt.Sprintf("%s:%s:%s", s.interfaceName, s.version, s.group)
}

func fromServiceFullKey(fullKey string) *serviceInfo {
	serviceInfoStrs := strings.Split(fullKey, ":")
	if len(serviceInfoStrs) != 4 {
		return nil
	}
	return &serviceInfo{
		interfaceName: serviceInfoStrs[1],
		version:       serviceInfoStrs[2],
		group:         serviceInfoStrs[3],
	}
}

func (n *nacosIntfListener) updateServiceList(serviceList []string) error {
	// add new service info and watch
	newServiceMap := make(map[string]struct{})

	for _, v := range serviceList {
		svcInfo := fromServiceFullKey(v)
		if svcInfo == nil {
			// invalid nacos dubbo service key
			continue
		}
		key := svcInfo.String()
		newServiceMap[key] = struct{}{}
		if _, ok := n.serviceInfoMap[key]; !ok {
			url, _ := dubboCommon.NewURL("mock://localhost:8848")
			url.SetParam(constant.InterfaceKey, svcInfo.interfaceName)
			url.SetParam(constant.GroupKey, svcInfo.group)
			url.SetParam(constant.VersionKey, svcInfo.version)
			l := newNacosSrvListener(url, n.client, n.adapterListener)
			l.wg.Add(1)

			svcInfo.listener = l
			n.serviceInfoMap[key] = svcInfo

			sub := &vo.SubscribeParam{
				ServiceName:       getSubscribeName(url),
				SubscribeCallback: l.Callback,
				GroupName:         n.regConf.Group,
			}

			// Retry subscription with exponential backoff
			var subscribeErr error
			for retry := 0; retry < MaxSubscribeRetry; retry++ {
				subscribeErr = n.client.Subscribe(sub)
				if subscribeErr == nil {
					logger.Infof("successfully subscribed to service %s", key)
					break
				}

				// Check if it's a "hosts is empty" error, which is expected during startup
				if strings.Contains(subscribeErr.Error(), "hosts is empty") {
					logger.Warnf("subscribe attempt %d/%d for service %s: hosts is empty, will retry after %v",
						retry+1, MaxSubscribeRetry, key, SubscribeRetryWait)
					if retry < MaxSubscribeRetry-1 {
						// Wait with exit signal check for graceful shutdown
						select {
						case <-n.exit:
							logger.Info("nacosIntfListener received exit signal during subscription retry, aborting")
							delete(n.serviceInfoMap, key)
							l.wg.Done()
							l.Close()
							return nil
						case <-time.After(SubscribeRetryWait):
						}
						continue
					}
					// On last retry, log as warning instead of error since service might register later
					logger.Warnf("subscribe to service %s still has no hosts after %d retries, will continue monitoring",
						key, MaxSubscribeRetry)
					subscribeErr = nil // Don't treat as fatal error
					break
				}

				// For other errors, log and retry
				logger.Warnf("subscribe attempt %d/%d for service %s failed: %s",
					retry+1, MaxSubscribeRetry, key, subscribeErr)
				if retry < MaxSubscribeRetry-1 {
					// Wait with exit signal check for graceful shutdown
					select {
					case <-n.exit:
						logger.Info("nacosIntfListener received exit signal during subscription retry, aborting")
						delete(n.serviceInfoMap, key)
						l.wg.Done()
						l.Close()
						return nil
					case <-time.After(SubscribeRetryWait):
					}
				}
			}

			if subscribeErr != nil {
				logger.Errorf("failed to subscribe to service %s after %d retries: %s",
					key, MaxSubscribeRetry, subscribeErr)
				// Clean up orphaned entry to prevent resource leak
				delete(n.serviceInfoMap, key)
				l.wg.Done()
				l.Close()
			}
		}
	}

	// handle deleted service
	for k, v := range n.serviceInfoMap {
		if _, ok := newServiceMap[k]; !ok {
			delete(n.serviceInfoMap, k)
			v.listener.Close()
		}
	}

	return nil
}

func getSubscribeName(url *dubboCommon.URL) string {
	var buffer bytes.Buffer
	buffer.Write([]byte(dubboCommon.DubboNodes[dubboCommon.PROVIDER]))
	appendParam(&buffer, url, constant.InterfaceKey)
	appendParam(&buffer, url, constant.VersionKey)
	appendParam(&buffer, url, constant.GroupKey)
	return buffer.String()
}

func appendParam(target *bytes.Buffer, url *dubboCommon.URL, key string) {
	value := url.GetParam(key, "")
	target.Write([]byte(constant.NacosServiceNameSeparator))
	if strings.TrimSpace(value) != "" {
		target.Write([]byte(value))
	}
}

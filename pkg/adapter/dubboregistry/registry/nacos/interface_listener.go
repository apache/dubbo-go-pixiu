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
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/remoting/zookeeper"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	MaxFailTimes = 2
	ConnDelay    = 3 * time.Second
)

var _ registry.Listener = new(nacosIntfListener)

type nacosIntfListener struct {
	exit            chan struct{}
	closeOnce       sync.Once
	watchOnce       sync.Once
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
	n.closeOnce.Do(func() {
		close(n.exit)
		n.wg.Wait()
		n.closeServiceListeners()
	})
}

func (n *nacosIntfListener) WatchAndHandle() {
	n.watchOnce.Do(func() {
		n.wg.Add(1)
		go n.watch()
	})
}

func (n *nacosIntfListener) watch() {
	defer n.wg.Done()
	var failTimes int64 = 0
	for {
		select {
		case <-n.exit:
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
			// Exit the watch if root node is in error
			if err == zookeeper.ErrNilNode {
				logger.Errorf("watching nacos services got errNilNode,so exit listen")
				return
			}
			if failTimes > MaxFailTimes {
				logger.Errorf("Error happens on nacos exceed max fail times: %s,so exit listen", MaxFailTimes)
				return
			}
			if waitForExit(n.exit, ConnDelay*time.Duration(failTimes)) {
				return
			}
			continue
		}
		failTimes = 0
		if err := n.updateServiceList(serviceList.Doms); err != nil {
			logger.Errorf("update service list failed %s", err)
		}
		if waitForExit(n.exit, time.Second*5) {
			return
		}
	}
}

func waitForExit(exit <-chan struct{}, delay time.Duration) bool {
	if delay <= 0 {
		select {
		case <-exit:
			return true
		default:
			return false
		}
	}
	timer := time.NewTimer(delay)
	defer timer.Stop()
	select {
	case <-exit:
		return true
	case <-timer.C:
		return false
	}
}

type serviceInfo struct {
	interfaceName string
	version       string
	group         string
	listener      *serviceListener
	subscription  *vo.SubscribeParam
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

			sub := &vo.SubscribeParam{
				ServiceName:       getSubscribeName(url),
				SubscribeCallback: l.Callback,
				GroupName:         n.regConf.Group,
			}

			if err := n.client.Subscribe(sub); err != nil {
				logger.Errorf("subscribe listener with interfaceKey = %s, error = %s", l, err)
				l.wg.Done()
				continue
			}
			l.wg.Done()
			svcInfo.listener = l
			svcInfo.subscription = sub
			n.serviceInfoMap[key] = svcInfo
		}
	}

	// handle deleted service
	for k, v := range n.serviceInfoMap {
		if _, ok := newServiceMap[k]; !ok {
			n.unsubscribeService(v)
			delete(n.serviceInfoMap, k)
		}
	}

	return nil
}

func (n *nacosIntfListener) closeServiceListeners() {
	for k, v := range n.serviceInfoMap {
		n.unsubscribeService(v)
		delete(n.serviceInfoMap, k)
	}
}

func (n *nacosIntfListener) unsubscribeService(svcInfo *serviceInfo) {
	if svcInfo == nil || svcInfo.listener == nil {
		return
	}
	if err := n.client.Unsubscribe(svcInfo.subscription); err != nil {
		logger.Errorf("unsubscribe listener with interfaceKey = %s, error = %s", svcInfo.listener, err)
	}
	svcInfo.listener.Close()
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

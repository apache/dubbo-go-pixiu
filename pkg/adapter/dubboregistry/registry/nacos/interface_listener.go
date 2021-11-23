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
	"fmt"
	"strings"
	"sync"
	"time"
)

import (
	dubboCommon "dubbo.apache.org/dubbo-go/v3/common"
	"dubbo.apache.org/dubbo-go/v3/common/constant"
	dubboConfig "dubbo.apache.org/dubbo-go/v3/config"
	dubboRegistry "dubbo.apache.org/dubbo-go/v3/registry"

	"github.com/apache/dubbo-go/common"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

import (
	common2 "github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/common"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/remoting/zookeeper"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

const (
	MaxFailTimes = 2
	ConnDelay    = 3 * time.Second
)

var _ registry.Listener = new(nacosIntfListener)

type nacosIntfListener struct {
	exit                 chan struct{}
	client               naming_client.INamingClient
	reg                  *NacosRegistry
	wg                   sync.WaitGroup
	dubbogoNacosRegistry dubboRegistry.Registry
	addr                 string
	adapterListener      common2.RegistryEventListener
}

// newNacosIntfListener returns a new nacosIntfListener with pre-defined path according to the registered type.
func newNacosIntfListener(client naming_client.INamingClient, addr string, reg *NacosRegistry, adapterListener common2.RegistryEventListener) registry.Listener {
	return &nacosIntfListener{
		exit:            make(chan struct{}),
		client:          client,
		reg:             reg,
		addr:            addr,
		adapterListener: adapterListener,
	}
}

func (z *nacosIntfListener) Close() {
	close(z.exit)
	z.wg.Wait()
}

func (z *nacosIntfListener) WatchAndHandle() {
	var err error
	z.dubbogoNacosRegistry, err = dubboConfig.NewRegistryConfigBuilder().
		SetProtocol("nacos").
		SetAddress(z.addr).
		Build().GetInstance(common.CONSUMER)
	if err != nil {
		logger.Errorf("create nacos registry with address = %s error = %s", z.addr, err)
		return
	}
	z.wg.Add(1)
	go z.watch()
}

func (z *nacosIntfListener) watch() {
	defer z.wg.Done()
	var (
		failTimes  int64 = 0
		delayTimer       = time.NewTimer(ConnDelay * time.Duration(failTimes))
	)
	defer delayTimer.Stop()
	for {
		serviceList, err := z.client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
			PageSize: 100,
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
			delayTimer.Reset(ConnDelay * time.Duration(failTimes))
			<-delayTimer.C
			continue
		}
		failTimes = 0
		if err := z.updateServiceList(serviceList.Doms); err != nil {
			logger.Errorf("update service list failed %s", err)
		}
		time.Sleep(time.Second * 3)
	}
}

type serviceInfo struct {
	interfaceName string
	version       string
	group         string
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

func (z *nacosIntfListener) updateServiceList(serviceList []string) error {
	// todo lock all svc listener

	allSvcListener := z.reg.GetAllSvcListener()
	subscribedServiceKeysMap := make(map[string]bool)
	for k := range allSvcListener {
		subscribedServiceKeysMap[k] = true
	}
	serviceNeedUpdate := make([]*serviceInfo, 0)
	for _, v := range serviceList {
		svcInfo := fromServiceFullKey(v)
		if svcInfo == nil {
			// invalid nacos dubbo service key
			continue
		}
		key := svcInfo.String()
		if _, ok := allSvcListener[key]; !ok {
			serviceNeedUpdate = append(serviceNeedUpdate, svcInfo)
		} else {
			delete(subscribedServiceKeysMap, key)
		}
	}
	if len(serviceNeedUpdate) == 0 && len(subscribedServiceKeysMap) == 0 {
		// there is no needs to update
		return nil
	}
	// subscribedServiceKeysMap is services needs to be removed
	// serviceNeedUpdate is services needs to be subscribed

	for _, v := range serviceNeedUpdate {
		url, _ := dubboCommon.NewURL("mock://localhost:8848")
		url.SetParam(constant.InterfaceKey, v.interfaceName)
		url.SetParam(constant.GroupKey, v.group)
		url.SetParam(constant.VersionKey, v.version)
		l := newNacosSrvListener(url, z.client, z.adapterListener)
		l.wg.Add(1)
		go func(v *serviceInfo) {
			defer l.wg.Done()
			z.reg.SetSvcListener(l.url.ServiceKey(), l)
			if err := z.dubbogoNacosRegistry.Subscribe(url, l); err != nil {
				logger.Errorf("subscribe listener with interfaceKey = %s, error = %s", l, err)
				z.reg.RemoveSvcListener(l.url.ServiceKey())
			}
		}(v)
	}
	// todo deal with subscribedServiceKeysMap services to be removed
	for k := range subscribedServiceKeysMap {
		z.reg.RemoveSvcListener(k)
	}
	return nil
}

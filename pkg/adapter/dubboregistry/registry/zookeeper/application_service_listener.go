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
	"encoding/json"
	"fmt"
	"strings"
	"sync"
	"time"
)

import (
	dubboCommon "dubbo.apache.org/dubbo-go/v3/common"
	dubboConst "dubbo.apache.org/dubbo-go/v3/common/constant"
	"dubbo.apache.org/dubbo-go/v3/metadata/definition"
	dr "dubbo.apache.org/dubbo-go/v3/registry"
	"dubbo.apache.org/dubbo-go/v3/registry/servicediscovery"
	"dubbo.apache.org/dubbo-go/v3/remoting/zookeeper/curator_discovery"
	"github.com/dubbo-go-pixiu/pixiu-api/pkg/api/config"
	"github.com/dubbogo/go-zookeeper/zk"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/common"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/registry"
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/dubboregistry/remoting/zookeeper"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
)

var _ registry.Listener = new(applicationServiceListener)

// applicationServiceListener normally monitors the /services/[:application]
type applicationServiceListener struct {
	urls        []*dubboCommon.URL
	servicePath string
	client      *zookeeper.ZooKeeperClient

	exit            chan struct{}
	wg              sync.WaitGroup
	adapterListener common.RegistryEventListener
}

// newApplicationServiceListener creates a new zk service listener
func newApplicationServiceListener(path string, client *zookeeper.ZooKeeperClient, adapterListener common.RegistryEventListener) *applicationServiceListener {
	return &applicationServiceListener{
		servicePath:     path,
		client:          client,
		exit:            make(chan struct{}),
		adapterListener: adapterListener,
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
		children, e, err := asl.client.GetChildrenW(asl.servicePath)
		if err != nil {
			failTimes++
			logger.Infof("watching (path{%s}) = error{%v}", asl.servicePath, err)
			// Exit the watch if root node is in error
			if err == zookeeper.ErrNilNode {
				logger.Errorf("watching (path{%s}) got errNilNode,so exit listen", asl.servicePath)
				return
			}
			if failTimes > MaxFailTimes {
				logger.Errorf(
					"Error happens on (path{%s}) exceed max fail times: %v,so exit listen",
					asl.servicePath, MaxFailTimes,
				)
				return
			}
			delayTimer.Reset(ConnDelay * time.Duration(failTimes))
			<-delayTimer.C
			continue
		}
		failTimes = 0
		if continueLoop := asl.waitEventAndHandlePeriod(children, e); !continueLoop {
			return
		}

	}
}

func (asl *applicationServiceListener) waitEventAndHandlePeriod(children []string, e <-chan zk.Event) bool {
	tickerTTL := defaultTTL
	ticker := time.NewTicker(tickerTTL)
	defer ticker.Stop()
	asl.handleEvent(children)
	for {
		select {
		case <-ticker.C:
			asl.handleEvent(children)
		case zkEvent := <-e:
			logger.Warnf(
				"get a zookeeper e{type:%s, server:%s, path:%s, state:%d-%s, err:%s}",
				zkEvent.Type.String(), zkEvent.Server, zkEvent.Path, zkEvent.State, zookeeper.StateToString(zkEvent.State), zkEvent.Err,
			)
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
	fetchChildren, err := asl.client.GetChildren(asl.servicePath)
	if err != nil {
		logger.Warnf("Error when retrieving newChildren in path: %s, Error:%s", asl.servicePath, err.Error())
		// disable the API
		for _, url := range asl.urls {
			bkConf, _, _, _ := registry.ParseDubboString(url.String())
			apiPattern := registry.GetAPIPattern(bkConf)

			if err := asl.adapterListener.OnDeleteRouter(config.Resource{Path: apiPattern}); err != nil {
				logger.Errorf("Error={%s} when try to remove API by path: %s", err.Error(), apiPattern)
			}
		}
		return
	}

	asl.urls = asl.getUrls(fetchChildren[0])
	for _, url := range asl.urls {
		bkConfig, _, location, err := registry.ParseDubboString(url.String())
		if err != nil {
			logger.Warnf("Parse dubbo interface provider %s failed; due to %s", url.String(), err.Error())
			continue
		}
		if len(bkConfig.ApplicationName) == 0 || len(bkConfig.Interface) == 0 {
			continue
		}
		methods, err := asl.getMethods(bkConfig.Interface)
		if err != nil {
			logger.Warnf("Get methods of interface %s failed; use prefix pattern to match url, due to %s", bkConfig.Interface, err.Error())
		}

		apiPattern := registry.GetAPIPattern(bkConfig)
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
		if len(methods) != 0 {
			for i := range methods {
				api := registry.CreateAPIConfig(apiPattern, location, bkConfig, methods[i], mappingParams)
				if err := asl.adapterListener.OnAddAPI(api); err != nil {
					logger.Errorf("Error={%s} happens when try to add api %s", err.Error(), api.Path)
				}
			}
		} else {
			// can't fetch methods, use http prefix pattern
			api := registry.CreateAPIConfig(apiPattern, location, bkConfig, constant.AnyValue, mappingParams)
			if err := asl.adapterListener.OnAddAPI(api); err != nil {
				logger.Errorf("Error={%s} happens when try to add api %s", err.Error(), api.Path)
			}
		}
	}
}

// getUrls return exported urls from instance
func (asl *applicationServiceListener) getUrls(path string) []*dubboCommon.URL {
	insPath := strings.Join([]string{asl.servicePath, path}, constant.PathSlash)
	data, err := asl.client.GetContent(insPath)
	if err != nil {
		logger.Errorf("Error when get content in path: %s, Error:%s", insPath, err.Error())
		return nil
	}

	// convert the data to service instance
	iss := &curator_discovery.ServiceInstance{}
	err = json.Unmarshal(data, iss)
	if err != nil {
		logger.Warnf("Parse service instance %s failed due to %s", insPath, err.Error())
		return nil
	}
	instance := toZookeeperInstance(iss)

	metaData := instance.GetMetadata()
	metadataInfo, err := servicediscovery.GetMetadataInfo(instance.GetServiceName(), instance, metaData[dubboConst.ExportedServicesRevisionPropertyName])
	if err != nil {
		logger.Errorf("get instance %s metadata info error %v", insPath, err.Error())
		return nil
	}
	instance.SetServiceMetadata(metadataInfo)
	urls := make([]*dubboCommon.URL, 0, len(metadataInfo.Services))
	for _, service := range metadataInfo.Services {
		urls = append(urls, instance.ToURLs(service)...)
	}
	return urls
}

// toZookeeperInstance convert to registry's service instance
func toZookeeperInstance(cris *curator_discovery.ServiceInstance) dr.ServiceInstance {
	pl, ok := cris.Payload.(map[string]any)
	if !ok {
		logger.Errorf("toZookeeperInstance{%s} payload is not map[string]interface{}", cris.ID)
		return nil
	}
	mdi, ok := pl["metadata"].(map[string]any)
	if !ok {
		logger.Errorf("toZookeeperInstance{%s} metadata is not map[string]interface{}", cris.ID)
		return nil
	}
	md := make(map[string]string, len(mdi))
	for k, v := range mdi {
		md[k] = fmt.Sprint(v)
	}
	return &dr.DefaultServiceInstance{
		ID:          cris.ID,
		ServiceName: cris.Name,
		Host:        cris.Address,
		Port:        cris.Port,
		Enable:      true,
		Healthy:     true,
		Metadata:    md,
	}
}

// getMethods return the methods of a service
func (asl *applicationServiceListener) getMethods(in string) ([]string, error) {
	path := strings.Join([]string{methodsRootPath, in}, constant.PathSlash)
	data, err := asl.client.GetContent(path)
	if err != nil {
		return nil, err
	}

	sd := &definition.ServiceDefinition{}
	err = json.Unmarshal(data, sd)
	if err != nil {
		return nil, err
	}

	methods := make([]string, len(sd.Methods))
	for i := range sd.Methods {
		methods[i] = sd.Methods[i].Name
	}
	return methods, nil
}

// Close closes this listener
func (asl *applicationServiceListener) Close() {
	close(asl.exit)
	asl.wg.Wait()
}

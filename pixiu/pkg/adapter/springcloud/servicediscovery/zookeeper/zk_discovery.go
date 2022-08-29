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
	"path"
	"strings"
	"sync"
	"time"
)

import (
	gxzookeeper "github.com/dubbogo/gost/database/kv/zk"
)

import (
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/springcloud/common"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/adapter/springcloud/servicediscovery"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/model"
	"github.com/apache/dubbo-go-pixiu/pixiu/pkg/remote/zookeeper"
)

const (
	ZKRootPath     = "/services"
	ZKName         = "SpringCloud-Zookeeper"
	StatUP         = "UP"
	MaxFailTimes   = 3
	DefaultTimeout = "3s"
	ConnDelay      = 3 * time.Second
	defaultTTL     = 3 * time.Second
)

type zookeeperDiscovery struct {
	basePath        string
	targetService   []string
	listener        servicediscovery.ServiceEventListener
	instanceMapLock sync.Mutex
	instanceMap     map[string]*servicediscovery.ServiceInstance
	zkListener      map[string]zookeeper.Listener
	clientFacade    *BaseZkClientFacade
}

func NewZKServiceDiscovery(targetService []string, config *model.RemoteConfig, listener servicediscovery.ServiceEventListener) (servicediscovery.ServiceDiscovery, error) {

	var err error

	if len(config.Timeout) == 0 {
		config.Timeout = DefaultTimeout
	}

	client, err := zookeeper.NewZkClient(ZKName, config)

	if err != nil {
		return nil, err
	}

	rootPath := ZKRootPath
	if len(strings.TrimSpace(config.Root)) > 0 {
		rootPath = strings.TrimSpace(config.Root)
	}

	z := &zookeeperDiscovery{
		basePath:      rootPath,
		listener:      listener,
		targetService: targetService,
		instanceMap:   make(map[string]*servicediscovery.ServiceInstance),
		zkListener:    map[string]zookeeper.Listener{},
		clientFacade: &BaseZkClientFacade{
			name:       ZKName,
			client:     client,
			conf:       config,
			clientLock: sync.Mutex{},
			wg:         sync.WaitGroup{},
			done:       make(chan struct{}),
		},
	}
	go zookeeper.HandleClientRestart(z.clientFacade)
	return z, err
}

func (sd *zookeeperDiscovery) QueryAllServices() ([]servicediscovery.ServiceInstance, error) {
	serviceNames, err := sd.queryForNames()
	logger.Debugf("%s get all services by root path %s, services %v", common.ZKLogDiscovery, sd.basePath, serviceNames)
	if err != nil {
		logger.Errorf("get all services error: %v", err.Error())
		return nil, err
	}
	return sd.QueryServicesByName(serviceNames)
}

func (sd *zookeeperDiscovery) QueryServicesByName(serviceNames []string) ([]servicediscovery.ServiceInstance, error) {

	var instancesAll []servicediscovery.ServiceInstance
	for _, s := range serviceNames {

		pathForName := sd.pathForName(s)
		ids, err := sd.getClient().GetChildren(pathForName)
		logger.Debugf("%s get services %s, services instanceIds %s", common.ZKLogDiscovery, s, ids)
		if err != nil {
			// todo refactor gost zk, make it return the definite err
			if strings.Contains(err.Error(), "none children") {
				logger.Debugf("%s get nodes from zookeeper fail: %s", common.ZKLogDiscovery, err.Error())
			} else {
				logger.Errorf("%s get services [%s] nodes from zookeeper fail: %s", common.ZKLogDiscovery, s, err.Error())
			}
			continue
		}

		for _, id := range ids {
			var instance *servicediscovery.ServiceInstance
			instance, err = sd.queryForInstance(s, id)
			if err != nil {
				return nil, err
			}
			instancesAll = append(instancesAll, *instance)
		}
	}

	for _, instance := range instancesAll {
		if sd.instanceMap[instance.ID] == nil {
			sd.instanceMap[instance.ID] = &instance
		}
	}

	return instancesAll, nil
}

func (sd *zookeeperDiscovery) Register() error {
	logger.Debugf("%s Register implement me!!", common.ZKLogDiscovery)
	return nil
}

func (sd *zookeeperDiscovery) UnRegister() error {
	logger.Debugf("%s UnRegister implement me!!", common.ZKLogDiscovery)
	return nil
}

func (sd *zookeeperDiscovery) getClient() *gxzookeeper.ZookeeperClient {
	if err := zookeeper.ValidateZookeeperClient(sd.clientFacade, "zka3"); err != nil {
		logger.Errorf("%s ValidateZookeeperClient error %s", common.ZKLogDiscovery, err)
	}
	return sd.clientFacade.ZkClient()
}

func (sd *zookeeperDiscovery) Subscribe() error {

	logger.Debugf("%s Subscribe ...", common.ZKLogDiscovery)

	sd.zkListener[sd.basePath] = newZkAppListener(sd)
	sd.zkListener[sd.basePath].WatchAndHandle()

	logger.Debugf("%s Subscribe Success!", common.ZKLogDiscovery)
	return nil
}

func (sd *zookeeperDiscovery) Unsubscribe() error {
	logger.Debugf("%s Unsubscribe ...", common.ZKLogDiscovery)

	for k, listener := range sd.zkListener {
		logger.Infof("Unsubscribe listener %s", k)
		listener.Close()
	}
	logger.Debugf("%s Unsubscribe Success!", common.ZKLogDiscovery)
	return nil
}

func (sd *zookeeperDiscovery) queryForInstance(name string, id string) (*servicediscovery.ServiceInstance, error) {
	path := sd.pathForInstance(name, id)
	data, _, err := sd.getClient().GetContent(path)
	if err != nil {
		return nil, err
	}
	sczk := &SpringCloudZKInstance{}
	instance := &servicediscovery.ServiceInstance{}
	err = json.Unmarshal(data, sczk)
	if err != nil {
		return nil, err
	}
	instance.Port = sczk.Port
	instance.ServiceName = sczk.Name
	instance.Host = sczk.Address
	instance.ID = sczk.ID
	instance.CLusterName = sczk.Name
	instance.Healthy = sczk.Payload.Metadata.InstanceStatus == StatUP
	return instance, nil
}

func (sd *zookeeperDiscovery) getServiceMap() map[string][]*servicediscovery.ServiceInstance {

	m := make(map[string][]*servicediscovery.ServiceInstance)

	for _, instance := range sd.instanceMap {

		if instances := m[instance.ServiceName]; instances == nil {
			m[instance.ServiceName] = []*servicediscovery.ServiceInstance{}
		}

		m[instance.ServiceName] = append(m[instance.ServiceName], instance)
	}

	return m
}

func (sd *zookeeperDiscovery) delServiceInstance(instance *servicediscovery.ServiceInstance) (bool, error) {

	if instance == nil {
		return true, nil
	}

	defer sd.instanceMapLock.Unlock()
	sd.instanceMapLock.Lock()
	if sd.instanceMap[instance.ID] != nil {
		sd.listener.OnDeleteServiceInstance(instance)
		delete(sd.instanceMap, instance.ID)
	}

	return true, nil
}

func (sd *zookeeperDiscovery) updateServiceInstance(instance *servicediscovery.ServiceInstance) (bool, error) {
	if instance == nil {
		return true, nil
	}

	defer sd.instanceMapLock.Unlock()
	sd.instanceMapLock.Lock()
	if sd.instanceMap[instance.ID] != nil {
		sd.listener.OnUpdateServiceInstance(instance)
		sd.instanceMap[instance.ID] = instance
	}

	return true, nil
}

func (sd *zookeeperDiscovery) addServiceInstance(instance *servicediscovery.ServiceInstance) (bool, error) {
	if instance == nil {
		return true, nil
	}

	defer sd.instanceMapLock.Unlock()
	sd.instanceMapLock.Lock()
	if sd.instanceMap[instance.ID] == nil {
		sd.listener.OnAddServiceInstance(instance)
		sd.instanceMap[instance.ID] = instance
	}

	return true, nil
}

func (sd *zookeeperDiscovery) queryForNames() ([]string, error) {
	children, err := sd.getClient().GetChildren(sd.basePath)
	// todo refactor gost zk, make it return the definite err
	if err != nil && strings.Contains(err.Error(), "none children") {
		logger.Debugf("%s get nodes from zookeeper fail: %s", common.ZKLogDiscovery, err.Error())
		return nil, nil
	}
	return children, err
}

func (sd *zookeeperDiscovery) pathForInstance(name, id string) string {
	return path.Join(sd.basePath, name, id)
}

func (sd *zookeeperDiscovery) pathForName(name string) string {
	return path.Join(sd.basePath, name)
}

type SpringCloudZKInstance struct {
	Name    string      `json:"name"`
	ID      string      `json:"id"`
	Address string      `json:"address"`
	Port    int         `json:"port"`
	SslPort interface{} `json:"sslPort"`
	Payload struct {
		Class    string `json:"@class"`
		ID       string `json:"id"`
		Name     string `json:"name"`
		Metadata struct {
			InstanceStatus string `json:"instance_status"`
		} `json:"metadata"`
	} `json:"payload"`
	RegistrationTimeUTC int64  `json:"registrationTimeUTC"`
	ServiceType         string `json:"serviceType"`
	URISpec             struct {
		Parts []struct {
			Value    string `json:"value"`
			Variable bool   `json:"variable"`
		} `json:"parts"`
	} `json:"uriSpec"`
}

type BaseZkClientFacade struct {
	name       string
	client     *gxzookeeper.ZookeeperClient
	clientLock sync.Mutex
	wg         sync.WaitGroup
	done       chan struct{}
	conf       *model.RemoteConfig
}

func (b *BaseZkClientFacade) Name() string {
	return b.name
}

func (b *BaseZkClientFacade) ZkClient() *gxzookeeper.ZookeeperClient {
	return b.client
}

func (b *BaseZkClientFacade) SetZkClient(client *gxzookeeper.ZookeeperClient) {
	b.client = client
}

func (b *BaseZkClientFacade) ZkClientLock() *sync.Mutex {
	return &b.clientLock
}

func (b *BaseZkClientFacade) WaitGroup() *sync.WaitGroup {
	return &b.wg
}

func (b *BaseZkClientFacade) Done() chan struct{} {
	return b.done
}

func (b *BaseZkClientFacade) RestartCallBack() bool {
	return true
}

func (b *BaseZkClientFacade) GetConfig() *model.RemoteConfig {
	return b.conf
}

func Keys(m map[string][]*servicediscovery.ServiceInstance) []string {
	j := 0
	keys := make([]string, len(m))
	for k := range m {
		keys[j] = k
		j++
	}
	return keys
}

func Diff(a, b []string) (diff []string) {
	m := make(map[string]bool)

	for _, item := range b {
		m[item] = true
	}

	for _, item := range a {
		if _, ok := m[item]; !ok {
			diff = append(diff, item)
		}
	}

	return diff
}

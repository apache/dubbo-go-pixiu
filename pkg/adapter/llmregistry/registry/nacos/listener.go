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
	"encoding/json"
	"reflect"
	"strconv"
	"strings"
	"sync"
	"time"
)

import (
	"dubbo.apache.org/dubbo-go/v3/remoting"

	"github.com/creasty/defaults"

	"github.com/hashicorp/go-uuid"

	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	nacosModel "github.com/nacos-group/nacos-sdk-go/model"
	"github.com/nacos-group/nacos-sdk-go/vo"
)

import (
	"github.com/apache/dubbo-go-pixiu/pkg/adapter/llmregistry/common"
	"github.com/apache/dubbo-go-pixiu/pkg/common/constant"
	"github.com/apache/dubbo-go-pixiu/pkg/logger"
	"github.com/apache/dubbo-go-pixiu/pkg/model"
)

const (
	// ServicePollingInterval How often to poll for the list of available services (to discover new ones).
	ServicePollingInterval = 30 * time.Second
)

// listener monitors Nacos for service changes.
type listener struct {
	client naming_client.INamingClient
	// Caches the last known instances for a given service. Key: ServiceName, Value: *sync.Map
	instanceCache sync.Map
	// Caches the set of services we are currently subscribed to. Key: ServiceName, Value: bool
	subscribedServices sync.Map
	regConf            *model.Registry
	adapterListener    common.RegistryEventListener // The callback to notify the gateway core.

	exit chan struct{}
	wg   sync.WaitGroup
}

// newNacosListener creates a new Nacos service listener.
func newNacosListener(client naming_client.INamingClient, regConf *model.Registry, adapterListener common.RegistryEventListener) *listener {
	return &listener{
		client:          client,
		exit:            make(chan struct{}),
		regConf:         regConf,
		adapterListener: adapterListener,
	}
}

// WatchAndHandle starts the background goroutine to watch for service changes.
func (l *listener) WatchAndHandle() {
	l.wg.Add(1)
	go l.watchForServices()
}

// watchForServices periodically polls Nacos to discover new services to subscribe to.
// The actual instance updates for subscribed services are push-based via the callback.
func (l *listener) watchForServices() {
	defer l.wg.Done()

	ticker := time.NewTicker(ServicePollingInterval)
	defer ticker.Stop()

	// Perform an initial check immediately.
	l.discoverAndSubscribe()

	for {
		select {
		case <-l.exit:
			logger.Info("Nacos listener is stopping...")
			l.unsubscribeAll()
			return
		case <-ticker.C:
			l.discoverAndSubscribe()
		}
	}
}

func (l *listener) discoverAndSubscribe() {
	serviceList, err := l.client.GetAllServicesInfo(vo.GetAllServiceInfoParam{
		GroupName: l.regConf.Group,
		NameSpace: l.regConf.Namespace,
	})
	if err != nil {
		logger.Warnf("Failed to get service list from Nacos: %v", err)
		return
	}

	currentServices := make(map[string]struct{})
	for _, serviceName := range serviceList.Doms {
		currentServices[serviceName] = struct{}{}
		// If we aren't already subscribed to this service, subscribe now.
		if _, loaded := l.subscribedServices.LoadOrStore(serviceName, true); !loaded {
			err := l.client.Subscribe(&vo.SubscribeParam{
				ServiceName:       serviceName,
				GroupName:         l.regConf.Group,
				SubscribeCallback: l.serviceCallback,
			})
			if err != nil {
				logger.Errorf("Failed to subscribe to Nacos service %s: %v", serviceName, err)
				l.subscribedServices.Delete(serviceName) // Remove from map to retry next time.
			} else {
				logger.Infof("Successfully subscribed to Nacos service: %s", serviceName)
			}
		}
	}

	// Unsubscribe from services that no longer exist.
	l.subscribedServices.Range(func(key, value any) bool {
		serviceName := key.(string)
		if _, exists := currentServices[serviceName]; !exists {
			err := l.client.Unsubscribe(&vo.SubscribeParam{
				ServiceName: serviceName,
				GroupName:   l.regConf.Group,
			})
			if err != nil {
				logger.Errorf("Failed to unsubscribe from Nacos service %s: %v", serviceName, err)
			} else {
				logger.Infof("Successfully unsubscribed from Nacos service: %s", serviceName)
				l.subscribedServices.Delete(serviceName)
				l.instanceCache.Delete(serviceName)
			}
		}
		return true
	})
}

// serviceCallback is the function that Nacos SDK invokes when there's a change
// in the instances of a subscribed service. This is the injected callback.
func (l *listener) serviceCallback(services []nacosModel.SubscribeService, err error) {
	if err != nil {
		logger.Errorf("Nacos subscribe callback received an error: %v", err)
		return
	}
	if len(services) == 0 {
		logger.Warn("Nacos callback received an empty list of services, which might indicate all instances are offline.")
		// The logic to handle removal of a service with zero instances
		// is handled by the polling `discoverAndSubscribe` loop.
		return
	}

	serviceName := services[0].ServiceName
	logger.Debugf("Received callback for service: %s with %d instances", serviceName, len(services))

	oldCache, _ := l.instanceCache.LoadOrStore(serviceName, &sync.Map{})
	oldInstanceMap := oldCache.(*sync.Map)

	newInstanceMap := &sync.Map{}
	newEndpoints := make(map[string]*model.Endpoint)

	// Process the new list from Nacos.
	for i := range services {
		// Also check for health
		if !services[i].Enable || !services[i].Healthy {
			continue
		}
		instance := generateInstance(services[i])
		endpoint := generateEndpoint(instance)
		if endpoint == nil {
			continue
		}
		key := serviceName + constant.At + endpoint.ID
		newInstanceMap.Store(key, instance)
		newEndpoints[key] = endpoint
	}

	// Check for added or updated instances.
	for key, endpoint := range newEndpoints {
		if oldRaw, ok := oldInstanceMap.Load(key); ok {
			oldInstance := oldRaw.(nacosModel.Instance)
			newInstance, _ := newInstanceMap.Load(key)
			if !reflect.DeepEqual(oldInstance, newInstance) {
				l.handle(endpoint, remoting.EventTypeUpdate)
			}
		} else {
			l.handle(endpoint, remoting.EventTypeAdd)
		}
	}

	// Check for removed instances.
	oldInstanceMap.Range(func(key, value any) bool {
		instanceKey := key.(string)
		if _, ok := newEndpoints[instanceKey]; !ok {
			instance := value.(nacosModel.Instance)
			endpoint := generateEndpoint(instance)
			l.handle(endpoint, remoting.EventTypeDel)
		}
		return true
	})

	// Update the cache with the new state.
	l.instanceCache.Store(serviceName, newInstanceMap)
}

func (l *listener) handle(endpoint *model.Endpoint, action remoting.EventType) {
	if endpoint == nil {
		return
	}
	logger.Infof("Handling endpoint event: %v for %s at %s", action, endpoint.Name, endpoint.Address.Address)
	switch action {
	case remoting.EventTypeAdd, remoting.EventTypeUpdate:
		if err := l.adapterListener.OnAddEndpoint(endpoint); err != nil {
			logger.Errorf("Failed to add/update endpoint %s: %s", endpoint.Name, err.Error())
		}
	case remoting.EventTypeDel:
		if err := l.adapterListener.OnRemoveEndpoint(endpoint); err != nil {
			logger.Errorf("Failed to remove endpoint %s: %s", endpoint.Name, err.Error())
		}
	}
}

// Close gracefully stops the listener.
func (l *listener) Close() {
	close(l.exit)
	l.wg.Wait()
}

func (l *listener) unsubscribeAll() {
	l.subscribedServices.Range(func(key, value any) bool {
		serviceName := key.(string)
		err := l.client.Unsubscribe(&vo.SubscribeParam{
			ServiceName: serviceName,
			GroupName:   l.regConf.Group,
		})
		if err != nil {
			logger.Errorf("Failed to unsubscribe from Nacos service %s on shutdown: %v", serviceName, err)
		} else {
			logger.Infof("Unsubscribed from Nacos service %s on shutdown.", serviceName)
		}
		return true
	})
}

func generateEndpoint(instance nacosModel.Instance) *model.Endpoint {
	if instance.Metadata == nil {
		logger.Warnf("Nacos instance metadata is empty, instance: %+v", instance)
		return nil
	}

	ret := &model.Endpoint{
		Address: model.SocketAddress{},
		LLMMeta: &model.LLMMeta{},
	}

	err := defaults.Set(ret)
	if err != nil {
		logger.Warnf("Failed to set default values for endpoint: %v", err)
		return nil
	}

	if ip, ok := instance.Metadata["ip"]; ok {
		ret.Address.Address = ip
	}

	if port, ok := instance.Metadata["port"]; ok {
		p, err := strconv.Atoi(port)
		if err != nil {
			logger.Warnf("Invalid port in metadata: %s, error: %v", port, err)
		}
		ret.Address.Port = p
	}

	if id, ok := instance.Metadata["id"]; ok {
		ret.ID = id
	} else {
		ret.ID, _ = uuid.GenerateUUID()
	}

	if name, ok := instance.Metadata["name"]; ok {
		ret.Name = name
	}
	if address, ok := instance.Metadata["address"]; ok {
		ret.Address.Domains = strings.Split(address, ",")
	}
	if apiKey, ok := instance.Metadata["llm-meta.api_key"]; ok {
		ret.LLMMeta.APIKey = apiKey
	}
	if retryPolicy, ok := instance.Metadata["llm-meta.retry_policy.name"]; ok {
		ret.LLMMeta.RetryPolicy.Name = model.RetryTypeValue[retryPolicy]
	}
	if retryConfig, ok := instance.Metadata["llm-meta.retry_policy.config"]; ok {
		err := json.Unmarshal([]byte(retryConfig), &ret.LLMMeta.RetryPolicy.Config)
		if err != nil {
			logger.Warnf("Failed to parse retry policy config JSON: %s, error: %v", retryConfig, err)
		}
	}
	if fallback, ok := instance.Metadata["llm-meta.fallback"]; ok {
		ret.LLMMeta.Fallback = strings.ToLower(strings.TrimSpace(fallback)) == "true"
	}

	ret.Metadata = instance.Metadata

	return ret
}

func generateInstance(ss nacosModel.SubscribeService) nacosModel.Instance {
	return nacosModel.Instance{
		InstanceId:  ss.InstanceId,
		Ip:          ss.Ip,
		Port:        ss.Port,
		ServiceName: ss.ServiceName,
		Valid:       ss.Valid,
		Enable:      ss.Enable,
		Weight:      ss.Weight,
		Metadata:    ss.Metadata,
		ClusterName: ss.ClusterName,
		Healthy:     ss.Healthy,
	}
}

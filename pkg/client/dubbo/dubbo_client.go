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

package dubbo

import (
	"context"
	"encoding/json"
	"io/ioutil"
	"strings"
	"sync"
	"time"
)

import (
	"github.com/apache/dubbo-go/common/constant"
	dg "github.com/apache/dubbo-go/config"
	"github.com/apache/dubbo-go/protocol/dubbo"
)

import (
	"github.com/dubbogo/dubbo-go-proxy/pkg/client"
	"github.com/dubbogo/dubbo-go-proxy/pkg/config"
	"github.com/dubbogo/dubbo-go-proxy/pkg/logger"
)

// TODO java class name elem
const (
	JavaStringClassName = "java.lang.String"
	JavaLangClassName   = "java.lang.Long"
)

var (
	_DubboClient *DubboClient
	onceClient   = sync.Once{}
	dgCfg        dg.ConsumerConfig
)

// DubboClient client to generic invoke dubbo
type DubboClient struct {
	mLock              sync.RWMutex
	GenericServicePool map[string]*dg.GenericService
}

// SingletonDubboClient singleton dubbo clent
func SingletonDubboClient() *DubboClient {
	if _DubboClient == nil {
		onceClient.Do(func() {
			_DubboClient = NewDubboClient()
		})
	}

	return _DubboClient
}

// NewDubboClient create dubbo client
func NewDubboClient() *DubboClient {
	return &DubboClient{
		mLock:              sync.RWMutex{},
		GenericServicePool: make(map[string]*dg.GenericService),
	}
}

// Init init dubbo, config mapping can do here
func (dc *DubboClient) Init() error {
	dgCfg = dg.GetConsumerConfig()
	//can change config here
	dg.SetConsumerConfig(dgCfg)
	dg.Load()
	dc.GenericServicePool = make(map[string]*dg.GenericService)
	return nil
}

// Close
func (dc *DubboClient) Close() error {
	return nil
}

// Call invoke service
func (dc *DubboClient) Call(r *client.Request) (resp client.Response, err error) {
	dm := r.API.Method.IntegrationRequest
	gs := dc.Get(dm.Interface, dm.DubboBackendConfig.Version, dm.Group, dm)

	var reqData []interface{}
	var indate []byte
	var e error

	indate, e = ioutil.ReadAll(r.IngressRequest.Body)

	if e != nil {
		return *client.EmptyResponse, err
	}

	l := len(dm.ParamTypes)
	switch {
	case l == 1:
		t := dm.ParamTypes[0]
		switch t {
		case JavaStringClassName:
			var s string
			if err := json.Unmarshal(indate, &s); err != nil {
				logger.Errorf("params parse error:%+v", err)
			} else {
				reqData = append(reqData, s)
			}
		case JavaLangClassName:
			var i int
			if err := json.Unmarshal(indate, &i); err != nil {
				logger.Errorf("params parse error:%+v", err)
			} else {
				reqData = append(reqData, i)
			}
		default:
			bodyMap := make(map[string]interface{})
			if err := json.Unmarshal(indate, &bodyMap); err != nil {
				return *client.EmptyResponse, err
			} else {
				reqData = append(reqData, bodyMap)
			}
		}
	case l > 1:
		if err = json.Unmarshal(indate, &reqData); err != nil {
			return *client.EmptyResponse, err
		}
	}

	logger.Debugf("[dubbogo proxy] invoke, method:%s, types:%s, reqData:%s", dm.Method, dm.ParamTypes, reqData)

	rst, err := gs.Invoke(context.Background(), []interface{}{dm.Method, dm.ParamTypes, reqData})
	if err != nil {
		return *client.EmptyResponse, err
	}
	if rst == nil {
		return client.Response{}, nil
	}
	resp = rst.(client.Response)
	logger.Debugf("[dubbogo proxy] dubbo client resp:%v", resp)
	return *NewDubboResponse(resp), nil
}

func (dc *DubboClient) get(key string) *dg.GenericService {
	dc.mLock.RLock()
	defer dc.mLock.RUnlock()
	return dc.GenericServicePool[key]
}

func (dc *DubboClient) check(key string) bool {
	dc.mLock.RLock()
	defer dc.mLock.RUnlock()
	if _, ok := dc.GenericServicePool[key]; ok {
		return true
	} else {
		return false
	}
}

func (dc *DubboClient) create(key string, irequest config.IntegrationRequest) *dg.GenericService {
	referenceConfig := dg.NewReferenceConfig(irequest.Interface, context.TODO())
	referenceConfig.InterfaceName = irequest.Interface
	referenceConfig.Cluster = constant.DEFAULT_CLUSTER
	var registers []string
	for k := range dgCfg.Registries {
		registers = append(registers, k)
	}
	referenceConfig.Registry = strings.Join(registers, ",")

	if len(irequest.DubboBackendConfig.Protocol) == 0 {
		referenceConfig.Protocol = dubbo.DUBBO
	} else {
		referenceConfig.Protocol = irequest.DubboBackendConfig.Protocol
	}

	referenceConfig.Version = irequest.DubboBackendConfig.Version
	referenceConfig.Group = irequest.Group
	referenceConfig.Generic = true
	if len(irequest.DubboBackendConfig.Retries) == 0 {
		referenceConfig.Retries = "3"
	} else {
		referenceConfig.Retries = irequest.DubboBackendConfig.Retries
	}
	dc.mLock.Lock()
	defer dc.mLock.Unlock()
	referenceConfig.GenericLoad(key)
	time.Sleep(200 * time.Millisecond) //sleep to wait invoker create
	clientService := referenceConfig.GetRPCService().(*dg.GenericService)

	dc.GenericServicePool[key] = clientService
	return clientService
}

// Get find a dubbo GenericService
func (dc *DubboClient) Get(interfaceName, version, group string, ir config.IntegrationRequest) *dg.GenericService {
	key := strings.Join([]string{ir.DubboBackendConfig.ApplicationName, interfaceName, version, group}, "_")
	if dc.check(key) {
		return dc.get(key)
	}

	return dc.create(key, ir)
}
